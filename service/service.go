package service

import "microservices/entity/config"

type Factory interface {
	Tencent() Tencent
	Aliyun() Aliyun
}

type factoryImpl struct {
	tencentOpt *config.TencentOptions
	aliyunOpt  *config.AliyunOptions
}

func (f *factoryImpl) Tencent() Tencent {
	return newTencent(f.tencentOpt)
}

func (f *factoryImpl) Aliyun() Aliyun {
	return newAliyun(f.aliyunOpt)
}

// NewFactory .
func NewFactory() Factory {
	return &factoryImpl{
		tencentOpt: config.NewTencentOptions(),
		aliyunOpt:  config.NewAliyunOptions(),
	}
}
