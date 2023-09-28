package notify

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type notifyController struct {
	logic logic.LogicInterface
}

func NewNotifyController(store store.DataFactory, cache store.CacheFactory, service service.ServiceFactory) *notifyController {
	return &notifyController{
		logic: logic.NewLogic(store, cache, service),
	}
}
