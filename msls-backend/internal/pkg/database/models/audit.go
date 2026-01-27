// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action being audited.
type AuditAction string

// Audit action constants.
const (
	AuditActionCreate         AuditAction = "create"
	AuditActionUpdate         AuditAction = "update"
	AuditActionDelete         AuditAction = "delete"
	AuditActionLogin          AuditAction = "login"
	AuditActionLogout         AuditAction = "logout"
	AuditActionLoginFailed    AuditAction = "login_failed"
	AuditActionPasswordChange AuditAction = "password_change"
	AuditActionPasswordReset  AuditAction = "password_reset"
	AuditActionEmailVerified  AuditAction = "email_verified"
	AuditActionRoleAssigned   AuditAction = "role_assigned"
	AuditActionRoleRemoved    AuditAction = "role_removed"
	AuditActionAccountLocked   AuditAction = "account_locked"
	AuditActionAccountUnlocked AuditAction = "account_unlocked"
	AuditAction2FASetup        AuditAction = "2fa_setup"
	AuditAction2FAEnabled      AuditAction = "2fa_enabled"
	AuditAction2FADisabled     AuditAction = "2fa_disabled"
	AuditAction2FAValidated    AuditAction = "2fa_validated"
	AuditAction2FAFailed       AuditAction = "2fa_failed"
	AuditActionBackupCodesUsed           AuditAction = "backup_code_used"
	AuditActionBackupCodesRegen          AuditAction = "backup_codes_regenerated"
	AuditActionAccountDeletionRequested  AuditAction = "account_deletion_requested"
	AuditActionProfileUpdated            AuditAction = "profile_updated"
	AuditActionAvatarUploaded            AuditAction = "avatar_uploaded"
	AuditActionPreferencesUpdated        AuditAction = "preferences_updated"
	AuditActionTokenRefresh              AuditAction = "token_refresh"
	AuditActionTokenRefreshFailed        AuditAction = "token_refresh_failed"
	AuditActionTokenRevoked              AuditAction = "token_revoked"
)

// String returns the string representation of the audit action.
func (a AuditAction) String() string {
	return string(a)
}

// AuditLog represents an audit trail entry.
type AuditLog struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID   *uuid.UUID      `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID     *uuid.UUID      `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action     AuditAction     `gorm:"type:varchar(100);not null;index" json:"action"`
	EntityType string          `gorm:"type:varchar(100);not null;index" json:"entity_type"`
	EntityID   *uuid.UUID      `gorm:"type:uuid;index" json:"entity_id,omitempty"`
	OldData    json.RawMessage `gorm:"type:jsonb" json:"old_data,omitempty"`
	NewData    json.RawMessage `gorm:"type:jsonb" json:"new_data,omitempty"`
	IPAddress  net.IP          `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent  string          `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt  time.Time       `gorm:"not null;default:now();index" json:"created_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for the AuditLog model.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// SetOldData sets the old data from any struct.
func (a *AuditLog) SetOldData(data interface{}) error {
	if data == nil {
		a.OldData = nil
		return nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	a.OldData = bytes
	return nil
}

// SetNewData sets the new data from any struct.
func (a *AuditLog) SetNewData(data interface{}) error {
	if data == nil {
		a.NewData = nil
		return nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	a.NewData = bytes
	return nil
}

// GetOldData unmarshals old data into the provided struct.
func (a *AuditLog) GetOldData(v interface{}) error {
	if a.OldData == nil {
		return nil
	}
	return json.Unmarshal(a.OldData, v)
}

// GetNewData unmarshals new data into the provided struct.
func (a *AuditLog) GetNewData(v interface{}) error {
	if a.NewData == nil {
		return nil
	}
	return json.Unmarshal(a.NewData, v)
}

// AuditLogBuilder provides a fluent interface for creating audit logs.
type AuditLogBuilder struct {
	log *AuditLog
}

// NewAuditLog creates a new AuditLogBuilder.
func NewAuditLog(action AuditAction, entityType string) *AuditLogBuilder {
	return &AuditLogBuilder{
		log: &AuditLog{
			Action:     action,
			EntityType: entityType,
			CreatedAt:  time.Now(),
		},
	}
}

// WithTenant sets the tenant ID.
func (b *AuditLogBuilder) WithTenant(tenantID uuid.UUID) *AuditLogBuilder {
	b.log.TenantID = &tenantID
	return b
}

// WithUser sets the user ID.
func (b *AuditLogBuilder) WithUser(userID uuid.UUID) *AuditLogBuilder {
	b.log.UserID = &userID
	return b
}

// WithEntity sets the entity ID.
func (b *AuditLogBuilder) WithEntity(entityID uuid.UUID) *AuditLogBuilder {
	b.log.EntityID = &entityID
	return b
}

// WithOldData sets the old data.
func (b *AuditLogBuilder) WithOldData(data interface{}) *AuditLogBuilder {
	_ = b.log.SetOldData(data)
	return b
}

// WithNewData sets the new data.
func (b *AuditLogBuilder) WithNewData(data interface{}) *AuditLogBuilder {
	_ = b.log.SetNewData(data)
	return b
}

// WithIPAddress sets the IP address.
func (b *AuditLogBuilder) WithIPAddress(ip net.IP) *AuditLogBuilder {
	b.log.IPAddress = ip
	return b
}

// WithUserAgent sets the user agent.
func (b *AuditLogBuilder) WithUserAgent(userAgent string) *AuditLogBuilder {
	b.log.UserAgent = userAgent
	return b
}

// Build returns the constructed AuditLog.
func (b *AuditLogBuilder) Build() *AuditLog {
	return b.log
}
