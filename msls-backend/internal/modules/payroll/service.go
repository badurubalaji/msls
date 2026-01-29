// Package payroll provides payroll processing functionality.
package payroll

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// Service provides business logic for payroll operations.
type Service struct {
	repo *Repository
}

// NewService creates a new payroll service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Pay Run Service Methods
// ========================================

// CreatePayRun creates a new pay run.
func (s *Service) CreatePayRun(ctx context.Context, dto CreatePayRunDTO) (*models.PayRun, error) {
	// Validate pay period
	if dto.PayPeriodMonth < 1 || dto.PayPeriodMonth > 12 {
		return nil, ErrInvalidMonth
	}

	if dto.PayPeriodYear < 2000 || dto.PayPeriodYear > 2100 {
		return nil, ErrInvalidYear
	}

	// Check for existing pay run
	_, err := s.repo.GetPayRunByPeriod(ctx, dto.TenantID, dto.PayPeriodMonth, dto.PayPeriodYear, dto.BranchID)
	if err == nil {
		return nil, ErrDuplicatePayRun
	}
	if err != ErrPayRunNotFound {
		return nil, err
	}

	payRun := &models.PayRun{
		ID:             uuid.New(),
		TenantID:       dto.TenantID,
		PayPeriodMonth: dto.PayPeriodMonth,
		PayPeriodYear:  dto.PayPeriodYear,
		BranchID:       dto.BranchID,
		Status:         models.PayRunStatusDraft,
		TotalStaff:     0,
		TotalGross:     decimal.Zero,
		TotalDeductions: decimal.Zero,
		TotalNet:       decimal.Zero,
		Notes:          dto.Notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      dto.CreatedBy,
	}

	if err := s.repo.CreatePayRun(ctx, payRun); err != nil {
		return nil, err
	}

	return s.repo.GetPayRunByID(ctx, dto.TenantID, payRun.ID)
}

// GetPayRunByID retrieves a pay run by ID.
func (s *Service) GetPayRunByID(ctx context.Context, tenantID, id uuid.UUID) (*models.PayRun, error) {
	return s.repo.GetPayRunByID(ctx, tenantID, id)
}

// ListPayRuns retrieves pay runs with filters.
func (s *Service) ListPayRuns(ctx context.Context, filter PayRunFilter) ([]models.PayRun, int64, error) {
	return s.repo.ListPayRuns(ctx, filter)
}

// DeletePayRun deletes a draft pay run.
func (s *Service) DeletePayRun(ctx context.Context, tenantID, id uuid.UUID) error {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if payRun.Status != models.PayRunStatusDraft {
		return ErrPayRunNotDraft
	}

	return s.repo.DeletePayRun(ctx, tenantID, id)
}

// CalculatePayroll calculates payroll for a pay run.
func (s *Service) CalculatePayroll(ctx context.Context, tenantID, payRunID uuid.UUID) (*models.PayRun, error) {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	if payRun.Status != models.PayRunStatusDraft && payRun.Status != models.PayRunStatusCalculated {
		return nil, ErrPayRunNotDraft
	}

	// Update status to processing
	payRun.Status = models.PayRunStatusProcessing
	payRun.UpdatedAt = time.Now()
	if err := s.repo.UpdatePayRun(ctx, payRun); err != nil {
		return nil, err
	}

	// Delete existing payslips if recalculating
	if err := s.repo.DeletePayslipsByPayRun(ctx, payRunID); err != nil {
		return nil, err
	}

	// Get all active staff with salary
	staffList, err := s.repo.GetStaffWithActiveSalary(ctx, tenantID, payRun.BranchID)
	if err != nil {
		return nil, err
	}

	if len(staffList) == 0 {
		return nil, ErrNoStaffForPayroll
	}

	// Calculate working days for the month
	workingDays := s.getWorkingDaysInMonth(payRun.PayPeriodYear, payRun.PayPeriodMonth)

	var payslips []models.Payslip
	totalGross := decimal.Zero
	totalDeductions := decimal.Zero
	totalNet := decimal.Zero

	for _, staff := range staffList {
		// Get current salary for staff
		salary, err := s.repo.GetStaffCurrentSalary(ctx, tenantID, staff.ID)
		if err != nil {
			return nil, err
		}
		if salary == nil {
			continue // Skip staff without salary
		}

		// Calculate payslip for this staff
		payslip := s.calculatePayslip(payRun, &staff, salary, workingDays)
		payslips = append(payslips, payslip)

		totalGross = totalGross.Add(payslip.GrossSalary)
		totalDeductions = totalDeductions.Add(payslip.TotalDeductions)
		totalNet = totalNet.Add(payslip.NetSalary)
	}

	// Create all payslips
	if err := s.repo.CreatePayslips(ctx, payslips); err != nil {
		return nil, err
	}

	// Update pay run with totals
	now := time.Now()
	payRun.Status = models.PayRunStatusCalculated
	payRun.TotalStaff = len(payslips)
	payRun.TotalGross = totalGross
	payRun.TotalDeductions = totalDeductions
	payRun.TotalNet = totalNet
	payRun.CalculatedAt = &now
	payRun.UpdatedAt = now

	if err := s.repo.UpdatePayRun(ctx, payRun); err != nil {
		return nil, err
	}

	return s.repo.GetPayRunByID(ctx, tenantID, payRunID)
}

