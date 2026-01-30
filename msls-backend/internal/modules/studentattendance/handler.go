// Package studentattendance provides student attendance management functionality.
package studentattendance

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

// Handler handles student attendance-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new student attendance handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// MarkAttendanceRequest represents the request body for marking attendance.
type MarkAttendanceRequest struct {
	Date    string                    `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Records []AttendanceRecordRequest `json:"records" binding:"required,min=1"`
}

// AttendanceRecordRequest represents a single attendance record in the request.
type AttendanceRecordRequest struct {
	StudentID       string  `json:"studentId" binding:"required,uuid"`
	Status          string  `json:"status" binding:"required,oneof=present absent late half_day"`
	LateArrivalTime *string `json:"lateArrivalTime"` // Format: HH:MM
	Remarks         string  `json:"remarks" binding:"max=500"`
}

// SettingsRequest represents the request body for updating settings.
type SettingsRequest struct {
	BranchID             string `json:"branchId" binding:"required,uuid"`
	EditWindowMinutes    int    `json:"editWindowMinutes" binding:"min=0,max=1440"`
	LateThresholdMinutes int    `json:"lateThresholdMinutes" binding:"min=0,max=120"`
	SMSOnAbsent          bool   `json:"smsOnAbsent"`
}

// GetMyClasses returns the classes available for the teacher to mark attendance.
// @Summary Get teacher's classes
// @Description Get classes assigned to the teacher for attendance marking
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=[]TeacherClassResponse}
// @Router /api/v1/student-attendance/my-classes [get]
func (h *Handler) GetMyClasses(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Get branch ID from query or user's branch
	branchID := tenantID // Placeholder - in production, get from user's assignment

	// Parse date, default to today
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	classes, err := h.service.GetTeacherSections(c.Request.Context(), tenantID, branchID, date)
	if err != nil {
		logger.Error("Failed to get teacher classes",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve classes"))
		return
	}

	response.OK(c, classes)
}

// GetClassAttendance returns students and attendance for a class section.
// @Summary Get class attendance
// @Description Get students and their attendance status for a class section
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Section ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=ClassAttendanceResponse}
// @Router /api/v1/student-attendance/class/{id} [get]
func (h *Handler) GetClassAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date, default to today
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	// Prevent future dates
	if date.After(time.Now().Truncate(24 * time.Hour)) {
		apperrors.Abort(c, apperrors.BadRequest("Cannot view attendance for future dates"))
		return
	}

	attendance, err := h.service.GetClassStudentsForAttendance(c.Request.Context(), tenantID, sectionID, date)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrNoStudentsInSection):
			apperrors.Abort(c, apperrors.BadRequest("No students found in this section"))
		default:
			logger.Error("Failed to get class attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve class attendance"))
		}
		return
	}

	response.OK(c, attendance)
}

// MarkClassAttendance marks attendance for all students in a class section.
// @Summary Mark class attendance
// @Description Mark attendance for all students in a class section
// @Tags Student Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Section ID"
// @Param request body MarkAttendanceRequest true "Attendance records"
// @Success 201 {object} response.Success{data=MarkAttendanceResult}
// @Router /api/v1/student-attendance/class/{id} [post]
func (h *Handler) MarkClassAttendance(c *gin.Context) {
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

	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid date format, use YYYY-MM-DD"))
		return
	}

	// Convert request records to DTO
	records := make([]StudentAttendanceRecord, len(req.Records))
	for i, r := range req.Records {
		studentID, err := uuid.Parse(r.StudentID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid student ID: "+r.StudentID))
			return
		}

		record := StudentAttendanceRecord{
			StudentID: studentID,
			Status:    r.Status,
			Remarks:   r.Remarks,
		}

		// Parse late arrival time if provided
		if r.LateArrivalTime != nil && *r.LateArrivalTime != "" {
			lateTime, err := time.Parse("15:04", *r.LateArrivalTime)
			if err == nil {
				fullLateTime := time.Date(date.Year(), date.Month(), date.Day(),
					lateTime.Hour(), lateTime.Minute(), 0, 0, time.Local)
				record.LateArrivalTime = &fullLateTime
			}
		}

		records[i] = record
	}

	dto := MarkClassAttendanceDTO{
		TenantID:  tenantID,
		SectionID: sectionID,
		Date:      date,
		Records:   records,
		MarkedBy:  userID,
	}

	result, err := h.service.MarkClassAttendance(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrFutureDate):
			apperrors.Abort(c, apperrors.BadRequest("Cannot mark attendance for future date"))
		case errors.Is(err, ErrEditWindowExpired):
			apperrors.Abort(c, apperrors.Forbidden("Attendance edit window has expired"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid attendance status"))
		case errors.Is(err, ErrEmptyAttendanceRecords):
			apperrors.Abort(c, apperrors.BadRequest("No attendance records provided"))
		default:
			logger.Error("Failed to mark class attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to mark attendance"))
		}
		return
	}

	response.Created(c, result)
}

// ListAttendance retrieves attendance records with filters.
// @Summary List attendance records
// @Description Get attendance records with optional filters
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param section_id query string false "Filter by section ID"
// @Param student_id query string false "Filter by student ID"
// @Param status query string false "Filter by status"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} response.Success{data=AttendanceListResponse}
// @Router /api/v1/student-attendance [get]
func (h *Handler) ListAttendance(c *gin.Context) {
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
	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		sectionID, err := uuid.Parse(sectionIDStr)
		if err == nil {
			filter.SectionID = &sectionID
		}
	}
	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		studentID, err := uuid.Parse(studentIDStr)
		if err == nil {
			filter.StudentID = &studentID
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

// GetSettings retrieves student attendance settings for a branch.
// @Summary Get attendance settings
// @Description Get student attendance settings for a branch
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string true "Branch ID"
// @Success 200 {object} response.Success{data=SettingsResponse}
// @Router /api/v1/student-attendance/settings [get]
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
				BranchID:             branchID.String(),
				EditWindowMinutes:    120,
				LateThresholdMinutes: 15,
				SMSOnAbsent:          false,
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

// UpdateSettings updates student attendance settings for a branch.
// @Summary Update attendance settings
// @Description Update student attendance settings for a branch
// @Tags Student Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body SettingsRequest true "Settings details"
// @Success 200 {object} response.Success{data=SettingsResponse}
// @Router /api/v1/student-attendance/settings [put]
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

	dto := SettingsDTO{
		TenantID:             tenantID,
		BranchID:             branchID,
		EditWindowMinutes:    req.EditWindowMinutes,
		LateThresholdMinutes: req.LateThresholdMinutes,
		SMSOnAbsent:          req.SMSOnAbsent,
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

func toAttendanceResponse(a *models.StudentAttendance) AttendanceResponse {
	resp := AttendanceResponse{
		ID:             a.ID.String(),
		StudentID:      a.StudentID.String(),
		SectionID:      a.SectionID.String(),
		AttendanceDate: a.AttendanceDate.Format("2006-01-02"),
		Status:         string(a.Status),
		StatusLabel:    AttendanceStatus(a.Status).Label(),
		Remarks:        a.Remarks,
		MarkedBy:       a.MarkedBy.String(),
		MarkedAt:       a.MarkedAt.Format(time.RFC3339),
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}

	if a.LateArrivalTime != nil {
		resp.LateArrivalTime = a.LateArrivalTime.Format("15:04")
	}

	// Add student info if loaded
	if a.Student.ID != uuid.Nil {
		resp.StudentName = a.Student.FullName()
		resp.AdmissionNumber = a.Student.AdmissionNumber
	}

	// Add section info if loaded
	if a.Section.ID != uuid.Nil {
		resp.SectionName = a.Section.Name
	}

	// Add marker info if loaded
	if a.MarkedByUser.ID != uuid.Nil {
		resp.MarkedByName = a.MarkedByUser.FirstName + " " + a.MarkedByUser.LastName
	}

	return resp
}

func toAttendanceResponses(attendanceList []models.StudentAttendance) []AttendanceResponse {
	responses := make([]AttendanceResponse, len(attendanceList))
	for i, a := range attendanceList {
		responses[i] = toAttendanceResponse(&a)
	}
	return responses
}

func toSettingsResponse(s *models.StudentAttendanceSettings) SettingsResponse {
	resp := SettingsResponse{
		ID:                   s.ID.String(),
		BranchID:             s.BranchID.String(),
		EditWindowMinutes:    s.EditWindowMinutes,
		LateThresholdMinutes: s.LateThresholdMinutes,
		SMSOnAbsent:          s.SMSOnAbsent,
		CreatedAt:            s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            s.UpdatedAt.Format(time.RFC3339),
	}

	// Add branch name if loaded
	if s.Branch.ID != uuid.Nil {
		resp.BranchName = s.Branch.Name
	}

	return resp
}

// ============================================================================
// Period-wise Attendance Handlers (Story 7.2)
// ============================================================================

// PeriodAttendanceRecordRequest represents a single student's period attendance in a request.
type PeriodAttendanceRecordRequest struct {
	StudentID       string  `json:"studentId" binding:"required,uuid"`
	Status          string  `json:"status" binding:"required,oneof=present absent late half_day"`
	LateArrivalTime *string `json:"lateArrivalTime"` // Format: HH:MM
	Remarks         string  `json:"remarks" binding:"max=500"`
}

// MarkPeriodAttendanceRequestBody represents the request body for marking period attendance.
type MarkPeriodAttendanceRequestBody struct {
	SectionID string                          `json:"sectionId" binding:"required,uuid"`
	Date      string                          `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Records   []PeriodAttendanceRecordRequest `json:"records" binding:"required,min=1"`
}

