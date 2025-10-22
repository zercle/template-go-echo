package repository

import (
	"context"
	"testing"

	"github.com/zercle/template-go-echo/internal/user/domain"
)

func TestMemoryRepositoryCreate(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	created, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, created.ID)
	}
}

func TestMemoryRepositoryCreateDuplicate(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	_, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("First Create failed: %v", err)
	}

	_, err = repo.Create(context.Background(), user)
	if err == nil {
		t.Error("expected error for duplicate user")
	}
}

func TestMemoryRepositoryGetByID(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	if _, err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := repo.GetByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrieved.Name != user.Name {
		t.Errorf("expected name %s, got %s", user.Name, retrieved.Name)
	}
}

func TestMemoryRepositoryGetByIDNotFound(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.GetByID(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestMemoryRepositoryGetByEmail(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	if _, err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := repo.GetByEmail(context.Background(), "john@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, retrieved.ID)
	}
}

func TestMemoryRepositoryUpdate(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	if _, err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	user.Name = "Jane"
	updated, err := repo.Update(context.Background(), user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != "Jane" {
		t.Errorf("expected name Jane, got %s", updated.Name)
	}
}

func TestMemoryRepositoryDelete(t *testing.T) {
	repo := NewMemoryRepository()
	user := &domain.User{ID: "1", Name: "John", Email: "john@example.com"}

	if _, err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err := repo.Delete(context.Background(), "1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(context.Background(), "1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestMemoryRepositoryList(t *testing.T) {
	repo := NewMemoryRepository()

	for i := 1; i <= 15; i++ {
		user := &domain.User{ID: string(rune(i)), Name: "User" + string(rune(i)), Email: "user" + string(rune(i)) + "@example.com"}
		if _, err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	users, total, err := repo.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) != 10 {
		t.Errorf("expected 10 users, got %d", len(users))
	}

	if total != 15 {
		t.Errorf("expected total 15, got %d", total)
	}
}
