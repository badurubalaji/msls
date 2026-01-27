// Package salary provides salary management functionality.
package salary

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// Service provides business logic for salary operations.
type Service struct {
	repo *Repository
}

// NewService creates a new salary service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Salary Component Service Methods
// ========================================

// CreateComponent creates a new salary component.
func (s *Service) CreateComponent(ctx context.Context, dto CreateComponentDTO) (*models.SalaryComponent, error) {
	// Validate component type
	if !dto.ComponentType.IsValid() {
		return nil, ErrInvalidComponentType
	}

	// Validate calculation type
	if !dto.CalculationType.IsValid() {
		return nil, ErrInvalidCalcType
	}

	// For percentage-based components, PercentageOfID is required
	if dto.CalculationType == models.CalculationTypePercentage && dto.PercentageOfID == nil {
		return nil, ErrPercentageOfRequired
	}

	component := &models.SalaryComponent{
		ID:              uuid.New(),
		TenantID:        dto.TenantID,
		Name:            dto.Name,
		Code:            dto.Code,
		Description:     dto.Description,
		ComponentType:   dto.ComponentType,
		CalculationType: dto.CalculationType,
		PercentageOfID:  dto.PercentageOfID,
		IsTaxable:       dto.IsTaxable,
		IsActive:        true,
		DisplayOrder:    dto.DisplayOrder,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateComponent(ctx, component); err != nil {
		return nil, err
	}

	return s.repo.GetComponentByID(ctx, dto.TenantID, component.ID)
}

// GetComponentByID retrieves a salary component by ID.
func (s *Service) GetComponentByID(ctx context.Context, tenantID, id uuid.UUID) (*models.SalaryComponent, error) {
	return s.repo.GetComponentByID(ctx, tenantID, id)
}

// UpdateComponent updates a salary component.
func (s *Service) UpdateComponent(ctx context.Context, tenantID, id uuid.UUID, dto UpdateComponentDTO) (*models.SalaryComponent, error) {
	component, err := s.repo.GetComponentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		component.Name = *dto.Name
	}
	if dto.Code != nil {
		component.Code = *dto.Code
	}
	if dto.Description != nil {
		component.Description = dto.Description
	}
	if dto.ComponentType != nil {
		if !dto.ComponentType.IsValid() {
			return nil, ErrInvalidComponentType
		}
		component.ComponentType = *dto.ComponentType
	}
	if dto.CalculationType != nil {
		if !dto.CalculationType.IsValid() {
			return nil, ErrInvalidCalcType
		}
		component.CalculationType = *dto.CalculationType
	}
	if dto.PercentageOfID != nil {
		component.PercentageOfID = dto.PercentageOfID
	}
	if dto.IsTaxable != nil {
		component.IsTaxable = *dto.IsTaxable
	}
	if dto.IsActive != nil {
		component.IsActive = *dto.IsActive
	}
	if dto.DisplayOrder != nil {
		component.DisplayOrder = *dto.DisplayOrder
	}
	component.UpdatedAt = time.Now()

	if err := s.repo.UpdateComponent(ctx, component); err != nil {
		return nil, err
	}

	return s.repo.GetComponentByID(ctx, tenantID, id)
}

// DeleteComponent deletes a salary component.
func (s *Service) DeleteComponent(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteComponent(ctx, tenantID, id)
}

// ListComponents retrieves salary components with filters.
func (s *Service) ListComponents(ctx context.Context, filter ComponentFilter) ([]models.SalaryComponent, int64, error) {
	return s.repo.ListComponents(ctx, filter)
}

// GetActiveComponents retrieves active components for dropdown.
func (s *Service) GetActiveComponents(ctx context.Context, tenantID uuid.UUID, componentType *models.ComponentType) ([]models.SalaryComponent, error) {
	return s.repo.GetActiveComponents(ctx, tenantID, componentType)
}

// ========================================
// Salary Structure Service Methods
// ========================================

