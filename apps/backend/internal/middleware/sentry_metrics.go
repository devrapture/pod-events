package middleware

import (
	"time"

	"github.com/devrapture/pod-events/internal/metrics"
	"github.com/getsentry/sentry-go/attribute"
	"github.com/gin-gonic/gin"
)

// MetricsRecorder records HTTP request latency and count metrics. It is safe
// to register unconditionally; the recorder is a no-op when metrics are
// disabled.
func MetricsRecorder(recorder metrics.Recorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		attrs := []attribute.Builder{
			attribute.String("route", route),
			attribute.String("method", c.Request.Method),
			attribute.Int("status_code", c.Writer.Status()),
		}
		recorder.Distribution(
			"http.request.duration",
			float64(time.Since(start).Milliseconds()),
			metrics.WithUnit(metrics.UnitMillisecond),
			metrics.WithAttributes(attrs...),
		)
		recorder.Count("http.request.count", 1, metrics.WithAttributes(attrs...))
	}
}
