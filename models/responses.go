package models

import "github.com/orov-io/lbasket/packages/checkout"

// Pong models a ping response.
type Pong struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Basket models and exposes a basket to the client.
type Basket struct {
	*checkout.Basket
}

// Product models and exposes a product to the client.
type Product struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price string `json:"price"`
}
