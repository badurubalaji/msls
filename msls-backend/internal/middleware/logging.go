// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"msls-backend/internal/pkg/logger"
)

// LoggingConfig holds configuration for logging middleware.
type LoggingConfig struct {
	// Logger is the logger instance to use.
	Logger *logger.Logger

	// SkipPaths lists paths to exclude from logging (e.g., health checks).
	SkipPaths []string

	// SkipHealthCheck enables automatic skipping of /health and /ready endpoints.
	SkipHealthCheck bool

	// LogRequestBody enables logging of request body (use with caution).
	LogRequestBody bool

	// LogResponseBody enables logging of response body (use with caution).
	LogResponseBody bool

	// SlowThreshold defines the duration after which a request is considered slow.
	SlowThreshold time.Duration
}

// DefaultLoggingConfig returns the default logging configuration.
func DefaultLoggingConfig(log *logger.Logger) LoggingConfig {
	return LoggingConfig{
		Logger:          log,
		SkipPaths:       []string{},
		SkipHealthCheck: true,
		LogRequestBody:  false,
		LogResponseBody: false,
		SlowThreshold:   500 * time.Millisecond,
	}
}

// Logging returns a middleware that logs HTTP requests and responses.
// It captures method, path, status code, and duration for each request.
func Logging(config LoggingConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	// Add default health check paths if enabled
	if config.SkipHealthCheck {
		skipPaths["/health"] = true
		skipPaths["/ready"] = true
		skipPaths["/metrics"] = true
	}

	return func(c *gin.Context) {
		// Skip logging for excluded paths
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Record start time
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate duration
		latency := time.Since(start)
		status := c.Writer.Status()

		// Get request and tenant IDs
		requestID := GetRequestID(c)
		tenantID := GetTenantID(c)

		// Build base fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("response_size", c.Writer.Size()),
		}

		// Add optional fields
		if rawQuery != "" {
			fields = append(fields, zap.String("query", rawQuery))
		}
		if tenantID != "" {
			fields = append(fields, zap.String("tenant_id", tenantID))
		}
		if c.Request.UserAgent() != "" {
			fields = append(fields, zap.String("user_agent", c.Request.UserAgent()))
		}
		if c.Request.Referer() != "" {
			fields = append(fields, zap.String("referer", c.Request.Referer()))
		}

		// Add error information if present
		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
		}

		// Choose log level based on status code and latency
		switch {
		case status >= 500:
			config.Logger.Error("request completed", fields...)
		case status >= 400:
			config.Logger.Warn("request completed", fields...)
		case latency > config.SlowThreshold:
			fields = append(fields, zap.Bool("slow_request", true))
			config.Logger.Warn("slow request", fields...)
		default:
			config.Logger.Info("request completed", fields...)
		}
	}
}

// LoggingDefault returns a middleware with default configuration.
func LoggingDefault(log *logger.Logger) gin.HandlerFunc {
	return Logging(DefaultLoggingConfig(log))
}

// LoggingWithSkipPaths returns a middleware that skips logging for specified paths.
func LoggingWithSkipPaths(log *logger.Logger, skipPaths []string) gin.HandlerFunc {
	config := DefaultLoggingConfig(log)
	config.SkipPaths = skipPaths
	return Logging(config)
}

// LoggingWithSlowThreshold returns a middleware with custom slow request threshold.
func LoggingWithSlowThreshold(log *logger.Logger, threshold time.Duration) gin.HandlerFunc {
	config := DefaultLoggingConfig(log)
	config.SlowThreshold = threshold
	return Logging(config)
}
