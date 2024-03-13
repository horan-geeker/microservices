package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"microservices/entity/consts"
	"time"
)

type Auth interface {
	SetSmsCode(ctx context.Context, phone string, smsCode string) error
	GetSmsCode(ctx context.Context, phone string) (string, error)
	DeleteSmsCode(ctx context.Context, phone string) error
	GetEmailCode(ctx context.Context, email string) (string, error)
	SetEmailCode(ctx context.Context, email string, emailCode string) error
	DeleteEmailCode(ctx context.Context, email string) error
}

type authImpl struct {
	rdb *redis.Client
}

func (a *authImpl) SetSmsCode(ctx context.Context, phone string, smsCode string) error {
	return a.rdb.Set(ctx, fmt.Sprintf(consts.RedisUserSmsKey, phone), smsCode, time.Minute).Err()
}

func (a *authImpl) GetSmsCode(ctx context.Context, phone string) (string, error) {
	return a.rdb.Get(ctx, fmt.Sprintf(consts.RedisUserSmsKey, phone)).Result()
}

func (a *authImpl) DeleteSmsCode(ctx context.Context, phone string) error {
	return a.rdb.Del(ctx, fmt.Sprintf(consts.RedisUserSmsKey, phone)).Err()
}

func (a *authImpl) GetEmailCode(ctx context.Context, email string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authImpl) SetEmailCode(ctx context.Context, email string, emailCode string) error {
	//TODO implement me
	panic("implement me")
}

func (a *authImpl) DeleteEmailCode(ctx context.Context, email string) error {
	//TODO implement me
	panic("implement me")
}

func newAuth(rdb *redis.Client) Auth {
	return &authImpl{rdb: rdb}
}
