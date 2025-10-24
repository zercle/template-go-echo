# Architecture: Go Echo Modular Monolith Template

## System Architecture Overview

### Architectural Pattern
**Clean Architecture with Domain-Driven Design (DDD)**
- Modular monolith structure designed for future microservices transition
- Clear separation of concerns across all layers
- Dependency inversion with interfaces and dependency injection
- Domain-centric organization with bounded contexts

### Core Architectural Principles
1. **Modularity**: Each domain is self-contained with its own layers
2. **Testability**: All components are mockable and independently testable
3. **Scalability**: Architecture supports both vertical and horizontal scaling
4. **Maintainability**: Clear boundaries and dependencies enable long-term maintenance
5. **Performance**: Optimized for high-throughput, low-latency operations

## Layer Architecture

### Handler Layer (Presentation)
**Location**: `internal/{domain}/handler/`
**Responsibilities**:
- HTTP request/response handling
- Input validation and DTOs
- Swagger/OpenAPI documentation generation
- Response formatting (JSend standard)
- Authentication and authorization checks

**Key Components**:
- `*Handler` structs implementing HTTP endpoints
- Request/Response DTOs
- Swagger annotations for API documentation
- Middleware integration

### Usecase Layer (Business Logic)
**Location**: `internal/{domain}/usecase/`
**Responsibilities**:
- Core business logic implementation
- Domain operations orchestration
- Transaction coordination
- Business rule enforcement
- Domain event handling

**Key Components**:
- `*Usecase` interfaces and implementations
- Business logic methods
- Input validation at business level
- Error handling and domain-specific errors

### Repository Layer (Data Access)
**Location**: `internal/{domain}/repository/`
**Responsibilities**:
- Database operations via sqlc generated code
- Transaction management
- Data mapping between domain entities and database
- Caching strategies implementation

**Key Components**:
- `*Repository` interfaces and implementations
- Database transaction handling
- sqlc generated code integration
- Query optimization

### Domain Layer (Core)
**Location**: `internal/{domain}/domain/`
**Responsibilities**:
- Domain entities and value objects
- Business contracts and interfaces
- Domain-specific types and constants
- Repository and usecase interface definitions

**Key Components**:
- Entity structs and methods
- Value objects
- Domain interfaces
- Business constants and enums

## Domain Structure

### Domain Organization
```
internal/
├── {domain}/                 # e.g., user, order, auth, product
│   ├── domain/              # Core domain logic
│   │   ├── entities.go      # Domain entities
│   │   ├── interfaces.go    # Repository and Usecase interfaces
│   │   ├── constants.go     # Domain constants
│   │   └── errors.go        # Domain-specific errors
│   ├── handler/             # HTTP handlers
│   │   ├── handlers.go      # HTTP endpoint implementations
│   │   ├── dto.go          # Data Transfer Objects
│   │   └── validators.go    # Request validation
│   ├── usecase/            # Business logic
│   │   ├── usecase.go      # Usecase interface
│   │   ├── usecase_impl.go # Usecase implementation
│   │   └── input.go        # Usecase input structures
│   ├── repository/         # Data access
│   │   ├── repository.go   # Repository interface
│   │   ├── repository_impl.go # Repository implementation
│   │   └── transactions.go # Transaction handling
│   ├── sqlc/              # Generated database code
│   │   ├── db.go          # Generated database types
│   │   ├── queries.go     # Generated query functions
│   │   └── models.go      # Generated models
│   └── test/              # Test files
│       ├── handler_test.go
│       ├── usecase_test.go
│       └── repository_test.go
```

## Component Relationships

### Data Flow Architecture
```
HTTP Request → Handler → Usecase → Repository → Database
    ↑            ↓         ↓          ↓         ↓
Response ← JSON ← DTO ← Domain ← SQL Result ← DB
```

### Dependency Flow
```
Handler depends on → Usecase Interface
Usecase depends on → Repository Interface
Repository depends on → Database/sqlc
Domain defines → All interfaces and entities
```

## Cross-Cutting Concerns

### Infrastructure Layer
**Location**: `internal/infrastructure/`
- Database connection management
- External service integrations
- Shared utilities and helpers

### Middleware Layer
**Location**: `internal/middleware/`
- Authentication/authorization
- Request logging and tracing
- Rate limiting and security
- CORS and security headers
- Error handling middleware

### Configuration Layer
**Location**: `internal/config/`
- Environment variable management
- Configuration validation
- Service discovery settings

### Utility Layer
**Location**: `pkg/`
- Response formatting utilities
- Custom error types
- Common validation functions
- Shared types and constants

## Database Architecture

### SQL Organization
```
sql/
├── queries/              # SQL queries for sqlc
│   ├── {domain}/
│   │   ├── create_*.sql
│   │   ├── get_*.sql
│   │   ├── update_*.sql
│   │   └── delete_*.sql
└── migrations/           # Database migrations
    ├── 001_initial_schema.up.sql
    ├── 001_initial_schema.down.sql
    └── ...
```

### Migration Strategy
- Version-controlled schema changes using golang-migrate
- Forward and backward migration support
- Database-agnostic migration scripts
- Production-safe migration workflows

## Security Architecture

### Authentication Flow
```
Client → Login Request → Handler → Usecase → Repository → User Check
    ↓                ↓         ↓          ↓         ↓
JWT Token ← Response ← DTO ← Domain ← SQL Result ← DB
```

### Authorization Layers
1. **Middleware Level**: Route-based access control
2. **Handler Level**: Resource-level permissions
3. **Usecase Level**: Business rule enforcement
4. **Repository Level**: Data access restrictions

## Deployment Architecture

### Container Strategy
- Single binary deployment
- Multi-stage Docker builds
- Health check endpoints
- Graceful shutdown handling

### Configuration Management
- Environment-based configuration
- Configuration validation at startup
- Secret management integration
- Runtime configuration updates

## Performance Architecture

### Caching Strategy
- Application-level caching for frequently accessed data
- Database query result caching
- Session and token caching
- Static asset optimization

### Connection Management
- Database connection pooling
- HTTP client connection reuse
- Keep-alive connections optimization
- Resource cleanup and management

## Testing Architecture

### Test Organization
```
internal/{domain}/test/
├── handler_test.go     # HTTP endpoint tests
├── usecase_test.go     # Business logic tests
├── repository_test.go  # Data access tests
└── integration_test.go # End-to-end tests
```

### Mock Strategy
- Interface-based mocking with uber-go/mock
- Database mocking with go-sqlmock
- HTTP client mocking for external services
- Test fixture management

## Architecture Decision Records (ADRs)

### ADR-001: Clean Architecture Adoption
**Decision**: Adopt Clean Architecture with DDD principles
**Rationale**: Better testability, maintainability, and scalability
**Consequences**: Additional initial complexity, long-term benefits

### ADR-002: SQL First Approach
**Decision**: Use SQL with sqlc for type-safe database access
**Rationale**: Type safety, better performance, explicit queries
**Consequences**: Learning curve for sqlc, migration management

### ADR-003: Modular Monolith Design
**Decision**: Start with modular monolith, prepare for microservices
**Rationale**: Simpler deployment, easier debugging, gradual evolution
**Consequences**: Shared database, need for strict domain boundaries

### ADR-004: Dependency Injection Container
**Decision**: Use samber/do for dependency injection
**Rationale**: Cleaner code, better testability, explicit dependencies
**Consequences**: Runtime dependency resolution, learning curve