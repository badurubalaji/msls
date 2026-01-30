package models

import (
	"time"

	"github.com/google/uuid"
)

// ExamStatus represents the status of an examination
type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "draft"
	ExamStatusScheduled ExamStatus = "scheduled"
	ExamStatusOngoing   ExamStatus = "ongoing"
	ExamStatusCompleted ExamStatus = "completed"
	ExamStatusCancelled ExamStatus = "cancelled"
)

// IsValid checks if the status is valid
func (s ExamStatus) IsValid() bool {
	switch s {
	case ExamStatusDraft, ExamStatusScheduled, ExamStatusOngoing, ExamStatusCompleted, ExamStatusCancelled:
		return true
	}
	return false
}

// CanTransitionTo checks if status can transition to target status
func (s ExamStatus) CanTransitionTo(target ExamStatus) bool {
	switch s {
	case ExamStatusDraft:
		return target == ExamStatusScheduled || target == ExamStatusCancelled
	case ExamStatusScheduled:
		return target == ExamStatusOngoing || target == ExamStatusCancelled || target == ExamStatusDraft
	case ExamStatusOngoing:
		return target == ExamStatusCompleted || target == ExamStatusCancelled
	case ExamStatusCompleted:
		return false // Cannot transition from completed
	case ExamStatusCancelled:
		return false // Cannot transition from cancelled
	}
	return false
}

// Examination represents an examination event
type Examination struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID       uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenantId"`
	Name           string       `gorm:"size:200;not null" json:"name"`
	ExamTypeID     uuid.UUID    `gorm:"type:uuid;not null" json:"examTypeId"`
	AcademicYearID uuid.UUID    `gorm:"type:uuid;not null" json:"academicYearId"`
	StartDate      time.Time    `gorm:"type:date;not null" json:"startDate"`
	EndDate        time.Time    `gorm:"type:date;not null" json:"endDate"`
	Status         ExamStatus   `gorm:"size:20;not null;default:'draft'" json:"status"`
	Description    *string      `gorm:"type:text" json:"description,omitempty"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
	CreatedBy      *uuid.UUID   `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy      *uuid.UUID   `gorm:"type:uuid" json:"updatedBy,omitempty"`

	// Relationships
	ExamType     *ExamType      `gorm:"foreignKey:ExamTypeID" json:"examType,omitempty"`
	AcademicYear *AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academicYear,omitempty"`
	Classes      []Class        `gorm:"many2many:examination_classes;foreignKey:ID;joinForeignKey:ExaminationID;References:ID;joinReferences:ClassID" json:"classes,omitempty"`
	Schedules    []ExamSchedule `gorm:"foreignKey:ExaminationID" json:"schedules,omitempty"`
}

// TableName returns the table name for Examination
func (Examination) TableName() string {
	return "examinations"
}

// ExaminationClass represents the many-to-many relationship between examinations and classes
type ExaminationClass struct {
	ExaminationID uuid.UUID `gorm:"type:uuid;primaryKey" json:"examinationId"`
	ClassID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"classId"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// TableName returns the table name for ExaminationClass
func (ExaminationClass) TableName() string {
	return "examination_classes"
}

// ExamSchedule represents a single exam schedule entry
type ExamSchedule struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	ExaminationID uuid.UUID `gorm:"type:uuid;not null;index" json:"examinationId"`
	SubjectID     uuid.UUID `gorm:"type:uuid;not null" json:"subjectId"`
	ExamDate      time.Time `gorm:"type:date;not null" json:"examDate"`
	StartTime     string    `gorm:"type:time;not null" json:"startTime"` // Store as string for TIME type
	EndTime       string    `gorm:"type:time;not null" json:"endTime"`
	MaxMarks      int       `gorm:"not null;default:100" json:"maxMarks"`
	PassingMarks  *int      `json:"passingMarks,omitempty"`
	Venue         *string   `gorm:"size:100" json:"venue,omitempty"`
	Notes         *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Examination *Examination `gorm:"foreignKey:ExaminationID" json:"examination,omitempty"`
	Subject     *Subject     `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

// TableName returns the table name for ExamSchedule
func (ExamSchedule) TableName() string {
	return "exam_schedules"
}
