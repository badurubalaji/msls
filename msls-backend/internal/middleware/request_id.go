// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/pkg/database"
)

const (
	// RequestIDHeader is the HTTP header for request ID.
	RequestIDHeader = "X-Request-ID"

	// RequestIDKey is the context key for storing request ID.
	RequestIDKey = "request_id"
)

// RequestIDConfig holds configuration for request ID middleware.
type RequestIDConfig struct {
	// Generator is a function that generates a new request ID.
	// If nil, UUID v4 is used.
	Generator func() string

	// PropagateHeader determines if the incoming X-Request-ID header should be used.
	PropagateHeader bool
}

// DefaultRequestIDConfig returns the default request ID configuration.
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		Generator:       nil,
		PropagateHeader: true,
	}
}

// RequestID returns a middleware that generates or propagates request IDs.
// It sets the request ID in both the Gin context and the response header.
func RequestID(config RequestIDConfig) gin.HandlerFunc {
	generator := config.Generator
	if generator == nil {
		generator = func() string {
			return uuid.New().String()
		}
	}

	return func(c *gin.Context) {
		var requestID string

		// Check for incoming request ID header
		if config.PropagateHeader {
			requestID = c.GetHeader(RequestIDHeader)
		}

		// Generate new ID if not provided
		if requestID == "" {
			requestID = generator()
		}

		// Set request ID in Gin context
		c.Set(RequestIDKey, requestID)

		// Update request context with request ID
		ctx := database.ContextWithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// Set response header
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// RequestIDDefault returns a middleware with default configuration.
func RequestIDDefault() gin.HandlerFunc {
	return RequestID(DefaultRequestIDConfig())
}

// RequestIDWithGenerator returns a middleware with a custom ID generator.
func RequestIDWithGenerator(generator func() string) gin.HandlerFunc {
	return RequestID(RequestIDConfig{
		Generator:       generator,
		PropagateHeader: true,
	})
}

// RequestIDNoPropagate returns a middleware that always generates new IDs.
func RequestIDNoPropagate() gin.HandlerFunc {
	return RequestID(RequestIDConfig{
		Generator:       nil,
		PropagateHeader: false,
	})
}

// GetRequestID retrieves the request ID from the Gin context.
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetRequestIDFromContext retrieves the request ID from the request context.
func GetRequestIDFromContext(c *gin.Context) string {
	return database.RequestID(c.Request.Context())
}
