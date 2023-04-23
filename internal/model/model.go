package model

import (
	"context"
	"gorm.io/gorm"
	"microservices/internal/config"
	"microservices/internal/mysql"
)

var instance *gorm.DB

// HandlerName 数据库名称
const HandlerName = "databaseName"

// GetDB .
func GetDB(ctx context.Context) *gorm.DB {
	if instance != nil {
		return instance
	}
	instance = mysql.ConnectDB(ctx, config.AppConfig.DBHost, HandlerName, config.AppConfig.DBUsername,
		config.AppConfig.DBPassword, config.AppConfig.DBPort)
	return instance
}
