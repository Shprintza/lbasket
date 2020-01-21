package checkout

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

// ProductDB is an interface that knows how to CRUD product catalog on actual db.
// This provide an abstraction level over the DB.
type ProductDB interface {
	SeedProducts(products []*Product) error
	GetProducts() ([]*Product, error)
	GetProduct(code string) (*Product, error)
	IsProductNotExistError(err error) bool
}

// ProductManager implements the ProductManager interface on top
// on the badger DB. This allow us to be thread-safe without an external DB.
type ProductManager struct {
	db ProductDB
}

// NewProductManager return a ProductManager  with the
// desired badger DB attached.
func NewProductManager(db ProductDB) *ProductManager {
	return &ProductManager{
		db: db,
	}
}

// SeedProducts fills the DB with available products
func (m *ProductManager) SeedProducts(products []*Product) error {
	return m.db.SeedProducts(products)
}

// GetAll fetches available products
func (m *ProductManager) GetAll() ([]*Product, error) {
	return m.db.GetProducts()
}

// IsProductAvailable checks if a product exist on DB by its code.
func (m *ProductManager) IsProductAvailable(code string) (bool, error) {
	products, err := m.GetAll()
	if err != nil {
		return false, err
	}

	return isProductIn(code, products), nil
}

// Get returns a product from database.
func (m *ProductManager) Get(code string) (*Product, error) {
	return m.db.GetProduct(code)
}

// IsProductNotExistError is raised when yo try to fetch a product that is not
// found in DB
func (m *ProductManager) IsProductNotExistError(err error) bool {
	return m.db.IsProductNotExistError(err)
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

func isProductIn(code string, products []*Product) bool {
	for _, product := range products {
		if product.Code == code {
			return true
		}
	}
	return false
}
