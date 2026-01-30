/**
 * MSLS Student Service
 *
 * Handles all student-related HTTP API calls with reactive state management.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, throwError, finalize } from 'rxjs';

import { ApiService } from '../../../core/services/api.service';
import {
  Student,
  CreateStudentRequest,
  UpdateStudentRequest,
  StudentListFilter,
  StudentListResponse,
  StudentStatus,
  BulkOperation,
  BulkStatusUpdateRequest,
  ExportRequest,
  ImportResult,
} from '../models/student.model';

/** API endpoint for students */
const STUDENTS_ENDPOINT = '/students';

/**
 * StudentService - Manages student data with reactive signals.
 *
 * Provides CRUD operations for students with built-in loading and error states.
 *
 * Usage:
 * ```typescript
 * private studentService = inject(StudentService);
 *
 * // Access reactive state
 * students = this.studentService.students;
 * loading = this.studentService.loading;
 *
 * // Fetch students
 * this.studentService.loadStudents().subscribe();
 * ```
 */
@Injectable({ providedIn: 'root' })
export class StudentService {
  private api = inject(ApiService);

  // =========================================================================
  // State Signals
  // =========================================================================

  /** List of students */
  private _students = signal<Student[]>([]);

  /** Currently selected student */
  private _selectedStudent = signal<Student | null>(null);

  /** Loading state */
  private _loading = signal<boolean>(false);

  /** Error state */
  private _error = signal<string | null>(null);

  /** Whether there are more results */
  private _hasMore = signal<boolean>(false);

  /** Total count of students */
  private _totalCount = signal<number>(0);

  /** Current filters */
  private _currentFilter = signal<StudentListFilter>({});

  /** Cursor for next page (UUID of last student) */
  private _nextCursor = signal<string | null>(null);

  // =========================================================================
  // Public Readonly Signals
  // =========================================================================

  /** Public readonly students list */
  readonly students = this._students.asReadonly();

  /** Public readonly selected student */
  readonly selectedStudent = this._selectedStudent.asReadonly();

  /** Public readonly loading state */
  readonly loading = this._loading.asReadonly();

  /** Public readonly error state */
  readonly error = this._error.asReadonly();

  /** Public readonly has more flag */
  readonly hasMoreResults = this._hasMore.asReadonly();

  /** Public readonly total count */
  readonly totalCount = this._totalCount.asReadonly();

  /** Whether there are more pages to load */
  readonly hasMore = computed(() => this._hasMore());

  /** Whether the list is empty */
  readonly isEmpty = computed(() => this._students().length === 0 && !this._loading());

  // =========================================================================
  // CRUD Operations
  // =========================================================================

