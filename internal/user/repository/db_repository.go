package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zercle/template-go-echo/internal/user/domain"
	"github.com/zercle/template-go-echo/internal/user/repository/sqlc"
)

// DatabaseRepository implements the UserRepository interface using sqlc-generated code.
type DatabaseRepository struct {
	queries sqlc.Querier
}

// NewDatabaseRepository creates a new database repository.
func NewDatabaseRepository(db *sql.DB) *DatabaseRepository {
	return &DatabaseRepository{
		queries: sqlc.New(db),
	}
}

// Create saves a new user.
func (r *DatabaseRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID.
func (r *DatabaseRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// GetByEmail retrieves a user by email.
func (r *DatabaseRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Update updates an existing user.
func (r *DatabaseRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user by ID.
func (r *DatabaseRepository) Delete(ctx context.Context, id string) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves all users with pagination.
func (r *DatabaseRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int, error) {
	users, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert sqlc users to domain users
	domainUsers := make([]*domain.User, len(users))
	for i, u := range users {
		domainUsers[i] = &domain.User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	return domainUsers, int(count), nil
}
