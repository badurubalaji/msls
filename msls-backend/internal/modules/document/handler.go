package document

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
)

// Handler handles document-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new document handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Document Type Handlers
// =============================================================================

// ListDocumentTypes godoc
// @Summary List document types
// @Description Get all active document types for the tenant
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param active_only query bool false "Only return active types (default true)"
// @Success 200 {object} response.Success{data=DocumentTypeListResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/document-types [get]
func (h *Handler) ListDocumentTypes(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	// Default to active only
	activeOnly := true
	if c.Query("active_only") == "false" {
		activeOnly = false
	}

	types, err := h.service.ListDocumentTypes(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		logger.Error("Failed to list document types",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to list document types"))
		return
	}

	response.OK(c, DocumentTypeListResponse{
		DocumentTypes: types,
		Total:         len(types),
	})
}

// CreateDocumentType godoc
// @Summary Create document type
// @Description Create a new document type
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateDocumentTypeRequest true "Document type data"
// @Success 201 {object} response.Success{data=DocumentType}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/document-types [post]
func (h *Handler) CreateDocumentType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateDocumentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dt, err := h.service.CreateDocumentType(c.Request.Context(), tenantID, req)
	if err != nil {
		logger.Error("Failed to create document type",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to create document type"))
		return
	}

	response.Created(c, dt)
}

// UpdateDocumentType godoc
// @Summary Update document type
// @Description Update a document type
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Document Type ID" format(uuid)
// @Param request body UpdateDocumentTypeRequest true "Document type updates"
// @Success 200 {object} response.Success{data=DocumentType}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/document-types/{id} [put]
func (h *Handler) UpdateDocumentType(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document type ID"))
		return
	}

	var req UpdateDocumentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dt, err := h.service.UpdateDocumentType(c.Request.Context(), tenantID, id, req)
	if err != nil {
		if errors.Is(err, ErrDocumentTypeNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Document type not found"))
			return
		}
		logger.Error("Failed to update document type",
			zap.String("tenant_id", tenantID.String()),
			zap.String("id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to update document type"))
		return
	}

	response.OK(c, dt)
}

// =============================================================================
// Student Document Handlers
// =============================================================================

// ListDocuments godoc
// @Summary List student documents
// @Description Get all documents for a student
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param status query string false "Filter by status"
// @Success 200 {object} response.Success{data=DocumentListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents [get]
func (h *Handler) ListDocuments(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	filter := DocumentFilter{}
	if statusStr := c.Query("status"); statusStr != "" {
		status := DocumentStatus(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	docs, err := h.service.ListDocuments(c.Request.Context(), tenantID, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		logger.Error("Failed to list documents",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to list documents"))
		return
	}

	response.OK(c, DocumentListResponse{
		Documents: docs,
		Total:     len(docs),
	})
}

// GetDocument godoc
// @Summary Get document
// @Description Get a document by ID
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param docId path string true "Document ID" format(uuid)
// @Success 200 {object} response.Success{data=StudentDocument}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents/{docId} [get]
func (h *Handler) GetDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), tenantID, studentID, docID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
			return
		}
		logger.Error("Failed to get document",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.String("doc_id", docID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to get document"))
		return
	}

	response.OK(c, doc)
}

// UploadDocument godoc
// @Summary Upload document
// @Description Upload a document for a student
// @Tags Documents
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param file formData file true "Document file"
// @Param documentTypeId formData string true "Document Type ID"
// @Param documentNumber formData string false "Document number"
// @Param issueDate formData string false "Issue date (YYYY-MM-DD)"
// @Param expiryDate formData string false "Expiry date (YYYY-MM-DD)"
// @Success 201 {object} response.Success{data=StudentDocument}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents [post]
func (h *Handler) UploadDocument(c *gin.Context) {
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

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	// Get file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("File is required"))
		return
	}
	defer file.Close()

	// Parse document type ID
	docTypeIDStr := c.PostForm("documentTypeId")
	if docTypeIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Document type ID is required"))
		return
	}
	docTypeID, err := uuid.Parse(docTypeIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document type ID"))
		return
	}

	req := UploadDocumentRequest{
		DocumentTypeID: docTypeID,
		DocumentNumber: c.PostForm("documentNumber"),
		IssueDate:      c.PostForm("issueDate"),
		ExpiryDate:     c.PostForm("expiryDate"),
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	doc, err := h.service.UploadDocument(
		c.Request.Context(),
		tenantID,
		studentID,
		userID,
		file,
		header.Filename,
		header.Size,
		mimeType,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		case errors.Is(err, ErrDocumentTypeNotFound):
			apperrors.Abort(c, apperrors.NotFound("Document type not found"))
		case errors.Is(err, ErrFileTooLarge):
			apperrors.Abort(c, apperrors.BadRequest("File size exceeds maximum allowed"))
		case errors.Is(err, ErrInvalidFileType):
			apperrors.Abort(c, apperrors.BadRequest("Invalid file type"))
		case errors.Is(err, ErrInvalidIssueDate):
			apperrors.Abort(c, apperrors.BadRequest("Issue date cannot be in the future"))
		case errors.Is(err, ErrInvalidExpiryDate):
			apperrors.Abort(c, apperrors.BadRequest("Expiry date must be in the future"))
		default:
			logger.Error("Failed to upload document",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to upload document"))
		}
		return
	}

	response.Created(c, doc)
}

