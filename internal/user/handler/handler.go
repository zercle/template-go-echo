package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/user/usecase"
)

// Handler handles user HTTP requests.
type Handler struct {
	service *usecase.UserService
}

// NewHandler creates a new user handler.
func NewHandler(service *usecase.UserService) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the user routes on the Echo router.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api/v1/users")

	g.POST("", h.CreateUser)
	g.GET("", h.ListUsers)
	g.GET("/:id", h.GetUser)
	g.PUT("/:id", h.UpdateUser)
	g.DELETE("/:id", h.DeleteUser)
}
