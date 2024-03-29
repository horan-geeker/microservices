package config

import (
	"gorm.io/gorm/logger"
	"microservices/pkg/mysql"
	"time"
)

// NewMysqlOptions create a `zero` value instance.
func NewMysqlOptions() *mysql.Options {
	env := GetEnvConfig()
	return &mysql.Options{
		Host:                  env.DBHost,
		Port:                  env.DBPort,
		Username:              env.DBUsername,
		Password:              env.DBPassword,
		Database:              "microservice",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              logger.Info, // show all log
		Location:              "Asia/Shanghai",
		Charset:               "utf8mb4",
	}
}
