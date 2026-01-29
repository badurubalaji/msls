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
