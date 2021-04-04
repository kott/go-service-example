package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistContext(t *testing.T) {
	s := gin.New()
	s.Use(PersistContext())
	s.GET("/", func(c *gin.Context) {
		c.JSON(200, struct{}{})
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	r.Header.Set("X-Request-Id", "some-req-id")
	r.Header.Set("Accept-Language", "some-locale")
	s.ServeHTTP(rr, r)

	assert.Equal(t, "some-req-id", rr.Header()["X-Request-Id"][0])
}

func TestPersistContextNoReqID(t *testing.T) {
	s := gin.New()
	s.Use(PersistContext())
	s.GET("/", func(c *gin.Context) {
		c.JSON(200, struct{}{})
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	s.ServeHTTP(rr, r)

	hReqID := rr.Header()["X-Request-Id"][0]
	assert.NotEmpty(t, hReqID)
}
