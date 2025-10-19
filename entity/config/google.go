package config

type GoogleOptions struct {
	ProxyURL string
}

// NewGoogleOptions .
func NewGoogleOptions() *GoogleOptions {
	return &GoogleOptions{}
}
