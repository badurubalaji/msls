// Package salary provides salary management functionality.
package salary

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for salary management.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new salary repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Salary Component Repository Methods
// ========================================

// CreateComponent creates a new salary component.
func (r *Repository) CreateComponent(ctx context.Context, component *models.SalaryComponent) error {
	if err := r.db.WithContext(ctx).Create(component).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_salary_component_code") {
			return ErrDuplicateCode
		}
		return fmt.Errorf("create salary component: %w", err)
	}
	return nil
}

// GetComponentByID retrieves a salary component by ID.
func (r *Repository) GetComponentByID(ctx context.Context, tenantID, id uuid.UUID) (*models.SalaryComponent, error) {
	var component models.SalaryComponent
	err := r.db.WithContext(ctx).
		Preload("PercentageOf").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&component).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrComponentNotFound
		}
		return nil, fmt.Errorf("get salary component by id: %w", err)
	}
	return &component, nil
}

// UpdateComponent updates a salary component.
func (r *Repository) UpdateComponent(ctx context.Context, component *models.SalaryComponent) error {
	result := r.db.WithContext(ctx).
		Model(component).
		Updates(map[string]interface{}{
			"name":             component.Name,
			"code":             component.Code,
			"description":      component.Description,
			"component_type":   component.ComponentType,
			"calculation_type": component.CalculationType,
			"percentage_of_id": component.PercentageOfID,
			"is_taxable":       component.IsTaxable,
			"is_active":        component.IsActive,
			"display_order":    component.DisplayOrder,
			"updated_at":       component.UpdatedAt,
		})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "uniq_salary_component_code") {
			return ErrDuplicateCode
		}
		return fmt.Errorf("update salary component: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrComponentNotFound
	}

	return nil
}

// DeleteComponent deletes a salary component.
func (r *Repository) DeleteComponent(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if component is used in any structure
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.SalaryStructureComponent{}).
		Where("component_id = ?", id).
		Count(&count).Error; err != nil {
		return fmt.Errorf("check component usage: %w", err)
	}
	if count > 0 {
		return ErrComponentInUse
	}

	// Check if component is used in any staff salary
	if err := r.db.WithContext(ctx).
		Model(&models.StaffSalaryComponent{}).
		Where("component_id = ?", id).
		Count(&count).Error; err != nil {
		return fmt.Errorf("check component usage in staff salary: %w", err)
	}
	if count > 0 {
		return ErrComponentInUse
	}

	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.SalaryComponent{})

	if result.Error != nil {
		return fmt.Errorf("delete salary component: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrComponentNotFound
	}

	return nil
}

// ListComponents retrieves salary components with filters.
func (r *Repository) ListComponents(ctx context.Context, filter ComponentFilter) ([]models.SalaryComponent, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.SalaryComponent{}).
		Preload("PercentageOf").
		Where("tenant_id = ?", filter.TenantID)

	if filter.ComponentType != nil {
		query = query.Where("component_type = ?", *filter.ComponentType)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count salary components: %w", err)
	}

	var components []models.SalaryComponent
	if err := query.Order("display_order ASC, name ASC").Find(&components).Error; err != nil {
		return nil, 0, fmt.Errorf("list salary components: %w", err)
	}

	return components, total, nil
}

// GetActiveComponents retrieves active components for dropdown.
func (r *Repository) GetActiveComponents(ctx context.Context, tenantID uuid.UUID, componentType *models.ComponentType) ([]models.SalaryComponent, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID)

	if componentType != nil {
		query = query.Where("component_type = ?", *componentType)
	}

	var components []models.SalaryComponent
	if err := query.Order("display_order ASC, name ASC").Find(&components).Error; err != nil {
		return nil, fmt.Errorf("get active components: %w", err)
	}

	return components, nil
}

// ========================================
// Salary Structure Repository Methods
// ========================================

// CreateStructure creates a new salary structure.
func (r *Repository) CreateStructure(ctx context.Context, structure *models.SalaryStructure) error {
	if err := r.db.WithContext(ctx).Create(structure).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_salary_structure_code") {
			return ErrDuplicateStructureCode
		}
		return fmt.Errorf("create salary structure: %w", err)
	}
	return nil
}

// GetStructureByID retrieves a salary structure by ID with components.
func (r *Repository) GetStructureByID(ctx context.Context, tenantID, id uuid.UUID) (*models.SalaryStructure, error) {
	var structure models.SalaryStructure
	err := r.db.WithContext(ctx).
		Preload("Designation").
		Preload("Components").
		Preload("Components.Component").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&structure).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStructureNotFound
		}
		return nil, fmt.Errorf("get salary structure by id: %w", err)
	}
	return &structure, nil
}

// UpdateStructure updates a salary structure.
func (r *Repository) UpdateStructure(ctx context.Context, structure *models.SalaryStructure) error {
	result := r.db.WithContext(ctx).
		Model(structure).
		Updates(map[string]interface{}{
			"name":           structure.Name,
			"code":           structure.Code,
			"description":    structure.Description,
			"designation_id": structure.DesignationID,
			"is_active":      structure.IsActive,
			"updated_at":     structure.UpdatedAt,
		})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "uniq_salary_structure_code") {
			return ErrDuplicateStructureCode
		}
		return fmt.Errorf("update salary structure: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStructureNotFound
	}

	return nil
}

