// Package academic provides academic structure management functionality.
package academic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// Handler handles academic-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new academic handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// Class Handlers
// ========================================

// ListClasses returns all classes for the tenant.
func (h *Handler) ListClasses(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ClassFilter{TenantID: tenantID}

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

	if hasStreamsStr := c.Query("has_streams"); hasStreamsStr != "" {
		hasStreams := hasStreamsStr == "true"
		filter.HasStreams = &hasStreams
	}

	if level := c.Query("level"); level != "" {
		filter.Level = &level
	}

	filter.Search = c.Query("search")

	classes, total, err := h.service.ListClasses(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list classes"))
		return
	}

	resp := ClassListResponse{
		Classes: make([]ClassResponse, len(classes)),
		Total:   total,
	}
	for i, class := range classes {
		resp.Classes[i] = ClassToResponse(&class)
	}

	response.OK(c, resp)
}

// GetClass returns a single class by ID.
func (h *Handler) GetClass(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	class, err := h.service.GetClassByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Class not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get class"))
		return
	}

	response.OK(c, ClassToResponse(class))
}

// CreateClass creates a new class.
func (h *Handler) CreateClass(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	class, err := h.service.CreateClass(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrClassCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Class code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create class"))
		return
	}

	response.Created(c, ClassToResponse(class))
}

// UpdateClass updates an existing class.
func (h *Handler) UpdateClass(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	var req UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	class, err := h.service.UpdateClass(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Class not found"))
			return
		}
		if errors.Is(err, ErrClassCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Class code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update class"))
		return
	}

	response.OK(c, ClassToResponse(class))
}

