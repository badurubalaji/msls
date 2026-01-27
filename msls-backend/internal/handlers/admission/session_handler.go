// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	admissionservice "msls-backend/internal/services/admission"
)

// SessionHandler handles admission session HTTP requests.
type SessionHandler struct {
	sessionService *admissionservice.SessionService
}

// NewSessionHandler creates a new SessionHandler.
func NewSessionHandler(sessionService *admissionservice.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// List returns all admission sessions for the tenant.
// @Summary List admission sessions
// @Description Get all admission sessions for the current tenant
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name or description"
// @Param status query string false "Filter by status (upcoming, open, closed)"
// @Param branchId query string false "Filter by branch ID"
// @Param academicYearId query string false "Filter by academic year ID"
// @Param includeSeats query bool false "Include seat configurations"
// @Success 200 {object} response.Success{data=SessionListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admission-sessions [get]
func (h *SessionHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := admissionservice.ListSessionFilter{
		TenantID:     tenantID,
		Search:       c.Query("search"),
		IncludeSeats: c.Query("includeSeats") == "true",
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.AdmissionSessionStatus(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	// Parse branch ID filter
	if branchIDStr := c.Query("branchId"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = &branchID
	}

	// Parse academic year ID filter
	if academicYearIDStr := c.Query("academicYearId"); academicYearIDStr != "" {
		academicYearID, err := uuid.Parse(academicYearIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
			return
		}
		filter.AcademicYearID = &academicYearID
	}

	sessions, err := h.sessionService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve admission sessions"))
		return
	}

	total, err := h.sessionService.Count(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to count admission sessions"))
		return
	}

	resp := SessionListResponse{
		Sessions: sessionsToResponses(sessions),
		Total:    total,
	}

	response.OK(c, resp)
}

// GetByID returns an admission session by ID.
// @Summary Get admission session by ID
// @Description Get an admission session with full details
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Success 200 {object} response.Success{data=SessionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id} [get]
func (h *SessionHandler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	session, err := h.sessionService.GetByID(c.Request.Context(), tenantID, id, true)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve admission session"))
		}
		return
	}

	// Get stats
	stats, _ := h.sessionService.GetStats(c.Request.Context(), tenantID, id)

	resp := sessionToResponse(session)
	if stats != nil {
		resp.Stats = &SessionStatsDTO{
			TotalApplications: stats.TotalApplications,
			ApprovedCount:     stats.ApprovedCount,
			PendingCount:      stats.PendingCount,
			RejectedCount:     stats.RejectedCount,
			TotalSeats:        stats.TotalSeats,
			FilledSeats:       stats.FilledSeats,
			AvailableSeats:    stats.AvailableSeats,
		}
	}

	response.OK(c, resp)
}

// Create creates a new admission session.
// @Summary Create admission session
// @Description Create a new admission session for the tenant
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateSessionRequest true "Session details"
// @Success 201 {object} response.Success{data=SessionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/admission-sessions [post]
func (h *SessionHandler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateSessionRequest
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

	// Parse optional UUIDs
	var branchID *uuid.UUID
	if req.BranchID != nil {
		bid, err := uuid.Parse(*req.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &bid
	}

	var academicYearID *uuid.UUID
	if req.AcademicYearID != nil {
		ayid, err := uuid.Parse(*req.AcademicYearID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
			return
		}
		academicYearID = &ayid
	}

	createReq := admissionservice.CreateSessionRequest{
		TenantID:          tenantID,
		BranchID:          branchID,
		AcademicYearID:    academicYearID,
		Name:              req.Name,
		Description:       req.Description,
		StartDate:         startDate,
		EndDate:           endDate,
		ApplicationFee:    req.ApplicationFee,
		RequiredDocuments: req.RequiredDocuments,
		Settings:          settingsToModel(req.Settings),
		CreatedBy:         &userID,
	}

	session, err := h.sessionService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Session name is required"))
		case admissionservice.ErrSessionNameExists:
			apperrors.Abort(c, apperrors.Conflict("Session with this name already exists for the academic year"))
		case admissionservice.ErrInvalidDateRange:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after or equal to start date"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create admission session"))
		}
		return
	}

	response.Created(c, sessionToResponse(session))
}

