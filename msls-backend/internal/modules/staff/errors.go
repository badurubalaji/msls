// Package staff provides staff management functionality.
package staff

import "errors"

// Repository errors.
var (
	// ErrStaffNotFound is returned when a staff member is not found.
	ErrStaffNotFound = errors.New("staff member not found")

	// ErrDuplicateEmployeeID is returned when the employee ID already exists.
	ErrDuplicateEmployeeID = errors.New("employee ID already exists")

	// ErrOptimisticLockConflict is returned when there's a concurrent edit conflict.
	ErrOptimisticLockConflict = errors.New("concurrent modification detected, please refresh and try again")

	// ErrDepartmentNotFound is returned when a department is not found.
	ErrDepartmentNotFound = errors.New("department not found")

	// ErrDesignationNotFound is returned when a designation is not found.
	ErrDesignationNotFound = errors.New("designation not found")
)

// Validation errors.
var (
	// ErrTenantIDRequired is returned when tenant ID is not provided.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrBranchIDRequired is returned when branch ID is not provided.
	ErrBranchIDRequired = errors.New("branch ID is required")

	// ErrFirstNameRequired is returned when first name is not provided.
	ErrFirstNameRequired = errors.New("first name is required")

	// ErrLastNameRequired is returned when last name is not provided.
	ErrLastNameRequired = errors.New("last name is required")

	// ErrDateOfBirthRequired is returned when date of birth is not provided.
	ErrDateOfBirthRequired = errors.New("date of birth is required")

	// ErrGenderRequired is returned when gender is not provided.
	ErrGenderRequired = errors.New("gender is required")

	// ErrInvalidGender is returned when an invalid gender value is provided.
	ErrInvalidGender = errors.New("invalid gender value")

	// ErrInvalidStatus is returned when an invalid status value is provided.
	ErrInvalidStatus = errors.New("invalid status value")

	// ErrInvalidStaffType is returned when an invalid staff type is provided.
	ErrInvalidStaffType = errors.New("invalid staff type value")

	// ErrBranchNotFound is returned when the branch is not found.
	ErrBranchNotFound = errors.New("branch not found")

	// ErrInvalidDateOfBirth is returned when date of birth is in the future.
	ErrInvalidDateOfBirth = errors.New("date of birth cannot be in the future")

	// ErrWorkEmailRequired is returned when work email is not provided.
	ErrWorkEmailRequired = errors.New("work email is required")

	// ErrWorkPhoneRequired is returned when work phone is not provided.
	ErrWorkPhoneRequired = errors.New("work phone is required")

	// ErrJoinDateRequired is returned when join date is not provided.
	ErrJoinDateRequired = errors.New("join date is required")

	// ErrStatusReasonRequired is returned when status change reason is not provided.
	ErrStatusReasonRequired = errors.New("status change reason is required")

	// ErrEffectiveDateRequired is returned when effective date is not provided.
	ErrEffectiveDateRequired = errors.New("effective date is required")

	// ErrReportingManagerNotFound is returned when the reporting manager is not found.
	ErrReportingManagerNotFound = errors.New("reporting manager not found")
)
