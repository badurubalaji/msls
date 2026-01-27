// Package academicyear provides HTTP handlers for academic year management endpoints.
package academicyear

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	academicyearservice "msls-backend/internal/services/academicyear"
)

// Handler handles academic year-related HTTP requests.
type Handler struct {
	academicYearService *academicyearservice.Service
}

// NewHandler creates a new academic year Handler.
func NewHandler(academicYearService *academicyearservice.Service) *Handler {
	return &Handler{academicYearService: academicYearService}
}

// ============================================================================
// Academic Year Handlers
// ============================================================================

// List returns all academic years for the tenant.
// @Summary List academic years
// @Description Get all academic years for the current tenant
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name"
// @Param is_current query bool false "Filter by current status"
// @Param is_active query bool false "Filter by active status"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} response.Success{data=AcademicYearListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/academic-years [get]
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := academicyearservice.ListAcademicYearFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	// Parse is_current filter
	if isCurrentStr := c.Query("is_current"); isCurrentStr != "" {
		isCurrent := isCurrentStr == "true"
		filter.IsCurrent = &isCurrent
	}

	// Parse is_active filter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	// Parse branch_id filter
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = &branchID
	}

	academicYears, err := h.academicYearService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve academic years"))
		return
	}

	total, err := h.academicYearService.Count(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to count academic years"))
		return
	}

	resp := AcademicYearListResponse{
		AcademicYears: academicYearsToResponses(academicYears),
		Total:         total,
	}

	response.OK(c, resp)
}

// GetByID returns an academic year by ID.
// @Summary Get academic year by ID
// @Description Get an academic year with full details including terms and holidays
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Success 200 {object} response.Success{data=AcademicYearResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	academicYear, err := h.academicYearService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve academic year"))
		}
		return
	}

	response.OK(c, academicYearToResponse(academicYear))
}

// GetCurrent returns the current academic year for the tenant.
// @Summary Get current academic year
// @Description Get the current active academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string false "Branch ID for branch-specific current year"
// @Success 200 {object} response.Success{data=AcademicYearResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/current [get]
func (h *Handler) GetCurrent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		parsedID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &parsedID
	}

	academicYear, err := h.academicYearService.GetCurrent(c.Request.Context(), tenantID, branchID)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("No current academic year set"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve current academic year"))
		}
		return
	}

	response.OK(c, academicYearToResponse(academicYear))
}

// Create creates a new academic year.
// @Summary Create academic year
// @Description Create a new academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateAcademicYearRequest true "Academic year details"
// @Success 201 {object} response.Success{data=AcademicYearResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/academic-years [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
		return
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if req.BranchID != nil {
		parsedID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &parsedID
	}

	createReq := academicyearservice.CreateAcademicYearRequest{
		TenantID:  tenantID,
		BranchID:  branchID,
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
		IsCurrent: req.IsCurrent,
		CreatedBy: &userID,
	}

	academicYear, err := h.academicYearService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Academic year name is required"))
		case academicyearservice.ErrAcademicYearNameExists:
			apperrors.Abort(c, apperrors.Conflict("Academic year with this name already exists"))
		case academicyearservice.ErrAcademicYearInvalidDates:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after start date"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create academic year"))
		}
		return
	}

	response.Created(c, academicYearToResponse(academicYear))
}

// Update updates an academic year.
// @Summary Update academic year
// @Description Update an existing academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param request body UpdateAcademicYearRequest true "Academic year updates"
// @Success 200 {object} response.Success{data=AcademicYearResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var req UpdateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := academicyearservice.UpdateAcademicYearRequest{
		Name:      req.Name,
		IsActive:  req.IsActive,
		UpdatedBy: &userID,
	}

	// Parse start date if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.StartDate = &startDate
	}

	// Parse end date if provided
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.EndDate = &endDate
	}

	academicYear, err := h.academicYearService.Update(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		case academicyearservice.ErrAcademicYearNameExists:
			apperrors.Abort(c, apperrors.Conflict("Academic year with this name already exists"))
		case academicyearservice.ErrAcademicYearInvalidDates:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after start date"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update academic year"))
		}
		return
	}

	response.OK(c, academicYearToResponse(academicYear))
}

// SetCurrent sets an academic year as the current year.
// @Summary Set academic year as current
// @Description Set an academic year as the current year for the tenant
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Success 200 {object} response.Success{data=AcademicYearResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/current [patch]
func (h *Handler) SetCurrent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	academicYear, err := h.academicYearService.SetCurrent(c.Request.Context(), tenantID, id, &userID)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		case academicyearservice.ErrCannotSetInactiveAsCurrent:
			apperrors.Abort(c, apperrors.BadRequest("Cannot set inactive academic year as current"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to set current academic year"))
		}
		return
	}

	response.OK(c, academicYearToResponse(academicYear))
}

