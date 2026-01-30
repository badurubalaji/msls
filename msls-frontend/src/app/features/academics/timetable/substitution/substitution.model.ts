/**
 * Substitution Management Models
 */

export interface Substitution {
  id: string;
  branchId: string;
  branchName?: string;
  originalStaffId: string;
  originalStaffName?: string;
  substituteStaffId: string;
  substituteStaffName?: string;
  substitutionDate: string;
  reason?: string;
  status: SubstitutionStatus;
  notes?: string;
  createdBy?: string;
  createdByName?: string;
  approvedBy?: string;
  approvedByName?: string;
  approvedAt?: string;
  periods?: SubstitutionPeriod[];
  createdAt: string;
  updatedAt: string;
}

export type SubstitutionStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled';

export interface SubstitutionPeriod {
  id: string;
  periodSlotId: string;
  periodSlotName?: string;
  startTime?: string;
  endTime?: string;
  timetableEntryId?: string;
  subjectId?: string;
  subjectName?: string;
  sectionId?: string;
  sectionName?: string;
  className?: string;
  roomNumber?: string;
  notes?: string;
}

export interface SubstitutionListResponse {
  substitutions: Substitution[];
  total: number;
}

export interface CreateSubstitutionRequest {
  branchId: string;
  originalStaffId: string;
  substituteStaffId: string;
  substitutionDate: string;
  reason?: string;
  notes?: string;
  periods: CreateSubstitutionPeriodInput[];
}

export interface CreateSubstitutionPeriodInput {
  periodSlotId: string;
  timetableEntryId?: string;
  subjectId?: string;
  sectionId?: string;
  roomNumber?: string;
  notes?: string;
}

export interface UpdateSubstitutionRequest {
  substituteStaffId?: string;
  reason?: string;
  notes?: string;
  status?: string;
}

export interface SubstitutionFilter {
  branchId?: string;
  originalStaffId?: string;
  substituteStaffId?: string;
  startDate?: string;
  endDate?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

export interface AvailableTeacher {
  staffId: string;
  staffName: string;
  departmentId?: string;
  department?: string;
  freePeriods: number;
  totalPeriods: number;
  hasConflict: boolean;
}

export interface AvailableTeachersResponse {
  teachers: AvailableTeacher[];
}

export interface TeacherPeriod {
  id: string;
  periodSlotId: string;
  periodSlotName?: string;
  startTime: string;
  endTime: string;
  subjectId?: string;
  subjectName?: string;
  sectionId?: string;
  sectionName?: string;
  className?: string;
  roomNumber?: string;
}

// Status configuration for UI
export const SUBSTITUTION_STATUS_CONFIG: Record<SubstitutionStatus, { label: string; color: string; icon: string }> = {
  pending: { label: 'Pending', color: 'amber', icon: 'clock' },
  confirmed: { label: 'Confirmed', color: 'blue', icon: 'check-circle' },
  completed: { label: 'Completed', color: 'green', icon: 'badge-check' },
  cancelled: { label: 'Cancelled', color: 'red', icon: 'x-circle' },
};
