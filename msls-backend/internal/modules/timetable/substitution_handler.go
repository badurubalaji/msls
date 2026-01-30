package timetable

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"msls-backend/internal/middleware"
	apperr "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterSubstitutionRoutes registers substitution routes.
func (h *Handler) RegisterSubstitutionRoutes(rg *gin.RouterGroup) {
	substitutions := rg.Group("/substitutions")
	{
		// View operations
		subsView := substitutions.Group("")
		subsView.Use(middleware.PermissionRequired("substitution:view"))
		{
			subsView.GET("", h.ListSubstitutions)
			subsView.GET("/:id", h.GetSubstitution)
			subsView.GET("/available-teachers", h.GetAvailableTeachers)
			subsView.GET("/teacher-periods", h.GetTeacherAbsencePeriods)
		}

		// Create operations
		subsCreate := substitutions.Group("")
		subsCreate.Use(middleware.PermissionRequired("substitution:create"))
		{
			subsCreate.POST("", h.CreateSubstitution)
		}

		// Update operations
		subsUpdate := substitutions.Group("")
		subsUpdate.Use(middleware.PermissionRequired("substitution:update"))
		{
			subsUpdate.PUT("/:id", h.UpdateSubstitution)
		}

		// Approve operations
		subsApprove := substitutions.Group("")
		subsApprove.Use(middleware.PermissionRequired("substitution:approve"))
		{
			subsApprove.POST("/:id/confirm", h.ConfirmSubstitution)
			subsApprove.POST("/:id/cancel", h.CancelSubstitution)
		}

		// Delete operations
		subsDelete := substitutions.Group("")
		subsDelete.Use(middleware.PermissionRequired("substitution:delete"))
		{
			subsDelete.DELETE("/:id", h.DeleteSubstitution)
		}
	}
}

// ========================================
// Substitution Handlers
// ========================================

// ListSubstitutions returns all substitutions for the tenant.
func (h *Handler) ListSubstitutions(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	filter := SubstitutionFilter{TenantID: tenantID}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}

	if staffIDStr := c.Query("original_staff_id"); staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err == nil {
			filter.OriginalStaffID = &staffID
		}
	}

	if staffIDStr := c.Query("substitute_staff_id"); staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err == nil {
			filter.SubstituteStaffID = &staffID
		}
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	substitutions, total, err := h.service.ListSubstitutions(c.Request.Context(), filter)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list substitutions"))
		return
	}

	resp := SubstitutionListResponse{
		Substitutions: make([]SubstitutionResponse, len(substitutions)),
		Total:         total,
	}
	for i, sub := range substitutions {
		resp.Substitutions[i] = SubstitutionToResponse(&sub)
	}

	response.OK(c, resp)
}

// GetSubstitution returns a single substitution by ID.
func (h *Handler) GetSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid substitution ID"))
		return
	}

	substitution, err := h.service.GetSubstitutionByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrSubstitutionNotFound) {
			apperr.Abort(c, apperr.NotFound("Substitution not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to get substitution"))
		return
	}

	response.OK(c, SubstitutionToResponse(substitution))
}