// ApprovePayRun approves a calculated pay run.
func (s *Service) ApprovePayRun(ctx context.Context, tenantID, payRunID uuid.UUID, approvedBy uuid.UUID) (*models.PayRun, error) {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	if payRun.Status != models.PayRunStatusCalculated {
		return nil, ErrPayRunNotCalculated
	}

	now := time.Now()
	payRun.Status = models.PayRunStatusApproved
	payRun.ApprovedAt = &now
	payRun.ApprovedBy = &approvedBy
	payRun.UpdatedAt = now

	if err := s.repo.UpdatePayRun(ctx, payRun); err != nil {
		return nil, err
	}

	return s.repo.GetPayRunByID(ctx, tenantID, payRunID)
}

// FinalizePayRun finalizes an approved pay run.
func (s *Service) FinalizePayRun(ctx context.Context, tenantID, payRunID uuid.UUID, finalizedBy uuid.UUID) (*models.PayRun, error) {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	if payRun.Status != models.PayRunStatusApproved {
		return nil, ErrPayRunNotApproved
	}

	now := time.Now()
	payRun.Status = models.PayRunStatusFinalized
	payRun.FinalizedAt = &now
	payRun.FinalizedBy = &finalizedBy
	payRun.UpdatedAt = now

	if err := s.repo.UpdatePayRun(ctx, payRun); err != nil {
		return nil, err
	}

	return s.repo.GetPayRunByID(ctx, tenantID, payRunID)
}

// GetPayRunSummary retrieves summary for a pay run.
func (s *Service) GetPayRunSummary(ctx context.Context, tenantID, payRunID uuid.UUID) (*PayRunSummaryResponse, error) {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	departmentSummary, err := s.repo.GetDepartmentSummary(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	return &PayRunSummaryResponse{
		PayRun:            ToPayRunResponse(payRun),
		DepartmentSummary: departmentSummary,
		TotalStaff:        payRun.TotalStaff,
		TotalCalculated:   payRun.TotalStaff,
		TotalExceptions:   0, // TODO: Implement exception tracking
	}, nil
}

// ========================================
// Payslip Service Methods
// ========================================

// GetPayslipByID retrieves a payslip by ID.
func (s *Service) GetPayslipByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Payslip, error) {
	return s.repo.GetPayslipByID(ctx, tenantID, id)
}

// ListPayslipsByPayRun retrieves all payslips for a pay run.
func (s *Service) ListPayslipsByPayRun(ctx context.Context, tenantID, payRunID uuid.UUID) ([]models.Payslip, int64, error) {
	return s.repo.ListPayslipsByPayRun(ctx, tenantID, payRunID)
}

// GetStaffPayslipHistory retrieves payslip history for a staff member.
func (s *Service) GetStaffPayslipHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.Payslip, int64, error) {
	return s.repo.GetStaffPayslipHistory(ctx, tenantID, staffID)
}

// AdjustPayslip adjusts a payslip before approval.
func (s *Service) AdjustPayslip(ctx context.Context, tenantID, payslipID uuid.UUID, dto AdjustPayslipDTO) (*models.Payslip, error) {
	payslip, err := s.repo.GetPayslipByID(ctx, tenantID, payslipID)
	if err != nil {
		return nil, err
	}

	// Check if payslip can be adjusted
	if payslip.Status != models.PayslipStatusCalculated && payslip.Status != models.PayslipStatusAdjusted {
		return nil, ErrPayslipNotAdjustable
	}

	// Create new components
	var components []models.PayslipComponent
	totalEarnings := decimal.Zero
	totalDeductions := decimal.Zero

	for _, adj := range dto.Components {
		// Get original component info from existing components
		var compName, compCode, compType string
		for _, c := range payslip.Components {
			if c.ComponentID == adj.ComponentID {
				compName = c.ComponentName
				compCode = c.ComponentCode
				compType = c.ComponentType
				break
			}
		}

		components = append(components, models.PayslipComponent{
			ID:            uuid.New(),
			PayslipID:     payslipID,
			ComponentID:   adj.ComponentID,
			ComponentName: compName,
			ComponentCode: compCode,
			ComponentType: compType,
			Amount:        adj.Amount,
			IsProrated:    false,
			CreatedAt:     time.Now(),
		})

		if compType == string(models.ComponentTypeEarning) {
			totalEarnings = totalEarnings.Add(adj.Amount)
		} else {
			totalDeductions = totalDeductions.Add(adj.Amount)
		}
	}

	// Update payslip totals
	payslip.TotalEarnings = totalEarnings
	payslip.TotalDeductions = totalDeductions.Add(payslip.LOPDeduction)
	payslip.GrossSalary = totalEarnings
	payslip.NetSalary = totalEarnings.Sub(payslip.TotalDeductions)
	payslip.Status = models.PayslipStatusAdjusted
	payslip.UpdatedAt = time.Now()

	// Replace components
	if err := s.repo.ReplacePayslipComponents(ctx, payslipID, components); err != nil {
		return nil, err
	}

	// Update payslip
	if err := s.repo.UpdatePayslip(ctx, payslip); err != nil {
		return nil, err
	}

	return s.repo.GetPayslipByID(ctx, tenantID, payslipID)
}

