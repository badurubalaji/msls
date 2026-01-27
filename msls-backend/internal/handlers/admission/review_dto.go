// Package admission provides HTTP handlers for admission management.
package admission

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/services/admission"
)

// =====================================================================
// Application Review DTOs
// =====================================================================

// CreateReviewRequest represents a request to create an application review.
type CreateReviewRequest struct {
	ReviewType string `json:"reviewType" binding:"required,oneof=initial_screening document_verification academic_review interview final_decision"`
	Status     string `json:"status" binding:"required,oneof=approved rejected pending_info escalated"`
	Comments   string `json:"comments" binding:"omitempty,max=2000"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *CreateReviewRequest) ToServiceRequest(tenantID, applicationID, reviewerID uuid.UUID) *admission.CreateReviewRequest {
	return &admission.CreateReviewRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		ReviewerID:    reviewerID,
		ReviewType:    admission.ReviewType(r.ReviewType),
		Status:        admission.ReviewStatus(r.Status),
		Comments:      r.Comments,
	}
}

// VerifyDocumentRequest represents a request to verify a document.
type VerifyDocumentRequest struct {
	Status  string `json:"status" binding:"required,oneof=pending verified rejected resubmit_required"`
	Remarks string `json:"remarks" binding:"omitempty,max=1000"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *VerifyDocumentRequest) ToServiceRequest(tenantID, applicationID, documentID, verifiedBy uuid.UUID) *admission.VerifyDocumentRequest {
	return &admission.VerifyDocumentRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		DocumentID:    documentID,
		VerifiedBy:    verifiedBy,
		Status:        admission.VerificationStatus(r.Status),
		Remarks:       r.Remarks,
	}
}

