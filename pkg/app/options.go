package app

type ServerOptions struct {
	Timeout int
}

func NewServerOptions(timeout int) *ServerOptions {
	return &ServerOptions{Timeout: timeout}
}
