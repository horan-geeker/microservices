package options

import (
	"gorm.io/gorm/logger"
	"time"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Host                  string          `json:"host,omitempty"                     mapstructure:"host"`
	ReplicaHost           string          `json:"replicaHost" mapstructure:"replicaHost"`
	Username              string          `json:"username,omitempty"                 mapstructure:"username"`
	Password              string          `json:"-"                                  mapstructure:"password"`
	Port                  int             `json:"port" mapstructure:"port"`
	Database              string          `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int             `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int             `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration   `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              logger.LogLevel `json:"log-level"                          mapstructure:"log-level"`
	Location              string          `json:"location" mapstructure:"location"`
	Charset               string          `json:"charset"`
}

// NewMySQLOptions create a `zero` value instance.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  3306,
		Username:              "root",
		Password:              "root",
		Database:              "microservice",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              logger.Info, // show all log
		Location:              "Asia/Shanghai",
		Charset:               "utf8mb4",
	}
}
