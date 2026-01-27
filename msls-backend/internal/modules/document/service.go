package document

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"msls-backend/internal/pkg/logger"
)

// Service handles document business logic.
type Service struct {
	repo       *Repository
	uploadPath string
}

// NewService creates a new document service.
func NewService(repo *Repository, uploadPath string) *Service {
	if uploadPath == "" {
		uploadPath = "uploads"
	}
	return &Service{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

// =============================================================================
// Document Type Operations
// =============================================================================

// CreateDocumentType creates a new document type.
func (s *Service) CreateDocumentType(ctx context.Context, tenantID uuid.UUID, req CreateDocumentTypeRequest) (*DocumentType, error) {
	// Set defaults
	allowedExtensions := req.AllowedExtensions
	if allowedExtensions == "" {
		allowedExtensions = "pdf,jpg,jpeg,png"
	}
	maxSizeMB := req.MaxSizeMB
	if maxSizeMB <= 0 {
		maxSizeMB = 5
	}

	dt := &DocumentType{
		TenantID:          tenantID,
		Code:              strings.ToLower(strings.TrimSpace(req.Code)),
		Name:              strings.TrimSpace(req.Name),
		Description:       strings.TrimSpace(req.Description),
		IsMandatory:       req.IsMandatory,
		HasExpiry:         req.HasExpiry,
		AllowedExtensions: allowedExtensions,
		MaxSizeMB:         maxSizeMB,
		SortOrder:         req.SortOrder,
		IsActive:          true,
	}

	if err := s.repo.CreateDocumentType(ctx, dt); err != nil {
		logger.Error("Failed to create document type",
			zap.String("tenant_id", tenantID.String()),
			zap.String("code", dt.Code),
			zap.Error(err))
		return nil, err
	}

	return dt, nil
}

// GetDocumentTypeByID gets a document type by ID.
func (s *Service) GetDocumentTypeByID(ctx context.Context, tenantID, id uuid.UUID) (*DocumentType, error) {
	return s.repo.GetDocumentTypeByID(ctx, tenantID, id)
}

// ListDocumentTypes lists all active document types.
func (s *Service) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]DocumentType, error) {
	return s.repo.ListDocumentTypes(ctx, tenantID, activeOnly)
}

// UpdateDocumentType updates a document type.
func (s *Service) UpdateDocumentType(ctx context.Context, tenantID, id uuid.UUID, req UpdateDocumentTypeRequest) (*DocumentType, error) {
	dt, err := s.repo.GetDocumentTypeByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		dt.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		dt.Description = strings.TrimSpace(*req.Description)
	}
	if req.IsMandatory != nil {
		dt.IsMandatory = *req.IsMandatory
	}
	if req.HasExpiry != nil {
		dt.HasExpiry = *req.HasExpiry
	}
	if req.AllowedExtensions != nil {
		dt.AllowedExtensions = *req.AllowedExtensions
	}
	if req.MaxSizeMB != nil {
		dt.MaxSizeMB = *req.MaxSizeMB
	}
	if req.SortOrder != nil {
		dt.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		dt.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateDocumentType(ctx, dt); err != nil {
		return nil, err
	}

	return dt, nil
}

// EnsureDefaultDocumentTypes ensures default document types exist for a tenant.
func (s *Service) EnsureDefaultDocumentTypes(ctx context.Context, tenantID uuid.UUID) error {
	return s.repo.SeedDefaultDocumentTypes(ctx, tenantID)
}

// =============================================================================
// Student Document Operations
// =============================================================================

