// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for CORS middleware.
type CORSConfig struct {
	// AllowOrigins is a list of origins that may access the resource.
	// Use "*" to allow all origins (not recommended for production).
	AllowOrigins []string

	// AllowMethods is a list of methods the client is allowed to use.
	AllowMethods []string

	// AllowHeaders is a list of headers the client is allowed to use.
	AllowHeaders []string

	// ExposeHeaders is a list of headers that are safe to expose.
	ExposeHeaders []string

	// AllowCredentials indicates whether the request can include credentials.
	AllowCredentials bool

	// MaxAge indicates how long the results of a preflight request can be cached.
	MaxAge time.Duration

	// AllowWildcard enables wildcard matching for origins.
	AllowWildcard bool

	// AllowOriginFunc is a custom function to validate origins.
	// If set, AllowOrigins is ignored.
	AllowOriginFunc func(origin string) bool
}

// DefaultCORSConfig returns the default CORS configuration.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Type",
			"Content-Language",
			"Origin",
			"Authorization",
			"X-Request-ID",
			"X-Tenant-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    false,
		AllowOriginFunc:  nil,
	}
}

// ProductionCORSConfig returns a more restrictive CORS configuration for production.
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Type",
			"Content-Language",
			"Origin",
			"Authorization",
			"X-Request-ID",
			"X-Tenant-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    false,
		AllowOriginFunc:  nil,
	}
}

// CORS returns a middleware that handles Cross-Origin Resource Sharing (CORS).
func CORS(config CORSConfig) gin.HandlerFunc {
	// Normalize configuration
	normalizeConfig(&config)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// If no origin header, skip CORS processing
		if origin == "" {
			c.Next()
			return
		}

		// Check if origin is allowed
		if !isOriginAllowed(origin, config) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set CORS headers
		setCORSHeaders(c, origin, config)

		// Handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSDefault returns a middleware with default configuration.
func CORSDefault() gin.HandlerFunc {
	return CORS(DefaultCORSConfig())
}

// CORSWithOrigins returns a middleware that allows specific origins.
func CORSWithOrigins(origins ...string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOrigins = origins
	return CORS(config)
}

// CORSWithCredentials returns a middleware that allows credentials.
func CORSWithCredentials(origins []string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOrigins = origins
	config.AllowCredentials = true
	return CORS(config)
}

// CORSWithCustomValidator returns a middleware with a custom origin validator.
func CORSWithCustomValidator(validator func(origin string) bool) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOriginFunc = validator
	return CORS(config)
}

// normalizeConfig normalizes the CORS configuration.
func normalizeConfig(config *CORSConfig) {
	// Ensure at least one origin is configured
	if len(config.AllowOrigins) == 0 && config.AllowOriginFunc == nil {
		config.AllowOrigins = []string{"*"}
	}

	// Ensure at least one method is configured
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{http.MethodGet, http.MethodPost}
	}

	// Normalize methods to uppercase
	for i, method := range config.AllowMethods {
		config.AllowMethods[i] = strings.ToUpper(method)
	}
}

// isOriginAllowed checks if the given origin is allowed.
func isOriginAllowed(origin string, config CORSConfig) bool {
	// Use custom validator if provided
	if config.AllowOriginFunc != nil {
		return config.AllowOriginFunc(origin)
	}

	// Check against allowed origins
	for _, allowed := range config.AllowOrigins {
		if allowed == "*" {
			return true
		}

		if config.AllowWildcard && strings.HasPrefix(allowed, "*.") {
			// Wildcard matching (e.g., *.example.com)
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) || origin == "https://"+strings.TrimPrefix(domain, ".") {
				return true
			}
		} else if allowed == origin {
			return true
		}
	}

	return false
}

// setCORSHeaders sets the appropriate CORS headers on the response.
func setCORSHeaders(c *gin.Context, origin string, config CORSConfig) {
	// Set Access-Control-Allow-Origin
	// Note: When credentials are allowed, we can't use "*"
	if config.AllowCredentials || !containsWildcard(config.AllowOrigins) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
	}

	// Set Access-Control-Allow-Credentials
	if config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// Set Access-Control-Expose-Headers
	if len(config.ExposeHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
	}

	// Preflight-specific headers
	if c.Request.Method == http.MethodOptions {
		// Set Access-Control-Allow-Methods
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))

		// Set Access-Control-Allow-Headers
		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		// Set Access-Control-Max-Age
		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", strconv.FormatInt(int64(config.MaxAge.Seconds()), 10))
		}
	}
}

// containsWildcard checks if the origins list contains a wildcard.
func containsWildcard(origins []string) bool {
	for _, origin := range origins {
		if origin == "*" {
			return true
		}
	}
	return false
}
