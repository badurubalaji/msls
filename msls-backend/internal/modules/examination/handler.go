package examination

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
)

// Handler handles HTTP requests for examinations
type Handler struct {
	service *Service
}

// NewHandler creates a new examination handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers examination routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	examinations := rg.Group("/examinations")
	{
		examinations.GET("", h.ListExaminations)
		examinations.GET("/:id", h.GetExamination)
		examinations.POST("", h.CreateExamination)
		examinations.PUT("/:id", h.UpdateExamination)
		examinations.DELETE("/:id", h.DeleteExamination)
		examinations.POST("/:id/publish", h.PublishExamination)
		examinations.POST("/:id/unpublish", h.UnPublishExamination)

		// Schedule routes
		examinations.GET("/:id/schedules", h.GetSchedules)
		examinations.POST("/:id/schedules", h.CreateSchedule)
		examinations.PUT("/:id/schedules/:scheduleId", h.UpdateSchedule)
		examinations.DELETE("/:id/schedules/:scheduleId", h.DeleteSchedule)
	}
}

// ListExaminations godoc
// @Summary List examinations
// @Tags Examinations
// @Produce json
// @Param academicYearId query string false "Filter by academic year"
// @Param examTypeId query string false "Filter by exam type"
// @Param classId query string false "Filter by class"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by name"
// @Success 200 {object} response.Response{data=[]ExaminationResponse}
// @Router /examinations [get]
func (h *Handler) ListExaminations(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var filter ExaminationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	exams, err := h.service.List(tenantID, filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list examinations"))
		return
	}

	response.OK(c, ToResponseList(exams))
}

// GetExamination godoc
// @Summary Get examination by ID
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} response.Response{data=ExaminationResponse}
// @Router /examinations/{id} [get]
func (h *Handler) GetExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	exam, err := h.service.GetByID(tenantID, id)
	if err != nil {
		if errors.Is(err, ErrExaminationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Examination not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get examination"))
		return
	}

	response.OK(c, ToResponse(exam))
}

// CreateExamination godoc
// @Summary Create examination
// @Tags Examinations
// @Accept json
// @Produce json
// @Param examination body CreateExaminationRequest true "Examination data"
// @Success 201 {object} response.Response{data=ExaminationResponse}
// @Router /examinations [post]
func (h *Handler) CreateExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User ID is required"))
		return
	}

	var req CreateExaminationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	exam, err := h.service.Create(tenantID, req, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Created(c, ToResponse(exam))
}

// UpdateExamination godoc
// @Summary Update examination
// @Tags Examinations
// @Accept json
// @Produce json
// @Param id path string true "Examination ID"
// @Param examination body UpdateExaminationRequest true "Examination data"
// @Success 200 {object} response.Response{data=ExaminationResponse}
// @Router /examinations/{id} [put]
func (h *Handler) UpdateExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	var req UpdateExaminationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	exam, err := h.service.Update(tenantID, id, req, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, ToResponse(exam))
}

// DeleteExamination godoc
// @Summary Delete examination
// @Tags Examinations
// @Param id path string true "Examination ID"
// @Success 204
// @Router /examinations/{id} [delete]
func (h *Handler) DeleteExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	if err := h.service.Delete(tenantID, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// PublishExamination godoc
// @Summary Publish examination
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} response.Response{data=ExaminationResponse}
// @Router /examinations/{id}/publish [post]
func (h *Handler) PublishExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	exam, err := h.service.Publish(tenantID, id, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, ToResponse(exam))
}

