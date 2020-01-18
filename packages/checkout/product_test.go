package checkout_test

import (
	"testing"

	"github.com/orov-io/lbasket/packages/checkout"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBadgerProductManger(t *testing.T) {
	Convey("Given a badger based product manager", t, func() {
		productManager := checkout.NewBadgerProductManager(getBadgerDB())

		Convey("A valid manager is created", func() {
			So(productManager, ShouldNotBeNil)
		})
	})
}

func TestBadgerBasketManager_SeedProducts(t *testing.T) {
	Convey("Given a badger based product manager", t, func() {
		productManager := checkout.NewBadgerProductManager(db)

		Convey("When we seed the database", func() {
			err := productManager.SeedProducts(checkout.GetProductSeed())
			Convey("Operation is successfully", func() {
				So(err, ShouldBeNil)
				So(keyExists(checkout.ProductsKey), ShouldBeTrue)
			})

			Convey("DB is seeded", func() {
				products, err := productManager.GetProducts()
				So(err, ShouldBeNil)
				So(len(products), ShouldEqual, 3)
			})
		})
	})
}
