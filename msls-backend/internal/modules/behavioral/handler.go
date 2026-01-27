package behavioral

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/response"
)

// Handler handles HTTP requests for behavioral incidents
type Handler struct {
	service *Service
}

// NewHandler creates a new behavioral handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Incident Handlers
// =============================================================================

// ListIncidents godoc
// @Summary List behavioral incidents for a student
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param type query string false "Incident type filter"
// @Param severity query string false "Severity filter"
// @Param dateFrom query string false "Date from (YYYY-MM-DD)"
// @Param dateTo query string false "Date to (YYYY-MM-DD)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response{data=IncidentListResponse}
// @Router /students/{id}/behavioral-incidents [get]
func (h *Handler) ListIncidents(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	var filter IncidentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.ListIncidents(c.Request.Context(), tenantID, studentID, filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// GetIncident godoc
// @Summary Get a behavioral incident
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Success 200 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/behavioral-incidents/{incidentId} [get]
func (h *Handler) GetIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	result, err := h.service.GetIncident(c.Request.Context(), tenantID, studentID, incidentID)
	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// CreateIncident godoc
// @Summary Create a behavioral incident
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body CreateIncidentRequest true "Incident data"
// @Success 201 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/behavioral-incidents [post]
func (h *Handler) CreateIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	var req CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.CreateIncident(c.Request.Context(), tenantID, studentID, userID, req)
	if err != nil {
		if errors.Is(err, ErrInvalidIncidentDate) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// UpdateIncident godoc
// @Summary Update a behavioral incident
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Param request body UpdateIncidentRequest true "Incident data"
// @Success 200 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/behavioral-incidents/{incidentId} [put]
func (h *Handler) UpdateIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	var req UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.UpdateIncident(c.Request.Context(), tenantID, studentID, incidentID, req)
	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidIncidentDate) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// DeleteIncident godoc
// @Summary Delete a behavioral incident
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Success 204 "No Content"
// @Router /students/{id}/behavioral-incidents/{incidentId} [delete]
func (h *Handler) DeleteIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	if err := h.service.DeleteIncident(c.Request.Context(), tenantID, studentID, incidentID); err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBehaviorSummary godoc
// @Summary Get behavior summary for a student
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} response.Response{data=BehaviorSummary}
// @Router /students/{id}/behavioral-summary [get]
func (h *Handler) GetBehaviorSummary(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid student ID")
		return
	}

	result, err := h.service.GetBehaviorSummary(c.Request.Context(), tenantID, studentID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// =============================================================================
// Follow-Up Handlers
// =============================================================================

// CreateFollowUp godoc
// @Summary Create a follow-up for an incident
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param incidentId path string true "Incident ID"
// @Param request body CreateFollowUpRequest true "Follow-up data"
// @Success 201 {object} response.Response{data=FollowUpResponse}
// @Router /behavioral-incidents/{incidentId}/follow-ups [post]
func (h *Handler) CreateFollowUp(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	var req CreateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.CreateFollowUp(c.Request.Context(), tenantID, incidentID, userID, req)
	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidScheduledDate) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// UpdateFollowUp godoc
// @Summary Update a follow-up
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param incidentId path string true "Incident ID"
// @Param followUpId path string true "Follow-up ID"
// @Param request body UpdateFollowUpRequest true "Follow-up data"
// @Success 200 {object} response.Response{data=FollowUpResponse}
// @Router /behavioral-incidents/{incidentId}/follow-ups/{followUpId} [put]
func (h *Handler) UpdateFollowUp(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	followUpID, err := uuid.Parse(c.Param("followUpId"))
	if err != nil {
		response.BadRequest(c, "invalid follow-up ID")
		return
	}

	var req UpdateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.UpdateFollowUp(c.Request.Context(), tenantID, incidentID, followUpID, userID, req)
	if err != nil {
		if errors.Is(err, ErrFollowUpNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidScheduledDate) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// DeleteFollowUp godoc
// @Summary Delete a follow-up
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param incidentId path string true "Incident ID"
// @Param followUpId path string true "Follow-up ID"
// @Success 204 "No Content"
// @Router /behavioral-incidents/{incidentId}/follow-ups/{followUpId} [delete]
func (h *Handler) DeleteFollowUp(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "invalid incident ID")
		return
	}

	followUpID, err := uuid.Parse(c.Param("followUpId"))
	if err != nil {
		response.BadRequest(c, "invalid follow-up ID")
		return
	}

	if err := h.service.DeleteFollowUp(c.Request.Context(), tenantID, incidentID, followUpID); err != nil {
		if errors.Is(err, ErrFollowUpNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPendingFollowUps godoc
// @Summary List all pending follow-ups
// @Tags Behavioral
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response{data=PendingFollowUpsResponse}
// @Router /follow-ups/pending [get]
func (h *Handler) ListPendingFollowUps(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	result, err := h.service.ListPendingFollowUps(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// =============================================================================
// Helpers
// =============================================================================

func (h *Handler) getTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get(middleware.TenantIDKey)
	if !exists {
		return uuid.Nil, errors.New("tenant ID not found in context")
	}
	// Handle both string and uuid.UUID types
	switch v := tenantID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, errors.New("invalid tenant ID type in context")
	}
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	// Handle both string and uuid.UUID types
	switch v := userID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, errors.New("invalid user ID type in context")
	}
}
