// Package rbac provides role-based access control services.
package rbac

import "errors"

// Service layer errors for RBAC operations.
var (
	// Role errors
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleNameRequired    = errors.New("role name is required")
	ErrRoleNameExists      = errors.New("role with this name already exists")
	ErrCannotDeleteSystem  = errors.New("cannot delete system role")
	ErrCannotModifySystem  = errors.New("cannot modify system role permissions")
	ErrRoleInUse           = errors.New("role is assigned to users and cannot be deleted")
	ErrInvalidRoleHierarchy = errors.New("invalid role hierarchy level")

	// Permission errors
	ErrPermissionNotFound     = errors.New("permission not found")
	ErrPermissionCodeRequired = errors.New("permission code is required")
	ErrPermissionCodeExists   = errors.New("permission with this code already exists")
	ErrInvalidPermissionCode  = errors.New("invalid permission code format")

	// User role errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserRoleExists     = errors.New("user already has this role")
	ErrUserRoleNotFound   = errors.New("user does not have this role")
	ErrTenantMismatch     = errors.New("role does not belong to user's tenant")
)

// System role names.
const (
	RoleSuperAdmin   = "SuperAdmin"
	RoleTenantAdmin  = "TenantAdmin"
	RolePrincipal    = "Principal"
	RoleTeacher      = "Teacher"
	RoleStaff        = "Staff"
	RoleParent       = "Parent"
	RoleStudent      = "Student"
)

// Permission categories/modules.
const (
	ModuleUsers     = "users"
	ModuleStudents  = "students"
	ModuleStaff     = "staff"
	ModuleFinance   = "finance"
	ModuleAcademics = "academics"
	ModuleRoles     = "roles"
	ModuleSettings  = "settings"
	ModuleReports   = "reports"
	ModuleBranches  = "branches"
)

// Permission actions.
const (
	ActionRead   = "read"
	ActionWrite  = "write"
	ActionDelete = "delete"
	ActionManage = "manage"
)

// RoleHierarchy defines the hierarchy level for each system role.
// Lower numbers indicate higher privilege levels.
var RoleHierarchy = map[string]int{
	RoleSuperAdmin:  0,
	RoleTenantAdmin: 1,
	RolePrincipal:   2,
	RoleTeacher:     3,
	RoleStaff:       3,
	RoleParent:      4,
	RoleStudent:     5,
}

// SystemRoles returns the list of predefined system roles.
func SystemRoles() []string {
	return []string{
		RoleSuperAdmin,
		RoleTenantAdmin,
		RolePrincipal,
		RoleTeacher,
		RoleStaff,
		RoleParent,
		RoleStudent,
	}
}

// IsSystemRole checks if a role name is a system role.
func IsSystemRole(name string) bool {
	for _, role := range SystemRoles() {
		if role == name {
			return true
		}
	}
	return false
}
