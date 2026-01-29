// Package staffdocument provides staff document management functionality.
package staffdocument

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/storage"
)

// Handler handles HTTP requests for staff documents.
type Handler struct {
	service *Service
	storage storage.Storage
}

// NewHandler creates a new staff document handler.
func NewHandler(service *Service, storage storage.Storage) *Handler {
	return &Handler{
		service: service,
		storage: storage,
	}
}

// Helper functions to get tenant and user IDs from context.
func getTenantID(c *gin.Context) uuid.UUID {
	tenantID, _ := middleware.GetCurrentTenantID(c)
	return tenantID
}

func getUserID(c *gin.Context) uuid.UUID {
	userID, _ := middleware.GetCurrentUserID(c)
	return userID
}

// ============================================================================
// Document Type Handlers
// ============================================================================

// ListDocumentTypes handles GET /document-types
func (h *Handler) ListDocumentTypes(c *gin.Context) {
	tenantID := getTenantID(c)
	activeOnly := c.Query("active_only") == "true"

	types, err := h.service.ListDocumentTypes(c.Request.Context(), tenantID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DocumentTypeListResponse{
		DocumentTypes: ToDocumentTypeResponses(types),
		Total:         int64(len(types)),
	})
}

// CreateDocumentType handles POST /document-types
func (h *Handler) CreateDocumentType(c *gin.Context) {
	tenantID := getTenantID(c)

	var dto CreateDocumentTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt, err := h.service.CreateDocumentType(c.Request.Context(), tenantID, dto)
	if err != nil {
		if errors.Is(err, ErrDuplicateDocTypeCode) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToDocumentTypeResponse(dt))
}

// GetDocumentType handles GET /document-types/:id
func (h *Handler) GetDocumentType(c *gin.Context) {
	tenantID := getTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type ID"})
		return
	}

	dt, err := h.service.GetDocumentType(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrDocumentTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentTypeResponse(dt))
}

// UpdateDocumentType handles PUT /document-types/:id
func (h *Handler) UpdateDocumentType(c *gin.Context) {
	tenantID := getTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type ID"})
		return
	}

	var dto UpdateDocumentTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt, err := h.service.UpdateDocumentType(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrDocumentTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentTypeResponse(dt))
}

// DeleteDocumentType handles DELETE /document-types/:id
func (h *Handler) DeleteDocumentType(c *gin.Context) {
	tenantID := getTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type ID"})
		return
	}

	if err := h.service.DeleteDocumentType(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrDocumentTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrDocumentTypeInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ============================================================================
// Document Handlers
// ============================================================================

// ListStaffDocuments handles GET /staff/:id/documents
func (h *Handler) ListStaffDocuments(c *gin.Context) {
	tenantID := getTenantID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}

	params := ListDocumentsParams{
		Cursor: c.Query("cursor"),
	}

	if docTypeID := c.Query("document_type_id"); docTypeID != "" {
		id, err := uuid.Parse(docTypeID)
		if err == nil {
			params.DocumentTypeID = &id
		}
	}
	if status := c.Query("status"); status != "" {
		params.Status = status
	}
	if isCurrent := c.Query("is_current"); isCurrent != "" {
		current := isCurrent == "true"
		params.IsCurrent = &current
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			params.Limit = l
		}
	}
	if params.Limit == 0 {
		params.Limit = 50
	}

	docs, total, err := h.service.ListStaffDocuments(c.Request.Context(), tenantID, staffID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor string
	if len(docs) == params.Limit {
		nextCursor = docs[len(docs)-1].ID.String()
	}

	c.JSON(http.StatusOK, DocumentListResponse{
		Documents:  ToDocumentResponses(docs),
		NextCursor: nextCursor,
		HasMore:    len(docs) == params.Limit,
		Total:      total,
	})
}

// UploadDocument handles POST /staff/:id/documents
func (h *Handler) UploadDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large or invalid form"})
		return
	}

	// Get file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// Get metadata from form
	documentTypeID, err := uuid.Parse(c.PostForm("document_type_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_type_id"})
		return
	}

	dto := CreateDocumentDTO{
		DocumentTypeID: documentTypeID,
	}
	if docNum := c.PostForm("document_number"); docNum != "" {
		dto.DocumentNumber = &docNum
	}
	if remarks := c.PostForm("remarks"); remarks != "" {
		dto.Remarks = &remarks
	}
	// Issue date and expiry date would need proper parsing

	// Detect MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	mimeType := http.DetectContentType(buffer)

	// Reset file reader
	if _, err := file.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process file"})
		return
	}

	// Validate file
	if err := h.service.ValidateFileUpload(header.Filename, int(header.Size), mimeType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate document ID and file path
	docID := uuid.New()
	filePath := h.service.GenerateFilePath(tenantID, staffID, docID, header.Filename)

	// Upload to storage
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	if err := h.storage.Upload(c.Request.Context(), filePath, fileBytes, mimeType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	// Create document record
	doc, err := h.service.CreateDocument(c.Request.Context(), tenantID, staffID, dto, header.Filename, filePath, int(header.Size), mimeType, &userID)
	if err != nil {
		// Cleanup uploaded file on failure
		_ = h.storage.Delete(c.Request.Context(), filePath)

		if errors.Is(err, ErrDocumentTypeNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToDocumentResponse(doc))
}

// GetDocument handles GET /staff/:id/documents/:docId
func (h *Handler) GetDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), tenantID, staffID, docID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentResponse(doc))
}

