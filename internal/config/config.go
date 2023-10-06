package config

import (
	"microservices/internal/pkg/meta"
	"microservices/pkg/config"
)

var conf *meta.Config

func NewConfig() *meta.Config {
	if conf == nil {
		conf = &meta.Config{}
		if err := config.NewConfig(conf); err != nil {
			panic(err)
		}
	}
	return conf
}
