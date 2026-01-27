/**
 * MSLS Source Analysis Component
 *
 * Displays enquiry source analysis with a horizontal bar chart
 * and percentage breakdown.
 */

import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MslsSpinnerComponent } from '../../../../../shared/components/spinner/spinner.component';

import { SourceAnalysis as SourceAnalysisData } from '../../report.model';

/**
 * SourceAnalysisComponent - Displays enquiry source breakdown.
 */
@Component({
  selector: 'app-source-analysis',
  standalone: true,
  imports: [
    CommonModule,
    MslsSpinnerComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="chart-card">
      <div class="card-header">
        <div class="header-icon">
          <i class="fa-solid fa-magnifying-glass-location"></i>
        </div>
        <h2 class="card-title">Enquiry Sources</h2>
        <div class="header-stats">
          <span class="total-label">Total:</span>
          <span class="total-value">{{ totalCount() | number }}</span>
        </div>
      </div>

      <div class="card-content">
        @if (loading()) {
          <div class="flex items-center justify-center h-48">
            <msls-spinner size="lg" />
          </div>
        } @else if (data().length === 0) {
          <div class="empty-state">
            <i class="fa-regular fa-chart-bar"></i>
            <p>No source data available</p>
          </div>
        } @else {
          <div class="sources-container">
            @for (source of sortedData(); track source.source) {
              <div class="source-item">
                <div class="source-header">
                  <div class="source-info">
                    <div class="source-icon" [style.background]="getSourceBgColor(source.source)">
                      <i [class]="getSourceIcon(source.source) + ' text-white'"></i>
                    </div>
                    <div class="source-label">
                      <p class="source-name">{{ source.label }}</p>
                    </div>
                  </div>
                  <div class="source-value">
                    <p class="value-number">{{ source.count | number }}</p>
                    <p class="value-percent" [style.color]="getSourceBarColor(source.source)">{{ source.percentage | number: '1.0-1' }}%</p>
                  </div>
                </div>
                <!-- Source Bar -->
                <div class="source-bar-wrapper">
                  <div class="source-bar" [style.width.%]="source.percentage" [style.background]="getSourceBarColor(source.source)"></div>
                </div>
              </div>
            }
          </div>

          <!-- Color Legend -->
          <div class="sources-legend">
            @for (source of topSources(); track source.source) {
              <div class="legend-item">
                <div class="legend-dot" [style.background]="getSourceBarColor(source.source)"></div>
                <span class="legend-label">{{ source.label }}</span>
              </div>
            }
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .chart-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
      box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .card-header {
      padding: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
      background: linear-gradient(to right, #f0f9ff, #f8fafc);
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .header-icon {
      font-size: 1.25rem;
      color: #4f46e5;
      flex-shrink: 0;
    }

    .card-title {
      font-size: 1rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
      flex: 1;
    }

    .header-stats {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      flex-shrink: 0;
    }

    .total-label {
      color: #64748b;
    }

    .total-value {
      font-weight: 600;
      color: #0f172a;
    }

    .card-content {
      padding: 1.5rem;
      min-height: 300px;
    }

    .sources-container {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .source-item {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .source-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .source-info {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      flex: 1;
    }

    .source-icon {
      width: 2rem;
      height: 2rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.875rem;
      flex-shrink: 0;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .source-label {
      flex: 1;
    }

    .source-name {
      font-size: 0.875rem;
      font-weight: 500;
      color: #0f172a;
      margin: 0;
    }

    .source-value {
      text-align: right;
      min-width: 70px;
    }

    .value-number {
      font-size: 1rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0;
    }

    .value-percent {
      font-size: 0.75rem;
      font-weight: 600;
      margin: 0.15rem 0 0 0;
    }

    .source-bar-wrapper {
      width: 100%;
      height: 0.375rem;
      background: #f1f5f9;
      border-radius: 9999px;
      overflow: hidden;
    }

    .source-bar {
      height: 100%;
      border-radius: 9999px;
      transition: all 0.3s ease;
      min-width: 20px;
    }

    .sources-legend {
      margin-top: 1.5rem;
      padding-top: 1rem;
      border-top: 1px solid #f1f5f9;
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 0.75rem;
    }

    .legend-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.8rem;
    }

    .legend-dot {
      width: 0.5rem;
      height: 0.5rem;
      border-radius: 0.25rem;
      flex-shrink: 0;
    }

    .legend-label {
      color: #64748b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      padding: 2rem;
      text-align: center;
    }

    .empty-state i {
      font-size: 2.5rem;
      color: #cbd5e1;
      margin-bottom: 1rem;
    }

    .empty-state p {
      color: #64748b;
      margin: 0.5rem 0;
    }
  `],
})
export class SourceAnalysisComponent {
  /** Source analysis data */
  readonly data = input<SourceAnalysisData[]>([]);

  /** Loading state */
  readonly loading = input<boolean>(false);

  /** Total count of enquiries */
  readonly totalCount = computed(() => {
    return this.data().reduce((sum, s) => sum + s.count, 0);
  });

  /** Sorted data by count descending */
  readonly sortedData = computed(() => {
    return [...this.data()].sort((a, b) => b.count - a.count);
  });

  /** Top 5 sources for legend */
  readonly topSources = computed(() => {
    return this.sortedData().slice(0, 5);
  });

  /** Get icon for source */
  getSourceIcon(source: string): string {
    const icons: Record<string, string> = {
      walk_in: 'fa-solid fa-person-walking',
      website: 'fa-solid fa-globe',
      referral: 'fa-solid fa-user-group',
      phone: 'fa-solid fa-phone',
      advertisement: 'fa-solid fa-bullhorn',
      social_media: 'fa-solid fa-share-nodes',
      newspaper: 'fa-solid fa-newspaper',
      other: 'fa-solid fa-ellipsis',
    };
    return icons[source] || 'fa-solid fa-circle';
  }

  /** Get background class for source icon */
  getSourceBgClass(source: string): string {
    const classes: Record<string, string> = {
      walk_in: 'bg-blue-100',
      website: 'bg-purple-100',
      referral: 'bg-emerald-100',
      phone: 'bg-amber-100',
      advertisement: 'bg-rose-100',
      social_media: 'bg-cyan-100',
      newspaper: 'bg-slate-100',
      other: 'bg-gray-100',
    };
    return classes[source] || 'bg-gray-100';
  }

  /** Get icon color class for source */
  getSourceIconClass(source: string): string {
    const classes: Record<string, string> = {
      walk_in: 'text-blue-600',
      website: 'text-purple-600',
      referral: 'text-emerald-600',
      phone: 'text-amber-600',
      advertisement: 'text-rose-600',
      social_media: 'text-cyan-600',
      newspaper: 'text-slate-600',
      other: 'text-gray-600',
    };
    return classes[source] || 'text-gray-600';
  }

  /** Get badge variant for source */
  getSourceVariant(source: string): 'info' | 'primary' | 'success' | 'warning' | 'danger' {
    const variants: Record<string, 'info' | 'primary' | 'success' | 'warning' | 'danger'> = {
      walk_in: 'info',
      website: 'primary',
      referral: 'success',
      phone: 'warning',
      advertisement: 'danger',
      social_media: 'info',
      newspaper: 'primary',
      other: 'info',
    };
    return variants[source] || 'info';
  }

  /** Get background color for source icon */
  getSourceBgColor(source: string): string {
    const colors: Record<string, string> = {
      walk_in: '#3b82f6',
      website: '#a855f7',
      referral: '#10b981',
      phone: '#f59e0b',
      advertisement: '#ef4444',
      social_media: '#06b6d4',
      newspaper: '#64748b',
      other: '#9ca3af',
    };
    return colors[source] || '#9ca3af';
  }

  /** Get progress bar color for source */
  getSourceBarColor(source: string): string {
    const colors: Record<string, string> = {
      walk_in: '#3b82f6',
      website: '#a855f7',
      referral: '#10b981',
      phone: '#f59e0b',
      advertisement: '#ef4444',
      social_media: '#06b6d4',
      newspaper: '#64748b',
      other: '#9ca3af',
    };
    return colors[source] || '#9ca3af';
  }
}
