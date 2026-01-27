/**
 * MSLS Entrance Test Service
 *
 * HTTP service for entrance test management API calls.
 */

import { Injectable, inject, signal } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { map, tap, catchError } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  EntranceTest,
  TestRegistration,
  CreateTestDto,
  UpdateTestDto,
  RegisterCandidateDto,
  SubmitResultsDto,
  BulkResultsDto,
  TestFilterParams,
  TestSubject,
  SubjectMarks,
  TestStatus,
  RegistrationStatus,
} from './entrance-test.model';

/**
 * API Response interfaces matching backend DTOs
 */
interface TestApiResponse {
  id: string;
  sessionId: string;
  testName: string;
  testDate: string;
  startTime: string;
  durationMinutes: number;
  venue?: string;
  classNames: string[];
  maxCandidates: number;
  status: string;
  subjects: { subject: string; maxMarks: number | string }[];
  registeredCount?: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * API Request interfaces matching backend DTOs
 */
interface CreateTestApiRequest {
  sessionId: string;
  testName: string;
  testDate: string;
  startTime: string;
  durationMinutes: number;
  venue?: string;
  classNames: string[];
  maxCandidates?: number;
  subjects: { subject: string; maxMarks: number }[];
}

interface TestListApiResponse {
  tests: TestApiResponse[];
  total: number;
}

interface RegistrationApiResponse {
  id: string;
  testId: string;
  applicationId: string;
  rollNumber?: string;
  status: string;
  marks?: { subjectName: string; maxMarks: number; obtainedMarks: number }[];
  totalMarks?: number;
  maxMarks?: number;
  percentage?: number;
  result?: string;
  remarks?: string;
  createdAt: string;
  updatedAt: string;
  application?: {
    id: string;
    applicationNumber: string;
    studentName: string;
    parentName?: string;
    parentPhone?: string;
    classApplying: string;
    status: string;
  };
}

/**
 * EntranceTestService - Handles all entrance test-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class EntranceTestService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/entrance-tests';

  /** Loading state */
  private readonly _loading = signal(false);
  readonly loading = this._loading.asReadonly();

  /** Current tests list */
  private readonly _tests = signal<EntranceTest[]>([]);
  readonly tests = this._tests.asReadonly();

  /** Current test registrations */
  private readonly _registrations = signal<TestRegistration[]>([]);
  readonly registrations = this._registrations.asReadonly();

  /** Selected test */
  private readonly _selectedTest = signal<EntranceTest | null>(null);
  readonly selectedTest = this._selectedTest.asReadonly();

