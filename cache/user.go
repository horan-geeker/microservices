package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	consts2 "microservices/entity/consts"
)

// User defines the user storage interface.
type User interface {
	SetToken(ctx context.Context, uid int, token string) error
	GetToken(ctx context.Context, uid int) (string, error)
	DeleteToken(ctx context.Context, uid int) error
}

type user struct {
	db  *gorm.DB
	rdb *redis.Client
}

func newUser(rdb *redis.Client) User {
	return &user{
		rdb: rdb,
	}
}

func (u *user) SetToken(ctx context.Context, id int, token string) error {
	return u.rdb.Set(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id), token, consts2.UserTokenExpiredIn).Err()
}

func (u *user) GetToken(ctx context.Context, id int) (string, error) {
	return u.rdb.Get(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id)).Result()
}

func (u *user) DeleteToken(ctx context.Context, id int) error {
	return u.rdb.Del(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id)).Err()
}
