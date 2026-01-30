// Package bulk provides bulk operation functionality.
package bulk

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles bulk operation HTTP requests.
type Handler struct {
	service       *Service
	importService *ImportService
}

// NewHandler creates a new bulk operation handler.
func NewHandler(service *Service, importService *ImportService) *Handler {
	return &Handler{
		service:       service,
		importService: importService,
	}
}

// BulkStatusUpdateRequest represents a request to bulk update student status.
type BulkStatusUpdateRequest struct {
	StudentIDs []string `json:"studentIds" binding:"required,min=1"`
	NewStatus  string   `json:"newStatus" binding:"required,oneof=active inactive transferred graduated"`
}

// ExportRequest represents a request to export students.
type ExportRequest struct {
	StudentIDs []string `json:"studentIds" binding:"required,min=1"`
	Format     string   `json:"format" binding:"required,oneof=xlsx csv"`
	Columns    []string `json:"columns"`
}

// BulkStatusUpdate performs a bulk status update on students.
// @Summary Bulk update student status
// @Description Update the status of multiple students at once
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body BulkStatusUpdateRequest true "Bulk status update details"
// @Success 200 {object} response.Success{data=BulkOperationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/students/bulk/status [post]
func (h *Handler) BulkStatusUpdate(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req BulkStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse student IDs
	studentIDs := make([]uuid.UUID, 0, len(req.StudentIDs))
	for _, idStr := range req.StudentIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid student ID: "+idStr))
			return
		}
		studentIDs = append(studentIDs, id)
	}

	newStatus := models.StudentStatus(req.NewStatus)
	if !newStatus.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		return
	}

	// Create bulk operation
	dto := CreateBulkOperationDTO{
		TenantID:   tenantID,
		StudentIDs: studentIDs,
		CreatedBy:  userID,
	}

	op, err := h.service.CreateBulkStatusUpdate(c.Request.Context(), dto, newStatus)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Process synchronously (for small batches)
	// For large batches, this should be done asynchronously
	if err := h.service.ProcessStatusUpdate(c.Request.Context(), tenantID, op.ID, newStatus); err != nil {
		logger.Error("Failed to process bulk status update",
			zap.String("operation_id", op.ID.String()),
			zap.Error(err))
	}

	// Get updated operation
	op, err = h.service.GetByID(c.Request.Context(), tenantID, op.ID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve operation status"))
		return
	}

	response.OK(c, ToBulkOperationResponse(op))
}

// Export exports students to Excel or CSV.
// @Summary Export students
// @Description Export selected students to Excel or CSV format
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body ExportRequest true "Export details"
// @Success 200 {object} response.Success{data=BulkOperationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/students/export [post]
func (h *Handler) Export(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse student IDs
	studentIDs := make([]uuid.UUID, 0, len(req.StudentIDs))
	for _, idStr := range req.StudentIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid student ID: "+idStr))
			return
		}
		studentIDs = append(studentIDs, id)
	}

	params := ExportParams{
		Format:  req.Format,
		Columns: req.Columns,
	}

	// Create bulk operation
	dto := CreateBulkOperationDTO{
		TenantID:   tenantID,
		StudentIDs: studentIDs,
		CreatedBy:  userID,
	}

	op, err := h.service.CreateExport(c.Request.Context(), dto, params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Process export synchronously
	_, err = h.service.ProcessExport(c.Request.Context(), tenantID, op.ID, params)
	if err != nil {
		logger.Error("Failed to process export",
			zap.String("operation_id", op.ID.String()),
			zap.Error(err))
	}

	// Get updated operation
	op, err = h.service.GetByID(c.Request.Context(), tenantID, op.ID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve operation status"))
		return
	}

	response.OK(c, ToBulkOperationResponse(op))
}

// GetOperation retrieves a bulk operation by ID.
// @Summary Get bulk operation status
// @Description Get the status and details of a bulk operation
// @Tags Bulk Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Operation ID" format(uuid)
// @Success 200 {object} response.Success{data=BulkOperationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/bulk-operations/{id} [get]
func (h *Handler) GetOperation(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid operation ID"))
		return
	}

	op, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrOperationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Operation not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve operation"))
		return
	}

	response.OK(c, ToBulkOperationResponse(op))
}

