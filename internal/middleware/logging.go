package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestLogger creates a logging middleware using slog
func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			slog.Info("HTTP request",
				slog.String("uri", values.URI),
				slog.String("method", values.Method),
				slog.Int("status", values.Status),
				slog.String("remote_ip", values.RemoteIP),
				slog.Duration("latency", values.Latency),
				slog.String("user_agent", values.UserAgent),
			)
			return nil
		},
		Skipper: func(c echo.Context) bool {
			// Skip logging for health check endpoints
			return c.Path() == "/health" || c.Path() == "/ready" || c.Path() == "/live"
		},
	})
}

// RequestID middleware adds a request ID to each request context
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return time.Now().Format("20060102150405-000")
		},
	})
}

// TimeoutMiddleware sets request timeout
func Timeout(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: timeout,
	})
}

// CORS middleware
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	})
}

// SecurityHeaders middleware
func SecurityHeaders() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	})
}
