/**
 * MSLS Enrollment History Component
 *
 * Displays a timeline of student enrollment records.
 */

import { Component, input, output, computed, inject, OnInit, effect } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';

import { EnrollmentService } from '../../services/enrollment.service';
import {
  Enrollment,
  getEnrollmentStatusBadgeVariant,
  getEnrollmentStatusLabel,
  getEnrollmentTimelineColor,
} from '../../models/enrollment.model';

@Component({
  selector: 'app-enrollment-history',
  standalone: true,
  imports: [CommonModule, DatePipe],
  template: `
    <div class="enrollment-history">
      <!-- Action Bar -->
      @if (canCreate()) {
        <div class="flex justify-end mb-4">
          <button
            type="button"
            class="enroll-btn enroll-btn--primary"
            (click)="onCreateClick()"
          >
            <i class="fa-solid fa-plus"></i>
            New Enrollment
          </button>
        </div>
      }

      <!-- Loading State -->
      @if (service.loading()) {
        <div class="loading-state">
          <div class="spinner"></div>
          <p>Loading enrollment history...</p>
        </div>
      }

      <!-- Empty State -->
      @if (service.isEmpty() && !service.loading()) {
        <div class="empty-state">
          <div class="empty-state__icon">
            <i class="fa-solid fa-graduation-cap"></i>
          </div>
          <h3>No enrollment records</h3>
          <p>This student has no enrollment history yet.</p>
          @if (canCreate()) {
            <button
              type="button"
              class="enroll-btn enroll-btn--primary"
              (click)="onCreateClick()"
            >
              Create First Enrollment
            </button>
          }
        </div>
      }

      <!-- Timeline -->
      @if (!service.isEmpty() && !service.loading()) {
        <div class="timeline">
          @for (enrollment of service.enrollments(); track enrollment.id; let i = $index; let last = $last) {
            <div class="timeline-item" [class.timeline-item--last]="last">
              <!-- Timeline dot -->
              <div class="timeline-dot" [ngClass]="getTimelineDotClass(enrollment.status)">
                <i [class]="getStatusIcon(enrollment.status)"></i>
              </div>

              <!-- Enrollment card -->
              <div
                class="enrollment-card"
                [class.enrollment-card--selected]="isSelected(enrollment)"
                [class.enrollment-card--active]="enrollment.status === 'active'"
                (click)="onEnrollmentClick(enrollment)"
              >
                <div class="enrollment-card__header">
                  <div class="enrollment-card__title">
                    <span class="year-badge">{{ enrollment.academicYear?.name || 'Unknown Year' }}</span>
                    @if (enrollment.status === 'active') {
                      <span class="current-badge">Current</span>
                    }
                  </div>
                  <span class="status-badge" [ngClass]="getStatusBadgeClasses(enrollment.status)">
                    {{ getStatusLabel(enrollment.status) }}
                  </span>
                </div>

                <div class="enrollment-card__body">
                  <!-- Class Info -->
                  @if (enrollment.classId || enrollment.sectionId || enrollment.rollNumber) {
                    <div class="info-row">
                      @if (enrollment.classId) {
                        <span class="info-tag">
                          <i class="fa-solid fa-chalkboard"></i>
                          {{ enrollment.classId }}
                        </span>
                      }
                      @if (enrollment.sectionId) {
                        <span class="info-tag">
                          <i class="fa-solid fa-users"></i>
                          {{ enrollment.sectionId }}
                        </span>
                      }
                      @if (enrollment.rollNumber) {
                        <span class="info-tag">
                          <i class="fa-solid fa-hashtag"></i>
                          {{ enrollment.rollNumber }}
                        </span>
                      }
                    </div>
                  }

                  <!-- Dates -->
                  <div class="dates-section">
                    <div class="date-item">
                      <i class="fa-solid fa-calendar-plus"></i>
                      <span>Enrolled: {{ enrollment.enrollmentDate | date:'mediumDate' }}</span>
                    </div>

                    @if (enrollment.completionDate) {
                      <div class="date-item date-item--success">
                        <i class="fa-solid fa-check-circle"></i>
                        <span>Completed: {{ enrollment.completionDate | date:'mediumDate' }}</span>
                      </div>
                    }

                    @if (enrollment.transferDate) {
                      <div class="date-item date-item--warning">
                        <i class="fa-solid fa-arrow-right-from-bracket"></i>
                        <span>Transferred: {{ enrollment.transferDate | date:'mediumDate' }}</span>
                      </div>
                    }

                    @if (enrollment.dropoutDate) {
                      <div class="date-item date-item--danger">
                        <i class="fa-solid fa-times-circle"></i>
                        <span>Dropped out: {{ enrollment.dropoutDate | date:'mediumDate' }}</span>
                      </div>
                    }
                  </div>

                  <!-- Reason -->
                  @if (enrollment.transferReason || enrollment.dropoutReason) {
                    <div class="reason-box">
                      <i class="fa-solid fa-quote-left"></i>
                      {{ enrollment.transferReason || enrollment.dropoutReason }}
                    </div>
                  }

                  <!-- Notes -->
                  @if (enrollment.notes) {
                    <div class="notes-box">
                      <i class="fa-solid fa-sticky-note"></i>
                      {{ enrollment.notes }}
                    </div>
                  }
                </div>

                <!-- Actions for active enrollment -->
                @if (enrollment.status === 'active' && showActions()) {
                  <div class="enrollment-card__actions">
                    <button
                      type="button"
                      class="action-btn action-btn--edit"
                      (click)="onEditClick(enrollment, $event)"
                    >
                      <i class="fa-solid fa-pen"></i>
                      Edit
                    </button>
                    <button
                      type="button"
                      class="action-btn action-btn--transfer"
                      (click)="onTransferClick(enrollment, $event)"
                    >
                      <i class="fa-solid fa-arrow-right-from-bracket"></i>
                      Transfer
                    </button>
                    <button
                      type="button"
                      class="action-btn action-btn--dropout"
                      (click)="onDropoutClick(enrollment, $event)"
                    >
                      <i class="fa-solid fa-times"></i>
                      Dropout
                    </button>
                  </div>
                }
              </div>
            </div>
          }
        </div>
      }

      <!-- Error State -->
      @if (service.error()) {
        <div class="error-state">
          <i class="fa-solid fa-exclamation-triangle"></i>
          <p>{{ service.error() }}</p>
        </div>
      }
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .enrollment-history {
      padding: 0.5rem 0;
    }

    /* Buttons */
    .enroll-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      font-size: 0.8125rem;
      font-weight: 600;
      border-radius: 0.5rem;
      border: none;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .enroll-btn--primary {
      background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
      color: white;
      box-shadow: 0 2px 8px rgba(139, 92, 246, 0.3);
    }

    .enroll-btn--primary:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
    }

    /* Loading State */
    .loading-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem 1rem;
      color: #64748b;
    }

    .spinner {
      width: 2.5rem;
      height: 2.5rem;
      border: 3px solid #e2e8f0;
      border-top-color: #8b5cf6;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
      margin-bottom: 1rem;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Empty State */
    .empty-state {
      text-align: center;
      padding: 3rem 1rem;
    }

    .empty-state__icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .empty-state__icon i {
      font-size: 1.5rem;
      color: #94a3b8;
    }

    .empty-state h3 {
      font-size: 0.9375rem;
      font-weight: 600;
      color: #1e293b;
      margin: 0 0 0.25rem;
    }

    .empty-state p {
      font-size: 0.8125rem;
      color: #64748b;
      margin: 0 0 1.5rem;
    }

    /* Timeline */
    .timeline {
      position: relative;
      padding-left: 2rem;
    }

    .timeline::before {
      content: '';
      position: absolute;
      left: 0.75rem;
      top: 0;
      bottom: 0;
      width: 2px;
      background: linear-gradient(to bottom, #e2e8f0, #f1f5f9);
      border-radius: 1px;
    }

    .timeline-item {
      position: relative;
      padding-bottom: 1.5rem;
    }

    .timeline-item--last {
      padding-bottom: 0;
    }

    .timeline-dot {
      position: absolute;
      left: -2rem;
      top: 0;
      width: 1.75rem;
      height: 1.75rem;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 3px solid white;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      z-index: 1;
    }

    .timeline-dot i {
      font-size: 0.625rem;
      color: white;
    }

    .timeline-dot--active {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    }

    .timeline-dot--completed {
      background: linear-gradient(135deg, #64748b 0%, #475569 100%);
    }

    .timeline-dot--transferred {
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    }

    .timeline-dot--dropout {
      background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    }

    /* Enrollment Card */
    .enrollment-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
      transition: all 0.2s ease;
      cursor: pointer;
    }

    .enrollment-card:hover {
      border-color: #cbd5e1;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .enrollment-card--selected {
      border-color: #8b5cf6;
      box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
    }

    .enrollment-card--active {
      border-left: 4px solid #10b981;
    }

    .enrollment-card__header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.875rem 1rem;
      background: linear-gradient(to right, #f8fafc, white);
      border-bottom: 1px solid #f1f5f9;
    }

    .enrollment-card__title {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .year-badge {
      font-size: 0.9375rem;
      font-weight: 600;
      color: #1e293b;
    }

    .current-badge {
      font-size: 0.625rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      padding: 0.125rem 0.5rem;
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      color: white;
      border-radius: 9999px;
    }

    .status-badge {
      font-size: 0.6875rem;
      font-weight: 600;
      padding: 0.25rem 0.625rem;
      border-radius: 9999px;
    }

    .status-badge--success {
      background: #dcfce7;
      color: #166534;
    }

    .status-badge--warning {
      background: #fef3c7;
      color: #92400e;
    }

    .status-badge--error {
      background: #fee2e2;
      color: #991b1b;
    }

    .status-badge--neutral {
      background: #f1f5f9;
      color: #475569;
    }

    .enrollment-card__body {
      padding: 1rem;
    }

    .info-row {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 0.75rem;
    }

    .info-tag {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
      font-weight: 500;
      padding: 0.25rem 0.625rem;
      background: #f1f5f9;
      color: #475569;
      border-radius: 0.375rem;
    }

    .info-tag i {
      font-size: 0.625rem;
      color: #94a3b8;
    }

    .dates-section {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .date-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.8125rem;
      color: #64748b;
    }

    .date-item i {
      width: 1rem;
      text-align: center;
      font-size: 0.75rem;
    }

    .date-item--success {
      color: #059669;
    }

    .date-item--warning {
      color: #d97706;
    }

    .date-item--danger {
      color: #dc2626;
    }

    .reason-box, .notes-box {
      margin-top: 0.75rem;
      padding: 0.75rem;
      background: #fafafa;
      border-radius: 0.5rem;
      font-size: 0.8125rem;
      color: #64748b;
      display: flex;
      gap: 0.5rem;
      font-style: italic;
    }

    .reason-box i, .notes-box i {
      color: #94a3b8;
      font-size: 0.75rem;
      flex-shrink: 0;
      margin-top: 0.125rem;
    }

    .notes-box {
      background: #fffbeb;
    }

    .notes-box i {
      color: #f59e0b;
    }

    /* Actions */
    .enrollment-card__actions {
      display: flex;
      gap: 0.5rem;
      padding: 0.75rem 1rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
    }

    .action-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 0.375rem;
      border: 1px solid transparent;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .action-btn i {
      font-size: 0.625rem;
    }

    .action-btn--edit {
      background: white;
      border-color: #e2e8f0;
      color: #475569;
    }

    .action-btn--edit:hover {
      border-color: #8b5cf6;
      color: #8b5cf6;
      background: #f5f3ff;
    }

    .action-btn--transfer {
      background: white;
      border-color: #e2e8f0;
      color: #d97706;
    }

    .action-btn--transfer:hover {
      border-color: #f59e0b;
      background: #fffbeb;
    }

    .action-btn--dropout {
      background: white;
      border-color: #e2e8f0;
      color: #dc2626;
    }

    .action-btn--dropout:hover {
      border-color: #ef4444;
      background: #fef2f2;
    }

    /* Error State */
    .error-state {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem;
      background: #fef2f2;
      border: 1px solid #fecaca;
      border-radius: 0.5rem;
      color: #dc2626;
      font-size: 0.875rem;
    }

    .error-state i {
      font-size: 1rem;
    }
  `]
})
export class EnrollmentHistoryComponent implements OnInit {
  protected service = inject(EnrollmentService);

