// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentAttendanceStatus represents the status of a student attendance record.
type StudentAttendanceStatus string

// StudentAttendanceStatus constants.
const (
	StudentAttendancePresent  StudentAttendanceStatus = "present"
	StudentAttendanceAbsent   StudentAttendanceStatus = "absent"
	StudentAttendanceLate     StudentAttendanceStatus = "late"
	StudentAttendanceHalfDay  StudentAttendanceStatus = "half_day"
)

// IsValid checks if the student attendance status is valid.
func (s StudentAttendanceStatus) IsValid() bool {
	switch s {
	case StudentAttendancePresent, StudentAttendanceAbsent, StudentAttendanceLate, StudentAttendanceHalfDay:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s StudentAttendanceStatus) String() string {
	return string(s)
}

// Label returns a human-readable label for the status.
func (s StudentAttendanceStatus) Label() string {
	switch s {
	case StudentAttendancePresent:
		return "Present"
	case StudentAttendanceAbsent:
		return "Absent"
	case StudentAttendanceLate:
		return "Late"
	case StudentAttendanceHalfDay:
		return "Half Day"
	default:
		return string(s)
	}
}

// ShortLabel returns a single character label for the status.
func (s StudentAttendanceStatus) ShortLabel() string {
	switch s {
	case StudentAttendancePresent:
		return "P"
	case StudentAttendanceAbsent:
		return "A"
	case StudentAttendanceLate:
		return "L"
	case StudentAttendanceHalfDay:
		return "H"
	default:
		return "?"
	}
}

// StudentAttendance represents a student attendance record.
type StudentAttendance struct {
	ID              uuid.UUID               `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID        uuid.UUID               `gorm:"type:uuid;not null;index" json:"tenantId"`
	StudentID       uuid.UUID               `gorm:"type:uuid;not null;index" json:"studentId"`
	SectionID       uuid.UUID               `gorm:"type:uuid;not null;index" json:"sectionId"`
	AttendanceDate  time.Time               `gorm:"type:date;not null;index" json:"attendanceDate"`
	Status          StudentAttendanceStatus `gorm:"type:varchar(20);not null;default:'present'" json:"status"`
	LateArrivalTime *time.Time              `gorm:"type:time" json:"lateArrivalTime,omitempty"`
	Remarks         string                  `gorm:"type:text" json:"remarks,omitempty"`
	MarkedBy        uuid.UUID               `gorm:"type:uuid;not null" json:"markedBy"`
	MarkedAt        time.Time               `gorm:"type:timestamptz;not null" json:"markedAt"`
	CreatedAt       time.Time               `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt       time.Time               `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Tenant   Tenant  `gorm:"foreignKey:TenantID" json:"-"`
	Student  Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Section  Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	MarkedByUser User `gorm:"foreignKey:MarkedBy" json:"markedByUser,omitempty"`
}

// TableName returns the table name for the StudentAttendance model.
func (StudentAttendance) TableName() string {
	return "student_attendance"
}

// BeforeCreate hook for StudentAttendance.
func (a *StudentAttendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Status == "" {
		a.Status = StudentAttendancePresent
	}
	if a.MarkedAt.IsZero() {
		a.MarkedAt = time.Now()
	}
	return nil
}

// StudentAttendanceSettings represents attendance settings per branch.
type StudentAttendanceSettings struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID             uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_student_attendance_settings_branch" json:"branchId"`
	EditWindowMinutes    int       `gorm:"not null;default:120" json:"editWindowMinutes"`
	LateThresholdMinutes int       `gorm:"not null;default:15" json:"lateThresholdMinutes"`
	SMSOnAbsent          bool      `gorm:"not null;default:false" json:"smsOnAbsent"`
	CreatedAt            time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
	Branch Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// TableName returns the table name for the StudentAttendanceSettings model.
func (StudentAttendanceSettings) TableName() string {
	return "student_attendance_settings"
}

// BeforeCreate hook for StudentAttendanceSettings.
func (s *StudentAttendanceSettings) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// StudentAttendanceWithHistory represents a student with their recent attendance history.
type StudentAttendanceWithHistory struct {
	Student        Student                   `json:"student"`
	TodayStatus    *StudentAttendanceStatus  `json:"todayStatus,omitempty"`
	TodayRemarks   string                    `json:"todayRemarks,omitempty"`
	Last5Days      []StudentAttendanceStatus `json:"last5Days"`
}
