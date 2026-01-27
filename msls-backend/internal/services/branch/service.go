// Package branch provides branch management services.
package branch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Service handles branch-related operations.
type Service struct {
	db *gorm.DB
}

// NewService creates a new Branch Service instance.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateRequest represents a request to create a branch.
type CreateRequest struct {
	TenantID     uuid.UUID
	Code         string
	Name         string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
	Phone        string
	Email        string
	LogoURL      string
	Timezone     string
	IsPrimary    bool
	Settings     map[string]interface{}
	CreatedBy    *uuid.UUID
}

// UpdateRequest represents a request to update a branch.
type UpdateRequest struct {
	Name         *string
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	State        *string
	PostalCode   *string
	Country      *string
	Phone        *string
	Email        *string
	LogoURL      *string
	Timezone     *string
	IsPrimary    *bool
	Settings     map[string]interface{}
	UpdatedBy    *uuid.UUID
}

// ListFilter contains filters for listing branches.
type ListFilter struct {
	TenantID  uuid.UUID
	Status    *models.Status
	IsPrimary *bool
	Search    string
}

// Create creates a new branch.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Branch, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.Name == "" {
		return nil, ErrBranchNameRequired
	}
	if req.Code == "" {
		return nil, ErrBranchCodeRequired
	}

	// Validate timezone if provided
	if req.Timezone != "" {
		if _, err := time.LoadLocation(req.Timezone); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidTimezone, req.Timezone)
		}
	}

	// Check if branch code already exists for this tenant
	var existing models.Branch
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", req.TenantID, req.Code).
		First(&existing).Error
	if err == nil {
		return nil, ErrBranchCodeExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing branch: %w", err)
	}

	// Set default country if not provided
	country := req.Country
	if country == "" {
		country = "India"
	}

	// Set default timezone if not provided
	timezone := req.Timezone
	if timezone == "" {
		timezone = "Asia/Kolkata"
	}

	// Build address
	address := models.BranchAddress{
		Street:     req.AddressLine1,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    country,
	}
	// Store AddressLine2 in the Street field if AddressLine1 is there
	if req.AddressLine2 != "" && req.AddressLine1 != "" {
		address.Street = req.AddressLine1 + "\n" + req.AddressLine2
	}

	// Build settings
	settings := models.BranchSettings{
		ContactPhone: req.Phone,
		ContactEmail: req.Email,
	}
	if req.Settings != nil {
		if features, ok := req.Settings["features"].(map[string]bool); ok {
			settings.Features = features
		}
		if hours, ok := req.Settings["operating_hours"].(map[string]string); ok {
			settings.OperatingHours = hours
		}
	}

	branch := &models.Branch{
		TenantID:  req.TenantID,
		Code:      req.Code,
		Name:      req.Name,
		Address:   address,
		Settings:  settings,
		IsPrimary: req.IsPrimary,
		Status:    models.StatusActive,
	}
	branch.CreatedBy = req.CreatedBy
	branch.UpdatedBy = req.CreatedBy

	// Create branch in a transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If this branch is primary, unset other primary branches
		if req.IsPrimary {
			if err := tx.Model(&models.Branch{}).
				Where("tenant_id = ? AND is_primary = ?", req.TenantID, true).
				Update("is_primary", false).Error; err != nil {
				return fmt.Errorf("failed to unset primary branch: %w", err)
			}
		}

		if err := tx.Create(branch).Error; err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, req.TenantID, branch.ID)
}

// GetByID retrieves a branch by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&branch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBranchNotFound
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return &branch, nil
}

// GetByCode retrieves a branch by code within a tenant.
func (s *Service) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.Branch, error) {
	var branch models.Branch
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&branch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBranchNotFound
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return &branch, nil
}

// List retrieves branches with optional filtering.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Branch, error) {
	query := s.db.WithContext(ctx).
		Model(&models.Branch{}).
		Where("tenant_id = ?", filter.TenantID).
		Order("is_primary DESC, name")

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filter.IsPrimary)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR address->>'city' ILIKE ?", search, search, search)
	}

	var branches []models.Branch
	if err := query.Find(&branches).Error; err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	return branches, nil
}

