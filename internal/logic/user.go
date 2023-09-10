package logic

import (
	"context"
	"microservices/internal/model"
	"microservices/internal/store"
)

// UserLogicInterface defines functions used to handle user api.
type UserLogicInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByUserName(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
}

type userLogic struct {
	store store.Factory
}

func newUsers(l *logic) *userLogic {
	return &userLogic{store: l.store}
}

// Create .
func (u *userLogic) Create(ctx context.Context, user *model.User) error {
	//return u.store.Users().Create(ctx, user)
	return nil
}

// GetByUserName .
func (u *userLogic) GetByUserName(ctx context.Context, username string) (*model.User, error) {
	return u.store.Users().GetByUserName(ctx, username)
}

// List .
func (u *userLogic) List(ctx context.Context) ([]model.User, error) {
	//return u.store.Users().List(ctx)
	return nil, nil
}
