package database

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/pkg/logger"

	"go.uber.org/zap"
	loggerEx "gorm.io/gorm/logger"
)

// GormLogger 实现 gorm.Logger 接口
type GormLogger struct {
	LogLevel loggerEx.LogLevel
}

// NewGormLogger 创建新的 GORM 日志适配器
func NewGormLogger(level string) *GormLogger {
	var logLevel loggerEx.LogLevel
	switch level {
	case "silent":
		logLevel = loggerEx.Silent
	case "error":
		logLevel = loggerEx.Error
	case "warn":
		logLevel = loggerEx.Warn
	case "info":
		logLevel = loggerEx.Info
	default:
		logLevel = loggerEx.Info
	}

	return &GormLogger{
		LogLevel: logLevel,
	}
}

// LogMode 实现 gorm.Logger 接口
func (l *GormLogger) LogMode(level loggerEx.LogLevel) loggerEx.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 实现 gorm.Logger 接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerEx.Info {
		logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 实现 gorm.Logger 接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerEx.Warn {
		logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 实现 gorm.Logger 接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerEx.Error {
		logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 实现 gorm.Logger 接口
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= loggerEx.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 根据不同情况记录日志
	switch {
	case err != nil && l.LogLevel >= loggerEx.Error:
		logger.Error("GORM SQL Error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	case elapsed > 200*time.Millisecond && l.LogLevel >= loggerEx.Warn:
		logger.Warn("GORM Slow SQL",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	case l.LogLevel >= loggerEx.Info:
		logger.Info("GORM SQL",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	}
}
