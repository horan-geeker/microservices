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
	instance = mysql.ConnectDB(ctx, config.Env.DBHost, HandlerName, config.Env.DBUsername,
		config.Env.DBPassword, config.Env.DBPort)
	return instance
}
