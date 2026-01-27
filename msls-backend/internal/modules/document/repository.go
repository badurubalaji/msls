package document

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles document database operations.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new document repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DocumentTypeModel represents the document_types table.
type DocumentTypeModel struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID          uuid.UUID      `gorm:"type:uuid;not null"`
	Code              string         `gorm:"type:varchar(50);not null"`
	Name              string         `gorm:"type:varchar(100);not null"`
	Description       sql.NullString `gorm:"type:text"`
	IsMandatory       bool           `gorm:"not null;default:false"`
	HasExpiry         bool           `gorm:"not null;default:false"`
	AllowedExtensions string         `gorm:"type:varchar(100);not null;default:'pdf,jpg,jpeg,png'"`
	MaxSizeMB         int            `gorm:"not null;default:5"`
	SortOrder         int            `gorm:"not null;default:0"`
	IsActive          bool           `gorm:"not null;default:true"`
	CreatedAt         time.Time      `gorm:"not null;default:now()"`
	UpdatedAt         time.Time      `gorm:"not null;default:now()"`
}

// TableName returns the table name for DocumentTypeModel.
func (DocumentTypeModel) TableName() string {
	return "document_types"
}

// StudentDocumentModel represents the student_documents table.
type StudentDocumentModel struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID        uuid.UUID      `gorm:"type:uuid;not null"`
	StudentID       uuid.UUID      `gorm:"type:uuid;not null"`
	DocumentTypeID  uuid.UUID      `gorm:"type:uuid;not null"`
	FileURL         string         `gorm:"type:varchar(500);not null"`
	FileName        string         `gorm:"type:varchar(255);not null"`
	FileSizeBytes   int            `gorm:"not null"`
	MimeType        string         `gorm:"type:varchar(100);not null"`
	DocumentNumber  sql.NullString `gorm:"type:varchar(100)"`
	IssueDate       sql.NullTime   `gorm:"type:date"`
	ExpiryDate      sql.NullTime   `gorm:"type:date"`
	Status          string         `gorm:"type:document_status;not null;default:'pending_verification'"`
	RejectionReason sql.NullString `gorm:"type:text"`
	VerifiedAt      sql.NullTime   `gorm:"type:timestamptz"`
	VerifiedBy      *uuid.UUID     `gorm:"type:uuid"`
	UploadedAt      time.Time      `gorm:"not null;default:now()"`
	UploadedBy      uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedAt       time.Time      `gorm:"not null;default:now()"`
	Version         int            `gorm:"not null;default:1"`

	// Associations
	DocumentType DocumentTypeModel `gorm:"foreignKey:DocumentTypeID"`
}

// TableName returns the table name for StudentDocumentModel.
func (StudentDocumentModel) TableName() string {
	return "student_documents"
}

// =============================================================================
// Document Type Operations
// =============================================================================

// CreateDocumentType creates a new document type.
func (r *Repository) CreateDocumentType(ctx context.Context, dt *DocumentType) error {
	model := &DocumentTypeModel{
		TenantID:          dt.TenantID,
		Code:              dt.Code,
		Name:              dt.Name,
		Description:       toNullString(dt.Description),
		IsMandatory:       dt.IsMandatory,
		HasExpiry:         dt.HasExpiry,
		AllowedExtensions: dt.AllowedExtensions,
		MaxSizeMB:         dt.MaxSizeMB,
		SortOrder:         dt.SortOrder,
		IsActive:          dt.IsActive,
	}

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}

	dt.ID = model.ID
	dt.CreatedAt = model.CreatedAt
	dt.UpdatedAt = model.UpdatedAt
	return nil
}

// GetDocumentTypeByID gets a document type by ID.
func (r *Repository) GetDocumentTypeByID(ctx context.Context, tenantID, id uuid.UUID) (*DocumentType, error) {
	var model DocumentTypeModel
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentTypeNotFound
		}
		return nil, result.Error
	}

	return toDocumentType(&model), nil
}

