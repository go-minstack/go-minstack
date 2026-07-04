package mysql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	log           *slog.Logger
	slowThreshold time.Duration
	logLevel      logger.LogLevel
}

func newGormLogger(log *slog.Logger, slowThreshold time.Duration) logger.Interface {
	return &gormLogger{log: log, slowThreshold: slowThreshold, logLevel: logger.Warn}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	n := *l
	n.logLevel = level
	return &n
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.InfoContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.WarnContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.ErrorContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		l.log.ErrorContext(ctx, "SQL error",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.Any("error", err),
		)
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		l.log.WarnContext(ctx, "Slow SQL query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
		)
	default:
		l.log.DebugContext(ctx, "SQL query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
		)
	}
}
