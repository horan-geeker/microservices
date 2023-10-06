package redis

import (
	"github.com/redis/go-redis/v9"
	"microservices/internal/store"
	redis2 "microservices/pkg/redis"
	"sync"
)

type redisInstance struct {
	rdb *redis.Client
}

func (c *redisInstance) Auth() store.AuthCache {
	return newAuth(c)
}

func (c *redisInstance) Users() store.UserCache {
	return newUser(c)
}

var (
	redisFactory store.CacheFactory
	once         sync.Once
)

// GetRedisInstance .
func GetRedisInstance(opts *redis2.Options) store.CacheFactory {
	once.Do(func() {
		rdb, err := redis2.NewRedis(opts)
		if err != nil {
			panic(err)
		}
		redisFactory = &redisInstance{rdb: rdb}
	})
	return redisFactory
}
