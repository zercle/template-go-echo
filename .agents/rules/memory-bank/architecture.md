# Architecture: template-go-echo

## Architectural Style
**Clean Architecture with Domain-Driven Design (DDD)**

The project implements layered architecture with explicit dependency inversion, separating concerns into horizontal layers while organizing code vertically by domain modules.

```
┌─────────────────────────────────────┐
│      HTTP Layer (Echo Server)       │
│    (Router, Middleware, Handlers)   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    Handler Layer (HTTP/REST)        │
│  (Requests, DTOs, Validation)       │
│  (Swagger Annotations)              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    Usecase Layer (Business Logic)   │
│  (Orchestration, DI, Business Rules)│
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Repository Layer (Data Access)    │
│  (sqlc Generated Code, Transactions)│
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Domain Layer (Business Entities)  │
│  (Interfaces, Entities, Contracts)  │
└─────────────────────────────────────┘
```

## Component Organization

### 1. **Domain Layer** (`internal/{domain}/domain`)
**Responsibility:** Core business logic and contracts

**Components:**
- **Entities:** Core business objects (e.g., User, Order)
- **Value Objects:** Immutable, business-meaningful values (e.g., Email, Money)
- **Repository Interfaces:** Contracts for data access
- **Error Types:** Domain-specific errors
- **Domain Services:** Pure business logic

**Key Rules:**
- No external dependencies (except Go stdlib)
- Language-agnostic business logic
- Define `interface{}` contracts for abstraction
- No knowledge of database or HTTP details

**Example Structure:**
```go
internal/user/domain/
├── user.go              # User entity
├── error.go             # Domain errors
├── repository.go        # Repository interface
└── service.go           # Domain services
```

### 2. **Repository Layer** (`internal/{domain}/repository`)
**Responsibility:** Data persistence abstraction

**Components:**
- **Repository Implementations:** Satisfy domain interfaces
- **sqlc Generated Code:** Type-safe SQL queries (`sqlc/` subdirectory)
- **Transaction Management:** Begin/Commit/Rollback
- **Query Building:** Construct complex queries
- **Error Translation:** Map DB errors to domain errors

**Key Rules:**
- Implement domain repository interfaces
- Return domain entities, not DTOs
- Handle transaction lifecycle
- Abstract database details from upper layers
- Use sqlc for type-safe generated code

**Example Structure:**
```go
internal/user/repository/
├── repository.go        # Repository implementation
├── query.go             # Custom queries
├── transaction.go       # Transaction handling
└── sqlc/               # Generated code
    ├── models.go
    ├── querier.go
    └── user.sql.go
```

### 3. **Usecase Layer** (`internal/{domain}/usecase`)
**Responsibility:** Business logic orchestration and application rules

**Components:**
- **Usecase Classes:** One per business operation (e.g., CreateUserUsecase)
- **DTOs:** Transfer objects for input/output
- **Dependency Injection:** Constructor injection via `samber/do`
- **Business Rules:** Applied before persistence
- **Logging:** Structured logging of operations

**Key Rules:**
- One usecase per business operation
- Depend on domain layer
- Constructor injection for dependencies
- Use domain entities, not DTOs
- Translate domain errors to meaningful responses

**Example Structure:**
```go
internal/user/usecase/
├── create.go            # CreateUser usecase
├── get.go               # GetUser usecase
├── update.go            # UpdateUser usecase
├── delete.go            # DeleteUser usecase
└── dto.go               # Data Transfer Objects
```

### 4. **Handler Layer** (`internal/{domain}/handler`)
**Responsibility:** HTTP request/response handling and validation

**Components:**
- **HTTP Handlers:** Echo handler functions
- **Request DTOs:** Parse and validate HTTP requests
- **Response DTOs:** Format HTTP responses
- **Swagger Annotations:** API documentation
- **Input Validation:** Parameter and body validation
- **Status Code Mapping:** HTTP status code selection

