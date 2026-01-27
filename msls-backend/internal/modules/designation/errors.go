// Package designation provides designation management functionality.
package designation

import "errors"

// Errors for designation operations.
var (
	ErrDesignationNotFound  = errors.New("designation not found")
	ErrDuplicateName        = errors.New("designation name already exists")
	ErrDesignationInUse     = errors.New("designation is in use and cannot be deleted")
	ErrInvalidDesignationID = errors.New("invalid designation ID")
	ErrInvalidLevel         = errors.New("invalid designation level (must be 1-10)")
)
