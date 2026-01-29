// Package assignment provides teacher subject assignment functionality.
package assignment

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateAssignmentDTO represents a request to create a teacher assignment.
type CreateAssignmentDTO struct {
	TenantID       uuid.UUID
	StaffID        uuid.UUID
	SubjectID      uuid.UUID
	ClassID        uuid.UUID
	SectionID      *uuid.UUID
	AcademicYearID uuid.UUID
	PeriodsPerWeek int
	IsClassTeacher bool
	EffectiveFrom  time.Time
	EffectiveTo    *time.Time
	Remarks        string
	CreatedBy      *uuid.UUID
}

// UpdateAssignmentDTO represents a request to update a teacher assignment.
type UpdateAssignmentDTO struct {
	PeriodsPerWeek *int
	IsClassTeacher *bool
	EffectiveFrom  *time.Time
	EffectiveTo    *time.Time
	Status         *models.AssignmentStatus
	Remarks        *string
}

// BulkCreateAssignmentDTO represents a request to create multiple assignments.
type BulkCreateAssignmentDTO struct {
	TenantID       uuid.UUID
	Assignments    []BulkAssignmentItem
	CreatedBy      *uuid.UUID
}

// BulkAssignmentItem represents a single assignment in a bulk create request.
type BulkAssignmentItem struct {
	StaffID        uuid.UUID
	SubjectID      uuid.UUID
	ClassID        uuid.UUID
	SectionID      *uuid.UUID
	AcademicYearID uuid.UUID
	PeriodsPerWeek int
	IsClassTeacher bool
	EffectiveFrom  time.Time
}

// WorkloadSettingsDTO represents workload settings configuration.
type WorkloadSettingsDTO struct {
	MinPeriodsPerWeek     int
	MaxPeriodsPerWeek     int
	MaxSubjectsPerTeacher *int
	MaxClassesPerTeacher  *int
}

// ListFilter contains filters for listing assignments.
type ListFilter struct {
	TenantID       uuid.UUID
	StaffID        *uuid.UUID
	SubjectID      *uuid.UUID
	ClassID        *uuid.UUID
	SectionID      *uuid.UUID
	AcademicYearID *uuid.UUID
	IsClassTeacher *bool
	Status         *models.AssignmentStatus
	Cursor         string
	Limit          int
}

