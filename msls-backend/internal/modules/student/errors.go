// Package student provides student management functionality.
package student

import "errors"

// Repository errors.
var (
	// ErrStudentNotFound is returned when a student is not found.
	ErrStudentNotFound = errors.New("student not found")

	// ErrDuplicateAdmissionNumber is returned when the admission number already exists.
	ErrDuplicateAdmissionNumber = errors.New("admission number already exists")

	// ErrOptimisticLockConflict is returned when there's a concurrent edit conflict.
	ErrOptimisticLockConflict = errors.New("concurrent modification detected, please refresh and try again")

	// ErrAddressNotFound is returned when an address is not found.
	ErrAddressNotFound = errors.New("address not found")
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

	// ErrBranchNotFound is returned when the branch is not found.
	ErrBranchNotFound = errors.New("branch not found")

	// ErrInvalidAddressType is returned when an invalid address type is provided.
	ErrInvalidAddressType = errors.New("invalid address type")

	// ErrInvalidDateOfBirth is returned when date of birth is in the future.
	ErrInvalidDateOfBirth = errors.New("date of birth cannot be in the future")
)
