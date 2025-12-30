package config

type StripeOptions struct {
	ApiKey        string
	WebhookSecret string
	Domain        string
}

func NewStripeOptions() *StripeOptions {
	env := GetEnvConfig()
	return &StripeOptions{
		ApiKey:        env.StripeSecretKey,
		WebhookSecret: env.StripeWebhookSecret,
		Domain:        env.AppDomain,
	}
}
