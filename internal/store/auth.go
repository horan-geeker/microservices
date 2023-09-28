package store

import "context"

type AuthCache interface {
	SetSmsCode(ctx context.Context, phone string, smsCode string) error
	GetSmsCode(ctx context.Context, phone string) (string, error)
	DeleteSmsCode(ctx context.Context, phone string) error
	GetEmailCode(ctx context.Context, email string) (string, error)
	SetEmailCode(ctx context.Context, email string, emailCode string) error
	DeleteEmailCode(ctx context.Context, email string) error
}
