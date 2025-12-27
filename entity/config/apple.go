package config

// AppleOptions captures configuration required for handling Apple in-app purchase callbacks.
type AppleOptions struct {
	SharedSecret string
}

func NewAppleOptions() *AppleOptions {
	c := GetEnvConfig()
	return &AppleOptions{
		SharedSecret: c.AppleSharedSecret,
	}
}
