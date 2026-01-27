// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/middleware"
	admissionservice "msls-backend/internal/services/admission"
)

// ExportHandler handles admission report export HTTP requests.
type ExportHandler struct {
	exportService *admissionservice.ExportService
}

// NewExportHandler creates a new admission ExportHandler.
func NewExportHandler(exportService *admissionservice.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// ExportQueryParams represents query parameters for export requests.
type ExportQueryParams struct {
	ReportType     string `form:"report_type" binding:"required,oneof=dashboard class-wise funnel source-analysis daily-trend"`
	Format         string `form:"format" binding:"required,oneof=excel pdf"`
	SessionID      string `form:"session_id"`
	BranchID       string `form:"branch_id"`
	AcademicYearID string `form:"academic_year_id"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
}

// Export generates and returns a report export file.
// @Summary Export admission report
// @Description Export admission report to Excel or PDF format
// @Tags Admissions
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Produce application/pdf
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param report_type query string true "Report type to export" Enums(dashboard, class-wise, funnel, source-analysis, daily-trend)
// @Param format query string true "Export format" Enums(excel, pdf)
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {file} binary
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/admissions/export [get]
func (h *ExportHandler) Export(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var params ExportQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Build filter
	filter := admissionservice.DashboardFilter{
		TenantID: tenantID,
	}

	// Parse session ID
	if params.SessionID != "" {
		sessionID, err := uuid.Parse(params.SessionID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid session_id format"))
			return
		}
		filter.SessionID = &sessionID
	}

	// Parse branch ID
	if params.BranchID != "" {
		branchID, err := uuid.Parse(params.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch_id format"))
			return
		}
		filter.BranchID = &branchID
	}

	// Parse dates
	startDate, err := parseDateWithError(params.StartDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid start_date format, expected YYYY-MM-DD"))
		return
	}
	filter.StartDate = startDate

	endDate, err := parseDateWithError(params.EndDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid end_date format, expected YYYY-MM-DD"))
		return
	}
	filter.EndDate = endDate

	// Build export request
	req := admissionservice.ExportRequest{
		ReportType: admissionservice.ExportReportType(params.ReportType),
		Format:     admissionservice.ExportFormat(params.Format),
		Filter:     filter,
	}

	// Generate export
	result, err := h.exportService.ExportReport(c.Request.Context(), req)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to generate export: "+err.Error()))
		}
		return
	}

	// Set response headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+result.Filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Data(200, result.ContentType, result.Data)
}
