/**
 * MSLS HasPermission Directive
 *
 * Structural directive for conditionally rendering elements based on user permissions.
 *
 * Usage:
 *
 * <!-- Single permission -->
 * <button *mslsHasPermission="'users:write'">Edit User</button>
 *
 * <!-- Multiple permissions (any) -->
 * <button *mslsHasPermission="['users:write', 'users:delete']">Manage User</button>
 *
 * <!-- Multiple permissions (all required) -->
 * <button *mslsHasPermission="['finance:read', 'finance:write']; all: true">Manage Finance</button>
 *
 * <!-- With else template -->
 * <button *mslsHasPermission="'users:write'; else noPermission">Edit</button>
 * <ng-template #noPermission>
 *   <span>No permission to edit</span>
 * </ng-template>
 */

import {
  Directive,
  Input,
  TemplateRef,
  ViewContainerRef,
  inject,
  effect,
  OnDestroy,
} from '@angular/core';

import { RbacService } from '../../core/services/rbac.service';
import { AuthService } from '../../core/services/auth.service';

@Directive({
  selector: '[mslsHasPermission]',
  standalone: true,
})
export class HasPermissionDirective implements OnDestroy {
  private templateRef = inject(TemplateRef<unknown>);
  private viewContainer = inject(ViewContainerRef);
  private rbacService = inject(RbacService);

  /** The permission(s) to check */
  private permissions: string[] = [];

  /** Whether all permissions are required (default: any) */
  private requireAll = false;

  /** The else template to show when permission check fails */
  private elseTemplateRef: TemplateRef<unknown> | null = null;

  /** Whether the view is currently showing */
  private hasView = false;

  /** Whether the else view is currently showing */
  private hasElseView = false;

  /**
   * Set the permission(s) to check
   */
  @Input()
  set mslsHasPermission(value: string | string[]) {
    this.permissions = Array.isArray(value) ? value : [value];
    this.updateView();
  }

  /**
   * Set whether all permissions are required
   */
  @Input()
  set mslsHasPermissionAll(value: boolean) {
    this.requireAll = value;
    this.updateView();
  }

  /**
   * Set the else template
   */
  @Input()
  set mslsHasPermissionElse(templateRef: TemplateRef<unknown> | null) {
    this.elseTemplateRef = templateRef;
    this.updateView();
  }

  constructor() {
    // React to permission changes using effect
    effect(() => {
      // Access the cached permissions signal to create dependency
      this.rbacService.cachedPermissions();
      this.updateView();
    });
  }

  ngOnDestroy(): void {
    this.viewContainer.clear();
  }

  /**
   * Update the view based on permission check
   */
  private updateView(): void {
    const hasPermission = this.checkPermission();

    if (hasPermission) {
      // Show main content
      if (!this.hasView) {
        this.viewContainer.clear();
        this.viewContainer.createEmbeddedView(this.templateRef);
        this.hasView = true;
        this.hasElseView = false;
      }
    } else {
      // Show else content or nothing
      if (this.hasView || !this.hasElseView) {
        this.viewContainer.clear();
        this.hasView = false;

        if (this.elseTemplateRef) {
          this.viewContainer.createEmbeddedView(this.elseTemplateRef);
          this.hasElseView = true;
        } else {
          this.hasElseView = false;
        }
      }
    }
  }

  /**
   * Check if the user has the required permission(s)
   */
  private checkPermission(): boolean {
    if (this.permissions.length === 0) {
      return true;
    }

    if (this.requireAll) {
      return this.rbacService.hasAllPermissions(this.permissions);
    }

    return this.rbacService.hasAnyPermission(this.permissions);
  }
}

/**
 * HasRole Directive
 *
 * Structural directive for conditionally rendering elements based on user roles.
 *
 * Usage:
 *
 * <!-- Single role -->
 * <div *mslsHasRole="'TenantAdmin'">Admin Content</div>
 *
 * <!-- Multiple roles (any) -->
 * <div *mslsHasRole="['TenantAdmin', 'Principal']">Leadership Content</div>
 */
@Directive({
  selector: '[mslsHasRole]',
  standalone: true,
})
export class HasRoleDirective implements OnDestroy {
  private templateRef = inject(TemplateRef<unknown>);
  private viewContainer = inject(ViewContainerRef);
  private rbacService = inject(RbacService);
  private authService = inject(AuthService);

  /** The role(s) to check */
  private roles: string[] = [];

  /** Whether all roles are required (default: any) */
  private requireAll = false;

  /** The else template to show when role check fails */
  private elseTemplateRef: TemplateRef<unknown> | null = null;

  /** Whether the view is currently showing */
  private hasView = false;

  /** Whether the else view is currently showing */
  private hasElseView = false;

  /**
   * Set the role(s) to check
   */
  @Input()
  set mslsHasRole(value: string | string[]) {
    this.roles = Array.isArray(value) ? value : [value];
    this.updateView();
  }

  /**
   * Set whether all roles are required
   */
  @Input()
  set mslsHasRoleAll(value: boolean) {
    this.requireAll = value;
    this.updateView();
  }

  /**
   * Set the else template
   */
  @Input()
  set mslsHasRoleElse(templateRef: TemplateRef<unknown> | null) {
    this.elseTemplateRef = templateRef;
    this.updateView();
  }

  constructor() {
    // React to permission changes using effect
    effect(() => {
      // Access the cached permissions signal to create dependency
      this.rbacService.cachedPermissions();
      this.updateView();
    });
  }

  ngOnDestroy(): void {
    this.viewContainer.clear();
  }

  /**
   * Update the view based on role check
   */
  private updateView(): void {
    const hasRole = this.checkRole();

    if (hasRole) {
      if (!this.hasView) {
        this.viewContainer.clear();
        this.viewContainer.createEmbeddedView(this.templateRef);
        this.hasView = true;
        this.hasElseView = false;
      }
    } else {
      if (this.hasView || !this.hasElseView) {
        this.viewContainer.clear();
        this.hasView = false;

        if (this.elseTemplateRef) {
          this.viewContainer.createEmbeddedView(this.elseTemplateRef);
          this.hasElseView = true;
        } else {
          this.hasElseView = false;
        }
      }
    }
  }

  /**
   * Check if the user has the required role(s)
   */
  private checkRole(): boolean {
    if (this.roles.length === 0) {
      return true;
    }

    if (this.requireAll) {
      return this.roles.every(role => this.authService.hasRole(role));
    }

    return this.roles.some(role => this.authService.hasRole(role));
  }
}
