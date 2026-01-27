/**
 * Promotion Service
 *
 * Handles API calls for student promotion/retention processing.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, catchError, throwError } from 'rxjs';

import { environment } from '../../../../environments/environment';
import {
  PromotionRule,
  CreateRuleRequest,
  RuleListResponse,
  PromotionBatch,
  CreateBatchRequest,
  BatchListResponse,
  PromotionRecord,
  RecordListResponse,
  RecordsSummary,
  UpdateRecordRequest,
  BulkUpdateRecordsRequest,
  ProcessBatchRequest,
  CancelBatchRequest,
  PromotionReportResponse,
  BatchStatus,
} from '../models/promotion.model';

@Injectable({
  providedIn: 'root',
})
export class PromotionService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/api/v1`;

  // =========================================================================
  // State
  // =========================================================================

  /** Loading state */
  readonly loading = signal<boolean>(false);

  /** Error state */
  readonly error = signal<string | null>(null);

  /** Promotion rules */
  readonly rules = signal<PromotionRule[]>([]);

  /** Promotion batches */
  readonly batches = signal<PromotionBatch[]>([]);

  /** Currently selected batch */
  readonly selectedBatch = signal<PromotionBatch | null>(null);

  /** Records for the selected batch */
  readonly records = signal<PromotionRecord[]>([]);

  /** Records summary */
  readonly summary = signal<RecordsSummary | null>(null);

  // =========================================================================
  // Computed
  // =========================================================================

  readonly hasRules = computed(() => this.rules().length > 0);
  readonly hasBatches = computed(() => this.batches().length > 0);
  readonly hasRecords = computed(() => this.records().length > 0);

  readonly pendingRecords = computed(() =>
    this.records().filter((r) => r.decision === 'pending')
  );
  readonly promoteRecords = computed(() =>
    this.records().filter((r) => r.decision === 'promote')
  );
  readonly retainRecords = computed(() =>
    this.records().filter((r) => r.decision === 'retain')
  );
  readonly transferRecords = computed(() =>
    this.records().filter((r) => r.decision === 'transfer')
  );

  // =========================================================================
  // Promotion Rules
  // =========================================================================

  loadRules(): Observable<RuleListResponse> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.get<RuleListResponse>(`${this.apiUrl}/promotion-rules`).pipe(
      tap((response) => {
        this.rules.set(response.rules);
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to load promotion rules');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  createOrUpdateRule(request: CreateRuleRequest): Observable<PromotionRule> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.post<PromotionRule>(`${this.apiUrl}/promotion-rules`, request).pipe(
      tap((rule) => {
        // Update or add to local state
        const existingIndex = this.rules().findIndex((r) => r.classId === rule.classId);
        if (existingIndex >= 0) {
          const updated = [...this.rules()];
          updated[existingIndex] = rule;
          this.rules.set(updated);
        } else {
          this.rules.set([...this.rules(), rule]);
        }
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to save promotion rule');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  deleteRule(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.delete<void>(`${this.apiUrl}/promotion-rules/${id}`).pipe(
      tap(() => {
        this.rules.set(this.rules().filter((r) => r.id !== id));
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to delete promotion rule');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  // =========================================================================
  // Promotion Batches
  // =========================================================================

  loadBatches(status?: BatchStatus): Observable<BatchListResponse> {
    this.loading.set(true);
    this.error.set(null);

    let url = `${this.apiUrl}/promotion-batches`;
    if (status) {
      url += `?status=${status}`;
    }

    return this.http.get<BatchListResponse>(url).pipe(
      tap((response) => {
        this.batches.set(response.batches);
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to load promotion batches');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  loadBatch(id: string): Observable<PromotionBatch> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.get<PromotionBatch>(`${this.apiUrl}/promotion-batches/${id}`).pipe(
      tap((batch) => {
        this.selectedBatch.set(batch);
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to load promotion batch');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  createBatch(request: CreateBatchRequest): Observable<PromotionBatch> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.post<PromotionBatch>(`${this.apiUrl}/promotion-batches`, request).pipe(
      tap((batch) => {
        this.batches.set([batch, ...this.batches()]);
        this.selectedBatch.set(batch);
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to create promotion batch');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  deleteBatch(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.delete<void>(`${this.apiUrl}/promotion-batches/${id}`).pipe(
      tap(() => {
        this.batches.set(this.batches().filter((b) => b.id !== id));
        if (this.selectedBatch()?.id === id) {
          this.selectedBatch.set(null);
        }
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to delete promotion batch');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  cancelBatch(id: string, request: CancelBatchRequest): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return this.http.post<void>(`${this.apiUrl}/promotion-batches/${id}/cancel`, request).pipe(
      tap(() => {
        // Refresh batch
        this.loadBatch(id).subscribe();
        this.loading.set(false);
      }),
      catchError((err) => {
        this.error.set(err.error?.error || 'Failed to cancel promotion batch');
        this.loading.set(false);
        return throwError(() => err);
      })
    );
  }

  // =========================================================================
  // Promotion Records
  // =========================================================================

  loadRecords(batchId: string): Observable<RecordListResponse> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .get<RecordListResponse>(`${this.apiUrl}/promotion-batches/${batchId}/records`)
      .pipe(
        tap((response) => {
          this.records.set(response.records);
          this.summary.set(response.summary || null);
          this.loading.set(false);
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to load promotion records');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  updateRecord(
    batchId: string,
    recordId: string,
    request: UpdateRecordRequest
  ): Observable<PromotionRecord> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .put<PromotionRecord>(
        `${this.apiUrl}/promotion-batches/${batchId}/records/${recordId}`,
        request
      )
      .pipe(
        tap((record) => {
          // Update local state
          const updated = this.records().map((r) => (r.id === record.id ? record : r));
          this.records.set(updated);
          this.loading.set(false);
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to update promotion record');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  bulkUpdateRecords(batchId: string, request: BulkUpdateRecordsRequest): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .post<void>(`${this.apiUrl}/promotion-batches/${batchId}/records/bulk`, request)
      .pipe(
        tap(() => {
          // Reload records to get updated state
          this.loadRecords(batchId).subscribe();
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to bulk update records');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  // =========================================================================
  // Auto-Decision & Processing
  // =========================================================================

  autoDecide(batchId: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .post<void>(`${this.apiUrl}/promotion-batches/${batchId}/auto-decide`, {})
      .pipe(
        tap(() => {
          // Reload records and batch to get updated state
          this.loadRecords(batchId).subscribe();
          this.loadBatch(batchId).subscribe();
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to auto-decide records');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  processBatch(batchId: string, request: ProcessBatchRequest): Observable<PromotionBatch> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .post<PromotionBatch>(`${this.apiUrl}/promotion-batches/${batchId}/process`, request)
      .pipe(
        tap((batch) => {
          this.selectedBatch.set(batch);
          // Update in batches list
          const updated = this.batches().map((b) => (b.id === batch.id ? batch : b));
          this.batches.set(updated);
          this.loading.set(false);
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to process promotion batch');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  // =========================================================================
  // Report
  // =========================================================================

  getReport(batchId: string): Observable<PromotionReportResponse> {
    this.loading.set(true);
    this.error.set(null);

    return this.http
      .get<PromotionReportResponse>(`${this.apiUrl}/promotion-batches/${batchId}/report`)
      .pipe(
        tap(() => {
          this.loading.set(false);
        }),
        catchError((err) => {
          this.error.set(err.error?.error || 'Failed to get promotion report');
          this.loading.set(false);
          return throwError(() => err);
        })
      );
  }

  // =========================================================================
  // Utilities
  // =========================================================================

  clearSelected(): void {
    this.selectedBatch.set(null);
    this.records.set([]);
    this.summary.set(null);
  }

  clearError(): void {
    this.error.set(null);
  }
}
