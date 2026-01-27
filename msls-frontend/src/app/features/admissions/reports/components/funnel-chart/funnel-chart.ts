/**
 * MSLS Funnel Chart Component
 *
 * Displays an admission conversion funnel with a proper funnel shape
 * showing the progression from enquiry to enrollment.
 */

import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MslsSpinnerComponent } from '../../../../../shared/components/spinner/spinner.component';

import { FunnelStage } from '../../report.model';

/**
 * FunnelChartComponent - Displays admission conversion funnel.
 */
@Component({
  selector: 'app-funnel-chart',
  standalone: true,
  imports: [
    CommonModule,
    MslsSpinnerComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="funnel-card">
      <!-- Header -->
      <div class="funnel-header">
        <div class="header-left">
          <div class="header-icon">
            <i class="fa-solid fa-filter"></i>
          </div>
          <div class="header-text">
            <h2 class="header-title">{{ title() }}</h2>
            <p class="header-subtitle">Admission pipeline overview</p>
          </div>
        </div>
        @if (data().length > 0) {
          <div class="header-badge">
            <span class="badge-value">{{ overallConversion() | number: '1.0-1' }}%</span>
            <span class="badge-label">Overall</span>
          </div>
        }
      </div>

      <!-- Content -->
      <div class="funnel-content">
        @if (loading()) {
          <div class="loading-state">
            <msls-spinner size="lg" />
          </div>
        } @else if (data().length === 0) {
          <div class="empty-state">
            <i class="fa-regular fa-filter"></i>
            <p class="empty-title">No funnel data</p>
            <p class="empty-text">Data will appear when there are enquiries</p>
          </div>
        } @else {
          <!-- Funnel Visual -->
          <div class="funnel-stages">
            @for (stage of data(); track stage.name; let i = $index; let isLast = $last) {
              <div class="stage-row">
                <!-- Stage Bar -->
                <div class="stage-bar-section">
                  <div
                    class="stage-bar"
                    [style.width.%]="getFunnelWidth(i)"
                    [style.background]="getStageGradient(stage.variant)"
                  >
                    <span class="bar-count">{{ stage.count | number }}</span>
                  </div>
                </div>

                <!-- Stage Label -->
                <div class="stage-label-section">
                  <div class="stage-icon" [style.background]="getStageColor(stage.variant)">
                    <i [class]="stage.icon"></i>
                  </div>
                  <div class="stage-text">
                    <span class="stage-name">{{ stage.name }}</span>
                    @if (i > 0) {
                      <span class="stage-percent">{{ stage.percentage | number: '1.0-1' }}% of total</span>
                    } @else {
                      <span class="stage-percent">Starting point</span>
                    }
                  </div>
                </div>

                <!-- Connector -->
                @if (!isLast) {
                  <div class="stage-connector">
                    <div class="connector-line"></div>
                    <div class="connector-badge" [style.background]="getConversionBadgeColor(i)">
                      <i class="fa-solid fa-arrow-down"></i>
                      {{ getConversionRate(i) | number: '1.0-0' }}%
                    </div>
                    <div class="connector-line"></div>
                  </div>
                }
              </div>
            }
          </div>

          <!-- Summary Cards -->
          <div class="summary-section">
            <h3 class="summary-title">Conversion Rates</h3>
            <div class="summary-grid">
              <div class="summary-card">
                <div class="card-icon success">
                  <i class="fa-solid fa-arrow-trend-up"></i>
                </div>
                <div class="card-body">
                  <span class="card-value">{{ getConversionRate(0) | number: '1.0-1' }}%</span>
                  <span class="card-label">Enquiry → Application</span>
                </div>
              </div>
              <div class="summary-card">
                <div class="card-icon warning">
                  <i class="fa-solid fa-clipboard-check"></i>
                </div>
                <div class="card-body">
                  <span class="card-value">{{ getConversionRate(1) | number: '1.0-1' }}%</span>
                  <span class="card-label">Application → Approved</span>
                </div>
              </div>
              <div class="summary-card">
                <div class="card-icon primary">
                  <i class="fa-solid fa-user-graduate"></i>
                </div>
                <div class="card-body">
                  <span class="card-value">{{ getConversionRate(2) | number: '1.0-1' }}%</span>
                  <span class="card-label">Approved → Enrolled</span>
                </div>
              </div>
            </div>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .funnel-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    /* Header */
    .funnel-header {
      padding: 16px 20px;
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
    }

    .header-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .header-icon {
      width: 40px;
      height: 40px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 18px;
      color: white;
    }

    .header-text {
      display: flex;
      flex-direction: column;
    }

    .header-title {
      font-size: 16px;
      font-weight: 600;
      color: white;
      margin: 0;
    }

    .header-subtitle {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.8);
      margin: 2px 0 0 0;
    }

    .header-badge {
      display: flex;
      flex-direction: column;
      align-items: center;
      background: rgba(255, 255, 255, 0.2);
      padding: 8px 16px;
      border-radius: 8px;
      min-width: 70px;
    }

    .badge-value {
      font-size: 20px;
      font-weight: 700;
      color: white;
    }

    .badge-label {
      font-size: 10px;
      color: rgba(255, 255, 255, 0.9);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    /* Content */
    .funnel-content {
      padding: 20px;
    }

    .loading-state {
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 300px;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-height: 250px;
      text-align: center;
    }

    .empty-state i {
      font-size: 48px;
      color: #cbd5e1;
      margin-bottom: 16px;
    }

    .empty-title {
      font-size: 16px;
      font-weight: 600;
      color: #475569;
      margin: 0 0 4px 0;
    }

    .empty-text {
      font-size: 14px;
      color: #94a3b8;
      margin: 0;
    }

    /* Funnel Stages */
    .funnel-stages {
      display: flex;
      flex-direction: column;
      margin-bottom: 24px;
    }

    .stage-row {
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    .stage-bar-section {
      width: 100%;
      display: flex;
      justify-content: center;
      padding: 0 20px;
    }

    .stage-bar {
      height: 48px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      transition: all 0.3s ease;
      cursor: pointer;
      position: relative;
    }

    .stage-bar::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 50%;
      background: linear-gradient(to bottom, rgba(255,255,255,0.25), transparent);
      border-radius: 8px 8px 0 0;
    }

    .stage-bar:hover {
      transform: scale(1.02);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
    }

    .bar-count {
      font-size: 18px;
      font-weight: 700;
      color: white;
      text-shadow: 0 1px 2px rgba(0,0,0,0.2);
      z-index: 1;
    }

    .stage-label-section {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-top: 10px;
    }

    .stage-icon {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 12px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.15);
    }

    .stage-text {
      display: flex;
      flex-direction: column;
    }

    .stage-name {
      font-size: 14px;
      font-weight: 600;
      color: #1e293b;
    }

    .stage-percent {
      font-size: 12px;
      color: #64748b;
    }

    .stage-connector {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 12px 0;
      gap: 8px;
      width: 100%;
    }

    .connector-line {
      height: 2px;
      width: 40px;
      background: #e2e8f0;
    }

    .connector-badge {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 11px;
      font-weight: 600;
      color: white;
    }

    .connector-badge i {
      font-size: 10px;
    }

    /* Summary Section */
    .summary-section {
      border-top: 1px solid #e2e8f0;
      padding-top: 20px;
      margin-top: 8px;
    }

    .summary-title {
      font-size: 14px;
      font-weight: 600;
      color: #475569;
      margin: 0 0 12px 0;
    }

    .summary-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 12px;
    }

    .summary-card {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 12px;
      background: #f8fafc;
      border-radius: 8px;
      border: 1px solid #e2e8f0;
    }

    .card-icon {
      width: 36px;
      height: 36px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 14px;
      flex-shrink: 0;
    }

    .card-icon.success {
      background: linear-gradient(135deg, #10b981, #059669);
    }

    .card-icon.warning {
      background: linear-gradient(135deg, #f59e0b, #d97706);
    }

    .card-icon.primary {
      background: linear-gradient(135deg, #6366f1, #4f46e5);
    }

    .card-body {
      display: flex;
      flex-direction: column;
      min-width: 0;
    }

    .card-value {
      font-size: 18px;
      font-weight: 700;
      color: #0f172a;
      line-height: 1.2;
    }

    .card-label {
      font-size: 11px;
      color: #64748b;
      line-height: 1.3;
    }

    /* Responsive */
    @media (max-width: 768px) {
      .summary-grid {
        grid-template-columns: 1fr;
      }

      .header-badge {
        padding: 6px 12px;
        min-width: 60px;
      }

      .badge-value {
        font-size: 16px;
      }

      .stage-bar {
        height: 40px;
      }

      .bar-count {
        font-size: 16px;
      }
    }

    @media (max-width: 480px) {
      .funnel-header {
        flex-wrap: wrap;
      }

      .header-subtitle {
        display: none;
      }

      .summary-card {
        padding: 10px;
      }

      .card-value {
        font-size: 16px;
      }
    }
  `],
})
export class FunnelChartComponent {
  /** Funnel data to display */
  readonly data = input<FunnelStage[]>([]);

  /** Loading state */
  readonly loading = input<boolean>(false);

  /** Card title */
  readonly title = input<string>('Conversion Funnel');

  /** Whether to show arrows between stages */
  readonly showArrows = input<boolean>(true);

  /** Calculate overall conversion rate */
  readonly overallConversion = computed(() => {
    const stages = this.data();
    if (stages.length < 2) return 0;
    const first = stages[0]?.count || 0;
    const last = stages[stages.length - 1]?.count || 0;
    return first > 0 ? (last / first) * 100 : 0;
  });

  /**
   * Get funnel width percentage based on stage index
   */
  getFunnelWidth(index: number): number {
    const widths = [100, 75, 55, 40];
    return widths[index] || 35;
  }

  /**
   * Get stage color
   */
  getStageColor(variant: FunnelStage['variant']): string {
    const colors: Record<FunnelStage['variant'], string> = {
      info: '#3b82f6',
      primary: '#8b5cf6',
      warning: '#f59e0b',
      success: '#10b981',
    };
    return colors[variant];
  }

  /**
   * Get stage gradient
   */
  getStageGradient(variant: FunnelStage['variant']): string {
    const gradients: Record<FunnelStage['variant'], string> = {
      info: 'linear-gradient(135deg, #3b82f6, #2563eb)',
      primary: 'linear-gradient(135deg, #8b5cf6, #7c3aed)',
      warning: 'linear-gradient(135deg, #f59e0b, #d97706)',
      success: 'linear-gradient(135deg, #10b981, #059669)',
    };
    return gradients[variant];
  }

  /**
   * Get conversion badge color based on rate
   */
  getConversionBadgeColor(index: number): string {
    const rate = this.getConversionRate(index);
    if (rate >= 70) return '#10b981';
    if (rate >= 40) return '#f59e0b';
    return '#ef4444';
  }

  /**
   * Get conversion rate between stages
   */
  getConversionRate(index: number): number {
    const stages = this.data();
    if (index >= stages.length - 1) return 0;
    const current = stages[index]?.count || 0;
    const next = stages[index + 1]?.count || 0;
    return current > 0 ? (next / current) * 100 : 0;
  }
}
