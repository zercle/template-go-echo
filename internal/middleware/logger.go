package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg/logger"
)

// Logger is a middleware that logs HTTP requests.
func Logger(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Generate request ID
			requestID := uuid.New().String()
			c.Set("request_id", requestID)

			// Add request ID to logger
			requestLogger := log.WithRequestID(requestID)

			// Start timer
			start := time.Now()

			// Log request
			requestLogger.Info("request started",
				"method", c.Request().Method,
				"path", c.Request().RequestURI,
				"ip", c.RealIP(),
			)

			// Call next handler
			err := next(c)

			// Log response
			duration := time.Since(start)
			statusCode := c.Response().Status
			requestLogger.Info("request completed",
				"method", c.Request().Method,
				"path", c.Request().RequestURI,
				"status_code", statusCode,
				"duration_ms", duration.Milliseconds(),
			)

			return err
		}
	}
}
