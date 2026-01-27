/**
 * MSLS Shared Directives - Barrel Export
 *
 * This file exports all shared directives.
 * Import from '@shared/directives' in your modules.
 */

// Feature Flag Directive
export {
  FeatureFlagDirective,
  FeatureFlagEnabledDirective,
} from './feature-flag.directive';
export type { FeatureFlagContext } from './feature-flag.directive';

// Permission Directives
export { HasPermissionDirective, HasRoleDirective } from './has-permission.directive';
