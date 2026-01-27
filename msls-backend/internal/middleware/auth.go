// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/services/auth"
)

// Context keys for authentication.
const (
	// UserKey is the context key for the authenticated user.
	UserKey = "user"
	// UserIDKey is the context key for the user ID.
	UserIDKey = "user_id"
	// TenantIDFromTokenKey is the context key for the tenant ID from the JWT token.
	TenantIDFromTokenKey = "tenant_id_from_token"
	// PermissionsKey is the context key for the user's permissions.
	PermissionsKey = "permissions"
	// ClaimsKey is the context key for the JWT claims.
	ClaimsKey = "claims"
)

// AuthConfig holds configuration for auth middleware.
type AuthConfig struct {
	JWTService *auth.JWTService
	// SkipPaths lists paths that don't require authentication.
	SkipPaths []string
}

// Auth returns a middleware that validates JWT tokens.
func Auth(config AuthConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// Skip authentication for excluded paths
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apperrors.Abort(c, apperrors.Unauthorized("Authorization header is required"))
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			apperrors.Abort(c, apperrors.Unauthorized("Invalid authorization header format"))
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			apperrors.Abort(c, apperrors.Unauthorized("Token is required"))
			return
		}

		// Validate token
		claims, err := config.JWTService.ValidateAccessToken(tokenString)
		if err != nil {
			switch err {
			case auth.ErrExpiredToken:
				apperrors.Abort(c, apperrors.Unauthorized("Token has expired"))
			case auth.ErrTokenNotYetValid:
				apperrors.Abort(c, apperrors.Unauthorized("Token is not yet valid"))
			default:
				apperrors.Abort(c, apperrors.Unauthorized("Invalid token"))
			}
			return
		}

		// Set user information in context
		c.Set(ClaimsKey, claims)
		c.Set(UserIDKey, claims.UserID)
		c.Set(TenantIDFromTokenKey, claims.TenantID)
		c.Set(PermissionsKey, claims.Permissions)

		c.Next()
	}
}

// AuthRequired returns a middleware that requires JWT authentication.
func AuthRequired(jwtService *auth.JWTService) gin.HandlerFunc {
	return Auth(AuthConfig{
		JWTService: jwtService,
		SkipPaths:  []string{},
	})
}

// AuthOptional returns a middleware that optionally validates JWT tokens.
// If a token is present, it will be validated. If not, the request proceeds without authentication.
func AuthOptional(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.Next()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			// Token is invalid, but continue without authentication
			c.Next()
			return
		}

		// Set user information in context
		c.Set(ClaimsKey, claims)
		c.Set(UserIDKey, claims.UserID)
		c.Set(TenantIDFromTokenKey, claims.TenantID)
		c.Set(PermissionsKey, claims.Permissions)

		c.Next()
	}
}

// PermissionRequired returns a middleware that checks if the user has the required permissions.
// The user must have at least one of the specified permissions.
func PermissionRequired(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user permissions from context
		permsValue, exists := c.Get(PermissionsKey)
		if !exists {
			apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
			return
		}

		userPerms, ok := permsValue.([]string)
		if !ok {
			apperrors.Abort(c, apperrors.InternalError("Invalid permissions in context"))
			return
		}

		// Check if user has at least one of the required permissions
		hasPermission := false
		for _, required := range permissions {
			for _, userPerm := range userPerms {
				if userPerm == required {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			apperrors.Abort(c, &apperrors.AppError{
				Type:   apperrors.TypeForbidden,
				Title:  "Forbidden",
				Status: http.StatusForbidden,
				Detail: "You don't have permission to access this resource",
			})
			return
		}

		c.Next()
	}
}

// PermissionRequiredAll returns a middleware that checks if the user has all required permissions.
func PermissionRequiredAll(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user permissions from context
		permsValue, exists := c.Get(PermissionsKey)
		if !exists {
			apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
			return
		}

		userPerms, ok := permsValue.([]string)
		if !ok {
			apperrors.Abort(c, apperrors.InternalError("Invalid permissions in context"))
			return
		}

		// Create a map for faster lookup
		userPermMap := make(map[string]bool)
		for _, perm := range userPerms {
			userPermMap[perm] = true
		}

		// Check if user has all required permissions
		for _, required := range permissions {
			if !userPermMap[required] {
				apperrors.Abort(c, &apperrors.AppError{
					Type:   apperrors.TypeForbidden,
					Title:  "Forbidden",
					Status: http.StatusForbidden,
					Detail: "You don't have permission to access this resource",
				})
				return
			}
		}

		c.Next()
	}
}

// GetCurrentUserID retrieves the current user's ID from the Gin context.
func GetCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	if userID, exists := c.Get(UserIDKey); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}

// GetCurrentTenantID retrieves the tenant ID from the JWT token in the Gin context.
func GetCurrentTenantID(c *gin.Context) (uuid.UUID, bool) {
	if tenantID, exists := c.Get(TenantIDFromTokenKey); exists {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}

// GetCurrentClaims retrieves the JWT claims from the Gin context.
func GetCurrentClaims(c *gin.Context) (*auth.Claims, bool) {
	if claims, exists := c.Get(ClaimsKey); exists {
		if cl, ok := claims.(*auth.Claims); ok {
			return cl, true
		}
	}
	return nil, false
}

// GetCurrentPermissions retrieves the user's permissions from the Gin context.
func GetCurrentPermissions(c *gin.Context) ([]string, bool) {
	if perms, exists := c.Get(PermissionsKey); exists {
		if p, ok := perms.([]string); ok {
			return p, true
		}
	}
	return nil, false
}

// HasPermission checks if the current user has a specific permission.
func HasPermission(c *gin.Context, permission string) bool {
	perms, ok := GetCurrentPermissions(c)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}
