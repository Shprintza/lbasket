package checkout_test

import (
	"testing"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/orov-io/lbasket/packages/checkout"
	"github.com/orov-io/lbasket/packages/lanabadger"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProductManager(t *testing.T) {
	Convey("Given a badger based product manager", t, func() {
		productManager := checkout.NewProductManager(getDB())

		Convey("A valid manager is created", func() {
			So(productManager, ShouldNotBeNil)
		})
	})
}

func TestBadgerProductManager_SeedProducts(t *testing.T) {
	Convey("Given a badger based product manager", t, func() {
		productManager := checkout.NewProductManager(db)

		Convey("When we seed the database", func() {
			err := productManager.SeedProducts(checkout.GetProductSeed())
			Convey("Operation is successfully", func() {
				So(err, ShouldBeNil)
				So(keyExists(lanabadger.ProductsKey), ShouldBeTrue)
			})

			Convey("DB is seeded", func() {
				products, err := productManager.GetAll()
				So(err, ShouldBeNil)
				So(len(products), ShouldEqual, 3)
				exits, err := productManager.IsProductAvailable(checkout.PenCode)
				So(err, ShouldBeNil)
				So(exits, ShouldBeTrue)
				exits, err = productManager.IsProductAvailable("ASDF")
				So(err, ShouldBeNil)
				So(exits, ShouldBeFalse)
			})
		})
	})
}

func TestBadgerProductManager_Get(t *testing.T) {
	Convey("Given a populate product db", t, func() {
		productManager := checkout.NewProductManager(db)
		productManager.SeedProducts(checkout.GetProductSeed())

		Convey("When we try to fetch an available product", func() {
			availableProductCode := checkout.GetProductSeed()[0].Code
			product, err := productManager.Get(availableProductCode)
			Convey("Operation is successfully", func() {
				So(err, ShouldBeNil)
				So(product.Code, ShouldEqual, availableProductCode)
			})
		})

		Convey("When we try to fetch an unavailable product", func() {
			unavailableProductCode := "BAD"
			product, err := productManager.Get(unavailableProductCode)
			Convey("Operation fails", func() {
				So(err, ShouldBeError)
				So(product, ShouldBeNil)
				So(err, ShouldBeError)
			})
		})
	})
}

func keyExists(key string) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})

	return err == nil
}
