// Package assignment provides teacher subject assignment functionality.
package assignment

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

// Handler handles assignment-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new assignment handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateAssignmentRequest represents the request body for creating an assignment.
type CreateAssignmentRequest struct {
	StaffID        string  `json:"staffId" binding:"required,uuid"`
	SubjectID      string  `json:"subjectId" binding:"required,uuid"`
	ClassID        string  `json:"classId" binding:"required,uuid"`
	SectionID      *string `json:"sectionId" binding:"omitempty,uuid"`
	AcademicYearID string  `json:"academicYearId" binding:"required,uuid"`
	PeriodsPerWeek int     `json:"periodsPerWeek" binding:"min=0,max=50"`
	IsClassTeacher bool    `json:"isClassTeacher"`
	EffectiveFrom  string  `json:"effectiveFrom" binding:"required"` // Format: YYYY-MM-DD
	EffectiveTo    *string `json:"effectiveTo"`                      // Format: YYYY-MM-DD
	Remarks        string  `json:"remarks" binding:"max=500"`
}

// UpdateAssignmentRequest represents the request body for updating an assignment.
type UpdateAssignmentRequest struct {
	PeriodsPerWeek *int    `json:"periodsPerWeek" binding:"omitempty,min=0,max=50"`
	IsClassTeacher *bool   `json:"isClassTeacher"`
	EffectiveFrom  *string `json:"effectiveFrom"`
	EffectiveTo    *string `json:"effectiveTo"`
	Status         *string `json:"status" binding:"omitempty,oneof=active inactive"`
	Remarks        *string `json:"remarks" binding:"omitempty,max=500"`
}

// BulkCreateRequest represents the request body for bulk creating assignments.
type BulkCreateRequest struct {
	Assignments []BulkAssignmentItemRequest `json:"assignments" binding:"required,min=1,max=100,dive"`
}

// BulkAssignmentItemRequest represents a single assignment in a bulk create request.
type BulkAssignmentItemRequest struct {
	StaffID        string  `json:"staffId" binding:"required,uuid"`
	SubjectID      string  `json:"subjectId" binding:"required,uuid"`
	ClassID        string  `json:"classId" binding:"required,uuid"`
	SectionID      *string `json:"sectionId" binding:"omitempty,uuid"`
	AcademicYearID string  `json:"academicYearId" binding:"required,uuid"`
	PeriodsPerWeek int     `json:"periodsPerWeek" binding:"min=0,max=50"`
	IsClassTeacher bool    `json:"isClassTeacher"`
	EffectiveFrom  string  `json:"effectiveFrom" binding:"required"`
}

// SetClassTeacherRequest represents the request body for setting a class teacher.
type SetClassTeacherRequest struct {
	StaffID string `json:"staffId" binding:"required,uuid"`
}

// WorkloadSettingsRequest represents the request body for updating workload settings.
type WorkloadSettingsRequest struct {
	MinPeriodsPerWeek     int  `json:"minPeriodsPerWeek" binding:"min=0,max=50"`
	MaxPeriodsPerWeek     int  `json:"maxPeriodsPerWeek" binding:"min=1,max=60"`
	MaxSubjectsPerTeacher *int `json:"maxSubjectsPerTeacher" binding:"omitempty,min=1,max=20"`
	MaxClassesPerTeacher  *int `json:"maxClassesPerTeacher" binding:"omitempty,min=1,max=20"`
}

// List returns all assignments for the tenant.
// @Summary List teacher assignments
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param staff_id query string false "Filter by staff ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param class_id query string false "Filter by class ID"
// @Param section_id query string false "Filter by section ID"
// @Param academic_year_id query string false "Filter by academic year ID"
// @Param is_class_teacher query boolean false "Filter by class teacher status"
// @Param status query string false "Filter by status (active, inactive)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} response.Success{data=AssignmentListResponse}
// @Router /api/v1/teacher-assignments [get]
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

	// Parse filters
	if staffIDStr := c.Query("staff_id"); staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
			return
		}
		filter.StaffID = &staffID
	}

	if subjectIDStr := c.Query("subject_id"); subjectIDStr != "" {
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid subject ID"))
			return
		}
		filter.SubjectID = &subjectID
	}

	if classIDStr := c.Query("class_id"); classIDStr != "" {
		classID, err := uuid.Parse(classIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
			return
		}
		filter.ClassID = &classID
	}

	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		sectionID, err := uuid.Parse(sectionIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
			return
		}
		filter.SectionID = &sectionID
	}

	if academicYearIDStr := c.Query("academic_year_id"); academicYearIDStr != "" {
		academicYearID, err := uuid.Parse(academicYearIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
			return
		}
		filter.AcademicYearID = &academicYearID
	}

	if isClassTeacherStr := c.Query("is_class_teacher"); isClassTeacherStr != "" {
		isClassTeacher := isClassTeacherStr == "true"
		filter.IsClassTeacher = &isClassTeacher
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.AssignmentStatus(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	assignments, nextCursor, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list assignments",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve assignments"))
		return
	}

	resp := AssignmentListResponse{
		Assignments: ToAssignmentResponses(assignments),
		NextCursor:  nextCursor,
		HasMore:     nextCursor != "",
		Total:       total,
	}

	response.OK(c, resp)
}

