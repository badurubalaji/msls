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

	// General errors
	ErrInvalidTimeRange = errors.New("end time must be after start time")
)
