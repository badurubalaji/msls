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
// Supports both daily attendance (period_id = NULL) and period-wise attendance.
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

	// Period-wise attendance fields (Story 7.2)
	// NULL for daily attendance, set for period-wise attendance
	PeriodID         *uuid.UUID `gorm:"type:uuid;index" json:"periodId,omitempty"`
	TimetableEntryID *uuid.UUID `gorm:"type:uuid;index" json:"timetableEntryId,omitempty"`

	// Relationships
	Tenant         Tenant          `gorm:"foreignKey:TenantID" json:"-"`
	Student        Student         `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Section        Section         `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	MarkedByUser   User            `gorm:"foreignKey:MarkedBy" json:"markedByUser,omitempty"`
	PeriodSlot     *PeriodSlot     `gorm:"foreignKey:PeriodID" json:"periodSlot,omitempty"`
	TimetableEntry *TimetableEntry `gorm:"foreignKey:TimetableEntryID" json:"timetableEntry,omitempty"`
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
	ID                       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID                 uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID                 uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_student_attendance_settings_branch" json:"branchId"`
	EditWindowMinutes        int       `gorm:"not null;default:120" json:"editWindowMinutes"`
	LateThresholdMinutes     int       `gorm:"not null;default:15" json:"lateThresholdMinutes"`
	SMSOnAbsent              bool      `gorm:"not null;default:false" json:"smsOnAbsent"`
	PeriodAttendanceEnabled  bool      `gorm:"not null;default:false" json:"periodAttendanceEnabled"`
	CreatedAt                time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt                time.Time `gorm:"not null;default:now()" json:"updatedAt"`

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

// PeriodAttendanceRecord represents attendance for a specific period.
type PeriodAttendanceRecord struct {
	StudentID       uuid.UUID               `json:"studentId"`
	StudentName     string                  `json:"studentName"`
	RollNumber      string                  `json:"rollNumber,omitempty"`
	PhotoURL        string                  `json:"photoUrl,omitempty"`
	Status          StudentAttendanceStatus `json:"status"`
	LateArrivalTime *time.Time              `json:"lateArrivalTime,omitempty"`
	Remarks         string                  `json:"remarks,omitempty"`
}

// PeriodInfo represents information about a period for attendance.
type PeriodInfo struct {
	PeriodSlotID     uuid.UUID  `json:"periodSlotId"`
	TimetableEntryID uuid.UUID  `json:"timetableEntryId"`
	PeriodName       string     `json:"periodName"`
	PeriodNumber     *int       `json:"periodNumber,omitempty"`
	StartTime        string     `json:"startTime"`
	EndTime          string     `json:"endTime"`
	SubjectID        *uuid.UUID `json:"subjectId,omitempty"`
	SubjectName      string     `json:"subjectName,omitempty"`
	SubjectCode      string     `json:"subjectCode,omitempty"`
	IsMarked         bool       `json:"isMarked"`
}

// DailySummaryRecord represents a student's attendance across all periods for a day.
type DailySummaryRecord struct {
	StudentID       uuid.UUID                          `json:"studentId"`
	StudentName     string                             `json:"studentName"`
	RollNumber      string                             `json:"rollNumber,omitempty"`
	PhotoURL        string                             `json:"photoUrl,omitempty"`
	PeriodStatuses  map[uuid.UUID]StudentAttendanceStatus `json:"periodStatuses"` // periodId -> status
	TotalPeriods    int                                `json:"totalPeriods"`
	PeriodsPresent  int                                `json:"periodsPresent"`
	PeriodsAbsent   int                                `json:"periodsAbsent"`
	OverallStatus   StudentAttendanceStatus            `json:"overallStatus"` // Derived: present if >50% periods attended
	AttendancePercentage float64                       `json:"attendancePercentage"`
}

// SubjectAttendanceStats represents attendance statistics for a subject.
type SubjectAttendanceStats struct {
	SubjectID             uuid.UUID `json:"subjectId"`
	SubjectName           string    `json:"subjectName"`
	SubjectCode           string    `json:"subjectCode"`
	TotalPeriods          int       `json:"totalPeriods"`
	PeriodsPresent        int       `json:"periodsPresent"`
	PeriodsAbsent         int       `json:"periodsAbsent"`
	PeriodsLate           int       `json:"periodsLate"`
	AttendancePercentage  float64   `json:"attendancePercentage"`
	MinimumRequired       float64   `json:"minimumRequired"` // e.g., 75%
	IsEligible            bool      `json:"isEligible"`      // Met minimum attendance
}

// ============================================================================
// Attendance Audit (Story 7.3)
// ============================================================================

// AttendanceChangeType represents the type of attendance change.
type AttendanceChangeType string

// AttendanceChangeType constants.
const (
	AttendanceChangeCreate AttendanceChangeType = "create"
	AttendanceChangeEdit   AttendanceChangeType = "edit"
)

// StudentAttendanceAudit represents an audit record for attendance changes.
type StudentAttendanceAudit struct {
	ID                      uuid.UUID                `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID                uuid.UUID                `gorm:"type:uuid;not null;index" json:"tenantId"`
	AttendanceID            uuid.UUID                `gorm:"type:uuid;not null;index" json:"attendanceId"`
	PreviousStatus          *StudentAttendanceStatus `gorm:"type:varchar(20)" json:"previousStatus,omitempty"`
	NewStatus               StudentAttendanceStatus  `gorm:"type:varchar(20);not null" json:"newStatus"`
	PreviousRemarks         *string                  `gorm:"type:text" json:"previousRemarks,omitempty"`
	NewRemarks              *string                  `gorm:"type:text" json:"newRemarks,omitempty"`
	PreviousLateArrivalTime *time.Time               `gorm:"type:time" json:"previousLateArrivalTime,omitempty"`
	NewLateArrivalTime      *time.Time               `gorm:"type:time" json:"newLateArrivalTime,omitempty"`
	ChangeType              AttendanceChangeType     `gorm:"type:varchar(20);not null;default:'edit'" json:"changeType"`
	ChangeReason            string                   `gorm:"type:text;not null" json:"changeReason"`
	ChangedBy               uuid.UUID                `gorm:"type:uuid;not null;index" json:"changedBy"`
	ChangedAt               time.Time                `gorm:"type:timestamptz;not null" json:"changedAt"`
	CreatedAt               time.Time                `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	Tenant       Tenant            `gorm:"foreignKey:TenantID" json:"-"`
	Attendance   StudentAttendance `gorm:"foreignKey:AttendanceID" json:"attendance,omitempty"`
	ChangedByUser User             `gorm:"foreignKey:ChangedBy" json:"changedByUser,omitempty"`
}

// TableName returns the table name for the StudentAttendanceAudit model.
func (StudentAttendanceAudit) TableName() string {
	return "student_attendance_audit"
}

// BeforeCreate hook for StudentAttendanceAudit.
func (a *StudentAttendanceAudit) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.ChangedAt.IsZero() {
		a.ChangedAt = time.Now()
	}
	return nil
}
