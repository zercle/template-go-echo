package usecase

import (
	"context"
	"testing"

	"github.com/zercle/template-go-echo/internal/user/repository"
	"github.com/zercle/template-go-echo/pkg/logger"
)

func newTestUserService() *UserService {
	repo := repository.NewMemoryRepository()
	log := logger.New("info")
	return NewUserService(repo, log)
}

func TestCreateUser(t *testing.T) {
	service := newTestUserService()

	req := CreateUserRequest{Name: "John Doe", Email: "john@example.com"}
	user, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, user.Name)
	}
	if user.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, user.Email)
	}
	if user.ID == "" {
		t.Error("expected ID to be generated")
	}
}

func TestCreateUserDuplicate(t *testing.T) {
	service := newTestUserService()

	req := CreateUserRequest{Name: "John Doe", Email: "john@example.com"}
	_, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("First CreateUser failed: %v", err)
	}

	_, err = service.CreateUser(context.Background(), req)
	if err == nil {
		t.Error("expected error for duplicate user")
	}
}

func TestGetUser(t *testing.T) {
	service := newTestUserService()

	req := CreateUserRequest{Name: "John Doe", Email: "john@example.com"}
	created, _ := service.CreateUser(context.Background(), req)

	user, err := service.GetUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, user.Name)
	}
}

func TestGetUserNotFound(t *testing.T) {
	service := newTestUserService()

	_, err := service.GetUser(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestUpdateUser(t *testing.T) {
	service := newTestUserService()

	req := CreateUserRequest{Name: "John Doe", Email: "john@example.com"}
	created, _ := service.CreateUser(context.Background(), req)

	updateReq := UpdateUserRequest{Name: "Jane Doe"}
	updated, err := service.UpdateUser(context.Background(), created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updated.Name != "Jane Doe" {
		t.Errorf("expected name Jane Doe, got %s", updated.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	service := newTestUserService()

	req := CreateUserRequest{Name: "John Doe", Email: "john@example.com"}
	created, _ := service.CreateUser(context.Background(), req)

	err := service.DeleteUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	_, err = service.GetUser(context.Background(), created.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestListUsers(t *testing.T) {
	service := newTestUserService()

	for i := 0; i < 5; i++ {
		_, err := service.CreateUser(context.Background(), CreateUserRequest{
			Name:  "User " + string(rune(i)),
			Email: "user" + string(rune(i)) + "@example.com",
		})
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
	}

	result, err := service.ListUsers(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(result.Users) != 5 {
		t.Errorf("expected 5 users, got %d", len(result.Users))
	}

	if result.Total != 5 {
		t.Errorf("expected total 5, got %d", result.Total)
	}
}
