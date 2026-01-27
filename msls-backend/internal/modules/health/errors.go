package health

import "errors"

var (
	// Health Profile errors
	ErrHealthProfileNotFound      = errors.New("health profile not found")
	ErrHealthProfileAlreadyExists = errors.New("health profile already exists for this student")

	// Allergy errors
	ErrAllergyNotFound     = errors.New("allergy not found")
	ErrInvalidAllergyType  = errors.New("invalid allergy type")
	ErrInvalidSeverity     = errors.New("invalid severity level")

	// Condition errors
	ErrConditionNotFound    = errors.New("chronic condition not found")
	ErrInvalidConditionType = errors.New("invalid condition type")

	// Medication errors
	ErrMedicationNotFound     = errors.New("medication not found")
	ErrInvalidFrequency       = errors.New("invalid medication frequency")
	ErrInvalidRoute           = errors.New("invalid medication route")
	ErrInvalidDateRange       = errors.New("end date must be after start date")

	// Vaccination errors
	ErrVaccinationNotFound  = errors.New("vaccination not found")
	ErrInvalidVaccineType   = errors.New("invalid vaccine type")
	ErrInvalidDoseNumber    = errors.New("dose number must be positive")

	// Incident errors
	ErrIncidentNotFound    = errors.New("medical incident not found")
	ErrInvalidIncidentType = errors.New("invalid incident type")

	// General errors
	ErrStudentNotFound = errors.New("student not found")
	ErrUnauthorized    = errors.New("unauthorized access to health records")
)
