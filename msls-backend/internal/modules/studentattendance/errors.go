// Package studentattendance provides student attendance management functionality.
package studentattendance

import "errors"

// Repository errors.
var (
	// ErrAttendanceNotFound is returned when an attendance record is not found.
	ErrAttendanceNotFound = errors.New("student attendance record not found")

	// ErrSettingsNotFound is returned when attendance settings are not found.
	ErrSettingsNotFound = errors.New("student attendance settings not found")

	// ErrDuplicateAttendance is returned when attendance already exists for the student on the date.
	ErrDuplicateAttendance = errors.New("attendance already marked for this student on this date")

	// ErrSectionNotFound is returned when a section is not found.
	ErrSectionNotFound = errors.New("section not found")

	// ErrStudentNotFound is returned when a student is not found.
	ErrStudentNotFound = errors.New("student not found")
)

// Validation errors.
var (
	// ErrTenantIDRequired is returned when tenant ID is not provided.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrSectionIDRequired is returned when section ID is not provided.
	ErrSectionIDRequired = errors.New("section ID is required")

	// ErrStudentIDRequired is returned when student ID is not provided.
	ErrStudentIDRequired = errors.New("student ID is required")

	// ErrBranchIDRequired is returned when branch ID is not provided.
	ErrBranchIDRequired = errors.New("branch ID is required")

	// ErrInvalidStatus is returned when an invalid attendance status is provided.
	ErrInvalidStatus = errors.New("invalid attendance status")

	// ErrDateRequired is returned when date is required but not provided.
	ErrDateRequired = errors.New("date is required")

	// ErrFutureDate is returned when trying to mark attendance for a future date.
	ErrFutureDate = errors.New("cannot mark attendance for future date")

	// ErrEditWindowExpired is returned when the edit window has expired.
	ErrEditWindowExpired = errors.New("attendance edit window has expired")

	// ErrNoStudentsInSection is returned when a section has no students.
	ErrNoStudentsInSection = errors.New("no students found in section")

	// ErrEmptyAttendanceRecords is returned when no attendance records are provided.
	ErrEmptyAttendanceRecords = errors.New("attendance records cannot be empty")

	// ErrUnauthorized is returned when the user is not authorized to perform the action.
	ErrUnauthorized = errors.New("not authorized to mark attendance for this section")
)
