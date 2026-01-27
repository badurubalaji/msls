package behavioral

import (
	"time"

	"github.com/google/uuid"
	"msls-backend/internal/pkg/database/models"
)

// =============================================================================
// Request DTOs
// =============================================================================

// CreateIncidentRequest represents a request to create a behavioral incident
type CreateIncidentRequest struct {
	IncidentType          models.BehavioralIncidentType     `json:"incidentType" binding:"required,oneof=positive_recognition minor_infraction major_violation"`
	Severity              models.BehavioralSeverity `json:"severity" binding:"omitempty,oneof=low medium high critical"`
	IncidentDate          string                  `json:"incidentDate" binding:"required"`
	IncidentTime          string                  `json:"incidentTime" binding:"required"`
	Location              *string                 `json:"location"`
	Description           string                  `json:"description" binding:"required"`
	Witnesses             []string                `json:"witnesses"`
	StudentResponse       *string                 `json:"studentResponse"`
	ActionTaken           string                  `json:"actionTaken" binding:"required"`
	ParentMeetingRequired bool                    `json:"parentMeetingRequired"`
}

// UpdateIncidentRequest represents a request to update a behavioral incident
type UpdateIncidentRequest struct {
	IncidentType          *models.BehavioralIncidentType     `json:"incidentType" binding:"omitempty,oneof=positive_recognition minor_infraction major_violation"`
	Severity              *models.BehavioralSeverity `json:"severity" binding:"omitempty,oneof=low medium high critical"`
	IncidentDate          *string                  `json:"incidentDate"`
	IncidentTime          *string                  `json:"incidentTime"`
	Location              *string                  `json:"location"`
	Description           *string                  `json:"description"`
	Witnesses             []string                 `json:"witnesses"`
	StudentResponse       *string                  `json:"studentResponse"`
	ActionTaken           *string                  `json:"actionTaken"`
	ParentMeetingRequired *bool                    `json:"parentMeetingRequired"`
	ParentNotified        *bool                    `json:"parentNotified"`
}

// CreateFollowUpRequest represents a request to create a follow-up
type CreateFollowUpRequest struct {
	ScheduledDate    string        `json:"scheduledDate" binding:"required"`
	ScheduledTime    *string       `json:"scheduledTime"`
	Participants     []Participant `json:"participants"`
	ExpectedOutcomes *string       `json:"expectedOutcomes"`
}

// UpdateFollowUpRequest represents a request to update a follow-up
type UpdateFollowUpRequest struct {
	ScheduledDate    *string                 `json:"scheduledDate"`
	ScheduledTime    *string                 `json:"scheduledTime"`
	Participants     []Participant           `json:"participants"`
	ExpectedOutcomes *string                 `json:"expectedOutcomes"`
	MeetingNotes     *string                 `json:"meetingNotes"`
	ActualOutcomes   *string                 `json:"actualOutcomes"`
	Status           *models.FollowUpStatus  `json:"status" binding:"omitempty,oneof=pending completed cancelled"`
}

// Participant represents a meeting participant
type Participant struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// IncidentResponse represents a behavioral incident in the response
type IncidentResponse struct {
	ID                    uuid.UUID               `json:"id"`
	StudentID             uuid.UUID               `json:"studentId"`
	IncidentType          models.BehavioralIncidentType     `json:"incidentType"`
	IncidentTypeLabel     string                  `json:"incidentTypeLabel"`
	Severity              models.BehavioralSeverity `json:"severity"`
	SeverityLabel         string                  `json:"severityLabel"`
	IncidentDate          string                  `json:"incidentDate"`
	IncidentTime          string                  `json:"incidentTime"`
	Location              *string                 `json:"location,omitempty"`
	Description           string                  `json:"description"`
	Witnesses             []string                `json:"witnesses,omitempty"`
	StudentResponse       *string                 `json:"studentResponse,omitempty"`
	ActionTaken           string                  `json:"actionTaken"`
	ParentMeetingRequired bool                    `json:"parentMeetingRequired"`
	ParentNotified        bool                    `json:"parentNotified"`
	ParentNotifiedAt      *time.Time              `json:"parentNotifiedAt,omitempty"`
	ReportedBy            uuid.UUID               `json:"reportedBy"`
	ReporterName          string                  `json:"reporterName,omitempty"`
	FollowUps             []FollowUpResponse      `json:"followUps,omitempty"`
	CreatedAt             time.Time               `json:"createdAt"`
	UpdatedAt             time.Time               `json:"updatedAt"`
}

