package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoMethod(t *testing.T) {
	s := gin.New()
	s.NoMethod(NoMethod())
	s.HandleMethodNotAllowed = true
	s.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	expectedResponse := `{"code":"METHOD_NOT_ALLOWED","description":"Method is not allowed for this resource.","field":"nil"}`
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}

func TestNoRoute(t *testing.T) {
	s := gin.New()
	s.NoRoute(NoRoute())

	s.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/does-not-exist", nil)
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	expectedResponse := `{"code":"NOT_FOUND","description":"Resource does not exist.","field":"nil"}`
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}
