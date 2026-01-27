// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	admissionservice "msls-backend/internal/services/admission"
)

// MeritHandler handles merit list HTTP requests.
type MeritHandler struct {
	meritService *admissionservice.MeritService
}

// NewMeritHandler creates a new MeritHandler.
func NewMeritHandler(meritService *admissionservice.MeritService) *MeritHandler {
	return &MeritHandler{meritService: meritService}
}

// =============================================================================
// Merit List Endpoints
// =============================================================================

// GenerateMeritList generates a merit list for a session and class.
// @Summary Generate merit list
// @Description Generate a merit list based on test scores for a session and class
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param request body GenerateMeritListRequest true "Merit list parameters"
// @Success 201 {object} response.Success{data=MeritListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/merit-list [post]
func (h *MeritHandler) GenerateMeritList(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	sessionIDParam := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	var req GenerateMeritListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse optional test ID
	var testID *uuid.UUID
	if req.TestID != nil {
		tid, err := uuid.Parse(*req.TestID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
			return
		}
		testID = &tid
	}

	generateReq := admissionservice.GenerateMeritListRequest{
		TenantID:    tenantID,
		SessionID:   sessionID,
		ClassName:   req.ClassName,
		TestID:      testID,
		CutoffScore: req.CutoffScore,
		GeneratedBy: &userID,
	}

	meritList, err := h.meritService.GenerateMeritList(c.Request.Context(), generateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		case admissionservice.ErrSessionIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Session ID is required"))
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Session not found"))
		case admissionservice.ErrClassNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Class name is required"))
		case admissionservice.ErrMeritListFinalized:
			apperrors.Abort(c, apperrors.BadRequest("Merit list is already finalized"))
		case admissionservice.ErrNoApplicantsForMeritList:
			apperrors.Abort(c, apperrors.BadRequest("No applicants found to generate merit list"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to generate merit list"))
		}
		return
	}

	response.Created(c, meritListToResponse(meritList))
}

// GetMeritList retrieves the merit list for a session.
// @Summary Get merit list
// @Description Get the merit list for a session with optional class filter
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param className query string false "Filter by class name"
// @Param testId query string false "Filter by test ID"
// @Success 200 {object} response.Success{data=MeritListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/merit-list [get]
func (h *MeritHandler) GetMeritList(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sessionIDParam := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	// Parse optional test ID
	var testID *uuid.UUID
	if testIDStr := c.Query("testId"); testIDStr != "" {
		tid, err := uuid.Parse(testIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
			return
		}
		testID = &tid
	}

	getReq := admissionservice.GetMeritListRequest{
		TenantID:  tenantID,
		SessionID: sessionID,
		ClassName: c.Query("className"),
		TestID:    testID,
	}

	meritList, err := h.meritService.GetMeritList(c.Request.Context(), getReq)
	if err != nil {
		switch err {
		case admissionservice.ErrMeritListNotFound:
			apperrors.Abort(c, apperrors.NotFound("Merit list not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve merit list"))
		}
		return
	}

	response.OK(c, meritListToResponse(meritList))
}

// ListMeritLists retrieves all merit lists for a session.
// @Summary List merit lists
// @Description Get all merit lists for a session
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Session ID" format(uuid)
// @Param className query string false "Filter by class name"
// @Param isFinal query boolean false "Filter by finalization status"
// @Success 200 {object} response.Success{data=MeritListListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/admission-sessions/{id}/merit-lists [get]
func (h *MeritHandler) ListMeritLists(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	sessionIDParam := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	filter := admissionservice.MeritListFilter{
		TenantID:  tenantID,
		SessionID: sessionID,
		ClassName: c.Query("className"),
	}

	// Parse optional isFinal filter
	if isFinalStr := c.Query("isFinal"); isFinalStr != "" {
		isFinal := isFinalStr == "true"
		filter.IsFinal = &isFinal
	}

	lists, err := h.meritService.ListMeritLists(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve merit lists"))
		return
	}

	resp := MeritListListResponse{
		MeritLists: meritListsToResponses(lists),
		Total:      len(lists),
	}

	response.OK(c, resp)
}

// FinalizeMeritList marks a merit list as final.
// @Summary Finalize merit list
// @Description Mark a merit list as final (cannot be regenerated)
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Merit List ID" format(uuid)
// @Success 200 {object} response.Success{data=MeritListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/merit-lists/{id}/finalize [post]
func (h *MeritHandler) FinalizeMeritList(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid merit list ID"))
		return
	}

	meritList, err := h.meritService.FinalizeMeritList(c.Request.Context(), tenantID, id, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrMeritListNotFound:
			apperrors.Abort(c, apperrors.NotFound("Merit list not found"))
		case admissionservice.ErrMeritListFinalized:
			apperrors.Abort(c, apperrors.BadRequest("Merit list is already finalized"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to finalize merit list"))
		}
		return
	}

	response.OK(c, meritListToResponse(meritList))
}

// UpdateCutoff updates the cutoff score of a merit list.
// @Summary Update cutoff score
// @Description Update the cutoff score of a merit list
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Merit List ID" format(uuid)
// @Param request body UpdateCutoffRequest true "Cutoff update"
// @Success 200 {object} response.Success{data=MeritListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/merit-lists/{id}/cutoff [patch]
func (h *MeritHandler) UpdateCutoff(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid merit list ID"))
		return
	}

	var req UpdateCutoffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	meritList, err := h.meritService.UpdateMeritListCutoff(c.Request.Context(), tenantID, id, req.CutoffScore)
	if err != nil {
		switch err {
		case admissionservice.ErrMeritListNotFound:
			apperrors.Abort(c, apperrors.NotFound("Merit list not found"))
		case admissionservice.ErrMeritListFinalized:
			apperrors.Abort(c, apperrors.BadRequest("Merit list is finalized"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update cutoff score"))
		}
		return
	}

	response.OK(c, meritListToResponse(meritList))
}

// DeleteMeritList deletes a merit list.
// @Summary Delete merit list
// @Description Delete a non-finalized merit list
// @Tags Merit Lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Merit List ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/merit-lists/{id} [delete]
func (h *MeritHandler) DeleteMeritList(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid merit list ID"))
		return
	}

	err = h.meritService.DeleteMeritList(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admissionservice.ErrMeritListNotFound:
			apperrors.Abort(c, apperrors.NotFound("Merit list not found"))
		case admissionservice.ErrMeritListFinalized:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete finalized merit list"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete merit list"))
		}
		return
	}

	response.NoContent(c)
}
