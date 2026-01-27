// Package department provides department management functionality.
package department

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for departments.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new department repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new department in the database.
func (r *Repository) Create(ctx context.Context, department *models.Department) error {
	if err := r.db.WithContext(ctx).Create(department).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_department_code") {
			return ErrDuplicateCode
		}
		return fmt.Errorf("create department: %w", err)
	}
	return nil
}

// GetByID retrieves a department by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Department, error) {
	var department models.Department
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Head").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepartmentNotFound
		}
		return nil, fmt.Errorf("get department by id: %w", err)
	}
	return &department, nil
}

// GetByCode retrieves a department by code within a branch.
func (r *Repository) GetByCode(ctx context.Context, tenantID, branchID uuid.UUID, code string) (*models.Department, error) {
	var department models.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND code = ?", tenantID, branchID, code).
		First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepartmentNotFound
		}
		return nil, fmt.Errorf("get department by code: %w", err)
	}
	return &department, nil
}

// Update updates a department in the database.
func (r *Repository) Update(ctx context.Context, department *models.Department) error {
	result := r.db.WithContext(ctx).
		Model(department).
		Updates(map[string]interface{}{
			"name":        department.Name,
			"code":        department.Code,
			"description": department.Description,
			"head_id":     department.HeadID,
			"is_active":   department.IsActive,
			"updated_at":  department.UpdatedAt,
		})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "uniq_department_code") {
			return ErrDuplicateCode
		}
		return fmt.Errorf("update department: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDepartmentNotFound
	}

	return nil
}

// Delete deletes a department.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if department is in use by any designations
	var designationCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Designation{}).
		Where("tenant_id = ? AND department_id = ?", tenantID, id).
		Count(&designationCount).Error; err != nil {
		return fmt.Errorf("check designation usage: %w", err)
	}
	if designationCount > 0 {
		return ErrDepartmentInUse
	}

	// Check if department is in use by any staff
	var staffCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ? AND department_id = ? AND deleted_at IS NULL", tenantID, id).
		Count(&staffCount).Error; err != nil {
		return fmt.Errorf("check staff usage: %w", err)
	}
	if staffCount > 0 {
		return ErrDepartmentInUse
	}

	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Department{})

	if result.Error != nil {
		return fmt.Errorf("delete department: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDepartmentNotFound
	}

	return nil
}

// List retrieves all departments with filters.
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Department, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Preload("Branch").
		Preload("Head").
		Where("departments.tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("departments.branch_id = ?", *filter.BranchID)
	}

	if filter.IsActive != nil {
		query = query.Where("departments.is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(departments.name) LIKE ? OR LOWER(departments.code) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count departments: %w", err)
	}

	// Get results
	var departments []models.Department
	if err := query.Order("departments.name ASC").Find(&departments).Error; err != nil {
		return nil, 0, fmt.Errorf("list departments: %w", err)
	}

	return departments, total, nil
}

// GetStaffCounts returns staff counts for departments.
func (r *Repository) GetStaffCounts(ctx context.Context, tenantID uuid.UUID, departmentIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(departmentIDs) == 0 {
		return make(map[uuid.UUID]int), nil
	}

	type countResult struct {
		DepartmentID uuid.UUID `gorm:"column:department_id"`
		Count        int       `gorm:"column:count"`
	}

	var results []countResult
	err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Select("department_id, COUNT(*) as count").
		Where("tenant_id = ? AND department_id IN ? AND deleted_at IS NULL", tenantID, departmentIDs).
		Group("department_id").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get staff counts: %w", err)
	}

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.DepartmentID] = r.Count
	}

	return counts, nil
}

// GetActiveForDropdown retrieves active departments for dropdown/select.
func (r *Repository) GetActiveForDropdown(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) ([]models.Department, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	var departments []models.Department
	err := query.Order("name ASC").Find(&departments).Error
	if err != nil {
		return nil, fmt.Errorf("get active departments: %w", err)
	}
	return departments, nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
