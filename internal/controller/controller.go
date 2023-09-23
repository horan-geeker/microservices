package controller

import (
	"microservices/internal/logic"
	"microservices/internal/store"
)

type Controller struct {
	logic logic.LogicInterface
}

// NewController creates a new controller
func NewController(store store.DataFactory, cache store.CacheFactory) *Controller {
	return &Controller{
		logic: logic.NewLogic(store, cache),
	}
}
