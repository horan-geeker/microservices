package system

import (
	"context"
	"microservices/pkg/app"
	"microservices/pkg/log"
)

func init() {
	app.RegisterOnce(HealthCheck)
}

func HealthCheck(ctx context.Context) error {
	log.Info(ctx, "command", map[string]any{
		"isHealth": "OK",
	})
	return nil
}
