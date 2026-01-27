/**
 * MSLS Enrollment Service
 *
 * Handles all enrollment-related HTTP API calls with reactive state management.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, throwError, finalize } from 'rxjs';

import { ApiService } from '../../../core/services/api.service';
import {
  Enrollment,
  CreateEnrollmentRequest,
  UpdateEnrollmentRequest,
  TransferRequest,
  DropoutRequest,
  EnrollmentHistoryResponse,
  EnrollmentStatusChange,
} from '../models/enrollment.model';

/**
 * EnrollmentService - Manages student enrollment data with reactive signals.
 *
 * Provides CRUD operations for enrollments with built-in loading and error states.
 *
 * Usage:
 * ```typescript
 * private enrollmentService = inject(EnrollmentService);
 *
 * // Access reactive state
 * enrollments = this.enrollmentService.enrollments;
 * loading = this.enrollmentService.loading;
 *
 * // Fetch enrollment history
 * this.enrollmentService.loadEnrollmentHistory(studentId).subscribe();
 * ```
 */
@Injectable({ providedIn: 'root' })
export class EnrollmentService {
  private api = inject(ApiService);

  // =========================================================================
  // State Signals
  // =========================================================================

  /** List of enrollments for current student */
  private _enrollments = signal<Enrollment[]>([]);

  /** Currently selected enrollment */
  private _selectedEnrollment = signal<Enrollment | null>(null);

  /** Current active enrollment */
  private _currentEnrollment = signal<Enrollment | null>(null);

  /** Status change history */
  private _statusHistory = signal<EnrollmentStatusChange[]>([]);

  /** Loading state */
  private _loading = signal<boolean>(false);

  /** Error state */
  private _error = signal<string | null>(null);

  /** Total count of enrollments */
  private _totalCount = signal<number>(0);

  // =========================================================================
  // Public Readonly Signals
  // =========================================================================

  /** Public readonly enrollments list */
  readonly enrollments = this._enrollments.asReadonly();

  /** Public readonly selected enrollment */
  readonly selectedEnrollment = this._selectedEnrollment.asReadonly();

  /** Public readonly current enrollment */
  readonly currentEnrollment = this._currentEnrollment.asReadonly();

  /** Public readonly status history */
  readonly statusHistory = this._statusHistory.asReadonly();

  /** Public readonly loading state */
  readonly loading = this._loading.asReadonly();

  /** Public readonly error state */
  readonly error = this._error.asReadonly();

  /** Public readonly total count */
  readonly totalCount = this._totalCount.asReadonly();

  /** Whether the student has an active enrollment */
  readonly hasActiveEnrollment = computed(() => this._currentEnrollment() !== null);

  /** Whether the enrollment list is empty */
  readonly isEmpty = computed(() => this._enrollments().length === 0 && !this._loading());

  // =========================================================================
  // CRUD Operations
  // =========================================================================

