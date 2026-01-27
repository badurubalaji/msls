// Package rbac provides role-based access control services.
package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// RoleService handles role-related operations.
type RoleService struct {
	db                *gorm.DB
	permissionService *PermissionService
}

// NewRoleService creates a new RoleService instance.
func NewRoleService(db *gorm.DB, permissionService *PermissionService) *RoleService {
	return &RoleService{
		db:                db,
		permissionService: permissionService,
	}
}

// CreateRoleRequest represents a request to create a role.
type CreateRoleRequest struct {
	TenantID      *uuid.UUID
	Name          string
	Description   string
	IsSystem      bool
	PermissionIDs []uuid.UUID
}

// UpdateRoleRequest represents a request to update a role.
type UpdateRoleRequest struct {
	Name        *string
	Description *string
}

// ListRolesFilter contains filters for listing roles.
type ListRolesFilter struct {
	TenantID     *uuid.UUID
	IncludeSystem bool
	Search       string
}

// AssignPermissionsRequest represents a request to assign permissions to a role.
type AssignPermissionsRequest struct {
	PermissionIDs []uuid.UUID
}

// Create creates a new role.
func (s *RoleService) Create(ctx context.Context, req CreateRoleRequest) (*models.Role, error) {
	if req.Name == "" {
		return nil, ErrRoleNameRequired
	}

	// Check if role name already exists for the tenant
	query := s.db.WithContext(ctx).Where("name = ?", req.Name)
	if req.TenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", req.TenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	var existing models.Role
	err := query.First(&existing).Error
	if err == nil {
		return nil, ErrRoleNameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get permissions if IDs are provided
	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		permissions, err = s.permissionService.GetByIDs(ctx, req.PermissionIDs)
		if err != nil {
			return nil, err
		}
		if len(permissions) != len(req.PermissionIDs) {
			return nil, ErrPermissionNotFound
		}
	}

	role := &models.Role{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    req.IsSystem,
		Permissions: permissions,
	}

	if err := s.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}

	// Reload with permissions
	return s.GetByID(ctx, role.ID)
}

// GetByID retrieves a role by ID with its permissions.
func (s *RoleService) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := s.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name within a tenant.
func (s *RoleService) GetByName(ctx context.Context, name string, tenantID *uuid.UUID) (*models.Role, error) {
	query := s.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name)

	if tenantID != nil {
		query = query.Where("tenant_id = ? OR tenant_id IS NULL", tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	var role models.Role
	err := query.First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// List retrieves roles with optional filtering.
func (s *RoleService) List(ctx context.Context, filter ListRolesFilter) ([]models.Role, error) {
	query := s.db.WithContext(ctx).
		Model(&models.Role{}).
		Preload("Permissions").
		Order("is_system DESC, name")

	if filter.TenantID != nil {
		if filter.IncludeSystem {
			query = query.Where("tenant_id = ? OR tenant_id IS NULL", filter.TenantID)
		} else {
			query = query.Where("tenant_id = ?", filter.TenantID)
		}
	} else if !filter.IncludeSystem {
		query = query.Where("tenant_id IS NOT NULL")
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	var roles []models.Role
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// ListByTenant retrieves all roles for a specific tenant including system roles.
func (s *RoleService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.Role, error) {
	return s.List(ctx, ListRolesFilter{
		TenantID:      &tenantID,
		IncludeSystem: true,
	})
}

// Update updates a role.
func (s *RoleService) Update(ctx context.Context, id uuid.UUID, req UpdateRoleRequest) (*models.Role, error) {
	role, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		// Check if new name conflicts with existing role
		if *req.Name != role.Name {
			query := s.db.WithContext(ctx).Where("name = ? AND id != ?", *req.Name, id)
			if role.TenantID != nil {
				query = query.Where("tenant_id = ? OR tenant_id IS NULL", role.TenantID)
			} else {
				query = query.Where("tenant_id IS NULL")
			}

			var existing models.Role
			err := query.First(&existing).Error
			if err == nil {
				return nil, ErrRoleNameExists
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(role).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetByID(ctx, id)
}

// Delete deletes a role.
func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	role, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Prevent deleting system roles
	if role.IsSystem {
		return ErrCannotDeleteSystem
	}

	// Check if role is assigned to any users
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&models.UserRole{}).
		Where("role_id = ?", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleInUse
	}

	// Clear permissions association
	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Clear(); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Delete(role).Error
}

// AssignPermissions assigns permissions to a role.
func (s *RoleService) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) (*models.Role, error) {
	role, err := s.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// System roles cannot have their permissions modified
	if role.IsSystem {
		return nil, ErrCannotModifySystem
	}

	// Get permissions
	permissions, err := s.permissionService.GetByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}
	if len(permissions) != len(permissionIDs) {
		return nil, ErrPermissionNotFound
	}

	// Append permissions (not replace)
	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Append(permissions); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, roleID)
}

// RemovePermissions removes permissions from a role.
func (s *RoleService) RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) (*models.Role, error) {
	role, err := s.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// System roles cannot have their permissions modified
	if role.IsSystem {
		return nil, ErrCannotModifySystem
	}

	// Get permissions
	permissions, err := s.permissionService.GetByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}

	// Remove permissions
	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Delete(permissions); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, roleID)
}

