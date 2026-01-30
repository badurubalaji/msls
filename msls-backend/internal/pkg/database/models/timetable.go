// Package models provides database models for the MSLS application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// PeriodSlotType represents the type of a period slot.
type PeriodSlotType string

const (
	PeriodSlotTypeRegular    PeriodSlotType = "regular"
	PeriodSlotTypeShort      PeriodSlotType = "short"
	PeriodSlotTypeAssembly   PeriodSlotType = "assembly"
	PeriodSlotTypeBreak      PeriodSlotType = "break"
	PeriodSlotTypeLunch      PeriodSlotType = "lunch"
	PeriodSlotTypeActivity   PeriodSlotType = "activity"
	PeriodSlotTypeZeroPeriod PeriodSlotType = "zero_period"
)

// Shift represents a school shift (morning, afternoon, etc.).
type Shift struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	BranchID  uuid.UUID `gorm:"type:uuid;not null"`
	Branch    *Branch   `gorm:"foreignKey:BranchID"`

	Name        string `gorm:"type:varchar(50);not null"`
	Code        string `gorm:"type:varchar(20);not null"`
	StartTime   string `gorm:"type:time;not null"` // Using string for TIME type
	EndTime     string `gorm:"type:time;not null"`
	Description string `gorm:"type:text"`
	DisplayOrder int   `gorm:"not null;default:0"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
}

// TableName returns the table name for Shift.
func (Shift) TableName() string {
	return "shifts"
}

// DayPattern represents a day pattern (regular, half-day, etc.).
type DayPattern struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID uuid.UUID `gorm:"type:uuid;not null"`

	Name         string `gorm:"type:varchar(50);not null"`
	Code         string `gorm:"type:varchar(20);not null"`
	Description  string `gorm:"type:text"`
	TotalPeriods int    `gorm:"not null;default:8"`
	DisplayOrder int    `gorm:"not null;default:0"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relationships
	PeriodSlots []PeriodSlot `gorm:"foreignKey:DayPatternID"`
}

// TableName returns the table name for DayPattern.
func (DayPattern) TableName() string {
	return "day_patterns"
}

// DayPatternAssignment represents the assignment of a day pattern to a day of week.
type DayPatternAssignment struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID     uuid.UUID   `gorm:"type:uuid;not null"`
	BranchID     uuid.UUID   `gorm:"type:uuid;not null"`
	Branch       *Branch     `gorm:"foreignKey:BranchID"`

	DayOfWeek    int         `gorm:"not null"` // 0=Sunday, 1=Monday, ..., 6=Saturday
	DayPatternID *uuid.UUID  `gorm:"type:uuid"`
	DayPattern   *DayPattern `gorm:"foreignKey:DayPatternID"`
	IsWorkingDay bool        `gorm:"not null;default:true"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName returns the table name for DayPatternAssignment.
func (DayPatternAssignment) TableName() string {
	return "day_pattern_assignments"
}

// GetDayName returns the name of the day for the assignment.
func (d DayPatternAssignment) GetDayName() string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if d.DayOfWeek >= 0 && d.DayOfWeek < 7 {
		return days[d.DayOfWeek]
	}
	return ""
}

// PeriodSlot represents a period slot in the timetable.
type PeriodSlot struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	BranchID uuid.UUID `gorm:"type:uuid;not null"`
	Branch   *Branch   `gorm:"foreignKey:BranchID"`

	Name            string         `gorm:"type:varchar(50);not null"`
	PeriodNumber    *int           `gorm:"type:int"` // Null for breaks/assembly
	SlotType        PeriodSlotType `gorm:"type:varchar(20);not null;default:'regular'"`
	StartTime       string         `gorm:"type:time;not null"`
	EndTime         string         `gorm:"type:time;not null"`
	DurationMinutes int            `gorm:"not null"`

	DayPatternID *uuid.UUID  `gorm:"type:uuid"`
	DayPattern   *DayPattern `gorm:"foreignKey:DayPatternID"`

	ShiftID *uuid.UUID `gorm:"type:uuid"`
	Shift   *Shift     `gorm:"foreignKey:ShiftID"`

	DisplayOrder int  `gorm:"not null;default:0"`
	IsActive     bool `gorm:"not null;default:true"`

	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
}

// TableName returns the table name for PeriodSlot.
func (PeriodSlot) TableName() string {
	return "period_slots"
}

// IsTeachingPeriod returns true if the slot is a teaching period.
func (p PeriodSlot) IsTeachingPeriod() bool {
	return p.SlotType == PeriodSlotTypeRegular || p.SlotType == PeriodSlotTypeShort
}

// TimetableStatus represents the status of a timetable.
type TimetableStatus string

const (
	TimetableStatusDraft     TimetableStatus = "draft"
	TimetableStatusPublished TimetableStatus = "published"
	TimetableStatusArchived  TimetableStatus = "archived"
)

// Timetable represents a section's timetable.
type Timetable struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID       uuid.UUID       `gorm:"type:uuid;not null"`
	BranchID       uuid.UUID       `gorm:"type:uuid;not null"`
	Branch         *Branch         `gorm:"foreignKey:BranchID"`
	SectionID      uuid.UUID       `gorm:"type:uuid;not null"`
	Section        *Section        `gorm:"foreignKey:SectionID"`
	AcademicYearID uuid.UUID       `gorm:"type:uuid;not null"`
	AcademicYear   *AcademicYear   `gorm:"foreignKey:AcademicYearID"`

	Name        string          `gorm:"type:varchar(100);not null"`
	Description string          `gorm:"type:text"`
	Status      TimetableStatus `gorm:"type:varchar(20);not null;default:'draft'"`

	EffectiveFrom *time.Time `gorm:"type:date"`
	EffectiveTo   *time.Time `gorm:"type:date"`
	PublishedAt   *time.Time `gorm:"type:timestamptz"`
	PublishedBy   *uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	DeletedAt *time.Time `gorm:"type:timestamptz"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`
	Version   int        `gorm:"not null;default:1"`

	// Relationships
	Entries []TimetableEntry `gorm:"foreignKey:TimetableID"`
}

// TableName returns the table name for Timetable.
func (Timetable) TableName() string {
	return "timetables"
}

// TimetableEntry represents a single entry in a timetable (period assignment).
type TimetableEntry struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID     uuid.UUID   `gorm:"type:uuid;not null"`
	TimetableID  uuid.UUID   `gorm:"type:uuid;not null"`
	Timetable    *Timetable  `gorm:"foreignKey:TimetableID"`
	DayOfWeek    int         `gorm:"not null"` // 0=Sunday, 1=Monday, ..., 6=Saturday
	PeriodSlotID uuid.UUID   `gorm:"type:uuid;not null"`
	PeriodSlot   *PeriodSlot `gorm:"foreignKey:PeriodSlotID"`
	SubjectID    *uuid.UUID  `gorm:"type:uuid"`
	Subject      *Subject    `gorm:"foreignKey:SubjectID"`
	StaffID      *uuid.UUID  `gorm:"type:uuid"`
	Staff        *Staff      `gorm:"foreignKey:StaffID"`
	RoomNumber   string      `gorm:"type:varchar(50)"`
	Notes        string      `gorm:"type:text"`
	IsFreePeriod bool        `gorm:"not null;default:false"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName returns the table name for TimetableEntry.
func (TimetableEntry) TableName() string {
	return "timetable_entries"
}

// GetDayName returns the name of the day for the entry.
func (e TimetableEntry) GetDayName() string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if e.DayOfWeek >= 0 && e.DayOfWeek < 7 {
		return days[e.DayOfWeek]
	}
	return ""
}
