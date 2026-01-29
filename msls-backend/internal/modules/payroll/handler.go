// Package payroll provides payroll processing functionality.
package payroll

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles payroll-related HTTP requests.
type Handler struct {
	service      *Service
	pdfGenerator *PDFGenerator
}

// NewHandler creates a new payroll handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service:      service,
		pdfGenerator: NewPDFGenerator("MSLS School", ""),
	}
}

// ========================================
// Pay Run Handlers
// ========================================

// CreatePayRunRequest represents the request body for creating a pay run.
type CreatePayRunRequest struct {
	PayPeriodMonth int     `json:"payPeriodMonth" binding:"required,min=1,max=12"`
	PayPeriodYear  int     `json:"payPeriodYear" binding:"required,min=2000,max=2100"`
	BranchID       *string `json:"branchId" binding:"omitempty,uuid"`
	Notes          *string `json:"notes"`
}

// ListPayRuns returns all pay runs.
func (h *Handler) ListPayRuns(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := PayRunFilter{
		TenantID: tenantID,
	}

	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			filter.Year = &year
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err == nil {
			filter.Month = &month
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.PayRunStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err == nil {
			filter.BranchID = &branchID
		}
	}

	payRuns, total, err := h.service.ListPayRuns(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list pay runs"))
		return
	}

	response.OK(c, PayRunListResponse{
		PayRuns: ToPayRunResponses(payRuns),
		Total:   total,
	})
}

// GetPayRun returns a pay run by ID.
func (h *Handler) GetPayRun(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	payRun, err := h.service.GetPayRunByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get pay run"))
		return
	}

	response.OK(c, ToPayRunResponse(payRun))
}

// CreatePayRun creates a new pay run.
func (h *Handler) CreatePayRun(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreatePayRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreatePayRunDTO{
		TenantID:       tenantID,
		PayPeriodMonth: req.PayPeriodMonth,
		PayPeriodYear:  req.PayPeriodYear,
		Notes:          req.Notes,
		CreatedBy:      &userID,
	}

	if req.BranchID != nil {
		branchID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		dto.BranchID = &branchID
	}

	payRun, err := h.service.CreatePayRun(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrDuplicatePayRun) {
			apperrors.Abort(c, apperrors.Conflict("Pay run already exists for this period"))
			return
		}
		if errors.Is(err, ErrInvalidMonth) || errors.Is(err, ErrInvalidYear) {
			apperrors.Abort(c, apperrors.BadRequest(err.Error()))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to create pay run"))
		return
	}

	response.Created(c, ToPayRunResponse(payRun))
}

