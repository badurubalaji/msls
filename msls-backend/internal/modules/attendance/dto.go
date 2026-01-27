// Package attendance provides staff attendance management functionality.
package attendance

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceStatus represents the status of an attendance record.
type AttendanceStatus string

const (
	StatusPresent  AttendanceStatus = "present"
	StatusAbsent   AttendanceStatus = "absent"
	StatusHalfDay  AttendanceStatus = "half_day"
	StatusOnLeave  AttendanceStatus = "on_leave"
	StatusHoliday  AttendanceStatus = "holiday"
)

// IsValid checks if the attendance status is valid.
func (s AttendanceStatus) IsValid() bool {
	switch s {
	case StatusPresent, StatusAbsent, StatusHalfDay, StatusOnLeave, StatusHoliday:
		return true
	}
	return false
}

// HalfDayType represents the type of half day.
type HalfDayType string

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

const (
	RegularizationPending  RegularizationStatus = "pending"
	RegularizationApproved RegularizationStatus = "approved"
	RegularizationRejected RegularizationStatus = "rejected"
)

// IsValid checks if the regularization status is valid.
func (s RegularizationStatus) IsValid() bool {
	switch s {
	case RegularizationPending, RegularizationApproved, RegularizationRejected:
		return true
	}
	return false
}

// CheckInDTO represents a request to mark check-in.
type CheckInDTO struct {
	TenantID    uuid.UUID
	StaffID     uuid.UUID
	HalfDayType HalfDayType
	Remarks     string
	MarkedBy    *uuid.UUID
}

// CheckOutDTO represents a request to mark check-out.
type CheckOutDTO struct {
	TenantID uuid.UUID
	StaffID  uuid.UUID
	Remarks  string
	MarkedBy *uuid.UUID
}

// MarkAttendanceDTO represents a request to mark attendance for a staff member (by HR).
type MarkAttendanceDTO struct {
	TenantID       uuid.UUID
	StaffID        uuid.UUID
	AttendanceDate time.Time
	Status         AttendanceStatus
	CheckInTime    *time.Time
	CheckOutTime   *time.Time
	HalfDayType    HalfDayType
	Remarks        string
	MarkedBy       *uuid.UUID
}

// RegularizationRequestDTO represents a request to submit a regularization.
type RegularizationRequestDTO struct {
	TenantID             uuid.UUID
	StaffID              uuid.UUID
	RequestDate          time.Time
	RequestedStatus      AttendanceStatus
	Reason               string
	SupportingDocumentURL string
}

// RegularizationReviewDTO represents a request to review a regularization.
type RegularizationReviewDTO struct {
	TenantID        uuid.UUID
	RegularizationID uuid.UUID
	Approved        bool
	RejectionReason string
	ReviewedBy      *uuid.UUID
}

// SettingsDTO represents attendance settings for a branch.
type SettingsDTO struct {
	TenantID                   uuid.UUID
	BranchID                   uuid.UUID
	WorkStartTime              time.Time
	WorkEndTime                time.Time
	LateThresholdMinutes       int
	HalfDayThresholdHours      float64
	AllowSelfCheckout          bool
	RequireRegularizationApproval bool
}

// ListFilter contains filters for listing attendance records.
type ListFilter struct {
	TenantID      uuid.UUID
	StaffID       *uuid.UUID
	BranchID      *uuid.UUID
	DepartmentID  *uuid.UUID
	Status        *AttendanceStatus
	DateFrom      *time.Time
	DateTo        *time.Time
	Cursor        string
	Limit         int
	SortBy        string
	SortOrder     string
}

// RegularizationFilter contains filters for listing regularization requests.
type RegularizationFilter struct {
	TenantID uuid.UUID
	StaffID  *uuid.UUID
	Status   *RegularizationStatus
	DateFrom *time.Time
	DateTo   *time.Time
	Cursor   string
	Limit    int
}

