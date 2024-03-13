package cache

import (
	"github.com/redis/go-redis/v9"
	"microservices/entity/config"
	redis2 "microservices/pkg/redis"
)

var (
	cache         Factory
	redisInstance = GetRedisInstance(config.NewRedisOptions())
)

type Factory interface {
	Auth() Auth
	User() User
}

type factory struct {
	rdb *redis.Client
}

func (f factory) Auth() Auth {
	return newAuth(f.rdb)
}

func (f factory) User() User {
	return newUser(f.rdb)
}

// NewFactory .
func NewFactory() Factory {
	if cache == nil {
		cache = &factory{rdb: redisInstance}
	}
	return cache
}

// GetRedisInstance .
func GetRedisInstance(opts *redis2.Options) *redis.Client {
	rdb, err := redis2.NewRedis(opts)
	if err != nil {
		panic(err)
	}
	return rdb
}
