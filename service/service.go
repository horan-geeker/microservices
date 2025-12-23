package service

import "microservices/entity/config"

var srv Factory

type Factory interface {
	Tencent() Tencent
	Aliyun() Aliyun
	Google() Google
}

type factory struct {
	tencentOpt *config.TencentOptions
	aliyunOpt  *config.AliyunOptions
	googleOpt  *config.GoogleOptions
}

func (f *factory) Tencent() Tencent {
	return newTencent(f.tencentOpt)
}

func (f *factory) Aliyun() Aliyun {
	return newAliyun(f.aliyunOpt)
}

func (f *factory) Google() Google {
	return newGoogle(f.googleOpt)
}

// NewFactory .
func NewFactory() Factory {
	if srv == nil {
		srv = &factory{
			tencentOpt: config.NewTencentOptions(),
			aliyunOpt:  config.NewAliyunOptions(),
			googleOpt:  config.NewGoogleOptions(),
		}
	}
	return srv
}
