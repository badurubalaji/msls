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

// ============================================================================
// Period-wise Attendance DTOs (Story 7.2)
// ============================================================================

// PeriodInfoResponse represents information about a period for attendance.
type PeriodInfoResponse struct {
	PeriodSlotID     string `json:"periodSlotId"`
	TimetableEntryID string `json:"timetableEntryId"`
	PeriodName       string `json:"periodName"`
	PeriodNumber     *int   `json:"periodNumber,omitempty"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	SubjectID        string `json:"subjectId,omitempty"`
	SubjectName      string `json:"subjectName,omitempty"`
	SubjectCode      string `json:"subjectCode,omitempty"`
	StaffID          string `json:"staffId,omitempty"`
	StaffName        string `json:"staffName,omitempty"`
	IsMarked         bool   `json:"isMarked"`
	MarkedCount      int    `json:"markedCount,omitempty"`
	TotalStudents    int    `json:"totalStudents,omitempty"`
}

// TeacherPeriodsResponse represents the list of periods for a section.
type TeacherPeriodsResponse struct {
	SectionID   string               `json:"sectionId"`
	SectionName string               `json:"sectionName"`
	ClassName   string               `json:"className"`
	Date        string               `json:"date"`
	DayOfWeek   int                  `json:"dayOfWeek"`
	DayName     string               `json:"dayName"`
	Periods     []PeriodInfoResponse `json:"periods"`
}

// PeriodAttendanceRecord represents a single student's attendance for a period.
type PeriodAttendanceRecord struct {
	StudentID       uuid.UUID  `json:"studentId"`
	Status          string     `json:"status"`
	LateArrivalTime *time.Time `json:"lateArrivalTime,omitempty"`
	Remarks         string     `json:"remarks,omitempty"`
}

// MarkPeriodAttendanceRequest represents a request to mark period attendance.
type MarkPeriodAttendanceRequest struct {
	SectionID string                   `json:"sectionId" binding:"required"`
	Date      string                   `json:"date" binding:"required"`
	Records   []PeriodAttendanceRecord `json:"records" binding:"required,min=1"`
}

// MarkPeriodAttendanceDTO represents validated data for marking period attendance.
type MarkPeriodAttendanceDTO struct {
	TenantID         uuid.UUID
	SectionID        uuid.UUID
	PeriodID         uuid.UUID
	TimetableEntryID uuid.UUID
	Date             time.Time
	Records          []PeriodAttendanceRecord
	MarkedBy         uuid.UUID
}

// StudentForPeriodAttendance represents a student for period attendance marking.
type StudentForPeriodAttendance struct {
	StudentID        string `json:"studentId"`
	AdmissionNumber  string `json:"admissionNumber"`
	RollNumber       *int   `json:"rollNumber,omitempty"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	FullName         string `json:"fullName"`
	PhotoURL         string `json:"photoUrl,omitempty"`
	Status           string `json:"status,omitempty"`
	LateArrivalTime  string `json:"lateArrivalTime,omitempty"`
	Remarks          string `json:"remarks,omitempty"`
}

// PeriodAttendanceResponse represents period attendance data.
type PeriodAttendanceResponse struct {
	SectionID        string                       `json:"sectionId"`
	SectionName      string                       `json:"sectionName"`
	ClassName        string                       `json:"className"`
	Date             string                       `json:"date"`
	PeriodSlotID     string                       `json:"periodSlotId"`
	PeriodName       string                       `json:"periodName"`
	PeriodNumber     *int                         `json:"periodNumber,omitempty"`
	StartTime        string                       `json:"startTime"`
	EndTime          string                       `json:"endTime"`
	SubjectID        string                       `json:"subjectId,omitempty"`
	SubjectName      string                       `json:"subjectName,omitempty"`
	SubjectCode      string                       `json:"subjectCode,omitempty"`
	TimetableEntryID string                       `json:"timetableEntryId,omitempty"`
	Students         []StudentForPeriodAttendance `json:"students"`
	IsMarked         bool                         `json:"isMarked"`
	CanEdit          bool                         `json:"canEdit"`
	MarkedAt         string                       `json:"markedAt,omitempty"`
	MarkedByName     string                       `json:"markedByName,omitempty"`
	Summary          AttendanceSummary            `json:"summary"`
}

