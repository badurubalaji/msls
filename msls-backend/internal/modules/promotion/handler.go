// Package promotion provides student promotion and retention processing functionality.
package promotion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for promotion endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new promotion handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the promotion routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Promotion rules
	rules := rg.Group("/promotion-rules")
	{
		rules.GET("", h.ListRules)
		rules.POST("", h.CreateOrUpdateRule)
		rules.GET("/:id", h.GetRule)
		rules.DELETE("/:id", h.DeleteRule)
	}

	// Promotion batches
	batches := rg.Group("/promotion-batches")
	{
		batches.GET("", h.ListBatches)
		batches.POST("", h.CreateBatch)
		batches.GET("/:id", h.GetBatch)
		batches.DELETE("/:id", h.DeleteBatch)
		batches.POST("/:id/cancel", h.CancelBatch)
		batches.GET("/:id/records", h.ListRecords)
		batches.PUT("/:id/records/:recordId", h.UpdateRecord)
		batches.POST("/:id/records/bulk", h.BulkUpdateRecords)
		batches.POST("/:id/auto-decide", h.AutoDecide)
		batches.POST("/:id/process", h.ProcessBatch)
		batches.GET("/:id/report", h.GetReport)
	}
}

// ========================================================================
// Context Helpers
// ========================================================================

func getTenantID(c *gin.Context) uuid.UUID {
	// Get tenant ID from context (set by auth middleware)
	if tenantID, exists := c.Get("tenantID"); exists {
		if id, ok := tenantID.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

func getUserID(c *gin.Context) *uuid.UUID {
	// Get user ID from context (set by auth middleware)
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return &id
		}
	}
	return nil
}

// ========================================================================
// Promotion Rules Handlers
// ========================================================================

// ListRules lists all promotion rules.
// @Summary List promotion rules
// @Tags Promotion
// @Accept json
// @Produce json
// @Success 200 {object} RuleListResponse
// @Router /promotion-rules [get]
func (h *Handler) ListRules(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rules, err := h.service.ListRules(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RuleListResponse{
		Rules: ToRuleResponses(rules),
		Total: len(rules),
	})
}

// CreateOrUpdateRule creates or updates a promotion rule.
// @Summary Create or update promotion rule
// @Tags Promotion
// @Accept json
// @Produce json
// @Param rule body CreateRuleRequest true "Rule details"
// @Success 200 {object} RuleResponse
// @Router /promotion-rules [post]
func (h *Handler) CreateOrUpdateRule(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	rule, err := h.service.CreateOrUpdateRule(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToRuleResponse(rule))
}

// GetRule retrieves a promotion rule by ID.
// @Summary Get promotion rule
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} RuleResponse
// @Router /promotion-rules/{id} [get]
func (h *Handler) GetRule(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	rule, err := h.service.GetRule(c.Request.Context(), tenantID, id)
	if err != nil {
		if err == ErrRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToRuleResponse(rule))
}

