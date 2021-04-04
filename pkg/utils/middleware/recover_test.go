package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	s := gin.New()
	s.Use(Recover())
	s.GET("/", func(c *gin.Context) {
		panic(errors.New("some error"))
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	s.ServeHTTP(rr, r)

	expectedResponse := `{"code":"INTERNAL_SERVER_ERROR","description":"Internal server error.","field":""}`
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}

func TestRecoverWithPanicMsg(t *testing.T) {
	s := gin.New()
	s.Use(Recover())
	s.GET("/", func(c *gin.Context) {
		panic("something bad happened")
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)
	s.ServeHTTP(rr, r)

	expectedResponse := `{"code":"INTERNAL_SERVER_ERROR","description":"Internal server error.","field":""}`
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}
