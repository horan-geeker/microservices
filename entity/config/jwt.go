package config

import (
	consts2 "microservices/entity/consts"
	"microservices/entity/jwt"
)

// NewJwtOptions creates a Options object with default parameters.
func NewJwtOptions() *jwt.Options {
	env := GetConfig()
	return &jwt.Options{
		Realm:      consts2.AppName,
		Key:        env.JWTSecret,
		Timeout:    consts2.UserTokenExpiredIn,
		MaxRefresh: consts2.UserTokenExpiredIn,
	}
}
