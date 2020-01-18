package main

import (
	"os"

	"github.com/orov-io/BlackBart/server"
	"github.com/orov-io/lbasket/packages/checkout"
	"github.com/orov-io/lbasket/service"
)

const (
	envKey  = "ENV"
	portKey = "PORT"
	local   = "local"
)

var log = server.GetLogger()

func main() {

	app, err := server.StartDefaultService()
	if err != nil {
		log.WithError(err).Panic("Can't initialize the service ...")
	}

	service.AddRoutes(app)
	err = seedProducts()
	if err != nil {
		log.WithError(err).Panic("Can't seed database")
	}

	environment := os.Getenv(envKey)

	if environment == local {
		err = app.Run(":" + server.GetEnvPort(portKey))
	} else {
		err = nil
		app.SetMode(server.ReleaseMode)
		app.RunAppEngine()
	}

	if err != nil {
		log.WithError(err).Panic("Can't start the server")
	}

}

// seedProducts seeds the database with available products.
func seedProducts() error {
	db, err := server.GetInternalDB()
	if err != nil {
		return err
	}
	productManager := checkout.NewBadgerProductManager(db)
	return productManager.SeedProducts(checkout.GetProductSeed())
}