// UpdateApplicationStatusRequest represents a request to update application status.
type UpdateApplicationStatusRequest struct {
	Status  string `json:"status" binding:"required,oneof=draft submitted under_review documents_pending test_scheduled test_completed shortlisted approved rejected waitlisted enrolled"`
	Remarks string `json:"remarks" binding:"omitempty,max=1000"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *UpdateApplicationStatusRequest) ToServiceRequest(tenantID, applicationID, updatedBy uuid.UUID) *admission.UpdateApplicationStatusRequest {
	return &admission.UpdateApplicationStatusRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		Status:        admission.ApplicationStatus(r.Status),
		UpdatedBy:     updatedBy,
		Remarks:       r.Remarks,
	}
}

// =====================================================================
// Response DTOs
// =====================================================================

// ApplicationReviewResponse represents an application review in the response.
type ApplicationReviewResponse struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	ApplicationID string `json:"applicationId"`
	ReviewerID    string `json:"reviewerId"`
	ReviewerName  string `json:"reviewerName,omitempty"`
	ReviewType    string `json:"reviewType"`
	Status        string `json:"status"`
	Comments      string `json:"comments,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

// NewApplicationReviewResponse creates an ApplicationReviewResponse from an entity.
func NewApplicationReviewResponse(review *admission.ApplicationReview) ApplicationReviewResponse {
	return ApplicationReviewResponse{
		ID:            review.ID.String(),
		TenantID:      review.TenantID.String(),
		ApplicationID: review.ApplicationID.String(),
		ReviewerID:    review.ReviewerID.String(),
		ReviewType:    string(review.ReviewType),
		Status:        string(review.Status),
		Comments:      review.Comments,
		CreatedAt:     review.CreatedAt.Format(time.RFC3339),
	}
}

// ApplicationDocumentResponse represents an application document in the response.
type ApplicationDocumentResponse struct {
	ID                  string  `json:"id"`
	TenantID            string  `json:"tenantId"`
	ApplicationID       string  `json:"applicationId"`
	DocumentType        string  `json:"documentType"`
	DocumentName        string  `json:"documentName"`
	FilePath            string  `json:"filePath"`
	FileSize            int     `json:"fileSize"`
	MimeType            string  `json:"mimeType,omitempty"`
	VerificationStatus  string  `json:"verificationStatus"`
	VerifiedBy          *string `json:"verifiedBy,omitempty"`
	VerifiedAt          *string `json:"verifiedAt,omitempty"`
	VerificationRemarks string  `json:"verificationRemarks,omitempty"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// NewApplicationDocumentResponse creates an ApplicationDocumentResponse from an entity.
func NewApplicationDocumentResponse(doc *admission.ApplicationDocument) ApplicationDocumentResponse {
	resp := ApplicationDocumentResponse{
		ID:                  doc.ID.String(),
		TenantID:            doc.TenantID.String(),
		ApplicationID:       doc.ApplicationID.String(),
		DocumentType:        doc.DocumentType,
		DocumentName:        doc.DocumentName,
		FilePath:            doc.FilePath,
		FileSize:            doc.FileSize,
		MimeType:            doc.MimeType,
		VerificationStatus:  string(doc.VerificationStatus),
		VerificationRemarks: doc.VerificationRemarks,
		CreatedAt:           doc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           doc.UpdatedAt.Format(time.RFC3339),
	}

	if doc.VerifiedBy != nil {
		s := doc.VerifiedBy.String()
		resp.VerifiedBy = &s
	}
	if doc.VerifiedAt != nil {
		s := doc.VerifiedAt.Format(time.RFC3339)
		resp.VerifiedAt = &s
	}

	return resp
}

// AdmissionApplicationResponse represents an admission application in the response.
type AdmissionApplicationResponse struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	SessionID         string `json:"sessionId"`
	BranchID          string `json:"branchId,omitempty"`
	ApplicationNumber string `json:"applicationNumber"`

	// Student Information
	StudentName string `json:"studentName"`
	DateOfBirth string `json:"dateOfBirth"`
	Gender      string `json:"gender"`
	BloodGroup  string `json:"bloodGroup,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Religion    string `json:"religion,omitempty"`
	Category    string `json:"category,omitempty"`
	AadharNumber string `json:"aadharNumber,omitempty"`

	// Academic Information
	ClassApplying      string  `json:"classApplying"`
	PreviousSchool     string  `json:"previousSchool,omitempty"`
	PreviousClass      string  `json:"previousClass,omitempty"`
	PreviousPercentage *string `json:"previousPercentage,omitempty"`

	// Contact Information
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`

	// Parent Information
	FatherName       string `json:"fatherName,omitempty"`
	FatherPhone      string `json:"fatherPhone,omitempty"`
	FatherEmail      string `json:"fatherEmail,omitempty"`
	FatherOccupation string `json:"fatherOccupation,omitempty"`
	MotherName       string `json:"motherName,omitempty"`
	MotherPhone      string `json:"motherPhone,omitempty"`
	MotherEmail      string `json:"motherEmail,omitempty"`
	MotherOccupation string `json:"motherOccupation,omitempty"`

	// Status
	Status      string  `json:"status"`
	SubmittedAt *string `json:"submittedAt,omitempty"`
	Remarks     string  `json:"remarks,omitempty"`
	FeePaid     bool    `json:"feePaid"`

	// Timestamps
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// Related data
	Documents []ApplicationDocumentResponse `json:"documents,omitempty"`
	Reviews   []ApplicationReviewResponse   `json:"reviews,omitempty"`

	// Computed fields
	DocumentSummary *DocumentSummary `json:"documentSummary,omitempty"`
}

// DocumentSummary represents a summary of document verification.
type DocumentSummary struct {
	Total    int `json:"total"`
	Verified int `json:"verified"`
	Pending  int `json:"pending"`
	Rejected int `json:"rejected"`
}

// NewAdmissionApplicationResponse creates an AdmissionApplicationResponse from an entity.
func NewAdmissionApplicationResponse(app *admission.AdmissionApplication) AdmissionApplicationResponse {
	resp := AdmissionApplicationResponse{
		ID:                app.ID.String(),
		TenantID:          app.TenantID.String(),
		SessionID:         app.SessionID.String(),
		ApplicationNumber: app.ApplicationNumber,
		StudentName:       app.StudentName,
		DateOfBirth:       app.DateOfBirth.Format("2006-01-02"),
		Gender:            app.Gender,
		BloodGroup:        app.BloodGroup,
		Nationality:       app.Nationality,
		Religion:          app.Religion,
		Category:          app.Category,
		AadharNumber:      app.AadharNumber,
		ClassApplying:     app.ClassApplying,
		PreviousSchool:    app.PreviousSchool,
		PreviousClass:     app.PreviousClass,
		AddressLine1:      app.AddressLine1,
		AddressLine2:      app.AddressLine2,
		City:              app.City,
		State:             app.State,
		PostalCode:        app.PostalCode,
		Country:           app.Country,
		FatherName:        app.FatherName,
		FatherPhone:       app.FatherPhone,
		FatherEmail:       app.FatherEmail,
		FatherOccupation:  app.FatherOccupation,
		MotherName:        app.MotherName,
		MotherPhone:       app.MotherPhone,
		MotherEmail:       app.MotherEmail,
		MotherOccupation:  app.MotherOccupation,
		Status:            string(app.Status),
		Remarks:           app.Remarks,
		FeePaid:           app.FeePaid,
		CreatedAt:         app.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         app.UpdatedAt.Format(time.RFC3339),
	}

	if app.BranchID != nil {
		resp.BranchID = app.BranchID.String()
	}
	if app.PreviousPercentage != nil {
		s := app.PreviousPercentage.String()
		resp.PreviousPercentage = &s
	}
	if app.SubmittedAt != nil {
		s := app.SubmittedAt.Format(time.RFC3339)
		resp.SubmittedAt = &s
	}

	// Convert documents
	if len(app.Documents) > 0 {
		resp.Documents = make([]ApplicationDocumentResponse, len(app.Documents))
		docSummary := &DocumentSummary{Total: len(app.Documents)}
		for i, doc := range app.Documents {
			resp.Documents[i] = NewApplicationDocumentResponse(&doc)
			switch doc.VerificationStatus {
			case admission.VerificationStatusVerified:
				docSummary.Verified++
			case admission.VerificationStatusPending:
				docSummary.Pending++
			case admission.VerificationStatusRejected, admission.VerificationStatusResubmitRequired:
				docSummary.Rejected++
			}
		}
		resp.DocumentSummary = docSummary
	}

	// Convert reviews
	if len(app.Reviews) > 0 {
		resp.Reviews = make([]ApplicationReviewResponse, len(app.Reviews))
		for i, review := range app.Reviews {
			resp.Reviews[i] = NewApplicationReviewResponse(&review)
		}
	}

	return resp
}

// ReviewListResponse represents a list of reviews.
type ReviewListResponse struct {
	Reviews []ApplicationReviewResponse `json:"reviews"`
	Total   int                         `json:"total"`
}

// DocumentListResponse represents a list of documents.
type DocumentListResponse struct {
	Documents []ApplicationDocumentResponse `json:"documents"`
	Total     int                           `json:"total"`
}
