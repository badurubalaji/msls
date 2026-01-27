// Package promotion provides student promotion and retention processing functionality.
package promotion

import (
	"time"

	"github.com/google/uuid"
)

// PromotionDecision represents the decision for a student's promotion.
type PromotionDecision string

const (
	DecisionPending  PromotionDecision = "pending"
	DecisionPromote  PromotionDecision = "promote"
	DecisionRetain   PromotionDecision = "retain"
	DecisionTransfer PromotionDecision = "transfer"
)

// BatchStatus represents the status of a promotion batch.
type BatchStatus string

const (
	BatchStatusDraft      BatchStatus = "draft"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusCancelled  BatchStatus = "cancelled"
)

// PromotionRule represents the promotion rules for a class.
type PromotionRule struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID              uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	ClassID               uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uniq_promotion_rule" json:"classId"`
	MinAttendancePct      *float64   `gorm:"type:decimal(5,2);default:75.00" json:"minAttendancePct"`
	MinOverallMarksPct    *float64   `gorm:"type:decimal(5,2);default:33.00" json:"minOverallMarksPct"`
	MinSubjectsPassed     int        `gorm:"default:0" json:"minSubjectsPassed"`
	AutoPromoteOnCriteria bool       `gorm:"default:true" json:"autoPromoteOnCriteria"`
	IsActive              bool       `gorm:"default:true" json:"isActive"`
	CreatedAt             time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy             *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy             *uuid.UUID `gorm:"type:uuid" json:"updatedBy,omitempty"`
}

// TableName returns the table name for the PromotionRule model.
func (PromotionRule) TableName() string {
	return "promotion_rules"
}

// PromotionBatch represents a group of students being promoted together.
type PromotionBatch struct {
	ID                 uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID           uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenantId"`
	FromAcademicYearID uuid.UUID   `gorm:"type:uuid;not null;index" json:"fromAcademicYearId"`
	ToAcademicYearID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"toAcademicYearId"`
	FromClassID        uuid.UUID   `gorm:"type:uuid;not null;index" json:"fromClassId"`
	FromSectionID      *uuid.UUID  `gorm:"type:uuid" json:"fromSectionId,omitempty"`
	ToClassID          *uuid.UUID  `gorm:"type:uuid" json:"toClassId,omitempty"`
	Status             BatchStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	TotalStudents      int         `gorm:"default:0" json:"totalStudents"`
	PromotedCount      int         `gorm:"default:0" json:"promotedCount"`
	RetainedCount      int         `gorm:"default:0" json:"retainedCount"`
	TransferredCount   int         `gorm:"default:0" json:"transferredCount"`
	ProcessedAt        *time.Time  `gorm:"type:timestamptz" json:"processedAt,omitempty"`
	ProcessedBy        *uuid.UUID  `gorm:"type:uuid" json:"processedBy,omitempty"`
	CancelledAt        *time.Time  `gorm:"type:timestamptz" json:"cancelledAt,omitempty"`
	CancelledBy        *uuid.UUID  `gorm:"type:uuid" json:"cancelledBy,omitempty"`
	CancellationReason string      `gorm:"type:text" json:"cancellationReason,omitempty"`
	Notes              string      `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt          time.Time   `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt          time.Time   `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy          *uuid.UUID  `gorm:"type:uuid" json:"createdBy,omitempty"`

	// Relationships (for preloading)
	FromAcademicYear *AcademicYearRef   `gorm:"foreignKey:FromAcademicYearID" json:"fromAcademicYear,omitempty"`
	ToAcademicYear   *AcademicYearRef   `gorm:"foreignKey:ToAcademicYearID" json:"toAcademicYear,omitempty"`
	Records          []PromotionRecord  `gorm:"foreignKey:BatchID" json:"records,omitempty"`
}

// TableName returns the table name for the PromotionBatch model.
func (PromotionBatch) TableName() string {
	return "promotion_batches"
}

// PromotionRecord represents a single student's promotion record within a batch.
type PromotionRecord struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID         uuid.UUID         `gorm:"type:uuid;not null;index" json:"tenantId"`
	BatchID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"batchId"`
	StudentID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"studentId"`
	FromEnrollmentID uuid.UUID         `gorm:"type:uuid;not null" json:"fromEnrollmentId"`
	ToEnrollmentID   *uuid.UUID        `gorm:"type:uuid" json:"toEnrollmentId,omitempty"`
	Decision         PromotionDecision `gorm:"type:varchar(20);not null;default:'pending'" json:"decision"`
	ToClassID        *uuid.UUID        `gorm:"type:uuid" json:"toClassId,omitempty"`
	ToSectionID      *uuid.UUID        `gorm:"type:uuid" json:"toSectionId,omitempty"`
	RollNumber       string            `gorm:"type:varchar(20)" json:"rollNumber,omitempty"`
	AutoDecided      bool              `gorm:"default:false" json:"autoDecided"`
	DecisionReason   string            `gorm:"type:text" json:"decisionReason,omitempty"`
	// Performance metrics
	AttendancePct    *float64   `gorm:"type:decimal(5,2)" json:"attendancePct,omitempty"`
	OverallMarksPct  *float64   `gorm:"type:decimal(5,2)" json:"overallMarksPct,omitempty"`
	SubjectsPassed   *int       `gorm:"type:int" json:"subjectsPassed,omitempty"`
	// Override tracking
	OverrideBy       *uuid.UUID `gorm:"type:uuid" json:"overrideBy,omitempty"`
	OverrideAt       *time.Time `gorm:"type:timestamptz" json:"overrideAt,omitempty"`
	OverrideReason   string     `gorm:"type:text" json:"overrideReason,omitempty"`
	// Status-specific fields
	RetentionReason     string `gorm:"type:text" json:"retentionReason,omitempty"`
	TransferDestination string `gorm:"type:text" json:"transferDestination,omitempty"`
	CreatedAt           time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships (for preloading)
	Student *StudentRef `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName returns the table name for the PromotionRecord model.
func (PromotionRecord) TableName() string {
	return "promotion_records"
}

// AcademicYearRef is a reference model for academic years (from another module).
type AcademicYearRef struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50)" json:"name"`
	StartDate time.Time `gorm:"type:date" json:"startDate"`
	EndDate   time.Time `gorm:"type:date" json:"endDate"`
	IsCurrent bool      `gorm:"default:false" json:"isCurrent"`
}

// TableName returns the table name for the AcademicYearRef model.
func (AcademicYearRef) TableName() string {
	return "academic_years"
}

// StudentRef is a reference model for students (from another module).
type StudentRef struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AdmissionNumber string    `gorm:"type:varchar(50)" json:"admissionNumber"`
	FirstName       string    `gorm:"type:varchar(100)" json:"firstName"`
	LastName        string    `gorm:"type:varchar(100)" json:"lastName"`
	PhotoURL        string    `gorm:"type:varchar(500)" json:"photoUrl,omitempty"`
}

// TableName returns the table name for the StudentRef model.
func (StudentRef) TableName() string {
	return "students"
}

// FullName returns the full name of the student.
func (s *StudentRef) FullName() string {
	if s == nil {
		return ""
	}
	return s.FirstName + " " + s.LastName
}