  /** Student ID */
  studentId = input.required<string>();

  /** Whether to show action buttons */
  showActions = input<boolean>(true);

  /** Whether user can create new enrollments */
  canCreate = input<boolean>(true);

  /** Event emitted when create button is clicked */
  create = output<void>();

  /** Event emitted when edit button is clicked */
  edit = output<Enrollment>();

  /** Event emitted when transfer button is clicked */
  transfer = output<Enrollment>();

  /** Event emitted when dropout button is clicked */
  dropout = output<Enrollment>();

  /** Event emitted when an enrollment is selected */
  enrollmentSelected = output<Enrollment>();

  constructor() {
    // Auto-load when studentId changes
    effect(() => {
      const id = this.studentId();
      if (id) {
        this.service.loadEnrollmentHistory(id).subscribe();
      }
    });
  }

  ngOnInit(): void {
    // Initial load handled by effect
  }

  // =========================================================================
  // Template Helpers
  // =========================================================================

  getTimelineColor(status: string): string {
    return getEnrollmentTimelineColor(status as any);
  }

  getStatusLabel(status: string): string {
    return getEnrollmentStatusLabel(status as any);
  }

  getStatusBadgeClasses(status: string): string {
    const variant = getEnrollmentStatusBadgeVariant(status as any);
    switch (variant) {
      case 'success':
        return 'status-badge--success';
      case 'warning':
        return 'status-badge--warning';
      case 'error':
        return 'status-badge--error';
      default:
        return 'status-badge--neutral';
    }
  }