// Delete deletes an academic year.
// @Summary Delete academic year
// @Description Delete an academic year and all its terms and holidays
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	err = h.academicYearService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		case academicyearservice.ErrAcademicYearHasDependencies:
			apperrors.Abort(c, apperrors.Conflict("Academic year has associated records and cannot be deleted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete academic year"))
		}
		return
	}

	response.NoContent(c)
}

// ============================================================================
// Term Handlers
// ============================================================================

// ListTerms returns all terms for an academic year.
// @Summary List terms
// @Description Get all terms for an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Success 200 {object} response.Success{data=TermListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/terms [get]
func (h *Handler) ListTerms(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	academicYearID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	// Verify academic year exists
	_, err = h.academicYearService.GetByID(c.Request.Context(), tenantID, academicYearID)
	if err != nil {
		if err == academicyearservice.ErrAcademicYearNotFound {
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to verify academic year"))
		return
	}

	terms, err := h.academicYearService.ListTerms(c.Request.Context(), tenantID, academicYearID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve terms"))
		return
	}

	resp := TermListResponse{
		Terms: termsToResponses(terms),
		Total: len(terms),
	}

	response.OK(c, resp)
}

// CreateTerm creates a new term for an academic year.
// @Summary Create term
// @Description Create a new term for an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param request body CreateTermRequest true "Term details"
// @Success 201 {object} response.Success{data=TermResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/terms [post]
func (h *Handler) CreateTerm(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	academicYearID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var req CreateTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
		return
	}

	createReq := academicyearservice.CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYearID,
		Name:           req.Name,
		StartDate:      startDate,
		EndDate:        endDate,
		Sequence:       req.Sequence,
	}

	term, err := h.academicYearService.CreateTerm(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		case academicyearservice.ErrTermNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Term name is required"))
		case academicyearservice.ErrTermNameExists:
			apperrors.Abort(c, apperrors.Conflict("Term with this name already exists"))
		case academicyearservice.ErrTermInvalidDates:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after start date"))
		case academicyearservice.ErrTermOutsideAcademicYear:
			apperrors.Abort(c, apperrors.BadRequest("Term dates must be within academic year dates"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create term"))
		}
		return
	}

	response.Created(c, termToResponse(term))
}

// UpdateTerm updates a term.
// @Summary Update term
// @Description Update an existing term
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param termId path string true "Term ID" format(uuid)
// @Param request body UpdateTermRequest true "Term updates"
// @Success 200 {object} response.Success{data=TermResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/terms/{termId} [put]
func (h *Handler) UpdateTerm(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	termID, err := uuid.Parse(c.Param("termId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid term ID"))
		return
	}

	var req UpdateTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := academicyearservice.UpdateTermRequest{
		Name:     req.Name,
		Sequence: req.Sequence,
	}

	// Parse start date if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.StartDate = &startDate
	}

	// Parse end date if provided
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.EndDate = &endDate
	}

	term, err := h.academicYearService.UpdateTerm(c.Request.Context(), tenantID, academicYearID, termID, updateReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrTermNotFound:
			apperrors.Abort(c, apperrors.NotFound("Term not found"))
		case academicyearservice.ErrTermNameExists:
			apperrors.Abort(c, apperrors.Conflict("Term with this name already exists"))
		case academicyearservice.ErrTermInvalidDates:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after start date"))
		case academicyearservice.ErrTermOutsideAcademicYear:
			apperrors.Abort(c, apperrors.BadRequest("Term dates must be within academic year dates"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update term"))
		}
		return
	}

	response.OK(c, termToResponse(term))
}

// DeleteTerm deletes a term.
// @Summary Delete term
// @Description Delete a term from an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param termId path string true "Term ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/terms/{termId} [delete]
func (h *Handler) DeleteTerm(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	termID, err := uuid.Parse(c.Param("termId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid term ID"))
		return
	}

	err = h.academicYearService.DeleteTerm(c.Request.Context(), tenantID, academicYearID, termID)
	if err != nil {
		switch err {
		case academicyearservice.ErrTermNotFound:
			apperrors.Abort(c, apperrors.NotFound("Term not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete term"))
		}
		return
	}

	response.NoContent(c)
}

// ============================================================================
// Holiday Handlers
// ============================================================================

