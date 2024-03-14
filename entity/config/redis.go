package config

import (
	"microservices/pkg/redis"
)

func NewRedisOptions() *redis.Options {
	env := GetEnvConfig()
	return &redis.Options{
		Host:     env.RedisHost,
		Password: env.RedisPassword,
		Port:     env.RedisPort,
		DB:       env.RedisDB,
	}
}