// GetPeriods returns the periods available for a section on a given date.
// @Summary Get periods for section
// @Description Get teaching periods available for marking attendance on a specific date
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param section_id query string true "Section ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=TeacherPeriodsResponse}
// @Router /api/v1/student-attendance/periods [get]
func (h *Handler) GetPeriods(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionIDStr := c.Query("section_id")
	if sectionIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Section ID is required"))
		return
	}

	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date, default to today
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	periods, err := h.service.GetTeacherPeriods(c.Request.Context(), tenantID, sectionID, date)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrNoTimetableForSection):
			apperrors.Abort(c, apperrors.NotFound("No timetable found for this section"))
		default:
			logger.Error("Failed to get periods",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve periods"))
		}
		return
	}

	response.OK(c, periods)
}

// GetPeriodAttendance returns students with their attendance for a specific period.
// @Summary Get period attendance
// @Description Get students and their attendance status for a specific period
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Period Slot ID"
// @Param section_id query string true "Section ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=PeriodAttendanceResponse}
// @Router /api/v1/student-attendance/period/{id} [get]
func (h *Handler) GetPeriodAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	periodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid period ID"))
		return
	}

	sectionIDStr := c.Query("section_id")
	if sectionIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Section ID is required"))
		return
	}

	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date, default to today
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	// Prevent future dates
	if date.After(time.Now().Truncate(24 * time.Hour)) {
		apperrors.Abort(c, apperrors.BadRequest("Cannot view attendance for future dates"))
		return
	}

	attendance, err := h.service.GetPeriodAttendance(c.Request.Context(), tenantID, sectionID, periodID, date)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrPeriodNotFound):
			apperrors.Abort(c, apperrors.NotFound("Period not found"))
		case errors.Is(err, ErrNoStudentsInSection):
			apperrors.Abort(c, apperrors.BadRequest("No students found in this section"))
		case errors.Is(err, ErrInvalidPeriodForSection):
			apperrors.Abort(c, apperrors.BadRequest("Period does not belong to this section's timetable"))
		default:
			logger.Error("Failed to get period attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.String("period_id", periodID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve period attendance"))
		}
		return
	}

	response.OK(c, attendance)
}

