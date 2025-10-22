package domain

import (
	"github.com/google/uuid"
)

// User represents a user entity.
type User struct {
	ID    string
	Name  string
	Email string
}

// NewUser creates a new user with a generated ID.
func NewUser(name, email string) *User {
	return &User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
}

// IsValid checks if the user is valid.
func (u *User) IsValid() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}

	if u.Email == "" {
		return ErrInvalidUserEmail
	}

	return nil
}
