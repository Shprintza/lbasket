package models

// Pong models a ping response.
type Pong struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
