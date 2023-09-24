package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"microservices/internal/pkg/consts"
	"microservices/internal/store"
)

type user struct {
	rdb *redis.Client
}

func (u *user) SetToken(ctx context.Context, id uint64, token string) error {
	return u.rdb.WithContext(ctx).Set(fmt.Sprintf(consts.RedisUserTokenKey, id), token, consts.UserTokenExpiredIn).Err()
}

func (u *user) GetToken(ctx context.Context, id uint64) (string, error) {
	return u.rdb.WithContext(ctx).Get(fmt.Sprintf(consts.RedisUserTokenKey, id)).Result()
}

func newUser(s *redisStore) store.UserCache {
	return &user{rdb: s.rdb}
}
