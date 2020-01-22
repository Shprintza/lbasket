package client

import (
	"fmt"
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
	basketEndpoint   = "/baskets"
	productsEndpoint = "/products"
)

const defaultPort = 8080

// Ping make a call to the is_alive endpoint.
func Ping() (*models.Pong, error) {
	client := getClient()
	resp, err := client.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := models.Pong{}
	err = response.ParseTo(resp, &pong)

	return &pong, err
}

func getClient() *api.Client {
	service := os.Getenv(serviceKey)
	port, err := strconv.Atoi(os.Getenv(portKey))
	if err != nil {
		port = defaultPort
	}
	return api.MakeNewClient().WithDefaultBasePath().WithPort(port).
		WithVersion(v1).ToService(service)
}

// NewBasket requests a new basket to the server.
func NewBasket() (*models.Basket, error) {
	client := getClient()
	resp, err := client.POST(basketEndpoint, nil)
	if err != nil {
		return nil, err
	}
	basket := models.Basket{}
	err = response.ParseTo(resp, &basket)

	return &basket, err
}

// GetAvailableProducts fetches a list of available products.
func GetAvailableProducts() ([]*models.Product, error) {
	client := getClient()
	resp, err := client.GET(productsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	products := make([]*models.Product, 0)
	err = response.ParseTo(resp, &products)

	return products, err
}

// AddProductToBasket requests a new basket to the server.
func AddProductToBasket(product, basket string) (*models.Basket, error) {
	client := getClient()
	uri := getAddProductURI(product, basket)
	resp, err := client.POST(uri, nil)
	if err != nil {
		return nil, err
	}
	NewBasket := models.Basket{}
	err = response.ParseTo(resp, &NewBasket)

	return &NewBasket, err
}

func getAddProductURI(product, basket string) string {
	return fmt.Sprintf(
		"%s/%s%s/%s",
		basketEndpoint,
		basket,
		productsEndpoint,
		product,
	)
}

// GetBasket fetches a list of available products.
func GetBasket(basket string) (*models.Basket, error) {
	client := getClient()
	resp, err := client.GET(getGetBasketURI(basket), nil)
	if err != nil {
		return nil, err
	}
	fetchedBasket := new(models.Basket)
	err = response.ParseTo(resp, &fetchedBasket)

	return fetchedBasket, err
}

func getGetBasketURI(basket string) string {
	return fmt.Sprintf(
		"%s/%s",
		basketEndpoint,
		basket,
	)
}

// DeleteBasket deletes a basket in the server.
func DeleteBasket(basket string) error {
	client := getClient()
	resp, err := client.DELETE(getGetBasketURI(basket), nil)
	if err != nil {
		return err
	}
	if !isOK(resp.StatusCode) {
		return fmt.Errorf("Can't delete basket: %v", basket)
	}

	return nil
}

func isOK(code int) bool {
	return code > 200 && code < 400
}
