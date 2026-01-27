// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentStatus represents the status of a student enrollment.
type EnrollmentStatus string

// EnrollmentStatus constants.
const (
	EnrollmentStatusActive      EnrollmentStatus = "active"
	EnrollmentStatusCompleted   EnrollmentStatus = "completed"
	EnrollmentStatusTransferred EnrollmentStatus = "transferred"
	EnrollmentStatusDropout     EnrollmentStatus = "dropout"
)

// IsValid checks if the enrollment status is a valid value.
func (s EnrollmentStatus) IsValid() bool {
	switch s {
	case EnrollmentStatusActive, EnrollmentStatusCompleted, EnrollmentStatusTransferred, EnrollmentStatusDropout:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s EnrollmentStatus) String() string {
	return string(s)
}

// StudentEnrollment represents a student's enrollment for an academic year.
type StudentEnrollment struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenantId"`
	StudentID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"studentId"`
	AcademicYearID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"academicYearId"`
	ClassID         *uuid.UUID       `gorm:"type:uuid;index" json:"classId,omitempty"`
	SectionID       *uuid.UUID       `gorm:"type:uuid" json:"sectionId,omitempty"`
	RollNumber      string           `gorm:"type:varchar(20)" json:"rollNumber,omitempty"`
	ClassTeacherID  *uuid.UUID       `gorm:"type:uuid" json:"classTeacherId,omitempty"`
	Status          EnrollmentStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	EnrollmentDate  time.Time        `gorm:"type:date;not null;default:CURRENT_DATE" json:"enrollmentDate"`
	CompletionDate  *time.Time       `gorm:"type:date" json:"completionDate,omitempty"`
	TransferDate    *time.Time       `gorm:"type:date" json:"transferDate,omitempty"`
	TransferReason  string           `gorm:"type:text" json:"transferReason,omitempty"`
	DropoutDate     *time.Time       `gorm:"type:date" json:"dropoutDate,omitempty"`
	DropoutReason   string           `gorm:"type:text" json:"dropoutReason,omitempty"`
	Notes           string           `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time        `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy       *uuid.UUID       `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy       *uuid.UUID       `gorm:"type:uuid" json:"updatedBy,omitempty"`

	// Relationships (for preloading)
	AcademicYear *AcademicYearRef `gorm:"foreignKey:AcademicYearID" json:"academicYear,omitempty"`
}

// TableName returns the table name for the StudentEnrollment model.
func (StudentEnrollment) TableName() string {
	return "student_enrollments"
}

// BeforeCreate hook for StudentEnrollment.
func (e *StudentEnrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Status == "" {
		e.Status = EnrollmentStatusActive
	}
	if e.EnrollmentDate.IsZero() {
		e.EnrollmentDate = time.Now()
	}
	return e.Validate()
}

// Validate performs validation on the StudentEnrollment model.
func (e *StudentEnrollment) Validate() error {
	if e.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if e.StudentID == uuid.Nil {
		return ErrStudentIDRequired
	}
	if e.AcademicYearID == uuid.Nil {
		return ErrAcademicYearIDRequired
	}
	if e.Status != "" && !e.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

// EnrollmentStatusChange represents a change in enrollment status.
type EnrollmentStatusChange struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID         `gorm:"type:uuid;not null;index" json:"tenantId"`
	EnrollmentID uuid.UUID         `gorm:"type:uuid;not null;index" json:"enrollmentId"`
	FromStatus   *EnrollmentStatus `gorm:"type:varchar(20)" json:"fromStatus,omitempty"`
	ToStatus     EnrollmentStatus  `gorm:"type:varchar(20);not null" json:"toStatus"`
	ChangeReason string            `gorm:"type:text" json:"changeReason,omitempty"`
	ChangeDate   time.Time         `gorm:"type:date;not null;default:CURRENT_DATE" json:"changeDate"`
	ChangedAt    time.Time         `gorm:"not null;default:now()" json:"changedAt"`
	ChangedBy    uuid.UUID         `gorm:"type:uuid;not null" json:"changedBy"`
}

// TableName returns the table name for the EnrollmentStatusChange model.
func (EnrollmentStatusChange) TableName() string {
	return "enrollment_status_changes"
}

// BeforeCreate hook for EnrollmentStatusChange.
func (c *EnrollmentStatusChange) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ChangeDate.IsZero() {
		c.ChangeDate = time.Now()
	}
	return c.Validate()
}

// Validate performs validation on the EnrollmentStatusChange model.
func (c *EnrollmentStatusChange) Validate() error {
	if c.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if c.EnrollmentID == uuid.Nil {
		return ErrEnrollmentIDRequired
	}
	if !c.ToStatus.IsValid() {
		return ErrInvalidStatus
	}
	if c.ChangedBy == uuid.Nil {
		return ErrChangedByRequired
	}
	return nil
}

// AcademicYearRef is a lightweight reference to an academic year for preloading.
type AcademicYearRef struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50)" json:"name"`
	StartDate time.Time `gorm:"type:date" json:"startDate"`
	EndDate   time.Time `gorm:"type:date" json:"endDate"`
	IsCurrent bool      `json:"isCurrent"`
}

// TableName returns the table name for the AcademicYearRef model.
func (AcademicYearRef) TableName() string {
	return "academic_years"
}
