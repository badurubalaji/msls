import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ApiService } from '../../../core/services/api.service';
import {
  Designation,
  DesignationListResponse,
  CreateDesignationRequest,
  UpdateDesignationRequest,
  DesignationDropdownItem,
} from './designation.model';

@Injectable({ providedIn: 'root' })
export class DesignationService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/designations';

  /**
   * Get all designations with optional filters
   */
  getDesignations(params?: {
    departmentId?: string;
    isActive?: boolean;
    search?: string;
  }): Observable<Designation[]> {
    const queryParams: Record<string, string> = {};
    if (params?.departmentId) queryParams['department_id'] = params.departmentId;
    if (params?.isActive !== undefined)
      queryParams['is_active'] = String(params.isActive);
    if (params?.search) queryParams['search'] = params.search;

    return this.apiService
      .get<DesignationListResponse>(this.basePath, { params: queryParams })
      .pipe(map(response => response.designations || []));
  }

  /**
   * Get a single designation by ID
   */
  getDesignation(id: string): Observable<Designation> {
    return this.apiService.get<Designation>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new designation
   */
  createDesignation(data: CreateDesignationRequest): Observable<Designation> {
    return this.apiService.post<Designation>(this.basePath, data);
  }

  /**
   * Update an existing designation
   */
  updateDesignation(
    id: string,
    data: UpdateDesignationRequest
  ): Observable<Designation> {
    return this.apiService.put<Designation>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a designation
   */
  deleteDesignation(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Get designations for dropdown/select
   */
  getDesignationsDropdown(departmentId?: string): Observable<DesignationDropdownItem[]> {
    const queryParams: Record<string, string> = {};
    if (departmentId) queryParams['department_id'] = departmentId;

    return this.apiService.get<DesignationDropdownItem[]>(
      `${this.basePath}/dropdown`,
      { params: queryParams }
    );
  }
}
