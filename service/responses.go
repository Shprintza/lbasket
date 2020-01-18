package service

import (
	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/lbasket/models"
	"github.com/orov-io/lbasket/packages/checkout"
)

func sendPong(c *gin.Context, pong models.Pong) {
	response.SendOK(c, pong)
}

func sendBaskedCreated(c *gin.Context, basket *checkout.Basket) {
	basketToSend := models.Basket{
		Basket: basket,
	}
	response.SendOK(c, basketToSend)
}
