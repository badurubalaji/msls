// Package staffdocument provides staff document management functionality.
package staffdocument

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"msls-backend/internal/pkg/database/models"
)

// MaxFileSize is the maximum file size allowed (5MB).
const MaxFileSize = 5 * 1024 * 1024

// AllowedMimeTypes contains the allowed file types.
var AllowedMimeTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
}

// Service handles business logic for staff documents.
type Service struct {
	repo *Repository
}

// NewService creates a new staff document service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ============================================================================
// Document Type Methods
// ============================================================================

// CreateDocumentType creates a new document type.
func (s *Service) CreateDocumentType(ctx context.Context, tenantID uuid.UUID, dto CreateDocumentTypeDTO) (*models.StaffDocumentType, error) {
	category := models.DocumentCategory(dto.Category)
	if !category.IsValid() {
		return nil, ErrInvalidCategory
	}

	dt := &models.StaffDocumentType{
		TenantID:              tenantID,
		Name:                  dto.Name,
		Code:                  strings.ToLower(strings.ReplaceAll(dto.Code, " ", "_")),
		Category:              category,
		Description:           dto.Description,
		IsMandatory:           dto.IsMandatory,
		HasExpiry:             dto.HasExpiry,
		DefaultValidityMonths: dto.DefaultValidityMonths,
		ApplicableTo:          dto.ApplicableTo,
		IsActive:              true,
		DisplayOrder:          dto.DisplayOrder,
	}

	if len(dt.ApplicableTo) == 0 {
		dt.ApplicableTo = []string{"teaching", "non_teaching"}
	}

	if err := s.repo.CreateDocumentType(ctx, dt); err != nil {
		return nil, err
	}
	return dt, nil
}

// GetDocumentType retrieves a document type by ID.
func (s *Service) GetDocumentType(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffDocumentType, error) {
	return s.repo.GetDocumentTypeByID(ctx, tenantID, id)
}

// UpdateDocumentType updates a document type.
func (s *Service) UpdateDocumentType(ctx context.Context, tenantID, id uuid.UUID, dto UpdateDocumentTypeDTO) (*models.StaffDocumentType, error) {
	dt, err := s.repo.GetDocumentTypeByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		dt.Name = *dto.Name
	}
	if dto.Description != nil {
		dt.Description = dto.Description
	}
	if dto.IsMandatory != nil {
		dt.IsMandatory = *dto.IsMandatory
	}
	if dto.HasExpiry != nil {
		dt.HasExpiry = *dto.HasExpiry
	}
	if dto.DefaultValidityMonths != nil {
		dt.DefaultValidityMonths = dto.DefaultValidityMonths
	}
	if len(dto.ApplicableTo) > 0 {
		dt.ApplicableTo = dto.ApplicableTo
	}
	if dto.IsActive != nil {
		dt.IsActive = *dto.IsActive
	}
	if dto.DisplayOrder != nil {
		dt.DisplayOrder = *dto.DisplayOrder
	}
	dt.UpdatedAt = time.Now()

	if err := s.repo.UpdateDocumentType(ctx, dt); err != nil {
		return nil, err
	}
	return dt, nil
}

// DeleteDocumentType deletes a document type.
func (s *Service) DeleteDocumentType(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteDocumentType(ctx, tenantID, id)
}

// ListDocumentTypes lists all document types for a tenant.
func (s *Service) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]models.StaffDocumentType, error) {
	return s.repo.ListDocumentTypes(ctx, tenantID, activeOnly)
}

// ============================================================================
// Document Methods
// ============================================================================

// ValidateFileUpload validates file upload parameters.
func (s *Service) ValidateFileUpload(fileName string, fileSize int, mimeType string) error {
	if fileSize > MaxFileSize {
		return ErrFileTooLarge
	}
	if !AllowedMimeTypes[mimeType] {
		return ErrInvalidFileType
	}
	return nil
}

// GenerateFilePath generates a storage path for a document.
func (s *Service) GenerateFilePath(tenantID, staffID, documentID uuid.UUID, fileName string) string {
	ext := filepath.Ext(fileName)
	return fmt.Sprintf("%s/staff/%s/documents/%s%s", tenantID.String(), staffID.String(), documentID.String(), ext)
}

