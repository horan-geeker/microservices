package router

import (
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	authService := service.NewAuthService(repositoryFactory)
	router.POST("/auth/login", authService.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authService.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authService.ChangePassword)
	router.POST("/auth/change-password-by-phone", authService.ChangePasswordByPhone)
}
