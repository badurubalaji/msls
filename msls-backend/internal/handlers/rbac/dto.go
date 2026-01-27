// Package rbac provides HTTP handlers for role and permission management endpoints.
package rbac

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Request DTOs
// ============================================================================

// CreateRoleRequest represents the request body for creating a role.
type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=100"`
	Description   string   `json:"description" binding:"max=500"`
	PermissionIDs []string `json:"permission_ids" binding:"dive,uuid"`
}

// UpdateRoleRequest represents the request body for updating a role.
type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

// AssignPermissionsRequest represents the request body for assigning permissions.
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1,dive,uuid"`
}

// RemovePermissionsRequest represents the request body for removing permissions.
type RemovePermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1,dive,uuid"`
}

// AssignRolesRequest represents the request body for assigning roles to a user.
type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required,min=1,dive,uuid"`
}

// RemoveRolesRequest represents the request body for removing roles from a user.
type RemoveRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required,min=1,dive,uuid"`
}

// CreatePermissionRequest represents the request body for creating a permission.
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required,min=3,max=100"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Module      string `json:"module" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// UpdatePermissionRequest represents the request body for updating a permission.
type UpdatePermissionRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// RoleDTO represents a role in API responses.
type RoleDTO struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    *uuid.UUID      `json:"tenant_id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	IsSystem    bool            `json:"is_system"`
	Permissions []PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// PermissionDTO represents a permission in API responses.
type PermissionDTO struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Module      string    `json:"module"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRolesDTO represents a user's roles in API responses.
type UserRolesDTO struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []RoleDTO `json:"roles"`
}

// ModulesDTO represents the list of permission modules.
type ModulesDTO struct {
	Modules []string `json:"modules"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// roleToDTO converts a Role model to a RoleDTO.
func roleToDTO(role *models.Role) RoleDTO {
	dto := RoleDTO{
		ID:          role.ID,
		TenantID:    role.TenantID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	if len(role.Permissions) > 0 {
		dto.Permissions = make([]PermissionDTO, len(role.Permissions))
		for i, perm := range role.Permissions {
			dto.Permissions[i] = permissionToDTO(&perm)
		}
	}

	return dto
}

// rolesToDTOs converts a slice of Role models to RoleDTOs.
func rolesToDTOs(roles []models.Role) []RoleDTO {
	dtos := make([]RoleDTO, len(roles))
	for i, role := range roles {
		dtos[i] = roleToDTO(&role)
	}
	return dtos
}

// permissionToDTO converts a Permission model to a PermissionDTO.
func permissionToDTO(perm *models.Permission) PermissionDTO {
	return PermissionDTO{
		ID:          perm.ID,
		Code:        perm.Code,
		Name:        perm.Name,
		Module:      perm.Module,
		Description: perm.Description,
		CreatedAt:   perm.CreatedAt,
	}
}

// permissionsToDTOs converts a slice of Permission models to PermissionDTOs.
func permissionsToDTOs(perms []models.Permission) []PermissionDTO {
	dtos := make([]PermissionDTO, len(perms))
	for i, perm := range perms {
		dtos[i] = permissionToDTO(&perm)
	}
	return dtos
}

// parseUUIDs converts a slice of string UUIDs to uuid.UUIDs.
func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids[i] = parsed
	}
	return uuids, nil
}