// UploadDocument uploads a new document for a student.
func (s *Service) UploadDocument(
	ctx context.Context,
	tenantID, studentID, userID uuid.UUID,
	file io.Reader,
	fileName string,
	fileSize int64,
	mimeType string,
	req UploadDocumentRequest,
) (*StudentDocument, error) {
	// Verify student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, fmt.Errorf("check student: %w", err)
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	// Get document type
	docType, err := s.repo.GetDocumentTypeByID(ctx, tenantID, req.DocumentTypeID)
	if err != nil {
		return nil, err
	}

	// Validate file size
	maxSize := int64(docType.MaxSizeMB) * 1024 * 1024
	if fileSize > maxSize {
		return nil, ErrFileTooLarge
	}

	// Validate file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	allowed := strings.Split(docType.AllowedExtensions, ",")
	isAllowed := false
	for _, a := range allowed {
		if strings.TrimSpace(a) == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, ErrInvalidFileType
	}

	// Parse dates
	var issueDate, expiryDate *time.Time
	if req.IssueDate != "" {
		t, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			if t.After(time.Now()) {
				return nil, ErrInvalidIssueDate
			}
			issueDate = &t
		}
	}
	if req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			if !docType.HasExpiry {
				// Ignore expiry date for documents that don't have expiry
			} else if t.Before(time.Now()) {
				return nil, ErrInvalidExpiryDate
			} else {
				expiryDate = &t
			}
		}
	}

	// Check if document already exists for this type
	existingDoc, err := s.repo.GetDocumentByType(ctx, tenantID, studentID, req.DocumentTypeID)
	if err != nil && err != ErrDocumentNotFound {
		return nil, fmt.Errorf("check existing document: %w", err)
	}

	// Save file
	uploadDir := filepath.Join(s.uploadPath, "documents", tenantID.String(), studentID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	uniqueFileName := fmt.Sprintf("%s_%s.%s", docType.Code, uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, uniqueFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("save file: %w", err)
	}

	// Build file URL
	fileURL := fmt.Sprintf("/uploads/documents/%s/%s/%s", tenantID.String(), studentID.String(), uniqueFileName)

	if existingDoc != nil {
		// Re-upload - update existing document
		existingDoc.FileURL = fileURL
		existingDoc.FileName = fileName
		existingDoc.FileSizeBytes = int(fileSize)
		existingDoc.MimeType = mimeType
		existingDoc.DocumentNumber = req.DocumentNumber
		existingDoc.IssueDate = issueDate
		existingDoc.ExpiryDate = expiryDate
		existingDoc.UploadedBy = userID

		if err := s.repo.UpdateDocumentWithFile(ctx, existingDoc); err != nil {
			os.Remove(filePath)
			return nil, err
		}

		// Delete old file if different
		if existingDoc.FileURL != fileURL {
			oldPath := filepath.Join(s.uploadPath, strings.TrimPrefix(existingDoc.FileURL, "/uploads/"))
			os.Remove(oldPath)
		}

		return s.repo.GetDocumentByID(ctx, tenantID, studentID, existingDoc.ID)
	}

	// Create new document
	doc := &StudentDocument{
		TenantID:       tenantID,
		StudentID:      studentID,
		DocumentTypeID: req.DocumentTypeID,
		FileURL:        fileURL,
		FileName:       fileName,
		FileSizeBytes:  int(fileSize),
		MimeType:       mimeType,
		DocumentNumber: req.DocumentNumber,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		Status:         StatusPendingVerification,
		UploadedBy:     userID,
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	// Fetch with document type
	return s.repo.GetDocumentByID(ctx, tenantID, studentID, doc.ID)
}

// GetDocument gets a document by ID.
func (s *Service) GetDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID) (*StudentDocument, error) {
	doc, err := s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
	if err != nil {
		return nil, err
	}

	// Get uploader name
	if doc.UploadedBy != uuid.Nil {
		name, _ := s.repo.GetUserFullName(ctx, doc.UploadedBy)
		doc.UploaderName = name
	}

	// Get verifier name
	if doc.VerifiedBy != nil {
		name, _ := s.repo.GetUserFullName(ctx, *doc.VerifiedBy)
		doc.VerifierName = name
	}

	return doc, nil
}

// ListDocuments lists all documents for a student.
func (s *Service) ListDocuments(ctx context.Context, tenantID, studentID uuid.UUID, filter DocumentFilter) ([]StudentDocument, error) {
	// Verify student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, fmt.Errorf("check student: %w", err)
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	return s.repo.ListDocumentsByStudent(ctx, tenantID, studentID, filter)
}

