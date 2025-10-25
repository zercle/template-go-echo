package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"github.com/zercle/template-go-echo/docs"
	"github.com/zercle/template-go-echo/internal/config"
	"github.com/zercle/template-go-echo/internal/infrastructure"
	"github.com/zercle/template-go-echo/internal/middleware"
)

// @title Go Echo Template API
// @version 1.0
// @description Production-ready Go Echo modular monolith template with clean architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize Swagger documentation
	docs.SwaggerInfo.Host = "localhost:8080"

	// Load configuration
	cfg := config.Load()

	// Create Echo instance
	e := echo.New()

	// Set error handler
	e.HTTPErrorHandler = middleware.ErrorHandler

	// Register middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Timeout(30 * time.Second))
	e.Use(middleware.CORS())
	e.Use(middleware.SecurityHeaders())

	// Register health check routes
	infrastructure.RegisterHealthRoutes(e)

	// Register Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	log.Fatal(e.Start(cfg.Server.Address))
}
