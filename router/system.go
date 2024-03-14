package router

import (
	"microservices/cache"
	"microservices/controller/system"
	"microservices/model"
	"microservices/service"
)

func init() {
	systemController := system.NewController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/system/health", systemController.Health)
}
