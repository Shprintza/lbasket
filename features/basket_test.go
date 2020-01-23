package features

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/google/uuid"
	"github.com/orov-io/lbasket/client"
	"github.com/orov-io/lbasket/models"
)

const zeroValue = "0.00€"

var basket *models.Basket
var newBasketError error

func iHaveANewBasketRequest() error {
	basket, newBasketError = client.NewBasket()
	return nil
}

func iReceiveTheResponse() error {
	if newBasketError != nil || basket == nil {
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

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a new basket request$`, iHaveANewBasketRequest)
	s.Step(`^I receive the response$`, iReceiveTheResponse)
	s.Step(`^I should receive a new empty basket$`, iShouldReceiveANewEmptyBasket)

	s.BeforeFeature(func(*gherkin.Feature) {
		upServer()
	})

	s.AfterFeature(func(*gherkin.Feature) {
		downServer()
	})
}

func upServer() {
	fmt.Println("Please, be sure that you run 'make up [logs]' before start tests")
}

func downServer() {
	fmt.Println("Please,run 'make down' to shutdown the server")
}
