package features

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"github.com/orov-io/lbasket/client"
	"github.com/orov-io/lbasket/models"
)

const zeroValue = "0.00€"

var basket *models.Basket
var gettedBasket *models.Basket
var stepError error
var invalidBasketUUID string

func iHaveANewBasketRequest() error {
	lanaBasket := client.NewWithDefaults()
	basket, stepError = lanaBasket.NewBasket()
	return nil
}

func iReceiveTheResponse() error {
	if stepError != nil || basket == nil {
		return fmt.Errorf("Unable to retrieve a new basket")
	}

	return nil
}

func iShouldReceiveANewEmptyBasket() error {
	if _, err := uuid.Parse(basket.UUID); err != nil {
		return fmt.Errorf("Bad identifier for new Basket")
	}

	if len(basket.Items) != 0 {
		return fmt.Errorf("Basket is not empty")
	}

	if basket.Total != zeroValue {

		return fmt.Errorf("Basket value is not valid: %v", basket.Total)
	}

	return nil
}

func iHaveAnInvalidBasket() error {
	invalidBasketUUID = uuid.New().String()
	return nil
}

func iCallToGetInvalidBasket() error {
	_, stepError = client.NewWithDefaults().GetBasket(invalidBasketUUID)
	return nil
}

func iShouldReceiveAnError() error {
	if stepError == nil {
		return fmt.Errorf("Basket exists")
	}
	return nil
}

func iHaveAValidBasket() error {
	var err error
	basket, err = client.NewWithDefaults().NewBasket()
	return err
}

func iCallToGetTheValidBasket() error {
	gettedBasket, stepError = client.NewWithDefaults().GetBasket(basket.UUID)
	return nil
}

func iShouldReceiveDesiredBasket() error {
	if stepError != nil {
		return fmt.Errorf("Bad response with error: %v", stepError)
	}

	if !areSameBaskets(basket, gettedBasket) {
		return fmt.Errorf("Basket fetched unsuccessfully")
	}

	return nil
}

func areSameBaskets(original, copied *models.Basket) bool {
	return original.UUID == copied.UUID &&
		len(original.Items) == len(copied.Items) &&
		original.Total == copied.Total
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a new basket request$`, iHaveANewBasketRequest)
	s.Step(`^I receive the response$`, iReceiveTheResponse)
	s.Step(`^I should receive a new empty basket$`, iShouldReceiveANewEmptyBasket)

	s.Step(`^I have an invalid basket$`, iHaveAnInvalidBasket)
	s.Step(`^I try to retrive the invalid basket$`, iCallToGetInvalidBasket)
	s.Step(`^I should receive a error message$`, iShouldReceiveAnError)

	s.Step(`^I have a valid basket$`, iHaveAValidBasket)
	s.Step(`^I try to retrive the basket$`, iCallToGetTheValidBasket)
	s.Step(`^I should receive desired basket$`, iShouldReceiveDesiredBasket)

	wasRunning := false
	s.BeforeSuite(func() {
		if !serviceAlreadyRunning() {
			upServer()
			return
		}
		fmt.Print("Server was running before run tests!")
		wasRunning = true
	})

	s.AfterSuite(func() {
		if !wasRunning {
			downServer()
		}
		fmt.Println("Leaving server living")
	})
}

const dockerComposeTimeOut = 120
const dockerCompose = "docker-compose"
const useCustomFile = "-f"
const customFile = "../docker-compose.yml"
const up = "up"
const detachedMode = "-d"
const logs = "logs"
const streamMode = "-f"
const top = "top"

func serviceAlreadyRunning() bool {
	top := exec.Command(dockerCompose, useCustomFile, customFile, top)
	output, err := top.Output()
	if err != nil {
		log.Fatal("Can't run docker-compose ps: ", err)
	}

	return len(output) > 0
}

func upServer() {

	fmt.Println("Starting service...")
	startSpinner := showSpinner()
	upService()
	serviceOutput := getServiceStreamOutput()

	waitToServiceAlive(serviceOutput)

	startSpinner.Stop()
	fmt.Println("Service is running :)")
}

func showSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Start()
	return s
}

func getServiceStreamOutput() io.ReadCloser {
	service := exec.Command(dockerCompose, useCustomFile, customFile, logs, streamMode)
	serviceOutput, err := service.StdoutPipe()
	if err != nil {
		log.Fatalf("Can't open stream with service logs: %v\n", err)
	}
	err = service.Start()
	if err != nil {
		log.Fatalf("Can't do 'make up logs': %v\n", err)
	}
	return serviceOutput
}

func waitToServiceAlive(stream io.ReadCloser) {
	serviceBuild := make(chan bool)
	// This also can be achieved without go routines
	go waitForDockerCompose(stream, serviceBuild)
	isRunning := <-serviceBuild
	if !isRunning {
		downServer()
		log.Fatal("Unable to start the service")
	}
}

func upService() {
	serviceUp := exec.Command(dockerCompose, useCustomFile, customFile, up, detachedMode)
	err := serviceUp.Run()
	if err != nil {
		log.Fatalf("Can't up docker images: %v\n", err)
	}
}

func waitForDockerCompose(serviceOutput io.ReadCloser, serviceBuild chan<- bool) {

	reader := bufio.NewReader(serviceOutput)
	go shutdownIfTimeout(serviceBuild)

	for {
		checkServiceAliveness(reader, serviceBuild)
		return
	}
}

func checkServiceAliveness(reader *bufio.Reader, serviceBuild chan<- bool) {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			serviceBuild <- false
			close(serviceBuild)
			log.Println("EOF reached")
			return
		}

		if err != nil {
			log.Fatal("Read Error:", err)
			return
		}
		if strings.Contains(line, "Build Failed") {
			serviceBuild <- false
			close(serviceBuild)
			return
		}
		if strings.Contains(line, "Running...") {
			time.Sleep(2 * time.Second)
			serviceBuild <- true
			close(serviceBuild)
			return
		}
	}
}
func shutdownIfTimeout(serviceBuild chan<- bool) {
	time.Sleep(dockerComposeTimeOut * time.Second)
	fmt.Printf("Service don't initialized in first %v seconds\n", dockerComposeTimeOut)
	downServer()
	log.Fatal("Unable to start service")
}

func downServer() {
	fmt.Println("Shutting down service...")
	serviceDown := exec.Command("docker-compose", "-f", "../docker-compose.yml", "down")
	err := serviceDown.Run()
	if err != nil {
		fmt.Println("Can't shutdown service: ", err)
	}
	fmt.Println("Service is down")
}
