// Package payroll provides payroll processing functionality.
package payroll

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for payroll management.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new payroll repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Pay Run Repository Methods
// ========================================

// CreatePayRun creates a new pay run.
func (r *Repository) CreatePayRun(ctx context.Context, payRun *models.PayRun) error {
	if err := r.db.WithContext(ctx).Create(payRun).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "pay_runs_tenant_id_pay_period") {
			return ErrDuplicatePayRun
		}
		return fmt.Errorf("create pay run: %w", err)
	}
	return nil
}

// GetPayRunByID retrieves a pay run by ID.
func (r *Repository) GetPayRunByID(ctx context.Context, tenantID, id uuid.UUID) (*models.PayRun, error) {
	var payRun models.PayRun
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Approver").
		Preload("Finalizer").
		Preload("Creator").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&payRun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPayRunNotFound
		}
		return nil, fmt.Errorf("get pay run by id: %w", err)
	}
	return &payRun, nil
}

// GetPayRunByPeriod retrieves a pay run by period.
func (r *Repository) GetPayRunByPeriod(ctx context.Context, tenantID uuid.UUID, month, year int, branchID *uuid.UUID) (*models.PayRun, error) {
	var payRun models.PayRun
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND pay_period_month = ? AND pay_period_year = ?", tenantID, month, year)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}

	err := query.First(&payRun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPayRunNotFound
		}
		return nil, fmt.Errorf("get pay run by period: %w", err)
	}
	return &payRun, nil
}

// UpdatePayRun updates a pay run.
func (r *Repository) UpdatePayRun(ctx context.Context, payRun *models.PayRun) error {
	result := r.db.WithContext(ctx).
		Model(payRun).
		Updates(map[string]interface{}{
			"status":           payRun.Status,
			"total_staff":      payRun.TotalStaff,
			"total_gross":      payRun.TotalGross,
			"total_deductions": payRun.TotalDeductions,
			"total_net":        payRun.TotalNet,
			"calculated_at":    payRun.CalculatedAt,
			"approved_at":      payRun.ApprovedAt,
			"approved_by":      payRun.ApprovedBy,
			"finalized_at":     payRun.FinalizedAt,
			"finalized_by":     payRun.FinalizedBy,
			"notes":            payRun.Notes,
			"updated_at":       payRun.UpdatedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("update pay run: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrPayRunNotFound
	}

	return nil
}

// DeletePayRun deletes a pay run.
func (r *Repository) DeletePayRun(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND status = ?", tenantID, id, models.PayRunStatusDraft).
		Delete(&models.PayRun{})

	if result.Error != nil {
		return fmt.Errorf("delete pay run: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrPayRunNotFound
	}

	return nil
}

// ListPayRuns retrieves pay runs with filters.
func (r *Repository) ListPayRuns(ctx context.Context, filter PayRunFilter) ([]models.PayRun, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.PayRun{}).
		Preload("Branch").
		Preload("Approver").
		Preload("Finalizer").
		Where("tenant_id = ?", filter.TenantID)

	if filter.Year != nil {
		query = query.Where("pay_period_year = ?", *filter.Year)
	}

	if filter.Month != nil {
		query = query.Where("pay_period_month = ?", *filter.Month)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count pay runs: %w", err)
	}

	var payRuns []models.PayRun
	if err := query.Order("pay_period_year DESC, pay_period_month DESC, created_at DESC").Find(&payRuns).Error; err != nil {
		return nil, 0, fmt.Errorf("list pay runs: %w", err)
	}

	return payRuns, total, nil
}

// ========================================
// Payslip Repository Methods
// ========================================

// CreatePayslips creates multiple payslips in a batch.
func (r *Repository) CreatePayslips(ctx context.Context, payslips []models.Payslip) error {
	if len(payslips) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Create(&payslips).Error; err != nil {
		return fmt.Errorf("create payslips: %w", err)
	}
	return nil
}

// GetPayslipByID retrieves a payslip by ID.
func (r *Repository) GetPayslipByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Payslip, error) {
	var payslip models.Payslip
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.Department").
		Preload("Staff.Designation").
		Preload("Components").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&payslip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPayslipNotFound
		}
		return nil, fmt.Errorf("get payslip by id: %w", err)
	}
	return &payslip, nil
}

// ListPayslipsByPayRun retrieves all payslips for a pay run.
func (r *Repository) ListPayslipsByPayRun(ctx context.Context, tenantID, payRunID uuid.UUID) ([]models.Payslip, int64, error) {
	var payslips []models.Payslip
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.Department").
		Preload("Staff.Designation").
		Where("tenant_id = ? AND pay_run_id = ?", tenantID, payRunID).
		Order("created_at ASC").
		Find(&payslips).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list payslips: %w", err)
	}
	return payslips, int64(len(payslips)), nil
}

