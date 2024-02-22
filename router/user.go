package router

import (
	"microservices/controller"
	"microservices/repository"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	userController := controller.NewUserController(repository.NewFactory(), service.NewFactory())
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
