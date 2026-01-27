// Package designation provides designation management functionality.
package designation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for designations.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new designation repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new designation in the database.
func (r *Repository) Create(ctx context.Context, designation *models.Designation) error {
	if err := r.db.WithContext(ctx).Create(designation).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_designation_name") {
			return ErrDuplicateName
		}
		return fmt.Errorf("create designation: %w", err)
	}
	return nil
}

// GetByID retrieves a designation by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Designation, error) {
	var designation models.Designation
	err := r.db.WithContext(ctx).
		Preload("Department").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&designation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDesignationNotFound
		}
		return nil, fmt.Errorf("get designation by id: %w", err)
	}
	return &designation, nil
}

// GetByName retrieves a designation by name.
func (r *Repository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*models.Designation, error) {
	var designation models.Designation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, name).
		First(&designation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDesignationNotFound
		}
		return nil, fmt.Errorf("get designation by name: %w", err)
	}
	return &designation, nil
}

// Update updates a designation in the database.
func (r *Repository) Update(ctx context.Context, designation *models.Designation) error {
	result := r.db.WithContext(ctx).
		Model(designation).
		Updates(map[string]interface{}{
			"name":          designation.Name,
			"level":         designation.Level,
			"department_id": designation.DepartmentID,
			"is_active":     designation.IsActive,
			"updated_at":    designation.UpdatedAt,
		})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "uniq_designation_name") {
			return ErrDuplicateName
		}
		return fmt.Errorf("update designation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDesignationNotFound
	}

	return nil
}

// Delete deletes a designation.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if designation is in use by any staff
	var staffCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ? AND designation_id = ? AND deleted_at IS NULL", tenantID, id).
		Count(&staffCount).Error; err != nil {
		return fmt.Errorf("check staff usage: %w", err)
	}
	if staffCount > 0 {
		return ErrDesignationInUse
	}

	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Designation{})

	if result.Error != nil {
		return fmt.Errorf("delete designation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDesignationNotFound
	}

	return nil
}

// List retrieves all designations with filters.
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Designation, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Designation{}).
		Preload("Department").
		Where("designations.tenant_id = ?", filter.TenantID)

	if filter.DepartmentID != nil {
		query = query.Where("designations.department_id = ?", *filter.DepartmentID)
	}

	if filter.IsActive != nil {
		query = query.Where("designations.is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(designations.name) LIKE ?", search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count designations: %w", err)
	}

	// Get results ordered by level (hierarchy)
	var designations []models.Designation
	if err := query.Order("designations.level ASC, designations.name ASC").Find(&designations).Error; err != nil {
		return nil, 0, fmt.Errorf("list designations: %w", err)
	}

	return designations, total, nil
}

// GetStaffCounts returns staff counts for designations.
func (r *Repository) GetStaffCounts(ctx context.Context, tenantID uuid.UUID, designationIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(designationIDs) == 0 {
		return make(map[uuid.UUID]int), nil
	}

	type countResult struct {
		DesignationID uuid.UUID `gorm:"column:designation_id"`
		Count         int       `gorm:"column:count"`
	}

	var results []countResult
	err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Select("designation_id, COUNT(*) as count").
		Where("tenant_id = ? AND designation_id IN ? AND deleted_at IS NULL", tenantID, designationIDs).
		Group("designation_id").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get staff counts: %w", err)
	}

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.DesignationID] = r.Count
	}

	return counts, nil
}

// GetActiveForDropdown retrieves active designations for dropdown/select.
func (r *Repository) GetActiveForDropdown(ctx context.Context, tenantID uuid.UUID, departmentID *uuid.UUID) ([]models.Designation, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID)

	if departmentID != nil {
		query = query.Where("department_id = ? OR department_id IS NULL", *departmentID)
	}

	var designations []models.Designation
	err := query.Order("level ASC, name ASC").Find(&designations).Error
	if err != nil {
		return nil, fmt.Errorf("get active designations: %w", err)
	}
	return designations, nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
