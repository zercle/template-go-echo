package pkg

import "github.com/labstack/echo/v4"

// JSendStatus represents JSend status values
type JSendStatus string

const (
	StatusSuccess JSendStatus = "success"
	StatusFail    JSendStatus = "fail"
	StatusError   JSendStatus = "error"
)

// JSendResponse is the standard response format (JSend specification)
type JSendResponse struct {
	Status  JSendStatus `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// Success returns a successful JSend response
func Success(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	})
}

// SuccessWithMessage returns a successful JSend response with a message
func SuccessWithMessage(c echo.Context, statusCode int, data interface{}, message string) error {
	return c.JSON(statusCode, JSendResponse{
		Status:  StatusSuccess,
		Data:    data,
		Message: message,
	})
}

// Fail returns a fail JSend response (validation errors, etc)
func Fail(c echo.Context, statusCode int, data interface{}, message string) error {
	return c.JSON(statusCode, JSendResponse{
		Status:  StatusFail,
		Data:    data,
		Message: message,
	})
}

// Error returns an error JSend response (server errors)
func Error(c echo.Context, statusCode int, message string, code string) error {
	return c.JSON(statusCode, JSendResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
	})
}
