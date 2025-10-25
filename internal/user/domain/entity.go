package domain

import "time"

// User represents a user entity in the domain
type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Name         string    `db:"name" json:"name"`
	PasswordHash string    `db:"password_hash" json:"-"` // Never expose password hash
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// IsDeleted checks if user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// UserSession represents a user session with refresh token
type UserSession struct {
	ID               string    `db:"id" json:"id"`
	UserID           string    `db:"user_id" json:"user_id"`
	RefreshTokenHash string    `db:"refresh_token_hash" json:"-"` // Never expose token hash
	IPAddress        string    `db:"ip_address" json:"ip_address"`
	UserAgent        string    `db:"user_agent" json:"user_agent"`
	ExpiresAt        time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// IsExpired checks if the session has expired
func (us *UserSession) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}
