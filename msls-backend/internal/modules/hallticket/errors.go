// Package hallticket provides hall ticket generation and management.
package hallticket

import "errors"

// Error definitions for hall ticket operations.
var (
	// ErrHallTicketNotFound is returned when a hall ticket is not found.
	ErrHallTicketNotFound = errors.New("hall ticket not found")

	// ErrTemplateNotFound is returned when a hall ticket template is not found.
	ErrTemplateNotFound = errors.New("hall ticket template not found")

	// ErrExaminationNotScheduled is returned when trying to generate hall tickets for an unscheduled exam.
	ErrExaminationNotScheduled = errors.New("examination must be in scheduled status to generate hall tickets")

	// ErrHallTicketAlreadyExists is returned when a hall ticket already exists for a student-exam pair.
	ErrHallTicketAlreadyExists = errors.New("hall ticket already exists for this student and examination")

	// ErrNoStudentsInClass is returned when no students are found for the selected class.
	ErrNoStudentsInClass = errors.New("no students found in the selected class")

	// ErrNoClassesForExam is returned when the examination has no assigned classes.
	ErrNoClassesForExam = errors.New("examination has no assigned classes")

	// ErrInvalidRollNumberPattern is returned when the roll number pattern is invalid.
	ErrInvalidRollNumberPattern = errors.New("invalid roll number pattern")

	// ErrTenantIDRequired is returned when tenant ID is not provided.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrExaminationIDRequired is returned when examination ID is not provided.
	ErrExaminationIDRequired = errors.New("examination ID is required")

	// ErrStudentIDRequired is returned when student ID is not provided.
	ErrStudentIDRequired = errors.New("student ID is required")

	// ErrTemplateNameRequired is returned when template name is not provided.
	ErrTemplateNameRequired = errors.New("template name is required")
)
