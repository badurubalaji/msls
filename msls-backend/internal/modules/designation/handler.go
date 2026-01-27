// Package designation provides designation management functionality.
package designation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// Handler handles designation-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new designation handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateDesignationRequest represents the request body for creating a designation.
type CreateDesignationRequest struct {
	Name         string  `json:"name" binding:"required,max=100"`
	Level        int     `json:"level" binding:"required,min=1,max=10"`
	DepartmentID *string `json:"departmentId" binding:"omitempty,uuid"`
	IsActive     *bool   `json:"isActive"`
}

// UpdateDesignationRequest represents the request body for updating a designation.
type UpdateDesignationRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Level        *int    `json:"level" binding:"omitempty,min=1,max=10"`
	DepartmentID *string `json:"departmentId" binding:"omitempty,uuid"`
	IsActive     *bool   `json:"isActive"`
}

// List returns all designations for the tenant.
// @Summary List designations
// @Description Get all designations for the current tenant
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param department_id query string false "Filter by department ID"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name"
// @Success 200 {object} response.Success{data=DesignationListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/designations [get]
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ListFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	// Parse department_id filter
	if departmentIDStr := c.Query("department_id"); departmentIDStr != "" {
		departmentID, err := uuid.Parse(departmentIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
			return
		}
		filter.DepartmentID = &departmentID
	}

	// Parse is_active filter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	designations, total, staffCounts, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list designations"))
		return
	}

	resp := DesignationListResponse{
		Designations: ToDesignationResponses(designations, staffCounts),
		Total:        total,
	}

	response.OK(c, resp)
}

// Get returns a designation by ID.
// @Summary Get designation
// @Description Get a designation by ID
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Designation ID"
// @Success 200 {object} response.Success{data=DesignationResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/designations/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid designation ID"))
		return
	}

	designation, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrDesignationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Designation not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get designation"))
		return
	}

	// Get staff count for this designation
	staffCounts, _ := h.service.repo.GetStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToDesignationResponse(designation, staffCount))
}

// Create creates a new designation.
// @Summary Create designation
// @Description Create a new designation
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param body body CreateDesignationRequest true "Designation data"
// @Success 201 {object} response.Success{data=DesignationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/designations [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateDesignationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreateDesignationDTO{
		TenantID: tenantID,
		Name:     req.Name,
		Level:    req.Level,
		IsActive: true,
	}

	if req.IsActive != nil {
		dto.IsActive = *req.IsActive
	}

	if req.DepartmentID != nil {
		departmentID, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
			return
		}
		dto.DepartmentID = &departmentID
	}

	designation, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrDuplicateName) {
			apperrors.Abort(c, apperrors.Conflict("Designation name already exists"))
			return
		}
		if errors.Is(err, ErrInvalidLevel) {
			apperrors.Abort(c, apperrors.BadRequest("Level must be between 1 and 10"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create designation"))
		return
	}

	response.Created(c, ToDesignationResponse(designation, 0))
}

// Update updates a designation.
// @Summary Update designation
// @Description Update a designation
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Designation ID"
// @Param body body UpdateDesignationRequest true "Designation data"
// @Success 200 {object} response.Success{data=DesignationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/designations/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid designation ID"))
		return
	}

	var req UpdateDesignationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateDesignationDTO{
		Name:     req.Name,
		Level:    req.Level,
		IsActive: req.IsActive,
	}

	if req.DepartmentID != nil {
		departmentID, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
			return
		}
		dto.DepartmentID = &departmentID
	}

	designation, err := h.service.Update(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrDesignationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Designation not found"))
			return
		}
		if errors.Is(err, ErrDuplicateName) {
			apperrors.Abort(c, apperrors.Conflict("Designation name already exists"))
			return
		}
		if errors.Is(err, ErrInvalidLevel) {
			apperrors.Abort(c, apperrors.BadRequest("Level must be between 1 and 10"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update designation"))
		return
	}

	// Get staff count for this designation
	staffCounts, _ := h.service.repo.GetStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToDesignationResponse(designation, staffCount))
}

// Delete deletes a designation.
// @Summary Delete designation
// @Description Delete a designation
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Designation ID"
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/designations/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid designation ID"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrDesignationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Designation not found"))
			return
		}
		if errors.Is(err, ErrDesignationInUse) {
			apperrors.Abort(c, apperrors.Conflict("Designation is in use and cannot be deleted"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete designation"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDropdown returns active designations for dropdown/select.
// @Summary Get designations for dropdown
// @Description Get active designations for dropdown/select
// @Tags Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param department_id query string false "Filter by department ID"
// @Success 200 {object} response.Success{data=[]DropdownItem}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/designations/dropdown [get]
func (h *Handler) GetDropdown(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var departmentID *uuid.UUID
	if departmentIDStr := c.Query("department_id"); departmentIDStr != "" {
		id, err := uuid.Parse(departmentIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
			return
		}
		departmentID = &id
	}

	designations, err := h.service.GetActiveForDropdown(c.Request.Context(), tenantID, departmentID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get designations"))
		return
	}

	response.OK(c, ToDropdownItems(designations))
}

// RegisterRoutes registers designation routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	designations := rg.Group("/designations")
	designations.Use(authMiddleware)
	{
		designations.GET("", h.List)
		designations.GET("/dropdown", h.GetDropdown)
		designations.GET("/:id", h.Get)
		designations.POST("", h.Create)
		designations.PUT("/:id", h.Update)
		designations.DELETE("/:id", h.Delete)
	}
}