// SetPermissions replaces all permissions for a role.
func (s *RoleService) SetPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) (*models.Role, error) {
	role, err := s.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// System roles cannot have their permissions modified
	if role.IsSystem {
		return nil, ErrCannotModifySystem
	}

	// Get permissions
	permissions, err := s.permissionService.GetByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}
	if len(permissions) != len(permissionIDs) {
		return nil, ErrPermissionNotFound
	}

	// Replace all permissions
	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, roleID)
}

// GetHierarchyLevel returns the hierarchy level for a role.
// Lower values indicate higher privilege.
func (s *RoleService) GetHierarchyLevel(role *models.Role) int {
	if level, ok := RoleHierarchy[role.Name]; ok {
		return level
	}
	// Custom roles default to a lower privilege than students
	return 100
}

// CanAssignRole checks if a user with sourceRole can assign targetRole to another user.
func (s *RoleService) CanAssignRole(sourceRole, targetRole *models.Role) bool {
	sourceLevel := s.GetHierarchyLevel(sourceRole)
	targetLevel := s.GetHierarchyLevel(targetRole)
	// Users can only assign roles at or below their own level
	return sourceLevel <= targetLevel
}

// SeedSystemRoles creates the predefined system roles with default permissions.
func (s *RoleService) SeedSystemRoles(ctx context.Context) error {
	// Define permissions for each system role
	rolePermissions := map[string][]string{
		RoleSuperAdmin: {
			"users:read", "users:write", "users:delete",
			"students:read", "students:write", "students:delete",
			"staff:read", "staff:write", "staff:delete",
			"finance:read", "finance:write", "finance:delete",
			"academics:read", "academics:write", "academics:delete",
			"roles:read", "roles:write", "roles:delete",
			"settings:read", "settings:write",
			"reports:read", "reports:write",
			"branches:read", "branches:create", "branches:update", "branches:delete",
		},
		RoleTenantAdmin: {
			"users:read", "users:write", "users:delete",
			"students:read", "students:write", "students:delete",
			"staff:read", "staff:write", "staff:delete",
			"finance:read", "finance:write", "finance:delete",
			"academics:read", "academics:write", "academics:delete",
			"roles:read", "roles:write",
			"settings:read", "settings:write",
			"reports:read", "reports:write",
			"branches:read", "branches:create", "branches:update", "branches:delete",
		},
		RolePrincipal: {
			"users:read", "users:write",
			"students:read", "students:write",
			"staff:read", "staff:write",
			"finance:read", "finance:write",
			"academics:read", "academics:write",
			"roles:read",
			"reports:read", "reports:write",
			"branches:read", "branches:update",
		},
		RoleTeacher: {
			"students:read",
			"academics:read", "academics:write",
			"reports:read",
			"branches:read",
		},
		RoleStaff: {
			"students:read",
			"staff:read",
			"academics:read",
			"branches:read",
		},
		RoleParent: {
			"students:read",
			"academics:read",
			"finance:read",
			"branches:read",
		},
		RoleStudent: {
			"academics:read",
			"branches:read",
		},
	}

	// Get all permissions
	allPermissions, err := s.permissionService.List(ctx, ListPermissionsFilter{})
	if err != nil {
		return err
	}

	// Create a map for quick lookup
	permMap := make(map[string]models.Permission)
	for _, p := range allPermissions {
		permMap[p.Code] = p
	}

	// Create each system role
	for _, roleName := range SystemRoles() {
		// Check if role already exists
		var existing models.Role
		err := s.db.WithContext(ctx).Where("name = ? AND is_system = ?", roleName, true).First(&existing).Error
		if err == nil {
			continue // Role already exists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check role %s: %w", roleName, err)
		}

		// Get permissions for this role
		permCodes := rolePermissions[roleName]
		var permissions []models.Permission
		for _, code := range permCodes {
			if p, ok := permMap[code]; ok {
				permissions = append(permissions, p)
			}
		}

		// Create role
		role := &models.Role{
			Name:        roleName,
			Description: fmt.Sprintf("System role: %s", roleName),
			IsSystem:    true,
			Permissions: permissions,
		}

		if err := s.db.WithContext(ctx).Create(role).Error; err != nil {
			return fmt.Errorf("failed to create role %s: %w", roleName, err)
		}
	}

	return nil
}
