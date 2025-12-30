package service

import "microservices/entity/config"

var srv Factory

type Factory interface {
	Tencent() Tencent
	Aliyun() Aliyun
	Google() Google
	Alipay() Alipay
	Apple() Apple
	Stripe() Stripe
}

type factory struct {
	tencentOpt *config.TencentOptions
	aliyunOpt  *config.AliyunOptions
	googleOpt  *config.GoogleOptions
	alipayOpt  *config.AlipayOptions
	appleOpt   *config.AppleOptions
	stripeOpt  *config.StripeOptions
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

func (f *factory) Alipay() Alipay {
	return newAlipay(f.alipayOpt)
}

func (f *factory) Apple() Apple {
	return newApple(f.appleOpt)
}

func (f *factory) Stripe() Stripe {
	return newStripe(f.stripeOpt)
}

// NewFactory .
func NewFactory() Factory {
	if srv == nil {
		srv = &factory{
			tencentOpt: config.NewTencentOptions(),
			aliyunOpt:  config.NewAliyunOptions(),
			googleOpt:  config.NewGoogleOptions(),
			alipayOpt:  config.NewAlipayOptions(),
			appleOpt:   config.NewAppleOptions(),
			stripeOpt:  config.NewStripeOptions(),
		}
	}
	return srv
}
