package service

import (
	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/BlackBart/server"
	"github.com/orov-io/lbasket/models"
	"github.com/orov-io/lbasket/packages/checkout"
	"github.com/orov-io/lbasket/packages/lanabadger"
)

func pong(c *gin.Context) {
	sendPong(c, models.Pong{
		Status:  "OK",
		Message: "pong",
	})
}

func newBasket(c *gin.Context) {
	db, err := getLanaDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basketManager := checkout.NewBasketManager(db)
	basket, err := basketManager.New()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	sendBasked(c, basket)
}

func getProducts(c *gin.Context) {
	db, err := getLanaDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	productManager := checkout.NewProductManager(db)
	products, err := productManager.GetAll()
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

	db, err := getLanaDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basketManager := checkout.NewBasketManager(db)
	basket, err := basketManager.Get(request.BasketUUID)
	if basketManager.IsBaskedNotExistError(err) {
		response.SendBadRequest(c, err)
		return
	} else if err != nil {
		response.SendInternalError(c, err)
		return
	}

	productManager := checkout.NewProductManager(db)
	product, err := productManager.Get(request.ProductCode)
	if productManager.IsProductNotExistError(err) {
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

func getBasket(c *gin.Context) {
	request := new(models.GetBasketRequest)
	if err := c.ShouldBindUri(&request); err != nil {
		response.SendBadRequest(c, err)
		return
	}

	db, err := getLanaDB()
	if err != nil {
		response.SendInternalError(c, err)
		return
	}

	basketManager := checkout.NewBasketManager(db)
	basket, err := basketManager.Get(request.BasketUUID)
	if basketManager.IsBaskedNotExistError(err) {
		response.SendBadRequest(c, err)
		return
	} else if err != nil {
		response.SendInternalError(c, err)
		return
	}

	sendBasked(c, basket)
}

func getLanaDB() (checkout.DB, error) {
	db, err := server.GetInternalDB()

	return lanabadger.New(db), err
}
