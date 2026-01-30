// Package examination provides examination scheduling functionality.
package examination

import "errors"

var (
	// ErrExaminationNotFound is returned when an examination is not found.
	ErrExaminationNotFound = errors.New("examination not found")

	// ErrInvalidDateRange is returned when end date is before start date.
	ErrInvalidDateRange = errors.New("end date must be on or after start date")

	// ErrInvalidStatus is returned when status is invalid.
	ErrInvalidStatus = errors.New("invalid examination status")

	// ErrInvalidStatusTransition is returned when status transition is not allowed.
	ErrInvalidStatusTransition = errors.New("cannot transition examination to the requested status")

	// ErrCannotDeletePublished is returned when trying to delete a non-draft examination.
	ErrCannotDeletePublished = errors.New("cannot delete a published or ongoing examination")

	// ErrCannotUpdatePublished is returned when trying to update a non-draft examination.
	ErrCannotUpdatePublished = errors.New("cannot update a published examination; change status to draft first")

	// ErrNoSchedulesForPublish is returned when trying to publish an exam without schedules.
	ErrNoSchedulesForPublish = errors.New("cannot publish examination without any schedules")

	// ErrScheduleNotFound is returned when an exam schedule is not found.
	ErrScheduleNotFound = errors.New("exam schedule not found")

	// ErrScheduleConflict is returned when a schedule conflicts with existing schedules.
	ErrScheduleConflict = errors.New("schedule conflicts with an existing schedule")

	// ErrSubjectAlreadyScheduled is returned when a subject is already scheduled.
	ErrSubjectAlreadyScheduled = errors.New("subject is already scheduled for this examination")

	// ErrScheduleOutsideExamDates is returned when schedule date is outside exam date range.
	ErrScheduleOutsideExamDates = errors.New("schedule date must be within examination date range")

	// ErrInvalidTimeRange is returned when end time is before start time.
	ErrInvalidTimeRange = errors.New("end time must be after start time")

	// ErrInvalidPassingMarks is returned when passing marks exceeds max marks.
	ErrInvalidPassingMarks = errors.New("passing marks must be less than or equal to max marks")

	// ErrNoClassesSpecified is returned when no classes are specified for examination.
	ErrNoClassesSpecified = errors.New("at least one class must be specified for the examination")

	// ErrExamTypeNotActive is returned when selected exam type is not active.
	ErrExamTypeNotActive = errors.New("the selected exam type is not active")
)
