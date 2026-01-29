// Package assignment provides teacher subject assignment functionality.
package assignment

import "errors"

// Assignment errors.
var (
	ErrAssignmentNotFound         = errors.New("assignment not found")
	ErrDuplicateAssignment        = errors.New("assignment already exists for this teacher-subject-class combination")
	ErrStaffNotFound              = errors.New("staff member not found")
	ErrSubjectNotFound            = errors.New("subject not found")
	ErrClassNotFound              = errors.New("class not found")
	ErrSectionNotFound            = errors.New("section not found")
	ErrAcademicYearNotFound       = errors.New("academic year not found")
	ErrTeacherOverAssigned        = errors.New("teacher exceeds maximum periods per week")
	ErrClassTeacherExists         = errors.New("class teacher already assigned for this class-section")
	ErrInvalidEffectiveDate       = errors.New("effective from date cannot be after effective to date")
	ErrAssignmentAlreadyInactive  = errors.New("assignment is already inactive")
	ErrWorkloadSettingsNotFound   = errors.New("workload settings not found")
	ErrInvalidStatus              = errors.New("invalid assignment status")
	ErrCannotModifyInactiveAssignment = errors.New("cannot modify inactive assignment")
)
