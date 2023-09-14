package auth

import (
	"microservices/internal/logic"
	"microservices/internal/store"
)

type AuthController struct {
	logic logic.LogicInterface
}

func NewAuthController(store store.DataFactory, cache store.CacheFactory) *AuthController {
	return &AuthController{
		logic: logic.NewLogic(store, cache),
	}
}
