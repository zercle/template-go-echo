package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"
	email := "john@example.com"

	user := NewUser(name, email)

	if user.ID == "" {
		t.Error("expected ID to be generated")
	}
	if user.Name != name {
		t.Errorf("expected name %s, got %s", name, user.Name)
	}
	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}
}

func TestUserIsValid(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    &User{ID: "1", Name: "John", Email: "john@example.com"},
			wantErr: false,
		},
		{
			name:    "missing name",
			user:    &User{ID: "1", Name: "", Email: "john@example.com"},
			wantErr: true,
		},
		{
			name:    "missing email",
			user:    &User{ID: "1", Name: "John", Email: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
