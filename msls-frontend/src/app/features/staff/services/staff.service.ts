/**
 * MSLS Staff Service
 *
 * Handles all staff-related HTTP API calls with reactive state management.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, throwError, finalize } from 'rxjs';

import { ApiService } from '../../../core/services/api.service';
import {
  Staff,
  CreateStaffRequest,
  UpdateStaffRequest,
  StatusUpdateRequest,
  StaffListFilter,
  StaffListResponse,
  StatusHistory,
} from '../models/staff.model';

/** API endpoint for staff */
const STAFF_ENDPOINT = '/staff';

/**
 * StaffService - Manages staff data with reactive signals.
 *
 * Provides CRUD operations for staff with built-in loading and error states.
 *
 * Usage:
 * ```typescript
 * private staffService = inject(StaffService);
 *
 * // Access reactive state
 * staff = this.staffService.staffList;
 * loading = this.staffService.loading;
 *
 * // Fetch staff
 * this.staffService.loadStaff().subscribe();
 * ```
 */
@Injectable({ providedIn: 'root' })
export class StaffService {
  private api = inject(ApiService);

  // =========================================================================
  // State Signals
  // =========================================================================

  /** List of staff */
  private _staffList = signal<Staff[]>([]);

  /** Currently selected staff */
  private _selectedStaff = signal<Staff | null>(null);

  /** Loading state */
  private _loading = signal<boolean>(false);

  /** Error state */
  private _error = signal<string | null>(null);

  /** Whether there are more results */
  private _hasMore = signal<boolean>(false);

  /** Total count of staff */
  private _totalCount = signal<number>(0);

  /** Current filters */
  private _currentFilter = signal<StaffListFilter>({});

  /** Cursor for next page (UUID of last staff) */
  private _nextCursor = signal<string | null>(null);

  /** Status history for currently selected staff */
  private _statusHistory = signal<StatusHistory[]>([]);

  // =========================================================================
  // Public Readonly Signals
  // =========================================================================

  /** Public readonly staff list */
  readonly staffList = this._staffList.asReadonly();

  /** Public readonly selected staff */
  readonly selectedStaff = this._selectedStaff.asReadonly();

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
  readonly isEmpty = computed(() => this._staffList().length === 0 && !this._loading());

  /** Public readonly status history */
  readonly statusHistory = this._statusHistory.asReadonly();

  // =========================================================================
  // CRUD Operations
  // =========================================================================

