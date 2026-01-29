// Package studentattendance provides student attendance management functionality.
package studentattendance

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceStatus represents the status of a student attendance record.
type AttendanceStatus string

const (
	StatusPresent  AttendanceStatus = "present"
	StatusAbsent   AttendanceStatus = "absent"
	StatusLate     AttendanceStatus = "late"
	StatusHalfDay  AttendanceStatus = "half_day"
)

// IsValid checks if the attendance status is valid.
func (s AttendanceStatus) IsValid() bool {
	switch s {
	case StatusPresent, StatusAbsent, StatusLate, StatusHalfDay:
		return true
	}
	return false
}

// Label returns a human-readable label for the status.
func (s AttendanceStatus) Label() string {
	switch s {
	case StatusPresent:
		return "Present"
	case StatusAbsent:
		return "Absent"
	case StatusLate:
		return "Late"
	case StatusHalfDay:
		return "Half Day"
	default:
		return string(s)
	}
}

// ShortLabel returns a single character label for the status.
func (s AttendanceStatus) ShortLabel() string {
	switch s {
	case StatusPresent:
		return "P"
	case StatusAbsent:
		return "A"
	case StatusLate:
		return "L"
	case StatusHalfDay:
		return "H"
	default:
		return "?"
	}
}

// StudentAttendanceRecord represents a single student's attendance for marking.
type StudentAttendanceRecord struct {
	StudentID       uuid.UUID  `json:"studentId"`
	Status          string     `json:"status"`
	LateArrivalTime *time.Time `json:"lateArrivalTime,omitempty"`
	Remarks         string     `json:"remarks,omitempty"`
}

// MarkClassAttendanceDTO represents a request to mark attendance for a class section.
type MarkClassAttendanceDTO struct {
	TenantID   uuid.UUID
	SectionID  uuid.UUID
	Date       time.Time
	Records    []StudentAttendanceRecord
	MarkedBy   uuid.UUID
}

// ListFilter contains filters for listing attendance records.
type ListFilter struct {
	TenantID   uuid.UUID
	SectionID  *uuid.UUID
	StudentID  *uuid.UUID
	Status     *AttendanceStatus
	DateFrom   *time.Time
	DateTo     *time.Time
	Cursor     string
	Limit      int
	SortBy     string
	SortOrder  string
}

// SettingsDTO represents student attendance settings for a branch.
type SettingsDTO struct {
	TenantID             uuid.UUID
	BranchID             uuid.UUID
	EditWindowMinutes    int
	LateThresholdMinutes int
	SMSOnAbsent          bool
}

// TeacherClassResponse represents a class section assigned to a teacher.
type TeacherClassResponse struct {
	SectionID      string `json:"sectionId"`
	SectionName    string `json:"sectionName"`
	SectionCode    string `json:"sectionCode"`
	ClassName      string `json:"className"`
	ClassCode      string `json:"classCode"`
	StudentCount   int    `json:"studentCount"`
	IsMarkedToday  bool   `json:"isMarkedToday"`
	MarkedCount    int    `json:"markedCount,omitempty"`
}

// StudentForAttendance represents a student record for attendance marking.
type StudentForAttendance struct {
	StudentID        string            `json:"studentId"`
	AdmissionNumber  string            `json:"admissionNumber"`
	RollNumber       *int              `json:"rollNumber,omitempty"`
	FirstName        string            `json:"firstName"`
	LastName         string            `json:"lastName"`
	FullName         string            `json:"fullName"`
	PhotoURL         string            `json:"photoUrl,omitempty"`
	Status           string            `json:"status,omitempty"`
	LateArrivalTime  string            `json:"lateArrivalTime,omitempty"`
	Remarks          string            `json:"remarks,omitempty"`
	Last5Days        []string          `json:"last5Days"`
}

// ClassAttendanceResponse represents attendance data for a class section.
type ClassAttendanceResponse struct {
	SectionID      string                 `json:"sectionId"`
	SectionName    string                 `json:"sectionName"`
	ClassName      string                 `json:"className"`
	Date           string                 `json:"date"`
	Students       []StudentForAttendance `json:"students"`
	IsMarked       bool                   `json:"isMarked"`
	CanEdit        bool                   `json:"canEdit"`
	MarkedAt       string                 `json:"markedAt,omitempty"`
	MarkedBy       string                 `json:"markedBy,omitempty"`
	MarkedByName   string                 `json:"markedByName,omitempty"`
	Summary        AttendanceSummary      `json:"summary"`
}

// AttendanceSummary represents a summary of attendance counts.
type AttendanceSummary struct {
	Total    int `json:"total"`
	Present  int `json:"present"`
	Absent   int `json:"absent"`
	Late     int `json:"late"`
	HalfDay  int `json:"halfDay"`
}

// AttendanceResponse represents a single attendance record in API responses.
type AttendanceResponse struct {
	ID              string `json:"id"`
	StudentID       string `json:"studentId"`
	StudentName     string `json:"studentName,omitempty"`
	AdmissionNumber string `json:"admissionNumber,omitempty"`
	SectionID       string `json:"sectionId"`
	SectionName     string `json:"sectionName,omitempty"`
	AttendanceDate  string `json:"attendanceDate"`
	Status          string `json:"status"`
	StatusLabel     string `json:"statusLabel"`
	LateArrivalTime string `json:"lateArrivalTime,omitempty"`
	Remarks         string `json:"remarks,omitempty"`
	MarkedBy        string `json:"markedBy,omitempty"`
	MarkedByName    string `json:"markedByName,omitempty"`
	MarkedAt        string `json:"markedAt"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// AttendanceListResponse represents a paginated list of attendance records.
type AttendanceListResponse struct {
	Attendance []AttendanceResponse `json:"attendance"`
	NextCursor string               `json:"nextCursor,omitempty"`
	HasMore    bool                 `json:"hasMore"`
	Total      int64                `json:"total,omitempty"`
}

// MarkAttendanceResult represents the result of marking attendance.
type MarkAttendanceResult struct {
	SectionID   string            `json:"sectionId"`
	Date        string            `json:"date"`
	Summary     AttendanceSummary `json:"summary"`
	MarkedAt    string            `json:"markedAt"`
	Message     string            `json:"message"`
}

// SettingsResponse represents student attendance settings in API responses.
type SettingsResponse struct {
	ID                   string `json:"id,omitempty"`
	BranchID             string `json:"branchId"`
	BranchName           string `json:"branchName,omitempty"`
	EditWindowMinutes    int    `json:"editWindowMinutes"`
	LateThresholdMinutes int    `json:"lateThresholdMinutes"`
	SMSOnAbsent          bool   `json:"smsOnAbsent"`
	CreatedAt            string `json:"createdAt,omitempty"`
	UpdatedAt            string `json:"updatedAt,omitempty"`
}

// StudentAttendanceHistoryEntry represents a single day's attendance for history.
type StudentAttendanceHistoryEntry struct {
	Date   time.Time
	Status AttendanceStatus
}
