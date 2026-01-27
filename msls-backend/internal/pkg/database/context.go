// Package database provides context propagation helpers for database operations.
package database

import (
	"context"

	"gorm.io/gorm"
)

// Context keys for database-related values.
const (
	// tenantIDKey is the context key for tenant ID.
	tenantIDKey contextKey = "tenant_id"

	// requestIDKey is the context key for request ID.
	requestIDKey contextKey = "request_id"

	// userIDKey is the context key for user ID.
	userIDKey contextKey = "user_id"
)

// TenantID retrieves the tenant ID from the context.
func TenantID(ctx context.Context) string {
	if id, ok := ctx.Value(tenantIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextWithTenantID adds a tenant ID to the context.
func ContextWithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// RequestID retrieves the request ID from the context.
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextWithRequestID adds a request ID to the context.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// UserID retrieves the user ID from the context.
func UserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextWithUserID adds a user ID to the context.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// DBWithTenantContext returns a DB connection with tenant context set.
// This sets the app.tenant_id session variable in PostgreSQL for RLS.
func DBWithTenantContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	tenantID := TenantID(ctx)
	if tenantID == "" {
		return db.WithContext(ctx)
	}

	return db.WithContext(ctx).Session(&gorm.Session{
		NewDB: true,
	}).Exec("SET LOCAL app.tenant_id = ?", tenantID)
}

// SetTenantContext sets the tenant context on an existing DB connection.
// This should be called at the beginning of each request that requires tenant isolation.
func SetTenantContext(db *gorm.DB, tenantID string) error {
	if tenantID == "" {
		return nil
	}
	return db.Exec("SET LOCAL app.tenant_id = ?", tenantID).Error
}

// ClearTenantContext clears the tenant context from the DB session.
func ClearTenantContext(db *gorm.DB) error {
	return db.Exec("RESET app.tenant_id").Error
}

// ScopedDB returns a new DB session with all context values propagated.
// This is useful for creating a scoped DB connection for a specific operation.
func ScopedDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	session := db.WithContext(ctx)

	// Set tenant context if available
	if tenantID := TenantID(ctx); tenantID != "" {
		// Use Session to avoid modifying the original DB
		session = session.Session(&gorm.Session{
			NewDB: true,
		})
		// Set the tenant context (ignoring error as it's a session setup)
		session.Exec("SET LOCAL app.tenant_id = ?", tenantID)
	}

	return session
}

// WithContext creates a context with all database-related values.
func WithContext(ctx context.Context, tenantID, requestID, userID string) context.Context {
	if tenantID != "" {
		ctx = ContextWithTenantID(ctx, tenantID)
	}
	if requestID != "" {
		ctx = ContextWithRequestID(ctx, requestID)
	}
	if userID != "" {
		ctx = ContextWithUserID(ctx, userID)
	}
	return ctx
}
