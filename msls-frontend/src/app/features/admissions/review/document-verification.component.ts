/**
 * MSLS Document Verification Component
 *
 * Component for verifying and rejecting application documents.
 */

import {
  Component,
  input,
  output,
  signal,
  inject,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsBadgeComponent } from '../../../shared/components/badge/badge.component';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { ToastService } from '../../../shared/services/toast.service';

import { ApplicationReviewService } from './application-review.service';
import {
  ApplicationDocument,
  DocumentStatus,
  DOCUMENT_STATUS_CONFIG,
  formatFileSize,
} from './application-review.model';

/**
 * DocumentVerificationComponent - Document verification UI.
 */
@Component({
  selector: 'msls-document-verification',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsButtonComponent,
    MslsBadgeComponent,
    MslsModalComponent,
  ],
  template: `
    <div class="document-verification">
      @if (documents().length === 0) {
        <div class="no-documents">
          <i class="fa-regular fa-file"></i>
          <p>No documents uploaded</p>
        </div>
      } @else {
        <div class="documents-list">
          @for (doc of documents(); track doc.id) {
            <div class="document-item" [class.verified]="doc.status === 'verified'" [class.rejected]="doc.status === 'rejected'">
              <div class="doc-icon">
                <i [class]="getFileIcon(doc.fileName)"></i>
              </div>
              <div class="doc-info">
                <div class="doc-name">{{ doc.name }}</div>
                <div class="doc-meta">
                  <span class="file-name">{{ doc.fileName }}</span>
                  <span class="file-size">{{ formatSize(doc.fileSize) }}</span>
                </div>
                @if (doc.status === 'rejected' && doc.rejectionReason) {
                  <div class="rejection-reason">
                    <i class="fa-solid fa-exclamation-triangle"></i>
                    {{ doc.rejectionReason }}
                  </div>
                }
                @if (doc.status === 'verified' && doc.verifiedAt) {
                  <div class="verified-info">
                    <i class="fa-solid fa-check"></i>
                    Verified by {{ doc.verifiedBy }} on {{ formatDate(doc.verifiedAt) }}
                  </div>
                }
              </div>
              <div class="doc-status">
                <msls-badge
                  [variant]="getStatusVariant(doc.status)"
                  size="sm"
                >
                  <i [class]="getStatusIcon(doc.status)"></i>
                  {{ getStatusLabel(doc.status) }}
                </msls-badge>
              </div>
              <div class="doc-actions">
                <button
                  class="action-btn view"
                  title="View Document"
                  (click)="viewDocument(doc)"
                >
                  <i class="fa-solid fa-eye"></i>
                </button>
                @if (doc.status === 'pending' || doc.status === 'rejected') {
                  <button
                    class="action-btn verify"
                    title="Verify Document"
                    (click)="verifyDoc(doc)"
                    [disabled]="processing()"
                  >
                    <i class="fa-solid fa-check"></i>
                  </button>
                }
                @if (doc.status === 'pending' || doc.status === 'verified') {
                  <button
                    class="action-btn reject"
                    title="Reject Document"
                    (click)="openRejectModal(doc)"
                    [disabled]="processing()"
                  >
                    <i class="fa-solid fa-times"></i>
                  </button>
                }
              </div>
            </div>
          }
        </div>
      }

      <!-- Reject Document Modal -->
      <msls-modal
        [isOpen]="rejectModalOpen"
        (closed)="rejectModalOpen = false"
        title="Reject Document"
        size="sm"
      >
        <ng-container modal-body>
          <div class="reject-doc-info">
            <i class="fa-solid fa-file-circle-xmark"></i>
            <span>{{ selectedDoc?.name }}</span>
          </div>
          <div class="form-group">
            <label>Reason for Rejection</label>
            <textarea
              class="form-textarea"
              rows="3"
              placeholder="Enter the reason for rejecting this document..."
              [(ngModel)]="rejectReason"
            ></textarea>
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="modal-actions">
            <msls-button variant="ghost" (click)="rejectModalOpen = false">Cancel</msls-button>
            <msls-button
              variant="danger"
              [loading]="processing()"
              [disabled]="!rejectReason.trim()"
              (click)="confirmReject()"
            >
              Reject Document
            </msls-button>
          </div>
        </ng-container>
      </msls-modal>

      <!-- View Document Modal -->
      <msls-modal
        [isOpen]="viewModalOpen"
        (closed)="viewModalOpen = false"
        title="View Document"
        size="lg"
      >
        <ng-container modal-body>
          <div class="document-preview">
            @if (selectedDoc) {
              @if (isImage(selectedDoc.fileName)) {
                <img [src]="selectedDoc.fileUrl" [alt]="selectedDoc.name" class="preview-image" />
              } @else if (isPdf(selectedDoc.fileName)) {
                <div class="pdf-preview">
                  <i class="fa-solid fa-file-pdf"></i>
                  <p>{{ selectedDoc.fileName }}</p>
                  <a [href]="selectedDoc.fileUrl" target="_blank" class="download-link">
                    <i class="fa-solid fa-download"></i>
                    Download PDF
                  </a>
                </div>
              } @else {
                <div class="file-preview">
                  <i class="fa-solid fa-file"></i>
                  <p>{{ selectedDoc.fileName }}</p>
                  <a [href]="selectedDoc.fileUrl" target="_blank" class="download-link">
                    <i class="fa-solid fa-download"></i>
                    Download File
                  </a>
                </div>
              }
            }
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="modal-actions">
            <msls-button variant="ghost" (click)="viewModalOpen = false">Close</msls-button>
            @if (selectedDoc?.status === 'pending') {
              <msls-button variant="danger" (click)="openRejectFromView()">
                <i class="fa-solid fa-times"></i>
                Reject
              </msls-button>
              <msls-button variant="primary" [loading]="processing()" (click)="verifyDoc(selectedDoc!)">
                <i class="fa-solid fa-check"></i>
                Verify
              </msls-button>
            }
          </div>
        </ng-container>
      </msls-modal>
    </div>
  `,
  styles: [`
    .document-verification {
      width: 100%;
    }

    .no-documents {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 3rem 1rem;
      color: #94a3b8;
    }

    .no-documents i {
      font-size: 3rem;
      margin-bottom: 1rem;
    }

    .documents-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .document-item {
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      padding: 1rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      transition: all 0.2s;
    }

    .document-item:hover {
      border-color: #cbd5e1;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }

    .document-item.verified {
      border-color: #86efac;
      background: linear-gradient(to right, #f0fdf4, white);
    }

    .document-item.rejected {
      border-color: #fca5a5;
      background: linear-gradient(to right, #fef2f2, white);
    }

    .doc-icon {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border-radius: 0.5rem;
      color: #64748b;
      font-size: 1.25rem;
      flex-shrink: 0;
    }

    .document-item.verified .doc-icon {
      background: #dcfce7;
      color: #16a34a;
    }

    .document-item.rejected .doc-icon {
      background: #fee2e2;
      color: #dc2626;
    }

    .doc-info {
      flex: 1;
      min-width: 0;
    }

    .doc-name {
      font-weight: 600;
      color: #0f172a;
      font-size: 0.875rem;
      margin-bottom: 0.25rem;
    }

    .doc-meta {
      display: flex;
      gap: 0.75rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .file-name {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 150px;
    }

    .rejection-reason {
      margin-top: 0.5rem;
      padding: 0.5rem;
      background: #fef2f2;
      border-radius: 0.375rem;
      font-size: 0.75rem;
      color: #991b1b;
      display: flex;
      align-items: center;
      gap: 0.375rem;
    }

    .verified-info {
      margin-top: 0.5rem;
      font-size: 0.75rem;
      color: #16a34a;
      display: flex;
      align-items: center;
      gap: 0.375rem;
    }

    .doc-status {
      flex-shrink: 0;
    }

    .doc-actions {
      display: flex;
      gap: 0.375rem;
      flex-shrink: 0;
    }

    .action-btn {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
      background: white;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 0.75rem;
    }

    .action-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .action-btn.view {
      color: #64748b;
    }

    .action-btn.view:hover:not(:disabled) {
      background: #f1f5f9;
      border-color: #cbd5e1;
      color: #334155;
    }

    .action-btn.verify {
      color: #16a34a;
    }

    .action-btn.verify:hover:not(:disabled) {
      background: #dcfce7;
      border-color: #86efac;
    }

    .action-btn.reject {
      color: #dc2626;
    }

    .action-btn.reject:hover:not(:disabled) {
      background: #fee2e2;
      border-color: #fca5a5;
    }

    /* Modal Styles */
    .reject-doc-info {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem;
      background: #f1f5f9;
      border-radius: 0.5rem;
      margin-bottom: 1rem;
    }

    .reject-doc-info i {
      font-size: 1.5rem;
      color: #ef4444;
    }

    .reject-doc-info span {
      font-weight: 500;
      color: #0f172a;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .form-group label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .form-textarea {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #0f172a;
      resize: vertical;
      min-height: 80px;
    }

    .form-textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .modal-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
    }

    .document-preview {
      min-height: 300px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .preview-image {
      max-width: 100%;
      max-height: 500px;
      border-radius: 0.5rem;
    }

    .pdf-preview,
    .file-preview {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 1rem;
      padding: 3rem;
      background: #f8fafc;
      border-radius: 0.75rem;
    }

    .pdf-preview i,
    .file-preview i {
      font-size: 4rem;
      color: #dc2626;
    }

    .file-preview i {
      color: #64748b;
    }

    .pdf-preview p,
    .file-preview p {
      font-weight: 500;
      color: #0f172a;
      margin: 0;
    }

    .download-link {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      background: #4f46e5;
      color: white;
      border-radius: 0.5rem;
      text-decoration: none;
      font-size: 0.875rem;
      font-weight: 500;
      transition: background 0.2s;
    }

    .download-link:hover {
      background: #4338ca;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DocumentVerificationComponent {
  private readonly reviewService = inject(ApplicationReviewService);
  private readonly toast = inject(ToastService);

  /** Input: List of documents */
  readonly documents = input<ApplicationDocument[]>([]);

  /** Input: Application ID */
  readonly applicationId = input<string>('');

  /** Output: Document verified event */
  readonly documentVerified = output<{ documentId: string; status: string }>();

  /** Processing state */
  readonly processing = signal(false);

  /** Modal states */
  rejectModalOpen = false;
  viewModalOpen = false;

  /** Selected document */
  selectedDoc: ApplicationDocument | null = null;
  rejectReason = '';

  getStatusVariant(status: DocumentStatus): 'success' | 'warning' | 'danger' {
    const config = DOCUMENT_STATUS_CONFIG[status];
    return config?.variant as 'success' | 'warning' | 'danger' || 'warning';
  }

  getStatusIcon(status: DocumentStatus): string {
    return DOCUMENT_STATUS_CONFIG[status]?.icon || 'fa-solid fa-clock';
  }

  getStatusLabel(status: DocumentStatus): string {
    return DOCUMENT_STATUS_CONFIG[status]?.label || status;
  }

  getFileIcon(fileName: string): string {
    const ext = fileName.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return 'fa-solid fa-file-pdf';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
        return 'fa-solid fa-file-image';
      case 'doc':
      case 'docx':
        return 'fa-solid fa-file-word';
      default:
        return 'fa-solid fa-file';
    }
  }

  formatSize(bytes: number): string {
    return formatFileSize(bytes);
  }

  formatDate(dateString?: string): string {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  isImage(fileName: string): boolean {
    const ext = fileName.split('.').pop()?.toLowerCase();
    return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext || '');
  }

  isPdf(fileName: string): boolean {
    return fileName.toLowerCase().endsWith('.pdf');
  }

  viewDocument(doc: ApplicationDocument): void {
    this.selectedDoc = doc;
    this.viewModalOpen = true;
  }

  verifyDoc(doc: ApplicationDocument): void {
    if (!this.applicationId()) return;

    this.processing.set(true);
    this.reviewService
      .verifyDocument(this.applicationId(), doc.id, { status: 'verified' })
      .subscribe({
        next: () => {
          this.toast.success(`${doc.name} verified successfully`);
          this.documentVerified.emit({ documentId: doc.id, status: 'verified' });
          this.processing.set(false);
          this.viewModalOpen = false;
        },
        error: (err) => {
          this.toast.error('Failed to verify document');
          console.error('Error verifying document:', err);
          this.processing.set(false);
        },
      });
  }

  openRejectModal(doc: ApplicationDocument): void {
    this.selectedDoc = doc;
    this.rejectReason = '';
    this.rejectModalOpen = true;
  }

  openRejectFromView(): void {
    this.viewModalOpen = false;
    this.rejectReason = '';
    this.rejectModalOpen = true;
  }

  confirmReject(): void {
    if (!this.applicationId() || !this.selectedDoc || !this.rejectReason.trim()) return;

    this.processing.set(true);
    this.reviewService
      .verifyDocument(this.applicationId(), this.selectedDoc.id, {
        status: 'rejected',
        rejectionReason: this.rejectReason,
      })
      .subscribe({
        next: () => {
          this.toast.success(`${this.selectedDoc?.name} rejected`);
          this.documentVerified.emit({ documentId: this.selectedDoc!.id, status: 'rejected' });
          this.processing.set(false);
          this.rejectModalOpen = false;
          this.selectedDoc = null;
        },
        error: (err) => {
          this.toast.error('Failed to reject document');
          console.error('Error rejecting document:', err);
          this.processing.set(false);
        },
      });
  }
}
