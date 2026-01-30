// Package exam provides examination management functionality.
package exam

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for exam entities.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new exam repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Exam Type Repository Methods
// ========================================

// ListExamTypes returns all exam types for a tenant with filters.
func (r *Repository) ListExamTypes(ctx context.Context, filter ExamTypeFilter) ([]models.ExamType, int64, error) {
	var examTypes []models.ExamType
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ExamType{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("display_order ASC, name ASC").
		Find(&examTypes).Error; err != nil {
		return nil, 0, err
	}

	return examTypes, total, nil
}

// GetExamTypeByID returns an exam type by ID.
func (r *Repository) GetExamTypeByID(ctx context.Context, tenantID, id uuid.UUID) (*models.ExamType, error) {
	var examType models.ExamType
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&examType).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrExamTypeNotFound
	}
	return &examType, err
}

// GetExamTypeByCode returns an exam type by code.
func (r *Repository) GetExamTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.ExamType, error) {
	var examType models.ExamType
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&examType).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Not found is not an error for checking existence
	}
	return &examType, err
}

// CreateExamType creates a new exam type.
func (r *Repository) CreateExamType(ctx context.Context, examType *models.ExamType) error {
	return r.db.WithContext(ctx).Create(examType).Error
}

// UpdateExamType updates an existing exam type.
func (r *Repository) UpdateExamType(ctx context.Context, examType *models.ExamType) error {
	return r.db.WithContext(ctx).Save(examType).Error
}

// DeleteExamType deletes an exam type by ID.
func (r *Repository) DeleteExamType(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.ExamType{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExamTypeNotFound
	}
	return nil
}

// GetMaxDisplayOrder returns the maximum display order for exam types.
func (r *Repository) GetMaxDisplayOrder(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).
		Model(&models.ExamType{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

// UpdateDisplayOrders updates display orders for multiple exam types.
func (r *Repository) UpdateDisplayOrders(ctx context.Context, tenantID uuid.UUID, items []DisplayOrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.ExamType{}).
				Where("tenant_id = ? AND id = ?", tenantID, item.ID).
				Update("display_order", item.DisplayOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ToggleActive toggles the active status of an exam type.
func (r *Repository) ToggleActive(ctx context.Context, tenantID, id uuid.UUID, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.ExamType{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Update("is_active", isActive)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExamTypeNotFound
	}
	return nil
}

// CountExamTypeUsage counts how many exams use this exam type.
// This is a placeholder - will be implemented when exams table exists.
func (r *Repository) CountExamTypeUsage(ctx context.Context, tenantID, examTypeID uuid.UUID) (int64, error) {
	// TODO: Implement when exams table is created in Story 8.2
	return 0, nil
}
