/**
 * Student Attendance Service
 * Handles all student attendance-related API calls
 */

import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, map } from 'rxjs';

import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../core/models/api-response.model';
import {
  TeacherClass,
  ClassAttendance,
  StudentAttendanceListResponse,
  MarkAttendanceResult,
  StudentAttendanceSettings,
  MarkClassAttendanceRequest,
  UpdateSettingsRequest,
  StudentAttendanceFilter,
} from './student-attendance.model';

@Injectable({ providedIn: 'root' })
export class StudentAttendanceService {
  private http = inject(HttpClient);
  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}/student-attendance`;

  /**
   * Get teacher's assigned classes for attendance marking
   */
  getTeacherClasses(date?: string): Observable<TeacherClass[]> {
    let params = new HttpParams();
    if (date) {
      params = params.set('date', date);
    }

    return this.http
      .get<ApiResponse<TeacherClass[]>>(`${this.apiUrl}/my-classes`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get class attendance for a section on a specific date
   */
  getClassAttendance(sectionId: string, date?: string): Observable<ClassAttendance> {
    let params = new HttpParams();
    if (date) {
      params = params.set('date', date);
    }

    return this.http
      .get<ApiResponse<ClassAttendance>>(`${this.apiUrl}/class/${sectionId}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Mark attendance for a class section
   */
  markClassAttendance(sectionId: string, request: MarkClassAttendanceRequest): Observable<MarkAttendanceResult> {
    return this.http
      .post<ApiResponse<MarkAttendanceResult>>(`${this.apiUrl}/class/${sectionId}`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get all student attendance records (admin view)
   */
  getAllAttendance(filter?: StudentAttendanceFilter): Observable<StudentAttendanceListResponse> {
    let params = new HttpParams();

    if (filter?.sectionId) {
      params = params.set('section_id', filter.sectionId);
    }
    if (filter?.studentId) {
      params = params.set('student_id', filter.studentId);
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
      .get<ApiResponse<StudentAttendanceListResponse>>(`${this.apiUrl}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get student attendance settings for a branch
   */
  getSettings(branchId: string): Observable<StudentAttendanceSettings> {
    const params = new HttpParams().set('branch_id', branchId);
    return this.http
      .get<ApiResponse<StudentAttendanceSettings>>(`${this.apiUrl}/settings`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Update student attendance settings for a branch
   */
  updateSettings(request: UpdateSettingsRequest): Observable<StudentAttendanceSettings> {
    return this.http
      .put<ApiResponse<StudentAttendanceSettings>>(`${this.apiUrl}/settings`, request)
      .pipe(map(response => response.data!));
  }
}
