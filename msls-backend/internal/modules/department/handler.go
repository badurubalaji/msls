// Package department provides department management functionality.
package department

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// Handler handles department-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new department handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateDepartmentRequest represents the request body for creating a department.
type CreateDepartmentRequest struct {
	BranchID    string  `json:"branchId" binding:"required,uuid"`
	Name        string  `json:"name" binding:"required,max=100"`
	Code        string  `json:"code" binding:"required,max=20"`
	Description *string `json:"description"`
	HeadID      *string `json:"headId" binding:"omitempty,uuid"`
	IsActive    *bool   `json:"isActive"`
}

// UpdateDepartmentRequest represents the request body for updating a department.
type UpdateDepartmentRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Code        *string `json:"code" binding:"omitempty,max=20"`
	Description *string `json:"description"`
	HeadID      *string `json:"headId" binding:"omitempty,uuid"`
	IsActive    *bool   `json:"isActive"`
}

// List returns all departments for the tenant.
// @Summary List departments
// @Description Get all departments for the current tenant
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or code"
// @Success 200 {object} response.Success{data=DepartmentListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/departments [get]
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

	// Parse branch_id filter
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = &branchID
	}

	// Parse is_active filter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	departments, total, staffCounts, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list departments"))
		return
	}

	resp := DepartmentListResponse{
		Departments: ToDepartmentResponses(departments, staffCounts),
		Total:       total,
	}

	response.OK(c, resp)
}

// Get returns a department by ID.
// @Summary Get department
// @Description Get a department by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Department ID"
// @Success 200 {object} response.Success{data=DepartmentResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/departments/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
		return
	}

	department, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrDepartmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Department not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get department"))
		return
	}

	// Get staff count for this department
	staffCounts, _ := h.service.repo.GetStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToDepartmentResponse(department, staffCount))
}

// Create creates a new department.
// @Summary Create department
// @Description Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param body body CreateDepartmentRequest true "Department data"
// @Success 201 {object} response.Success{data=DepartmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/departments [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	branchID, _ := uuid.Parse(req.BranchID)

	dto := CreateDepartmentDTO{
		TenantID:    tenantID,
		BranchID:    branchID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    true,
	}

	if req.IsActive != nil {
		dto.IsActive = *req.IsActive
	}

	if req.HeadID != nil {
		headID, err := uuid.Parse(*req.HeadID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid head ID"))
			return
		}
		dto.HeadID = &headID
	}

	department, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrDuplicateCode) {
			apperrors.Abort(c, apperrors.Conflict("Department code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create department"))
		return
	}

	response.Created(c, ToDepartmentResponse(department, 0))
}

// Update updates a department.
// @Summary Update department
// @Description Update a department
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Department ID"
// @Param body body UpdateDepartmentRequest true "Department data"
// @Success 200 {object} response.Success{data=DepartmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/departments/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
		return
	}

	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateDepartmentDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if req.HeadID != nil {
		headID, err := uuid.Parse(*req.HeadID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid head ID"))
			return
		}
		dto.HeadID = &headID
	}

	department, err := h.service.Update(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrDepartmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Department not found"))
			return
		}
		if errors.Is(err, ErrDuplicateCode) {
			apperrors.Abort(c, apperrors.Conflict("Department code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update department"))
		return
	}

	// Get staff count for this department
	staffCounts, _ := h.service.repo.GetStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToDepartmentResponse(department, staffCount))
}

// Delete deletes a department.
// @Summary Delete department
// @Description Delete a department
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Department ID"
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/departments/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrDepartmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Department not found"))
			return
		}
		if errors.Is(err, ErrDepartmentInUse) {
			apperrors.Abort(c, apperrors.Conflict("Department is in use and cannot be deleted"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete department"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDropdown returns active departments for dropdown/select.
// @Summary Get departments for dropdown
// @Description Get active departments for dropdown/select
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} response.Success{data=[]DropdownItem}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/departments/dropdown [get]
func (h *Handler) GetDropdown(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		id, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &id
	}

	departments, err := h.service.GetActiveForDropdown(c.Request.Context(), tenantID, branchID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get departments"))
		return
	}

	response.OK(c, ToDropdownItems(departments))
}

// RegisterRoutes registers department routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	departments := rg.Group("/departments")
	departments.Use(authMiddleware)
	{
		departments.GET("", h.List)
		departments.GET("/dropdown", h.GetDropdown)
		departments.GET("/:id", h.Get)
		departments.POST("", h.Create)
		departments.PUT("/:id", h.Update)
		departments.DELETE("/:id", h.Delete)
	}
}
