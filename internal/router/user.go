package router

import (
	"microservices/internal/controller/user"
	"microservices/internal/middleware"
)

func init() {
	userController := user.NewUserController(dataFactory, cacheFactory)
	router.GET("/users/:id", userController.Get)
	router.POST("/users/edit", middleware.Authenticate(), userController.Edit)
}