// MarkPeriodAttendanceResult represents the result of marking period attendance.
type MarkPeriodAttendanceResult struct {
	SectionID   string            `json:"sectionId"`
	PeriodID    string            `json:"periodId"`
	Date        string            `json:"date"`
	Summary     AttendanceSummary `json:"summary"`
	MarkedAt    string            `json:"markedAt"`
	Message     string            `json:"message"`
}

// DailySummaryStudent represents a student's attendance across all periods.
type DailySummaryStudent struct {
	StudentID            string            `json:"studentId"`
	AdmissionNumber      string            `json:"admissionNumber"`
	RollNumber           *int              `json:"rollNumber,omitempty"`
	FullName             string            `json:"fullName"`
	PhotoURL             string            `json:"photoUrl,omitempty"`
	PeriodStatuses       map[string]string `json:"periodStatuses"` // periodId -> status
	TotalPeriods         int               `json:"totalPeriods"`
	PeriodsPresent       int               `json:"periodsPresent"`
	PeriodsAbsent        int               `json:"periodsAbsent"`
	PeriodsLate          int               `json:"periodsLate"`
	AttendancePercentage float64           `json:"attendancePercentage"`
	OverallStatus        string            `json:"overallStatus"` // Derived: present if >50%
}

// DailySummaryResponse represents the daily attendance summary for all periods.
type DailySummaryResponse struct {
	SectionID   string                `json:"sectionId"`
	SectionName string                `json:"sectionName"`
	ClassName   string                `json:"className"`
	Date        string                `json:"date"`
	DayName     string                `json:"dayName"`
	Periods     []PeriodInfoResponse  `json:"periods"`
	Students    []DailySummaryStudent `json:"students"`
	Summary     DailySummarySummary   `json:"summary"`
}

// DailySummarySummary represents overall daily summary statistics.
type DailySummarySummary struct {
	TotalStudents    int     `json:"totalStudents"`
	TotalPeriods     int     `json:"totalPeriods"`
	AverageAttendance float64 `json:"averageAttendance"`
	FullPresentCount int     `json:"fullPresentCount"` // Students present all periods
	AbsentCount      int     `json:"absentCount"`      // Students absent >50% periods
}

// SubjectAttendanceRequest represents a request for subject attendance analytics.
type SubjectAttendanceRequest struct {
	StudentID uuid.UUID
	SubjectID uuid.UUID
	DateFrom  *time.Time
	DateTo    *time.Time
}

// SubjectAttendanceResponse represents subject-wise attendance statistics.
type SubjectAttendanceResponse struct {
	StudentID            string  `json:"studentId"`
	StudentName          string  `json:"studentName,omitempty"`
	SubjectID            string  `json:"subjectId"`
	SubjectName          string  `json:"subjectName"`
	SubjectCode          string  `json:"subjectCode"`
	TotalPeriods         int     `json:"totalPeriods"`
	PeriodsPresent       int     `json:"periodsPresent"`
	PeriodsAbsent        int     `json:"periodsAbsent"`
	PeriodsLate          int     `json:"periodsLate"`
	AttendancePercentage float64 `json:"attendancePercentage"`
	MinimumRequired      float64 `json:"minimumRequired"` // e.g., 75%
	IsEligible           bool    `json:"isEligible"`      // Met minimum attendance
}

// ============================================================================
// Attendance Edit & Audit DTOs (Story 7.3)
// ============================================================================

// EditAttendanceRequest represents a request to edit an attendance record.
type EditAttendanceRequest struct {
	Status          string `json:"status,omitempty"`
	Remarks         string `json:"remarks,omitempty"`
	LateArrivalTime string `json:"lateArrivalTime,omitempty"`
	Reason          string `json:"reason" binding:"required"`
}

// EditAttendanceDTO represents validated data for editing attendance.
type EditAttendanceDTO struct {
	TenantID        uuid.UUID
	AttendanceID    uuid.UUID
	Status          string
	Remarks         *string
	LateArrivalTime *time.Time
	Reason          string
	EditedBy        uuid.UUID
	IsAdmin         bool
}

// EditAttendanceResult represents the result of editing attendance.
type EditAttendanceResult struct {
	AttendanceID string `json:"attendanceId"`
	StudentID    string `json:"studentId"`
	Date         string `json:"date"`
	Status       string `json:"status"`
	EditedAt     string `json:"editedAt"`
	EditedBy     string `json:"editedBy"`
	Message      string `json:"message"`
}

