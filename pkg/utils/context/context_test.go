package context

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetGetRequestLogger(t *testing.T) {
	ctx := context.Background()
	newLogger := logrus.New().WithField("some_key", "some_value")
	ctx = SetRequestLogger(ctx, newLogger)
	assert.Equal(t, newLogger, GetRequestLogger(ctx))
}

func TestSetGetReqID(t *testing.T) {
	ctx := context.Background()
	assert.Empty(t, GetReqID(ctx))

	id := "some_id"
	ctx = SetReqID(ctx, id)
	assert.Equal(t, id, GetReqID(ctx))
}
