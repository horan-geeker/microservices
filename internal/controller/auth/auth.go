package auth

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type AuthController struct {
	logic logic.LogicInterface
}

func NewAuthController(store store.DataFactory, cache store.CacheFactory, srv service.ServiceFactory) *AuthController {
	return &AuthController{
		logic: logic.NewLogic(store, cache, srv),
	}
}