// DeleteRule deletes a promotion rule.
// @Summary Delete promotion rule
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 204
// @Router /promotion-rules/{id} [delete]
func (h *Handler) DeleteRule(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	if err := h.service.DeleteRule(c.Request.Context(), tenantID, id); err != nil {
		if err == ErrRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================================================
// Promotion Batches Handlers
// ========================================================================

// ListBatches lists all promotion batches.
// @Summary List promotion batches
// @Tags Promotion
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {object} BatchListResponse
// @Router /promotion-batches [get]
func (h *Handler) ListBatches(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var status *BatchStatus
	if s := c.Query("status"); s != "" {
		bs := BatchStatus(s)
		status = &bs
	}

	batches, err := h.service.ListBatches(c.Request.Context(), tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, BatchListResponse{
		Batches: ToBatchResponses(batches),
		Total:   len(batches),
	})
}

// CreateBatch creates a new promotion batch.
// @Summary Create promotion batch
// @Tags Promotion
// @Accept json
// @Produce json
// @Param batch body CreateBatchRequest true "Batch details"
// @Success 201 {object} BatchResponse
// @Router /promotion-batches [post]
func (h *Handler) CreateBatch(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	batch, err := h.service.CreateBatch(c.Request.Context(), tenantID, req, userID)
	if err != nil {
		switch err {
		case ErrSameAcademicYear:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case ErrNoStudentsInClass:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, ToBatchResponse(batch))
}

// GetBatch retrieves a promotion batch by ID.
// @Summary Get promotion batch
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Success 200 {object} BatchResponse
// @Router /promotion-batches/{id} [get]
func (h *Handler) GetBatch(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	batch, err := h.service.GetBatch(c.Request.Context(), tenantID, id)
	if err != nil {
		if err == ErrBatchNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToBatchResponse(batch))
}

// DeleteBatch deletes a draft promotion batch.
// @Summary Delete promotion batch
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Success 204
// @Router /promotion-batches/{id} [delete]
func (h *Handler) DeleteBatch(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	if err := h.service.DeleteBatch(c.Request.Context(), tenantID, id); err != nil {
		switch err {
		case ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrBatchNotDraft:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// CancelBatch cancels a promotion batch.
// @Summary Cancel promotion batch
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Param request body CancelBatchRequest true "Cancellation reason"
// @Success 204
// @Router /promotion-batches/{id}/cancel [post]
func (h *Handler) CancelBatch(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	var req CancelBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	if err := h.service.CancelBatch(c.Request.Context(), tenantID, id, req.Reason, userID); err != nil {
		switch err {
		case ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrBatchNotDraft:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================================================
// Promotion Records Handlers
// ========================================================================

// ListRecords lists all promotion records for a batch.
// @Summary List promotion records
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Success 200 {object} RecordListResponse
// @Router /promotion-batches/{id}/records [get]
func (h *Handler) ListRecords(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	records, summary, err := h.service.GetRecordsByBatch(c.Request.Context(), tenantID, batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RecordListResponse{
		Records: ToRecordResponses(records),
		Total:   len(records),
		Summary: summary,
	})
}

// UpdateRecord updates a promotion record.
// @Summary Update promotion record
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Param recordId path string true "Record ID"
// @Param record body UpdateRecordRequest true "Record updates"
// @Success 200 {object} RecordResponse
// @Router /promotion-batches/{id}/records/{recordId} [put]
func (h *Handler) UpdateRecord(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	recordID, err := uuid.Parse(c.Param("recordId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record ID"})
		return
	}

	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	record, err := h.service.UpdateRecord(c.Request.Context(), tenantID, batchID, recordID, req, userID)
	if err != nil {
		switch err {
		case ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrBatchNotDraft:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ToRecordResponse(record))
}

// BulkUpdateRecords updates multiple promotion records.
// @Summary Bulk update promotion records
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Param request body BulkUpdateRecordsRequest true "Bulk update request"
// @Success 204
// @Router /promotion-batches/{id}/records/bulk [post]
func (h *Handler) BulkUpdateRecords(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	var req BulkUpdateRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	if err := h.service.BulkUpdateRecords(c.Request.Context(), tenantID, batchID, req, userID); err != nil {
		switch err {
		case ErrBatchNotDraft:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ========================================================================
// Processing Handlers
// ========================================================================

// AutoDecide applies promotion rules to auto-decide pending records.
// @Summary Auto-decide promotion records
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Success 204
// @Router /promotion-batches/{id}/auto-decide [post]
func (h *Handler) AutoDecide(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	if err := h.service.AutoDecide(c.Request.Context(), tenantID, batchID); err != nil {
		switch err {
		case ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrBatchNotDraft:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ProcessBatch processes a promotion batch.
// @Summary Process promotion batch
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Param request body ProcessBatchRequest true "Process options"
// @Success 200 {object} BatchResponse
// @Router /promotion-batches/{id}/process [post]
func (h *Handler) ProcessBatch(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	var req ProcessBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Default to not generating roll numbers if no body
		req = ProcessBatchRequest{GenerateRollNumbers: false}
	}

	userID := getUserID(c)
	if err := h.service.ProcessBatch(c.Request.Context(), tenantID, batchID, req.GenerateRollNumbers, userID); err != nil {
		switch err {
		case ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrBatchNotDraft, ErrBatchNotProcessable, ErrPendingDecisions, ErrMissingTargetClass:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case ErrBatchAlreadyProcessed, ErrBatchCancelled:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Fetch updated batch
	batch, err := h.service.GetBatch(c.Request.Context(), tenantID, batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToBatchResponse(batch))
}

// GetReport generates a promotion report.
// @Summary Get promotion report
// @Tags Promotion
// @Accept json
// @Produce json
// @Param id path string true "Batch ID"
// @Success 200 {array} PromotionReportRow
// @Router /promotion-batches/{id}/report [get]
func (h *Handler) GetReport(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	batchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch ID"})
		return
	}

	rows, err := h.service.GetPromotionReport(c.Request.Context(), tenantID, batchID)
	if err != nil {
		if err == ErrBatchNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rows": rows, "total": len(rows)})
}
