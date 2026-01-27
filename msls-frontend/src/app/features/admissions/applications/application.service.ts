/**
 * MSLS Application Service
 *
 * HTTP service for admission application management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  AdmissionApplication,
  ApplicationParent,
  ApplicationDocument,
  CreateApplicationRequest,
  UpdateApplicationRequest,
  ParentRequest,
  UpdateStageRequest,
  ApplicationFilterParams,
  ApplicationStage,
  DocumentType,
  StatusCheckRequest,
  StatusCheckResponse,
  Gender,
  BloodGroup,
  ParentRelation,
} from './application.model';

/**
 * API Response interfaces matching backend DTOs
 */
interface ApplicationApiResponse {
  id: string;
  applicationNumber: string;
  studentName: string;
  className: string;
  status: string;
  source?: string;
  applicantDetails?: {
    gender?: string;
    dateOfBirth?: string;
    nationality?: string;
    religion?: string;
    category?: string;
    bloodGroup?: string;
    address?: string;
    city?: string;
    state?: string;
    pinCode?: string;
    previousSchool?: string;
    transferReason?: string;
  };
  parentInfo?: {
    fatherName?: string;
    fatherOccupation?: string;
    fatherPhone?: string;
    fatherEmail?: string;
    motherName?: string;
    motherOccupation?: string;
    motherPhone?: string;
    motherEmail?: string;
    guardianName?: string;
    guardianRelation?: string;
    guardianPhone?: string;
    guardianEmail?: string;
  };
  applicationFee?: number;
  feePaid?: boolean;
  paymentDate?: string;
  paymentReference?: string;
  submittedAt?: string;
  reviewedAt?: string;
  reviewNotes?: string;
  approvedAt?: string;
  enrolledAt?: string;
  waitlistPosition?: number;
  session?: { id: string; name: string };
  parents?: ParentApiResponse[];
  documents?: DocumentApiResponse[];
  createdAt: string;
  updatedAt: string;
}

interface ParentApiResponse {
  id: string;
  relation: string;
  name: string;
  phone: string;
  email: string;
  occupation: string;
  education: string;
  annualIncome: string;
  createdAt: string;
  updatedAt: string;
}

interface DocumentApiResponse {
  id: string;
  documentType: string;
  fileUrl: string;
  fileName: string;
  fileSize: number;
  mimeType: string;
  isVerified: boolean;
  verifiedAt?: string;
  rejectionReason?: string;
  createdAt: string;
  updatedAt: string;
}

interface ApplicationListApiResponse {
  applications: ApplicationApiResponse[];
  total: number;
}