// MarkPeriodAttendance marks attendance for a specific period.
// @Summary Mark period attendance
// @Description Mark attendance for all students in a specific period
// @Tags Student Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Period Slot ID"
// @Param request body MarkPeriodAttendanceRequestBody true "Period attendance records"
// @Success 201 {object} response.Success{data=MarkPeriodAttendanceResult}
// @Router /api/v1/student-attendance/period/{id} [post]
func (h *Handler) MarkPeriodAttendance(c *gin.Context) {
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

	periodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid period ID"))
		return
	}

	var req MarkPeriodAttendanceRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	sectionID, err := uuid.Parse(req.SectionID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid date format, use YYYY-MM-DD"))
		return
	}

	// Convert request records to DTO
	records := make([]PeriodAttendanceRecord, len(req.Records))
	for i, r := range req.Records {
		studentID, err := uuid.Parse(r.StudentID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid student ID: "+r.StudentID))
			return
		}

		record := PeriodAttendanceRecord{
			StudentID: studentID,
			Status:    r.Status,
			Remarks:   r.Remarks,
		}

		// Parse late arrival time if provided
		if r.LateArrivalTime != nil && *r.LateArrivalTime != "" {
			lateTime, err := time.Parse("15:04", *r.LateArrivalTime)
			if err == nil {
				fullLateTime := time.Date(date.Year(), date.Month(), date.Day(),
					lateTime.Hour(), lateTime.Minute(), 0, 0, time.Local)
				record.LateArrivalTime = &fullLateTime
			}
		}

		records[i] = record
	}

	// Get timetable entry ID for this period
	// We need to fetch it from the timetable entries
	entries, err := h.service.repo.GetTimetableEntriesForSection(c.Request.Context(), tenantID, sectionID, int(date.Weekday()))
	if err != nil {
		logger.Error("Failed to get timetable entries",
			zap.String("tenant_id", tenantID.String()),
			zap.String("section_id", sectionID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve timetable"))
		return
	}

	var timetableEntryID uuid.UUID
	for _, entry := range entries {
		if entry.PeriodSlotID == periodID {
			timetableEntryID = entry.ID
			break
		}
	}

	if timetableEntryID == uuid.Nil {
		apperrors.Abort(c, apperrors.BadRequest("Period does not belong to this section's timetable"))
		return
	}

	dto := MarkPeriodAttendanceDTO{
		TenantID:         tenantID,
		SectionID:        sectionID,
		PeriodID:         periodID,
		TimetableEntryID: timetableEntryID,
		Date:             date,
		Records:          records,
		MarkedBy:         userID,
	}

	result, err := h.service.MarkPeriodAttendance(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrPeriodNotFound):
			apperrors.Abort(c, apperrors.NotFound("Period not found"))
		case errors.Is(err, ErrFutureDate):
			apperrors.Abort(c, apperrors.BadRequest("Cannot mark attendance for future date"))
		case errors.Is(err, ErrEditWindowExpired):
			apperrors.Abort(c, apperrors.Forbidden("Attendance edit window has expired"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid attendance status"))
		case errors.Is(err, ErrEmptyAttendanceRecords):
			apperrors.Abort(c, apperrors.BadRequest("No attendance records provided"))
		default:
			logger.Error("Failed to mark period attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.String("period_id", periodID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to mark attendance"))
		}
		return
	}

	response.Created(c, result)
}

