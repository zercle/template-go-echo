package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/user/usecase"
	"github.com/zercle/template-go-echo/pkg/response"
)

// UpdateUser handles PUT /api/v1/users/:id
// @Summary Update a user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body usecase.UpdateUserRequest true "Updated user data"
// @Success 200 {object} usecase.UserResponse
// @Failure 400 {object} response.JSend
// @Failure 404 {object} response.JSend
// @Failure 500 {object} response.JSend
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "user id is required")
	}

	var req usecase.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	user, err := h.service.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		switch err.Error() {
		case "USER_NOT_FOUND: user not found":
			return response.NotFound(c, "user not found")
		default:
			return response.InternalError(c, "failed to update user")
		}
	}

	return response.OK(c, user)
}
