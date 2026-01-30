// Package exam provides examination management functionality.
package exam

import "errors"

var (
	// ErrExamTypeNotFound is returned when an exam type is not found.
	ErrExamTypeNotFound = errors.New("exam type not found")

	// ErrExamTypeCodeExists is returned when an exam type code already exists.
	ErrExamTypeCodeExists = errors.New("exam type code already exists")

	// ErrInvalidWeightage is returned when weightage is out of range.
	ErrInvalidWeightage = errors.New("weightage must be between 0 and 100")

	// ErrInvalidMaxMarks is returned when max marks is invalid.
	ErrInvalidMaxMarks = errors.New("maximum marks must be greater than 0")

	// ErrInvalidPassingMarks is returned when passing marks exceeds max marks.
	ErrInvalidPassingMarks = errors.New("passing marks cannot exceed maximum marks")

	// ErrExamTypeInUse is returned when trying to delete an exam type that is in use.
	ErrExamTypeInUse = errors.New("exam type is in use and cannot be deleted")
)
