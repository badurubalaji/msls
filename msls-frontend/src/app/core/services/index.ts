/**
 * MSLS Core Services - Barrel Export
 *
 * This file exports all core services.
 * Import from '@core/services' in your modules.
 */

// Layout Service
export { LayoutService } from './layout.service';

// Storage Service
export { StorageService, STORAGE_KEYS } from './storage.service';
export type { StorageKey } from './storage.service';

// Loading Service
export { LoadingService } from './loading.service';

// Tenant Service
export { TenantService } from './tenant.service';
export type { Tenant, TenantConfig, TenantStatus } from './tenant.service';

// Auth Service
export { AuthService } from './auth.service';

// API Service
export { ApiService } from './api.service';
export type { RequestOptions } from './api.service';

// Profile Service
export { ProfileService } from './profile.service';
export type { UploadProgress } from './profile.service';

// RBAC Service
export { RbacService } from './rbac.service';

// Feature Flag Service
export { FeatureFlagService, FEATURE_FLAGS } from './feature-flag.service';
export type { FeatureFlagState, FeatureFlagKey, FeatureFlagsResponse } from './feature-flag.service';
