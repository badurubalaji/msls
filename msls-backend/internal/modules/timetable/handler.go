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

	// Timetable routes
	timetables := rg.Group("/timetables")
	{
		// My timetable - available to all authenticated users
		timetables.GET("/teacher/me", h.GetMySchedule)

		// View operations
		timetablesView := timetables.Group("")
		timetablesView.Use(middleware.PermissionRequired("timetables:read"))
		{
			timetablesView.GET("", h.ListTimetables)
			timetablesView.GET("/:id", h.GetTimetable)
			timetablesView.GET("/:id/entries", h.GetTimetableEntries)
			timetablesView.GET("/conflicts", h.CheckConflicts)
			timetablesView.GET("/teacher/:staffId", h.GetTeacherSchedule)
		}

		// Manage operations
		timetablesManage := timetables.Group("")
		timetablesManage.Use(middleware.PermissionRequired("timetables:create"))
		{
			timetablesManage.POST("", h.CreateTimetable)
		}

		timetablesUpdate := timetables.Group("")
		timetablesUpdate.Use(middleware.PermissionRequired("timetables:update"))
		{
			timetablesUpdate.PUT("/:id", h.UpdateTimetable)
			timetablesUpdate.POST("/:id/entries", h.UpsertTimetableEntry)
			timetablesUpdate.POST("/:id/entries/bulk", h.BulkUpsertTimetableEntries)
			timetablesUpdate.DELETE("/:id/entries/:entryId", h.DeleteTimetableEntry)
		}

		timetablesPublish := timetables.Group("")
		timetablesPublish.Use(middleware.PermissionRequired("timetables:publish"))
		{
			timetablesPublish.POST("/:id/publish", h.PublishTimetable)
			timetablesPublish.POST("/:id/archive", h.ArchiveTimetable)
		}

		timetablesDelete := timetables.Group("")
		timetablesDelete.Use(middleware.PermissionRequired("timetables:delete"))
		{
			timetablesDelete.DELETE("/:id", h.DeleteTimetable)
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

// ========================================
// Timetable Handlers
// ========================================

// ListTimetables returns all timetables for the tenant.
func (h *Handler) ListTimetables(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	filter := TimetableFilter{TenantID: tenantID}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}

	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		sectionID, err := uuid.Parse(sectionIDStr)
		if err == nil {
			filter.SectionID = &sectionID
		}
	}

	if academicYearIDStr := c.Query("academic_year_id"); academicYearIDStr != "" {
		academicYearID, err := uuid.Parse(academicYearIDStr)
		if err == nil {
			filter.AcademicYearID = &academicYearID
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	timetables, total, err := h.service.ListTimetables(c.Request.Context(), filter)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to list timetables"))
		return
	}

	resp := TimetableListResponse{
		Timetables: make([]TimetableResponse, len(timetables)),
		Total:      total,
	}
	for i, tt := range timetables {
		resp.Timetables[i] = TimetableToResponse(&tt)
	}

	response.OK(c, resp)
}

// GetTimetable returns a single timetable by ID with entries.
func (h *Handler) GetTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	timetable, err := h.service.GetTimetableByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to get timetable"))
		return
	}

	response.OK(c, TimetableToResponse(timetable))
}

// CreateTimetable creates a new timetable.
func (h *Handler) CreateTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	timetable, err := h.service.CreateTimetable(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to create timetable"))
		return
	}

	response.Created(c, TimetableToResponse(timetable))
}

// UpdateTimetable updates an existing timetable.
func (h *Handler) UpdateTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	var req UpdateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	timetable, err := h.service.UpdateTimetable(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableNotDraft) {
			apperr.Abort(c, apperr.Conflict("Only draft timetables can be updated"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to update timetable"))
		return
	}

	response.OK(c, TimetableToResponse(timetable))
}

// PublishTimetable publishes a draft timetable.
func (h *Handler) PublishTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	timetable, err := h.service.PublishTimetable(c.Request.Context(), tenantID, id, userID)
	if err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableAlreadyPublished) {
			apperr.Abort(c, apperr.Conflict("Timetable is already published"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to publish timetable"))
		return
	}

	response.OK(c, TimetableToResponse(timetable))
}

// ArchiveTimetable archives a published timetable.
func (h *Handler) ArchiveTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	timetable, err := h.service.ArchiveTimetable(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to archive timetable"))
		return
	}

	response.OK(c, TimetableToResponse(timetable))
}

