package config

import (
	"microservices/entity/meta"
	"microservices/pkg/config"
)

var conf *meta.Config

func GetConfig() *meta.Config {
	if conf == nil {
		conf = &meta.Config{}
		if err := config.NewConfig(conf); err != nil {
			panic(err)
		}
	}
	return conf
}
