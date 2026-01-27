// Package rbac provides HTTP handlers for role and permission management endpoints.
package rbac

import (
	"github.com/gin-gonic/gin"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	rbacservice "msls-backend/internal/services/rbac"
)

// PermissionHandler handles permission-related HTTP requests.
type PermissionHandler struct {
	permissionService *rbacservice.PermissionService
}

// NewPermissionHandler creates a new PermissionHandler.
func NewPermissionHandler(permissionService *rbacservice.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

// List returns all permissions.
// @Summary List permissions
// @Description Get all permissions with optional filtering
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param module query string false "Filter by module"
// @Param search query string false "Search by code or name"
// @Success 200 {object} response.Success{data=[]PermissionDTO}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) List(c *gin.Context) {
	module := c.Query("module")
	search := c.Query("search")

	filter := rbacservice.ListPermissionsFilter{
		Module: module,
		Search: search,
	}

	permissions, err := h.permissionService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve permissions"))
		return
	}

	response.OK(c, permissionsToDTOs(permissions))
}

// GetModules returns all permission modules.
// @Summary List permission modules
// @Description Get all distinct permission modules/categories
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} response.Success{data=ModulesDTO}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/permissions/modules [get]
func (h *PermissionHandler) GetModules(c *gin.Context) {
	modules, err := h.permissionService.GetModules(c.Request.Context())
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve modules"))
		return
	}

	response.OK(c, ModulesDTO{Modules: modules})
}

// GetByModule returns all permissions for a specific module.
// @Summary List permissions by module
// @Description Get all permissions for a specific module
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param module path string true "Module name"
// @Success 200 {object} response.Success{data=[]PermissionDTO}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/permissions/modules/{module} [get]
func (h *PermissionHandler) GetByModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		apperrors.Abort(c, apperrors.BadRequest("Module is required"))
		return
	}

	permissions, err := h.permissionService.ListByModule(c.Request.Context(), module)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve permissions"))
		return
	}

	response.OK(c, permissionsToDTOs(permissions))
}
