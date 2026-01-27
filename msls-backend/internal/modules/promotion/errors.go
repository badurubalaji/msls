// Package promotion provides student promotion and retention processing functionality.
package promotion

import "errors"

// Domain errors for promotion module.
var (
	// ErrRuleNotFound is returned when a promotion rule is not found.
	ErrRuleNotFound = errors.New("promotion rule not found")

	// ErrBatchNotFound is returned when a promotion batch is not found.
	ErrBatchNotFound = errors.New("promotion batch not found")

	// ErrRecordNotFound is returned when a promotion record is not found.
	ErrRecordNotFound = errors.New("promotion record not found")

	// ErrBatchNotDraft is returned when trying to modify a non-draft batch.
	ErrBatchNotDraft = errors.New("batch is not in draft status")

	// ErrBatchNotProcessable is returned when batch cannot be processed.
	ErrBatchNotProcessable = errors.New("batch cannot be processed")

	// ErrPendingDecisions is returned when there are unresolved pending decisions.
	ErrPendingDecisions = errors.New("there are pending decisions that must be resolved")

	// ErrNoStudentsInClass is returned when no students are found for the class/section.
	ErrNoStudentsInClass = errors.New("no students found in the specified class/section")

	// ErrInvalidDecision is returned when an invalid decision is provided.
	ErrInvalidDecision = errors.New("invalid promotion decision")

	// ErrMissingTargetClass is returned when promoted student has no target class.
	ErrMissingTargetClass = errors.New("target class must be specified for promotion")

	// ErrDuplicateBatch is returned when a batch already exists for the same source.
	ErrDuplicateBatch = errors.New("a batch already exists for this source")

	// ErrBatchAlreadyProcessed is returned when trying to process an already processed batch.
	ErrBatchAlreadyProcessed = errors.New("batch has already been processed")

	// ErrBatchCancelled is returned when trying to operate on a cancelled batch.
	ErrBatchCancelled = errors.New("batch has been cancelled")

	// ErrInvalidAcademicYears is returned when from and to academic years are invalid.
	ErrInvalidAcademicYears = errors.New("target academic year must be after source academic year")

	// ErrSameAcademicYear is returned when from and to academic years are the same.
	ErrSameAcademicYear = errors.New("source and target academic years cannot be the same")
)
