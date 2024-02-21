package router

import (
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	userController := service.NewUserController(repositoryFactory)
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
