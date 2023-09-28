package log

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func Log(ctx context.Context, event string, ext map[string]any) {
	traceId, _ := ctx.Value("traceId").(string)
	spanId, _ := ctx.Value("spanId").(string)
	fields := log.Fields{
		"traceId": traceId,
		"spanId":  spanId,
		"event":   event,
	}
	for k, v := range ext {
		fields[k] = v
	}
	log.WithFields(fields).Info()
}
