package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

func init() {
	l := logic.NewLogic(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	c := controller.NewCallbackController(l)
	router.GET("/callback/google-auth", c.GoogleAuthCallback)
}
