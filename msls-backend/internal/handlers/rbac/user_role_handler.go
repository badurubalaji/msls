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

// UserRoleHandler handles user-role assignment HTTP requests.
type UserRoleHandler struct {
	userRoleService *rbacservice.UserRoleService
}

// NewUserRoleHandler creates a new UserRoleHandler.
func NewUserRoleHandler(userRoleService *rbacservice.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{userRoleService: userRoleService}
}

// GetUserRoles returns roles for a user.
// @Summary Get user roles
// @Description Get all roles assigned to a user
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} response.Success{data=UserRolesDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/users/{id}/roles [get]
func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	roles, err := h.userRoleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case rbacservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve user roles"))
		}
		return
	}

	response.OK(c, UserRolesDTO{
		UserID: userID,
		Roles:  rolesToDTOs(roles),
	})
}

// AssignRoles assigns roles to a user.
// @Summary Assign roles to user
// @Description Add roles to a user
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID" format(uuid)
// @Param request body AssignRolesRequest true "Roles to assign"
// @Success 200 {object} response.Success{data=UserRolesDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/users/{id}/roles [post]
func (h *UserRoleHandler) AssignRoles(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	roleIDs, err := parseUUIDs(req.RoleIDs)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID format"))
		return
	}

	roles, err := h.userRoleService.AssignRoles(c.Request.Context(), userID, roleIDs)
	if err != nil {
		switch err {
		case rbacservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		case rbacservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("One or more roles not found"))
		case rbacservice.ErrTenantMismatch:
			apperrors.Abort(c, apperrors.BadRequest("Role does not belong to user's tenant"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to assign roles"))
		}
		return
	}

	response.OK(c, UserRolesDTO{
		UserID: userID,
		Roles:  rolesToDTOs(roles),
	})
}

// RemoveRoles removes roles from a user.
// @Summary Remove roles from user
// @Description Remove roles from a user
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID" format(uuid)
// @Param request body RemoveRolesRequest true "Roles to remove"
// @Success 200 {object} response.Success{data=UserRolesDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/users/{id}/roles [delete]
func (h *UserRoleHandler) RemoveRoles(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	var req RemoveRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	roleIDs, err := parseUUIDs(req.RoleIDs)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid role ID format"))
		return
	}

	roles, err := h.userRoleService.RemoveRoles(c.Request.Context(), userID, roleIDs)
	if err != nil {
		switch err {
		case rbacservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to remove roles"))
		}
		return
	}

	response.OK(c, UserRolesDTO{
		UserID: userID,
		Roles:  rolesToDTOs(roles),
	})
}

// GetMyRoles returns roles for the current user.
// @Summary Get my roles
// @Description Get all roles assigned to the authenticated user
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} response.Success{data=UserRolesDTO}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/users/me/roles [get]
func (h *UserRoleHandler) GetMyRoles(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	roles, err := h.userRoleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case rbacservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve user roles"))
		}
		return
	}

	response.OK(c, UserRolesDTO{
		UserID: userID,
		Roles:  rolesToDTOs(roles),
	})
}
