package service

import "microservices/entity/config"

var srv Factory

type Factory interface {
	Tencent() Tencent
	Aliyun() Aliyun
}

type factory struct {
	tencentOpt *config.TencentOptions
	aliyunOpt  *config.AliyunOptions
}

func (f *factory) Tencent() Tencent {
	return newTencent(f.tencentOpt)
}

func (f *factory) Aliyun() Aliyun {
	return newAliyun(f.aliyunOpt)
}

// NewFactory .
func NewFactory() Factory {
	if srv == nil {
		srv = &factory{
			tencentOpt: config.NewTencentOptions(),
			aliyunOpt:  config.NewAliyunOptions(),
		}
	}
	return srv
}
