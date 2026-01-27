/**
 * MSLS Merit List Service
 *
 * HTTP service for merit list generation, admission decisions,
 * and enrollment operations.
 */

import { Injectable, inject } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, of, catchError, throwError } from 'rxjs';
import { delay, map } from 'rxjs/operators';
import { environment } from '../../../../environments/environment';

import { ApiService } from '../../../core/services/api.service';
import {
  MeritList,
  MeritListEntry,
  AdmissionDecision,
  GenerateMeritListRequest,
  MakeDecisionRequest,
  BulkDecisionRequest,
  BulkDecisionResponse,
  AcceptOfferRequest,
  EnrollmentRequest,
  EnrollmentResponse,
  ApplicationStatus,
} from './merit.model';
import { CLASS_NAMES } from '../sessions/admission-session.model';

/**
 * API Response interfaces
 */
interface MeritListApiResponse {
  id: string;
  sessionId: string;
  className: string;
  testId?: string;
  generatedAt: string;
  generatedBy?: string;
  cutoffScore?: number;
  entries: MeritListEntryApiResponse[];
  isFinal: boolean;
  totalCount: number;
  aboveCutoff: number;
  createdAt: string;
}

interface MeritListEntryApiResponse {
  rank: number;
  applicationId: string;
  studentName: string;
  score: number;
  testScore?: number;
  interviewScore?: number;
  previousMarks?: number;
  status: string;
  parentPhone: string;
  parentEmail?: string;
}

interface DecisionApiResponse {
  id: string;
  applicationId: string;
  decision: string;
  decisionDate: string;
  decidedBy?: string;
  sectionAssigned?: string;
  waitlistPosition?: number;
  rejectionReason?: string;
  offerLetterUrl?: string;
  offerValidUntil?: string;
  offerAccepted?: boolean;
  offerAcceptedAt?: string;
  remarks?: string;
  createdAt: string;
  updatedAt: string;
}

interface OfferLetterApiResponse {
  url: string;
  validUntil: string;
  generatedAt: string;
  applicationId: string;
}

interface EnrollmentApiResponse {
  applicationId: string;
  studentId?: string;
  enrollmentNumber?: string;
  admissionDate: string;
  className: string;
  sectionName?: string;
  rollNumber?: string;
  status: string;
}

/**
 * MeritService - Handles merit list and admission decision API operations.
 */
@Injectable({ providedIn: 'root' })
export class MeritService {
  private readonly apiService = inject(ApiService);
  private readonly basePathSessions = '/admission-sessions';
  private readonly basePathApplications = '/applications';
  private readonly basePathMeritLists = '/merit-lists';

  // Only use mock data in development when backend is unavailable (fallback)
  private readonly useMockData = !environment.production;

  // Mock data for development
  private mockMeritLists: Map<string, MeritList> = new Map();
  private mockEntries: MeritListEntry[] = this.generateMockEntries();