// ========================================
// Export Methods
// ========================================

// GetBankTransferExport generates bank transfer export data.
func (s *Service) GetBankTransferExport(ctx context.Context, tenantID, payRunID uuid.UUID) (*BankTransferExport, error) {
	payRun, err := s.repo.GetPayRunByID(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	payslips, _, err := s.repo.ListPayslipsByPayRun(ctx, tenantID, payRunID)
	if err != nil {
		return nil, err
	}

	var items []BankTransferItem
	totalAmount := decimal.Zero

	for _, ps := range payslips {
		if ps.Staff == nil {
			continue
		}

		item := BankTransferItem{
			StaffName:    ps.Staff.FullName(),
			EmployeeCode: ps.Staff.EmployeeID,
			BankName:     "", // Bank fields not yet implemented in Staff model
			BankAccount:  "",
			IFSC:         "",
			NetAmount:    ps.NetSalary.StringFixed(2),
		}
		items = append(items, item)
		totalAmount = totalAmount.Add(ps.NetSalary)
	}

	return &BankTransferExport{
		PayPeriod:    getMonthName(payRun.PayPeriodMonth) + " " + formatYear(payRun.PayPeriodYear),
		TotalRecords: len(items),
		TotalAmount:  totalAmount.StringFixed(2),
		Items:        items,
	}, nil
}

// ========================================
// Helper Methods
// ========================================

// calculatePayslip calculates a payslip for a staff member.
func (s *Service) calculatePayslip(payRun *models.PayRun, staff *models.Staff, salary *models.StaffSalary, workingDays int) models.Payslip {
	// For now, assume all staff are present for all working days (no attendance integration yet)
	// LOP calculation would be: (Daily rate) × (Absent days without approved leave)
	presentDays := workingDays
	leaveDays := 0
	absentDays := 0
	lopDays := 0

	// Calculate prorated amounts if needed (full month for now)
	totalEarnings := decimal.Zero
	totalDeductions := decimal.Zero

	var components []models.PayslipComponent

	for _, comp := range salary.Components {
		if comp.Component == nil {
			continue
		}

		amount := comp.Amount

		pc := models.PayslipComponent{
			ID:            uuid.New(),
			ComponentID:   comp.ComponentID,
			ComponentName: comp.Component.Name,
			ComponentCode: comp.Component.Code,
			ComponentType: string(comp.Component.ComponentType),
			Amount:        amount,
			IsProrated:    false,
			CreatedAt:     time.Now(),
		}

		if comp.Component.ComponentType == models.ComponentTypeEarning {
			totalEarnings = totalEarnings.Add(amount)
		} else {
			totalDeductions = totalDeductions.Add(amount)
		}

		components = append(components, pc)
	}

	// Calculate LOP deduction
	lopDeduction := decimal.Zero
	if lopDays > 0 && workingDays > 0 {
		dailyRate := totalEarnings.Div(decimal.NewFromInt(int64(workingDays)))
		lopDeduction = dailyRate.Mul(decimal.NewFromInt(int64(lopDays)))
		totalDeductions = totalDeductions.Add(lopDeduction)
	}

	netSalary := totalEarnings.Sub(totalDeductions)

	payslip := models.Payslip{
		ID:              uuid.New(),
		TenantID:        payRun.TenantID,
		PayRunID:        payRun.ID,
		StaffID:         staff.ID,
		StaffSalaryID:   &salary.ID,
		WorkingDays:     workingDays,
		PresentDays:     presentDays,
		LeaveDays:       leaveDays,
		AbsentDays:      absentDays,
		LOPDays:         lopDays,
		GrossSalary:     totalEarnings,
		TotalEarnings:   totalEarnings,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
		LOPDeduction:    lopDeduction,
		Status:          models.PayslipStatusCalculated,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Components:      components,
	}

	return payslip
}

// getWorkingDaysInMonth calculates working days in a month (excluding weekends).
func (s *Service) getWorkingDaysInMonth(year, month int) int {
	// Get the first day of the month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// Get the last day of the month
	lastDay := firstDay.AddDate(0, 1, -1)

	workingDays := 0
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		// Count weekdays (Monday to Saturday in Indian context)
		if d.Weekday() != time.Sunday {
			workingDays++
		}
	}

	return workingDays
}

