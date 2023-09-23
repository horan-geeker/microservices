package store

import (
	"context"
	"microservices/internal/model"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	GetByUid(ctx context.Context, id uint64) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	//List(ctx context.Context) ([]model.User, error)
}

type UserCache interface {
	SetToken(ctx context.Context, uid uint64, token string) error
	GetToken(ctx context.Context, uid uint64) (string, error)
}