// CalculatePayroll calculates payroll for a pay run.
func (h *Handler) CalculatePayroll(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	payRun, err := h.service.CalculatePayroll(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		if errors.Is(err, ErrPayRunNotDraft) {
			apperrors.Abort(c, apperrors.BadRequest("Pay run is not in draft status"))
			return
		}
		if errors.Is(err, ErrNoStaffForPayroll) {
			apperrors.Abort(c, apperrors.BadRequest("No staff found with active salary for this period"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to calculate payroll"))
		return
	}

	response.OK(c, ToPayRunResponse(payRun))
}

// ApprovePayRun approves a pay run.
func (h *Handler) ApprovePayRun(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	payRun, err := h.service.ApprovePayRun(c.Request.Context(), tenantID, id, userID)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		if errors.Is(err, ErrPayRunNotCalculated) {
			apperrors.Abort(c, apperrors.BadRequest("Pay run has not been calculated"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to approve pay run"))
		return
	}

	response.OK(c, ToPayRunResponse(payRun))
}

// FinalizePayRun finalizes a pay run.
func (h *Handler) FinalizePayRun(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	payRun, err := h.service.FinalizePayRun(c.Request.Context(), tenantID, id, userID)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		if errors.Is(err, ErrPayRunNotApproved) {
			apperrors.Abort(c, apperrors.BadRequest("Pay run has not been approved"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to finalize pay run"))
		return
	}

	response.OK(c, ToPayRunResponse(payRun))
}

// DeletePayRun deletes a draft pay run.
func (h *Handler) DeletePayRun(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	if err := h.service.DeletePayRun(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		if errors.Is(err, ErrPayRunNotDraft) {
			apperrors.Abort(c, apperrors.BadRequest("Only draft pay runs can be deleted"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete pay run"))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPayRunSummary returns summary for a pay run.
func (h *Handler) GetPayRunSummary(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	summary, err := h.service.GetPayRunSummary(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get pay run summary"))
		return
	}

	response.OK(c, summary)
}

// ========================================
// Payslip Handlers
// ========================================

// ListPayslips returns all payslips for a pay run.
func (h *Handler) ListPayslips(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	payRunID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	payslips, total, err := h.service.ListPayslipsByPayRun(c.Request.Context(), tenantID, payRunID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list payslips"))
		return
	}

	response.OK(c, PayslipListResponse{
		Payslips: ToPayslipResponses(payslips),
		Total:    total,
	})
}

// GetPayslip returns a payslip by ID.
func (h *Handler) GetPayslip(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid payslip ID"))
		return
	}

	payslip, err := h.service.GetPayslipByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayslipNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Payslip not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get payslip"))
		return
	}

	response.OK(c, ToPayslipResponse(payslip, true))
}

// AdjustPayslip adjusts a payslip.
func (h *Handler) AdjustPayslip(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid payslip ID"))
		return
	}

	var dto AdjustPayslipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	payslip, err := h.service.AdjustPayslip(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		if errors.Is(err, ErrPayslipNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Payslip not found"))
			return
		}
		if errors.Is(err, ErrPayslipNotAdjustable) {
			apperrors.Abort(c, apperrors.BadRequest("Payslip cannot be adjusted in current state"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to adjust payslip"))
		return
	}

	response.OK(c, ToPayslipResponse(payslip, true))
}

// GetStaffPayslipHistory returns payslip history for a staff member.
func (h *Handler) GetStaffPayslipHistory(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	staffID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	payslips, total, err := h.service.GetStaffPayslipHistory(c.Request.Context(), tenantID, staffID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get payslip history"))
		return
	}

	response.OK(c, PayslipListResponse{
		Payslips: ToPayslipResponses(payslips),
		Total:    total,
	})
}

// ========================================
// Export Handlers
// ========================================

// DownloadPayslipPDF downloads a payslip as PDF.
func (h *Handler) DownloadPayslipPDF(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid payslip ID"))
		return
	}

	payslip, err := h.service.GetPayslipByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayslipNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Payslip not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get payslip"))
		return
	}

	// Get pay run for period info
	payRun, err := h.service.GetPayRunByID(c.Request.Context(), tenantID, payslip.PayRunID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get pay run"))
		return
	}

	// Generate PDF
	pdfBytes, err := h.pdfGenerator.GeneratePayslipPDF(payslip, payRun)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to generate PDF"))
		return
	}

	// Set response headers with month_year filename format
	filename := GetPayslipFilename(payRun)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportBankTransfer exports bank transfer data.
func (h *Handler) ExportBankTransfer(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid pay run ID"))
		return
	}

	export, err := h.service.GetBankTransferExport(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrPayRunNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Pay run not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to export bank transfer"))
		return
	}

	response.OK(c, export)
}

// RegisterRoutes registers payroll routes.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Pay Runs
	payRuns := rg.Group("/payroll/runs")
	{
		// Read operations - require payroll.view permission
		payRunsRead := payRuns.Group("")
		payRunsRead.Use(middleware.PermissionRequired("payroll.view"))
		{
			payRunsRead.GET("", h.ListPayRuns)
			payRunsRead.GET("/:id", h.GetPayRun)
			payRunsRead.GET("/:id/payslips", h.ListPayslips)
			payRunsRead.GET("/:id/summary", h.GetPayRunSummary)
		}

		// Create operations - require payroll.create permission
		payRunsCreate := payRuns.Group("")
		payRunsCreate.Use(middleware.PermissionRequired("payroll.create"))
		{
			payRunsCreate.POST("", h.CreatePayRun)
		}

		// Calculate operations - require payroll.calculate permission
		payRunsCalculate := payRuns.Group("")
		payRunsCalculate.Use(middleware.PermissionRequired("payroll.calculate"))
		{
			payRunsCalculate.POST("/:id/calculate", h.CalculatePayroll)
		}

		// Approve operations - require payroll.approve permission
		payRunsApprove := payRuns.Group("")
		payRunsApprove.Use(middleware.PermissionRequired("payroll.approve"))
		{
			payRunsApprove.POST("/:id/approve", h.ApprovePayRun)
		}

		// Finalize operations - require payroll.finalize permission
		payRunsFinalize := payRuns.Group("")
		payRunsFinalize.Use(middleware.PermissionRequired("payroll.finalize"))
		{
			payRunsFinalize.POST("/:id/finalize", h.FinalizePayRun)
		}

		// Delete operations - require payroll.delete permission
		payRunsDelete := payRuns.Group("")
		payRunsDelete.Use(middleware.PermissionRequired("payroll.delete"))
		{
			payRunsDelete.DELETE("/:id", h.DeletePayRun)
		}

		// Export operations - require payroll.export permission
		payRunsExport := payRuns.Group("")
		payRunsExport.Use(middleware.PermissionRequired("payroll.export"))
		{
			payRunsExport.GET("/:id/export", h.ExportBankTransfer)
		}
	}

	// Payslips
	payslips := rg.Group("/payroll/payslips")
	{
		// Read operations - require payroll.view permission
		payslipsRead := payslips.Group("")
		payslipsRead.Use(middleware.PermissionRequired("payroll.view"))
		{
			payslipsRead.GET("/:id", h.GetPayslip)
			payslipsRead.GET("/:id/pdf", h.DownloadPayslipPDF)
		}

		// Adjust operations - require payroll.adjust permission
		payslipsAdjust := payslips.Group("")
		payslipsAdjust.Use(middleware.PermissionRequired("payroll.adjust"))
		{
			payslipsAdjust.PUT("/:id", h.AdjustPayslip)
		}
	}
}

// RegisterStaffPayslipRoutes registers staff payslip routes.
func (h *Handler) RegisterStaffPayslipRoutes(staffGroup *gin.RouterGroup) {
	staffGroup.GET("/:id/payslips", h.GetStaffPayslipHistory)
}
