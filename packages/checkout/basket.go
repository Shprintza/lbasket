package checkout

import (
	"math"

	"github.com/google/uuid"
)

// Basket models a lana checkout basket.
type Basket struct {
	UUID  string        `json:"uuid"`
	Items []*BasketItem `json:"items"`
	Total int
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

func (b *Basket) calcTotal() *Basket {
	var total int
	for _, item := range b.Items {
		total += item.calcTotal()
	}

	b.Total = total
	return b
}

// BasketItem models a chunk of same products
type BasketItem struct {
	Product *Product `json:"product"`
	Amount  int      `json:"amount"`
}

func (ib *BasketItem) isProduct(code string) bool {
	return ib.Product.Code == code
}

func (ib *BasketItem) calcTotal() int {
	switch ib.Product.Code {
	case PenCode:
		return ib.Product.Price * int(math.Ceil(float64(ib.Amount)*0.5))
	case TShirtCode:
		pricePerUnit := ib.Product.Price
		if ib.Amount >= 3 {
			pricePerUnit = (ib.Product.Price * 3) / 4
		}
		return pricePerUnit * ib.Amount

	default:
		return ib.Product.Price * ib.Amount
	}
}

// BasketDB is an interface that knows how to CRUD a basket on actual db
// This provide an abstraction level over the DB
type BasketDB interface {
	SaveBasket(basket *Basket) error
	GetBasket(uuid string) (*Basket, error)
	IsBaskedNotExistError(error) bool
}

// BasketManager implements a manager on top of the BasketDB interface
// This dependency injection of a interface allow us to change or DB with
// little effort
type BasketManager struct {
	db BasketDB
}

// NewBasketManager returns a BasketManager with the desired
// badger DB attached.
func NewBasketManager(db BasketDB) *BasketManager {
	return &BasketManager{
		db: db,
	}
}

// New returns a new basket. This function assert that new basket is saved.
func (m *BasketManager) New() (*Basket, error) {
	basket := newBasket()
	err := m.db.SaveBasket(basket)

	return basket, err
}

// Get fetches and returns desired basket if exists.
func (m *BasketManager) Get(uuid string) (*Basket, error) {
	basket, err := m.db.GetBasket(uuid)

	return basket.calcTotal(), err
}

// Save stores a Basket into database
func (m *BasketManager) Save(basket *Basket) error {
	return m.db.SaveBasket(basket)

}

// AddProductToBasket adds a product to basket saving changes in database.
func (m *BasketManager) AddProductToBasket(
	product *Product,
	basket *Basket,
) (*Basket, error) {
	basket.push(product)
	basket.calcTotal()
	return basket, m.db.SaveBasket(basket)
}

// IsBaskedNotExistError is raised when yo try to fetch a basket that is not
// found in DB
func (m *BasketManager) IsBaskedNotExistError(err error) bool {
	return m.db.IsBaskedNotExistError(err)
}
