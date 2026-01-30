package timetable

import (
	"time"

	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
)

// ========================================
// Substitution DTOs
// ========================================

// SubstitutionFilter contains filter parameters for listing substitutions.
type SubstitutionFilter struct {
	TenantID          uuid.UUID
	BranchID          *uuid.UUID
	OriginalStaffID   *uuid.UUID
	SubstituteStaffID *uuid.UUID
	StartDate         *time.Time
	EndDate           *time.Time
	Status            *string
	Limit             int
	Offset            int
}

// CreateSubstitutionRequest contains data for creating a substitution.
type CreateSubstitutionRequest struct {
	BranchID          uuid.UUID                       `json:"branchId" binding:"required"`
	OriginalStaffID   uuid.UUID                       `json:"originalStaffId" binding:"required"`
	SubstituteStaffID uuid.UUID                       `json:"substituteStaffId" binding:"required"`
	SubstitutionDate  string                          `json:"substitutionDate" binding:"required"`
	Reason            string                          `json:"reason"`
	Notes             string                          `json:"notes"`
	Periods           []CreateSubstitutionPeriodInput `json:"periods" binding:"required,min=1"`
}

// CreateSubstitutionPeriodInput contains period data for a substitution.
type CreateSubstitutionPeriodInput struct {
	PeriodSlotID     uuid.UUID  `json:"periodSlotId" binding:"required"`
	TimetableEntryID *uuid.UUID `json:"timetableEntryId"`
	SubjectID        *uuid.UUID `json:"subjectId"`
	SectionID        *uuid.UUID `json:"sectionId"`
	RoomNumber       string     `json:"roomNumber"`
	Notes            string     `json:"notes"`
}

// UpdateSubstitutionRequest contains data for updating a substitution.
type UpdateSubstitutionRequest struct {
	SubstituteStaffID *uuid.UUID `json:"substituteStaffId"`
	Reason            *string    `json:"reason"`
	Notes             *string    `json:"notes"`
	Status            *string    `json:"status"`
}

// SubstitutionResponse is the API response for a substitution.
type SubstitutionResponse struct {
	ID                  uuid.UUID                    `json:"id"`
	BranchID            uuid.UUID                    `json:"branchId"`
	BranchName          string                       `json:"branchName,omitempty"`
	OriginalStaffID     uuid.UUID                    `json:"originalStaffId"`
	OriginalStaffName   string                       `json:"originalStaffName,omitempty"`
	SubstituteStaffID   uuid.UUID                    `json:"substituteStaffId"`
	SubstituteStaffName string                       `json:"substituteStaffName,omitempty"`
	SubstitutionDate    string                       `json:"substitutionDate"`
	Reason              string                       `json:"reason,omitempty"`
	Status              string                       `json:"status"`
	Notes               string                       `json:"notes,omitempty"`
	CreatedBy           *uuid.UUID                   `json:"createdBy,omitempty"`
	CreatedByName       string                       `json:"createdByName,omitempty"`
	ApprovedBy          *uuid.UUID                   `json:"approvedBy,omitempty"`
	ApprovedByName      string                       `json:"approvedByName,omitempty"`
	ApprovedAt          *string                      `json:"approvedAt,omitempty"`
	Periods             []SubstitutionPeriodResponse `json:"periods,omitempty"`
	CreatedAt           string                       `json:"createdAt"`
	UpdatedAt           string                       `json:"updatedAt"`
}

// SubstitutionPeriodResponse is the API response for a substitution period.
type SubstitutionPeriodResponse struct {
	ID               uuid.UUID  `json:"id"`
	PeriodSlotID     uuid.UUID  `json:"periodSlotId"`
	PeriodSlotName   string     `json:"periodSlotName,omitempty"`
	StartTime        string     `json:"startTime,omitempty"`
	EndTime          string     `json:"endTime,omitempty"`
	TimetableEntryID *uuid.UUID `json:"timetableEntryId,omitempty"`
	SubjectID        *uuid.UUID `json:"subjectId,omitempty"`
	SubjectName      string     `json:"subjectName,omitempty"`
	SectionID        *uuid.UUID `json:"sectionId,omitempty"`
	SectionName      string     `json:"sectionName,omitempty"`
	ClassName        string     `json:"className,omitempty"`
	RoomNumber       string     `json:"roomNumber,omitempty"`
	Notes            string     `json:"notes,omitempty"`
}

