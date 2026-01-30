// Package timetable provides timetable management functionality.
package timetable

import "errors"

// Domain errors for the timetable module.
var (
	// Shift errors
	ErrShiftNotFound   = errors.New("shift not found")
	ErrShiftCodeExists = errors.New("shift code already exists")
	ErrShiftInUse      = errors.New("cannot delete shift that is in use")

	// Day Pattern errors
	ErrDayPatternNotFound   = errors.New("day pattern not found")
	ErrDayPatternCodeExists = errors.New("day pattern code already exists")
	ErrDayPatternInUse      = errors.New("cannot delete day pattern that is in use")

	// Period Slot errors
	ErrPeriodSlotNotFound = errors.New("period slot not found")
	ErrPeriodSlotOverlap  = errors.New("period slot overlaps with existing slot")

	// Timetable errors
	ErrTimetableNotFound         = errors.New("timetable not found")
	ErrTimetableAlreadyPublished = errors.New("timetable is already published")
	ErrTimetableNotDraft         = errors.New("only draft timetables can be modified")
	ErrTimetableHasConflicts     = errors.New("timetable has teacher conflicts")
	ErrPublishedTimetableExists  = errors.New("a published timetable already exists for this section")

	// Timetable Entry errors
	ErrTimetableEntryNotFound = errors.New("timetable entry not found")
	ErrTeacherConflict        = errors.New("teacher is already assigned to another class at this time")

	// Substitution errors
	ErrSubstitutionNotFound     = errors.New("substitution not found")
	ErrSubstitutionConflict     = errors.New("substitution already exists for this teacher and period")
	ErrSubstituteConflict       = errors.New("substitute teacher has a conflict at this time")
	ErrSubstitutionNotPending   = errors.New("only pending substitutions can be modified")
	ErrSubstitutionNotCancellable = errors.New("only pending or confirmed substitutions can be cancelled")

	// General errors
	ErrInvalidTimeRange = errors.New("end time must be after start time")
)
