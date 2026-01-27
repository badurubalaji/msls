// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttendanceStatus represents the status of an attendance record.
type AttendanceStatus string

// AttendanceStatus constants.
const (
	AttendanceStatusPresent  AttendanceStatus = "present"
	AttendanceStatusAbsent   AttendanceStatus = "absent"
	AttendanceStatusHalfDay  AttendanceStatus = "half_day"
	AttendanceStatusOnLeave  AttendanceStatus = "on_leave"
	AttendanceStatusHoliday  AttendanceStatus = "holiday"
)

// IsValid checks if the attendance status is valid.
func (s AttendanceStatus) IsValid() bool {
	switch s {
	case AttendanceStatusPresent, AttendanceStatusAbsent, AttendanceStatusHalfDay, AttendanceStatusOnLeave, AttendanceStatusHoliday:
		return true
	}
	return false
}

// HalfDayType represents the type of half day.
type HalfDayType string

// HalfDayType constants.
const (
	HalfDayFirstHalf  HalfDayType = "first_half"
	HalfDaySecondHalf HalfDayType = "second_half"
)

// IsValid checks if the half day type is valid.
func (h HalfDayType) IsValid() bool {
	switch h {
	case HalfDayFirstHalf, HalfDaySecondHalf:
		return true
	}
	return false
}

// RegularizationStatus represents the status of a regularization request.
type RegularizationStatus string

// RegularizationStatus constants.
const (
	RegularizationStatusPending  RegularizationStatus = "pending"
	RegularizationStatusApproved RegularizationStatus = "approved"
	RegularizationStatusRejected RegularizationStatus = "rejected"
)

// IsValid checks if the regularization status is valid.
func (s RegularizationStatus) IsValid() bool {
	switch s {
	case RegularizationStatusPending, RegularizationStatusApproved, RegularizationStatusRejected:
		return true
	}
	return false
}

// StaffAttendance represents a staff attendance record.
type StaffAttendance struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenantId"`
	StaffID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"staffId"`
	AttendanceDate time.Time        `gorm:"type:date;not null;index" json:"attendanceDate"`

	Status       AttendanceStatus `gorm:"type:varchar(20);not null;default:'present'" json:"status"`
	CheckInTime  *time.Time       `gorm:"type:timestamptz" json:"checkInTime,omitempty"`
	CheckOutTime *time.Time       `gorm:"type:timestamptz" json:"checkOutTime,omitempty"`

	IsLate      bool `gorm:"not null;default:false" json:"isLate"`
	LateMinutes int  `gorm:"default:0" json:"lateMinutes"`

	HalfDayType HalfDayType `gorm:"type:varchar(20)" json:"halfDayType,omitempty"`
	Remarks     string      `gorm:"type:text" json:"remarks,omitempty"`

	MarkedBy *uuid.UUID `gorm:"type:uuid" json:"markedBy,omitempty"`
	MarkedAt time.Time  `gorm:"not null;default:now()" json:"markedAt"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
	Staff  Staff  `gorm:"foreignKey:StaffID" json:"-"`
}

// TableName returns the table name for the StaffAttendance model.
func (StaffAttendance) TableName() string {
	return "staff_attendance"
}

// BeforeCreate hook for StaffAttendance.
func (a *StaffAttendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Status == "" {
		a.Status = AttendanceStatusPresent
	}
	a.MarkedAt = time.Now()
	return nil
}

// StaffAttendanceRegularization represents a regularization request.
type StaffAttendanceRegularization struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	StaffID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"staffId"`
	AttendanceID *uuid.UUID `gorm:"type:uuid" json:"attendanceId,omitempty"`

	RequestDate           time.Time        `gorm:"type:date;not null;index" json:"requestDate"`
	RequestedStatus       AttendanceStatus `gorm:"type:varchar(20);not null" json:"requestedStatus"`
	Reason                string           `gorm:"type:text;not null" json:"reason"`
	SupportingDocumentURL string           `gorm:"type:text" json:"supportingDocumentUrl,omitempty"`

	Status          RegularizationStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	ReviewedBy      *uuid.UUID           `gorm:"type:uuid" json:"reviewedBy,omitempty"`
	ReviewedAt      *time.Time           `gorm:"type:timestamptz" json:"reviewedAt,omitempty"`
	RejectionReason string               `gorm:"type:text" json:"rejectionReason,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Tenant     Tenant           `gorm:"foreignKey:TenantID" json:"-"`
	Staff      Staff            `gorm:"foreignKey:StaffID" json:"-"`
	Attendance *StaffAttendance `gorm:"foreignKey:AttendanceID" json:"-"`
}

// TableName returns the table name for the StaffAttendanceRegularization model.
func (StaffAttendanceRegularization) TableName() string {
	return "staff_attendance_regularization"
}

// BeforeCreate hook for StaffAttendanceRegularization.
func (r *StaffAttendanceRegularization) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.Status == "" {
		r.Status = RegularizationStatusPending
	}
	return nil
}

// StaffAttendanceSettings represents attendance settings per branch.
type StaffAttendanceSettings struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_attendance_settings" json:"branchId"`

	WorkStartTime         time.Time `gorm:"type:time;not null;default:'09:00'" json:"workStartTime"`
	WorkEndTime           time.Time `gorm:"type:time;not null;default:'17:00'" json:"workEndTime"`
	LateThresholdMinutes  int       `gorm:"not null;default:15" json:"lateThresholdMinutes"`
	HalfDayThresholdHours float64   `gorm:"type:decimal(4,2);not null;default:4.0" json:"halfDayThresholdHours"`

	AllowSelfCheckout             bool `gorm:"not null;default:true" json:"allowSelfCheckout"`
	RequireRegularizationApproval bool `gorm:"not null;default:true" json:"requireRegularizationApproval"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
	Branch Branch `gorm:"foreignKey:BranchID" json:"-"`
}

// TableName returns the table name for the StaffAttendanceSettings model.
func (StaffAttendanceSettings) TableName() string {
	return "staff_attendance_settings"
}

// BeforeCreate hook for StaffAttendanceSettings.
func (s *StaffAttendanceSettings) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