// UpdateDocument handles PUT /staff/:id/documents/:docId
func (h *Handler) UpdateDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var dto UpdateDocumentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.UpdateDocument(c.Request.Context(), tenantID, staffID, docID, dto)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentResponse(doc))
}

// DeleteDocument handles DELETE /staff/:id/documents/:docId
func (h *Handler) DeleteDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	doc, err := h.service.DeleteDocument(c.Request.Context(), tenantID, staffID, docID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete file from storage
	if err := h.storage.Delete(c.Request.Context(), doc.FilePath); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to delete file %s: %v\n", doc.FilePath, err)
	}

	c.JSON(http.StatusNoContent, nil)
}

// DownloadDocument handles GET /staff/:id/documents/:docId/download
func (h *Handler) DownloadDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), tenantID, staffID, docID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get presigned URL or download file
	url, err := h.storage.GetPresignedURL(c.Request.Context(), doc.FilePath, doc.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// ============================================================================
// Verification Handlers
// ============================================================================

// VerifyDocument handles PUT /staff/:id/documents/:docId/verify
func (h *Handler) VerifyDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var dto VerifyDocumentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.VerifyDocument(c.Request.Context(), tenantID, staffID, docID, userID, dto)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrDocumentAlreadyVerified) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentResponse(doc))
}

// RejectDocument handles PUT /staff/:id/documents/:docId/reject
func (h *Handler) RejectDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}
	docID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var dto RejectDocumentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.service.RejectDocument(c.Request.Context(), tenantID, staffID, docID, userID, dto)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrDocumentAlreadyRejected) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrMissingRejectionReason) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToDocumentResponse(doc))
}

// ============================================================================
// Expiry & Compliance Handlers
// ============================================================================

// GetExpiringDocuments handles GET /documents/expiring
func (h *Handler) GetExpiringDocuments(c *gin.Context) {
	tenantID := getTenantID(c)

	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	docs, err := h.service.GetExpiringDocuments(c.Request.Context(), tenantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs, "days": days})
}

// GetComplianceReport handles GET /documents/compliance
func (h *Handler) GetComplianceReport(c *gin.Context) {
	tenantID := getTenantID(c)
	includeDetails := c.Query("include_staff_details") == "true"

	report, err := h.service.GetComplianceReport(c.Request.Context(), tenantID, includeDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// ============================================================================
// Route Registration
// ============================================================================

// RegisterDocumentTypeRoutes registers document type routes.
func (h *Handler) RegisterDocumentTypeRoutes(rg *gin.RouterGroup) {
	docTypes := rg.Group("/staff-document-types")
	{
		docTypes.GET("", h.ListDocumentTypes)
		docTypes.POST("", h.CreateDocumentType)
		docTypes.GET("/:id", h.GetDocumentType)
		docTypes.PUT("/:id", h.UpdateDocumentType)
		docTypes.DELETE("/:id", h.DeleteDocumentType)
	}
}

// RegisterStaffDocumentRoutes registers staff document routes.
func (h *Handler) RegisterStaffDocumentRoutes(staffRoutes *gin.RouterGroup) {
	staffRoutes.GET("/:id/documents", h.ListStaffDocuments)
	staffRoutes.POST("/:id/documents", h.UploadDocument)
	staffRoutes.GET("/:id/documents/:docId", h.GetDocument)
	staffRoutes.PUT("/:id/documents/:docId", h.UpdateDocument)
	staffRoutes.DELETE("/:id/documents/:docId", h.DeleteDocument)
	staffRoutes.GET("/:id/documents/:docId/download", h.DownloadDocument)
	staffRoutes.PUT("/:id/documents/:docId/verify", h.VerifyDocument)
	staffRoutes.PUT("/:id/documents/:docId/reject", h.RejectDocument)
}

// RegisterGlobalDocumentRoutes registers global document routes.
func (h *Handler) RegisterGlobalDocumentRoutes(rg *gin.RouterGroup) {
	docs := rg.Group("/staff-documents")
	{
		docs.GET("/expiring", h.GetExpiringDocuments)
		docs.GET("/compliance", h.GetComplianceReport)
	}
}
