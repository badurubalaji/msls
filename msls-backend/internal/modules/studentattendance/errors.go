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

// Period-wise attendance errors (Story 7.2).
var (
	// ErrPeriodNotFound is returned when a period slot is not found.
	ErrPeriodNotFound = errors.New("period slot not found")

	// ErrTimetableEntryNotFound is returned when a timetable entry is not found.
	ErrTimetableEntryNotFound = errors.New("timetable entry not found")

	// ErrNoTimetableForSection is returned when no published timetable exists for a section.
	ErrNoTimetableForSection = errors.New("no published timetable found for this section")

	// ErrPeriodAttendanceNotEnabled is returned when period attendance is not enabled.
	ErrPeriodAttendanceNotEnabled = errors.New("period-wise attendance is not enabled for this branch")

	// ErrInvalidPeriodForSection is returned when the period doesn't belong to the section's timetable.
	ErrInvalidPeriodForSection = errors.New("period does not belong to this section's timetable")

	// ErrPeriodIDRequired is returned when period ID is not provided.
	ErrPeriodIDRequired = errors.New("period ID is required")

	// ErrSubjectIDRequired is returned when subject ID is not provided.
	ErrSubjectIDRequired = errors.New("subject ID is required")
)

// Edit and audit errors (Story 7.3).
var (
	// ErrEditReasonRequired is returned when edit reason is not provided.
	ErrEditReasonRequired = errors.New("edit reason is required")

	// ErrNotOriginalMarker is returned when user is not the original marker and not admin.
	ErrNotOriginalMarker = errors.New("only the original marker or admin can edit within edit window")

	// ErrAdminOnlyEdit is returned when edit window expired and user is not admin.
	ErrAdminOnlyEdit = errors.New("edit window expired - admin approval required")

	// ErrAuditNotFound is returned when audit record is not found.
	ErrAuditNotFound = errors.New("audit record not found")

	// ErrAttendanceIDRequired is returned when attendance ID is not provided.
	ErrAttendanceIDRequired = errors.New("attendance ID is required")
)
