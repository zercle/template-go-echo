# Technology Stack: Go Echo Modular Monolith Template

## Core Technologies

### Runtime Environment
- **Go**: 1.25.3 (minimum version)
- **Platform**: Cross-platform (Linux, macOS, Windows)
- **Architecture**: 64-bit (amd64, arm64)

### Web Framework
- **Echo**: v4.12.0+ (high-performance HTTP framework)
  - Key Features: Middleware support, route groups, context handling
  - Benefits: Performance, simplicity, extensibility
  - Alternatives: Gin, Fiber (not recommended for this template)

### Database Stack

#### Primary Database Options
1. **MariaDB**: 11.0+ (default for small-to-medium systems)
   - Benefits: MySQL compatibility, performance, open-source
   - Use Case: General-purpose applications, startups, SMB

2. **PostgreSQL**: 18.0+ (recommended for large systems)
   - Benefits: Advanced features, extensibility, robustness
   - Use Case: Enterprise applications, complex queries, high scalability

#### Alternative Storage
3. **FerretDB**: 2.5+ (MongoDB-compatible document storage)
   - Benefits: Open-source, PostgreSQL backend, MongoDB compatibility
   - Use Case: Document-based data, flexible schemas, rapid prototyping

4. **Valkey**: 9.0+ (Redis-compatible in-memory storage)
   - Benefits: Performance, persistence, clustering
   - Use Case: Caching, session storage, real-time features

## Database Tools

### Code Generation
- **sqlc**: v1.28.0+ (type-safe SQL code generation)
  - Benefits: Compile-time safety, performance, no runtime reflection
  - Integration: Generated Go code from SQL queries
  - Workflow: SQL queries → Go types → type-safe database operations

### Migration Management
- **golang-migrate/migrate**: v4.17.0+ (database version control)
  - Benefits: Versioned migrations, rollback support, multiple databases
  - Workflow: Migration files → Schema changes → Version tracking

## Authentication & Security

### Primary Authentication
- **golang-jwt/jwt**: v5.2.0+ (JWT token management)
  - Features: Access tokens, refresh tokens, validation
  - Benefits: Stateless, standard-based, scalable

### Alternative Authentication
- **zitadel/oidc**: v2.0.0+ (OpenID Connect)
  - Features: External identity providers, SSO support
  - Use Case: Enterprise SSO, external authentication

### Security Libraries
- **Rate Limiting**: golang.org/x/time rate limiting
- **Input Validation**: Custom validation with struct tags
- **CORS**: Echo CORS middleware
- **Security Headers**: Echo security middleware

## Dependency Management

### Dependency Injection
- **samber/do**: v2.0.0+ (dependency injection container)
  - Benefits: Clean code, testability, explicit dependencies
  - Features: Auto-wiring, health checks, lifecycle management

### Utility Libraries
- **samber/lo**: v2.6.0+ (synchronous helpers for sequences)
  - Benefits: Functional programming helpers, data manipulation
  - Use Case: Data transformation, collection operations

- **samber/ro**: v1.3.0+ (event-driven infinite streams)
  - Benefits: Reactive programming, stream processing
  - Use Case: Event handling, real-time data processing

## Development Tools

### Code Quality
- **golangci-lint**: v1.55.0+ (Go linting tool)
  - Configuration: Custom ruleset for clean architecture
  - Integration: Pre-commit hooks, CI/CD pipeline
  - Benefits: Code consistency, bug prevention, best practices

### Mock Generation
- **uber-go/mock**: v1.0.0+ (interface mocking)
  - Integration: `//go:generate` annotations for auto-generation
  - Benefits: Easy testing, interface-based design, TDD support

### Database Testing
- **DATA-DOG/go-sqlmock**: v1.5.0+ (database mocking)
  - Benefits: Isolated tests, no real database required
  - Use Case: Repository layer testing, transaction testing

## API Documentation

### Swagger/OpenAPI Stack
- **swaggo/swag**: v1.16.6+ (Swagger/OpenAPI code generation)
  - Features: Automated documentation from Go annotations
  - Integration: Generates JSON/YAML specifications and Go code
  - Benefits: Type-safe documentation, interactive UI, version control
- **swaggo/echo-swagger**: v1.4.1+ (Echo middleware for Swagger)
  - Features: Seamless Echo integration for Swagger UI
  - Benefits: Easy setup, interactive API testing
- **swaggo/files**: v1.0.1+ (Static file serving for Swagger UI)
  - Features: Static file serving for generated documentation
  - Benefits: Efficient serving, cache support

### Documentation Features
- **Interactive UI**: Swagger UI provides interactive API testing
- **Authentication Support**: JWT Bearer token authentication in documentation
- **Auto-Generation**: Documentation automatically generated from handler annotations
- **Version Control**: API documentation versioned with codebase
- **Multiple Formats**: Support for JSON, YAML, and interactive HTML formats

## Configuration Management

