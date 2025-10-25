package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/config"
	"github.com/zercle/template-go-echo/pkg"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(cfg *config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return pkg.Error(c, http.StatusUnauthorized, "missing authorization header", pkg.ErrCodeUnauthorized)
			}

			// Extract token from Bearer scheme
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return pkg.Error(c, http.StatusUnauthorized, "invalid authorization header format", pkg.ErrCodeUnauthorized)
			}

			token := parts[1]

			// Parse and validate token
			claims := &Claims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Secret), nil
			})

			if err != nil || !parsedToken.Valid {
				slog.Warn("invalid token",
					slog.String("error", err.Error()),
				)
				return pkg.Error(c, http.StatusUnauthorized, "invalid token", pkg.ErrCodeUnauthorized)
			}

			// Store claims in context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// OptionalJWTAuth is an optional JWT middleware that doesn't fail if no token is present
func OptionalJWTAuth(cfg *config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No token, continue without authentication
				return next(c)
			}

			// Extract token from Bearer scheme
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, continue without authentication
				return next(c)
			}

			token := parts[1]

			// Parse and validate token
			claims := &Claims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Secret), nil
			})

			if err == nil && parsedToken.Valid {
				// Token is valid, store claims
				c.Set("user_id", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("claims", claims)
			}

			return next(c)
		}
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

// GetClaims extracts JWT claims from context
func GetClaims(c echo.Context) *Claims {
	claims, ok := c.Get("claims").(*Claims)
	if !ok {
		return nil
	}
	return claims
}
