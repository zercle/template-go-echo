package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/user/usecase"
	"github.com/zercle/template-go-echo/pkg/response"
)

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user with the provided name and email
// @Tags users
// @Accept json
// @Produce json
// @Param request body usecase.CreateUserRequest true "User data"
// @Success 201 {object} usecase.UserResponse
// @Failure 400 {object} response.JSend
// @Failure 409 {object} response.JSend
// @Failure 500 {object} response.JSend
// @Router /api/v1/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var req usecase.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Name == "" || req.Email == "" {
		return response.BadRequest(c, "name and email are required")
	}

	user, err := h.service.CreateUser(c.Request().Context(), req)
	if err != nil {
		switch err.Error() {
		case "USER_ALREADY_EXISTS: user already exists":
			return response.Conflict(c, "user with this email already exists")
		default:
			return response.InternalError(c, "failed to create user")
		}
	}

	return response.Created(c, user)
}
