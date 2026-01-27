// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account belonging to a tenant.
type User struct {
	AuditModel
	TenantID                   uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Email                      *string    `gorm:"type:varchar(255)" json:"email,omitempty"`
	Phone                      *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	PasswordHash               *string    `gorm:"type:varchar(255)" json:"-"`
	FirstName                  string     `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName                   string     `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	AvatarURL                  *string    `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	Bio                        *string    `gorm:"type:text" json:"bio,omitempty"`
	Timezone                   string     `gorm:"type:varchar(50);default:'UTC'" json:"timezone"`
	Locale                     string     `gorm:"type:varchar(10);default:'en'" json:"locale"`
	NotificationPreferences    []byte     `gorm:"type:jsonb;default:'{\"email\": true, \"push\": true, \"sms\": false}'" json:"-"`
	Status                     Status     `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	TwoFactorEnabled           bool       `gorm:"not null;default:false" json:"two_factor_enabled"`
	TwoFactorSecret            *string    `gorm:"type:varchar(255)" json:"-"`
	TOTPVerifiedAt             *time.Time `gorm:"type:timestamptz" json:"totp_verified_at,omitempty"`
	EmailVerifiedAt            *time.Time `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`
	PhoneVerifiedAt            *time.Time `gorm:"type:timestamptz" json:"phone_verified_at,omitempty"`
	LastLoginAt                *time.Time `gorm:"type:timestamptz" json:"last_login_at,omitempty"`
	LastLoginIP                *string    `gorm:"type:inet" json:"-"`
	LockedUntil                *time.Time `gorm:"type:timestamptz" json:"locked_until,omitempty"`
	FailedLoginAttempts        int        `gorm:"not null;default:0" json:"-"`
	AccountDeletionRequestedAt *time.Time `gorm:"type:timestamptz" json:"account_deletion_requested_at,omitempty"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Roles  []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName returns the table name for the User model.
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.AuditModel.BeforeCreate(tx); err != nil {
		return err
	}
	if u.Status == "" {
		u.Status = StatusActive
	}
	return nil
}

// Validate performs validation on the User model.
func (u *User) Validate() error {
	if u.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if u.Email == nil && u.Phone == nil {
		return ErrEmailOrPhoneRequired
	}
	if !u.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

// IsEmailVerified returns true if the user's email is verified.
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsPhoneVerified returns true if the user's phone is verified.
func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerifiedAt != nil
}

// HasPassword returns true if the user has a password set.
func (u *User) HasPassword() bool {
	return u.PasswordHash != nil && *u.PasswordHash != ""
}

// SetEmail sets the user's email address.
func (u *User) SetEmail(email string) {
	u.Email = &email
	u.EmailVerifiedAt = nil
}

// SetPhone sets the user's phone number.
func (u *User) SetPhone(phone string) {
	u.Phone = &phone
	u.PhoneVerifiedAt = nil
}

// VerifyEmail marks the user's email as verified.
func (u *User) VerifyEmail() {
	now := time.Now()
	u.EmailVerifiedAt = &now
}

// VerifyPhone marks the user's phone as verified.
func (u *User) VerifyPhone() {
	now := time.Now()
	u.PhoneVerifiedAt = &now
}

// RecordLogin updates the last login timestamp.
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// Enable2FA enables two-factor authentication with the given secret.
func (u *User) Enable2FA(secret string) {
	u.TwoFactorEnabled = true
	u.TwoFactorSecret = &secret
}

// Disable2FA disables two-factor authentication.
func (u *User) Disable2FA() {
	u.TwoFactorEnabled = false
	u.TwoFactorSecret = nil
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

// IsLocked returns true if the user account is currently locked.
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// Lock locks the user account for the specified duration.
func (u *User) Lock(duration time.Duration) {
	lockedUntil := time.Now().Add(duration)
	u.LockedUntil = &lockedUntil
}

// Unlock unlocks the user account.
func (u *User) Unlock() {
	u.LockedUntil = nil
	u.FailedLoginAttempts = 0
}

// IncrementFailedAttempts increments the failed login attempts counter.
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
}

// ResetFailedAttempts resets the failed login attempts counter.
func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
}

// IsActive returns true if the user account is active.
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has a specific permission through any of their roles.
func (u *User) HasPermission(permCode string) bool {
	for _, r := range u.Roles {
		if r.HasPermission(permCode) {
			return true
		}
	}
	return false
}

// GetPermissions returns all unique permission codes for this user.
func (u *User) GetPermissions() []string {
	permMap := make(map[string]bool)
	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			permMap[p.Code] = true
		}
	}
	perms := make([]string, 0, len(permMap))
	for code := range permMap {
		perms = append(perms, code)
	}
	return perms
}