// Get retrieves an assignment by ID.
// @Summary Get assignment by ID
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Assignment ID"
// @Success 200 {object} response.Success{data=AssignmentResponse}
// @Router /api/v1/teacher-assignments/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid assignment ID"))
		return
	}

	assignment, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrAssignmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Assignment not found"))
			return
		}
		logger.Error("Failed to get assignment",
			zap.String("tenant_id", tenantID.String()),
			zap.String("assignment_id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve assignment"))
		return
	}

	response.OK(c, ToAssignmentResponse(assignment))
}

// Create creates a new assignment.
// @Summary Create teacher assignment
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateAssignmentRequest true "Assignment details"
// @Success 201 {object} response.Success{data=AssignmentResponse}
// @Router /api/v1/teacher-assignments [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, _ := uuid.Parse(req.StaffID)
	subjectID, _ := uuid.Parse(req.SubjectID)
	classID, _ := uuid.Parse(req.ClassID)
	academicYearID, _ := uuid.Parse(req.AcademicYearID)

	effectiveFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid effective from date format, use YYYY-MM-DD"))
		return
	}

	dto := CreateAssignmentDTO{
		TenantID:       tenantID,
		StaffID:        staffID,
		SubjectID:      subjectID,
		ClassID:        classID,
		AcademicYearID: academicYearID,
		PeriodsPerWeek: req.PeriodsPerWeek,
		IsClassTeacher: req.IsClassTeacher,
		EffectiveFrom:  effectiveFrom,
		Remarks:        req.Remarks,
		CreatedBy:      &userID,
	}

	if req.SectionID != nil {
		sectionID, err := uuid.Parse(*req.SectionID)
		if err == nil {
			dto.SectionID = &sectionID
		}
	}

	if req.EffectiveTo != nil {
		effectiveTo, err := time.Parse("2006-01-02", *req.EffectiveTo)
		if err == nil {
			dto.EffectiveTo = &effectiveTo
		}
	}

	assignment, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStaffNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Staff member not found"))
		case errors.Is(err, ErrSubjectNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Subject not found"))
		case errors.Is(err, ErrClassNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Class not found"))
		case errors.Is(err, ErrSectionNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Section not found"))
		case errors.Is(err, ErrAcademicYearNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Academic year not found"))
		case errors.Is(err, ErrDuplicateAssignment):
			apperrors.Abort(c, apperrors.Conflict("Assignment already exists for this teacher-subject-class combination"))
		case errors.Is(err, ErrTeacherOverAssigned):
			apperrors.Abort(c, apperrors.BadRequest("Teacher would exceed maximum periods per week"))
		case errors.Is(err, ErrInvalidEffectiveDate):
			apperrors.Abort(c, apperrors.BadRequest("Effective from date cannot be after effective to date"))
		default:
			logger.Error("Failed to create assignment",
				zap.String("tenant_id", tenantID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to create assignment"))
		}
		return
	}

	response.Created(c, ToAssignmentResponse(assignment))
}

// Update updates an assignment.
// @Summary Update teacher assignment
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Assignment ID"
// @Param request body UpdateAssignmentRequest true "Assignment details"
// @Success 200 {object} response.Success{data=AssignmentResponse}
// @Router /api/v1/teacher-assignments/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid assignment ID"))
		return
	}

	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateAssignmentDTO{
		PeriodsPerWeek: req.PeriodsPerWeek,
		IsClassTeacher: req.IsClassTeacher,
		Remarks:        req.Remarks,
	}

	if req.EffectiveFrom != nil {
		effectiveFrom, err := time.Parse("2006-01-02", *req.EffectiveFrom)
		if err == nil {
			dto.EffectiveFrom = &effectiveFrom
		}
	}

	if req.EffectiveTo != nil {
		effectiveTo, err := time.Parse("2006-01-02", *req.EffectiveTo)
		if err == nil {
			dto.EffectiveTo = &effectiveTo
		}
	}

	if req.Status != nil {
		status := models.AssignmentStatus(*req.Status)
		dto.Status = &status
	}

	assignment, err := h.service.Update(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrAssignmentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Assignment not found"))
		case errors.Is(err, ErrTeacherOverAssigned):
			apperrors.Abort(c, apperrors.BadRequest("Teacher would exceed maximum periods per week"))
		case errors.Is(err, ErrInvalidEffectiveDate):
			apperrors.Abort(c, apperrors.BadRequest("Effective from date cannot be after effective to date"))
		case errors.Is(err, ErrCannotModifyInactiveAssignment):
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify inactive assignment"))
		default:
			logger.Error("Failed to update assignment",
				zap.String("tenant_id", tenantID.String()),
				zap.String("assignment_id", id.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update assignment"))
		}
		return
	}

	response.OK(c, ToAssignmentResponse(assignment))
}

// Delete deletes an assignment.
// @Summary Delete teacher assignment
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Assignment ID"
// @Success 204 "No Content"
// @Router /api/v1/teacher-assignments/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid assignment ID"))
		return
	}

	err = h.service.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrAssignmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Assignment not found"))
			return
		}
		logger.Error("Failed to delete assignment",
			zap.String("tenant_id", tenantID.String()),
			zap.String("assignment_id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to delete assignment"))
		return
	}

	response.NoContent(c)
}

