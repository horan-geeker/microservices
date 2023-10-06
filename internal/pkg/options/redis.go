package options

import (
	"microservices/internal/config"
	"microservices/pkg/redis"
)

func NewRedisOptions() *redis.Options {
	env := config.NewConfig()
	return &redis.Options{
		Host:     env.RedisHost,
		Password: env.RedisPassword,
		Port:     env.RedisPort,
		DB:       env.RedisDB,
	}
}
