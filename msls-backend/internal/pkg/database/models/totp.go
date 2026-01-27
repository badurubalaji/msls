// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// BackupCode represents a backup code for 2FA account recovery.
type BackupCode struct {
	BaseModel
	UserID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	CodeHash string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	UsedAt   *time.Time `gorm:"type:timestamptz" json:"used_at,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the BackupCode model.
func (BackupCode) TableName() string {
	return "backup_codes"
}

// IsUsed returns true if the backup code has been used.
func (b *BackupCode) IsUsed() bool {
	return b.UsedAt != nil
}

// MarkUsed marks the backup code as used.
func (b *BackupCode) MarkUsed() {
	now := time.Now()
	b.UsedAt = &now
}

// TOTPAttempt represents a TOTP validation attempt for rate limiting.
type TOTPAttempt struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	IPAddress net.IP    `gorm:"type:inet;not null;index" json:"ip_address"`
	Success   bool      `gorm:"not null;default:false" json:"success"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the TOTPAttempt model.
func (TOTPAttempt) TableName() string {
	return "totp_attempts"
}
