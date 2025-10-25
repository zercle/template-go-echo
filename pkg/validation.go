package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator provides common validation utilities
type Validator struct {
	errors map[string][]string
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string][]string),
	}
}

// AddError adds a validation error for a field
func (v *Validator) AddError(field, message string) *Validator {
	v.errors[field] = append(v.errors[field], message)
	return v
}

// IsEmpty checks if a string is empty
func (v *Validator) IsEmpty(field, value string) bool {
	if strings.TrimSpace(value) == "" {
		_ = v.AddError(field, fmt.Sprintf("%s is required", field))
		return true
	}
	return false
}

// IsValidEmail checks if an email is valid
func (v *Validator) IsValidEmail(field, email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		_ = v.AddError(field, fmt.Sprintf("%s is not a valid email", field))
		return false
	}
	return true
}

// IsMinLength checks if a string has minimum length
func (v *Validator) IsMinLength(field, value string, minLength int) bool {
	if len(strings.TrimSpace(value)) < minLength {
		_ = v.AddError(field, fmt.Sprintf("%s must be at least %d characters", field, minLength))
		return false
	}
	return true
}

// IsMaxLength checks if a string has maximum length
func (v *Validator) IsMaxLength(field, value string, maxLength int) bool {
	if len(value) > maxLength {
		_ = v.AddError(field, fmt.Sprintf("%s must be at most %d characters", field, maxLength))
		return false
	}
	return true
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string][]string {
	return v.errors
}

// Error returns the first error message
func (v *Validator) Error() string {
	for field, messages := range v.errors {
		if len(messages) > 0 {
			return fmt.Sprintf("%s: %s", field, messages[0])
		}
	}
	return ""
}