// GetStaffAssignments retrieves all assignments for a staff member.
// @Summary Get staff assignments
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Param academic_year_id query string false "Filter by academic year ID"
// @Success 200 {object} response.Success{data=[]AssignmentResponse}
// @Router /api/v1/staff/{id}/assignments [get]
func (h *Handler) GetStaffAssignments(c *gin.Context) {
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

	var academicYearID *uuid.UUID
	if academicYearIDStr := c.Query("academic_year_id"); academicYearIDStr != "" {
		id, err := uuid.Parse(academicYearIDStr)
		if err == nil {
			academicYearID = &id
		}
	}

	assignments, err := h.service.GetStaffAssignments(c.Request.Context(), tenantID, staffID, academicYearID)
	if err != nil {
		logger.Error("Failed to get staff assignments",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", staffID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve staff assignments"))
		return
	}

	response.OK(c, ToAssignmentResponses(assignments))
}

// GetWorkloadReport retrieves the workload report for all teachers.
// @Summary Get workload report
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param academic_year_id query string true "Academic year ID"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} response.Success{data=WorkloadReportResponse}
// @Router /api/v1/teacher-assignments/workload [get]
func (h *Handler) GetWorkloadReport(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		id, err := uuid.Parse(branchIDStr)
		if err == nil {
			branchID = &id
		}
	}

	report, err := h.service.GetWorkloadReport(c.Request.Context(), tenantID, academicYearID, branchID)
	if err != nil {
		logger.Error("Failed to get workload report",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to generate workload report"))
		return
	}

	response.OK(c, report)
}

