package usecase

// CreateUserRequest is the request DTO for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest is the request DTO for updating a user.
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
}

// UserResponse is the response DTO for a user.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ListUsersResponse is the response DTO for listing users.
type ListUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int             `json:"total"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
}
