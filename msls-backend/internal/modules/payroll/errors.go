// Package payroll provides payroll processing functionality.
package payroll

import "errors"

// Pay run errors.
var (
	ErrPayRunNotFound      = errors.New("pay run not found")
	ErrDuplicatePayRun     = errors.New("pay run already exists for this period")
	ErrPayRunNotDraft      = errors.New("pay run is not in draft status")
	ErrPayRunNotCalculated = errors.New("pay run has not been calculated")
	ErrPayRunNotApproved   = errors.New("pay run has not been approved")
	ErrPayRunFinalized     = errors.New("pay run is already finalized")
	ErrNoStaffForPayroll   = errors.New("no staff found with active salary for this period")
)

// Payslip errors.
var (
	ErrPayslipNotFound      = errors.New("payslip not found")
	ErrPayslipNotAdjustable = errors.New("payslip cannot be adjusted in current state")
	ErrInvalidPayslipStatus = errors.New("invalid payslip status")
)

// Validation errors.
var (
	ErrInvalidPayPeriod = errors.New("invalid pay period")
	ErrInvalidMonth     = errors.New("month must be between 1 and 12")
	ErrInvalidYear      = errors.New("invalid year")
)
