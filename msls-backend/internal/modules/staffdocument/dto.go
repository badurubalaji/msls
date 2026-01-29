// Package staffdocument provides staff document management functionality.
package staffdocument

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Document Type DTOs
// ============================================================================

// CreateDocumentTypeDTO represents the data for creating a document type.
type CreateDocumentTypeDTO struct {
	Name                  string   `json:"name" binding:"required,max=100"`
	Code                  string   `json:"code" binding:"required,max=50"`
	Category              string   `json:"category" binding:"required"`
	Description           *string  `json:"description"`
	IsMandatory           bool     `json:"is_mandatory"`
	HasExpiry             bool     `json:"has_expiry"`
	DefaultValidityMonths *int     `json:"default_validity_months"`
	ApplicableTo          []string `json:"applicable_to"` // teaching, non_teaching
	DisplayOrder          int      `json:"display_order"`
}

// UpdateDocumentTypeDTO represents the data for updating a document type.
type UpdateDocumentTypeDTO struct {
	Name                  *string  `json:"name" binding:"omitempty,max=100"`
	Description           *string  `json:"description"`
	IsMandatory           *bool    `json:"is_mandatory"`
	HasExpiry             *bool    `json:"has_expiry"`
	DefaultValidityMonths *int     `json:"default_validity_months"`
	ApplicableTo          []string `json:"applicable_to"`
	IsActive              *bool    `json:"is_active"`
	DisplayOrder          *int     `json:"display_order"`
}

// DocumentTypeResponse represents a document type in API responses.
type DocumentTypeResponse struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Code                  string    `json:"code"`
	Category              string    `json:"category"`
	Description           *string   `json:"description,omitempty"`
	IsMandatory           bool      `json:"is_mandatory"`
	HasExpiry             bool      `json:"has_expiry"`
	DefaultValidityMonths *int      `json:"default_validity_months,omitempty"`
	ApplicableTo          []string  `json:"applicable_to"`
	IsActive              bool      `json:"is_active"`
	DisplayOrder          int       `json:"display_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ToDocumentTypeResponse converts a model to a response.
func ToDocumentTypeResponse(dt *models.StaffDocumentType) *DocumentTypeResponse {
	if dt == nil {
		return nil
	}
	return &DocumentTypeResponse{
		ID:                    dt.ID,
		Name:                  dt.Name,
		Code:                  dt.Code,
		Category:              string(dt.Category),
		Description:           dt.Description,
		IsMandatory:           dt.IsMandatory,
		HasExpiry:             dt.HasExpiry,
		DefaultValidityMonths: dt.DefaultValidityMonths,
		ApplicableTo:          dt.ApplicableTo,
		IsActive:              dt.IsActive,
		DisplayOrder:          dt.DisplayOrder,
		CreatedAt:             dt.CreatedAt,
		UpdatedAt:             dt.UpdatedAt,
	}
}

// ToDocumentTypeResponses converts multiple models to responses.
func ToDocumentTypeResponses(dts []models.StaffDocumentType) []DocumentTypeResponse {
	responses := make([]DocumentTypeResponse, len(dts))
	for i, dt := range dts {
		responses[i] = *ToDocumentTypeResponse(&dt)
	}
	return responses
}

// ============================================================================
// Document DTOs
// ============================================================================

// CreateDocumentDTO represents the data for creating a document.
type CreateDocumentDTO struct {
	DocumentTypeID uuid.UUID  `json:"document_type_id" binding:"required"`
	DocumentNumber *string    `json:"document_number"`
	IssueDate      *time.Time `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	Remarks        *string    `json:"remarks"`
}

