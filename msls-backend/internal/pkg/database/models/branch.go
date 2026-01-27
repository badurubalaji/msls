// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BranchAddress represents the physical address of a branch.
type BranchAddress struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
}

// BranchSettings represents configurable settings for a branch.
type BranchSettings struct {
	// OperatingHours specifies the branch operating hours.
	OperatingHours map[string]string `json:"operating_hours,omitempty"`
	// ContactPhone is the branch contact phone number.
	ContactPhone string `json:"contact_phone,omitempty"`
	// ContactEmail is the branch contact email address.
	ContactEmail string `json:"contact_email,omitempty"`
	// Features contains feature flags specific to this branch.
	Features map[string]bool `json:"features,omitempty"`
}

// Branch represents a physical or logical location belonging to a tenant.
type Branch struct {
	AuditModel
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"type:varchar(50);not null" json:"code"`
	Address   BranchAddress  `gorm:"type:jsonb;serializer:json;not null;default:'{}'" json:"address"`
	Settings  BranchSettings `gorm:"type:jsonb;serializer:json;not null;default:'{}'" json:"settings"`
	IsPrimary bool           `gorm:"not null;default:false" json:"is_primary"`
	Status    Status         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TableName returns the table name for the Branch model.
func (Branch) TableName() string {
	return "branches"
}

// BeforeCreate hook for Branch.
func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	if err := b.AuditModel.BeforeCreate(tx); err != nil {
		return err
	}
	if b.Status == "" {
		b.Status = StatusActive
	}
	return nil
}

// Validate performs validation on the Branch model.
func (b *Branch) Validate() error {
	if b.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if b.Name == "" {
		return ErrBranchNameRequired
	}
	if b.Code == "" {
		return ErrBranchCodeRequired
	}
	if !b.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

// MarshalAddress converts BranchAddress to JSON bytes.
func (b *Branch) MarshalAddress() ([]byte, error) {
	return json.Marshal(b.Address)
}

// UnmarshalAddress parses JSON bytes into BranchAddress.
func (b *Branch) UnmarshalAddress(data []byte) error {
	return json.Unmarshal(data, &b.Address)
}

// MarshalSettings converts BranchSettings to JSON bytes.
func (b *Branch) MarshalSettings() ([]byte, error) {
	return json.Marshal(b.Settings)
}

// UnmarshalSettings parses JSON bytes into BranchSettings.
func (b *Branch) UnmarshalSettings(data []byte) error {
	return json.Unmarshal(data, &b.Settings)
}