// Update updates an admission session.
// @Summary Update admission session
// @Description Update an existing admission session
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param request body UpdateSessionRequest true "Session updates"
// @Success 200 {object} response.Success{data=SessionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id} [put]
func (h *SessionHandler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := admissionservice.UpdateSessionRequest{
		Name:              req.Name,
		Description:       req.Description,
		ApplicationFee:    req.ApplicationFee,
		RequiredDocuments: req.RequiredDocuments,
		UpdatedBy:         &userID,
	}

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.StartDate = &startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
		updateReq.EndDate = &endDate
	}

	if req.Settings != nil {
		settings := settingsToModel(req.Settings)
		updateReq.Settings = &settings
	}

	session, err := h.sessionService.Update(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		case admissionservice.ErrSessionNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Session name is required"))
		case admissionservice.ErrSessionNameExists:
			apperrors.Abort(c, apperrors.Conflict("Session with this name already exists"))
		case admissionservice.ErrInvalidDateRange:
			apperrors.Abort(c, apperrors.BadRequest("End date must be after or equal to start date"))
		case admissionservice.ErrCannotModifyClosedSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify a closed admission session"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update admission session"))
		}
		return
	}

	response.OK(c, sessionToResponse(session))
}

// ChangeStatus changes the status of an admission session.
// @Summary Change admission session status
// @Description Change the status of an admission session (upcoming, open, closed)
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param request body ChangeStatusRequest true "Status change"
// @Success 200 {object} response.Success{data=SessionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/status [patch]
func (h *SessionHandler) ChangeStatus(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	newStatus := models.AdmissionSessionStatus(req.Status)
	if !newStatus.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		return
	}

	session, err := h.sessionService.ChangeStatus(c.Request.Context(), tenantID, id, newStatus, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		case admissionservice.ErrInvalidStatus:
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		case admissionservice.ErrInvalidStatusTransition:
			apperrors.Abort(c, apperrors.BadRequest("Invalid status transition"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to change session status"))
		}
		return
	}

	response.OK(c, sessionToResponse(session))
}

// ExtendDeadline extends the end date of an admission session.
// @Summary Extend admission session deadline
// @Description Extend the end date of an admission session
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param request body ExtendDeadlineRequest true "New end date"
// @Success 200 {object} response.Success{data=SessionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/extend [patch]
func (h *SessionHandler) ExtendDeadline(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	var req ExtendDeadlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	newEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD"))
		return
	}

	session, err := h.sessionService.ExtendDeadline(c.Request.Context(), tenantID, id, newEndDate, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		case admissionservice.ErrInvalidDateRange:
			apperrors.Abort(c, apperrors.BadRequest("New end date must be after start date"))
		case admissionservice.ErrCannotModifyClosedSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify a closed admission session"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to extend deadline"))
		}
		return
	}

	response.OK(c, sessionToResponse(session))
}

// Delete deletes an admission session.
// @Summary Delete admission session
// @Description Delete an admission session (cannot delete open sessions)
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id} [delete]
func (h *SessionHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	err = h.sessionService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		case admissionservice.ErrCannotDeleteOpenSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete an open admission session"))
		case admissionservice.ErrSessionHasApplications:
			apperrors.Abort(c, apperrors.Conflict("Session has applications and cannot be deleted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete admission session"))
		}
		return
	}

	response.NoContent(c)
}

// =============================================================================
// Seat Endpoints
// =============================================================================

// ListSeats returns all seat configurations for a session.
// @Summary List seat configurations
// @Description Get all seat configurations for an admission session
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Success 200 {object} response.Success{data=SeatListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/seats [get]
func (h *SessionHandler) ListSeats(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	sessionID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	seats, err := h.sessionService.ListSeats(c.Request.Context(), tenantID, sessionID)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve seat configurations"))
		}
		return
	}

	resp := SeatListResponse{
		Seats: seatsToResponses(seats),
		Total: len(seats),
	}

	response.OK(c, resp)
}

