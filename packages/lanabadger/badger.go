package lanabadger

import "github.com/dgraph-io/badger/v2"

// DB implements both Basket and Product DB interfaces.
type DB struct {
	*badger.DB
}

// New returns a new BadgerDB with desired DB attached.
func New(db *badger.DB) *DB {
	return &DB{
		DB: db,
	}
}
