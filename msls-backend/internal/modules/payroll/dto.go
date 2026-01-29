// Package payroll provides payroll processing functionality.
package payroll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Pay Run DTOs
// ========================================

// CreatePayRunDTO represents a request to create a pay run.
type CreatePayRunDTO struct {
	TenantID       uuid.UUID
	PayPeriodMonth int
	PayPeriodYear  int
	BranchID       *uuid.UUID
	Notes          *string
	CreatedBy      *uuid.UUID
}

// PayRunFilter contains filter options for listing pay runs.
type PayRunFilter struct {
	TenantID uuid.UUID
	Year     *int
	Month    *int
	Status   *models.PayRunStatus
	BranchID *uuid.UUID
}

// PayRunResponse represents a pay run in API responses.
type PayRunResponse struct {
	ID              string `json:"id"`
	PayPeriodMonth  int    `json:"payPeriodMonth"`
	PayPeriodYear   int    `json:"payPeriodYear"`
	PayPeriodLabel  string `json:"payPeriodLabel"`
	BranchID        string `json:"branchId,omitempty"`
	BranchName      string `json:"branchName,omitempty"`
	Status          string `json:"status"`
	TotalStaff      int    `json:"totalStaff"`
	TotalGross      string `json:"totalGross"`
	TotalDeductions string `json:"totalDeductions"`
	TotalNet        string `json:"totalNet"`
	CalculatedAt    string `json:"calculatedAt,omitempty"`
	ApprovedAt      string `json:"approvedAt,omitempty"`
	ApprovedByName  string `json:"approvedByName,omitempty"`
	FinalizedAt     string `json:"finalizedAt,omitempty"`
	FinalizedByName string `json:"finalizedByName,omitempty"`
	Notes           string `json:"notes,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// PayRunListResponse represents a list of pay runs.
type PayRunListResponse struct {
	PayRuns []PayRunResponse `json:"payRuns"`
	Total   int64            `json:"total"`
}

// PayRunSummaryResponse represents summary for a pay run.
type PayRunSummaryResponse struct {
	PayRun              PayRunResponse            `json:"payRun"`
	DepartmentSummary   []DepartmentSummaryItem   `json:"departmentSummary"`
	TotalStaff          int                       `json:"totalStaff"`
	TotalCalculated     int                       `json:"totalCalculated"`
	TotalExceptions     int                       `json:"totalExceptions"`
}

// DepartmentSummaryItem represents department-wise payroll summary.
type DepartmentSummaryItem struct {
	DepartmentID   string `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	StaffCount     int    `json:"staffCount"`
	TotalGross     string `json:"totalGross"`
	TotalDeductions string `json:"totalDeductions"`
	TotalNet       string `json:"totalNet"`
}

// ToPayRunResponse converts a PayRun model to a response.
func ToPayRunResponse(pr *models.PayRun) PayRunResponse {
	resp := PayRunResponse{
		ID:              pr.ID.String(),
		PayPeriodMonth:  pr.PayPeriodMonth,
		PayPeriodYear:   pr.PayPeriodYear,
		PayPeriodLabel:  getMonthName(pr.PayPeriodMonth) + " " + string(rune(pr.PayPeriodYear)),
		Status:          string(pr.Status),
		TotalStaff:      pr.TotalStaff,
		TotalGross:      pr.TotalGross.StringFixed(2),
		TotalDeductions: pr.TotalDeductions.StringFixed(2),
		TotalNet:        pr.TotalNet.StringFixed(2),
		CreatedAt:       pr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       pr.UpdatedAt.Format(time.RFC3339),
	}

	// Format pay period label properly
	resp.PayPeriodLabel = getMonthName(pr.PayPeriodMonth) + " " + formatYear(pr.PayPeriodYear)

	if pr.BranchID != nil {
		resp.BranchID = pr.BranchID.String()
	}

	if pr.Branch != nil {
		resp.BranchName = pr.Branch.Name
	}

	if pr.CalculatedAt != nil {
		resp.CalculatedAt = pr.CalculatedAt.Format(time.RFC3339)
	}

	if pr.ApprovedAt != nil {
		resp.ApprovedAt = pr.ApprovedAt.Format(time.RFC3339)
	}

	if pr.Approver != nil {
		resp.ApprovedByName = pr.Approver.FullName()
	}

	if pr.FinalizedAt != nil {
		resp.FinalizedAt = pr.FinalizedAt.Format(time.RFC3339)
	}

	if pr.Finalizer != nil {
		resp.FinalizedByName = pr.Finalizer.FullName()
	}

	if pr.Notes != nil {
		resp.Notes = *pr.Notes
	}

	return resp
}

// ToPayRunResponses converts a slice of PayRun models to responses.
func ToPayRunResponses(payRuns []models.PayRun) []PayRunResponse {
	responses := make([]PayRunResponse, len(payRuns))
	for i, pr := range payRuns {
		responses[i] = ToPayRunResponse(&pr)
	}
	return responses
}

// ========================================
// Payslip DTOs
// ========================================

// AdjustPayslipDTO represents a request to adjust a payslip.
type AdjustPayslipDTO struct {
	Components []PayslipComponentAdjustment `json:"components"`
}

// PayslipComponentAdjustment represents an adjustment to a payslip component.
type PayslipComponentAdjustment struct {
	ComponentID uuid.UUID       `json:"componentId"`
	Amount      decimal.Decimal `json:"amount"`
}

