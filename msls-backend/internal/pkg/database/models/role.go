// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a system permission for access control.
type Permission struct {
	BaseModel
	Code        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Module      string `gorm:"type:varchar(100);not null;index" json:"module"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	// Relationships
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName returns the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}

// Role represents a role that can be assigned to users.
type Role struct {
	BaseModel
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool       `gorm:"not null;default:false;index" json:"is_system"`

	// Relationships
	Tenant      *Tenant      `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName returns the table name for the Role model.
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate hook for Role.
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	return r.BaseModel.BeforeCreate(tx)
}

// Validate performs validation on the Role model.
func (r *Role) Validate() error {
	if r.Name == "" {
		return ErrRoleNameRequired
	}
	return nil
}

// IsSystemRole returns true if this is a system-defined role.
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// IsTenantRole returns true if this role belongs to a specific tenant.
func (r *Role) IsTenantRole() bool {
	return r.TenantID != nil
}

// HasPermission checks if the role has a specific permission.
func (r *Role) HasPermission(code string) bool {
	for _, p := range r.Permissions {
		if p.Code == code {
			return true
		}
	}
	return false
}

// GetPermissionCodes returns all permission codes for this role.
func (r *Role) GetPermissionCodes() []string {
	codes := make([]string, len(r.Permissions))
	for i, p := range r.Permissions {
		codes[i] = p.Code
	}
	return codes
}

// RolePermission represents the junction table between roles and permissions.
type RolePermission struct {
	BaseModel
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;index" json:"permission_id"`
}

// TableName returns the table name for the RolePermission model.
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole represents the junction table between users and roles.
type UserRole struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
}

// TableName returns the table name for the UserRole model.
func (UserRole) TableName() string {
	return "user_roles"
}
