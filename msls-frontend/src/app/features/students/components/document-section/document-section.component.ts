/**
 * Document Section Component
 *
 * Displays document checklist and management for a student.
 */

import { Component, OnInit, inject, input, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MslsSpinnerComponent } from '../../../../shared/components/spinner/spinner.component';
import { MslsBadgeComponent } from '../../../../shared/components/badge/badge.component';
import { DocumentService } from '../../services/document.service';
import {
  DocumentType,
  StudentDocument,
  DocumentChecklistItem,
  DocumentChecklistResponse,
  DocumentStatus,
  getDocumentStatusBadgeVariant,
  getDocumentStatusLabel,
  formatFileSize,
  getFileIcon,
  isExtensionAllowed,
} from '../../models/document.model';

@Component({
  selector: 'msls-document-section',
  standalone: true,
  imports: [CommonModule, MslsSpinnerComponent, MslsBadgeComponent],
  templateUrl: './document-section.component.html',
  styleUrl: './document-section.component.scss',
})
export class DocumentSectionComponent implements OnInit {
  private documentService = inject(DocumentService);

  /** Student ID input */
  studentId = input.required<string>();

  /** Loading state */
  loading = this.documentService.loading;

  /** Checklist data */
  checklist = this.documentService.checklist;

  /** Whether checklist has items */
  hasItems = computed(() => (this.checklist()?.items?.length ?? 0) > 0);

  /** Show upload modal */
  showUploadModal = signal(false);

  /** Document type being uploaded */
  uploadingType = signal<DocumentType | null>(null);

  /** Show preview modal */
  showPreviewModal = signal(false);

  /** Document being previewed */
  previewDocument = signal<StudentDocument | null>(null);

  /** Show verify modal */
  showVerifyModal = signal(false);

  /** Document being verified */
  verifyingDocument = signal<StudentDocument | null>(null);

  /** Show reject modal */
  showRejectModal = signal(false);

  /** Document being rejected */
  rejectingDocument = signal<StudentDocument | null>(null);

  /** Rejection reason */
  rejectionReason = signal('');

  /** Upload form state */
  uploadForm = signal({
    documentNumber: '',
    issueDate: '',
    expiryDate: '',
  });

  /** Selected file for upload */
  selectedFile = signal<File | null>(null);

  /** Upload error message */
  uploadError = signal<string | null>(null);

  ngOnInit(): void {
    this.loadChecklist();
  }

  private loadChecklist(): void {
    this.documentService.loadChecklist(this.studentId()).subscribe();
  }

  // =========================================================================
  // Template Helpers
  // =========================================================================

  getStatusBadgeVariant(status: DocumentStatus): 'success' | 'warning' | 'danger' | 'neutral' {
    return getDocumentStatusBadgeVariant(status);
  }

  getStatusLabel(status: DocumentStatus): string {
    return getDocumentStatusLabel(status);
  }

  formatSize(bytes: number): string {
    return formatFileSize(bytes);
  }

  getIcon(mimeType: string): string {
    return getFileIcon(mimeType);
  }

  formatDate(dateStr: string | undefined): string {
    if (!dateStr) return '';
    return new Date(dateStr).toLocaleDateString();
  }

  getStatusIcon(item: DocumentChecklistItem): string {
    if (!item.document) return 'fa-circle';
    switch (item.document.status) {
      case 'verified':
        return 'fa-check-circle';
      case 'rejected':
        return 'fa-times-circle';
      case 'pending_verification':
        return 'fa-clock';
      default:
        return 'fa-circle';
    }
  }

  getStatusClass(item: DocumentChecklistItem): string {
    if (!item.document) return 'status--missing';
    switch (item.document.status) {
      case 'verified':
        return 'status--verified';
      case 'rejected':
        return 'status--rejected';
      case 'pending_verification':
        return 'status--pending';
      default:
        return 'status--missing';
    }
  }

  getAcceptedExtensions(): string {
    const extensions = this.uploadingType()?.allowedExtensions;
    if (!extensions) return '';
    return extensions.split(',').map((e: string) => '.' + e.trim()).join(',');
  }

  // =========================================================================
  // Upload Handling
  // =========================================================================

  openUploadModal(docType: DocumentType): void {
    this.uploadingType.set(docType);
    this.selectedFile.set(null);
    this.uploadError.set(null);
    this.uploadForm.set({
      documentNumber: '',
      issueDate: '',
      expiryDate: '',
    });
    this.showUploadModal.set(true);
  }

