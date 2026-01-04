package config

type CloudflareOptions struct {
	AccountId       string
	AccessKeyId     string
	SecretAccessKey string
	Bucket          string
	PublicDomain    string
}

func NewCloudflareOptions() *CloudflareOptions {
	env := GetEnvConfig()
	return &CloudflareOptions{
		AccountId:       env.CloudflareAccountId,
		AccessKeyId:     env.CloudflareAccessKeyId,
		SecretAccessKey: env.CloudflareSecretAccessKey,
		Bucket:          env.CloudflareBucket,
		PublicDomain:    env.CloudflarePublicDomain,
	}
}
