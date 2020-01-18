package service

import (
	"fmt"

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

func sendProducts(c *gin.Context, products []*checkout.Product) {
	products2send := parseProducts(products)
	response.SendOK(c, products2send)
}

func parseProducts(products []*checkout.Product) []*models.Product {
	parsedProducts := make([]*models.Product, 0)
	for _, product := range products {
		parsedProducts = append(parsedProducts, parseProduct(product))
	}

	return parsedProducts
}

func parseProduct(product *checkout.Product) *models.Product {
	return &models.Product{
		Code:  product.Code,
		Name:  product.Name,
		Price: fmt.Sprintf("%0.2f", float64(product.Price)/10.0),
	}
}
