/**
 * MSLS Feature Flag Guard
 *
 * Protects routes based on feature flags.
 * Redirects to dashboard or error page if required flag is disabled.
 */

import { inject } from '@angular/core';
import { Router, CanActivateFn, ActivatedRouteSnapshot } from '@angular/router';

import { FeatureFlagService, FeatureFlagKey } from '../services/feature-flag.service';

/**
 * Feature Flag Guard Factory
 *
 * Creates a guard that checks if a specific feature flag is enabled.
 * If the flag is disabled, redirects to the dashboard or a custom URL.
 *
 * Usage in routes:
 * ```typescript
 * {
 *   path: 'transport',
 *   component: TransportComponent,
 *   canActivate: [authGuard, featureFlagGuard('transport_tracking')]
 * }
 *
 * // With custom redirect
 * {
 *   path: 'ai-insights',
 *   component: AiInsightsComponent,
 *   canActivate: [authGuard, featureFlagGuard('ai_insights', '/upgrade')]
 * }
 * ```
 */
export function featureFlagGuard(
  requiredFlag: FeatureFlagKey,
  redirectUrl: string = '/dashboard'
): CanActivateFn {
  return (): boolean => {
    const featureFlagService = inject(FeatureFlagService);
    const router = inject(Router);

    // Check if the flag is enabled
    if (featureFlagService.isEnabled(requiredFlag)) {
      return true;
    }

    // Redirect to the specified URL
    router.navigate([redirectUrl], {
      queryParams: { feature: requiredFlag, reason: 'disabled' },
    });
    return false;
  };
}

/**
 * Feature Flag Any Guard Factory
 *
 * Creates a guard that allows access if ANY of the specified flags are enabled.
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'communication',
 *   component: CommunicationComponent,
 *   canActivate: [authGuard, featureFlagAnyGuard(['parent_messaging', 'student_portal'])]
 * }
 * ```
 */
export function featureFlagAnyGuard(
  requiredFlags: FeatureFlagKey[],
  redirectUrl: string = '/dashboard'
): CanActivateFn {
  return (): boolean => {
    const featureFlagService = inject(FeatureFlagService);
    const router = inject(Router);

    // Check if any flag is enabled
    if (featureFlagService.isAnyEnabled(requiredFlags)) {
      return true;
    }

    // Redirect if no flags are enabled
    router.navigate([redirectUrl], {
      queryParams: { features: requiredFlags.join(','), reason: 'all_disabled' },
    });
    return false;
  };
}

/**
 * Feature Flag All Guard Factory
 *
 * Creates a guard that allows access only if ALL specified flags are enabled.
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'advanced-analytics',
 *   component: AdvancedAnalyticsComponent,
 *   canActivate: [authGuard, featureFlagAllGuard(['ai_insights', 'student_portal'])]
 * }
 * ```
 */
export function featureFlagAllGuard(
  requiredFlags: FeatureFlagKey[],
  redirectUrl: string = '/dashboard'
): CanActivateFn {
  return (): boolean => {
    const featureFlagService = inject(FeatureFlagService);
    const router = inject(Router);

    // Check if all flags are enabled
    if (featureFlagService.areAllEnabled(requiredFlags)) {
      return true;
    }

    // Find which flags are disabled
    const disabledFlags = requiredFlags.filter(flag => !featureFlagService.isEnabled(flag));

    // Redirect if any flag is disabled
    router.navigate([redirectUrl], {
      queryParams: { features: disabledFlags.join(','), reason: 'some_disabled' },
    });
    return false;
  };
}

/**
 * Feature Flag From Route Guard
 *
 * A guard that reads the required flag from route data.
 * Useful when you want to configure flags in the route definition.
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'feature',
 *   component: FeatureComponent,
 *   canActivate: [authGuard, featureFlagFromRouteGuard],
 *   data: { featureFlag: 'my_feature', featureFlagRedirect: '/not-available' }
 * }
 * ```
 */
export const featureFlagFromRouteGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot
): boolean => {
  const featureFlagService = inject(FeatureFlagService);
  const router = inject(Router);

  // Get flag from route data
  const requiredFlag = route.data['featureFlag'] as FeatureFlagKey;
  const redirectUrl = (route.data['featureFlagRedirect'] as string) || '/dashboard';

  if (!requiredFlag) {
    console.warn('featureFlagFromRouteGuard: No featureFlag specified in route data');
    return true;
  }

  // Check if the flag is enabled
  if (featureFlagService.isEnabled(requiredFlag)) {
    return true;
  }

  // Redirect
  router.navigate([redirectUrl], {
    queryParams: { feature: requiredFlag, reason: 'disabled' },
  });
  return false;
};

/**
 * Feature Unavailable Guard
 *
 * Guard that always denies access and redirects.
 * Useful for marking routes as completely unavailable.
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'deprecated-feature',
 *   component: DeprecatedComponent,
 *   canActivate: [featureUnavailableGuard('/new-feature')]
 * }
 * ```
 */
export function featureUnavailableGuard(redirectUrl: string = '/dashboard'): CanActivateFn {
  return (): boolean => {
    const router = inject(Router);
    router.navigate([redirectUrl], {
      queryParams: { reason: 'unavailable' },
    });
    return false;
  };
}
