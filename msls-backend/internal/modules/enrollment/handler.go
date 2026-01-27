// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// Handler handles enrollment-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new enrollment handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateEnrollmentRequest represents the request body for creating an enrollment.
type CreateEnrollmentRequest struct {
	AcademicYearID string  `json:"academicYearId" binding:"required,uuid"`
	ClassID        string  `json:"classId" binding:"omitempty,uuid"`
	SectionID      string  `json:"sectionId" binding:"omitempty,uuid"`
	RollNumber     string  `json:"rollNumber" binding:"max=20"`
	ClassTeacherID string  `json:"classTeacherId" binding:"omitempty,uuid"`
	EnrollmentDate string  `json:"enrollmentDate"` // Format: YYYY-MM-DD
	Notes          string  `json:"notes"`
}

// UpdateEnrollmentRequest represents the request body for updating an enrollment.
type UpdateEnrollmentRequest struct {
	ClassID        *string `json:"classId" binding:"omitempty,uuid"`
	SectionID      *string `json:"sectionId" binding:"omitempty,uuid"`
	RollNumber     *string `json:"rollNumber" binding:"omitempty,max=20"`
	ClassTeacherID *string `json:"classTeacherId" binding:"omitempty,uuid"`
	Notes          *string `json:"notes"`
}

// TransferRequest represents the request body for processing a transfer.
type TransferRequest struct {
	TransferDate   string `json:"transferDate" binding:"required"` // Format: YYYY-MM-DD
	TransferReason string `json:"transferReason" binding:"required"`
}

// DropoutRequest represents the request body for processing a dropout.
type DropoutRequest struct {
	DropoutDate   string `json:"dropoutDate" binding:"required"` // Format: YYYY-MM-DD
	DropoutReason string `json:"dropoutReason" binding:"required"`
}

// ListEnrollments returns all enrollments for a student (enrollment history).
// @Summary List enrollment history
// @Description Get all enrollment records for a student
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=EnrollmentHistoryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments [get]
func (h *Handler) ListEnrollments(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	enrollments, err := h.service.ListByStudent(c.Request.Context(), tenantID, studentID)
	if err != nil {
		logger.Error("Failed to list enrollments",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enrollments"))
		return
	}

	resp := EnrollmentHistoryResponse{
		Enrollments: ToEnrollmentResponses(enrollments),
		Total:       len(enrollments),
	}

	response.OK(c, resp)
}

// GetCurrentEnrollment returns the current active enrollment for a student.
// @Summary Get current enrollment
// @Description Get the current active enrollment for a student
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments/current [get]
func (h *Handler) GetCurrentEnrollment(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	enrollment, err := h.service.GetCurrentByStudent(c.Request.Context(), tenantID, studentID)
	if err != nil {
		if errors.Is(err, ErrActiveEnrollmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("No active enrollment found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enrollment"))
		return
	}

	response.OK(c, ToEnrollmentResponse(enrollment))
}

// CreateEnrollment creates a new enrollment for a student.
// @Summary Create enrollment
// @Description Create a new enrollment record for a student
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param request body CreateEnrollmentRequest true "Enrollment details"
// @Success 201 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments [post]
func (h *Handler) CreateEnrollment(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	dto := CreateEnrollmentDTO{
		TenantID:       tenantID,
		StudentID:      studentID,
		AcademicYearID: academicYearID,
		RollNumber:     req.RollNumber,
		Notes:          req.Notes,
		CreatedBy:      &userID,
	}

	// Parse optional UUIDs
	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
			return
		}
		dto.ClassID = &classID
	}

	if req.SectionID != "" {
		sectionID, err := uuid.Parse(req.SectionID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
			return
		}
		dto.SectionID = &sectionID
	}

	if req.ClassTeacherID != "" {
		classTeacherID, err := uuid.Parse(req.ClassTeacherID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class teacher ID"))
			return
		}
		dto.ClassTeacherID = &classTeacherID
	}

	// Parse enrollment date if provided
	if req.EnrollmentDate != "" {
		enrollmentDate, err := time.Parse("2006-01-02", req.EnrollmentDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid enrollment date format, use YYYY-MM-DD"))
			return
		}
		dto.EnrollmentDate = &enrollmentDate
	}

	enrollment, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateEnrollment):
			apperrors.Abort(c, apperrors.Conflict("Student already enrolled in this academic year"))
		case errors.Is(err, ErrActiveEnrollmentExists):
			apperrors.Abort(c, apperrors.Conflict("Student already has an active enrollment"))
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		case errors.Is(err, ErrAcademicYearNotFound):
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		default:
			logger.Error("Failed to create enrollment", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to create enrollment"))
		}
		return
	}

	response.Created(c, ToEnrollmentResponse(enrollment))
}

// UpdateEnrollment updates an enrollment.
// @Summary Update enrollment
// @Description Update an existing enrollment record
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param eid path string true "Enrollment ID" format(uuid)
// @Param request body UpdateEnrollmentRequest true "Enrollment updates"
// @Success 200 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments/{eid} [put]
func (h *Handler) UpdateEnrollment(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	enrollmentIDParam := c.Param("eid")
	enrollmentID, err := uuid.Parse(enrollmentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enrollment ID"))
		return
	}

	var req UpdateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateEnrollmentDTO{
		RollNumber: req.RollNumber,
		Notes:      req.Notes,
		UpdatedBy:  &userID,
	}

	// Parse optional UUIDs
	if req.ClassID != nil {
		classID, err := uuid.Parse(*req.ClassID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
			return
		}
		dto.ClassID = &classID
	}

	if req.SectionID != nil {
		sectionID, err := uuid.Parse(*req.SectionID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
			return
		}
		dto.SectionID = &sectionID
	}

	if req.ClassTeacherID != nil {
		classTeacherID, err := uuid.Parse(*req.ClassTeacherID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class teacher ID"))
			return
		}
		dto.ClassTeacherID = &classTeacherID
	}

	enrollment, err := h.service.Update(c.Request.Context(), tenantID, enrollmentID, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrEnrollmentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Enrollment not found"))
		case errors.Is(err, ErrCannotModifyInactiveEnrollment):
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify inactive enrollment"))
		default:
			logger.Error("Failed to update enrollment", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update enrollment"))
		}
		return
	}

	response.OK(c, ToEnrollmentResponse(enrollment))
}

