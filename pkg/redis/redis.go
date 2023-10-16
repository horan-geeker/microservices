package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedis(opts *Options) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       opts.DB,
	})
	conn.AddHook(NewRedisLogHook())
	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
