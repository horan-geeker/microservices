package http

type Options struct {
	Timeout int
}

func NewOptions(timeout int) *Options {
    return &Options{
        Timeout: timeout,
    }
}