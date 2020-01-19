package checkout

import (
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

// Basket models a lana checkout basket.
type Basket struct {
	UUID  string        `json:"uuid"`
	Items []*BasketItem `json:"items"`
}

func newBasket() *Basket {
	basketUUID := uuid.New()
	return &Basket{
		UUID: basketUUID.String(),
	}
}

func (b *Basket) push(product *Product) {
	for key, item := range b.Items {
		if item.isProduct(product.Code) {
			b.Items[key].Amount++
			return
		}
	}

	newItem := &BasketItem{
		Product: product,
		Amount:  1,
	}

	b.Items = append(b.Items, newItem)
	return
}

// BasketItem models a chunk of same products
type BasketItem struct {
	Product *Product `json:"product"`
	Amount  int      `json:"amount"`
}

func (ib *BasketItem) isProduct(code string) bool {
	return ib.Product.Code == code
}

// BasketManager is an interface that knows how to manages the
// basket entities livecycle.
// This provide an abstraction level over the DB
type BasketManager interface {
	New() (*Basket, error)
	Exists() (bool, error)
	AddProductToBasket(string, string) (*Basket, error)
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

// Get fetches and returns desired basket if exists.
func (m *BadgerBasketManager) Get(uuid string) (*Basket, error) {
	basket := new(Basket)
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uuid))
		if isBadgerKeyNotFoundError(err) {
			return NewBaskedNotExistError(uuid)
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

// AddProductToBasket adds a product to basket saving changes in database.
func (m *BadgerBasketManager) AddProductToBasket(
	product *Product,
	basket *Basket,
) (*Basket, error) {
	basket.push(product)
	return basket, m.Save(basket)
}

// Save updates or create a basket in database.
func (m *BadgerBasketManager) Save(basket *Basket) error {
	basket2store, err := json.Marshal(basket)
	if err != nil {
		return err
	}
	err = m.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(
			[]byte(basket.UUID),
			[]byte(basket2store),
		)
		return err
	})

	return err
}

// BaskedNotExistError is used when we try to get a basket that does
// not exists in DB
type BaskedNotExistError struct {
	BasketUUID string
}

func (e *BaskedNotExistError) Error() string {
	return fmt.Sprintf("Basket %v does not exists", e.BasketUUID)
}

// NewBaskedNotExistError returns a new BaskedNotExistErrorError error.
func NewBaskedNotExistError(uuid string) error {
	return &BaskedNotExistError{
		BasketUUID: uuid,
	}
}

// IsBaskedNotExistError checks if the error is a BaskedNotExistError error.
func IsBaskedNotExistError(err error) bool {
	_, ok := err.(*BaskedNotExistError)
	return ok
}

func isBadgerKeyNotFoundError(err error) bool {
	return err == badger.ErrKeyNotFound
}
