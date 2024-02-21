package options

import (
	"microservices/entity/config"
	"microservices/pkg/redis"
)

func NewRedisOptions() *redis.Options {
	env := config.GetConfig()
	return &redis.Options{
		Host:     env.RedisHost,
		Password: env.RedisPassword,
		Port:     env.RedisPort,
		DB:       env.RedisDB,
	}
}
