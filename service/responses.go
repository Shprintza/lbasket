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

func sendBasked(c *gin.Context, basket *checkout.Basket) {
	items2send := make([]*models.BasketItem, 0)
	basket2send := new(models.Basket)
	basket2send.Items = items2send
	basket2send.UUID = basket.UUID
	basket2send.Total = iToEuro(basket.Total)
	for _, item := range basket.Items {
		basket2send.Items = append(basket2send.Items, parseItem(item))
	}
	response.SendOK(c, basket2send)
}

func sendProducts(c *gin.Context, products []*checkout.Product) {
	products2send := parseProducts(products)
	response.SendOK(c, products2send)
}
func parseItem(item *checkout.BasketItem) *models.BasketItem {
	return &models.BasketItem{
		Product: parseProduct(item.Product),
		Amount:  item.Amount,
	}
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
		Price: iToEuro(product.Price),
	}
}

func iToEuro(value int) string {
	return fmt.Sprintf("%0.2f€", float64(value)/100.0)
}