// DeleteTimetable deletes a draft timetable.
func (h *Handler) DeleteTimetable(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	if err := h.service.DeleteTimetable(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableNotDraft) {
			apperr.Abort(c, apperr.Conflict("Only draft timetables can be deleted"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete timetable"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Timetable Entry Handlers
// ========================================

// GetTimetableEntries returns all entries for a timetable.
func (h *Handler) GetTimetableEntries(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	entries, err := h.service.GetTimetableEntries(c.Request.Context(), id)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to get timetable entries"))
		return
	}

	resp := make([]TimetableEntryResponse, len(entries))
	for i, e := range entries {
		resp[i] = TimetableEntryToResponse(&e)
	}

	response.OK(c, resp)
}

// UpsertTimetableEntry creates or updates a timetable entry.
func (h *Handler) UpsertTimetableEntry(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	timetableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	var req CreateTimetableEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	entry, err := h.service.UpsertTimetableEntry(c.Request.Context(), tenantID, timetableID, req)
	if err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableNotDraft) {
			apperr.Abort(c, apperr.Conflict("Only draft timetables can be modified"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to save timetable entry"))
		return
	}

	response.OK(c, TimetableEntryToResponse(entry))
}

// BulkUpsertTimetableEntries creates or updates multiple entries.
func (h *Handler) BulkUpsertTimetableEntries(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	timetableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	var req BulkTimetableEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Abort(c, apperr.BadRequest(err.Error()))
		return
	}

	if err := h.service.BulkUpsertTimetableEntries(c.Request.Context(), tenantID, timetableID, req); err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableNotDraft) {
			apperr.Abort(c, apperr.Conflict("Only draft timetables can be modified"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to save timetable entries"))
		return
	}

	response.OK(c, gin.H{"message": "Entries saved successfully"})
}

// DeleteTimetableEntry deletes a timetable entry.
func (h *Handler) DeleteTimetableEntry(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	timetableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid timetable ID"))
		return
	}

	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid entry ID"))
		return
	}

	if err := h.service.DeleteTimetableEntry(c.Request.Context(), tenantID, timetableID, entryID); err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			apperr.Abort(c, apperr.NotFound("Timetable not found"))
			return
		}
		if errors.Is(err, ErrTimetableNotDraft) {
			apperr.Abort(c, apperr.Conflict("Only draft timetables can be modified"))
			return
		}
		apperr.Abort(c, apperr.InternalError("Failed to delete timetable entry"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Conflict Detection Handlers
// ========================================

// CheckConflicts checks for teacher conflicts.
func (h *Handler) CheckConflicts(c *gin.Context) {
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

	dayOfWeek, err := strconv.Atoi(c.Query("day_of_week"))
	if err != nil || dayOfWeek < 0 || dayOfWeek > 6 {
		apperr.Abort(c, apperr.BadRequest("Invalid day of week"))
		return
	}

	periodSlotIDStr := c.Query("period_slot_id")
	if periodSlotIDStr == "" {
		apperr.Abort(c, apperr.BadRequest("Period slot ID is required"))
		return
	}

	periodSlotID, err := uuid.Parse(periodSlotIDStr)
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid period slot ID"))
		return
	}

	var excludeTimetableID *uuid.UUID
	if excludeStr := c.Query("exclude_timetable_id"); excludeStr != "" {
		id, err := uuid.Parse(excludeStr)
		if err == nil {
			excludeTimetableID = &id
		}
	}

	conflicts, err := h.service.CheckTeacherConflicts(c.Request.Context(), tenantID, staffID, dayOfWeek, periodSlotID, excludeTimetableID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to check conflicts"))
		return
	}

	response.OK(c, ConflictCheckResponse{
		HasConflicts: len(conflicts) > 0,
		Conflicts:    conflicts,
	})
}

// GetTeacherSchedule returns a teacher's full schedule.
func (h *Handler) GetTeacherSchedule(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	staffID, err := uuid.Parse(c.Param("staffId"))
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid staff ID"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperr.Abort(c, apperr.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid academic year ID"))
		return
	}

	entries, err := h.service.GetTeacherSchedule(c.Request.Context(), tenantID, staffID, academicYearID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to get teacher schedule"))
		return
	}

	resp := make([]TimetableEntryResponse, len(entries))
	for i, e := range entries {
		resp[i] = TimetableEntryToResponse(&e)
	}

	response.OK(c, TeacherScheduleResponse{
		StaffID: staffID,
		Entries: resp,
	})
}

// GetMySchedule returns the current user's teaching schedule.
func (h *Handler) GetMySchedule(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperr.Abort(c, apperr.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperr.Abort(c, apperr.Unauthorized("User not authenticated"))
		return
	}

	// Get staff ID from user ID
	staffID, err := h.service.GetStaffIDByUserID(c.Request.Context(), tenantID, userID)
	if err != nil {
		apperr.Abort(c, apperr.NotFound("Staff profile not found for this user"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperr.Abort(c, apperr.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperr.Abort(c, apperr.BadRequest("Invalid academic year ID"))
		return
	}

	entries, err := h.service.GetTeacherSchedule(c.Request.Context(), tenantID, staffID, academicYearID)
	if err != nil {
		apperr.Abort(c, apperr.InternalError("Failed to get schedule"))
		return
	}

	resp := make([]TimetableEntryResponse, len(entries))
	for i, e := range entries {
		resp[i] = TimetableEntryToResponse(&e)
	}

	response.OK(c, TeacherScheduleResponse{
		StaffID: staffID,
		Entries: resp,
	})
}