// Update updates a branch.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, req UpdateRequest) (*models.Branch, error) {
	branch, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	// Update address fields
	address := branch.Address
	addressUpdated := false

	if req.AddressLine1 != nil {
		address.Street = *req.AddressLine1
		if req.AddressLine2 != nil && *req.AddressLine2 != "" {
			address.Street = *req.AddressLine1 + "\n" + *req.AddressLine2
		}
		addressUpdated = true
	}
	if req.City != nil {
		address.City = *req.City
		addressUpdated = true
	}
	if req.State != nil {
		address.State = *req.State
		addressUpdated = true
	}
	if req.PostalCode != nil {
		address.PostalCode = *req.PostalCode
		addressUpdated = true
	}
	if req.Country != nil {
		address.Country = *req.Country
		addressUpdated = true
	}

	if addressUpdated {
		updates["address"] = address
	}

	// Update settings fields
	settings := branch.Settings
	settingsUpdated := false

	if req.Phone != nil {
		settings.ContactPhone = *req.Phone
		settingsUpdated = true
	}
	if req.Email != nil {
		settings.ContactEmail = *req.Email
		settingsUpdated = true
	}
	if req.Settings != nil {
		if features, ok := req.Settings["features"].(map[string]bool); ok {
			settings.Features = features
			settingsUpdated = true
		}
		if hours, ok := req.Settings["operating_hours"].(map[string]string); ok {
			settings.OperatingHours = hours
			settingsUpdated = true
		}
	}

	if settingsUpdated {
		updates["settings"] = settings
	}

	// Validate and update timezone
	if req.Timezone != nil {
		if _, err := time.LoadLocation(*req.Timezone); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidTimezone, *req.Timezone)
		}
		// Note: timezone is stored in settings, not as a separate column
		// If you want to store it separately, add it to the branch model
	}

	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(branch).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update branch: %w", err)
		}
	}

	// Handle isPrimary change - this needs special handling as it affects other branches
	if req.IsPrimary != nil && *req.IsPrimary && !branch.IsPrimary {
		// Use SetPrimary to properly handle unsetting other primary branches
		return s.SetPrimary(ctx, tenantID, id, req.UpdatedBy)
	}

	return s.GetByID(ctx, tenantID, id)
}

// SetPrimary sets a branch as the primary branch for the tenant.
func (s *Service) SetPrimary(ctx context.Context, tenantID, id uuid.UUID, updatedBy *uuid.UUID) (*models.Branch, error) {
	branch, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if branch.IsPrimary {
		// Already primary, return as-is
		return branch, nil
	}

	// Check if branch is active
	if branch.Status != models.StatusActive {
		return nil, errors.New("cannot set inactive branch as primary")
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset current primary branch
		if err := tx.Model(&models.Branch{}).
			Where("tenant_id = ? AND is_primary = ?", tenantID, true).
			Updates(map[string]interface{}{
				"is_primary": false,
				"updated_by": updatedBy,
			}).Error; err != nil {
			return fmt.Errorf("failed to unset primary branch: %w", err)
		}

		// Set new primary branch
		if err := tx.Model(branch).Updates(map[string]interface{}{
			"is_primary": true,
			"updated_by": updatedBy,
		}).Error; err != nil {
			return fmt.Errorf("failed to set primary branch: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// SetStatus sets the status (active/inactive) of a branch.
func (s *Service) SetStatus(ctx context.Context, tenantID, id uuid.UUID, status models.Status, updatedBy *uuid.UUID) (*models.Branch, error) {
	branch, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Cannot deactivate primary branch
	if branch.IsPrimary && status != models.StatusActive {
		return nil, ErrCannotDeactivatePrimaryBranch
	}

	updates := map[string]interface{}{
		"status":     status,
		"updated_by": updatedBy,
	}

	if err := s.db.WithContext(ctx).Model(branch).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update branch status: %w", err)
	}

	return s.GetByID(ctx, tenantID, id)
}

// Delete deletes a branch.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	branch, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Cannot delete primary branch
	if branch.IsPrimary {
		return ErrCannotDeletePrimaryBranch
	}

	// TODO: Check if branch has associated records (students, staff, etc.)
	// This would require checking other tables that reference branch_id
	// For now, we'll do a hard delete since there are no dependent tables yet

	if err := s.db.WithContext(ctx).Delete(branch).Error; err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// GetPrimaryBranch retrieves the primary branch for a tenant.
func (s *Service) GetPrimaryBranch(ctx context.Context, tenantID uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND is_primary = ?", tenantID, true).
		First(&branch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBranchNotFound
		}
		return nil, fmt.Errorf("failed to get primary branch: %w", err)
	}
	return &branch, nil
}

// Count returns the total number of branches for a tenant.
func (s *Service) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.Branch{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count branches: %w", err)
	}
	return count, nil
}
