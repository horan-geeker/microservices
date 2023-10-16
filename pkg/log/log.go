package log

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func Info(ctx context.Context, event string, ext map[string]any) {
	log.WithFields(makeFields(ctx, event, ext)).Info()
}

func Warning(ctx context.Context, event string, message string, ext map[string]any) {
	ext["message"] = message
	log.WithFields(makeFields(ctx, event, ext)).Error()
}

func Error(ctx context.Context, event string, err error, ext map[string]any) {
	ext["err"] = err
	log.WithFields(makeFields(ctx, event, ext)).Error()
}

func makeFields(ctx context.Context, event string, ext map[string]any) log.Fields {
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
	return fields
}
