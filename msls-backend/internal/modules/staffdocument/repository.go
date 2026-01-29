// Package staffdocument provides staff document management functionality.
package staffdocument

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for staff documents.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new staff document repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// Document Type Methods
// ============================================================================

// CreateDocumentType creates a new document type.
func (r *Repository) CreateDocumentType(ctx context.Context, dt *models.StaffDocumentType) error {
	if err := r.db.WithContext(ctx).Create(dt).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_document_type_code") {
			return ErrDuplicateDocTypeCode
		}
		return fmt.Errorf("create document type: %w", err)
	}
	return nil
}

// GetDocumentTypeByID retrieves a document type by ID.
func (r *Repository) GetDocumentTypeByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffDocumentType, error) {
	var dt models.StaffDocumentType
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&dt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentTypeNotFound
		}
		return nil, fmt.Errorf("get document type: %w", err)
	}
	return &dt, nil
}

// GetDocumentTypeByCode retrieves a document type by code.
func (r *Repository) GetDocumentTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.StaffDocumentType, error) {
	var dt models.StaffDocumentType
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&dt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentTypeNotFound
		}
		return nil, fmt.Errorf("get document type by code: %w", err)
	}
	return &dt, nil
}

// UpdateDocumentType updates a document type.
func (r *Repository) UpdateDocumentType(ctx context.Context, dt *models.StaffDocumentType) error {
	if err := r.db.WithContext(ctx).Save(dt).Error; err != nil {
		return fmt.Errorf("update document type: %w", err)
	}
	return nil
}

// DeleteDocumentType deletes a document type.
func (r *Repository) DeleteDocumentType(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.StaffDocumentType{})
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "violates foreign key") {
			return ErrDocumentTypeInUse
		}
		return fmt.Errorf("delete document type: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDocumentTypeNotFound
	}
	return nil
}

// ListDocumentTypes retrieves all document types for a tenant.
func (r *Repository) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]models.StaffDocumentType, error) {
	var types []models.StaffDocumentType
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if activeOnly {
		query = query.Where("is_active = true")
	}
	if err := query.Order("display_order ASC, name ASC").Find(&types).Error; err != nil {
		return nil, fmt.Errorf("list document types: %w", err)
	}
	return types, nil
}

// GetMandatoryDocumentTypes retrieves mandatory document types for a staff type.
func (r *Repository) GetMandatoryDocumentTypes(ctx context.Context, tenantID uuid.UUID, staffType string) ([]models.StaffDocumentType, error) {
	var types []models.StaffDocumentType
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_mandatory = true AND is_active = true", tenantID).
		Where("? = ANY(applicable_to)", staffType).
		Order("display_order ASC").
		Find(&types).Error
	if err != nil {
		return nil, fmt.Errorf("get mandatory document types: %w", err)
	}
	return types, nil
}

// ============================================================================
// Document Methods
// ============================================================================

// CreateDocument creates a new staff document.
func (r *Repository) CreateDocument(ctx context.Context, doc *models.StaffDocument) error {
	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	return nil
}

// GetDocumentByID retrieves a document by ID.
func (r *Repository) GetDocumentByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffDocument, error) {
	var doc models.StaffDocument
	err := r.db.WithContext(ctx).
		Preload("DocumentType").
		Preload("Verifier").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &doc, nil
}

