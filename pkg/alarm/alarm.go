package alarm

import (
	"context"
	"microservices/pkg/http"
)

type alarm struct {
	wecomBotUrl string
}

func (a *alarm) SendWeComBot(ctx context.Context, msgType string, message string) error {
	return http.NewClient().Post(ctx, a.wecomBotUrl, nil, map[string]any{
		"msgtype":  msgType,
		"markdown": map[string]any{"content": message},
	}, nil)
}

func NewAlarm(opts *Options) *alarm {
	return &alarm{
		wecomBotUrl: opts.WeComBot,
	}
}