// DeleteStructure deletes a salary structure.
func (r *Repository) DeleteStructure(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if structure is in use
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.StaffSalary{}).
		Where("structure_id = ?", id).
		Count(&count).Error; err != nil {
		return fmt.Errorf("check structure usage: %w", err)
	}
	if count > 0 {
		return ErrStructureInUse
	}

	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.SalaryStructure{})

	if result.Error != nil {
		return fmt.Errorf("delete salary structure: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStructureNotFound
	}

	return nil
}

// ListStructures retrieves salary structures with filters.
func (r *Repository) ListStructures(ctx context.Context, filter StructureFilter) ([]models.SalaryStructure, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.SalaryStructure{}).
		Preload("Designation").
		Preload("Components").
		Where("tenant_id = ?", filter.TenantID)

	if filter.DesignationID != nil {
		query = query.Where("designation_id = ?", *filter.DesignationID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count salary structures: %w", err)
	}

	var structures []models.SalaryStructure
	if err := query.Order("name ASC").Find(&structures).Error; err != nil {
		return nil, 0, fmt.Errorf("list salary structures: %w", err)
	}

	return structures, total, nil
}

// GetActiveStructures retrieves active structures for dropdown.
func (r *Repository) GetActiveStructures(ctx context.Context, tenantID uuid.UUID) ([]models.SalaryStructure, error) {
	var structures []models.SalaryStructure
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID).
		Order("name ASC").
		Find(&structures).Error; err != nil {
		return nil, fmt.Errorf("get active structures: %w", err)
	}
	return structures, nil
}

// GetStructureStaffCounts returns staff counts for structures.
func (r *Repository) GetStructureStaffCounts(ctx context.Context, tenantID uuid.UUID, structureIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(structureIDs) == 0 {
		return make(map[uuid.UUID]int), nil
	}

	type countResult struct {
		StructureID uuid.UUID `gorm:"column:structure_id"`
		Count       int       `gorm:"column:count"`
	}

	var results []countResult
	err := r.db.WithContext(ctx).
		Model(&models.StaffSalary{}).
		Select("structure_id, COUNT(*) as count").
		Where("tenant_id = ? AND structure_id IN ? AND is_current = true", tenantID, structureIDs).
		Group("structure_id").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get structure staff counts: %w", err)
	}

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.StructureID] = r.Count
	}

	return counts, nil
}

// ReplaceStructureComponents replaces all components in a structure.
func (r *Repository) ReplaceStructureComponents(ctx context.Context, structureID uuid.UUID, components []models.SalaryStructureComponent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing components
		if err := tx.Where("structure_id = ?", structureID).Delete(&models.SalaryStructureComponent{}).Error; err != nil {
			return fmt.Errorf("delete existing components: %w", err)
		}

		// Insert new components
		if len(components) > 0 {
			if err := tx.Create(&components).Error; err != nil {
				return fmt.Errorf("create structure components: %w", err)
			}
		}

		return nil
	})
}

// ========================================
// Staff Salary Repository Methods
// ========================================

// CreateStaffSalary creates a new staff salary assignment.
func (r *Repository) CreateStaffSalary(ctx context.Context, salary *models.StaffSalary) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Mark previous salary as not current
		if err := tx.Model(&models.StaffSalary{}).
			Where("staff_id = ? AND is_current = true", salary.StaffID).
			Updates(map[string]interface{}{
				"is_current":   false,
				"effective_to": salary.EffectiveFrom,
			}).Error; err != nil {
			return fmt.Errorf("update previous salary: %w", err)
		}

		// Create new salary
		if err := tx.Create(salary).Error; err != nil {
			return fmt.Errorf("create staff salary: %w", err)
		}

		return nil
	})
}

// GetCurrentStaffSalary retrieves the current salary for a staff member.
func (r *Repository) GetCurrentStaffSalary(ctx context.Context, tenantID, staffID uuid.UUID) (*models.StaffSalary, error) {
	var salary models.StaffSalary
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Structure").
		Preload("Components").
		Preload("Components.Component").
		Where("tenant_id = ? AND staff_id = ? AND is_current = true", tenantID, staffID).
		First(&salary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffSalaryNotFound
		}
		return nil, fmt.Errorf("get current staff salary: %w", err)
	}
	return &salary, nil
}

// GetStaffSalaryHistory retrieves salary history for a staff member.
func (r *Repository) GetStaffSalaryHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.StaffSalary, int64, error) {
	var salaries []models.StaffSalary
	err := r.db.WithContext(ctx).
		Preload("Structure").
		Preload("Components").
		Preload("Components.Component").
		Where("tenant_id = ? AND staff_id = ?", tenantID, staffID).
		Order("effective_from DESC").
		Find(&salaries).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get staff salary history: %w", err)
	}

	return salaries, int64(len(salaries)), nil
}

// DB returns the underlying database connection.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
