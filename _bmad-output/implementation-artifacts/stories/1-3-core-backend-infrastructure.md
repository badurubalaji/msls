# Story 1.3: Core Backend Infrastructure

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As a **developer**,
I want **reusable backend infrastructure components**,
So that **all modules follow consistent patterns for database, logging, and error handling**.

## Acceptance Criteria

**Given** a developer needs to create a new module
**When** they use the core infrastructure
**Then** the following components are available:
- Database connection pool with context propagation
- Transaction helper with automatic rollback on panic
- Tenant context middleware that sets `app.tenant_id` session variable
- Request ID middleware for tracing
- Structured JSON logging with correlation ID
- Standard error types (NotFound, Validation, Unauthorized, Forbidden)
- RFC 7807 error response formatter
**And** API rate limiting middleware exists (configurable per endpoint)
**And** Request/response logging middleware captures method, path, status, duration
**And** Panic recovery middleware prevents server crashes

## Technical Requirements

### Package Structure
```
msls-backend/internal/
├── middleware/
│   ├── tenant.go          # Tenant context middleware
│   ├── request_id.go      # Request ID middleware
│   ├── logging.go         # Request logging middleware
│   ├── recovery.go        # Panic recovery middleware
│   ├── rate_limit.go      # Rate limiting middleware
│   └── cors.go            # CORS middleware
├── pkg/
│   ├── database/
│   │   ├── connection.go  # Connection pool setup
│   │   ├── transaction.go # Transaction helper
│   │   └── context.go     # Context propagation
│   ├── errors/
│   │   ├── types.go       # Error types
│   │   └── handler.go     # Error handler middleware
│   └── validator/
│       └── validator.go   # Custom validators
```

### Middleware Implementations

#### Tenant Middleware
```go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            // Return error
        }
        // Set in context
        c.Set("tenant_id", tenantID)
        // Set PostgreSQL session variable
        db.Exec("SET app.tenant_id = ?", tenantID)
        c.Next()
    }
}
```

#### Request ID Middleware
- Generate UUID v4 for each request
- Set in context and response header
- Propagate to all log entries

#### Rate Limiting
- Use token bucket algorithm
- Configurable per route/endpoint
- Return 429 Too Many Requests with Retry-After header

### Error Types
```go
type AppError struct {
    Type       string `json:"type"`
    Title      string `json:"title"`
    Status     int    `json:"status"`
    Detail     string `json:"detail"`
    Instance   string `json:"instance"`
    Extensions map[string]interface{} `json:"extensions,omitempty"`
}
```

### Transaction Helper
```go
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
    tx := db.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

## Tasks

1. [ ] Create database connection pool with config
2. [ ] Create transaction helper with rollback
3. [ ] Create tenant context middleware
4. [ ] Create request ID middleware
5. [ ] Create structured logging middleware
6. [ ] Create panic recovery middleware
7. [ ] Create rate limiting middleware
8. [ ] Create CORS middleware
9. [ ] Create standard error types
10. [ ] Create RFC 7807 error formatter
11. [ ] Create error handler middleware
12. [ ] Create custom validators (phone, email, etc.)
13. [ ] Update main.go to use all middleware
14. [ ] Write unit tests for middleware
15. [ ] Write integration tests for database helpers

## Definition of Done

- [ ] All middleware components implemented
- [ ] Database connection pool properly configured
- [ ] Transaction helper tested with rollback scenarios
- [ ] Error types follow RFC 7807
- [ ] Unit tests pass with >80% coverage
- [ ] Integration tests verify middleware chain