// GetDocumentTypeByCode gets a document type by code.
func (r *Repository) GetDocumentTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*DocumentType, error) {
	var model DocumentTypeModel
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentTypeNotFound
		}
		return nil, result.Error
	}

	return toDocumentType(&model), nil
}

// ListDocumentTypes lists all active document types for a tenant.
func (r *Repository) ListDocumentTypes(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]DocumentType, error) {
	var models []DocumentTypeModel

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	result := query.Order("sort_order ASC, name ASC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	types := make([]DocumentType, len(models))
	for i, m := range models {
		types[i] = *toDocumentType(&m)
	}

	return types, nil
}

// UpdateDocumentType updates a document type.
func (r *Repository) UpdateDocumentType(ctx context.Context, dt *DocumentType) error {
	updates := map[string]interface{}{
		"name":               dt.Name,
		"description":        toNullString(dt.Description),
		"is_mandatory":       dt.IsMandatory,
		"has_expiry":         dt.HasExpiry,
		"allowed_extensions": dt.AllowedExtensions,
		"max_size_mb":        dt.MaxSizeMB,
		"sort_order":         dt.SortOrder,
		"is_active":          dt.IsActive,
		"updated_at":         time.Now(),
	}

	result := r.db.WithContext(ctx).
		Model(&DocumentTypeModel{}).
		Where("tenant_id = ? AND id = ?", dt.TenantID, dt.ID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDocumentTypeNotFound
	}

	return nil
}

// =============================================================================
// Student Document Operations
// =============================================================================

// CreateDocument creates a new student document.
func (r *Repository) CreateDocument(ctx context.Context, doc *StudentDocument) error {
	model := &StudentDocumentModel{
		TenantID:        doc.TenantID,
		StudentID:       doc.StudentID,
		DocumentTypeID:  doc.DocumentTypeID,
		FileURL:         doc.FileURL,
		FileName:        doc.FileName,
		FileSizeBytes:   doc.FileSizeBytes,
		MimeType:        doc.MimeType,
		DocumentNumber:  toNullString(doc.DocumentNumber),
		IssueDate:       toNullTime(doc.IssueDate),
		ExpiryDate:      toNullTime(doc.ExpiryDate),
		Status:          string(doc.Status),
		UploadedBy:      doc.UploadedBy,
	}

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}

	doc.ID = model.ID
	doc.UploadedAt = model.UploadedAt
	doc.UpdatedAt = model.UpdatedAt
	doc.Version = model.Version
	return nil
}

// GetDocumentByID gets a document by ID with document type.
func (r *Repository) GetDocumentByID(ctx context.Context, tenantID, studentID, docID uuid.UUID) (*StudentDocument, error) {
	var model StudentDocumentModel
	result := r.db.WithContext(ctx).
		Preload("DocumentType").
		Where("tenant_id = ? AND student_id = ? AND id = ?", tenantID, studentID, docID).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, result.Error
	}

	return toStudentDocument(&model), nil
}

// GetDocumentByType gets a document by student and type.
func (r *Repository) GetDocumentByType(ctx context.Context, tenantID, studentID, docTypeID uuid.UUID) (*StudentDocument, error) {
	var model StudentDocumentModel
	result := r.db.WithContext(ctx).
		Preload("DocumentType").
		Where("tenant_id = ? AND student_id = ? AND document_type_id = ?", tenantID, studentID, docTypeID).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, result.Error
	}

	return toStudentDocument(&model), nil
}

// ListDocumentsByStudent lists all documents for a student.
func (r *Repository) ListDocumentsByStudent(ctx context.Context, tenantID, studentID uuid.UUID, filter DocumentFilter) ([]StudentDocument, error) {
	var models []StudentDocumentModel

	query := r.db.WithContext(ctx).
		Preload("DocumentType").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)

	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.DocumentTypeID != nil {
		query = query.Where("document_type_id = ?", *filter.DocumentTypeID)
	}

	result := query.Order("uploaded_at DESC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	docs := make([]StudentDocument, len(models))
	for i, m := range models {
		docs[i] = *toStudentDocument(&m)
	}

	return docs, nil
}