// ListHolidays returns all holidays for an academic year.
// @Summary List holidays
// @Description Get all holidays for an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Success 200 {object} response.Success{data=HolidayListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/holidays [get]
func (h *Handler) ListHolidays(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	academicYearID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	// Verify academic year exists
	_, err = h.academicYearService.GetByID(c.Request.Context(), tenantID, academicYearID)
	if err != nil {
		if err == academicyearservice.ErrAcademicYearNotFound {
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to verify academic year"))
		return
	}

	holidays, err := h.academicYearService.ListHolidays(c.Request.Context(), tenantID, academicYearID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve holidays"))
		return
	}

	resp := HolidayListResponse{
		Holidays: holidaysToResponses(holidays),
		Total:    len(holidays),
	}

	response.OK(c, resp)
}

// CreateHoliday creates a new holiday for an academic year.
// @Summary Create holiday
// @Description Create a new holiday for an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param request body CreateHolidayRequest true "Holiday details"
// @Success 201 {object} response.Success{data=HolidayResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/holidays [post]
func (h *Handler) CreateHoliday(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	academicYearID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	var req CreateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid date format. Use YYYY-MM-DD"))
		return
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if req.BranchID != nil {
		parsedID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &parsedID
	}

	// Parse holiday type
	holidayType := models.HolidayTypePublic
	if req.Type != "" {
		holidayType = models.HolidayType(req.Type)
	}

	createReq := academicyearservice.CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYearID,
		BranchID:       branchID,
		Name:           req.Name,
		Date:           date,
		Type:           holidayType,
		IsOptional:     req.IsOptional,
	}

	holiday, err := h.academicYearService.CreateHoliday(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrAcademicYearNotFound:
			apperrors.Abort(c, apperrors.NotFound("Academic year not found"))
		case academicyearservice.ErrHolidayNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Holiday name is required"))
		case academicyearservice.ErrHolidayOutsideAcademicYear:
			apperrors.Abort(c, apperrors.BadRequest("Holiday date must be within academic year dates"))
		case academicyearservice.ErrHolidayInvalidType:
			apperrors.Abort(c, apperrors.BadRequest("Invalid holiday type"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create holiday"))
		}
		return
	}

	response.Created(c, holidayToResponse(holiday))
}

// UpdateHoliday updates a holiday.
// @Summary Update holiday
// @Description Update an existing holiday
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param holidayId path string true "Holiday ID" format(uuid)
// @Param request body UpdateHolidayRequest true "Holiday updates"
// @Success 200 {object} response.Success{data=HolidayResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/holidays/{holidayId} [put]
func (h *Handler) UpdateHoliday(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	holidayID, err := uuid.Parse(c.Param("holidayId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid holiday ID"))
		return
	}

	var req UpdateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := academicyearservice.UpdateHolidayRequest{
		Name:       req.Name,
		IsOptional: req.IsOptional,
	}

	// Parse date if provided
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.Date = &date
	}

	// Parse type if provided
	if req.Type != nil {
		holidayType := models.HolidayType(*req.Type)
		updateReq.Type = &holidayType
	}

	holiday, err := h.academicYearService.UpdateHoliday(c.Request.Context(), tenantID, academicYearID, holidayID, updateReq)
	if err != nil {
		switch err {
		case academicyearservice.ErrHolidayNotFound:
			apperrors.Abort(c, apperrors.NotFound("Holiday not found"))
		case academicyearservice.ErrHolidayOutsideAcademicYear:
			apperrors.Abort(c, apperrors.BadRequest("Holiday date must be within academic year dates"))
		case academicyearservice.ErrHolidayInvalidType:
			apperrors.Abort(c, apperrors.BadRequest("Invalid holiday type"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update holiday"))
		}
		return
	}

	response.OK(c, holidayToResponse(holiday))
}

// DeleteHoliday deletes a holiday.
// @Summary Delete holiday
// @Description Delete a holiday from an academic year
// @Tags Academic Years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Academic Year ID" format(uuid)
// @Param holidayId path string true "Holiday ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/academic-years/{id}/holidays/{holidayId} [delete]
func (h *Handler) DeleteHoliday(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	academicYearID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	holidayID, err := uuid.Parse(c.Param("holidayId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid holiday ID"))
		return
	}

	err = h.academicYearService.DeleteHoliday(c.Request.Context(), tenantID, academicYearID, holidayID)
	if err != nil {
		switch err {
		case academicyearservice.ErrHolidayNotFound:
			apperrors.Abort(c, apperrors.NotFound("Holiday not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete holiday"))
		}
		return
	}

	response.NoContent(c)
}
