package app

type ServerOptions struct {
	Host    string
	Port    int
	Timeout int
	Env     string
}

func NewServerOptions(env, host string, port, timeout int) *ServerOptions {
	return &ServerOptions{
		Host:    host,
		Port:    port,
		Env:     env,
		Timeout: timeout,
	}
}
