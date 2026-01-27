// Package response provides standardized HTTP response structures.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success represents a successful API response.
type Success struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains pagination and additional metadata.
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Problem represents an RFC 7807 Problem Details response.
// See: https://datatracker.ietf.org/doc/html/rfc7807
type Problem struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   string      `json:"detail,omitempty"`
	Instance string      `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Common problem types (URIs for identifying error types).
const (
	ProblemTypeBadRequest          = "https://httpstatuses.com/400"
	ProblemTypeUnauthorized        = "https://httpstatuses.com/401"
	ProblemTypeForbidden           = "https://httpstatuses.com/403"
	ProblemTypeNotFound            = "https://httpstatuses.com/404"
	ProblemTypeConflict            = "https://httpstatuses.com/409"
	ProblemTypeUnprocessableEntity = "https://httpstatuses.com/422"
	ProblemTypeInternalError       = "https://httpstatuses.com/500"
)

// OK sends a 200 OK response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Success{
		Success: true,
		Data:    data,
	})
}

// OKWithMeta sends a 200 OK response with data and metadata.
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Success{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Success{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request error response.
func BadRequest(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeBadRequest,
		Title:    "Bad Request",
		Status:   http.StatusBadRequest,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusBadRequest, problem)
}

// BadRequestWithErrors sends a 400 Bad Request error with field errors.
func BadRequestWithErrors(c *gin.Context, detail string, errors []FieldError) {
	problem := Problem{
		Type:     ProblemTypeBadRequest,
		Title:    "Bad Request",
		Status:   http.StatusBadRequest,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		Errors:   errors,
	}
	c.JSON(http.StatusBadRequest, problem)
}

// Unauthorized sends a 401 Unauthorized error response.
func Unauthorized(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeUnauthorized,
		Title:    "Unauthorized",
		Status:   http.StatusUnauthorized,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusUnauthorized, problem)
}

// Forbidden sends a 403 Forbidden error response.
func Forbidden(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeForbidden,
		Title:    "Forbidden",
		Status:   http.StatusForbidden,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusForbidden, problem)
}

// NotFound sends a 404 Not Found error response.
func NotFound(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeNotFound,
		Title:    "Not Found",
		Status:   http.StatusNotFound,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusNotFound, problem)
}

// Conflict sends a 409 Conflict error response.
func Conflict(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeConflict,
		Title:    "Conflict",
		Status:   http.StatusConflict,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusConflict, problem)
}

// UnprocessableEntity sends a 422 Unprocessable Entity error response.
func UnprocessableEntity(c *gin.Context, detail string, errors []FieldError) {
	problem := Problem{
		Type:     ProblemTypeUnprocessableEntity,
		Title:    "Unprocessable Entity",
		Status:   http.StatusUnprocessableEntity,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		Errors:   errors,
	}
	c.JSON(http.StatusUnprocessableEntity, problem)
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(c *gin.Context, detail string) {
	problem := Problem{
		Type:     ProblemTypeInternalError,
		Title:    "Internal Server Error",
		Status:   http.StatusInternalServerError,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(http.StatusInternalServerError, problem)
}

// Error sends a custom error response with the given status code.
func Error(c *gin.Context, status int, problemType, title, detail string) {
	problem := Problem{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	c.JSON(status, problem)
}
