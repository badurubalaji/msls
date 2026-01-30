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
  // Private Helper Methods
  // ========================================

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
