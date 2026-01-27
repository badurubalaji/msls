/**
 * MSLS Enrollment Model
 *
 * Defines enrollment entities and related types for student enrollment management.
 */

/** Enrollment status options */
export type EnrollmentStatus = 'active' | 'completed' | 'transferred' | 'dropout';

/**
 * Academic year reference for enrollment display
 */
export interface AcademicYearRef {
  /** Unique academic year identifier */
  id: string;

  /** Academic year name (e.g., "2024-2025") */
  name: string;

  /** Start date of the academic year */
  startDate: string;

  /** End date of the academic year */
  endDate: string;

  /** Whether this is the current academic year */
  isCurrent: boolean;
}

/**
 * Student enrollment interface representing an enrollment record
 */
export interface Enrollment {
  /** Unique enrollment identifier */
  id: string;

  /** Student ID */
  studentId: string;

  /** Academic year reference */
  academicYear?: AcademicYearRef;

  /** Class ID (optional until Epic 6) */
  classId?: string;

  /** Section ID (optional) */
  sectionId?: string;

  /** Roll number within the class/section */
  rollNumber?: string;

  /** Class teacher ID (optional until Epic 5) */
  classTeacherId?: string;

  /** Enrollment status */
  status: EnrollmentStatus;

  /** Enrollment date */
  enrollmentDate: string;

  /** Completion date (when academic year ended) */
  completionDate?: string;

  /** Transfer date (if transferred) */
  transferDate?: string;

  /** Transfer reason */
  transferReason?: string;

  /** Dropout date (if dropped out) */
  dropoutDate?: string;

  /** Dropout reason */
  dropoutReason?: string;

  /** Notes */
  notes?: string;

  /** Record creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;
}

/**
 * DTO for creating a new enrollment
 */
export interface CreateEnrollmentRequest {
  /** Academic year ID */
  academicYearId: string;

  /** Class ID (optional) */
  classId?: string;

  /** Section ID (optional) */
  sectionId?: string;

  /** Roll number (optional) */
  rollNumber?: string;

  /** Class teacher ID (optional) */
  classTeacherId?: string;

  /** Enrollment date (optional, defaults to current date) */
  enrollmentDate?: string;

  /** Notes (optional) */
  notes?: string;
}

/**
 * DTO for updating an enrollment
 */
export interface UpdateEnrollmentRequest {
  /** Class ID (optional) */
  classId?: string;

  /** Section ID (optional) */
  sectionId?: string;

  /** Roll number (optional) */
  rollNumber?: string;

  /** Class teacher ID (optional) */
  classTeacherId?: string;

  /** Notes (optional) */
  notes?: string;
}

/**
 * DTO for processing a transfer
 */
export interface TransferRequest {
  /** Transfer date (required) */
  transferDate: string;

  /** Transfer reason (required) */
  transferReason: string;
}

/**
 * DTO for processing a dropout
 */
export interface DropoutRequest {
  /** Dropout date (required) */
  dropoutDate: string;

  /** Dropout reason (required) */
  dropoutReason: string;
}

/**
 * Enrollment history response
 */
export interface EnrollmentHistoryResponse {
  /** List of enrollments */
  enrollments: Enrollment[];

  /** Total count */
  total: number;
}

/**
 * Enrollment status change log entry
 */
export interface EnrollmentStatusChange {
  /** Unique identifier */
  id: string;

  /** Enrollment ID */
  enrollmentId: string;

  /** Previous status */
  fromStatus?: EnrollmentStatus;

  /** New status */
  toStatus: EnrollmentStatus;

  /** Reason for the change */
  changeReason?: string;

  /** Date of the change */
  changeDate: string;

  /** Timestamp of the change */
  changedAt: string;

  /** ID of the user who made the change */
  changedBy: string;
}

/**
 * Get status badge variant based on enrollment status
 */
export function getEnrollmentStatusBadgeVariant(status: EnrollmentStatus): 'success' | 'warning' | 'error' | 'neutral' {
  switch (status) {
    case 'active':
      return 'success';
    case 'completed':
      return 'neutral';
    case 'transferred':
      return 'warning';
    case 'dropout':
      return 'error';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for enrollment status
 */
export function getEnrollmentStatusLabel(status: EnrollmentStatus): string {
  switch (status) {
    case 'active':
      return 'Active';
    case 'completed':
      return 'Completed';
    case 'transferred':
      return 'Transferred';
    case 'dropout':
      return 'Dropout';
    default:
      return status;
  }
}

/**
 * Get timeline dot color class based on enrollment status
 */
export function getEnrollmentTimelineColor(status: EnrollmentStatus): string {
  switch (status) {
    case 'active':
      return 'bg-green-500';
    case 'completed':
      return 'bg-gray-400';
    case 'transferred':
      return 'bg-orange-500';
    case 'dropout':
      return 'bg-red-500';
    default:
      return 'bg-gray-400';
  }
}