  closeUploadModal(): void {
    this.showUploadModal.set(false);
    this.uploadingType.set(null);
    this.selectedFile.set(null);
    this.uploadError.set(null);
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      const file = input.files[0];
      const docType = this.uploadingType();

      if (!docType) return;

      // Validate file extension
      if (!isExtensionAllowed(file.name, docType.allowedExtensions)) {
        this.uploadError.set(`Invalid file type. Allowed: ${docType.allowedExtensions}`);
        this.selectedFile.set(null);
        return;
      }

      // Validate file size
      const maxSizeBytes = docType.maxSizeMb * 1024 * 1024;
      if (file.size > maxSizeBytes) {
        this.uploadError.set(`File too large. Maximum: ${docType.maxSizeMb}MB`);
        this.selectedFile.set(null);
        return;
      }

      this.uploadError.set(null);
      this.selectedFile.set(file);
    }
  }

  onDocumentNumberInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.uploadForm.update((f) => ({ ...f, documentNumber: input.value }));
  }

  onIssueDateChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.uploadForm.update((f) => ({ ...f, issueDate: input.value }));
  }

  onExpiryDateChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.uploadForm.update((f) => ({ ...f, expiryDate: input.value }));
  }

  uploadDocument(): void {
    const file = this.selectedFile();
    const docType = this.uploadingType();

    if (!file || !docType) return;

    const form = this.uploadForm();
    const metadata: { documentNumber?: string; issueDate?: string; expiryDate?: string } = {};
    if (form.documentNumber) metadata.documentNumber = form.documentNumber;
    if (form.issueDate) metadata.issueDate = form.issueDate;
    if (form.expiryDate) metadata.expiryDate = form.expiryDate;

    this.documentService.uploadDocument(this.studentId(), file, docType.id, metadata).subscribe({
      next: () => {
        this.closeUploadModal();
        this.loadChecklist();
      },
      error: (err) => {
        this.uploadError.set(err.message || 'Failed to upload document');
      },
    });
  }

  // =========================================================================
  // Preview Handling
  // =========================================================================

  openPreview(doc: StudentDocument): void {
    this.previewDocument.set(doc);
    this.showPreviewModal.set(true);
  }

  closePreview(): void {
    this.showPreviewModal.set(false);
    this.previewDocument.set(null);
  }

  downloadDocument(doc: StudentDocument): void {
    window.open(doc.fileUrl, '_blank');
  }

  // =========================================================================
  // Verification Handling
  // =========================================================================

  openVerifyModal(doc: StudentDocument): void {
    this.verifyingDocument.set(doc);
    this.showVerifyModal.set(true);
  }

  closeVerifyModal(): void {
    this.showVerifyModal.set(false);
    this.verifyingDocument.set(null);
  }

  confirmVerify(): void {
    const doc = this.verifyingDocument();
    if (!doc) return;

    this.documentService.verifyDocument(this.studentId(), doc.id, { version: doc.version }).subscribe({
      next: () => {
        this.closeVerifyModal();
        this.loadChecklist();
      },
    });
  }

  openRejectModal(doc: StudentDocument): void {
    this.rejectingDocument.set(doc);
    this.rejectionReason.set('');
    this.showRejectModal.set(true);
  }

  closeRejectModal(): void {
    this.showRejectModal.set(false);
    this.rejectingDocument.set(null);
    this.rejectionReason.set('');
  }

  onRejectReasonInput(event: Event): void {
    const input = event.target as HTMLTextAreaElement;
    this.rejectionReason.set(input.value);
  }

  confirmReject(): void {
    const doc = this.rejectingDocument();
    const reason = this.rejectionReason().trim();
    if (!doc || !reason) return;

    this.documentService.rejectDocument(this.studentId(), doc.id, { reason, version: doc.version }).subscribe({
      next: () => {
        this.closeRejectModal();
        this.loadChecklist();
      },
    });
  }

  // =========================================================================
  // Delete Handling
  // =========================================================================

  deleteDocument(doc: StudentDocument): void {
    if (!confirm('Are you sure you want to delete this document?')) return;

    this.documentService.deleteDocument(this.studentId(), doc.id).subscribe({
      next: () => {
        this.loadChecklist();
      },
    });
  }
}
