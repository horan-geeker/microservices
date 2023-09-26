package system

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type SystemController struct {
	logic logic.LogicInterface
}

func NewSystemController(store store.DataFactory, cache store.CacheFactory, srv service.ServiceFactory) *SystemController {
	return &SystemController{
		logic: logic.NewLogic(store, cache, srv),
	}
}
