package client

import (
	"os"
	"strconv"

	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/BlackBart/server"
	api "github.com/orov-io/BlackBeard"
	"github.com/orov-io/lbasket/models"
)

const (
	portKey      = "PORT"
	serviceKey   = "SERVICE_BASE_PATH"
	v1           = "v1"
	pingEndpoint = "/ping"
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
	server.GetLogger().Info("PONG: ", pong)
	err = response.ParseTo(resp, &pong)

	return &pong, err
}
