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
  TeacherPeriodsResponse,
  PeriodAttendance,
  MarkPeriodAttendanceRequest,
  MarkPeriodAttendanceResult,
  DailySummary,
  SubjectAttendanceStats,
  EditAttendanceRequest,
  EditAttendanceResult,
  AttendanceAuditTrail,
  EditWindowStatus,
  MonthlyCalendar,
  StudentSummary,
  ClassReport,
  MonthlyClassReport,
  LowAttendanceDashboard,
  UnmarkedAttendance,
  DailyReportSummary,
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

  // ============================================================================
  // Period-wise Attendance Methods (Story 7.2)
  // ============================================================================

  /**
   * Get periods available for a section on a specific date
   */
  getTeacherPeriods(sectionId: string, date?: string): Observable<TeacherPeriodsResponse> {
    let params = new HttpParams().set('section_id', sectionId);
    if (date) {
      params = params.set('date', date);
    }

    return this.http
      .get<ApiResponse<TeacherPeriodsResponse>>(`${this.apiUrl}/periods`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get attendance for a specific period
   */
  getPeriodAttendance(periodId: string, sectionId: string, date?: string): Observable<PeriodAttendance> {
    let params = new HttpParams().set('section_id', sectionId);
    if (date) {
      params = params.set('date', date);
    }

    return this.http
      .get<ApiResponse<PeriodAttendance>>(`${this.apiUrl}/period/${periodId}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Mark attendance for a specific period
   */
  markPeriodAttendance(periodId: string, request: MarkPeriodAttendanceRequest): Observable<MarkPeriodAttendanceResult> {
    return this.http
      .post<ApiResponse<MarkPeriodAttendanceResult>>(`${this.apiUrl}/period/${periodId}`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get daily summary of period attendance for a section
   */
  getDailySummary(sectionId: string, date?: string): Observable<DailySummary> {
    let params = new HttpParams().set('section_id', sectionId);
    if (date) {
      params = params.set('date', date);
    }

    return this.http
      .get<ApiResponse<DailySummary>>(`${this.apiUrl}/daily-summary`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get subject-wise attendance statistics for a student
   */
  getSubjectAttendance(
    subjectId: string,
    studentId: string,
    dateFrom?: string,
    dateTo?: string
  ): Observable<SubjectAttendanceStats> {
    let params = new HttpParams().set('student_id', studentId);
    if (dateFrom) {
      params = params.set('date_from', dateFrom);
    }
    if (dateTo) {
      params = params.set('date_to', dateTo);
    }

    return this.http
      .get<ApiResponse<SubjectAttendanceStats>>(`${this.apiUrl}/subject/${subjectId}`, { params })
      .pipe(map(response => response.data!));
  }

  // ============================================================================
  // Attendance Edit & Audit Methods (Story 7.3)
  // ============================================================================

  /**
   * Edit an attendance record with reason
   */
  editAttendance(attendanceId: string, request: EditAttendanceRequest): Observable<EditAttendanceResult> {
    return this.http
      .put<ApiResponse<EditAttendanceResult>>(`${this.apiUrl}/${attendanceId}`, request)
      .pipe(map(response => response.data!));
  }

  /**
   * Get audit trail (edit history) for an attendance record
   */
  getAttendanceHistory(attendanceId: string): Observable<AttendanceAuditTrail> {
    return this.http
      .get<ApiResponse<AttendanceAuditTrail>>(`${this.apiUrl}/${attendanceId}/history`)
      .pipe(map(response => response.data!));
  }

  /**
   * Get edit window status for an attendance record
   */
  getEditWindowStatus(attendanceId: string): Observable<EditWindowStatus> {
    return this.http
      .get<ApiResponse<EditWindowStatus>>(`${this.apiUrl}/${attendanceId}/edit-status`)
      .pipe(map(response => response.data!));
  }

  // ============================================================================
  // Calendar & Reports Methods (Stories 7.4-7.8)
  // ============================================================================

  /**
   * Get student's monthly attendance calendar
   */
  getStudentCalendar(studentId: string, year?: number, month?: number): Observable<MonthlyCalendar> {
    let params = new HttpParams();
    if (year) params = params.set('year', year.toString());
    if (month) params = params.set('month', month.toString());

    return this.http
      .get<ApiResponse<MonthlyCalendar>>(`${this.apiUrl}/calendar/${studentId}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get student's attendance summary with trend
   */
  getStudentSummaryReport(studentId: string, dateFrom?: string, dateTo?: string): Observable<StudentSummary> {
    let params = new HttpParams();
    if (dateFrom) params = params.set('date_from', dateFrom);
    if (dateTo) params = params.set('date_to', dateTo);

    return this.http
      .get<ApiResponse<StudentSummary>>(`${this.apiUrl}/summary/${studentId}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get class attendance report for a date
   */
  getClassReport(sectionId: string, date?: string): Observable<ClassReport> {
    let params = new HttpParams();
    if (date) params = params.set('date', date);

    return this.http
      .get<ApiResponse<ClassReport>>(`${this.apiUrl}/reports/class/${sectionId}`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get monthly class attendance report
   */
  getMonthlyClassReport(sectionId: string, year?: number, month?: number): Observable<MonthlyClassReport> {
    let params = new HttpParams();
    if (year) params = params.set('year', year.toString());
    if (month) params = params.set('month', month.toString());

    return this.http
      .get<ApiResponse<MonthlyClassReport>>(`${this.apiUrl}/reports/class/${sectionId}/monthly`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get daily attendance report
   */
  getDailyReport(date?: string): Observable<DailyReportSummary> {
    let params = new HttpParams();
    if (date) params = params.set('date', date);

    return this.http
      .get<ApiResponse<DailyReportSummary>>(`${this.apiUrl}/reports/daily`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get low attendance dashboard
   */
  getLowAttendanceDashboard(dateFrom?: string, dateTo?: string, threshold?: number): Observable<LowAttendanceDashboard> {
    let params = new HttpParams();
    if (dateFrom) params = params.set('date_from', dateFrom);
    if (dateTo) params = params.set('date_to', dateTo);
    if (threshold) params = params.set('threshold', threshold.toString());

    return this.http
      .get<ApiResponse<LowAttendanceDashboard>>(`${this.apiUrl}/alerts/low-attendance`, { params })
      .pipe(map(response => response.data!));
  }

  /**
   * Get unmarked attendance classes
   */
  getUnmarkedAttendance(date?: string): Observable<UnmarkedAttendance> {
    let params = new HttpParams();
    if (date) params = params.set('date', date);

    return this.http
      .get<ApiResponse<UnmarkedAttendance>>(`${this.apiUrl}/alerts/unmarked`, { params })
      .pipe(map(response => response.data!));
  }
}
