// Package timetable provides timetable management functionality.
package timetable

import (
	"errors"
	"net/http"
	"strconv"

	"msls-backend/internal/middleware"
	apperr "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for timetable entities.
type Handler struct {
	service *Service
}

// NewHandler creates a new timetable handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers timetable routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Shift routes
	shifts := rg.Group("/shifts")
	{
		// View operations
		shiftsView := shifts.Group("")
		shiftsView.Use(middleware.PermissionRequired("shift:view"))
		{
			shiftsView.GET("", h.ListShifts)
			shiftsView.GET("/:id", h.GetShift)
		}

		// Manage operations
		shiftsManage := shifts.Group("")
		shiftsManage.Use(middleware.PermissionRequired("shift:manage"))
		{
			shiftsManage.POST("", h.CreateShift)
			shiftsManage.PUT("/:id", h.UpdateShift)
			shiftsManage.DELETE("/:id", h.DeleteShift)
		}
	}

	// Day Pattern routes
	dayPatterns := rg.Group("/day-patterns")
	{
		// View operations
		dayPatternsView := dayPatterns.Group("")
		dayPatternsView.Use(middleware.PermissionRequired("timetable:view"))
		{
			dayPatternsView.GET("", h.ListDayPatterns)
			dayPatternsView.GET("/:id", h.GetDayPattern)
		}

		// Manage operations
		dayPatternsManage := dayPatterns.Group("")
		dayPatternsManage.Use(middleware.PermissionRequired("timetable:manage"))
		{
			dayPatternsManage.POST("", h.CreateDayPattern)
			dayPatternsManage.PUT("/:id", h.UpdateDayPattern)
			dayPatternsManage.DELETE("/:id", h.DeleteDayPattern)
		}
	}

	// Day Pattern Assignment routes
	dayAssignments := rg.Group("/day-pattern-assignments")
	{
		// View operations
		dayAssignmentsView := dayAssignments.Group("")
		dayAssignmentsView.Use(middleware.PermissionRequired("timetable:view"))
		{
			dayAssignmentsView.GET("", h.ListDayPatternAssignments)
		}

		// Manage operations
		dayAssignmentsManage := dayAssignments.Group("")
		dayAssignmentsManage.Use(middleware.PermissionRequired("timetable:manage"))
		{
			dayAssignmentsManage.PUT("/:dayOfWeek", h.UpdateDayPatternAssignment)
		}
	}

	// Period Slot routes
	periodSlots := rg.Group("/period-slots")
	{
		// View operations
		periodSlotsView := periodSlots.Group("")
		periodSlotsView.Use(middleware.PermissionRequired("timetable:view"))
		{
			periodSlotsView.GET("", h.ListPeriodSlots)
			periodSlotsView.GET("/:id", h.GetPeriodSlot)
		}

		// Manage operations
		periodSlotsManage := periodSlots.Group("")
		periodSlotsManage.Use(middleware.PermissionRequired("timetable:manage"))
		{
			periodSlotsManage.POST("", h.CreatePeriodSlot)
			periodSlotsManage.PUT("/:id", h.UpdatePeriodSlot)
			periodSlotsManage.DELETE("/:id", h.DeletePeriodSlot)
		}
	}
}

// ========================================
// Shift Handlers
// ========================================

// ListShifts returns all shifts for the tenant.
func (h *Handler) ListShifts(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	filter := ShiftFilter{TenantID: tenantID}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	shifts, total, err := h.service.ListShifts(c.Request.Context(), filter)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list shifts"))
		return
	}

	resp := ShiftListResponse{
		Shifts: make([]ShiftResponse, len(shifts)),
		Total:  total,
	}
	for i, shift := range shifts {
		resp.Shifts[i] = ShiftToResponse(&shift)
	}

	response.OK(c, resp)
}

// GetShift returns a single shift by ID.
func (h *Handler) GetShift(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid shift ID"))
		return
	}

	shift, err := h.service.GetShiftByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrShiftNotFound) {
			apperr.Abort(c, apperr.NotFound("Shift not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to get shift"))
		return
	}

	response.OK(c, ShiftToResponse(shift))
}

