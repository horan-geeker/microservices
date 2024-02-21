package router

import (
	"microservices/controller"
	"microservices/pkg/app"
	"microservices/repository"
	"microservices/router/middleware"
)

func init() {
	authService := controller.NewAuthController(repository.NewFactory())
	router := app.GetApp()
	router.POST("/auth/login", authService.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authService.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authService.ChangePassword)
	router.POST("/auth/change-password-by-phone", authService.ChangePasswordByPhone)
}
