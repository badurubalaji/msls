// Package models contains database model definitions.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ComponentType represents the type of salary component.
type ComponentType string

const (
	ComponentTypeEarning   ComponentType = "earning"
	ComponentTypeDeduction ComponentType = "deduction"
)

// IsValid checks if the component type is valid.
func (c ComponentType) IsValid() bool {
	switch c {
	case ComponentTypeEarning, ComponentTypeDeduction:
		return true
	}
	return false
}

// CalculationType represents how a component is calculated.
type CalculationType string

const (
	CalculationTypeFixed      CalculationType = "fixed"
	CalculationTypePercentage CalculationType = "percentage"
)

// IsValid checks if the calculation type is valid.
func (c CalculationType) IsValid() bool {
	switch c {
	case CalculationTypeFixed, CalculationTypePercentage:
		return true
	}
	return false
}

// SalaryComponent represents a salary component (earning or deduction).
type SalaryComponent struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID        uuid.UUID       `gorm:"type:uuid;not null;index"`
	Name            string          `gorm:"type:varchar(100);not null"`
	Code            string          `gorm:"type:varchar(20);not null"`
	Description     *string         `gorm:"type:text"`
	ComponentType   ComponentType   `gorm:"type:varchar(20);not null"`
	CalculationType CalculationType `gorm:"type:varchar(20);not null"`
	PercentageOfID  *uuid.UUID      `gorm:"type:uuid"`
	IsTaxable       bool            `gorm:"not null;default:true"`
	IsActive        bool            `gorm:"not null;default:true"`
	DisplayOrder    int             `gorm:"not null;default:0"`
	CreatedAt       time.Time       `gorm:"not null;default:now()"`
	UpdatedAt       time.Time       `gorm:"not null;default:now()"`

	// Relations
	PercentageOf *SalaryComponent `gorm:"foreignKey:PercentageOfID"`
}

// TableName returns the table name for SalaryComponent.
func (SalaryComponent) TableName() string {
	return "salary_components"
}

// SalaryStructure represents a salary structure template.
type SalaryStructure struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	Name          string     `gorm:"type:varchar(100);not null"`
	Code          string     `gorm:"type:varchar(20);not null"`
	Description   *string    `gorm:"type:text"`
	DesignationID *uuid.UUID `gorm:"type:uuid;index"`
	IsActive      bool       `gorm:"not null;default:true"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()"`

	// Relations
	Designation *Designation               `gorm:"foreignKey:DesignationID"`
	Components  []SalaryStructureComponent `gorm:"foreignKey:StructureID"`
}

// TableName returns the table name for SalaryStructure.
func (SalaryStructure) TableName() string {
	return "salary_structures"
}

// SalaryStructureComponent represents a component in a salary structure.
type SalaryStructureComponent struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	StructureID uuid.UUID        `gorm:"type:uuid;not null;index"`
	ComponentID uuid.UUID        `gorm:"type:uuid;not null"`
	Amount      *decimal.Decimal `gorm:"type:decimal(12,2)"`
	Percentage  *decimal.Decimal `gorm:"type:decimal(5,2)"`
	CreatedAt   time.Time        `gorm:"not null;default:now()"`

	// Relations
	Component *SalaryComponent `gorm:"foreignKey:ComponentID"`
}

// TableName returns the table name for SalaryStructureComponent.
func (SalaryStructureComponent) TableName() string {
	return "salary_structure_components"
}

// StaffSalary represents a staff member's salary assignment.
type StaffSalary struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID       uuid.UUID        `gorm:"type:uuid;not null;index"`
	StaffID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	StructureID    *uuid.UUID       `gorm:"type:uuid;index"`
	EffectiveFrom  time.Time        `gorm:"type:date;not null"`
	EffectiveTo    *time.Time       `gorm:"type:date"`
	GrossSalary    decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	NetSalary      decimal.Decimal  `gorm:"type:decimal(12,2);not null"`
	CTC            *decimal.Decimal `gorm:"type:decimal(12,2)"`
	RevisionReason *string          `gorm:"type:text"`
	IsCurrent      bool             `gorm:"not null;default:true"`
	CreatedAt      time.Time        `gorm:"not null;default:now()"`
	UpdatedAt      time.Time        `gorm:"not null;default:now()"`
	CreatedBy      *uuid.UUID       `gorm:"type:uuid"`

	// Relations
	Staff      *Staff                   `gorm:"foreignKey:StaffID"`
	Structure  *SalaryStructure         `gorm:"foreignKey:StructureID"`
	Components []StaffSalaryComponent   `gorm:"foreignKey:StaffSalaryID"`
	Creator    *User                    `gorm:"foreignKey:CreatedBy"`
}

// TableName returns the table name for StaffSalary.
func (StaffSalary) TableName() string {
	return "staff_salaries"
}

// StaffSalaryComponent represents a component value in a staff's salary.
type StaffSalaryComponent struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	StaffSalaryID uuid.UUID       `gorm:"type:uuid;not null;index"`
	ComponentID   uuid.UUID       `gorm:"type:uuid;not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(12,2);not null"`
	IsOverridden  bool            `gorm:"not null;default:false"`
	CreatedAt     time.Time       `gorm:"not null;default:now()"`

	// Relations
	Component *SalaryComponent `gorm:"foreignKey:ComponentID"`
}

// TableName returns the table name for StaffSalaryComponent.
func (StaffSalaryComponent) TableName() string {
	return "staff_salary_components"
}
