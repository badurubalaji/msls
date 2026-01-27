// Package academicyear provides academic year management services.
package academicyear

import "errors"

// Service-level errors for academic year operations.
var (
	// ErrAcademicYearNotFound is returned when an academic year is not found.
	ErrAcademicYearNotFound = errors.New("academic year not found")

	// ErrAcademicYearNameRequired is returned when academic year name is missing.
	ErrAcademicYearNameRequired = errors.New("academic year name is required")

	// ErrAcademicYearNameExists is returned when academic year name already exists.
	ErrAcademicYearNameExists = errors.New("academic year with this name already exists")

	// ErrAcademicYearStartDateRequired is returned when start date is missing.
	ErrAcademicYearStartDateRequired = errors.New("start date is required")

	// ErrAcademicYearEndDateRequired is returned when end date is missing.
	ErrAcademicYearEndDateRequired = errors.New("end date is required")

	// ErrAcademicYearInvalidDates is returned when end date is not after start date.
	ErrAcademicYearInvalidDates = errors.New("end date must be after start date")

	// ErrAcademicYearOverlap is returned when dates overlap with existing academic year.
	ErrAcademicYearOverlap = errors.New("academic year dates overlap with existing year")

	// ErrAcademicYearHasDependencies is returned when academic year has associated records.
	ErrAcademicYearHasDependencies = errors.New("academic year has associated records and cannot be deleted")

	// ErrCannotSetInactiveAsCurrent is returned when trying to set an inactive year as current.
	ErrCannotSetInactiveAsCurrent = errors.New("cannot set inactive academic year as current")

	// ErrTenantIDRequired is returned when tenant ID is missing.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrTermNotFound is returned when an academic term is not found.
	ErrTermNotFound = errors.New("academic term not found")

	// ErrTermNameRequired is returned when term name is missing.
	ErrTermNameRequired = errors.New("term name is required")

	// ErrTermNameExists is returned when term name already exists for the academic year.
	ErrTermNameExists = errors.New("term with this name already exists in the academic year")

	// ErrTermStartDateRequired is returned when term start date is missing.
	ErrTermStartDateRequired = errors.New("term start date is required")

	// ErrTermEndDateRequired is returned when term end date is missing.
	ErrTermEndDateRequired = errors.New("term end date is required")

	// ErrTermInvalidDates is returned when term end date is not after start date.
	ErrTermInvalidDates = errors.New("term end date must be after start date")

	// ErrTermOutsideAcademicYear is returned when term dates are outside academic year bounds.
	ErrTermOutsideAcademicYear = errors.New("term dates must be within academic year dates")

	// ErrTermOverlap is returned when term dates overlap with existing term.
	ErrTermOverlap = errors.New("term dates overlap with existing term")

	// ErrHolidayNotFound is returned when a holiday is not found.
	ErrHolidayNotFound = errors.New("holiday not found")

	// ErrHolidayNameRequired is returned when holiday name is missing.
	ErrHolidayNameRequired = errors.New("holiday name is required")

	// ErrHolidayDateRequired is returned when holiday date is missing.
	ErrHolidayDateRequired = errors.New("holiday date is required")

	// ErrHolidayOutsideAcademicYear is returned when holiday date is outside academic year.
	ErrHolidayOutsideAcademicYear = errors.New("holiday date must be within academic year dates")

	// ErrHolidayInvalidType is returned when an invalid holiday type is provided.
	ErrHolidayInvalidType = errors.New("invalid holiday type")
)
