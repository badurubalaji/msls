// Package exam provides examination management functionality.
package exam

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// RegisterRoutes registers all exam routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	examTypes := rg.Group("/exam-types")
	{
		examTypes.GET("", h.ListExamTypes)
		examTypes.GET("/:id", h.GetExamType)
		examTypes.POST("", h.CreateExamType)
		examTypes.PUT("/:id", h.UpdateExamType)
		examTypes.DELETE("/:id", h.DeleteExamType)
		examTypes.PATCH("/:id/active", h.ToggleExamTypeActive)
		examTypes.PUT("/order", h.UpdateDisplayOrders)
	}
}

// Handler handles exam-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new exam handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// Exam Type Handlers
// ========================================

// ListExamTypes returns all exam types for the tenant.
// @Summary List exam types
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or code"
// @Success 200 {object} response.Response{data=ExamTypeListResponse}
// @Router /api/v1/exam-types [get]
func (h *Handler) ListExamTypes(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ExamTypeFilter{TenantID: tenantID}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	filter.Search = c.Query("search")

	examTypes, total, err := h.service.ListExamTypes(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list exam types"))
		return
	}

	response.OK(c, ToExamTypeListResponse(examTypes, total))
}

// GetExamType returns a single exam type by ID.
// @Summary Get exam type
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param id path string true "Exam Type ID"
// @Success 200 {object} response.Response{data=ExamTypeResponse}
// @Router /api/v1/exam-types/{id} [get]
func (h *Handler) GetExamType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid exam type ID"))
		return
	}

	examType, err := h.service.GetExamType(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrExamTypeNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Exam type not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get exam type"))
		return
	}

	response.OK(c, ToExamTypeResponse(examType))
}

// CreateExamType creates a new exam type.
// @Summary Create exam type
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param request body CreateExamTypeRequest true "Exam type data"
// @Success 201 {object} response.Response{data=ExamTypeResponse}
// @Router /api/v1/exam-types [post]
func (h *Handler) CreateExamType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateExamTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	examType, err := h.service.CreateExamType(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamTypeCodeExists):
			apperrors.Abort(c, apperrors.Conflict("Exam type code already exists"))
		case errors.Is(err, ErrInvalidWeightage):
			apperrors.Abort(c, apperrors.BadRequest("Weightage must be between 0 and 100"))
		case errors.Is(err, ErrInvalidMaxMarks):
			apperrors.Abort(c, apperrors.BadRequest("Maximum marks must be greater than 0"))
		case errors.Is(err, ErrInvalidPassingMarks):
			apperrors.Abort(c, apperrors.BadRequest("Passing marks cannot exceed maximum marks"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create exam type"))
		}
		return
	}

	response.Created(c, ToExamTypeResponse(examType))
}

// UpdateExamType updates an existing exam type.
// @Summary Update exam type
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param id path string true "Exam Type ID"
// @Param request body UpdateExamTypeRequest true "Exam type data"
// @Success 200 {object} response.Response{data=ExamTypeResponse}
// @Router /api/v1/exam-types/{id} [put]
func (h *Handler) UpdateExamType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid exam type ID"))
		return
	}

	var req UpdateExamTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	examType, err := h.service.UpdateExamType(c.Request.Context(), tenantID, id, req, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamTypeNotFound):
			apperrors.Abort(c, apperrors.NotFound("Exam type not found"))
		case errors.Is(err, ErrExamTypeCodeExists):
			apperrors.Abort(c, apperrors.Conflict("Exam type code already exists"))
		case errors.Is(err, ErrInvalidWeightage):
			apperrors.Abort(c, apperrors.BadRequest("Weightage must be between 0 and 100"))
		case errors.Is(err, ErrInvalidMaxMarks):
			apperrors.Abort(c, apperrors.BadRequest("Maximum marks must be greater than 0"))
		case errors.Is(err, ErrInvalidPassingMarks):
			apperrors.Abort(c, apperrors.BadRequest("Passing marks cannot exceed maximum marks"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update exam type"))
		}
		return
	}

	response.OK(c, ToExamTypeResponse(examType))
}

// DeleteExamType deletes an exam type.
// @Summary Delete exam type
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param id path string true "Exam Type ID"
// @Success 204 "No Content"
// @Router /api/v1/exam-types/{id} [delete]
func (h *Handler) DeleteExamType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid exam type ID"))
		return
	}

	if err := h.service.DeleteExamType(c.Request.Context(), tenantID, id); err != nil {
		switch {
		case errors.Is(err, ErrExamTypeNotFound):
			apperrors.Abort(c, apperrors.NotFound("Exam type not found"))
		case errors.Is(err, ErrExamTypeInUse):
			apperrors.Abort(c, apperrors.Conflict("Exam type is in use and cannot be deleted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete exam type"))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleExamTypeActive toggles the active status of an exam type.
// @Summary Toggle exam type active status
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param id path string true "Exam Type ID"
// @Param request body ToggleActiveRequest true "Active status"
// @Success 200 {object} response.Response
// @Router /api/v1/exam-types/{id}/active [patch]
func (h *Handler) ToggleExamTypeActive(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid exam type ID"))
		return
	}

	var req ToggleActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	if err := h.service.ToggleExamTypeActive(c.Request.Context(), tenantID, id, req.IsActive); err != nil {
		if errors.Is(err, ErrExamTypeNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Exam type not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to toggle exam type status"))
		return
	}

	response.OK(c, gin.H{"message": "Exam type status updated"})
}

// UpdateDisplayOrders updates the display order of exam types.
// @Summary Update exam type display orders
// @Tags Exam Types
// @Accept json
// @Produce json
// @Param request body UpdateDisplayOrderRequest true "Display order data"
// @Success 200 {object} response.Response
// @Router /api/v1/exam-types/order [put]
func (h *Handler) UpdateDisplayOrders(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req UpdateDisplayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	if err := h.service.UpdateDisplayOrders(c.Request.Context(), tenantID, req); err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to update display orders"))
		return
	}

	response.OK(c, gin.H{"message": "Display orders updated"})
}
