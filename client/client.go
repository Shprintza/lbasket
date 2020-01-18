package client

import (
	"os"
	"strconv"

	"github.com/orov-io/BlackBart/response"
	api "github.com/orov-io/BlackBeard"
	"github.com/orov-io/lbasket/models"
)

const (
	portKey          = "PORT"
	serviceKey       = "SERVICE_BASE_PATH"
	v1               = "v1"
	pingEndpoint     = "/ping"
	basketEndpoint   = "/basket"
	productsEndpoint = "/products"
)

var service = os.Getenv(serviceKey)
var port int

func init() {
	var err error
	port, err = strconv.Atoi(os.Getenv(portKey))
	if err != nil {
		port = 8080
	}
}

// Ping make a call to the is_alive endpoint.
func Ping() (*models.Pong, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithPort(port).
		WithVersion(v1).ToService(service)
	resp, err := client.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := models.Pong{}
	err = response.ParseTo(resp, &pong)

	return &pong, err
}

// NewBasket requests a new basket to the server.
func NewBasket() (*models.Basket, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithPort(port).
		WithVersion(v1).ToService(service)
	resp, err := client.POST(basketEndpoint, nil)
	if err != nil {
		return nil, err
	}
	basket := models.Basket{}
	err = response.ParseTo(resp, &basket)

	return &basket, err
}

// GetAvailableProducts fetches a list of available products
func GetAvailableProducts() ([]*models.Product, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithPort(port).
		WithVersion(v1).ToService(service)
	resp, err := client.GET(productsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	products := make([]*models.Product, 0)
	err = response.ParseTo(resp, &products)

	return products, err
}
