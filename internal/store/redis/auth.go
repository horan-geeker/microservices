package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"microservices/internal/pkg/consts"
	"time"
)

type auth struct {
	rdb *redis.Client
}

func newAuth(s *redisInstance) *auth {
	return &auth{rdb: s.rdb}
}

func (a *auth) SetSmsCode(ctx context.Context, uid uint64, smsCode string) error {
	return a.rdb.Set(ctx, fmt.Sprintf(consts.RedisUserSmsKey, uid), smsCode, time.Minute).Err()
}

func (a *auth) GetSmsCode(ctx context.Context, uid uint64) (string, error) {
	return a.rdb.Get(ctx, fmt.Sprintf(consts.RedisUserSmsKey, uid)).Result()
}

func (a *auth) DeleteSmsCode(ctx context.Context, uid uint64) error {
	return a.rdb.Del(ctx, fmt.Sprintf(consts.RedisUserSmsKey, uid)).Err()
}

func (a *auth) GetEmailCode(ctx context.Context, uid uint64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (a *auth) SetEmailCode(ctx context.Context, uid uint64, emailCode string) error {
	//TODO implement me
	panic("implement me")
}

func (a *auth) DeleteEmailCode(ctx context.Context, uid uint64) error {
	//TODO implement me
	panic("implement me")
}
