// Package enrollment provides student enrollment management functionality.
package enrollment

import "errors"

// Enrollment validation and domain errors.
var (
	// Validation errors
	ErrTenantIDRequired       = errors.New("tenant ID is required")
	ErrStudentIDRequired      = errors.New("student ID is required")
	ErrAcademicYearIDRequired = errors.New("academic year ID is required")
	ErrEnrollmentIDRequired   = errors.New("enrollment ID is required")
	ErrChangedByRequired      = errors.New("changed by user ID is required")
	ErrInvalidStatus          = errors.New("invalid enrollment status")
	ErrTransferDateRequired   = errors.New("transfer date is required")
	ErrTransferReasonRequired = errors.New("transfer reason is required")
	ErrDropoutDateRequired    = errors.New("dropout date is required")
	ErrDropoutReasonRequired  = errors.New("dropout reason is required")

	// Domain errors
	ErrEnrollmentNotFound           = errors.New("enrollment not found")
	ErrActiveEnrollmentNotFound     = errors.New("active enrollment not found")
	ErrDuplicateEnrollment          = errors.New("student already enrolled in this academic year")
	ErrActiveEnrollmentExists       = errors.New("student already has an active enrollment")
	ErrInvalidStatusTransition      = errors.New("invalid status transition")
	ErrEnrollmentAlreadyCompleted   = errors.New("enrollment is already completed")
	ErrEnrollmentAlreadyTransferred = errors.New("enrollment is already transferred")
	ErrAcademicYearNotFound         = errors.New("academic year not found")
	ErrStudentNotFound              = errors.New("student not found")
	ErrCannotModifyInactiveEnrollment = errors.New("cannot modify inactive enrollment")
)
