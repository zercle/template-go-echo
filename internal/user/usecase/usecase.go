package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zercle/template-go-echo/internal/user/domain"
	"github.com/zercle/template-go-echo/pkg"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase implements domain.UserUsecase
type UserUsecase struct {
	repo   domain.UserRepository
	tokenTTL int
}

// New creates a new user usecase
func New(repo domain.UserRepository, tokenTTL int) *UserUsecase {
	return &UserUsecase{
		repo:     repo,
		tokenTTL: tokenTTL,
	}
}

// RegisterUser creates a new user with validation
func (u *UserUsecase) RegisterUser(ctx context.Context, email, name, password string) (*domain.User, error) {
	// Validate inputs
	validator := pkg.NewValidator()
	if validator.IsEmpty("email", email) || !validator.IsValidEmail("email", email) {
		return nil, domain.ErrInvalidEmail
	}
	if validator.IsEmpty("name", name) || !validator.IsMaxLength("name", name, domain.MaxNameLength) {
		return nil, domain.ErrInvalidName
	}
	if len(password) < domain.MinPasswordLength || len(password) > domain.MaxPasswordLength {
		return nil, domain.ErrInvalidPassword
	}

	// Check if user already exists
	existingUser, err := u.repo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil && !existingUser.IsDeleted() {
		slog.Warn("attempted to register existing email", slog.String("email", email))
		return nil, domain.ErrUserExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, pkg.ErrInternalError
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: string(passwordHash),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		slog.Error("failed to create user", slog.String("error", err.Error()))
		return nil, pkg.ErrInternalError
	}

	slog.Info("user registered successfully", slog.String("user_id", user.ID), slog.String("email", user.Email))
	return user, nil
}

// LoginUser authenticates a user and returns tokens
func (u *UserUsecase) LoginUser(ctx context.Context, email, password string, ipAddress, userAgent string) (*domain.User, string, string, error) {
	// Get user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil || user.IsDeleted() {
		slog.Warn("login failed: user not found", slog.String("email", email))
		return nil, "", "", domain.ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.Warn("login failed: invalid password", slog.String("email", email))
		return nil, "", "", domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		slog.Warn("login failed: user inactive", slog.String("user_id", user.ID))
		return nil, "", "", domain.ErrUnauthorized
	}

	// Generate tokens
	accessToken := u.generateToken(user.ID, user.Email)
	refreshToken := u.generateRefreshToken(user.ID)

	// Create session
	refreshTokenHash := u.hashToken(refreshToken)
	session := &domain.UserSession{
		ID:               uuid.New().String(),
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(time.Hour * domain.SessionDurationHours),
		CreatedAt:        time.Now(),
	}

	if err := u.repo.CreateSession(ctx, session); err != nil {
		slog.Error("failed to create session", slog.String("error", err.Error()))
		return nil, "", "", pkg.ErrInternalError
	}

	slog.Info("user logged in successfully", slog.String("user_id", user.ID))
	return user, accessToken, refreshToken, nil
}

// GetUser retrieves a user by ID
func (u *UserUsecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil || user == nil || user.IsDeleted() {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil || user.IsDeleted() {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

// UpdateUserProfile updates user information
func (u *UserUsecase) UpdateUserProfile(ctx context.Context, id, name, email string) (*domain.User, error) {
	// Validate inputs
	if name == "" || len(name) > domain.MaxNameLength {
		return nil, domain.ErrInvalidName
	}
	validator := pkg.NewValidator()
	if !validator.IsValidEmail("email", email) {
		return nil, domain.ErrInvalidEmail
	}

	// Get existing user
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil || user == nil || user.IsDeleted() {
		return nil, domain.ErrUserNotFound
	}

	// Check if email is already in use by another user
	if email != user.Email {
		existingUser, _ := u.repo.GetUserByEmail(ctx, email)
		if existingUser != nil && !existingUser.IsDeleted() && existingUser.ID != id {
			return nil, domain.ErrUserExists
		}
	}

	// Update user
	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now()

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		slog.Error("failed to update user", slog.String("error", err.Error()))
		return nil, pkg.ErrInternalError
	}

	slog.Info("user profile updated", slog.String("user_id", id))
	return user, nil
}

// ChangePassword changes user password
func (u *UserUsecase) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	// Get user
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil || user == nil || user.IsDeleted() {
		return domain.ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		slog.Warn("password change failed: invalid old password", slog.String("user_id", id))
		return pkg.NewDomainError(domain.ErrCodeInvalidPassword, domain.ValidationMessages["old_password_invalid"])
	}

	// Validate new password
	if len(newPassword) < domain.MinPasswordLength || len(newPassword) > domain.MaxPasswordLength {
		return domain.ErrInvalidPassword
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash new password", slog.String("error", err.Error()))
		return pkg.ErrInternalError
	}

	// Update password
	user.PasswordHash = string(passwordHash)
	user.UpdatedAt = time.Now()

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		slog.Error("failed to update password", slog.String("error", err.Error()))
		return pkg.ErrInternalError
	}

	slog.Info("password changed successfully", slog.String("user_id", id))
	return nil
}

