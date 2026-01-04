package config

import (
	"microservices/pkg/config"

	"github.com/mitchellh/mapstructure"
)

// Config 环境变量映射结构体
type Config struct {
	AppEnv    string `mapstructure:"APP_ENV" json:"app_env" default:"development"`
	AppDomain string `mapstructure:"APP_DOMAIN" json:"app_domain" default:"localhost"`

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
	COSBucketDomain           string `mapstructure:"COS_BUCKET_DOMAIN" json:"cos_bucket_domain"`
	COSSecretId               string `mapstructure:"COS_SECRET_ID" json:"cos_secret_id"`
	COSSecretKey              string `mapstructure:"COS_SECRET_KEY" json:"cos_secret_key"`

	GoogleClientId     string `mapstructure:"GOOGLE_CLIENT_ID" json:"google_client_id"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET" json:"google_client_secret"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL" json:"google_redirect_url"`

	AliyunAccessKeyId     string `mapstructure:"ALIYUN_ACCESS_KEY_ID" json:"aliyun_access_key_id"`
	AliyunAccessKeySecret string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET" json:"aliyun_access_key_secret"`
	AliyunSmsSignName     string `mapstructure:"ALIYUN_SMS_SIGN_NAME" json:"aliyun_sms_sign_name"`
	AliyunSmsTemplateCode string `mapstructure:"ALIYUN_SMS_TEMPLATE_CODE" json:"aliyun_sms_template_code"`

	AlipayAppId             string `mapstructure:"ALIPAY_APP_ID" json:"alipay_app_id"`
	AlipayPrivateKey        string `mapstructure:"ALIPAY_PRIVATE_KEY" json:"alipay_private_key"`
	AlipayPublicKey         string `mapstructure:"ALIPAY_PUBLIC_KEY" json:"alipay_public_key"`
	AlipayAppPublicKey      string `mapstructure:"ALIPAY_APP_PUBLIC_KEY" json:"alipay_app_public_key"`
	AlipayAppPublicCertPath string `mapstructure:"ALIPAY_APP_PUBLIC_CERT_PATH" json:"alipay_app_public_cert_path"`
	AlipayRootCertPath      string `mapstructure:"ALIPAY_ROOT_CERT_PATH" json:"alipay_root_cert_path"`
	AlipayPublicCertPath    string `mapstructure:"ALIPAY_PUBLIC_CERT_PATH" json:"alipay_public_cert_path"`
	AlipayNotifyUrl         string `mapstructure:"ALIPAY_NOTIFY_URL" json:"alipay_notify_url"`
	AlipayReturnUrl         string `mapstructure:"ALIPAY_RETURN_URL" json:"alipay_return_url"`
	AlipayIsProduction      bool   `mapstructure:"ALIPAY_IS_PRODUCTION" json:"alipay_is_production"`

	StripeSecretKey     string `mapstructure:"STRIPE_SECRET_KEY" json:"stripe_secret_key"`
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET" json:"stripe_webhook_secret"`

	AppleSharedSecret string `mapstructure:"APPLE_SHARED_SECRET" json:"apple_shared_secret"`

	JWTSecret string `mapstructure:"JWT_SECRET" json:"jwt_secret"`

	CloudflareAccountId       string `mapstructure:"CLOUDFLARE_ACCOUNT_ID" json:"cloudflare_account_id"`
	CloudflareAccessKeyId     string `mapstructure:"CLOUDFLARE_ACCESS_KEY_ID" json:"cloudflare_access_key_id"`
	CloudflareSecretAccessKey string `mapstructure:"CLOUDFLARE_SECRET_ACCESS_KEY" json:"cloudflare_secret_access_key"`
	CloudflareBucket          string `mapstructure:"CLOUDFLARE_BUCKET" json:"cloudflare_bucket"`
	CloudflarePublicDomain    string `mapstructure:"CLOUDFLARE_PUBLIC_DOMAIN" json:"cloudflare_public_domain"`

	FalKey string `mapstructure:"FAL_KEY" json:"fal_key"`
}

var conf *Config

// GetEnvConfig .
func GetEnvConfig() *Config {
	if conf == nil {
		fileConfigInstance := config.FileConfigInstance{}
		envConfigInstance := config.EnvConfigInstance{}
		fileConfig, err := fileConfigInstance.GetContent()
		if err != nil {
			panic(err)
		}
		envConfig, err := envConfigInstance.GetContent()
		if err != nil {
			panic(err)
		}
		for k, v := range envConfig {
			fileConfig[k] = v
		}
		if err := mapstructure.WeakDecode(fileConfig, &conf); err != nil {
			panic(err)
		}
	}
	return conf
}
