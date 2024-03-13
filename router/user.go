package router

import (
	"microservices/cache"
	"microservices/controller/user"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	userController := user.NewController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
	router.POST("/users/register", userController.Register)
}
