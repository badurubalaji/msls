// Package attendance provides staff attendance management functionality.
package attendance

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles attendance-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new attendance handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// MarkAttendanceRequest represents the request body for marking attendance (HR).
type MarkAttendanceRequest struct {
	StaffID        string  `json:"staffId" binding:"required,uuid"`
	AttendanceDate string  `json:"attendanceDate" binding:"required"` // Format: YYYY-MM-DD
	Status         string  `json:"status" binding:"required,oneof=present absent half_day on_leave holiday"`
	CheckInTime    *string `json:"checkInTime"`  // Format: HH:MM
	CheckOutTime   *string `json:"checkOutTime"` // Format: HH:MM
	HalfDayType    string  `json:"halfDayType" binding:"omitempty,oneof=first_half second_half"`
	Remarks        string  `json:"remarks" binding:"max=500"`
}

// RegularizationReviewRequest represents the request body for reviewing regularization.
type RegularizationReviewRequest struct {
	RejectionReason string `json:"rejectionReason" binding:"max=1000"`
}

// SettingsRequest represents the request body for updating attendance settings.
type SettingsRequest struct {
	BranchID                      string  `json:"branchId" binding:"required,uuid"`
	WorkStartTime                 string  `json:"workStartTime" binding:"required"` // Format: HH:MM
	WorkEndTime                   string  `json:"workEndTime" binding:"required"`   // Format: HH:MM
	LateThresholdMinutes          int     `json:"lateThresholdMinutes" binding:"min=0,max=120"`
	HalfDayThresholdHours         float64 `json:"halfDayThresholdHours" binding:"min=0,max=12"`
	AllowSelfCheckout             bool    `json:"allowSelfCheckout"`
	RequireRegularizationApproval bool    `json:"requireRegularizationApproval"`
}

// CheckInRequest with staff ID for self check-in
type CheckInSelfRequest struct {
	StaffID     string `json:"staffId" binding:"required,uuid"`
	HalfDayType string `json:"halfDayType" binding:"omitempty,oneof=first_half second_half"`
	Remarks     string `json:"remarks" binding:"max=500"`
}

// CheckIn handles staff check-in.
// @Summary Check in
// @Description Mark check-in for a staff member (self service)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CheckInSelfRequest true "Check-in details"
// @Success 200 {object} response.Success{data=AttendanceResponse}
// @Router /api/v1/attendance/check-in [post]
func (h *Handler) CheckIn(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User not authenticated"))
		return
	}

	var req CheckInSelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	dto := CheckInDTO{
		TenantID:    tenantID,
		StaffID:     staffID,
		HalfDayType: HalfDayType(req.HalfDayType),
		Remarks:     req.Remarks,
		MarkedBy:    &userID,
	}

	attendance, err := h.service.CheckIn(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlreadyCheckedIn):
			apperrors.Abort(c, apperrors.Conflict("Already checked in for today"))
		case errors.Is(err, ErrStaffNotFound):
			apperrors.Abort(c, apperrors.NotFound("Staff profile not found"))
		default:
			logger.Error("Failed to check in",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", staffID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to check in"))
		}
		return
	}

	response.OK(c, toAttendanceResponse(attendance))
}

// CheckOutSelfRequest for self check-out
type CheckOutSelfRequest struct {
	StaffID string `json:"staffId" binding:"required,uuid"`
	Remarks string `json:"remarks" binding:"max=500"`
}

// CheckOut handles staff check-out.
// @Summary Check out
// @Description Mark check-out for a staff member (self service)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CheckOutSelfRequest true "Check-out details"
// @Success 200 {object} response.Success{data=AttendanceResponse}
// @Router /api/v1/attendance/check-out [post]
func (h *Handler) CheckOut(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User not authenticated"))
		return
	}

	var req CheckOutSelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	dto := CheckOutDTO{
		TenantID: tenantID,
		StaffID:  staffID,
		Remarks:  req.Remarks,
		MarkedBy: &userID,
	}

	attendance, err := h.service.CheckOut(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotCheckedIn):
			apperrors.Abort(c, apperrors.BadRequest("Not checked in yet"))
		case errors.Is(err, ErrAlreadyCheckedOut):
			apperrors.Abort(c, apperrors.Conflict("Already checked out for today"))
		default:
			logger.Error("Failed to check out",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", staffID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to check out"))
		}
		return
	}

	response.OK(c, toAttendanceResponse(attendance))
}

