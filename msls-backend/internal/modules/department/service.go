// Package department provides department management functionality.
package department

import (
	"context"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// Service provides business logic for department operations.
type Service struct {
	repo *Repository
}

// NewService creates a new department service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new department.
func (s *Service) Create(ctx context.Context, dto CreateDepartmentDTO) (*models.Department, error) {
	department := &models.Department{
		ID:          uuid.New(),
		TenantID:    dto.TenantID,
		BranchID:    dto.BranchID,
		Name:        dto.Name,
		Code:        dto.Code,
		Description: dto.Description,
		HeadID:      dto.HeadID,
		IsActive:    dto.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, department); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.repo.GetByID(ctx, dto.TenantID, department.ID)
}

// GetByID retrieves a department by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Department, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Update updates a department.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateDepartmentDTO) (*models.Department, error) {
	department, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		department.Name = *dto.Name
	}
	if dto.Code != nil {
		department.Code = *dto.Code
	}
	if dto.Description != nil {
		department.Description = dto.Description
	}
	if dto.HeadID != nil {
		department.HeadID = dto.HeadID
	}
	if dto.IsActive != nil {
		department.IsActive = *dto.IsActive
	}
	department.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, department); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.repo.GetByID(ctx, tenantID, id)
}

// Delete deletes a department.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List retrieves departments with filters.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Department, int64, map[uuid.UUID]int, error) {
	departments, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, nil, err
	}

	// Get staff counts for each department
	var departmentIDs []uuid.UUID
	for _, d := range departments {
		departmentIDs = append(departmentIDs, d.ID)
	}

	staffCounts, err := s.repo.GetStaffCounts(ctx, filter.TenantID, departmentIDs)
	if err != nil {
		return nil, 0, nil, err
	}

	return departments, total, staffCounts, nil
}

// GetActiveForDropdown retrieves active departments for dropdown/select.
func (s *Service) GetActiveForDropdown(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) ([]models.Department, error) {
	return s.repo.GetActiveForDropdown(ctx, tenantID, branchID)
}
