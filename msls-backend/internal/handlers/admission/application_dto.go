// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// =============================================================================
// Application DTOs
// =============================================================================

// ApplicationResponse represents an application in API responses.
type ApplicationResponse struct {
	ID                string                      `json:"id"`
	ApplicationNumber string                      `json:"applicationNumber"`
	StudentName       string                      `json:"studentName"`
	ClassName         string                      `json:"className"`
	Status            string                      `json:"status"`
	Source            string                      `json:"source"`
	ApplicantDetails  ApplicantDetailsDTO         `json:"applicantDetails"`
	ParentInfo        ParentGuardianInfoDTO       `json:"parentInfo"`
	ApplicationFee    decimal.Decimal             `json:"applicationFee"`
	FeePaid           bool                        `json:"feePaid"`
	PaymentDate       *time.Time                  `json:"paymentDate,omitempty"`
	PaymentReference  string                      `json:"paymentReference,omitempty"`
	SubmittedAt       *time.Time                  `json:"submittedAt,omitempty"`
	ReviewedAt        *time.Time                  `json:"reviewedAt,omitempty"`
	ReviewNotes       string                      `json:"reviewNotes,omitempty"`
	ApprovedAt        *time.Time                  `json:"approvedAt,omitempty"`
	EnrolledAt        *time.Time                  `json:"enrolledAt,omitempty"`
	WaitlistPosition  *int                        `json:"waitlistPosition,omitempty"`
	Session           *SessionSummaryDTO          `json:"session,omitempty"`
	Parents           []ApplicationParentResponse `json:"parents,omitempty"`
	Documents         []AppDocumentResponse `json:"documents,omitempty"`
	CreatedAt         time.Time                   `json:"createdAt"`
	UpdatedAt         time.Time                   `json:"updatedAt"`
}

// ApplicationListResponse represents the response for listing applications.
type ApplicationListResponse struct {
	Applications []ApplicationResponse `json:"applications"`
	Total        int64                 `json:"total"`
}

// SessionSummaryDTO represents a session summary in application responses.
type SessionSummaryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ApplicantDetailsDTO represents applicant details in requests/responses.
type ApplicantDetailsDTO struct {
	Gender         string `json:"gender,omitempty"`
	DateOfBirth    string `json:"dateOfBirth,omitempty"`
	Nationality    string `json:"nationality,omitempty"`
	Religion       string `json:"religion,omitempty"`
	Category       string `json:"category,omitempty"`
	BloodGroup     string `json:"bloodGroup,omitempty"`
	Address        string `json:"address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	PinCode        string `json:"pinCode,omitempty"`
	PreviousSchool string `json:"previousSchool,omitempty"`
	TransferReason string `json:"transferReason,omitempty"`
}

// ParentGuardianInfoDTO represents parent/guardian info in requests/responses.
type ParentGuardianInfoDTO struct {
	FatherName       string `json:"fatherName,omitempty"`
	FatherOccupation string `json:"fatherOccupation,omitempty"`
	FatherPhone      string `json:"fatherPhone,omitempty"`
	FatherEmail      string `json:"fatherEmail,omitempty"`
	MotherName       string `json:"motherName,omitempty"`
	MotherOccupation string `json:"motherOccupation,omitempty"`
	MotherPhone      string `json:"motherPhone,omitempty"`
	MotherEmail      string `json:"motherEmail,omitempty"`
	GuardianName     string `json:"guardianName,omitempty"`
	GuardianRelation string `json:"guardianRelation,omitempty"`
	GuardianPhone    string `json:"guardianPhone,omitempty"`
	GuardianEmail    string `json:"guardianEmail,omitempty"`
}

// CreateApplicationRequest represents a request to create an application.
type CreateApplicationRequest struct {
	SessionID        string               `json:"sessionId" binding:"required"`
	EnquiryID        *string              `json:"enquiryId,omitempty"`
	StudentName      string               `json:"studentName" binding:"required"`
	ClassName        string               `json:"className" binding:"required"`
	Source           string               `json:"source,omitempty"`
	ApplicantDetails *ApplicantDetailsDTO `json:"applicantDetails,omitempty"`
	ParentInfo       *ParentGuardianInfoDTO `json:"parentInfo,omitempty"`
}

// UpdateApplicationRequest represents a request to update an application.
type UpdateApplicationRequest struct {
	StudentName      *string               `json:"studentName,omitempty"`
	ClassName        *string               `json:"className,omitempty"`
	ApplicantDetails *ApplicantDetailsDTO  `json:"applicantDetails,omitempty"`
	ParentInfo       *ParentGuardianInfoDTO `json:"parentInfo,omitempty"`
}

// UpdateStageRequest represents a request to update the application stage.
type UpdateStageRequest struct {
	NewStage string `json:"newStage" binding:"required"`
	Remarks  string `json:"remarks,omitempty"`
}

// =============================================================================
// Parent DTOs
// =============================================================================

// ApplicationParentResponse represents a parent in API responses.
type ApplicationParentResponse struct {
	ID           string    `json:"id"`
	Relation     string    `json:"relation"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Occupation   string    `json:"occupation"`
	Education    string    `json:"education"`
	AnnualIncome string    `json:"annualIncome"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// CreateParentRequest represents a request to add a parent.
type CreateParentRequest struct {
	Relation     string `json:"relation" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Occupation   string `json:"occupation,omitempty"`
	Education    string `json:"education,omitempty"`
	AnnualIncome string `json:"annualIncome,omitempty"`
}

