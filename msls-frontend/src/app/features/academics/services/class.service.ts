/**
 * MSLS Class Service
 *
 * HTTP service for class management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  Class,
  ClassListResponse,
  ClassFilter,
  CreateClassRequest,
  UpdateClassRequest,
  ClassStructureResponse,
  Section,
  SectionListResponse,
} from '../academic.model';

/**
 * ClassService - Handles all class-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class ClassService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/classes';

  /**
   * Get all classes with optional filters
   */
  getClasses(filter?: ClassFilter): Observable<Class[]> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<ClassListResponse>(this.basePath, { params }).pipe(
      map(response => response.classes || [])
    );
  }

  /**
   * Get classes with total count
   */
  getClassesWithTotal(filter?: ClassFilter): Observable<ClassListResponse> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<ClassListResponse>(this.basePath, { params });
  }

  /**
   * Get a single class by ID
   */
  getClass(id: string): Observable<Class> {
    return this.apiService.get<Class>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new class
   */
  createClass(data: CreateClassRequest): Observable<Class> {
    return this.apiService.post<Class>(this.basePath, data);
  }

  /**
   * Update an existing class
   */
  updateClass(id: string, data: UpdateClassRequest): Observable<Class> {
    return this.apiService.put<Class>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a class
   */
  deleteClass(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Get sections for a specific class
   */
  getClassSections(classId: string): Observable<Section[]> {
    return this.apiService.get<SectionListResponse>(`${this.basePath}/${classId}/sections`).pipe(
      map(response => response.sections || [])
    );
  }

  /**
   * Get hierarchical class-section structure
   */
  getClassStructure(branchId?: string, academicYearId?: string): Observable<ClassStructureResponse> {
    const params: Record<string, string> = {};
    if (branchId) params['branch_id'] = branchId;
    if (academicYearId) params['academic_year_id'] = academicYearId;
    return this.apiService.get<ClassStructureResponse>(`${this.basePath}/structure`, { params });
  }

  /**
   * Build query parameters from filter object
   */
  private buildFilterParams(filter?: ClassFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.branchId) params['branch_id'] = filter.branchId;
    if (filter.level) params['level'] = filter.level;
    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);
    if (filter.hasStreams !== undefined) params['has_streams'] = String(filter.hasStreams);
    if (filter.search) params['search'] = filter.search;

    return params;
  }
}
