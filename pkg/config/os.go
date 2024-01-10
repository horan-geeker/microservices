package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

type EnvConfigInstance struct {
}

func (f *EnvConfigInstance) GetContent() (map[string]string, error) {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			envMap[kv[0]] = kv[1]
		}
	}
	return envMap, nil
}

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