// AssignmentResponse represents an assignment in API responses.
type AssignmentResponse struct {
	ID             string `json:"id"`
	StaffID        string `json:"staffId"`
	StaffName      string `json:"staffName"`
	StaffEmployeeID string `json:"staffEmployeeId,omitempty"`
	SubjectID      string `json:"subjectId"`
	SubjectName    string `json:"subjectName"`
	SubjectCode    string `json:"subjectCode"`
	ClassID        string `json:"classId"`
	ClassName      string `json:"className"`
	ClassCode      string `json:"classCode"`
	SectionID      string `json:"sectionId,omitempty"`
	SectionName    string `json:"sectionName,omitempty"`
	AcademicYearID string `json:"academicYearId"`
	AcademicYearName string `json:"academicYearName"`
	PeriodsPerWeek int    `json:"periodsPerWeek"`
	IsClassTeacher bool   `json:"isClassTeacher"`
	EffectiveFrom  string `json:"effectiveFrom"`
	EffectiveTo    string `json:"effectiveTo,omitempty"`
	Status         string `json:"status"`
	Remarks        string `json:"remarks,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// AssignmentListResponse represents a paginated list of assignments.
type AssignmentListResponse struct {
	Assignments []AssignmentResponse `json:"assignments"`
	NextCursor  string               `json:"nextCursor,omitempty"`
	HasMore     bool                 `json:"hasMore"`
	Total       int64                `json:"total,omitempty"`
}

// WorkloadSummary represents a teacher's workload summary.
type WorkloadSummary struct {
	StaffID          string `json:"staffId"`
	StaffName        string `json:"staffName"`
	StaffEmployeeID  string `json:"staffEmployeeId"`
	DepartmentName   string `json:"departmentName,omitempty"`
	TotalPeriods     int    `json:"totalPeriods"`
	TotalSubjects    int    `json:"totalSubjects"`
	TotalClasses     int    `json:"totalClasses"`
	IsClassTeacher   bool   `json:"isClassTeacher"`
	ClassTeacherFor  string `json:"classTeacherFor,omitempty"`
	WorkloadStatus   string `json:"workloadStatus"` // "under", "normal", "over"
	MinPeriods       int    `json:"minPeriods"`
	MaxPeriods       int    `json:"maxPeriods"`
}

// WorkloadReportResponse represents the workload report for all teachers.
type WorkloadReportResponse struct {
	Teachers       []WorkloadSummary `json:"teachers"`
	TotalTeachers  int               `json:"totalTeachers"`
	OverAssigned   int               `json:"overAssigned"`
	UnderAssigned  int               `json:"underAssigned"`
	NormalAssigned int               `json:"normalAssigned"`
}

// UnassignedSubject represents a subject without a teacher assigned.
type UnassignedSubject struct {
	SubjectID   string `json:"subjectId"`
	SubjectName string `json:"subjectName"`
	SubjectCode string `json:"subjectCode"`
	ClassID     string `json:"classId"`
	ClassName   string `json:"className"`
	SectionID   string `json:"sectionId,omitempty"`
	SectionName string `json:"sectionName,omitempty"`
}

// UnassignedSubjectsResponse represents subjects without teachers.
type UnassignedSubjectsResponse struct {
	Subjects []UnassignedSubject `json:"subjects"`
	Total    int                 `json:"total"`
}

// ClassTeacherResponse represents the class teacher for a class-section.
type ClassTeacherResponse struct {
	ClassID     string `json:"classId"`
	ClassName   string `json:"className"`
	SectionID   string `json:"sectionId,omitempty"`
	SectionName string `json:"sectionName,omitempty"`
	TeacherID   string `json:"teacherId,omitempty"`
	TeacherName string `json:"teacherName,omitempty"`
	IsAssigned  bool   `json:"isAssigned"`
}

// WorkloadSettingsResponse represents workload settings in API responses.
type WorkloadSettingsResponse struct {
	ID                    string `json:"id"`
	BranchID              string `json:"branchId"`
	BranchName            string `json:"branchName,omitempty"`
	MinPeriodsPerWeek     int    `json:"minPeriodsPerWeek"`
	MaxPeriodsPerWeek     int    `json:"maxPeriodsPerWeek"`
	MaxSubjectsPerTeacher *int   `json:"maxSubjectsPerTeacher,omitempty"`
	MaxClassesPerTeacher  *int   `json:"maxClassesPerTeacher,omitempty"`
}

// ToAssignmentResponse converts a TeacherSubjectAssignment model to an AssignmentResponse.
func ToAssignmentResponse(a *models.TeacherSubjectAssignment) AssignmentResponse {
	resp := AssignmentResponse{
		ID:             a.ID.String(),
		StaffID:        a.StaffID.String(),
		SubjectID:      a.SubjectID.String(),
		ClassID:        a.ClassID.String(),
		AcademicYearID: a.AcademicYearID.String(),
		PeriodsPerWeek: a.PeriodsPerWeek,
		IsClassTeacher: a.IsClassTeacher,
		EffectiveFrom:  a.EffectiveFrom.Format("2006-01-02"),
		Status:         string(a.Status),
		Remarks:        a.Remarks,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}

	if a.EffectiveTo != nil {
		resp.EffectiveTo = a.EffectiveTo.Format("2006-01-02")
	}

	if a.SectionID != nil {
		resp.SectionID = a.SectionID.String()
	}

	// Add staff details if loaded
	if a.Staff.ID != uuid.Nil {
		resp.StaffName = a.Staff.FullName()
		resp.StaffEmployeeID = a.Staff.EmployeeID
	}

	// Add subject details if loaded
	if a.Subject.ID != uuid.Nil {
		resp.SubjectName = a.Subject.Name
		resp.SubjectCode = a.Subject.Code
	}

	// Add class details if loaded
	if a.Class.ID != uuid.Nil {
		resp.ClassName = a.Class.Name
		resp.ClassCode = a.Class.Code
	}

	// Add section details if loaded
	if a.Section != nil && a.Section.ID != uuid.Nil {
		resp.SectionName = a.Section.Name
	}

	// Add academic year details if loaded
	if a.AcademicYear.ID != uuid.Nil {
		resp.AcademicYearName = a.AcademicYear.Name
	}

	return resp
}

// ToAssignmentResponses converts a slice of assignments to responses.
func ToAssignmentResponses(assignments []models.TeacherSubjectAssignment) []AssignmentResponse {
	responses := make([]AssignmentResponse, len(assignments))
	for i := range assignments {
		responses[i] = ToAssignmentResponse(&assignments[i])
	}
	return responses
}

// ToWorkloadSettingsResponse converts a TeacherWorkloadSettings model to a response.
func ToWorkloadSettingsResponse(w *models.TeacherWorkloadSettings) WorkloadSettingsResponse {
	resp := WorkloadSettingsResponse{
		ID:                    w.ID.String(),
		BranchID:              w.BranchID.String(),
		MinPeriodsPerWeek:     w.MinPeriodsPerWeek,
		MaxPeriodsPerWeek:     w.MaxPeriodsPerWeek,
		MaxSubjectsPerTeacher: w.MaxSubjectsPerTeacher,
		MaxClassesPerTeacher:  w.MaxClassesPerTeacher,
	}

	if w.Branch.ID != uuid.Nil {
		resp.BranchName = w.Branch.Name
	}

	return resp
}
