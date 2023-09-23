package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
)

var rdb *redis.Client

type redisStore struct {
	rdb *redis.Client
}

func (c *redisStore) Auth() store.AuthCache {
	return newAuth(c)
}

func (c *redisStore) Users() store.UserCache {
	return newUser(c)
}

func (c *redisStore) Del(ctx context.Context, key string) error {
	return c.rdb.WithContext(ctx).Del(key).Err()
}

// GetRedisInstance .
func GetRedisInstance(opts *options.RedisOptions) store.CacheFactory {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       opts.DB,
	})

	_, err := conn.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &redisStore{rdb: conn}
}