// AttendanceResponse represents an attendance record in API responses.
type AttendanceResponse struct {
	ID             string `json:"id"`
	StaffID        string `json:"staffId"`
	StaffName      string `json:"staffName,omitempty"`
	EmployeeID     string `json:"employeeId,omitempty"`
	AttendanceDate string `json:"attendanceDate"`
	Status         string `json:"status"`
	CheckInTime    string `json:"checkInTime,omitempty"`
	CheckOutTime   string `json:"checkOutTime,omitempty"`
	IsLate         bool   `json:"isLate"`
	LateMinutes    int    `json:"lateMinutes"`
	HalfDayType    string `json:"halfDayType,omitempty"`
	Remarks        string `json:"remarks,omitempty"`
	MarkedBy       string `json:"markedBy,omitempty"`
	MarkedAt       string `json:"markedAt"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// TodayAttendanceResponse represents today's attendance status.
type TodayAttendanceResponse struct {
	Status       string `json:"status"` // not_marked, checked_in, checked_out
	Attendance   *AttendanceResponse `json:"attendance,omitempty"`
	CanCheckIn   bool   `json:"canCheckIn"`
	CanCheckOut  bool   `json:"canCheckOut"`
}

// AttendanceSummaryResponse represents monthly attendance summary.
type AttendanceSummaryResponse struct {
	Month         string `json:"month"`
	Year          int    `json:"year"`
	TotalDays     int    `json:"totalDays"`
	PresentDays   int    `json:"presentDays"`
	AbsentDays    int    `json:"absentDays"`
	HalfDays      int    `json:"halfDays"`
	LeaveDays     int    `json:"leaveDays"`
	HolidayDays   int    `json:"holidayDays"`
	LateDays      int    `json:"lateDays"`
	TotalLateMinutes int `json:"totalLateMinutes"`
}

// AttendanceListResponse represents a paginated list of attendance records.
type AttendanceListResponse struct {
	Attendance []AttendanceResponse `json:"attendance"`
	NextCursor string               `json:"nextCursor,omitempty"`
	HasMore    bool                 `json:"hasMore"`
	Total      int64                `json:"total,omitempty"`
}

// RegularizationResponse represents a regularization request in API responses.
type RegularizationResponse struct {
	ID                    string `json:"id"`
	StaffID               string `json:"staffId"`
	StaffName             string `json:"staffName,omitempty"`
	EmployeeID            string `json:"employeeId,omitempty"`
	AttendanceID          string `json:"attendanceId,omitempty"`
	RequestDate           string `json:"requestDate"`
	RequestedStatus       string `json:"requestedStatus"`
	Reason                string `json:"reason"`
	SupportingDocumentURL string `json:"supportingDocumentUrl,omitempty"`
	Status                string `json:"status"`
	ReviewedBy            string `json:"reviewedBy,omitempty"`
	ReviewedAt            string `json:"reviewedAt,omitempty"`
	RejectionReason       string `json:"rejectionReason,omitempty"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// RegularizationListResponse represents a paginated list of regularization requests.
type RegularizationListResponse struct {
	Regularizations []RegularizationResponse `json:"regularizations"`
	NextCursor      string                   `json:"nextCursor,omitempty"`
	HasMore         bool                     `json:"hasMore"`
	Total           int64                    `json:"total,omitempty"`
}

// SettingsResponse represents attendance settings in API responses.
type SettingsResponse struct {
	ID                         string  `json:"id"`
	BranchID                   string  `json:"branchId"`
	BranchName                 string  `json:"branchName,omitempty"`
	WorkStartTime              string  `json:"workStartTime"`
	WorkEndTime                string  `json:"workEndTime"`
	LateThresholdMinutes       int     `json:"lateThresholdMinutes"`
	HalfDayThresholdHours      float64 `json:"halfDayThresholdHours"`
	AllowSelfCheckout          bool    `json:"allowSelfCheckout"`
	RequireRegularizationApproval bool `json:"requireRegularizationApproval"`
	CreatedAt                  string  `json:"createdAt"`
	UpdatedAt                  string  `json:"updatedAt"`
}
