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
