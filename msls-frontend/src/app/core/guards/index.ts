/**
 * MSLS Core Guards - Barrel Export
 *
 * This file exports all route guards from the core/guards directory.
 * Guards are functional (CanActivateFn) following Angular 17+ patterns.
 */

// Auth Guards
export { authGuard, guestGuard, permissionGuard, roleGuard } from './auth.guard';

// Tenant Guards
export { tenantGuard, tenantRequiredGuard, tenantFeatureGuard } from './tenant.guard';

// Feature Flag Guards
export {
  featureFlagGuard,
  featureFlagAnyGuard,
  featureFlagAllGuard,
  featureFlagFromRouteGuard,
  featureUnavailableGuard,
} from './feature-flag.guard';
