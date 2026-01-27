// Package admission provides HTTP handlers for admission management.
package admission

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/services/admission"
)

// =====================================================================
// Entrance Test DTOs
// =====================================================================

// SubjectMarksDTO represents marks configuration for a subject.
type SubjectMarksDTO struct {
	Subject  string          `json:"subject" binding:"required"`
	MaxMarks decimal.Decimal `json:"maxMarks" binding:"required"`
}

// CreateTestRequest represents a request to create an entrance test.
type CreateTestRequest struct {
	SessionID         string            `json:"sessionId" binding:"required,uuid"`
	TestName          string            `json:"testName" binding:"required,min=3,max=200"`
	TestDate          string            `json:"testDate" binding:"required"`
	StartTime         string            `json:"startTime" binding:"required"`
	DurationMinutes   int               `json:"durationMinutes" binding:"omitempty,min=15,max=480"`
	Venue             string            `json:"venue" binding:"omitempty,max=200"`
	ClassNames        []string          `json:"classNames" binding:"required,min=1"`
	MaxCandidates     int               `json:"maxCandidates" binding:"omitempty,min=1"`
	Subjects          []SubjectMarksDTO `json:"subjects" binding:"required,min=1,dive"`
	Instructions      string            `json:"instructions" binding:"omitempty"`
	PassingPercentage *decimal.Decimal  `json:"passingPercentage" binding:"omitempty"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *CreateTestRequest) ToServiceRequest(tenantID uuid.UUID, userID *uuid.UUID) (*admission.CreateTestRequest, error) {
	sessionID, err := uuid.Parse(r.SessionID)
	if err != nil {
		return nil, err
	}

	testDate, err := time.Parse("2006-01-02", r.TestDate)
	if err != nil {
		return nil, err
	}

	subjects := make([]admission.SubjectMarks, len(r.Subjects))
	for i, s := range r.Subjects {
		subjects[i] = admission.SubjectMarks{
			Subject:  s.Subject,
			MaxMarks: s.MaxMarks,
		}
	}

	passingPct := decimal.NewFromFloat(33.0)
	if r.PassingPercentage != nil {
		passingPct = *r.PassingPercentage
	}

	return &admission.CreateTestRequest{
		TenantID:          tenantID,
		SessionID:         sessionID,
		TestName:          r.TestName,
		TestDate:          testDate,
		StartTime:         r.StartTime,
		DurationMinutes:   r.DurationMinutes,
		Venue:             r.Venue,
		ClassNames:        r.ClassNames,
		MaxCandidates:     r.MaxCandidates,
		Subjects:          subjects,
		Instructions:      r.Instructions,
		PassingPercentage: passingPct,
		CreatedBy:         userID,
	}, nil
}

// UpdateTestRequest represents a request to update an entrance test.
type UpdateTestRequest struct {
	TestName          *string           `json:"testName" binding:"omitempty,min=3,max=200"`
	TestDate          *string           `json:"testDate" binding:"omitempty"`
	StartTime         *string           `json:"startTime" binding:"omitempty"`
	DurationMinutes   *int              `json:"durationMinutes" binding:"omitempty,min=15,max=480"`
	Venue             *string           `json:"venue" binding:"omitempty,max=200"`
	ClassNames        []string          `json:"classNames" binding:"omitempty"`
	MaxCandidates     *int              `json:"maxCandidates" binding:"omitempty,min=1"`
	Subjects          []SubjectMarksDTO `json:"subjects" binding:"omitempty,dive"`
	Instructions      *string           `json:"instructions" binding:"omitempty"`
	PassingPercentage *decimal.Decimal  `json:"passingPercentage" binding:"omitempty"`
	Status            *string           `json:"status" binding:"omitempty,oneof=scheduled in_progress completed cancelled"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *UpdateTestRequest) ToServiceRequest(userID *uuid.UUID) (*admission.UpdateTestRequest, error) {
	req := &admission.UpdateTestRequest{
		UpdatedBy: userID,
	}

	if r.TestName != nil {
		req.TestName = r.TestName
	}
	if r.TestDate != nil {
		testDate, err := time.Parse("2006-01-02", *r.TestDate)
		if err != nil {
			return nil, err
		}
		req.TestDate = &testDate
	}
	if r.StartTime != nil {
		req.StartTime = r.StartTime
	}
	if r.DurationMinutes != nil {
		req.DurationMinutes = r.DurationMinutes
	}
	if r.Venue != nil {
		req.Venue = r.Venue
	}
	if r.ClassNames != nil {
		req.ClassNames = r.ClassNames
	}
	if r.MaxCandidates != nil {
		req.MaxCandidates = r.MaxCandidates
	}
	if r.Subjects != nil {
		subjects := make([]admission.SubjectMarks, len(r.Subjects))
		for i, s := range r.Subjects {
			subjects[i] = admission.SubjectMarks{
				Subject:  s.Subject,
				MaxMarks: s.MaxMarks,
			}
		}
		req.Subjects = subjects
	}
	if r.Instructions != nil {
		req.Instructions = r.Instructions
	}
	if r.PassingPercentage != nil {
		req.PassingPercentage = r.PassingPercentage
	}
	if r.Status != nil {
		status := admission.EntranceTestStatus(*r.Status)
		req.Status = &status
	}

	return req, nil
}

