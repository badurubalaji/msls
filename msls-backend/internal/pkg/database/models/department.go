// Package models contains database model definitions.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Department represents a department in the organization.
type Department struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	BranchID    uuid.UUID  `gorm:"type:uuid;not null;index"`

	Name        string  `gorm:"type:varchar(100);not null"`
	Code        string  `gorm:"type:varchar(20);not null"`
	Description *string `gorm:"type:text"`
	HeadID      *uuid.UUID `gorm:"type:uuid"` // References staff(id)

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`

	// Relations
	Branch *Branch `gorm:"foreignKey:BranchID"`
	Head   *Staff  `gorm:"foreignKey:HeadID"`
}

// TableName returns the table name for Department.
func (Department) TableName() string {
	return "departments"
}

// Designation represents a job designation in the organization.
type Designation struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name         string     `gorm:"type:varchar(100);not null"`
	Level        int        `gorm:"not null;default:1"` // 1 = highest (CEO), 10 = lowest
	DepartmentID *uuid.UUID `gorm:"type:uuid;index"`

	IsActive bool `gorm:"not null;default:true"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`

	// Relations
	Department *Department `gorm:"foreignKey:DepartmentID"`
}

// TableName returns the table name for Designation.
func (Designation) TableName() string {
	return "designations"
}