// UpdateDocumentDTO represents the data for updating a document.
type UpdateDocumentDTO struct {
	DocumentNumber *string    `json:"document_number"`
	IssueDate      *time.Time `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	Remarks        *string    `json:"remarks"`
}

// VerifyDocumentDTO represents the data for verifying a document.
type VerifyDocumentDTO struct {
	Notes *string `json:"notes"`
}

// RejectDocumentDTO represents the data for rejecting a document.
type RejectDocumentDTO struct {
	Reason string  `json:"reason" binding:"required"`
	Notes  *string `json:"notes"`
}

// DocumentResponse represents a document in API responses.
type DocumentResponse struct {
	ID                 uuid.UUID             `json:"id"`
	StaffID            uuid.UUID             `json:"staff_id"`
	DocumentType       *DocumentTypeResponse `json:"document_type,omitempty"`
	DocumentTypeID     uuid.UUID             `json:"document_type_id"`
	DocumentNumber     *string               `json:"document_number,omitempty"`
	IssueDate          *time.Time            `json:"issue_date,omitempty"`
	ExpiryDate         *time.Time            `json:"expiry_date,omitempty"`
	FileName           string                `json:"file_name"`
	FileSize           int                   `json:"file_size"`
	MimeType           string                `json:"mime_type"`
	VerificationStatus string                `json:"verification_status"`
	VerifiedBy         *uuid.UUID            `json:"verified_by,omitempty"`
	VerifiedAt         *time.Time            `json:"verified_at,omitempty"`
	VerificationNotes  *string               `json:"verification_notes,omitempty"`
	RejectionReason    *string               `json:"rejection_reason,omitempty"`
	Remarks            *string               `json:"remarks,omitempty"`
	IsCurrent          bool                  `json:"is_current"`
	IsExpired          bool                  `json:"is_expired"`
	IsExpiringSoon     bool                  `json:"is_expiring_soon"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

// ToDocumentResponse converts a model to a response.
func ToDocumentResponse(doc *models.StaffDocument) *DocumentResponse {
	if doc == nil {
		return nil
	}
	resp := &DocumentResponse{
		ID:                 doc.ID,
		StaffID:            doc.StaffID,
		DocumentTypeID:     doc.DocumentTypeID,
		DocumentNumber:     doc.DocumentNumber,
		IssueDate:          doc.IssueDate,
		ExpiryDate:         doc.ExpiryDate,
		FileName:           doc.FileName,
		FileSize:           doc.FileSize,
		MimeType:           doc.MimeType,
		VerificationStatus: string(doc.VerificationStatus),
		VerifiedBy:         doc.VerifiedBy,
		VerifiedAt:         doc.VerifiedAt,
		VerificationNotes:  doc.VerificationNotes,
		RejectionReason:    doc.RejectionReason,
		Remarks:            doc.Remarks,
		IsCurrent:          doc.IsCurrent,
		IsExpired:          doc.IsExpired(),
		IsExpiringSoon:     doc.IsExpiringSoon(30),
		CreatedAt:          doc.CreatedAt,
		UpdatedAt:          doc.UpdatedAt,
	}
	if doc.DocumentType != nil {
		resp.DocumentType = ToDocumentTypeResponse(doc.DocumentType)
	}
	return resp
}

// ToDocumentResponses converts multiple models to responses.
func ToDocumentResponses(docs []models.StaffDocument) []DocumentResponse {
	responses := make([]DocumentResponse, len(docs))
	for i, doc := range docs {
		responses[i] = *ToDocumentResponse(&doc)
	}
	return responses
}

// ============================================================================
// Compliance & Report DTOs
// ============================================================================

// ExpiringDocumentResponse represents a document that is expiring.
type ExpiringDocumentResponse struct {
	Document    *DocumentResponse `json:"document"`
	StaffName   string            `json:"staff_name"`
	EmployeeID  string            `json:"employee_id"`
	DaysToExpiry int              `json:"days_to_expiry"`
}

// ComplianceStats represents compliance statistics.
type ComplianceStats struct {
	TotalStaff            int     `json:"total_staff"`
	DocumentsSubmitted    int     `json:"documents_submitted"`
	PendingVerification   int     `json:"pending_verification"`
	Verified              int     `json:"verified"`
	Rejected              int     `json:"rejected"`
	Expired               int     `json:"expired"`
	ExpiringIn30Days      int     `json:"expiring_in_30_days"`
	ExpiringIn60Days      int     `json:"expiring_in_60_days"`
	ExpiringIn90Days      int     `json:"expiring_in_90_days"`
	CompliancePercentage  float64 `json:"compliance_percentage"`
}

// StaffComplianceDetail represents compliance details for a staff member.
type StaffComplianceDetail struct {
	StaffID           uuid.UUID `json:"staff_id"`
	StaffName         string    `json:"staff_name"`
	EmployeeID        string    `json:"employee_id"`
	TotalRequired     int       `json:"total_required"`
	Submitted         int       `json:"submitted"`
	Verified          int       `json:"verified"`
	Pending           int       `json:"pending"`
	Rejected          int       `json:"rejected"`
	Expired           int       `json:"expired"`
	MissingDocuments  []string  `json:"missing_documents"`
	CompliancePercent float64   `json:"compliance_percent"`
}

// ComplianceReportResponse represents the compliance report.
type ComplianceReportResponse struct {
	Stats         *ComplianceStats         `json:"stats"`
	ByDocumentType []DocumentTypeCompliance `json:"by_document_type"`
	StaffDetails  []StaffComplianceDetail  `json:"staff_details,omitempty"`
}

// DocumentTypeCompliance represents compliance for a document type.
type DocumentTypeCompliance struct {
	DocumentType       *DocumentTypeResponse `json:"document_type"`
	Required           int                   `json:"required"`
	Submitted          int                   `json:"submitted"`
	Verified           int                   `json:"verified"`
	Pending            int                   `json:"pending"`
	Rejected           int                   `json:"rejected"`
	Expired            int                   `json:"expired"`
	CompliancePercent  float64               `json:"compliance_percent"`
}

// ============================================================================
// List Response DTOs
// ============================================================================

// DocumentListResponse represents a paginated list of documents.
type DocumentListResponse struct {
	Documents  []DocumentResponse `json:"documents"`
	NextCursor string             `json:"next_cursor,omitempty"`
	HasMore    bool               `json:"has_more"`
	Total      int64              `json:"total"`
}

// DocumentTypeListResponse represents a list of document types.
type DocumentTypeListResponse struct {
	DocumentTypes []DocumentTypeResponse `json:"document_types"`
	Total         int64                  `json:"total"`
}
