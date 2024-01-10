package config

import (
	"github.com/mitchellh/mapstructure"
	"microservices/internal/pkg/meta"
	"microservices/pkg/config"
)

var (
	conf *meta.Config
)

func GetConfig() *meta.Config {
	if conf == nil {
		config.RegisterProvider(&config.FileConfigInstance{})
		content, err := config.GetConfig()
		if err != nil {
			panic(err)
		}
		if err := mapstructure.WeakDecode(content, &conf); err != nil {
			panic(err)
		}
	}
	return conf
}
