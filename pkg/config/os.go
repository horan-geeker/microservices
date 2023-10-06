package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

func ParseConfigFromOSEnv[T any](conf *T) error {
	envTpl := ""
	envs := os.Environ()
	for _, e := range envs {
		envTpl += e + "\n"
	}
	viper.SetConfigType("env")
	if err := viper.ReadConfig(strings.NewReader(envTpl)); err != nil {
		return err
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return err
	}
	return nil
}