/**
 * ApplicationService - Handles all admission application-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class ApplicationService {
  private readonly apiService = inject(ApiService);
  private readonly http = inject(HttpClient);
  private readonly basePath = '/applications';

  /**
   * Map API response to AdmissionApplication model
   */
  private mapApplicationResponse(response: ApplicationApiResponse): AdmissionApplication {
    // Parse student name into first/last
    const nameParts = (response.studentName || '').split(' ');
    const firstName = nameParts[0] || '';
    const lastName = nameParts.slice(1).join(' ') || '';

    // Map gender safely
    const genderValue = response.applicantDetails?.gender;
    const gender: Gender = (genderValue === 'male' || genderValue === 'female' || genderValue === 'other')
      ? genderValue
      : 'male';

    // Map blood group safely
    const bloodGroupValue = response.applicantDetails?.bloodGroup;
    const validBloodGroups = ['A+', 'A-', 'B+', 'B-', 'O+', 'O-', 'AB+', 'AB-'];
    const bloodGroup: BloodGroup | undefined = validBloodGroups.includes(bloodGroupValue || '')
      ? bloodGroupValue as BloodGroup
      : undefined;

    // Convert parentInfo (flat) to parents array, or use existing parents array
    const parents = this.mapParentInfoToArray(response.parentInfo, response.id, response.parents);

    return {
      id: response.id,
      tenantId: '',
      sessionId: response.session?.id || '',
      sessionName: response.session?.name,
      applicationNumber: response.applicationNumber,
      currentStage: this.mapStatusToStage(response.status),
      stageHistory: [],
      classApplying: response.className,
      firstName,
      lastName,
      dateOfBirth: response.applicantDetails?.dateOfBirth || '',
      gender,
      bloodGroup,
      nationality: response.applicantDetails?.nationality || 'Indian',
      religion: response.applicantDetails?.religion,
      category: response.applicantDetails?.category,
      addressLine1: response.applicantDetails?.address,
      city: response.applicantDetails?.city,
      state: response.applicantDetails?.state,
      postalCode: response.applicantDetails?.pinCode,
      previousSchool: response.applicantDetails?.previousSchool,
      parents,
      documents: (response.documents || []).map(d => this.mapDocumentResponse(d, response.id)),
      submittedAt: response.submittedAt,
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  }

  /**
   * Convert parentInfo (flat structure) to parents array
   */
  private mapParentInfoToArray(
    parentInfo: ApplicationApiResponse['parentInfo'],
    applicationId: string,
    existingParents?: ParentApiResponse[]
  ): ApplicationParent[] {
    // If there are existing parents from the parents array, use those
    if (existingParents && existingParents.length > 0) {
      return existingParents.map(p => this.mapParentResponse(p, applicationId));
    }

    // Otherwise, convert parentInfo to parents array
    const parents: ApplicationParent[] = [];

    if (parentInfo?.fatherName) {
      parents.push({
        id: `father-${applicationId}`,
        applicationId,
        relation: 'father',
        name: parentInfo.fatherName,
        phone: parentInfo.fatherPhone || '',
        email: parentInfo.fatherEmail || '',
        occupation: parentInfo.fatherOccupation || '',
        createdAt: new Date().toISOString(),
      });
    }

    if (parentInfo?.motherName) {
      parents.push({
        id: `mother-${applicationId}`,
        applicationId,
        relation: 'mother',
        name: parentInfo.motherName,
        phone: parentInfo.motherPhone || '',
        email: parentInfo.motherEmail || '',
        occupation: parentInfo.motherOccupation || '',
        createdAt: new Date().toISOString(),
      });
    }

    if (parentInfo?.guardianName) {
      parents.push({
        id: `guardian-${applicationId}`,
        applicationId,
        relation: 'guardian',
        name: parentInfo.guardianName,
        phone: parentInfo.guardianPhone || '',
        email: parentInfo.guardianEmail || '',
        createdAt: new Date().toISOString(),
      });
    }

    return parents;
  }

  /**
   * Map backend status to frontend stage
   */
  private mapStatusToStage(status: string): ApplicationStage {
    const statusMap: Record<string, ApplicationStage> = {
      draft: 'draft',
      submitted: 'submitted',
      under_review: 'under_review',
      documents_pending: 'documents_pending',
      documents_verified: 'documents_verified',
      interview_scheduled: 'interview_scheduled',
      interview_completed: 'interview_completed',
      approved: 'approved',
      waitlisted: 'waitlisted',
      rejected: 'rejected',
      enrolled: 'admitted', // Map enrolled to admitted
      admitted: 'admitted',
      withdrawn: 'rejected', // Map withdrawn to rejected for now
    };
    return statusMap[status] || 'draft';
  }

  /**
   * Map parent API response
   */
  private mapParentResponse(response: ParentApiResponse, applicationId: string): ApplicationParent {
    // Map relation safely
    const relationValue = response.relation?.toLowerCase();
    const relation: ParentRelation = (relationValue === 'father' || relationValue === 'mother' || relationValue === 'guardian')
      ? relationValue
      : 'guardian';

    return {
      id: response.id,
      applicationId,
      relation,
      name: response.name,
      phone: response.phone,
      email: response.email,
      occupation: response.occupation,
      education: response.education,
      annualIncome: response.annualIncome,
      createdAt: response.createdAt,
    };
  }

  /**
   * Map document API response
   */
  private mapDocumentResponse(response: DocumentApiResponse, applicationId: string): ApplicationDocument {
    return {
      id: response.id,
      applicationId,
      documentType: response.documentType as DocumentType,
      fileUrl: response.fileUrl,
      fileName: response.fileName,
      isVerified: response.isVerified,
      verifiedAt: response.verifiedAt,
      createdAt: response.createdAt,
    };
  }

  /**
   * Get all applications with optional filters
   */
  getApplications(filters?: ApplicationFilterParams): Observable<AdmissionApplication[]> {
    const params: Record<string, string> = {};
    if (filters?.stage) params['status'] = filters.stage;
    if (filters?.classApplying) params['className'] = filters.classApplying;
    if (filters?.sessionId) params['sessionId'] = filters.sessionId;
    if (filters?.search) params['search'] = filters.search;

    return this.apiService.get<ApplicationListApiResponse>(this.basePath, { params }).pipe(
      map(response => (response.applications || []).map(a => this.mapApplicationResponse(a))),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Get a single application by ID
   */
  getApplication(id: string): Observable<AdmissionApplication> {
    return this.apiService.get<ApplicationApiResponse>(`${this.basePath}/${id}`).pipe(
      map(response => this.mapApplicationResponse(response)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Create a new application
   */
  createApplication(data: CreateApplicationRequest): Observable<AdmissionApplication> {
    const payload: Record<string, unknown> = {
      sessionId: data.sessionId,
      studentName: `${data.firstName} ${data.lastName || ''}`.trim(),
      className: data.classApplying,
      applicantDetails: {
        gender: data.gender,
        dateOfBirth: data.dateOfBirth,
        nationality: data.nationality,
        religion: data.religion,
        category: data.category,
        bloodGroup: data.bloodGroup,
        address: data.addressLine1,
        city: data.city,
        state: data.state,
        pinCode: data.postalCode,
        previousSchool: data.previousSchool,
      },
    };

    // Add parent info if present
    const parentInfo: Record<string, unknown> = {};
    if (data.fatherName) parentInfo['fatherName'] = data.fatherName;
    if (data.fatherPhone) parentInfo['fatherPhone'] = data.fatherPhone;
    if (data.fatherEmail) parentInfo['fatherEmail'] = data.fatherEmail;
    if (data.fatherOccupation) parentInfo['fatherOccupation'] = data.fatherOccupation;
    if (data.motherName) parentInfo['motherName'] = data.motherName;
    if (data.motherPhone) parentInfo['motherPhone'] = data.motherPhone;
    if (data.motherEmail) parentInfo['motherEmail'] = data.motherEmail;
    if (data.motherOccupation) parentInfo['motherOccupation'] = data.motherOccupation;
    if (data.guardianName) parentInfo['guardianName'] = data.guardianName;
    if (data.guardianPhone) parentInfo['guardianPhone'] = data.guardianPhone;
    if (data.guardianEmail) parentInfo['guardianEmail'] = data.guardianEmail;
    if (data.guardianRelation) parentInfo['guardianRelation'] = data.guardianRelation;

    if (Object.keys(parentInfo).length > 0) {
      payload['parentInfo'] = parentInfo;
    }

    return this.apiService.post<ApplicationApiResponse>(this.basePath, payload).pipe(
      map(response => this.mapApplicationResponse(response)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Update an existing application
   */
  updateApplication(id: string, data: UpdateApplicationRequest): Observable<AdmissionApplication> {
    const payload: Record<string, unknown> = {};

    if (data.firstName || data.lastName) {
      payload['studentName'] = `${data.firstName || ''} ${data.lastName || ''}`.trim();
    }
    if (data.classApplying) {
      payload['className'] = data.classApplying;
    }

    // Build applicant details if any are present
    const applicantDetails: Record<string, unknown> = {};
    if (data.gender) applicantDetails['gender'] = data.gender;
    if (data.dateOfBirth) applicantDetails['dateOfBirth'] = data.dateOfBirth;
    if (data.nationality) applicantDetails['nationality'] = data.nationality;
    if (data.religion) applicantDetails['religion'] = data.religion;
    if (data.category) applicantDetails['category'] = data.category;
    if (data.bloodGroup) applicantDetails['bloodGroup'] = data.bloodGroup;
    if (data.addressLine1) applicantDetails['address'] = data.addressLine1;
    if (data.city) applicantDetails['city'] = data.city;
    if (data.state) applicantDetails['state'] = data.state;
    if (data.postalCode) applicantDetails['pinCode'] = data.postalCode;
    if (data.previousSchool) applicantDetails['previousSchool'] = data.previousSchool;

    if (Object.keys(applicantDetails).length > 0) {
      payload['applicantDetails'] = applicantDetails;
    }

    // Build parent info if any are present
    const parentInfo: Record<string, unknown> = {};
    if (data.fatherName) parentInfo['fatherName'] = data.fatherName;
    if (data.fatherPhone) parentInfo['fatherPhone'] = data.fatherPhone;
    if (data.fatherEmail) parentInfo['fatherEmail'] = data.fatherEmail;
    if (data.fatherOccupation) parentInfo['fatherOccupation'] = data.fatherOccupation;
    if (data.motherName) parentInfo['motherName'] = data.motherName;
    if (data.motherPhone) parentInfo['motherPhone'] = data.motherPhone;
    if (data.motherEmail) parentInfo['motherEmail'] = data.motherEmail;
    if (data.motherOccupation) parentInfo['motherOccupation'] = data.motherOccupation;
    if (data.guardianName) parentInfo['guardianName'] = data.guardianName;
    if (data.guardianPhone) parentInfo['guardianPhone'] = data.guardianPhone;
    if (data.guardianEmail) parentInfo['guardianEmail'] = data.guardianEmail;
    if (data.guardianRelation) parentInfo['guardianRelation'] = data.guardianRelation;

    if (Object.keys(parentInfo).length > 0) {
      payload['parentInfo'] = parentInfo;
    }

    return this.apiService.put<ApplicationApiResponse>(`${this.basePath}/${id}`, payload).pipe(
      map(response => this.mapApplicationResponse(response)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Submit application for review
   */
  submitApplication(id: string): Observable<AdmissionApplication> {
    return this.apiService.post<ApplicationApiResponse>(`${this.basePath}/${id}/submit`, {}).pipe(
      map(response => this.mapApplicationResponse(response)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Update application stage
   */
  updateStage(id: string, data: UpdateStageRequest): Observable<AdmissionApplication> {
    // Map frontend stage to backend status (admitted -> enrolled)
    const backendStage = data.stage === 'admitted' ? 'enrolled' : data.stage;
    const payload = {
      newStage: backendStage,
      remarks: data.remarks,
    };
    return this.apiService.patch<ApplicationApiResponse>(`${this.basePath}/${id}/stage`, payload).pipe(
      map(response => this.mapApplicationResponse(response)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Add parent to application
   */
  addParent(applicationId: string, data: ParentRequest): Observable<ApplicationParent> {
    return this.apiService.post<ParentApiResponse>(`${this.basePath}/${applicationId}/parents`, data).pipe(
      map(response => this.mapParentResponse(response, applicationId)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Update parent details
   */
  updateParent(applicationId: string, parentId: string, data: ParentRequest): Observable<ApplicationParent> {
    return this.apiService.put<ParentApiResponse>(
      `${this.basePath}/${applicationId}/parents/${parentId}`,
      data
    ).pipe(
      map(response => this.mapParentResponse(response, applicationId)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Delete parent from application
   */
  deleteParent(applicationId: string, parentId: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${applicationId}/parents/${parentId}`);
  }

  /**
   * Get documents for an application
   */
  getDocuments(applicationId: string): Observable<ApplicationDocument[]> {
    return this.apiService.get<DocumentApiResponse[]>(`${this.basePath}/${applicationId}/documents`).pipe(
      map(response => (response || []).map(d => this.mapDocumentResponse(d, applicationId))),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Upload document
   */
  uploadDocument(applicationId: string, documentType: DocumentType, file: File): Observable<ApplicationDocument> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('documentType', documentType);

    // Use http directly since apiService doesn't handle FormData properly
    // Backend returns: { success: true, data: DocumentApiResponse }
    return this.http.post<{ success: boolean; data: DocumentApiResponse }>(
      `/api${this.basePath}/${applicationId}/documents`,
      formData
    ).pipe(
      map(response => {
        if (response.success && response.data) {
          return this.mapDocumentResponse(response.data, applicationId);
        }
        throw new Error('Upload failed');
      }),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Delete document
   */
  deleteDocument(applicationId: string, documentId: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${applicationId}/documents/${documentId}`);
  }

  /**
   * Verify document
   */
  verifyDocument(applicationId: string, documentId: string, verified: boolean): Observable<ApplicationDocument> {
    return this.apiService.patch<DocumentApiResponse>(
      `${this.basePath}/${applicationId}/documents/${documentId}/verify`,
      { isVerified: verified }
    ).pipe(
      map(response => this.mapDocumentResponse(response, applicationId)),
      catchError(err => throwError(() => err))
    );
  }

  /**
   * Get available admission sessions
   */
  getAvailableSessions(): Observable<{ id: string; name: string; status: string }[]> {
    interface SessionResponse {
      sessions: { id: string; name: string; status: string }[];
      total: number;
    }
    return this.apiService.get<SessionResponse>('/admission-sessions', {
      params: { status: 'open' }
    }).pipe(
      map(response => response.sessions || []),
      catchError(() => {
        // Return empty array on error
        return [];
      })
    );
  }

  /**
   * Delete application (admin only)
   */
  deleteApplication(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Check application status (public API - no authentication required)
   * Allows parents to check their application status using application number and phone
   */
  checkStatus(request: StatusCheckRequest): Observable<StatusCheckResponse> {
    return this.http.post<StatusCheckResponse>('/api/v1/public/applications/status', request).pipe(
      catchError(err => throwError(() => err))
    );
  }
}
