/**
 * MSLS Application Review Component
 *
 * Main page for reviewing admission applications with document verification,
 * status updates, and review comments.
 */

import {
  Component,
  inject,
  signal,
  computed,
  OnInit,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsBadgeComponent } from '../../../shared/components/badge/badge.component';
import { MslsCardComponent } from '../../../shared/components/card/card.component';
import { MslsSpinnerComponent } from '../../../shared/components/spinner/spinner.component';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { ToastService } from '../../../shared/services/toast.service';

import { ApplicationReviewService } from './application-review.service';
import { DocumentVerificationComponent } from './document-verification.component';
import {
  AdmissionApplication,
  ApplicationStatus,
  APPLICATION_STATUS_CONFIG,
  getDocumentSummary,
  AddReviewDto,
} from './application-review.model';
import { EntranceTestService } from '../tests/entrance-test.service';
import { EntranceTest } from '../tests/entrance-test.model';

/**
 * ApplicationReviewComponent - Application review page with document verification.
 */
@Component({
  selector: 'msls-application-review',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsButtonComponent,
    MslsBadgeComponent,
    MslsCardComponent,
    MslsSpinnerComponent,
    MslsModalComponent,
    DocumentVerificationComponent,
  ],
  template: `
    <div class="application-review">
      <!-- Header -->
      <div class="page-header">
        <div class="header-left">
          <button class="back-btn" (click)="goBack()">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-info">
            <h1 class="page-title">Application Review</h1>
            @if (application()) {
              <span class="app-number">{{ application()?.applicationNumber }}</span>
            }
          </div>
        </div>
        @if (application()) {
          <div class="header-actions">
            <msls-button variant="outline" (click)="openReviewModal()">
              <i class="fa-solid fa-comment"></i>
              Add Review
            </msls-button>
            <msls-button
              variant="danger"
              [disabled]="application()?.status === 'rejected'"
              (click)="openRejectModal()"
            >
              <i class="fa-solid fa-times"></i>
              Reject
            </msls-button>
            <msls-button
              variant="primary"
              [disabled]="!canApprove()"
              (click)="approveApplication()"
            >
              <i class="fa-solid fa-check"></i>
              Approve
            </msls-button>
          </div>
        }
      </div>

      @if (loading()) {
        <div class="loading-container">
          <msls-spinner size="lg" />
          <p>Loading application...</p>
        </div>
      } @else if (application()) {
        <div class="review-content">
          <!-- Left Column - Application Details -->
          <div class="details-column">
            <!-- Status Card -->
            <msls-card variant="elevated" padding="md">
              <div card-header>
                <div class="status-header">
                  <span>Application Status</span>
                  <msls-badge [variant]="getStatusConfig(application()!.status).variant">
                    <i [class]="getStatusConfig(application()!.status).icon"></i>
                    {{ getStatusConfig(application()!.status).label }}
                  </msls-badge>
                </div>
              </div>
              <div card-body>
                <p class="status-description">
                  {{ getStatusConfig(application()!.status).description }}
                </p>
                <div class="status-actions">
                  @if (application()?.status === 'documents_verified') {
                    <msls-button variant="primary" size="sm" (click)="openScheduleTestModal()">
                      <i class="fa-solid fa-calendar-plus"></i>
                      Schedule Test
                    </msls-button>
                  }
                </div>
              </div>
            </msls-card>

            <!-- Student Details -->
            <msls-card variant="default" padding="md">
              <div card-header>
                <i class="fa-solid fa-user section-icon"></i>
                Student Information
              </div>
              <div card-body>
                <div class="info-grid">
                  <div class="info-item">
                    <label>Full Name</label>
                    <span>{{ application()?.studentName }}</span>
                  </div>
                  <div class="info-item">
                    <label>Date of Birth</label>
                    <span>{{ formatDate(application()?.dateOfBirth) }}</span>
                  </div>
                  <div class="info-item">
                    <label>Gender</label>
                    <span>{{ application()?.gender || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>Blood Group</label>
                    <span>{{ application()?.bloodGroup || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>Class Applying For</label>
                    <span class="highlight">{{ application()?.classApplying }}</span>
                  </div>
                  <div class="info-item">
                    <label>Academic Year</label>
                    <span>{{ application()?.academicYear }}</span>
                  </div>
                  @if (application()?.previousSchool) {
                    <div class="info-item full-width">
                      <label>Previous School</label>
                      <span>{{ application()?.previousSchool }} ({{ application()?.previousClass }})</span>
                    </div>
                  }
                </div>
              </div>
            </msls-card>

            <!-- Parent/Guardian Details -->
            <msls-card variant="default" padding="md">
              <div card-header>
                <i class="fa-solid fa-users section-icon"></i>
                Parent/Guardian Information
              </div>
              <div card-body>
                <div class="info-grid">
                  <div class="info-item">
                    <label>Father's Name</label>
                    <span>{{ application()?.fatherName || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>Mother's Name</label>
                    <span>{{ application()?.motherName || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>Contact Number</label>
                    <span>
                      <a href="tel:{{ application()?.parentPhone }}" class="contact-link">
                        <i class="fa-solid fa-phone"></i>
                        {{ application()?.parentPhone }}
                      </a>
                    </span>
                  </div>
                  <div class="info-item">
                    <label>Email</label>
                    <span>
                      @if (application()?.parentEmail) {
                        <a href="mailto:{{ application()?.parentEmail }}" class="contact-link">
                          <i class="fa-solid fa-envelope"></i>
                          {{ application()?.parentEmail }}
                        </a>
                      } @else {
                        -
                      }
                    </span>
                  </div>
                  <div class="info-item full-width">
                    <label>Address</label>
                    <span>
                      {{ application()?.address }}
                      @if (application()?.city) {
                        , {{ application()?.city }}
                      }
                      @if (application()?.state) {
                        , {{ application()?.state }}
                      }
                      @if (application()?.pincode) {
                        - {{ application()?.pincode }}
                      }
                    </span>
                  </div>
                </div>
              </div>
            </msls-card>

            <!-- Review History -->
            <msls-card variant="default" padding="md">
              <div card-header>
                <i class="fa-solid fa-history section-icon"></i>
                Review History
              </div>
              <div card-body>
                @if (application()?.reviews?.length) {
                  <div class="reviews-list">
                    @for (review of application()?.reviews; track review.id) {
                      <div class="review-item">
                        <div class="review-header">
                          <span class="reviewer">{{ review.reviewerName }}</span>
                          <span class="review-date">{{ formatDateTime(review.createdAt) }}</span>
                        </div>
                        <div class="review-type">
                          <msls-badge
                            [variant]="review.status === 'approved' ? 'success' : review.status === 'rejected' ? 'danger' : 'warning'"
                            size="sm"
                          >
                            {{ getReviewTypeLabel(review.reviewType) }}
                          </msls-badge>
                        </div>
                        @if (review.comments) {
                          <p class="review-comments">{{ review.comments }}</p>
                        }
                      </div>
                    }
                  </div>
                } @else {
                  <div class="empty-reviews">
                    <i class="fa-regular fa-comment-dots"></i>
                    <p>No reviews yet</p>
                  </div>
                }
              </div>
            </msls-card>
          </div>

          <!-- Right Column - Documents -->
          <div class="documents-column">
            <msls-card variant="elevated" padding="md">
              <div card-header>
                <div class="docs-header">
                  <span>
                    <i class="fa-solid fa-folder-open section-icon"></i>
                    Documents
                  </span>
                  <div class="docs-summary">
                    <span class="doc-stat verified">
                      <i class="fa-solid fa-check-circle"></i>
                      {{ docSummary().verified }}
                    </span>
                    <span class="doc-stat pending">
                      <i class="fa-solid fa-clock"></i>
                      {{ docSummary().pending }}
                    </span>
                    <span class="doc-stat rejected">
                      <i class="fa-solid fa-times-circle"></i>
                      {{ docSummary().rejected }}
                    </span>
                  </div>
                </div>
              </div>
              <div card-body>
                <msls-document-verification
                  [documents]="application()?.documents || []"
                  [applicationId]="application()?.id || ''"
                  (documentVerified)="onDocumentVerified($event)"
                />
              </div>
            </msls-card>
          </div>
        </div>
      } @else {
        <div class="error-container">
          <i class="fa-solid fa-exclamation-circle"></i>
          <p>Application not found</p>
          <msls-button variant="primary" (click)="goBack()">Go Back</msls-button>
        </div>
      }

      <!-- Add Review Modal -->
      <msls-modal
        [isOpen]="reviewModalOpen"
        (closed)="reviewModalOpen = false"
        title="Add Review"
        size="md"
      >
        <ng-container modal-body>
          <div class="review-form">
            <div class="form-group">
              <label>Review Type</label>
              <select class="form-select" [(ngModel)]="reviewForm.reviewType">
                <option value="document_verification">Document Verification</option>
                <option value="academic_review">Academic Review</option>
                <option value="interview">Interview</option>
                <option value="final_decision">Final Decision</option>
              </select>
            </div>
            <div class="form-group">
              <label>Status</label>
              <select class="form-select" [(ngModel)]="reviewForm.status">
                <option value="pending">Pending</option>
                <option value="approved">Approved</option>
                <option value="rejected">Rejected</option>
              </select>
            </div>
            <div class="form-group">
              <label>Comments</label>
              <textarea
                class="form-textarea"
                rows="4"
                placeholder="Enter your review comments..."
                [(ngModel)]="reviewForm.comments"
              ></textarea>
            </div>
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="modal-actions">
            <msls-button variant="ghost" (click)="reviewModalOpen = false">Cancel</msls-button>
            <msls-button variant="primary" [loading]="submitting()" (click)="submitReview()">
              Submit Review
            </msls-button>
          </div>
        </ng-container>
      </msls-modal>

      <!-- Reject Modal -->
      <msls-modal
        [isOpen]="rejectModalOpen"
        (closed)="rejectModalOpen = false"
        title="Reject Application"
        size="md"
      >
        <ng-container modal-body>
          <div class="reject-warning">
            <i class="fa-solid fa-exclamation-triangle"></i>
            <p>Are you sure you want to reject this application?</p>
          </div>
          <div class="form-group">
            <label>Reason for Rejection</label>
            <textarea
              class="form-textarea"
              rows="4"
              placeholder="Enter the reason for rejection..."
              [(ngModel)]="rejectReason"
            ></textarea>
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="modal-actions">
            <msls-button variant="ghost" (click)="rejectModalOpen = false">Cancel</msls-button>
            <msls-button
              variant="danger"
              [loading]="submitting()"
              [disabled]="!rejectReason.trim()"
              (click)="confirmReject()"
            >
              Reject Application
            </msls-button>
          </div>
        </ng-container>
      </msls-modal>

      <!-- Schedule Test Modal -->
      <msls-modal
        [isOpen]="scheduleTestModal"
        (closed)="scheduleTestModal = false"
        title="Schedule Entrance Test"
        size="md"
      >
        <ng-container modal-body>
          @if (loadingTests()) {
            <div class="loading-tests">
              <div class="spinner"></div>
              <p>Loading available tests...</p>
            </div>
          } @else if (availableTests().length === 0) {
            <div class="no-tests">
              <i class="fa-solid fa-calendar-xmark"></i>
              <p>No available tests found</p>
              <small>There are no scheduled tests for {{ application()?.classApplying }} with available seats.</small>
            </div>
          } @else {
            <p class="schedule-info">
              Select an entrance test to schedule for this applicant.
            </p>
            <div class="form-group">
              <label>Available Tests</label>
              <select class="form-select" [(ngModel)]="selectedTestId">
                <option value="">Select a test...</option>
                @for (test of availableTests(); track test.id) {
                  <option [value]="test.id">
                    {{ test.testName }} - {{ formatTestDate(test.testDate) }} {{ formatTestTime(test.startTime) }}
                    ({{ test.registeredCount || 0 }}/{{ test.maxCandidates }} seats)
                  </option>
                }
              </select>
            </div>
          }
        </ng-container>
        <ng-container modal-footer>
          <div class="modal-actions">
            <msls-button variant="ghost" (click)="scheduleTestModal = false">Cancel</msls-button>
            <msls-button
              variant="primary"
              [loading]="submitting()"
              [disabled]="!selectedTestId || loadingTests()"
              (click)="confirmScheduleTest()"
            >
              Schedule Test
            </msls-button>
          </div>
        </ng-container>
      </msls-modal>
    </div>
  `,
  styles: [`
    .application-review {
      min-height: 100%;
      background: #f8fafc;
      padding: 1.5rem;
    }

    .page-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 1.5rem;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .header-left {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
      color: #64748b;
    }

    .back-btn:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
      color: #334155;
    }

    .page-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .app-number {
      font-size: 0.875rem;
      color: #64748b;
      font-family: monospace;
    }

    .header-actions {
      display: flex;
      gap: 0.75rem;
      flex-wrap: wrap;
    }

    .loading-container,
    .error-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-height: 400px;
      gap: 1rem;
      color: #64748b;
    }

    .error-container i {
      font-size: 3rem;
      color: #ef4444;
    }

    .review-content {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }

    @media (min-width: 1024px) {
      .review-content {
        grid-template-columns: 1fr 1fr;
      }
    }

    @media (min-width: 1280px) {
      .review-content {
        grid-template-columns: 3fr 2fr;
      }
    }

    .details-column,
    .documents-column {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .section-icon {
      color: #4f46e5;
      margin-right: 0.5rem;
    }

    .status-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
    }

    .status-description {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0 0 1rem 0;
    }

    .status-actions {
      display: flex;
      gap: 0.5rem;
    }

    .info-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .info-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .info-item.full-width {
      grid-column: span 2;
    }

    .info-item label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .info-item span {
      font-size: 0.875rem;
      color: #0f172a;
    }

    .info-item .highlight {
      font-weight: 600;
      color: #4f46e5;
    }

    .contact-link {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      color: #4f46e5;
      text-decoration: none;
    }

    .contact-link:hover {
      text-decoration: underline;
    }

    .docs-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
    }

    .docs-summary {
      display: flex;
      gap: 0.75rem;
    }

    .doc-stat {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      font-size: 0.875rem;
      font-weight: 500;
    }

    .doc-stat.verified {
      color: #10b981;
    }

    .doc-stat.pending {
      color: #f59e0b;
    }

    .doc-stat.rejected {
      color: #ef4444;
    }

    .reviews-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .review-item {
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      border-left: 3px solid #4f46e5;
    }

    .review-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }

    .reviewer {
      font-weight: 600;
      color: #0f172a;
      font-size: 0.875rem;
    }

    .review-date {
      font-size: 0.75rem;
      color: #94a3b8;
    }

    .review-type {
      margin-bottom: 0.5rem;
    }

    .review-comments {
      font-size: 0.875rem;
      color: #475569;
      margin: 0;
      line-height: 1.5;
    }

    .empty-reviews {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 2rem;
      color: #94a3b8;
    }

    .empty-reviews i {
      font-size: 2rem;
      margin-bottom: 0.5rem;
    }

    /* Modal Styles */
    .review-form,
    .form-group {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .form-group label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .form-select,
    .form-textarea {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #0f172a;
      transition: all 0.2s;
    }

    .form-select:focus,
    .form-textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .form-textarea {
      resize: vertical;
      min-height: 100px;
    }

    .modal-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
    }

    .reject-warning {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem;
      background: #fef2f2;
      border-radius: 0.5rem;
      margin-bottom: 1rem;
    }

    .reject-warning i {
      color: #ef4444;
      font-size: 1.5rem;
    }

    .reject-warning p {
      color: #991b1b;
      margin: 0;
      font-weight: 500;
    }

    .schedule-info {
      color: #64748b;
      margin: 0 0 1rem 0;
    }

    .loading-tests,
    .no-tests {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 2rem;
      gap: 0.75rem;
      color: #64748b;
      text-align: center;
    }

    .loading-tests .spinner {
      width: 2rem;
      height: 2rem;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .no-tests i {
      font-size: 2.5rem;
      color: #cbd5e1;
    }

    .no-tests small {
      color: #94a3b8;
      font-size: 0.75rem;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ApplicationReviewComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly reviewService = inject(ApplicationReviewService);
  private readonly testService = inject(EntranceTestService);
  private readonly toast = inject(ToastService);

  readonly loading = this.reviewService.loading;
  readonly application = this.reviewService.selectedApplication;
  readonly submitting = signal(false);

  // Available tests for scheduling
  availableTests = signal<EntranceTest[]>([]);
  loadingTests = signal(false);

  // Document summary computed
  readonly docSummary = computed(() => {
    const docs = this.application()?.documents || [];
    return getDocumentSummary(docs);
  });

  // Modal states
  reviewModalOpen = false;
  rejectModalOpen = false;
  scheduleTestModal = false;

  // Form data
  reviewForm: AddReviewDto = {
    reviewType: 'document_verification',
    status: 'pending',
    comments: '',
  };
  rejectReason = '';
  selectedTestId = '';

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadApplication(id);
    }
  }

  private loadApplication(id: string): void {
    this.reviewService.getApplication(id).subscribe({
      error: (err) => {
        this.toast.error('Failed to load application');
        console.error('Error loading application:', err);
      },
    });
  }

  goBack(): void {
    this.router.navigate(['/admissions/applications']);
  }

  getStatusConfig(status: ApplicationStatus) {
    return APPLICATION_STATUS_CONFIG[status] || APPLICATION_STATUS_CONFIG.submitted;
  }

  getReviewTypeLabel(type: string): string {
    const labels: Record<string, string> = {
      document_verification: 'Document Verification',
      academic_review: 'Academic Review',
      interview: 'Interview',
      final_decision: 'Final Decision',
    };
    return labels[type] || type;
  }

  formatDate(dateString?: string): string {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  formatDateTime(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  canApprove(): boolean {
    const app = this.application();
    if (!app) return false;
    // Can approve if documents are verified or test is completed
    return ['documents_verified', 'test_completed'].includes(app.status);
  }

  openReviewModal(): void {
    this.reviewForm = {
      reviewType: 'document_verification',
      status: 'pending',
      comments: '',
    };
    this.reviewModalOpen = true;
  }

  openRejectModal(): void {
    this.rejectReason = '';
    this.rejectModalOpen = true;
  }

  submitReview(): void {
    if (!this.application()) return;

    this.submitting.set(true);
    this.reviewService.addReview(this.application()!.id, this.reviewForm).subscribe({
      next: () => {
        this.toast.success('Review added successfully');
        this.reviewModalOpen = false;
        this.submitting.set(false);
      },
      error: (err) => {
        this.toast.error('Failed to add review');
        console.error('Error adding review:', err);
        this.submitting.set(false);
      },
    });
  }

  confirmReject(): void {
    if (!this.application() || !this.rejectReason.trim()) return;

    this.submitting.set(true);
    this.reviewService
      .updateStatus(this.application()!.id, {
        status: 'rejected',
        comments: this.rejectReason,
      })
      .subscribe({
        next: () => {
          this.toast.success('Application rejected');
          this.rejectModalOpen = false;
          this.submitting.set(false);
        },
        error: (err) => {
          this.toast.error('Failed to reject application');
          console.error('Error rejecting application:', err);
          this.submitting.set(false);
        },
      });
  }

  approveApplication(): void {
    if (!this.application()) return;

    this.submitting.set(true);
    this.reviewService
      .updateStatus(this.application()!.id, { status: 'approved' })
      .subscribe({
        next: () => {
          this.toast.success('Application approved');
          this.submitting.set(false);
        },
        error: (err) => {
          this.toast.error('Failed to approve application');
          console.error('Error approving application:', err);
          this.submitting.set(false);
        },
      });
  }

  openScheduleTestModal(): void {
    this.scheduleTestModal = true;
    this.selectedTestId = '';
    this.loadAvailableTests();
  }

  loadAvailableTests(): void {
    const app = this.application();
    if (!app) return;

    this.loadingTests.set(true);
    this.testService.getTests({ status: 'scheduled' }).subscribe({
      next: tests => {
        // Filter tests that match the class the student is applying for
        const filtered = tests.filter(t =>
          t.classNames.includes(app.classApplying) &&
          (t.registeredCount || 0) < t.maxCandidates
        );
        this.availableTests.set(filtered);
        this.loadingTests.set(false);
      },
      error: err => {
        console.error('Failed to load tests:', err);
        this.toast.error('Failed to load available tests');
        this.loadingTests.set(false);
      },
    });
  }

  confirmScheduleTest(): void {
    if (!this.application() || !this.selectedTestId) return;

    this.submitting.set(true);

    // Register the candidate for the test
    this.testService.registerCandidate(this.selectedTestId, {
      applicationId: this.application()!.id,
    }).subscribe({
      next: () => {
        // Update application status to test_scheduled
        this.reviewService.scheduleTest(this.application()!.id, this.selectedTestId).subscribe({
          next: () => {
            this.toast.success('Test scheduled successfully');
            this.scheduleTestModal = false;
            this.selectedTestId = '';
            this.submitting.set(false);
          },
          error: err => {
            this.toast.error('Test registered but failed to update application status');
            console.error('Error updating status:', err);
            this.submitting.set(false);
          },
        });
      },
      error: err => {
        this.toast.error('Failed to schedule test');
        console.error('Error scheduling test:', err);
        this.submitting.set(false);
      },
    });
  }

  formatTestDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  formatTestTime(timeStr: string): string {
    const [hours, minutes] = timeStr.split(':');
    const hour = parseInt(hours, 10);
    const ampm = hour >= 12 ? 'PM' : 'AM';
    const hour12 = hour % 12 || 12;
    return `${hour12}:${minutes} ${ampm}`;
  }

  onDocumentVerified(event: { documentId: string; status: string }): void {
    // Document verification is handled by the child component
    // Just refresh the application if needed
  }
}
