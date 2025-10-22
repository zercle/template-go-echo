# Technology Stack: template-go-echo

## Language & Runtime
**Go Version:** 1.25.3 (minimum 1.25+)
- Fast compilation
- Excellent performance
- Strong concurrency support
- Static typing with good inference
- Simple deployment (single binary)

**Go Modules:** go.mod + go.sum (module-based dependency management)

## Core Web Framework
**Echo v4+** (Ultra-fast and minimalist web framework)
- **Package:** github.com/labstack/echo/v4
- **Why:** Lightweight, high-performance, excellent routing
- **Features:**
  - Fast HTTP router
  - Built-in middleware
  - Data binding and validation
  - HTTP error handling
  - Request/response context

## Database Layer

### Database Backends (Choose One)
- **MariaDB 11+** (default, recommended for small-medium systems)
- **PostgreSQL 18+** (recommended for large systems)
- **FerretDB 2.5+** (document-based, MongoDB-compatible)
- **Valkey 9+** (Redis replacement, in-memory key-value)

### Database Tools
**sqlc** - Type-safe SQL query generation
- **Purpose:** Generate Go code from SQL queries
- **Benefits:** Type safety, compile-time verification, zero runtime overhead
- **Usage:** Write SQL in `/sql/queries/`, run `sqlc generate`
- **Version:** Latest stable (v1.x)

**golang-migrate/migrate** - Database migration management
- **Purpose:** Version control for database schema
- **Benefits:** Reproducible deployments, rollback capability
- **Usage:** Migration files in `/sql/migrations/`
- **Version:** Latest stable

## Dependency Injection
**samber/do v2** - Simple DI container
- **Package:** github.com/samber/do/v2
- **Why:** Type-safe, minimal magic, excellent performance
- **Usage:** Constructor injection throughout application
- **Benefits:** Explicit dependencies, easy testing, no reflection overhead

## Authentication & Authorization
**Primary:** golang-jwt/jwt (JWT token validation)
- **Package:** github.com/golang-jwt/jwt/v5
- **Usage:** Token generation, validation, refresh
- **Scope:** Authorization header token validation

**Alternative (Not Primary):** zitadel/oidc (OIDC support for future)
- **Package:** github.com/zitadel/oidc/v3
- **Status:** Optional, for enterprise SSO scenarios
- **Note:** Not required for initial implementation

## Logging
**slog** (Standard Go logging)
- **Package:** Built-in log/slog (Go 1.21+)
- **Why:** Standard library, structured logging, context propagation
- **Wrapper:** Create internal/logger wrapper for consistency
- **Features:**
  - Structured logging (JSON compatible)
  - Context propagation (request IDs)
  - Multiple log levels
  - Custom handlers possible

## Utilities & Helpers

### samber Libraries
**samber/lo** - Utility functions for finite sequences
- **Package:** github.com/samber/lo
- **Usage:** Filter, map, reduce, other functional helpers
- **Benefit:** Functional programming utilities without bloat

**samber/ro** - Event-driven streams for infinite sequences
- **Package:** github.com/samber/ro
- **Status:** Optional, for advanced streaming scenarios
- **Note:** Consider for background job processing

### Configuration Management
**spf13/viper** - Configuration management
- **Package:** github.com/spf13/viper
- **Why:** Flexible config loading from environment, files, or flags
- **Usage:** Load `.env` or environment variables on startup
- **Integration:** Create internal/config wrapper
- **Features:**
  - Environment variable support
  - File-based config
  - Default values
  - Type conversion
  - Validation support

### ID Generation
**google/uuid** (Standard library compatible)
- **Package:** github.com/google/uuid
- **Usage:** Generate UUIDs for entities
- **Alternative:** Consider UUIDv7 for database-friendly ordering
- **Scope:** User IDs, correlation IDs, request tracing

## Testing Infrastructure

### Test Framework
**Go Standard Testing** (built-in)
- **Package:** testing
- **Why:** Built-in, no external dependency
- **Pattern:** `*_test.go` files in `_test` package
- **Assertions:** Manual assertions or popular assertion libraries

### Mocking
**uber-go/mock** (Interface-based mocking)
- **Package:** github.com/uber-go/mock/gomock
- **Tool:** `mockgen` (code generation)
- **Usage:** Generate mocks from interfaces via `//go:generate`
- **Benefits:** Type-safe mocks, compile-time verification
- **Pattern:** Auto-generated in `domain/mock/` directories

