package meta

// Config 环境变量映射结构体
type Config struct {
	AppEnv string `mapstructure:"APP_ENV" default:"development"`

	ServerHost    string `mapstructure:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort    int    `mapstructure:"SERVER_PORT" default:"80"`
	ServerTimeout int    `mapstructure:"SERVER_TIMEOUT" default:"10"`

	DBName     string `mapstructure:"DB_NAME"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisHost     string `mapstructure:"REDIS_HOST" default:"127.0.0.1"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" default:""`
	RedisPort     int    `mapstructure:"REDIS_PORT" default:"6379"`
	RedisDB       int    `mapstructure:"REDIS_DB" default:"0"`

	TencentMailServerAddress  string `mapstructure:"TENCENT_MAIL_SERVER_ADDRESS"`
	TencentMailServerPassword string `mapstructure:"TENCENT_MAIL_SERVER_PASSWORD"`

	AliyunAccessKeyId     string `mapstructure:"ALIYUN_ACCESS_KEY_ID"`
	AliyunAccessKeySecret string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET"`
	AliyunSmsSignName     string `mapstructure:"ALIYUN_SMS_SIGN_NAME"`
	AliyunSmsTemplateCode string `mapstructure:"ALIYUN_SMS_TEMPLATE_CODE"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
}