  /**
   * Load students with optional filtering and pagination.
   * @param filter - Filter options
   * @param append - Whether to append to existing list (for infinite scroll)
   */
  loadStudents(filter?: StudentListFilter, append = false): Observable<StudentListResponse> {
    this._loading.set(true);
    this._error.set(null);

    const mergedFilter = { ...this._currentFilter(), ...filter };
    this._currentFilter.set(mergedFilter);

    const params: Record<string, string | number> = {};
    if (mergedFilter.branchId) params['branch_id'] = mergedFilter.branchId;
    if (mergedFilter.classId) params['class_id'] = mergedFilter.classId;
    if (mergedFilter.sectionId) params['section_id'] = mergedFilter.sectionId;
    if (mergedFilter.status) params['status'] = mergedFilter.status;
    if (mergedFilter.gender) params['gender'] = mergedFilter.gender;
    if (mergedFilter.admissionFrom) params['admission_from'] = mergedFilter.admissionFrom;
    if (mergedFilter.admissionTo) params['admission_to'] = mergedFilter.admissionTo;
    if (mergedFilter.search) params['search'] = mergedFilter.search;
    if (mergedFilter.cursor) params['cursor'] = mergedFilter.cursor;
    if (mergedFilter.limit) params['limit'] = mergedFilter.limit;
    if (mergedFilter.sortBy) params['sort_by'] = mergedFilter.sortBy;
    if (mergedFilter.sortOrder) params['sort_order'] = mergedFilter.sortOrder;

    return this.api.get<StudentListResponse>(STUDENTS_ENDPOINT, { params }).pipe(
      tap((response) => {
        if (append) {
          this._students.update((current) => [...current, ...response.students]);
        } else {
          this._students.set(response.students);
        }
        this._hasMore.set(response.hasMore);
        this._totalCount.set(response.total);
        this._nextCursor.set(response.nextCursor ?? null);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load students');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Load more students (for infinite scroll/pagination).
   */
  loadMore(): Observable<StudentListResponse> | null {
    if (!this._hasMore() || !this._nextCursor()) return null;

    const limit = this._currentFilter().limit || 20;
    return this.loadStudents({ limit, cursor: this._nextCursor()! }, true);
  }

  /**
   * Refresh the student list with current filters.
   */
  refresh(): Observable<StudentListResponse> {
    const filter = { ...this._currentFilter(), cursor: undefined };
    return this.loadStudents(filter, false);
  }

  /**
   * Get a single student by ID.
   * @param id - Student ID
   */
  getStudent(id: string): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<Student>(`${STUDENTS_ENDPOINT}/${id}`).pipe(
      tap((student) => {
        this._selectedStudent.set(student);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load student');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Create a new student.
   * @param data - Student creation data
   */
  createStudent(data: CreateStudentRequest): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Student>(STUDENTS_ENDPOINT, data).pipe(
      tap((student) => {
        // Add to the beginning of the list
        this._students.update((current) => [student, ...current]);
        this._totalCount.update((count) => count + 1);
        this._selectedStudent.set(student);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to create student');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update an existing student.
   * @param id - Student ID
   * @param data - Student update data
   */
  updateStudent(id: string, data: UpdateStudentRequest): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.put<Student>(`${STUDENTS_ENDPOINT}/${id}`, data).pipe(
      tap((student) => {
        // Update in the list
        this._students.update((current) =>
          current.map((s) => (s.id === id ? student : s))
        );
        if (this._selectedStudent()?.id === id) {
          this._selectedStudent.set(student);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update student');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Delete a student.
   * @param id - Student ID
   */
  deleteStudent(id: string): Observable<void> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.delete<void>(`${STUDENTS_ENDPOINT}/${id}`).pipe(
      tap(() => {
        // Remove from the list
        this._students.update((current) => current.filter((s) => s.id !== id));
        this._totalCount.update((count) => Math.max(0, count - 1));
        if (this._selectedStudent()?.id === id) {
          this._selectedStudent.set(null);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to delete student');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update student status.
   * @param id - Student ID
   * @param status - New status
   */
  updateStatus(id: string, status: StudentStatus): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.patch<Student>(`${STUDENTS_ENDPOINT}/${id}/status`, { status }).pipe(
      tap((student) => {
        this._students.update((current) =>
          current.map((s) => (s.id === id ? student : s))
        );
        if (this._selectedStudent()?.id === id) {
          this._selectedStudent.set(student);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update status');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Upload student photo.
   * @param id - Student ID
   * @param file - Photo file
   */
  uploadPhoto(id: string, file: File): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.uploadFile<Student>(`${STUDENTS_ENDPOINT}/${id}/photo`, file, 'photo').pipe(
      tap((student) => {
        this._students.update((current) =>
          current.map((s) => (s.id === id ? student : s))
        );
        if (this._selectedStudent()?.id === id) {
          this._selectedStudent.set(student);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to upload photo');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Upload student document (birth certificate, etc.).
   * @param id - Student ID
   * @param file - Document file
   * @param documentType - Type of document
   */
  uploadDocument(id: string, file: File, documentType: string): Observable<Student> {
    this._loading.set(true);
    this._error.set(null);

    return this.api
      .uploadFile<Student>(`${STUDENTS_ENDPOINT}/${id}/documents`, file, 'document', {
        documentType,
      })
      .pipe(
        tap((student) => {
          this._students.update((current) =>
            current.map((s) => (s.id === id ? student : s))
          );
          if (this._selectedStudent()?.id === id) {
            this._selectedStudent.set(student);
          }
        }),
        catchError((error) => {
          this._error.set(error.message || 'Failed to upload document');
          return throwError(() => error);
        }),
        finalize(() => this._loading.set(false))
      );
  }

  /**
   * Get the next admission number preview.
   * @param branchId - Branch ID
   */
  getNextAdmissionNumber(branchId: string): Observable<{ admissionNumber: string }> {
    return this.api.get<{ admissionNumber: string }>(
      `${STUDENTS_ENDPOINT}/next-admission-number`,
      { params: { branchId } }
    );
  }

  // =========================================================================
  // State Management Helpers
  // =========================================================================

  /**
   * Select a student.
   * @param student - Student to select
   */
  selectStudent(student: Student | null): void {
    this._selectedStudent.set(student);
  }

  /**
   * Clear the selected student.
   */
  clearSelection(): void {
    this._selectedStudent.set(null);
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
    this._students.set([]);
    this._selectedStudent.set(null);
    this._loading.set(false);
    this._error.set(null);
    this._hasMore.set(false);
    this._totalCount.set(0);
    this._currentFilter.set({});
    this._nextCursor.set(null);
  }

  // =========================================================================
  // Bulk Operations
  // =========================================================================

  /**
   * Bulk update student status.
   * @param request - Bulk status update request
   */
  bulkUpdateStatus(request: BulkStatusUpdateRequest): Observable<BulkOperation> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<BulkOperation>(`${STUDENTS_ENDPOINT}/bulk/status`, request).pipe(
      tap((operation) => {
        // Refresh the list after bulk update
        if (operation.status === 'completed') {
          this.refresh().subscribe();
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update students');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Export students to file.
   * @param request - Export request
   */
  exportStudents(request: ExportRequest): Observable<BulkOperation> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<BulkOperation>(`${STUDENTS_ENDPOINT}/export`, request).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to export students');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get bulk operation status.
   * @param id - Operation ID
   */
  getBulkOperation(id: string): Observable<BulkOperation> {
    return this.api.get<BulkOperation>(`/bulk-operations/${id}`).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to get operation status');
        return throwError(() => error);
      })
    );
  }

  /**
   * List bulk operations for current user.
   * @param limit - Number of operations to fetch
   */
  listBulkOperations(limit = 20): Observable<BulkOperation[]> {
    return this.api.get<BulkOperation[]>('/bulk-operations', { params: { limit } }).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to list operations');
        return throwError(() => error);
      })
    );
  }

  /**
   * Download export result.
   * @param operationId - Operation ID
   */
  downloadExportResult(operationId: string): void {
    window.open(`/api/v1/bulk-operations/${operationId}/result`, '_blank');
  }

  // =========================================================================
  // Import Operations
  // =========================================================================

  /**
   * Download the student import template.
   */
  downloadImportTemplate(): void {
    this.api.downloadFile(`${STUDENTS_ENDPOINT}/import/template`, 'student_import_template.xlsx');
  }

  /**
   * Import students from a file.
   * @param file - Excel or CSV file
   * @param branchId - Branch ID
   * @param academicYearId - Academic year ID
   */
  importStudents(
    file: File,
    branchId: string,
    academicYearId: string
  ): Observable<ImportResult> {
    this._loading.set(true);
    this._error.set(null);

    const formData = new FormData();
    formData.append('file', file);
    formData.append('branch_id', branchId);
    formData.append('academic_year_id', academicYearId);

    return this.api.post<ImportResult>(`${STUDENTS_ENDPOINT}/import`, formData).pipe(
      tap(() => {
        // Refresh the list after import
        this.refresh().subscribe();
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to import students');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }
}
