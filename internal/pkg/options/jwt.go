package options

import (
	"microservices/internal/config"
	"microservices/internal/pkg/consts"
	"microservices/pkg/jwt"
)

// NewJwtOptions creates a Options object with default parameters.
func NewJwtOptions() *jwt.Options {
	env := config.NewConfig()
	return &jwt.Options{
		Realm:      consts.AppName,
		Key:        env.JWTSecret,
		Timeout:    consts.UserTokenExpiredIn,
		MaxRefresh: consts.UserTokenExpiredIn,
	}
}
