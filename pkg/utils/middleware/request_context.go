package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	rcontext "github.com/kott/go-service-example/pkg/utils/context"
	"github.com/kott/go-service-example/pkg/utils/log"
)

const reqIDHeader = "X-Request-Id"

// PersistContext sets any values we want persisted throughout the life of a request
func PersistContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := currentReqID(c)
		ctxLogger := log.New().WithField("reqID", reqID)

		ctx := context.Background()
		ctx = rcontext.SetRequestLogger(ctx, ctxLogger)
		ctx = rcontext.SetReqID(ctx, reqID)
		rcontext.SetReqCtx(ctx, c)

		c.Header(reqIDHeader, reqID)
		c.Next()
	}
}

func currentReqID(c *gin.Context) string {
	var reqID string
	if reqID = c.GetHeader(reqIDHeader); reqID == "" {
		reqID = rcontext.GenerateReqID()
	}
	return reqID
}
