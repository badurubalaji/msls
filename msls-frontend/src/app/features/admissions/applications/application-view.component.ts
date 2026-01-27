/**
 * MSLS Application View Component
 *
 * Read-only view for displaying admission application details.
 * Separate from the form component for cleaner code and better performance.
 */

import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';

import { ApplicationService } from './application.service';
import { AdmissionApplication, ApplicationDocument, DOCUMENT_TYPE_LABELS } from './application.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'app-application-view',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="view-page">
      <!-- Header -->
      <div class="view-header">
        <button class="back-btn" (click)="goBack()">
          <i class="fa-solid fa-arrow-left"></i>
        </button>
        <div class="view-header__content">
          <h1 class="view-header__title">Application Details</h1>
          <p class="view-header__subtitle">View admission application information</p>
        </div>
        <div class="view-header__actions">
          <button class="btn btn-primary" (click)="goToEdit()">
            <i class="fa-solid fa-pen-to-square"></i>
            Edit Application
          </button>
        </div>
      </div>

      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <p>Loading application...</p>
        </div>
      } @else if (application()) {
        <div class="view-content">
          <!-- Status Badge -->
          <div class="status-section">
            <span class="status-badge" [class]="'status-badge--' + application()!.currentStage">
              {{ getStatusLabel(application()!.currentStage) }}
            </span>
            <span class="application-id">Application #{{ application()!.applicationNumber }}</span>
          </div>

          <!-- Student Details -->
          <section class="detail-section">
            <h2 class="section-title">
              <i class="fa-solid fa-user-graduate"></i>
              Student Details
            </h2>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">Full Name</span>
                <span class="detail-value">{{ getFullName() }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Class Applying For</span>
                <span class="detail-value">{{ application()!.classApplying }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Date of Birth</span>
                <span class="detail-value">{{ application()!.dateOfBirth | date:'mediumDate' }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Gender</span>
                <span class="detail-value">{{ application()!.gender | titlecase }}</span>
              </div>
              @if (application()!.bloodGroup) {
                <div class="detail-item">
                  <span class="detail-label">Blood Group</span>
                  <span class="detail-value">{{ application()!.bloodGroup }}</span>
                </div>
              }
              @if (application()!.nationality) {
                <div class="detail-item">
                  <span class="detail-label">Nationality</span>
                  <span class="detail-value">{{ application()!.nationality }}</span>
                </div>
              }
              @if (application()!.religion) {
                <div class="detail-item">
                  <span class="detail-label">Religion</span>
                  <span class="detail-value">{{ application()!.religion }}</span>
                </div>
              }
              @if (application()!.category) {
                <div class="detail-item">
                  <span class="detail-label">Category</span>
                  <span class="detail-value">{{ application()!.category | titlecase }}</span>
                </div>
              }
              @if (application()!.aadhaarNumber) {
                <div class="detail-item">
                  <span class="detail-label">Aadhaar Number</span>
                  <span class="detail-value">{{ maskAadhaar(application()!.aadhaarNumber!) }}</span>
                </div>
              }
            </div>

            <!-- Address -->
            @if (hasAddress()) {
              <h3 class="subsection-title">Address</h3>
              <div class="detail-grid">
                @if (application()!.addressLine1) {
                  <div class="detail-item detail-item--full">
                    <span class="detail-label">Address</span>
                    <span class="detail-value">
                      {{ application()!.addressLine1 }}
                      @if (application()!.addressLine2) {
                        , {{ application()!.addressLine2 }}
                      }
                    </span>
                  </div>
                }
                @if (application()!.city) {
                  <div class="detail-item">
                    <span class="detail-label">City</span>
                    <span class="detail-value">{{ application()!.city }}</span>
                  </div>
                }
                @if (application()!.state) {
                  <div class="detail-item">
                    <span class="detail-label">State</span>
                    <span class="detail-value">{{ application()!.state }}</span>
                  </div>
                }
                @if (application()!.postalCode) {
                  <div class="detail-item">
                    <span class="detail-label">Postal Code</span>
                    <span class="detail-value">{{ application()!.postalCode }}</span>
                  </div>
                }
              </div>
            }
          </section>

          <!-- Parents/Guardians -->
          <section class="detail-section">
            <h2 class="section-title">
              <i class="fa-solid fa-users"></i>
              Parents / Guardians
            </h2>
            @if (application()!.parents && application()!.parents!.length > 0) {
              <div class="parents-grid">
                @for (parent of application()!.parents; track parent.id) {
                  <div class="parent-card">
                    <div class="parent-card__header">
                      <span class="parent-relation">{{ parent.relation | titlecase }}</span>
                    </div>
                    <div class="parent-card__body">
                      <div class="detail-item">
                        <span class="detail-label">Name</span>
                        <span class="detail-value">{{ parent.name }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="detail-label">Phone</span>
                        <span class="detail-value">{{ parent.phone }}</span>
                      </div>
                      @if (parent.email) {
                        <div class="detail-item">
                          <span class="detail-label">Email</span>
                          <span class="detail-value">{{ parent.email }}</span>
                        </div>
                      }
                      @if (parent.occupation) {
                        <div class="detail-item">
                          <span class="detail-label">Occupation</span>
                          <span class="detail-value">{{ parent.occupation }}</span>
                        </div>
                      }
                      @if (parent.education) {
                        <div class="detail-item">
                          <span class="detail-label">Education</span>
                          <span class="detail-value">{{ parent.education }}</span>
                        </div>
                      }
                      @if (parent.annualIncome) {
                        <div class="detail-item">
                          <span class="detail-label">Annual Income</span>
                          <span class="detail-value">{{ parent.annualIncome }}</span>
                        </div>
                      }
                    </div>
                  </div>
                }
              </div>
            } @else {
              <div class="empty-state">
                <i class="fa-solid fa-user-slash"></i>
                <p>No parent/guardian information added</p>
              </div>
            }
          </section>

          <!-- Previous School -->
          @if (hasPreviousSchool()) {
            <section class="detail-section">
              <h2 class="section-title">
                <i class="fa-solid fa-school"></i>
                Previous School
              </h2>
              <div class="detail-grid">
                @if (application()!.previousSchool) {
                  <div class="detail-item detail-item--full">
                    <span class="detail-label">School Name</span>
                    <span class="detail-value">{{ application()!.previousSchool }}</span>
                  </div>
                }
                @if (application()!.previousClass) {
                  <div class="detail-item">
                    <span class="detail-label">Previous Class</span>
                    <span class="detail-value">{{ application()!.previousClass }}</span>
                  </div>
                }
                @if (application()!.previousPercentage) {
                  <div class="detail-item">
                    <span class="detail-label">Percentage</span>
                    <span class="detail-value">{{ application()!.previousPercentage }}%</span>
                  </div>
                }
              </div>
            </section>
          }

          <!-- Documents -->
          <section class="detail-section">
            <h2 class="section-title">
              <i class="fa-solid fa-folder-open"></i>
              Documents
            </h2>
            @if (application()!.documents && application()!.documents!.length > 0) {
              <div class="documents-grid">
                @for (doc of application()!.documents; track doc.id) {
                  <div class="document-card">
                    <div class="document-card__icon">
                      <i class="fa-solid fa-file-pdf"></i>
                    </div>
                    <div class="document-card__content">
                      <span class="document-type">{{ getDocumentTypeLabel(doc.documentType) }}</span>
                      <span class="document-name">{{ doc.fileName }}</span>
                    </div>
                    <div class="document-card__status">
                      @if (doc.isVerified) {
                        <span class="doc-status doc-status--verified">
                          <i class="fa-solid fa-circle-check"></i>
                          Verified
                        </span>
                      } @else {
                        <span class="doc-status doc-status--pending">
                          <i class="fa-solid fa-clock"></i>
                          Pending
                        </span>
                      }
                    </div>
                    <a [href]="doc.fileUrl" target="_blank" class="document-card__action" title="View Document">
                      <i class="fa-solid fa-external-link"></i>
                    </a>
                  </div>
                }
              </div>
            } @else {
              <div class="empty-state">
                <i class="fa-solid fa-folder-open"></i>
                <p>No documents uploaded</p>
              </div>
            }
          </section>

          <!-- Timeline -->
          <section class="detail-section">
            <h2 class="section-title">
              <i class="fa-solid fa-clock-rotate-left"></i>
              Timeline
            </h2>
            <div class="timeline">
              <div class="timeline-item">
                <div class="timeline-marker"></div>
                <div class="timeline-content">
                  <span class="timeline-label">Created</span>
                  <span class="timeline-date">{{ application()!.createdAt | date:'medium' }}</span>
                </div>
              </div>
              @if (application()!.submittedAt) {
                <div class="timeline-item">
                  <div class="timeline-marker"></div>
                  <div class="timeline-content">
                    <span class="timeline-label">Submitted</span>
                    <span class="timeline-date">{{ application()!.submittedAt | date:'medium' }}</span>
                  </div>
                </div>
              }
              <div class="timeline-item">
                <div class="timeline-marker"></div>
                <div class="timeline-content">
                  <span class="timeline-label">Last Updated</span>
                  <span class="timeline-date">{{ application()!.updatedAt | date:'medium' }}</span>
                </div>
              </div>
            </div>
          </section>
        </div>
      } @else {
        <div class="error-container">
          <i class="fa-solid fa-exclamation-triangle"></i>
          <p>Application not found</p>
          <button class="btn btn-primary" (click)="goBack()">Go Back</button>
        </div>
      }
    </div>
  `,
  styles: [`
    .view-page {
      padding: 1.5rem;
      max-width: 1000px;
      margin: 0 auto;
    }

    /* Header */
    .view-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.5rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid #e5e7eb;
    }

    .back-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 2.5rem;
      height: 2.5rem;
      border: 1px solid #e5e7eb;
      border-radius: 0.5rem;
      background: white;
      color: #6b7280;
      cursor: pointer;
      transition: all 0.2s;
    }

    .back-btn:hover {
      background: #f3f4f6;
      color: #374151;
    }

    .view-header__content {
      flex: 1;
    }

    .view-header__title {
      font-size: 1.5rem;
      font-weight: 600;
      color: #111827;
      margin: 0;
    }

    .view-header__subtitle {
      font-size: 0.875rem;
      color: #6b7280;
      margin: 0.25rem 0 0 0;
    }

    .view-header__actions {
      display: flex;
      gap: 0.75rem;
    }

    /* Status Section */
    .status-section {
      display: flex;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      padding: 0.375rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
    }

    .status-badge--draft { background: #f3f4f6; color: #6b7280; }
    .status-badge--submitted { background: #dbeafe; color: #1d4ed8; }
    .status-badge--under_review { background: #fef3c7; color: #d97706; }
    .status-badge--documents_pending { background: #fee2e2; color: #dc2626; }
    .status-badge--documents_verified { background: #d1fae5; color: #059669; }
    .status-badge--approved { background: #d1fae5; color: #059669; }
    .status-badge--rejected { background: #fee2e2; color: #dc2626; }
    .status-badge--waitlisted { background: #e0e7ff; color: #4f46e5; }

    .application-id {
      font-size: 0.875rem;
      color: #6b7280;
    }

    /* Detail Sections */
    .detail-section {
      background: white;
      border: 1px solid #e5e7eb;
      border-radius: 0.75rem;
      padding: 1.5rem;
      margin-bottom: 1rem;
    }

    .section-title {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 1.125rem;
      font-weight: 600;
      color: #111827;
      margin: 0 0 1rem 0;
      padding-bottom: 0.75rem;
      border-bottom: 1px solid #f3f4f6;
    }

    .section-title i {
      color: #4f46e5;
    }

    .subsection-title {
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
      margin: 1.5rem 0 0.75rem 0;
    }

    .detail-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 1rem;
    }

    .detail-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .detail-item--full {
      grid-column: 1 / -1;
    }

    .detail-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #6b7280;
      text-transform: uppercase;
      letter-spacing: 0.025em;
    }

    .detail-value {
      font-size: 0.9375rem;
      color: #111827;
      font-weight: 500;
    }

    /* Parents Grid */
    .parents-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1rem;
    }

    .parent-card {
      background: #f9fafb;
      border: 1px solid #e5e7eb;
      border-radius: 0.5rem;
      overflow: hidden;
    }

    .parent-card__header {
      background: #4f46e5;
      color: white;
      padding: 0.75rem 1rem;
    }

    .parent-relation {
      font-weight: 600;
      font-size: 0.875rem;
    }

    .parent-card__body {
      padding: 1rem;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    /* Documents Grid */
    .documents-grid {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .document-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem;
      background: #f9fafb;
      border: 1px solid #e5e7eb;
      border-radius: 0.5rem;
    }

    .document-card__icon {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #fee2e2;
      color: #dc2626;
      border-radius: 0.375rem;
      font-size: 1.25rem;
    }

    .document-card__content {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .document-type {
      font-size: 0.875rem;
      font-weight: 600;
      color: #111827;
    }

    .document-name {
      font-size: 0.75rem;
      color: #6b7280;
    }

    .document-card__status {
      margin-right: 0.5rem;
    }

    .doc-status {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      padding: 0.25rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .doc-status--verified {
      background: #d1fae5;
      color: #059669;
    }

    .doc-status--pending {
      background: #fef3c7;
      color: #d97706;
    }

    .document-card__action {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border-radius: 0.375rem;
      color: #6b7280;
      transition: all 0.2s;
    }

    .document-card__action:hover {
      background: #e5e7eb;
      color: #374151;
    }

    /* Timeline */
    .timeline {
      display: flex;
      flex-direction: column;
      gap: 0;
    }

    .timeline-item {
      display: flex;
      gap: 1rem;
      padding: 0.75rem 0;
      position: relative;
    }

    .timeline-item:not(:last-child)::after {
      content: '';
      position: absolute;
      left: 0.4375rem;
      top: 2rem;
      bottom: 0;
      width: 2px;
      background: #e5e7eb;
    }

    .timeline-marker {
      width: 0.875rem;
      height: 0.875rem;
      background: #4f46e5;
      border-radius: 50%;
      flex-shrink: 0;
      margin-top: 0.125rem;
    }

    .timeline-content {
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .timeline-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .timeline-date {
      font-size: 0.75rem;
      color: #6b7280;
    }

    /* Empty State */
    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.5rem;
      padding: 2rem;
      color: #9ca3af;
      text-align: center;
    }

    .empty-state i {
      font-size: 2rem;
    }

    .empty-state p {
      margin: 0;
      font-size: 0.875rem;
    }

    /* Loading & Error */
    .loading-container,
    .error-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 4rem;
      color: #6b7280;
    }

    .error-container i {
      font-size: 3rem;
      color: #f59e0b;
    }

    .spinner {
      width: 2.5rem;
      height: 2.5rem;
      border: 3px solid #e5e7eb;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      border: none;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    /* Responsive */
    @media (max-width: 640px) {
      .view-header {
        flex-wrap: wrap;
      }

      .view-header__actions {
        width: 100%;
        margin-top: 0.5rem;
      }

      .view-header__actions .btn {
        flex: 1;
        justify-content: center;
      }

      .detail-grid {
        grid-template-columns: 1fr;
      }

      .parents-grid {
        grid-template-columns: 1fr;
      }
    }
  `]
})
export class ApplicationViewComponent implements OnInit {
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private applicationService = inject(ApplicationService);
  private toastService = inject(ToastService);

  loading = signal(true);
  application = signal<AdmissionApplication | null>(null);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadApplication(id);
    } else {
      this.loading.set(false);
    }
  }

  private loadApplication(id: string): void {
    this.applicationService.getApplication(id).subscribe({
      next: (app) => {
        this.application.set(app);
        this.loading.set(false);
      },
      error: (err) => {
        this.toastService.error('Failed to load application');
        this.loading.set(false);
        console.error('Failed to load application:', err);
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/admissions/applications']);
  }

  goToEdit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.router.navigate(['/admissions/applications', id, 'edit']);
    }
  }

  getFullName(): string {
    const app = this.application();
    if (!app) return '';
    const parts = [app.firstName, app.middleName, app.lastName].filter(Boolean);
    return parts.join(' ');
  }

  getStatusLabel(stage: string): string {
    const labels: Record<string, string> = {
      draft: 'Draft',
      submitted: 'Submitted',
      under_review: 'Under Review',
      documents_pending: 'Documents Pending',
      documents_verified: 'Documents Verified',
      interview_scheduled: 'Interview Scheduled',
      interview_completed: 'Interview Completed',
      approved: 'Approved',
      rejected: 'Rejected',
      waitlisted: 'Waitlisted',
      enrolled: 'Enrolled'
    };
    return labels[stage] || stage;
  }

  getDocumentTypeLabel(type: string): string {
    return DOCUMENT_TYPE_LABELS[type as keyof typeof DOCUMENT_TYPE_LABELS] || type;
  }

  maskAadhaar(aadhaar: string): string {
    if (!aadhaar || aadhaar.length < 4) return aadhaar;
    return 'XXXX-XXXX-' + aadhaar.slice(-4);
  }

  hasAddress(): boolean {
    const app = this.application();
    return !!(app?.addressLine1 || app?.city || app?.state || app?.postalCode);
  }

  hasPreviousSchool(): boolean {
    const app = this.application();
    return !!(app?.previousSchool || app?.previousClass || app?.previousPercentage);
  }
}
