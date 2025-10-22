package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg/errors"
	"github.com/zercle/template-go-echo/pkg/response"
)

// ErrorHandler is a middleware that handles panics and errors.
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					code := errors.GetCode(e)
					message := errors.GetMessage(e)
					statusCode := errors.HTTPStatusCode(e)
					_ = response.Error(c, statusCode, message, code)
					return
				}

				_ = response.InternalError(c, "an unexpected error occurred")
			}
		}()

		err := next(c)
		if err != nil {
			// Handle Echo's HTTPError
			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == http.StatusMethodNotAllowed {
					return response.Error(c, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
				}
				if he.Code == http.StatusNotFound {
					return response.Fail(c, http.StatusNotFound, "endpoint not found", "NOT_FOUND")
				}
				return response.InternalError(c, he.Message.(string))
			}

			// Handle domain errors
			code := errors.GetCode(err)
			message := errors.GetMessage(err)
			statusCode := errors.HTTPStatusCode(err)

			return response.Error(c, statusCode, message, code)
		}

		return nil
	}
}
