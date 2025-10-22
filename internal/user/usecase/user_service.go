package usecase

import (
	"context"

	"github.com/zercle/template-go-echo/internal/user/domain"
	"github.com/zercle/template-go-echo/pkg/logger"
)

// UserService handles user business logic.
type UserService struct {
	repo domain.UserRepository
	log  *logger.Logger
}

// NewUserService creates a new user service.
func NewUserService(repo domain.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	user := domain.NewUser(req.Name, req.Email)

	if err := user.IsValid(); err != nil {
		s.log.Warn("invalid user data", "error", err)
		return nil, err
	}

	// Check if user with this email already exists
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		s.log.Warn("user already exists", "email", user.Email)
		return nil, domain.ErrUserAlreadyExists
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error("failed to create user", "error", err)
		return nil, err
	}

	s.log.Info("user created", "id", created.ID, "email", created.Email)

	return toUserResponse(created), nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("user not found", "id", id, "error", err)
		return nil, domain.ErrUserNotFound
	}

	return toUserResponse(user), nil
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("user not found", "id", id)
		return nil, domain.ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if err := user.IsValid(); err != nil {
		s.log.Warn("invalid user data", "id", id, "error", err)
		return nil, err
	}

	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		s.log.Error("failed to update user", "id", id, "error", err)
		return nil, err
	}

	s.log.Info("user updated", "id", updated.ID)

	return toUserResponse(updated), nil
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error("failed to delete user", "id", id, "error", err)
		return err
	}

	s.log.Info("user deleted", "id", id)
	return nil
}

// ListUsers retrieves all users.
func (s *UserService) ListUsers(ctx context.Context, offset, limit int) (*ListUsersResponse, error) {
	users, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		s.log.Error("failed to list users", "error", err)
		return nil, err
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = toUserResponse(user)
	}

	return &ListUsersResponse{
		Users:  userResponses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

// Helper function to convert domain.User to UserResponse.
func toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
