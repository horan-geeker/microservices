package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	userController := controller.NewUserController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
