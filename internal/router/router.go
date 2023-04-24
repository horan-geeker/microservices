package router

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

// ToDo 在框架初始化的时候通过反射获取类型同时注册路由，这样就不需要在controller里每次获取参数映射，而变成了函数参数

// MicroserviceHandlerFunc .
type MicroserviceHandlerFunc func(*gin.Context) (*entity.Response, error)

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
