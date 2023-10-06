package system

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type SystemController struct {
	logic logic.LogicInterface
}

func NewSystemController(store store.Factory, cache store.CacheFactory, srv service.Factory) *SystemController {
	return &SystemController{
		logic: logic.NewLogic(store, cache, srv),
	}
}
