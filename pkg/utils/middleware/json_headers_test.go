package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONHeaders(t *testing.T) {
	s := gin.New()
	s.Use(JSONResponseHeader())
	s.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestForceJSONReject(t *testing.T) {
	s := gin.New()
	s.Use(ForceJSON())
	s.POST("/", func(c *gin.Context) {
		c.Status(201)
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/", strings.NewReader("{}"))
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	expectedResponse := `{"code":"UNSUPPORTED_MEDIA_TYPE","description":"The server does not support the media type transmitted in the request.","field":"nil"}`
	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}

func TestForceJSONPass(t *testing.T) {
	s := gin.New()
	s.Use(ForceJSON())
	s.POST("/", func(c *gin.Context) {
		c.JSON(201, struct{}{})
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/", strings.NewReader("{}"))
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	expectedResponse := `{}`
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}

func TestForceJSONGet(t *testing.T) {
	s := gin.New()
	s.Use(ForceJSON())
	s.GET("/", func(c *gin.Context) {
		c.JSON(200, struct{}{})
	})

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	s.ServeHTTP(rr, r)
	expectedResponse := `{}`
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedResponse, string(rr.Body.Bytes()))
}
