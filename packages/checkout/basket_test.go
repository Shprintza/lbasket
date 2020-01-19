package checkout_test

import (
	"os"
	"testing"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
	"github.com/orov-io/lbasket/packages/checkout"
	. "github.com/smartystreets/goconvey/convey"
)

var db *badger.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	initDB()
}

func initDB() {
	opt := badger.DefaultOptions("").WithInMemory(true)
	var err error
	db, err = badger.Open(opt)
	if err != nil {
		panic(err)
	}
}

func shutdown() {
	db.Close()
}

func TestNewBadgerBasketManager(t *testing.T) {
	Convey("Given a badger based basket manager", t, func() {
		badgerManager := checkout.NewBadgerBasketManager(db)

		Convey("A valid manager is created", func() {
			So(badgerManager, ShouldNotBeNil)
		})
	})
}

func TestBadgerBasketManager_New(t *testing.T) {
	Convey("Given a new basket request", t, func() {
		basketManager := checkout.NewBadgerBasketManager(db)
		basket, err := basketManager.New()

		Convey("Operation is successfully", func() {
			So(err, ShouldBeNil)
			So(isUUID(basket.UUID), ShouldBeTrue)
			basket, err := basketManager.Get(basket.UUID)
			So(err, ShouldBeNil)
			So(basket.Items, ShouldBeEmpty)
			So(basket.Total, ShouldEqual, 0)
		})
	})
}

func TestBadgerBasketManager_Get(t *testing.T) {
	Convey("Given a basket", t, func() {
		basketManager := checkout.NewBadgerBasketManager(db)
		basket := &checkout.Basket{
			UUID: uuid.New().String(),
		}
		anotherBasket := &checkout.Basket{
			UUID: uuid.New().String(),
		}
		basketManager.Save(basket)

		Convey("When you try to fetch it", func() {
			fetchedBasket, err := basketManager.Get(basket.UUID)
			Convey("Operation is successfully", func() {
				So(err, ShouldBeNil)
				So(fetchedBasket.UUID, ShouldEqual, basket.UUID)
			})

			Convey("Another basket is not saved", func() {
				_, err := basketManager.Get(anotherBasket.UUID)
				So(checkout.IsBaskedNotExistError(err), ShouldBeTrue)

			})
		})
	})
}

func TestBadgerBasketManager__AddProductToBasket(t *testing.T) {
	Convey("Given a new basket", t, func() {
		basketManager := checkout.NewBadgerBasketManager(db)
		basket, _ := basketManager.New()

		Convey("When you add a new products to the basket", func() {
			product := checkout.GetProductSeed()[0]
			basket, err := basketManager.AddProductToBasket(product, basket)

			Convey("Operation is successfully", func() {
				So(err, ShouldBeNil)
				basket, err := basketManager.Get(basket.UUID)
				So(err, ShouldBeNil)
				So(len(basket.Items), ShouldEqual, 1)
				So(basket.Items[0].Product.Code, ShouldEqual, product.Code)
				So(basket.Total, ShouldEqual, product.Price)
			})
		})

		Convey("When you add a two identical new products to the basket", func() {
			product := checkout.GetProductSeed()[0]
			basket, err1 := basketManager.AddProductToBasket(product, basket)
			basket, err2 := basketManager.AddProductToBasket(product, basket)

			Convey("Operation is successfully", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				basket, err := basketManager.Get(basket.UUID)
				So(err, ShouldBeNil)
				So(len(basket.Items), ShouldEqual, 1)
				So(basket.Items[0].Product.Code, ShouldEqual, product.Code)
				So(basket.Items[0].Amount, ShouldEqual, 2)
				So(basket.Total, ShouldEqual, product.Price*2)
			})
		})

		Convey("When you add a two different new products to the basket", func() {
			product := checkout.GetProductSeed()[0]
			product2 := checkout.GetProductSeed()[1]
			basket, err1 := basketManager.AddProductToBasket(product, basket)
			basket, err2 := basketManager.AddProductToBasket(product2, basket)

			Convey("Operation is successfully", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				basket, err := basketManager.Get(basket.UUID)
				So(err, ShouldBeNil)
				So(len(basket.Items), ShouldEqual, 2)
				So(basket.Items[0].Amount, ShouldEqual, 1)
				So(basket.Items[1].Amount, ShouldEqual, 1)
				So(basket.Total, ShouldEqual, product.Price+product2.Price)
			})
		})
	})
}

func isUUID(candidate string) bool {
	_, error := uuid.Parse(candidate)
	return error == nil
}

func getBadgerDB() *badger.DB {
	return db
}