// RegisterCandidateRequest represents a request to register a candidate.
type RegisterCandidateRequest struct {
	ApplicationID string `json:"applicationId" binding:"required,uuid"`
}

// SubmitResultRequest represents a request to submit a single result.
type SubmitResultRequest struct {
	RegistrationID string             `json:"registrationId" binding:"required,uuid"`
	Marks          map[string]float64 `json:"marks" binding:"required"`
	Remarks        string             `json:"remarks" binding:"omitempty"`
}

// BulkSubmitResultsRequest represents a request to submit multiple results.
type BulkSubmitResultsRequest struct {
	Results []SubmitResultRequest `json:"results" binding:"required,min=1,dive"`
}

// =====================================================================
// Response DTOs
// =====================================================================

// TestResponse represents an entrance test in the response.
type TestResponse struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenantId"`
	SessionID         string            `json:"sessionId"`
	TestName          string            `json:"testName"`
	TestDate          string            `json:"testDate"`
	StartTime         string            `json:"startTime"`
	DurationMinutes   int               `json:"durationMinutes"`
	Venue             string            `json:"venue,omitempty"`
	ClassNames        []string          `json:"classNames"`
	MaxCandidates     int               `json:"maxCandidates"`
	Status            string            `json:"status"`
	Subjects          []SubjectMarksDTO `json:"subjects"`
	Instructions      string            `json:"instructions,omitempty"`
	PassingPercentage string            `json:"passingPercentage"`
	RegisteredCount   int64             `json:"registeredCount"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

// NewTestResponse creates a TestResponse from an EntranceTest entity.
func NewTestResponse(test *admission.EntranceTest, registeredCount int64) TestResponse {
	subjects := make([]SubjectMarksDTO, len(test.Subjects))
	for i, s := range test.Subjects {
		subjects[i] = SubjectMarksDTO{
			Subject:  s.Subject,
			MaxMarks: s.MaxMarks,
		}
	}

	return TestResponse{
		ID:                test.ID.String(),
		TenantID:          test.TenantID.String(),
		SessionID:         test.SessionID.String(),
		TestName:          test.TestName,
		TestDate:          test.TestDate.Format("2006-01-02"),
		StartTime:         test.StartTime,
		DurationMinutes:   test.DurationMinutes,
		Venue:             test.Venue,
		ClassNames:        test.ClassNames,
		MaxCandidates:     test.MaxCandidates,
		Status:            string(test.Status),
		Subjects:          subjects,
		Instructions:      test.Instructions,
		PassingPercentage: test.PassingPercentage.String(),
		RegisteredCount:   registeredCount,
		CreatedAt:         test.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         test.UpdatedAt.Format(time.RFC3339),
	}
}

// RegistrationApplicationResponse represents application info in registration.
type RegistrationApplicationResponse struct {
	ID                string `json:"id"`
	ApplicationNumber string `json:"applicationNumber"`
	StudentName       string `json:"studentName"`
	FatherName        string `json:"fatherName,omitempty"`
	FatherPhone       string `json:"fatherPhone,omitempty"`
	ClassApplying     string `json:"classApplying"`
	Status            string `json:"status"`
}

// TestRegistrationResponse represents a test registration in the response.
type TestRegistrationResponse struct {
	ID                    string                           `json:"id"`
	TenantID              string                           `json:"tenantId"`
	TestID                string                           `json:"testId"`
	ApplicationID         string                           `json:"applicationId"`
	RollNumber            string                           `json:"rollNumber,omitempty"`
	Status                string                           `json:"status"`
	Marks                 map[string]string                `json:"marks,omitempty"`
	TotalMarks            *string                          `json:"totalMarks,omitempty"`
	MaxMarks              *string                          `json:"maxMarks,omitempty"`
	Percentage            *string                          `json:"percentage,omitempty"`
	Result                *string                          `json:"result,omitempty"`
	Remarks               string                           `json:"remarks,omitempty"`
	HallTicketGeneratedAt *string                          `json:"hallTicketGeneratedAt,omitempty"`
	CreatedAt             string                           `json:"createdAt"`
	UpdatedAt             string                           `json:"updatedAt"`
	Application           *RegistrationApplicationResponse `json:"application,omitempty"`
}

// NewTestRegistrationResponse creates a TestRegistrationResponse from a TestRegistration entity.
func NewTestRegistrationResponse(reg *admission.TestRegistration) TestRegistrationResponse {
	resp := TestRegistrationResponse{
		ID:            reg.ID.String(),
		TenantID:      reg.TenantID.String(),
		TestID:        reg.TestID.String(),
		ApplicationID: reg.ApplicationID.String(),
		RollNumber:    reg.RollNumber,
		Status:        string(reg.Status),
		Remarks:       reg.Remarks,
		CreatedAt:     reg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     reg.UpdatedAt.Format(time.RFC3339),
	}

	// Convert marks map
	if len(reg.Marks) > 0 {
		resp.Marks = make(map[string]string)
		for subject, marks := range reg.Marks {
			resp.Marks[subject] = marks.String()
		}
	}

	// Convert decimal fields
	if reg.TotalMarks != nil {
		s := reg.TotalMarks.String()
		resp.TotalMarks = &s
	}
	if reg.MaxMarks != nil {
		s := reg.MaxMarks.String()
		resp.MaxMarks = &s
	}
	if reg.Percentage != nil {
		s := reg.Percentage.String()
		resp.Percentage = &s
	}
	if reg.Result != nil {
		s := string(*reg.Result)
		resp.Result = &s
	}
	if reg.HallTicketGeneratedAt != nil {
		s := reg.HallTicketGeneratedAt.Format(time.RFC3339)
		resp.HallTicketGeneratedAt = &s
	}

	// Add application info if available
	if reg.Application != nil {
		resp.Application = &RegistrationApplicationResponse{
			ID:                reg.Application.ID.String(),
			ApplicationNumber: reg.Application.ApplicationNumber,
			StudentName:       reg.Application.StudentName,
			FatherName:        reg.Application.FatherName,
			FatherPhone:       reg.Application.FatherPhone,
			ClassApplying:     reg.Application.ClassApplying,
			Status:            string(reg.Application.Status),
		}
	}

	return resp
}

// HallTicketResponse represents hall ticket data in the response.
type HallTicketResponse struct {
	Registration TestRegistrationResponse `json:"registration"`
	Test         TestResponse             `json:"test"`
	Student      struct {
		Name        string `json:"name"`
		DateOfBirth string `json:"dateOfBirth"`
		Gender      string `json:"gender"`
		Photo       string `json:"photo,omitempty"`
	} `json:"student"`
	Parent struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email,omitempty"`
	} `json:"parent"`
	GeneratedAt string `json:"generatedAt"`
}