// FollowUpResponse represents a follow-up in the response
type FollowUpResponse struct {
	ID               uuid.UUID              `json:"id"`
	IncidentID       uuid.UUID              `json:"incidentId"`
	ScheduledDate    string                 `json:"scheduledDate"`
	ScheduledTime    *string                `json:"scheduledTime,omitempty"`
	Participants     []Participant          `json:"participants,omitempty"`
	ExpectedOutcomes *string                `json:"expectedOutcomes,omitempty"`
	MeetingNotes     *string                `json:"meetingNotes,omitempty"`
	ActualOutcomes   *string                `json:"actualOutcomes,omitempty"`
	Status           models.FollowUpStatus  `json:"status"`
	StatusLabel      string                 `json:"statusLabel"`
	CompletedAt      *time.Time             `json:"completedAt,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
}

// IncidentListResponse represents a list of incidents
type IncidentListResponse struct {
	Incidents []IncidentResponse `json:"incidents"`
	Total     int                `json:"total"`
}

// BehaviorSummary represents behavioral statistics for a student
type BehaviorSummary struct {
	TotalIncidents       int    `json:"totalIncidents"`
	PositiveCount        int    `json:"positiveCount"`
	MinorInfractionCount int    `json:"minorInfractionCount"`
	MajorViolationCount  int    `json:"majorViolationCount"`
	ThisMonthCount       int    `json:"thisMonthCount"`
	LastMonthCount       int    `json:"lastMonthCount"`
	Trend                string `json:"trend"` // improving, declining, stable
	PendingFollowUps     int    `json:"pendingFollowUps"`
}

// PendingFollowUpsResponse represents a list of pending follow-ups
type PendingFollowUpsResponse struct {
	FollowUps []PendingFollowUpItem `json:"followUps"`
	Total     int                   `json:"total"`
}

// PendingFollowUpItem represents a pending follow-up with student info
type PendingFollowUpItem struct {
	FollowUpResponse
	StudentID     uuid.UUID `json:"studentId"`
	StudentName   string    `json:"studentName"`
	IncidentType  string    `json:"incidentType"`
	IncidentDate  string    `json:"incidentDate"`
}

// =============================================================================
// Filter DTOs
// =============================================================================

// IncidentFilter represents filters for listing incidents
type IncidentFilter struct {
	IncidentType *models.BehavioralIncidentType     `form:"type"`
	Severity     *models.BehavioralSeverity `form:"severity"`
	DateFrom     *string                  `form:"dateFrom"`
	DateTo       *string                  `form:"dateTo"`
	Limit        int                      `form:"limit"`
	Offset       int                      `form:"offset"`
}

// =============================================================================
// Helper functions
// =============================================================================

// GetIncidentTypeLabel returns a human-readable label for incident type
func GetIncidentTypeLabel(t models.BehavioralIncidentType) string {
	switch t {
	case models.BehavioralIncidentTypePositiveRecognition:
		return "Positive Recognition"
	case models.BehavioralIncidentTypeMinorInfraction:
		return "Minor Infraction"
	case models.BehavioralIncidentTypeMajorViolation:
		return "Major Violation"
	default:
		return string(t)
	}
}

// GetSeverityLabel returns a human-readable label for severity
func GetSeverityLabel(s models.BehavioralSeverity) string {
	switch s {
	case models.BehavioralSeverityLow:
		return "Low"
	case models.BehavioralSeverityMedium:
		return "Medium"
	case models.BehavioralSeverityHigh:
		return "High"
	case models.BehavioralSeverityCritical:
		return "Critical"
	default:
		return string(s)
	}
}

// GetFollowUpStatusLabel returns a human-readable label for follow-up status
func GetFollowUpStatusLabel(s models.FollowUpStatus) string {
	switch s {
	case models.FollowUpStatusPending:
		return "Pending"
	case models.FollowUpStatusCompleted:
		return "Completed"
	case models.FollowUpStatusCancelled:
		return "Cancelled"
	default:
		return string(s)
	}
}
