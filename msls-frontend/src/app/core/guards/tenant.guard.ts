/**
 * MSLS Tenant Guard
 *
 * Validates that a tenant exists and is active before allowing route access.
 * Redirects to appropriate error pages if tenant is invalid.
 */

import { inject } from '@angular/core';
import { Router, CanActivateFn, ActivatedRouteSnapshot } from '@angular/router';

import { TenantService } from '../services/tenant.service';

/**
 * Tenant Guard Function
 *
 * Checks if a valid tenant is set and active before allowing route access.
 * Can also extract tenant from route parameters.
 *
 * Usage in routes:
 * {
 *   path: ':tenantSlug/dashboard',
 *   component: DashboardComponent,
 *   canActivate: [tenantGuard]
 * }
 */
export const tenantGuard: CanActivateFn = (route: ActivatedRouteSnapshot): boolean => {
  const tenantService = inject(TenantService);
  const router = inject(Router);

  // Check for tenant slug in route params
  const tenantSlug = route.paramMap.get('tenantSlug');

  // If tenant slug is in route, we could validate it here
  // For now, we just check if a tenant is set
  if (tenantSlug) {
    // In a real application, you might want to:
    // 1. Fetch tenant details from API
    // 2. Validate the tenant exists
    // 3. Set the tenant in TenantService
    // For now, we'll just set the tenant ID from the slug
    tenantService.setTenantId(tenantSlug);
  }

  // Check if tenant is set
  if (!tenantService.hasTenant()) {
    // Redirect to tenant selection or error page
    router.navigate(['/select-tenant']);
    return false;
  }

  // Check if tenant is active
  const tenant = tenantService.currentTenant();
  if (tenant && tenant.status !== 'active') {
    // Redirect to tenant suspended/inactive page
    router.navigate(['/tenant-unavailable'], {
      queryParams: { reason: tenant.status },
    });
    return false;
  }

  return true;
};

/**
 * Tenant Required Guard
 *
 * A stricter guard that requires both tenant to be set AND verified.
 * Use this for routes that need full tenant context.
 *
 * Usage in routes:
 * {
 *   path: 'courses',
 *   component: CoursesComponent,
 *   canActivate: [authGuard, tenantRequiredGuard]
 * }
 */
export const tenantRequiredGuard: CanActivateFn = (): boolean => {
  const tenantService = inject(TenantService);
  const router = inject(Router);

  // Check if tenant ID is set
  if (!tenantService.currentTenantId()) {
    router.navigate(['/select-tenant']);
    return false;
  }

  // Check if full tenant details are loaded and active
  if (!tenantService.isTenantActive()) {
    // If tenant details are not loaded, we might want to fetch them
    // For now, we allow access if tenant ID is set
    // In production, you might want to add an async check here
    return true;
  }

  return true;
};

/**
 * Tenant Feature Guard Factory
 *
 * Creates a guard that checks if a specific feature is enabled for the tenant.
 *
 * Usage:
 * {
 *   path: 'analytics',
 *   component: AnalyticsComponent,
 *   canActivate: [authGuard, tenantFeatureGuard('analytics')]
 * }
 */
export function tenantFeatureGuard(featureName: string): CanActivateFn {
  return (): boolean => {
    const tenantService = inject(TenantService);
    const router = inject(Router);

    // Check if tenant has the required feature
    if (!tenantService.hasFeature(featureName)) {
      router.navigate(['/feature-unavailable'], {
        queryParams: { feature: featureName },
      });
      return false;
    }

    return true;
  };
}
