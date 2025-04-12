package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// 初始化日志
func InitLogger(cfg *LogConfig) *zap.Logger {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Filename)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建Core
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台输出
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 文件输出
	fileHandle, _ := os.OpenFile(cfg.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(fileHandle),
		level,
	)

	// 创建Logger
	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	Log = logger
	return logger
}

// 以下是封装的常用日志方法
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debug(fmt.Sprintf(format, args...))
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	Log.Info(fmt.Sprintf(format, args...))
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warn(fmt.Sprintf(format, args...))
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	Log.Error(fmt.Sprintf(format, args...))
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatal(fmt.Sprintf(format, args...))
}

func Sync() {
	Log.Sync()
}
