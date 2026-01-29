/**
 * Teacher Assignment Service
 * Story 5.7: Teacher Subject Assignment
 *
 * HTTP service for teacher assignment management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import {
  Assignment,
  AssignmentListResponse,
  CreateAssignmentRequest,
  UpdateAssignmentRequest,
  BulkCreateRequest,
  WorkloadReportResponse,
  UnassignedSubjectsResponse,
  ClassTeacherResponse,
  WorkloadSettings,
  UpdateWorkloadSettingsRequest,
} from './assignment.model';

@Injectable({ providedIn: 'root' })
export class AssignmentService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/teacher-assignments';

  // ========================================
  // Assignment CRUD Methods
  // ========================================

  /**
   * Get all assignments with optional filters
   */
  getAssignments(params?: {
    staffId?: string;
    subjectId?: string;
    classId?: string;
    sectionId?: string;
    academicYearId?: string;
    isClassTeacher?: boolean;
    status?: string;
    cursor?: string;
    limit?: number;
  }): Observable<AssignmentListResponse> {
    const queryParams: Record<string, string> = {};
    if (params?.staffId) queryParams['staff_id'] = params.staffId;
    if (params?.subjectId) queryParams['subject_id'] = params.subjectId;
    if (params?.classId) queryParams['class_id'] = params.classId;
    if (params?.sectionId) queryParams['section_id'] = params.sectionId;
    if (params?.academicYearId) queryParams['academic_year_id'] = params.academicYearId;
    if (params?.isClassTeacher !== undefined) queryParams['is_class_teacher'] = String(params.isClassTeacher);
    if (params?.status) queryParams['status'] = params.status;
    if (params?.cursor) queryParams['cursor'] = params.cursor;
    if (params?.limit) queryParams['limit'] = String(params.limit);

    return this.apiService.get<AssignmentListResponse>(this.basePath, { params: queryParams });
  }

  /**
   * Get a single assignment by ID
   */
  getAssignment(id: string): Observable<Assignment> {
    return this.apiService.get<Assignment>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new assignment
   */
  createAssignment(data: CreateAssignmentRequest): Observable<Assignment> {
    return this.apiService.post<Assignment>(this.basePath, data);
  }

  /**
   * Update an assignment
   */
  updateAssignment(id: string, data: UpdateAssignmentRequest): Observable<Assignment> {
    return this.apiService.put<Assignment>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete an assignment
   */
  deleteAssignment(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Bulk create assignments
   */
  bulkCreateAssignments(data: BulkCreateRequest): Observable<{
    created: Assignment[];
    errors: number;
    total: number;
  }> {
    return this.apiService.post<{ created: Assignment[]; errors: number; total: number }>(
      `${this.basePath}/bulk`,
      data
    );
  }

  // ========================================
  // Workload Methods
  // ========================================

  /**
   * Get workload report for all teachers
   */
  getWorkloadReport(academicYearId: string, branchId?: string): Observable<WorkloadReportResponse> {
    const params: Record<string, string> = { academic_year_id: academicYearId };
    if (branchId) params['branch_id'] = branchId;
    return this.apiService.get<WorkloadReportResponse>(`${this.basePath}/workload`, { params });
  }

  /**
   * Get unassigned subjects
   */
  getUnassignedSubjects(academicYearId: string): Observable<UnassignedSubjectsResponse> {
    return this.apiService.get<UnassignedSubjectsResponse>(`${this.basePath}/unassigned`, {
      params: { academic_year_id: academicYearId },
    });
  }

  // ========================================
  // Staff Assignment Methods
  // ========================================

  /**
   * Get all assignments for a specific staff member
   */
  getStaffAssignments(staffId: string, academicYearId?: string): Observable<Assignment[]> {
    const params: Record<string, string> = {};
    if (academicYearId) params['academic_year_id'] = academicYearId;
    return this.apiService.get<Assignment[]>(`/staff/${staffId}/assignments`, { params });
  }

  // ========================================
  // Class Teacher Methods
  // ========================================

  /**
   * Get class teacher for a class-section
   */
  getClassTeacher(classId: string, academicYearId: string, sectionId?: string): Observable<ClassTeacherResponse> {
    const params: Record<string, string> = { academic_year_id: academicYearId };
    if (sectionId) params['section_id'] = sectionId;
    return this.apiService.get<ClassTeacherResponse>(`/classes/${classId}/class-teacher`, { params });
  }

  /**
   * Set class teacher for a class-section
   */
  setClassTeacher(
    classId: string,
    academicYearId: string,
    staffId: string,
    sectionId?: string
  ): Observable<ClassTeacherResponse> {
    const params: Record<string, string> = { academic_year_id: academicYearId };
    if (sectionId) params['section_id'] = sectionId;
    return this.apiService.put<ClassTeacherResponse>(
      `/classes/${classId}/class-teacher`,
      { staffId },
      { params }
    );
  }

  // ========================================
  // Workload Settings Methods
  // ========================================

  /**
   * Get workload settings for a branch
   */
  getWorkloadSettings(branchId: string): Observable<WorkloadSettings> {
    return this.apiService.get<WorkloadSettings>(`/workload-settings/${branchId}`);
  }

  /**
   * Update workload settings for a branch
   */
  updateWorkloadSettings(branchId: string, data: UpdateWorkloadSettingsRequest): Observable<WorkloadSettings> {
    return this.apiService.put<WorkloadSettings>(`/workload-settings/${branchId}`, data);
  }
}
