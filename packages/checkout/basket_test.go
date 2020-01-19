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
			exist := keyExists(basket.UUID)
			So(exist, ShouldBeTrue)
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

func keyExists(key string) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})

	return err == nil
}