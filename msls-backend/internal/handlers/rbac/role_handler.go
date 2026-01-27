// Package rbac provides HTTP handlers for role and permission management endpoints.
package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	rbacservice "msls-backend/internal/services/rbac"
)

// RoleHandler handles role-related HTTP requests.
type RoleHandler struct {
	roleService *rbacservice.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleService *rbacservice.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// List returns all roles for the tenant.
// @Summary List roles
// @Description Get all roles for the current tenant
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name or description"
// @Param include_system query bool false "Include system roles" default(true)
// @Success 200 {object} response.Success{data=[]RoleDTO}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	search := c.Query("search")
	includeSystem := c.DefaultQuery("include_system", "true") == "true"

	filter := rbacservice.ListRolesFilter{
		TenantID:      &tenantID,
		IncludeSystem: includeSystem,
		Search:        search,
	}

	roles, err := h.roleService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve roles"))
		return
	}

	response.OK(c, rolesToDTOs(roles))
}

// GetByID returns a role by ID.
// @Summary Get role by ID
// @Description Get a role with its permissions
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Role ID" format(uuid)
// @Success 200 {object} response.Success{data=RoleDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID"))
		return
	}

	role, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve role"))
		}
		return
	}

	// Verify tenant access (system roles are accessible by all tenants)
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}
	if role.TenantID != nil && *role.TenantID != tenantID {
		apperrors.Abort(c, apperrors.NotFound("Role not found"))
		return
	}

	response.OK(c, roleToDTO(role))
}

// Create creates a new role.
// @Summary Create role
// @Description Create a new custom role for the tenant
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateRoleRequest true "Role details"
// @Success 201 {object} response.Success{data=RoleDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Parse permission IDs if provided
	var permissionIDs []uuid.UUID
	if len(req.PermissionIDs) > 0 {
		var err error
		permissionIDs, err = parseUUIDs(req.PermissionIDs)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid permission ID format"))
			return
		}
	}

	createReq := rbacservice.CreateRoleRequest{
		TenantID:      &tenantID,
		Name:          req.Name,
		Description:   req.Description,
		IsSystem:      false, // Custom roles are never system roles
		PermissionIDs: permissionIDs,
	}

	role, err := h.roleService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Role name is required"))
		case rbacservice.ErrRoleNameExists:
			apperrors.Abort(c, apperrors.Conflict("Role with this name already exists"))
		case rbacservice.ErrPermissionNotFound:
			apperrors.Abort(c, apperrors.NotFound("One or more permissions not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create role"))
		}
		return
	}

	response.Created(c, roleToDTO(role))
}

// Update updates a role.
// @Summary Update role
// @Description Update an existing role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Role ID" format(uuid)
// @Param request body UpdateRoleRequest true "Role updates"
// @Success 200 {object} response.Success{data=RoleDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID"))
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Verify tenant access
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Check role belongs to tenant
	existingRole, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == rbacservice.ErrRoleNotFound {
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve role"))
		return
	}
	if existingRole.TenantID != nil && *existingRole.TenantID != tenantID {
		apperrors.Abort(c, apperrors.NotFound("Role not found"))
		return
	}

	updateReq := rbacservice.UpdateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	role, err := h.roleService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
		case rbacservice.ErrRoleNameExists:
			apperrors.Abort(c, apperrors.Conflict("Role with this name already exists"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update role"))
		}
		return
	}

	response.OK(c, roleToDTO(role))
}

// Delete deletes a role.
// @Summary Delete role
// @Description Delete a custom role (system roles cannot be deleted)
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Role ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID"))
		return
	}

	// Verify tenant access
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Check role belongs to tenant
	existingRole, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == rbacservice.ErrRoleNotFound {
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve role"))
		return
	}
	if existingRole.TenantID != nil && *existingRole.TenantID != tenantID {
		apperrors.Abort(c, apperrors.NotFound("Role not found"))
		return
	}

	err = h.roleService.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
		case rbacservice.ErrCannotDeleteSystem:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete system role"))
		case rbacservice.ErrRoleInUse:
			apperrors.Abort(c, apperrors.Conflict("Role is assigned to users and cannot be deleted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete role"))
		}
		return
	}

	response.NoContent(c)
}

// AssignPermissions assigns permissions to a role.
// @Summary Assign permissions to role
// @Description Add permissions to a role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Role ID" format(uuid)
// @Param request body AssignPermissionsRequest true "Permissions to assign"
// @Success 200 {object} response.Success{data=RoleDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID"))
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	permissionIDs, err := parseUUIDs(req.PermissionIDs)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid permission ID format"))
		return
	}

	// Verify tenant access
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Check role belongs to tenant
	existingRole, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == rbacservice.ErrRoleNotFound {
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve role"))
		return
	}
	if existingRole.TenantID != nil && *existingRole.TenantID != tenantID {
		apperrors.Abort(c, apperrors.NotFound("Role not found"))
		return
	}

	role, err := h.roleService.AssignPermissions(c.Request.Context(), id, permissionIDs)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
		case rbacservice.ErrPermissionNotFound:
			apperrors.Abort(c, apperrors.NotFound("One or more permissions not found"))
		case rbacservice.ErrCannotModifySystem:
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify system role permissions"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to assign permissions"))
		}
		return
	}

	response.OK(c, roleToDTO(role))
}

// RemovePermissions removes permissions from a role.
// @Summary Remove permissions from role
// @Description Remove permissions from a role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Role ID" format(uuid)
// @Param request body RemovePermissionsRequest true "Permissions to remove"
// @Success 200 {object} response.Success{data=RoleDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/roles/{id}/permissions [delete]
func (h *RoleHandler) RemovePermissions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID"))
		return
	}

	var req RemovePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	permissionIDs, err := parseUUIDs(req.PermissionIDs)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid permission ID format"))
		return
	}

	// Verify tenant access
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Check role belongs to tenant
	existingRole, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == rbacservice.ErrRoleNotFound {
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve role"))
		return
	}
	if existingRole.TenantID != nil && *existingRole.TenantID != tenantID {
		apperrors.Abort(c, apperrors.NotFound("Role not found"))
		return
	}

	role, err := h.roleService.RemovePermissions(c.Request.Context(), id, permissionIDs)
	if err != nil {
		switch err {
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("Role not found"))
		case rbacservice.ErrCannotModifySystem:
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify system role permissions"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to remove permissions"))
		}
		return
	}

	response.OK(c, roleToDTO(role))
}
