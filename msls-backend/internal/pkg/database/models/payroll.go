// Package models contains database model definitions.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PayRunStatus represents the status of a pay run.
type PayRunStatus string

const (
	PayRunStatusDraft      PayRunStatus = "draft"
	PayRunStatusProcessing PayRunStatus = "processing"
	PayRunStatusCalculated PayRunStatus = "calculated"
	PayRunStatusApproved   PayRunStatus = "approved"
	PayRunStatusFinalized  PayRunStatus = "finalized"
	PayRunStatusReversed   PayRunStatus = "reversed"
)

// IsValid checks if the pay run status is valid.
func (s PayRunStatus) IsValid() bool {
	switch s {
	case PayRunStatusDraft, PayRunStatusProcessing, PayRunStatusCalculated,
		PayRunStatusApproved, PayRunStatusFinalized, PayRunStatusReversed:
		return true
	}
	return false
}

// PayslipStatus represents the status of a payslip.
type PayslipStatus string

const (
	PayslipStatusCalculated PayslipStatus = "calculated"
	PayslipStatusAdjusted   PayslipStatus = "adjusted"
	PayslipStatusApproved   PayslipStatus = "approved"
	PayslipStatusPaid       PayslipStatus = "paid"
)

// IsValid checks if the payslip status is valid.
func (s PayslipStatus) IsValid() bool {
	switch s {
	case PayslipStatusCalculated, PayslipStatusAdjusted, PayslipStatusApproved, PayslipStatusPaid:
		return true
	}
	return false
}

// PayRun represents a monthly payroll batch.
type PayRun struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID       uuid.UUID        `gorm:"type:uuid;not null;index"`
	PayPeriodMonth int              `gorm:"not null"`
	PayPeriodYear  int              `gorm:"not null"`
	BranchID       *uuid.UUID       `gorm:"type:uuid;index"`
	Status         PayRunStatus     `gorm:"type:varchar(20);not null;default:'draft'"`
	TotalStaff     int              `gorm:"not null;default:0"`
	TotalGross     decimal.Decimal  `gorm:"type:decimal(14,2);not null;default:0"`
	TotalDeductions decimal.Decimal `gorm:"type:decimal(14,2);not null;default:0"`
	TotalNet       decimal.Decimal  `gorm:"type:decimal(14,2);not null;default:0"`
	CalculatedAt   *time.Time       `gorm:"type:timestamptz"`
	ApprovedAt     *time.Time       `gorm:"type:timestamptz"`
	ApprovedBy     *uuid.UUID       `gorm:"type:uuid"`
	FinalizedAt    *time.Time       `gorm:"type:timestamptz"`
	FinalizedBy    *uuid.UUID       `gorm:"type:uuid"`
	Notes          *string          `gorm:"type:text"`
	CreatedAt      time.Time        `gorm:"not null;default:now()"`
	UpdatedAt      time.Time        `gorm:"not null;default:now()"`
	CreatedBy      *uuid.UUID       `gorm:"type:uuid"`

	// Relations
	Branch       *Branch    `gorm:"foreignKey:BranchID"`
	Approver     *User      `gorm:"foreignKey:ApprovedBy"`
	Finalizer    *User      `gorm:"foreignKey:FinalizedBy"`
	Creator      *User      `gorm:"foreignKey:CreatedBy"`
	Payslips     []Payslip  `gorm:"foreignKey:PayRunID"`
}

// TableName returns the table name for PayRun.
func (PayRun) TableName() string {
	return "pay_runs"
}

// Payslip represents an individual staff pay record for a pay run.
type Payslip struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	PayRunID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	StaffID         uuid.UUID        `gorm:"type:uuid;not null;index"`
	StaffSalaryID   *uuid.UUID       `gorm:"type:uuid"`
	WorkingDays     int              `gorm:"not null;default:0"`
	PresentDays     int              `gorm:"not null;default:0"`
	LeaveDays       int              `gorm:"not null;default:0"`
	AbsentDays      int              `gorm:"not null;default:0"`
	LOPDays         int              `gorm:"not null;default:0"`
	GrossSalary     decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	TotalEarnings   decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	TotalDeductions decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	NetSalary       decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	LOPDeduction    decimal.Decimal  `gorm:"type:decimal(12,2);not null;default:0"`
	Status          PayslipStatus    `gorm:"type:varchar(20);not null;default:'calculated'"`
	PaymentDate     *time.Time       `gorm:"type:date"`
	PaymentReference *string         `gorm:"type:varchar(100)"`
	CreatedAt       time.Time        `gorm:"not null;default:now()"`
	UpdatedAt       time.Time        `gorm:"not null;default:now()"`

	// Relations
	PayRun      *PayRun             `gorm:"foreignKey:PayRunID"`
	Staff       *Staff              `gorm:"foreignKey:StaffID"`
	StaffSalary *StaffSalary        `gorm:"foreignKey:StaffSalaryID"`
	Components  []PayslipComponent  `gorm:"foreignKey:PayslipID"`
}

// TableName returns the table name for Payslip.
func (Payslip) TableName() string {
	return "payslips"
}

// PayslipComponent represents a breakdown of a payslip.
type PayslipComponent struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	PayslipID     uuid.UUID       `gorm:"type:uuid;not null;index"`
	ComponentID   uuid.UUID       `gorm:"type:uuid;not null"`
	ComponentName string          `gorm:"type:varchar(100);not null"`
	ComponentCode string          `gorm:"type:varchar(20);not null"`
	ComponentType string          `gorm:"type:varchar(20);not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(12,2);not null"`
	IsProrated    bool            `gorm:"not null;default:false"`
	CreatedAt     time.Time       `gorm:"not null;default:now()"`

	// Relations
	Component *SalaryComponent `gorm:"foreignKey:ComponentID"`
}

// TableName returns the table name for PayslipComponent.
func (PayslipComponent) TableName() string {
	return "payslip_components"
}