// GetDocumentByStaffAndID retrieves a document for a specific staff member.
func (r *Repository) GetDocumentByStaffAndID(ctx context.Context, tenantID, staffID, docID uuid.UUID) (*models.StaffDocument, error) {
	var doc models.StaffDocument
	err := r.db.WithContext(ctx).
		Preload("DocumentType").
		Preload("Verifier").
		Where("tenant_id = ? AND staff_id = ? AND id = ?", tenantID, staffID, docID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &doc, nil
}

// UpdateDocument updates a document.
func (r *Repository) UpdateDocument(ctx context.Context, doc *models.StaffDocument) error {
	doc.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(doc).Error; err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	return nil
}

// DeleteDocument deletes a document.
func (r *Repository) DeleteDocument(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.StaffDocument{})
	if result.Error != nil {
		return fmt.Errorf("delete document: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDocumentNotFound
	}
	return nil
}

// ListStaffDocuments retrieves documents for a staff member.
func (r *Repository) ListStaffDocuments(ctx context.Context, tenantID, staffID uuid.UUID, params ListDocumentsParams) ([]models.StaffDocument, int64, error) {
	var docs []models.StaffDocument
	var total int64

	query := r.db.WithContext(ctx).
		Preload("DocumentType").
		Where("tenant_id = ? AND staff_id = ?", tenantID, staffID)

	if params.DocumentTypeID != nil {
		query = query.Where("document_type_id = ?", *params.DocumentTypeID)
	}
	if params.Status != "" {
		query = query.Where("verification_status = ?", params.Status)
	}
	if params.IsCurrent != nil {
		query = query.Where("is_current = ?", *params.IsCurrent)
	}

	// Count total
	if err := query.Model(&models.StaffDocument{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count documents: %w", err)
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Cursor != "" {
		cursorID, _ := uuid.Parse(params.Cursor)
		query = query.Where("id > ?", cursorID)
	}

	if err := query.Order("created_at DESC").Find(&docs).Error; err != nil {
		return nil, 0, fmt.Errorf("list documents: %w", err)
	}

	return docs, total, nil
}

// ListDocumentsParams contains parameters for listing documents.
type ListDocumentsParams struct {
	DocumentTypeID *uuid.UUID
	Status         string
	IsCurrent      *bool
	Cursor         string
	Limit          int
}

// GetExistingDocument checks if a document of this type exists for staff.
func (r *Repository) GetExistingDocument(ctx context.Context, tenantID, staffID, documentTypeID uuid.UUID) (*models.StaffDocument, error) {
	var doc models.StaffDocument
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND staff_id = ? AND document_type_id = ? AND is_current = true", tenantID, staffID, documentTypeID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get existing document: %w", err)
	}
	return &doc, nil
}

// MarkPreviousDocumentsAsNotCurrent marks old documents as not current.
func (r *Repository) MarkPreviousDocumentsAsNotCurrent(ctx context.Context, tenantID, staffID, documentTypeID, excludeID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&models.StaffDocument{}).
		Where("tenant_id = ? AND staff_id = ? AND document_type_id = ? AND id != ?", tenantID, staffID, documentTypeID, excludeID).
		Update("is_current", false).Error
	if err != nil {
		return fmt.Errorf("mark previous documents: %w", err)
	}
	return nil
}

// ============================================================================
// Expiry & Compliance Methods
// ============================================================================

// GetExpiringDocuments retrieves documents expiring within the given days.
func (r *Repository) GetExpiringDocuments(ctx context.Context, tenantID uuid.UUID, days int) ([]ExpiringDocumentResult, error) {
	var results []ExpiringDocumentResult
	threshold := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Table("staff_documents d").
		Select(`d.*, s.first_name || ' ' || COALESCE(s.last_name, '') as staff_name, s.employee_id,
			EXTRACT(DAY FROM (d.expiry_date - CURRENT_DATE)) as days_to_expiry`).
		Joins("JOIN staff s ON s.id = d.staff_id").
		Where("d.tenant_id = ?", tenantID).
		Where("d.expiry_date IS NOT NULL").
		Where("d.expiry_date <= ?", threshold).
		Where("d.is_current = true").
		Where("d.verification_status != 'rejected'").
		Order("d.expiry_date ASC").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get expiring documents: %w", err)
	}
	return results, nil
}

