/**
 * MSLS Admission Dashboard Component
 *
 * Main dashboard for admission reports and analytics.
 * Displays summary cards, conversion funnel, class-wise seat availability,
 * and source analysis with a consistent, modern design.
 */

import { Component, OnInit, inject, signal, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsSelectComponent, SelectOption } from '../../../shared/components/select/select.component';
import { MslsSpinnerComponent } from '../../../shared/components/spinner/spinner.component';

import { AdmissionReportService } from './report.service';
import {
  DashboardStats,
  FunnelStage,
  ClassWiseReport,
  SourceAnalysis,
  DailyTrendPoint,
  ReportFilterParams,
  AdmissionSessionOption,
} from './report.model';
import { ClassWiseReportComponent } from './class-wise-report.component';
import { SourceAnalysisComponent } from './components/source-analysis/source-analysis';
import { FunnelChartComponent } from './components/funnel-chart/funnel-chart';

/**
 * AdmissionDashboardComponent - Main dashboard for admission analytics.
 */
@Component({
  selector: 'app-admission-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsSelectComponent,
    MslsSpinnerComponent,
    ClassWiseReportComponent,
    SourceAnalysisComponent,
    FunnelChartComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="dashboard-container">
      <!-- Header -->
      <header class="dashboard-header">
        <div class="header-content">
          <div class="header-title-section">
            <div class="header-icon-wrapper">
              <i class="fa-solid fa-chart-pie"></i>
            </div>
            <div class="header-text">
              <h1>Admission Dashboard</h1>
              <p>Track admission progress and analyze conversion metrics</p>
            </div>
          </div>
          <div class="header-actions">
            <div class="session-filter">
              <label class="filter-label">Session</label>
              <msls-select
                [options]="sessionOptions()"
                placeholder="All Sessions"
                [searchable]="true"
                (valueChange)="onSessionChange($event)"
                class="session-select"
              ></msls-select>
            </div>
            <button class="btn-refresh" (click)="refreshData()" [disabled]="loading()">
              <i class="fa-solid fa-arrows-rotate" [class.fa-spin]="loading()"></i>
            </button>
          </div>
        </div>
      </header>

      <!-- Main Content -->
      <main class="dashboard-main">
        <!-- Loading Overlay -->
        @if (initialLoading()) {
          <div class="loading-container">
            <div class="loading-content">
              <msls-spinner size="lg" />
              <p>Loading dashboard...</p>
            </div>
          </div>
        } @else if (hasError()) {
          <!-- Error State -->
          <div class="error-container">
            <div class="error-content">
              <i class="fa-solid fa-exclamation-triangle"></i>
              <h3>Unable to load dashboard</h3>
              <p>{{ errorMessage() }}</p>
              <button class="btn-retry" (click)="refreshData()">
                <i class="fa-solid fa-arrows-rotate"></i>
                Try Again
              </button>
            </div>
          </div>
        } @else {
          <!-- Stats Cards -->
          <section class="stats-section">
            <div class="stats-grid">
              <!-- Enquiries Card -->
              <div class="stat-card stat-enquiries">
                <div class="stat-icon-bg">
                  <i class="fa-solid fa-phone"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">Total Enquiries</span>
                  <span class="stat-value">{{ stats()?.totalEnquiries || 0 | number }}</span>
                  <span class="stat-trend">
                    <i class="fa-solid fa-arrow-trend-up"></i>
                    Active pipeline
                  </span>
                </div>
              </div>

              <!-- Applications Card -->
              <div class="stat-card stat-applications">
                <div class="stat-icon-bg">
                  <i class="fa-solid fa-file-lines"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">Applications</span>
                  <span class="stat-value">{{ stats()?.totalApplications || 0 | number }}</span>
                  <span class="stat-trend trend-info">
                    <i class="fa-solid fa-percent"></i>
                    {{ stats()?.conversionRates?.enquiryToApplication || 0 | number:'1.0-1' }}% from enquiries
                  </span>
                </div>
              </div>

              <!-- Approved Card -->
              <div class="stat-card stat-approved">
                <div class="stat-icon-bg">
                  <i class="fa-solid fa-circle-check"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">Approved</span>
                  <span class="stat-value">{{ stats()?.approved || 0 | number }}</span>
                  <span class="stat-trend trend-warning">
                    <i class="fa-solid fa-clock"></i>
                    {{ stats()?.pending || 0 }} pending review
                  </span>
                </div>
              </div>

              <!-- Enrolled Card -->
              <div class="stat-card stat-enrolled">
                <div class="stat-icon-bg">
                  <i class="fa-solid fa-graduation-cap"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">Enrolled</span>
                  <span class="stat-value">{{ stats()?.enrolled || 0 | number }}</span>
                  <span class="stat-trend trend-success">
                    <i class="fa-solid fa-check"></i>
                    {{ stats()?.conversionRates?.approvedToEnrolled || 0 | number:'1.0-1' }}% of approved
                  </span>
                </div>
              </div>
            </div>
          </section>

          <!-- Charts Section -->
          <section class="charts-section">
            <div class="charts-row">
              <!-- Conversion Funnel -->
              <div class="chart-card funnel-card">
                <div class="card-header">
                  <i class="fa-solid fa-filter header-icon"></i>
                  <h2>Conversion Funnel</h2>
                </div>
                <div class="card-body">
                  @if (funnelLoading()) {
                    <div class="card-loading">
                      <msls-spinner size="md" />
                    </div>
                  } @else if (funnelData().length === 0) {
                    <div class="card-empty">
                      <i class="fa-regular fa-chart-bar"></i>
                      <p>No funnel data available</p>
                    </div>
                  } @else {
                    <app-funnel-chart [data]="funnelData()" [loading]="funnelLoading()" />
                  }
                </div>
              </div>

              <!-- Conversion Rates -->
              <div class="chart-card rates-card">
                <div class="card-header">
                  <i class="fa-solid fa-chart-line header-icon"></i>
                  <h2>Conversion Rates</h2>
                </div>
                <div class="card-body">
                  @if (stats()) {
                    <div class="rates-container">
                      <!-- Rate Items -->
                      <div class="rate-item">
                        <div class="rate-info">
                          <div class="rate-label-row">
                            <span class="rate-label">Enquiry → Application</span>
                            <span class="rate-percentage">{{ stats()!.conversionRates.enquiryToApplication | number:'1.1-1' }}%</span>
                          </div>
                          <div class="rate-bar-bg">
                            <div class="rate-bar rate-bar-blue" [style.width.%]="stats()!.conversionRates.enquiryToApplication"></div>
                          </div>
                        </div>
                      </div>

                      <div class="rate-item">
                        <div class="rate-info">
                          <div class="rate-label-row">
                            <span class="rate-label">Application → Approved</span>
                            <span class="rate-percentage">{{ stats()!.conversionRates.applicationToApproved | number:'1.1-1' }}%</span>
                          </div>
                          <div class="rate-bar-bg">
                            <div class="rate-bar rate-bar-amber" [style.width.%]="stats()!.conversionRates.applicationToApproved"></div>
                          </div>
                        </div>
                      </div>

                      <div class="rate-item">
                        <div class="rate-info">
                          <div class="rate-label-row">
                            <span class="rate-label">Approved → Enrolled</span>
                            <span class="rate-percentage">{{ stats()!.conversionRates.approvedToEnrolled | number:'1.1-1' }}%</span>
                          </div>
                          <div class="rate-bar-bg">
                            <div class="rate-bar rate-bar-emerald" [style.width.%]="stats()!.conversionRates.approvedToEnrolled"></div>
                          </div>
                        </div>
                      </div>

                      <!-- Overall Conversion -->
                      <div class="overall-rate">
                        <div class="overall-label">
                          <i class="fa-solid fa-bullseye"></i>
                          Overall Conversion
                        </div>
                        <div class="overall-value">{{ overallConversion() | number:'1.1-1' }}%</div>
                        <div class="overall-subtitle">Enquiry to Enrollment</div>
                      </div>
                    </div>
                  }
                </div>
              </div>
            </div>
          </section>

          <!-- Analytics Section -->
          <section class="analytics-section">
            <div class="analytics-row">
              <!-- Source Analysis -->
              <app-source-analysis [data]="sourceData()" [loading]="sourceLoading()" />

              <!-- Daily Trend -->
              <div class="chart-card trend-card">
                <div class="card-header">
                  <i class="fa-solid fa-chart-area header-icon"></i>
                  <h2>Daily Trend</h2>
                  <span class="header-badge">Last 30 Days</span>
                </div>
                <div class="card-body">
                  @if (trendLoading()) {
                    <div class="card-loading">
                      <msls-spinner size="md" />
                    </div>
                  } @else if (trendData().length === 0) {
                    <div class="card-empty">
                      <i class="fa-regular fa-calendar"></i>
                      <p>No trend data available</p>
                    </div>
                  } @else {
                    <div class="trend-content">
                      <!-- Summary Stats -->
                      <div class="trend-summary">
                        <div class="trend-stat">
                          <span class="trend-stat-value text-blue">{{ trendTotals().enquiries | number }}</span>
                          <span class="trend-stat-label">Enquiries</span>
                        </div>
                        <div class="trend-stat">
                          <span class="trend-stat-value text-purple">{{ trendTotals().applications | number }}</span>
                          <span class="trend-stat-label">Applications</span>
                        </div>
                        <div class="trend-stat">
                          <span class="trend-stat-value text-emerald">{{ trendTotals().enrollments | number }}</span>
                          <span class="trend-stat-label">Enrollments</span>
                        </div>
                      </div>

                      <!-- Mini Chart -->
                      <div class="trend-chart-container">
                        <div class="trend-chart">
                          @for (point of trendDataSubset(); track point.date; let i = $index) {
                            <div class="trend-bar-group" [title]="point.date">
                              <div class="trend-bar trend-bar-enquiry" [style.height.px]="getBarHeight(point.enquiries, trendMax())"></div>
                              <div class="trend-bar trend-bar-app" [style.height.px]="getBarHeight(point.applications, trendMax())"></div>
                            </div>
                          }
                        </div>
                      </div>

                      <!-- Legend -->
                      <div class="trend-legend">
                        <span class="legend-item"><span class="legend-dot bg-blue"></span>Enquiries</span>
                        <span class="legend-item"><span class="legend-dot bg-purple"></span>Applications</span>
                      </div>
                    </div>
                  }
                </div>
              </div>
            </div>
          </section>

          <!-- Class-wise Section -->
          <section class="class-section">
            <app-class-wise-report [data]="classWiseData()" [loading]="classWiseLoading()" />
          </section>
        }
      </main>
    </div>
  `,
  styles: [`
    /* Container */
    .dashboard-container {
      min-height: 100vh;
      background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    }

    /* Header */
    .dashboard-header {
      background: white;
      border-bottom: 1px solid #e2e8f0;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
      position: sticky;
      top: 0;
      z-index: 10;
    }

    .header-content {
      max-width: 1400px;
      margin: 0 auto;
      padding: 1rem 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    @media (min-width: 768px) {
      .header-content {
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        padding: 1.25rem 2rem;
      }
    }

    .header-title-section {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon-wrapper {
      width: 3rem;
      height: 3rem;
      background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 1.25rem;
      box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.3);
    }

    .header-text h1 {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0;
    }

    .header-text p {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0.25rem 0 0 0;
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .session-filter {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .filter-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .session-select {
      min-width: 200px;
    }

    .btn-refresh {
      width: 2.5rem;
      height: 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .btn-refresh:hover:not(:disabled) {
      background: #f1f5f9;
      color: #4f46e5;
      border-color: #4f46e5;
    }

    .btn-refresh:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    /* Main Content */
    .dashboard-main {
      max-width: 1400px;
      margin: 0 auto;
      padding: 1.5rem;
    }

    @media (min-width: 768px) {
      .dashboard-main {
        padding: 2rem;
      }
    }

    /* Loading & Error States */
    .loading-container,
    .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 60vh;
    }

    .loading-content,
    .error-content {
      text-align: center;
    }

    .loading-content p {
      margin-top: 1rem;
      color: #64748b;
    }

    .error-content i {
      font-size: 3rem;
      color: #ef4444;
      margin-bottom: 1rem;
    }

    .error-content h3 {
      font-size: 1.25rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .error-content p {
      color: #64748b;
      margin: 0 0 1.5rem 0;
    }

    .btn-retry {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      background: #4f46e5;
      color: white;
      border: none;
      border-radius: 0.5rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.2s;
    }

    .btn-retry:hover {
      background: #4338ca;
    }

    /* Stats Section */
    .stats-section {
      margin-bottom: 2rem;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    @media (min-width: 768px) {
      .stats-grid {
        grid-template-columns: repeat(4, 1fr);
        gap: 1.5rem;
      }
    }

    .stat-card {
      background: white;
      border-radius: 1rem;
      padding: 1.25rem;
      display: flex;
      align-items: flex-start;
      gap: 1rem;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      transition: all 0.2s;
    }

    .stat-card:hover {
      box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.1);
      transform: translateY(-2px);
    }

    .stat-icon-bg {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      flex-shrink: 0;
    }

    .stat-enquiries .stat-icon-bg {
      background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
      color: #2563eb;
    }

    .stat-applications .stat-icon-bg {
      background: linear-gradient(135deg, #f3e8ff 0%, #e9d5ff 100%);
      color: #9333ea;
    }

    .stat-approved .stat-icon-bg {
      background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
      color: #d97706;
    }

    .stat-enrolled .stat-icon-bg {
      background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
      color: #059669;
    }

    .stat-content {
      flex: 1;
      min-width: 0;
    }

    .stat-label {
      display: block;
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .stat-value {
      display: block;
      font-size: 1.75rem;
      font-weight: 700;
      color: #0f172a;
      line-height: 1.2;
      margin-top: 0.25rem;
    }

    .stat-trend {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.5rem;
    }

    .stat-trend i {
      font-size: 0.625rem;
    }

    .trend-info {
      color: #3b82f6;
    }

    .trend-warning {
      color: #f59e0b;
    }

    .trend-success {
      color: #10b981;
    }

    /* Charts Section */
    .charts-section,
    .analytics-section,
    .class-section {
      margin-bottom: 2rem;
    }

    .charts-row,
    .analytics-row {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }

    @media (min-width: 1024px) {
      .charts-row,
      .analytics-row {
        grid-template-columns: 1fr 1fr;
      }
    }

    .chart-card {
      background: white;
      border-radius: 1rem;
      overflow: hidden;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
    }

    .card-header {
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #f1f5f9;
      display: flex;
      align-items: center;
      gap: 0.75rem;
      background: linear-gradient(to right, #fafbfc, #ffffff);
    }

    .header-icon {
      color: #4f46e5;
      font-size: 1.125rem;
    }

    .card-header h2 {
      font-size: 1rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
      flex: 1;
    }

    .header-badge {
      font-size: 0.75rem;
      padding: 0.25rem 0.625rem;
      background: #f1f5f9;
      color: #64748b;
      border-radius: 9999px;
    }

    .card-body {
      padding: 1.5rem;
      min-height: 280px;
    }

    .card-loading,
    .card-empty {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 240px;
      color: #94a3b8;
    }

    .card-empty i {
      font-size: 2.5rem;
      margin-bottom: 0.75rem;
    }

    .card-empty p {
      font-size: 0.875rem;
      margin: 0;
    }

    /* Rates Card */
    .rates-container {
      display: flex;
      flex-direction: column;
      gap: 1.25rem;
    }

    .rate-item {
      padding-bottom: 1.25rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .rate-item:last-of-type {
      border-bottom: none;
      padding-bottom: 0;
    }

    .rate-label-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }

    .rate-label {
      font-size: 0.875rem;
      color: #475569;
    }

    .rate-percentage {
      font-size: 0.875rem;
      font-weight: 600;
      color: #0f172a;
    }

    .rate-bar-bg {
      height: 0.5rem;
      background: #f1f5f9;
      border-radius: 9999px;
      overflow: hidden;
    }

    .rate-bar {
      height: 100%;
      border-radius: 9999px;
      transition: width 0.5s ease;
    }

    .rate-bar-blue {
      background: linear-gradient(to right, #3b82f6, #60a5fa);
    }

    .rate-bar-amber {
      background: linear-gradient(to right, #f59e0b, #fbbf24);
    }

    .rate-bar-emerald {
      background: linear-gradient(to right, #10b981, #34d399);
    }

    .overall-rate {
      background: linear-gradient(135deg, #f0f9ff 0%, #eff6ff 100%);
      border-radius: 0.75rem;
      padding: 1.25rem;
      text-align: center;
      margin-top: 0.5rem;
    }

    .overall-label {
      font-size: 0.875rem;
      color: #64748b;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
    }

    .overall-label i {
      color: #4f46e5;
    }

    .overall-value {
      font-size: 2rem;
      font-weight: 700;
      color: #4f46e5;
      line-height: 1;
    }

    .overall-subtitle {
      font-size: 0.75rem;
      color: #94a3b8;
      margin-top: 0.25rem;
    }

    /* Trend Card */
    .trend-content {
      display: flex;
      flex-direction: column;
      gap: 1.25rem;
    }

    .trend-summary {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 1rem;
    }

    .trend-stat {
      text-align: center;
      padding: 0.75rem;
      background: #f8fafc;
      border-radius: 0.5rem;
    }

    .trend-stat-value {
      display: block;
      font-size: 1.5rem;
      font-weight: 700;
      line-height: 1;
    }

    .trend-stat-label {
      display: block;
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.25rem;
    }

    .text-blue { color: #3b82f6; }
    .text-purple { color: #9333ea; }
    .text-emerald { color: #10b981; }

    .trend-chart-container {
      background: #f8fafc;
      border-radius: 0.75rem;
      padding: 1rem;
      overflow-x: auto;
    }

    .trend-chart {
      display: flex;
      align-items: flex-end;
      gap: 0.25rem;
      height: 100px;
      min-width: max-content;
    }

    .trend-bar-group {
      display: flex;
      gap: 2px;
      align-items: flex-end;
    }

    .trend-bar {
      width: 8px;
      border-radius: 2px 2px 0 0;
      transition: height 0.3s ease;
      min-height: 2px;
    }

    .trend-bar-enquiry {
      background: linear-gradient(to top, #3b82f6, #60a5fa);
    }

    .trend-bar-app {
      background: linear-gradient(to top, #9333ea, #a855f7);
    }

    .trend-legend {
      display: flex;
      justify-content: center;
      gap: 1.5rem;
      padding-top: 0.5rem;
    }

    .legend-item {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .legend-dot {
      width: 0.5rem;
      height: 0.5rem;
      border-radius: 50%;
    }

    .bg-blue { background: #3b82f6; }
    .bg-purple { background: #9333ea; }

    /* Responsive Adjustments */
    @media (max-width: 640px) {
      .header-text h1 {
        font-size: 1.25rem;
      }

      .header-text p {
        display: none;
      }

      .stat-card {
        padding: 1rem;
      }

      .stat-icon-bg {
        width: 2.5rem;
        height: 2.5rem;
        font-size: 1rem;
      }

      .stat-value {
        font-size: 1.5rem;
      }

      .card-body {
        padding: 1rem;
        min-height: 200px;
      }

      .trend-summary {
        grid-template-columns: 1fr;
        gap: 0.5rem;
      }

      .trend-stat {
        display: flex;
        justify-content: space-between;
        align-items: center;
        text-align: left;
      }

      .trend-stat-value {
        font-size: 1.25rem;
      }

      .trend-stat-label {
        margin-top: 0;
      }
    }
  `],
})
export class AdmissionDashboardComponent implements OnInit {
  private readonly reportService = inject(AdmissionReportService);

  // State signals
  readonly initialLoading = signal(true);
  readonly loading = signal(false);
  readonly hasError = signal(false);
  readonly errorMessage = signal('');
  readonly classWiseLoading = signal(true);
  readonly funnelLoading = signal(true);
  readonly sourceLoading = signal(true);
  readonly trendLoading = signal(true);
  readonly stats = signal<DashboardStats | null>(null);
  readonly funnelData = signal<FunnelStage[]>([]);
  readonly classWiseData = signal<ClassWiseReport[]>([]);
  readonly sourceData = signal<SourceAnalysis[]>([]);
  readonly trendData = signal<DailyTrendPoint[]>([]);
  readonly sessions = signal<AdmissionSessionOption[]>([]);
  readonly selectedSessionId = signal<string | null>(null);

  // Filter state
  readonly filters = signal<ReportFilterParams>({});

  // Computed values
  readonly sessionOptions = computed<SelectOption[]>(() => {
    return this.sessions().map(session => ({
      value: session.id,
      label: `${session.name} (${session.status})`,
    }));
  });

  readonly overallConversion = computed(() => {
    const statsData = this.stats();
    if (!statsData || statsData.totalEnquiries === 0) return 0;
    return (statsData.enrolled / statsData.totalEnquiries) * 100;
  });

  readonly trendTotals = computed(() => {
    const data = this.trendData();
    return {
      enquiries: data.reduce((sum, d) => sum + d.enquiries, 0),
      applications: data.reduce((sum, d) => sum + d.applications, 0),
      enrollments: data.reduce((sum, d) => sum + d.enrollments, 0),
    };
  });

  readonly trendMax = computed(() => {
    const data = this.trendData();
    const maxEnquiries = Math.max(...data.map(d => d.enquiries), 1);
    const maxApplications = Math.max(...data.map(d => d.applications), 1);
    return Math.max(maxEnquiries, maxApplications);
  });

  readonly trendDataSubset = computed(() => {
    return this.trendData().slice(-20);
  });

  ngOnInit(): void {
    this.loadSessions();
    this.loadDashboardData();
  }

  /**
   * Load available sessions for filter
   */
  private loadSessions(): void {
    this.reportService.getSessions().subscribe({
      next: (sessions) => {
        this.sessions.set(sessions);
        const activeSession = sessions.find(s => s.status === 'open');
        if (activeSession) {
          this.selectedSessionId.set(activeSession.id);
        }
      },
      error: (err) => console.error('Failed to load sessions:', err),
    });
  }

  /**
   * Load all dashboard data
   */
  private loadDashboardData(): void {
    this.loading.set(true);
    this.hasError.set(false);
    this.classWiseLoading.set(true);
    this.funnelLoading.set(true);
    this.sourceLoading.set(true);
    this.trendLoading.set(true);

    const currentFilters = this.filters();

    // Load dashboard stats
    this.reportService.getDashboardStats(currentFilters).subscribe({
      next: (stats) => {
        this.stats.set(stats);
        this.loading.set(false);
        this.initialLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load dashboard stats:', err);
        this.hasError.set(true);
        this.errorMessage.set('Failed to load dashboard statistics. Please try again.');
        this.loading.set(false);
        this.initialLoading.set(false);
      },
    });

    // Load funnel data
    this.reportService.getFunnelData(currentFilters).subscribe({
      next: (data) => {
        this.funnelData.set(data);
        this.funnelLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load funnel data:', err);
        this.funnelData.set([]);
        this.funnelLoading.set(false);
      },
    });

    // Load class-wise report
    this.reportService.getClassWiseReport(currentFilters).subscribe({
      next: (response) => {
        this.classWiseData.set(response.classes);
        this.classWiseLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load class-wise report:', err);
        this.classWiseData.set([]);
        this.classWiseLoading.set(false);
      },
    });

    // Load source analysis
    this.reportService.getSourceAnalysis(currentFilters).subscribe({
      next: (data) => {
        this.sourceData.set(data);
        this.sourceLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load source analysis:', err);
        this.sourceData.set([]);
        this.sourceLoading.set(false);
      },
    });

    // Load daily trend
    this.reportService.getDailyTrend(currentFilters).subscribe({
      next: (data) => {
        this.trendData.set(data);
        this.trendLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load daily trend:', err);
        this.trendData.set([]);
        this.trendLoading.set(false);
      },
    });
  }

  /**
   * Handle session filter change
   */
  onSessionChange(value: string | number | (string | number)[] | null): void {
    const sessionId = value as string | null;
    this.selectedSessionId.set(sessionId);
    this.filters.update(f => ({ ...f, sessionId: sessionId || undefined }));
    this.loadDashboardData();
  }

  /**
   * Refresh all dashboard data
   */
  refreshData(): void {
    this.loadDashboardData();
  }

  /**
   * Calculate bar height for trend visualization
   */
  getBarHeight(value: number, maxValue: number): number {
    if (maxValue === 0) return 2;
    const maxHeight = 80;
    return Math.max(2, (value / maxValue) * maxHeight);
  }
}
