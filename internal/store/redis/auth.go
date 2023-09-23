package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"microservices/internal/pkg/consts"
	"time"
)

type auth struct {
	rdb *redis.Client
}

func newAuth(s *redisStore) *auth {
	return &auth{rdb: s.rdb}
}

func (a *auth) SetSmsCode(ctx context.Context, uid uint64, smsCode string) error {
	return a.rdb.WithContext(ctx).Set(fmt.Sprintf(consts.RedisUserSmsKey, uid), smsCode, time.Minute).Err()
}

func (a *auth) GetSmsCode(ctx context.Context, uid uint64) (string, error) {
	return a.rdb.WithContext(ctx).Get(fmt.Sprintf(consts.RedisUserSmsKey, uid)).Result()
}

func (a *auth) DeleteSmsCode(ctx context.Context, uid uint64) error {
	return a.rdb.WithContext(ctx).Del(fmt.Sprintf(consts.RedisUserSmsKey, uid)).Err()
}