### Database Mocking
**DATA-DOG/go-sqlmock** - Database mock for testing
- **Package:** github.com/DATA-DOG/go-sqlmock
- **Usage:** Test repository layer without real database
- **Pattern:** Mock SQL queries and responses
- **Benefits:** Fast tests, no external database needed

### Test Assertions (Optional)
**testify/assert** (Popular assertion library)
- **Package:** github.com/stretchr/testify/assert
- **Usage:** Simplified assertions in tests
- **Status:** Optional, nice-to-have (not required)

## Code Quality & Linting

### Linting
**golangci-lint v2** - Comprehensive Go linter
- **Tool:** golangci-lint
- **Configuration:** `.golangci.yml` (pre-configured)
- **Strictness:** ALL checks enabled
- **Coverage:** Lint all code before commit
- **CI/CD:** Required to pass in pipeline

### Code Coverage
**go test -cover** (Built-in coverage)
- **Target:** > 80% coverage across all layers
- **Tool:** `go test -cover ./...`
- **Report:** Optional coverage reports (html, xml)

### Performance Profiling
**pprof** (Built-in Go profiling)
- **Usage:** Profile endpoints in main.go
- **Optional:** For performance optimization phase

## API Documentation

### Swagger/OpenAPI
**swaggo/swag** (Swagger code generation)
- **Package:** github.com/swaggo/swag/cmd/swag
- **Tool:** `swag init` generates docs from annotations
- **Usage:** Add annotations to handler functions
- **Output:** Swagger UI available at `/swagger/`
- **Benefits:** Auto-generated API documentation

**Alternative:** Skip swag, use raw OpenAPI YAML
- **Status:** Lower priority, can add later

## Container & Deployment

### Docker
**Dockerfile** (Multi-stage production image)
- **Base:** Official Go image for builder
- **Final:** Minimal alpine or distroless image
- **Size:** < 20MB final image (typical)
- **Features:**
  - Build stage optimization
  - Security (non-root user)
  - Health check support
  - Signal handling

### Docker Compose (Optional)
**docker-compose.yml** (For local development)
- **Services:** API + Database (MariaDB/PostgreSQL)
- **Volumes:** Database persistence
- **Networks:** Service communication
- **Status:** Optional, helpful for development

## Development Tools

### Make
**Makefile** (Task automation)
- **Commands:** build, run, test, lint, generate, docker-build, migrate-up, migrate-down
- **Purpose:** Consistent command interface
- **Benefits:** Self-documenting, consistent across team

### Environment Files
**.env.example** (Template for environment variables)
- **Usage:** Copy to .env for local development
- **Contents:** Database URL, server port, log level, JWT secret, etc.
- **Security:** Never commit .env with real secrets

## Project Template Dependencies (Initial go.mod)

```go
module github.com/zercle/template-go-echo

go 1.25

require (
    github.com/labstack/echo/v4 v4.11.0+
    github.com/golang-jwt/jwt/v5 v5.0.0+
    github.com/samber/do/v2 v2.0.0+
    github.com/samber/lo v1.44.0+
    github.com/spf13/viper v1.18.0+
    github.com/google/uuid v1.5.0+
)

require (
    // For testing
    github.com/DATA-DOG/go-sqlmock v1.5.0+
    github.com/stretchr/testify v1.8.0+
    github.com/golang/mock v1.6.0+
)

require (
    // For database (add based on choice)
    // Option 1: MySQL/MariaDB
    github.com/go-sql-driver/mysql v1.7.0+
    // Option 2: PostgreSQL
    github.com/lib/pq v1.10.0+
    // Option 3: MongoDB (via FerretDB)
    go.mongodb.org/mongo-driver v1.13.0+
)
```

**Note:** Exact versions will be locked based on dependency compatibility

## Indirect Dependencies

These are typically pulled in automatically:

### Echo Dependencies
- github.com/valyala/fasthttp (HTTP internals)
- github.com/mattn/go-colorable (Terminal colors)

### JWT Dependencies
- github.com/golang-jwt/jwt/v5 (JWT handling)

### Viper Dependencies
- github.com/spf13/pflag (Flag parsing)
- github.com/mitchellh/mapstructure (Config unmarshaling)

### samber Dependencies
- None major (minimal dependencies)

## Optional Advanced Dependencies (For Future)

### Event Processing
- **message-db** (Event sourcing)
- **NATS** (Pub/Sub messaging)

### Advanced Logging
- **logrus** (Alternative structured logging)
- **zerolog** (High-performance structured logging)

