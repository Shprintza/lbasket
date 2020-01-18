package service

import (
	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/BlackBart/server"
	"github.com/orov-io/lbasket/models"
	"github.com/orov-io/lbasket/packages/checkout"
)

func pong(c *gin.Context) {
	sendPong(c, models.Pong{
		Status:  "OK",
		Message: "pong",
	})
}

func newBasket(c *gin.Context) {
	db, err := server.GetInternalDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basketManager := checkout.NewBadgerBasketManager(db)
	basket, err := basketManager.New()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	sendBaskedCreated(c, basket)
}
