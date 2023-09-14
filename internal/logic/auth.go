package logic

import (
	"context"
	"fmt"
	"math/rand"
	"microservices/consts"
	"microservices/internal/store"
	"strconv"
	"time"
)

// AuthLogicInterface defines functions used to handle user api.
type AuthLogicInterface interface {
	AuthorizeCookie(ctx context.Context, uid int) (string, error)
	ChangePassword(ctx context.Context, uid int, oldPassword, newPassword string) error
}

type authLogic struct {
	store store.DataFactory
	cache store.CacheFactory
}

// AuthorizeCookie .
func (a *authLogic) AuthorizeCookie(ctx context.Context, uid int) (string, error) {
	a.cache.SetKey(ctx, fmt.Sprintf(consts.RedisUserTokenKey, uid), strconv.Itoa(rand.Int()), time.Hour*2)
	return "", nil
}

// ChangePassword .
func (a *authLogic) ChangePassword(ctx context.Context, uid int, oldPassword, newPassword string) error {
	return nil
}

func newAuth(l *logic) AuthLogicInterface {
	return &authLogic{
		store: l.store,
		cache: l.cache,
	}
}
