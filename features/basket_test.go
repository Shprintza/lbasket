package features

import (
	"fmt"

	"github.com/DATA-DOG/godog"
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
	basket, stepError = client.NewBasket()
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
	_, stepError = client.GetBasket(invalidBasketUUID)
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
	basket, err = client.NewBasket()
	return err
}

func iCallToGetTheValidBasket() error {
	gettedBasket, stepError = client.GetBasket(basket.UUID)
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

	s.BeforeSuite(func() {
		upServer()
	})

	s.AfterSuite(func() {
		downServer()
	})
}

func upServer() {
	fmt.Println("Please, be sure that you run 'make up [logs]' before start tests")
}

func downServer() {
	fmt.Println("Please,run 'make down' to shutdown the server")
}
