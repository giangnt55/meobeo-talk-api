package middleware

import (
	"meobeo-talk-api/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", latency,
			"ip", c.ClientIP(),
		)
	}
}
