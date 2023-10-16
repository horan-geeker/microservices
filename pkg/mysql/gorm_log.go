package mysql

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"microservices/pkg/log"
	"time"
)

type gormLogger struct {
	logger.Config
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Info(ctx, "mysql-info", map[string]any{
			"msg": fmt.Sprintf(msg, data...),
		})
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Warning(ctx, "mysql-warn", fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)), nil)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Error(ctx, "mysql-error", errors.New(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...))), nil)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	log.Info(ctx, "mysql", map[string]any{
		"sql":          sql,
		"rowsAffected": rowsAffected,
		"timeCost":     elapsed.Milliseconds(),
	})
}

func NewGormCustomLogger(logLevel logger.LogLevel) logger.Interface {
	return &gormLogger{
		logger.Config{
			LogLevel:                  logLevel, // Log level
			IgnoreRecordNotFoundError: true,     // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,    // Disable color
		},
	}
}
