package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	authService := controller.NewAuthController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/auth/login", authService.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authService.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authService.ChangePassword)
	router.POST("/auth/change-password-by-phone", authService.ChangePasswordByPhone)
}
