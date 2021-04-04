package context

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ctxKey defines how we label context in our requests
const ctxKey = "ctx"

// GetReqCtx obtains the context from the http request in gin
func GetReqCtx(c *gin.Context) context.Context {
	rCtx, exists := c.Get(ctxKey)
	if !exists {
		return context.Background()
	}
	return rCtx.(context.Context)
}

// SetReqCtx sets the general context to gin's context
func SetReqCtx(ctx context.Context, c *gin.Context) {
	c.Set(ctxKey, ctx)
}
