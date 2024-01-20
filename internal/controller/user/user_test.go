package user

import (
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"microservices/internal/mock"
	"microservices/internal/pkg/options"
	"microservices/internal/service"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"testing"
)

func TestUserLogic_GetByUid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言方法是否被调用
	mockUserStore := mock.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().GetByUid(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	userController := NewUserController(
		mysql.GetMysqlInstance(options.NewMysqlOptions()),
		redis.GetRedisInstance(options.NewRedisOptions()),
		service.GetServiceInstance(options.NewTencentOptions(), options.NewAliyunOptions()))
	userController.Get(nil, 1)
}
