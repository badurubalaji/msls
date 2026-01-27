// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database"
	apperrors "msls-backend/internal/pkg/errors"
)

const (
	// TenantIDHeader is the HTTP header for tenant identification.
	TenantIDHeader = "X-Tenant-ID"

	// TenantIDKey is the context key for storing tenant ID.
	TenantIDKey = "tenant_id"
)

// TenantConfig holds configuration for tenant middleware.
type TenantConfig struct {
	// Required determines if tenant ID is mandatory for all requests.
	Required bool

	// SkipPaths lists paths that don't require tenant ID (e.g., health checks).
	SkipPaths []string

	// DB is the database connection for setting PostgreSQL session variable.
	DB *gorm.DB
}

// DefaultTenantConfig returns the default tenant middleware configuration.
func DefaultTenantConfig() TenantConfig {
	return TenantConfig{
		Required:  true,
		SkipPaths: []string{"/health", "/ready", "/metrics"},
		DB:        nil,
	}
}

// Tenant returns a middleware that extracts tenant ID from the X-Tenant-ID header.
// It sets the tenant ID in the Gin context and optionally in the PostgreSQL session.
func Tenant(config TenantConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// Skip tenant check for excluded paths
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		tenantID := c.GetHeader(TenantIDHeader)

		// Check if tenant ID is required
		if tenantID == "" && config.Required {
			apperrors.Abort(c, &apperrors.AppError{
				Type:   apperrors.TypeBadRequest,
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "X-Tenant-ID header is required",
			})
			return
		}

		if tenantID != "" {
			// Validate tenant ID format (UUID)
			if !isValidUUID(tenantID) {
				apperrors.Abort(c, &apperrors.AppError{
					Type:   apperrors.TypeBadRequest,
					Title:  "Bad Request",
					Status: http.StatusBadRequest,
					Detail: "Invalid tenant ID format",
				})
				return
			}

			// Set tenant ID in Gin context
			c.Set(TenantIDKey, tenantID)

			// Update request context with tenant ID
			ctx := database.ContextWithTenantID(c.Request.Context(), tenantID)
			c.Request = c.Request.WithContext(ctx)

			// Set PostgreSQL session variable if DB is configured
			if config.DB != nil {
				if err := setDBTenantContext(config.DB, tenantID); err != nil {
					apperrors.Abort(c, apperrors.InternalError("Failed to set tenant context"))
					return
				}
			}
		}

		c.Next()
	}
}

// TenantRequired returns a middleware that requires tenant ID for all requests.
func TenantRequired() gin.HandlerFunc {
	return Tenant(DefaultTenantConfig())
}

// TenantOptional returns a middleware that makes tenant ID optional.
func TenantOptional() gin.HandlerFunc {
	config := DefaultTenantConfig()
	config.Required = false
	return Tenant(config)
}

// TenantWithDB returns a middleware that sets tenant context in PostgreSQL.
func TenantWithDB(db *gorm.DB) gin.HandlerFunc {
	config := DefaultTenantConfig()
	config.DB = db
	return Tenant(config)
}

// GetTenantID retrieves the tenant ID from the Gin context.
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(TenantIDKey); exists {
		if id, ok := tenantID.(string); ok {
			return id
		}
	}
	return ""
}

// setDBTenantContext sets the app.tenant_id session variable in PostgreSQL.
func setDBTenantContext(db *gorm.DB, tenantID string) error {
	return db.Exec("SET LOCAL app.tenant_id = ?", tenantID).Error
}

// isValidUUID validates if a string is a valid UUID format.
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