// GetStaffPayslipHistory retrieves payslip history for a staff member.
func (r *Repository) GetStaffPayslipHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.Payslip, int64, error) {
	var payslips []models.Payslip
	err := r.db.WithContext(ctx).
		Preload("PayRun").
		Preload("Components").
		Where("tenant_id = ? AND staff_id = ?", tenantID, staffID).
		Order("created_at DESC").
		Find(&payslips).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get staff payslip history: %w", err)
	}
	return payslips, int64(len(payslips)), nil
}

// UpdatePayslip updates a payslip.
func (r *Repository) UpdatePayslip(ctx context.Context, payslip *models.Payslip) error {
	result := r.db.WithContext(ctx).
		Model(payslip).
		Updates(map[string]interface{}{
			"gross_salary":      payslip.GrossSalary,
			"total_earnings":    payslip.TotalEarnings,
			"total_deductions":  payslip.TotalDeductions,
			"net_salary":        payslip.NetSalary,
			"lop_deduction":     payslip.LOPDeduction,
			"status":            payslip.Status,
			"payment_date":      payslip.PaymentDate,
			"payment_reference": payslip.PaymentReference,
			"updated_at":        payslip.UpdatedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("update payslip: %w", result.Error)
	}

	return nil
}

// DeletePayslipsByPayRun deletes all payslips for a pay run.
func (r *Repository) DeletePayslipsByPayRun(ctx context.Context, payRunID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("pay_run_id = ?", payRunID).
		Delete(&models.Payslip{}).Error; err != nil {
		return fmt.Errorf("delete payslips: %w", err)
	}
	return nil
}

// ReplacePayslipComponents replaces all components in a payslip.
func (r *Repository) ReplacePayslipComponents(ctx context.Context, payslipID uuid.UUID, components []models.PayslipComponent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing components
		if err := tx.Where("payslip_id = ?", payslipID).Delete(&models.PayslipComponent{}).Error; err != nil {
			return fmt.Errorf("delete existing components: %w", err)
		}

		// Insert new components
		if len(components) > 0 {
			if err := tx.Create(&components).Error; err != nil {
				return fmt.Errorf("create payslip components: %w", err)
			}
		}

		return nil
	})
}

// ========================================
// Staff Query Methods (for payroll calculation)
// ========================================

// GetStaffWithActiveSalary retrieves staff with active salary for payroll.
func (r *Repository) GetStaffWithActiveSalary(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) ([]models.Staff, error) {
	query := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Designation").
		Where("tenant_id = ? AND status = ?", tenantID, models.StaffStatusActive)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	var staff []models.Staff
	if err := query.Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("get staff with active salary: %w", err)
	}

	return staff, nil
}

// GetStaffCurrentSalary retrieves current salary for a staff member.
func (r *Repository) GetStaffCurrentSalary(ctx context.Context, tenantID, staffID uuid.UUID) (*models.StaffSalary, error) {
	var salary models.StaffSalary
	err := r.db.WithContext(ctx).
		Preload("Components").
		Preload("Components.Component").
		Where("tenant_id = ? AND staff_id = ? AND is_current = true", tenantID, staffID).
		First(&salary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No salary assigned
		}
		return nil, fmt.Errorf("get staff current salary: %w", err)
	}
	return &salary, nil
}

// GetDepartmentSummary retrieves department-wise summary for a pay run.
func (r *Repository) GetDepartmentSummary(ctx context.Context, tenantID, payRunID uuid.UUID) ([]DepartmentSummaryItem, error) {
	var results []struct {
		DepartmentID   uuid.UUID `gorm:"column:department_id"`
		DepartmentName string    `gorm:"column:department_name"`
		StaffCount     int       `gorm:"column:staff_count"`
		TotalGross     string    `gorm:"column:total_gross"`
		TotalDeductions string   `gorm:"column:total_deductions"`
		TotalNet       string    `gorm:"column:total_net"`
	}

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				s.department_id,
				d.name as department_name,
				COUNT(p.id) as staff_count,
				COALESCE(SUM(p.gross_salary), 0)::text as total_gross,
				COALESCE(SUM(p.total_deductions), 0)::text as total_deductions,
				COALESCE(SUM(p.net_salary), 0)::text as total_net
			FROM payslips p
			JOIN staff s ON p.staff_id = s.id
			LEFT JOIN departments d ON s.department_id = d.id
			WHERE p.tenant_id = ? AND p.pay_run_id = ?
			GROUP BY s.department_id, d.name
			ORDER BY d.name
		`, tenantID, payRunID).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("get department summary: %w", err)
	}

	items := make([]DepartmentSummaryItem, len(results))
	for i, r := range results {
		items[i] = DepartmentSummaryItem{
			DepartmentID:    r.DepartmentID.String(),
			DepartmentName:  r.DepartmentName,
			StaffCount:      r.StaffCount,
			TotalGross:      r.TotalGross,
			TotalDeductions: r.TotalDeductions,
			TotalNet:        r.TotalNet,
		}
	}

	return items, nil
}

// DB returns the underlying database connection.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
