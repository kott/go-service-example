package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kott/go-service-example/pkg/errors"
)

// NoRoute sets default 404 errors to a JSON AppError
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, errors.NewAppError(errors.NotFound,
			errors.Descriptions[errors.NotFound], "nil"))
	}
}

// NoMethod sets default 405 errors to a JSON AppError
func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, errors.NewAppError(errors.MethodNotAllowed,
			errors.Descriptions[errors.MethodNotAllowed], "nil"))
	}
}