### Configuration Library
- **spf13/viper**: v1.17.0+ (configuration management)
  - Features: Environment variables, config files, validation
  - Integration: Environment-based configuration with defaults
  - Benefits: Flexible configuration, multiple sources, validation

### Environment Variables
- **Required**: Database connection, JWT secrets, server ports
- **Optional**: Logging levels, feature flags, external service URLs
- **Development**: Local environment configuration
- **Production**: Container environment configuration

## Logging & Monitoring

### Structured Logging
- **slog**: Go 1.21+ standard library (contextual logging)
  - Benefits: Performance, context propagation, structured output
  - Integration: Request tracing, error logging, audit trails
  - Formats: JSON, text output configurable

### Health Checks
- **Echo health endpoints**: `/health`, `/ready`, `/live`
- **Database health**: Connection pool status, query validation
- **Service health**: External dependency checks

## Performance Optimization

### Performance Guidelines
- **Patterns**: Based on goperf.dev best practices
- **Networking**: Connection pooling, keep-alives, timeouts
- **Memory**: Efficient data structures, garbage collection tuning
- **Concurrency**: Goroutines, channels, sync primitives

### Optimization Tools
- **Profiling**: pprof integration for performance analysis
- **Metrics**: OpenTelemetry integration (optional)
- **Tracing**: Request tracing with context propagation

## Container & Deployment

### Containerization
- **Docker**: Multi-stage builds for production images
- **Base Image**: alpine:3.18+ (minimal, secure)
- **Features**: Health checks, graceful shutdown, non-root user

### Orchestration Support
- **Docker Compose**: Development environment setup
- **Kubernetes**: Deployment manifests, health checks
- **Environment**: Configuration via environment variables, secrets

## Development Environment

### Required Tools
- **Go**: 1.25.3+ (runtime and toolchain)
- **git**: Version control
- **Docker**: Container development and testing
- **sqlc**: SQL code generation
- **golang-migrate**: Database migrations
- **golangci-lint**: Code quality checks
- **swag**: Swagger/OpenAPI documentation generation

### Development Scripts
- **Makefile**: Common development tasks
  - `make build`: Build application
  - `make test`: Run tests with coverage
  - `make lint`: Run code quality checks
  - `make migrate`: Run database migrations
  - `make generate`: Generate code (mocks, sqlc)
  - `make swagger`: Generate Swagger documentation
  - `make generate-all`: Generate all code (mocks, sqlc, swagger)

### IDE Integration
- **VS Code**: Go extension, Docker extension
- **GoLand**: JetBrains Go IDE
- **Development Features**: Debugging, testing, refactoring

## Version Control

### Git Configuration
- **.gitignore**: Comprehensive Go project ignore rules
- **Branch Strategy**: GitFlow or similar branching model
- **Hooks**: Pre-commit linting and testing

### Dependency Management
- **go.mod**: Module definition and version control
- **go.sum**: Dependency checksums
- **Updates**: Regular dependency updates and security patches

## Testing Strategy

### Test Organization
- **Unit Tests**: Individual component testing
- **Integration Tests**: Database and external service testing
- **End-to-End Tests**: Full application workflow testing
- **Benchmark Tests**: Performance testing and optimization

### Test Tools
- **Testing Framework**: Go standard testing package
- **Coverage**: Built-in coverage reporting
- **Assertions**: Custom assertion helpers
- **Test Data**: Fixtures and test data management

## Security Considerations

### Input Validation
- **HTTP Validation**: Request body validation
- **SQL Injection**: sqlc provides protection via parameterized queries
- **XSS Protection**: Output encoding, CSP headers
- **CSRF Protection**: Token-based CSRF protection

### Authentication Security
- **JWT Security**: Secure token generation and validation
- **Password Security**: Hashing, salt, rate limiting
- **Session Management**: Secure session handling
- **API Security**: Rate limiting, API keys, authentication

## Environment Requirements

### Development Environment
- **CPU**: 2+ cores recommended
- **Memory**: 8GB+ RAM recommended
- **Storage**: 10GB+ free space
- **Network**: Internet access for dependencies

### Production Environment
- **CPU**: 4+ cores for production workloads
- **Memory**: 16GB+ RAM recommended
- **Storage**: SSD storage preferred
- **Network**: Low latency, high bandwidth

### Scaling Requirements
- **Horizontal Scaling**: Stateless design for load balancing
- **Vertical Scaling**: Performance optimization and resource management
- **Database Scaling**: Read replicas, sharding strategies
- **Caching**: Multi-level caching for performance

## Monitoring & Observability

### Application Monitoring
- **Health Endpoints**: Application health and readiness
- **Metrics Collection**: Custom metrics, performance counters
- **Error Tracking**: Structured error logging and aggregation
- **Request Tracing**: Distributed tracing for debugging

### Infrastructure Monitoring
- **System Metrics**: CPU, memory, disk usage
- **Database Monitoring**: Query performance, connection pools
- **Network Monitoring**: Latency, throughput, error rates
- **Container Monitoring**: Resource usage, health checks