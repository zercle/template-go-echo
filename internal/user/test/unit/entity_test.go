package unit_test

import (
	"testing"
	"time"

	"github.com/zercle/template-go-echo/internal/user/domain"
)

func TestUserIsDeleted(t *testing.T) {
	// User without deletion timestamp
	user := &domain.User{
		ID:        "user-1",
		Email:     "user@example.com",
		DeletedAt: nil,
	}
	if user.IsDeleted() {
		t.Error("expected user to not be deleted")
	}

	// User with deletion timestamp
	now := time.Now()
	user.DeletedAt = &now
	if !user.IsDeleted() {
		t.Error("expected user to be deleted")
	}
}

func TestUserSessionIsExpired(t *testing.T) {
	// Future expiration
	session := &domain.UserSession{
		ID:        "session-1",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if session.IsExpired() {
		t.Error("expected session to not be expired")
	}

	// Past expiration
	session.ExpiresAt = time.Now().Add(-time.Hour)
	if !session.IsExpired() {
		t.Error("expected session to be expired")
	}
}