// GetMyToday retrieves today's attendance for a staff member.
// @Summary Get today's attendance
// @Description Get today's attendance status for a staff member
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string true "Staff ID"
// @Success 200 {object} response.Success{data=TodayAttendanceResponse}
// @Router /api/v1/attendance/my/today [get]
func (h *Handler) GetMyToday(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffIDStr := c.Query("staff_id")
	if staffIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Staff ID is required"))
		return
	}

	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	attendance, err := h.service.GetTodayAttendance(c.Request.Context(), tenantID, staffID)

	resp := TodayAttendanceResponse{
		CanCheckIn:  true,
		CanCheckOut: false,
	}

	if err != nil {
		if errors.Is(err, ErrAttendanceNotFound) {
			resp.Status = "not_marked"
			response.OK(c, resp)
			return
		}
		logger.Error("Failed to get today's attendance",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", staffID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to get today's attendance"))
		return
	}

	attendanceResp := toAttendanceResponse(attendance)
	resp.Attendance = &attendanceResp

	if attendance.CheckInTime != nil && attendance.CheckOutTime == nil {
		resp.Status = "checked_in"
		resp.CanCheckIn = false
		resp.CanCheckOut = true
	} else if attendance.CheckOutTime != nil {
		resp.Status = "checked_out"
		resp.CanCheckIn = false
		resp.CanCheckOut = false
	} else {
		resp.Status = "not_marked"
	}

	response.OK(c, resp)
}

// GetMyAttendance retrieves attendance records for a staff member.
// @Summary Get my attendance
// @Description Get attendance records for a staff member
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string true "Staff ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} response.Success{data=AttendanceListResponse}
// @Router /api/v1/attendance/my [get]
func (h *Handler) GetMyAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffIDStr := c.Query("staff_id")
	if staffIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Staff ID is required"))
		return
	}

	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	filter := ListFilter{
		TenantID: tenantID,
		StaffID:  &staffID,
		Cursor:   c.Query("cursor"),
	}

	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			filter.DateFrom = &dateFrom
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			filter.DateTo = &dateTo
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	attendanceList, nextCursor, total, err := h.service.ListAttendance(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list my attendance",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", staffID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve attendance records"))
		return
	}

	resp := AttendanceListResponse{
		Attendance: toAttendanceResponses(attendanceList),
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}

	response.OK(c, resp)
}

// GetMySummary retrieves monthly attendance summary for a staff member.
// @Summary Get my monthly summary
// @Description Get monthly attendance summary for a staff member
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string true "Staff ID"
// @Param year query int false "Year (default: current year)"
// @Param month query int false "Month (default: current month)"
// @Success 200 {object} response.Success{data=AttendanceSummaryResponse}
// @Router /api/v1/attendance/my/summary [get]
func (h *Handler) GetMySummary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffIDStr := c.Query("staff_id")
	if staffIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Staff ID is required"))
		return
	}

	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}

	summary, err := h.service.GetMonthlySummary(c.Request.Context(), tenantID, staffID, year, month)
	if err != nil {
		logger.Error("Failed to get monthly summary",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", staffID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve monthly summary"))
		return
	}

	response.OK(c, summary)
}

// List retrieves all attendance records (HR view).
// @Summary List all attendance
// @Description Get all staff attendance records with filters (HR view)
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string false "Filter by staff ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param department_id query string false "Filter by department ID"
// @Param status query string false "Filter by status"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} response.Success{data=AttendanceListResponse}
// @Router /api/v1/attendance [get]
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ListFilter{
		TenantID: tenantID,
		Cursor:   c.Query("cursor"),
	}

	// Parse UUID filters
	if staffIDStr := c.Query("staff_id"); staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err == nil {
			filter.StaffID = &staffID
		}
	}
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		deptID, err := uuid.Parse(deptIDStr)
		if err == nil {
			filter.DepartmentID = &deptID
		}
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := AttendanceStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			filter.DateFrom = &dateFrom
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			filter.DateTo = &dateTo
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	attendanceList, nextCursor, total, err := h.service.ListAttendance(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list attendance",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve attendance records"))
		return
	}

	resp := AttendanceListResponse{
		Attendance: toAttendanceResponses(attendanceList),
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}

	response.OK(c, resp)
}

