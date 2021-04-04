package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kott/go-service-example/pkg/errors"
)

func TestError(t *testing.T) {
	appError := errors.AppError{
		Code:        "some_code",
		Description: "some_description",
		Field:       "some_field",
	}
	assert.Equal(t, "some_code (some_field) some_description", appError.Error())
}

func TestNewAppError(t *testing.T) {
	code := errors.ErrorCode("some_code")
	description := "some_description"
	field := "some_field"
	appError := errors.NewAppError(code, description, field)
	assert.Equal(t, "some_code (some_field) some_description", appError.Error())
}
