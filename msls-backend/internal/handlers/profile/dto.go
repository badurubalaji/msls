// Package profile provides HTTP handlers for profile management endpoints.
package profile

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
	profileservice "msls-backend/internal/services/profile"
)

// UpdateProfileRequest represents the profile update request body.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Timezone  *string `json:"timezone,omitempty"`
	Locale    *string `json:"locale,omitempty"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *UpdateProfileRequest) ToServiceRequest() profileservice.UpdateProfileRequest {
	return profileservice.UpdateProfileRequest{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Phone:     r.Phone,
		Bio:       r.Bio,
		Timezone:  r.Timezone,
		Locale:    r.Locale,
	}
}

// ChangePasswordRequest represents the password change request body.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *ChangePasswordRequest) ToServiceRequest() profileservice.ChangePasswordRequest {
	return profileservice.ChangePasswordRequest{
		CurrentPassword: r.CurrentPassword,
		NewPassword:     r.NewPassword,
		ConfirmPassword: r.ConfirmPassword,
	}
}

// UpdatePreferencesRequest represents the notification preferences update request body.
type UpdatePreferencesRequest struct {
	Email *bool `json:"email,omitempty"`
	Push  *bool `json:"push,omitempty"`
	SMS   *bool `json:"sms,omitempty"`
}

// ToServiceRequest converts the DTO to a service request.
func (r *UpdatePreferencesRequest) ToServiceRequest() profileservice.UpdatePreferencesRequest {
	return profileservice.UpdatePreferencesRequest{
		Email: r.Email,
		Push:  r.Push,
		SMS:   r.SMS,
	}
}

// NotificationPreferencesResponse represents the notification preferences response.
type NotificationPreferencesResponse struct {
	Email bool `json:"email"`
	Push  bool `json:"push"`
	SMS   bool `json:"sms"`
}

// ProfileResponse represents a user profile in API responses.
type ProfileResponse struct {
	ID                      uuid.UUID                        `json:"id"`
	TenantID                uuid.UUID                        `json:"tenant_id"`
	Email                   string                           `json:"email,omitempty"`
	Phone                   string                           `json:"phone,omitempty"`
	FirstName               string                           `json:"first_name,omitempty"`
	LastName                string                           `json:"last_name,omitempty"`
	FullName                string                           `json:"full_name,omitempty"`
	AvatarURL               string                           `json:"avatar_url,omitempty"`
	Bio                     string                           `json:"bio,omitempty"`
	Timezone                string                           `json:"timezone"`
	Locale                  string                           `json:"locale"`
	NotificationPreferences *NotificationPreferencesResponse `json:"notification_preferences,omitempty"`
	Status                  string                           `json:"status"`
	EmailVerifiedAt         *time.Time                       `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt         *time.Time                       `json:"phone_verified_at,omitempty"`
	LastLoginAt             *time.Time                       `json:"last_login_at,omitempty"`
	TwoFactorEnabled        bool                             `json:"two_factor_enabled"`
	CreatedAt               time.Time                        `json:"created_at"`
	Roles                   []RoleDTO                        `json:"roles,omitempty"`
	Permissions             []string                         `json:"permissions,omitempty"`
}

// RoleDTO represents a role in API responses.
type RoleDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
}

// AvatarUploadResponse represents the avatar upload response.
type AvatarUploadResponse struct {
	AvatarURL string `json:"avatar_url"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// UserPreferenceDTO represents a user preference in API responses.
type UserPreferenceDTO struct {
	Category string      `json:"category"`
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
}

// UserPreferencesResponse represents a list of user preferences.
type UserPreferencesResponse struct {
	Preferences []UserPreferenceDTO `json:"preferences"`
}

// SetPreferenceRequest represents a request to set a user preference.
type SetPreferenceRequest struct {
	Category string      `json:"category" binding:"required"`
	Key      string      `json:"key" binding:"required"`
	Value    interface{} `json:"value" binding:"required"`
}

// UserToProfileResponse converts a User model to a ProfileResponse.
func UserToProfileResponse(user *models.User) ProfileResponse {
	resp := ProfileResponse{
		ID:               user.ID,
		TenantID:         user.TenantID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		FullName:         user.FullName(),
		Timezone:         user.Timezone,
		Locale:           user.Locale,
		Status:           string(user.Status),
		EmailVerifiedAt:  user.EmailVerifiedAt,
		PhoneVerifiedAt:  user.PhoneVerifiedAt,
		LastLoginAt:      user.LastLoginAt,
		TwoFactorEnabled: user.TwoFactorEnabled,
		CreatedAt:        user.CreatedAt,
		Permissions:      user.GetPermissions(),
	}

	if user.Email != nil {
		resp.Email = *user.Email
	}
	if user.Phone != nil {
		resp.Phone = *user.Phone
	}
	if user.AvatarURL != nil {
		resp.AvatarURL = *user.AvatarURL
	}
	if user.Bio != nil {
		resp.Bio = *user.Bio
	}

	// Parse notification preferences
	if user.NotificationPreferences != nil {
		var prefs profileservice.NotificationPreferences
		if err := json.Unmarshal(user.NotificationPreferences, &prefs); err == nil {
			resp.NotificationPreferences = &NotificationPreferencesResponse{
				Email: prefs.Email,
				Push:  prefs.Push,
				SMS:   prefs.SMS,
			}
		}
	}

	// Convert roles
	resp.Roles = make([]RoleDTO, len(user.Roles))
	for i, role := range user.Roles {
		resp.Roles[i] = RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
		}
	}

	return resp
}
