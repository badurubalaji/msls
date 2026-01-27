// Package salary provides salary management functionality.
package salary

import "errors"

// Salary component errors.
var (
	ErrComponentNotFound    = errors.New("salary component not found")
	ErrDuplicateCode        = errors.New("component code already exists")
	ErrComponentInUse       = errors.New("component is in use")
	ErrInvalidComponentType = errors.New("invalid component type")
	ErrInvalidCalcType      = errors.New("invalid calculation type")
	ErrPercentageOfRequired = errors.New("percentage_of_id is required for percentage-based components")
)

// Salary structure errors.
var (
	ErrStructureNotFound     = errors.New("salary structure not found")
	ErrDuplicateStructureCode = errors.New("structure code already exists")
	ErrStructureInUse        = errors.New("structure is in use")
	ErrNoComponentsInStructure = errors.New("structure must have at least one component")
)

// Staff salary errors.
var (
	ErrStaffSalaryNotFound  = errors.New("staff salary not found")
	ErrSalaryAlreadyExists  = errors.New("active salary already exists for this staff")
	ErrInvalidEffectiveDate = errors.New("effective date cannot be in the past")
)
