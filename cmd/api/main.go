// @title			template-go-echo API
// @version		1.0.0
// @description	Production-ready Go microservice template using Clean Architecture and DDD
// @host			localhost:8080
// @basePath		/
// @schemes		http https
// @securityDefinitions.apiKey	Bearer
// @in						header
// @name					Authorization
// @description			JWT Bearer token
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/zercle/template-go-echo/docs"
	"github.com/zercle/template-go-echo/internal/config"
	"github.com/zercle/template-go-echo/internal/middleware"
	"github.com/zercle/template-go-echo/internal/user/handler"
	"github.com/zercle/template-go-echo/internal/user/repository"
	"github.com/zercle/template-go-echo/internal/user/usecase"
	"github.com/zercle/template-go-echo/pkg/database"
	"github.com/zercle/template-go-echo/pkg/logger"
	"github.com/zercle/template-go-echo/pkg/response"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level)
	log.Info("Starting application",
		"environment", cfg.Server.Env,
		"port", cfg.Server.Port,
	)

	// Initialize database connection
	db, err := database.Open(database.Config{
		Driver: cfg.DB.Driver,
		DSN:    cfg.DB.DSN,
	})
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close(db) //nolint:errcheck

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Register middleware
	e.Use(middleware.Logger(log))
	e.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	e.Use(middleware.ErrorHandler)

	// Register routes
	if err := registerRoutes(e, cfg, log, db); err != nil {
		log.Error("Failed to register routes", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := ":" + cfg.Server.Port
		log.Info("Listening on", "address", addr)

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("Shutdown error", "error", err)
		os.Exit(1)
	}

	log.Info("Server stopped")
}

func registerRoutes(e *echo.Echo, cfg *config.Config, log *logger.Logger, db *sql.DB) error {
	// Set swagger configuration based on environment
	docs.SwaggerInfo.Host = "localhost:" + cfg.Server.Port
	if cfg.IsProduction() {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	// Swagger documentation route
	// @Summary Swagger UI
	// @Description API documentation using Swagger UI
	// @Tags System
	// @Produce json
	// @Router /swagger/* [get]
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check endpoint
	// @Summary Health Check
	// @Description Returns the health status of the service
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	e.GET("/health", func(c echo.Context) error {
		return response.OK(c, map[string]string{
			"status": "healthy",
		})
	})

	// Ready check endpoint
	// @Summary Readiness Check
	// @Description Returns the readiness status of the service
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /ready [get]
	e.GET("/ready", func(c echo.Context) error {
		return response.OK(c, map[string]string{
			"status": "ready",
		})
	})

	// API root endpoint
	// @Summary API Root
	// @Description Returns welcome message and API version
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router / [get]
	e.GET("/", func(c echo.Context) error {
		return response.OK(c, map[string]string{
			"message": "Welcome to template-go-echo API",
			"version": "1.0.0",
		})
	})

	// User domain routes
	if err := registerUserRoutes(e, log, db); err != nil {
		return err
	}

	return nil
}

func registerUserRoutes(e *echo.Echo, log *logger.Logger, db *sql.DB) error {
	// Initialize user repository (database-based)
	userRepo := repository.NewDatabaseRepository(db)

	// Initialize user service
	userService := usecase.NewUserService(userRepo, log)

	// Initialize user handler
	userHandler := handler.NewHandler(userService)

	// Register user routes
	userHandler.RegisterRoutes(e)

	return nil
}
