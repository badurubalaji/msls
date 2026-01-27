package behavioral

import "errors"

var (
	// ErrIncidentNotFound is returned when a behavioral incident is not found
	ErrIncidentNotFound = errors.New("behavioral incident not found")

	// ErrFollowUpNotFound is returned when a follow-up is not found
	ErrFollowUpNotFound = errors.New("follow-up not found")

	// ErrInvalidIncidentDate is returned when the incident date is invalid
	ErrInvalidIncidentDate = errors.New("invalid incident date format, expected YYYY-MM-DD")

	// ErrInvalidScheduledDate is returned when the scheduled date is invalid
	ErrInvalidScheduledDate = errors.New("invalid scheduled date format, expected YYYY-MM-DD")

	// ErrUnauthorized is returned when user doesn't have permission
	ErrUnauthorized = errors.New("unauthorized to perform this action")
)
