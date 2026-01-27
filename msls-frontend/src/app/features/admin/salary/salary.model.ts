/**
 * Salary model interfaces for the frontend
 */

// Component Types
export type ComponentType = 'earning' | 'deduction';
export type CalculationType = 'fixed' | 'percentage';

// Salary Component
export interface SalaryComponent {
  id: string;
  name: string;
  code: string;
  description?: string;
  componentType: ComponentType;
  calculationType: CalculationType;
  percentageOfId?: string;
  percentageOfName?: string;
  isTaxable: boolean;
  isActive: boolean;
  displayOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface ComponentListResponse {
  components: SalaryComponent[];
  total: number;
}

export interface CreateComponentRequest {
  name: string;
  code: string;
  description?: string;
  componentType: ComponentType;
  calculationType: CalculationType;
  percentageOfId?: string;
  isTaxable: boolean;
  displayOrder: number;
}

export interface UpdateComponentRequest {
  name?: string;
  code?: string;
  description?: string;
  componentType?: ComponentType;
  calculationType?: CalculationType;
  percentageOfId?: string;
  isTaxable?: boolean;
  isActive?: boolean;
  displayOrder?: number;
}

export interface ComponentDropdownItem {
  id: string;
  name: string;
  code: string;
  componentType: ComponentType;
}

// Salary Structure
export interface SalaryStructure {
  id: string;
  name: string;
  code: string;
  description?: string;
  designationId?: string;
  designationName?: string;
  isActive: boolean;
  componentCount: number;
  staffCount: number;
  components?: StructureComponent[];
  createdAt: string;
  updatedAt: string;
}

export interface StructureComponent {
  id: string;
  componentId: string;
  componentName: string;
  componentCode: string;
  componentType: ComponentType;
  amount?: string;
  percentage?: string;
}

export interface StructureListResponse {
  structures: SalaryStructure[];
  total: number;
}

export interface CreateStructureRequest {
  name: string;
  code: string;
  description?: string;
  designationId?: string;
  components: StructureComponentInput[];
}

export interface StructureComponentInput {
  componentId: string;
  amount?: string;
  percentage?: string;
}

export interface UpdateStructureRequest {
  name?: string;
  code?: string;
  description?: string;
  designationId?: string;
  isActive?: boolean;
  components?: StructureComponentInput[];
}

export interface StructureDropdownItem {
  id: string;
  name: string;
  code: string;
}

// Staff Salary
export interface StaffSalary {
  id: string;
  staffId: string;
  staffName?: string;
  structureId?: string;
  structureName?: string;
  effectiveFrom: string;
  effectiveTo?: string;
  grossSalary: string;
  netSalary: string;
  ctc?: string;
  revisionReason?: string;
  isCurrent: boolean;
  components?: StaffSalaryComponent[];
  createdAt: string;
  updatedAt: string;
}

export interface StaffSalaryComponent {
  id: string;
  componentId: string;
  componentName: string;
  componentCode: string;
  componentType: ComponentType;
  amount: string;
  isOverridden: boolean;
}

export interface AssignSalaryRequest {
  staffId: string;
  structureId?: string;
  effectiveFrom: string;
  components: StaffSalaryComponentInput[];
  revisionReason?: string;
}

export interface StaffSalaryComponentInput {
  componentId: string;
  amount: string;
  isOverridden: boolean;
}

export interface StaffSalaryHistoryResponse {
  history: StaffSalary[];
  total: number;
}
