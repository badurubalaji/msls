// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
)

// Handler handles HTTP requests for hall tickets.
type Handler struct {
	service *Service
}

// NewHandler creates a new hall ticket handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the hall ticket routes.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Hall ticket routes under examinations
	exams := r.Group("/examinations/:examId/hall-tickets")
	exams.Use(authMiddleware)
	{
		exams.GET("", middleware.PermissionRequired("hall-ticket:view"), h.ListHallTickets)
		exams.POST("/generate", middleware.PermissionRequired("hall-ticket:generate"), h.GenerateHallTickets)
		exams.GET("/pdf", middleware.PermissionRequired("hall-ticket:download"), h.DownloadBatchPDF)
		exams.GET("/:ticketId", middleware.PermissionRequired("hall-ticket:view"), h.GetHallTicket)
		exams.GET("/:ticketId/pdf", middleware.PermissionRequired("hall-ticket:download"), h.DownloadSinglePDF)
		exams.DELETE("/:ticketId", middleware.PermissionRequired("hall-ticket:generate"), h.DeleteHallTicket)
	}

	// Public verification endpoint
	r.GET("/hall-tickets/verify/:qrCode", h.VerifyHallTicket)

	// Template routes
	templates := r.Group("/hall-ticket-templates")
	templates.Use(authMiddleware)
	{
		templates.GET("", middleware.PermissionRequired("hall-ticket:view"), h.ListTemplates)
		templates.POST("", middleware.PermissionRequired("hall-ticket:template-manage"), h.CreateTemplate)
		templates.GET("/:id", middleware.PermissionRequired("hall-ticket:view"), h.GetTemplate)
		templates.PUT("/:id", middleware.PermissionRequired("hall-ticket:template-manage"), h.UpdateTemplate)
		templates.DELETE("/:id", middleware.PermissionRequired("hall-ticket:template-manage"), h.DeleteTemplate)
	}
}

// getTenantUUID extracts and parses the tenant ID from context.
func getTenantUUID(c *gin.Context) (uuid.UUID, bool) {
	tenantIDStr := middleware.GetTenantID(c)
	if tenantIDStr == "" {
		return uuid.Nil, false
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, false
	}
	return tenantID, true
}

// ListHallTickets handles GET /examinations/:examId/hall-tickets
func (h *Handler) ListHallTickets(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	examID, err := uuid.Parse(c.Param("examId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid examination ID"})
		return
	}

	filter := ListFilter{
		TenantID:      tenantID,
		ExaminationID: examID,
		Search:        c.Query("search"),
	}

	if classID := c.Query("classId"); classID != "" {
		id, err := uuid.Parse(classID)
		if err == nil {
			filter.ClassID = &id
		}
	}

	if sectionID := c.Query("sectionId"); sectionID != "" {
		id, err := uuid.Parse(sectionID)
		if err == nil {
			filter.SectionID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		s := HallTicketStatus(status)
		filter.Status = &s
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	} else {
		filter.Limit = 50
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	tickets, total, err := h.service.ListHallTickets(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tickets,
		"meta": gin.H{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

// GenerateHallTickets handles POST /examinations/:examId/hall-tickets/generate
func (h *Handler) GenerateHallTickets(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	examID, err := uuid.Parse(c.Param("examId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid examination ID"})
		return
	}

	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body - will generate for all classes
	}

	req.TenantID = tenantID
	req.ExaminationID = examID

	if userID, exists := middleware.GetCurrentUserID(c); exists {
		req.CreatedBy = &userID
	}

	result, err := h.service.GenerateHallTickets(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrExaminationNotScheduled, ErrNoClassesForExam, ErrNoStudentsInClass:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetHallTicket handles GET /examinations/:examId/hall-tickets/:ticketId
func (h *Handler) GetHallTicket(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	ticket, err := h.service.GetHallTicket(c.Request.Context(), tenantID, ticketID)
	if err != nil {
		if err == ErrHallTicketNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// DownloadSinglePDF handles GET /examinations/:examId/hall-tickets/:ticketId/pdf
func (h *Handler) DownloadSinglePDF(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	examID, err := uuid.Parse(c.Param("examId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid examination ID"})
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	pdf, filename, err := h.service.GetHallTicketPDF(c.Request.Context(), tenantID, examID, ticketID)
	if err != nil {
		if err == ErrHallTicketNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", pdf)
}

// DownloadBatchPDF handles GET /examinations/:examId/hall-tickets/pdf
func (h *Handler) DownloadBatchPDF(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	examID, err := uuid.Parse(c.Param("examId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid examination ID"})
		return
	}

	var classID *uuid.UUID
	if cid := c.Query("classId"); cid != "" {
		id, err := uuid.Parse(cid)
		if err == nil {
			classID = &id
		}
	}

	pdf, filename, err := h.service.GetBatchPDF(c.Request.Context(), tenantID, examID, classID)
	if err != nil {
		if err == ErrHallTicketNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "no hall tickets found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", pdf)
}

// DeleteHallTicket handles DELETE /examinations/:examId/hall-tickets/:ticketId
func (h *Handler) DeleteHallTicket(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	if err := h.service.DeleteHallTicket(c.Request.Context(), tenantID, ticketID); err != nil {
		if err == ErrHallTicketNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// VerifyHallTicket handles GET /hall-tickets/verify/:qrCode
func (h *Handler) VerifyHallTicket(c *gin.Context) {
	qrCode := c.Param("qrCode")
	if qrCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code is required"})
		return
	}

	result, err := h.service.VerifyHallTicket(c.Request.Context(), qrCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Template handlers

// ListTemplates handles GET /hall-ticket-templates
func (h *Handler) ListTemplates(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	templates, err := h.service.ListTemplates(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// CreateTemplate handles POST /hall-ticket-templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = tenantID

	if userID, exists := middleware.GetCurrentUserID(c); exists {
		req.CreatedBy = &userID
	}

	template, err := h.service.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": template})
}

// GetTemplate handles GET /hall-ticket-templates/:id
func (h *Handler) GetTemplate(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	template, err := h.service.GetTemplate(c.Request.Context(), tenantID, id)
	if err != nil {
		if err == ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// UpdateTemplate handles PUT /hall-ticket-templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID, exists := middleware.GetCurrentUserID(c); exists {
		req.UpdatedBy = &userID
	}

	template, err := h.service.UpdateTemplate(c.Request.Context(), tenantID, id, &req)
	if err != nil {
		if err == ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// DeleteTemplate handles DELETE /hall-ticket-templates/:id
func (h *Handler) DeleteTemplate(c *gin.Context) {
	tenantID, ok := getTenantUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	if err := h.service.DeleteTemplate(c.Request.Context(), tenantID, id); err != nil {
		if err == ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
