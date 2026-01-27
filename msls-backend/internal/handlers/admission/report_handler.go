// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	admissionservice "msls-backend/internal/services/admission"
)

// ReportHandler handles admission report-related HTTP requests.
type ReportHandler struct {
	reportService *admissionservice.ReportService
}

// NewReportHandler creates a new admission ReportHandler.
func NewReportHandler(reportService *admissionservice.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetDashboard returns admission dashboard statistics.
// @Summary Get admission dashboard
// @Description Get admission dashboard statistics including totals and conversion rates
// @Tags Admissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {object} response.Success{data=DashboardStatsResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admissions/dashboard [get]
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter, err := h.parseQueryParams(c, tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	stats, err := h.reportService.GetDashboardStats(c.Request.Context(), filter)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve dashboard statistics"))
		}
		return
	}

	response.OK(c, dashboardStatsToResponse(stats))
}

// GetFunnel returns the admission conversion funnel.
// @Summary Get admission funnel
// @Description Get admission conversion funnel from enquiry to enrollment
// @Tags Admissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} response.Success{data=FunnelReportResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admissions/reports/funnel [get]
func (h *ReportHandler) GetFunnel(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter, err := h.parseQueryParams(c, tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	report, err := h.reportService.GetFunnelReport(c.Request.Context(), filter)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve funnel report"))
		}
		return
	}

	response.OK(c, funnelReportToResponse(report))
}

// GetClassWise returns the class-wise admission report.
// @Summary Get class-wise admission report
// @Description Get admission statistics broken down by class
// @Tags Admissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} response.Success{data=ClassWiseReportResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admissions/reports/class-wise [get]
func (h *ReportHandler) GetClassWise(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter, err := h.parseQueryParams(c, tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	report, err := h.reportService.GetClassWiseReport(c.Request.Context(), filter)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve class-wise report"))
		}
		return
	}

	response.OK(c, classWiseReportToResponse(report))
}

// GetSourceAnalysis returns the enquiry source analysis.
// @Summary Get enquiry source analysis
// @Description Get analysis of enquiry sources (walk-in, website, referral, etc.)
// @Tags Admissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {object} response.Success{data=SourceAnalysisResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admissions/reports/source-analysis [get]
func (h *ReportHandler) GetSourceAnalysis(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter, err := h.parseQueryParams(c, tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	report, err := h.reportService.GetSourceAnalysis(c.Request.Context(), filter)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve source analysis"))
		}
		return
	}

	response.OK(c, sourceAnalysisToResponse(report))
}

// GetDailyTrend returns the daily application trend.
// @Summary Get daily application trend
// @Description Get daily trend of applications and enquiries
// @Tags Admissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param session_id query string false "Filter by admission session ID"
// @Param branch_id query string false "Filter by branch ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param days query int false "Number of days to include (default 30)"
// @Success 200 {object} response.Success{data=DailyTrendResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admissions/reports/daily-trend [get]
func (h *ReportHandler) GetDailyTrend(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter, err := h.parseQueryParams(c, tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse days parameter
	var params ReportQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	days := params.Days
	if days <= 0 {
		days = 30 // Default to 30 days
	}

	report, err := h.reportService.GetDailyTrend(c.Request.Context(), filter, days)
	if err != nil {
		switch err {
		case admissionservice.ErrTenantIDRequired:
			apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve daily trend"))
		}
		return
	}

	response.OK(c, dailyTrendToResponse(report))
}

// parseQueryParams extracts and validates query parameters for report filters.
func (h *ReportHandler) parseQueryParams(c *gin.Context, tenantID uuid.UUID) (admissionservice.DashboardFilter, error) {
	var params ReportQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return admissionservice.DashboardFilter{}, err
	}

	filter := admissionservice.DashboardFilter{
		TenantID: tenantID,
	}

	// Parse session ID
	if params.SessionID != "" {
		sessionID, err := uuid.Parse(params.SessionID)
		if err != nil {
			return admissionservice.DashboardFilter{}, apperrors.BadRequest("Invalid session_id format")
		}
		filter.SessionID = &sessionID
	}

	// Parse branch ID
	if params.BranchID != "" {
		branchID, err := uuid.Parse(params.BranchID)
		if err != nil {
			return admissionservice.DashboardFilter{}, apperrors.BadRequest("Invalid branch_id format")
		}
		filter.BranchID = &branchID
	}

	// Parse academic year ID
	if params.AcademicYearID != "" {
		academicYearID, err := uuid.Parse(params.AcademicYearID)
		if err != nil {
			return admissionservice.DashboardFilter{}, apperrors.BadRequest("Invalid academic_year_id format")
		}
		filter.AcademicYearID = &academicYearID
	}

	// Parse dates
	startDate, err := parseDateWithError(params.StartDate)
	if err != nil {
		return admissionservice.DashboardFilter{}, apperrors.BadRequest("Invalid start_date format, expected YYYY-MM-DD")
	}
	filter.StartDate = startDate

	endDate, err := parseDateWithError(params.EndDate)
	if err != nil {
		return admissionservice.DashboardFilter{}, apperrors.BadRequest("Invalid end_date format, expected YYYY-MM-DD")
	}
	filter.EndDate = endDate

	return filter, nil
}
