// Package salary provides salary management functionality.
package salary

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles salary-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new salary handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// Salary Component Handlers
// ========================================

// CreateComponentRequest represents the request body for creating a component.
type CreateComponentRequest struct {
	Name            string  `json:"name" binding:"required,max=100"`
	Code            string  `json:"code" binding:"required,max=20"`
	Description     *string `json:"description"`
	ComponentType   string  `json:"componentType" binding:"required,oneof=earning deduction"`
	CalculationType string  `json:"calculationType" binding:"required,oneof=fixed percentage"`
	PercentageOfID  *string `json:"percentageOfId" binding:"omitempty,uuid"`
	IsTaxable       *bool   `json:"isTaxable"`
	DisplayOrder    *int    `json:"displayOrder"`
}

// UpdateComponentRequest represents the request body for updating a component.
type UpdateComponentRequest struct {
	Name            *string `json:"name" binding:"omitempty,max=100"`
	Code            *string `json:"code" binding:"omitempty,max=20"`
	Description     *string `json:"description"`
	ComponentType   *string `json:"componentType" binding:"omitempty,oneof=earning deduction"`
	CalculationType *string `json:"calculationType" binding:"omitempty,oneof=fixed percentage"`
	PercentageOfID  *string `json:"percentageOfId" binding:"omitempty,uuid"`
	IsTaxable       *bool   `json:"isTaxable"`
	IsActive        *bool   `json:"isActive"`
	DisplayOrder    *int    `json:"displayOrder"`
}

// ListComponents returns all salary components.
func (h *Handler) ListComponents(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ComponentFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	if typeStr := c.Query("type"); typeStr != "" {
		compType := models.ComponentType(typeStr)
		if compType.IsValid() {
			filter.ComponentType = &compType
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	components, total, err := h.service.ListComponents(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list salary components"))
		return
	}

	response.OK(c, ComponentListResponse{
		Components: ToComponentResponses(components),
		Total:      total,
	})
}

// GetComponent returns a salary component by ID.
func (h *Handler) GetComponent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid component ID"))
		return
	}

	component, err := h.service.GetComponentByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrComponentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary component not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get salary component"))
		return
	}

	response.OK(c, ToComponentResponse(component))
}

// CreateComponent creates a new salary component.
func (h *Handler) CreateComponent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreateComponentDTO{
		TenantID:        tenantID,
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		ComponentType:   models.ComponentType(req.ComponentType),
		CalculationType: models.CalculationType(req.CalculationType),
		IsTaxable:       true,
		DisplayOrder:    0,
	}

	if req.IsTaxable != nil {
		dto.IsTaxable = *req.IsTaxable
	}

	if req.DisplayOrder != nil {
		dto.DisplayOrder = *req.DisplayOrder
	}

	if req.PercentageOfID != nil {
		id, err := uuid.Parse(*req.PercentageOfID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid percentageOfId"))
			return
		}
		dto.PercentageOfID = &id
	}

	component, err := h.service.CreateComponent(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrDuplicateCode) {
			apperrors.Abort(c, apperrors.Conflict("Component code already exists"))
			return
		}
		if errors.Is(err, ErrPercentageOfRequired) {
			apperrors.Abort(c, apperrors.BadRequest("percentageOfId is required for percentage-based components"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create salary component"))
		return
	}

	response.Created(c, ToComponentResponse(component))
}

// UpdateComponent updates a salary component.
func (h *Handler) UpdateComponent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid component ID"))
		return
	}

	var req UpdateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateComponentDTO{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		IsTaxable:    req.IsTaxable,
		IsActive:     req.IsActive,
		DisplayOrder: req.DisplayOrder,
	}

	if req.ComponentType != nil {
		compType := models.ComponentType(*req.ComponentType)
		dto.ComponentType = &compType
	}

	if req.CalculationType != nil {
		calcType := models.CalculationType(*req.CalculationType)
		dto.CalculationType = &calcType
	}

	if req.PercentageOfID != nil {
		pid, err := uuid.Parse(*req.PercentageOfID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid percentageOfId"))
			return
		}
		dto.PercentageOfID = &pid
	}

	component, err := h.service.UpdateComponent(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrComponentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary component not found"))
			return
		}
		if errors.Is(err, ErrDuplicateCode) {
			apperrors.Abort(c, apperrors.Conflict("Component code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update salary component"))
		return
	}

	response.OK(c, ToComponentResponse(component))
}

