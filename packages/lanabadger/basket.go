package lanabadger

import (
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/orov-io/lbasket/packages/checkout"
)

// SaveBasket stores a new basket in the badger DB.
func (db *DB) SaveBasket(basket *checkout.Basket) error {
	basket2store, err := json.Marshal(basket)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(
			[]byte(basket.UUID),
			[]byte(basket2store),
		)
		return err
	})

	return err
}

// GetBasket try to retrieve a Basket from the badger DB. If basket not exists,
// a not key found error is throwed.
func (db *DB) GetBasket(uuid string) (*checkout.Basket, error) {
	basket := new(checkout.Basket)
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uuid))
		if err == badger.ErrKeyNotFound {
			return db.newBaskedNotExistError(uuid)
		} else if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &basket)
			return err
		})
		return err
	})

	return basket, err
}

// IsBaskedNotExistError is raised when yo try to fetch a basket that is not
// found in DB
func (db *DB) IsBaskedNotExistError(err error) bool {
	_, ok := err.(*baskedNotExistError)
	return ok
}

// newBaskedNotExistError returns a new BaskedNotExistErrorError error.
func (db *DB) newBaskedNotExistError(uuid string) error {
	return &baskedNotExistError{
		BasketUUID: uuid,
	}
}

// BaskedNotExistError is used when we try to get a basket that does
// not exists in DB
type baskedNotExistError struct {
	BasketUUID string
}

func (e *baskedNotExistError) Error() string {
	return fmt.Sprintf("Basket %v does not exists", e.BasketUUID)
}
