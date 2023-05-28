package router

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

// ToDo 在框架初始化的时候通过反射获取类型同时注册路由，这样就不需要在controller里每次获取参数映射，而变成了函数参数

// MicroserviceHandlerFunc .
type MicroserviceHandlerFunc func(req any) (*entity.Response, error)

var routes []router

type router struct {
	Method        string
	Path          string
	Request       any
	BeforeHandler MicroserviceHandlerFunc
	Func          MicroserviceHandlerFunc
	AfterHandler  MicroserviceHandlerFunc
}

// Register router register
func Register(r *gin.Engine) {
	for _, route := range routes {
		r.Handle(route.Method, route.Path, wrapperGin(route))
	}
}

func wrapperGin(route router) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, _ := c.Request.Context().Value("requestId").(string)
		if err := c.ShouldBindJSON(&route.Request); err != nil {
			c.JSON(400, gin.H{
				"request_id": requestId,
				"code":       400,
				"message":    err.Error(),
				"data":       nil,
			})
			return
		}
		response, err := route.Func(route.Request)
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
