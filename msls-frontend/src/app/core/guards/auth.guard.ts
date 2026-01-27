/**
 * MSLS Auth Guard
 *
 * Protects routes that require authentication.
 * Redirects unauthenticated users to the login page.
 */

import { inject } from '@angular/core';
import { Router, CanActivateFn, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

import { AuthService } from '../services/auth.service';

/**
 * Auth Guard Function
 *
 * Checks if the user is authenticated before allowing route access.
 * If not authenticated, redirects to the login page with a return URL.
 *
 * Usage in routes:
 * {
 *   path: 'dashboard',
 *   component: DashboardComponent,
 *   canActivate: [authGuard]
 * }
 */
export const authGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
): boolean => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // Check if user is authenticated
  if (authService.isAuthenticated()) {
    return true;
  }

  // Store the attempted URL for redirecting after login
  const returnUrl = state.url;

  // Redirect to login page with return URL
  router.navigate(['/login'], {
    queryParams: { returnUrl },
  });

  return false;
};

/**
 * Guest Guard Function
 *
 * Protects routes that should only be accessible to non-authenticated users
 * (e.g., login, register pages). Redirects authenticated users to dashboard.
 *
 * Usage in routes:
 * {
 *   path: 'login',
 *   component: LoginComponent,
 *   canActivate: [guestGuard]
 * }
 */
export const guestGuard: CanActivateFn = (): boolean => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // If user is authenticated, redirect to dashboard
  if (authService.isAuthenticated()) {
    router.navigate(['/dashboard']);
    return false;
  }

  return true;
};

/**
 * Permission Guard Factory
 *
 * Creates a guard that checks for specific permissions.
 *
 * Usage:
 * {
 *   path: 'admin/users',
 *   component: UsersComponent,
 *   canActivate: [authGuard, permissionGuard(['users:read'])]
 * }
 */
export function permissionGuard(requiredPermissions: string[]): CanActivateFn {
  return (): boolean => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // First check if authenticated
    if (!authService.isAuthenticated()) {
      router.navigate(['/login']);
      return false;
    }

    // Check if user has all required permissions
    const hasAllPermissions = requiredPermissions.every(permission =>
      authService.hasPermission(permission)
    );

    if (!hasAllPermissions) {
      // Redirect to unauthorized page or dashboard
      router.navigate(['/unauthorized']);
      return false;
    }

    return true;
  };
}

/**
 * Role Guard Factory
 *
 * Creates a guard that checks for specific roles.
 *
 * Usage:
 * {
 *   path: 'admin',
 *   component: AdminComponent,
 *   canActivate: [authGuard, roleGuard(['admin', 'manager'])]
 * }
 */
export function roleGuard(allowedRoles: string[]): CanActivateFn {
  return (): boolean => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // First check if authenticated
    if (!authService.isAuthenticated()) {
      router.navigate(['/login']);
      return false;
    }

    // Check if user has any of the allowed roles
    const hasAllowedRole = allowedRoles.some(role => authService.hasRole(role));

    if (!hasAllowedRole) {
      // Redirect to unauthorized page or dashboard
      router.navigate(['/unauthorized']);
      return false;
    }

    return true;
  };
}