// GetDailySummary returns aggregated attendance for all periods in a day.
// @Summary Get daily summary
// @Description Get aggregated attendance summary for all periods in a day
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param section_id query string true "Section ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=DailySummaryResponse}
// @Router /api/v1/student-attendance/daily-summary [get]
func (h *Handler) GetDailySummary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionIDStr := c.Query("section_id")
	if sectionIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Section ID is required"))
		return
	}

	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date, default to today
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	summary, err := h.service.GetDailySummary(c.Request.Context(), tenantID, sectionID, date)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		case errors.Is(err, ErrNoTimetableForSection):
			apperrors.Abort(c, apperrors.NotFound("No timetable found for this section"))
		default:
			logger.Error("Failed to get daily summary",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve daily summary"))
		}
		return
	}

	response.OK(c, summary)
}

// GetSubjectAttendance returns attendance statistics for a student by subject.
// @Summary Get subject attendance
// @Description Get attendance statistics for a student by subject
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Subject ID"
// @Param student_id query string true "Student ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Success{data=SubjectAttendanceResponse}
// @Router /api/v1/student-attendance/subject/{id} [get]
func (h *Handler) GetSubjectAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	subjectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid subject ID"))
		return
	}

	studentIDStr := c.Query("student_id")
	if studentIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Student ID is required"))
		return
	}

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	// Parse optional date range
	var dateFrom, dateTo *time.Time
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		parsed, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			dateFrom = &parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		parsed, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			dateTo = &parsed
		}
	}

	stats, err := h.service.GetSubjectAttendance(c.Request.Context(), tenantID, studentID, subjectID, dateFrom, dateTo)
	if err != nil {
		logger.Error("Failed to get subject attendance",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.String("subject_id", subjectID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve subject attendance"))
		return
	}

	response.OK(c, stats)
}

// ============================================================================
// Attendance Edit & Audit Handlers (Story 7.3)
// ============================================================================

// EditAttendance updates an attendance record with audit trail.
// @Summary Edit attendance
// @Description Edit an existing attendance record with reason for change
// @Tags Student Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Attendance ID"
// @Param request body EditAttendanceRequest true "Edit details with reason"
// @Success 200 {object} response.Success{data=EditAttendanceResult}
// @Router /api/v1/student-attendance/{id} [put]
func (h *Handler) EditAttendance(c *gin.Context) {
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

	attendanceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid attendance ID"))
		return
	}

	var req EditAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Check if user is admin (simplified check - in production, check roles)
	isAdmin := middleware.HasPermission(c, "attendance:admin")

	// Build DTO
	dto := EditAttendanceDTO{
		TenantID:     tenantID,
		AttendanceID: attendanceID,
		Status:       req.Status,
		Reason:       req.Reason,
		EditedBy:     userID,
		IsAdmin:      isAdmin,
	}

	if req.Remarks != "" {
		dto.Remarks = &req.Remarks
	}

	// Parse late arrival time if provided
	if req.LateArrivalTime != "" {
		lateTime, err := time.Parse("15:04", req.LateArrivalTime)
		if err == nil {
			dto.LateArrivalTime = &lateTime
		}
	}

	result, err := h.service.EditAttendance(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrAttendanceNotFound):
			apperrors.Abort(c, apperrors.NotFound("Attendance record not found"))
		case errors.Is(err, ErrEditReasonRequired):
			apperrors.Abort(c, apperrors.BadRequest("Edit reason is required"))
		case errors.Is(err, ErrNotOriginalMarker):
			apperrors.Abort(c, apperrors.Forbidden("Only the original marker can edit within edit window"))
		case errors.Is(err, ErrAdminOnlyEdit):
			apperrors.Abort(c, apperrors.Forbidden("Edit window expired - admin approval required"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid attendance status"))
		default:
			logger.Error("Failed to edit attendance",
				zap.String("tenant_id", tenantID.String()),
				zap.String("attendance_id", attendanceID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to edit attendance"))
		}
		return
	}

	response.OK(c, result)
}

// GetAttendanceHistory returns the audit trail for an attendance record.
// @Summary Get attendance history
// @Description Get the complete edit history (audit trail) for an attendance record
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Attendance ID"
// @Success 200 {object} response.Success{data=AttendanceAuditTrailResponse}
// @Router /api/v1/student-attendance/{id}/history [get]
func (h *Handler) GetAttendanceHistory(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	attendanceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid attendance ID"))
		return
	}

	auditTrail, err := h.service.GetAttendanceAuditTrail(c.Request.Context(), tenantID, attendanceID)
	if err != nil {
		switch {
		case errors.Is(err, ErrAttendanceNotFound):
			apperrors.Abort(c, apperrors.NotFound("Attendance record not found"))
		default:
			logger.Error("Failed to get attendance history",
				zap.String("tenant_id", tenantID.String()),
				zap.String("attendance_id", attendanceID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve attendance history"))
		}
		return
	}

	response.OK(c, auditTrail)
}

// GetEditWindowStatus returns the edit window status for an attendance record.
// @Summary Get edit window status
// @Description Get the current edit window status and whether the user can edit
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Attendance ID"
// @Success 200 {object} response.Success{data=EditWindowStatusResponse}
// @Router /api/v1/student-attendance/{id}/edit-status [get]
func (h *Handler) GetEditWindowStatus(c *gin.Context) {
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

	attendanceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid attendance ID"))
		return
	}

	// Check if user is admin
	isAdmin := middleware.HasPermission(c, "attendance:admin")

	status, err := h.service.GetEditWindowStatus(c.Request.Context(), tenantID, attendanceID, userID, isAdmin)
	if err != nil {
		switch {
		case errors.Is(err, ErrAttendanceNotFound):
			apperrors.Abort(c, apperrors.NotFound("Attendance record not found"))
		default:
			logger.Error("Failed to get edit window status",
				zap.String("tenant_id", tenantID.String()),
				zap.String("attendance_id", attendanceID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve edit window status"))
		}
		return
	}

	response.OK(c, status)
}

// ============================================================================
// Calendar & Reports Handlers (Stories 7.4-7.8)
// ============================================================================

// GetStudentCalendar returns a student's monthly attendance calendar.
// @Summary Get student attendance calendar
// @Description Get monthly attendance calendar for a student
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID"
// @Param year query int false "Year (defaults to current year)"
// @Param month query int false "Month (1-12, defaults to current month)"
// @Success 200 {object} response.Success{data=MonthlyCalendarResponse}
// @Router /api/v1/student-attendance/calendar/{studentId} [get]
func (h *Handler) GetStudentCalendar(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	// Parse year and month, default to current
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

	calendar, err := h.service.GetStudentCalendar(c.Request.Context(), tenantID, studentID, year, month)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		default:
			logger.Error("Failed to get student calendar",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve calendar"))
		}
		return
	}

	response.OK(c, calendar)
}

// GetStudentSummary returns a student's attendance summary with trend.
// @Summary Get student attendance summary
// @Description Get attendance summary with trend for a student
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Success{data=StudentSummaryResponse}
// @Router /api/v1/student-attendance/summary/{studentId} [get]
func (h *Handler) GetStudentSummary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	// Parse date range, default to current month
	now := time.Now()
	dateFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	dateTo := now

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = parsed
		}
	}

	summary, err := h.service.GetStudentSummary(c.Request.Context(), tenantID, studentID, dateFrom, dateTo)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		default:
			logger.Error("Failed to get student summary",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve summary"))
		}
		return
	}

	response.OK(c, summary)
}