  /**
   * Generate mock merit list entries for testing
   */
  private generateMockEntries(): MeritListEntry[] {
    const names = [
      { student: 'Aarav Sharma', parent: 'Rajesh Sharma', phone: '9876543210' },
      { student: 'Priya Patel', parent: 'Vikram Patel', phone: '9876543211' },
      { student: 'Arjun Singh', parent: 'Manpreet Singh', phone: '9876543212' },
      { student: 'Ananya Gupta', parent: 'Sanjay Gupta', phone: '9876543213' },
      { student: 'Vihaan Kumar', parent: 'Anil Kumar', phone: '9876543214' },
      { student: 'Diya Reddy', parent: 'Krishna Reddy', phone: '9876543215' },
      { student: 'Aditya Joshi', parent: 'Mohan Joshi', phone: '9876543216' },
      { student: 'Ishita Verma', parent: 'Rakesh Verma', phone: '9876543217' },
      { student: 'Reyansh Mehta', parent: 'Deepak Mehta', phone: '9876543218' },
      { student: 'Saanvi Iyer', parent: 'Ramesh Iyer', phone: '9876543219' },
      { student: 'Kabir Das', parent: 'Suresh Das', phone: '9876543220' },
      { student: 'Myra Kapoor', parent: 'Amit Kapoor', phone: '9876543221' },
      { student: 'Vivaan Shah', parent: 'Rahul Shah', phone: '9876543222' },
      { student: 'Aisha Khan', parent: 'Salman Khan', phone: '9876543223' },
      { student: 'Dhruv Mishra', parent: 'Prakash Mishra', phone: '9876543224' },
    ];

    return names.map((name, index) => {
      const totalScore = Math.floor(Math.random() * 40) + 60; // 60-100
      const maxScore = 100;
      const statuses: ApplicationStatus[] = ['test_completed', 'selected', 'waitlisted', 'offer_sent'];
      const status = statuses[Math.floor(Math.random() * statuses.length)];

      return {
        id: `entry-${index + 1}`,
        applicationId: `app-${index + 1}`,
        studentName: name.student,
        parentName: name.parent,
        parentPhone: name.phone,
        applicationNumber: `APP-2026-${String(index + 1).padStart(4, '0')}`,
        classApplying: 'Class 1',
        totalScore,
        maxScore,
        percentage: (totalScore / maxScore) * 100,
        rank: 0, // Will be set when sorted
        status,
        testId: 'test-1',
        testName: 'Entrance Test 2026',
        subjectScores: [
          { subjectName: 'English', score: Math.floor(Math.random() * 15) + 10, maxScore: 25 },
          { subjectName: 'Mathematics', score: Math.floor(Math.random() * 15) + 10, maxScore: 25 },
          { subjectName: 'General Knowledge', score: Math.floor(Math.random() * 15) + 10, maxScore: 25 },
          { subjectName: 'Reasoning', score: Math.floor(Math.random() * 15) + 10, maxScore: 25 },
        ],
        sectionAssigned: status === 'selected' || status === 'offer_sent' ? 'Section A' : undefined,
        waitlistPosition: status === 'waitlisted' ? Math.floor(Math.random() * 5) + 1 : undefined,
      };
    })
      .sort((a, b) => b.totalScore - a.totalScore)
      .map((entry, index) => ({ ...entry, rank: index + 1 }));
  }

  /**
   * Map API response to MeritList model
   */
  private mapMeritListResponse(response: MeritListApiResponse): MeritList {
    return {
      id: response.id,
      tenantId: '',
      sessionId: response.sessionId,
      className: response.className,
      testId: response.testId,
      testName: '', // API doesn't return test name, would need separate call
      generatedAt: response.generatedAt,
      cutoffScore: response.cutoffScore,
      entries: response.entries.map((e, index) => this.mapMeritListEntryResponse(e, index)),
      isFinal: response.isFinal,
      createdAt: response.createdAt,
    };
  }

  /**
   * Map API response to MeritListEntry model
   */
  private mapMeritListEntryResponse(response: MeritListEntryApiResponse, index: number): MeritListEntry {
    return {
      id: `entry-${index}`,
      applicationId: response.applicationId,
      studentName: response.studentName,
      parentName: '',
      parentPhone: response.parentPhone,
      parentEmail: response.parentEmail,
      applicationNumber: '',
      classApplying: '',
      totalScore: response.score,
      maxScore: 100,
      percentage: response.score,
      rank: response.rank,
      status: this.mapStatusFromApi(response.status),
      testId: '',
      testName: '',
      subjectScores: [],
    };
  }

  /**
   * Map API status to frontend ApplicationStatus
   */
  private mapStatusFromApi(apiStatus: string): ApplicationStatus {
    const statusMap: Record<string, ApplicationStatus> = {
      'submitted': 'submitted',
      'under_review': 'under_review',
      'approved': 'selected',
      'waitlisted': 'waitlisted',
      'rejected': 'rejected',
      'enrolled': 'enrolled',
    };
    return statusMap[apiStatus] || 'test_completed';
  }

