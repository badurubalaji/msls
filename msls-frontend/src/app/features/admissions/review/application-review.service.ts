/**
 * MSLS Application Review Service
 *
 * HTTP service for application review and document verification.
 */

import { Injectable, inject, signal } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay, tap } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  AdmissionApplication,
  ApplicationReview,
  AddReviewDto,
  VerifyDocumentDto,
  UpdateStatusDto,
  ApplicationFilterParams,
  ApplicationDocument,
} from './application-review.model';

/**
 * ApplicationReviewService - Handles application review API operations.
 */
@Injectable({ providedIn: 'root' })
export class ApplicationReviewService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/applications';

  /** Loading state */
  private readonly _loading = signal(false);
  readonly loading = this._loading.asReadonly();

  /** Current applications list */
  private readonly _applications = signal<AdmissionApplication[]>([]);
  readonly applications = this._applications.asReadonly();

  /** Selected application */
  private readonly _selectedApplication = signal<AdmissionApplication | null>(null);
  readonly selectedApplication = this._selectedApplication.asReadonly();

  // Mock data for development
  private mockApplications: AdmissionApplication[] = [
    {
      id: 'app-1',
      tenantId: '1',
      sessionId: '1',
      applicationNumber: 'ADM-2026-001',
      studentName: 'Rahul Sharma',
      dateOfBirth: '2020-05-15',
      gender: 'Male',
      bloodGroup: 'B+',
      nationality: 'Indian',
      religion: 'Hindu',
      category: 'General',
      previousSchool: 'ABC Public School',
      previousClass: 'UKG',
      classApplying: 'Class 1',
      academicYear: '2026-27',
      parentName: 'Amit Sharma',
      parentPhone: '9876543210',
      parentEmail: 'amit.sharma@email.com',
      parentOccupation: 'Software Engineer',
      fatherName: 'Amit Sharma',
      motherName: 'Priya Sharma',
      address: '123, Green Valley, Sector 15',
      city: 'Bangalore',
      state: 'Karnataka',
      pincode: '560001',
      documents: [
        {
          id: 'd1',
          name: 'Birth Certificate',
          type: 'birth_certificate',
          fileName: 'birth_cert_rahul.pdf',
          fileUrl: '/uploads/birth_cert_rahul.pdf',
          fileSize: 245000,
          uploadedAt: '2026-01-15T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-16T11:00:00Z',
        },
        {
          id: 'd2',
          name: 'Transfer Certificate',
          type: 'transfer_certificate',
          fileName: 'tc_rahul.pdf',
          fileUrl: '/uploads/tc_rahul.pdf',
          fileSize: 180000,
          uploadedAt: '2026-01-15T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-16T11:00:00Z',
        },
        {
          id: 'd3',
          name: 'Address Proof',
          type: 'address_proof',
          fileName: 'address_rahul.pdf',
          fileUrl: '/uploads/address_rahul.pdf',
          fileSize: 320000,
          uploadedAt: '2026-01-15T10:00:00Z',
          status: 'pending',
        },
        {
          id: 'd4',
          name: 'Passport Photos',
          type: 'photos',
          fileName: 'photos_rahul.jpg',
          fileUrl: '/uploads/photos_rahul.jpg',
          fileSize: 150000,
          uploadedAt: '2026-01-15T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-16T11:00:00Z',
        },
      ],
      status: 'under_review',
      reviews: [
        {
          id: 'r1',
          tenantId: '1',
          applicationId: 'app-1',
          reviewerId: 'user-1',
          reviewerName: 'Admin User',
          reviewType: 'document_verification',
          status: 'pending',
          comments: 'Started document verification. Birth certificate and TC verified.',
          createdAt: '2026-01-16T11:00:00Z',
        },
      ],
      submittedAt: '2026-01-15T10:00:00Z',
      createdAt: '2026-01-15T10:00:00Z',
      updatedAt: '2026-01-16T11:00:00Z',
    },
    {
      id: 'app-2',
      tenantId: '1',
      sessionId: '1',
      applicationNumber: 'ADM-2026-002',
      studentName: 'Priya Patel',
      dateOfBirth: '2020-08-22',
      gender: 'Female',
      bloodGroup: 'O+',
      nationality: 'Indian',
      religion: 'Hindu',
      category: 'General',
      previousSchool: 'XYZ School',
      previousClass: 'Class 5',
      classApplying: 'Class 6',
      academicYear: '2026-27',
      parentName: 'Rajesh Patel',
      parentPhone: '9876543211',
      parentEmail: 'rajesh.patel@email.com',
      parentOccupation: 'Business',
      fatherName: 'Rajesh Patel',
      motherName: 'Meera Patel',
      address: '456, Lake View Apartments',
      city: 'Mumbai',
      state: 'Maharashtra',
      pincode: '400001',
      documents: [
        {
          id: 'd5',
          name: 'Birth Certificate',
          type: 'birth_certificate',
          fileName: 'birth_cert_priya.pdf',
          fileUrl: '/uploads/birth_cert_priya.pdf',
          fileSize: 220000,
          uploadedAt: '2026-01-14T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-15T10:00:00Z',
        },
        {
          id: 'd6',
          name: 'Transfer Certificate',
          type: 'transfer_certificate',
          fileName: 'tc_priya.pdf',
          fileUrl: '/uploads/tc_priya.pdf',
          fileSize: 190000,
          uploadedAt: '2026-01-14T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-15T10:00:00Z',
        },
        {
          id: 'd7',
          name: 'Report Card',
          type: 'report_card',
          fileName: 'report_priya.pdf',
          fileUrl: '/uploads/report_priya.pdf',
          fileSize: 450000,
          uploadedAt: '2026-01-14T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-15T10:00:00Z',
        },
      ],
      status: 'documents_verified',
      reviews: [
        {
          id: 'r2',
          tenantId: '1',
          applicationId: 'app-2',
          reviewerId: 'user-1',
          reviewerName: 'Admin User',
          reviewType: 'document_verification',
          status: 'approved',
          comments: 'All documents verified successfully.',
          createdAt: '2026-01-15T10:00:00Z',
        },
      ],
      submittedAt: '2026-01-14T10:00:00Z',
      createdAt: '2026-01-14T10:00:00Z',
      updatedAt: '2026-01-15T10:00:00Z',
    },
    {
      id: 'app-3',
      tenantId: '1',
      sessionId: '1',
      applicationNumber: 'ADM-2026-003',
      studentName: 'Vikram Singh',
      dateOfBirth: '2019-11-10',
      gender: 'Male',
      bloodGroup: 'A+',
      nationality: 'Indian',
      religion: 'Sikh',
      category: 'General',
      classApplying: 'LKG',
      academicYear: '2026-27',
      parentName: 'Harpreet Singh',
      parentPhone: '9876543212',
      parentEmail: 'harpreet@email.com',
      parentOccupation: 'Doctor',
      fatherName: 'Harpreet Singh',
      motherName: 'Jasleen Kaur',
      address: '789, Model Town',
      city: 'Delhi',
      state: 'Delhi',
      pincode: '110001',
      documents: [
        {
          id: 'd8',
          name: 'Birth Certificate',
          type: 'birth_certificate',
          fileName: 'birth_cert_vikram.pdf',
          fileUrl: '/uploads/birth_cert_vikram.pdf',
          fileSize: 235000,
          uploadedAt: '2026-01-13T10:00:00Z',
          status: 'rejected',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-14T10:00:00Z',
          rejectionReason: 'Document is not clear. Please upload a clear scanned copy.',
        },
        {
          id: 'd9',
          name: 'Address Proof',
          type: 'address_proof',
          fileName: 'address_vikram.pdf',
          fileUrl: '/uploads/address_vikram.pdf',
          fileSize: 280000,
          uploadedAt: '2026-01-13T10:00:00Z',
          status: 'pending',
        },
      ],
      status: 'documents_pending',
      reviews: [
        {
          id: 'r3',
          tenantId: '1',
          applicationId: 'app-3',
          reviewerId: 'user-1',
          reviewerName: 'Admin User',
          reviewType: 'document_verification',
          status: 'rejected',
          comments: 'Birth certificate rejected - unclear document.',
          createdAt: '2026-01-14T10:00:00Z',
        },
      ],
      submittedAt: '2026-01-13T10:00:00Z',
      createdAt: '2026-01-13T10:00:00Z',
      updatedAt: '2026-01-14T10:00:00Z',
    },
    {
      id: 'app-4',
      tenantId: '1',
      sessionId: '1',
      applicationNumber: 'ADM-2026-004',
      studentName: 'Ananya Gupta',
      dateOfBirth: '2020-03-25',
      gender: 'Female',
      bloodGroup: 'AB+',
      nationality: 'Indian',
      classApplying: 'Class 1',
      academicYear: '2026-27',
      parentName: 'Sunil Gupta',
      parentPhone: '9876543213',
      parentEmail: 'sunil.gupta@email.com',
      fatherName: 'Sunil Gupta',
      motherName: 'Neha Gupta',
      address: '321, Sunshine Colony',
      city: 'Pune',
      state: 'Maharashtra',
      pincode: '411001',
      documents: [
        {
          id: 'd10',
          name: 'Birth Certificate',
          type: 'birth_certificate',
          fileName: 'birth_cert_ananya.pdf',
          fileUrl: '/uploads/birth_cert_ananya.pdf',
          fileSize: 210000,
          uploadedAt: '2026-01-18T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-19T10:00:00Z',
        },
        {
          id: 'd11',
          name: 'Transfer Certificate',
          type: 'transfer_certificate',
          fileName: 'tc_ananya.pdf',
          fileUrl: '/uploads/tc_ananya.pdf',
          fileSize: 175000,
          uploadedAt: '2026-01-18T10:00:00Z',
          status: 'verified',
          verifiedBy: 'Admin',
          verifiedAt: '2026-01-19T10:00:00Z',
        },
      ],
      status: 'test_scheduled',
      testId: '1',
      reviews: [
        {
          id: 'r4',
          tenantId: '1',
          applicationId: 'app-4',
          reviewerId: 'user-1',
          reviewerName: 'Admin User',
          reviewType: 'document_verification',
          status: 'approved',
          comments: 'Documents verified. Scheduled for entrance test.',
          createdAt: '2026-01-19T10:00:00Z',
        },
      ],
      submittedAt: '2026-01-18T10:00:00Z',
      createdAt: '2026-01-18T10:00:00Z',
      updatedAt: '2026-01-20T10:00:00Z',
    },
    {
      id: 'app-5',
      tenantId: '1',
      sessionId: '1',
      applicationNumber: 'ADM-2026-005',
      studentName: 'Arjun Reddy',
      dateOfBirth: '2020-07-08',
      gender: 'Male',
      classApplying: 'Class 1',
      academicYear: '2026-27',
      parentName: 'Krishna Reddy',
      parentPhone: '9876543214',
      fatherName: 'Krishna Reddy',
      motherName: 'Lakshmi Reddy',
      address: '555, Tech Park Road',
      city: 'Hyderabad',
      state: 'Telangana',
      pincode: '500001',
      documents: [
        {
          id: 'd12',
          name: 'Birth Certificate',
          type: 'birth_certificate',
          fileName: 'birth_cert_arjun.pdf',
          fileUrl: '/uploads/birth_cert_arjun.pdf',
          fileSize: 225000,
          uploadedAt: '2026-01-17T10:00:00Z',
          status: 'pending',
        },
      ],
      status: 'submitted',
      reviews: [],
      submittedAt: '2026-01-17T10:00:00Z',
      createdAt: '2026-01-17T10:00:00Z',
      updatedAt: '2026-01-17T10:00:00Z',
    },
  ];

  /**
   * Get all applications with optional filters
   */
  getApplications(filters?: ApplicationFilterParams): Observable<AdmissionApplication[]> {
    this._loading.set(true);
    // TODO: Replace with actual API call
    // return this.apiService.get<AdmissionApplication[]>(this.basePath, { params: filters });
    let filtered = [...this.mockApplications];

    if (filters?.sessionId) {
      filtered = filtered.filter(a => a.sessionId === filters.sessionId);
    }
    if (filters?.status) {
      filtered = filtered.filter(a => a.status === filters.status);
    }
    if (filters?.classApplying) {
      filtered = filtered.filter(a => a.classApplying === filters.classApplying);
    }
    if (filters?.search) {
      const search = filters.search.toLowerCase();
      filtered = filtered.filter(
        a =>
          a.studentName.toLowerCase().includes(search) ||
          a.applicationNumber.toLowerCase().includes(search) ||
          a.parentName.toLowerCase().includes(search)
      );
    }

    return of(filtered).pipe(
      delay(500),
      tap(() => {
        this._applications.set(filtered);
        this._loading.set(false);
      })
    );
  }

  /**
   * Get a single application by ID
   */
  getApplication(id: string): Observable<AdmissionApplication> {
    this._loading.set(true);
    // TODO: Replace with actual API call
    // return this.apiService.get<AdmissionApplication>(`${this.basePath}/${id}`);
    const application = this.mockApplications.find(a => a.id === id);
    if (!application) {
      throw new Error('Application not found');
    }
    return of({ ...application }).pipe(
      delay(300),
      tap(app => {
        this._selectedApplication.set(app);
        this._loading.set(false);
      })
    );
  }

  /**
   * Add a review to an application
   */
  addReview(applicationId: string, data: AddReviewDto): Observable<ApplicationReview> {
    // TODO: Replace with actual API call
    // return this.apiService.post<ApplicationReview>(`${this.basePath}/${applicationId}/review`, data);
    const index = this.mockApplications.findIndex(a => a.id === applicationId);
    if (index === -1) {
      throw new Error('Application not found');
    }

    const newReview: ApplicationReview = {
      id: String(Date.now()),
      tenantId: '1',
      applicationId,
      reviewerId: 'user-1',
      reviewerName: 'Admin User',
      reviewType: data.reviewType,
      status: data.status,
      comments: data.comments,
      createdAt: new Date().toISOString(),
    };

    this.mockApplications[index].reviews.push(newReview);
    this.mockApplications[index].updatedAt = new Date().toISOString();

    return of({ ...newReview }).pipe(
      delay(500),
      tap(() => {
        this._selectedApplication.set({ ...this.mockApplications[index] });
      })
    );
  }

  /**
   * Get reviews for an application
   */
  getReviews(applicationId: string): Observable<ApplicationReview[]> {
    // TODO: Replace with actual API call
    // return this.apiService.get<ApplicationReview[]>(`${this.basePath}/${applicationId}/reviews`);
    const application = this.mockApplications.find(a => a.id === applicationId);
    return of(application?.reviews || []).pipe(delay(300));
  }

  /**
   * Verify or reject a document
   */
  verifyDocument(
    applicationId: string,
    documentId: string,
    data: VerifyDocumentDto
  ): Observable<ApplicationDocument> {
    // TODO: Replace with actual API call
    // return this.apiService.patch<ApplicationDocument>(
    //   `${this.basePath}/${applicationId}/documents/${documentId}/verify`,
    //   data
    // );
    const appIndex = this.mockApplications.findIndex(a => a.id === applicationId);
    if (appIndex === -1) {
      throw new Error('Application not found');
    }

    const docIndex = this.mockApplications[appIndex].documents.findIndex(d => d.id === documentId);
    if (docIndex === -1) {
      throw new Error('Document not found');
    }

    this.mockApplications[appIndex].documents[docIndex] = {
      ...this.mockApplications[appIndex].documents[docIndex],
      status: data.status,
      rejectionReason: data.rejectionReason,
      verifiedBy: data.status !== 'pending' ? 'Admin' : undefined,
      verifiedAt: data.status !== 'pending' ? new Date().toISOString() : undefined,
    };

    // Update application status based on document verification
    const docs = this.mockApplications[appIndex].documents;
    const allVerified = docs.every(d => d.status === 'verified');
    const hasRejected = docs.some(d => d.status === 'rejected');
    const hasPending = docs.some(d => d.status === 'pending');

    if (allVerified) {
      this.mockApplications[appIndex].status = 'documents_verified';
    } else if (hasRejected || hasPending) {
      this.mockApplications[appIndex].status = 'documents_pending';
    }

    this.mockApplications[appIndex].updatedAt = new Date().toISOString();

    return of({ ...this.mockApplications[appIndex].documents[docIndex] }).pipe(
      delay(500),
      tap(() => {
        this._selectedApplication.set({ ...this.mockApplications[appIndex] });
        this._applications.set([...this.mockApplications]);
      })
    );
  }

  /**
   * Update application status
   */
  updateStatus(applicationId: string, data: UpdateStatusDto): Observable<AdmissionApplication> {
    // TODO: Replace with actual API call
    // return this.apiService.patch<AdmissionApplication>(`${this.basePath}/${applicationId}/status`, data);
    const index = this.mockApplications.findIndex(a => a.id === applicationId);
    if (index === -1) {
      throw new Error('Application not found');
    }

    this.mockApplications[index] = {
      ...this.mockApplications[index],
      status: data.status,
      updatedAt: new Date().toISOString(),
    };

    // If updating to rejected, add a review
    if (data.status === 'rejected' && data.comments) {
      const review: ApplicationReview = {
        id: String(Date.now()),
        tenantId: '1',
        applicationId,
        reviewerId: 'user-1',
        reviewerName: 'Admin User',
        reviewType: 'final_decision',
        status: 'rejected',
        comments: data.comments,
        createdAt: new Date().toISOString(),
      };
      this.mockApplications[index].reviews.push(review);
    }

    return of({ ...this.mockApplications[index] }).pipe(
      delay(500),
      tap(updated => {
        this._selectedApplication.set(updated);
        this._applications.set([...this.mockApplications]);
      })
    );
  }

  /**
   * Schedule test for application
   */
  scheduleTest(applicationId: string, testId: string): Observable<AdmissionApplication> {
    const index = this.mockApplications.findIndex(a => a.id === applicationId);
    if (index === -1) {
      throw new Error('Application not found');
    }

    this.mockApplications[index] = {
      ...this.mockApplications[index],
      status: 'test_scheduled',
      testId,
      updatedAt: new Date().toISOString(),
    };

    return of({ ...this.mockApplications[index] }).pipe(
      delay(500),
      tap(updated => {
        this._selectedApplication.set(updated);
        this._applications.set([...this.mockApplications]);
      })
    );
  }

  /**
   * Clear selected application
   */
  clearSelection(): void {
    this._selectedApplication.set(null);
  }
}
