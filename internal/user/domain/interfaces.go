package domain

import "context"

//go:generate go run github.com/uber-go/mock/cmd/mockgen -destination=../mock/mock_repository.go -package=mock github.com/zercle/template-go-echo/internal/user/domain UserRepository
//go:generate go run github.com/uber-go/mock/cmd/mockgen -destination=../mock/mock_usecase.go -package=mock github.com/zercle/template-go-echo/internal/user/domain UserUsecase

// UserRepository defines database operations for users
type UserRepository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *User) error

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser soft deletes a user
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)

	// GetUserCount returns the total count of non-deleted users
	GetUserCount(ctx context.Context) (int, error)

	// CreateSession creates a new user session
	CreateSession(ctx context.Context, session *UserSession) error

	// GetSessionByID retrieves a session by ID
	GetSessionByID(ctx context.Context, id string) (*UserSession, error)

	// GetSessionsByUserID retrieves all active sessions for a user
	GetSessionsByUserID(ctx context.Context, userID string) ([]*UserSession, error)

	// DeleteSession deletes a session
	DeleteSession(ctx context.Context, id string) error

	// DeleteExpiredSessions deletes all expired sessions
	DeleteExpiredSessions(ctx context.Context) error

	// GetSessionByTokenHash retrieves a session by token hash
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*UserSession, error)
}

// UserUsecase defines business logic for users
type UserUsecase interface {
	// RegisterUser creates a new user with validation
	RegisterUser(ctx context.Context, email, name, password string) (*User, error)

	// LoginUser authenticates a user and returns tokens
	LoginUser(ctx context.Context, email, password string, ipAddress, userAgent string) (*User, string, string, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id string) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// UpdateUserProfile updates user information
	UpdateUserProfile(ctx context.Context, id, name, email string) (*User, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error

	// DeleteUser deletes a user account
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error)

	// RefreshToken generates a new access token from refresh token
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// LogoutUser invalidates a session
	LogoutUser(ctx context.Context, sessionID string) error

	// LogoutAllSessions invalidates all sessions for a user
	LogoutAllSessions(ctx context.Context, userID string) error
}
