package mocks

import (
	"context"

	"github.com/zercle/template-go-echo/internal/user/domain"
)

// MockUserRepository is a simple mock for testing
type MockUserRepository struct {
	users    map[string]*domain.User
	sessions map[string]*domain.UserSession
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockUserRepository {
	return &MockUserRepository{
		users:    make(map[string]*domain.User),
		sessions: make(map[string]*domain.UserSession),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return m.users[id], nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email && !user.IsDeleted() {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	user := m.users[id]
	if user != nil {
		user.DeletedAt = nil // Mock soft delete
	}
	return nil
}

func (m *MockUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.users {
		if !user.IsDeleted() {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int, error) {
	count := 0
	for _, user := range m.users {
		if !user.IsDeleted() {
			count++
		}
	}
	return count, nil
}

func (m *MockUserRepository) CreateSession(ctx context.Context, session *domain.UserSession) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MockUserRepository) GetSessionByID(ctx context.Context, id string) (*domain.UserSession, error) {
	return m.sessions[id], nil
}

func (m *MockUserRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*domain.UserSession, error) {
	var sessions []*domain.UserSession
	for _, session := range m.sessions {
		if session.UserID == userID && !session.IsExpired() {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *MockUserRepository) DeleteSession(ctx context.Context, id string) error {
	delete(m.sessions, id)
	return nil
}

func (m *MockUserRepository) DeleteExpiredSessions(ctx context.Context) error {
	return nil
}

func (m *MockUserRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*domain.UserSession, error) {
	for _, session := range m.sessions {
		if session.RefreshTokenHash == tokenHash && !session.IsExpired() {
			return session, nil
		}
	}
	return nil, nil
}
