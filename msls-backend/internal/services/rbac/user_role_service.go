// Package rbac provides role-based access control services.
package rbac

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// UserRoleService handles user-role assignment operations.
type UserRoleService struct {
	db          *gorm.DB
	roleService *RoleService
}

// NewUserRoleService creates a new UserRoleService instance.
func NewUserRoleService(db *gorm.DB, roleService *RoleService) *UserRoleService {
	return &UserRoleService{
		db:          db,
		roleService: roleService,
	}
}

// UserRolesResponse represents a user with their roles.
type UserRolesResponse struct {
	UserID uuid.UUID     `json:"user_id"`
	Roles  []models.Role `json:"roles"`
}

// AssignRolesRequest represents a request to assign roles to a user.
type AssignRolesRequest struct {
	UserID  uuid.UUID
	RoleIDs []uuid.UUID
}

// GetUserRoles retrieves all roles for a user.
func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user.Roles, nil
}

// AssignRoles assigns roles to a user.
func (s *UserRoleService) AssignRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) ([]models.Role, error) {
	// Get user
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get roles
	var roles []models.Role
	if err := s.db.WithContext(ctx).
		Where("id IN ?", roleIDs).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	if len(roles) != len(roleIDs) {
		return nil, ErrRoleNotFound
	}

	// Verify roles belong to user's tenant or are system roles
	for _, role := range roles {
		if role.TenantID != nil && *role.TenantID != user.TenantID {
			return nil, ErrTenantMismatch
		}
	}

	// Append roles
	if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Append(roles); err != nil {
		return nil, err
	}

	return s.GetUserRoles(ctx, userID)
}

// RemoveRoles removes roles from a user.
func (s *UserRoleService) RemoveRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) ([]models.Role, error) {
	// Get user
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get roles
	var roles []models.Role
	if err := s.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	// Remove roles
	if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Delete(roles); err != nil {
		return nil, err
	}

	return s.GetUserRoles(ctx, userID)
}

// SetRoles replaces all roles for a user.
func (s *UserRoleService) SetRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) ([]models.Role, error) {
	// Get user
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get roles
	var roles []models.Role
	if err := s.db.WithContext(ctx).
		Where("id IN ?", roleIDs).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	if len(roles) != len(roleIDs) {
		return nil, ErrRoleNotFound
	}

	// Verify roles belong to user's tenant or are system roles
	for _, role := range roles {
		if role.TenantID != nil && *role.TenantID != user.TenantID {
			return nil, ErrTenantMismatch
		}
	}

	// Replace all roles
	if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Replace(roles); err != nil {
		return nil, err
	}

	return s.GetUserRoles(ctx, userID)
}

// HasRole checks if a user has a specific role.
func (s *UserRoleService) HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

// HasPermission checks if a user has a specific permission through any of their roles.
func (s *UserRoleService) HasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		for _, perm := range role.Permissions {
			if perm.Code == permissionCode {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserPermissions retrieves all permissions for a user through their roles.
func (s *UserRoleService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]models.Permission, error) {
	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Use a map to deduplicate permissions
	permMap := make(map[uuid.UUID]models.Permission)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permMap[perm.ID] = perm
		}
	}

	// Convert map to slice
	permissions := make([]models.Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// GetUserPermissionCodes retrieves all permission codes for a user.
func (s *UserRoleService) GetUserPermissionCodes(ctx context.Context, userID uuid.UUID) ([]string, error) {
	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(permissions))
	for i, perm := range permissions {
		codes[i] = perm.Code
	}

	return codes, nil
}

// GetUsersWithRole retrieves all users with a specific role.
func (s *UserRoleService) GetUsersWithRole(ctx context.Context, roleID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Preload("Roles.Permissions").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
