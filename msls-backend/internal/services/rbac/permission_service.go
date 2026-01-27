// Package rbac provides role-based access control services.
package rbac

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// PermissionService handles permission-related operations.
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService creates a new PermissionService instance.
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// permissionCodeRegex validates permission code format: "module:action"
var permissionCodeRegex = regexp.MustCompile(`^[a-z_]+:[a-z_]+$`)

// CreatePermissionRequest represents a request to create a permission.
type CreatePermissionRequest struct {
	Code        string
	Name        string
	Module      string
	Description string
}

// UpdatePermissionRequest represents a request to update a permission.
type UpdatePermissionRequest struct {
	Name        *string
	Description *string
}

// ListPermissionsFilter contains filters for listing permissions.
type ListPermissionsFilter struct {
	Module string
	Search string
}

// Create creates a new permission.
func (s *PermissionService) Create(ctx context.Context, req CreatePermissionRequest) (*models.Permission, error) {
	// Validate code format
	if req.Code == "" {
		return nil, ErrPermissionCodeRequired
	}
	if !permissionCodeRegex.MatchString(req.Code) {
		return nil, ErrInvalidPermissionCode
	}

	// Check if permission code already exists
	var existing models.Permission
	err := s.db.WithContext(ctx).Where("code = ?", req.Code).First(&existing).Error
	if err == nil {
		return nil, ErrPermissionCodeExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Extract module from code if not provided
	module := req.Module
	if module == "" {
		parts := strings.Split(req.Code, ":")
		if len(parts) == 2 {
			module = parts[0]
		}
	}

	permission := &models.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Module:      module,
		Description: req.Description,
	}

	if err := s.db.WithContext(ctx).Create(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

// GetByID retrieves a permission by ID.
func (s *PermissionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := s.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// GetByCode retrieves a permission by code.
func (s *PermissionService) GetByCode(ctx context.Context, code string) (*models.Permission, error) {
	var permission models.Permission
	err := s.db.WithContext(ctx).First(&permission, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// List retrieves all permissions with optional filtering.
func (s *PermissionService) List(ctx context.Context, filter ListPermissionsFilter) ([]models.Permission, error) {
	query := s.db.WithContext(ctx).Model(&models.Permission{}).Order("module, code")

	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ?", search, search)
	}

	var permissions []models.Permission
	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// ListByModule retrieves all permissions for a specific module.
func (s *PermissionService) ListByModule(ctx context.Context, module string) ([]models.Permission, error) {
	return s.List(ctx, ListPermissionsFilter{Module: module})
}

// GetModules retrieves all distinct permission modules.
func (s *PermissionService) GetModules(ctx context.Context) ([]string, error) {
	var modules []string
	err := s.db.WithContext(ctx).
		Model(&models.Permission{}).
		Distinct("module").
		Order("module").
		Pluck("module", &modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// Update updates a permission.
func (s *PermissionService) Update(ctx context.Context, id uuid.UUID, req UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(permission).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetByID(ctx, id)
}

// Delete deletes a permission.
func (s *PermissionService) Delete(ctx context.Context, id uuid.UUID) error {
	permission, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Remove from all roles first (clear the many-to-many relationship)
	if err := s.db.WithContext(ctx).Model(permission).Association("Roles").Clear(); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Delete(permission).Error
}

// GetByIDs retrieves permissions by their IDs.
func (s *PermissionService) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetByCodes retrieves permissions by their codes.
func (s *PermissionService) GetByCodes(ctx context.Context, codes []string) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.WithContext(ctx).Where("code IN ?", codes).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// SeedDefaultPermissions creates the default system permissions.
func (s *PermissionService) SeedDefaultPermissions(ctx context.Context) error {
	permissions := []CreatePermissionRequest{
		// Users module
		{Code: "users:read", Name: "View Users", Module: ModuleUsers, Description: "View user accounts"},
		{Code: "users:write", Name: "Manage Users", Module: ModuleUsers, Description: "Create and update user accounts"},
		{Code: "users:delete", Name: "Delete Users", Module: ModuleUsers, Description: "Delete user accounts"},

		// Students module
		{Code: "students:read", Name: "View Students", Module: ModuleStudents, Description: "View student records"},
		{Code: "students:write", Name: "Manage Students", Module: ModuleStudents, Description: "Create and update student records"},
		{Code: "students:delete", Name: "Delete Students", Module: ModuleStudents, Description: "Delete student records"},

		// Staff module
		{Code: "staff:read", Name: "View Staff", Module: ModuleStaff, Description: "View staff records"},
		{Code: "staff:write", Name: "Manage Staff", Module: ModuleStaff, Description: "Create and update staff records"},
		{Code: "staff:delete", Name: "Delete Staff", Module: ModuleStaff, Description: "Delete staff records"},

		// Finance module
		{Code: "finance:read", Name: "View Finance", Module: ModuleFinance, Description: "View financial records"},
		{Code: "finance:write", Name: "Manage Finance", Module: ModuleFinance, Description: "Create and update financial records"},
		{Code: "finance:delete", Name: "Delete Finance", Module: ModuleFinance, Description: "Delete financial records"},

		// Academics module
		{Code: "academics:read", Name: "View Academics", Module: ModuleAcademics, Description: "View academic records"},
		{Code: "academics:write", Name: "Manage Academics", Module: ModuleAcademics, Description: "Create and update academic records"},
		{Code: "academics:delete", Name: "Delete Academics", Module: ModuleAcademics, Description: "Delete academic records"},

		// Roles module
		{Code: "roles:read", Name: "View Roles", Module: ModuleRoles, Description: "View roles and permissions"},
		{Code: "roles:write", Name: "Manage Roles", Module: ModuleRoles, Description: "Create and update roles"},
		{Code: "roles:delete", Name: "Delete Roles", Module: ModuleRoles, Description: "Delete roles"},

		// Settings module
		{Code: "settings:read", Name: "View Settings", Module: ModuleSettings, Description: "View system settings"},
		{Code: "settings:write", Name: "Manage Settings", Module: ModuleSettings, Description: "Update system settings"},

		// Reports module
		{Code: "reports:read", Name: "View Reports", Module: ModuleReports, Description: "View reports"},
		{Code: "reports:write", Name: "Generate Reports", Module: ModuleReports, Description: "Generate and export reports"},

		// Branches module
		{Code: "branches:read", Name: "View Branches", Module: ModuleBranches, Description: "View branch information"},
		{Code: "branches:create", Name: "Create Branches", Module: ModuleBranches, Description: "Create new branches"},
		{Code: "branches:update", Name: "Update Branches", Module: ModuleBranches, Description: "Update branch information"},
		{Code: "branches:delete", Name: "Delete Branches", Module: ModuleBranches, Description: "Delete branches"},
	}

	for _, req := range permissions {
		// Check if permission already exists
		var existing models.Permission
		err := s.db.WithContext(ctx).Where("code = ?", req.Code).First(&existing).Error
		if err == nil {
			continue // Permission already exists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check permission %s: %w", req.Code, err)
		}

		// Create permission
		if _, err := s.Create(ctx, req); err != nil {
			return fmt.Errorf("failed to create permission %s: %w", req.Code, err)
		}
	}

	return nil
}
