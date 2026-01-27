package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BehavioralIncidentType represents the type of behavioral incident
type BehavioralIncidentType string

const (
	BehavioralIncidentTypePositiveRecognition BehavioralIncidentType = "positive_recognition"
	BehavioralIncidentTypeMinorInfraction     BehavioralIncidentType = "minor_infraction"
	BehavioralIncidentTypeMajorViolation      BehavioralIncidentType = "major_violation"
)

// BehavioralSeverity represents the severity level of a behavioral incident
type BehavioralSeverity string

const (
	BehavioralSeverityLow      BehavioralSeverity = "low"
	BehavioralSeverityMedium   BehavioralSeverity = "medium"
	BehavioralSeverityHigh     BehavioralSeverity = "high"
	BehavioralSeverityCritical BehavioralSeverity = "critical"
)

// FollowUpStatus represents the status of a follow-up
type FollowUpStatus string

const (
	FollowUpStatusPending   FollowUpStatus = "pending"
	FollowUpStatusCompleted FollowUpStatus = "completed"
	FollowUpStatusCancelled FollowUpStatus = "cancelled"
)

// StudentBehavioralIncident represents a behavioral incident record
type StudentBehavioralIncident struct {
	ID                    uuid.UUID              `gorm:"type:uuid;primaryKey"`
	TenantID              uuid.UUID              `gorm:"type:uuid;not null;index"`
	StudentID             uuid.UUID              `gorm:"type:uuid;not null;index"`
	IncidentType          BehavioralIncidentType `gorm:"type:varchar(30);not null"`
	Severity              BehavioralSeverity     `gorm:"type:varchar(20);not null;default:'medium'"`
	IncidentDate          time.Time        `gorm:"type:date;not null"`
	IncidentTime          string           `gorm:"type:time;not null"`
	Location              *string          `gorm:"type:varchar(200)"`
	Description           string           `gorm:"type:text;not null"`
	Witnesses             json.RawMessage  `gorm:"type:jsonb"`
	StudentResponse       *string          `gorm:"type:text"`
	ActionTaken           string           `gorm:"type:text;not null"`
	ParentMeetingRequired bool             `gorm:"not null;default:false"`
	ParentNotified        bool             `gorm:"not null;default:false"`
	ParentNotifiedAt      *time.Time       `gorm:"type:timestamptz"`
	ReportedBy            uuid.UUID        `gorm:"type:uuid;not null"`
	CreatedAt             time.Time        `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt             time.Time        `gorm:"type:timestamptz;not null;default:now()"`

	// Relations
	Student   *Student           `gorm:"foreignKey:StudentID"`
	Reporter  *User              `gorm:"foreignKey:ReportedBy"`
	FollowUps []IncidentFollowUp `gorm:"foreignKey:IncidentID"`
}

func (StudentBehavioralIncident) TableName() string {
	return "student_behavioral_incidents"
}

// IncidentFollowUp represents a scheduled follow-up for an incident
type IncidentFollowUp struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey"`
	TenantID         uuid.UUID       `gorm:"type:uuid;not null;index"`
	IncidentID       uuid.UUID       `gorm:"type:uuid;not null;index"`
	ScheduledDate    time.Time       `gorm:"type:date;not null"`
	ScheduledTime    *string         `gorm:"type:time"`
	Participants     json.RawMessage `gorm:"type:jsonb"`
	ExpectedOutcomes *string         `gorm:"type:text"`
	MeetingNotes     *string         `gorm:"type:text"`
	ActualOutcomes   *string         `gorm:"type:text"`
	Status           FollowUpStatus  `gorm:"type:varchar(20);not null;default:'pending'"`
	CompletedAt      *time.Time      `gorm:"type:timestamptz"`
	CompletedBy      *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt        time.Time       `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy        *uuid.UUID      `gorm:"type:uuid"`

	// Relations
	Incident *StudentBehavioralIncident `gorm:"foreignKey:IncidentID"`
}

func (IncidentFollowUp) TableName() string {
	return "incident_follow_ups"
}

// Participant represents a follow-up meeting participant
type Participant struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
