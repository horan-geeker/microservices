package notify

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type notifyController struct {
	logic logic.LogicInterface
}

func NewNotifyController(store store.Factory, cache store.CacheFactory, service service.Factory) *notifyController {
	return &notifyController{
		logic: logic.NewLogic(store, cache, service),
	}
}
