// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"time"

	"github.com/google/uuid"
)

// CreateEnrollmentDTO represents a request to create an enrollment.
type CreateEnrollmentDTO struct {
	TenantID       uuid.UUID
	StudentID      uuid.UUID
	AcademicYearID uuid.UUID
	ClassID        *uuid.UUID
	SectionID      *uuid.UUID
	RollNumber     string
	ClassTeacherID *uuid.UUID
	EnrollmentDate *time.Time
	Notes          string
	CreatedBy      *uuid.UUID
}

// UpdateEnrollmentDTO represents a request to update an enrollment.
type UpdateEnrollmentDTO struct {
	ClassID        *uuid.UUID
	SectionID      *uuid.UUID
	RollNumber     *string
	ClassTeacherID *uuid.UUID
	Notes          *string
	UpdatedBy      *uuid.UUID
}

// TransferDTO represents a request to process a student transfer.
type TransferDTO struct {
	TransferDate   time.Time
	TransferReason string
	UpdatedBy      *uuid.UUID
}

// DropoutDTO represents a request to process a student dropout.
type DropoutDTO struct {
	DropoutDate   time.Time
	DropoutReason string
	UpdatedBy     *uuid.UUID
}

// CompleteEnrollmentDTO represents a request to complete an enrollment.
type CompleteEnrollmentDTO struct {
	CompletionDate *time.Time
	UpdatedBy      *uuid.UUID
}

// EnrollmentResponse represents an enrollment in API responses.
type EnrollmentResponse struct {
	ID              string               `json:"id"`
	StudentID       string               `json:"studentId"`
	AcademicYear    *AcademicYearRefDTO  `json:"academicYear"`
	ClassID         string               `json:"classId,omitempty"`
	SectionID       string               `json:"sectionId,omitempty"`
	RollNumber      string               `json:"rollNumber,omitempty"`
	ClassTeacherID  string               `json:"classTeacherId,omitempty"`
	Status          string               `json:"status"`
	EnrollmentDate  string               `json:"enrollmentDate"`
	CompletionDate  string               `json:"completionDate,omitempty"`
	TransferDate    string               `json:"transferDate,omitempty"`
	TransferReason  string               `json:"transferReason,omitempty"`
	DropoutDate     string               `json:"dropoutDate,omitempty"`
	DropoutReason   string               `json:"dropoutReason,omitempty"`
	Notes           string               `json:"notes,omitempty"`
	CreatedAt       string               `json:"createdAt"`
	UpdatedAt       string               `json:"updatedAt"`
}

// AcademicYearRefDTO represents an academic year reference in responses.
type AcademicYearRefDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	IsCurrent bool   `json:"isCurrent"`
}

// EnrollmentHistoryResponse represents a list of enrollment history items.
type EnrollmentHistoryResponse struct {
	Enrollments []EnrollmentResponse `json:"enrollments"`
	Total       int                  `json:"total"`
}

// StatusChangeResponse represents a status change log entry in responses.
type StatusChangeResponse struct {
	ID           string `json:"id"`
	EnrollmentID string `json:"enrollmentId"`
	FromStatus   string `json:"fromStatus,omitempty"`
	ToStatus     string `json:"toStatus"`
	ChangeReason string `json:"changeReason,omitempty"`
	ChangeDate   string `json:"changeDate"`
	ChangedAt    string `json:"changedAt"`
	ChangedBy    string `json:"changedBy"`
}

// ToEnrollmentResponse converts a StudentEnrollment model to an EnrollmentResponse.
func ToEnrollmentResponse(e *StudentEnrollment) EnrollmentResponse {
	resp := EnrollmentResponse{
		ID:             e.ID.String(),
		StudentID:      e.StudentID.String(),
		Status:         string(e.Status),
		EnrollmentDate: e.EnrollmentDate.Format("2006-01-02"),
		CreatedAt:      e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      e.UpdatedAt.Format(time.RFC3339),
	}

	// Add academic year reference if loaded
	if e.AcademicYear != nil {
		resp.AcademicYear = &AcademicYearRefDTO{
			ID:        e.AcademicYear.ID.String(),
			Name:      e.AcademicYear.Name,
			StartDate: e.AcademicYear.StartDate.Format("2006-01-02"),
			EndDate:   e.AcademicYear.EndDate.Format("2006-01-02"),
			IsCurrent: e.AcademicYear.IsCurrent,
		}
	}

	// Add optional fields
	if e.ClassID != nil {
		resp.ClassID = e.ClassID.String()
	}
	if e.SectionID != nil {
		resp.SectionID = e.SectionID.String()
	}
	if e.RollNumber != "" {
		resp.RollNumber = e.RollNumber
	}
	if e.ClassTeacherID != nil {
		resp.ClassTeacherID = e.ClassTeacherID.String()
	}
	if e.CompletionDate != nil {
		resp.CompletionDate = e.CompletionDate.Format("2006-01-02")
	}
	if e.TransferDate != nil {
		resp.TransferDate = e.TransferDate.Format("2006-01-02")
	}
	if e.TransferReason != "" {
		resp.TransferReason = e.TransferReason
	}
	if e.DropoutDate != nil {
		resp.DropoutDate = e.DropoutDate.Format("2006-01-02")
	}
	if e.DropoutReason != "" {
		resp.DropoutReason = e.DropoutReason
	}
	if e.Notes != "" {
		resp.Notes = e.Notes
	}

	return resp
}

// ToEnrollmentResponses converts a slice of StudentEnrollment models to EnrollmentResponses.
func ToEnrollmentResponses(enrollments []StudentEnrollment) []EnrollmentResponse {
	responses := make([]EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		responses[i] = ToEnrollmentResponse(&e)
	}
	return responses
}

// ToStatusChangeResponse converts an EnrollmentStatusChange to a StatusChangeResponse.
func ToStatusChangeResponse(c *EnrollmentStatusChange) StatusChangeResponse {
	resp := StatusChangeResponse{
		ID:           c.ID.String(),
		EnrollmentID: c.EnrollmentID.String(),
		ToStatus:     string(c.ToStatus),
		ChangeDate:   c.ChangeDate.Format("2006-01-02"),
		ChangedAt:    c.ChangedAt.Format(time.RFC3339),
		ChangedBy:    c.ChangedBy.String(),
	}

	if c.FromStatus != nil {
		resp.FromStatus = string(*c.FromStatus)
	}
	if c.ChangeReason != "" {
		resp.ChangeReason = c.ChangeReason
	}

	return resp
}

// ToStatusChangeResponses converts a slice of EnrollmentStatusChange to StatusChangeResponses.
func ToStatusChangeResponses(changes []EnrollmentStatusChange) []StatusChangeResponse {
	responses := make([]StatusChangeResponse, len(changes))
	for i, c := range changes {
		responses[i] = ToStatusChangeResponse(&c)
	}
	return responses
}
