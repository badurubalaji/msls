// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FeatureFlagMetadata contains additional configuration for a feature flag.
type FeatureFlagMetadata struct {
	// Category groups related feature flags
	Category string `json:"category,omitempty"`
	// RequiresSetup indicates if the feature needs additional configuration
	RequiresSetup bool `json:"requires_setup,omitempty"`
	// Beta indicates if this is a beta feature
	Beta bool `json:"beta,omitempty"`
	// RolloutPercentage for gradual rollout (0-100)
	RolloutPercentage *int `json:"rollout_percentage,omitempty"`
}

// FeatureFlag represents a system-wide feature flag.
type FeatureFlag struct {
	ID           uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	Key          string              `gorm:"type:varchar(100);not null;uniqueIndex" json:"key"`
	Name         string              `gorm:"type:varchar(255);not null" json:"name"`
	Description  string              `gorm:"type:text" json:"description,omitempty"`
	DefaultValue bool                `gorm:"not null;default:false" json:"default_value"`
	Metadata     FeatureFlagMetadata `gorm:"type:jsonb;serializer:json;not null;default:'{}'" json:"metadata"`
	CreatedAt    time.Time           `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time           `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	TenantOverrides []TenantFeatureFlag `gorm:"foreignKey:FlagID" json:"tenant_overrides,omitempty"`
	UserOverrides   []UserFeatureFlag   `gorm:"foreignKey:FlagID" json:"user_overrides,omitempty"`
}

// TableName returns the table name for the FeatureFlag model.
func (FeatureFlag) TableName() string {
	return "feature_flags"
}

// BeforeCreate hook for FeatureFlag.
func (f *FeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// Validate performs validation on the FeatureFlag model.
func (f *FeatureFlag) Validate() error {
	if f.Key == "" {
		return ErrFeatureFlagKeyRequired
	}
	if f.Name == "" {
		return ErrFeatureFlagNameRequired
	}
	return nil
}

// TenantFeatureFlag represents a tenant-specific feature flag override.
type TenantFeatureFlag struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	FlagID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"flag_id"`
	Enabled     bool            `gorm:"not null" json:"enabled"`
	CustomValue json.RawMessage `gorm:"type:jsonb" json:"custom_value,omitempty"`
	CreatedAt   time.Time       `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	Tenant      Tenant      `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	FeatureFlag FeatureFlag `gorm:"foreignKey:FlagID" json:"feature_flag,omitempty"`
}

// TableName returns the table name for the TenantFeatureFlag model.
func (TenantFeatureFlag) TableName() string {
	return "tenant_feature_flags"
}

// BeforeCreate hook for TenantFeatureFlag.
func (tf *TenantFeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if tf.ID == uuid.Nil {
		tf.ID = uuid.New()
	}
	return nil
}

// Validate performs validation on the TenantFeatureFlag model.
func (tf *TenantFeatureFlag) Validate() error {
	if tf.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if tf.FlagID == uuid.Nil {
		return ErrFeatureFlagIDRequired
	}
	return nil
}

// GetCustomValueAs unmarshals the custom value into the provided type.
func (tf *TenantFeatureFlag) GetCustomValueAs(v interface{}) error {
	if tf.CustomValue == nil {
		return nil
	}
	return json.Unmarshal(tf.CustomValue, v)
}

// SetCustomValue marshals the provided value into the custom value field.
func (tf *TenantFeatureFlag) SetCustomValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tf.CustomValue = data
	return nil
}

// UserFeatureFlag represents a user-specific feature flag override for beta testing.
type UserFeatureFlag struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	FlagID    uuid.UUID `gorm:"type:uuid;not null;index" json:"flag_id"`
	Enabled   bool      `gorm:"not null" json:"enabled"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FeatureFlag FeatureFlag `gorm:"foreignKey:FlagID" json:"feature_flag,omitempty"`
}

// TableName returns the table name for the UserFeatureFlag model.
func (UserFeatureFlag) TableName() string {
	return "user_feature_flags"
}

// BeforeCreate hook for UserFeatureFlag.
func (uf *UserFeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if uf.ID == uuid.Nil {
		uf.ID = uuid.New()
	}
	return nil
}

// Validate performs validation on the UserFeatureFlag model.
func (uf *UserFeatureFlag) Validate() error {
	if uf.UserID == uuid.Nil {
		return ErrUserIDRequired
	}
	if uf.FlagID == uuid.Nil {
		return ErrFeatureFlagIDRequired
	}
	return nil
}

// FeatureFlagState represents the computed state of a feature flag for a specific context.
type FeatureFlagState struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Enabled     bool            `json:"enabled"`
	CustomValue json.RawMessage `json:"custom_value,omitempty"`
	Source      string          `json:"source"` // "default", "tenant", "user"
}