// CreateSeat creates a new seat configuration for a session.
// @Summary Create seat configuration
// @Description Create a new seat configuration for an admission session
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param request body CreateSeatRequest true "Seat configuration"
// @Success 201 {object} response.Success{data=SeatResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/seats [post]
func (h *SessionHandler) CreateSeat(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	sessionID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	var req CreateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	createReq := admissionservice.CreateSeatRequest{
		TenantID:      tenantID,
		SessionID:     sessionID,
		ClassName:     req.ClassName,
		TotalSeats:    req.TotalSeats,
		WaitlistLimit: req.WaitlistLimit,
		ReservedSeats: models.ReservedSeats(req.ReservedSeats),
	}

	seat, err := h.sessionService.CreateSeat(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		case admissionservice.ErrClassNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Class name is required"))
		case admissionservice.ErrClassAlreadyExists:
			apperrors.Abort(c, apperrors.Conflict("Seat configuration for this class already exists"))
		case admissionservice.ErrInvalidTotalSeats:
			apperrors.Abort(c, apperrors.BadRequest("Total seats must be greater than or equal to zero"))
		case admissionservice.ErrCannotModifyClosedSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot add seats to a closed session"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create seat configuration"))
		}
		return
	}

	response.Created(c, seatToResponse(seat))
}

// UpdateSeat updates a seat configuration.
// @Summary Update seat configuration
// @Description Update an existing seat configuration
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param seatId path string true "Seat ID" format(uuid)
// @Param request body UpdateSeatRequest true "Seat updates"
// @Success 200 {object} response.Success{data=SeatResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/seats/{seatId} [put]
func (h *SessionHandler) UpdateSeat(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	seatIDParam := c.Param("seatId")
	seatID, err := uuid.Parse(seatIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid seat ID"))
		return
	}

	var req UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := admissionservice.UpdateSeatRequest{
		TotalSeats:    req.TotalSeats,
		WaitlistLimit: req.WaitlistLimit,
	}

	if req.ReservedSeats != nil {
		reservedSeats := models.ReservedSeats(*req.ReservedSeats)
		updateReq.ReservedSeats = &reservedSeats
	}

	seat, err := h.sessionService.UpdateSeat(c.Request.Context(), tenantID, seatID, updateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrSeatNotFound:
			apperrors.Abort(c, apperrors.NotFound("Seat configuration not found"))
		case admissionservice.ErrInvalidTotalSeats:
			apperrors.Abort(c, apperrors.BadRequest("Total seats must be greater than or equal to zero"))
		case admissionservice.ErrFilledExceedsTotal:
			apperrors.Abort(c, apperrors.BadRequest("Filled seats cannot exceed total seats"))
		case admissionservice.ErrCannotModifyClosedSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot modify seats in a closed session"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update seat configuration"))
		}
		return
	}

	response.OK(c, seatToResponse(seat))
}

// DeleteSeat deletes a seat configuration.
// @Summary Delete seat configuration
// @Description Delete a seat configuration from an admission session
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param seatId path string true "Seat ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/seats/{seatId} [delete]
func (h *SessionHandler) DeleteSeat(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	seatIDParam := c.Param("seatId")
	seatID, err := uuid.Parse(seatIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid seat ID"))
		return
	}

	err = h.sessionService.DeleteSeat(c.Request.Context(), tenantID, seatID)
	if err != nil {
		switch err {
		case admissionservice.ErrSeatNotFound:
			apperrors.Abort(c, apperrors.NotFound("Seat configuration not found"))
		case admissionservice.ErrCannotModifyClosedSession:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete seats from a closed session"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete seat configuration"))
		}
		return
	}

	response.NoContent(c)
}

// GetStats returns statistics for an admission session.
// @Summary Get admission session statistics
// @Description Get statistics for an admission session (applications, seats)
// @Tags Admission Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Success 200 {object} response.Success{data=SessionStatsDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/stats [get]
func (h *SessionHandler) GetStats(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	stats, err := h.sessionService.GetStats(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Admission session not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve session statistics"))
		}
		return
	}

	resp := SessionStatsDTO{
		TotalApplications: stats.TotalApplications,
		ApprovedCount:     stats.ApprovedCount,
		PendingCount:      stats.PendingCount,
		RejectedCount:     stats.RejectedCount,
		TotalSeats:        stats.TotalSeats,
		FilledSeats:       stats.FilledSeats,
		AvailableSeats:    stats.AvailableSeats,
	}

	response.OK(c, resp)
}

// Ensure decimal.Decimal is imported and used to avoid "imported and not used" error.
var _ decimal.Decimal
