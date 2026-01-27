/**
 * MSLS Class-wise Admission Report Component
 *
 * Displays class-wise seat availability with detailed breakdown
 * of applications, approvals, enrollments, and vacancies.
 */

import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MslsBadgeComponent, BadgeVariant } from '../../../shared/components/badge/badge.component';
import { MslsSpinnerComponent } from '../../../shared/components/spinner/spinner.component';

import { ClassWiseReport, getFillStatusConfig } from './report.model';

/**
 * ClassWiseReportComponent - Displays class-wise seat availability table.
 */
@Component({
  selector: 'app-class-wise-report',
  standalone: true,
  imports: [
    CommonModule,
    MslsBadgeComponent,
    MslsSpinnerComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="chart-card">
      <div class="card-header">
        <div class="header-icon">
          <i class="fa-solid fa-chairs"></i>
        </div>
        <h2 class="card-title">Class-wise Seat Availability</h2>
        <div class="header-stats">
          <span class="total-label">Classes:</span>
          <span class="total-value">{{ data().length }}</span>
        </div>
      </div>

      <div class="card-content">
        @if (loading()) {
          <div class="flex items-center justify-center h-64">
            <msls-spinner size="lg" />
          </div>
        } @else if (data().length === 0) {
          <div class="empty-state">
            <i class="fa-regular fa-folder-open"></i>
            <p>No class data available</p>
          </div>
        } @else {
          <div class="classes-container">
            @for (row of data(); track row.className) {
              <div class="class-card">
                <div class="class-header">
                  <div class="class-info">
                    <div class="class-icon">
                      <i class="fa-solid fa-chalkboard"></i>
                    </div>
                    <div class="class-label">
                      <p class="class-name">{{ row.className }}</p>
                      <p class="class-meta">{{ row.totalSeats }} total seats</p>
                    </div>
                  </div>
                  <div class="class-fill-status">
                    <div class="fill-percentage">{{ row.fillPercentage | number: '1.0-0' }}%</div>
                    <msls-badge [variant]="getFillStatus(row.fillPercentage).variant" size="sm">
                      {{ getFillStatus(row.fillPercentage).label }}
                    </msls-badge>
                  </div>
                </div>

                <!-- Fill Bar -->
                <div class="fill-bar-wrapper">
                  <div class="fill-bar" [style.width.%]="row.fillPercentage" [style.background]="getFillBarColor(row.fillPercentage)"></div>
                </div>

                <!-- Stats Grid -->
                <div class="class-stats">
                  <div class="stat-item">
                    <span class="stat-label">Applications</span>
                    <span class="stat-value" [class.bg-amber]="row.applications > row.totalSeats">{{ row.applications }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">Approved</span>
                    <span class="stat-value bg-purple">{{ row.approved }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">Enrolled</span>
                    <span class="stat-value bg-emerald">{{ row.enrolled }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">Vacant</span>
                    <span class="stat-value" [class.bg-red]="row.vacant === 0" [class.bg-green]="row.vacant > 0">{{ row.vacant }}</span>
                  </div>
                  @if (row.waitlisted > 0) {
                    <div class="stat-item">
                      <span class="stat-label">Waitlisted</span>
                      <span class="stat-value bg-orange">{{ row.waitlisted }}</span>
                    </div>
                  }
                </div>
              </div>
            }
          </div>

          <!-- Summary Footer -->
          <div class="summary-footer">
            <div class="footer-item">
              <span class="footer-label">Total Seats:</span>
              <span class="footer-value">{{ totals().totalSeats | number }}</span>
            </div>
            <div class="footer-item">
              <span class="footer-label">Total Enrolled:</span>
              <span class="footer-value fg-emerald">{{ totals().enrolled | number }}</span>
            </div>
            <div class="footer-item">
              <span class="footer-label">Total Vacant:</span>
              <span class="footer-value fg-blue">{{ totals().vacant | number }}</span>
            </div>
            <div class="footer-item">
              <span class="footer-label">Overall Fill Rate:</span>
              <span class="footer-value fg-primary">{{ totals().overallFillRate | number: '1.0-1' }}%</span>
            </div>
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
    }

    .classes-container {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .class-card {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      padding: 1rem;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .class-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
    }

    .class-info {
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      flex: 1;
    }

    .class-icon {
      width: 2rem;
      height: 2rem;
      background: #e0e7ff;
      border-radius: 0.375rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #4f46e5;
      font-size: 0.875rem;
      flex-shrink: 0;
    }

    .class-label {
      flex: 1;
      min-width: 0;
    }

    .class-name {
      font-size: 0.95rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .class-meta {
      font-size: 0.75rem;
      color: #64748b;
      margin: 0.15rem 0 0 0;
    }

    .class-fill-status {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 0.25rem;
    }

    .fill-percentage {
      font-size: 1.25rem;
      font-weight: 700;
      color: #0f172a;
    }

    .fill-bar-wrapper {
      width: 100%;
      height: 0.5rem;
      background: #e2e8f0;
      border-radius: 9999px;
      overflow: hidden;
    }

    .fill-bar {
      height: 100%;
      border-radius: 9999px;
      transition: all 0.3s ease;
      min-width: 5px;
    }

    .class-stats {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(60px, 1fr));
      gap: 0.5rem;
    }

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.25rem;
      padding: 0.5rem;
      background: white;
      border-radius: 0.375rem;
      border: 1px solid #f1f5f9;
    }

    .stat-label {
      font-size: 0.65rem;
      color: #64748b;
      text-align: center;
      line-height: 1.2;
    }

    .stat-value {
      font-size: 0.9rem;
      font-weight: 700;
      color: white;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      min-width: 28px;
      text-align: center;
    }

    .stat-value.bg-blue {
      background: #3b82f6;
    }

    .stat-value.bg-purple {
      background: #a855f7;
    }

    .stat-value.bg-emerald {
      background: #10b981;
    }

    .stat-value.bg-orange {
      background: #f97316;
    }

    .stat-value.bg-green {
      background: #22c55e;
    }

    .stat-value.bg-red {
      background: #ef4444;
    }

    .stat-value.bg-amber {
      background: #f59e0b;
    }

    .summary-footer {
      border-top: 1px solid #e2e8f0;
      padding-top: 1rem;
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 1rem;
    }

    .footer-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .footer-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .footer-value {
      font-size: 1.25rem;
      font-weight: 700;
      color: #0f172a;
    }

    .footer-value.fg-emerald {
      color: #10b981;
    }

    .footer-value.fg-blue {
      color: #3b82f6;
    }

    .footer-value.fg-primary {
      color: #4f46e5;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 300px;
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

    @media (max-width: 768px) {
      .classes-container {
        grid-template-columns: 1fr;
      }

      .summary-footer {
        grid-template-columns: repeat(2, 1fr);
      }
    }
  `],
})
export class ClassWiseReportComponent {
  /** Class-wise report data */
  readonly data = input<ClassWiseReport[]>([]);

  /** Loading state */
  readonly loading = input<boolean>(false);

  /** Computed totals for summary */
  readonly totals = computed(() => {
    const classes = this.data();
    const totalSeats = classes.reduce((sum, c) => sum + c.totalSeats, 0);
    const enrolled = classes.reduce((sum, c) => sum + c.enrolled, 0);
    const vacant = classes.reduce((sum, c) => sum + c.vacant, 0);
    const overallFillRate = totalSeats > 0 ? (enrolled / totalSeats) * 100 : 0;

    return {
      totalSeats,
      enrolled,
      vacant,
      overallFillRate,
    };
  });

  /**
   * Get fill status configuration based on percentage
   */
  getFillStatus(fillPercentage: number): { label: string; variant: BadgeVariant } {
    return getFillStatusConfig(fillPercentage);
  }

  /**
   * Get progress bar color class based on fill percentage
   */
  getFillBarColor(fillPercentage: number): string {
    if (fillPercentage >= 90) {
      return '#ef4444'; // Red for almost full
    } else if (fillPercentage >= 70) {
      return '#f59e0b'; // Amber for high
    } else if (fillPercentage >= 40) {
      return '#3b82f6'; // Blue for medium
    } else {
      return '#10b981'; // Green for low
    }
  }
}
