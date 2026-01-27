// Package branch provides HTTP handlers for branch management endpoints.
package branch

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	branchservice "msls-backend/internal/services/branch"
)

// Handler handles branch-related HTTP requests.
type Handler struct {
	branchService *branchservice.Service
}

// NewHandler creates a new branch Handler.
func NewHandler(branchService *branchservice.Service) *Handler {
	return &Handler{branchService: branchService}
}

// List returns all branches for the tenant.
// @Summary List branches
// @Description Get all branches for the current tenant
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name, code, or city"
// @Param status query string false "Filter by status (active, inactive)"
// @Param is_primary query bool false "Filter by primary status"
// @Success 200 {object} response.Success{data=BranchListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/branches [get]
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := branchservice.ListFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.Status(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	// Parse is_primary filter
	if isPrimaryStr := c.Query("is_primary"); isPrimaryStr != "" {
		isPrimary := isPrimaryStr == "true"
		filter.IsPrimary = &isPrimary
	}

	branches, err := h.branchService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve branches"))
		return
	}

	total, err := h.branchService.Count(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to count branches"))
		return
	}

	resp := BranchListResponse{
		Branches: branchesToResponses(branches),
		Total:    total,
	}

	response.OK(c, resp)
}

// GetByID returns a branch by ID.
// @Summary Get branch by ID
// @Description Get a branch with full details
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Branch ID" format(uuid)
// @Success 200 {object} response.Success{data=BranchResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/branches/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	branch, err := h.branchService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNotFound:
			apperrors.Abort(c, apperrors.NotFound("Branch not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve branch"))
		}
		return
	}

	response.OK(c, branchToResponse(branch))
}

// Create creates a new branch.
// @Summary Create branch
// @Description Create a new branch for the tenant
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateBranchRequest true "Branch details"
// @Success 201 {object} response.Success{data=BranchResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/branches [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	createReq := branchservice.CreateRequest{
		TenantID:     tenantID,
		Code:         req.Code,
		Name:         req.Name,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		Phone:        req.Phone,
		Email:        req.Email,
		LogoURL:      req.LogoURL,
		Timezone:     req.Timezone,
		IsPrimary:    req.IsPrimary,
		Settings:     req.Settings,
		CreatedBy:    &userID,
	}

	branch, err := h.branchService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Branch name is required"))
		case branchservice.ErrBranchCodeRequired:
			apperrors.Abort(c, apperrors.BadRequest("Branch code is required"))
		case branchservice.ErrBranchCodeExists:
			apperrors.Abort(c, apperrors.Conflict("Branch with this code already exists"))
		case branchservice.ErrInvalidTimezone:
			apperrors.Abort(c, apperrors.BadRequest("Invalid timezone"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create branch"))
		}
		return
	}

	response.Created(c, branchToResponse(branch))
}

// Update updates a branch.
// @Summary Update branch
// @Description Update an existing branch
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Branch ID" format(uuid)
// @Param request body UpdateBranchRequest true "Branch updates"
// @Success 200 {object} response.Success{data=BranchResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/branches/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	var req UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := branchservice.UpdateRequest{
		Name:         req.Name,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		Phone:        req.Phone,
		Email:        req.Email,
		LogoURL:      req.LogoURL,
		Timezone:     req.Timezone,
		IsPrimary:    req.IsPrimary,
		Settings:     req.Settings,
		UpdatedBy:    &userID,
	}

	branch, err := h.branchService.Update(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNotFound:
			apperrors.Abort(c, apperrors.NotFound("Branch not found"))
		case branchservice.ErrInvalidTimezone:
			apperrors.Abort(c, apperrors.BadRequest("Invalid timezone"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update branch"))
		}
		return
	}

	response.OK(c, branchToResponse(branch))
}

// SetPrimary sets a branch as the primary branch.
// @Summary Set branch as primary
// @Description Set a branch as the primary branch for the tenant
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Branch ID" format(uuid)
// @Success 200 {object} response.Success{data=BranchResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/branches/{id}/primary [patch]
func (h *Handler) SetPrimary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	branch, err := h.branchService.SetPrimary(c.Request.Context(), tenantID, id, &userID)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNotFound:
			apperrors.Abort(c, apperrors.NotFound("Branch not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to set primary branch"))
		}
		return
	}

	response.OK(c, branchToResponse(branch))
}

// SetStatus sets the status of a branch.
// @Summary Set branch status
// @Description Activate or deactivate a branch
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Branch ID" format(uuid)
// @Param request body SetStatusRequest true "Status update"
// @Success 200 {object} response.Success{data=BranchResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/branches/{id}/status [patch]
func (h *Handler) SetStatus(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	var req SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	status := models.StatusInactive
	if req.IsActive {
		status = models.StatusActive
	}

	branch, err := h.branchService.SetStatus(c.Request.Context(), tenantID, id, status, &userID)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNotFound:
			apperrors.Abort(c, apperrors.NotFound("Branch not found"))
		case branchservice.ErrCannotDeactivatePrimaryBranch:
			apperrors.Abort(c, apperrors.BadRequest("Cannot deactivate primary branch"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update branch status"))
		}
		return
	}

	response.OK(c, branchToResponse(branch))
}

// Delete deletes a branch.
// @Summary Delete branch
// @Description Delete a branch (cannot delete primary branch)
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Branch ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/branches/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	err = h.branchService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case branchservice.ErrBranchNotFound:
			apperrors.Abort(c, apperrors.NotFound("Branch not found"))
		case branchservice.ErrCannotDeletePrimaryBranch:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete primary branch"))
		case branchservice.ErrBranchHasDependencies:
			apperrors.Abort(c, apperrors.Conflict("Branch has associated records and cannot be deleted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete branch"))
		}
		return
	}

	response.NoContent(c)
}
