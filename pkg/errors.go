package pkg

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(details map[string]interface{}) *DomainError {
	e.Details = details
	return e
}

// Common error codes
const (
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeInternalError  = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
)

// Common errors
var (
	ErrNotFound       = NewDomainError(ErrCodeNotFound, "Resource not found")
	ErrUnauthorized   = NewDomainError(ErrCodeUnauthorized, "Unauthorized")
	ErrForbidden      = NewDomainError(ErrCodeForbidden, "Forbidden")
	ErrInternalError  = NewDomainError(ErrCodeInternalError, "Internal server error")
	ErrBadRequest     = NewDomainError(ErrCodeBadRequest, "Bad request")
)
