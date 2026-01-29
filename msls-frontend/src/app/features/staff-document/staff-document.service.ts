/**
 * Staff Document Service
 * Story 5.8: Staff Document Management
 *
 * HTTP service for staff document management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import {
  DocumentType,
  DocumentTypeListResponse,
  CreateDocumentTypeRequest,
  UpdateDocumentTypeRequest,
  StaffDocument,
  DocumentListResponse,
  UpdateDocumentRequest,
  VerifyDocumentRequest,
  RejectDocumentRequest,
  ExpiringDocument,
  ComplianceReportResponse,
} from './staff-document.model';

@Injectable({ providedIn: 'root' })
export class StaffDocumentService {
  private readonly apiService = inject(ApiService);

  // ========================================
  // Document Type Methods
  // ========================================

  /**
   * Get all document types
   */
  getDocumentTypes(activeOnly: boolean = false): Observable<DocumentTypeListResponse> {
    const params: Record<string, string> = {};
    if (activeOnly) params['active_only'] = 'true';
    return this.apiService.get<DocumentTypeListResponse>('/staff-document-types', { params });
  }

  /**
   * Get a document type by ID
   */
  getDocumentType(id: string): Observable<DocumentType> {
    return this.apiService.get<DocumentType>(`/staff-document-types/${id}`);
  }

  /**
   * Create a new document type
   */
  createDocumentType(data: CreateDocumentTypeRequest): Observable<DocumentType> {
    return this.apiService.post<DocumentType>('/staff-document-types', data);
  }

  /**
   * Update a document type
   */
  updateDocumentType(id: string, data: UpdateDocumentTypeRequest): Observable<DocumentType> {
    return this.apiService.put<DocumentType>(`/staff-document-types/${id}`, data);
  }

  /**
   * Delete a document type
   */
  deleteDocumentType(id: string): Observable<void> {
    return this.apiService.delete<void>(`/staff-document-types/${id}`);
  }

  // ========================================
  // Staff Document Methods
  // ========================================

  /**
   * Get documents for a staff member
   */
  getStaffDocuments(staffId: string, params?: {
    document_type_id?: string;
    status?: string;
    is_current?: boolean;
    cursor?: string;
    limit?: number;
  }): Observable<DocumentListResponse> {
    const queryParams: Record<string, string> = {};
    if (params?.document_type_id) queryParams['document_type_id'] = params.document_type_id;
    if (params?.status) queryParams['status'] = params.status;
    if (params?.is_current !== undefined) queryParams['is_current'] = String(params.is_current);
    if (params?.cursor) queryParams['cursor'] = params.cursor;
    if (params?.limit) queryParams['limit'] = String(params.limit);
    return this.apiService.get<DocumentListResponse>(`/staff/${staffId}/documents`, { params: queryParams });
  }

  /**
   * Get a single document
   */
  getDocument(staffId: string, documentId: string): Observable<StaffDocument> {
    return this.apiService.get<StaffDocument>(`/staff/${staffId}/documents/${documentId}`);
  }

  /**
   * Upload a document for a staff member
   */
  uploadDocument(staffId: string, documentTypeId: string, file: File, metadata?: {
    document_number?: string;
    issue_date?: string;
    expiry_date?: string;
    remarks?: string;
  }): Observable<StaffDocument> {
    const additionalData: Record<string, string> = {
      document_type_id: documentTypeId,
    };
    if (metadata?.document_number) additionalData['document_number'] = metadata.document_number;
    if (metadata?.issue_date) additionalData['issue_date'] = metadata.issue_date;
    if (metadata?.expiry_date) additionalData['expiry_date'] = metadata.expiry_date;
    if (metadata?.remarks) additionalData['remarks'] = metadata.remarks;

    return this.apiService.uploadFile<StaffDocument>(`/staff/${staffId}/documents`, file, 'file', additionalData);
  }

  /**
   * Update document metadata
   */
  updateDocument(staffId: string, documentId: string, data: UpdateDocumentRequest): Observable<StaffDocument> {
    return this.apiService.put<StaffDocument>(`/staff/${staffId}/documents/${documentId}`, data);
  }

  /**
   * Delete a document
   */
  deleteDocument(staffId: string, documentId: string): Observable<void> {
    return this.apiService.delete<void>(`/staff/${staffId}/documents/${documentId}`);
  }

  /**
   * Get download URL for a document
   */
  getDocumentDownloadUrl(staffId: string, documentId: string): Observable<{ download_url: string }> {
    return this.apiService.get<{ download_url: string }>(`/staff/${staffId}/documents/${documentId}/download`);
  }

  // ========================================
  // Verification Methods
  // ========================================

  /**
   * Verify a document
   */
  verifyDocument(staffId: string, documentId: string, data: VerifyDocumentRequest): Observable<StaffDocument> {
    return this.apiService.put<StaffDocument>(`/staff/${staffId}/documents/${documentId}/verify`, data);
  }

  /**
   * Reject a document
   */
  rejectDocument(staffId: string, documentId: string, data: RejectDocumentRequest): Observable<StaffDocument> {
    return this.apiService.put<StaffDocument>(`/staff/${staffId}/documents/${documentId}/reject`, data);
  }

  // ========================================
  // Expiry & Compliance Methods
  // ========================================

  /**
   * Get documents expiring within given days
   */
  getExpiringDocuments(days: number = 30): Observable<{ documents: ExpiringDocument[]; days: number }> {
    return this.apiService.get<{ documents: ExpiringDocument[]; days: number }>('/staff-documents/expiring', {
      params: { days: String(days) },
    });
  }

  /**
   * Get compliance report
   */
  getComplianceReport(includeStaffDetails: boolean = false): Observable<ComplianceReportResponse> {
    const params: Record<string, string> = {};
    if (includeStaffDetails) params['include_staff_details'] = 'true';
    return this.apiService.get<ComplianceReportResponse>('/staff-documents/compliance', { params });
  }
}
