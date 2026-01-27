import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ApiService } from '../../../core/services/api.service';
import {
  Department,
  DepartmentListResponse,
  CreateDepartmentRequest,
  UpdateDepartmentRequest,
  DepartmentDropdownItem,
} from './department.model';

@Injectable({ providedIn: 'root' })
export class DepartmentService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/departments';

  /**
   * Get all departments with optional filters
   */
  getDepartments(params?: {
    branchId?: string;
    isActive?: boolean;
    search?: string;
  }): Observable<Department[]> {
    const queryParams: Record<string, string> = {};
    if (params?.branchId) queryParams['branch_id'] = params.branchId;
    if (params?.isActive !== undefined)
      queryParams['is_active'] = String(params.isActive);
    if (params?.search) queryParams['search'] = params.search;

    return this.apiService
      .get<DepartmentListResponse>(this.basePath, { params: queryParams })
      .pipe(map(response => response.departments || []));
  }

  /**
   * Get a single department by ID
   */
  getDepartment(id: string): Observable<Department> {
    return this.apiService.get<Department>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new department
   */
  createDepartment(data: CreateDepartmentRequest): Observable<Department> {
    return this.apiService.post<Department>(this.basePath, data);
  }

  /**
   * Update an existing department
   */
  updateDepartment(
    id: string,
    data: UpdateDepartmentRequest
  ): Observable<Department> {
    return this.apiService.put<Department>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a department
   */
  deleteDepartment(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Get departments for dropdown/select
   */
  getDepartmentsDropdown(branchId?: string): Observable<DepartmentDropdownItem[]> {
    const queryParams: Record<string, string> = {};
    if (branchId) queryParams['branch_id'] = branchId;

    return this.apiService.get<DepartmentDropdownItem[]>(
      `${this.basePath}/dropdown`,
      { params: queryParams }
    );
  }
}
