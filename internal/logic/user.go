package logic

import (
	"context"
	"microservices/internal/store"
)

// UserLogicInterface defines functions used to handle user request.
type UserLogicInterface interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
}

type userLogic struct {
	store store.Factory
}

func newUsers(l *logic) *userLogic {
	return &userLogic{store: l.store}
}

// Create .
func (u *userLogic) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.store.Users().Create(ctx, user, opts)
}

// Update .
func (u *userLogic) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.store.Users().Update(ctx, user, opts)
}

// Delete .
func (u *userLogic) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	return u.store.Users().Delete(ctx, username, opts)
}

// Get .
func (u *userLogic) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	return u.store.Users().Get(ctx, username, opts)
}

// List .
func (u *userLogic) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	return u.store.Users().List(ctx, opts)
}