// DeleteUser deletes a user account
func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	// Get user
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil || user == nil || user.IsDeleted() {
		return domain.ErrUserNotFound
	}

	// Soft delete user
	if err := u.repo.DeleteUser(ctx, id); err != nil {
		slog.Error("failed to delete user", slog.String("error", err.Error()))
		return pkg.ErrInternalError
	}

	// Logout all sessions
	_ = u.LogoutAllSessions(ctx, id)

	slog.Info("user deleted", slog.String("user_id", id))
	return nil
}

// ListUsers retrieves a paginated list of users
func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := u.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list users", slog.String("error", err.Error()))
		return nil, 0, pkg.ErrInternalError
	}

	count, err := u.repo.GetUserCount(ctx)
	if err != nil {
		slog.Error("failed to get user count", slog.String("error", err.Error()))
		return nil, 0, pkg.ErrInternalError
	}

	return users, count, nil
}

// RefreshToken generates a new access token from refresh token
func (u *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Hash the refresh token to find the session
	tokenHash := u.hashToken(refreshToken)

	// Find session by token hash
	session, err := u.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil || session == nil {
		return "", domain.ErrSessionNotFound
	}

	// Check if session is expired
	if session.IsExpired() {
		_ = u.repo.DeleteSession(ctx, session.ID)
		return "", domain.ErrSessionExpired
	}

	// Get user
	user, err := u.repo.GetUserByID(ctx, session.UserID)
	if err != nil || user == nil || user.IsDeleted() {
		return "", domain.ErrUserNotFound
	}

	// Generate new access token
	accessToken := u.generateToken(user.ID, user.Email)

	slog.Info("token refreshed", slog.String("user_id", user.ID))
	return accessToken, nil
}

// LogoutUser invalidates a session
func (u *UserUsecase) LogoutUser(ctx context.Context, sessionID string) error {
	// Get session to verify it exists
	session, err := u.repo.GetSessionByID(ctx, sessionID)
	if err != nil || session == nil {
		return domain.ErrSessionNotFound
	}

	// Delete session
	if err := u.repo.DeleteSession(ctx, sessionID); err != nil {
		slog.Error("failed to delete session", slog.String("error", err.Error()))
		return pkg.ErrInternalError
	}

	slog.Info("user logged out", slog.String("user_id", session.UserID))
	return nil
}

// LogoutAllSessions invalidates all sessions for a user
func (u *UserUsecase) LogoutAllSessions(ctx context.Context, userID string) error {
	// Get all active sessions for user
	sessions, err := u.repo.GetSessionsByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user sessions", slog.String("error", err.Error()))
		return pkg.ErrInternalError
	}

	// Delete all sessions
	for _, session := range sessions {
		_ = u.repo.DeleteSession(ctx, session.ID)
	}

	slog.Info("all sessions deleted for user", slog.String("user_id", userID), slog.Int("count", len(sessions)))
	return nil
}

// generateToken creates a simple JWT token (in production, use proper JWT library)
func (u *UserUsecase) generateToken(userID, email string) string {
	// This is a placeholder - in production, use golang-jwt/jwt
	// For now, return a simple token format
	return fmt.Sprintf("%s.%s.%d", userID, email, time.Now().Unix())
}

// generateRefreshToken creates a refresh token
func (u *UserUsecase) generateRefreshToken(userID string) string {
	// Generate a unique refresh token
	return uuid.New().String() + "-" + userID
}

// hashToken hashes a token for storage
func (u *UserUsecase) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
