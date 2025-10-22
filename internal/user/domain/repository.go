package domain

import "context"

//go:generate go run github.com/golang/mock/cmd/mockgen -destination=../mock/mock_repository.go -package=mock . UserRepository

// UserRepository defines the interface for user data access.
type UserRepository interface {
	// Create saves a new user and returns the created user.
	Create(ctx context.Context, user *User) (*User, error)

	// GetByID retrieves a user by ID.
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByEmail retrieves a user by email.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user.
	Update(ctx context.Context, user *User) (*User, error)

	// Delete deletes a user by ID.
	Delete(ctx context.Context, id string) error

	// List retrieves all users with pagination.
	List(ctx context.Context, offset, limit int) ([]*User, int, error)
}