// UpdateDocument updates document metadata.
func (s *Service) UpdateDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID, req UpdateDocumentRequest) (*StudentDocument, error) {
	doc, err := s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
	if err != nil {
		return nil, err
	}

	if doc.Version != req.Version {
		return nil, ErrOptimisticLockConflict
	}

	if req.DocumentNumber != nil {
		doc.DocumentNumber = *req.DocumentNumber
	}
	if req.IssueDate != nil {
		if *req.IssueDate == "" {
			doc.IssueDate = nil
		} else {
			t, err := time.Parse("2006-01-02", *req.IssueDate)
			if err == nil {
				if t.After(time.Now()) {
					return nil, ErrInvalidIssueDate
				}
				doc.IssueDate = &t
			}
		}
	}
	if req.ExpiryDate != nil {
		if *req.ExpiryDate == "" {
			doc.ExpiryDate = nil
		} else {
			t, err := time.Parse("2006-01-02", *req.ExpiryDate)
			if err == nil {
				if t.Before(time.Now()) {
					return nil, ErrInvalidExpiryDate
				}
				doc.ExpiryDate = &t
			}
		}
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
}

// VerifyDocument marks a document as verified.
func (s *Service) VerifyDocument(ctx context.Context, tenantID, studentID, docID, verifiedBy uuid.UUID, version int) (*StudentDocument, error) {
	doc, err := s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
	if err != nil {
		return nil, err
	}

	if doc.Status == StatusVerified {
		return nil, ErrDocumentAlreadyVerified
	}

	if err := s.repo.VerifyDocument(ctx, tenantID, studentID, docID, verifiedBy, version); err != nil {
		return nil, err
	}

	return s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
}

// RejectDocument marks a document as rejected.
func (s *Service) RejectDocument(ctx context.Context, tenantID, studentID, docID, rejectedBy uuid.UUID, reason string, version int) (*StudentDocument, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, ErrRejectionReasonRequired
	}

	doc, err := s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
	if err != nil {
		return nil, err
	}

	if doc.Status == StatusVerified {
		return nil, ErrDocumentAlreadyVerified
	}

	if err := s.repo.RejectDocument(ctx, tenantID, studentID, docID, rejectedBy, reason, version); err != nil {
		return nil, err
	}

	return s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
}

// DeleteDocument deletes a document.
func (s *Service) DeleteDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID) error {
	// Get document to find file path
	doc, err := s.repo.GetDocumentByID(ctx, tenantID, studentID, docID)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.DeleteDocument(ctx, tenantID, studentID, docID); err != nil {
		return err
	}

	// Delete file
	filePath := filepath.Join(s.uploadPath, strings.TrimPrefix(doc.FileURL, "/uploads/"))
	if err := os.Remove(filePath); err != nil {
		logger.Warn("Failed to delete document file",
			zap.String("file_path", filePath),
			zap.Error(err))
	}

	return nil
}

// GetDocumentChecklist gets the document checklist for a student.
func (s *Service) GetDocumentChecklist(ctx context.Context, tenantID, studentID uuid.UUID) (*DocumentChecklistResponse, error) {
	// Verify student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, fmt.Errorf("check student: %w", err)
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	// Ensure default document types exist
	if err := s.repo.SeedDefaultDocumentTypes(ctx, tenantID); err != nil {
		logger.Warn("Failed to seed default document types",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
	}

	// Get all active document types
	docTypes, err := s.repo.ListDocumentTypes(ctx, tenantID, true)
	if err != nil {
		return nil, fmt.Errorf("list document types: %w", err)
	}

	// Get all documents for the student
	docs, err := s.repo.ListDocumentsByStudent(ctx, tenantID, studentID, DocumentFilter{})
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}

	// Create a map for quick lookup
	docMap := make(map[uuid.UUID]*StudentDocument)
	for i := range docs {
		docMap[docs[i].DocumentTypeID] = &docs[i]
	}

	// Build checklist
	items := make([]ChecklistItem, len(docTypes))
	totalRequired := 0
	totalUploaded := 0
	totalVerified := 0
	totalPending := 0
	totalRejected := 0

	for i, dt := range docTypes {
		doc := docMap[dt.ID]
		isUploaded := doc != nil
		isVerified := doc != nil && doc.Status == StatusVerified

		items[i] = ChecklistItem{
			DocumentType: dt,
			Document:     doc,
			IsRequired:   dt.IsMandatory,
			IsUploaded:   isUploaded,
			IsVerified:   isVerified,
		}

		if dt.IsMandatory {
			totalRequired++
		}
		if isUploaded {
			totalUploaded++
			switch doc.Status {
			case StatusVerified:
				totalVerified++
			case StatusPendingVerification:
				totalPending++
			case StatusRejected:
				totalRejected++
			}
		}
	}

	// Calculate completion percentage (based on mandatory documents)
	completionPercent := 0.0
	if totalRequired > 0 {
		// Count how many mandatory docs are uploaded
		mandatoryUploaded := 0
		for _, item := range items {
			if item.IsRequired && item.IsUploaded {
				mandatoryUploaded++
			}
		}
		completionPercent = float64(mandatoryUploaded) / float64(totalRequired) * 100
	} else if totalUploaded > 0 {
		completionPercent = 100.0
	}

	return &DocumentChecklistResponse{
		Items:             items,
		TotalRequired:     totalRequired,
		TotalUploaded:     totalUploaded,
		TotalVerified:     totalVerified,
		TotalPending:      totalPending,
		TotalRejected:     totalRejected,
		CompletionPercent: completionPercent,
	}, nil
}
