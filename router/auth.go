package router

import (
	"microservices/controller"
	"microservices/repository"
	"microservices/router/middleware"
)

func init() {
	authService := controller.NewAuthController(repository.NewFactory())
	router.POST("/auth/login", authService.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authService.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authService.ChangePassword)
	router.POST("/auth/change-password-by-phone", authService.ChangePasswordByPhone)
}