  /**
   * Load enrollment history for a student.
   * @param studentId - Student ID
   */
  loadEnrollmentHistory(studentId: string): Observable<EnrollmentHistoryResponse> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<EnrollmentHistoryResponse>(`/students/${studentId}/enrollments`).pipe(
      tap((response) => {
        this._enrollments.set(response.enrollments);
        this._totalCount.set(response.total);
        // Find the active enrollment
        const active = response.enrollments.find((e) => e.status === 'active');
        this._currentEnrollment.set(active ?? null);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load enrollment history');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get the current active enrollment for a student.
   * @param studentId - Student ID
   */
  getCurrentEnrollment(studentId: string): Observable<Enrollment> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<Enrollment>(`/students/${studentId}/enrollments/current`).pipe(
      tap((enrollment) => {
        this._currentEnrollment.set(enrollment);
      }),
      catchError((error) => {
        // Not having an active enrollment is not necessarily an error
        if (error.status === 404) {
          this._currentEnrollment.set(null);
        }
        this._error.set(error.message || 'Failed to load current enrollment');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Create a new enrollment for a student.
   * @param studentId - Student ID
   * @param data - Enrollment creation data
   */
  createEnrollment(studentId: string, data: CreateEnrollmentRequest): Observable<Enrollment> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Enrollment>(`/students/${studentId}/enrollments`, data).pipe(
      tap((enrollment) => {
        // Add to the beginning of the list
        this._enrollments.update((current) => [enrollment, ...current]);
        this._totalCount.update((count) => count + 1);
        if (enrollment.status === 'active') {
          this._currentEnrollment.set(enrollment);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to create enrollment');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update an existing enrollment.
   * @param studentId - Student ID
   * @param enrollmentId - Enrollment ID
   * @param data - Enrollment update data
   */
  updateEnrollment(studentId: string, enrollmentId: string, data: UpdateEnrollmentRequest): Observable<Enrollment> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.put<Enrollment>(`/students/${studentId}/enrollments/${enrollmentId}`, data).pipe(
      tap((enrollment) => {
        // Update in the list
        this._enrollments.update((current) =>
          current.map((e) => (e.id === enrollmentId ? enrollment : e))
        );
        if (this._selectedEnrollment()?.id === enrollmentId) {
          this._selectedEnrollment.set(enrollment);
        }
        if (this._currentEnrollment()?.id === enrollmentId) {
          this._currentEnrollment.set(enrollment);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update enrollment');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Process a student transfer.
   * @param studentId - Student ID
   * @param enrollmentId - Enrollment ID
   * @param data - Transfer data
   */
  processTransfer(studentId: string, enrollmentId: string, data: TransferRequest): Observable<Enrollment> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Enrollment>(`/students/${studentId}/enrollments/${enrollmentId}/transfer`, data).pipe(
      tap((enrollment) => {
        // Update in the list
        this._enrollments.update((current) =>
          current.map((e) => (e.id === enrollmentId ? enrollment : e))
        );
        if (this._selectedEnrollment()?.id === enrollmentId) {
          this._selectedEnrollment.set(enrollment);
        }
        // Clear current enrollment since it's no longer active
        if (this._currentEnrollment()?.id === enrollmentId) {
          this._currentEnrollment.set(null);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to process transfer');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Process a student dropout.
   * @param studentId - Student ID
   * @param enrollmentId - Enrollment ID
   * @param data - Dropout data
   */
  processDropout(studentId: string, enrollmentId: string, data: DropoutRequest): Observable<Enrollment> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Enrollment>(`/students/${studentId}/enrollments/${enrollmentId}/dropout`, data).pipe(
      tap((enrollment) => {
        // Update in the list
        this._enrollments.update((current) =>
          current.map((e) => (e.id === enrollmentId ? enrollment : e))
        );
        if (this._selectedEnrollment()?.id === enrollmentId) {
          this._selectedEnrollment.set(enrollment);
        }
        // Clear current enrollment since it's no longer active
        if (this._currentEnrollment()?.id === enrollmentId) {
          this._currentEnrollment.set(null);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to process dropout');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Load status change history for an enrollment.
   * @param studentId - Student ID
   * @param enrollmentId - Enrollment ID
   */
  loadStatusHistory(studentId: string, enrollmentId: string): Observable<EnrollmentStatusChange[]> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<EnrollmentStatusChange[]>(`/students/${studentId}/enrollments/${enrollmentId}/status-history`).pipe(
      tap((changes) => {
        this._statusHistory.set(changes);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load status history');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get enrollments by class.
   * @param classId - Class ID
   */
  getEnrollmentsByClass(classId: string): Observable<EnrollmentHistoryResponse> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<EnrollmentHistoryResponse>(`/enrollments/by-class/${classId}`).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to load enrollments by class');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get enrollments by section.
   * @param sectionId - Section ID
   */
  getEnrollmentsBySection(sectionId: string): Observable<EnrollmentHistoryResponse> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<EnrollmentHistoryResponse>(`/enrollments/by-section/${sectionId}`).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to load enrollments by section');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  // =========================================================================
  // State Management Helpers
  // =========================================================================

  /**
   * Select an enrollment.
   * @param enrollment - Enrollment to select
   */
  selectEnrollment(enrollment: Enrollment | null): void {
    this._selectedEnrollment.set(enrollment);
  }

  /**
   * Clear the selected enrollment.
   */
  clearSelection(): void {
    this._selectedEnrollment.set(null);
  }

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
    this._enrollments.set([]);
    this._selectedEnrollment.set(null);
    this._currentEnrollment.set(null);
    this._statusHistory.set([]);
    this._loading.set(false);
    this._error.set(null);
    this._totalCount.set(0);
  }
}