// UpdateParentRequest represents a request to update a parent.
type UpdateParentRequest struct {
	Relation     *string `json:"relation,omitempty"`
	Name         *string `json:"name,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	Occupation   *string `json:"occupation,omitempty"`
	Education    *string `json:"education,omitempty"`
	AnnualIncome *string `json:"annualIncome,omitempty"`
}

// =============================================================================
// Document DTOs
// =============================================================================

// AppDocumentResponse represents a document in API responses (for application module).
type AppDocumentResponse struct {
	ID              string     `json:"id"`
	DocumentType    string     `json:"documentType"`
	FileURL         string     `json:"fileUrl"`
	FileName        string     `json:"fileName"`
	FileSize        int64      `json:"fileSize"`
	MimeType        string     `json:"mimeType"`
	IsVerified      bool       `json:"isVerified"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	RejectionReason string     `json:"rejectionReason,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// CreateDocumentRequest represents a request to add a document.
type CreateDocumentRequest struct {
	DocumentType string `json:"documentType" binding:"required"`
	FileURL      string `json:"fileUrl" binding:"required"`
	FileName     string `json:"fileName" binding:"required"`
	FileSize     int64  `json:"fileSize"`
	MimeType     string `json:"mimeType"`
}

// AppVerifyDocumentRequest represents a request to verify a document (for application module).
type AppVerifyDocumentRequest struct {
	IsVerified      bool   `json:"isVerified"`
	RejectionReason string `json:"rejectionReason,omitempty"`
}

// =============================================================================
// Public Status Check DTOs
// =============================================================================

// StatusCheckRequest represents a request to check application status.
type StatusCheckRequest struct {
	ApplicationNumber string `json:"applicationNumber" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
}

// StatusCheckResponse represents the response for status check.
type StatusCheckResponse struct {
	ApplicationNumber string     `json:"applicationNumber"`
	StudentName       string     `json:"studentName"`
	ClassName         string     `json:"className"`
	Status            string     `json:"status"`
	SubmittedAt       *time.Time `json:"submittedAt,omitempty"`
	ReviewNotes       string     `json:"reviewNotes,omitempty"`
}

// =============================================================================
// Conversion helpers
// =============================================================================

