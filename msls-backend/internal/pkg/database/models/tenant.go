// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// TenantSettings represents configurable settings for a tenant.
type TenantSettings struct {
	// Timezone specifies the default timezone for the tenant.
	Timezone string `json:"timezone,omitempty"`
	// Currency specifies the default currency code (e.g., "USD", "INR").
	Currency string `json:"currency,omitempty"`
	// Locale specifies the default locale (e.g., "en-US", "hi-IN").
	Locale string `json:"locale,omitempty"`
	// Features contains feature flags for the tenant.
	Features map[string]bool `json:"features,omitempty"`
}

// Tenant represents an organization using the MSLS system.
// Each tenant has complete data isolation via Row-Level Security.
type Tenant struct {
	BaseModel
	Name     string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug     string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	Settings TenantSettings `gorm:"type:jsonb;serializer:json;not null;default:'{}'" json:"settings"`
	Status   Status         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// Relationships
	Branches []Branch `gorm:"foreignKey:TenantID" json:"branches,omitempty"`
	Users    []User   `gorm:"foreignKey:TenantID" json:"users,omitempty"`
}

// TableName returns the table name for the Tenant model.
func (Tenant) TableName() string {
	return "tenants"
}

// BeforeCreate hook for Tenant.
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if err := t.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	if t.Status == "" {
		t.Status = StatusActive
	}
	return nil
}

// Validate performs validation on the Tenant model.
func (t *Tenant) Validate() error {
	if t.Name == "" {
		return ErrTenantNameRequired
	}
	if t.Slug == "" {
		return ErrTenantSlugRequired
	}
	if !t.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

// MarshalSettings converts TenantSettings to JSON bytes.
func (t *Tenant) MarshalSettings() ([]byte, error) {
	return json.Marshal(t.Settings)
}

// UnmarshalSettings parses JSON bytes into TenantSettings.
func (t *Tenant) UnmarshalSettings(data []byte) error {
	return json.Unmarshal(data, &t.Settings)
}
