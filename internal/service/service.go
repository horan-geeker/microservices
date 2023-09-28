package service

import (
	"microservices/internal/pkg/options"
)

type Factory interface {
	Tencent() Tencent
	Aliyun() Aliyun
}

type serviceInstance struct {
	aliyunOptions  *options.AliyunOptions
	tencentOptions *options.TencentOptions
}

func (s *serviceInstance) Tencent() Tencent {
	return newTencent(s)
}

func (s *serviceInstance) Aliyun() Aliyun {
	return newAliyun(s)
}

func GetServiceInstance(tencentOpt *options.TencentOptions, aliyunOpt *options.AliyunOptions) Factory {
	return &serviceInstance{
		tencentOptions: tencentOpt,
		aliyunOptions:  aliyunOpt,
	}
}