// MarkAttendance marks attendance for a staff member (HR action).
// @Summary Mark attendance (HR)
// @Description Mark attendance for a staff member (HR action)
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body MarkAttendanceRequest true "Attendance details"
// @Success 201 {object} response.Success{data=AttendanceResponse}
// @Router /api/v1/attendance/mark [post]
func (h *Handler) MarkAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	attendanceDate, err := time.Parse("2006-01-02", req.AttendanceDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid attendance date format, use YYYY-MM-DD"))
		return
	}

	dto := MarkAttendanceDTO{
		TenantID:       tenantID,
		StaffID:        staffID,
		AttendanceDate: attendanceDate,
		Status:         AttendanceStatus(req.Status),
		HalfDayType:    HalfDayType(req.HalfDayType),
		Remarks:        req.Remarks,
		MarkedBy:       &userID,
	}

	// Parse check-in time if provided
	if req.CheckInTime != nil {
		checkIn, err := time.Parse("15:04", *req.CheckInTime)
		if err == nil {
			fullCheckIn := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(),
				checkIn.Hour(), checkIn.Minute(), 0, 0, time.Local)
			dto.CheckInTime = &fullCheckIn
		}
	}

	// Parse check-out time if provided
	if req.CheckOutTime != nil {
		checkOut, err := time.Parse("15:04", *req.CheckOutTime)
		if err == nil {
			fullCheckOut := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(),
				checkOut.Hour(), checkOut.Minute(), 0, 0, time.Local)
			dto.CheckOutTime = &fullCheckOut
		}
	}

	attendance, err := h.service.MarkAttendance(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStaffNotFound):
			apperrors.Abort(c, apperrors.NotFound("Staff not found"))
		case errors.Is(err, ErrFutureDate):
			apperrors.Abort(c, apperrors.BadRequest("Cannot mark attendance for future date"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid attendance status"))
		default:
			logger.Error("Failed to mark attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", staffID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to mark attendance"))
		}
		return
	}

	response.Created(c, toAttendanceResponse(attendance))
}

// RegularizationSelfRequest includes staff ID for self-service
type RegularizationSelfRequest struct {
	StaffID               string `json:"staffId" binding:"required,uuid"`
	RequestDate           string `json:"requestDate" binding:"required"` // Format: YYYY-MM-DD
	RequestedStatus       string `json:"requestedStatus" binding:"required,oneof=present half_day"`
	Reason                string `json:"reason" binding:"required,max=1000"`
	SupportingDocumentURL string `json:"supportingDocumentUrl"`
}

// SubmitRegularization submits a regularization request.
// @Summary Submit regularization request
// @Description Submit a regularization request for attendance
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body RegularizationSelfRequest true "Regularization details"
// @Success 201 {object} response.Success{data=RegularizationResponse}
// @Router /api/v1/attendance/regularization [post]
func (h *Handler) SubmitRegularization(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req RegularizationSelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	requestDate, err := time.Parse("2006-01-02", req.RequestDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid request date format, use YYYY-MM-DD"))
		return
	}

	dto := RegularizationRequestDTO{
		TenantID:              tenantID,
		StaffID:               staffID,
		RequestDate:           requestDate,
		RequestedStatus:       AttendanceStatus(req.RequestedStatus),
		Reason:                req.Reason,
		SupportingDocumentURL: req.SupportingDocumentURL,
	}

	regularization, err := h.service.RequestRegularization(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrFutureDate):
			apperrors.Abort(c, apperrors.BadRequest("Cannot request regularization for future date"))
		case errors.Is(err, ErrCannotRegularizePendingRequest):
			apperrors.Abort(c, apperrors.Conflict("Pending regularization request already exists for this date"))
		default:
			logger.Error("Failed to submit regularization",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", staffID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to submit regularization request"))
		}
		return
	}

	response.Created(c, toRegularizationResponse(regularization))
}

// ListRegularizations lists regularization requests.
// @Summary List regularization requests
// @Description List regularization requests with filters (HR view)
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string false "Filter by staff ID"
// @Param status query string false "Filter by status (pending, approved, rejected)"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} response.Success{data=RegularizationListResponse}
// @Router /api/v1/attendance/regularization [get]
func (h *Handler) ListRegularizations(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := RegularizationFilter{
		TenantID: tenantID,
		Cursor:   c.Query("cursor"),
	}

	// Parse UUID filters
	if staffIDStr := c.Query("staff_id"); staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err == nil {
			filter.StaffID = &staffID
		}
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := RegularizationStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			filter.DateFrom = &dateFrom
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			filter.DateTo = &dateTo
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	regularizations, nextCursor, total, err := h.service.ListRegularizations(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list regularizations",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve regularization requests"))
		return
	}

	resp := RegularizationListResponse{
		Regularizations: toRegularizationResponses(regularizations),
		NextCursor:      nextCursor,
		HasMore:         nextCursor != "",
		Total:           total,
	}

	response.OK(c, resp)
}