// UpdateDocument godoc
// @Summary Update document metadata
// @Description Update document metadata (not the file)
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param docId path string true "Document ID" format(uuid)
// @Param request body UpdateDocumentRequest true "Document updates"
// @Success 200 {object} response.Success{data=StudentDocument}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents/{docId} [put]
func (h *Handler) UpdateDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	doc, err := h.service.UpdateDocument(c.Request.Context(), tenantID, studentID, docID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDocumentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		case errors.Is(err, ErrOptimisticLockConflict):
			apperrors.Abort(c, apperrors.Conflict("Document was modified by another user"))
		case errors.Is(err, ErrInvalidIssueDate):
			apperrors.Abort(c, apperrors.BadRequest("Issue date cannot be in the future"))
		case errors.Is(err, ErrInvalidExpiryDate):
			apperrors.Abort(c, apperrors.BadRequest("Expiry date must be in the future"))
		default:
			logger.Error("Failed to update document",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.String("doc_id", docID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update document"))
		}
		return
	}

	response.OK(c, doc)
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete a document
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param docId path string true "Document ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents/{docId} [delete]
func (h *Handler) DeleteDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), tenantID, studentID, docID); err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
			return
		}
		logger.Error("Failed to delete document",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.String("doc_id", docID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to delete document"))
		return
	}

	response.NoContent(c)
}

// VerifyDocument godoc
// @Summary Verify document
// @Description Mark a document as verified
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param docId path string true "Document ID" format(uuid)
// @Param request body VerifyDocumentRequest true "Verify request"
// @Success 200 {object} response.Success{data=StudentDocument}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents/{docId}/verify [post]
func (h *Handler) VerifyDocument(c *gin.Context) {
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

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	var req VerifyDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	doc, err := h.service.VerifyDocument(c.Request.Context(), tenantID, studentID, docID, userID, req.Version)
	if err != nil {
		switch {
		case errors.Is(err, ErrDocumentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		case errors.Is(err, ErrDocumentAlreadyVerified):
			apperrors.Abort(c, apperrors.BadRequest("Document is already verified"))
		case errors.Is(err, ErrOptimisticLockConflict):
			apperrors.Abort(c, apperrors.Conflict("Document was modified by another user"))
		default:
			logger.Error("Failed to verify document",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.String("doc_id", docID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to verify document"))
		}
		return
	}

	response.OK(c, doc)
}

// RejectDocument godoc
// @Summary Reject document
// @Description Mark a document as rejected
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param docId path string true "Document ID" format(uuid)
// @Param request body RejectDocumentRequest true "Reject request"
// @Success 200 {object} response.Success{data=StudentDocument}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{id}/documents/{docId}/reject [post]
func (h *Handler) RejectDocument(c *gin.Context) {
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

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	var req RejectDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	doc, err := h.service.RejectDocument(c.Request.Context(), tenantID, studentID, docID, userID, req.Reason, req.Version)
	if err != nil {
		switch {
		case errors.Is(err, ErrDocumentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		case errors.Is(err, ErrDocumentAlreadyVerified):
			apperrors.Abort(c, apperrors.BadRequest("Cannot reject a verified document"))
		case errors.Is(err, ErrRejectionReasonRequired):
			apperrors.Abort(c, apperrors.BadRequest("Rejection reason is required"))
		case errors.Is(err, ErrOptimisticLockConflict):
			apperrors.Abort(c, apperrors.Conflict("Document was modified by another user"))
		default:
			logger.Error("Failed to reject document",
				zap.String("tenant_id", tenantID.String()),
				zap.String("student_id", studentID.String()),
				zap.String("doc_id", docID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to reject document"))
		}
		return
	}

	response.OK(c, doc)
}

// GetDocumentChecklist godoc
// @Summary Get document checklist
// @Description Get document checklist with status for a student
// @Tags Documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=DocumentChecklistResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/document-checklist [get]
func (h *Handler) GetDocumentChecklist(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	checklist, err := h.service.GetDocumentChecklist(c.Request.Context(), tenantID, studentID)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		logger.Error("Failed to get document checklist",
			zap.String("tenant_id", tenantID.String()),
			zap.String("student_id", studentID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to get document checklist"))
		return
	}

	response.OK(c, checklist)
}

// RegisterRoutes registers the document routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	// Document types routes (require authentication)
	docTypes := r.Group("/document-types")
	docTypes.Use(authMiddleware)
	{
		docTypes.GET("", h.ListDocumentTypes)
		docTypes.POST("", h.CreateDocumentType)
		docTypes.PUT("/:id", h.UpdateDocumentType)
	}

	// Student documents routes
	students := r.Group("/students/:id")
	students.Use(authMiddleware)
	{
		students.GET("/documents", h.ListDocuments)
		students.POST("/documents", h.UploadDocument)
		students.GET("/documents/:docId", h.GetDocument)
		students.PUT("/documents/:docId", h.UpdateDocument)
		students.DELETE("/documents/:docId", h.DeleteDocument)
		students.POST("/documents/:docId/verify", h.VerifyDocument)
		students.POST("/documents/:docId/reject", h.RejectDocument)
		students.GET("/document-checklist", h.GetDocumentChecklist)
	}
}
