package meta

// Config 环境变量映射结构体
type Config struct {
	AppEnv string `mapstructure:"APP_ENV" json:"app_env" default:"development"`

	ServerHost    string `mapstructure:"SERVER_HOST" json:"server_host" default:"127.0.0.1"`
	ServerPort    int    `mapstructure:"SERVER_PORT" json:"server_port,string" default:"80"`
	ServerTimeout int    `mapstructure:"SERVER_TIMEOUT" json:"server_timeout,string" default:"10"`

	DBName     string `mapstructure:"DB_NAME" json:"db_name"`
	DBHost     string `mapstructure:"DB_HOST" json:"db_host"`
	DBPort     int    `mapstructure:"DB_PORT" json:"db_port,string"`
	DBUsername string `mapstructure:"DB_USERNAME" json:"db_username"`
	DBPassword string `mapstructure:"DB_PASSWORD" json:"db_password"`

	RedisHost     string `mapstructure:"REDIS_HOST" json:"redis_host" default:"127.0.0.1"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" json:"redis_password" default:""`
	RedisPort     int    `mapstructure:"REDIS_PORT" json:"redis_port,string" default:"6379"`
	RedisDB       int    `mapstructure:"REDIS_DB" json:"redis_db,string" default:"0"`

	TencentMailServerAddress  string `mapstructure:"TENCENT_MAIL_SERVER_ADDRESS" json:"tencent_mail_server_address"`
	TencentMailServerPassword string `mapstructure:"TENCENT_MAIL_SERVER_PASSWORD" json:"tencent_mail_server_password"`

	AliyunAccessKeyId     string `mapstructure:"ALIYUN_ACCESS_KEY_ID" json:"aliyun_access_key_id"`
	AliyunAccessKeySecret string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET" json:"aliyun_access_key_secret"`
	AliyunSmsSignName     string `mapstructure:"ALIYUN_SMS_SIGN_NAME" json:"aliyun_sms_sign_name"`
	AliyunSmsTemplateCode string `mapstructure:"ALIYUN_SMS_TEMPLATE_CODE" json:"aliyun_sms_template_code"`

	JWTSecret string `mapstructure:"JWT_SECRET" json:"jwt_secret"`
}
