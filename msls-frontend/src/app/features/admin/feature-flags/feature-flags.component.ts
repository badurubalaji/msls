/**
 * MSLS Admin Feature Flags Component
 *
 * Administrative interface for managing feature flags.
 * Allows viewing, enabling/disabling flags for tenants.
 */

import { Component, ChangeDetectionStrategy, inject, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

import { ApiResponse } from '../../../core/models/api-response.model';
import { environment } from '../../../../environments/environment';

/**
 * Feature flag from admin API
 */
interface FeatureFlag {
  id: string;
  key: string;
  name: string;
  description?: string;
  default_value: boolean;
  metadata: {
    category?: string;
    requires_setup?: boolean;
    beta?: boolean;
    rollout_percentage?: number;
  };
  created_at: string;
  updated_at: string;
}

/**
 * Tenant feature flag override
 */
interface TenantFeatureFlag {
  id: string;
  tenant_id: string;
  flag_key: string;
  flag_name: string;
  enabled: boolean;
  custom_value?: unknown;
  created_at: string;
  updated_at: string;
}

/**
 * List feature flags response
 */
interface ListFlagsResponse {
  flags: FeatureFlag[];
}

/**
 * AdminFeatureFlagsComponent - Manage feature flags.
 */
@Component({
  selector: 'msls-admin-feature-flags',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="feature-flags-page">
      <!-- Header -->
      <div class="page-header">
        <h1 class="page-title">Feature Flags</h1>
        <p class="page-subtitle">Manage feature flags for gradual feature rollout</p>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <span>Loading feature flags...</span>
        </div>
      }

      <!-- Error State -->
      @if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button (click)="loadFlags()" class="retry-btn">Try again</button>
        </div>
      }

      <!-- Feature Flags Table -->
      @if (!loading() && !error() && flags().length > 0) {
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Flag</th>
                <th>Key</th>
                <th>Category</th>
                <th>Default</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              @for (flag of flags(); track flag.id) {
                <tr>
                  <td>
                    <div class="flag-info">
                      <span class="flag-name">
                        {{ flag.name }}
                        @if (flag.metadata.beta) {
                          <span class="badge badge-purple">Beta</span>
                        }
                      </span>
                      @if (flag.description) {
                        <span class="flag-description">{{ flag.description }}</span>
                      }
                    </div>
                  </td>
                  <td>
                    <code class="flag-key">{{ flag.key }}</code>
                  </td>
                  <td class="category-cell">{{ flag.metadata.category || '-' }}</td>
                  <td>
                    @if (flag.default_value) {
                      <span class="badge badge-green">Enabled</span>
                    } @else {
                      <span class="badge badge-gray">Disabled</span>
                    }
                  </td>
                  <td>
                    @if (flag.metadata.requires_setup) {
                      <span class="badge badge-yellow">Requires Setup</span>
                    } @else {
                      <span class="badge badge-blue">Ready</span>
                    }
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      }

      <!-- Empty State -->
      @if (!loading() && !error() && flags().length === 0) {
        <div class="empty-container">
          <i class="fa-regular fa-flag"></i>
          <h3>No feature flags</h3>
          <p>No feature flags have been configured yet.</p>
        </div>
      }

      <!-- Help Section -->
      <div class="help-section">
        <h3><i class="fa-solid fa-circle-info"></i> About Feature Flags</h3>
        <ul>
          <li>Feature flags allow gradual feature rollout and A/B testing.</li>
          <li>Default values apply to all tenants unless overridden.</li>
          <li>User-level overrides (for beta testing) take highest priority.</li>
          <li>Contact your administrator to enable features for your organization.</li>
        </ul>
      </div>
    </div>
  `,
  styles: [`
    .feature-flags-page {
      padding: 1.5rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .page-header {
      margin-bottom: 1.5rem;
    }

    .page-title {
      font-size: 1.75rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .page-subtitle {
      color: #64748b;
      font-size: 0.9rem;
      margin: 0;
    }

    /* Loading */
    .loading-container {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 3rem;
      gap: 0.75rem;
      color: #64748b;
    }

    .spinner {
      width: 1.5rem;
      height: 1.5rem;
      border: 2px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Error */
    .error-container {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem;
      background: #fef2f2;
      border: 1px solid #fecaca;
      border-radius: 0.5rem;
      color: #dc2626;
      margin-bottom: 1.5rem;
    }

    .retry-btn {
      margin-left: auto;
      padding: 0.375rem 0.75rem;
      font-size: 0.875rem;
      color: #dc2626;
      background: transparent;
      border: 1px solid #dc2626;
      border-radius: 0.375rem;
      cursor: pointer;
    }

    .retry-btn:hover {
      background: #dc2626;
      color: white;
    }

    /* Table */
    .table-container {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table thead {
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table th {
      padding: 0.75rem 1rem;
      text-align: left;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #334155;
      font-size: 0.875rem;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .data-table tbody tr:last-child td {
      border-bottom: none;
    }

    .flag-info {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .flag-name {
      font-weight: 500;
      color: #0f172a;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .flag-description {
      font-size: 0.8125rem;
      color: #64748b;
    }

    .flag-key {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.8125rem;
      background: #f1f5f9;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      color: #475569;
    }

    .category-cell {
      color: #64748b;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      align-items: center;
      padding: 0.25rem 0.625rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .badge-purple {
      background: #f3e8ff;
      color: #7c3aed;
    }

    .badge-green {
      background: #dcfce7;
      color: #16a34a;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #64748b;
    }

    .badge-yellow {
      background: #fef3c7;
      color: #d97706;
    }

    .badge-blue {
      background: #dbeafe;
      color: #2563eb;
    }

    /* Empty State */
    .empty-container {
      text-align: center;
      padding: 3rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
    }

    .empty-container i {
      font-size: 3rem;
      color: #cbd5e1;
      margin-bottom: 1rem;
    }

    .empty-container h3 {
      font-size: 1rem;
      font-weight: 500;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .empty-container p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
    }

    /* Help Section */
    .help-section {
      margin-top: 2rem;
      padding: 1rem 1.25rem;
      background: #eff6ff;
      border: 1px solid #bfdbfe;
      border-radius: 0.5rem;
    }

    .help-section h3 {
      font-size: 0.875rem;
      font-weight: 600;
      color: #1e40af;
      margin: 0 0 0.75rem 0;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .help-section ul {
      margin: 0;
      padding-left: 1.25rem;
      color: #1e40af;
      font-size: 0.875rem;
    }

    .help-section li {
      margin-bottom: 0.25rem;
    }

    .help-section li:last-child {
      margin-bottom: 0;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdminFeatureFlagsComponent implements OnInit {
  private http = inject(HttpClient);

  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}`;

  /** Feature flags */
  flags = signal<FeatureFlag[]>([]);

  /** Loading state */
  loading = signal<boolean>(false);

  /** Error state */
  error = signal<string | null>(null);

  ngOnInit(): void {
    this.loadFlags();
  }

  /**
   * Load all feature flags from the admin API
   */
  loadFlags(): void {
    this.loading.set(true);
    this.error.set(null);

    this.http
      .get<ApiResponse<ListFlagsResponse>>(`${this.apiUrl}/admin/feature-flags`)
      .subscribe({
        next: (response) => {
          if (response.success && response.data) {
            this.flags.set(response.data.flags);
          } else {
            this.error.set('Failed to load feature flags');
          }
          this.loading.set(false);
        },
        error: (err) => {
          console.error('Failed to load feature flags:', err);
          this.error.set(err.error?.detail || 'Failed to load feature flags');
          this.loading.set(false);
        },
      });
  }
}