**Key Rules:**
- One handler per HTTP endpoint
- Validate all inputs
- Include Swagger annotations
- Return JSend-formatted responses
- Delegate to usecases
- Translate domain errors to HTTP responses

**Example Structure:**
```go
internal/user/handler/
├── handler.go           # User handler definition
├── create.go            # POST /users
├── get.go               # GET /users/{id}
├── update.go            # PUT /users/{id}
├── delete.go            # DELETE /users/{id}
└── dto.go               # Request/Response DTOs
```

### 5. **Middleware Layer** (`internal/middleware`)
**Responsibility:** Cross-cutting concerns for HTTP requests

**Components:**
- **Authentication:** JWT verification
- **Authorization:** Role/permission checks
- **Logging:** Request/response logging
- **CORS:** Cross-origin request handling
- **Rate Limiting:** Request throttling
- **Error Handling:** Global error handler
- **Request ID:** Correlation ID propagation

**Key Rules:**
- Chainable middleware pattern
- Context propagation
- Structured logging integration
- Graceful error handling

**Example Structure:**
```go
internal/middleware/
├── auth.go              # JWT authentication
├── cors.go              # CORS handling
├── logger.go            # Request logging
├── rate_limit.go        # Rate limiting
└── error.go             # Error handling
```

### 6. **Configuration Layer** (`internal/config`)
**Responsibility:** Environment and application configuration

**Components:**
- **Config Struct:** All application settings
- **Validation:** Config validation on startup
- **Environment Loading:** Load from .env or environment
- **Defaults:** Sensible defaults per environment
- **Secrets Management:** Safe secret handling

**Key Rules:**
- Load on startup, validate immediately
- Fail fast if invalid
- Use environment-specific defaults
- Structure by concern (database, server, logging, etc.)

**Example Structure:**
```go
internal/config/
├── config.go            # Main config struct
├── loader.go            # Load from environment
└── validator.go         # Validate config
```

### 7. **Shared Utilities** (`pkg/`)
**Responsibility:** Reusable utilities across domains

**Components:**
- **Response Formatting:** JSend response wrapper
- **Custom Types:** Domain-agnostic custom types
- **Error Utilities:** Error handling helpers
- **Logging Utilities:** Structured logging setup
- **Database Utilities:** Connection pooling, helpers

**Key Rules:**
- Zero domain-specific logic
- Reusable across services
- Well-documented and tested

**Example Structure:**
```go
pkg/
├── response/            # JSend response format
├── errors/              # Error handling
├── logger/              # Logging setup
├── database/            # DB utilities
└── validate/            # Validation helpers
```

## Dependency Flow (Dependency Inversion)

```
Handler → Usecase → Repository → Domain
  ↓         ↓           ↓          ↓
  └─────────┴───────────┴──────────┘
         ↑
    Config & Middleware
```

**Direction:** Dependencies flow inward (toward Domain layer)
- Handlers depend on Usecases, not vice versa
- Usecases depend on Domain, not vice versa
- All layers can depend on Config and Utilities

## Module Organization

### Domain Module Structure
Each domain (e.g., `user`, `order`, `product`) follows this structure:

```
internal/{domain}/
├── domain/              # Domain layer
│   ├── *.go
│   └── *_mock.go        # Generated mocks
├── handler/             # Handler layer
│   ├── *.go
│   └── *_test.go
├── usecase/             # Usecase layer
│   ├── *.go
│   └── *_test.go
├── repository/          # Repository layer
│   ├── *.go
│   ├── sqlc/           # Generated sqlc code
│   └── *_test.go
└── mock/               # Generated mocks
    └── *.go
```

### Root Level Structure
```
template-go-echo/
├── cmd/api/             # Application entry point
│   └── main.go
├── internal/            # Private application code
│   ├── {domain}/        # Domain modules
│   ├── middleware/      # Shared middleware
│   └── config/          # Configuration
├── pkg/                 # Shared utilities
├── sql/                 # Database files
│   ├── queries/         # sqlc queries
│   └── migrations/      # Migrations
├── docs/                # API & Architecture docs
├── .agents/             # Agent configuration
│   └── rules/
│       └── memory-bank/ # Project documentation
└── [Config files]
```

