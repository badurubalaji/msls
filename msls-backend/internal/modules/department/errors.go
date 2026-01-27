// Package department provides department management functionality.
package department

import "errors"

// Errors for department operations.
var (
	ErrDepartmentNotFound  = errors.New("department not found")
	ErrDuplicateCode       = errors.New("department code already exists")
	ErrDepartmentInUse     = errors.New("department is in use and cannot be deleted")
	ErrInvalidDepartmentID = errors.New("invalid department ID")
)
