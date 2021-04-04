package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/kott/go-service-example/pkg/errors"
	rcontext "github.com/kott/go-service-example/pkg/utils/context"
	"github.com/kott/go-service-example/pkg/utils/log"
)

// Recover is a middleware that recovers from panics, logs the panic (and a
// backtrace) using wallet-api logger, and returns a HTTP 500 status.
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := rcontext.GetReqCtx(c)
		defer func() {
			if r := recover(); r != nil {
				logPanic(ctx, r)
				err := errors.NewAppError(errors.InternalServerError,
					errors.Descriptions[errors.InternalServerError], "")
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			}
		}()
		c.Next()
	}
}

func logPanic(ctx context.Context, r interface{}) {
	if err, ok := r.(error); ok {
		log.Error(ctx, "[panic]", err.Error(), "stacktrace", string(debug.Stack()))
	} else {
		log.Error(ctx, "[panic]", r, "stacktrace", string(debug.Stack()))
	}
}
