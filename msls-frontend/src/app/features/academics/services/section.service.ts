/**
 * MSLS Section Service
 *
 * HTTP service for section management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  Section,
  SectionListResponse,
  SectionFilter,
  CreateSectionRequest,
  UpdateSectionRequest,
} from '../academic.model';

/**
 * SectionService - Handles all section-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class SectionService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/sections';

  /**
   * Get all sections with optional filters
   */
  getSections(filter?: SectionFilter): Observable<Section[]> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<SectionListResponse>(this.basePath, { params }).pipe(
      map(response => response.sections || [])
    );
  }

  /**
   * Get sections with total count
   */
  getSectionsWithTotal(filter?: SectionFilter): Observable<SectionListResponse> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<SectionListResponse>(this.basePath, { params });
  }

  /**
   * Get a single section by ID
   */
  getSection(id: string): Observable<Section> {
    return this.apiService.get<Section>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new section
   */
  createSection(data: CreateSectionRequest): Observable<Section> {
    return this.apiService.post<Section>(this.basePath, data);
  }

  /**
   * Update an existing section
   */
  updateSection(id: string, data: UpdateSectionRequest): Observable<Section> {
    return this.apiService.put<Section>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a section
   */
  deleteSection(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Assign class teacher to a section
   */
  assignClassTeacher(sectionId: string, classTeacherId: string | null): Observable<Section> {
    return this.apiService.put<Section>(`${this.basePath}/${sectionId}/class-teacher`, {
      classTeacherId,
    });
  }

  /**
   * Build query parameters from filter object
   */
  private buildFilterParams(filter?: SectionFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.classId) params['class_id'] = filter.classId;
    if (filter.academicYearId) params['academic_year_id'] = filter.academicYearId;
    if (filter.streamId) params['stream_id'] = filter.streamId;
    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);
    if (filter.search) params['search'] = filter.search;

    return params;
  }
}
