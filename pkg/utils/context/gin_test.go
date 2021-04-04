package context

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ctxTest string

func TestGetReqCtx(t *testing.T) {
	gCtx := &gin.Context{}
	assert.Empty(t, GetReqCtx(gCtx))
	assert.Equal(t, context.Background(), GetReqCtx(gCtx))
	ctx := context.WithValue(context.Background(), ctxTest("some_key"), ctxTest("some_value"))
	gCtx.Keys = map[string]interface{}{ctxKey: ctx}
	assert.Equal(t, ctx, GetReqCtx(gCtx))
}

func TestSetReqCtx(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxTest("some_key"), ctxTest("some_value"))
	gCtx := &gin.Context{}
	assert.Empty(t, GetReqCtx(gCtx))
	SetReqCtx(ctx, gCtx)
	assert.Equal(t, ctx, GetReqCtx(gCtx))
}
