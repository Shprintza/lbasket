package models

import "github.com/orov-io/lbasket/packages/checkout"

// Pong models a ping response.
type Pong struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Basket models and exposes a
type Basket struct {
	*checkout.Basket
}
