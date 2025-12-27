package config

type AlipayOptions struct {
	AppId                string
	PrivateKey           string
	AlipayPublicKey      string
	AlipayAppPublicKey   string
	AppPublicCertPath    string
	AlipayRootCertPath   string
	AlipayPublicCertPath string
	NotifyUrl            string
	ReturnUrl            string
	IsProduction         bool
}

func NewAlipayOptions() *AlipayOptions {
	c := GetEnvConfig()
	return &AlipayOptions{
		AppId:                c.AlipayAppId,
		PrivateKey:           c.AlipayPrivateKey,
		AlipayPublicKey:      c.AlipayPublicKey,
		AlipayAppPublicKey:   c.AlipayAppPublicKey,
		AppPublicCertPath:    c.AlipayAppPublicCertPath,
		AlipayRootCertPath:   c.AlipayRootCertPath,
		AlipayPublicCertPath: c.AlipayPublicCertPath,
		NotifyUrl:            c.AlipayNotifyUrl,
		ReturnUrl:            c.AlipayReturnUrl,
		IsProduction:         c.AlipayIsProduction,
	}
}
