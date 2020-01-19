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

	sendBasked(c, basket)
}

func getProducts(c *gin.Context) {
	db, err := server.GetInternalDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	productManager := checkout.NewBadgerProductManager(db)
	products, err := productManager.GetProducts()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	sendProducts(c, products)
}

func addProduct(c *gin.Context) {
	request := new(models.AddProductRequest)
	if err := c.ShouldBindUri(&request); err != nil {
		response.SendBadRequest(c, err)
		return
	}

	db, err := server.GetInternalDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basketManager := checkout.NewBadgerBasketManager(db)
	basket, err := basketManager.Get(request.BasketUUID)
	if checkout.IsBaskedNotExistError(err) {
		response.SendBadRequest(c, err)
		return
	} else if err != nil {
		response.SendInternalError(c, err)
		return
	}

	productManager := checkout.NewBadgerProductManager(db)
	product, err := productManager.Get(request.ProductCode)
	if checkout.IsProductNotExistError(err) {
		response.SendBadRequest(c, err)
		return
	} else if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basket, err = basketManager.AddProductToBasket(product, basket)
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	sendBasked(c, basket)
}
