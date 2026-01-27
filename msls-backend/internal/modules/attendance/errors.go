// Package attendance provides staff attendance management functionality.
package attendance

import "errors"

// Repository errors.
var (
	// ErrAttendanceNotFound is returned when an attendance record is not found.
	ErrAttendanceNotFound = errors.New("attendance record not found")

	// ErrRegularizationNotFound is returned when a regularization request is not found.
	ErrRegularizationNotFound = errors.New("regularization request not found")

	// ErrSettingsNotFound is returned when attendance settings are not found.
	ErrSettingsNotFound = errors.New("attendance settings not found")

	// ErrDuplicateAttendance is returned when attendance already exists for the date.
	ErrDuplicateAttendance = errors.New("attendance already marked for this date")

	// ErrAlreadyCheckedIn is returned when staff has already checked in.
	ErrAlreadyCheckedIn = errors.New("already checked in for today")

	// ErrNotCheckedIn is returned when trying to check out without checking in.
	ErrNotCheckedIn = errors.New("not checked in yet")

	// ErrAlreadyCheckedOut is returned when staff has already checked out.
	ErrAlreadyCheckedOut = errors.New("already checked out for today")
)

// Validation errors.
var (
	// ErrTenantIDRequired is returned when tenant ID is not provided.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrStaffIDRequired is returned when staff ID is not provided.
	ErrStaffIDRequired = errors.New("staff ID is required")

	// ErrStaffNotFound is returned when staff is not found.
	ErrStaffNotFound = errors.New("staff not found")

	// ErrBranchIDRequired is returned when branch ID is not provided.
	ErrBranchIDRequired = errors.New("branch ID is required")

	// ErrInvalidStatus is returned when an invalid attendance status is provided.
	ErrInvalidStatus = errors.New("invalid attendance status")

	// ErrInvalidHalfDayType is returned when an invalid half day type is provided.
	ErrInvalidHalfDayType = errors.New("invalid half day type")

	// ErrReasonRequired is returned when reason is required but not provided.
	ErrReasonRequired = errors.New("reason is required")

	// ErrDateRequired is returned when date is required but not provided.
	ErrDateRequired = errors.New("date is required")

	// ErrFutureDate is returned when trying to mark attendance for a future date.
	ErrFutureDate = errors.New("cannot mark attendance for future date")

	// ErrInvalidRegularizationStatus is returned when an invalid regularization status is provided.
	ErrInvalidRegularizationStatus = errors.New("invalid regularization status")

	// ErrRegularizationAlreadyProcessed is returned when trying to process an already processed request.
	ErrRegularizationAlreadyProcessed = errors.New("regularization request already processed")

	// ErrCannotRegularizePendingRequest is returned when trying to regularize with a pending request.
	ErrCannotRegularizePendingRequest = errors.New("cannot request regularization for date with pending request")
)