// DeleteComponent deletes a salary component.
func (h *Handler) DeleteComponent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid component ID"))
		return
	}

	if err := h.service.DeleteComponent(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrComponentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary component not found"))
			return
		}
		if errors.Is(err, ErrComponentInUse) {
			apperrors.Abort(c, apperrors.Conflict("Salary component is in use"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete salary component"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetComponentsDropdown returns active components for dropdown.
func (h *Handler) GetComponentsDropdown(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var compType *models.ComponentType
	if typeStr := c.Query("type"); typeStr != "" {
		ct := models.ComponentType(typeStr)
		if ct.IsValid() {
			compType = &ct
		}
	}

	components, err := h.service.GetActiveComponents(c.Request.Context(), tenantID, compType)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get salary components"))
		return
	}

	response.OK(c, ToComponentDropdownItems(components))
}

// ========================================
// Salary Structure Handlers
// ========================================

// CreateStructureRequest represents the request body for creating a structure.
type CreateStructureRequest struct {
	Name          string                     `json:"name" binding:"required,max=100"`
	Code          string                     `json:"code" binding:"required,max=20"`
	Description   *string                    `json:"description"`
	DesignationID *string                    `json:"designationId" binding:"omitempty,uuid"`
	Components    []StructureComponentInput  `json:"components" binding:"required,min=1,dive"`
}

// StructureComponentInput represents a component in a structure request.
type StructureComponentInput struct {
	ComponentID string  `json:"componentId" binding:"required,uuid"`
	Amount      *string `json:"amount"`
	Percentage  *string `json:"percentage"`
}

// UpdateStructureRequest represents the request body for updating a structure.
type UpdateStructureRequest struct {
	Name          *string                    `json:"name" binding:"omitempty,max=100"`
	Code          *string                    `json:"code" binding:"omitempty,max=20"`
	Description   *string                    `json:"description"`
	DesignationID *string                    `json:"designationId" binding:"omitempty,uuid"`
	IsActive      *bool                      `json:"isActive"`
	Components    []StructureComponentInput  `json:"components" binding:"omitempty,dive"`
}

// ListStructures returns all salary structures.
func (h *Handler) ListStructures(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := StructureFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	if desigIDStr := c.Query("designation_id"); desigIDStr != "" {
		desigID, err := uuid.Parse(desigIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid designation ID"))
			return
		}
		filter.DesignationID = &desigID
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	structures, total, staffCounts, err := h.service.ListStructures(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list salary structures"))
		return
	}

	response.OK(c, StructureListResponse{
		Structures: ToStructureResponses(structures, staffCounts),
		Total:      total,
	})
}

// GetStructure returns a salary structure by ID.
func (h *Handler) GetStructure(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid structure ID"))
		return
	}

	structure, err := h.service.GetStructureByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStructureNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary structure not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get salary structure"))
		return
	}

	// Get staff count
	staffCounts, _ := h.service.repo.GetStructureStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToStructureResponse(structure, staffCount, true))
}

