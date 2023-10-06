package alarm

type Options struct {
	WeComBot string
}

func NewOptions(weComBot string) *Options {
	return &Options{
		WeComBot: weComBot,
	}
}
