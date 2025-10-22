package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg/response"
)

// GetUser handles GET /api/v1/users/:id
// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} usecase.UserResponse
// @Failure 404 {object} response.JSend
// @Failure 500 {object} response.JSend
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "user id is required")
	}

	user, err := h.service.GetUser(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	return response.OK(c, user)
}