## Data Flow Example: Create User API

```
1. HTTP Request (POST /users)
           ↓
2. Handler (handler/create.go)
   - Parse JSON to CreateUserRequest DTO
   - Validate request
   - Call usecase.CreateUser()
           ↓
3. Usecase (usecase/create.go)
   - Apply business rules
   - Call repository.SaveUser(domain.User)
           ↓
4. Repository (repository/create.go)
   - Translate domain.User to DB model
   - Execute sqlc generated Insert
   - Handle transaction
   - Translate DB model back to domain.User
           ↓
5. Domain Layer
   - User entity encapsulates business logic
   - Repository interface defines contract
           ↓
6. Response
   - Handler translates domain.User to response DTO
   - JSend wrapper around response
   - HTTP 201 Created
```

## Key Architectural Decisions (ADRs)

### ADR-001: Clean Architecture Layers
**Decision:** Use 5-layer architecture (Handler → Usecase → Repository → Domain)
**Rationale:** Clear separation of concerns, testability, independence from frameworks
**Trade-offs:** Extra files/code vs. maintainability and flexibility

### ADR-002: Domain-Driven Design Organization
**Decision:** Organize code vertically by domain (not by layer)
**Rationale:** Each domain is self-contained, easier to scale to microservices
**Trade-offs:** More duplication of layer structure vs. better organization

### ADR-003: Dependency Injection
**Decision:** Use `samber/do v2` for DI container
**Rationale:** Type-safe, performant, minimal boilerplate
**Trade-offs:** Manual DI wiring vs. reduced magic/reflection

### ADR-004: SQL Code Generation
**Decision:** Use `sqlc` for type-safe SQL query generation
**Rationale:** Type safety at compile time, performance, no ORM overhead
**Trade-offs:** Manual SQL writing vs. type safety

### ADR-005: JSend Response Format
**Decision:** Use JSend for all JSON responses
**Rationale:** Standardized, explicit status field, consistent error format
**Trade-offs:** Not traditional REST vs. clarity and consistency

### ADR-006: Interface Mocking
**Decision:** Generate mocks via `//go:generate mockgen`
**Rationale:** Consistent mocks, compile-time verification
**Trade-offs:** Setup overhead vs. reliable mocking

## Scalability Considerations

### Horizontal Scaling
- **Statelessness:** No session storage in service, use external session store if needed
- **Health Checks:** `/health` endpoint for load balancer checks
- **Graceful Shutdown:** Allow in-flight requests to complete
- **Database Connection Pooling:** Configured per environment

### Vertical Scaling
- **Async Operations:** Background jobs for long-running tasks
- **Caching:** Redis/Valkey for frequently accessed data
- **Rate Limiting:** Protect APIs from abuse
- **Request Timeouts:** Prevent hung requests

### Evolution to Microservices
- **Domain Modules:** Each domain can become independent service
- **APIs:** Domain modules expose clear interfaces
- **Database per Domain:** Each domain has own database
- **Event-Driven:** Events between services (future evolution)

## Security Architecture

- **Authentication:** JWT tokens in Authorization header
- **Validation:** Input validation at handler layer
- **Authorization:** RBAC middleware for protected endpoints
- **Error Messages:** Sanitized error responses (no internal details)
- **Headers:** Security headers (CORS, CSP, etc.)
- **SQL Injection:** Protection via sqlc type-safe queries
- **Rate Limiting:** Per-endpoint rate limits

## Testing Strategy

- **Unit Tests:** Domain layer (pure functions)
- **Integration Tests:** Usecase + Repository with sqlmock
- **Handler Tests:** Mock usecase, test HTTP behavior
- **Mock Generation:** Automatic via `go generate`
- **Coverage Target:** > 80% across all layers
