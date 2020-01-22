package service

import (
	"fmt"
	"os"

	"github.com/orov-io/BlackBart/server"
)

// Service params
const (
	BasePathKey      = "SERVICE_BASE_PATH"
	VersionKey       = "SERVICE_VERSION"
	pingEndpoint     = "ping"
	basketEndpoint   = "baskets"
	productsEndpoint = "products"
)

var relativePath = os.Getenv(BasePathKey)
var version = os.Getenv(VersionKey)
var servicePath = fmt.Sprintf("/%v/%v", version, relativePath)

// AddRoutes add service handlers to the service
func AddRoutes(service *server.Service) {
	addPong(service)
	addBasket(service)
	addProducts(service)
}

func addPong(service *server.Service) {
	pingGroup := service.Group(getPathTo(pingEndpoint))
	{
		pingGroup.GET("", pong)
	}
}

func addBasket(service *server.Service) {
	basketGroup := service.Group(getPathTo(basketEndpoint))
	{
		basketGroup.GET("/:BasketUUID", getBasket)
		basketGroup.POST("", newBasket)
		basketGroup.POST("/:BasketUUID/products/:ProductCode", addProduct)
		basketGroup.DELETE("/:BasketUUID", deleteBasket)
	}
}

func addProducts(service *server.Service) {
	productGroup := service.Group(getPathTo(productsEndpoint))
	{
		productGroup.GET("", getProducts)
	}
}

func getPathTo(endpoint string) string {
	return fmt.Sprintf("%v/%v", servicePath, endpoint)
}
