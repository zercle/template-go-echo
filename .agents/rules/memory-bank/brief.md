# Project Brief: Go Echo Modular Monolith Template

## Project Identity
**Project Name:** template-go-echo
**Type:** Go Echo modular monolith Template
**Architecture:** Clean Architecture with Domain Driven Design
**Purpose:** A lightweight and performant foundation for building modular monolith applications that can evolve into microservices.

## Executive Summary
A production-ready Go backend template using the Echo framework and a Clean Architecture. It is designed for simplicity, performance, and ease of use, making it ideal for developing modular monolith applications with clear domain boundaries.

## Core Objectives
1.  **High Performance**: Leverage Echo's speed for low-latency services.
2.  **Simplicity**: Offer a straightforward structure that is easy to understand and extend.
3.  **Production-Ready**: Include essentials like logging, configuration, and error handling.
4.  **Developer Experience**: Ensure a minimal learning curve and quick setup.
5.  **Testability**: Structure the code to facilitate unit and integration testing.

## Project Scope

### In Scope
**Core Infrastructure:** Echo framework, Clean Architecture with DDD (Handler, Usecase, Repository, Domain), Database integration (sqlc + golang-migrate), JWT Authentication, Swagger/OpenAPI Documentation, Environment-based Config, Structured Logging, Centralized Error Handling, Health Checks, Graceful Shutdown, JSend Response Format, Interface Mocking.
**Development Tools:** Docker support, Makefile commands, Linting (golangci-lint v2), sqlc code generation, Database migrations.
**Documentation:** README, API (Swagger), Architecture Decision Records (ADRs), Example domain implementation.

## Technical Requirements

### Architecture Principles
**Clean Architecture with Domain-Driven Design:**
-   **Handler Layer:** Manages HTTP requests/responses, DTOs, and Swagger annotations.
-   **Usecase Layer:** Contains business logic and orchestrates domain operations.
-   **Repository Layer:** Handles database interactions and transactions via sqlc generated code.
-   **Domain Layer:** Defines interfaces, entities, value objects, and domain contracts.
-   **Mock Generation:** All interfaces must include `//go:generate` annotations for mockgen.
**Modular Monolith:** Start with modular structure that can transition to microservices.
**Dependency Injection:** Use `samber/do v2` for DI container management across all layers.

### Technology Stack
**Core Framework:** Go 1.25+, Echo v4+.
**Database Options:**
- MariaDB 11+ (small-medium systems, default)
- PostgreSQL 18+ (large systems)
- FerretDB 2.5+ (document-based storage)
- Valkey 9+ (in-memory key-value, Redis replacement)
**Database Tools:** `sqlc` for type-safe queries, `golang-migrate/migrate` for migrations.
**Authentication:** `golang-jwt/jwt` (primary), `zitadel/oidc` (alternative).
**Configuration:** Environment variables with validation via `spf13/viper`.
**Logging:** `slog` (structured logging with context propagation).
**Utilities:**
- `samber/do v2` (Dependency Injection)
- `samber/lo` (Synchronous helpers for finite sequences)
- `samber/ro` (Event-driven infinite data streams)
**Testing:** `uber-go/mock` (interface mocking), `DATA-DOG/go-sqlmock` (DB mocking).
**Code Quality:** `golangci-lint v2`, UUIDv7 for database-friendly unique IDs.

### Project Structure
```
template-go-echo/
├── cmd/api/                  # Application entry point (main.go)
├── internal/                 # Private application code (domain modules)
│   ├── {domain}/             # Per-domain organization (e.g., user, order)
│   │   ├── domain/           # Domain interfaces and contracts
│   │   ├── handler/          # HTTP handlers with Swagger annotations
│   │   ├── usecase/          # Business logic layer
│   │   ├── repository/       # Data access via sqlc (handles transactions)
│   │   ├── mock/             # Generated mock implementations
│   │   └── test/             # Test files using _test package convention
│   ├── infrastructure/       # Project infrastructure
│   │   └── sqlc/             # Generated SQL code from sqlc
│   ├── middleware/           # HTTP middleware
│   └── config/               # Configuration loading and validation
├── pkg/                      # Shared utilities (response, errors, custom types)
├── sql/                      # SQL source files
│   ├── queries/              # SQL queries for sqlc code generation
│   └── migrations/           # Database migration files (golang-migrate)
├── docs/                     # Swagger/OpenAPI specs and ADRs
├── .agents/                  # Agent rules and memory bank
│   └── rules/
│       └── memory-bank/      # Project context and documentation
├── bin/                      # Binary releases
├── tmp/                      # Temporary files and artifacts
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Development tasks (build, test, lint, migrate)
├── compose.yml               # Container compose
├── Dockerfile                # Container definition
├── .env.example              # Environment variable template
└── README.md                 # Project documentation
```

## Key Features
1.  **RESTful API Foundation**: Basic routing, request validation, and JSend-formatted JSON responses.
2.  **Authentication**: JWT-based middleware for securing endpoints (with OIDC alternative).
3.  **Type-Safe Database Access**: `sqlc` for generating type-safe Go code from SQL queries.
4.  **Database Migrations**: `golang-migrate/migrate` for version-controlled schema changes.
5.  **Configuration Management**: Environment-variable-driven with validation and defaults.
6.  **API Documentation**: Comprehensive Swagger/OpenAPI documentation with interactive UI, authentication support, and automated generation from handler annotations.
7.  **Testing Infrastructure**: Mock generation via `go:generate`, sqlmock for DB tests.
8.  **Containerization**: Production-ready Dockerfile with graceful shutdown support.
9.  **Performance Optimization**: Following goperf.dev patterns for networking and common operations.
10. **Security**: Rate limiting, input validation/sanitization, secure headers, HTTPS support.

## Constraints & Requirements
1.  **Code Quality:** Must pass `golangci-lint v2` with all checks enabled.
2.  **Testing:** High test coverage across all layers (handlers, usecases, repositories).
3.  **Mock Generation:** Every interface file must have `//go:generate` annotations.
4.  **Swagger Documentation:** Every HTTP handler must include Swagger annotations.
5.  **Security:** Rate limiting, input validation, sanitization, secure headers, HTTPS.
6.  **Error Handling:** Structured error types per domain with proper propagation.
7.  **Logging:** Structured logging with context propagation across services.
8.  **Monitoring:** Application metrics and health check endpoints.
9.  **Statelessness:** Services must be stateless for horizontal scaling.
10. **Container Support:** All services must support container environments.
11. **Graceful Shutdown:** Implement graceful shutdown for all services.

## Development Workflow
1.  **SQL Development**: Write queries in `/sql/queries` with sqlc annotations.
2.  **Code Generation**: Run `sqlc generate` to create Go code in domain `/sqlc` directories.
3.  **Database Migrations**: Use `migrate -path /sql/migrations -database <db_url> up`.
4.  **Testing**: Generate mocks with `go generate ./...`, write tests with sqlmock.
5.  **CI/CD Pipeline**: Automated testing, linting, security scanning, container builds.

## Documentation Requirements
1.  **API Documentation**: Maintain Swagger/OpenAPI specs via code annotations.
2.  **Architecture Decision Records (ADRs)**: Document significant design decisions.
3.  **Contribution Guidelines**: Clear setup and contribution instructions.
4.  **README**: Comprehensive project overview with quick-start guide.
 
## Quality Assurance
1. Generate mocks and tests with `go generate ./...`
2. Lint with `golangci-lint run --fix ./...`
3. Run tests with `go clean -testcache && go test -v -race ./...`