// GetUnassignedSubjects retrieves subjects without teachers.
// @Summary Get unassigned subjects
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param academic_year_id query string true "Academic year ID"
// @Success 200 {object} response.Success{data=UnassignedSubjectsResponse}
// @Router /api/v1/teacher-assignments/unassigned [get]
func (h *Handler) GetUnassignedSubjects(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	subjects, err := h.service.GetUnassignedSubjects(c.Request.Context(), tenantID, academicYearID)
	if err != nil {
		logger.Error("Failed to get unassigned subjects",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve unassigned subjects"))
		return
	}

	response.OK(c, subjects)
}

// BulkCreate creates multiple assignments at once.
// @Summary Bulk create teacher assignments
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body BulkCreateRequest true "Assignments"
// @Success 201 {object} response.Success{data=[]AssignmentResponse}
// @Router /api/v1/teacher-assignments/bulk [post]
func (h *Handler) BulkCreate(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req BulkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	items := make([]BulkAssignmentItem, len(req.Assignments))
	for i, item := range req.Assignments {
		staffID, _ := uuid.Parse(item.StaffID)
		subjectID, _ := uuid.Parse(item.SubjectID)
		classID, _ := uuid.Parse(item.ClassID)
		academicYearID, _ := uuid.Parse(item.AcademicYearID)
		effectiveFrom, _ := time.Parse("2006-01-02", item.EffectiveFrom)

		items[i] = BulkAssignmentItem{
			StaffID:        staffID,
			SubjectID:      subjectID,
			ClassID:        classID,
			AcademicYearID: academicYearID,
			PeriodsPerWeek: item.PeriodsPerWeek,
			IsClassTeacher: item.IsClassTeacher,
			EffectiveFrom:  effectiveFrom,
		}

		if item.SectionID != nil {
			sectionID, err := uuid.Parse(*item.SectionID)
			if err == nil {
				items[i].SectionID = &sectionID
			}
		}
	}

	dto := BulkCreateAssignmentDTO{
		TenantID:    tenantID,
		Assignments: items,
		CreatedBy:   &userID,
	}

	created, errs := h.service.BulkCreate(c.Request.Context(), dto)
	if len(errs) > 0 {
		// Log errors but return partial success
		for _, err := range errs {
			logger.Error("Failed to create assignment in bulk",
				zap.String("tenant_id", tenantID.String()),
				zap.Error(err))
		}
	}

	response.Created(c, map[string]interface{}{
		"created": ToAssignmentResponses(created),
		"errors":  len(errs),
		"total":   len(req.Assignments),
	})
}

// GetClassTeacher retrieves the class teacher for a class-section.
// @Summary Get class teacher
// @Tags Assignments
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Class ID"
// @Param section_id query string false "Section ID"
// @Param academic_year_id query string true "Academic year ID"
// @Success 200 {object} response.Success{data=ClassTeacherResponse}
// @Router /api/v1/classes/{id}/class-teacher [get]
func (h *Handler) GetClassTeacher(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var sectionID *uuid.UUID
	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		id, err := uuid.Parse(sectionIDStr)
		if err == nil {
			sectionID = &id
		}
	}

	assignment, err := h.service.GetClassTeacher(c.Request.Context(), tenantID, classID, sectionID, academicYearID)
	if err != nil {
		logger.Error("Failed to get class teacher",
			zap.String("tenant_id", tenantID.String()),
			zap.String("class_id", classID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve class teacher"))
		return
	}

	resp := ClassTeacherResponse{
		ClassID:    classID.String(),
		IsAssigned: assignment != nil,
	}

	if sectionID != nil {
		resp.SectionID = sectionID.String()
	}

	if assignment != nil {
		resp.TeacherID = assignment.StaffID.String()
		resp.TeacherName = assignment.Staff.FullName()
		resp.ClassName = assignment.Class.Name
		if assignment.Section != nil {
			resp.SectionName = assignment.Section.Name
		}
	}

	response.OK(c, resp)
}

// SetClassTeacher sets a teacher as the class teacher.
// @Summary Set class teacher
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Class ID"
// @Param section_id query string false "Section ID"
// @Param academic_year_id query string true "Academic year ID"
// @Param request body SetClassTeacherRequest true "Teacher to set"
// @Success 200 {object} response.Success{data=ClassTeacherResponse}
// @Router /api/v1/classes/{id}/class-teacher [put]
func (h *Handler) SetClassTeacher(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	academicYearIDStr := c.Query("academic_year_id")
	if academicYearIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Academic year ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var sectionID *uuid.UUID
	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		id, err := uuid.Parse(sectionIDStr)
		if err == nil {
			sectionID = &id
		}
	}

	var req SetClassTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	staffID, _ := uuid.Parse(req.StaffID)

	err = h.service.SetClassTeacher(c.Request.Context(), tenantID, classID, sectionID, academicYearID, staffID)
	if err != nil {
		logger.Error("Failed to set class teacher",
			zap.String("tenant_id", tenantID.String()),
			zap.String("class_id", classID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	resp := ClassTeacherResponse{
		ClassID:    classID.String(),
		TeacherID:  staffID.String(),
		IsAssigned: true,
	}
	if sectionID != nil {
		resp.SectionID = sectionID.String()
	}

	response.OK(c, resp)
}

// RegisterRoutes registers assignment routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	assignments := rg.Group("/teacher-assignments")
	{
		assignments.GET("", middleware.PermissionRequired("assignment:view"), h.List)
		assignments.POST("", middleware.PermissionRequired("assignment:create"), h.Create)
		assignments.GET("/workload", middleware.PermissionRequired("assignment:workload"), h.GetWorkloadReport)
		assignments.GET("/unassigned", middleware.PermissionRequired("assignment:view"), h.GetUnassignedSubjects)
		assignments.POST("/bulk", middleware.PermissionRequired("assignment:create"), h.BulkCreate)
		assignments.GET("/:id", middleware.PermissionRequired("assignment:view"), h.Get)
		assignments.PUT("/:id", middleware.PermissionRequired("assignment:update"), h.Update)
		assignments.DELETE("/:id", middleware.PermissionRequired("assignment:delete"), h.Delete)
	}
}

// RegisterStaffRoutes registers staff-related assignment routes.
func (h *Handler) RegisterStaffRoutes(staffGroup *gin.RouterGroup) {
	staffGroup.GET("/:id/assignments", middleware.PermissionRequired("assignment:view"), h.GetStaffAssignments)
}

// RegisterClassRoutes registers class-related assignment routes.
func (h *Handler) RegisterClassRoutes(rg *gin.RouterGroup) {
	classes := rg.Group("/classes")
	{
		classes.GET("/:id/class-teacher", middleware.PermissionRequired("assignment:view"), h.GetClassTeacher)
		classes.PUT("/:id/class-teacher", middleware.PermissionRequired("assignment:class_teacher"), h.SetClassTeacher)
	}
}
