package checkout

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

// Basket models a lana checkout basket.
type Basket struct {
	UUID string `json:"uuid"`
}

func newBasket() *Basket {
	basketUUID := uuid.New()
	return &Basket{basketUUID.String()}
}

// BasketManager is an interface that knows how to manages the
// basket entities livecycle.
// This provide an abstraction level over the DB
type BasketManager interface {
	NewBasket() (*Basket, error)
}

// BadgerBasketManager implements the BasketManager interface on top
// on the badger DB. This allow us to be thread-safe without an external DB.
type BadgerBasketManager struct {
	db *badger.DB
}

// NewBadgerBasketManager returns a BadgerBasketManager with the desired
// badger DB attached.
func NewBadgerBasketManager(db *badger.DB) *BadgerBasketManager {
	return &BadgerBasketManager{
		db: db,
	}
}

// New returns a new basket. This function assert that new basket is saved.
func (m *BadgerBasketManager) New() (*Basket, error) {
	basket := newBasket()
	basket2store, err := json.Marshal(basket)
	if err != nil {
		return nil, err
	}
	err = m.db.Update(func(txn *badger.Txn) error {
		txn.Set(
			[]byte(basket.UUID),
			[]byte(basket2store),
		)
		return nil
	})

	return basket, err
}
