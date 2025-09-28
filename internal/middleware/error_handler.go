package middleware

import (
	"meobeo-talk-api/internal/domain/common"
	"meobeo-talk-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Error("Request error", err)

			common.InternalServerErrorResponse(c)
			return
		}
	}
}