  /**
   * Load staff with optional filtering and pagination.
   * @param filter - Filter options
   * @param append - Whether to append to existing list (for infinite scroll)
   */
  loadStaff(filter?: StaffListFilter, append = false): Observable<StaffListResponse> {
    this._loading.set(true);
    this._error.set(null);

    const mergedFilter = { ...this._currentFilter(), ...filter };
    this._currentFilter.set(mergedFilter);

    const params: Record<string, string | number> = {};
    if (mergedFilter.branchId) params['branch_id'] = mergedFilter.branchId;
    if (mergedFilter.departmentId) params['department_id'] = mergedFilter.departmentId;
    if (mergedFilter.designationId) params['designation_id'] = mergedFilter.designationId;
    if (mergedFilter.staffType) params['staff_type'] = mergedFilter.staffType;
    if (mergedFilter.status) params['status'] = mergedFilter.status;
    if (mergedFilter.gender) params['gender'] = mergedFilter.gender;
    if (mergedFilter.joinDateFrom) params['join_date_from'] = mergedFilter.joinDateFrom;
    if (mergedFilter.joinDateTo) params['join_date_to'] = mergedFilter.joinDateTo;
    if (mergedFilter.search) params['search'] = mergedFilter.search;
    if (mergedFilter.cursor) params['cursor'] = mergedFilter.cursor;
    if (mergedFilter.limit) params['limit'] = mergedFilter.limit;
    if (mergedFilter.sortBy) params['sort_by'] = mergedFilter.sortBy;
    if (mergedFilter.sortOrder) params['sort_order'] = mergedFilter.sortOrder;

    return this.api.get<StaffListResponse>(STAFF_ENDPOINT, { params }).pipe(
      tap((response) => {
        if (append) {
          this._staffList.update((current) => [...current, ...response.staff]);
        } else {
          this._staffList.set(response.staff);
        }
        this._hasMore.set(response.hasMore);
        this._totalCount.set(response.total);
        this._nextCursor.set(response.nextCursor ?? null);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load staff');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Load more staff (for infinite scroll/pagination).
   */
  loadMore(): Observable<StaffListResponse> | null {
    if (!this._hasMore() || !this._nextCursor()) return null;

    const limit = this._currentFilter().limit || 20;
    return this.loadStaff({ limit, cursor: this._nextCursor()! }, true);
  }

  /**
   * Refresh the staff list with current filters.
   */
  refresh(): Observable<StaffListResponse> {
    const filter = { ...this._currentFilter(), cursor: undefined };
    return this.loadStaff(filter, false);
  }

  /**
   * Get a single staff member by ID.
   * @param id - Staff ID
   */
  getStaff(id: string): Observable<Staff> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<Staff>(`${STAFF_ENDPOINT}/${id}`).pipe(
      tap((staff) => {
        this._selectedStaff.set(staff);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load staff');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Create a new staff member.
   * @param data - Staff creation data
   */
  createStaff(data: CreateStaffRequest): Observable<Staff> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Staff>(STAFF_ENDPOINT, data).pipe(
      tap((staff) => {
        // Add to the beginning of the list
        this._staffList.update((current) => [staff, ...current]);
        this._totalCount.update((count) => count + 1);
        this._selectedStaff.set(staff);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to create staff');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update an existing staff member.
   * @param id - Staff ID
   * @param data - Staff update data
   */
  updateStaff(id: string, data: UpdateStaffRequest): Observable<Staff> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.put<Staff>(`${STAFF_ENDPOINT}/${id}`, data).pipe(
      tap((staff) => {
        // Update in the list
        this._staffList.update((current) =>
          current.map((s) => (s.id === id ? staff : s))
        );
        if (this._selectedStaff()?.id === id) {
          this._selectedStaff.set(staff);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to update staff');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Delete a staff member.
   * @param id - Staff ID
   */
  deleteStaff(id: string): Observable<void> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.delete<void>(`${STAFF_ENDPOINT}/${id}`).pipe(
      tap(() => {
        // Remove from the list
        this._staffList.update((current) => current.filter((s) => s.id !== id));
        this._totalCount.update((count) => Math.max(0, count - 1));
        if (this._selectedStaff()?.id === id) {
          this._selectedStaff.set(null);
        }
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to delete staff');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update staff status.
   * @param id - Staff ID
   * @param data - Status update data
   */
  updateStatus(id: string, data: StatusUpdateRequest): Observable<Staff> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.patch<Staff>(`${STAFF_ENDPOINT}/${id}/status`, data).pipe(
      tap((staff) => {
        this._staffList.update((current) =>
          current.map((s) => (s.id === id ? staff : s))
        );
        if (this._selectedStaff()?.id === id) {
          this._selectedStaff.set(staff);
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
   * Get status history for a staff member.
   * @param id - Staff ID
   */
  getStatusHistory(id: string): Observable<StatusHistory[]> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<StatusHistory[]>(`${STAFF_ENDPOINT}/${id}/status-history`).pipe(
      tap((history) => {
        this._statusHistory.set(history);
      }),
      catchError((error) => {
        this._error.set(error.message || 'Failed to load status history');
        return throwError(() => error);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Upload staff photo.
   * @param id - Staff ID
   * @param file - Photo file
   */
  uploadPhoto(id: string, file: File): Observable<Staff> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.uploadFile<Staff>(`${STAFF_ENDPOINT}/${id}/photo`, file, 'photo').pipe(
      tap((staff) => {
        this._staffList.update((current) =>
          current.map((s) => (s.id === id ? staff : s))
        );
        if (this._selectedStaff()?.id === id) {
          this._selectedStaff.set(staff);
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
   * Preview the next employee ID.
   */
  previewEmployeeId(): Observable<{ employeeId: string }> {
    return this.api.get<{ employeeId: string }>(`${STAFF_ENDPOINT}/employee-id/preview`).pipe(
      catchError((error) => {
        this._error.set(error.message || 'Failed to preview employee ID');
        return throwError(() => error);
      })
    );
  }

  // =========================================================================
  // State Management Helpers
  // =========================================================================

  /**
   * Select a staff member.
   * @param staff - Staff to select
   */
  selectStaff(staff: Staff | null): void {
    this._selectedStaff.set(staff);
  }

  /**
   * Clear the selected staff.
   */
  clearSelection(): void {
    this._selectedStaff.set(null);
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
    this._staffList.set([]);
    this._selectedStaff.set(null);
    this._loading.set(false);
    this._error.set(null);
    this._hasMore.set(false);
    this._totalCount.set(0);
    this._currentFilter.set({});
    this._nextCursor.set(null);
    this._statusHistory.set([]);
  }
}
