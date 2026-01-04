package config

type FalOptions struct {
	Key string
}

func NewFalOptions() *FalOptions {
	env := GetEnvConfig()
	return &FalOptions{
		Key: env.FalKey,
	}
}