// ProcessTransfer processes a student transfer.
// @Summary Process transfer
// @Description Process a student transfer (marks enrollment as transferred and student as inactive)
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param eid path string true "Enrollment ID" format(uuid)
// @Param request body TransferRequest true "Transfer details"
// @Success 200 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments/{eid}/transfer [post]
func (h *Handler) ProcessTransfer(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	transferDate, err := time.Parse("2006-01-02", req.TransferDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid transfer date format, use YYYY-MM-DD"))
		return
	}

	dto := TransferDTO{
		TransferDate:   transferDate,
		TransferReason: req.TransferReason,
		UpdatedBy:      &userID,
	}

	enrollment, err := h.service.ProcessTransfer(c.Request.Context(), tenantID, studentID, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrActiveEnrollmentNotFound):
			apperrors.Abort(c, apperrors.NotFound("No active enrollment found"))
		case errors.Is(err, ErrTransferDateRequired):
			apperrors.Abort(c, apperrors.BadRequest("Transfer date is required"))
		case errors.Is(err, ErrTransferReasonRequired):
			apperrors.Abort(c, apperrors.BadRequest("Transfer reason is required"))
		default:
			logger.Error("Failed to process transfer", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to process transfer"))
		}
		return
	}

	response.OK(c, ToEnrollmentResponse(enrollment))
}

// ProcessDropout processes a student dropout.
// @Summary Process dropout
// @Description Process a student dropout (marks enrollment as dropout and student as inactive)
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param eid path string true "Enrollment ID" format(uuid)
// @Param request body DropoutRequest true "Dropout details"
// @Success 200 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments/{eid}/dropout [post]
func (h *Handler) ProcessDropout(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req DropoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dropoutDate, err := time.Parse("2006-01-02", req.DropoutDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid dropout date format, use YYYY-MM-DD"))
		return
	}

	dto := DropoutDTO{
		DropoutDate:   dropoutDate,
		DropoutReason: req.DropoutReason,
		UpdatedBy:     &userID,
	}

	enrollment, err := h.service.ProcessDropout(c.Request.Context(), tenantID, studentID, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrActiveEnrollmentNotFound):
			apperrors.Abort(c, apperrors.NotFound("No active enrollment found"))
		case errors.Is(err, ErrDropoutDateRequired):
			apperrors.Abort(c, apperrors.BadRequest("Dropout date is required"))
		case errors.Is(err, ErrDropoutReasonRequired):
			apperrors.Abort(c, apperrors.BadRequest("Dropout reason is required"))
		default:
			logger.Error("Failed to process dropout", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to process dropout"))
		}
		return
	}

	response.OK(c, ToEnrollmentResponse(enrollment))
}

// ListByClass returns all enrollments for a class.
// @Summary List enrollments by class
// @Description Get all active enrollments for a class
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param classId path string true "Class ID" format(uuid)
// @Success 200 {object} response.Success{data=EnrollmentHistoryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/enrollments/by-class/{classId} [get]
func (h *Handler) ListByClass(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	classIDParam := c.Param("classId")
	classID, err := uuid.Parse(classIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
		return
	}

	enrollments, err := h.service.ListByClass(c.Request.Context(), tenantID, classID)
	if err != nil {
		logger.Error("Failed to list enrollments by class", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enrollments"))
		return
	}

	resp := EnrollmentHistoryResponse{
		Enrollments: ToEnrollmentResponses(enrollments),
		Total:       len(enrollments),
	}

	response.OK(c, resp)
}

// ListBySection returns all enrollments for a section.
// @Summary List enrollments by section
// @Description Get all active enrollments for a section
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sectionId path string true "Section ID" format(uuid)
// @Success 200 {object} response.Success{data=EnrollmentHistoryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/enrollments/by-section/{sectionId} [get]
func (h *Handler) ListBySection(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sectionIDParam := c.Param("sectionId")
	sectionID, err := uuid.Parse(sectionIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
		return
	}

	enrollments, err := h.service.ListBySection(c.Request.Context(), tenantID, sectionID)
	if err != nil {
		logger.Error("Failed to list enrollments by section", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enrollments"))
		return
	}

	resp := EnrollmentHistoryResponse{
		Enrollments: ToEnrollmentResponses(enrollments),
		Total:       len(enrollments),
	}

	response.OK(c, resp)
}

// GetStatusHistory returns the status change history for an enrollment.
// @Summary Get enrollment status history
// @Description Get the status change history for an enrollment
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param eid path string true "Enrollment ID" format(uuid)
// @Success 200 {object} response.Success{data=[]StatusChangeResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/enrollments/{eid}/status-history [get]
func (h *Handler) GetStatusHistory(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	enrollmentIDParam := c.Param("eid")
	enrollmentID, err := uuid.Parse(enrollmentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enrollment ID"))
		return
	}

	changes, err := h.service.GetStatusHistory(c.Request.Context(), tenantID, enrollmentID)
	if err != nil {
		if errors.Is(err, ErrEnrollmentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Enrollment not found"))
			return
		}
		logger.Error("Failed to get status history", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve status history"))
		return
	}

	response.OK(c, ToStatusChangeResponses(changes))
}
