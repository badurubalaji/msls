// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/services/admission"
)

// ============================================================================
// Request DTOs
// ============================================================================

// CreateEnquiryRequest represents the request body for creating an enquiry.
type CreateEnquiryRequest struct {
	BranchID        *string `json:"branchId"`
	SessionID       *string `json:"sessionId"`
	StudentName     string  `json:"studentName" binding:"required,max=200"`
	DateOfBirth     *string `json:"dateOfBirth"`
	Gender          string  `json:"gender" binding:"omitempty,oneof=male female other"`
	ClassApplying   string  `json:"classApplying" binding:"required,max=50"`
	ParentName      string  `json:"parentName" binding:"required,max=200"`
	ParentPhone     string  `json:"parentPhone" binding:"required,max=20"`
	ParentEmail     string  `json:"parentEmail" binding:"omitempty,email,max=255"`
	Source          string  `json:"source" binding:"omitempty,oneof=walk_in phone website referral advertisement social_media other"`
	ReferralDetails string  `json:"referralDetails" binding:"max=500"`
	Remarks         string  `json:"remarks" binding:"max=1000"`
	FollowUpDate    *string `json:"followUpDate"`
	AssignedTo      *string `json:"assignedTo"`
}

// UpdateEnquiryRequest represents the request body for updating an enquiry.
type UpdateEnquiryRequest struct {
	BranchID        *string `json:"branchId"`
	SessionID       *string `json:"sessionId"`
	StudentName     *string `json:"studentName" binding:"omitempty,max=200"`
	DateOfBirth     *string `json:"dateOfBirth"`
	Gender          *string `json:"gender" binding:"omitempty,oneof=male female other ''"`
	ClassApplying   *string `json:"classApplying" binding:"omitempty,max=50"`
	ParentName      *string `json:"parentName" binding:"omitempty,max=200"`
	ParentPhone     *string `json:"parentPhone" binding:"omitempty,max=20"`
	ParentEmail     *string `json:"parentEmail" binding:"omitempty,email,max=255"`
	Source          *string `json:"source" binding:"omitempty,oneof=walk_in phone website referral advertisement social_media other"`
	ReferralDetails *string `json:"referralDetails" binding:"omitempty,max=500"`
	Remarks         *string `json:"remarks" binding:"omitempty,max=1000"`
	Status          *string `json:"status" binding:"omitempty,oneof=new contacted interested converted closed"`
	FollowUpDate    *string `json:"followUpDate"`
	AssignedTo      *string `json:"assignedTo"`
}

// CreateFollowUpRequest represents the request body for creating a follow-up.
type CreateFollowUpRequest struct {
	FollowUpDate string  `json:"followUpDate" binding:"required"`
	ContactMode  string  `json:"contactMode" binding:"omitempty,oneof=phone email whatsapp in_person other"`
	Notes        string  `json:"notes" binding:"max=2000"`
	Outcome      *string `json:"outcome" binding:"omitempty,oneof=interested not_interested follow_up_required converted no_response"`
	NextFollowUp *string `json:"nextFollowUp"`
}

// ConvertEnquiryRequest represents the request body for converting an enquiry.
type ConvertEnquiryRequest struct {
	SessionID string `json:"sessionId,omitempty" binding:"omitempty,uuid"`
	BranchID  string `json:"branchId,omitempty" binding:"omitempty,uuid"`
	Remarks   string `json:"remarks,omitempty"`
}

// ListEnquiriesParams represents query parameters for listing enquiries.
type ListEnquiriesParams struct {
	BranchID      string `form:"branchId"`
	SessionID     string `form:"sessionId"`
	Status        string `form:"status"`
	Source        string `form:"source"`
	ClassApplying string `form:"classApplying"`
	Search        string `form:"search"`
	StartDate     string `form:"startDate"`
	EndDate       string `form:"endDate"`
	AssignedTo    string `form:"assignedTo"`
	Page          int    `form:"page"`
	PageSize      int    `form:"pageSize"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// EnquiryResponse represents an enquiry in API responses.
type EnquiryResponse struct {
	ID                     string  `json:"id"`
	EnquiryNumber          string  `json:"enquiryNumber"`
	BranchID               *string `json:"branchId,omitempty"`
	SessionID              *string `json:"sessionId,omitempty"`
	StudentName            string  `json:"studentName"`
	DateOfBirth            *string `json:"dateOfBirth,omitempty"`
	Gender                 string  `json:"gender,omitempty"`
	ClassApplying          string  `json:"classApplying"`
	ParentName             string  `json:"parentName"`
	ParentPhone            string  `json:"parentPhone"`
	ParentEmail            string  `json:"parentEmail,omitempty"`
	Source                 string  `json:"source"`
	ReferralDetails        string  `json:"referralDetails,omitempty"`
	Remarks                string  `json:"remarks,omitempty"`
	Status                 string  `json:"status"`
	FollowUpDate           *string `json:"followUpDate,omitempty"`
	AssignedTo             *string `json:"assignedTo,omitempty"`
	ConvertedApplicationID *string `json:"convertedApplicationId,omitempty"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
	CreatedBy              *string `json:"createdBy,omitempty"`
	UpdatedBy              *string `json:"updatedBy,omitempty"`
}

