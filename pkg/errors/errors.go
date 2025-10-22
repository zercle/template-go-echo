package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// DomainError represents a domain-level error.
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// New creates a new domain error.
func New(code, message string) error {
	return &DomainError{Code: code, Message: message}
}

// Newf creates a new domain error with formatted message.
func Newf(code, format string, args ...any) error {
	return &DomainError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Wrap creates a new domain error wrapping an underlying error.
func Wrap(code string, err error) error {
	return &DomainError{Code: code, Message: "", Err: err}
}

// Wrapf creates a new domain error wrapping an underlying error with formatted message.
func Wrapf(code string, err error, format string, args ...any) error {
	return &DomainError{Code: code, Message: fmt.Sprintf(format, args...), Err: err}
}

// Domain-specific error codes.
var (
	// NotFound is returned when a resource is not found.
	NotFound = New("NOT_FOUND", "resource not found")

	// AlreadyExists is returned when a resource already exists.
	AlreadyExists = New("ALREADY_EXISTS", "resource already exists")

	// InvalidInput is returned when input validation fails.
	InvalidInput = New("INVALID_INPUT", "invalid input")

	// Unauthorized is returned when authentication fails.
	Unauthorized = New("UNAUTHORIZED", "unauthorized")

	// Forbidden is returned when authorization fails.
	Forbidden = New("FORBIDDEN", "forbidden")

	// InternalError is returned for internal server errors.
	InternalError = New("INTERNAL_ERROR", "internal server error")
)

// Is checks if an error is a specific domain error by code.
func Is(err error, code string) bool {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Code == code
	}
	return false
}

// GetCode retrieves the error code from a domain error.
func GetCode(err error) string {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Code
	}
	return "INTERNAL_ERROR"
}

// GetMessage retrieves the error message from a domain error.
func GetMessage(err error) string {
	var de *DomainError
	if errors.As(err, &de) {
		if de.Message != "" {
			return de.Message
		}
		if de.Err != nil {
			return de.Err.Error()
		}
	}
	return "An error occurred"
}

// HTTPStatusCode returns the HTTP status code for a domain error.
func HTTPStatusCode(err error) int {
	code := GetCode(err)
	switch code {
	case "NOT_FOUND":
		return http.StatusNotFound
	case "ALREADY_EXISTS":
		return http.StatusConflict
	case "INVALID_INPUT":
		return http.StatusBadRequest
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
