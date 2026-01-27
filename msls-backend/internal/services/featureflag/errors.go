// Package featureflag provides feature flag management services.
package featureflag

import "errors"

// Service errors.
var (
	// ErrFlagNotFound is returned when a feature flag is not found.
	ErrFlagNotFound = errors.New("feature flag not found")

	// ErrFlagKeyExists is returned when trying to create a flag with an existing key.
	ErrFlagKeyExists = errors.New("feature flag key already exists")

	// ErrTenantOverrideNotFound is returned when a tenant override is not found.
	ErrTenantOverrideNotFound = errors.New("tenant feature flag override not found")

	// ErrUserOverrideNotFound is returned when a user override is not found.
	ErrUserOverrideNotFound = errors.New("user feature flag override not found")

	// ErrInvalidFlagKey is returned when the flag key format is invalid.
	ErrInvalidFlagKey = errors.New("invalid feature flag key format")

	// ErrCannotDeleteSystemFlag is returned when trying to delete a system flag.
	ErrCannotDeleteSystemFlag = errors.New("cannot delete system feature flag")
)
