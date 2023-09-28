package mysql

import (
	"gorm.io/gorm/logger"
	"time"
)

// Options defines consts for mysql database.
type Options struct {
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
