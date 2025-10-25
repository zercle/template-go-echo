.PHONY: help build test lint clean migrate run docker-build docker-up docker-down generate

## help: Show this help message
help:
	@echo "Template Go Echo - Development Commands"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application binary
build:
	@echo "Building application..."
	@go build -o tmp/template-go-echo ./cmd/api/
	@echo "✓ Build complete: tmp/template-go-echo"

## test: Run all tests with race detection
test:
	@echo "Running tests..."
	@go clean -testcache
	@go test -v -race ./...
	@echo "✓ Tests complete"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go clean -testcache
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

## lint: Run linters (golangci-lint)
lint:
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "✓ Linting complete"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .
	@echo "✓ Code formatted"

## generate: Run go generate for mocks and code generation
generate:
	@echo "Generating code..."
	@go generate ./...
	@echo "✓ Code generation complete"

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs/
	@echo "✓ Swagger documentation generated at docs/"

## generate-all: Generate all code (mocks, sqlc, swagger)
generate-all: generate swagger
	@echo "✓ All code generation complete"

## clean: Clean build artifacts and temporary files
clean:
	@echo "Cleaning..."
	@go clean -cache -testcache
	@rm -f tmp/template-go-echo coverage.out coverage.html
	@echo "✓ Clean complete"

## migrate-create: Create a new migration (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@timestamp=$$(date +%s); \
	mkdir -p sql/migrations; \
	touch sql/migrations/$${timestamp}_$(NAME).up.sql; \
	touch sql/migrations/$${timestamp}_$(NAME).down.sql; \
	echo "Created migration: sql/migrations/$${timestamp}_$(NAME).(up|down).sql"

## migrate-up: Run database migrations (requires MIGRATION_DSN env var)
migrate-up:
	@if [ -z "$(MIGRATION_DSN)" ]; then \
		echo "Error: MIGRATION_DSN environment variable is required"; \
		echo "Example: make migrate-up MIGRATION_DSN='mysql://user:pass@tcp(localhost:3306)/dbname'"; \
		exit 1; \
	fi
	@echo "Running migrations up..."
	@migrate -path sql/migrations -database "$(MIGRATION_DSN)" up
	@echo "✓ Migrations complete"

## migrate-down: Rollback database migrations (requires MIGRATION_DSN env var)
migrate-down:
	@if [ -z "$(MIGRATION_DSN)" ]; then \
		echo "Error: MIGRATION_DSN environment variable is required"; \
		exit 1; \
	fi
	@echo "Reverting migrations..."
	@migrate -path sql/migrations -database "$(MIGRATION_DSN)" down
	@echo "✓ Rollback complete"

## run: Run the application locally
run: build
	@echo "Starting application on :8080..."
	@./tmp/template-go-echo

## run-dev: Run the application with hot reload (requires air)
run-dev:
	@echo "Starting application with hot reload..."
	@air

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t template-go-echo:latest .
	@echo "✓ Docker image built"

## docker-up: Start Docker containers (requires docker-compose)
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "✓ Containers started. Database available at localhost:3306"

## docker-down: Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "✓ Containers stopped"

## docker-logs: View Docker container logs
docker-logs:
	@docker-compose logs -f

## quality: Run full quality checks (test, lint, build)
quality: test lint build
	@echo "✓ All quality checks passed"

## all: Build, test, lint (complete check)
all: clean generate quality
	@echo "✓ All checks complete - ready for production"

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/uber-go/mock/cmd/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✓ Tools installed"

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies tidied"

## version: Show version information
version:
	@echo "Go version: $$(go version)"
	@echo "Module: $$(go list -m)"
	@go list -m all | head -20
