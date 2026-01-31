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
  HallTicket,
  HallTicketListResponse,
  HallTicketFilter,
  HallTicketTemplate,
  GenerateHallTicketsRequest,
  GenerateHallTicketsResponse,
  CreateHallTicketTemplateRequest,
  UpdateHallTicketTemplateRequest,
  VerifyHallTicketResponse,
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
  // Hall Ticket Methods
  // ========================================

  /**
   * Get hall tickets for an examination.
   */
  getHallTickets(examinationId: string, filter?: HallTicketFilter): Observable<HallTicketListResponse> {
    const params = this.buildHallTicketFilterParams(filter);
    return this.apiService.get<HallTicketListResponse>(`/examinations/${examinationId}/hall-tickets`, { params });
  }

  /**
   * Get a single hall ticket by ID.
   */
  getHallTicket(examinationId: string, ticketId: string): Observable<HallTicket> {
    return this.apiService.get<{ data: HallTicket }>(`/examinations/${examinationId}/hall-tickets/${ticketId}`).pipe(
      map(response => response.data)
    );
  }

  /**
   * Generate hall tickets for an examination.
   */
  generateHallTickets(examinationId: string, request?: GenerateHallTicketsRequest): Observable<GenerateHallTicketsResponse> {
    const payload: Record<string, unknown> = {};
    if (request?.classId) payload['classId'] = request.classId;
    if (request?.sectionId) payload['sectionId'] = request.sectionId;
    if (request?.rollNumberPrefix) payload['rollNumberPrefix'] = request.rollNumberPrefix;
    return this.apiService.post<{ data: GenerateHallTicketsResponse }>(`/examinations/${examinationId}/hall-tickets/generate`, payload).pipe(
      map(response => response.data)
    );
  }

  /**
   * Delete a hall ticket.
   */
  deleteHallTicket(examinationId: string, ticketId: string): Observable<void> {
    return this.apiService.delete<void>(`/examinations/${examinationId}/hall-tickets/${ticketId}`);
  }

  /**
   * Download a single hall ticket PDF.
   */
  downloadHallTicketPdf(examinationId: string, ticketId: string): Observable<Blob> {
    return this.apiService.getBlob(`/examinations/${examinationId}/hall-tickets/${ticketId}/pdf`);
  }

  /**
   * Download batch hall tickets PDF.
   */
  downloadBatchHallTicketsPdf(examinationId: string, classId?: string): Observable<Blob> {
    const params: Record<string, string> = {};
    if (classId) params['classId'] = classId;
    return this.apiService.getBlob(`/examinations/${examinationId}/hall-tickets/pdf`, { params });
  }

  /**
   * Verify a hall ticket by QR code.
   */
  verifyHallTicket(qrCode: string): Observable<VerifyHallTicketResponse> {
    return this.apiService.get<{ data: VerifyHallTicketResponse }>(`/hall-tickets/verify/${encodeURIComponent(qrCode)}`).pipe(
      map(response => response.data)
    );
  }

  // ========================================
  // Hall Ticket Template Methods
  // ========================================

  /**
   * Get all hall ticket templates.
   */
  getHallTicketTemplates(): Observable<HallTicketTemplate[]> {
    return this.apiService.get<{ data: HallTicketTemplate[] }>('/hall-ticket-templates').pipe(
      map(response => response.data || [])
    );
  }

  /**
   * Get a single hall ticket template by ID.
   */
  getHallTicketTemplate(id: string): Observable<HallTicketTemplate> {
    return this.apiService.get<{ data: HallTicketTemplate }>(`/hall-ticket-templates/${id}`).pipe(
      map(response => response.data)
    );
  }

  /**
   * Create a new hall ticket template.
   */
  createHallTicketTemplate(data: CreateHallTicketTemplateRequest): Observable<HallTicketTemplate> {
    return this.apiService.post<{ data: HallTicketTemplate }>('/hall-ticket-templates', data).pipe(
      map(response => response.data)
    );
  }

  /**
   * Update a hall ticket template.
   */
  updateHallTicketTemplate(id: string, data: UpdateHallTicketTemplateRequest): Observable<HallTicketTemplate> {
    return this.apiService.put<{ data: HallTicketTemplate }>(`/hall-ticket-templates/${id}`, data).pipe(
      map(response => response.data)
    );
  }

  /**
   * Delete a hall ticket template.
   */
  deleteHallTicketTemplate(id: string): Observable<void> {
    return this.apiService.delete<void>(`/hall-ticket-templates/${id}`);
  }

  // ========================================
  // Private Helper Methods
  // ========================================

  private buildHallTicketFilterParams(filter?: HallTicketFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.classId) params['classId'] = filter.classId;
    if (filter.sectionId) params['sectionId'] = filter.sectionId;
    if (filter.status) params['status'] = filter.status;
    if (filter.search) params['search'] = filter.search;
    if (filter.limit) params['limit'] = String(filter.limit);
    if (filter.offset) params['offset'] = String(filter.offset);

    return params;
  }

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
