package store

import (
	"context"
	"microservices/internal/pkg/meta"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *meta.User) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	GetByUid(ctx context.Context, id uint64) (*meta.User, error)
	GetByName(ctx context.Context, name string) (*meta.User, error)
	GetByEmail(ctx context.Context, email string) (*meta.User, error)
	GetByPhone(ctx context.Context, phone string) (*meta.User, error)
	//List(ctx context.Context) ([]model.User, error)
}

type UserCache interface {
	SetToken(ctx context.Context, uid uint64, token string) error
	GetToken(ctx context.Context, uid uint64) (string, error)
	DeleteToken(ctx context.Context, uid uint64) error
}
