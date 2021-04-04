package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	_ contextKey = iota
	requestLoggerKey
	reqIDKey
)

//SetRequestLogger sets the logger on context
func SetRequestLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, logger)
}

//GetRequestLogger returns a Logger
func GetRequestLogger(ctx context.Context) logrus.FieldLogger {
	logger := ctx.Value(requestLoggerKey)
	if logger != nil {
		return logger.(logrus.FieldLogger)
	}
	return logrus.New()
}

//SetReqID sets the reqID
func SetReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

//GetReqID retrieves the reqID
func GetReqID(ctx context.Context) string {
	return getContextStringValue(ctx, reqIDKey)
}

func getContextStringValue(ctx context.Context, key contextKey) string {
	value, ok := ctx.Value(key).(string)
	if ok {
		return value
	}
	return ""
}