// ListOperations lists bulk operations for the current user.
// @Summary List bulk operations
// @Description Get a list of bulk operations created by the current user
// @Tags Bulk Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param limit query int false "Number of results (default 20)"
// @Success 200 {object} response.Success{data=[]BulkOperationResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/bulk-operations [get]
func (h *Handler) ListOperations(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		var parsed int
		for _, ch := range limitStr {
			if ch >= '0' && ch <= '9' {
				parsed = parsed*10 + int(ch-'0')
			}
		}
		if parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ops, err := h.service.ListByUser(c.Request.Context(), tenantID, userID, limit)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve operations"))
		return
	}

	respList := make([]BulkOperationResponse, len(ops))
	for i, op := range ops {
		respList[i] = ToBulkOperationResponse(&op)
	}

	response.OK(c, respList)
}

// DownloadResult redirects to the export file download.
// @Summary Download export result
// @Description Download the result file of an export operation
// @Tags Bulk Operations
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Operation ID" format(uuid)
// @Success 302 "Redirect to file"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/bulk-operations/{id}/result [get]
func (h *Handler) DownloadResult(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid operation ID"))
		return
	}

	op, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrOperationNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Operation not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve operation"))
		return
	}

	if op.ResultURL == "" {
		apperrors.Abort(c, apperrors.NotFound("No result file available"))
		return
	}

	c.Redirect(http.StatusFound, op.ResultURL)
}

// DownloadTemplate generates and downloads the student import template.
// @Summary Download student import template
// @Description Download an Excel template for bulk student import
// @Tags Students
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {file} binary "Excel template file"
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/students/import/template [get]
func (h *Handler) DownloadTemplate(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	templateBytes, err := h.importService.GenerateTemplate(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error("Failed to generate import template", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to generate template"))
		return
	}

	c.Header("Content-Disposition", "attachment; filename=student_import_template.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", templateBytes)
}

// ImportStudents imports students from an uploaded file.
// @Summary Import students from file
// @Description Bulk import students from Excel or CSV file
// @Tags Students
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param file formance file true "Excel or CSV file"
// @Param branch_id formData string true "Branch ID"
// @Param academic_year_id formData string true "Academic Year ID"
// @Success 200 {object} response.Success{data=ImportResult}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/students/import [post]
func (h *Handler) ImportStudents(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	// Parse form
	branchIDStr := c.PostForm("branch_id")
	if branchIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Branch ID is required"))
		return
	}
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	academicYearIDStr := c.PostForm("academic_year_id")
	if academicYearIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Academic Year ID is required"))
		return
	}
	academicYearID, err := uuid.Parse(academicYearIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid academic year ID"))
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("File is required"))
		return
	}
	defer file.Close()

	// Determine file type
	fileType := "xlsx"
	if strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		fileType = "csv"
	}

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Failed to read file"))
		return
	}

	// Process import
	result, err := h.importService.ImportStudents(c.Request.Context(), tenantID, branchID, academicYearID, userID, fileContent, fileType)
	if err != nil {
		logger.Error("Failed to import students", zap.Error(err))
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	response.OK(c, result)
}

// handleServiceError converts service errors to API errors.
func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNoStudentsProvided):
		apperrors.Abort(c, apperrors.BadRequest("No students provided"))
	case errors.Is(err, ErrTooManyStudents):
		apperrors.Abort(c, apperrors.BadRequest("Too many students provided"))
	case errors.Is(err, ErrInvalidExportFormat):
		apperrors.Abort(c, apperrors.BadRequest("Invalid export format"))
	case errors.Is(err, ErrInvalidOperationType):
		apperrors.Abort(c, apperrors.BadRequest("Invalid operation type"))
	default:
		logger.Error("Bulk operation error", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to process bulk operation"))
	}
}
