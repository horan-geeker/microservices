package config

import (
	"github.com/spf13/viper"
)

// EnvConfig 环境变量映射结构体
type EnvConfig struct {
	AppEnv string `mapstructure:"APP_ENV" default:"development"`

	ServerHost string `mapstructure:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort int    `mapstructure:"SERVER_PORT" default:"80"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisHost     string `mapstructure:"REDIS_HOST" default:"127.0.0.1"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" default:""`
	RedisPort     int    `mapstructure:"REDIS_PORT" default:"6379"`
	RedisDB       int    `mapstructure:"REDIS_DB" default:"0"`

	MailServerAddress  string `mapstructure:"MAIL_SERVER_ADDRESS"`
	MailServerPassword string `mapstructure:"MAIL_SERVER_PASSWORD"`

	AliyunAccessKeyId     string `mapstructure:"ALIYUN_ACCESS_KEY_ID"`
	AliyunAccessKeySecret string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET"`
	AliyunSmsSignName     string `mapstructure:"ALIYUN_SMS_SIGN_NAME"`
	AliyunSmsTemplateCode string `mapstructure:"ALIYUN_SMS_TEMPLATE_CODE"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
}

// Env 系统配置
var env *EnvConfig

func NewEnvConfig() *EnvConfig {
	if env == nil {
		viper.AddConfigPath(".")
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&env); err != nil {
			panic(err)
		}
	}
	return env
}
