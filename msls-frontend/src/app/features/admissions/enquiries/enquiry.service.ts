/**
 * MSLS Enquiry Service
 *
 * HTTP service for managing admission enquiries.
 * Handles CRUD operations, follow-ups, and conversion to applications.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, of, finalize } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiService } from '../../../core/services/api.service';
import { PaginationParams } from '../../../core/models/api-response.model';
import {
  Enquiry,
  EnquiryFollowUp,
  CreateEnquiryDto,
  UpdateEnquiryDto,
  CreateFollowUpDto,
  ConvertEnquiryDto,
  EnquiryFilterParams,
} from './enquiry.model';

/** Backend response format for enquiry list */
interface EnquiryListResponse {
  enquiries: Enquiry[];
  total: number;
  page: number;
  pageSize: number;
}

/** Backend response format for follow-up list */
interface FollowUpListResponse {
  followUps: EnquiryFollowUp[];
  total: number;
}

/** API endpoints */
const ENDPOINTS = {
  enquiries: '/enquiries',
  enquiry: (id: string) => `/v1/enquiries/${id}`,
  followUps: (id: string) => `/v1/enquiries/${id}/follow-ups`,
  convert: (id: string) => `/v1/enquiries/${id}/convert`,
};

/**
 * EnquiryService - Manages admission enquiries data and operations.
 */
@Injectable({ providedIn: 'root' })
export class EnquiryService {
  private readonly api = inject(ApiService);

  // State signals
  private readonly _enquiries = signal<Enquiry[]>([]);
  private readonly _selectedEnquiry = signal<Enquiry | null>(null);
  private readonly _followUps = signal<EnquiryFollowUp[]>([]);
  private readonly _loading = signal<boolean>(false);
  private readonly _error = signal<string | null>(null);
  private readonly _totalItems = signal<number>(0);
  private readonly _currentPage = signal<number>(1);
  private readonly _pageSize = signal<number>(20);

  // Public read-only signals
  readonly enquiries = this._enquiries.asReadonly();
  readonly selectedEnquiry = this._selectedEnquiry.asReadonly();
  readonly followUps = this._followUps.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();
  readonly totalItems = this._totalItems.asReadonly();
  readonly currentPage = this._currentPage.asReadonly();
  readonly pageSize = this._pageSize.asReadonly();

  // Computed values
  readonly totalPages = computed(() =>
    Math.ceil(this._totalItems() / this._pageSize())
  );

  readonly hasNextPage = computed(() =>
    this._currentPage() < this.totalPages()
  );

  readonly hasPreviousPage = computed(() =>
    this._currentPage() > 1
  );

  /**
   * Load enquiries with optional filters and pagination
   */
  loadEnquiries(
    filters?: EnquiryFilterParams,
    pagination?: PaginationParams
  ): Observable<Enquiry[]> {
    this._loading.set(true);
    this._error.set(null);

    const params = this.buildQueryParams(filters);
    if (pagination?.page) params['page'] = String(pagination.page);
    if (pagination?.pageSize) params['pageSize'] = String(pagination.pageSize);

    return this.api
      .get<EnquiryListResponse>(ENDPOINTS.enquiries, { params })
      .pipe(
        map((response) => {
          const enquiries = response.enquiries || [];
          this._enquiries.set(enquiries);
          this._totalItems.set(response.total || 0);
          this._currentPage.set(response.page || 1);
          this._pageSize.set(response.pageSize || 20);
          return enquiries;
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to load enquiries');
          this._enquiries.set([]);
          return of([]);
        }),
        finalize(() => this._loading.set(false))
      );
  }

