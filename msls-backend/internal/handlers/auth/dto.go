// Package auth provides HTTP handlers for authentication endpoints.
package auth

import (
	"time"

	"github.com/google/uuid"
)

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
}

// LoginResponse represents the login response.
type LoginResponse struct {
	User             UserDTO `json:"user"`
	AccessToken      string  `json:"access_token"`
	RefreshToken     string  `json:"refresh_token"`
	ExpiresIn        int     `json:"expires_in"`          // seconds
	TwoFactorEnabled bool    `json:"two_factor_enabled"`  // indicates if 2FA is enabled
	RequiresTwoFactor bool    `json:"requires_two_factor"` // indicates if 2FA verification is needed
	PartialToken     string  `json:"partial_token,omitempty"` // temporary token for 2FA verification
}

// RegisterRequest represents the user registration request body.
type RegisterRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=8"`
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	RoleIDs   []string `json:"role_ids" binding:"required,min=1,dive,uuid"`
}

// RefreshTokenRequest represents the refresh token request body.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the refresh token response.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// LogoutRequest represents the logout request body.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// VerifyEmailRequest represents the email verification request body.
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// ForgotPasswordRequest represents the forgot password request body.
type ForgotPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
}

// ResetPasswordRequest represents the reset password request body.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserDTO represents a user in API responses.
type UserDTO struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	Email            string     `json:"email,omitempty"`
	Phone            string     `json:"phone,omitempty"`
	FirstName        string     `json:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty"`
	FullName         string     `json:"full_name,omitempty"`
	Status           string     `json:"status"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt  *time.Time `json:"phone_verified_at,omitempty"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	Roles            []RoleDTO  `json:"roles,omitempty"`
	Permissions      []string   `json:"permissions,omitempty"`
}

// RoleDTO represents a role in API responses.
type RoleDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
}

// MeResponse represents the current user info response.
type MeResponse struct {
	User        UserDTO  `json:"user"`
	Permissions []string `json:"permissions"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// ==================== 2FA DTOs ====================

// TwoFactorSetupResponse represents the 2FA setup response.
type TwoFactorSetupResponse struct {
	Secret        string `json:"secret"`
	QRCodeDataURL string `json:"qr_code_data_url"`
	ManualEntry   string `json:"manual_entry"`
}

// TwoFactorVerifyRequest represents the request to verify and enable 2FA.
type TwoFactorVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// TwoFactorVerifyResponse represents the response after enabling 2FA.
type TwoFactorVerifyResponse struct {
	Message     string   `json:"message"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFactorValidateRequest represents the request to validate 2FA during login.
type TwoFactorValidateRequest struct {
	PartialToken string `json:"partial_token" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// TwoFactorDisableRequest represents the request to disable 2FA.
type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"required"`
}

// TwoFactorStatusResponse represents the 2FA status response.
type TwoFactorStatusResponse struct {
	Enabled          bool   `json:"enabled"`
	BackupCodesCount int64  `json:"backup_codes_count"`
}

// BackupCodesResponse represents the backup codes response.
type BackupCodesResponse struct {
	BackupCodes      []string `json:"backup_codes,omitempty"`
	BackupCodesCount int64    `json:"backup_codes_count"`
}

// RegenerateBackupCodesRequest represents the request to regenerate backup codes.
type RegenerateBackupCodesRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// ==================== OTP DTOs ====================

// OTPRequestRequest represents the request to send an OTP.
type OTPRequestRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Phone number or email
	Type       string `json:"type" binding:"required,oneof=sms email"` // Delivery channel
	TenantID   string `json:"tenant_id,omitempty"` // Optional tenant ID
}

// OTPRequestResponse represents the response after sending an OTP.
type OTPRequestResponse struct {
	Message          string `json:"message"`
	ExpiresIn        int    `json:"expires_in"` // seconds
	MaskedIdentifier string `json:"masked_identifier"`
}

// OTPVerifyRequest represents the request to verify an OTP.
type OTPVerifyRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Phone number or email
	Code       string `json:"code" binding:"required,len=6"` // 6-digit OTP
	TenantID   string `json:"tenant_id,omitempty"` // Optional tenant ID
}

// OTPVerifyResponse represents the response after verifying an OTP.
type OTPVerifyResponse struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"` // seconds
}

// OTPResendRequest represents the request to resend an OTP.
type OTPResendRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Phone number or email
	Type       string `json:"type" binding:"required,oneof=sms email"` // Delivery channel
	TenantID   string `json:"tenant_id,omitempty"` // Optional tenant ID
}
