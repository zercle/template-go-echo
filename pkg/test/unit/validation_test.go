package unit_test

import (
	"testing"

	"github.com/zercle/template-go-echo/pkg"
)

func TestValidatorIsEmpty(t *testing.T) {
	v := pkg.NewValidator()

	// Test empty string
	if !v.IsEmpty("email", "") {
		t.Error("expected IsEmpty to return true for empty string")
	}

	// Test non-empty string
	v2 := pkg.NewValidator()
	if v2.IsEmpty("email", "test@example.com") {
		t.Error("expected IsEmpty to return false for non-empty string")
	}
}

func TestValidatorIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name+tag@example.co.uk", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"test@.com", false},
	}

	for _, tt := range tests {
		v := pkg.NewValidator()
		result := v.IsValidEmail("email", tt.email)
		if result != tt.valid {
			t.Errorf("email %s: expected %v, got %v", tt.email, tt.valid, result)
		}
	}
}

func TestValidatorIsMinLength(t *testing.T) {
	v := pkg.NewValidator()
	if v.IsMinLength("password", "short", 10) {
		t.Error("expected IsMinLength to return false for string shorter than minimum")
	}

	v2 := pkg.NewValidator()
	if !v2.IsMinLength("password", "longpassword", 10) {
		t.Error("expected IsMinLength to return true for string meeting minimum")
	}
}

func TestValidatorIsMaxLength(t *testing.T) {
	v := pkg.NewValidator()
	if !v.IsMaxLength("name", "a", 255) {
		t.Error("expected IsMaxLength to return true for string within limit")
	}

	v2 := pkg.NewValidator()
	longString := string(make([]byte, 300))
	if v2.IsMaxLength("name", longString, 255) {
		t.Error("expected IsMaxLength to return false for string exceeding limit")
	}
}

func TestValidatorIsValid(t *testing.T) {
	v := pkg.NewValidator()
	if !v.IsValid() {
		t.Error("expected new validator to be valid")
	}

	_ = v.AddError("field", "error message")
	if v.IsValid() {
		t.Error("expected validator with errors to be invalid")
	}
}

func TestValidatorErrors(t *testing.T) {
	v := pkg.NewValidator()
	_ = v.AddError("email", "email is required")
	_ = v.AddError("email", "email is invalid")
	_ = v.AddError("password", "password is required")

	errors := v.Errors()
	if len(errors) != 2 {
		t.Errorf("expected 2 fields with errors, got %d", len(errors))
	}

	if len(errors["email"]) != 2 {
		t.Errorf("expected 2 errors for email, got %d", len(errors["email"]))
	}
}