// UnPublishExamination godoc
// @Summary Unpublish examination (revert to draft)
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} response.Response{data=ExaminationResponse}
// @Router /examinations/{id}/unpublish [post]
func (h *Handler) UnPublishExamination(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("User ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	exam, err := h.service.UnPublish(tenantID, id, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, ToResponse(exam))
}

// ========================================
// Schedule Handlers
// ========================================

// GetSchedules godoc
// @Summary Get examination schedules
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} response.Response{data=[]ExamScheduleResponse}
// @Router /examinations/{id}/schedules [get]
func (h *Handler) GetSchedules(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	schedules, err := h.service.GetSchedules(tenantID, examID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Map to response
	responses := make([]ExamScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = ExamScheduleResponse{
			ID:           schedule.ID,
			SubjectID:    schedule.SubjectID,
			ExamDate:     schedule.ExamDate.Format("2006-01-02"),
			StartTime:    schedule.StartTime,
			EndTime:      schedule.EndTime,
			MaxMarks:     schedule.MaxMarks,
			PassingMarks: schedule.PassingMarks,
			Venue:        schedule.Venue,
			Notes:        schedule.Notes,
		}
		if schedule.Subject != nil {
			responses[i].SubjectName = schedule.Subject.Name
			responses[i].SubjectCode = schedule.Subject.Code
		}
	}

	response.OK(c, responses)
}

// CreateSchedule godoc
// @Summary Create exam schedule
// @Tags Examinations
// @Accept json
// @Produce json
// @Param id path string true "Examination ID"
// @Param schedule body CreateScheduleRequest true "Schedule data"
// @Success 201 {object} response.Response{data=ExamScheduleResponse}
// @Router /examinations/{id}/schedules [post]
func (h *Handler) CreateSchedule(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	schedule, err := h.service.CreateSchedule(tenantID, examID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	resp := ExamScheduleResponse{
		ID:           schedule.ID,
		SubjectID:    schedule.SubjectID,
		ExamDate:     schedule.ExamDate.Format("2006-01-02"),
		StartTime:    schedule.StartTime,
		EndTime:      schedule.EndTime,
		MaxMarks:     schedule.MaxMarks,
		PassingMarks: schedule.PassingMarks,
		Venue:        schedule.Venue,
		Notes:        schedule.Notes,
	}
	if schedule.Subject != nil {
		resp.SubjectName = schedule.Subject.Name
		resp.SubjectCode = schedule.Subject.Code
	}

	response.Created(c, resp)
}

// UpdateSchedule godoc
// @Summary Update exam schedule
// @Tags Examinations
// @Accept json
// @Produce json
// @Param id path string true "Examination ID"
// @Param scheduleId path string true "Schedule ID"
// @Param schedule body UpdateScheduleRequest true "Schedule data"
// @Success 200 {object} response.Response{data=ExamScheduleResponse}
// @Router /examinations/{id}/schedules/{scheduleId} [put]
func (h *Handler) UpdateSchedule(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	scheduleID, err := uuid.Parse(c.Param("scheduleId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid schedule ID"))
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	schedule, err := h.service.UpdateSchedule(tenantID, examID, scheduleID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	resp := ExamScheduleResponse{
		ID:           schedule.ID,
		SubjectID:    schedule.SubjectID,
		ExamDate:     schedule.ExamDate.Format("2006-01-02"),
		StartTime:    schedule.StartTime,
		EndTime:      schedule.EndTime,
		MaxMarks:     schedule.MaxMarks,
		PassingMarks: schedule.PassingMarks,
		Venue:        schedule.Venue,
		Notes:        schedule.Notes,
	}
	if schedule.Subject != nil {
		resp.SubjectName = schedule.Subject.Name
		resp.SubjectCode = schedule.Subject.Code
	}

	response.OK(c, resp)
}

// DeleteSchedule godoc
// @Summary Delete exam schedule
// @Tags Examinations
// @Param id path string true "Examination ID"
// @Param scheduleId path string true "Schedule ID"
// @Success 204
// @Router /examinations/{id}/schedules/{scheduleId} [delete]
func (h *Handler) DeleteSchedule(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid examination ID"))
		return
	}

	scheduleID, err := uuid.Parse(c.Param("scheduleId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid schedule ID"))
		return
	}

	if err := h.service.DeleteSchedule(tenantID, examID, scheduleID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleServiceError maps service errors to appropriate HTTP responses
func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrExaminationNotFound):
		apperrors.Abort(c, apperrors.NotFound("Examination not found"))
	case errors.Is(err, ErrScheduleNotFound):
		apperrors.Abort(c, apperrors.NotFound("Schedule not found"))
	case errors.Is(err, ErrInvalidDateRange):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrInvalidTimeRange):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrInvalidPassingMarks):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrCannotUpdatePublished):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrCannotDeletePublished):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrInvalidStatusTransition):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrNoSchedulesForPublish):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrScheduleOutsideExamDates):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrSubjectAlreadyScheduled):
		apperrors.Abort(c, apperrors.Conflict(err.Error()))
	case errors.Is(err, ErrNoClassesSpecified):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	case errors.Is(err, ErrExamTypeNotActive):
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
	default:
		apperrors.Abort(c, apperrors.InternalError("An error occurred"))
	}
}
