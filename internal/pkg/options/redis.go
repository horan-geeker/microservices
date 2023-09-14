package options

type RedisOptions struct {
	Host string `json:"host,omitempty"                     mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host: "127.0.0.1",
		Port: 6379,
	}
}
