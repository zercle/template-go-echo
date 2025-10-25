package middleware

import (
	"net/http"
	"sync"

	ecmiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"github.com/zercle/template-go-echo/pkg"
	"golang.org/x/time/rate"
)

// RateLimiter represents a simple rate limiter using token bucket algorithm
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// Limit implements IP-based rate limiting
func (rl *RateLimiter) Limit(rps float64, burst int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			rl.mu.RLock()
			limiter, exists := rl.limiters[ip]
			rl.mu.RUnlock()

			if !exists {
				limiter = rate.NewLimiter(rate.Limit(rps), burst)

				rl.mu.Lock()
				rl.limiters[ip] = limiter
				rl.mu.Unlock()
			}

			if !limiter.Allow() {
				return pkg.Error(c, http.StatusTooManyRequests, "rate limit exceeded", "RATE_LIMIT_EXCEEDED")
			}

			return next(c)
		}
	}
}

// BodyLimitMiddleware limits request body size
func BodyLimitMiddleware(limit string) echo.MiddlewareFunc {
	return ecmiddleware.BodyLimit(limit)
}
