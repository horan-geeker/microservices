package options

import (
	"microservices/internal/config"
)

type AliyunOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	SmsSignName     string
	SmsTemplateCode string
}

func NewAliyunOptions() *AliyunOptions {
	env := config.NewConfig()
	return &AliyunOptions{
		AccessKeyId:     env.AliyunAccessKeyId,
		AccessKeySecret: env.AliyunAccessKeySecret,
		SmsSignName:     env.AliyunSmsSignName,
		SmsTemplateCode: env.AliyunSmsTemplateCode,
	}
}
