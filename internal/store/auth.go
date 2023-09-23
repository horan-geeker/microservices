package store

import "context"

type AuthCache interface {
	SetSmsCode(ctx context.Context, uid uint64, smsCode string) error
	GetSmsCode(ctx context.Context, uid uint64) (string, error)
	DeleteSmsCode(ctx context.Context, uid uint64) error
}
