package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/zercle/template-go-echo/internal/infrastructure/sqlc"
	"github.com/zercle/template-go-echo/internal/user/domain"
)

// UserRepository implements domain.UserRepository using sqlc generated code
type UserRepository struct {
	q sqlc.Querier
}

// New creates a new user repository with sqlc querier
func New(q sqlc.Querier) *UserRepository {
	return &UserRepository{q: q}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	params := sqlc.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		IsActive:     sql.NullBool{Bool: user.IsActive, Valid: true},
	}

	err := r.q.CreateUser(ctx, params)
	if err != nil {
		slog.Error("failed to create user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	sqlcUser, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get user by id", slog.String("error", err.Error()))
		return nil, err
	}

	return sqlcUserToDomain(&sqlcUser), nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sqlcUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	return sqlcUserToDomain(&sqlcUser), nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	params := sqlc.UpdateUserParams{
		Email: user.Email,
		Name:  user.Name,
		ID:    user.ID,
	}

	err := r.q.UpdateUser(ctx, params)
	if err != nil {
		slog.Error("failed to update user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		slog.Error("failed to delete user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	sqlcUsers, err := r.q.ListUsers(ctx, params)
	if err != nil {
		slog.Error("failed to list users", slog.String("error", err.Error()))
		return nil, err
	}

	users := make([]*domain.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		users[i] = sqlcUserToDomain(&sqlcUser)
	}

	return users, nil
}

// GetUserCount returns the total count of non-deleted users
func (r *UserRepository) GetUserCount(ctx context.Context) (int, error) {
	count, err := r.q.GetUserCount(ctx)
	if err != nil {
		slog.Error("failed to get user count", slog.String("error", err.Error()))
		return 0, err
	}

	return int(count), nil
}

// CreateSession creates a new user session
func (r *UserRepository) CreateSession(ctx context.Context, session *domain.UserSession) error {
	params := sqlc.CreateSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		IpAddress:        sql.NullString{String: session.IPAddress, Valid: session.IPAddress != ""},
		UserAgent:        sql.NullString{String: session.UserAgent, Valid: session.UserAgent != ""},
		ExpiresAt:        session.ExpiresAt,
	}

	err := r.q.CreateSession(ctx, params)
	if err != nil {
		slog.Error("failed to create session", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// GetSessionByID retrieves a session by ID
func (r *UserRepository) GetSessionByID(ctx context.Context, id string) (*domain.UserSession, error) {
	sqlcSession, err := r.q.GetSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get session by id", slog.String("error", err.Error()))
		return nil, err
	}

	return sqlcSessionToDomain(&sqlcSession), nil
}

// GetSessionsByUserID retrieves all active sessions for a user
func (r *UserRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*domain.UserSession, error) {
	sqlcSessions, err := r.q.GetSessionByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get sessions by user id", slog.String("error", err.Error()))
		return nil, err
	}

	sessions := make([]*domain.UserSession, len(sqlcSessions))
	for i, sqlcSession := range sqlcSessions {
		sessions[i] = sqlcSessionToDomain(&sqlcSession)
	}

	return sessions, nil
}

// DeleteSession deletes a session
func (r *UserRepository) DeleteSession(ctx context.Context, id string) error {
	err := r.q.DeleteSession(ctx, id)
	if err != nil {
		slog.Error("failed to delete session", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// DeleteExpiredSessions deletes all expired sessions
func (r *UserRepository) DeleteExpiredSessions(ctx context.Context) error {
	err := r.q.DeleteExpiredSessions(ctx)
	if err != nil {
		slog.Error("failed to delete expired sessions", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// GetSessionByTokenHash retrieves a session by token hash
func (r *UserRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*domain.UserSession, error) {
	sqlcSession, err := r.q.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get session by token hash", slog.String("error", err.Error()))
		return nil, err
	}

	return sqlcSessionToDomain(&sqlcSession), nil
}

// Helper functions to convert sqlc types to domain types

func sqlcUserToDomain(sqlcUser *sqlc.Users) *domain.User {
	user := &domain.User{
		ID:           sqlcUser.ID,
		Email:        sqlcUser.Email,
		Name:         sqlcUser.Name,
		PasswordHash: sqlcUser.PasswordHash,
		IsActive:     sqlcUser.IsActive.Bool,
	}

	if sqlcUser.CreatedAt.Valid {
		user.CreatedAt = sqlcUser.CreatedAt.Time
	}

	if sqlcUser.UpdatedAt.Valid {
		user.UpdatedAt = sqlcUser.UpdatedAt.Time
	}

	if sqlcUser.DeletedAt.Valid {
		user.DeletedAt = &sqlcUser.DeletedAt.Time
	}

	return user
}

func sqlcSessionToDomain(sqlcSession *sqlc.UserSessions) *domain.UserSession {
	session := &domain.UserSession{
		ID:               sqlcSession.ID,
		UserID:           sqlcSession.UserID,
		RefreshTokenHash: sqlcSession.RefreshTokenHash,
		ExpiresAt:        sqlcSession.ExpiresAt,
	}

	if sqlcSession.IpAddress.Valid {
		session.IPAddress = sqlcSession.IpAddress.String
	}

	if sqlcSession.UserAgent.Valid {
		session.UserAgent = sqlcSession.UserAgent.String
	}

	if sqlcSession.CreatedAt.Valid {
		session.CreatedAt = sqlcSession.CreatedAt.Time
	}

	return session
}
