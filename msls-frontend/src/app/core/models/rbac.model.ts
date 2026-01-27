/**
 * MSLS RBAC Models
 *
 * Defines Role-Based Access Control entities for the frontend.
 */

/**
 * Permission entity
 */
export interface Permission {
  /** Unique permission identifier */
  id: string;

  /** Permission code in format "module:action" (e.g., "users:read") */
  code: string;

  /** Human-readable permission name */
  name: string;

  /** Permission module/category */
  module: string;

  /** Description of what this permission allows */
  description?: string;

  /** Creation timestamp */
  createdAt: string;
}

/**
 * Role entity
 */
export interface Role {
  /** Unique role identifier */
  id: string;

  /** Tenant ID (null for system roles) */
  tenantId?: string;

  /** Role name */
  name: string;

  /** Role description */
  description?: string;

  /** Whether this is a system-defined role */
  isSystem: boolean;

  /** Permissions assigned to this role */
  permissions?: Permission[];

  /** Creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;
}

/**
 * User roles response
 */
export interface UserRoles {
  /** User ID */
  userId: string;

  /** Roles assigned to the user */
  roles: Role[];
}

/**
 * Permission modules response
 */
export interface PermissionModules {
  /** List of permission module names */
  modules: string[];
}

/**
 * Create role request
 */
export interface CreateRoleRequest {
  /** Role name */
  name: string;

  /** Role description */
  description?: string;

  /** Permission IDs to assign */
  permissionIds?: string[];
}

/**
 * Update role request
 */
export interface UpdateRoleRequest {
  /** New role name */
  name?: string;

  /** New role description */
  description?: string;
}

/**
 * Assign/remove permissions request
 */
export interface PermissionsRequest {
  /** Permission IDs */
  permissionIds: string[];
}

/**
 * Assign/remove roles request
 */
export interface RolesRequest {
  /** Role IDs */
  roleIds: string[];
}

/**
 * System role names
 */
export const SystemRoles = {
  SUPER_ADMIN: 'SuperAdmin',
  TENANT_ADMIN: 'TenantAdmin',
  PRINCIPAL: 'Principal',
  TEACHER: 'Teacher',
  STAFF: 'Staff',
  PARENT: 'Parent',
  STUDENT: 'Student',
} as const;

export type SystemRoleName = (typeof SystemRoles)[keyof typeof SystemRoles];

/**
 * Permission modules
 */
export const PermissionModules = {
  USERS: 'users',
  STUDENTS: 'students',
  STAFF: 'staff',
  FINANCE: 'finance',
  ACADEMICS: 'academics',
  ROLES: 'roles',
  SETTINGS: 'settings',
  REPORTS: 'reports',
} as const;

export type PermissionModule = (typeof PermissionModules)[keyof typeof PermissionModules];

/**
 * Permission actions
 */
export const PermissionActions = {
  READ: 'read',
  WRITE: 'write',
  DELETE: 'delete',
  MANAGE: 'manage',
} as const;

export type PermissionAction = (typeof PermissionActions)[keyof typeof PermissionActions];

/**
 * Build a permission code
 */
export function buildPermissionCode(module: PermissionModule, action: PermissionAction): string {
  return `${module}:${action}`;
}

/**
 * Parse a permission code into module and action
 */
export function parsePermissionCode(code: string): { module: string; action: string } | null {
  const parts = code.split(':');
  if (parts.length !== 2) return null;
  return { module: parts[0], action: parts[1] };
}

/**
 * Check if a role is a system role
 */
export function isSystemRole(role: Role): boolean {
  return role.isSystem;
}

/**
 * Check if a role can be modified
 */
export function canModifyRole(role: Role): boolean {
  return !role.isSystem;
}

/**
 * Check if a role can be deleted
 */
export function canDeleteRole(role: Role): boolean {
  return !role.isSystem;
}
