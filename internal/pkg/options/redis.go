package options

import "microservices/internal/config"

type RedisOptions struct {
	Host     string
	Password string
	Port     int
	DB       int
}

func NewRedisOptions() *RedisOptions {
	env := config.NewEnvConfig()
	return &RedisOptions{
		Host:     env.RedisHost,
		Password: env.RedisPassword,
		Port:     env.RedisPort,
		DB:       env.RedisDB,
	}
}
