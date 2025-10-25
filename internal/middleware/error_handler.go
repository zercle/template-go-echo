package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg"
)

// ErrorHandler handles errors from Echo handlers
func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"
	errorCode := pkg.ErrCodeInternalError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	} else if de, ok := err.(*pkg.DomainError); ok {
		errorCode = de.Code
		message = de.Message

		// Map domain error codes to HTTP status codes
		switch de.Code {
		case pkg.ErrCodeNotFound:
			code = http.StatusNotFound
		case pkg.ErrCodeValidation:
			code = http.StatusBadRequest
		case pkg.ErrCodeConflict:
			code = http.StatusConflict
		case pkg.ErrCodeUnauthorized:
			code = http.StatusUnauthorized
		case pkg.ErrCodeForbidden:
			code = http.StatusForbidden
		case pkg.ErrCodeBadRequest:
			code = http.StatusBadRequest
		}
	}

	// Log the error
	slog.Error("HTTP request error",
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.Int("code", code),
		slog.String("error", err.Error()),
	)

	// Send response
	if !c.Response().Committed {
		_ = pkg.Error(c, code, message, errorCode)
	}
}