  /**
   * Map API response to EntranceTest model
   * Backend returns 'subject' field, frontend uses 'name'
   */
  private mapTestResponse(response: TestApiResponse): EntranceTest {
    return {
      id: response.id,
      tenantId: '',
      sessionId: response.sessionId,
      testName: response.testName,
      testDate: response.testDate,
      startTime: response.startTime,
      durationMinutes: response.durationMinutes,
      venue: response.venue,
      classNames: response.classNames,
      maxCandidates: response.maxCandidates,
      status: response.status as TestStatus,
      subjects: (response.subjects || []).map(s => ({
        name: s.subject,
        maxMarks: typeof s.maxMarks === 'string' ? parseFloat(s.maxMarks) : s.maxMarks,
      })),
      registeredCount: response.registeredCount,
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  }

  /**
   * Transform CreateTestDto to API request format
   * Frontend uses 'name' for subjects, backend expects 'subject'
   */
  private toApiRequest(data: CreateTestDto): CreateTestApiRequest {
    return {
      sessionId: data.sessionId,
      testName: data.testName,
      testDate: data.testDate,
      startTime: data.startTime,
      durationMinutes: data.durationMinutes,
      venue: data.venue,
      classNames: data.classNames,
      maxCandidates: data.maxCandidates,
      subjects: data.subjects.map(s => ({
        subject: s.name,
        maxMarks: s.maxMarks,
      })),
    };
  }

  /**
   * Map API response to TestRegistration model
   */
  private mapRegistrationResponse(response: RegistrationApiResponse): TestRegistration {
    return {
      id: response.id,
      tenantId: '',
      testId: response.testId,
      applicationId: response.applicationId,
      rollNumber: response.rollNumber,
      status: response.status as RegistrationStatus,
      marks: response.marks || [],
      totalMarks: response.totalMarks,
      maxMarks: response.maxMarks,
      percentage: response.percentage,
      result: response.result as 'pass' | 'fail' | 'pending' | undefined,
      remarks: response.remarks,
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
      application: response.application ? {
        id: response.application.id,
        applicationNumber: response.application.applicationNumber,
        studentName: response.application.studentName,
        parentName: response.application.parentName || '',
        parentPhone: response.application.parentPhone || '',
        classApplying: response.application.classApplying,
        status: response.application.status,
      } : undefined,
    };
  }

  /**
   * Get all entrance tests
   */
  getTests(filters?: TestFilterParams): Observable<EntranceTest[]> {
    this._loading.set(true);
    const params: Record<string, string> = {};
    if (filters?.sessionId) params['sessionId'] = filters.sessionId;
    if (filters?.status) params['status'] = filters.status;
    if (filters?.className) params['className'] = filters.className;

    return this.apiService.get<TestListApiResponse>(this.basePath, { params }).pipe(
      map(response => (response.tests || []).map(t => this.mapTestResponse(t))),
      tap(tests => {
        this._tests.set(tests);
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Get a single test by ID
   */
  getTest(id: string): Observable<EntranceTest> {
    this._loading.set(true);
    return this.apiService.get<TestApiResponse>(`${this.basePath}/${id}`).pipe(
      map(response => this.mapTestResponse(response)),
      tap(test => {
        this._selectedTest.set(test);
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Create a new entrance test
   */
  createTest(data: CreateTestDto): Observable<EntranceTest> {
    this._loading.set(true);
    const apiRequest = this.toApiRequest(data);
    return this.apiService.post<TestApiResponse>(this.basePath, apiRequest).pipe(
      map(response => this.mapTestResponse(response)),
      tap(test => {
        this._tests.update(tests => [test, ...tests]);
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Update an existing test
   */
  updateTest(id: string, data: UpdateTestDto): Observable<EntranceTest> {
    this._loading.set(true);
    // Transform subjects if present
    const apiRequest: Record<string, unknown> = { ...data };
    if (data.subjects) {
      apiRequest['subjects'] = data.subjects.map(s => ({
        subject: s.name,
        maxMarks: s.maxMarks,
      }));
    }
    return this.apiService.put<TestApiResponse>(`${this.basePath}/${id}`, apiRequest).pipe(
      map(response => this.mapTestResponse(response)),
      tap(updated => {
        this._tests.update(tests => tests.map(t => t.id === id ? updated : t));
        this._selectedTest.set(updated);
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Delete an entrance test
   */
  deleteTest(id: string): Observable<void> {
    this._loading.set(true);
    return this.apiService.delete<void>(`${this.basePath}/${id}`).pipe(
      tap(() => {
        this._tests.update(tests => tests.filter(t => t.id !== id));
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Get registrations for a test
   */
  getRegistrations(testId: string): Observable<TestRegistration[]> {
    this._loading.set(true);
    return this.apiService.get<RegistrationApiResponse[]>(`${this.basePath}/${testId}/registrations`).pipe(
      map(response => (response || []).map(r => this.mapRegistrationResponse(r))),
      tap(regs => {
        this._registrations.set(regs);
        this._loading.set(false);
      }),
      catchError(err => {
        this._loading.set(false);
        return throwError(() => err);
      })
    );
  }

  /**
   * Register a candidate for a test
   */
  registerCandidate(testId: string, data: RegisterCandidateDto): Observable<TestRegistration> {
    return this.apiService.post<RegistrationApiResponse>(`${this.basePath}/${testId}/register`, data).pipe(
      map(response => this.mapRegistrationResponse(response)),
      tap(reg => {
        this._registrations.update(regs => [...regs, reg]);
        // Update test registered count
        this._tests.update(tests => tests.map(t => {
          if (t.id === testId) {
            return { ...t, registeredCount: (t.registeredCount || 0) + 1 };
          }
          return t;
        }));
      })
    );
  }

  /**
   * Submit results for a registration
   */
  submitResult(testId: string, data: SubmitResultsDto): Observable<TestRegistration> {
    // Convert marks array to object format expected by backend
    const marksObj: Record<string, number> = {};
    data.marks.forEach(m => {
      marksObj[m.subjectName] = m.obtainedMarks;
    });

    const payload = {
      registrationId: data.registrationId,
      marks: marksObj,
      remarks: data.remarks,
    };

    return this.apiService.post<RegistrationApiResponse>(`${this.basePath}/${testId}/results`, payload).pipe(
      map(response => this.mapRegistrationResponse(response)),
      tap(reg => {
        this._registrations.update(regs => regs.map(r => r.id === reg.id ? reg : r));
      })
    );
  }

  /**
   * Submit bulk results
   */
  submitBulkResults(testId: string, data: BulkResultsDto): Observable<TestRegistration[]> {
    // Convert marks arrays to object format
    const results = data.results.map(r => {
      const marksObj: Record<string, number> = {};
      r.marks.forEach(m => {
        marksObj[m.subjectName] = m.obtainedMarks;
      });
      return {
        registrationId: r.registrationId,
        marks: marksObj,
        remarks: r.remarks,
      };
    });

    return this.apiService.post<RegistrationApiResponse[]>(
      `${this.basePath}/${testId}/results/bulk`,
      { results }
    ).pipe(
      map(response => (response || []).map(r => this.mapRegistrationResponse(r))),
      tap(regs => {
        this._registrations.update(current => {
          const regMap = new Map(regs.map(r => [r.id, r]));
          return current.map(r => regMap.get(r.id) || r);
        });
      })
    );
  }

  /**
   * Generate hall tickets for a test
   */
  generateHallTickets(testId: string): Observable<Blob> {
    return this.apiService.get<Blob>(`${this.basePath}/${testId}/hall-tickets`, {
      responseType: 'blob' as 'json'
    });
  }

  /**
   * Mark test as completed
   */
  completeTest(testId: string): Observable<EntranceTest> {
    return this.updateTest(testId, { status: 'completed' } as UpdateTestDto);
  }

  /**
   * Cancel a test
   */
  cancelTest(testId: string): Observable<EntranceTest> {
    return this.updateTest(testId, { status: 'cancelled' } as UpdateTestDto);
  }
}
