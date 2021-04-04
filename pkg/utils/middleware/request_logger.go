package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	rcontext "github.com/kott/go-service-example/pkg/utils/context"
)

// RequestLogger logs before and after an HTTP request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := rcontext.GetRequestLogger(rcontext.GetReqCtx(c))

		start := time.Now().UTC()
		path := c.Request.URL.Path

		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info()

		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"duration":   latency,
			"user_agent": c.Request.UserAgent(),
		}).Info()
	}
}
