// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token stored in the database.
type RefreshToken struct {
	BaseModel
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null;index" json:"expires_at"`
	RevokedAt *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the RefreshToken model.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired returns true if the refresh token has expired.
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsRevoked returns true if the refresh token has been revoked.
func (r *RefreshToken) IsRevoked() bool {
	return r.RevokedAt != nil
}

// IsValid returns true if the refresh token is valid (not expired and not revoked).
func (r *RefreshToken) IsValid() bool {
	return !r.IsExpired() && !r.IsRevoked()
}

// Revoke marks the refresh token as revoked.
func (r *RefreshToken) Revoke() {
	now := time.Now()
	r.RevokedAt = &now
}

// VerificationTokenType represents the type of verification token.
type VerificationTokenType string

const (
	// VerificationTokenTypeEmailVerify is used for email verification.
	VerificationTokenTypeEmailVerify VerificationTokenType = "email_verify"
	// VerificationTokenTypePasswordReset is used for password reset.
	VerificationTokenTypePasswordReset VerificationTokenType = "password_reset"
	// VerificationTokenTypePhoneVerify is used for phone verification.
	VerificationTokenTypePhoneVerify VerificationTokenType = "phone_verify"
)

// IsValid checks if the verification token type is valid.
func (t VerificationTokenType) IsValid() bool {
	switch t {
	case VerificationTokenTypeEmailVerify, VerificationTokenTypePasswordReset, VerificationTokenTypePhoneVerify:
		return true
	}
	return false
}

// String returns the string representation of the verification token type.
func (t VerificationTokenType) String() string {
	return string(t)
}

// VerificationToken represents a verification token (email verify, password reset, etc.).
type VerificationToken struct {
	BaseModel
	UserID    uuid.UUID             `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string                `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	Type      VerificationTokenType `gorm:"type:varchar(50);not null;index" json:"type"`
	ExpiresAt time.Time             `gorm:"type:timestamptz;not null;index" json:"expires_at"`
	UsedAt    *time.Time            `gorm:"type:timestamptz" json:"used_at,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the VerificationToken model.
func (VerificationToken) TableName() string {
	return "verification_tokens"
}

// IsExpired returns true if the verification token has expired.
func (v *VerificationToken) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsUsed returns true if the verification token has been used.
func (v *VerificationToken) IsUsed() bool {
	return v.UsedAt != nil
}

// IsValid returns true if the verification token is valid (not expired and not used).
func (v *VerificationToken) IsValid() bool {
	return !v.IsExpired() && !v.IsUsed()
}

// MarkUsed marks the verification token as used.
func (v *VerificationToken) MarkUsed() {
	now := time.Now()
	v.UsedAt = &now
}

// Validate performs validation on the VerificationToken model.
func (v *VerificationToken) Validate() error {
	if v.UserID == uuid.Nil {
		return ErrUserIDRequired
	}
	if v.TokenHash == "" {
		return ErrTokenHashRequired
	}
	if !v.Type.IsValid() {
		return ErrInvalidTokenType
	}
	return nil
}

// LoginAttempt represents a login attempt record for security monitoring.
type LoginAttempt struct {
	BaseModel
	UserID        *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Email         string     `gorm:"type:varchar(255);index" json:"email,omitempty"`
	IPAddress     net.IP     `gorm:"type:inet;not null;index" json:"ip_address"`
	UserAgent     string     `gorm:"type:text" json:"user_agent,omitempty"`
	Success       bool       `gorm:"not null;default:false" json:"success"`
	FailureReason string     `gorm:"type:varchar(100)" json:"failure_reason,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the LoginAttempt model.
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// FailureReasons for login attempts.
const (
	LoginFailureInvalidCredentials = "invalid_credentials"
	LoginFailureAccountLocked      = "account_locked"
	LoginFailureAccountInactive    = "account_inactive"
	LoginFailureEmailNotVerified   = "email_not_verified"
	LoginFailureTenantInactive     = "tenant_inactive"
	LoginFailureUserNotFound       = "user_not_found"
)
