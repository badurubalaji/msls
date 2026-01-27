/**
 * Department model interfaces for the frontend
 */

export interface Department {
  id: string;
  name: string;
  code: string;
  description?: string;
  branchId: string;
  branchName?: string;
  headId?: string;
  headName?: string;
  isActive: boolean;
  staffCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface DepartmentListResponse {
  departments: Department[];
  total: number;
}

export interface CreateDepartmentRequest {
  branchId: string;
  name: string;
  code: string;
  description?: string;
  headId?: string;
  isActive?: boolean;
}

export interface UpdateDepartmentRequest {
  name?: string;
  code?: string;
  description?: string;
  headId?: string;
  isActive?: boolean;
}

export interface DepartmentDropdownItem {
  id: string;
  name: string;
  code: string;
}
