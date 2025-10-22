package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// JSend represents a JSend-formatted response.
// https://labs.omniti.com/labs/jsend
type JSend struct {
	Status string      `json:"status"`
	Data   any         `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
	Code   string      `json:"code,omitempty"`
}

// Success returns a successful JSend response.
func Success(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, JSend{
		Status: "success",
		Data:   data,
	})
}

// Error returns an error JSend response.
func Error(c echo.Context, statusCode int, message string, code string) error {
	return c.JSON(statusCode, JSend{
		Status: "error",
		Error:  message,
		Code:   code,
	})
}

// Fail returns a fail JSend response (client error).
func Fail(c echo.Context, statusCode int, message string, code string) error {
	return c.JSON(statusCode, JSend{
		Status: "fail",
		Error:  message,
		Code:   code,
	})
}

// OK returns a 200 OK success response.
func OK(c echo.Context, data any) error {
	return Success(c, http.StatusOK, data)
}

// Created returns a 201 Created success response.
func Created(c echo.Context, data any) error {
	return Success(c, http.StatusCreated, data)
}

// NoContent returns a 204 No Content response.
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// BadRequest returns a 400 Bad Request error response.
func BadRequest(c echo.Context, message string) error {
	return Fail(c, http.StatusBadRequest, message, "BAD_REQUEST")
}

// Unauthorized returns a 401 Unauthorized error response.
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, "UNAUTHORIZED")
}

// Forbidden returns a 403 Forbidden error response.
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, "FORBIDDEN")
}

// NotFound returns a 404 Not Found error response.
func NotFound(c echo.Context, message string) error {
	return Fail(c, http.StatusNotFound, message, "NOT_FOUND")
}

// Conflict returns a 409 Conflict error response.
func Conflict(c echo.Context, message string) error {
	return Fail(c, http.StatusConflict, message, "CONFLICT")
}

// InternalError returns a 500 Internal Server Error response.
func InternalError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message, "INTERNAL_ERROR")
}
