package infrastructure

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/health", Health)
	e.GET("/ready", Ready)
	e.GET("/live", Live)
}

// Health is the basic health check endpoint
// @Summary Health check
// @Description Basic health check endpoint
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(c echo.Context) error {
	return pkg.Success(c, http.StatusOK, HealthResponse{
		Status:  "healthy",
		Service: "template-go-echo",
	})
}

// Ready is the readiness check endpoint
// @Summary Readiness check
// @Description Check if service is ready to accept requests
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ready [get]
func Ready(c echo.Context) error {
	return pkg.Success(c, http.StatusOK, HealthResponse{
		Status:  "ready",
		Service: "template-go-echo",
	})
}

// Live is the liveness check endpoint
// @Summary Liveness check
// @Description Check if service is still alive
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /live [get]
func Live(c echo.Context) error {
	return pkg.Success(c, http.StatusOK, HealthResponse{
		Status:  "alive",
		Service: "template-go-echo",
	})
}
