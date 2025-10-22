package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg/response"
)

// DeleteUser handles DELETE /api/v1/users/:id
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} response.JSend
// @Failure 500 {object} response.JSend
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "user id is required")
	}

	err := h.service.DeleteUser(c.Request().Context(), id)
	if err != nil {
		switch err.Error() {
		case "USER_NOT_FOUND: user not found":
			return response.NotFound(c, "user not found")
		default:
			return response.InternalError(c, "failed to delete user")
		}
	}

	return response.NoContent(c)
}