// FollowUpResponse represents a follow-up in API responses.
type FollowUpResponse struct {
	ID           string  `json:"id"`
	EnquiryID    string  `json:"enquiryId"`
	FollowUpDate string  `json:"followUpDate"`
	ContactMode  string  `json:"contactMode"`
	Notes        string  `json:"notes,omitempty"`
	Outcome      *string `json:"outcome,omitempty"`
	NextFollowUp *string `json:"nextFollowUp,omitempty"`
	CreatedAt    string  `json:"createdAt"`
	CreatedBy    *string `json:"createdBy,omitempty"`
}

// EnquiryListResponse represents a list of enquiries with pagination.
type EnquiryListResponse struct {
	Enquiries []EnquiryResponse `json:"enquiries"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"pageSize"`
}

// FollowUpListResponse represents a list of follow-ups.
type FollowUpListResponse struct {
	FollowUps []FollowUpResponse `json:"followUps"`
	Total     int                `json:"total"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// enquiryToResponse converts an AdmissionEnquiry model to EnquiryResponse.
func enquiryToResponse(e *admission.AdmissionEnquiry) EnquiryResponse {
	resp := EnquiryResponse{
		ID:              e.ID.String(),
		EnquiryNumber:   e.EnquiryNumber,
		StudentName:     e.StudentName,
		Gender:          e.Gender,
		ClassApplying:   e.ClassApplying,
		ParentName:      e.ParentName,
		ParentPhone:     e.ParentPhone,
		ParentEmail:     e.ParentEmail,
		Source:          string(e.Source),
		ReferralDetails: e.ReferralDetails,
		Remarks:         e.Remarks,
		Status:          string(e.Status),
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
	}

	if e.BranchID != nil {
		branchID := e.BranchID.String()
		resp.BranchID = &branchID
	}

	if e.SessionID != nil {
		sessionID := e.SessionID.String()
		resp.SessionID = &sessionID
	}

	if e.DateOfBirth != nil {
		dob := e.DateOfBirth.Format("2006-01-02")
		resp.DateOfBirth = &dob
	}

	if e.FollowUpDate != nil {
		followUp := e.FollowUpDate.Format("2006-01-02")
		resp.FollowUpDate = &followUp
	}

	if e.AssignedTo != nil {
		assignedTo := e.AssignedTo.String()
		resp.AssignedTo = &assignedTo
	}

	if e.ConvertedApplicationID != nil {
		appID := e.ConvertedApplicationID.String()
		resp.ConvertedApplicationID = &appID
	}

	if e.CreatedBy != nil {
		createdBy := e.CreatedBy.String()
		resp.CreatedBy = &createdBy
	}

	if e.UpdatedBy != nil {
		updatedBy := e.UpdatedBy.String()
		resp.UpdatedBy = &updatedBy
	}

	return resp
}

// enquiriesToResponses converts a slice of AdmissionEnquiry models to EnquiryResponses.
func enquiriesToResponses(enquiries []admission.AdmissionEnquiry) []EnquiryResponse {
	responses := make([]EnquiryResponse, len(enquiries))
	for i, e := range enquiries {
		responses[i] = enquiryToResponse(&e)
	}
	return responses
}

// followUpToResponse converts an EnquiryFollowUp model to FollowUpResponse.
func followUpToResponse(f *admission.EnquiryFollowUp) FollowUpResponse {
	resp := FollowUpResponse{
		ID:           f.ID.String(),
		EnquiryID:    f.EnquiryID.String(),
		FollowUpDate: f.FollowUpDate.Format("2006-01-02"),
		ContactMode:  string(f.ContactMode),
		Notes:        f.Notes,
		CreatedAt:    f.CreatedAt.Format(time.RFC3339),
	}

	if f.Outcome != nil {
		outcome := string(*f.Outcome)
		resp.Outcome = &outcome
	}

	if f.NextFollowUp != nil {
		nextFollowUp := f.NextFollowUp.Format("2006-01-02")
		resp.NextFollowUp = &nextFollowUp
	}

	if f.CreatedBy != nil {
		createdBy := f.CreatedBy.String()
		resp.CreatedBy = &createdBy
	}

	return resp
}

// followUpsToResponses converts a slice of EnquiryFollowUp models to FollowUpResponses.
func followUpsToResponses(followUps []admission.EnquiryFollowUp) []FollowUpResponse {
	responses := make([]FollowUpResponse, len(followUps))
	for i, f := range followUps {
		responses[i] = followUpToResponse(&f)
	}
	return responses
}

// parseUUID parses a string to UUID, returning nil if empty.
func parseUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

// parseDate parses a date string (YYYY-MM-DD) to time.Time, returning nil if empty.
func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}