// NewHallTicketResponse creates a HallTicketResponse from HallTicketData.
func NewHallTicketResponse(data *admission.HallTicketData) HallTicketResponse {
	resp := HallTicketResponse{
		Registration: NewTestRegistrationResponse(&data.Registration),
		Test:         NewTestResponse(&data.Test, 0),
		GeneratedAt:  data.GeneratedAt.Format(time.RFC3339),
	}

	// Add student info
	resp.Student.Name = data.Application.StudentName
	resp.Student.DateOfBirth = data.Application.DateOfBirth.Format("2006-01-02")
	resp.Student.Gender = data.Application.Gender

	// Add parent info
	resp.Parent.Name = data.Application.FatherName
	if resp.Parent.Name == "" {
		resp.Parent.Name = data.Application.MotherName
	}
	resp.Parent.Phone = data.Application.FatherPhone
	if resp.Parent.Phone == "" {
		resp.Parent.Phone = data.Application.MotherPhone
	}
	resp.Parent.Email = data.Application.FatherEmail
	if resp.Parent.Email == "" {
		resp.Parent.Email = data.Application.MotherEmail
	}

	return resp
}

// TestListResponse represents a list of tests.
type TestListResponse struct {
	Tests []TestResponse `json:"tests"`
	Total int            `json:"total"`
}