  /**
   * Get a single enquiry by ID
   */
  getEnquiry(id: string): Observable<Enquiry> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<Enquiry>(ENDPOINTS.enquiry(id)).pipe(
      tap((enquiry) => {
        this._selectedEnquiry.set(enquiry);
      }),
      catchError((err) => {
        this._error.set(err.message || 'Failed to load enquiry');
        throw err;
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Create a new enquiry
   */
  createEnquiry(dto: CreateEnquiryDto): Observable<Enquiry> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.post<Enquiry>(ENDPOINTS.enquiries, dto).pipe(
      tap((enquiry) => {
        this._enquiries.update((list) => [enquiry, ...list]);
        this._totalItems.update((n) => n + 1);
      }),
      catchError((err) => {
        this._error.set(err.message || 'Failed to create enquiry');
        throw err;
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Update an existing enquiry
   */
  updateEnquiry(id: string, dto: UpdateEnquiryDto): Observable<Enquiry> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.put<Enquiry>(ENDPOINTS.enquiry(id), dto).pipe(
      tap((enquiry) => {
        this._enquiries.update((list) =>
          list.map((e) => (e.id === id ? enquiry : e))
        );
        if (this._selectedEnquiry()?.id === id) {
          this._selectedEnquiry.set(enquiry);
        }
      }),
      catchError((err) => {
        this._error.set(err.message || 'Failed to update enquiry');
        throw err;
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Delete an enquiry
   */
  deleteEnquiry(id: string): Observable<void> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.delete<void>(ENDPOINTS.enquiry(id)).pipe(
      tap(() => {
        this._enquiries.update((list) => list.filter((e) => e.id !== id));
        this._totalItems.update((n) => Math.max(0, n - 1));
        if (this._selectedEnquiry()?.id === id) {
          this._selectedEnquiry.set(null);
        }
      }),
      catchError((err) => {
        this._error.set(err.message || 'Failed to delete enquiry');
        throw err;
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Load follow-ups for an enquiry
   */
  loadFollowUps(enquiryId: string): Observable<EnquiryFollowUp[]> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<FollowUpListResponse>(ENDPOINTS.followUps(enquiryId)).pipe(
      map((response) => {
        const followUps = response.followUps || [];
        this._followUps.set(followUps);
        return followUps;
      }),
      catchError((err) => {
        this._error.set(err.message || 'Failed to load follow-ups');
        this._followUps.set([]);
        return of([]);
      }),
      finalize(() => this._loading.set(false))
    );
  }

  /**
   * Add a follow-up to an enquiry
   */
  addFollowUp(
    enquiryId: string,
    dto: CreateFollowUpDto
  ): Observable<EnquiryFollowUp> {
    this._loading.set(true);
    this._error.set(null);

    return this.api
      .post<EnquiryFollowUp>(ENDPOINTS.followUps(enquiryId), dto)
      .pipe(
        tap((followUp) => {
          this._followUps.update((list) => [followUp, ...list]);
          // Update the enquiry's follow-up date if nextFollowUp is set
          if (dto.nextFollowUp) {
            this._enquiries.update((list) =>
              list.map((e) =>
                e.id === enquiryId
                  ? { ...e, followUpDate: dto.nextFollowUp }
                  : e
              )
            );
          }
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to add follow-up');
          throw err;
        }),
        finalize(() => this._loading.set(false))
      );
  }

  /**
   * Convert an enquiry to an application
   */
  convertToApplication(
    enquiryId: string,
    dto?: ConvertEnquiryDto
  ): Observable<{ applicationId: string }> {
    this._loading.set(true);
    this._error.set(null);

    return this.api
      .post<{ applicationId: string }>(ENDPOINTS.convert(enquiryId), dto || {})
      .pipe(
        tap((result) => {
          // Update enquiry status to converted
          this._enquiries.update((list) =>
            list.map((e) =>
              e.id === enquiryId
                ? {
                    ...e,
                    status: 'converted' as const,
                    convertedApplicationId: result.applicationId,
                  }
                : e
            )
          );
          if (this._selectedEnquiry()?.id === enquiryId) {
            this._selectedEnquiry.update((e) =>
              e
                ? {
                    ...e,
                    status: 'converted' as const,
                    convertedApplicationId: result.applicationId,
                  }
                : null
            );
          }
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to convert enquiry');
          throw err;
        }),
        finalize(() => this._loading.set(false))
      );
  }

  /**
   * Update enquiry status
   */
  updateStatus(id: string, status: Enquiry['status']): Observable<Enquiry> {
    return this.updateEnquiry(id, { status });
  }

  /**
   * Select an enquiry for viewing/editing
   */
  selectEnquiry(enquiry: Enquiry | null): void {
    this._selectedEnquiry.set(enquiry);
    if (enquiry) {
      this._followUps.set([]);
    }
  }

  /**
   * Clear error state
   */
  clearError(): void {
    this._error.set(null);
  }

  /**
   * Build query params from filter object
   */
  private buildQueryParams(
    filters?: EnquiryFilterParams
  ): Record<string, string> {
    if (!filters) return {};

    const params: Record<string, string> = {};

    if (filters.status) params['status'] = filters.status;
    if (filters.classApplying) params['classApplying'] = filters.classApplying;
    if (filters.source) params['source'] = filters.source;
    if (filters.fromDate) params['startDate'] = filters.fromDate;
    if (filters.toDate) params['endDate'] = filters.toDate;
    if (filters.search) params['search'] = filters.search;
    if (filters.assignedTo) params['assignedTo'] = filters.assignedTo;
    if (filters.branchId) params['branchId'] = filters.branchId;
    if (filters.sessionId) params['sessionId'] = filters.sessionId;

    return params;
  }
}