// AttendanceAuditEntry represents a single audit entry for an attendance change.
type AttendanceAuditEntry struct {
	ID              string `json:"id"`
	ChangeType      string `json:"changeType"`
	PreviousStatus  string `json:"previousStatus,omitempty"`
	NewStatus       string `json:"newStatus"`
	PreviousRemarks string `json:"previousRemarks,omitempty"`
	NewRemarks      string `json:"newRemarks,omitempty"`
	ChangeReason    string `json:"changeReason"`
	ChangedByID     string `json:"changedById"`
	ChangedByName   string `json:"changedByName"`
	ChangedAt       string `json:"changedAt"`
}

// AttendanceAuditTrailResponse represents the full audit trail for an attendance record.
type AttendanceAuditTrailResponse struct {
	AttendanceID string                 `json:"attendanceId"`
	StudentID    string                 `json:"studentId"`
	StudentName  string                 `json:"studentName"`
	Date         string                 `json:"date"`
	AuditEntries []AttendanceAuditEntry `json:"auditEntries"`
	TotalChanges int                    `json:"totalChanges"`
}

// EditWindowStatusResponse represents the edit window status for an attendance record.
type EditWindowStatusResponse struct {
	AttendanceID      string `json:"attendanceId"`
	MarkedAt          string `json:"markedAt"`
	WindowEndAt       string `json:"windowEndAt"`
	WindowMinutes     int    `json:"windowMinutes"`
	RemainingMinutes  int    `json:"remainingMinutes"`
	IsWithinWindow    bool   `json:"isWithinWindow"`
	CanEdit           bool   `json:"canEdit"`
	EditDeniedReason  string `json:"editDeniedReason,omitempty"`
	IsOriginalMarker  bool   `json:"isOriginalMarker"`
	RequiresAdminEdit bool   `json:"requiresAdminEdit"`
}

// ============================================================================
// Calendar & Reports DTOs (Stories 7.4-7.8)
// ============================================================================

// CalendarDayResponse represents a single day in the attendance calendar.
type CalendarDayResponse struct {
	Date       string `json:"date"`
	DayOfWeek  int    `json:"dayOfWeek"`
	Status     string `json:"status,omitempty"`     // present, absent, late, half_day
	IsHoliday  bool   `json:"isHoliday"`
	IsWeekend  bool   `json:"isWeekend"`
	HolidayName string `json:"holidayName,omitempty"`
	Remarks    string `json:"remarks,omitempty"`
}

// MonthlyCalendarResponse represents a student's monthly attendance calendar.
type MonthlyCalendarResponse struct {
	StudentID    string                `json:"studentId"`
	StudentName  string                `json:"studentName"`
	Year         int                   `json:"year"`
	Month        int                   `json:"month"`
	MonthName    string                `json:"monthName"`
	Days         []CalendarDayResponse `json:"days"`
	Summary      MonthlySummary        `json:"summary"`
	ClassAverage float64               `json:"classAverage"`
	Trend        string                `json:"trend"` // improving, declining, stable
}

// MonthlySummary represents monthly attendance summary statistics.
type MonthlySummary struct {
	WorkingDays  int     `json:"workingDays"`
	Present      int     `json:"present"`
	Absent       int     `json:"absent"`
	Late         int     `json:"late"`
	HalfDay      int     `json:"halfDay"`
	Holidays     int     `json:"holidays"`
	Percentage   float64 `json:"percentage"`
}

// StudentSummaryResponse represents a student's attendance summary with trend.
type StudentSummaryResponse struct {
	StudentID           string          `json:"studentId"`
	StudentName         string          `json:"studentName"`
	SectionID           string          `json:"sectionId"`
	SectionName         string          `json:"sectionName"`
	DateFrom            string          `json:"dateFrom"`
	DateTo              string          `json:"dateTo"`
	Summary             MonthlySummary  `json:"summary"`
	ClassAverage        float64         `json:"classAverage"`
	Trend               string          `json:"trend"`
	TrendPercentage     float64         `json:"trendPercentage"`
	MonthlyBreakdown    []MonthBreakdown `json:"monthlyBreakdown,omitempty"`
}

// MonthBreakdown represents attendance for a single month in trend analysis.
type MonthBreakdown struct {
	Year       int     `json:"year"`
	Month      int     `json:"month"`
	MonthName  string  `json:"monthName"`
	Percentage float64 `json:"percentage"`
}

