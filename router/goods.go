package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	goodsController := controller.NewGoodsController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/goods", middleware.Authenticate(), goodsController.GetList)
}
