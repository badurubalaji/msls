import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ApiService } from '../../../core/services/api.service';
import {
  SalaryComponent,
  ComponentListResponse,
  CreateComponentRequest,
  UpdateComponentRequest,
  ComponentDropdownItem,
  SalaryStructure,
  StructureListResponse,
  CreateStructureRequest,
  UpdateStructureRequest,
  StructureDropdownItem,
  StaffSalary,
  StaffSalaryHistoryResponse,
  AssignSalaryRequest,
} from './salary.model';

@Injectable({ providedIn: 'root' })
export class SalaryService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/salary';

  // ========================================
  // Salary Component Methods
  // ========================================

  /**
   * Get all salary components with optional filters
   */
  getComponents(params?: {
    componentType?: string;
    isActive?: boolean;
    search?: string;
  }): Observable<SalaryComponent[]> {
    const queryParams: Record<string, string> = {};
    if (params?.componentType) queryParams['component_type'] = params.componentType;
    if (params?.isActive !== undefined) queryParams['is_active'] = String(params.isActive);
    if (params?.search) queryParams['search'] = params.search;

    return this.apiService
      .get<ComponentListResponse>(`${this.basePath}/components`, { params: queryParams })
      .pipe(map(response => response.components || []));
  }

  /**
   * Get a single salary component by ID
   */
  getComponent(id: string): Observable<SalaryComponent> {
    return this.apiService.get<SalaryComponent>(`${this.basePath}/components/${id}`);
  }

  /**
   * Create a new salary component
   */
  createComponent(data: CreateComponentRequest): Observable<SalaryComponent> {
    return this.apiService.post<SalaryComponent>(`${this.basePath}/components`, data);
  }

  /**
   * Update an existing salary component
   */
  updateComponent(id: string, data: UpdateComponentRequest): Observable<SalaryComponent> {
    return this.apiService.put<SalaryComponent>(`${this.basePath}/components/${id}`, data);
  }

  /**
   * Delete a salary component
   */
  deleteComponent(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/components/${id}`);
  }

  /**
   * Get active components for dropdown
   */
  getComponentsDropdown(componentType?: string): Observable<ComponentDropdownItem[]> {
    const queryParams: Record<string, string> = {};
    if (componentType) queryParams['component_type'] = componentType;

    return this.apiService.get<ComponentDropdownItem[]>(
      `${this.basePath}/components/dropdown`,
      { params: queryParams }
    );
  }

  // ========================================
  // Salary Structure Methods
  // ========================================

  /**
   * Get all salary structures with optional filters
   */
  getStructures(params?: {
    designationId?: string;
    isActive?: boolean;
    search?: string;
  }): Observable<SalaryStructure[]> {
    const queryParams: Record<string, string> = {};
    if (params?.designationId) queryParams['designation_id'] = params.designationId;
    if (params?.isActive !== undefined) queryParams['is_active'] = String(params.isActive);
    if (params?.search) queryParams['search'] = params.search;

    return this.apiService
      .get<StructureListResponse>(`${this.basePath}/structures`, { params: queryParams })
      .pipe(map(response => response.structures || []));
  }

  /**
   * Get a single salary structure by ID
   */
  getStructure(id: string): Observable<SalaryStructure> {
    return this.apiService.get<SalaryStructure>(`${this.basePath}/structures/${id}`);
  }

  /**
   * Create a new salary structure
   */
  createStructure(data: CreateStructureRequest): Observable<SalaryStructure> {
    return this.apiService.post<SalaryStructure>(`${this.basePath}/structures`, data);
  }

  /**
   * Update an existing salary structure
   */
  updateStructure(id: string, data: UpdateStructureRequest): Observable<SalaryStructure> {
    return this.apiService.put<SalaryStructure>(`${this.basePath}/structures/${id}`, data);
  }

  /**
   * Delete a salary structure
   */
  deleteStructure(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/structures/${id}`);
  }

  /**
   * Get active structures for dropdown
   */
  getStructuresDropdown(): Observable<StructureDropdownItem[]> {
    return this.apiService.get<StructureDropdownItem[]>(`${this.basePath}/structures/dropdown`);
  }

  // ========================================
  // Staff Salary Methods
  // ========================================

  /**
   * Get current salary for a staff member
   */
  getStaffSalary(staffId: string): Observable<StaffSalary> {
    return this.apiService.get<StaffSalary>(`/staff/${staffId}/salary`);
  }

  /**
   * Get salary history for a staff member
   */
  getStaffSalaryHistory(staffId: string): Observable<StaffSalary[]> {
    return this.apiService
      .get<StaffSalaryHistoryResponse>(`/staff/${staffId}/salary/history`)
      .pipe(map(response => response.history || []));
  }

  /**
   * Assign or revise salary for a staff member
   */
  assignSalary(data: AssignSalaryRequest): Observable<StaffSalary> {
    return this.apiService.post<StaffSalary>(`/staff/${data.staffId}/salary`, data);
  }
}
