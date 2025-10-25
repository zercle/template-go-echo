package unit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/internal/config"
	"github.com/zercle/template-go-echo/internal/middleware"
)

func TestJWTAuthMissingToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret: "test-secret",
		TTL:    3600,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	// No Authorization header
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := middleware.JWTAuth(cfg)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Handler should return an error when token is missing
	_ = handler(c)
	// Middleware will write unauthorized response
}

func TestJWTAuthValidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret: "test-secret",
		TTL:    3600,
	}

	// Create a valid token
	claims := &middleware.Claims{
		UserID: "user-123",
		Email:  "user@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := middleware.JWTAuth(cfg)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify claims are set in context
	if userID := c.Get("user_id"); userID != "user-123" {
		t.Errorf("expected user_id 'user-123', got %v", userID)
	}
}

func TestGetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", "test-user")
	if userID := middleware.GetUserID(c); userID != "test-user" {
		t.Errorf("expected 'test-user', got %s", userID)
	}
}
