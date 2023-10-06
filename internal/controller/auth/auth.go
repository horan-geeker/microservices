package auth

import (
	"microservices/internal/logic"
	"microservices/internal/service"
	"microservices/internal/store"
)

type AuthController struct {
	logic logic.LogicInterface
}

func NewAuthController(store store.Factory, cache store.CacheFactory, srv service.Factory) *AuthController {
	return &AuthController{
		logic: logic.NewLogic(store, cache, srv),
	}
}
