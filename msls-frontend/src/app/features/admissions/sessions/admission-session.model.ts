/**
 * MSLS Admission Session Models
 *
 * TypeScript interfaces for admission session management.
 */

/**
 * Admission session status types
 */
export type SessionStatus = 'upcoming' | 'open' | 'closed';

/**
 * Admission session entity returned from API
 */
export interface AdmissionSession {
  id: string;
  name: string;
  academicYearId: string;
  academicYearName?: string;
  branchId?: string;
  branchName?: string;
  startDate: string;
  endDate: string;
  status: SessionStatus;
  applicationFee: number;
  requiredDocuments: string[];
  settings: Record<string, unknown>;
  totalApplications: number;
  totalSeats: number;
  filledSeats: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * Admission seat configuration per class
 */
export interface AdmissionSeat {
  id: string;
  sessionId: string;
  className: string;
  totalSeats: number;
  filledSeats: number;
  waitlistLimit: number;
  reservedSeats: ReservedSeats;
  createdAt: string;
  updatedAt: string;
}

/**
 * Reserved seats configuration by category
 */
export interface ReservedSeats {
  general?: number;
  sc?: number;
  st?: number;
  obc?: number;
  ews?: number;
  management?: number;
  [key: string]: number | undefined;
}

/**
 * Request payload for creating an admission session
 */
export interface CreateSessionRequest {
  name: string;
  academicYearId: string;
  branchId?: string;
  startDate: string;
  endDate: string;
  applicationFee?: number;
  requiredDocuments?: string[];
  settings?: Record<string, unknown>;
}

/**
 * Request payload for updating an admission session
 */
export type UpdateSessionRequest = Partial<CreateSessionRequest>;

/**
 * Request payload for changing session status
 */
export interface ChangeStatusRequest {
  status: SessionStatus;
}

/**
 * Request payload for creating/updating seat configuration
 */
export interface SeatConfigRequest {
  className: string;
  totalSeats: number;
  waitlistLimit?: number;
  reservedSeats?: ReservedSeats;
}

/**
 * Academic year for dropdown
 */
export interface AcademicYear {
  id: string;
  name: string;
  startDate: string;
  endDate: string;
  isCurrent: boolean;
}

/**
 * Common required documents for schools
 */
export const COMMON_DOCUMENTS: string[] = [
  'Birth Certificate',
  'Transfer Certificate',
  'Previous School Report Card',
  'Address Proof',
  'Passport Size Photos',
  'Aadhar Card (Student)',
  'Aadhar Card (Parent/Guardian)',
  'Caste Certificate',
  'Income Certificate',
  'Medical Certificate',
];

/**
 * Common class names for Indian schools
 */
export const CLASS_NAMES: string[] = [
  'Nursery',
  'LKG',
  'UKG',
  'Class 1',
  'Class 2',
  'Class 3',
  'Class 4',
  'Class 5',
  'Class 6',
  'Class 7',
  'Class 8',
  'Class 9',
  'Class 10',
  'Class 11',
  'Class 12',
];

/**
 * Get status badge configuration
 */
export function getStatusConfig(status: SessionStatus): {
  label: string;
  class: string;
  icon: string;
} {
  switch (status) {
    case 'upcoming':
      return {
        label: 'Upcoming',
        class: 'badge-blue',
        icon: 'fa-clock',
      };
    case 'open':
      return {
        label: 'Open',
        class: 'badge-green',
        icon: 'fa-door-open',
      };
    case 'closed':
      return {
        label: 'Closed',
        class: 'badge-gray',
        icon: 'fa-door-closed',
      };
    default:
      return {
        label: status,
        class: 'badge-gray',
        icon: 'fa-circle',
      };
  }
}
