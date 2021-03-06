package client_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orov-io/lbasket/client"
	"github.com/orov-io/lbasket/packages/checkout"
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

const (
	PenCode    = checkout.PenCode
	MugCode    = checkout.MugCode
	TShirtCode = checkout.TShirtCode
)

func TestPing(t *testing.T) {
	Convey(givenAClient, t, func() {
		lanaClient := client.NewWithDefaults()

		Convey(callHandlerByService, func() {
			pong, err := lanaClient.Ping()
			Convey(responseShouldBeOK, func() {

				So(err, ShouldBeNil)
				So(pong.Message, ShouldNotBeEmpty)
			})
		})
	})
}

func TestNewBasket(t *testing.T) {
	Convey(givenAClient, t, func() {
		lanaClient := client.NewWithDefaults()

		Convey(newBasketCall, func() {
			newBasket, err := lanaClient.NewBasket()
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
		lanaClient := client.NewWithDefaults()

		Convey(tryToGetAListOfProducts, func() {
			products, err := lanaClient.GetAvailableProducts()
			Convey(availableProductsAreReturned, func() {

				So(err, ShouldBeNil)
				So(len(products), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestAddProduct(t *testing.T) {
	Convey(givenAEmptyBasket, t, func() {
		lanaClient := client.NewWithDefaults()
		newBasket, _ := lanaClient.NewBasket()
		products, _ := lanaClient.GetAvailableProducts()
		product := products[0].Code
		updatedBasket, err := lanaClient.AddProductToBasket(
			product,
			newBasket.UUID,
		)

		Convey(tryToAddExistingProduct, func() {
			Convey(productIsAdded, func() {
				So(err, ShouldBeNil)
				So(updatedBasket.UUID, ShouldEqual, newBasket.UUID)
				So(len(updatedBasket.Items), ShouldEqual, 1)
				So(updatedBasket.Items[0].Product.Code, ShouldEqual, product)
			})
		})
	})
}

func TestGetBasket(t *testing.T) {
	testData := getGetBasketTestData()
	Convey(givenAEmptyBasket, t, func() {
		lanaClient := client.NewWithDefaults()
		basket, _ := lanaClient.NewBasket()

		for value, products := range testData {
			Convey(fmt.Sprintf("When we fill it with %s", products), func() {
				fillBasketWithProducts(basket.UUID, products)

				Convey("Total amount is correct", func() {
					updatedBasket, err := lanaClient.GetBasket(basket.UUID)
					So(err, ShouldBeNil)
					So(updatedBasket.UUID, ShouldEqual, basket.UUID)
					So(updatedBasket.Total, ShouldEqual, value)
				})
			})
		}
	})
}

func TestDeleteBasket(t *testing.T) {
	Convey("Given a basket that already exists", t, func() {
		lanaClient := client.NewWithDefaults()
		basket, _ := lanaClient.NewBasket()

		Convey("When we delete it", func() {
			err := lanaClient.DeleteBasket(basket.UUID)

			Convey("Operation is done", func() {
				So(err, ShouldBeNil)
			})

			Convey("We can't retrieve it anymore", func() {
				_, err := lanaClient.GetBasket(basket.UUID)
				So(err, ShouldBeError)
			})
		})
	})
}

func fillBasketWithProducts(basket string, products []string) {
	lanaClient := client.NewWithDefaults()
	for _, product := range products {
		lanaClient.AddProductToBasket(product, basket)
	}

}

func getGetBasketTestData() map[string][]string {
	data := make(map[string][]string, 0)
	data["32.50€"] = []string{PenCode, TShirtCode, MugCode}
	data["25.00€"] = []string{PenCode, TShirtCode, PenCode}
	data["65.00€"] = []string{TShirtCode, TShirtCode, TShirtCode, PenCode, TShirtCode}
	data["62.50€"] = []string{PenCode, TShirtCode, PenCode, PenCode, MugCode, TShirtCode, TShirtCode}
	return data
}
