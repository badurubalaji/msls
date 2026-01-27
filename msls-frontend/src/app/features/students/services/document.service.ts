/**
 * MSLS Document Service
 *
 * Handles all document-related HTTP API calls with reactive state management.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, throwError, finalize } from 'rxjs';

import { ApiService } from '../../../core/services/api.service';
import {
  DocumentType,
  StudentDocument,
  DocumentChecklistResponse,
  DocumentTypeListResponse,
  DocumentListResponse,
  UpdateDocumentRequest,
  VerifyDocumentRequest,
  RejectDocumentRequest,
} from '../models/document.model';

/**
 * DocumentService - Manages document data with reactive signals.
 *
 * Provides operations for student document management.
 */
@Injectable({ providedIn: 'root' })
export class DocumentService {
  private api = inject(ApiService);

  // =========================================================================
  // State Signals
  // =========================================================================

  /** List of document types */
  private _documentTypes = signal<DocumentType[]>([]);

  /** List of student documents */
  private _documents = signal<StudentDocument[]>([]);

  /** Document checklist */
  private _checklist = signal<DocumentChecklistResponse | null>(null);

  /** Loading state */
  private _loading = signal<boolean>(false);

  /** Error state */
  private _error = signal<string | null>(null);

  // =========================================================================
  // Public Readonly Signals
  // =========================================================================

  /** Public readonly document types */
  readonly documentTypes = this._documentTypes.asReadonly();

  /** Public readonly documents list */
  readonly documents = this._documents.asReadonly();

  /** Public readonly checklist */
  readonly checklist = this._checklist.asReadonly();

  /** Public readonly loading state */
  readonly loading = this._loading.asReadonly();

  /** Public readonly error state */
  readonly error = this._error.asReadonly();

  /** Whether there are any documents */
  readonly hasDocuments = computed(() => this._documents().length > 0);

  // =========================================================================
  // Document Type Operations
  // =========================================================================

  /**
   * Load document types.
   * @param activeOnly - Whether to only load active types
   */
  loadDocumentTypes(activeOnly = true): Observable<DocumentTypeListResponse> {
    this._loading.set(true);
    this._error.set(null);

    const params: Record<string, string> = {};
    if (!activeOnly) params['active_only'] = 'false';

    return this.api.get<DocumentTypeListResponse>('/document-types', { params }).pipe(
      tap((response) => {
        this._documentTypes.set(response.documentTypes);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load document types');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  // =========================================================================
  // Student Document Operations
  // =========================================================================

  /**
   * Load documents for a student.
   * @param studentId - Student ID
   */
  loadDocuments(studentId: string): Observable<DocumentListResponse> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<DocumentListResponse>(`/students/${studentId}/documents`).pipe(
      tap((response) => {
        this._documents.set(response.documents);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load documents');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get a single document.
   * @param studentId - Student ID
   * @param documentId - Document ID
   */
  getDocument(studentId: string, documentId: string): Observable<StudentDocument> {
    return this.api.get<StudentDocument>(`/students/${studentId}/documents/${documentId}`);
  }

  /**
   * Upload a document.
   * @param studentId - Student ID
   * @param file - File to upload
   * @param documentTypeId - Document type ID
   * @param metadata - Additional metadata
   */
  uploadDocument(
    studentId: string,
    file: File,
    documentTypeId: string,
    metadata?: {
      documentNumber?: string;
      issueDate?: string;
      expiryDate?: string;
    }
  ): Observable<StudentDocument> {
    this._loading.set(true);
    this._error.set(null);

    const formData = new FormData();
    formData.append('file', file);
    formData.append('documentTypeId', documentTypeId);
    if (metadata?.documentNumber) {
      formData.append('documentNumber', metadata.documentNumber);
    }
    if (metadata?.issueDate) {
      formData.append('issueDate', metadata.issueDate);
    }
    if (metadata?.expiryDate) {
      formData.append('expiryDate', metadata.expiryDate);
    }

    return this.api.post<StudentDocument>(`/students/${studentId}/documents`, formData).pipe(
      tap((document) => {
        // Update documents list
        this._documents.update((docs) => {
          const existing = docs.findIndex((d) => d.documentTypeId === document.documentTypeId);
          if (existing >= 0) {
            // Replace existing
            return docs.map((d, i) => (i === existing ? document : d));
          }
          // Add new
          return [document, ...docs];
        });
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to upload document');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update document metadata.
   * @param studentId - Student ID
   * @param documentId - Document ID
   * @param data - Update data
   */
  updateDocument(studentId: string, documentId: string, data: UpdateDocumentRequest): Observable<StudentDocument> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.put<StudentDocument>(`/students/${studentId}/documents/${documentId}`, data).pipe(
      tap((document) => {
        this._documents.update((docs) => docs.map((d) => (d.id === documentId ? document : d)));
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update document');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Delete a document.
   * @param studentId - Student ID
   * @param documentId - Document ID
   */
  deleteDocument(studentId: string, documentId: string): Observable<void> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.delete<void>(`/students/${studentId}/documents/${documentId}`).pipe(
      tap(() => {
        this._documents.update((docs) => docs.filter((d) => d.id !== documentId));
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to delete document');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Verify a document.
   * @param studentId - Student ID
   * @param documentId - Document ID
   * @param data - Verify request data
   */
  verifyDocument(studentId: string, documentId: string, data: VerifyDocumentRequest): Observable<StudentDocument> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<StudentDocument>(`/students/${studentId}/documents/${documentId}/verify`, data).pipe(
      tap((document) => {
        this._documents.update((docs) => docs.map((d) => (d.id === documentId ? document : d)));
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to verify document');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Reject a document.
   * @param studentId - Student ID
   * @param documentId - Document ID
   * @param data - Reject request data
   */
  rejectDocument(studentId: string, documentId: string, data: RejectDocumentRequest): Observable<StudentDocument> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<StudentDocument>(`/students/${studentId}/documents/${documentId}/reject`, data).pipe(
      tap((document) => {
        this._documents.update((docs) => docs.map((d) => (d.id === documentId ? document : d)));
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to reject document');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Load document checklist for a student.
   * @param studentId - Student ID
   */
  loadChecklist(studentId: string): Observable<DocumentChecklistResponse> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<DocumentChecklistResponse>(`/students/${studentId}/document-checklist`).pipe(
      tap((checklist) => {
        this._checklist.set(checklist);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load document checklist');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  // =========================================================================
  // State Management Helpers
  // =========================================================================

  /**
   * Clear error state.
   */
  clearError(): void {
    this._error.set(null);
  }

  /**
   * Reset all state.
   */
  reset(): void {
    this._documentTypes.set([]);
    this._documents.set([]);
    this._checklist.set(null);
    this._loading.set(false);
    this._error.set(null);
  }
}
