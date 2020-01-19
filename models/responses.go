package models

// Pong models a ping response.
type Pong struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Basket models and exposes a basket to the client.
type Basket struct {
	UUID  string        `json:"uuid"`
	Items []*BasketItem `json:"items"`
	Total string        `json:"total"`
}

// BasketItem models a chunk of same products
type BasketItem struct {
	Product *Product `json:"product"`
	Amount  int      `json:"amount"`
}

// Product models and exposes a product to the client.
type Product struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price string `json:"price"`
}
