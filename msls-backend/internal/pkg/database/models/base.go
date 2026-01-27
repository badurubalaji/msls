// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all database models.
// It includes the primary key and timestamp fields.
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// AuditModel extends BaseModel with audit trail fields.
// Use this for tables that track who created/updated records.
type AuditModel struct {
	BaseModel
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// TenantModel extends AuditModel with tenant isolation.
// Use this for all multi-tenant tables.
type TenantModel struct {
	AuditModel
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
}

// BeforeCreate generates a new UUID v7 if the ID is not set.
// This is a fallback in case the database default is not triggered.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Status represents the status of an entity.
type Status string

// Status constants for entities.
const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
	StatusPending   Status = "pending"
)

// IsValid checks if the status is a valid value.
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusSuspended, StatusPending:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}
