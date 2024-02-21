package router

import (
	"microservices/controller"
	"microservices/router/middleware"
)

func init() {
	userController := controller.NewUserController(repositoryFactory)
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
