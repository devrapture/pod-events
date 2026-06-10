package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs all HTTP requests with structured fields.
// This is much more useful than the default Gin logger in production.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process the request
		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logFn := logger.Info
		if statusCode >= 500 {
			logFn = logger.Error
		} else if statusCode >= 400 {
			logFn = logger.Warn
		}

		logFn("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	}
}
