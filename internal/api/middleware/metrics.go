package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
)

func MetricsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			status,
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(time.Since(start).Seconds())
	}
}
