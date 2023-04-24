package router

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/controller/auth"
)

// Register .
func Register(r *gin.Engine) {
	route := r.Group("/user")
	{
		// wrapperResponse 将 controller 返回固定的结构体和错误，结合 gin 框架来处理
		route.POST("/login", wrapperResponse(auth.Login))
		route.POST("/logout", wrapperResponse(auth.Logout))
	}
}
