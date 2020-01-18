package service

import (
	"fmt"
	"os"

	"github.com/orov-io/BlackBart/server"
)

// Service params
const (
	BasePathKey    = "SERVICE_BASE_PATH"
	VersionKey     = "SERVICE_VERSION"
	pingEndpoint   = "ping"
	basketEndpoint = "basket"
)

var relativePath = os.Getenv(BasePathKey)
var version = os.Getenv(VersionKey)
var servicePath = fmt.Sprintf("/%v/%v", version, relativePath)

// AddRoutes add service handlers to the service
func AddRoutes(service *server.Service) {
	addPong(service)
	addBasket(service)
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
		basketGroup.POST("", newBasket)
	}
}

func getPathTo(endpoint string) string {
	return fmt.Sprintf("%v/%v", servicePath, endpoint)
}
