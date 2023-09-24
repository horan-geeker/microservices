package options

import (
	"microservices/internal/config"
	"microservices/internal/pkg/consts"
	"time"
)

// JwtOptions contains configuration items related to API server features.
type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"maxRefresh" mapstructure:"max-refresh"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewJwtOptions() *JwtOptions {
	env := config.NewEnvConfig()

	return &JwtOptions{
		Realm:      consts.AppName,
		Key:        env.JWTSecret,
		Timeout:    consts.UserTokenExpiredIn,
		MaxRefresh: consts.UserTokenExpiredIn,
	}
}