### Metrics & Monitoring
- **prometheus/client_golang** (Metrics collection)
- **opentelemetry-go** (Distributed tracing)

### Advanced Authentication
- **zitadel/oidc** (OIDC/OpenID Connect)
- **ory/hydra** (OAuth2/OIDC server)

### Caching
- **redis/go-redis** (Redis client)
- **valkey-io/valkey-go** (Valkey client)

## Dependency Update Strategy

**Philosophy:** Conservative, tested updates
1. **Critical Security:** Apply immediately
2. **Minor/Patch:** Regular intervals (monthly)
3. **Major:** Evaluate and test carefully
4. **Breaking Changes:** Only upgrade if necessary

**Process:**
```bash
go get -u ./...              # Check updates
go mod tidy                  # Cleanup
go test ./...               # Verify tests pass
golangci-lint run          # Verify linting
```

## Environment Configuration

### Local Development (.env.example)
```env
# Server
PORT=8080
LOG_LEVEL=debug
ENV=development

# Database (MariaDB example)
DB_DRIVER=mysql
DB_DSN=user:password@tcp(localhost:3306)/template_go_echo

# JWT
JWT_SECRET=your-secret-key-here-change-in-production

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Production Environment Variables
```env
# Server
PORT=8080 (via container)
LOG_LEVEL=info
ENV=production

# Database
DB_DRIVER=mysql
DB_DSN=from-secrets-manager

# JWT
JWT_SECRET=from-secrets-manager

# CORS
CORS_ALLOWED_ORIGINS=https://api.example.com
```

## Build & Deployment

### Local Build
```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

### Docker Build
```bash
docker build -t template-go-echo:latest .
docker run -p 8080:8080 -e PORT=8080 template-go-echo:latest
```

### Development Dependencies Installation
Required tools to install (one-time):
```bash
# Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/google/go-bindata/go-bindata@latest

# Code generation
go install github.com/uber-go/mock/cmd/mockgen@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Go Best Practices Applied

1. **Interface Segregation:** Small, focused interfaces
2. **Dependency Injection:** Constructor injection pattern
3. **Error Handling:** Explicit error returns, structured errors
4. **Concurrency:** Goroutines with context awareness
5. **Logging:** Structured logs with context propagation
6. **Testing:** High coverage with table-driven tests
7. **Documentation:** Comments on exported symbols
8. **Package Organization:** Logical grouping by concern
9. **Code Style:** gofmt + golangci-lint compliant
10. **Performance:** Minimal allocations, connection pooling

## Compatibility Matrix

| Component | Min Version | Recommended | Max Version |
|-----------|------------|-------------|------------|
| Go | 1.25 | 1.25.3 | Latest |
| Echo | v4.0 | v4.11+ | v5.x? |
| sqlc | v1.18 | Latest | v2.x? |
| golangci-lint | v1.54 | Latest | Latest |
| Database Drivers | Native | Latest | Latest |

## Security Considerations

### Built-In Protections
- **Type Safety:** Go's type system prevents many errors
- **Memory Safety:** Go manages memory automatically
- **SQL Safety:** sqlc prevents SQL injection
- **Concurrency Safety:** Go's race detector available

### Additional Security Measures
- **JWT Validation:** Every protected endpoint validates token
- **Input Validation:** All handler inputs validated
- **CORS Configuration:** Restrict cross-origin requests
- **Security Headers:** HSTS, X-Frame-Options, etc.
- **Rate Limiting:** Prevent brute force attacks
- **Secrets Management:** Use environment variables (never in code)

## Performance Characteristics

### Expected Performance (Single Instance)
- **Throughput:** 10,000+ req/sec (simple endpoint)
- **Latency:** < 50ms P95 (typical query)
- **Memory:** 50-200MB base (varies by goroutines)
- **CPU:** Low-moderate on standard workload

### Optimization Opportunities (Future)
- Connection pooling (sqlc generates optimized queries)
- Caching layer (Redis/Valkey)
- Async processing (background workers)
- Database indexing (application-specific)

## Resource Requirements

### Minimum (Development)
- **CPU:** 2 cores
- **Memory:** 1GB
- **Storage:** 2GB
- **Go:** 1.25+

### Recommended (Production)
- **CPU:** 4+ cores
- **Memory:** 2-4GB
- **Storage:** 10GB+ (depends on data)
- **Database:** Dedicated instance

### Container Image
- **Base:** golang:1.25 (build)
- **Runtime:** alpine:latest or distroless
- **Size:** ~10-20MB (final image)