// SubstitutionListResponse is the API response for listing substitutions.
type SubstitutionListResponse struct {
	Substitutions []SubstitutionResponse `json:"substitutions"`
	Total         int64                  `json:"total"`
}

// AvailableTeacherResponse contains data about an available substitute teacher.
type AvailableTeacherResponse struct {
	StaffID       uuid.UUID `json:"staffId"`
	StaffName     string    `json:"staffName"`
	DepartmentID  *uuid.UUID `json:"departmentId,omitempty"`
	Department    string    `json:"department,omitempty"`
	FreePeriods   int       `json:"freePeriods"`
	TotalPeriods  int       `json:"totalPeriods"`
	HasConflict   bool      `json:"hasConflict"`
}

// AvailableTeachersResponse is the response for available teachers endpoint.
type AvailableTeachersResponse struct {
	Teachers []AvailableTeacherResponse `json:"teachers"`
}

// ========================================
// Response Converters
// ========================================

// SubstitutionToResponse converts a Substitution model to response DTO.
func SubstitutionToResponse(s *models.Substitution) SubstitutionResponse {
	resp := SubstitutionResponse{
		ID:                s.ID,
		BranchID:          s.BranchID,
		OriginalStaffID:   s.OriginalStaffID,
		SubstituteStaffID: s.SubstituteStaffID,
		SubstitutionDate:  s.SubstitutionDate.Format("2006-01-02"),
		Reason:            s.Reason,
		Status:            string(s.Status),
		Notes:             s.Notes,
		CreatedBy:         s.CreatedBy,
		ApprovedBy:        s.ApprovedBy,
		CreatedAt:         s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         s.UpdatedAt.Format(time.RFC3339),
	}

	if s.Branch != nil {
		resp.BranchName = s.Branch.Name
	}

	if s.OriginalStaff != nil {
		resp.OriginalStaffName = s.OriginalStaff.FirstName
		if s.OriginalStaff.LastName != "" {
			resp.OriginalStaffName += " " + s.OriginalStaff.LastName
		}
	}

	if s.SubstituteStaff != nil {
		resp.SubstituteStaffName = s.SubstituteStaff.FirstName
		if s.SubstituteStaff.LastName != "" {
			resp.SubstituteStaffName += " " + s.SubstituteStaff.LastName
		}
	}

	if s.Creator != nil {
		resp.CreatedByName = s.Creator.FullName()
	}

	if s.Approver != nil {
		resp.ApprovedByName = s.Approver.FullName()
	}

	if s.ApprovedAt != nil {
		approvedAtStr := s.ApprovedAt.Format(time.RFC3339)
		resp.ApprovedAt = &approvedAtStr
	}

	if len(s.Periods) > 0 {
		resp.Periods = make([]SubstitutionPeriodResponse, len(s.Periods))
		for i, p := range s.Periods {
			resp.Periods[i] = SubstitutionPeriodToResponse(&p)
		}
	}

	return resp
}

// SubstitutionPeriodToResponse converts a SubstitutionPeriod model to response DTO.
func SubstitutionPeriodToResponse(p *models.SubstitutionPeriod) SubstitutionPeriodResponse {
	resp := SubstitutionPeriodResponse{
		ID:               p.ID,
		PeriodSlotID:     p.PeriodSlotID,
		TimetableEntryID: p.TimetableEntryID,
		SubjectID:        p.SubjectID,
		SectionID:        p.SectionID,
		RoomNumber:       p.RoomNumber,
		Notes:            p.Notes,
	}

	if p.PeriodSlot != nil {
		resp.PeriodSlotName = p.PeriodSlot.Name
		resp.StartTime = p.PeriodSlot.StartTime
		resp.EndTime = p.PeriodSlot.EndTime
	}

	if p.Subject != nil {
		resp.SubjectName = p.Subject.Name
	}

	if p.Section != nil {
		resp.SectionName = p.Section.Name
		if p.Section.Class.ID != uuid.Nil {
			resp.ClassName = p.Section.Class.Name
		}
	}

	return resp
}
