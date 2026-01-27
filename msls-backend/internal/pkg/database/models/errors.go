// Package models provides GORM model definitions for the MSLS database.
package models

import "errors"

// Model validation errors.
var (
	// Common errors
	ErrInvalidStatus    = errors.New("invalid status value")
	ErrTenantIDRequired = errors.New("tenant_id is required")

	// Tenant errors
	ErrTenantNameRequired = errors.New("tenant name is required")
	ErrTenantSlugRequired = errors.New("tenant slug is required")

	// Branch errors
	ErrBranchNameRequired = errors.New("branch name is required")
	ErrBranchCodeRequired = errors.New("branch code is required")

	// User errors
	ErrEmailOrPhoneRequired = errors.New("email or phone is required")
	ErrUserIDRequired       = errors.New("user_id is required")

	// Role errors
	ErrRoleNameRequired = errors.New("role name is required")

	// Token errors
	ErrTokenHashRequired = errors.New("token hash is required")
	ErrInvalidTokenType  = errors.New("invalid token type")

	// Feature flag errors
	ErrFeatureFlagKeyRequired  = errors.New("feature flag key is required")
	ErrFeatureFlagNameRequired = errors.New("feature flag name is required")
	ErrFeatureFlagIDRequired   = errors.New("feature flag id is required")

	// OTP errors
	ErrIdentifierRequired = errors.New("identifier is required")
	ErrCodeHashRequired   = errors.New("code hash is required")
	ErrInvalidOTPType     = errors.New("invalid OTP type")
	ErrInvalidOTPChannel  = errors.New("invalid OTP channel")

	// Academic Year errors
	ErrAcademicYearNameRequired      = errors.New("academic year name is required")
	ErrAcademicYearStartDateRequired = errors.New("academic year start date is required")
	ErrAcademicYearEndDateRequired   = errors.New("academic year end date is required")
	ErrAcademicYearInvalidDates      = errors.New("academic year end date must be after start date")
	ErrAcademicYearIDRequired        = errors.New("academic year id is required")

	// Academic Term errors
	ErrAcademicTermNameRequired      = errors.New("academic term name is required")
	ErrAcademicTermStartDateRequired = errors.New("academic term start date is required")
	ErrAcademicTermEndDateRequired   = errors.New("academic term end date is required")
	ErrAcademicTermInvalidDates      = errors.New("academic term end date must be after start date")
	ErrAcademicTermInvalidSequence   = errors.New("academic term sequence must be positive")

	// Holiday errors
	ErrHolidayNameRequired = errors.New("holiday name is required")
	ErrHolidayDateRequired = errors.New("holiday date is required")
	ErrHolidayInvalidType  = errors.New("invalid holiday type")
)
