/**
 * Payroll Service
 * Story 5.6: Payroll Processing
 *
 * HTTP service for payroll management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import {
  PayRun,
  PayRunListResponse,
  CreatePayRunRequest,
  Payslip,
  PayslipListResponse,
  AdjustPayslipRequest,
  PayRunSummary,
} from './payroll.model';

@Injectable({ providedIn: 'root' })
export class PayrollService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/payroll';

  // ========================================
  // Pay Run Methods
  // ========================================

  /**
   * Get all pay runs with optional filters
   */
  getPayRuns(params?: {
    year?: number;
    month?: number;
    status?: string;
    branchId?: string;
  }): Observable<PayRun[]> {
    const queryParams: Record<string, string> = {};
    if (params?.year) queryParams['year'] = String(params.year);
    if (params?.month) queryParams['month'] = String(params.month);
    if (params?.status) queryParams['status'] = params.status;
    if (params?.branchId) queryParams['branch_id'] = params.branchId;

    return this.apiService
      .get<PayRunListResponse>(`${this.basePath}/runs`, { params: queryParams })
      .pipe(map(response => response.payRuns || []));
  }

  /**
   * Get a single pay run by ID
   */
  getPayRun(id: string): Observable<PayRun> {
    return this.apiService.get<PayRun>(`${this.basePath}/runs/${id}`);
  }

  /**
   * Create a new pay run
   */
  createPayRun(data: CreatePayRunRequest): Observable<PayRun> {
    return this.apiService.post<PayRun>(`${this.basePath}/runs`, data);
  }

  /**
   * Calculate payroll for a pay run
   */
  calculatePayroll(id: string): Observable<PayRun> {
    return this.apiService.post<PayRun>(`${this.basePath}/runs/${id}/calculate`, {});
  }

  /**
   * Approve a pay run
   */
  approvePayRun(id: string): Observable<PayRun> {
    return this.apiService.post<PayRun>(`${this.basePath}/runs/${id}/approve`, {});
  }

  /**
   * Finalize a pay run
   */
  finalizePayRun(id: string): Observable<PayRun> {
    return this.apiService.post<PayRun>(`${this.basePath}/runs/${id}/finalize`, {});
  }

  /**
   * Delete a draft pay run
   */
  deletePayRun(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/runs/${id}`);
  }

  /**
   * Get pay run summary (department-wise)
   */
  getPayRunSummary(id: string): Observable<PayRunSummary> {
    return this.apiService.get<PayRunSummary>(`${this.basePath}/runs/${id}/summary`);
  }

  /**
   * Export bank transfer file
   */
  exportBankFile(id: string): Observable<Blob> {
    return this.apiService.getBlob(`${this.basePath}/runs/${id}/export`);
  }

  // ========================================
  // Payslip Methods
  // ========================================

  /**
   * Get all payslips for a pay run
   */
  getPayslips(payRunId: string, params?: {
    search?: string;
    departmentId?: string;
  }): Observable<Payslip[]> {
    const queryParams: Record<string, string> = {};
    if (params?.search) queryParams['search'] = params.search;
    if (params?.departmentId) queryParams['department_id'] = params.departmentId;

    return this.apiService
      .get<PayslipListResponse>(`${this.basePath}/runs/${payRunId}/payslips`, { params: queryParams })
      .pipe(map(response => response.payslips || []));
  }

  /**
   * Get a single payslip by ID
   */
  getPayslip(id: string): Observable<Payslip> {
    return this.apiService.get<Payslip>(`${this.basePath}/payslips/${id}`);
  }

  /**
   * Adjust a payslip (before approval)
   */
  adjustPayslip(id: string, data: AdjustPayslipRequest): Observable<Payslip> {
    return this.apiService.put<Payslip>(`${this.basePath}/payslips/${id}`, data);
  }

  /**
   * Download payslip PDF
   */
  downloadPayslipPdf(id: string): Observable<Blob> {
    return this.apiService.getBlob(`${this.basePath}/payslips/${id}/pdf`);
  }

  // ========================================
  // Staff Payslip History
  // ========================================

  /**
   * Get payslip history for a staff member
   */
  getStaffPayslips(staffId: string): Observable<Payslip[]> {
    return this.apiService
      .get<PayslipListResponse>(`/staff/${staffId}/payslips`)
      .pipe(map(response => response.payslips || []));
  }
}
