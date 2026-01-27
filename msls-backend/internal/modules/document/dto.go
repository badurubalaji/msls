package document

import (
	"time"

	"github.com/google/uuid"
)

// DocumentStatus represents the verification status of a document.
type DocumentStatus string

const (
	StatusPendingVerification DocumentStatus = "pending_verification"
	StatusVerified            DocumentStatus = "verified"
	StatusRejected            DocumentStatus = "rejected"
)

// IsValid checks if the document status is valid.
func (s DocumentStatus) IsValid() bool {
	switch s {
	case StatusPendingVerification, StatusVerified, StatusRejected:
		return true
	}
	return false
}

// DocumentType represents a configurable document type.
type DocumentType struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"-"`
	Code              string     `json:"code"`
	Name              string     `json:"name"`
	Description       string     `json:"description,omitempty"`
	IsMandatory       bool       `json:"isMandatory"`
	HasExpiry         bool       `json:"hasExpiry"`
	AllowedExtensions string     `json:"allowedExtensions"`
	MaxSizeMB         int        `json:"maxSizeMb"`
	SortOrder         int        `json:"sortOrder"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// StudentDocument represents a document uploaded for a student.
type StudentDocument struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"-"`
	StudentID       uuid.UUID       `json:"studentId"`
	DocumentTypeID  uuid.UUID       `json:"documentTypeId"`
	DocumentType    *DocumentType   `json:"documentType,omitempty"`
	FileURL         string          `json:"fileUrl"`
	FileName        string          `json:"fileName"`
	FileSizeBytes   int             `json:"fileSizeBytes"`
	MimeType        string          `json:"mimeType"`
	DocumentNumber  string          `json:"documentNumber,omitempty"`
	IssueDate       *time.Time      `json:"issueDate,omitempty"`
	ExpiryDate      *time.Time      `json:"expiryDate,omitempty"`
	Status          DocumentStatus  `json:"status"`
	RejectionReason string          `json:"rejectionReason,omitempty"`
	VerifiedAt      *time.Time      `json:"verifiedAt,omitempty"`
	VerifiedBy      *uuid.UUID      `json:"verifiedBy,omitempty"`
	VerifierName    string          `json:"verifierName,omitempty"`
	UploadedAt      time.Time       `json:"uploadedAt"`
	UploadedBy      uuid.UUID       `json:"uploadedBy"`
	UploaderName    string          `json:"uploaderName,omitempty"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Version         int             `json:"version"`
}

// CreateDocumentTypeRequest represents the request to create a document type.
type CreateDocumentTypeRequest struct {
	Code              string `json:"code" binding:"required,max=50"`
	Name              string `json:"name" binding:"required,max=100"`
	Description       string `json:"description"`
	IsMandatory       bool   `json:"isMandatory"`
	HasExpiry         bool   `json:"hasExpiry"`
	AllowedExtensions string `json:"allowedExtensions"`
	MaxSizeMB         int    `json:"maxSizeMb"`
	SortOrder         int    `json:"sortOrder"`
}

// UpdateDocumentTypeRequest represents the request to update a document type.
type UpdateDocumentTypeRequest struct {
	Name              *string `json:"name" binding:"omitempty,max=100"`
	Description       *string `json:"description"`
	IsMandatory       *bool   `json:"isMandatory"`
	HasExpiry         *bool   `json:"hasExpiry"`
	AllowedExtensions *string `json:"allowedExtensions"`
	MaxSizeMB         *int    `json:"maxSizeMb"`
	SortOrder         *int    `json:"sortOrder"`
	IsActive          *bool   `json:"isActive"`
}

// UploadDocumentRequest represents the metadata for uploading a document.
type UploadDocumentRequest struct {
	DocumentTypeID uuid.UUID `form:"documentTypeId" binding:"required"`
	DocumentNumber string    `form:"documentNumber"`
	IssueDate      string    `form:"issueDate"`  // YYYY-MM-DD
	ExpiryDate     string    `form:"expiryDate"` // YYYY-MM-DD
}

// UpdateDocumentRequest represents the request to update document metadata.
type UpdateDocumentRequest struct {
	DocumentNumber *string `json:"documentNumber"`
	IssueDate      *string `json:"issueDate"`  // YYYY-MM-DD
	ExpiryDate     *string `json:"expiryDate"` // YYYY-MM-DD
	Version        int     `json:"version" binding:"required"`
}

// VerifyDocumentRequest represents the request to verify a document.
type VerifyDocumentRequest struct {
	Version int `json:"version" binding:"required"`
}

// RejectDocumentRequest represents the request to reject a document.
type RejectDocumentRequest struct {
	Reason  string `json:"reason" binding:"required,max=500"`
	Version int    `json:"version" binding:"required"`
}

// DocumentTypeListResponse represents a list of document types.
type DocumentTypeListResponse struct {
	DocumentTypes []DocumentType `json:"documentTypes"`
	Total         int            `json:"total"`
}

// DocumentListResponse represents a list of student documents.
type DocumentListResponse struct {
	Documents []StudentDocument `json:"documents"`
	Total     int               `json:"total"`
}

// ChecklistItem represents a document in the checklist.
type ChecklistItem struct {
	DocumentType DocumentType     `json:"documentType"`
	Document     *StudentDocument `json:"document,omitempty"`
	IsRequired   bool             `json:"isRequired"`
	IsUploaded   bool             `json:"isUploaded"`
	IsVerified   bool             `json:"isVerified"`
}

// DocumentChecklistResponse represents the document checklist for a student.
type DocumentChecklistResponse struct {
	Items             []ChecklistItem `json:"items"`
	TotalRequired     int             `json:"totalRequired"`
	TotalUploaded     int             `json:"totalUploaded"`
	TotalVerified     int             `json:"totalVerified"`
	TotalPending      int             `json:"totalPending"`
	TotalRejected     int             `json:"totalRejected"`
	CompletionPercent float64         `json:"completionPercent"`
}

// DocumentFilter contains filters for listing documents.
type DocumentFilter struct {
	Status         *DocumentStatus
	DocumentTypeID *uuid.UUID
}
