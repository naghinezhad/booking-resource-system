package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/internal/utils/concurrency"
)

func RequestLimitMiddleware(l *concurrency.RequestLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !l.Acquire() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "server overloaded, try again",
			})
			c.Abort()
			return
		}

		defer l.Release()
		c.Next()
	}
}