// ClassReportResponse represents class-level attendance report.
type ClassReportResponse struct {
	SectionID       string               `json:"sectionId"`
	SectionName     string               `json:"sectionName"`
	ClassName       string               `json:"className"`
	Date            string               `json:"date"`
	Students        []StudentReportEntry `json:"students"`
	Summary         AttendanceSummary    `json:"summary"`
	AttendanceRate  float64              `json:"attendanceRate"`
}

// StudentReportEntry represents a student entry in class report.
type StudentReportEntry struct {
	StudentID       string `json:"studentId"`
	AdmissionNumber string `json:"admissionNumber"`
	FullName        string `json:"fullName"`
	RollNumber      *int   `json:"rollNumber,omitempty"`
	Status          string `json:"status"`
	StatusLabel     string `json:"statusLabel"`
	Remarks         string `json:"remarks,omitempty"`
}

// MonthlyClassReportResponse represents monthly class attendance report.
type MonthlyClassReportResponse struct {
	SectionID    string                     `json:"sectionId"`
	SectionName  string                     `json:"sectionName"`
	ClassName    string                     `json:"className"`
	Year         int                        `json:"year"`
	Month        int                        `json:"month"`
	MonthName    string                     `json:"monthName"`
	WorkingDays  int                        `json:"workingDays"`
	Dates        []string                   `json:"dates"`
	Students     []MonthlyStudentReport     `json:"students"`
	Summary      ClassMonthlySummary        `json:"summary"`
}

// MonthlyStudentReport represents a student's monthly attendance in grid format.
type MonthlyStudentReport struct {
	StudentID       string            `json:"studentId"`
	AdmissionNumber string            `json:"admissionNumber"`
	FullName        string            `json:"fullName"`
	RollNumber      *int              `json:"rollNumber,omitempty"`
	DailyStatus     map[string]string `json:"dailyStatus"` // date -> status
	Present         int               `json:"present"`
	Absent          int               `json:"absent"`
	Late            int               `json:"late"`
	Percentage      float64           `json:"percentage"`
}

// ClassMonthlySummary represents class-level monthly summary.
type ClassMonthlySummary struct {
	TotalStudents       int     `json:"totalStudents"`
	AverageAttendance   float64 `json:"averageAttendance"`
	StudentsAbove90     int     `json:"studentsAbove90"`
	StudentsBelow75     int     `json:"studentsBelow75"`
	StudentsBelow60     int     `json:"studentsBelow60"`
}

// ClassComparisonResponse represents attendance comparison across sections.
type ClassComparisonResponse struct {
	DateFrom    string              `json:"dateFrom"`
	DateTo      string              `json:"dateTo"`
	Sections    []SectionComparison `json:"sections"`
	Threshold   float64             `json:"threshold"`
}

// SectionComparison represents a section's attendance for comparison.
type SectionComparison struct {
	SectionID       string  `json:"sectionId"`
	SectionName     string  `json:"sectionName"`
	ClassName       string  `json:"className"`
	TotalStudents   int     `json:"totalStudents"`
	AttendanceRate  float64 `json:"attendanceRate"`
	BelowThreshold  bool    `json:"belowThreshold"`
	Trend           string  `json:"trend"`
}

// IndividualReportResponse represents individual student attendance report.
type IndividualReportResponse struct {
	StudentID       string                 `json:"studentId"`
	AdmissionNumber string                 `json:"admissionNumber"`
	FullName        string                 `json:"fullName"`
	ClassName       string                 `json:"className"`
	SectionName     string                 `json:"sectionName"`
	DateFrom        string                 `json:"dateFrom"`
	DateTo          string                 `json:"dateTo"`
	Summary         MonthlySummary         `json:"summary"`
	MonthlyDetails  []MonthlyDetail        `json:"monthlyDetails"`
	SubjectWise     []SubjectAttendanceResponse `json:"subjectWise,omitempty"`
}

// MonthlyDetail represents attendance details for a month.
type MonthlyDetail struct {
	Year        int                   `json:"year"`
	Month       int                   `json:"month"`
	MonthName   string                `json:"monthName"`
	Days        []CalendarDayResponse `json:"days"`
	Summary     MonthlySummary        `json:"summary"`
}

// AttendanceCertificateRequest represents request for attendance certificate.
type AttendanceCertificateRequest struct {
	StudentID string `json:"studentId" binding:"required,uuid"`
	DateFrom  string `json:"dateFrom" binding:"required"`
	DateTo    string `json:"dateTo" binding:"required"`
	Purpose   string `json:"purpose"`
}