// CreateSubstitution creates a new substitution.
func (h *Handler) CreateSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateSubstitutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	substitution, err := h.service.CreateSubstitution(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrSubstitutionConflict) {
			apperr.Abort(c, apperr.Conflict("Substitution already exists for this teacher and period"))
			return
		}
		if errors.Is(err, ErrSubstituteConflict) {
			apperr.Abort(c, apperr.Conflict("Substitute teacher has a conflict at this time"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to create substitution"))
		return
	}

	response.Created(c, SubstitutionToResponse(substitution))
}

// UpdateSubstitution updates an existing substitution.
func (h *Handler) UpdateSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid substitution ID"))
		return
	}

	var req UpdateSubstitutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	substitution, err := h.service.UpdateSubstitution(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrSubstitutionNotFound) {
			apperr.Abort(c, apperr.NotFound("Substitution not found"))
			return
		}
		if errors.Is(err, ErrSubstitutionNotPending) {
			apperr.Abort(c, apperr.Conflict("Only pending substitutions can be modified"))
			return
		}
		if errors.Is(err, ErrSubstituteConflict) {
			apperr.Abort(c, apperr.Conflict("Substitute teacher has a conflict at this time"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to update substitution"))
		return
	}

	response.OK(c, SubstitutionToResponse(substitution))
}

// ConfirmSubstitution confirms a pending substitution.
func (h *Handler) ConfirmSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid substitution ID"))
		return
	}

	substitution, err := h.service.ConfirmSubstitution(c.Request.Context(), tenantID, id, userID)
	if err != nil {
		if errors.Is(err, ErrSubstitutionNotFound) {
			apperr.Abort(c, apperr.NotFound("Substitution not found"))
			return
		}
		if errors.Is(err, ErrSubstitutionNotPending) {
			apperr.Abort(c, apperr.Conflict("Only pending substitutions can be confirmed"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to confirm substitution"))
		return
	}

	response.OK(c, SubstitutionToResponse(substitution))
}

// CancelSubstitution cancels a substitution.
func (h *Handler) CancelSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid substitution ID"))
		return
	}

	substitution, err := h.service.CancelSubstitution(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrSubstitutionNotFound) {
			apperr.Abort(c, apperr.NotFound("Substitution not found"))
			return
		}
		if errors.Is(err, ErrSubstitutionNotCancellable) {
			apperr.Abort(c, apperr.Conflict("Only pending or confirmed substitutions can be cancelled"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to cancel substitution"))
		return
	}

	response.OK(c, SubstitutionToResponse(substitution))
}

// DeleteSubstitution deletes a substitution.
func (h *Handler) DeleteSubstitution(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid substitution ID"))
		return
	}

	if err := h.service.DeleteSubstitution(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrSubstitutionNotFound) {
			apperr.Abort(c, apperr.NotFound("Substitution not found"))
			return
		}
		if errors.Is(err, ErrSubstitutionNotPending) {
			apperr.Abort(c, apperr.Conflict("Only pending substitutions can be deleted"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete substitution"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAvailableTeachers returns teachers available for substitution.
func (h *Handler) GetAvailableTeachers(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	branchIDStr := c.Query("branch_id")
	if branchIDStr == "" {
		apperr.Abort(c, apperr.BadRequest("Branch ID is required"))
		return
	}

	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid branch ID"))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		apperr.Abort(c, apperr.BadRequest("Date is required"))
		return
	}

	periodSlotIDsStr := c.Query("period_slot_ids")
	if periodSlotIDsStr == "" {
		apperr.Abort(c, apperr.BadRequest("Period slot IDs are required"))
		return
	}

	periodSlotIDs := make([]uuid.UUID, 0)
	for _, idStr := range strings.Split(periodSlotIDsStr, ",") {
		id, err := uuid.Parse(strings.TrimSpace(idStr))
		if err == nil {
			periodSlotIDs = append(periodSlotIDs, id)
		}
	}

	if len(periodSlotIDs) == 0 {
		apperr.Abort(c, apperr.BadRequest("At least one valid period slot ID is required"))
		return
	}

	excludeStaffIDStr := c.Query("exclude_staff_id")
	var excludeStaffID uuid.UUID
	if excludeStaffIDStr != "" {
		excludeStaffID, _ = uuid.Parse(excludeStaffIDStr)
	}

	teachers, err := h.service.GetAvailableTeachers(c.Request.Context(), tenantID, branchID, dateStr, periodSlotIDs, excludeStaffID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to get available teachers"))
		return
	}

	response.OK(c, AvailableTeachersResponse{Teachers: teachers})
}

// GetTeacherAbsencePeriods returns the timetable entries for an absent teacher.
func (h *Handler) GetTeacherAbsencePeriods(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	staffIDStr := c.Query("staff_id")
	if staffIDStr == "" {
		apperr.Abort(c, apperr.BadRequest("Staff ID is required"))
		return
	}

	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid staff ID"))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		apperr.Abort(c, apperr.BadRequest("Date is required"))
		return
	}

	entries, err := h.service.GetTeacherAbsencePeriods(c.Request.Context(), tenantID, staffID, dateStr)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to get teacher periods"))
		return
	}

	response.OK(c, entries)
}
