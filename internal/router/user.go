package router

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/controller/auth"
)

// Register .
func Register(r *gin.Engine) {
	route := r.Group("/user")
	{
		// todo 将 controller 返回固定的结构体和错误，需要结合 gin 框架来处理
		route.POST("/login", auth.Login)
		route.POST("/logout", auth.Logout)
	}
}