// applicationToResponse converts a model to response DTO.
func applicationToResponse(app *models.AdmissionApplication) ApplicationResponse {
	// Format date of birth for response
	var dobStr string
	if app.DateOfBirth != nil {
		dobStr = app.DateOfBirth.Format("2006-01-02")
	}

	resp := ApplicationResponse{
		ID:                app.ID.String(),
		ApplicationNumber: app.ApplicationNumber,
		StudentName:       app.StudentName,
		ClassName:         app.ClassApplying,
		Status:            string(app.Status),
		Source:            "", // Source field removed from model
		ApplicantDetails: ApplicantDetailsDTO{
			Gender:         app.Gender,
			DateOfBirth:    dobStr,
			Nationality:    app.Nationality,
			Religion:       app.Religion,
			Category:       app.Category,
			BloodGroup:     app.BloodGroup,
			Address:        app.AddressLine1,
			City:           app.City,
			State:          app.State,
			PinCode:        app.PostalCode,
			PreviousSchool: app.PreviousSchool,
		},
		ParentInfo: ParentGuardianInfoDTO{
			FatherName:       app.FatherName,
			FatherOccupation: app.FatherOccupation,
			FatherPhone:      app.FatherPhone,
			FatherEmail:      app.FatherEmail,
			MotherName:       app.MotherName,
			MotherOccupation: app.MotherOccupation,
			MotherPhone:      app.MotherPhone,
			MotherEmail:      app.MotherEmail,
			GuardianName:     app.GuardianName,
			GuardianRelation: app.GuardianRelation,
			GuardianPhone:    app.GuardianPhone,
			GuardianEmail:    app.GuardianEmail,
		},
		ApplicationFee:   decimal.Zero, // Fee not stored in model anymore
		FeePaid:          app.FeePaid,
		PaymentDate:      app.PaymentDate,
		PaymentReference: app.PaymentReference,
		SubmittedAt:      app.SubmittedAt,
		CreatedAt:        app.CreatedAt,
		UpdatedAt:        app.UpdatedAt,
	}

	if app.Session != nil {
		resp.Session = &SessionSummaryDTO{
			ID:   app.Session.ID.String(),
			Name: app.Session.Name,
		}
	}

	return resp
}

// applicationToResponseWithRelations converts a model to response DTO including parents and documents.
func applicationToResponseWithRelations(app *models.AdmissionApplication, parents []models.ApplicationParent, documents []models.ApplicationDocument) ApplicationResponse {
	resp := applicationToResponse(app)

	if len(parents) > 0 {
		resp.Parents = make([]ApplicationParentResponse, len(parents))
		for i, p := range parents {
			resp.Parents[i] = parentToResponse(&p)
		}
	}

	if len(documents) > 0 {
		resp.Documents = make([]AppDocumentResponse, len(documents))
		for i, d := range documents {
			resp.Documents[i] = documentToResponse(&d)
		}
	}

	return resp
}

// applicationsToResponses converts a slice of models to response DTOs.
func applicationsToResponses(apps []models.AdmissionApplication) []ApplicationResponse {
	responses := make([]ApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = applicationToResponse(&app)
	}
	return responses
}

// parentToResponse converts a parent model to response DTO.
func parentToResponse(p *models.ApplicationParent) ApplicationParentResponse {
	return ApplicationParentResponse{
		ID:           p.ID.String(),
		Relation:     string(p.Relation),
		Name:         p.Name,
		Phone:        p.Phone,
		Email:        p.Email,
		Occupation:   p.Occupation,
		Education:    p.Education,
		AnnualIncome: p.AnnualIncome,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// documentToResponse converts a document model to response DTO.
func documentToResponse(d *models.ApplicationDocument) AppDocumentResponse {
	return AppDocumentResponse{
		ID:              d.ID.String(),
		DocumentType:    string(d.DocumentType),
		FileURL:         d.FileURL,
		FileName:        d.FileName,
		FileSize:        d.FileSize,
		MimeType:        d.MimeType,
		IsVerified:      d.IsVerified(),
		VerifiedAt:      d.VerifiedAt,
		RejectionReason: d.RejectionReason(),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// Note: applicantDetailsToDTO and parentInfoToDTO are now inline in applicationToResponse
// since the model uses flat fields instead of nested JSONB structures.

// Ensure uuid.UUID is used to avoid import error.
var _ uuid.UUID
