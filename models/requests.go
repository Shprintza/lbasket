package models

import "github.com/gin-gonic/gin"

// PingRequestURI models the uri parameters for ping.
type PingRequestURI struct {
	PingUUID string `uri:"pingUUID" binding:"required"`
}

// PingRequest models the request data for ping.
type PingRequest struct {
	*PingRequestURI
	PingRequestBody
}

// PingRequestBody models the body parameters for pings.
type PingRequestBody struct {
	Ping string `json:"ping,omitempty" binding:"required"`
}

// Bind bind the context paramenters to a ping request.
func (request *PingRequest) Bind(c *gin.Context) error {
	if err := c.ShouldBindUri(request); err != nil {
		return err
	}

	if err := c.ShouldBindJSON(request); err != nil {
		return err
	}

	return nil
}

// AddProductRequest models a request to add a product to a basket
type AddProductRequest struct {
	BasketUUID  string `uri:"BasketUUID" binding:"required,uuid"`
	ProductCode string `uri:"ProductCode" binding:"required"`
}

// BasketRequests models a request to get or delete a single basket
type BasketRequests struct {
	BasketUUID string `uri:"BasketUUID" binding:"required,uuid"`
}