// CreateStructure creates a new salary structure.
func (s *Service) CreateStructure(ctx context.Context, dto CreateStructureDTO) (*models.SalaryStructure, error) {
	if len(dto.Components) == 0 {
		return nil, ErrNoComponentsInStructure
	}

	structure := &models.SalaryStructure{
		ID:            uuid.New(),
		TenantID:      dto.TenantID,
		Name:          dto.Name,
		Code:          dto.Code,
		Description:   dto.Description,
		DesignationID: dto.DesignationID,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create structure components
	for _, comp := range dto.Components {
		structure.Components = append(structure.Components, models.SalaryStructureComponent{
			ID:          uuid.New(),
			StructureID: structure.ID,
			ComponentID: comp.ComponentID,
			Amount:      comp.Amount,
			Percentage:  comp.Percentage,
			CreatedAt:   time.Now(),
		})
	}

	if err := s.repo.CreateStructure(ctx, structure); err != nil {
		return nil, err
	}

	return s.repo.GetStructureByID(ctx, dto.TenantID, structure.ID)
}

// GetStructureByID retrieves a salary structure by ID.
func (s *Service) GetStructureByID(ctx context.Context, tenantID, id uuid.UUID) (*models.SalaryStructure, error) {
	return s.repo.GetStructureByID(ctx, tenantID, id)
}

// UpdateStructure updates a salary structure.
func (s *Service) UpdateStructure(ctx context.Context, tenantID, id uuid.UUID, dto UpdateStructureDTO) (*models.SalaryStructure, error) {
	structure, err := s.repo.GetStructureByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		structure.Name = *dto.Name
	}
	if dto.Code != nil {
		structure.Code = *dto.Code
	}
	if dto.Description != nil {
		structure.Description = dto.Description
	}
	if dto.DesignationID != nil {
		structure.DesignationID = dto.DesignationID
	}
	if dto.IsActive != nil {
		structure.IsActive = *dto.IsActive
	}
	structure.UpdatedAt = time.Now()

	if err := s.repo.UpdateStructure(ctx, structure); err != nil {
		return nil, err
	}

	// Update components if provided
	if dto.Components != nil {
		var components []models.SalaryStructureComponent
		for _, comp := range dto.Components {
			components = append(components, models.SalaryStructureComponent{
				ID:          uuid.New(),
				StructureID: id,
				ComponentID: comp.ComponentID,
				Amount:      comp.Amount,
				Percentage:  comp.Percentage,
				CreatedAt:   time.Now(),
			})
		}

		if err := s.repo.ReplaceStructureComponents(ctx, id, components); err != nil {
			return nil, err
		}
	}

	return s.repo.GetStructureByID(ctx, tenantID, id)
}

// DeleteStructure deletes a salary structure.
func (s *Service) DeleteStructure(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteStructure(ctx, tenantID, id)
}

// ListStructures retrieves salary structures with filters.
func (s *Service) ListStructures(ctx context.Context, filter StructureFilter) ([]models.SalaryStructure, int64, map[uuid.UUID]int, error) {
	structures, total, err := s.repo.ListStructures(ctx, filter)
	if err != nil {
		return nil, 0, nil, err
	}

	// Get staff counts
	var structureIDs []uuid.UUID
	for _, s := range structures {
		structureIDs = append(structureIDs, s.ID)
	}

	staffCounts, err := s.repo.GetStructureStaffCounts(ctx, filter.TenantID, structureIDs)
	if err != nil {
		return nil, 0, nil, err
	}

	return structures, total, staffCounts, nil
}

// GetActiveStructures retrieves active structures for dropdown.
func (s *Service) GetActiveStructures(ctx context.Context, tenantID uuid.UUID) ([]models.SalaryStructure, error) {
	return s.repo.GetActiveStructures(ctx, tenantID)
}

// ========================================
// Staff Salary Service Methods
// ========================================

// AssignSalary assigns or revises salary for a staff member.
func (s *Service) AssignSalary(ctx context.Context, dto AssignSalaryDTO) (*models.StaffSalary, error) {
	// Calculate gross and net salary from components
	grossSalary := decimal.Zero
	netSalary := decimal.Zero

	for _, comp := range dto.Components {
		// Get component to check type
		component, err := s.repo.GetComponentByID(ctx, dto.TenantID, comp.ComponentID)
		if err != nil {
			return nil, err
		}

		if component.ComponentType == models.ComponentTypeEarning {
			grossSalary = grossSalary.Add(comp.Amount)
		} else {
			// Deduction
			netSalary = netSalary.Sub(comp.Amount)
		}
	}

	netSalary = grossSalary.Add(netSalary) // netSalary was negative for deductions

	salary := &models.StaffSalary{
		ID:             uuid.New(),
		TenantID:       dto.TenantID,
		StaffID:        dto.StaffID,
		StructureID:    dto.StructureID,
		EffectiveFrom:  dto.EffectiveFrom,
		GrossSalary:    grossSalary,
		NetSalary:      netSalary,
		RevisionReason: dto.RevisionReason,
		IsCurrent:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      dto.CreatedBy,
	}

	// Create salary components
	for _, comp := range dto.Components {
		salary.Components = append(salary.Components, models.StaffSalaryComponent{
			ID:            uuid.New(),
			StaffSalaryID: salary.ID,
			ComponentID:   comp.ComponentID,
			Amount:        comp.Amount,
			IsOverridden:  comp.IsOverridden,
			CreatedAt:     time.Now(),
		})
	}

	if err := s.repo.CreateStaffSalary(ctx, salary); err != nil {
		return nil, err
	}

	return s.repo.GetCurrentStaffSalary(ctx, dto.TenantID, dto.StaffID)
}

// GetCurrentStaffSalary retrieves the current salary for a staff member.
func (s *Service) GetCurrentStaffSalary(ctx context.Context, tenantID, staffID uuid.UUID) (*models.StaffSalary, error) {
	return s.repo.GetCurrentStaffSalary(ctx, tenantID, staffID)
}

// GetStaffSalaryHistory retrieves salary history for a staff member.
func (s *Service) GetStaffSalaryHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.StaffSalary, int64, error) {
	return s.repo.GetStaffSalaryHistory(ctx, tenantID, staffID)
}