// DeleteClass deletes a class.
func (h *Handler) DeleteClass(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	if err := h.service.DeleteClass(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrClassNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Class not found"))
			return
		}
		if errors.Is(err, ErrClassHasSections) {
			apperrors.Abort(c, apperrors.Conflict("Cannot delete class with sections"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete class"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetClassSections returns all sections for a class.
func (h *Handler) GetClassSections(c *gin.Context) {
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

	filter := SectionFilter{
		TenantID: tenantID,
		ClassID:  &classID,
	}

	sections, total, err := h.service.ListSections(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list sections"))
		return
	}

	resp := SectionListResponse{
		Sections: make([]SectionResponse, len(sections)),
		Total:    total,
	}
	for i, section := range sections {
		resp.Sections[i] = SectionToResponse(&section)
	}

	response.OK(c, resp)
}

// ========================================
// Section Handlers
// ========================================

// ListSections returns all sections for the tenant.
func (h *Handler) ListSections(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := SectionFilter{TenantID: tenantID}

	if classIDStr := c.Query("class_id"); classIDStr != "" {
		classID, err := uuid.Parse(classIDStr)
		if err == nil {
			filter.ClassID = &classID
		}
	}

	if academicYearIDStr := c.Query("academic_year_id"); academicYearIDStr != "" {
		academicYearID, err := uuid.Parse(academicYearIDStr)
		if err == nil {
			filter.AcademicYearID = &academicYearID
		}
	}

	if streamIDStr := c.Query("stream_id"); streamIDStr != "" {
		streamID, err := uuid.Parse(streamIDStr)
		if err == nil {
			filter.StreamID = &streamID
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	filter.Search = c.Query("search")

	sections, total, err := h.service.ListSections(c.Request.Context(), filter)
	if err != nil {
		fmt.Printf("\n\n[ERROR] ListSections failed: %v\n\n", err)
		apperrors.Abort(c, apperrors.InternalError("Failed to list sections"))
		return
	}

	resp := SectionListResponse{
		Sections: make([]SectionResponse, len(sections)),
		Total:    total,
	}
	for i, section := range sections {
		resp.Sections[i] = SectionToResponse(&section)
	}

	response.OK(c, resp)
}

// GetSection returns a single section by ID.
func (h *Handler) GetSection(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	section, err := h.service.GetSectionByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrSectionNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get section"))
		return
	}

	response.OK(c, SectionToResponse(section))
}

// CreateSection creates a new section.
func (h *Handler) CreateSection(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	section, err := h.service.CreateSection(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrSectionCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Section code already exists for this class"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create section"))
		return
	}

	response.Created(c, SectionToResponse(section))
}

// UpdateSection updates an existing section.
func (h *Handler) UpdateSection(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	var req UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	section, err := h.service.UpdateSection(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrSectionNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
			return
		}
		if errors.Is(err, ErrSectionCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Section code already exists for this class"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update section"))
		return
	}

	response.OK(c, SectionToResponse(section))
}

// DeleteSection deletes a section.
func (h *Handler) DeleteSection(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	if err := h.service.DeleteSection(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrSectionNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
			return
		}
		if errors.Is(err, ErrSectionHasStudents) {
			apperrors.Abort(c, apperrors.Conflict("Cannot delete section with enrolled students"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete section"))
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignClassTeacher assigns a class teacher to a section.
func (h *Handler) AssignClassTeacher(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	var req struct {
		ClassTeacherID *uuid.UUID `json:"classTeacherId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	section, err := h.service.AssignClassTeacher(c.Request.Context(), tenantID, id, req.ClassTeacherID)
	if err != nil {
		if errors.Is(err, ErrSectionNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Section not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to assign class teacher"))
		return
	}

	response.OK(c, SectionToResponse(section))
}

// ========================================
// Stream Handlers
// ========================================

// ListStreams returns all streams for the tenant.
func (h *Handler) ListStreams(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := StreamFilter{TenantID: tenantID}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	filter.Search = c.Query("search")

	streams, total, err := h.service.ListStreams(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list streams"))
		return
	}

	resp := StreamListResponse{
		Streams: make([]StreamResponse, len(streams)),
		Total:   total,
	}
	for i, stream := range streams {
		resp.Streams[i] = StreamToResponse(&stream)
	}

	response.OK(c, resp)
}

// GetStream returns a single stream by ID.
func (h *Handler) GetStream(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid stream ID"))
		return
	}

	stream, err := h.service.GetStreamByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStreamNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Stream not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get stream"))
		return
	}

	response.OK(c, StreamToResponse(stream))
}

// CreateStream creates a new stream.
func (h *Handler) CreateStream(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	stream, err := h.service.CreateStream(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrStreamCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Stream code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create stream"))
		return
	}

	response.Created(c, StreamToResponse(stream))
}

// UpdateStream updates an existing stream.
func (h *Handler) UpdateStream(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid stream ID"))
		return
	}

	var req UpdateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	stream, err := h.service.UpdateStream(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrStreamNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Stream not found"))
			return
		}
		if errors.Is(err, ErrStreamCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Stream code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update stream"))
		return
	}

	response.OK(c, StreamToResponse(stream))
}

// DeleteStream deletes a stream.
func (h *Handler) DeleteStream(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid stream ID"))
		return
	}

	if err := h.service.DeleteStream(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrStreamNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Stream not found"))
			return
		}
		if errors.Is(err, ErrStreamInUse) {
			apperrors.Abort(c, apperrors.Conflict("Cannot delete stream that is in use"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete stream"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Subject Handlers
// ========================================

// ListSubjects returns all subjects for the tenant.
func (h *Handler) ListSubjects(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := SubjectFilter{TenantID: tenantID}

	if subjectType := c.Query("subject_type"); subjectType != "" {
		filter.SubjectType = &subjectType
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	filter.Search = c.Query("search")

	subjects, total, err := h.service.ListSubjects(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list subjects"))
		return
	}

	resp := SubjectListResponse{
		Subjects: make([]SubjectResponse, len(subjects)),
		Total:    total,
	}
	for i, subject := range subjects {
		resp.Subjects[i] = SubjectToResponse(&subject)
	}

	response.OK(c, resp)
}

// GetSubject returns a single subject by ID.
func (h *Handler) GetSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid subject ID"))
		return
	}

	subject, err := h.service.GetSubjectByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Subject not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get subject"))
		return
	}

	response.OK(c, SubjectToResponse(subject))
}

// CreateSubject creates a new subject.
func (h *Handler) CreateSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	subject, err := h.service.CreateSubject(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		if errors.Is(err, ErrSubjectCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Subject code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create subject"))
		return
	}

	response.Created(c, SubjectToResponse(subject))
}

// UpdateSubject updates an existing subject.
func (h *Handler) UpdateSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid subject ID"))
		return
	}

	var req UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	subject, err := h.service.UpdateSubject(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Subject not found"))
			return
		}
		if errors.Is(err, ErrSubjectCodeExists) {
			apperrors.Abort(c, apperrors.Conflict("Subject code already exists"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update subject"))
		return
	}

	response.OK(c, SubjectToResponse(subject))
}

// DeleteSubject deletes a subject.
func (h *Handler) DeleteSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid subject ID"))
		return
	}

	if err := h.service.DeleteSubject(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Subject not found"))
			return
		}
		if errors.Is(err, ErrSubjectInUse) {
			apperrors.Abort(c, apperrors.Conflict("Cannot delete subject that is in use"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete subject"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Class-Subject Handlers
// ========================================

// ListClassSubjects returns all subjects for a class.
func (h *Handler) ListClassSubjects(c *gin.Context) {
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

	filter := ClassSubjectFilter{
		TenantID: tenantID,
		ClassID:  &classID,
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	classSubjects, total, err := h.service.ListClassSubjects(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list class subjects"))
		return
	}

	resp := ClassSubjectListResponse{
		ClassSubjects: make([]ClassSubjectResponse, len(classSubjects)),
		Total:         total,
	}
	for i, cs := range classSubjects {
		resp.ClassSubjects[i] = ClassSubjectToResponse(&cs)
	}

	response.OK(c, resp)
}

// CreateClassSubject assigns a subject to a class.
func (h *Handler) CreateClassSubject(c *gin.Context) {
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

	var req CreateClassSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	classSubject, err := h.service.CreateClassSubject(c.Request.Context(), tenantID, classID, req)
	if err != nil {
		if errors.Is(err, ErrClassSubjectExists) {
			apperrors.Abort(c, apperrors.Conflict("Subject already assigned to this class"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to assign subject to class"))
		return
	}

	response.Created(c, ClassSubjectToResponse(classSubject))
}

// UpdateClassSubject updates a class-subject mapping.
func (h *Handler) UpdateClassSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("subjectId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class-subject ID"))
		return
	}

	var req UpdateClassSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	classSubject, err := h.service.UpdateClassSubject(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrClassSubjectNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Class-subject mapping not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update class-subject mapping"))
		return
	}

	response.OK(c, ClassSubjectToResponse(classSubject))
}

// DeleteClassSubject removes a subject from a class.
func (h *Handler) DeleteClassSubject(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("subjectId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class-subject ID"))
		return
	}

	if err := h.service.DeleteClassSubject(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrClassSubjectNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Class-subject mapping not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to remove subject from class"))
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================
// Structure Handlers
// ========================================

// GetClassStructure returns the hierarchical class-section structure.
func (h *Handler) GetClassStructure(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		id, err := uuid.Parse(branchIDStr)
		if err == nil {
			branchID = &id
		}
	}

	var academicYearID *uuid.UUID
	if academicYearIDStr := c.Query("academic_year_id"); academicYearIDStr != "" {
		id, err := uuid.Parse(academicYearIDStr)
		if err == nil {
			academicYearID = &id
		}
	}

	structure, err := h.service.GetClassStructure(c.Request.Context(), tenantID, branchID, academicYearID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get class structure"))
		return
	}

	response.OK(c, ClassStructureResponse{Classes: structure})
}

// ========================================
// Route Registration
// ========================================

// RegisterRoutes registers all academic routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Class routes
	classes := rg.Group("/classes")
	{
		// View operations
		classesView := classes.Group("")
		classesView.Use(middleware.PermissionRequired("class:view"))
		{
			classesView.GET("", h.ListClasses)
			classesView.GET("/structure", h.GetClassStructure)
			classesView.GET("/:id", h.GetClass)
			classesView.GET("/:id/sections", h.GetClassSections)
		}

		// Create operations
		classesCreate := classes.Group("")
		classesCreate.Use(middleware.PermissionRequired("class:create"))
		{
			classesCreate.POST("", h.CreateClass)
		}

		// Update operations
		classesUpdate := classes.Group("")
		classesUpdate.Use(middleware.PermissionRequired("class:update"))
		{
			classesUpdate.PUT("/:id", h.UpdateClass)
		}

		// Delete operations
		classesDelete := classes.Group("")
		classesDelete.Use(middleware.PermissionRequired("class:delete"))
		{
			classesDelete.DELETE("/:id", h.DeleteClass)
		}
	}

	// Section routes
	sections := rg.Group("/sections")
	{
		// View operations
		sectionsView := sections.Group("")
		sectionsView.Use(middleware.PermissionRequired("section:view"))
		{
			sectionsView.GET("", h.ListSections)
			sectionsView.GET("/:id", h.GetSection)
		}

		// Create operations
		sectionsCreate := sections.Group("")
		sectionsCreate.Use(middleware.PermissionRequired("section:create"))
		{
			sectionsCreate.POST("", h.CreateSection)
		}

		// Update operations
		sectionsUpdate := sections.Group("")
		sectionsUpdate.Use(middleware.PermissionRequired("section:update"))
		{
			sectionsUpdate.PUT("/:id", h.UpdateSection)
			sectionsUpdate.PUT("/:id/class-teacher", h.AssignClassTeacher)
		}

		// Delete operations
		sectionsDelete := sections.Group("")
		sectionsDelete.Use(middleware.PermissionRequired("section:delete"))
		{
			sectionsDelete.DELETE("/:id", h.DeleteSection)
		}
	}

	// Stream routes
	streams := rg.Group("/streams")
	{
		// View operations
		streamsView := streams.Group("")
		streamsView.Use(middleware.PermissionRequired("stream:view"))
		{
			streamsView.GET("", h.ListStreams)
			streamsView.GET("/:id", h.GetStream)
		}

		// Create operations
		streamsCreate := streams.Group("")
		streamsCreate.Use(middleware.PermissionRequired("stream:create"))
		{
			streamsCreate.POST("", h.CreateStream)
		}

		// Update operations
		streamsUpdate := streams.Group("")
		streamsUpdate.Use(middleware.PermissionRequired("stream:update"))
		{
			streamsUpdate.PUT("/:id", h.UpdateStream)
		}

		// Delete operations
		streamsDelete := streams.Group("")
		streamsDelete.Use(middleware.PermissionRequired("stream:delete"))
		{
			streamsDelete.DELETE("/:id", h.DeleteStream)
		}
	}

	// Subject routes
	subjects := rg.Group("/subjects")
	{
		// View operations
		subjectsView := subjects.Group("")
		subjectsView.Use(middleware.PermissionRequired("subject:view"))
		{
			subjectsView.GET("", h.ListSubjects)
			subjectsView.GET("/:id", h.GetSubject)
		}

		// Create operations
		subjectsCreate := subjects.Group("")
		subjectsCreate.Use(middleware.PermissionRequired("subject:create"))
		{
			subjectsCreate.POST("", h.CreateSubject)
		}

		// Update operations
		subjectsUpdate := subjects.Group("")
		subjectsUpdate.Use(middleware.PermissionRequired("subject:update"))
		{
			subjectsUpdate.PUT("/:id", h.UpdateSubject)
		}

		// Delete operations
		subjectsDelete := subjects.Group("")
		subjectsDelete.Use(middleware.PermissionRequired("subject:delete"))
		{
			subjectsDelete.DELETE("/:id", h.DeleteSubject)
		}
	}

	// Class-Subject routes (nested under classes)
	classSubjects := classes.Group("/:id/subjects")
	{
		// View operations
		classSubjectsView := classSubjects.Group("")
		classSubjectsView.Use(middleware.PermissionRequired("class:view"))
		{
			classSubjectsView.GET("", h.ListClassSubjects)
		}

		// Create operations
		classSubjectsCreate := classSubjects.Group("")
		classSubjectsCreate.Use(middleware.PermissionRequired("class:update"))
		{
			classSubjectsCreate.POST("", h.CreateClassSubject)
		}

		// Update operations
		classSubjectsUpdate := classSubjects.Group("")
		classSubjectsUpdate.Use(middleware.PermissionRequired("class:update"))
		{
			classSubjectsUpdate.PUT("/:subjectId", h.UpdateClassSubject)
		}

		// Delete operations
		classSubjectsDelete := classSubjects.Group("")
		classSubjectsDelete.Use(middleware.PermissionRequired("class:update"))
		{
			classSubjectsDelete.DELETE("/:subjectId", h.DeleteClassSubject)
		}
	}
}
