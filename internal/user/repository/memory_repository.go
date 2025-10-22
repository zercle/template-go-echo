package repository

import (
	"context"
	"sync"

	"github.com/zercle/template-go-echo/internal/user/domain"
)

// MemoryRepository is an in-memory implementation of the UserRepository interface.
// Used for development and testing before database integration.
type MemoryRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

// NewMemoryRepository creates a new memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*domain.User),
	}
}

// Create saves a new user.
func (r *MemoryRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return nil, domain.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	return user, nil
}

// GetByID retrieves a user by ID.
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *MemoryRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// Update updates an existing user.
func (r *MemoryRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return nil, domain.ErrUserNotFound
	}

	r.users[user.ID] = user
	return user, nil
}

// Delete deletes a user by ID.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// List retrieves all users with pagination.
func (r *MemoryRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.users)

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Simple pagination
	if offset > len(users) {
		return []*domain.User{}, total, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], total, nil
}
