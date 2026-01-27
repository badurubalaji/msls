/**
 * MSLS User Model
 *
 * Defines the User entity and related types for authentication and authorization.
 */

/** User account status */
export type UserStatus = 'active' | 'inactive' | 'pending' | 'suspended';

/** User role within the system */
export type UserRole = 'admin' | 'manager' | 'instructor' | 'student' | 'guest';

/** System permission identifiers */
export type Permission =
  | 'users:read'
  | 'users:write'
  | 'users:delete'
  | 'courses:read'
  | 'courses:write'
  | 'courses:delete'
  | 'reports:read'
  | 'reports:write'
  | 'settings:read'
  | 'settings:write';

/**
 * User interface representing an authenticated user
 */
export interface User {
  /** Unique user identifier */
  id: string;

  /** Tenant ID for multi-tenancy support */
  tenantId: string;

  /** User's email address */
  email: string;

  /** User's phone number (optional) */
  phone?: string;

  /** User's first name */
  firstName: string;

  /** User's last name */
  lastName: string;

  /** User's assigned roles */
  roles: UserRole[];

  /** User's granted permissions */
  permissions: Permission[];

  /** Account status */
  status: UserStatus;

  /** Account creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;
}

/**
 * Computed full name helper
 */
export function getUserFullName(user: User): string {
  return `${user.firstName} ${user.lastName}`.trim();
}

/**
 * Check if user has a specific permission
 */
export function hasPermission(user: User, permission: Permission): boolean {
  return user.permissions.includes(permission);
}

/**
 * Check if user has a specific role
 */
export function hasRole(user: User, role: UserRole): boolean {
  return user.roles.includes(role);
}
