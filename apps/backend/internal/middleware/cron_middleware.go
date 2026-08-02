package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CronMiddleware(secret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		provided := c.GetHeader("X-Cron-Secret")
		if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
			logger.Warn(
				"unauthorized cron attempt",
				zap.String("ip", c.ClientIP()),
				zap.Bool("header_present", provided != ""),
			)
			response.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
