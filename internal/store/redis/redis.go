package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
	"microservices/pkg/log"
)

var redisFactory store.CacheFactory

type redisInstance struct {
	rdb *redis.Client
}

func (c *redisInstance) Auth() store.AuthCache {
	return newAuth(c)
}

func (c *redisInstance) Users() store.UserCache {
	return newUser(c)
}

// GetRedisInstance .
func GetRedisInstance(opts *options.RedisOptions) store.CacheFactory {
	if redisFactory != nil {
		return redisFactory
	}
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       opts.DB,
	})
	conn.AddHook(log.NewRedisLogHook())
	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	redisFactory = &redisInstance{rdb: conn}
	return redisFactory
}