  /**
   * Map decision API response to AdmissionDecision model
   */
  private mapDecisionResponse(response: DecisionApiResponse): AdmissionDecision {
    return {
      id: response.id,
      tenantId: '',
      applicationId: response.applicationId,
      decision: response.decision as 'selected' | 'waitlisted' | 'rejected',
      decisionDate: response.decisionDate,
      decidedBy: response.decidedBy,
      sectionAssigned: response.sectionAssigned,
      waitlistPosition: response.waitlistPosition,
      rejectionReason: response.rejectionReason,
      offerLetterUrl: response.offerLetterUrl,
      offerValidUntil: response.offerValidUntil,
      offerAccepted: response.offerAccepted,
      offerAcceptedAt: response.offerAcceptedAt,
      remarks: response.remarks,
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  }

  /**
   * Get merit list for a session and class
   */
  getMeritList(sessionId: string, className: string): Observable<MeritList | null> {
    return this.apiService.get<MeritListApiResponse>(
      `${this.basePathSessions}/${sessionId}/merit-list`,
      { params: { className } }
    ).pipe(
      map(response => this.mapMeritListResponse(response)),
      catchError((error: HttpErrorResponse) => {
        // 404 means no merit list exists yet - this is expected
        if (error.status === 404) {
          return of(null);
        }
        // For other errors, return null but log
        console.error('Failed to fetch merit list:', error);
        return of(null);
      })
    );
  }

  private getMeritListMock(sessionId: string, className: string): Observable<MeritList | null> {
    const key = `${sessionId}-${className}`;
    const existing = this.mockMeritLists.get(key);
    return of(existing ? { ...existing } : null).pipe(delay(300));
  }

  /**
   * List all merit lists for a session
   */
  listMeritLists(sessionId: string): Observable<MeritList[]> {
    return this.apiService.get<{ meritLists: MeritListApiResponse[]; total: number }>(
      `${this.basePathSessions}/${sessionId}/merit-lists`
    ).pipe(
      map(response => response.meritLists.map(ml => this.mapMeritListResponse(ml))),
      catchError(() => of([]))
    );
  }

  /**
   * Generate merit list for a session
   */
  generateMeritList(sessionId: string, request: GenerateMeritListRequest): Observable<MeritList> {
    return this.apiService.post<MeritListApiResponse>(
      `${this.basePathSessions}/${sessionId}/merit-list`,
      request
    ).pipe(
      map(response => this.mapMeritListResponse(response)),
      catchError((error: HttpErrorResponse) => {
        // Extract error message from backend response
        const errorMessage = error.error?.error?.message ||
          error.error?.message ||
          'Failed to generate merit list';
        console.error('Failed to generate merit list:', error);
        return throwError(() => new Error(errorMessage));
      })
    );
  }

  private generateMeritListMock(sessionId: string, request: GenerateMeritListRequest): Observable<MeritList> {
    let entries = [...this.mockEntries].filter(e => e.classApplying === request.className);

    // If no entries for this class, use all entries as mock
    if (entries.length === 0) {
      entries = [...this.mockEntries].map(e => ({ ...e, classApplying: request.className }));
    }

    // Apply cutoff filter if provided
    if (request.cutoffScore !== undefined && request.cutoffScore > 0) {
      entries = entries.filter(e => e.percentage >= request.cutoffScore!);
    }

    // Re-rank after filtering
    entries = entries
      .sort((a, b) => b.totalScore - a.totalScore)
      .map((entry, index) => ({ ...entry, rank: index + 1 }));

    const meritList: MeritList = {
      id: `merit-${Date.now()}`,
      tenantId: 'tenant-1',
      sessionId,
      className: request.className,
      testId: request.testId,
      testName: 'Entrance Test 2026',
      generatedAt: new Date().toISOString(),
      cutoffScore: request.cutoffScore,
      entries,
      isFinal: false,
      createdAt: new Date().toISOString(),
    };

    const key = `${sessionId}-${request.className}`;
    this.mockMeritLists.set(key, meritList);

    return of({ ...meritList }).pipe(delay(500));
  }

  /**
   * Finalize a merit list
   */
  finalizeMeritList(meritListId: string): Observable<MeritList> {
    return this.apiService.post<MeritListApiResponse>(
      `${this.basePathMeritLists}/${meritListId}/finalize`,
      {}
    ).pipe(
      map(response => this.mapMeritListResponse(response))
    );
  }

  /**
   * Update cutoff score for a merit list
   */
  updateCutoff(meritListId: string, cutoffScore: number | null): Observable<MeritList> {
    return this.apiService.patch<MeritListApiResponse>(
      `${this.basePathMeritLists}/${meritListId}/cutoff`,
      { cutoffScore }
    ).pipe(
      map(response => this.mapMeritListResponse(response))
    );
  }

  /**
   * Make admission decision for a single application
   */
  makeDecision(applicationId: string, request: MakeDecisionRequest): Observable<AdmissionDecision> {
    // Map frontend decision type to backend
    const backendDecision = request.decision === 'selected' ? 'approved' : request.decision;

    return this.apiService.post<DecisionApiResponse>(
      `${this.basePathApplications}/${applicationId}/decision`,
      { ...request, decision: backendDecision }
    ).pipe(
      map(response => this.mapDecisionResponse(response)),
      catchError(() => {
        if (this.useMockData) {
          return this.makeDecisionMock(applicationId, request);
        }
        throw new Error('Failed to make decision');
      })
    );
  }

  private makeDecisionMock(applicationId: string, request: MakeDecisionRequest): Observable<AdmissionDecision> {
    const statusMap: Record<string, ApplicationStatus> = {
      selected: 'offer_sent',
      waitlisted: 'waitlisted',
      rejected: 'rejected',
    };

    // Update mock entry status
    const entryIndex = this.mockEntries.findIndex(e => e.applicationId === applicationId);
    if (entryIndex !== -1) {
      this.mockEntries[entryIndex] = {
        ...this.mockEntries[entryIndex],
        status: statusMap[request.decision],
        sectionAssigned: request.sectionAssigned,
        waitlistPosition: request.waitlistPosition,
        decisionDate: new Date().toISOString(),
        remarks: request.remarks,
      };
    }

    const decision: AdmissionDecision = {
      id: `decision-${Date.now()}`,
      tenantId: 'tenant-1',
      applicationId,
      decision: request.decision,
      decisionDate: new Date().toISOString().split('T')[0],
      sectionAssigned: request.sectionAssigned,
      waitlistPosition: request.waitlistPosition,
      rejectionReason: request.rejectionReason,
      offerValidUntil: request.offerValidUntil,
      remarks: request.remarks,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    return of(decision).pipe(delay(500));
  }

  /**
   * Make bulk admission decisions
   */
  makeBulkDecision(request: BulkDecisionRequest): Observable<BulkDecisionResponse> {
    // Map frontend decision type to backend
    const backendDecision = request.decision === 'selected' ? 'approved' : request.decision;

    return this.apiService.post<BulkDecisionResponse>(
      `${this.basePathApplications}/bulk-decision`,
      { ...request, decision: backendDecision }
    ).pipe(
      catchError(() => {
        if (this.useMockData) {
          return this.makeBulkDecisionMock(request);
        }
        throw new Error('Failed to make bulk decision');
      })
    );
  }

  private makeBulkDecisionMock(request: BulkDecisionRequest): Observable<BulkDecisionResponse> {
    const statusMap: Record<string, ApplicationStatus> = {
      selected: 'offer_sent',
      waitlisted: 'waitlisted',
      rejected: 'rejected',
    };

    const decisions: AdmissionDecision[] = request.applicationIds.map((appId, index) => {
      // Update mock entry status
      const entryIndex = this.mockEntries.findIndex(e => e.applicationId === appId);
      if (entryIndex !== -1) {
        this.mockEntries[entryIndex] = {
          ...this.mockEntries[entryIndex],
          status: statusMap[request.decision],
          sectionAssigned: request.sectionAssigned,
          waitlistPosition: request.decision === 'waitlisted' ? index + 1 : undefined,
          decisionDate: new Date().toISOString(),
          remarks: request.remarks,
        };
      }

      return {
        id: `decision-${Date.now()}-${index}`,
        tenantId: 'tenant-1',
        applicationId: appId,
        decision: request.decision,
        decisionDate: new Date().toISOString().split('T')[0],
        sectionAssigned: request.sectionAssigned,
        waitlistPosition: request.decision === 'waitlisted' ? index + 1 : undefined,
        offerValidUntil: request.offerValidUntil,
        remarks: request.remarks,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
    });

    const response: BulkDecisionResponse = {
      successful: decisions.length,
      failed: 0,
      decisions,
      errors: [],
    };

    return of(response).pipe(delay(500));
  }

  /**
   * Get decision for an application
   */
  getDecision(applicationId: string): Observable<AdmissionDecision | null> {
    return this.apiService.get<DecisionApiResponse>(
      `${this.basePathApplications}/${applicationId}/decision`
    ).pipe(
      map(response => this.mapDecisionResponse(response)),
      catchError(() => of(null))
    );
  }

  /**
   * Generate offer letter for an application
   */
  generateOfferLetter(applicationId: string, validUntil?: string): Observable<{ url: string; validUntil: string }> {
    return this.apiService.post<OfferLetterApiResponse>(
      `${this.basePathApplications}/${applicationId}/offer-letter`,
      { validUntil }
    ).pipe(
      map(response => ({ url: response.url, validUntil: response.validUntil })),
      catchError(() => {
        if (this.useMockData) {
          return this.generateOfferLetterMock(applicationId);
        }
        throw new Error('Failed to generate offer letter');
      })
    );
  }

  private generateOfferLetterMock(applicationId: string): Observable<{ url: string; validUntil: string }> {
    const url = `/api/v1/applications/${applicationId}/offer-letter.pdf`;
    const validUntil = new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
    return of({ url, validUntil }).pipe(delay(500));
  }

  /**
   * Accept an offer
   */
  acceptOffer(applicationId: string, request?: AcceptOfferRequest): Observable<AdmissionDecision> {
    return this.apiService.post<DecisionApiResponse>(
      `${this.basePathApplications}/${applicationId}/accept-offer`,
      request || {}
    ).pipe(
      map(response => this.mapDecisionResponse(response)),
      catchError(() => {
        if (this.useMockData) {
          return this.acceptOfferMock(applicationId);
        }
        throw new Error('Failed to accept offer');
      })
    );
  }

  private acceptOfferMock(applicationId: string): Observable<AdmissionDecision> {
    // Update mock entry status
    const entryIndex = this.mockEntries.findIndex(e => e.applicationId === applicationId);
    if (entryIndex !== -1) {
      this.mockEntries[entryIndex] = {
        ...this.mockEntries[entryIndex],
        status: 'offer_accepted',
      };
    }

    const decision: AdmissionDecision = {
      id: `decision-${Date.now()}`,
      tenantId: 'tenant-1',
      applicationId,
      decision: 'selected',
      decisionDate: new Date().toISOString().split('T')[0],
      offerAccepted: true,
      offerAcceptedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    return of(decision).pipe(delay(500));
  }

  /**
   * Complete enrollment for an application
   */
  enroll(applicationId: string, request?: EnrollmentRequest): Observable<EnrollmentResponse> {
    return this.apiService.post<EnrollmentApiResponse>(
      `${this.basePathApplications}/${applicationId}/enroll`,
      request || {}
    ).pipe(
      map(response => ({
        applicationId: response.applicationId,
        studentId: response.studentId,
        enrollmentNumber: response.enrollmentNumber,
        admissionDate: response.admissionDate,
        className: response.className,
        sectionName: response.sectionName,
        rollNumber: response.rollNumber,
      })),
      catchError(() => {
        if (this.useMockData) {
          return this.enrollMock(applicationId, request);
        }
        throw new Error('Failed to complete enrollment');
      })
    );
  }

  private enrollMock(applicationId: string, request?: EnrollmentRequest): Observable<EnrollmentResponse> {
    // Update mock entry status
    const entryIndex = this.mockEntries.findIndex(e => e.applicationId === applicationId);
    let entry: MeritListEntry | undefined;
    if (entryIndex !== -1) {
      this.mockEntries[entryIndex] = {
        ...this.mockEntries[entryIndex],
        status: 'enrolled',
      };
      entry = this.mockEntries[entryIndex];
    }

    const response: EnrollmentResponse = {
      applicationId,
      studentId: `student-${Date.now()}`,
      enrollmentNumber: `ENR-2026-${String(Math.floor(Math.random() * 1000)).padStart(4, '0')}`,
      admissionDate: request?.admissionDate || new Date().toISOString().split('T')[0],
      className: entry?.classApplying || 'Class 1',
      sectionName: request?.sectionId || entry?.sectionAssigned,
      rollNumber: request?.rollNumber,
    };

    return of(response).pipe(delay(500));
  }

  /**
   * Promote from waitlist to approved
   */
  promoteFromWaitlist(applicationId: string, sectionAssigned?: string): Observable<AdmissionDecision> {
    return this.apiService.post<DecisionApiResponse>(
      `${this.basePathApplications}/${applicationId}/promote`,
      { sectionAssigned }
    ).pipe(
      map(response => this.mapDecisionResponse(response))
    );
  }

  /**
   * Update waitlist position
   */
  updateWaitlistPosition(applicationId: string, position: number): Observable<AdmissionDecision> {
    return this.apiService.patch<DecisionApiResponse>(
      `${this.basePathApplications}/${applicationId}/waitlist-position`,
      { position }
    ).pipe(
      map(response => this.mapDecisionResponse(response))
    );
  }

  /**
   * Get available classes for merit list generation
   */
  getAvailableClasses(sessionId: string): Observable<string[]> {
    // For now, return a static list - could be enhanced to fetch from session seats
    return of(CLASS_NAMES.slice(0, 8)).pipe(delay(200));
  }

  /**
   * Get available sections for a class
   */
  getSections(className: string): Observable<{ id: string; name: string }[]> {
    // For now, return a static list - could be enhanced to fetch from class configuration
    const sections = ['Section A', 'Section B', 'Section C', 'Section D'].map((name, i) => ({
      id: `section-${i + 1}`,
      name,
    }));
    return of(sections).pipe(delay(200));
  }
}
