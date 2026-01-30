package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubstitutionStatus represents the status of a substitution.
type SubstitutionStatus string

const (
	SubstitutionStatusPending   SubstitutionStatus = "pending"
	SubstitutionStatusConfirmed SubstitutionStatus = "confirmed"
	SubstitutionStatusCompleted SubstitutionStatus = "completed"
	SubstitutionStatusCancelled SubstitutionStatus = "cancelled"
)

// Substitution represents a teacher substitution record.
type Substitution struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID          uuid.UUID          `gorm:"type:uuid;not null"`
	BranchID          uuid.UUID          `gorm:"type:uuid;not null"`
	OriginalStaffID   uuid.UUID          `gorm:"type:uuid;not null"`
	SubstituteStaffID uuid.UUID          `gorm:"type:uuid;not null"`
	SubstitutionDate  time.Time          `gorm:"type:date;not null"`
	Reason            string             `gorm:"type:varchar(255)"`
	Status            SubstitutionStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	Notes             string             `gorm:"type:text"`
	CreatedBy         *uuid.UUID         `gorm:"type:uuid"`
	ApprovedBy        *uuid.UUID         `gorm:"type:uuid"`
	ApprovedAt        *time.Time         `gorm:"type:timestamptz"`
	CreatedAt         time.Time          `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt         time.Time          `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt         gorm.DeletedAt     `gorm:"type:timestamptz;index"`

	// Relationships
	Tenant          *Tenant              `gorm:"foreignKey:TenantID"`
	Branch          *Branch              `gorm:"foreignKey:BranchID"`
	OriginalStaff   *Staff               `gorm:"foreignKey:OriginalStaffID"`
	SubstituteStaff *Staff               `gorm:"foreignKey:SubstituteStaffID"`
	Creator         *User                `gorm:"foreignKey:CreatedBy"`
	Approver        *User                `gorm:"foreignKey:ApprovedBy"`
	Periods         []SubstitutionPeriod `gorm:"foreignKey:SubstitutionID"`
}

// TableName returns the table name for Substitution.
func (Substitution) TableName() string {
	return "substitutions"
}

// SubstitutionPeriod represents a specific period covered by a substitution.
type SubstitutionPeriod struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubstitutionID   uuid.UUID  `gorm:"type:uuid;not null"`
	TimetableEntryID *uuid.UUID `gorm:"type:uuid"`
	PeriodSlotID     uuid.UUID  `gorm:"type:uuid;not null"`
	SubjectID        *uuid.UUID `gorm:"type:uuid"`
	SectionID        *uuid.UUID `gorm:"type:uuid"`
	RoomNumber       string     `gorm:"type:varchar(50)"`
	Notes            string     `gorm:"type:text"`
	CreatedAt        time.Time  `gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Substitution   *Substitution   `gorm:"foreignKey:SubstitutionID"`
	TimetableEntry *TimetableEntry `gorm:"foreignKey:TimetableEntryID"`
	PeriodSlot     *PeriodSlot     `gorm:"foreignKey:PeriodSlotID"`
	Subject        *Subject        `gorm:"foreignKey:SubjectID"`
	Section        *Section        `gorm:"foreignKey:SectionID"`
}

// TableName returns the table name for SubstitutionPeriod.
func (SubstitutionPeriod) TableName() string {
	return "substitution_periods"
}
