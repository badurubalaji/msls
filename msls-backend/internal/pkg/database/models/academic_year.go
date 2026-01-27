// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AcademicYear represents a school academic year.
type AcademicYear struct {
	BaseModel
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	BranchID  *uuid.UUID `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	Name      string     `gorm:"type:varchar(50);not null" json:"name"`
	StartDate time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time  `gorm:"type:date;not null" json:"end_date"`
	IsCurrent bool       `gorm:"not null;default:false" json:"is_current"`
	IsActive  bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Tenant   Tenant          `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Branch   *Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Terms    []AcademicTerm  `gorm:"foreignKey:AcademicYearID" json:"terms,omitempty"`
	Holidays []Holiday       `gorm:"foreignKey:AcademicYearID" json:"holidays,omitempty"`
}

// TableName returns the table name for the AcademicYear model.
func (AcademicYear) TableName() string {
	return "academic_years"
}

// BeforeCreate hook for AcademicYear.
func (a *AcademicYear) BeforeCreate(tx *gorm.DB) error {
	if err := a.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// Validate performs validation on the AcademicYear model.
func (a *AcademicYear) Validate() error {
	if a.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if a.Name == "" {
		return ErrAcademicYearNameRequired
	}
	if a.StartDate.IsZero() {
		return ErrAcademicYearStartDateRequired
	}
	if a.EndDate.IsZero() {
		return ErrAcademicYearEndDateRequired
	}
	if !a.EndDate.After(a.StartDate) {
		return ErrAcademicYearInvalidDates
	}
	return nil
}

// AcademicTerm represents a term or semester within an academic year.
type AcademicTerm struct {
	BaseModel
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	StartDate      time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate        time.Time `gorm:"type:date;not null" json:"end_date"`
	Sequence       int       `gorm:"not null;default:1" json:"sequence"`

	// Relationships
	Tenant       Tenant       `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	AcademicYear AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
}

// TableName returns the table name for the AcademicTerm model.
func (AcademicTerm) TableName() string {
	return "academic_terms"
}

// BeforeCreate hook for AcademicTerm.
func (t *AcademicTerm) BeforeCreate(tx *gorm.DB) error {
	if err := t.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// Validate performs validation on the AcademicTerm model.
func (t *AcademicTerm) Validate() error {
	if t.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if t.AcademicYearID == uuid.Nil {
		return ErrAcademicYearIDRequired
	}
	if t.Name == "" {
		return ErrAcademicTermNameRequired
	}
	if t.StartDate.IsZero() {
		return ErrAcademicTermStartDateRequired
	}
	if t.EndDate.IsZero() {
		return ErrAcademicTermEndDateRequired
	}
	if !t.EndDate.After(t.StartDate) {
		return ErrAcademicTermInvalidDates
	}
	if t.Sequence < 1 {
		return ErrAcademicTermInvalidSequence
	}
	return nil
}

// HolidayType represents the type of holiday.
type HolidayType string

// HolidayType constants.
const (
	HolidayTypePublic   HolidayType = "public"
	HolidayTypeReligious HolidayType = "religious"
	HolidayTypeNational HolidayType = "national"
	HolidayTypeSchool   HolidayType = "school"
	HolidayTypeOther    HolidayType = "other"
)

// IsValid checks if the holiday type is valid.
func (h HolidayType) IsValid() bool {
	switch h {
	case HolidayTypePublic, HolidayTypeReligious, HolidayTypeNational, HolidayTypeSchool, HolidayTypeOther:
		return true
	}
	return false
}

// String returns the string representation of the holiday type.
func (h HolidayType) String() string {
	return string(h)
}

// Holiday represents a holiday within an academic year.
type Holiday struct {
	BaseModel
	TenantID       uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AcademicYearID uuid.UUID   `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	BranchID       *uuid.UUID  `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	Name           string      `gorm:"type:varchar(200);not null" json:"name"`
	Date           time.Time   `gorm:"type:date;not null" json:"date"`
	Type           HolidayType `gorm:"type:varchar(50);not null;default:'public'" json:"type"`
	IsOptional     bool        `gorm:"not null;default:false" json:"is_optional"`

	// Relationships
	Tenant       Tenant        `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	AcademicYear AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	Branch       *Branch       `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// TableName returns the table name for the Holiday model.
func (Holiday) TableName() string {
	return "holidays"
}

// BeforeCreate hook for Holiday.
func (h *Holiday) BeforeCreate(tx *gorm.DB) error {
	if err := h.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	if h.Type == "" {
		h.Type = HolidayTypePublic
	}
	return nil
}

// Validate performs validation on the Holiday model.
func (h *Holiday) Validate() error {
	if h.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if h.AcademicYearID == uuid.Nil {
		return ErrAcademicYearIDRequired
	}
	if h.Name == "" {
		return ErrHolidayNameRequired
	}
	if h.Date.IsZero() {
		return ErrHolidayDateRequired
	}
	if !h.Type.IsValid() {
		return ErrHolidayInvalidType
	}
	return nil
}