// GetClassReport returns class-level attendance report for a date.
// @Summary Get class attendance report
// @Description Get attendance report for a class section on a specific date
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sectionId path string true "Section ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=ClassReportResponse}
// @Router /api/v1/student-attendance/reports/class/{sectionId} [get]
func (h *Handler) GetClassReport(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionID, err := uuid.Parse(c.Param("sectionId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse date, default to today
	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	report, err := h.service.GetClassReport(c.Request.Context(), tenantID, sectionID, date)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		default:
			logger.Error("Failed to get class report",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve class report"))
		}
		return
	}

	response.OK(c, report)
}

// GetMonthlyClassReport returns monthly class attendance report.
// @Summary Get monthly class attendance report
// @Description Get monthly attendance report for a class section in grid format
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sectionId path string true "Section ID"
// @Param year query int false "Year (defaults to current year)"
// @Param month query int false "Month (1-12, defaults to current month)"
// @Success 200 {object} response.Success{data=MonthlyClassReportResponse}
// @Router /api/v1/student-attendance/reports/class/{sectionId}/monthly [get]
func (h *Handler) GetMonthlyClassReport(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionID, err := uuid.Parse(c.Param("sectionId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	// Parse year and month
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

	report, err := h.service.GetMonthlyClassReport(c.Request.Context(), tenantID, sectionID, year, month)
	if err != nil {
		switch {
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
		default:
			logger.Error("Failed to get monthly class report",
				zap.String("tenant_id", tenantID.String()),
				zap.String("section_id", sectionID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve monthly report"))
		}
		return
	}

	response.OK(c, report)
}

// GetLowAttendanceDashboard returns dashboard data for low attendance students.
// @Summary Get low attendance dashboard
// @Description Get dashboard showing students with low attendance
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param threshold query number false "Attendance threshold percentage (default 75)"
// @Success 200 {object} response.Success{data=LowAttendanceDashboardResponse}
// @Router /api/v1/student-attendance/alerts/low-attendance [get]
func (h *Handler) GetLowAttendanceDashboard(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Parse date range
	now := time.Now()
	dateFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	dateTo := now

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = parsed
		}
	}

	// Parse thresholds
	threshold := 75.0
	criticalThreshold := 60.0
	if thresholdStr := c.Query("threshold"); thresholdStr != "" {
		if t, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = t
		}
	}

	dashboard, err := h.service.GetLowAttendanceDashboard(c.Request.Context(), tenantID, dateFrom, dateTo, threshold, criticalThreshold)
	if err != nil {
		logger.Error("Failed to get low attendance dashboard",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve dashboard"))
		return
	}

	response.OK(c, dashboard)
}

// GetUnmarkedAttendance returns classes with unmarked attendance.
// @Summary Get unmarked attendance
// @Description Get classes that haven't marked attendance for a date
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=UnmarkedAttendanceResponse}
// @Router /api/v1/student-attendance/alerts/unmarked [get]
func (h *Handler) GetUnmarkedAttendance(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Parse date
	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	// Default deadline time
	deadlineTime := "10:00"

	unmarked, err := h.service.GetUnmarkedAttendance(c.Request.Context(), tenantID, date, deadlineTime)
	if err != nil {
		logger.Error("Failed to get unmarked attendance",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve unmarked attendance"))
		return
	}

	response.OK(c, unmarked)
}

// GetDailyReport returns end-of-day attendance summary report.
// @Summary Get daily attendance report
// @Description Get end-of-day attendance summary report
// @Tags Student Attendance
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} response.Success{data=DailyReportSummaryResponse}
// @Router /api/v1/student-attendance/reports/daily [get]
func (h *Handler) GetDailyReport(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Parse date
	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	report, err := h.service.GetDailyReport(c.Request.Context(), tenantID, date)
	if err != nil {
		logger.Error("Failed to get daily report",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve daily report"))
		return
	}

	response.OK(c, report)
}
