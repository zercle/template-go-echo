package domain

import (
	"github.com/zercle/template-go-echo/pkg/errors"
)

// Domain-specific error variables.
var (
	ErrUserNotFound      = errors.New("USER_NOT_FOUND", "user not found")
	ErrUserAlreadyExists = errors.New("USER_ALREADY_EXISTS", "user already exists")
	ErrInvalidUserName   = errors.New("INVALID_USER_NAME", "invalid user name")
	ErrInvalidUserEmail  = errors.New("INVALID_USER_EMAIL", "invalid user email")
)
