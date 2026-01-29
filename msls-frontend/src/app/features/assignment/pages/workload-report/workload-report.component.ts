/**
 * Workload Report Component
 * Story 5.7: Teacher Subject Assignment
 *
 * Displays teacher workload summary and identifies over/under-assigned teachers.
 */

import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AssignmentService } from '../../assignment.service';
import {
  WorkloadSummary,
  WorkloadReportResponse,
  UnassignedSubjectsResponse,
  getWorkloadStatusClass,
  getWorkloadStatusLabel,
} from '../../assignment.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-workload-report',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <button class="back-btn" routerLink="/assignments">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-icon">
            <i class="fa-solid fa-chart-bar"></i>
          </div>
          <div class="header-text">
            <h1>Workload Report</h1>
            <p>Analyze teacher workload distribution and identify imbalances</p>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <label class="filter-label">Academic Year</label>
          <select
            class="filter-select"
            [ngModel]="academicYearFilter()"
            (ngModelChange)="academicYearFilter.set($event); loadReport()"
          >
            @for (year of academicYears(); track year.id) {
              <option [value]="year.id">{{ year.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Branch</label>
          <select
            class="filter-select"
            [ngModel]="branchFilter()"
            (ngModelChange)="branchFilter.set($event); loadReport()"
          >
            <option value="">All Branches</option>
            @for (branch of branches(); track branch.id) {
              <option [value]="branch.id">{{ branch.name }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Summary Cards -->
      @if (!loading() && report()) {
        <div class="summary-grid">
          <div class="summary-card">
            <div class="summary-icon total">
              <i class="fa-solid fa-users"></i>
            </div>
            <div class="summary-content">
              <span class="summary-value">{{ report()?.totalTeachers }}</span>
              <span class="summary-label">Total Teachers</span>
            </div>
          </div>
          <div class="summary-card">
            <div class="summary-icon normal">
              <i class="fa-solid fa-check-circle"></i>
            </div>
            <div class="summary-content">
              <span class="summary-value">{{ report()?.normalAssigned }}</span>
              <span class="summary-label">Normal Workload</span>
            </div>
          </div>
          <div class="summary-card">
            <div class="summary-icon under">
              <i class="fa-solid fa-arrow-down"></i>
            </div>
            <div class="summary-content">
              <span class="summary-value">{{ report()?.underAssigned }}</span>
              <span class="summary-label">Under-assigned</span>
            </div>
          </div>
          <div class="summary-card">
            <div class="summary-icon over">
              <i class="fa-solid fa-arrow-up"></i>
            </div>
            <div class="summary-content">
              <span class="summary-value">{{ report()?.overAssigned }}</span>
              <span class="summary-label">Over-assigned</span>
            </div>
          </div>
        </div>
      }

      <!-- Tabs -->
      <div class="tabs">
        <button
          class="tab"
          [class.active]="activeTab() === 'teachers'"
          (click)="activeTab.set('teachers')"
        >
          <i class="fa-solid fa-chalkboard-user"></i>
          Teacher Workload
        </button>
        <button
          class="tab"
          [class.active]="activeTab() === 'unassigned'"
          (click)="activeTab.set('unassigned'); loadUnassignedSubjects()"
        >
          <i class="fa-solid fa-exclamation-triangle"></i>
          Unassigned Subjects
          @if (unassignedCount() > 0) {
            <span class="tab-badge">{{ unassignedCount() }}</span>
          }
        </button>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading report...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadReport()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          @if (activeTab() === 'teachers') {
            <table class="data-table">
              <thead>
                <tr>
                  <th>Teacher</th>
                  <th>Department</th>
                  <th style="text-align: center;">Periods</th>
                  <th style="text-align: center;">Subjects</th>
                  <th style="text-align: center;">Classes</th>
                  <th>Class Teacher</th>
                  <th>Status</th>
                  <th style="width: 200px;">Workload</th>
                </tr>
              </thead>
              <tbody>
                @for (teacher of report()?.teachers || []; track teacher.staffId) {
                  <tr [class.warning-row]="teacher.workloadStatus !== 'normal'">
                    <td class="teacher-cell">
                      <div class="teacher-wrapper">
                        <div class="teacher-avatar">
                          {{ getInitials(teacher.staffName) }}
                        </div>
                        <div class="teacher-info">
                          <span class="teacher-name">{{ teacher.staffName }}</span>
                          <span class="teacher-id">{{ teacher.staffEmployeeId }}</span>
                        </div>
                      </div>
                    </td>
                    <td>{{ teacher.departmentName || '-' }}</td>
                    <td class="count-cell">
                      <span class="count-value" [class.warning]="teacher.totalPeriods > teacher.maxPeriods">
                        {{ teacher.totalPeriods }}
                      </span>
                    </td>
                    <td class="count-cell">{{ teacher.totalSubjects }}</td>
                    <td class="count-cell">{{ teacher.totalClasses }}</td>
                    <td>
                      @if (teacher.isClassTeacher) {
                        <span class="badge badge-blue">
                          <i class="fa-solid fa-star"></i>
                          {{ teacher.classTeacherFor || 'Yes' }}
                        </span>
                      } @else {
                        <span class="text-muted">-</span>
                      }
                    </td>
                    <td>
                      <span class="badge" [class]="getWorkloadStatusClass(teacher.workloadStatus)">
                        {{ getWorkloadStatusLabel(teacher.workloadStatus) }}
                      </span>
                    </td>
                    <td>
                      <div class="workload-bar">
                        <div
                          class="workload-fill"
                          [class.over]="teacher.totalPeriods > teacher.maxPeriods"
                          [class.under]="teacher.totalPeriods < teacher.minPeriods"
                          [style.width.%]="getWorkloadPercentage(teacher)"
                        ></div>
                        <span class="workload-label">
                          {{ teacher.totalPeriods }}/{{ teacher.maxPeriods }}
                        </span>
                      </div>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="8" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No teachers found</p>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          } @else {
            <!-- Unassigned Subjects Tab -->
            @if (loadingUnassigned()) {
              <div class="loading-container">
                <div class="spinner"></div>
                <span>Loading unassigned subjects...</span>
              </div>
            } @else {
              <table class="data-table">
                <thead>
                  <tr>
                    <th>Subject</th>
                    <th>Code</th>
                    <th>Class</th>
                    <th>Section</th>
                    <th style="width: 120px; text-align: right;">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  @for (subject of unassignedSubjects()?.subjects || []; track subject.subjectId + subject.classId + (subject.sectionId || '')) {
                    <tr>
                      <td>
                        <span class="subject-name">{{ subject.subjectName }}</span>
                      </td>
                      <td>
                        <span class="subject-code">{{ subject.subjectCode }}</span>
                      </td>
                      <td>{{ subject.className }}</td>
                      <td>{{ subject.sectionName || 'All Sections' }}</td>
                      <td class="actions-cell">
                        <button
                          class="btn btn-primary btn-sm"
                          routerLink="/assignments/new"
                          [queryParams]="{
                            subjectId: subject.subjectId,
                            classId: subject.classId,
                            sectionId: subject.sectionId
                          }"
                        >
                          <i class="fa-solid fa-plus"></i>
                          Assign
                        </button>
                      </td>
                    </tr>
                  } @empty {
                    <tr>
                      <td colspan="5" class="empty-cell">
                        <div class="empty-state empty-state--success">
                          <i class="fa-solid fa-check-circle"></i>
                          <p>All subjects are assigned!</p>
                        </div>
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            }
          }
        }
      </div>
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
      background: white;
      color: #64748b;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.2s;
    }

    .back-btn:hover {
      background: #f1f5f9;
      color: #1e293b;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: #fef3c7;
      color: #d97706;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .filter-group {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .filter-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
    }

    .filter-select {
      padding: 0.5rem 2rem 0.5rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
      min-width: 160px;
    }

    .filter-select:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .summary-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .summary-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.25rem;
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .summary-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .summary-icon.total {
      background: #dbeafe;
      color: #2563eb;
    }

    .summary-icon.normal {
      background: #dcfce7;
      color: #16a34a;
    }

    .summary-icon.under {
      background: #fef3c7;
      color: #d97706;
    }

    .summary-icon.over {
      background: #fef2f2;
      color: #dc2626;
    }

    .summary-content {
      display: flex;
      flex-direction: column;
    }

    .summary-value {
      font-size: 1.5rem;
      font-weight: 700;
      color: #1e293b;
    }

    .summary-label {
      font-size: 0.75rem;
      color: #64748b;
    }

    .tabs {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }

    .tab {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.25rem;
      border: none;
      background: transparent;
      color: #64748b;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      border-radius: 0.5rem;
      transition: all 0.2s;
    }

    .tab:hover {
      background: #f1f5f9;
    }

    .tab.active {
      background: #4f46e5;
      color: white;
    }

    .tab-badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 1.25rem;
      height: 1.25rem;
      padding: 0 0.375rem;
      background: #dc2626;
      color: white;
      border-radius: 9999px;
      font-size: 0.625rem;
      font-weight: 600;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .loading-container,
    .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table th {
      text-align: left;
      padding: 0.875rem 1rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #374151;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .warning-row {
      background: #fffbeb;
    }

    .warning-row:hover {
      background: #fef3c7 !important;
    }

    .teacher-cell {
      min-width: 200px;
    }

    .teacher-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .teacher-avatar {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 50%;
      background: linear-gradient(135deg, #4f46e5, #7c3aed);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 600;
      flex-shrink: 0;
    }

    .teacher-info {
      display: flex;
      flex-direction: column;
    }

    .teacher-name {
      font-weight: 500;
      color: #1e293b;
    }

    .teacher-id {
      font-size: 0.75rem;
      color: #64748b;
    }

    .count-cell {
      text-align: center;
    }

    .count-value {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      height: 1.75rem;
      padding: 0 0.5rem;
      background: #f1f5f9;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #475569;
    }

    .count-value.warning {
      background: #fef2f2;
      color: #dc2626;
    }

    .text-muted {
      color: #94a3b8;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge i {
      font-size: 0.625rem;
    }

    .badge-gray { background: #f1f5f9; color: #64748b; }
    .badge-blue { background: #dbeafe; color: #1e40af; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-yellow { background: #fef3c7; color: #92400e; }
    .badge-red { background: #fef2f2; color: #991b1b; }

    .workload-bar {
      position: relative;
      height: 1.5rem;
      background: #f1f5f9;
      border-radius: 0.375rem;
      overflow: hidden;
    }

    .workload-fill {
      height: 100%;
      background: linear-gradient(90deg, #22c55e, #16a34a);
      border-radius: 0.375rem;
      transition: width 0.3s ease;
    }

    .workload-fill.under {
      background: linear-gradient(90deg, #fbbf24, #d97706);
    }

    .workload-fill.over {
      background: linear-gradient(90deg, #f87171, #dc2626);
    }

    .workload-label {
      position: absolute;
      inset: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 600;
      color: #1e293b;
    }

    .subject-name {
      font-weight: 500;
    }

    .subject-code {
      font-size: 0.75rem;
      color: #64748b;
      font-family: monospace;
    }

    .actions-cell {
      text-align: right;
    }

    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    .empty-state--success i {
      color: #22c55e;
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    @media (max-width: 1024px) {
      .summary-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    @media (max-width: 640px) {
      .summary-grid {
        grid-template-columns: 1fr;
      }

      .filters-bar {
        flex-direction: column;
      }

      .tabs {
        flex-direction: column;
      }
    }
  `],
})
export class WorkloadReportComponent implements OnInit {
  private assignmentService = inject(AssignmentService);
  private toastService = inject(ToastService);

  // State signals
  report = signal<WorkloadReportResponse | null>(null);
  unassignedSubjects = signal<UnassignedSubjectsResponse | null>(null);
  loading = signal(true);
  loadingUnassigned = signal(false);
  error = signal<string | null>(null);
  activeTab = signal<'teachers' | 'unassigned'>('teachers');

  // Filter signals
  academicYearFilter = signal<string>('');
  branchFilter = signal<string>('');

  // Reference data
  academicYears = signal<{ id: string; name: string }[]>([]);
  branches = signal<{ id: string; name: string }[]>([]);

  // Computed
  unassignedCount = computed(() => this.unassignedSubjects()?.total || 0);

  ngOnInit(): void {
    this.loadReferenceData();
    this.loadReport();
  }

  loadReferenceData(): void {
    // TODO: Load from API
    this.academicYears.set([
      { id: '1', name: '2024-25' },
      { id: '2', name: '2025-26' },
    ]);
    this.branches.set([
      { id: '1', name: 'Main Branch' },
      { id: '2', name: 'City Branch' },
    ]);

    if (this.academicYears().length > 0) {
      this.academicYearFilter.set(this.academicYears()[0].id);
    }
  }

  loadReport(): void {
    const academicYearId = this.academicYearFilter();
    if (!academicYearId) return;

    this.loading.set(true);
    this.error.set(null);

    this.assignmentService
      .getWorkloadReport(academicYearId, this.branchFilter() || undefined)
      .subscribe({
        next: report => {
          this.report.set(report);
          this.loading.set(false);
        },
        error: () => {
          this.error.set('Failed to load workload report. Please try again.');
          this.loading.set(false);
        },
      });
  }

  loadUnassignedSubjects(): void {
    if (this.unassignedSubjects()) return; // Already loaded

    const academicYearId = this.academicYearFilter();
    if (!academicYearId) return;

    this.loadingUnassigned.set(true);

    this.assignmentService.getUnassignedSubjects(academicYearId).subscribe({
      next: subjects => {
        this.unassignedSubjects.set(subjects);
        this.loadingUnassigned.set(false);
      },
      error: () => {
        this.toastService.error('Failed to load unassigned subjects');
        this.loadingUnassigned.set(false);
      },
    });
  }

  getWorkloadPercentage(teacher: WorkloadSummary): number {
    if (teacher.maxPeriods === 0) return 0;
    return Math.min(100, (teacher.totalPeriods / teacher.maxPeriods) * 100);
  }

  getInitials(name: string): string {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .substring(0, 2)
      .toUpperCase();
  }

  getWorkloadStatusClass = getWorkloadStatusClass;
  getWorkloadStatusLabel = getWorkloadStatusLabel;
}
