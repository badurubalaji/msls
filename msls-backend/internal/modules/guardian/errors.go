// Package guardian provides guardian and emergency contact management functionality.
package guardian

import "errors"

// Guardian-related errors.
var (
	// ErrGuardianNotFound is returned when a guardian is not found.
	ErrGuardianNotFound = errors.New("guardian not found")

	// ErrStudentNotFound is returned when the student is not found.
	ErrStudentNotFound = errors.New("student not found")

	// ErrEmergencyContactNotFound is returned when an emergency contact is not found.
	ErrEmergencyContactNotFound = errors.New("emergency contact not found")

	// ErrPrimaryGuardianExists is returned when trying to set a second primary guardian.
	ErrPrimaryGuardianExists = errors.New("student already has a primary guardian")

	// ErrPriorityConflict is returned when emergency contact priority conflicts.
	ErrPriorityConflict = errors.New("emergency contact priority already exists for this student")

	// ErrInvalidRelation is returned when guardian relation is invalid.
	ErrInvalidRelation = errors.New("invalid guardian relation")

	// ErrInvalidPriority is returned when emergency contact priority is invalid.
	ErrInvalidPriority = errors.New("priority must be between 1 and 5")

	// ErrFirstNameRequired is returned when first name is not provided.
	ErrFirstNameRequired = errors.New("first name is required")

	// ErrLastNameRequired is returned when last name is not provided.
	ErrLastNameRequired = errors.New("last name is required")

	// ErrPhoneRequired is returned when phone is not provided.
	ErrPhoneRequired = errors.New("phone number is required")

	// ErrNameRequired is returned when name is not provided for emergency contact.
	ErrNameRequired = errors.New("name is required")

	// ErrRelationRequired is returned when relation is not provided.
	ErrRelationRequired = errors.New("relation is required")
)