// CreateDocument creates a new document record.
func (s *Service) CreateDocument(ctx context.Context, tenantID, staffID uuid.UUID, dto CreateDocumentDTO, fileName string, filePath string, fileSize int, mimeType string, createdBy *uuid.UUID) (*models.StaffDocument, error) {
	// Validate document type exists
	_, err := s.repo.GetDocumentTypeByID(ctx, tenantID, dto.DocumentTypeID)
	if err != nil {
		return nil, err
	}

	// Check for existing document of same type
	existing, err := s.repo.GetExistingDocument(ctx, tenantID, staffID, dto.DocumentTypeID)
	if err != nil {
		return nil, err
	}

	doc := &models.StaffDocument{
		TenantID:           tenantID,
		StaffID:            staffID,
		DocumentTypeID:     dto.DocumentTypeID,
		DocumentNumber:     dto.DocumentNumber,
		IssueDate:          dto.IssueDate,
		ExpiryDate:         dto.ExpiryDate,
		FileName:           fileName,
		FilePath:           filePath,
		FileSize:           fileSize,
		MimeType:           mimeType,
		VerificationStatus: models.VerificationStatusPending,
		Remarks:            dto.Remarks,
		IsCurrent:          true,
		CreatedBy:          createdBy,
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	// Mark previous documents as not current if exists
	if existing != nil {
		if err := s.repo.MarkPreviousDocumentsAsNotCurrent(ctx, tenantID, staffID, dto.DocumentTypeID, doc.ID); err != nil {
			return nil, err
		}
	}

	return doc, nil
}

// GetDocument retrieves a document by ID.
func (s *Service) GetDocument(ctx context.Context, tenantID, staffID, docID uuid.UUID) (*models.StaffDocument, error) {
	return s.repo.GetDocumentByStaffAndID(ctx, tenantID, staffID, docID)
}

// UpdateDocument updates a document's metadata.
func (s *Service) UpdateDocument(ctx context.Context, tenantID, staffID, docID uuid.UUID, dto UpdateDocumentDTO) (*models.StaffDocument, error) {
	doc, err := s.repo.GetDocumentByStaffAndID(ctx, tenantID, staffID, docID)
	if err != nil {
		return nil, err
	}

	if dto.DocumentNumber != nil {
		doc.DocumentNumber = dto.DocumentNumber
	}
	if dto.IssueDate != nil {
		doc.IssueDate = dto.IssueDate
	}
	if dto.ExpiryDate != nil {
		doc.ExpiryDate = dto.ExpiryDate
	}
	if dto.Remarks != nil {
		doc.Remarks = dto.Remarks
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// DeleteDocument deletes a document.
func (s *Service) DeleteDocument(ctx context.Context, tenantID, staffID, docID uuid.UUID) (*models.StaffDocument, error) {
	doc, err := s.repo.GetDocumentByStaffAndID(ctx, tenantID, staffID, docID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteDocument(ctx, tenantID, docID); err != nil {
		return nil, err
	}
	return doc, nil
}

// ListStaffDocuments lists documents for a staff member.
func (s *Service) ListStaffDocuments(ctx context.Context, tenantID, staffID uuid.UUID, params ListDocumentsParams) ([]models.StaffDocument, int64, error) {
	return s.repo.ListStaffDocuments(ctx, tenantID, staffID, params)
}

// ============================================================================
// Verification Methods
// ============================================================================

// VerifyDocument marks a document as verified.
func (s *Service) VerifyDocument(ctx context.Context, tenantID, staffID, docID uuid.UUID, verifierID uuid.UUID, dto VerifyDocumentDTO) (*models.StaffDocument, error) {
	doc, err := s.repo.GetDocumentByStaffAndID(ctx, tenantID, staffID, docID)
	if err != nil {
		return nil, err
	}

	if doc.VerificationStatus == models.VerificationStatusVerified {
		return nil, ErrDocumentAlreadyVerified
	}

	now := time.Now()
	doc.VerificationStatus = models.VerificationStatusVerified
	doc.VerifiedBy = &verifierID
	doc.VerifiedAt = &now
	doc.VerificationNotes = dto.Notes
	doc.RejectionReason = nil

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// RejectDocument marks a document as rejected.
func (s *Service) RejectDocument(ctx context.Context, tenantID, staffID, docID uuid.UUID, verifierID uuid.UUID, dto RejectDocumentDTO) (*models.StaffDocument, error) {
	if dto.Reason == "" {
		return nil, ErrMissingRejectionReason
	}

	doc, err := s.repo.GetDocumentByStaffAndID(ctx, tenantID, staffID, docID)
	if err != nil {
		return nil, err
	}

	if doc.VerificationStatus == models.VerificationStatusRejected {
		return nil, ErrDocumentAlreadyRejected
	}

	now := time.Now()
	doc.VerificationStatus = models.VerificationStatusRejected
	doc.VerifiedBy = &verifierID
	doc.VerifiedAt = &now
	doc.VerificationNotes = dto.Notes
	doc.RejectionReason = &dto.Reason

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// ============================================================================
// Expiry & Compliance Methods
// ============================================================================

// GetExpiringDocuments retrieves documents expiring within the given days.
func (s *Service) GetExpiringDocuments(ctx context.Context, tenantID uuid.UUID, days int) ([]ExpiringDocumentResponse, error) {
	results, err := s.repo.GetExpiringDocuments(ctx, tenantID, days)
	if err != nil {
		return nil, err
	}

	responses := make([]ExpiringDocumentResponse, len(results))
	for i, r := range results {
		responses[i] = ExpiringDocumentResponse{
			Document:     ToDocumentResponse(&r.StaffDocument),
			StaffName:    r.StaffName,
			EmployeeID:   r.EmployeeID,
			DaysToExpiry: r.DaysToExpiry,
		}
	}
	return responses, nil
}

// GetComplianceReport generates a compliance report.
func (s *Service) GetComplianceReport(ctx context.Context, tenantID uuid.UUID, includeStaffDetails bool) (*ComplianceReportResponse, error) {
	stats, err := s.repo.GetComplianceStats(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	byType, err := s.repo.GetComplianceByDocumentType(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	typeCompliance := make([]DocumentTypeCompliance, len(byType))
	for i, t := range byType {
		compliancePercent := 0.0
		if t.Required > 0 {
			compliancePercent = float64(t.Verified) / float64(t.Required) * 100
		}
		typeCompliance[i] = DocumentTypeCompliance{
			DocumentType: &DocumentTypeResponse{
				ID:          t.DocumentTypeID,
				Name:        t.DocumentTypeName,
				Code:        t.DocumentTypeCode,
				Category:    t.Category,
				IsMandatory: t.IsMandatory,
			},
			Required:          t.Required,
			Submitted:         t.Submitted,
			Verified:          t.Verified,
			Pending:           t.Pending,
			Rejected:          t.Rejected,
			Expired:           t.Expired,
			CompliancePercent: compliancePercent,
		}
	}

	report := &ComplianceReportResponse{
		Stats:          stats,
		ByDocumentType: typeCompliance,
	}

	if includeStaffDetails {
		staffDetails, err := s.repo.GetStaffComplianceDetails(ctx, tenantID, 100, 0)
		if err != nil {
			return nil, err
		}

		details := make([]StaffComplianceDetail, len(staffDetails))
		for i, sd := range staffDetails {
			compliancePercent := 0.0
			if sd.TotalRequired > 0 {
				compliancePercent = float64(sd.Verified) / float64(sd.TotalRequired) * 100
			}
			details[i] = StaffComplianceDetail{
				StaffID:           sd.StaffID,
				StaffName:         sd.StaffName,
				EmployeeID:        sd.EmployeeID,
				TotalRequired:     sd.TotalRequired,
				Submitted:         sd.Submitted,
				Verified:          sd.Verified,
				Pending:           sd.Pending,
				Rejected:          sd.Rejected,
				Expired:           sd.Expired,
				MissingDocuments:  sd.MissingDocuments,
				CompliancePercent: compliancePercent,
			}
		}
		report.StaffDetails = details
	}

	return report, nil
}

// ============================================================================
// Notification Methods
// ============================================================================

// SendExpiryNotifications sends notifications for expiring documents.
func (s *Service) SendExpiryNotifications(ctx context.Context, tenantID uuid.UUID) error {
	// Get documents expiring in 30 days
	docs30, err := s.repo.GetExpiringDocuments(ctx, tenantID, 30)
	if err != nil {
		return err
	}

	for _, doc := range docs30 {
		if doc.DaysToExpiry <= 7 {
			// Send 7-day notification if not already sent
			if err := s.sendNotificationIfNeeded(ctx, tenantID, doc.ID, models.NotificationTypeExpiry7Days); err != nil {
				continue // Log error but continue
			}
		} else if doc.DaysToExpiry <= 30 {
			// Send 30-day notification if not already sent
			if err := s.sendNotificationIfNeeded(ctx, tenantID, doc.ID, models.NotificationTypeExpiry30Days); err != nil {
				continue
			}
		}
	}

	// Get expired documents
	docsExpired, err := s.repo.GetExpiringDocuments(ctx, tenantID, 0)
	if err != nil {
		return err
	}

	for _, doc := range docsExpired {
		if doc.DaysToExpiry < 0 {
			if err := s.sendNotificationIfNeeded(ctx, tenantID, doc.ID, models.NotificationTypeExpired); err != nil {
				continue
			}
		}
	}

	return nil
}

func (s *Service) sendNotificationIfNeeded(ctx context.Context, tenantID, documentID uuid.UUID, notificationType models.NotificationType) error {
	// Check if notification already sent
	lastNotification, err := s.repo.GetLastNotification(ctx, documentID, notificationType)
	if err != nil {
		return err
	}

	// If notification sent within last 24 hours, skip
	if lastNotification != nil && time.Since(lastNotification.SentAt) < 24*time.Hour {
		return nil
	}

	// Create notification record
	notification := &models.StaffDocumentNotification{
		TenantID:         tenantID,
		DocumentID:       documentID,
		NotificationType: notificationType,
		SentAt:           time.Now(),
		SentTo:           pq.StringArray{}, // Would be populated with actual recipient IDs
	}

	return s.repo.CreateNotification(ctx, notification)
}
