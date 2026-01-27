/**
 * Attendance Service
 * Handles all attendance-related API calls
 */

import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, map } from 'rxjs';

import { environment } from '../../../../environments/environment';
import { ApiResponse } from '../../../core/models/api-response.model';
import {
  Attendance,
  TodayAttendance,
  AttendanceSummary,
  AttendanceListResponse,
  Regularization,
  RegularizationListResponse,
  AttendanceSettings,
  CheckInRequest,
  CheckOutRequest,
  MarkAttendanceRequest,
  RegularizationRequest,
  RegularizationReviewRequest,
  UpdateSettingsRequest,
  AttendanceFilter,
  RegularizationFilter,
} from '../models/attendance.model';

@Injectable({ providedIn: 'root' })
export class AttendanceService {
  private http = inject(HttpClient);
  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}/attendance`;

  /**
   * Check in for a staff member
   */
  checkIn(request: CheckInRequest): Observable<Attendance> {
    return this.http
      .post<ApiResponse<Attendance>>(`${this.apiUrl}/check-in`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Check out for a staff member
   */
  checkOut(request: CheckOutRequest): Observable<Attendance> {
    return this.http
      .post<ApiResponse<Attendance>>(`${this.apiUrl}/check-out`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get today's attendance status for a staff member
   */
  getTodayAttendance(staffId: string): Observable<TodayAttendance> {
    const params = new HttpParams().set('staff_id', staffId);
    return this.http
      .get<ApiResponse<TodayAttendance>>(`${this.apiUrl}/my/today`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get attendance records for a staff member
   */
  getMyAttendance(staffId: string, filter?: Partial<AttendanceFilter>): Observable<AttendanceListResponse> {
    let params = new HttpParams().set('staff_id', staffId);

    if (filter?.dateFrom) {
      params = params.set('date_from', filter.dateFrom);
    }
    if (filter?.dateTo) {
      params = params.set('date_to', filter.dateTo);
    }
    if (filter?.cursor) {
      params = params.set('cursor', filter.cursor);
    }
    if (filter?.limit) {
      params = params.set('limit', filter.limit.toString());
    }

    return this.http
      .get<ApiResponse<AttendanceListResponse>>(`${this.apiUrl}/my`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get monthly attendance summary for a staff member
   */
  getMonthlySummary(staffId: string, year?: number, month?: number): Observable<AttendanceSummary> {
    let params = new HttpParams().set('staff_id', staffId);

    if (year) {
      params = params.set('year', year.toString());
    }
    if (month) {
      params = params.set('month', month.toString());
    }

    return this.http
      .get<ApiResponse<AttendanceSummary>>(`${this.apiUrl}/my/summary`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get all attendance records (HR view)
   */
  getAllAttendance(filter?: AttendanceFilter): Observable<AttendanceListResponse> {
    let params = new HttpParams();

    if (filter?.staffId) {
      params = params.set('staff_id', filter.staffId);
    }
    if (filter?.branchId) {
      params = params.set('branch_id', filter.branchId);
    }
    if (filter?.departmentId) {
      params = params.set('department_id', filter.departmentId);
    }
    if (filter?.status) {
      params = params.set('status', filter.status);
    }
    if (filter?.dateFrom) {
      params = params.set('date_from', filter.dateFrom);
    }
    if (filter?.dateTo) {
      params = params.set('date_to', filter.dateTo);
    }
    if (filter?.cursor) {
      params = params.set('cursor', filter.cursor);
    }
    if (filter?.limit) {
      params = params.set('limit', filter.limit.toString());
    }

    return this.http
      .get<ApiResponse<AttendanceListResponse>>(`${this.apiUrl}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Mark attendance for a staff member (HR action)
   */
  markAttendance(request: MarkAttendanceRequest): Observable<Attendance> {
    return this.http
      .post<ApiResponse<Attendance>>(`${this.apiUrl}/mark`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Submit a regularization request
   */
  submitRegularization(request: RegularizationRequest): Observable<Regularization> {
    return this.http
      .post<ApiResponse<Regularization>>(`${this.apiUrl}/regularization`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get regularization requests
   */
  getRegularizations(filter?: RegularizationFilter): Observable<RegularizationListResponse> {
    let params = new HttpParams();

    if (filter?.staffId) {
      params = params.set('staff_id', filter.staffId);
    }
    if (filter?.status) {
      params = params.set('status', filter.status);
    }
    if (filter?.dateFrom) {
      params = params.set('date_from', filter.dateFrom);
    }
    if (filter?.dateTo) {
      params = params.set('date_to', filter.dateTo);
    }
    if (filter?.cursor) {
      params = params.set('cursor', filter.cursor);
    }
    if (filter?.limit) {
      params = params.set('limit', filter.limit.toString());
    }

    return this.http
      .get<ApiResponse<RegularizationListResponse>>(`${this.apiUrl}/regularization`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Approve a regularization request
   */
  approveRegularization(id: string): Observable<Regularization> {
    return this.http
      .put<ApiResponse<Regularization>>(`${this.apiUrl}/regularization/${id}/approve`, {})
      .pipe(map(response => response.data!));
  }

  /**
   * Reject a regularization request
   */
  rejectRegularization(id: string, request: RegularizationReviewRequest): Observable<Regularization> {
    return this.http
      .put<ApiResponse<Regularization>>(`${this.apiUrl}/regularization/${id}/reject`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get attendance settings for a branch
   */
  getSettings(branchId: string): Observable<AttendanceSettings> {
    const params = new HttpParams().set('branch_id', branchId);
    return this.http
      .get<ApiResponse<AttendanceSettings>>(`${this.apiUrl}/settings`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Update attendance settings for a branch
   */
  updateSettings(request: UpdateSettingsRequest): Observable<AttendanceSettings> {
    return this.http
      .put<ApiResponse<AttendanceSettings>>(`${this.apiUrl}/settings`, request)
      .pipe(map(response => response.data!));
  }
}
