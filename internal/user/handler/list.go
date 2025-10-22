package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg/response"
)

// ListUsers handles GET /api/v1/users
// @Summary List all users
// @Description Get a paginated list of all users
// @Tags users
// @Produce json
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Success 200 {object} usecase.ListUsersResponse
// @Failure 500 {object} response.JSend
// @Router /api/v1/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	offset := 0
	limit := 10

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	users, err := h.service.ListUsers(c.Request().Context(), offset, limit)
	if err != nil {
		return response.InternalError(c, "failed to list users")
	}

	return response.OK(c, users)
}
