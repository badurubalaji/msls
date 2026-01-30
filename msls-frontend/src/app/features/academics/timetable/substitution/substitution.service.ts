import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, tap, finalize } from 'rxjs';
import { environment } from '../../../../../environments/environment';
import {
  Substitution,
  SubstitutionListResponse,
  CreateSubstitutionRequest,
  UpdateSubstitutionRequest,
  SubstitutionFilter,
  AvailableTeachersResponse,
  TeacherPeriod,
} from './substitution.model';

@Injectable({
  providedIn: 'root',
})
export class SubstitutionService {
  private http = inject(HttpClient);
  private baseUrl = `${environment.apiUrl}/timetables/substitutions`;

  // State signals
  private _substitutions = signal<Substitution[]>([]);
  private _currentSubstitution = signal<Substitution | null>(null);
  private _total = signal<number>(0);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);
  private _availableTeachers = signal<AvailableTeachersResponse | null>(null);
  private _teacherPeriods = signal<TeacherPeriod[]>([]);

  // Public computed signals
  readonly substitutions = computed(() => this._substitutions());
  readonly currentSubstitution = computed(() => this._currentSubstitution());
  readonly total = computed(() => this._total());
  readonly loading = computed(() => this._loading());
  readonly error = computed(() => this._error());
  readonly availableTeachers = computed(() => this._availableTeachers());
  readonly teacherPeriods = computed(() => this._teacherPeriods());

  /**
   * List substitutions with optional filters
   */
  loadSubstitutions(filter: SubstitutionFilter = {}): Observable<SubstitutionListResponse> {
    this._loading.set(true);
    this._error.set(null);

    let params = new HttpParams();
    if (filter.branchId) params = params.set('branch_id', filter.branchId);
    if (filter.originalStaffId) params = params.set('original_staff_id', filter.originalStaffId);
    if (filter.substituteStaffId) params = params.set('substitute_staff_id', filter.substituteStaffId);
    if (filter.startDate) params = params.set('start_date', filter.startDate);
    if (filter.endDate) params = params.set('end_date', filter.endDate);
    if (filter.status) params = params.set('status', filter.status);
    if (filter.limit) params = params.set('limit', filter.limit.toString());
    if (filter.offset) params = params.set('offset', filter.offset.toString());

    return this.http.get<SubstitutionListResponse>(this.baseUrl, { params }).pipe(
      tap((response) => {
        this._substitutions.set(response.substitutions);
        this._total.set(response.total);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get a single substitution by ID
   */
  getSubstitution(id: string): Observable<Substitution> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.get<Substitution>(`${this.baseUrl}/${id}`).pipe(
      tap((substitution) => this._currentSubstitution.set(substitution)),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Create a new substitution
   */
  createSubstitution(request: CreateSubstitutionRequest): Observable<Substitution> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.post<Substitution>(this.baseUrl, request).pipe(
      tap((substitution) => {
        this._substitutions.update((list) => [substitution, ...list]);
        this._total.update((t) => t + 1);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update an existing substitution
   */
  updateSubstitution(id: string, request: UpdateSubstitutionRequest): Observable<Substitution> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.put<Substitution>(`${this.baseUrl}/${id}`, request).pipe(
      tap((substitution) => {
        this._substitutions.update((list) =>
          list.map((s) => (s.id === id ? substitution : s))
        );
        if (this._currentSubstitution()?.id === id) {
          this._currentSubstitution.set(substitution);
        }
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Confirm a pending substitution
   */
  confirmSubstitution(id: string): Observable<Substitution> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.post<Substitution>(`${this.baseUrl}/${id}/confirm`, {}).pipe(
      tap((substitution) => {
        this._substitutions.update((list) =>
          list.map((s) => (s.id === id ? substitution : s))
        );
        if (this._currentSubstitution()?.id === id) {
          this._currentSubstitution.set(substitution);
        }
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Cancel a substitution
   */
  cancelSubstitution(id: string): Observable<Substitution> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.post<Substitution>(`${this.baseUrl}/${id}/cancel`, {}).pipe(
      tap((substitution) => {
        this._substitutions.update((list) =>
          list.map((s) => (s.id === id ? substitution : s))
        );
        if (this._currentSubstitution()?.id === id) {
          this._currentSubstitution.set(substitution);
        }
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Delete a substitution
   */
  deleteSubstitution(id: string): Observable<void> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.delete<void>(`${this.baseUrl}/${id}`).pipe(
      tap(() => {
        this._substitutions.update((list) => list.filter((s) => s.id !== id));
        this._total.update((t) => t - 1);
        if (this._currentSubstitution()?.id === id) {
          this._currentSubstitution.set(null);
        }
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get available teachers for substitution
   */
  getAvailableTeachers(
    branchId: string,
    date: string,
    periodSlotIds: string[],
    excludeStaffId?: string
  ): Observable<AvailableTeachersResponse> {
    this._loading.set(true);
    this._error.set(null);

    let params = new HttpParams()
      .set('branch_id', branchId)
      .set('date', date)
      .set('period_slot_ids', periodSlotIds.join(','));

    if (excludeStaffId) {
      params = params.set('exclude_staff_id', excludeStaffId);
    }

    return this.http.get<AvailableTeachersResponse>(`${this.baseUrl}/available-teachers`, { params }).pipe(
      tap((response) => this._availableTeachers.set(response)),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Get teacher's timetable periods for a date
   */
  getTeacherPeriods(staffId: string, date: string): Observable<TeacherPeriod[]> {
    this._loading.set(true);
    this._error.set(null);

    const params = new HttpParams()
      .set('staff_id', staffId)
      .set('date', date);

    return this.http.get<TeacherPeriod[]>(`${this.baseUrl}/teacher-periods`, { params }).pipe(
      tap((periods) => this._teacherPeriods.set(periods)),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Clear current substitution
   */
  clearCurrentSubstitution(): void {
    this._currentSubstitution.set(null);
  }

  /**
   * Clear error
   */
  clearError(): void {
    this._error.set(null);
  }
}
