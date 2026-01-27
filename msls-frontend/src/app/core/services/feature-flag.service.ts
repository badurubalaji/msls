/**
 * MSLS Feature Flag Service
 *
 * Manages feature flags with signal-based reactivity.
 * Fetches flags on app init and provides reactive access to flag states.
 */

import { Injectable, computed, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, catchError, tap, map } from 'rxjs';

import { ApiResponse } from '../models/api-response.model';
import { AuthService } from './auth.service';
import { TenantService } from './tenant.service';
import { environment } from '../../../environments/environment';

/**
 * Feature flag state from the API
 */
export interface FeatureFlagState {
  /** Unique key for the flag */
  key: string;
  /** Human-readable name */
  name: string;
  /** Description of what the flag controls */
  description?: string;
  /** Whether the flag is enabled for the current user/tenant */
  enabled: boolean;
  /** Optional custom JSON value */
  customValue?: unknown;
  /** Source of the flag value: 'default', 'tenant', or 'user' */
  source: 'default' | 'tenant' | 'user';
}

/**
 * Response from the feature flags API
 */
export interface FeatureFlagsResponse {
  flags: FeatureFlagState[];
}

/**
 * Pre-defined feature flag keys for type safety
 */
export const FEATURE_FLAGS = {
  ONLINE_ADMISSIONS: 'online_admissions',
  TRANSPORT_TRACKING: 'transport_tracking',
  AI_INSIGHTS: 'ai_insights',
  PARENT_MESSAGING: 'parent_messaging',
  STUDENT_PORTAL: 'student_portal',
} as const;

export type FeatureFlagKey = (typeof FEATURE_FLAGS)[keyof typeof FEATURE_FLAGS] | string;

/**
 * FeatureFlagService - Manages feature flags with Angular Signals.
 *
 * Features:
 * - Signal-based reactive state
 * - Automatic refresh on auth/tenant changes
 * - Cached flag states for fast access
 * - Type-safe flag key constants
 *
 * Usage:
 * constructor(private featureFlagService: FeatureFlagService) {}
 *
 * // Check if a flag is enabled
 * if (this.featureFlagService.isEnabled('ai_insights')) {
 *   // Show AI features
 * }
 *
 * // In templates with computed signal:
 * @if (featureFlagService.isFeatureEnabled('transport_tracking')()) {
 *   <app-transport-tracking />
 * }
 */
@Injectable({ providedIn: 'root' })
export class FeatureFlagService {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private tenantService = inject(TenantService);

