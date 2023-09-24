package meta

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type RedisLogHook struct{}

func (r *RedisLogHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		traceId, _ := ctx.Value("traceId").(string)
		spanId, _ := ctx.Value("spanId").(string)
		logFields := logrus.Fields{
			"event":   "redis",
			"traceId": traceId,
			"spanId":  spanId,
			"network": network,
			"address": addr,
		}
		conn, err := next(ctx, network, addr)
		if err != nil {
			logrus.WithFields(logFields).Error()
			return conn, err
		}
		logrus.WithFields(logFields).Info()
		return conn, nil
	}
}

func (r *RedisLogHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		traceId, _ := ctx.Value("traceId").(string)
		spanId, _ := ctx.Value("spanId").(string)
		begin := time.Now()
		err := next(ctx, cmd)
		logFields := logrus.Fields{
			"event":    "redis",
			"traceId":  traceId,
			"spanId":   spanId,
			"cmd":      cmd.String(),
			"timeCost": time.Since(begin).Milliseconds(),
		}
		if err != nil {
			logFields["error"] = err.Error()
			logrus.WithFields(logFields).Error()
		} else {
			logrus.WithFields(logFields).Info()
		}
		return err
	}
}

func (r *RedisLogHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func NewRedisLogHook() *RedisLogHook {
	return &RedisLogHook{}
}
