// Package branch provides branch management services.
package branch

import "errors"

// Service-level errors for branch operations.
var (
	// ErrBranchNotFound is returned when a branch is not found.
	ErrBranchNotFound = errors.New("branch not found")

	// ErrBranchNameRequired is returned when branch name is missing.
	ErrBranchNameRequired = errors.New("branch name is required")

	// ErrBranchCodeRequired is returned when branch code is missing.
	ErrBranchCodeRequired = errors.New("branch code is required")

	// ErrBranchCodeExists is returned when branch code already exists for tenant.
	ErrBranchCodeExists = errors.New("branch code already exists")

	// ErrCannotDeletePrimaryBranch is returned when trying to delete the primary branch.
	ErrCannotDeletePrimaryBranch = errors.New("cannot delete primary branch")

	// ErrCannotDeactivatePrimaryBranch is returned when trying to deactivate the primary branch.
	ErrCannotDeactivatePrimaryBranch = errors.New("cannot deactivate primary branch")

	// ErrBranchHasDependencies is returned when branch has associated records.
	ErrBranchHasDependencies = errors.New("branch has associated records and cannot be deleted")

	// ErrInvalidTimezone is returned when an invalid timezone is provided.
	ErrInvalidTimezone = errors.New("invalid timezone")

	// ErrTenantIDRequired is returned when tenant ID is missing.
	ErrTenantIDRequired = errors.New("tenant ID is required")
)