  /** API base URL from environment configuration */
  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}`;

  /** Internal signal for all flags */
  private _flags = signal<Map<string, FeatureFlagState>>(new Map());

  /** Loading state */
  private _loading = signal<boolean>(false);

  /** Error state */
  private _error = signal<string | null>(null);

  /** Whether flags have been loaded at least once */
  private _initialized = signal<boolean>(false);

  /** Public readonly signals */
  readonly flags = this._flags.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();
  readonly initialized = this._initialized.asReadonly();

  /** Computed: all flags as an array */
  readonly allFlags = computed(() => Array.from(this._flags().values()));

  /** Computed: only enabled flags */
  readonly enabledFlags = computed(() => this.allFlags().filter(f => f.enabled));

  constructor() {
    // Note: Initial load should be triggered by the app initialization
    // or when the user authenticates
  }

  /**
   * Load all feature flags for the current user/tenant context.
   * Should be called after authentication and when tenant changes.
   */
  loadFlags(): Observable<FeatureFlagState[]> {
    // Don't load if not authenticated
    if (!this.authService.isAuthenticated()) {
      this._flags.set(new Map());
      this._initialized.set(true);
      return of([]);
    }

    this._loading.set(true);
    this._error.set(null);

    return this.http
      .get<ApiResponse<FeatureFlagsResponse>>(`${this.apiUrl}/feature-flags`)
      .pipe(
        map(response => {
          if (response.success && response.data) {
            return response.data.flags;
          }
          throw new Error('Failed to load feature flags');
        }),
        tap(flags => {
          const flagMap = new Map<string, FeatureFlagState>();
          flags.forEach(flag => flagMap.set(flag.key, flag));
          this._flags.set(flagMap);
          this._initialized.set(true);
          this._loading.set(false);
        }),
        catchError(error => {
          console.error('Failed to load feature flags:', error);
          this._error.set(error.message || 'Failed to load feature flags');
          this._loading.set(false);
          this._initialized.set(true);
          return of([]);
        })
      );
  }

  /**
   * Check if a feature flag is enabled (synchronous).
   * Returns false if flags haven't been loaded or flag doesn't exist.
   *
   * @param flagKey - The flag key to check
   * @returns boolean indicating if the flag is enabled
   */
  isEnabled(flagKey: FeatureFlagKey): boolean {
    const flag = this._flags().get(flagKey);
    return flag?.enabled ?? false;
  }

  /**
   * Get a computed signal for a specific feature flag's enabled state.
   * Useful for reactive template bindings.
   *
   * @param flagKey - The flag key to track
   * @returns Computed signal that updates when flags change
   */
  isFeatureEnabled(flagKey: FeatureFlagKey) {
    return computed(() => {
      const flag = this._flags().get(flagKey);
      return flag?.enabled ?? false;
    });
  }

  /**
   * Get the full state of a feature flag.
   *
   * @param flagKey - The flag key
   * @returns The flag state or undefined if not found
   */
  getFlag(flagKey: FeatureFlagKey): FeatureFlagState | undefined {
    return this._flags().get(flagKey);
  }

  /**
   * Get a computed signal for a specific feature flag.
   *
   * @param flagKey - The flag key to track
   * @returns Computed signal with the full flag state
   */
  getFlagSignal(flagKey: FeatureFlagKey) {
    return computed(() => this._flags().get(flagKey));
  }

  /**
   * Get the custom value for a feature flag.
   *
   * @param flagKey - The flag key
   * @returns The custom value or undefined
   */
  getCustomValue<T = unknown>(flagKey: FeatureFlagKey): T | undefined {
    const flag = this._flags().get(flagKey);
    return flag?.customValue as T | undefined;
  }

  /**
   * Check if multiple flags are all enabled.
   *
   * @param flagKeys - Array of flag keys to check
   * @returns true if all flags are enabled
   */
  areAllEnabled(flagKeys: FeatureFlagKey[]): boolean {
    return flagKeys.every(key => this.isEnabled(key));
  }

  /**
   * Check if any of the specified flags are enabled.
   *
   * @param flagKeys - Array of flag keys to check
   * @returns true if at least one flag is enabled
   */
  isAnyEnabled(flagKeys: FeatureFlagKey[]): boolean {
    return flagKeys.some(key => this.isEnabled(key));
  }

  /**
   * Get a computed signal that checks if all specified flags are enabled.
   *
   * @param flagKeys - Array of flag keys to track
   * @returns Computed signal
   */
  areAllEnabledSignal(flagKeys: FeatureFlagKey[]) {
    return computed(() => flagKeys.every(key => this._flags().get(key)?.enabled ?? false));
  }

  /**
   * Get a computed signal that checks if any specified flag is enabled.
   *
   * @param flagKeys - Array of flag keys to track
   * @returns Computed signal
   */
  isAnyEnabledSignal(flagKeys: FeatureFlagKey[]) {
    return computed(() => flagKeys.some(key => this._flags().get(key)?.enabled ?? false));
  }

  /**
   * Refresh feature flags from the server.
   * Call this when tenant or user context changes.
   */
  refresh(): Observable<FeatureFlagState[]> {
    return this.loadFlags();
  }

  /**
   * Clear all cached flags.
   * Called on logout.
   */
  clearFlags(): void {
    this._flags.set(new Map());
    this._initialized.set(false);
    this._error.set(null);
  }

  /**
   * Check a specific flag by making an API call.
   * Useful for critical features where you need real-time accuracy.
   *
   * @param flagKey - The flag key to check
   * @returns Observable with the flag state
   */
  checkFlag(flagKey: FeatureFlagKey): Observable<FeatureFlagState | null> {
    if (!this.authService.isAuthenticated()) {
      return of(null);
    }

    return this.http
      .get<ApiResponse<FeatureFlagState>>(`${this.apiUrl}/feature-flags/${flagKey}`)
      .pipe(
        map(response => {
          if (response.success && response.data) {
            // Update local cache
            const currentFlags = new Map(this._flags());
            currentFlags.set(flagKey, response.data);
            this._flags.set(currentFlags);
            return response.data;
          }
          return null;
        }),
        catchError(() => of(null))
      );
  }
}
