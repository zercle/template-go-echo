.PHONY: help build run test lint clean cover generate swagger fmt docker-build docker-run install-tools check

help:
	@echo "template-go-echo Makefile commands:"
	@echo ""
	@echo "  make build              Build the application binary"
	@echo "  make run                Run the application locally"
	@echo "  make test               Run all tests"
	@echo "  make cover              Generate test coverage report"
	@echo "  make lint               Run golangci-lint checks"
	@echo "  make fmt                Format code with gofmt"
	@echo "  make clean              Clean up build artifacts"
	@echo "  make generate           Run code generation (mocks, sqlc)"
	@echo "  make swagger            Generate Swagger documentation"
	@echo "  make docker-build       Build Docker image"
	@echo "  make docker-run         Run Docker container"
	@echo "  make help               Show this help message"

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "Build complete: ./bin/api"

# Run the application
run: build
	@echo "Starting application..."
	@./bin/api

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate test coverage report
cover:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run golangci-lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -w -s .
	@go mod tidy
	@echo "Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

# Generate code (mocks, etc.)
generate:
	@echo "Generating code..."
	@go generate ./...
	@echo "Code generation complete"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go
	@echo "Swagger documentation generated"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t template-go-echo:latest .
	@echo "Docker image built: template-go-echo:latest"

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 \
		-e PORT=8080 \
		-e ENV=development \
		-e LOG_LEVEL=debug \
		template-go-echo:latest

# Install development dependencies
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/cmd/mockgen@latest
	@echo "Tools installed"

# Quick check (format, lint, test)
check: fmt lint test
	@echo "✓ All checks passed"

.DEFAULT_GOAL := help
