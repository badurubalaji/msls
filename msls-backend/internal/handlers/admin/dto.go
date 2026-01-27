// Package admin provides HTTP handlers for administrative endpoints.
package admin

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FeatureFlagDTO represents a feature flag in API responses.
type FeatureFlagDTO struct {
	ID           uuid.UUID              `json:"id"`
	Key          string                 `json:"key"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	DefaultValue bool                   `json:"default_value"`
	Metadata     FeatureFlagMetadataDTO `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// FeatureFlagMetadataDTO represents feature flag metadata in API responses.
type FeatureFlagMetadataDTO struct {
	Category          string `json:"category,omitempty"`
	RequiresSetup     bool   `json:"requires_setup,omitempty"`
	Beta              bool   `json:"beta,omitempty"`
	RolloutPercentage *int   `json:"rollout_percentage,omitempty"`
}

// CreateFeatureFlagRequest represents a request to create a feature flag.
type CreateFeatureFlagRequest struct {
	Key          string                  `json:"key" binding:"required,min=1,max=100"`
	Name         string                  `json:"name" binding:"required,min=1,max=255"`
	Description  string                  `json:"description"`
	DefaultValue bool                    `json:"default_value"`
	Metadata     *FeatureFlagMetadataDTO `json:"metadata"`
}

// UpdateFeatureFlagRequest represents a request to update a feature flag.
type UpdateFeatureFlagRequest struct {
	Name         *string                 `json:"name" binding:"omitempty,min=1,max=255"`
	Description  *string                 `json:"description"`
	DefaultValue *bool                   `json:"default_value"`
	Metadata     *FeatureFlagMetadataDTO `json:"metadata"`
}

// TenantFeatureFlagDTO represents a tenant feature flag override in API responses.
type TenantFeatureFlagDTO struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	FlagKey     string          `json:"flag_key"`
	FlagName    string          `json:"flag_name"`
	Enabled     bool            `json:"enabled"`
	CustomValue json.RawMessage `json:"custom_value,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// SetTenantFlagsRequest represents a request to set tenant feature flags.
type SetTenantFlagsRequest struct {
	Flags []TenantFlagOverride `json:"flags" binding:"required,dive"`
}

// TenantFlagOverride represents a single flag override in the request.
type TenantFlagOverride struct {
	Key         string          `json:"key" binding:"required"`
	Enabled     bool            `json:"enabled"`
	CustomValue json.RawMessage `json:"custom_value,omitempty"`
}

// UserFeatureFlagDTO represents a user feature flag override in API responses.
type UserFeatureFlagDTO struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FlagKey   string    `json:"flag_key"`
	FlagName  string    `json:"flag_name"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetUserFlagsRequest represents a request to set user feature flags.
type SetUserFlagsRequest struct {
	Flags []UserFlagOverride `json:"flags" binding:"required,dive"`
}

// UserFlagOverride represents a single flag override for a user.
type UserFlagOverride struct {
	Key     string `json:"key" binding:"required"`
	Enabled bool   `json:"enabled"`
}

// FeatureFlagStateDTO represents the state of a feature flag for a user.
type FeatureFlagStateDTO struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Enabled     bool            `json:"enabled"`
	CustomValue json.RawMessage `json:"custom_value,omitempty"`
	Source      string          `json:"source"` // "default", "tenant", "user"
}

// ListFeatureFlagsResponse represents the response for listing feature flags.
type ListFeatureFlagsResponse struct {
	Flags []FeatureFlagDTO `json:"flags"`
}

// CurrentFlagsResponse represents the response for getting current user's flags.
type CurrentFlagsResponse struct {
	Flags []FeatureFlagStateDTO `json:"flags"`
}
