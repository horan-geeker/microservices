package meta

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"os"
	"time"
)

type gormLogger struct {
	Logger *logrus.Logger
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
		l.Logger.Printf("fffff", data)
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.Printf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.Printf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	requestId, _ := ctx.Value("requestId").(string)
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.Logger.WithFields(logrus.Fields{
		"requestId": requestId,
		"sql":       sql,
		"rows":      rows,
		"timeCost":  elapsed.Milliseconds(),
	}).Info()
}

func NewGormCustomLogger(logLevel logger.LogLevel) logger.Interface {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.WithField("event", "mysql-query")
	return &gormLogger{
		log,
		logger.Config{
			LogLevel:                  logLevel, // Log level
			IgnoreRecordNotFoundError: true,     // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,    // Disable color
		},
	}
}
