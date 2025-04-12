package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware 是一个Gin的日志中间件，用于记录每个请求的基本信息
func LoggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Content-Type
		contentType := c.Request.Header.Get("Content-Type")
		// HTTP 协议版本
		proto := c.Request.Proto
		// User-Agent
		userAgent := c.Request.UserAgent()

		// 处理请求
		c.Next()

		// Context Keys
		keys := c.Keys
		// Request ID
		requestID := c.Request.Header.Get("X-Request-Id")
		// 记录请求结束时间
		end := time.Now()
		// 计算请求处理时间
		latency := end.Sub(start).Milliseconds()
		// 获取客户端IP地址
		clientIP := c.ClientIP()
		// 获取请求方法和请求路径
		method := c.Request.Method
		// 获取响应状态码
		statusCode := c.Writer.Status()
		// ErrorMessage
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		// BodySize
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// 优化日志输出格式，使用更具可读性的格式
		log.Info(
			"request",
			zap.String("method", method),
			zap.Int("status", statusCode),
			zap.String("path", path),
			zap.String("cip", clientIP),
			zap.String("requestId", requestID),
			zap.String("ua", userAgent),
			zap.String("proto", proto),
			zap.String("contentType", contentType),
			zap.Int("len", bodySize),
			zap.String("err", errorMessage),
			zap.Any("keys", keys),
			zap.String("received", start.Format("2006/01/02 15:04:05")),
			zap.Int64("latency(ms)", latency),
		)
	}
}
