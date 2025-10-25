package domain

const (
	// Password constraints
	MinPasswordLength = 8
	MaxPasswordLength = 128

	// Name constraints
	MinNameLength = 1
	MaxNameLength = 255

	// Email constraints
	MaxEmailLength = 255

	// Session constraints
	SessionDurationHours = 24 // 24-hour session duration
	TokenExpiryHours     = 1  // 1-hour token expiry
)

// ValidationMessages provides domain-specific validation messages
var ValidationMessages = map[string]string{
	"email_required":      "Email is required",
	"email_invalid":       "Email format is invalid",
	"email_exists":        "Email is already registered",
	"password_required":   "Password is required",
	"password_too_short":  "Password must be at least 8 characters",
	"password_too_long":   "Password must be at most 128 characters",
	"password_weak":       "Password must contain uppercase, lowercase, number, and special character",
	"name_required":       "Name is required",
	"name_too_short":      "Name must be at least 1 character",
	"name_too_long":       "Name must be at most 255 characters",
	"old_password_invalid": "Old password is incorrect",
}
