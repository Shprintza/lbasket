package lanabadger

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/orov-io/lbasket/packages/checkout"
)

// ProductsKey is the kew to store all available products
const ProductsKey = "lanaProducts"

// SeedProducts fills the badger DB with available products
func (db *DB) SeedProducts(products []*checkout.Product) error {
	products2store, err := json.Marshal(products)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		txn.Set(
			[]byte(ProductsKey),
			[]byte(products2store),
		)
		return nil
	})

	return err
}

// GetProducts fetches available products
func (db *DB) GetProducts() ([]*checkout.Product, error) {
	products := make([]*checkout.Product, 0)
	err := db.View(func(txn *badger.Txn) error {
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

// GetProduct fetches a product by this code returning their struct.
func (db *DB) GetProduct(code string) (*checkout.Product, error) {
	products, err := db.GetProducts()
	if err != nil {
		return nil, err
	}

	product := getProductFromHaystack(code, products)
	if product == nil {
		return nil, newProductNotExistError(code)
	}

	return product, nil
}

func getProductFromHaystack(code string, products []*checkout.Product) *checkout.Product {
	for _, product := range products {
		if product.Code == code {
			return product
		}
	}
	return nil
}

// newProductNotExistError returns a new ProductNotExistErrorError error.
func newProductNotExistError(code string) error {
	return &productNotExistError{
		ProductCode: code,
	}
}

// IsProductNotExistError is raised when yo try to fetch a product that is not
// found in DB
func (db *DB) IsProductNotExistError(err error) bool {
	_, ok := err.(*productNotExistError)
	return ok
}

// ProductNotExistError is used when we try to get a product that does
// not exists in DB
type productNotExistError struct {
	ProductCode string
}

func (e *productNotExistError) Error() string {
	return fmt.Sprintf("Product %v does not exists", e.ProductCode)
}
