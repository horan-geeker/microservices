package options

import (
	"microservices/internal/config"
)

type TencentOptions struct {
	MailServerAddress  string
	MailServerPassword string
}

func NewTencentOptions() *TencentOptions {
	env := config.GetConfig()
	return &TencentOptions{
		MailServerAddress:  env.TencentMailServerAddress,
		MailServerPassword: env.TencentMailServerPassword,
	}
}
