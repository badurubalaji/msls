// Package errors provides standardized error types for the application.
// All error types follow RFC 7807 Problem Details specification.
package errors

import (
	"fmt"
	"net/http"
)

// Common problem type URIs following RFC 7807.
const (
	TypeBadRequest          = "https://httpstatuses.com/400"
	TypeUnauthorized        = "https://httpstatuses.com/401"
	TypeForbidden           = "https://httpstatuses.com/403"
	TypeNotFound            = "https://httpstatuses.com/404"
	TypeConflict            = "https://httpstatuses.com/409"
	TypeUnprocessableEntity = "https://httpstatuses.com/422"
	TypeTooManyRequests     = "https://httpstatuses.com/429"
	TypeInternalError       = "https://httpstatuses.com/500"
)

// AppError represents an application error following RFC 7807 Problem Details.
type AppError struct {
	Type       string                 `json:"type"`
	Title      string                 `json:"title"`
	Status     int                    `json:"status"`
	Detail     string                 `json:"detail,omitempty"`
	Instance   string                 `json:"instance,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Title, e.Detail)
	}
	return e.Title
}

// WithInstance adds the request path to the error.
func (e *AppError) WithInstance(instance string) *AppError {
	e.Instance = instance
	return e
}

// WithExtension adds an extension field to the error.
func (e *AppError) WithExtension(key string, value interface{}) *AppError {
	if e.Extensions == nil {
		e.Extensions = make(map[string]interface{})
	}
	e.Extensions[key] = value
	return e
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError represents a collection of field validation errors.
type ValidationError struct {
	*AppError
	Errors []FieldError `json:"errors,omitempty"`
}

// NewValidationError creates a new validation error with field errors.
func NewValidationError(detail string, errors []FieldError) *ValidationError {
	return &ValidationError{
		AppError: &AppError{
			Type:   TypeUnprocessableEntity,
			Title:  "Validation Error",
			Status: http.StatusUnprocessableEntity,
			Detail: detail,
		},
		Errors: errors,
	}
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return e.AppError.Error()
}

// NotFoundError creates a new not found error.
func NotFoundError(resource string) *AppError {
	return &AppError{
		Type:   TypeNotFound,
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("%s not found", resource),
	}
}

// NotFound creates a not found error with a custom message.
func NotFound(detail string) *AppError {
	return &AppError{
		Type:   TypeNotFound,
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Detail: detail,
	}
}

// ValidationFailed creates a validation error with a single message.
func ValidationFailed(detail string) *AppError {
	return &AppError{
		Type:   TypeUnprocessableEntity,
		Title:  "Validation Error",
		Status: http.StatusUnprocessableEntity,
		Detail: detail,
	}
}

// Unauthorized creates an unauthorized error.
func Unauthorized(detail string) *AppError {
	return &AppError{
		Type:   TypeUnauthorized,
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
		Detail: detail,
	}
}

// UnauthorizedDefault creates an unauthorized error with default message.
func UnauthorizedDefault() *AppError {
	return Unauthorized("Authentication required")
}

// Forbidden creates a forbidden error.
func Forbidden(detail string) *AppError {
	return &AppError{
		Type:   TypeForbidden,
		Title:  "Forbidden",
		Status: http.StatusForbidden,
		Detail: detail,
	}
}

// ForbiddenDefault creates a forbidden error with default message.
func ForbiddenDefault() *AppError {
	return Forbidden("You don't have permission to access this resource")
}

// BadRequest creates a bad request error.
func BadRequest(detail string) *AppError {
	return &AppError{
		Type:   TypeBadRequest,
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

// Conflict creates a conflict error.
func Conflict(detail string) *AppError {
	return &AppError{
		Type:   TypeConflict,
		Title:  "Conflict",
		Status: http.StatusConflict,
		Detail: detail,
	}
}

// TooManyRequests creates a rate limit error.
func TooManyRequests(retryAfter int) *AppError {
	err := &AppError{
		Type:   TypeTooManyRequests,
		Title:  "Too Many Requests",
		Status: http.StatusTooManyRequests,
		Detail: "Rate limit exceeded. Please try again later.",
	}
	if retryAfter > 0 {
		err.WithExtension("retry_after", retryAfter)
	}
	return err
}

// InternalError creates an internal server error.
func InternalError(detail string) *AppError {
	return &AppError{
		Type:   TypeInternalError,
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: detail,
	}
}

// InternalErrorDefault creates an internal server error with default message.
func InternalErrorDefault() *AppError {
	return InternalError("An unexpected error occurred")
}

// Wrap wraps an existing error into an AppError.
func Wrap(err error, status int, title string) *AppError {
	return &AppError{
		Type:   fmt.Sprintf("https://httpstatuses.com/%d", status),
		Title:  title,
		Status: status,
		Detail: err.Error(),
	}
}

// Is checks if the error is an AppError with the specified status code.
func Is(err error, status int) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status == status
	}
	if valErr, ok := err.(*ValidationError); ok {
		return valErr.Status == status
	}
	return false
}

// IsNotFound checks if the error is a not found error.
func IsNotFound(err error) bool {
	return Is(err, http.StatusNotFound)
}

// IsUnauthorized checks if the error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return Is(err, http.StatusUnauthorized)
}

// IsForbidden checks if the error is a forbidden error.
func IsForbidden(err error) bool {
	return Is(err, http.StatusForbidden)
}

// IsValidation checks if the error is a validation error.
func IsValidation(err error) bool {
	if _, ok := err.(*ValidationError); ok {
		return true
	}
	return Is(err, http.StatusUnprocessableEntity)
}
