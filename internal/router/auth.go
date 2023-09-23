package router

import (
	"microservices/internal/controller/auth"
)

func init() {
	authController := auth.NewAuthController(dataFactory, cacheFactory)
	router.POST("/auth/login", authController.Login)
	router.POST("/auth/logout", authController.Logout)
	router.POST("/auth/change-password", authController.ChangePassword)
}
