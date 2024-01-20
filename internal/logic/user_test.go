package logic

import (
	"context"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"microservices/internal/mock"
	"microservices/internal/store"
	"testing"
)

func TestUsers_GetByUid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := store.NewMockFactory(ctrl)
	mockUserStore := mock.NewMockUserStore(ctrl)
	mockStore.EXPECT().Users().Return(mockUserStore)
	mockUserStore.EXPECT().GetByUid(gomock.Any(), uint64(1)).Return(nil, gorm.ErrRecordNotFound)

	logic := NewLogic(mockStore, nil, nil)
	logic.Users().GetByUid(context.Background(), 1)
}
