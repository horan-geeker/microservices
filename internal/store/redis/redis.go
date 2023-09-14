package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"microservices/internal/config"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
	"time"
)

var rdb *redis.Client

type redisStore struct {
	rdb *redis.Client
}

func (c *redisStore) GetKey(ctx context.Context, key string) (string, error) {
	return c.rdb.WithContext(ctx).Get(key).Result()
}

func (c *redisStore) SetKey(ctx context.Context, key, value string, expire time.Duration) error {
	return c.rdb.WithContext(ctx).Set(key, value, expire).Err()
}

// GetRedisInstance .
func GetRedisInstance(opts *options.RedisOptions) store.CacheFactory {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: config.Env.RedisPassword,
		DB:       config.Env.RedisDB,
	})

	_, err := conn.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &redisStore{rdb: conn}
}