// CreateStructure creates a new salary structure.
func (h *Handler) CreateStructure(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreateStructureDTO{
		TenantID:    tenantID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if req.DesignationID != nil {
		desigID, err := uuid.Parse(*req.DesignationID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid designationId"))
			return
		}
		dto.DesignationID = &desigID
	}

	for _, comp := range req.Components {
		compID, _ := uuid.Parse(comp.ComponentID)
		scDTO := StructureComponentDTO{ComponentID: compID}

		if comp.Amount != nil {
			amt, err := decimal.NewFromString(*comp.Amount)
			if err != nil {
				apperrors.Abort(c, apperrors.BadRequest("Invalid amount for component"))
				return
			}
			scDTO.Amount = &amt
		}

		if comp.Percentage != nil {
			pct, err := decimal.NewFromString(*comp.Percentage)
			if err != nil {
				apperrors.Abort(c, apperrors.BadRequest("Invalid percentage for component"))
				return
			}
			scDTO.Percentage = &pct
		}

		dto.Components = append(dto.Components, scDTO)
	}

	structure, err := h.service.CreateStructure(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrDuplicateStructureCode) {
			apperrors.Abort(c, apperrors.Conflict("Structure code already exists"))
			return
		}
		if errors.Is(err, ErrNoComponentsInStructure) {
			apperrors.Abort(c, apperrors.BadRequest("Structure must have at least one component"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create salary structure"))
		return
	}

	response.Created(c, ToStructureResponse(structure, 0, true))
}

// UpdateStructure updates a salary structure.
func (h *Handler) UpdateStructure(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid structure ID"))
		return
	}

	var req UpdateStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateStructureDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if req.DesignationID != nil {
		desigID, err := uuid.Parse(*req.DesignationID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid designationId"))
			return
		}
		dto.DesignationID = &desigID
	}

	if req.Components != nil {
		for _, comp := range req.Components {
			compID, _ := uuid.Parse(comp.ComponentID)
			scDTO := StructureComponentDTO{ComponentID: compID}

			if comp.Amount != nil {
				amt, err := decimal.NewFromString(*comp.Amount)
				if err != nil {
					apperrors.Abort(c, apperrors.BadRequest("Invalid amount for component"))
					return
				}
				scDTO.Amount = &amt
			}

			if comp.Percentage != nil {
				pct, err := decimal.NewFromString(*comp.Percentage)
				if err != nil {
					apperrors.Abort(c, apperrors.BadRequest("Invalid percentage for component"))
					return
				}
				scDTO.Percentage = &pct
			}

			dto.Components = append(dto.Components, scDTO)
		}
	}

	structure, err := h.service.UpdateStructure(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrStructureNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary structure not found"))
			return
		}
		if errors.Is(err, ErrDuplicateStructureCode) {
			apperrors.Abort(c, apperrors.Conflict("Structure code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update salary structure"))
		return
	}

	// Get staff count
	staffCounts, _ := h.service.repo.GetStructureStaffCounts(c.Request.Context(), tenantID, []uuid.UUID{id})
	staffCount := 0
	if staffCounts != nil {
		staffCount = staffCounts[id]
	}

	response.OK(c, ToStructureResponse(structure, staffCount, true))
}

// DeleteStructure deletes a salary structure.
func (h *Handler) DeleteStructure(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid structure ID"))
		return
	}

	if err := h.service.DeleteStructure(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrStructureNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Salary structure not found"))
			return
		}
		if errors.Is(err, ErrStructureInUse) {
			apperrors.Abort(c, apperrors.Conflict("Salary structure is in use"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete salary structure"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetStructuresDropdown returns active structures for dropdown.
func (h *Handler) GetStructuresDropdown(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	structures, err := h.service.GetActiveStructures(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get salary structures"))
		return
	}

	response.OK(c, ToStructureDropdownItems(structures))
}

// ========================================
// Staff Salary Handlers
// ========================================

// AssignSalaryRequest represents the request body for assigning salary.
type AssignSalaryRequest struct {
	StructureID    *string                    `json:"structureId" binding:"omitempty,uuid"`
	EffectiveFrom  string                     `json:"effectiveFrom" binding:"required"`
	Components     []StaffComponentInput      `json:"components" binding:"required,min=1,dive"`
	RevisionReason *string                    `json:"revisionReason"`
}

// StaffComponentInput represents a component in staff salary request.
type StaffComponentInput struct {
	ComponentID  string `json:"componentId" binding:"required,uuid"`
	Amount       string `json:"amount" binding:"required"`
	IsOverridden bool   `json:"isOverridden"`
}

// GetStaffSalary returns the current salary for a staff member.
func (h *Handler) GetStaffSalary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	salary, err := h.service.GetCurrentStaffSalary(c.Request.Context(), tenantID, staffID)
	if err != nil {
		if errors.Is(err, ErrStaffSalaryNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Staff salary not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get staff salary"))
		return
	}

	response.OK(c, ToStaffSalaryResponse(salary, true))
}

// AssignStaffSalary assigns or revises salary for a staff member.
func (h *Handler) AssignStaffSalary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req AssignSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	effectiveFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid effectiveFrom date format (use YYYY-MM-DD)"))
		return
	}

	dto := AssignSalaryDTO{
		TenantID:       tenantID,
		StaffID:        staffID,
		EffectiveFrom:  effectiveFrom,
		RevisionReason: req.RevisionReason,
		CreatedBy:      &userID,
	}

	if req.StructureID != nil {
		structID, err := uuid.Parse(*req.StructureID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid structureId"))
			return
		}
		dto.StructureID = &structID
	}

	for _, comp := range req.Components {
		compID, _ := uuid.Parse(comp.ComponentID)
		amt, err := decimal.NewFromString(comp.Amount)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid amount for component"))
			return
		}

		dto.Components = append(dto.Components, StaffComponentDTO{
			ComponentID:  compID,
			Amount:       amt,
			IsOverridden: comp.IsOverridden,
		})
	}

	salary, err := h.service.AssignSalary(c.Request.Context(), dto)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to assign staff salary"))
		return
	}

	response.Created(c, ToStaffSalaryResponse(salary, true))
}

// GetStaffSalaryHistory returns salary history for a staff member.
func (h *Handler) GetStaffSalaryHistory(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	salaries, total, err := h.service.GetStaffSalaryHistory(c.Request.Context(), tenantID, staffID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get staff salary history"))
		return
	}

	response.OK(c, StaffSalaryHistoryResponse{
		History: ToStaffSalaryResponses(salaries),
		Total:   total,
	})
}

// RegisterRoutes registers salary routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Salary Components
	components := rg.Group("/salary-components")
	components.Use(authMiddleware)
	{
		components.GET("", h.ListComponents)
		components.GET("/dropdown", h.GetComponentsDropdown)
		components.GET("/:id", h.GetComponent)
		components.POST("", h.CreateComponent)
		components.PUT("/:id", h.UpdateComponent)
		components.DELETE("/:id", h.DeleteComponent)
	}

	// Salary Structures
	structures := rg.Group("/salary-structures")
	structures.Use(authMiddleware)
	{
		structures.GET("", h.ListStructures)
		structures.GET("/dropdown", h.GetStructuresDropdown)
		structures.GET("/:id", h.GetStructure)
		structures.POST("", h.CreateStructure)
		structures.PUT("/:id", h.UpdateStructure)
		structures.DELETE("/:id", h.DeleteStructure)
	}
}

// RegisterStaffSalaryRoutes registers staff salary routes.
func (h *Handler) RegisterStaffSalaryRoutes(staffGroup *gin.RouterGroup) {
	staffGroup.GET("/:id/salary", h.GetStaffSalary)
	staffGroup.POST("/:id/salary", h.AssignStaffSalary)
	staffGroup.GET("/:id/salary/history", h.GetStaffSalaryHistory)
}