// ApproveRegularization approves a regularization request.
// @Summary Approve regularization
// @Description Approve a regularization request
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Regularization ID"
// @Success 200 {object} response.Success{data=RegularizationResponse}
// @Router /api/v1/attendance/regularization/{id}/approve [put]
func (h *Handler) ApproveRegularization(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid regularization ID"))
		return
	}

	dto := RegularizationReviewDTO{
		TenantID:         tenantID,
		RegularizationID: id,
		Approved:         true,
		ReviewedBy:       &userID,
	}

	regularization, err := h.service.ApproveRegularization(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrRegularizationNotFound):
			apperrors.Abort(c, apperrors.NotFound("Regularization request not found"))
		case errors.Is(err, ErrRegularizationAlreadyProcessed):
			apperrors.Abort(c, apperrors.Conflict("Regularization request already processed"))
		default:
			logger.Error("Failed to approve regularization",
				zap.String("tenant_id", tenantID.String()),
				zap.String("regularization_id", id.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to approve regularization"))
		}
		return
	}

	response.OK(c, toRegularizationResponse(regularization))
}

// RejectRegularization rejects a regularization request.
// @Summary Reject regularization
// @Description Reject a regularization request
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Regularization ID"
// @Param request body RegularizationReviewRequest true "Rejection details"
// @Success 200 {object} response.Success{data=RegularizationResponse}
// @Router /api/v1/attendance/regularization/{id}/reject [put]
func (h *Handler) RejectRegularization(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid regularization ID"))
		return
	}

	var req RegularizationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := RegularizationReviewDTO{
		TenantID:         tenantID,
		RegularizationID: id,
		Approved:         false,
		RejectionReason:  req.RejectionReason,
		ReviewedBy:       &userID,
	}

	regularization, err := h.service.RejectRegularization(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrRegularizationNotFound):
			apperrors.Abort(c, apperrors.NotFound("Regularization request not found"))
		case errors.Is(err, ErrRegularizationAlreadyProcessed):
			apperrors.Abort(c, apperrors.Conflict("Regularization request already processed"))
		case errors.Is(err, ErrReasonRequired):
			apperrors.Abort(c, apperrors.BadRequest("Rejection reason is required"))
		default:
			logger.Error("Failed to reject regularization",
				zap.String("tenant_id", tenantID.String()),
				zap.String("regularization_id", id.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to reject regularization"))
		}
		return
	}

	response.OK(c, toRegularizationResponse(regularization))
}

