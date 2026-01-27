// Package designation provides designation management functionality.
package designation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// Service provides business logic for designation operations.
type Service struct {
	repo *Repository
}

// NewService creates a new designation service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new designation.
func (s *Service) Create(ctx context.Context, dto CreateDesignationDTO) (*models.Designation, error) {
	// Validate level
	if dto.Level < 1 || dto.Level > 10 {
		return nil, ErrInvalidLevel
	}

	designation := &models.Designation{
		ID:           uuid.New(),
		TenantID:     dto.TenantID,
		Name:         dto.Name,
		Level:        dto.Level,
		DepartmentID: dto.DepartmentID,
		IsActive:     dto.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, designation); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.repo.GetByID(ctx, dto.TenantID, designation.ID)
}

// GetByID retrieves a designation by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Designation, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Update updates a designation.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateDesignationDTO) (*models.Designation, error) {
	designation, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		designation.Name = *dto.Name
	}
	if dto.Level != nil {
		if *dto.Level < 1 || *dto.Level > 10 {
			return nil, ErrInvalidLevel
		}
		designation.Level = *dto.Level
	}
	if dto.DepartmentID != nil {
		designation.DepartmentID = dto.DepartmentID
	}
	if dto.IsActive != nil {
		designation.IsActive = *dto.IsActive
	}
	designation.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, designation); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.repo.GetByID(ctx, tenantID, id)
}

// Delete deletes a designation.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List retrieves designations with filters.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Designation, int64, map[uuid.UUID]int, error) {
	designations, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, nil, err
	}

	// Get staff counts for each designation
	var designationIDs []uuid.UUID
	for _, d := range designations {
		designationIDs = append(designationIDs, d.ID)
	}

	staffCounts, err := s.repo.GetStaffCounts(ctx, filter.TenantID, designationIDs)
	if err != nil {
		return nil, 0, nil, err
	}

	return designations, total, staffCounts, nil
}

// GetActiveForDropdown retrieves active designations for dropdown/select.
func (s *Service) GetActiveForDropdown(ctx context.Context, tenantID uuid.UUID, departmentID *uuid.UUID) ([]models.Designation, error) {
	return s.repo.GetActiveForDropdown(ctx, tenantID, departmentID)
}
