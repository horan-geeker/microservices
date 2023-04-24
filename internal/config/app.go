package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 环境变量映射结构体
type Config struct {
	AppEnv string `mapstructure:"APP_ENV" default:"development"`

	ServerHost string `mapstructure:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort int    `mapstructure:"SERVER_PORT" default:"80"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisHost     string `mappstructure:"REDIS_HOST" default:"127.0.0.1"`
	RedisPassword string `mappstructure:"REDIS_PASSWORD" default:""`
	RedisPort     int    `mappstructure:"REDIS_PORT" default:"6379"`
	RedisDB       int    `mappstructure:"REDIS_DB" default:"0"`
}

// Env 系统配置
var Env Config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Env); err != nil {
		panic(err)
	}
	fmt.Print(Env)
}
