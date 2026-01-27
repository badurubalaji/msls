// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
)

// OTPType represents the type of OTP code.
type OTPType string

// OTPType constants.
const (
	OTPTypeLogin       OTPType = "login"
	OTPTypeVerify      OTPType = "verify"
	OTPTypePhoneVerify OTPType = "phone_verify"
)

// IsValid checks if the OTP type is valid.
func (t OTPType) IsValid() bool {
	switch t {
	case OTPTypeLogin, OTPTypeVerify, OTPTypePhoneVerify:
		return true
	}
	return false
}

// String returns the string representation of the OTP type.
func (t OTPType) String() string {
	return string(t)
}

// OTPChannel represents the delivery channel for OTP.
type OTPChannel string

// OTPChannel constants.
const (
	OTPChannelSMS   OTPChannel = "sms"
	OTPChannelEmail OTPChannel = "email"
)

// IsValid checks if the OTP channel is valid.
func (c OTPChannel) IsValid() bool {
	switch c {
	case OTPChannelSMS, OTPChannelEmail:
		return true
	}
	return false
}

// String returns the string representation of the OTP channel.
func (c OTPChannel) String() string {
	return string(c)
}

// OTPCode represents an OTP code for passwordless authentication.
type OTPCode struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Identifier string     `gorm:"type:varchar(255);not null;index" json:"identifier"`
	CodeHash   string     `gorm:"type:varchar(255);not null" json:"-"`
	Type       OTPType    `gorm:"type:otp_type;not null" json:"type"`
	Channel    OTPChannel `gorm:"type:otp_channel;not null" json:"channel"`
	ExpiresAt  time.Time  `gorm:"type:timestamptz;not null;index" json:"expires_at"`
	VerifiedAt *time.Time `gorm:"type:timestamptz" json:"verified_at,omitempty"`
	Attempts   int        `gorm:"not null;default:0" json:"attempts"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the OTPCode model.
func (OTPCode) TableName() string {
	return "otp_codes"
}

// IsExpired returns true if the OTP has expired.
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsVerified returns true if the OTP has been verified.
func (o *OTPCode) IsVerified() bool {
	return o.VerifiedAt != nil
}

// IsValid returns true if the OTP is valid (not expired, not verified, and within attempt limit).
func (o *OTPCode) IsValid() bool {
	return !o.IsExpired() && !o.IsVerified() && o.Attempts < MaxOTPAttempts
}

// IncrementAttempts increments the attempt counter.
func (o *OTPCode) IncrementAttempts() {
	o.Attempts++
}

// MarkVerified marks the OTP as verified.
func (o *OTPCode) MarkVerified() {
	now := time.Now()
	o.VerifiedAt = &now
}

// HasExceededMaxAttempts returns true if max attempts have been exceeded.
func (o *OTPCode) HasExceededMaxAttempts() bool {
	return o.Attempts >= MaxOTPAttempts
}

// Validate performs validation on the OTPCode model.
func (o *OTPCode) Validate() error {
	if o.Identifier == "" {
		return ErrIdentifierRequired
	}
	if o.CodeHash == "" {
		return ErrCodeHashRequired
	}
	if !o.Type.IsValid() {
		return ErrInvalidOTPType
	}
	if !o.Channel.IsValid() {
		return ErrInvalidOTPChannel
	}
	return nil
}

// OTPRateLimit tracks OTP request rate limits per identifier.
type OTPRateLimit struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	Identifier    string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_otp_rate_limits_identifier_channel" json:"identifier"`
	Channel       OTPChannel `gorm:"type:otp_channel;not null;uniqueIndex:idx_otp_rate_limits_identifier_channel" json:"channel"`
	RequestCount  int        `gorm:"not null;default:1" json:"request_count"`
	WindowStart   time.Time  `gorm:"type:timestamptz;not null;default:now();index" json:"window_start"`
	LastRequestAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"last_request_at"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName returns the table name for the OTPRateLimit model.
func (OTPRateLimit) TableName() string {
	return "otp_rate_limits"
}

// IsWithinWindow returns true if the current time is within the rate limit window.
func (r *OTPRateLimit) IsWithinWindow(windowDuration time.Duration) bool {
	return time.Now().Before(r.WindowStart.Add(windowDuration))
}

// HasExceededLimit returns true if the rate limit has been exceeded.
func (r *OTPRateLimit) HasExceededLimit(maxRequests int, windowDuration time.Duration) bool {
	if !r.IsWithinWindow(windowDuration) {
		return false
	}
	return r.RequestCount >= maxRequests
}

// IncrementCount increments the request count.
func (r *OTPRateLimit) IncrementCount() {
	r.RequestCount++
	r.LastRequestAt = time.Now()
	r.UpdatedAt = time.Now()
}

// ResetWindow resets the rate limit window.
func (r *OTPRateLimit) ResetWindow() {
	r.RequestCount = 1
	r.WindowStart = time.Now()
	r.LastRequestAt = time.Now()
	r.UpdatedAt = time.Now()
}

// CanRequestNewOTP returns true if a new OTP can be requested (respects cooldown).
func (r *OTPRateLimit) CanRequestNewOTP(cooldownDuration time.Duration) bool {
	return time.Now().After(r.LastRequestAt.Add(cooldownDuration))
}

// OTP configuration constants.
const (
	MaxOTPAttempts      = 5               // Maximum verification attempts per OTP
	MaxOTPRequestsPerHr = 3               // Maximum OTP requests per hour per identifier
	OTPExpiryDuration   = 5 * time.Minute // OTP validity duration
	OTPCooldownDuration = 60 * time.Second // Minimum time between OTP requests
	OTPRateLimitWindow  = time.Hour       // Rate limit window duration
)
