package integration_test

import (
	"context"
	"testing"

	"github.com/zercle/template-go-echo/internal/user/domain"
	"github.com/zercle/template-go-echo/internal/user/usecase"
	"github.com/zercle/template-go-echo/internal/user/test/mocks"
)

func TestRegisterUserSuccess(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	user, err := uc.RegisterUser(context.Background(), "test@example.com", "Test User", "SecurePass123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if user.Name != "Test User" {
		t.Errorf("expected name Test User, got %s", user.Name)
	}

	if !user.IsActive {
		t.Error("expected user to be active")
	}
}

func TestRegisterUserInvalidEmail(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	_, err := uc.RegisterUser(context.Background(), "invalid-email", "Test User", "SecurePass123")
	if err != domain.ErrInvalidEmail {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestRegisterUserInvalidPassword(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	_, err := uc.RegisterUser(context.Background(), "test@example.com", "Test User", "short")
	if err != domain.ErrInvalidPassword {
		t.Errorf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestGetUserSuccess(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	// Create user first
	user, _ := uc.RegisterUser(context.Background(), "test@example.com", "Test User", "SecurePass123")

	// Get user
	retrieved, err := uc.GetUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected user to be retrieved")
	}

	if retrieved.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestGetUserNotFound(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	_, err := uc.GetUser(context.Background(), "non-existent-id")
	if err != domain.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestListUsers(t *testing.T) {
	repo := mocks.NewMockRepository()
	uc := usecase.New(repo, 3600)

	// Create multiple users
	for i := 1; i <= 3; i++ {
		email := "user" + string(rune('0'+i)) + "@example.com"
		_, err := uc.RegisterUser(context.Background(), email, "User "+string(rune('0'+i)), "SecurePass123")
		if err != nil {
			t.Errorf("failed to register user: %v", err)
		}
	}

	// List users
	users, total, err := uc.ListUsers(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if total != 3 {
		t.Errorf("expected 3 users, got %d", total)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users in list, got %d", len(users))
	}
}