// PayslipResponse represents a payslip in API responses.
type PayslipResponse struct {
	ID               string                      `json:"id"`
	PayRunID         string                      `json:"payRunId"`
	StaffID          string                      `json:"staffId"`
	StaffName        string                      `json:"staffName"`
	EmployeeID       string                      `json:"employeeId,omitempty"`
	DepartmentName   string                      `json:"departmentName,omitempty"`
	DesignationName  string                      `json:"designationName,omitempty"`
	WorkingDays      int                         `json:"workingDays"`
	PresentDays      int                         `json:"presentDays"`
	LeaveDays        int                         `json:"leaveDays"`
	AbsentDays       int                         `json:"absentDays"`
	LOPDays          int                         `json:"lopDays"`
	GrossSalary      string                      `json:"grossSalary"`
	TotalEarnings    string                      `json:"totalEarnings"`
	TotalDeductions  string                      `json:"totalDeductions"`
	NetSalary        string                      `json:"netSalary"`
	LOPDeduction     string                      `json:"lopDeduction"`
	Status           string                      `json:"status"`
	PaymentDate      string                      `json:"paymentDate,omitempty"`
	PaymentReference string                      `json:"paymentReference,omitempty"`
	Components       []PayslipComponentResponse  `json:"components,omitempty"`
	CreatedAt        string                      `json:"createdAt"`
	UpdatedAt        string                      `json:"updatedAt"`
}

// PayslipComponentResponse represents a payslip component in API responses.
type PayslipComponentResponse struct {
	ID            string `json:"id"`
	ComponentID   string `json:"componentId"`
	ComponentName string `json:"componentName"`
	ComponentCode string `json:"componentCode"`
	ComponentType string `json:"componentType"`
	Amount        string `json:"amount"`
	IsProrated    bool   `json:"isProrated"`
}

// PayslipListResponse represents a list of payslips.
type PayslipListResponse struct {
	Payslips []PayslipResponse `json:"payslips"`
	Total    int64             `json:"total"`
}

// ToPayslipResponse converts a Payslip model to a response.
func ToPayslipResponse(ps *models.Payslip, includeComponents bool) PayslipResponse {
	resp := PayslipResponse{
		ID:              ps.ID.String(),
		PayRunID:        ps.PayRunID.String(),
		StaffID:         ps.StaffID.String(),
		WorkingDays:     ps.WorkingDays,
		PresentDays:     ps.PresentDays,
		LeaveDays:       ps.LeaveDays,
		AbsentDays:      ps.AbsentDays,
		LOPDays:         ps.LOPDays,
		GrossSalary:     ps.GrossSalary.StringFixed(2),
		TotalEarnings:   ps.TotalEarnings.StringFixed(2),
		TotalDeductions: ps.TotalDeductions.StringFixed(2),
		NetSalary:       ps.NetSalary.StringFixed(2),
		LOPDeduction:    ps.LOPDeduction.StringFixed(2),
		Status:          string(ps.Status),
		CreatedAt:       ps.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       ps.UpdatedAt.Format(time.RFC3339),
	}

	if ps.Staff != nil {
		resp.StaffName = ps.Staff.FullName()
		resp.EmployeeID = ps.Staff.EmployeeID
		if ps.Staff.Department != nil {
			resp.DepartmentName = ps.Staff.Department.Name
		}
		if ps.Staff.Designation != nil {
			resp.DesignationName = ps.Staff.Designation.Name
		}
	}

	if ps.PaymentDate != nil {
		resp.PaymentDate = ps.PaymentDate.Format("2006-01-02")
	}

	if ps.PaymentReference != nil {
		resp.PaymentReference = *ps.PaymentReference
	}

	if includeComponents && len(ps.Components) > 0 {
		resp.Components = make([]PayslipComponentResponse, len(ps.Components))
		for i, comp := range ps.Components {
			resp.Components[i] = ToPayslipComponentResponse(&comp)
		}
	}

	return resp
}

// ToPayslipComponentResponse converts a PayslipComponent to a response.
func ToPayslipComponentResponse(pc *models.PayslipComponent) PayslipComponentResponse {
	return PayslipComponentResponse{
		ID:            pc.ID.String(),
		ComponentID:   pc.ComponentID.String(),
		ComponentName: pc.ComponentName,
		ComponentCode: pc.ComponentCode,
		ComponentType: pc.ComponentType,
		Amount:        pc.Amount.StringFixed(2),
		IsProrated:    pc.IsProrated,
	}
}

// ToPayslipResponses converts a slice of Payslip models to responses.
func ToPayslipResponses(payslips []models.Payslip) []PayslipResponse {
	responses := make([]PayslipResponse, len(payslips))
	for i, ps := range payslips {
		responses[i] = ToPayslipResponse(&ps, false)
	}
	return responses
}

// ========================================
// Export DTOs
// ========================================

// BankTransferItem represents a single bank transfer entry for export.
type BankTransferItem struct {
	StaffName     string `json:"staffName"`
	EmployeeCode  string `json:"employeeCode"`
	BankName      string `json:"bankName"`
	BankAccount   string `json:"bankAccount"`
	IFSC          string `json:"ifsc"`
	NetAmount     string `json:"netAmount"`
}

// BankTransferExport represents the bank transfer export data.
type BankTransferExport struct {
	PayPeriod    string             `json:"payPeriod"`
	TotalRecords int                `json:"totalRecords"`
	TotalAmount  string             `json:"totalAmount"`
	Items        []BankTransferItem `json:"items"`
}

// ========================================
// Helper functions
// ========================================

func getMonthName(month int) string {
	months := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if month >= 1 && month <= 12 {
		return months[month]
	}
	return ""
}

func formatYear(year int) string {
	return fmt.Sprintf("%d", year)
}
