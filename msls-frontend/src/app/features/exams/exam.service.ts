/**
 * MSLS Exam Service
 *
 * HTTP service for exam management API calls including exam types.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../core/services/api.service';
import {
  ExamType,
  ExamTypeListResponse,
  ExamTypeFilter,
  CreateExamTypeRequest,
  UpdateExamTypeRequest,
  UpdateDisplayOrderRequest,
  ToggleActiveRequest,
  Examination,
  ExaminationFilter,
  CreateExaminationRequest,
  UpdateExaminationRequest,
  ExamSchedule,
  CreateScheduleRequest,
  UpdateScheduleRequest,
} from './exam.model';

/**
 * ExamService - Handles all exam-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class ExamService {
  private readonly apiService = inject(ApiService);

  // ========================================
  // Exam Type Methods
  // ========================================

  /**
   * Get all exam types with optional filters.
   */
  getExamTypes(filter?: ExamTypeFilter): Observable<ExamType[]> {
    const params = this.buildExamTypeFilterParams(filter);
    return this.apiService.get<ExamTypeListResponse>('/exam-types', { params }).pipe(
      map(response => response.items || [])
    );
  }

  /**
   * Get exam types with total count.
   */
  getExamTypesWithTotal(filter?: ExamTypeFilter): Observable<ExamTypeListResponse> {
    const params = this.buildExamTypeFilterParams(filter);
    return this.apiService.get<ExamTypeListResponse>('/exam-types', { params });
  }

  /**
   * Get a single exam type by ID.
   */
  getExamType(id: string): Observable<ExamType> {
    return this.apiService.get<ExamType>(`/exam-types/${id}`);
  }

  /**
   * Create a new exam type.
   */
  createExamType(data: CreateExamTypeRequest): Observable<ExamType> {
    return this.apiService.post<ExamType>('/exam-types', this.transformToSnakeCase(data));
  }

  /**
   * Update an existing exam type.
   */
  updateExamType(id: string, data: UpdateExamTypeRequest): Observable<ExamType> {
    return this.apiService.put<ExamType>(`/exam-types/${id}`, this.transformToSnakeCase(data));
  }

  /**
   * Delete an exam type.
   */
  deleteExamType(id: string): Observable<void> {
    return this.apiService.delete<void>(`/exam-types/${id}`);
  }

  /**
   * Toggle the active status of an exam type.
   */
  toggleExamTypeActive(id: string, isActive: boolean): Observable<void> {
    const request: ToggleActiveRequest = { isActive };
    return this.apiService.patch<void>(`/exam-types/${id}/active`, { is_active: isActive });
  }

  /**
   * Update display order for multiple exam types.
   */
  updateDisplayOrder(request: UpdateDisplayOrderRequest): Observable<void> {
    const transformedItems = request.items.map(item => ({
      id: item.id,
      display_order: item.displayOrder,
    }));
    return this.apiService.put<void>('/exam-types/order', { items: transformedItems });
  }

  // ========================================
  // Examination Methods
  // ========================================

  /**
   * Get all examinations with optional filters.
   */
  getExaminations(filter?: ExaminationFilter): Observable<Examination[]> {
    const params = this.buildExaminationFilterParams(filter);
    return this.apiService.get<Examination[]>('/examinations', { params });
  }

  /**
   * Get a single examination by ID.
   */
  getExamination(id: string): Observable<Examination> {
    return this.apiService.get<Examination>(`/examinations/${id}`);
  }

  /**
   * Create a new examination.
   */
  createExamination(data: CreateExaminationRequest): Observable<Examination> {
    return this.apiService.post<Examination>('/examinations', {
      name: data.name,
      examTypeId: data.examTypeId,
      academicYearId: data.academicYearId,
      startDate: data.startDate,
      endDate: data.endDate,
      description: data.description,
      classIds: data.classIds,
    });
  }

  /**
   * Update an existing examination.
   */
  updateExamination(id: string, data: UpdateExaminationRequest): Observable<Examination> {
    const payload: Record<string, unknown> = {};
    if (data.name !== undefined) payload['name'] = data.name;
    if (data.examTypeId !== undefined) payload['examTypeId'] = data.examTypeId;
    if (data.academicYearId !== undefined) payload['academicYearId'] = data.academicYearId;
    if (data.startDate !== undefined) payload['startDate'] = data.startDate;
    if (data.endDate !== undefined) payload['endDate'] = data.endDate;
    if (data.description !== undefined) payload['description'] = data.description;
    if (data.classIds !== undefined) payload['classIds'] = data.classIds;
    return this.apiService.put<Examination>(`/examinations/${id}`, payload);
  }

  /**
   * Delete an examination.
   */
  deleteExamination(id: string): Observable<void> {
    return this.apiService.delete<void>(`/examinations/${id}`);
  }

  /**
   * Publish an examination (change status to scheduled).
   */
  publishExamination(id: string): Observable<Examination> {
    return this.apiService.post<Examination>(`/examinations/${id}/publish`, {});
  }

  /**
   * Unpublish an examination (revert to draft).
   */
  unpublishExamination(id: string): Observable<Examination> {
    return this.apiService.post<Examination>(`/examinations/${id}/unpublish`, {});
  }

  // ========================================
  // Exam Schedule Methods
  // ========================================

  /**
   * Get all schedules for an examination.
   */
  getSchedules(examinationId: string): Observable<ExamSchedule[]> {
    return this.apiService.get<ExamSchedule[]>(`/examinations/${examinationId}/schedules`);
  }

  /**
   * Create a new schedule for an examination.
   */
  createSchedule(examinationId: string, data: CreateScheduleRequest): Observable<ExamSchedule> {
    return this.apiService.post<ExamSchedule>(`/examinations/${examinationId}/schedules`, {
      subjectId: data.subjectId,
      examDate: data.examDate,
      startTime: data.startTime,
      endTime: data.endTime,
      maxMarks: data.maxMarks,
      passingMarks: data.passingMarks,
      venue: data.venue,
      notes: data.notes,
    });
  }

  /**
   * Update a schedule.
   */
  updateSchedule(examinationId: string, scheduleId: string, data: UpdateScheduleRequest): Observable<ExamSchedule> {
    const payload: Record<string, unknown> = {};
    if (data.subjectId !== undefined) payload['subjectId'] = data.subjectId;
    if (data.examDate !== undefined) payload['examDate'] = data.examDate;
    if (data.startTime !== undefined) payload['startTime'] = data.startTime;
    if (data.endTime !== undefined) payload['endTime'] = data.endTime;
    if (data.maxMarks !== undefined) payload['maxMarks'] = data.maxMarks;
    if (data.passingMarks !== undefined) payload['passingMarks'] = data.passingMarks;
    if (data.venue !== undefined) payload['venue'] = data.venue;
    if (data.notes !== undefined) payload['notes'] = data.notes;
    return this.apiService.put<ExamSchedule>(`/examinations/${examinationId}/schedules/${scheduleId}`, payload);
  }

  /**
   * Delete a schedule.
   */
  deleteSchedule(examinationId: string, scheduleId: string): Observable<void> {
    return this.apiService.delete<void>(`/examinations/${examinationId}/schedules/${scheduleId}`);
  }

  // ========================================
  // Private Helper Methods
  // ========================================

  private buildExaminationFilterParams(filter?: ExaminationFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.academicYearId) params['academicYearId'] = filter.academicYearId;
    if (filter.examTypeId) params['examTypeId'] = filter.examTypeId;
    if (filter.classId) params['classId'] = filter.classId;
    if (filter.status) params['status'] = filter.status;
    if (filter.search) params['search'] = filter.search;

    return params;
  }

  private buildExamTypeFilterParams(filter?: ExamTypeFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);
    if (filter.search) params['search'] = filter.search;

    return params;
  }

  private transformToSnakeCase(data: CreateExamTypeRequest | UpdateExamTypeRequest): Record<string, unknown> {
    const result: Record<string, unknown> = {};

    if ('name' in data && data.name !== undefined) result['name'] = data.name;
    if ('code' in data && data.code !== undefined) result['code'] = data.code;
    if ('description' in data) result['description'] = data.description;
    if ('weightage' in data && data.weightage !== undefined) result['weightage'] = data.weightage;
    if ('evaluationType' in data && data.evaluationType !== undefined) result['evaluation_type'] = data.evaluationType;
    if ('defaultMaxMarks' in data && data.defaultMaxMarks !== undefined) result['default_max_marks'] = data.defaultMaxMarks;
    if ('defaultPassingMarks' in data) result['default_passing_marks'] = data.defaultPassingMarks;

    return result;
  }
}
