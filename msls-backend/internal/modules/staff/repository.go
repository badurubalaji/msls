// Package staff provides staff management functionality.
package staff

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

// Repository handles database operations for staff.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new staff repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new staff member in the database.
func (r *Repository) Create(ctx context.Context, staff *models.Staff) error {
	if err := r.db.WithContext(ctx).Create(staff).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "employee_id") {
				return ErrDuplicateEmployeeID
			}
		}
		return fmt.Errorf("create staff: %w", err)
	}
	return nil
}

// CreateWithTx creates a new staff member in the database within a transaction.
func (r *Repository) CreateWithTx(ctx context.Context, tx *gorm.DB, staff *models.Staff) error {
	if err := tx.WithContext(ctx).Create(staff).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "employee_id") {
				return ErrDuplicateEmployeeID
			}
		}
		return fmt.Errorf("create staff: %w", err)
	}
	return nil
}

// GetByID retrieves a staff member by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Department").
		Preload("Designation").
		Preload("ReportingManager").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("get staff by id: %w", err)
	}
	return &staff, nil
}

// GetByEmployeeID retrieves a staff member by employee ID.
func (r *Repository) GetByEmployeeID(ctx context.Context, tenantID uuid.UUID, employeeID string) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Department").
		Preload("Designation").
		Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("get staff by employee id: %w", err)
	}
	return &staff, nil
}

// Update updates a staff member in the database.
func (r *Repository) Update(ctx context.Context, staff *models.Staff) error {
	result := r.db.WithContext(ctx).
		Model(staff).
		Where("version = ?", staff.Version). // Optimistic locking
		Updates(staff)

	if result.Error != nil {
		return fmt.Errorf("update staff: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// UpdateWithTx updates a staff member within a transaction.
func (r *Repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, staff *models.Staff) error {
	result := tx.WithContext(ctx).
		Model(staff).
		Where("version = ?", staff.Version). // Optimistic locking
		Updates(staff)

	if result.Error != nil {
		return fmt.Errorf("update staff: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// Delete soft deletes a staff member.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Staff{})

	if result.Error != nil {
		return fmt.Errorf("delete staff: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStaffNotFound
	}

	return nil
}

// List retrieves staff with optional filtering and pagination.
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Staff, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Preload("Branch").
		Preload("Department").
		Preload("Designation").
		Where("staff.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.BranchID != nil {
		query = query.Where("staff.branch_id = ?", *filter.BranchID)
	}

	if filter.DepartmentID != nil {
		query = query.Where("staff.department_id = ?", *filter.DepartmentID)
	}

	if filter.DesignationID != nil {
		query = query.Where("staff.designation_id = ?", *filter.DesignationID)
	}

	if filter.StaffType != nil {
		query = query.Where("staff.staff_type = ?", *filter.StaffType)
	}

	if filter.Status != nil {
		query = query.Where("staff.status = ?", *filter.Status)
	}

	if filter.Gender != nil {
		query = query.Where("staff.gender = ?", *filter.Gender)
	}

	// Filter by join date range
	if filter.JoinDateFrom != nil {
		query = query.Where("staff.join_date >= ?", *filter.JoinDateFrom)
	}
	if filter.JoinDateTo != nil {
		query = query.Where("staff.join_date <= ?", *filter.JoinDateTo)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(staff.first_name) LIKE ? OR LOWER(staff.last_name) LIKE ? OR LOWER(staff.first_name || ' ' || staff.last_name) LIKE ? OR LOWER(staff.employee_id) LIKE ? OR LOWER(staff.work_phone) LIKE ?",
			search, search, search, search, search,
		)
	}

	// Get total count
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, "", 0, fmt.Errorf("count staff: %w", err)
	}

	// Apply cursor pagination
	if filter.Cursor != "" {
		cursorID, err := uuid.Parse(filter.Cursor)
		if err == nil {
			query = query.Where("staff.id > ?", cursorID)
		}
	}

	// Apply sorting
	sortBy := "staff.last_name, staff.first_name"
	if filter.SortBy != "" {
		sortOrder := "ASC"
		if strings.ToLower(filter.SortOrder) == "desc" {
			sortOrder = "DESC"
		}
		sortBy = fmt.Sprintf("staff.%s %s", filter.SortBy, sortOrder)
	}
	query = query.Order(sortBy)

	// Apply limit
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit + 1)

	var staffList []models.Staff
	if err := query.Find(&staffList).Error; err != nil {
		return nil, "", 0, fmt.Errorf("list staff: %w", err)
	}

	// Calculate next cursor
	var nextCursor string
	if len(staffList) > limit {
		staffList = staffList[:limit]
		nextCursor = staffList[len(staffList)-1].ID.String()
	}

	return staffList, nextCursor, total, nil
}

// Count returns the total number of staff for a tenant.
func (r *Repository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count staff: %w", err)
	}
	return count, nil
}

// GetNextSequence gets and increments the employee ID sequence.
func (r *Repository) GetNextSequence(ctx context.Context, tx *gorm.DB, tenantID uuid.UUID, prefix string) (int, error) {
	var seq models.StaffEmployeeSequence

	// Try to find existing sequence
	err := tx.WithContext(ctx).
		Where("tenant_id = ? AND prefix = ?", tenantID, prefix).
		First(&seq).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new sequence
		seq = models.StaffEmployeeSequence{
			TenantID:     tenantID,
			Prefix:       prefix,
			LastSequence: 1,
		}
		if err := tx.WithContext(ctx).Create(&seq).Error; err != nil {
			return 0, fmt.Errorf("create sequence: %w", err)
		}
		return seq.LastSequence, nil
	}

	if err != nil {
		return 0, fmt.Errorf("get sequence: %w", err)
	}

	// Increment sequence
	seq.LastSequence++
	seq.UpdatedAt = time.Now()
	if err := tx.WithContext(ctx).Save(&seq).Error; err != nil {
		return 0, fmt.Errorf("update sequence: %w", err)
	}

	return seq.LastSequence, nil
}

// UpdateStatus updates the status of a staff member.
func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status models.StaffStatus, reason string, updatedBy *uuid.UUID) error {
	updates := map[string]interface{}{
		"status":        status,
		"status_reason": reason,
		"updated_by":    updatedBy,
		"updated_at":    time.Now(),
	}

	// If terminated, set termination date
	if status == models.StaffStatusTerminated {
		updates["termination_date"] = time.Now()
	}

	result := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStaffNotFound
	}

	return nil
}

// UpdatePhoto updates the photo URL for a staff member.
func (r *Repository) UpdatePhoto(ctx context.Context, tenantID, id uuid.UUID, photoURL string, updatedBy *uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.Staff{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{
			"photo_url":  photoURL,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("update photo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStaffNotFound
	}

	return nil
}

// CreateStatusHistory creates a status history record.
func (r *Repository) CreateStatusHistory(ctx context.Context, tx *gorm.DB, history *models.StaffStatusHistory) error {
	if err := tx.WithContext(ctx).Create(history).Error; err != nil {
		return fmt.Errorf("create status history: %w", err)
	}
	return nil
}

// GetStatusHistory retrieves the status history for a staff member.
func (r *Repository) GetStatusHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.StaffStatusHistory, error) {
	var histories []models.StaffStatusHistory
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND staff_id = ?", tenantID, staffID).
		Order("changed_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, fmt.Errorf("get status history: %w", err)
	}
	return histories, nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
