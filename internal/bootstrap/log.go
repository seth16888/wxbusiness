package bootstrap

import (
	"github.com/seth16888/wxbusiness/pkg/logger"
	"go.uber.org/zap"
)

// InitLogger 初始化日志
func InitLogger(conf *logger.LogConfig) *zap.Logger {
	return logger.InitLogger(conf)
}
