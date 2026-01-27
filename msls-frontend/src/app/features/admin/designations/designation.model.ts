/**
 * Designation model interfaces for the frontend
 */

export interface Designation {
  id: string;
  name: string;
  level: number;
  departmentId?: string;
  departmentName?: string;
  isActive: boolean;
  staffCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface DesignationListResponse {
  designations: Designation[];
  total: number;
}

export interface CreateDesignationRequest {
  name: string;
  level: number;
  departmentId?: string;
  isActive?: boolean;
}

export interface UpdateDesignationRequest {
  name?: string;
  level?: number;
  departmentId?: string;
  isActive?: boolean;
}

export interface DesignationDropdownItem {
  id: string;
  name: string;
  level: number;
}

/**
 * Designation levels represent hierarchy where 1 is highest (e.g., CEO)
 * and 10 is lowest (e.g., entry level)
 */
export const DESIGNATION_LEVELS = [
  { value: 1, label: 'Level 1 - Executive' },
  { value: 2, label: 'Level 2 - Senior Management' },
  { value: 3, label: 'Level 3 - Management' },
  { value: 4, label: 'Level 4 - Senior Professional' },
  { value: 5, label: 'Level 5 - Professional' },
  { value: 6, label: 'Level 6 - Senior Associate' },
  { value: 7, label: 'Level 7 - Associate' },
  { value: 8, label: 'Level 8 - Junior' },
  { value: 9, label: 'Level 9 - Trainee' },
  { value: 10, label: 'Level 10 - Entry' },
];
