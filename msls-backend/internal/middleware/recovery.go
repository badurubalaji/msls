// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"msls-backend/internal/pkg/logger"
)

// RecoveryConfig holds configuration for recovery middleware.
type RecoveryConfig struct {
	// Logger is the logger instance for logging panics.
	Logger *logger.Logger

	// PrintStack enables printing of stack traces in logs.
	PrintStack bool

	// DisableStackAll disables all stack trace logging (even in debug mode).
	DisableStackAll bool

	// OnPanic is an optional callback function called when a panic occurs.
	// It receives the Gin context and the recovered panic value.
	OnPanic func(c *gin.Context, recovered interface{})
}

// DefaultRecoveryConfig returns the default recovery configuration.
func DefaultRecoveryConfig(log *logger.Logger) RecoveryConfig {
	return RecoveryConfig{
		Logger:          log,
		PrintStack:      true,
		DisableStackAll: false,
		OnPanic:         nil,
	}
}

// Recovery returns a middleware that recovers from panics and returns a 500 error.
// It logs the panic with stack trace for debugging purposes.
func Recovery(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// Get request and tenant IDs
				requestID := GetRequestID(c)
				tenantID := GetTenantID(c)

				// Build log fields
				fields := []zap.Field{
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.Any("panic", recovered),
				}

				if tenantID != "" {
					fields = append(fields, zap.String("tenant_id", tenantID))
				}

				// Add stack trace if enabled
				if config.PrintStack && !config.DisableStackAll {
					fields = append(fields, zap.String("stack", string(debug.Stack())))
				}

				// Log the panic
				config.Logger.Error("panic recovered", fields...)

				// Call optional panic handler
				if config.OnPanic != nil {
					config.OnPanic(c, recovered)
				}

				// Return 500 Internal Server Error with RFC 7807 format
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"type":     "https://httpstatuses.com/500",
					"title":    "Internal Server Error",
					"status":   http.StatusInternalServerError,
					"detail":   "An unexpected error occurred",
					"instance": c.Request.URL.Path,
				})
			}
		}()

		c.Next()
	}
}

// RecoveryDefault returns a middleware with default configuration.
func RecoveryDefault(log *logger.Logger) gin.HandlerFunc {
	return Recovery(DefaultRecoveryConfig(log))
}

// RecoveryWithCallback returns a middleware with a custom panic callback.
func RecoveryWithCallback(log *logger.Logger, callback func(c *gin.Context, recovered interface{})) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.OnPanic = callback
	return Recovery(config)
}

// RecoveryWithoutStack returns a middleware that doesn't log stack traces.
func RecoveryWithoutStack(log *logger.Logger) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.PrintStack = false
	return Recovery(config)
}

// CustomRecoveryWithWriter is a variant that writes recovery info with RFC 7807 format.
func CustomRecoveryWithWriter(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for broken pipe errors (client disconnected)
				var brokenPipe bool
				if errStr := fmt.Sprintf("%v", err); errStr != "" {
					brokenPipe = strings.Contains(errStr, "broken pipe") ||
						strings.Contains(errStr, "connection reset by peer")
				}

				requestID := GetRequestID(c)

				// Log the error
				log.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Any("error", err),
					zap.Bool("broken_pipe", brokenPipe),
					zap.String("stack", string(debug.Stack())),
				)

				// If the connection is dead, we can't write a status to it
				if brokenPipe {
					c.Abort()
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"type":     "https://httpstatuses.com/500",
					"title":    "Internal Server Error",
					"status":   http.StatusInternalServerError,
					"detail":   "An unexpected error occurred",
					"instance": c.Request.URL.Path,
				})
			}
		}()

		c.Next()
	}
}
