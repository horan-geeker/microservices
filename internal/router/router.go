package router

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

type MicroserviceHandlerFunc func(*gin.Context) (entity.Response, error)

func wrapperResponse(function MicroserviceHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, _ := c.Request.Context().Value("requestId").(string)
		response, err := function(c)
		if err != nil {
			c.JSON(500, gin.H{
				"request_id": requestId,
				"code":       500,
				"message":    err.Error(),
				"data":       nil,
			})
			return
		}
		response.RequestId = requestId
		c.JSON(200, response)
	}
}
