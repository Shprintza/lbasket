package checkout

import (
	"encoding/json"
	badger "github.com/dgraph-io/badger/v2"
)

const ProductsKey = "lanaProducts"

// Available products
const (
	PenCode  = "PEN"
	PenName  = "Lana Pen"
	PenPrice = 500

	TShirtCode  = "TSHIRT"
	TShirtName  = "Lana T-Shirt"
	TShirtPrice = 2000

	MugCode  = "MUG"
	MugName  = "Lana Coffee Mug"
	MugPrice = 750
)

// Product models a single lana product
type Product struct {
	Code  string
	Name  string
	Price int
}

// ProductManager is an interface that knows how to manages the
// product entities livecycle.
// This provide an abstraction level over the DB.
type ProductManager interface {
	SeedProducts(product *Product) error
}

// BadgerProductManager implements the ProductManager interface on top
// on the badger DB. This allow us to be thread-safe without an external DB.
type BadgerProductManager struct {
	db *badger.DB
}

// NewBadgerProductManager return a BadgerProductManager  with the
// desired badger DB attached.
func NewBadgerProductManager(db *badger.DB) *BadgerProductManager {
	return &BadgerProductManager{
		db: db,
	}
}

// SeedProducts fills the DB with available products
func (m *BadgerProductManager) SeedProducts(products []*Product) error {
	products2store, err := json.Marshal(products)
	if err != nil {
		return err
	}
	err = m.db.Update(func(txn *badger.Txn) error {
		txn.Set(
			[]byte(ProductsKey),
			[]byte(products2store),
		)
		return nil
	})

	return err
}

// GetProducts fetches available products
func (m *BadgerProductManager) GetProducts() ([]*Product, error) {
	products := make([]*Product, 0)
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ProductsKey))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &products)
			return err
		})
		return err
	})

	return products, err
}

// GetProductSeed returns the seed for the database.
//
// Disclaimer: I know that this is not the best site to store this.
// I do it in order to no develop a migration client for badger and store seed
// in the migrations dir.
func GetProductSeed() []*Product {
	return []*Product{
		&Product{
			Code:  PenCode,
			Name:  PenName,
			Price: PenPrice,
		},
		&Product{
			Code:  TShirtCode,
			Name:  TShirtName,
			Price: TShirtPrice,
		},
		&Product{
			Code:  MugCode,
			Name:  MugName,
			Price: MugPrice,
		},
	}
}
