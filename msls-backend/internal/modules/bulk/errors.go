// Package bulk provides bulk operation functionality.
package bulk

import "errors"

// Errors returned by the bulk package.
var (
	ErrOperationNotFound    = errors.New("bulk operation not found")
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrNoStudentsProvided   = errors.New("no students provided for bulk operation")
	ErrTooManyStudents      = errors.New("too many students for bulk operation")
	ErrInvalidExportFormat  = errors.New("invalid export format")
	ErrExportFailed         = errors.New("export failed")
	ErrOperationInProgress  = errors.New("operation is already in progress")
	ErrOperationCancelled   = errors.New("operation was cancelled")
)

// MaxExportRecords is the maximum number of records allowed in an export.
const MaxExportRecords = 10000

// MaxBulkStudents is the maximum number of students for a bulk operation.
const MaxBulkStudents = 1000
