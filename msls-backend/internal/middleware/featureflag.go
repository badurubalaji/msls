// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/services/featureflag"
)

// Context keys for feature flags.
const (
	// FeatureFlagsKey is the context key for the feature flags map.
	FeatureFlagsKey = "feature_flags"
	// FeatureFlagServiceKey is the context key for the feature flag service.
	FeatureFlagServiceKey = "feature_flag_service"
)

// FeatureFlagConfig holds configuration for feature flag middleware.
type FeatureFlagConfig struct {
	// Service is the feature flag service instance.
	Service *featureflag.Service
	// LoadOnRequest loads all flags into context on each request.
	LoadOnRequest bool
}

// FeatureFlag returns a middleware that loads feature flags into the request context.
// This middleware should be placed after authentication middleware to have access to user/tenant IDs.
func FeatureFlag(config FeatureFlagConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store service in context for handler access
		c.Set(FeatureFlagServiceKey, config.Service)

		if config.LoadOnRequest {
			// Get user and tenant IDs from context (set by auth middleware)
			userID, _ := GetCurrentUserID(c)
			tenantID, _ := GetCurrentTenantID(c)

			// Load all flags for this context
			flags := config.Service.GetAllFlagsForContext(c.Request.Context(), tenantID, userID)

			// Convert to a map for easy access
			flagMap := make(map[string]bool)
			for _, flag := range flags {
				flagMap[flag.Key] = flag.Enabled
			}
			c.Set(FeatureFlagsKey, flagMap)
		}

		c.Next()
	}
}

// FeatureFlagDefault returns a middleware with default configuration.
func FeatureFlagDefault(service *featureflag.Service) gin.HandlerFunc {
	return FeatureFlag(FeatureFlagConfig{
		Service:       service,
		LoadOnRequest: true,
	})
}

// FeatureFlagRequired returns a middleware that requires a specific feature flag to be enabled.
// Returns 403 Forbidden if the flag is disabled.
func FeatureFlagRequired(flagKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsFeatureEnabled(c, flagKey) {
			apperrors.Abort(c, apperrors.Forbidden("This feature is not available"))
			return
		}
		c.Next()
	}
}

// FeatureFlagAnyRequired returns a middleware that requires at least one of the specified feature flags to be enabled.
func FeatureFlagAnyRequired(flagKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, key := range flagKeys {
			if IsFeatureEnabled(c, key) {
				c.Next()
				return
			}
		}
		apperrors.Abort(c, apperrors.Forbidden("This feature is not available"))
	}
}

// FeatureFlagAllRequired returns a middleware that requires all specified feature flags to be enabled.
func FeatureFlagAllRequired(flagKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, key := range flagKeys {
			if !IsFeatureEnabled(c, key) {
				apperrors.Abort(c, apperrors.Forbidden("This feature is not available"))
				return
			}
		}
		c.Next()
	}
}

// IsFeatureEnabled checks if a feature flag is enabled for the current request.
// This function uses the cached flags from context first, then falls back to the service.
func IsFeatureEnabled(c *gin.Context, flagKey string) bool {
	// First, try to get from cached flags in context
	if flags, exists := c.Get(FeatureFlagsKey); exists {
		if flagMap, ok := flags.(map[string]bool); ok {
			if enabled, ok := flagMap[flagKey]; ok {
				return enabled
			}
		}
	}

	// Fall back to service if available
	if svc, exists := c.Get(FeatureFlagServiceKey); exists {
		if service, ok := svc.(*featureflag.Service); ok {
			userID, _ := GetCurrentUserID(c)
			tenantID, _ := GetCurrentTenantID(c)
			return service.IsEnabled(c.Request.Context(), flagKey, tenantID, userID)
		}
	}

	return false
}

// GetFeatureFlagService retrieves the feature flag service from the context.
func GetFeatureFlagService(c *gin.Context) (*featureflag.Service, bool) {
	if svc, exists := c.Get(FeatureFlagServiceKey); exists {
		if service, ok := svc.(*featureflag.Service); ok {
			return service, true
		}
	}
	return nil, false
}

// GetFeatureFlags retrieves all feature flags from the context.
func GetFeatureFlags(c *gin.Context) (map[string]bool, bool) {
	if flags, exists := c.Get(FeatureFlagsKey); exists {
		if flagMap, ok := flags.(map[string]bool); ok {
			return flagMap, true
		}
	}
	return nil, false
}

// CheckFeatureForTenant checks if a feature is enabled for a specific tenant.
// This is useful for admin operations that need to check flags for other tenants.
func CheckFeatureForTenant(c *gin.Context, flagKey string, tenantID uuid.UUID) bool {
	if service, ok := GetFeatureFlagService(c); ok {
		return service.IsEnabled(c.Request.Context(), flagKey, tenantID, uuid.Nil)
	}
	return false
}

// CheckFeatureForUser checks if a feature is enabled for a specific user.
// This is useful for admin operations that need to check flags for other users.
func CheckFeatureForUser(c *gin.Context, flagKey string, tenantID, userID uuid.UUID) bool {
	if service, ok := GetFeatureFlagService(c); ok {
		return service.IsEnabled(c.Request.Context(), flagKey, tenantID, userID)
	}
	return false
}