  getTimelineDotClass(status: string): string {
    switch (status) {
      case 'active':
        return 'timeline-dot--active';
      case 'completed':
        return 'timeline-dot--completed';
      case 'transferred':
        return 'timeline-dot--transferred';
      case 'dropout':
        return 'timeline-dot--dropout';
      default:
        return 'timeline-dot--completed';
    }
  }

  getStatusIcon(status: string): string {
    switch (status) {
      case 'active':
        return 'fa-solid fa-play';
      case 'completed':
        return 'fa-solid fa-check';
      case 'transferred':
        return 'fa-solid fa-arrow-right';
      case 'dropout':
        return 'fa-solid fa-times';
      default:
        return 'fa-solid fa-circle';
    }
  }

  isSelected(enrollment: Enrollment): boolean {
    return this.service.selectedEnrollment()?.id === enrollment.id;
  }

  // =========================================================================
  // Event Handlers
  // =========================================================================

  onCreateClick(): void {
    this.create.emit();
  }

  onEnrollmentClick(enrollment: Enrollment): void {
    this.service.selectEnrollment(enrollment);
    this.enrollmentSelected.emit(enrollment);
  }

  onEditClick(enrollment: Enrollment, event: Event): void {
    event.stopPropagation();
    this.edit.emit(enrollment);
  }

  onTransferClick(enrollment: Enrollment, event: Event): void {
    event.stopPropagation();
    this.transfer.emit(enrollment);
  }

  onDropoutClick(enrollment: Enrollment, event: Event): void {
    event.stopPropagation();
    this.dropout.emit(enrollment);
  }
}
