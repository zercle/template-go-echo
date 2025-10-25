# Template Go Echo

A **production-ready Go backend template** using the Echo framework with Clean Architecture and Domain-Driven Design (DDD) principles.

> Build fast, scalable REST APIs with clear separation of concerns and maintainable architecture from day one.

## 🎯 Features

- **Clean Architecture**: Handler → Usecase → Repository → Domain layers
- **Domain-Driven Design**: Modular domain-organized structure
- **Echo Framework**: High-performance HTTP framework with middleware support
- **Type-Safe Database**: sqlc for compile-time safe SQL operations
- **JWT Authentication**: Built-in JWT middleware with refresh token support
- **Database Migrations**: Automatic schema versioning with golang-migrate
- **Error Handling**: Structured error types with JSend response format
- **Logging**: Structured logging with slog
- **Rate Limiting**: IP-based request throttling
- **Health Checks**: `/health`, `/ready`, `/live` endpoints
- **Testing**: Comprehensive test infrastructure with mocks
- **Docker Support**: Multi-stage Dockerfile + docker-compose
- **Development Tools**: Makefile, code generation, linting

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (optional, for database)
- MySQL 8.0+ or MariaDB 11+

### 1. Clone the Template

```bash
git clone <repo-url> my-api
cd my-api
```

### 2. Install Dependencies

```bash
make deps
make install-tools
```

### 3. Start Database

```bash
# Using Docker Compose (recommended)
make docker-up

# Or run MariaDB locally
mysql -u root -p < /dev/null  # Create local MariaDB instance
```

### 4. Run Migrations

```bash
export MIGRATION_DSN="mysql://appuser:app-password@tcp(localhost:3306)/template_go_echo"
make migrate-up
```

### 5. Run Application

```bash
make run
```

Application starts on `http://localhost:8080`

## 📁 Project Structure

```
template-go-echo/
├── cmd/api/                      # Application entry point
│   └── main.go                  # Server initialization
├── internal/                     # Private application code
│   ├── config/                  # Configuration management
│   ├── middleware/              # HTTP middleware
│   ├── infrastructure/          # Database, cache, external services
│   │   └── database/           # Database connection pooling
│   ├── user/                    # Example domain
│   │   ├── domain/             # Interfaces, entities, errors
│   │   ├── usecase/            # Business logic
│   │   ├── repository/         # Data access layer
│   │   ├── handler/            # HTTP handlers
│   │   └── test/               # Test files
├── pkg/                         # Shared utilities
│   ├── response.go             # JSend response format
│   ├── errors.go               # Domain error types
│   └── validation.go           # Input validation
├── sql/                        # Database code
│   ├── migrations/             # Version-controlled schema
│   └── queries/                # SQL queries for sqlc
├── docs/                       # API documentation
├── Makefile                    # Development commands
├── Dockerfile                  # Container image
├── compose.yml                 # Docker Compose configuration
└── README.md                   # This file
```

## 🔧 Development Workflow

### Build

```bash
make build              # Build binary to tmp/template-go-echo
```

### Testing

```bash
make test              # Run all tests with race detection
make test-coverage     # Generate coverage report
make quality           # Run complete quality checks
```

### Linting

```bash
make lint              # Run golangci-lint
make fmt               # Format code
```

### Code Generation

```bash
make generate          # Run go:generate for mocks
```

### Database

```bash
# Create migration
make migrate-create NAME=add_users_table

# Run migrations
make migrate-up MIGRATION_DSN="mysql://user:pass@tcp(localhost:3306)/dbname"

# Rollback
make migrate-down MIGRATION_DSN="mysql://user:pass@tcp(localhost:3306)/dbname"
```

### Running

```bash
make run               # Run with binary build
make run-dev          # Run with hot reload (requires air)
```

### Docker

```bash
make docker-build      # Build Docker image
make docker-up         # Start containers
make docker-down       # Stop containers
make docker-logs       # View logs
```

## 📚 API Endpoints

### Authentication

- `POST /api/v1/users/register` - Create new user account
- `POST /api/v1/users/login` - Login and get tokens
- `POST /api/v1/users/token/refresh` - Refresh access token

### Users (Protected)

- `GET /api/v1/users` - List all users (paginated)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user profile
- `POST /api/v1/users/:id/password` - Change password
- `DELETE /api/v1/users/:id` - Delete user
- `POST /api/v1/users/logout` - Logout current session
- `POST /api/v1/users/logout-all` - Logout all sessions

### Health

- `GET /health` - Health status
- `GET /ready` - Readiness check
- `GET /live` - Liveness check

## 🔐 Configuration

Environment variables (see `.env.example`):

