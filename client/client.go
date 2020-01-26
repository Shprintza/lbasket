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
	portKey          = "CLIENT_LBASKET_SERVER_PORT"
	service          = "lana_basket"
	hostKey          = "CLIENT_LBASKET_SERVER_HOST"
	apiKeyKey        = "CLIENT_LBASKET_APIKEY"
	v1               = "v1"
	pingEndpoint     = "/ping"
	basketEndpoint   = "/baskets"
	productsEndpoint = "/products"
)

const defaultPort = 8080

// Client provides functionality to easily call server endpoints
type Client struct {
	c *api.Client
}

// New returns a fresh client pointing to desired host
func New(host string) *Client {
	return &Client{api.MakeNewClient().WithBasePath(host).
		WithVersion(v1).ToService(service)}
}

// NewWithDefaults returns a new client initialized with params from next env
// values:
// host: $LBASKET_SERVER_HOST
// port: $LBASKET_SERVER_PORT
func NewWithDefaults() *Client {

	host := os.Getenv(hostKey)
	key := os.Getenv(apiKeyKey)
	client := New(host)
	port, err := strconv.Atoi(os.Getenv(portKey))
	if err == nil && port != 0 {
		client = client.WithPort(port)
	}

	return client.WithAPIKey(key)
}

// WithPort attaches desired port to underlying API Client.
func (client *Client) WithPort(port int) *Client {
	client.c = client.c.WithPort(port)
	return client
}

// WithAPIKey forces the client to send provided api e¡key in each call to the
// server.
func (client *Client) WithAPIKey(key string) *Client {
	client.c = client.c.WithAPIKey(key)
	return client
}

// Ping make a call to the is_alive endpoint.
func (client *Client) Ping() (*models.Pong, error) {
	resp, err := client.c.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := models.Pong{}
	err = response.ParseTo(resp, &pong)

	return &pong, err
}

// NewBasket requests a new basket to the server.
func (client *Client) NewBasket() (*models.Basket, error) {
	resp, err := client.c.POST(basketEndpoint, nil)
	if err != nil {
		return nil, err
	}
	basket := models.Basket{}
	err = response.ParseTo(resp, &basket)

	return &basket, err
}

// GetAvailableProducts fetches a list of available products.
func (client *Client) GetAvailableProducts() ([]*models.Product, error) {
	resp, err := client.c.GET(productsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	products := make([]*models.Product, 0)
	err = response.ParseTo(resp, &products)

	return products, err
}

// AddProductToBasket requests a new basket to the server.
func (client *Client) AddProductToBasket(product, basket string) (*models.Basket, error) {
	uri := getAddProductURI(product, basket)
	resp, err := client.c.POST(uri, nil)
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
func (client *Client) GetBasket(basket string) (*models.Basket, error) {
	resp, err := client.c.GET(getGetBasketURI(basket), nil)
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
func (client *Client) DeleteBasket(basket string) error {
	resp, err := client.c.DELETE(getGetBasketURI(basket), nil)
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
