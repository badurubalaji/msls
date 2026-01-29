// Package academic provides academic structure management functionality.
package academic

import "errors"

// Domain errors for the academic module.
var (
	// Class errors
	ErrClassNotFound      = errors.New("class not found")
	ErrClassCodeExists    = errors.New("class code already exists")
	ErrClassHasSections   = errors.New("cannot delete class with sections")
	ErrClassHasStudents   = errors.New("cannot delete class with enrolled students")

	// Section errors
	ErrSectionNotFound    = errors.New("section not found")
	ErrSectionCodeExists  = errors.New("section code already exists for this class")
	ErrSectionHasStudents = errors.New("cannot delete section with enrolled students")

	// Stream errors
	ErrStreamNotFound     = errors.New("stream not found")
	ErrStreamCodeExists   = errors.New("stream code already exists")
	ErrStreamInUse        = errors.New("cannot delete stream that is in use")

	// Subject errors
	ErrSubjectNotFound       = errors.New("subject not found")
	ErrSubjectCodeExists     = errors.New("subject code already exists")
	ErrSubjectInUse          = errors.New("cannot delete subject that is in use")

	// Class-Subject errors
	ErrClassSubjectNotFound  = errors.New("class-subject mapping not found")
	ErrClassSubjectExists    = errors.New("subject already assigned to this class")

	// General errors
	ErrInvalidBranch       = errors.New("invalid branch")
	ErrInvalidAcademicYear = errors.New("invalid academic year")
	ErrInvalidClassTeacher = errors.New("invalid class teacher")
)