// CreateShift creates a new shift.
func (h *Handler) CreateShift(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	shift, err := h.service.CreateShift(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrShiftCodeExists) {
			apperr.Abort(c, apperr.Conflict("Shift code already exists"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to create shift"))
		return
	}

	response.Created(c, ShiftToResponse(shift))
}

// UpdateShift updates an existing shift.
func (h *Handler) UpdateShift(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid shift ID"))
		return
	}

	var req UpdateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	shift, err := h.service.UpdateShift(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrShiftNotFound) {
			apperr.Abort(c, apperr.NotFound("Shift not found"))
			return
		}
		if errors.Is(err, ErrShiftCodeExists) {
			apperr.Abort(c, apperr.Conflict("Shift code already exists"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to update shift"))
		return
	}

	response.OK(c, ShiftToResponse(shift))
}

// DeleteShift deletes a shift.
func (h *Handler) DeleteShift(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid shift ID"))
		return
	}

	if err := h.service.DeleteShift(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrShiftNotFound) {
			apperr.Abort(c, apperr.NotFound("Shift not found"))
			return
		}
		if errors.Is(err, ErrShiftInUse) {
			apperr.Abort(c, apperr.Conflict("Cannot delete shift that is in use"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete shift"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Day Pattern Handlers
// ========================================

// ListDayPatterns returns all day patterns for the tenant.
func (h *Handler) ListDayPatterns(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	filter := DayPatternFilter{TenantID: tenantID}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	patterns, total, err := h.service.ListDayPatterns(c.Request.Context(), filter)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list day patterns"))
		return
	}

	resp := DayPatternListResponse{
		DayPatterns: make([]DayPatternResponse, len(patterns)),
		Total:       total,
	}
	for i, pattern := range patterns {
		resp.DayPatterns[i] = DayPatternToResponse(&pattern)
	}

	response.OK(c, resp)
}

// GetDayPattern returns a single day pattern by ID.
func (h *Handler) GetDayPattern(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid day pattern ID"))
		return
	}

	pattern, err := h.service.GetDayPatternByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrDayPatternNotFound) {
			apperr.Abort(c, apperr.NotFound("Day pattern not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to get day pattern"))
		return
	}

	response.OK(c, DayPatternToResponse(pattern))
}

// CreateDayPattern creates a new day pattern.
func (h *Handler) CreateDayPattern(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateDayPatternRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	pattern, err := h.service.CreateDayPattern(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrDayPatternCodeExists) {
			apperr.Abort(c, apperr.Conflict("Day pattern code already exists"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to create day pattern"))
		return
	}

	response.Created(c, DayPatternToResponse(pattern))
}

// UpdateDayPattern updates an existing day pattern.
func (h *Handler) UpdateDayPattern(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid day pattern ID"))
		return
	}

	var req UpdateDayPatternRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	pattern, err := h.service.UpdateDayPattern(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrDayPatternNotFound) {
			apperr.Abort(c, apperr.NotFound("Day pattern not found"))
			return
		}
		if errors.Is(err, ErrDayPatternCodeExists) {
			apperr.Abort(c, apperr.Conflict("Day pattern code already exists"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to update day pattern"))
		return
	}

	response.OK(c, DayPatternToResponse(pattern))
}

// DeleteDayPattern deletes a day pattern.
func (h *Handler) DeleteDayPattern(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid day pattern ID"))
		return
	}

	if err := h.service.DeleteDayPattern(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrDayPatternNotFound) {
			apperr.Abort(c, apperr.NotFound("Day pattern not found"))
			return
		}
		if errors.Is(err, ErrDayPatternInUse) {
			apperr.Abort(c, apperr.Conflict("Cannot delete day pattern that is in use"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete day pattern"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Day Pattern Assignment Handlers
// ========================================

// ListDayPatternAssignments returns all day pattern assignments for a branch.
func (h *Handler) ListDayPatternAssignments(c *gin.Context) {
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

	assignments, err := h.service.ListDayPatternAssignments(c.Request.Context(), tenantID, branchID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list day pattern assignments"))
		return
	}

	resp := DayPatternAssignmentListResponse{
		Assignments: make([]DayPatternAssignmentResponse, len(assignments)),
	}
	for i, assignment := range assignments {
		resp.Assignments[i] = DayPatternAssignmentToResponse(&assignment)
	}

	response.OK(c, resp)
}

// UpdateDayPatternAssignment updates a day pattern assignment.
func (h *Handler) UpdateDayPatternAssignment(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	dayOfWeek, err := strconv.Atoi(c.Param("dayOfWeek"))
	if err != nil || dayOfWeek < 0 || dayOfWeek > 6 {
		apperr.Abort(c, apperr.BadRequest("Invalid day of week (must be 0-6)"))
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

	var req UpdateDayPatternAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	assignment, err := h.service.UpdateDayPatternAssignment(c.Request.Context(), tenantID, branchID, dayOfWeek, req)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to update day pattern assignment"))
		return
	}

	response.OK(c, DayPatternAssignmentToResponse(assignment))
}

// ========================================
// Period Slot Handlers
// ========================================

// ListPeriodSlots returns all period slots for the tenant.
func (h *Handler) ListPeriodSlots(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	filter := PeriodSlotFilter{TenantID: tenantID}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}

	if dayPatternIDStr := c.Query("day_pattern_id"); dayPatternIDStr != "" {
		dayPatternID, err := uuid.Parse(dayPatternIDStr)
		if err == nil {
			filter.DayPatternID = &dayPatternID
		}
	}

	if shiftIDStr := c.Query("shift_id"); shiftIDStr != "" {
		shiftID, err := uuid.Parse(shiftIDStr)
		if err == nil {
			filter.ShiftID = &shiftID
		}
	}

	if slotType := c.Query("slot_type"); slotType != "" {
		filter.SlotType = &slotType
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	slots, total, err := h.service.ListPeriodSlots(c.Request.Context(), filter)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list period slots"))
		return
	}

	resp := PeriodSlotListResponse{
		PeriodSlots: make([]PeriodSlotResponse, len(slots)),
		Total:       total,
	}
	for i, slot := range slots {
		resp.PeriodSlots[i] = PeriodSlotToResponse(&slot)
	}

	response.OK(c, resp)
}

// GetPeriodSlot returns a single period slot by ID.
func (h *Handler) GetPeriodSlot(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid period slot ID"))
		return
	}

	slot, err := h.service.GetPeriodSlotByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPeriodSlotNotFound) {
			apperr.Abort(c, apperr.NotFound("Period slot not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to get period slot"))
		return
	}

	response.OK(c, PeriodSlotToResponse(slot))
}

// CreatePeriodSlot creates a new period slot.
func (h *Handler) CreatePeriodSlot(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreatePeriodSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	slot, err := h.service.CreatePeriodSlot(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to create period slot"))
		return
	}

	response.Created(c, PeriodSlotToResponse(slot))
}

// UpdatePeriodSlot updates an existing period slot.
func (h *Handler) UpdatePeriodSlot(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid period slot ID"))
		return
	}

	var req UpdatePeriodSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	slot, err := h.service.UpdatePeriodSlot(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrPeriodSlotNotFound) {
			apperr.Abort(c, apperr.NotFound("Period slot not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to update period slot"))
		return
	}

	response.OK(c, PeriodSlotToResponse(slot))
}

// DeletePeriodSlot deletes a period slot.
func (h *Handler) DeletePeriodSlot(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid period slot ID"))
		return
	}

	if err := h.service.DeletePeriodSlot(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrPeriodSlotNotFound) {
			apperr.Abort(c, apperr.NotFound("Period slot not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete period slot"))
		return
	}

	c.Status(http.StatusNoContent)
}
