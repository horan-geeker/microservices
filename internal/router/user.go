package router

import (
	"microservices/internal/controller/user"
)

func init() {
	userController := user.NewUserController(dataFactory, cacheFactory)
	router.GET("/users/:id", userController.GetUserinfo)
}
