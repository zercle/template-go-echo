package domain

import "github.com/zercle/template-go-echo/pkg"

// User domain-specific error codes
const (
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeUserExists         = "USER_ALREADY_EXISTS"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeInvalidEmail       = "INVALID_EMAIL"
	ErrCodeInvalidPassword    = "INVALID_PASSWORD"
	ErrCodeInvalidName        = "INVALID_NAME"
	ErrCodeSessionNotFound    = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired     = "SESSION_EXPIRED"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
)

// User domain errors
var (
	ErrUserNotFound = pkg.NewDomainError(
		ErrCodeUserNotFound,
		"user not found",
	)

	ErrUserExists = pkg.NewDomainError(
		ErrCodeUserExists,
		"user with this email already exists",
	)

	ErrInvalidCredentials = pkg.NewDomainError(
		ErrCodeInvalidCredentials,
		"invalid email or password",
	)

	ErrInvalidEmail = pkg.NewDomainError(
		ErrCodeInvalidEmail,
		"email format is invalid",
	)

	ErrInvalidPassword = pkg.NewDomainError(
		ErrCodeInvalidPassword,
		"password does not meet requirements",
	)

	ErrInvalidName = pkg.NewDomainError(
		ErrCodeInvalidName,
		"name is required and must be between 1 and 255 characters",
	)

	ErrSessionNotFound = pkg.NewDomainError(
		ErrCodeSessionNotFound,
		"session not found",
	)

	ErrSessionExpired = pkg.NewDomainError(
		ErrCodeSessionExpired,
		"session has expired",
	)

	ErrUnauthorized = pkg.NewDomainError(
		ErrCodeUnauthorized,
		"unauthorized access",
	)
)