// ExpiringDocumentResult contains document with staff info for expiry queries.
type ExpiringDocumentResult struct {
	models.StaffDocument
	StaffName    string `gorm:"column:staff_name"`
	EmployeeID   string `gorm:"column:employee_id"`
	DaysToExpiry int    `gorm:"column:days_to_expiry"`
}

// GetComplianceStats retrieves compliance statistics.
func (r *Repository) GetComplianceStats(ctx context.Context, tenantID uuid.UUID) (*ComplianceStats, error) {
	stats := &ComplianceStats{}

	// Total staff count
	var staffCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ? AND status = 'active'", tenantID).
		Count(&staffCount).Error; err != nil {
		return nil, fmt.Errorf("count staff: %w", err)
	}
	stats.TotalStaff = int(staffCount)

	// Document counts
	var counts struct {
		Total    int
		Pending  int
		Verified int
		Rejected int
		Expired  int
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE verification_status = 'pending') as pending,
			COUNT(*) FILTER (WHERE verification_status = 'verified') as verified,
			COUNT(*) FILTER (WHERE verification_status = 'rejected') as rejected,
			COUNT(*) FILTER (WHERE expiry_date < CURRENT_DATE AND is_current = true) as expired
		FROM staff_documents
		WHERE tenant_id = ? AND is_current = true
	`, tenantID).Scan(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("get document counts: %w", err)
	}

	stats.DocumentsSubmitted = counts.Total
	stats.PendingVerification = counts.Pending
	stats.Verified = counts.Verified
	stats.Rejected = counts.Rejected
	stats.Expired = counts.Expired

	// Expiring counts
	now := time.Now()
	for days, field := range map[int]*int{30: &stats.ExpiringIn30Days, 60: &stats.ExpiringIn60Days, 90: &stats.ExpiringIn90Days} {
		threshold := now.AddDate(0, 0, days)
		var count int64
		if err := r.db.WithContext(ctx).
			Model(&models.StaffDocument{}).
			Where("tenant_id = ? AND is_current = true", tenantID).
			Where("expiry_date IS NOT NULL AND expiry_date <= ? AND expiry_date > ?", threshold, now).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("count expiring documents: %w", err)
		}
		*field = int(count)
	}

	// Calculate compliance percentage (verified / total submitted * 100)
	if stats.DocumentsSubmitted > 0 {
		stats.CompliancePercentage = float64(stats.Verified) / float64(stats.DocumentsSubmitted) * 100
	}

	return stats, nil
}

// GetComplianceByDocumentType retrieves compliance per document type.
func (r *Repository) GetComplianceByDocumentType(ctx context.Context, tenantID uuid.UUID) ([]DocumentTypeComplianceResult, error) {
	var results []DocumentTypeComplianceResult

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			dt.id as document_type_id,
			dt.name as document_type_name,
			dt.code as document_type_code,
			dt.category,
			dt.is_mandatory,
			(SELECT COUNT(*) FROM staff WHERE tenant_id = dt.tenant_id AND status = 'active') as required,
			COUNT(d.id) as submitted,
			COUNT(d.id) FILTER (WHERE d.verification_status = 'verified') as verified,
			COUNT(d.id) FILTER (WHERE d.verification_status = 'pending') as pending,
			COUNT(d.id) FILTER (WHERE d.verification_status = 'rejected') as rejected,
			COUNT(d.id) FILTER (WHERE d.expiry_date < CURRENT_DATE) as expired
		FROM staff_document_types dt
		LEFT JOIN staff_documents d ON d.document_type_id = dt.id AND d.is_current = true
		WHERE dt.tenant_id = ? AND dt.is_active = true
		GROUP BY dt.id, dt.name, dt.code, dt.category, dt.is_mandatory, dt.display_order
		ORDER BY dt.display_order ASC
	`, tenantID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get compliance by document type: %w", err)
	}
	return results, nil
}

