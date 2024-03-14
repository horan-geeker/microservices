package config

type Provider interface {
	GetContent() (map[string]string, error)
}

var provider Provider

func GetConfig() (map[string]string, error) {
	return provider.GetContent()
}

func RegisterProvider(instance Provider) {
	provider = instance
}
