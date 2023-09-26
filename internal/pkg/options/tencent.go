package options

import "microservices/internal/config"

type TencentOptions struct {
	MailServerAddress  string
	MailServerPassword string
}

func NewTencentOptions() *TencentOptions {
	env := config.NewEnvConfig()
	return &TencentOptions{
		MailServerAddress:  env.TencentMailServerAddress,
		MailServerPassword: env.TencentMailServerPassword,
	}
}
