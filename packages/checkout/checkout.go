package checkout

// DB is an interface to satisfy all needs of the checkout managers.
type DB interface {
	BasketDB
	ProductDB
}
