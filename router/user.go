package router

import (
	"microservices/controller"
	"microservices/pkg/app"
	"microservices/repository"
	"microservices/router/middleware"
)

func init() {
	userController := controller.NewUserController(repository.NewFactory())
	router := app.GetApp()
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
