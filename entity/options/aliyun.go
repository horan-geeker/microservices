package options

import (
	"microservices/entity/config"
)

type AliyunOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	SmsSignName     string
	SmsTemplateCode string
}

func NewAliyunOptions() *AliyunOptions {
	env := config.GetConfig()
	return &AliyunOptions{
		AccessKeyId:     env.AliyunAccessKeyId,
		AccessKeySecret: env.AliyunAccessKeySecret,
		SmsSignName:     env.AliyunSmsSignName,
		SmsTemplateCode: env.AliyunSmsTemplateCode,
	}
}
