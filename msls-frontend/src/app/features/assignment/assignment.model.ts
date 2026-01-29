/**
 * Teacher Assignment model interfaces for the frontend
 * Story 5.7: Teacher Subject Assignment
 */

// Assignment Status
export type AssignmentStatus = 'active' | 'inactive';

// Workload Status
export type WorkloadStatus = 'under' | 'normal' | 'over';

/**
 * Teacher Subject Assignment interface
 */
export interface Assignment {
  id: string;
  staffId: string;
  staffName?: string;
  staffEmployeeId?: string;
  subjectId: string;
  subjectName?: string;
  subjectCode?: string;
  classId: string;
  className?: string;
  classCode?: string;
  sectionId?: string;
  sectionName?: string;
  academicYearId: string;
  academicYearName?: string;
  periodsPerWeek: number;
  isClassTeacher: boolean;
  effectiveFrom: string;
  effectiveTo?: string;
  status: AssignmentStatus;
  remarks?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Assignment list response
 */
export interface AssignmentListResponse {
  assignments: Assignment[];
  nextCursor?: string;
  hasMore: boolean;
  total: number;
}

/**
 * Create assignment request DTO
 */
export interface CreateAssignmentRequest {
  staffId: string;
  subjectId: string;
  classId: string;
  sectionId?: string;
  academicYearId: string;
  periodsPerWeek: number;
  isClassTeacher?: boolean;
  effectiveFrom: string;
  effectiveTo?: string;
  remarks?: string;
}

/**
 * Update assignment request DTO
 */
export interface UpdateAssignmentRequest {
  periodsPerWeek?: number;
  isClassTeacher?: boolean;
  effectiveFrom?: string;
  effectiveTo?: string;
  status?: AssignmentStatus;
  remarks?: string;
}

/**
 * Bulk create assignment request
 */
export interface BulkCreateRequest {
  assignments: BulkAssignmentItem[];
}

/**
 * Single item in bulk create request
 */
export interface BulkAssignmentItem {
  staffId: string;
  subjectId: string;
  classId: string;
  sectionId?: string;
  academicYearId: string;
  periodsPerWeek: number;
  isClassTeacher?: boolean;
  effectiveFrom: string;
}

/**
 * Workload summary for a teacher
 */
export interface WorkloadSummary {
  staffId: string;
  staffName: string;
  staffEmployeeId: string;
  departmentName?: string;
  totalPeriods: number;
  totalSubjects: number;
  totalClasses: number;
  isClassTeacher: boolean;
  classTeacherFor?: string;
  workloadStatus: WorkloadStatus;
  minPeriods: number;
  maxPeriods: number;
}

/**
 * Workload report response
 */
export interface WorkloadReportResponse {
  teachers: WorkloadSummary[];
  totalTeachers: number;
  overAssigned: number;
  underAssigned: number;
  normalAssigned: number;
}

/**
 * Unassigned subject item
 */
export interface UnassignedSubject {
  subjectId: string;
  subjectName: string;
  subjectCode: string;
  classId: string;
  className: string;
  sectionId?: string;
  sectionName?: string;
}

/**
 * Unassigned subjects response
 */
export interface UnassignedSubjectsResponse {
  subjects: UnassignedSubject[];
  total: number;
}

/**
 * Class teacher response
 */
export interface ClassTeacherResponse {
  classId: string;
  className?: string;
  sectionId?: string;
  sectionName?: string;
  teacherId?: string;
  teacherName?: string;
  isAssigned: boolean;
}

/**
 * Workload settings
 */
export interface WorkloadSettings {
  id: string;
  branchId: string;
  branchName?: string;
  minPeriodsPerWeek: number;
  maxPeriodsPerWeek: number;
  maxSubjectsPerTeacher?: number;
  maxClassesPerTeacher?: number;
}

/**
 * Update workload settings request
 */
export interface UpdateWorkloadSettingsRequest {
  minPeriodsPerWeek: number;
  maxPeriodsPerWeek: number;
  maxSubjectsPerTeacher?: number;
  maxClassesPerTeacher?: number;
}

/**
 * Get status badge variant based on assignment status
 */
export function getAssignmentStatusVariant(status: AssignmentStatus): 'success' | 'neutral' {
  return status === 'active' ? 'success' : 'neutral';
}

/**
 * Get display label for assignment status
 */
export function getAssignmentStatusLabel(status: AssignmentStatus): string {
  switch (status) {
    case 'active':
      return 'Active';
    case 'inactive':
      return 'Inactive';
    default:
      return status;
  }
}

/**
 * Get workload status badge class
 */
export function getWorkloadStatusClass(status: WorkloadStatus): string {
  switch (status) {
    case 'over':
      return 'badge-red';
    case 'under':
      return 'badge-yellow';
    case 'normal':
      return 'badge-green';
    default:
      return 'badge-gray';
  }
}

/**
 * Get workload status label
 */
export function getWorkloadStatusLabel(status: WorkloadStatus): string {
  switch (status) {
    case 'over':
      return 'Over-assigned';
    case 'under':
      return 'Under-assigned';
    case 'normal':
      return 'Normal';
    default:
      return status;
  }
}

/**
 * Format class-section display string
 */
export function formatClassSection(className: string, sectionName?: string): string {
  return sectionName ? `${className} - ${sectionName}` : className;
}
