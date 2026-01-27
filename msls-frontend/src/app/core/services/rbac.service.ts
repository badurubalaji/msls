/**
 * MSLS RBAC Service
 *
 * Manages role-based access control operations including:
 * - Role CRUD operations
 * - Permission management
 * - User role assignments
 * - Permission checking with caching
 */

import { Injectable, computed, effect, inject, signal } from '@angular/core';
import { Observable, tap, shareReplay, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

import { ApiService } from './api.service';
import { AuthService } from './auth.service';
import {
  Role,
  RbacPermission,
  UserRoles,
  PermissionModulesResponse,
  CreateRoleRequest,
  UpdateRoleRequest,
  PermissionsRequest,
  RolesRequest,
} from '../models';

/** Cache duration in milliseconds (5 minutes) */
const CACHE_DURATION = 5 * 60 * 1000;

/**
 * Backend API response interfaces (snake_case format)
 */

/** Permission as returned by the backend API */
interface BackendPermission {
  id: string;
  code: string;
  name: string;
  module: string;
  description?: string;
  created_at: string;
}

/** Role as returned by the backend API */
interface BackendRole {
  id: string;
  tenant_id?: string;
  name: string;
  description?: string;
  is_system: boolean;
  permissions?: BackendPermission[];
  created_at: string;
  updated_at: string;
}

/**
 * RbacService - Manages role-based access control.
 *
 * Features:
 * - Signal-based reactive state
 * - Permission caching for performance
 * - Role and permission management
 * - User role assignment
 *
 * Usage:
 * constructor(private rbacService: RbacService) {}
 *
 * // Check permission
 * if (rbacService.hasPermission('users:write')) {
 *   // Show edit button
 * }
 *
 * // In templates with signals
 * @if (rbacService.can('users:write')) {
 *   <button>Edit</button>
 * }
 */
@Injectable({ providedIn: 'root' })
export class RbacService {
  private apiService = inject(ApiService);
  private authService = inject(AuthService);

  /** Internal signal for cached permissions */
  private _cachedPermissions = signal<Set<string>>(new Set());

  /** Internal signal for cache timestamp */
  private _cacheTimestamp = signal<number>(0);

  /** Internal signal for all permissions (for admin UI) */
  private _allPermissions = signal<RbacPermission[]>([]);

  /** Internal signal for all roles (for admin UI) */
  private _allRoles = signal<Role[]>([]);

  /** Internal signal for permission modules */
  private _permissionModules = signal<string[]>([]);

  /** Observable cache for permissions list */
  private permissionsCache$: Observable<RbacPermission[]> | null = null;

  /** Observable cache for roles list */
  private rolesCache$: Observable<Role[]> | null = null;

  /** Public readonly signal for cached permissions */
  readonly cachedPermissions = this._cachedPermissions.asReadonly();

  /** Public readonly signal for all permissions */
  readonly allPermissions = this._allPermissions.asReadonly();

  /** Public readonly signal for all roles */
  readonly allRoles = this._allRoles.asReadonly();

  /** Public readonly signal for permission modules */
  readonly permissionModules = this._permissionModules.asReadonly();

  /** Computed signal: whether current user can manage roles */
  readonly canManageRoles = computed(() => this.hasPermission('roles:write'));

  /** Computed signal: whether current user can view roles */
  readonly canViewRoles = computed(() => this.hasPermission('roles:read'));

  /** Computed signal: whether current user can delete roles */
  readonly canDeleteRoles = computed(() => this.hasPermission('roles:delete'));

  constructor() {
    // Initialize permissions from current user
    this.initializeFromUser();

    // Re-initialize when user changes (e.g., after login)
    effect(() => {
      const user = this.authService.currentUser();
      if (user?.permissions) {
        this._cachedPermissions.set(new Set(user.permissions as string[]));
        this._cacheTimestamp.set(Date.now());
      } else {
        this._cachedPermissions.set(new Set());
      }
    });
  }

  /**
   * Initialize cached permissions from the current user
   */
  private initializeFromUser(): void {
    const user = this.authService.currentUser();
    if (user?.permissions) {
      this._cachedPermissions.set(new Set(user.permissions as string[]));
      this._cacheTimestamp.set(Date.now());
    }
  }

  /**
   * Check if the current user has a specific permission
   * @param permission - Permission code (e.g., 'users:write')
   */
  hasPermission(permission: string): boolean {
    // First check cached permissions
    if (this._cachedPermissions().has(permission)) {
      return true;
    }

    // Fallback to auth service
    return this.authService.hasPermission(permission);
  }

  /**
   * Signal-based permission check for templates
   * @param permission - Permission code (e.g., 'users:write')
   */
  can(permission: string): boolean {
    return this.hasPermission(permission);
  }

  /**
   * Check if the current user has any of the specified permissions
   * @param permissions - Array of permission codes
   */
  hasAnyPermission(permissions: string[]): boolean {
    return permissions.some(p => this.hasPermission(p));
  }

  /**
   * Check if the current user has all of the specified permissions
   * @param permissions - Array of permission codes
   */
  hasAllPermissions(permissions: string[]): boolean {
    return permissions.every(p => this.hasPermission(p));
  }

  /**
   * Refresh the cached permissions from the server
   */
  refreshPermissions(): Observable<string[]> {
    return this.apiService.get<UserRoles>('/users/me/roles').pipe(
      map(response => {
        const permissions = new Set<string>();
        response.roles.forEach(role => {
          role.permissions?.forEach(perm => {
            permissions.add(perm.code);
          });
        });
        this._cachedPermissions.set(permissions);
        this._cacheTimestamp.set(Date.now());
        return Array.from(permissions);
      })
    );
  }

  /**
   * Clear the permission cache
   */
  clearCache(): void {
    this._cachedPermissions.set(new Set());
    this._cacheTimestamp.set(0);
    this.permissionsCache$ = null;
    this.rolesCache$ = null;
  }

  // =========================================================================
  // Role Management
  // =========================================================================

  /**
   * Get all roles for the current tenant
   * @param includeSystem - Whether to include system roles
   * @param search - Optional search term
   */
  getRoles(includeSystem: boolean = true, search?: string): Observable<Role[]> {
    const params: Record<string, string | boolean> = { include_system: includeSystem };
    if (search) {
      params['search'] = search;
    }

    return this.apiService.get<BackendRole[]>('/roles', { params }).pipe(
      map(rawRoles => rawRoles.map(r => this.mapRoleResponse(r))),
      tap(roles => this._allRoles.set(roles))
    );
  }

  /**
   * Map backend role response (snake_case) to frontend model (camelCase)
   */
  private mapRoleResponse(data: BackendRole): Role {
    // Map permissions from backend format
    const permissions: RbacPermission[] = data.permissions
      ? data.permissions.map(p => this.mapPermissionResponse(p))
      : [];

    return {
      id: data.id,
      tenantId: data.tenant_id,
      name: data.name,
      description: data.description,
      isSystem: data.is_system,
      permissions,
      createdAt: data.created_at,
      updatedAt: data.updated_at,
    };
  }

  /**
   * Map backend permission response (snake_case) to frontend model (camelCase)
   */
  private mapPermissionResponse(data: BackendPermission): RbacPermission {
    return {
      id: data.id,
      code: data.code,
      name: data.name,
      module: data.module,
      description: data.description,
      createdAt: data.created_at,
    };
  }

  /**
   * Get roles with caching
   */
  getCachedRoles(): Observable<Role[]> {
    if (!this.rolesCache$) {
      this.rolesCache$ = this.getRoles().pipe(
        shareReplay({ bufferSize: 1, refCount: true })
      );
    }
    return this.rolesCache$;
  }

  /**
   * Get a role by ID
   * @param id - Role ID
   */
  getRole(id: string): Observable<Role> {
    return this.apiService.get<BackendRole>(`/v1/roles/${id}`).pipe(
      map(raw => this.mapRoleResponse(raw))
    );
  }

  /**
   * Create a new role
   * @param request - Role creation request
   */
  createRole(request: CreateRoleRequest): Observable<Role> {
    return this.apiService.post<BackendRole>('/roles', {
      name: request.name,
      description: request.description,
      permission_ids: request.permissionIds,
    }).pipe(
      map(raw => this.mapRoleResponse(raw)),
      tap(() => {
        // Invalidate cache
        this.rolesCache$ = null;
      })
    );
  }

  /**
   * Update a role
   * @param id - Role ID
   * @param request - Role update request
   */
  updateRole(id: string, request: UpdateRoleRequest): Observable<Role> {
    return this.apiService.put<BackendRole>(`/v1/roles/${id}`, request).pipe(
      map(raw => this.mapRoleResponse(raw)),
      tap(() => {
        // Invalidate cache
        this.rolesCache$ = null;
      })
    );
  }

  /**
   * Delete a role
   * @param id - Role ID
   */
  deleteRole(id: string): Observable<void> {
    return this.apiService.delete<void>(`/v1/roles/${id}`).pipe(
      tap(() => {
        // Invalidate cache
        this.rolesCache$ = null;
      })
    );
  }

  /**
   * Assign permissions to a role
   * @param roleId - Role ID
   * @param request - Permissions to assign
   */
  assignPermissions(roleId: string, request: PermissionsRequest): Observable<Role> {
    return this.apiService.post<BackendRole>(`/v1/roles/${roleId}/permissions`, {
      permission_ids: request.permissionIds,
    }).pipe(
      map(raw => this.mapRoleResponse(raw)),
      tap(() => {
        // Invalidate cache
        this.rolesCache$ = null;
      })
    );
  }

  /**
   * Remove permissions from a role
   * @param roleId - Role ID
   * @param request - Permissions to remove
   */
  removePermissions(roleId: string, request: PermissionsRequest): Observable<Role> {
    return this.apiService.delete<BackendRole>(`/v1/roles/${roleId}/permissions`, {
      params: { permission_ids: request.permissionIds },
    }).pipe(
      map(raw => this.mapRoleResponse(raw)),
      tap(() => {
        // Invalidate cache
        this.rolesCache$ = null;
      })
    );
  }

  // =========================================================================
  // Permission Management
  // =========================================================================

  /**
   * Get all permissions
   * @param module - Optional module filter
   * @param search - Optional search term
   */
  getPermissions(module?: string, search?: string): Observable<RbacPermission[]> {
    const params: Record<string, string> = {};
    if (module) params['module'] = module;
    if (search) params['search'] = search;

    return this.apiService.get<BackendPermission[]>('/permissions', { params }).pipe(
      map(rawPerms => rawPerms.map(p => this.mapPermissionResponse(p))),
      tap(permissions => this._allPermissions.set(permissions))
    );
  }

  /**
   * Get permissions with caching
   */
  getCachedPermissions(): Observable<RbacPermission[]> {
    if (!this.permissionsCache$) {
      this.permissionsCache$ = this.getPermissions().pipe(
        shareReplay({ bufferSize: 1, refCount: true })
      );
    }
    return this.permissionsCache$;
  }

  /**
   * Get all permission modules
   */
  getPermissionModules(): Observable<string[]> {
    return this.apiService.get<PermissionModulesResponse>('/permissions/modules').pipe(
      map(response => response.modules),
      tap(modules => this._permissionModules.set(modules))
    );
  }

  /**
   * Get permissions by module
   * @param module - Module name
   */
  getPermissionsByModule(module: string): Observable<RbacPermission[]> {
    return this.apiService.get<BackendPermission[]>(`/v1/permissions/modules/${module}`).pipe(
      map(rawPerms => rawPerms.map(p => this.mapPermissionResponse(p)))
    );
  }

  // =========================================================================
  // User Role Management
  // =========================================================================

  /**
   * Get current user's roles
   */
  getMyRoles(): Observable<UserRoles> {
    return this.apiService.get<UserRoles>('/users/me/roles');
  }

  /**
   * Get a user's roles
   * @param userId - User ID
   */
  getUserRoles(userId: string): Observable<UserRoles> {
    return this.apiService.get<UserRoles>(`/v1/users/${userId}/roles`);
  }

  /**
   * Assign roles to a user
   * @param userId - User ID
   * @param request - Roles to assign
   */
  assignUserRoles(userId: string, request: RolesRequest): Observable<UserRoles> {
    return this.apiService.post<UserRoles>(`/v1/users/${userId}/roles`, {
      role_ids: request.roleIds,
    });
  }

  /**
   * Remove roles from a user
   * @param userId - User ID
   * @param request - Roles to remove
   */
  removeUserRoles(userId: string, request: RolesRequest): Observable<UserRoles> {
    return this.apiService.delete<UserRoles>(`/v1/users/${userId}/roles`, {
      params: { role_ids: request.roleIds },
    });
  }

  // =========================================================================
  // Utility Methods
  // =========================================================================

  /**
   * Group permissions by module
   * @param permissions - Array of permissions
   */
  groupPermissionsByModule(permissions: RbacPermission[]): Map<string, RbacPermission[]> {
    const grouped = new Map<string, RbacPermission[]>();
    permissions.forEach(perm => {
      const existing = grouped.get(perm.module) || [];
      existing.push(perm);
      grouped.set(perm.module, existing);
    });
    return grouped;
  }

  /**
   * Check if a permission code is valid format
   * @param code - Permission code
   */
  isValidPermissionCode(code: string): boolean {
    return /^[a-z_]+:[a-z_]+$/.test(code);
  }
}
