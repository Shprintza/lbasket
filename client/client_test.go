package client_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/orov-io/lbasket/client"
	. "github.com/smartystreets/goconvey/convey"
)

// convey phrases
const (
	givenAClient                 = "Given a lbasket client"
	callHandlerByService         = "When call is handler by the service"
	newBasketCall                = "When a new basket call is done."
	responseShouldBeOK           = "Then response should be OK"
	givenAEmptyBasket            = "Given a empty basket"
	tryToAddExistingProduct      = "When you try to add a product that exists to your basket"
	productIsAdded               = "Then product is added"
	tryToAddInexistendProduct    = "When you try to add a product that does not exits to your basket"
	responseShouldBeKO           = "Then response should be KO"
	productIsNotAdded            = "Then product is not added"
	tryToGetAListOfProducts      = "When you request to get a list of available products"
	availableProductsAreReturned = "Then a list of available products is returned"
)

func TestPing(t *testing.T) {
	Convey(givenAClient, t, func() {

		Convey(callHandlerByService, func() {
			pong, err := client.Ping()
			Convey(responseShouldBeOK, func() {

				So(err, ShouldBeNil)
				So(pong.Message, ShouldNotBeEmpty)
			})
		})
	})
}

func TestNewBasket(t *testing.T) {
	Convey(givenAClient, t, func() {

		Convey(newBasketCall, func() {
			newBasket, err := client.NewBasket()
			Convey(responseShouldBeOK, func() {

				So(err, ShouldBeNil)
				So(isUUID(newBasket.UUID), ShouldBeTrue)
			})
		})
	})
}

func isUUID(candidate string) bool {
	_, error := uuid.Parse(candidate)
	return error == nil
}

func TestGetProducts(t *testing.T) {
	Convey(givenAClient, t, func() {

		Convey(tryToGetAListOfProducts, func() {
			products, err := client.GetAvailableProducts()
			Convey(availableProductsAreReturned, func() {

				So(err, ShouldBeNil)
				So(len(products), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestAddProduct(t *testing.T) {
	Convey(givenAEmptyBasket, t, func() {
		newBasket, err := client.NewBasket()
		updatedBasket, err := client.AddProduct()

		Convey(tryToAddExistingProduct, func() {
			Convey(productIsAdded, func() {

				So(err, ShouldBeNil)
				So(isUUID(newBasket.UUID), ShouldBeTrue)
			})
		})
	})
}
