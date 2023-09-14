package redis

import (
	"github.com/go-redis/redis"
)

type auth struct {
	rdb *redis.Client
}

func newAuth(s redisStore) *auth {
	return &auth{rdb: s.rdb}
}
