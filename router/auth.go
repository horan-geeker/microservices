package router

import (
	"microservices/controller"
	"microservices/router/middleware"
)

func init() {
	authService := controller.NewAuthService(repositoryFactory)
	router.POST("/auth/login", authService.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authService.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authService.ChangePassword)
	router.POST("/auth/change-password-by-phone", authService.ChangePasswordByPhone)
}