```bash
# Server
SERVER_ADDRESS=:8080                    # Server address:port
SERVER_PORT=8080                        # Server port (deprecated)
SERVER_DEBUG=false                      # Debug mode

# Database
DB_DRIVER=mysql                         # mysql or postgres
DB_DSN=root:password@tcp(localhost:3306)/template_go_echo
DB_MAX_CONNS=10                        # Connection pool size

# JWT
JWT_SECRET=your-secret-key             # CHANGE IN PRODUCTION!
JWT_TTL=3600                           # Token expiry in seconds
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Run specific package
go test -v ./internal/user/...
```

### Integration Tests

Tests use mock repositories for isolation. See `internal/user/usecase/usecase_test.go` for examples.

### Test Coverage

```bash
make test-coverage
open coverage.html
```

## 🏗️ Architecture Decisions

### Clean Architecture

```
HTTP Request
    ↓
Handler (validation, binding)
    ↓
Usecase (business logic)
    ↓
Repository (data access)
    ↓
Database
```

### Domain-Driven Design

- Each domain (user, product, order) is independently organized
- Interfaces define contracts between layers
- Errors and constants are domain-specific
- Tests are co-located with implementation

### Dependency Injection

- Minimal DI (no framework for simplicity)
- Constructor injection for testability
- Interfaces for all external dependencies

## 📖 Adding a New Domain

### 1. Create Domain Structure

```bash
mkdir -p internal/myfeature/{domain,usecase,repository,handler,test}
```

### 2. Define Domain

Create `internal/myfeature/domain/`:

- `entity.go` - Data models
- `interfaces.go` - Repository & Usecase interfaces
- `errors.go` - Domain-specific errors
- `constants.go` - Domain constants

### 3. Implement Business Logic

Create `internal/myfeature/usecase/usecase.go`:

```go
type MyFeatureUsecase struct {
    repo domain.MyRepository
}

func (u *MyFeatureUsecase) SomeOperation(ctx context.Context) error {
    // Business logic here
}
```

### 4. Implement Data Access

Create `internal/myfeature/repository/repository.go`:

```go
type MyRepository struct {
    db *sql.DB
}

func (r *MyRepository) GetItem(ctx context.Context) (*domain.Item, error) {
    // Database operations
}
```

### 5. Create HTTP Handlers

Create `internal/myfeature/handler/`:

- `handler.go` - HTTP endpoints
- `dto.go` - Request/Response DTOs
- `handler_test.go` - Handler tests

### 6. Register Routes

In `cmd/api/main.go`:

```go
handler := handler.New(usecase)
handler.RegisterRoutes(e)
```

### 7. Write Tests

Create comprehensive tests in `internal/myfeature/test/`:

- Unit tests for business logic
- Mock repositories for isolation
- Handler tests for HTTP layer

## 🚢 Deployment

### Environment Preparation

```bash
# Production environment variables
export SERVER_DEBUG=false
export JWT_SECRET=$(openssl rand -base64 32)  # Generate secure secret
export DB_DSN=user:password@tcp(prod-db:3306)/myapp
```

### Build for Production

```bash
make docker-build
docker push myregistry/template-go-echo:latest
```

### Run with Docker Compose

```yaml
# production-compose.yml
services:
  app:
    image: myregistry/template-go-echo:latest
    environment:
      DB_DSN: ${DB_DSN}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - "8080:8080"
```

### Health Checks

Kubernetes/Docker will use:

- `/health` - Startup check
- `/ready` - Readiness probe
- `/live` - Liveness probe

## 🔍 Monitoring

### Structured Logging

All logs are structured with slog:

```json
{"time":"...","level":"INFO","msg":"user logged in","user_id":"...","email":"..."}
```

### Metrics

Implement custom metrics by adding HTTP handlers:

```go
e.GET("/metrics", prometheusHandler)
```

## 🤝 Contributing

1. Create feature branch: `git checkout -b feature/myfeature`
2. Implement feature with tests
3. Run quality checks: `make quality`
4. Submit pull request

## 📝 License

MIT License - see LICENSE file

## 🆘 Support

- **Docs**: Check `docs/` directory
- **Issues**: Report bugs via GitHub
- **Templates**: See `internal/user/` for working example

## ✨ What's Next?

1. **Customize Configuration**: Update `.env` for your needs
2. **Add Domains**: Use `internal/user/` as template
3. **Write Tests**: Ensure > 80% coverage
4. **Setup CI/CD**: Integrate with GitHub Actions or similar
5. **Document APIs**: Add Swagger annotations to handlers
6. **Deploy**: Use Docker Compose or Kubernetes

---

**Built with Go 1.25+ and Echo Framework** 🚀
