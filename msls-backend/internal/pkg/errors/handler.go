// Package errors provides error handling middleware for Gin applications.
package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"msls-backend/internal/pkg/logger"
)

// Handler is a middleware that handles errors and returns RFC 7807 responses.
func Handler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Get request ID for logging
		requestID, _ := c.Get("request_id")

		// Handle different error types
		switch e := err.(type) {
		case *AppError:
			e.Instance = c.Request.URL.Path
			logError(log, requestID, e)
			c.JSON(e.Status, e)

		case *ValidationError:
			e.Instance = c.Request.URL.Path
			logError(log, requestID, e.AppError)
			c.JSON(e.Status, e)

		default:
			// Log the original error
			log.Error("unhandled error",
				zap.Any("request_id", requestID),
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)

			// Return a generic internal server error
			appErr := InternalErrorDefault()
			appErr.Instance = c.Request.URL.Path
			c.JSON(http.StatusInternalServerError, appErr)
		}

		// Abort to prevent further processing
		c.Abort()
	}
}

// logError logs the error with appropriate level based on status code.
func logError(log *logger.Logger, requestID interface{}, err *AppError) {
	fields := []zap.Field{
		zap.Any("request_id", requestID),
		zap.String("type", err.Type),
		zap.String("title", err.Title),
		zap.Int("status", err.Status),
		zap.String("detail", err.Detail),
		zap.String("instance", err.Instance),
	}

	// Log level based on status code
	switch {
	case err.Status >= 500:
		log.Error("server error", fields...)
	case err.Status >= 400:
		log.Warn("client error", fields...)
	default:
		log.Info("error response", fields...)
	}
}

// Abort aborts the request with an AppError.
func Abort(c *gin.Context, err *AppError) {
	err.Instance = c.Request.URL.Path
	c.AbortWithStatusJSON(err.Status, err)
}

// AbortWithValidation aborts the request with a ValidationError.
func AbortWithValidation(c *gin.Context, err *ValidationError) {
	err.Instance = c.Request.URL.Path
	c.AbortWithStatusJSON(err.Status, err)
}

// AbortNotFound aborts with a not found error.
func AbortNotFound(c *gin.Context, resource string) {
	Abort(c, NotFoundError(resource))
}

// AbortUnauthorized aborts with an unauthorized error.
func AbortUnauthorized(c *gin.Context, detail string) {
	Abort(c, Unauthorized(detail))
}

// AbortForbidden aborts with a forbidden error.
func AbortForbidden(c *gin.Context, detail string) {
	Abort(c, Forbidden(detail))
}

// AbortBadRequest aborts with a bad request error.
func AbortBadRequest(c *gin.Context, detail string) {
	Abort(c, BadRequest(detail))
}

// AbortInternalError aborts with an internal server error.
func AbortInternalError(c *gin.Context) {
	Abort(c, InternalErrorDefault())
}

// AbortTooManyRequests aborts with a rate limit error.
func AbortTooManyRequests(c *gin.Context, retryAfter int) {
	err := TooManyRequests(retryAfter)
	c.Header("Retry-After", string(rune(retryAfter)))
	Abort(c, err)
}

// JSON sends an error response without aborting.
func JSON(c *gin.Context, err *AppError) {
	err.Instance = c.Request.URL.Path
	c.JSON(err.Status, err)
}

// JSONValidation sends a validation error response without aborting.
func JSONValidation(c *gin.Context, err *ValidationError) {
	err.Instance = c.Request.URL.Path
	c.JSON(err.Status, err)
}
