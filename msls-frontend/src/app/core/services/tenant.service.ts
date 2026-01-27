/**
 * MSLS Tenant Service
 *
 * Manages tenant context for multi-tenant operations.
 * Uses Angular Signals for reactive state management.
 */

import { Injectable, computed, inject, signal } from '@angular/core';
import { StorageService, STORAGE_KEYS } from './storage.service';

/** Tenant status */
export type TenantStatus = 'active' | 'inactive' | 'suspended' | 'pending';

/**
 * Tenant interface representing a tenant/organization
 */
export interface Tenant {
  /** Unique tenant identifier */
  id: string;

  /** Tenant name */
  name: string;

  /** Tenant slug (URL-friendly identifier) */
  slug: string;

  /** Tenant status */
  status: TenantStatus;

  /** Tenant configuration */
  config?: TenantConfig;
}

/**
 * Tenant configuration options
 */
export interface TenantConfig {
  /** Primary brand color */
  primaryColor?: string;

  /** Logo URL */
  logoUrl?: string;

  /** Enabled features */
  features?: string[];

  /** Custom domain */
  customDomain?: string;
}

/**
 * TenantService - Manages tenant context for multi-tenant operations.
 *
 * Usage:
 * constructor(private tenantService: TenantService) {}
 *
 * // Get current tenant ID
 * const tenantId = this.tenantService.currentTenantId();
 *
 * // Set tenant
 * this.tenantService.setTenant(tenant);
 */
@Injectable({ providedIn: 'root' })
export class TenantService {
  private storageService = inject(StorageService);

  /** Internal signal for current tenant */
  private _currentTenant = signal<Tenant | null>(null);

  /** Internal signal for current tenant ID */
  private _currentTenantId = signal<string | null>(null);

  /** Public readonly signal for current tenant */
  readonly currentTenant = this._currentTenant.asReadonly();

  /** Public readonly signal for current tenant ID */
  readonly currentTenantId = this._currentTenantId.asReadonly();

  /** Computed signal indicating if tenant is set */
  readonly hasTenant = computed(() => !!this._currentTenantId());

  /** Computed signal indicating if tenant is active */
  readonly isTenantActive = computed(() => {
    const tenant = this._currentTenant();
    return tenant?.status === 'active';
  });

  constructor() {
    this.initializeFromStorage();
  }

  /**
   * Initialize tenant from storage
   */
  private initializeFromStorage(): void {
    const tenantId = this.storageService.getItem<string>(STORAGE_KEYS.TENANT_ID);
    if (tenantId) {
      this._currentTenantId.set(tenantId);
    }
  }

  /**
   * Set the current tenant
   * @param tenant - Tenant to set as current
   */
  setTenant(tenant: Tenant): void {
    this._currentTenant.set(tenant);
    this._currentTenantId.set(tenant.id);
    this.storageService.setItem(STORAGE_KEYS.TENANT_ID, tenant.id);
  }

  /**
   * Set only the tenant ID (without full tenant details)
   * @param tenantId - Tenant ID to set
   */
  setTenantId(tenantId: string): void {
    this._currentTenantId.set(tenantId);
    this.storageService.setItem(STORAGE_KEYS.TENANT_ID, tenantId);
  }

  /**
   * Clear the current tenant
   */
  clearTenant(): void {
    this._currentTenant.set(null);
    this._currentTenantId.set(null);
    this.storageService.removeItem(STORAGE_KEYS.TENANT_ID);
  }

  /**
   * Check if the current tenant has a specific feature enabled
   * @param feature - Feature name to check
   */
  hasFeature(feature: string): boolean {
    const tenant = this._currentTenant();
    return tenant?.config?.features?.includes(feature) ?? false;
  }

  /**
   * Get tenant configuration value
   * @param key - Configuration key
   */
  getConfig<K extends keyof TenantConfig>(key: K): TenantConfig[K] | undefined {
    return this._currentTenant()?.config?.[key];
  }
}
