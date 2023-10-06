package alarm

import (
	"context"
	"microservices/pkg/http"
)

type alarm struct {
	wecomBotUrl string
}

func (a *alarm) SendWeComBot(ctx context.Context, msgType string, message string) error {
	_, err := http.NewHttp[any](http.NewOptions(5)).Post(ctx, a.wecomBotUrl, nil, map[string]any{
		"msgtype":  msgType,
		"markdown": map[string]any{"content": message},
	}, nil)
	return err
}

func NewAlarm(opts *Options) *alarm {
	return &alarm{
		wecomBotUrl: opts.WeComBot,
	}
}