// LowAttendanceStudent represents a student with low attendance.
type LowAttendanceStudent struct {
	StudentID       string  `json:"studentId"`
	AdmissionNumber string  `json:"admissionNumber"`
	FullName        string  `json:"fullName"`
	ClassName       string  `json:"className"`
	SectionName     string  `json:"sectionName"`
	AttendanceRate  float64 `json:"attendanceRate"`
	DaysAbsent      int     `json:"daysAbsent"`
	LastPresent     string  `json:"lastPresent,omitempty"`
	ConsecutiveAbsent int   `json:"consecutiveAbsent"`
}

// LowAttendanceDashboardResponse represents low attendance dashboard data.
type LowAttendanceDashboardResponse struct {
	DateFrom              string                 `json:"dateFrom"`
	DateTo                string                 `json:"dateTo"`
	Threshold             float64                `json:"threshold"`
	CriticalThreshold     float64                `json:"criticalThreshold"`
	TotalStudents         int                    `json:"totalStudents"`
	BelowThreshold        int                    `json:"belowThreshold"`
	ChronicAbsentees      int                    `json:"chronicAbsentees"`
	OverallAttendanceRate float64                `json:"overallAttendanceRate"`
	Students              []LowAttendanceStudent `json:"students"`
	TrendData             []TrendDataPoint       `json:"trendData"`
	ClassBreakdown        []ClassAttendanceBreakdown `json:"classBreakdown"`
}

// TrendDataPoint represents a data point for attendance trend.
type TrendDataPoint struct {
	Date           string  `json:"date"`
	AttendanceRate float64 `json:"attendanceRate"`
}

// ClassAttendanceBreakdown represents attendance breakdown by class.
type ClassAttendanceBreakdown struct {
	ClassName       string  `json:"className"`
	SectionName     string  `json:"sectionName"`
	SectionID       string  `json:"sectionId"`
	TotalStudents   int     `json:"totalStudents"`
	AttendanceRate  float64 `json:"attendanceRate"`
	BelowThreshold  int     `json:"belowThreshold"`
}

// UnmarkedAttendanceResponse represents classes with unmarked attendance.
type UnmarkedAttendanceResponse struct {
	Date         string                   `json:"date"`
	Deadline     string                   `json:"deadline"`
	IsPostDeadline bool                   `json:"isPostDeadline"`
	UnmarkedClasses []UnmarkedClassInfo  `json:"unmarkedClasses"`
	TotalClasses int                      `json:"totalClasses"`
	MarkedClasses int                     `json:"markedClasses"`
}

// UnmarkedClassInfo represents a class with unmarked attendance.
type UnmarkedClassInfo struct {
	SectionID     string `json:"sectionId"`
	SectionName   string `json:"sectionName"`
	ClassName     string `json:"className"`
	TeacherID     string `json:"teacherId,omitempty"`
	TeacherName   string `json:"teacherName,omitempty"`
	StudentCount  int    `json:"studentCount"`
	IsEscalated   bool   `json:"isEscalated"`
}

// DeadlineSettingsResponse represents attendance deadline settings.
type DeadlineSettingsResponse struct {
	BranchID              string   `json:"branchId"`
	DeadlineTime          string   `json:"deadlineTime"`
	LateMarkingAllowed    bool     `json:"lateMarkingAllowed"`
	EscalationAfterMinutes int     `json:"escalationAfterMinutes"`
	NotifyTeacher         bool     `json:"notifyTeacher"`
	NotifyHOD             bool     `json:"notifyHod"`
	Holidays              []HolidayInfo `json:"holidays,omitempty"`
}

// HolidayInfo represents a holiday entry.
type HolidayInfo struct {
	Date        string `json:"date"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsRecurring bool   `json:"isRecurring"`
}

// DailyReportSummaryResponse represents end-of-day attendance summary.
type DailyReportSummaryResponse struct {
	Date                  string                     `json:"date"`
	OverallAttendanceRate float64                    `json:"overallAttendanceRate"`
	TotalStudents         int                        `json:"totalStudents"`
	TotalPresent          int                        `json:"totalPresent"`
	TotalAbsent           int                        `json:"totalAbsent"`
	ClassesMarked         int                        `json:"classesMarked"`
	ClassesTotal          int                        `json:"classesTotal"`
	LowAttendanceClasses  []ClassAttendanceBreakdown `json:"lowAttendanceClasses"`
	GeneratedAt           string                     `json:"generatedAt"`
}
