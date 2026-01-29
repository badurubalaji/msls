/**
 * MSLS Subject Service
 *
 * HTTP service for subject management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  Subject,
  SubjectListResponse,
  SubjectFilter,
  CreateSubjectRequest,
  UpdateSubjectRequest,
  ClassSubject,
  ClassSubjectListResponse,
  CreateClassSubjectRequest,
  UpdateClassSubjectRequest,
} from '../academic.model';

/**
 * SubjectService - Handles all subject-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class SubjectService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/subjects';

  /**
   * Get all subjects with optional filters
   */
  getSubjects(filter?: SubjectFilter): Observable<Subject[]> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<SubjectListResponse>(this.basePath, { params }).pipe(
      map(response => response.subjects || [])
    );
  }

  /**
   * Get subjects with total count
   */
  getSubjectsWithTotal(filter?: SubjectFilter): Observable<SubjectListResponse> {
    const params = this.buildFilterParams(filter);
    return this.apiService.get<SubjectListResponse>(this.basePath, { params });
  }

  /**
   * Get a single subject by ID
   */
  getSubject(id: string): Observable<Subject> {
    return this.apiService.get<Subject>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new subject
   */
  createSubject(data: CreateSubjectRequest): Observable<Subject> {
    return this.apiService.post<Subject>(this.basePath, data);
  }

  /**
   * Update an existing subject
   */
  updateSubject(id: string, data: UpdateSubjectRequest): Observable<Subject> {
    return this.apiService.put<Subject>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a subject
   */
  deleteSubject(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Get subjects assigned to a class
   */
  getClassSubjects(classId: string): Observable<ClassSubject[]> {
    return this.apiService.get<ClassSubjectListResponse>(`/classes/${classId}/subjects`).pipe(
      map(response => response.classSubjects || [])
    );
  }

  /**
   * Assign a subject to a class
   */
  createClassSubject(classId: string, data: CreateClassSubjectRequest): Observable<ClassSubject> {
    return this.apiService.post<ClassSubject>(`/classes/${classId}/subjects`, data);
  }

  /**
   * Update a class-subject mapping
   */
  updateClassSubject(classId: string, classSubjectId: string, data: UpdateClassSubjectRequest): Observable<ClassSubject> {
    return this.apiService.put<ClassSubject>(`/classes/${classId}/subjects/${classSubjectId}`, data);
  }

  /**
   * Remove a subject from a class
   */
  deleteClassSubject(classId: string, classSubjectId: string): Observable<void> {
    return this.apiService.delete<void>(`/classes/${classId}/subjects/${classSubjectId}`);
  }

  /**
   * Build query parameters from filter object
   */
  private buildFilterParams(filter?: SubjectFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.subjectType) params['subject_type'] = filter.subjectType;
    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);
    if (filter.search) params['search'] = filter.search;

    return params;
  }
}