// UpdateDocument updates a document.
func (r *Repository) UpdateDocument(ctx context.Context, doc *StudentDocument) error {
	updates := map[string]interface{}{
		"document_number": toNullString(doc.DocumentNumber),
		"issue_date":      toNullTime(doc.IssueDate),
		"expiry_date":     toNullTime(doc.ExpiryDate),
	}

	result := r.db.WithContext(ctx).
		Model(&StudentDocumentModel{}).
		Where("tenant_id = ? AND student_id = ? AND id = ? AND version = ?",
			doc.TenantID, doc.StudentID, doc.ID, doc.Version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	doc.Version++
	return nil
}

// UpdateDocumentWithFile updates a document including the file info.
func (r *Repository) UpdateDocumentWithFile(ctx context.Context, doc *StudentDocument) error {
	updates := map[string]interface{}{
		"file_url":        doc.FileURL,
		"file_name":       doc.FileName,
		"file_size_bytes": doc.FileSizeBytes,
		"mime_type":       doc.MimeType,
		"document_number": toNullString(doc.DocumentNumber),
		"issue_date":      toNullTime(doc.IssueDate),
		"expiry_date":     toNullTime(doc.ExpiryDate),
		"status":          string(StatusPendingVerification),
		"rejection_reason": sql.NullString{Valid: false},
		"verified_at":      sql.NullTime{Valid: false},
		"verified_by":      nil,
		"uploaded_at":     time.Now(),
		"uploaded_by":     doc.UploadedBy,
	}

	result := r.db.WithContext(ctx).
		Model(&StudentDocumentModel{}).
		Where("tenant_id = ? AND student_id = ? AND id = ? AND version = ?",
			doc.TenantID, doc.StudentID, doc.ID, doc.Version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	doc.Status = StatusPendingVerification
	doc.RejectionReason = ""
	doc.VerifiedAt = nil
	doc.VerifiedBy = nil
	doc.Version++
	return nil
}

// VerifyDocument marks a document as verified.
func (r *Repository) VerifyDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID, verifiedBy uuid.UUID, version int) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":           string(StatusVerified),
		"verified_at":      now,
		"verified_by":      verifiedBy,
		"rejection_reason": sql.NullString{Valid: false},
	}

	result := r.db.WithContext(ctx).
		Model(&StudentDocumentModel{}).
		Where("tenant_id = ? AND student_id = ? AND id = ? AND version = ?",
			tenantID, studentID, docID, version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// RejectDocument marks a document as rejected.
func (r *Repository) RejectDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID, rejectedBy uuid.UUID, reason string, version int) error {
	updates := map[string]interface{}{
		"status":           string(StatusRejected),
		"verified_at":      sql.NullTime{Valid: false},
		"verified_by":      rejectedBy,
		"rejection_reason": sql.NullString{String: reason, Valid: true},
	}

	result := r.db.WithContext(ctx).
		Model(&StudentDocumentModel{}).
		Where("tenant_id = ? AND student_id = ? AND id = ? AND version = ?",
			tenantID, studentID, docID, version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// DeleteDocument deletes a document.
func (r *Repository) DeleteDocument(ctx context.Context, tenantID, studentID, docID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ? AND id = ?", tenantID, studentID, docID).
		Delete(&StudentDocumentModel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDocumentNotFound
	}

	return nil
}

// StudentExists checks if a student exists.
func (r *Repository) StudentExists(ctx context.Context, tenantID, studentID uuid.UUID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Table("students").
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, studentID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// GetUserFullName gets a user's full name.
func (r *Repository) GetUserFullName(ctx context.Context, userID uuid.UUID) (string, error) {
	var name string
	result := r.db.WithContext(ctx).
		Table("users").
		Select("COALESCE(first_name || ' ' || COALESCE(last_name, ''), email) as name").
		Where("id = ?", userID).
		Scan(&name)

	if result.Error != nil {
		return "", result.Error
	}

	return name, nil
}

// =============================================================================
// Helper Functions
// =============================================================================

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func toDocumentType(m *DocumentTypeModel) *DocumentType {
	return &DocumentType{
		ID:                m.ID,
		TenantID:          m.TenantID,
		Code:              m.Code,
		Name:              m.Name,
		Description:       m.Description.String,
		IsMandatory:       m.IsMandatory,
		HasExpiry:         m.HasExpiry,
		AllowedExtensions: m.AllowedExtensions,
		MaxSizeMB:         m.MaxSizeMB,
		SortOrder:         m.SortOrder,
		IsActive:          m.IsActive,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func toStudentDocument(m *StudentDocumentModel) *StudentDocument {
	doc := &StudentDocument{
		ID:             m.ID,
		TenantID:       m.TenantID,
		StudentID:      m.StudentID,
		DocumentTypeID: m.DocumentTypeID,
		FileURL:        m.FileURL,
		FileName:       m.FileName,
		FileSizeBytes:  m.FileSizeBytes,
		MimeType:       m.MimeType,
		Status:         DocumentStatus(m.Status),
		UploadedAt:     m.UploadedAt,
		UploadedBy:     m.UploadedBy,
		UpdatedAt:      m.UpdatedAt,
		Version:        m.Version,
	}

	if m.DocumentNumber.Valid {
		doc.DocumentNumber = m.DocumentNumber.String
	}
	if m.IssueDate.Valid {
		doc.IssueDate = &m.IssueDate.Time
	}
	if m.ExpiryDate.Valid {
		doc.ExpiryDate = &m.ExpiryDate.Time
	}
	if m.RejectionReason.Valid {
		doc.RejectionReason = m.RejectionReason.String
	}
	if m.VerifiedAt.Valid {
		doc.VerifiedAt = &m.VerifiedAt.Time
	}
	if m.VerifiedBy != nil {
		doc.VerifiedBy = m.VerifiedBy
	}

	// Convert embedded DocumentType if loaded
	if m.DocumentType.ID != uuid.Nil {
		doc.DocumentType = toDocumentType(&m.DocumentType)
	}

	return doc
}

// SeedDefaultDocumentTypes creates default document types for a tenant.
func (r *Repository) SeedDefaultDocumentTypes(ctx context.Context, tenantID uuid.UUID) error {
	defaults := []DocumentTypeModel{
		{TenantID: tenantID, Code: "birth_certificate", Name: "Birth Certificate", IsMandatory: true, SortOrder: 1, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
		{TenantID: tenantID, Code: "aadhaar", Name: "Aadhaar Card", IsMandatory: true, SortOrder: 2, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
		{TenantID: tenantID, Code: "transfer_certificate", Name: "Transfer Certificate", IsMandatory: false, SortOrder: 3, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
		{TenantID: tenantID, Code: "caste_certificate", Name: "Caste Certificate", IsMandatory: false, SortOrder: 4, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
		{TenantID: tenantID, Code: "income_certificate", Name: "Income Certificate", IsMandatory: false, SortOrder: 5, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
		{TenantID: tenantID, Code: "photo_id", Name: "Photo ID", IsMandatory: false, SortOrder: 6, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 2, IsActive: true},
		{TenantID: tenantID, Code: "medical_certificate", Name: "Medical Certificate", HasExpiry: true, IsMandatory: false, SortOrder: 7, AllowedExtensions: "pdf,jpg,jpeg,png", MaxSizeMB: 5, IsActive: true},
	}

	for _, dt := range defaults {
		// Check if already exists
		var count int64
		r.db.WithContext(ctx).
			Model(&DocumentTypeModel{}).
			Where("tenant_id = ? AND code = ?", tenantID, dt.Code).
			Count(&count)

		if count == 0 {
			if err := r.db.WithContext(ctx).Create(&dt).Error; err != nil {
				return fmt.Errorf("seed document type %s: %w", dt.Code, err)
			}
		}
	}

	return nil
}