// GetSettings retrieves attendance settings for a branch.
// @Summary Get attendance settings
// @Description Get attendance settings for a branch
// @Tags Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string true "Branch ID"
// @Success 200 {object} response.Success{data=SettingsResponse}
// @Router /api/v1/attendance/settings [get]
func (h *Handler) GetSettings(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	branchIDStr := c.Query("branch_id")
	if branchIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Branch ID is required"))
		return
	}

	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	settings, err := h.service.GetSettings(c.Request.Context(), tenantID, branchID)
	if err != nil {
		if errors.Is(err, ErrSettingsNotFound) {
			// Return default settings
			response.OK(c, SettingsResponse{
				BranchID:                      branchID.String(),
				WorkStartTime:                 "09:00",
				WorkEndTime:                   "17:00",
				LateThresholdMinutes:          15,
				HalfDayThresholdHours:         4.0,
				AllowSelfCheckout:             true,
				RequireRegularizationApproval: true,
			})
			return
		}
		logger.Error("Failed to get attendance settings",
			zap.String("tenant_id", tenantID.String()),
			zap.String("branch_id", branchID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve attendance settings"))
		return
	}

	response.OK(c, toSettingsResponse(settings))
}

// UpdateSettings updates attendance settings for a branch.
// @Summary Update attendance settings
// @Description Update attendance settings for a branch
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body SettingsRequest true "Settings details"
// @Success 200 {object} response.Success{data=SettingsResponse}
// @Router /api/v1/attendance/settings [put]
func (h *Handler) UpdateSettings(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req SettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	workStartTime, err := time.Parse("15:04", req.WorkStartTime)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid work start time format, use HH:MM"))
		return
	}

	workEndTime, err := time.Parse("15:04", req.WorkEndTime)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid work end time format, use HH:MM"))
		return
	}

	dto := SettingsDTO{
		TenantID:                      tenantID,
		BranchID:                      branchID,
		WorkStartTime:                 workStartTime,
		WorkEndTime:                   workEndTime,
		LateThresholdMinutes:          req.LateThresholdMinutes,
		HalfDayThresholdHours:         req.HalfDayThresholdHours,
		AllowSelfCheckout:             req.AllowSelfCheckout,
		RequireRegularizationApproval: req.RequireRegularizationApproval,
	}

	settings, err := h.service.UpdateSettings(c.Request.Context(), dto)
	if err != nil {
		logger.Error("Failed to update attendance settings",
			zap.String("tenant_id", tenantID.String()),
			zap.String("branch_id", branchID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to update attendance settings"))
		return
	}

	response.OK(c, toSettingsResponse(settings))
}

// Helper functions for response conversion

func toAttendanceResponse(a *models.StaffAttendance) AttendanceResponse {
	resp := AttendanceResponse{
		ID:             a.ID.String(),
		StaffID:        a.StaffID.String(),
		AttendanceDate: a.AttendanceDate.Format("2006-01-02"),
		Status:         string(a.Status),
		IsLate:         a.IsLate,
		LateMinutes:    a.LateMinutes,
		HalfDayType:    string(a.HalfDayType),
		Remarks:        a.Remarks,
		MarkedAt:       a.MarkedAt.Format(time.RFC3339),
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}

	if a.CheckInTime != nil {
		resp.CheckInTime = a.CheckInTime.Format(time.RFC3339)
	}
	if a.CheckOutTime != nil {
		resp.CheckOutTime = a.CheckOutTime.Format(time.RFC3339)
	}
	if a.MarkedBy != nil {
		resp.MarkedBy = a.MarkedBy.String()
	}

	// Add staff name if loaded
	if a.Staff.ID != uuid.Nil {
		resp.StaffName = a.Staff.FullName()
		resp.EmployeeID = a.Staff.EmployeeID
	}

	return resp
}

func toAttendanceResponses(attendanceList []models.StaffAttendance) []AttendanceResponse {
	responses := make([]AttendanceResponse, len(attendanceList))
	for i, a := range attendanceList {
		responses[i] = toAttendanceResponse(&a)
	}
	return responses
}

func toRegularizationResponse(r *models.StaffAttendanceRegularization) RegularizationResponse {
	resp := RegularizationResponse{
		ID:                    r.ID.String(),
		StaffID:               r.StaffID.String(),
		RequestDate:           r.RequestDate.Format("2006-01-02"),
		RequestedStatus:       string(r.RequestedStatus),
		Reason:                r.Reason,
		SupportingDocumentURL: r.SupportingDocumentURL,
		Status:                string(r.Status),
		RejectionReason:       r.RejectionReason,
		CreatedAt:             r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             r.UpdatedAt.Format(time.RFC3339),
	}

	if r.AttendanceID != nil {
		resp.AttendanceID = r.AttendanceID.String()
	}
	if r.ReviewedBy != nil {
		resp.ReviewedBy = r.ReviewedBy.String()
	}
	if r.ReviewedAt != nil {
		resp.ReviewedAt = r.ReviewedAt.Format(time.RFC3339)
	}

	// Add staff name if loaded
	if r.Staff.ID != uuid.Nil {
		resp.StaffName = r.Staff.FullName()
		resp.EmployeeID = r.Staff.EmployeeID
	}

	return resp
}

func toRegularizationResponses(regularizations []models.StaffAttendanceRegularization) []RegularizationResponse {
	responses := make([]RegularizationResponse, len(regularizations))
	for i, r := range regularizations {
		responses[i] = toRegularizationResponse(&r)
	}
	return responses
}

func toSettingsResponse(s *models.StaffAttendanceSettings) SettingsResponse {
	resp := SettingsResponse{
		ID:                            s.ID.String(),
		BranchID:                      s.BranchID.String(),
		WorkStartTime:                 s.WorkStartTime.Format("15:04"),
		WorkEndTime:                   s.WorkEndTime.Format("15:04"),
		LateThresholdMinutes:          s.LateThresholdMinutes,
		HalfDayThresholdHours:         s.HalfDayThresholdHours,
		AllowSelfCheckout:             s.AllowSelfCheckout,
		RequireRegularizationApproval: s.RequireRegularizationApproval,
		CreatedAt:                     s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                     s.UpdatedAt.Format(time.RFC3339),
	}

	// Add branch name if loaded
	if s.Branch.ID != uuid.Nil {
		resp.BranchName = s.Branch.Name
	}

	return resp
}
