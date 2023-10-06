package logic

import (
	"context"
	"microservices/internal/pkg/meta"
	"microservices/internal/store"
)

// UserLogicInterface defines functions used to handle user api.
type UserLogicInterface interface {
	Create(ctx context.Context, user *meta.User) error
	GetByUid(ctx context.Context, uid uint64) (*meta.User, error)
	List(ctx context.Context) ([]meta.User, error)
	Edit(ctx context.Context, id uint64, name, email, phone *string) error
}

type userLogic struct {
	store store.Factory
}

func newUsers(l *logic) *userLogic {
	return &userLogic{store: l.store}
}

// Create .
func (u *userLogic) Create(ctx context.Context, user *meta.User) error {
	//return u.store.Users().Create(ctx, user)
	return nil
}

// GetByUid .
func (u *userLogic) GetByUid(ctx context.Context, uid uint64) (*meta.User, error) {
	return u.store.Users().GetByUid(ctx, uid)
}

// List .
func (u *userLogic) List(ctx context.Context) ([]meta.User, error) {
	//return u.store.Users().List(ctx)
	return nil, nil
}

// Edit .
func (u *userLogic) Edit(ctx context.Context, id uint64, name, email, phone *string) error {
	user, err := u.store.Users().GetByUid(ctx, id)
	if err != nil {
		return err
	}
	data := make(map[string]any)
	if name != nil {
		data["name"] = *name
	}
	if email != nil {
		data["email"] = *email
	}
	if phone != nil {
		data["phone"] = *phone
	}
	return u.store.Users().Update(ctx, user.ID, data)
}
