package router

import (
	"microservices/internal/controller/auth"
	"microservices/internal/middleware"
)

func init() {
	authController := auth.NewAuthController(dataFactory, cacheFactory, serviceFactory)
	router.POST("/auth/login", authController.Login)
	router.POST("/auth/logout", middleware.Authenticate(), authController.Logout)
	router.POST("/auth/change-password", middleware.Authenticate(), authController.ChangePassword)
	router.POST("/auth/change-password-by-phone", authController.ChangePasswordByPhone)
}
