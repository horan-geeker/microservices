package config

type TencentOptions struct {
	MailServerAddress  string
	MailServerPassword string
}

func NewTencentOptions() *TencentOptions {
	env := GetEnvConfig()
	return &TencentOptions{
		MailServerAddress:  env.TencentMailServerAddress,
		MailServerPassword: env.TencentMailServerPassword,
	}
}
