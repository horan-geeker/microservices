package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"microservices/internal/pkg/consts"
	"microservices/internal/store"
)

type user struct {
	rdb *redis.Client
}

func (u *user) SetToken(ctx context.Context, id uint64, token string) error {
	return u.rdb.Set(ctx, fmt.Sprintf(consts.RedisUserTokenKey, id), token, consts.UserTokenExpiredIn).Err()
}

func (u *user) GetToken(ctx context.Context, id uint64) (string, error) {
	return u.rdb.Get(ctx, fmt.Sprintf(consts.RedisUserTokenKey, id)).Result()
}

func (u *user) DeleteToken(ctx context.Context, id uint64) error {
	return u.rdb.Del(ctx, fmt.Sprintf(consts.RedisUserTokenKey, id)).Err()
}

func newUser(s *redisInstance) store.UserCache {
	return &user{rdb: s.rdb}
}