// DocumentTypeComplianceResult contains compliance data per document type.
type DocumentTypeComplianceResult struct {
	DocumentTypeID   uuid.UUID `gorm:"column:document_type_id"`
	DocumentTypeName string    `gorm:"column:document_type_name"`
	DocumentTypeCode string    `gorm:"column:document_type_code"`
	Category         string    `gorm:"column:category"`
	IsMandatory      bool      `gorm:"column:is_mandatory"`
	Required         int       `gorm:"column:required"`
	Submitted        int       `gorm:"column:submitted"`
	Verified         int       `gorm:"column:verified"`
	Pending          int       `gorm:"column:pending"`
	Rejected         int       `gorm:"column:rejected"`
	Expired          int       `gorm:"column:expired"`
}

// GetStaffComplianceDetails retrieves compliance details per staff member.
func (r *Repository) GetStaffComplianceDetails(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]StaffComplianceRow, error) {
	var results []StaffComplianceRow

	query := r.db.WithContext(ctx).Raw(`
		WITH mandatory_types AS (
			SELECT id, name FROM staff_document_types
			WHERE tenant_id = ? AND is_mandatory = true AND is_active = true
		),
		staff_doc_counts AS (
			SELECT
				s.id as staff_id,
				s.first_name || ' ' || COALESCE(s.last_name, '') as staff_name,
				s.employee_id,
				(SELECT COUNT(*) FROM mandatory_types) as total_required,
				COUNT(d.id) as submitted,
				COUNT(d.id) FILTER (WHERE d.verification_status = 'verified') as verified,
				COUNT(d.id) FILTER (WHERE d.verification_status = 'pending') as pending,
				COUNT(d.id) FILTER (WHERE d.verification_status = 'rejected') as rejected,
				COUNT(d.id) FILTER (WHERE d.expiry_date < CURRENT_DATE) as expired,
				ARRAY_AGG(dt.name) FILTER (WHERE d.id IS NULL) as missing_documents
			FROM staff s
			CROSS JOIN mandatory_types mt
			LEFT JOIN staff_documents d ON d.staff_id = s.id AND d.document_type_id = mt.id AND d.is_current = true
			LEFT JOIN staff_document_types dt ON dt.id = mt.id
			WHERE s.tenant_id = ? AND s.status = 'active'
			GROUP BY s.id, s.first_name, s.last_name, s.employee_id
		)
		SELECT * FROM staff_doc_counts
		ORDER BY staff_name
		LIMIT ? OFFSET ?
	`, tenantID, tenantID, limit, offset).Scan(&results)

	if query.Error != nil {
		return nil, fmt.Errorf("get staff compliance details: %w", query.Error)
	}
	return results, nil
}

// StaffComplianceRow contains compliance data for a staff member.
type StaffComplianceRow struct {
	StaffID          uuid.UUID `gorm:"column:staff_id"`
	StaffName        string    `gorm:"column:staff_name"`
	EmployeeID       string    `gorm:"column:employee_id"`
	TotalRequired    int       `gorm:"column:total_required"`
	Submitted        int       `gorm:"column:submitted"`
	Verified         int       `gorm:"column:verified"`
	Pending          int       `gorm:"column:pending"`
	Rejected         int       `gorm:"column:rejected"`
	Expired          int       `gorm:"column:expired"`
	MissingDocuments []string  `gorm:"column:missing_documents;type:text[]"`
}

// ============================================================================
// Notification Methods
// ============================================================================

// CreateNotification creates a notification record.
func (r *Repository) CreateNotification(ctx context.Context, notification *models.StaffDocumentNotification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

// GetLastNotification gets the last notification of a type for a document.
func (r *Repository) GetLastNotification(ctx context.Context, documentID uuid.UUID, notificationType models.NotificationType) (*models.StaffDocumentNotification, error) {
	var notification models.StaffDocumentNotification
	err := r.db.WithContext(ctx).
		Where("document_id = ? AND notification_type = ?", documentID, notificationType).
		Order("sent_at DESC").
		First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last notification: %w", err)
	}
	return &notification, nil
}
