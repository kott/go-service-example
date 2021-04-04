package errors

import "fmt"

const (
	// InternalServerError ...
	InternalServerError = "INTERNAL_SERVER_ERROR"

	// BadGateway ...
	BadGateway = "BAD_GATEWAY"

	// ServiceUnavailable ...
	ServiceUnavailable = "SERVICE_UNAVAILABLE"

	// BadRequest ...
	BadRequest = "BAD_REQUEST"

	// NotFound ...
	NotFound = "NOT_FOUND"

	// MethodNotAllowed ...
	MethodNotAllowed = "METHOD_NOT_ALLOWED"

	// MissingPermissions ...
	MissingPermissions = "MISSING_PERMISSIONS"

	// InvalidAuthorizationHeader ...
	InvalidAuthorizationHeader = "INVALID_AUTHORIZATION_HEADER"

	// MissingAuthorizationHeader ...
	MissingAuthorizationHeader = "MISSING_AUTHORIZATION_HEADER"

	// ExpiredAuthorizationHeader ...
	ExpiredAuthorizationHeader = "EXPIRED_AUTHORIZATION_HEADER"

	// ForbiddenAction ...
	ForbiddenAction = "FORBIDDEN_ACTION"

	// PreconditionFailed ...
	PreconditionFailed = "PRECONDITION_FAILED"

	// UnsupportedMediaType ...
	UnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
)

// ErrorCode is the string representation of an HTTP error
type ErrorCode string

// Descriptions maps error codes to a message associated with it.
var Descriptions = map[ErrorCode]string{
	InternalServerError:        "Internal server error.",
	BadGateway:                 "The server encountered a temporary error and could not complete your request.",
	ServiceUnavailable:         "The requested service is currently unreachable.",
	BadRequest:                 "The application sent a request that this server could not understand.",
	NotFound:                   "Resource does not exist.",
	MethodNotAllowed:           "Method is not allowed for this resource.",
	MissingPermissions:         "You do not have permissions to access this endpoint.",
	InvalidAuthorizationHeader: "The request contains an invalid Authorization header.",
	MissingAuthorizationHeader: "Request requires a JWT Authorization header.",
	ExpiredAuthorizationHeader: "The request contains an expired Authorization header.",
	ForbiddenAction:            "The action being performed is forbidden",
	PreconditionFailed:         "Precondition failed.",
	UnsupportedMediaType:       "The server does not support the media type transmitted in the request.",
}

// AppError application specific error
type AppError struct {
	Code        ErrorCode `json:"code"`
	Description string    `json:"description"`
	Field       string    `json:"field"`
}

// AppErrors is a collection of AppError
type AppErrors struct {
	Errors []AppError `json:"errors"`
}

// Error is a friendly formatted message
func (e *AppError) Error() string {
	return fmt.Sprintf("%s (%s) %s", e.Code, e.Field, e.Description)
}

// NewAppError makes it easy to construct new errors with everything populated
func NewAppError(code ErrorCode, description, field string) error {
	appErr := &AppError{
		Code:        code,
		Description: description,
		Field:       field,
	}
	return appErr
}
