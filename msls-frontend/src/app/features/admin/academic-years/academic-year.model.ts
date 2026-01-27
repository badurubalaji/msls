/**
 * MSLS Academic Year Models
 *
 * TypeScript interfaces for academic year management including terms and holidays.
 */

/**
 * Academic Term entity
 */
export interface AcademicTerm {
  id: string;
  academicYearId: string;
  name: string;
  startDate: string;
  endDate: string;
  sequence: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * Holiday entity
 */
export interface Holiday {
  id: string;
  academicYearId: string;
  branchId?: string;
  name: string;
  date: string;
  type: HolidayType;
  isOptional: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * Holiday types
 */
export type HolidayType = 'public' | 'religious' | 'school' | 'exam' | 'other';

/**
 * Academic Year entity returned from API
 */
export interface AcademicYear {
  id: string;
  tenantId: string;
  branchId?: string;
  name: string;
  startDate: string;
  endDate: string;
  isCurrent: boolean;
  isActive: boolean;
  terms?: AcademicTerm[];
  holidays?: Holiday[];
  createdAt: string;
  updatedAt: string;
}

/**
 * Request payload for creating an academic year
 */
export interface CreateAcademicYearRequest {
  name: string;
  startDate: string;
  endDate: string;
  isCurrent?: boolean;
  branchId?: string;
}

/**
 * Request payload for updating an academic year
 */
export type UpdateAcademicYearRequest = Partial<CreateAcademicYearRequest>;

/**
 * Request payload for creating a term
 */
export interface CreateTermRequest {
  name: string;
  startDate: string;
  endDate: string;
  sequence?: number;
}

/**
 * Request payload for updating a term
 */
export type UpdateTermRequest = Partial<CreateTermRequest>;

/**
 * Request payload for creating a holiday
 */
export interface CreateHolidayRequest {
  name: string;
  date: string;
  type?: HolidayType;
  isOptional?: boolean;
  branchId?: string;
}

/**
 * Request payload for updating a holiday
 */
export type UpdateHolidayRequest = Partial<CreateHolidayRequest>;

/**
 * Holiday type options for dropdown
 */
export const HOLIDAY_TYPES: { value: HolidayType; label: string }[] = [
  { value: 'public', label: 'Public Holiday' },
  { value: 'religious', label: 'Religious Holiday' },
  { value: 'school', label: 'School Event' },
  { value: 'exam', label: 'Examination' },
  { value: 'other', label: 'Other' },
];

/**
 * Status options for academic year
 */
export const ACADEMIC_YEAR_STATUS: { value: boolean; label: string }[] = [
  { value: true, label: 'Active' },
  { value: false, label: 'Inactive' },
];
