/**
 * MSLS Staff Model
 *
 * Defines staff entities and related types for staff management.
 */

/** Staff status options */
export type StaffStatus = 'active' | 'inactive' | 'terminated' | 'on_leave';

/** Staff type options */
export type StaffType = 'teaching' | 'non_teaching';

/** Gender options */
export type Gender = 'male' | 'female' | 'other';

/**
 * Staff address interface
 */
export interface StaffAddress {
  /** Primary address line */
  addressLine1: string;

  /** Secondary address line (optional) */
  addressLine2?: string;

  /** City */
  city: string;

  /** State/Province */
  state: string;

  /** Postal/ZIP code */
  pincode: string;

  /** Country */
  country: string;
}

/**
 * Staff reference interface for nested responses
 */
export interface StaffRef {
  /** Staff ID */
  id: string;

  /** Staff name */
  name: string;

  /** Photo URL (optional) */
  photoUrl?: string;
}

/**
 * Staff interface representing a staff profile
 */
export interface Staff {
  /** Unique staff identifier */
  id: string;

  /** Auto-generated employee ID (e.g., EMP00001) */
  employeeId: string;

  /** Employee ID prefix */
  employeeIdPrefix: string;

  /** Staff's first name */
  firstName: string;

  /** Staff's middle name (optional) */
  middleName?: string;

  /** Staff's last name */
  lastName: string;

  /** Full name (computed) */
  fullName: string;

  /** Initials for avatar fallback (computed) */
  initials: string;

  /** Date of birth (ISO date string) */
  dateOfBirth: string;

  /** Gender */
  gender: Gender;

  /** Blood group (optional) */
  bloodGroup?: string;

  /** Nationality */
  nationality?: string;

  /** Religion (optional) */
  religion?: string;

  /** Marital status (optional) */
  maritalStatus?: string;

  /** Personal email (optional) */
  personalEmail?: string;

  /** Work email */
  workEmail: string;

  /** Personal phone (optional) */
  personalPhone?: string;

  /** Work phone */
  workPhone: string;

  /** Emergency contact name (optional) */
  emergencyContactName?: string;

  /** Emergency contact phone (optional) */
  emergencyContactPhone?: string;

  /** Emergency contact relation (optional) */
  emergencyContactRelation?: string;

  /** Current address (optional) */
  currentAddress?: StaffAddress;

  /** Permanent address (optional) */
  permanentAddress?: StaffAddress;

  /** Whether permanent address is same as current */
  sameAsCurrent: boolean;

  /** Staff type */
  staffType: StaffType;

  /** Department ID (optional) */
  departmentId?: string;

  /** Department name (optional) */
  departmentName?: string;

  /** Designation ID (optional) */
  designationId?: string;

  /** Designation name (optional) */
  designationName?: string;

  /** Reporting manager ID (optional) */
  reportingManagerId?: string;

  /** Reporting manager info (optional) */
  reportingManager?: StaffRef;

  /** Join date (ISO date string) */
  joinDate: string;

  /** Confirmation date (optional) */
  confirmationDate?: string;

  /** Probation end date (optional) */
  probationEndDate?: string;

  /** Staff status */
  status: StaffStatus;

  /** Status reason (optional) */
  statusReason?: string;

  /** Termination date (optional) */
  terminationDate?: string;

  /** Photo URL (optional) */
  photoUrl?: string;

  /** Bio (optional) */
  bio?: string;

  /** Branch ID */
  branchId: string;

  /** Branch name (optional) */
  branchName?: string;

  /** Record creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;

  /** Optimistic lock version */
  version: number;
}

/**
 * DTO for creating a new staff member
 */
export interface CreateStaffRequest {
  /** Branch ID */
  branchId: string;

  /** Staff's first name */
  firstName: string;

  /** Staff's middle name (optional) */
  middleName?: string;

  /** Staff's last name */
  lastName: string;

  /** Date of birth (ISO date string) */
  dateOfBirth: string;

  /** Gender */
  gender: Gender;

  /** Blood group (optional) */
  bloodGroup?: string;

  /** Nationality (optional) */
  nationality?: string;

  /** Religion (optional) */
  religion?: string;

  /** Marital status (optional) */
  maritalStatus?: string;

  /** Personal email (optional) */
  personalEmail?: string;

  /** Work email */
  workEmail: string;

  /** Personal phone (optional) */
  personalPhone?: string;

  /** Work phone */
  workPhone: string;

  /** Emergency contact name (optional) */
  emergencyContactName?: string;

  /** Emergency contact phone (optional) */
  emergencyContactPhone?: string;

  /** Emergency contact relation (optional) */
  emergencyContactRelation?: string;

  /** Current address (optional) */
  currentAddress?: AddressRequest;

  /** Permanent address (optional) */
  permanentAddress?: AddressRequest;

  /** Whether permanent address is same as current */
  sameAsCurrent?: boolean;

  /** Staff type */
  staffType: StaffType;

  /** Department ID (optional) */
  departmentId?: string;

  /** Designation ID (optional) */
  designationId?: string;

  /** Reporting manager ID (optional) */
  reportingManagerId?: string;

  /** Join date (ISO date string) */
  joinDate: string;

  /** Confirmation date (optional) */
  confirmationDate?: string;

  /** Probation end date (optional) */
  probationEndDate?: string;

  /** Bio (optional) */
  bio?: string;
}

/**
 * DTO for updating a staff member
 */
export interface UpdateStaffRequest {
  /** Staff's first name (optional) */
  firstName?: string;

  /** Staff's middle name (optional) */
  middleName?: string;

  /** Staff's last name (optional) */
  lastName?: string;

  /** Date of birth (optional) */
  dateOfBirth?: string;

  /** Gender (optional) */
  gender?: Gender;

  /** Blood group (optional) */
  bloodGroup?: string;

  /** Nationality (optional) */
  nationality?: string;

  /** Religion (optional) */
  religion?: string;

  /** Marital status (optional) */
  maritalStatus?: string;

  /** Personal email (optional) */
  personalEmail?: string;

  /** Work email (optional) */
  workEmail?: string;

  /** Personal phone (optional) */
  personalPhone?: string;

  /** Work phone (optional) */
  workPhone?: string;

  /** Emergency contact name (optional) */
  emergencyContactName?: string;

  /** Emergency contact phone (optional) */
  emergencyContactPhone?: string;

  /** Emergency contact relation (optional) */
  emergencyContactRelation?: string;

  /** Current address (optional) */
  currentAddress?: AddressRequest;

  /** Permanent address (optional) */
  permanentAddress?: AddressRequest;

  /** Whether permanent address is same as current (optional) */
  sameAsCurrent?: boolean;

  /** Staff type (optional) */
  staffType?: StaffType;

  /** Department ID (optional) */
  departmentId?: string;

  /** Designation ID (optional) */
  designationId?: string;

  /** Reporting manager ID (optional) */
  reportingManagerId?: string;

  /** Confirmation date (optional) */
  confirmationDate?: string;

  /** Probation end date (optional) */
  probationEndDate?: string;

  /** Bio (optional) */
  bio?: string;

  /** Photo URL (optional) */
  photoUrl?: string;

  /** Version for optimistic locking */
  version: number;
}

/**
 * DTO for updating staff status
 */
export interface StatusUpdateRequest {
  /** New status */
  status: StaffStatus;

  /** Reason for status change */
  reason: string;

  /** Effective date (ISO date string) */
  effectiveDate: string;
}

/**
 * Address request DTO
 */
export interface AddressRequest {
  /** Primary address line */
  addressLine1: string;

  /** Secondary address line (optional) */
  addressLine2?: string;

  /** City */
  city: string;

  /** State/Province */
  state: string;

  /** Postal/ZIP code */
  pincode: string;

  /** Country (defaults to 'India') */
  country?: string;
}

/**
 * Staff list filter options
 */
export interface StaffListFilter {
  /** Filter by branch ID */
  branchId?: string;

  /** Filter by department ID */
  departmentId?: string;

  /** Filter by designation ID */
  designationId?: string;

  /** Filter by staff type */
  staffType?: StaffType;

  /** Filter by status */
  status?: StaffStatus;

  /** Filter by gender */
  gender?: Gender;

  /** Filter by join date from (YYYY-MM-DD) */
  joinDateFrom?: string;

  /** Filter by join date to (YYYY-MM-DD) */
  joinDateTo?: string;

  /** Search by name, employee ID, or work phone */
  search?: string;

  /** Cursor for pagination */
  cursor?: string;

  /** Number of items per page */
  limit?: number;

  /** Sort by field */
  sortBy?: 'name' | 'employee_id' | 'join_date' | 'created_at';

  /** Sort order */
  sortOrder?: 'asc' | 'desc';
}

/**
 * Staff list response with pagination
 */
export interface StaffListResponse {
  /** List of staff */
  staff: Staff[];

  /** Cursor for the next page (UUID of the last staff) */
  nextCursor?: string;

  /** Whether there are more results */
  hasMore: boolean;

  /** Total count of staff matching the filter */
  total: number;
}

/**
 * Status history entry
 */
export interface StatusHistory {
  /** Unique identifier */
  id: string;

  /** Old status (optional for first entry) */
  oldStatus?: string;

  /** New status */
  newStatus: string;

  /** Reason for change (optional) */
  reason?: string;

  /** Effective date */
  effectiveDate: string;

  /** Who made the change (optional) */
  changedBy?: string;

  /** When the change was made */
  changedAt: string;
}

/**
 * Get status badge variant based on staff status
 */
export function getStatusBadgeVariant(status: StaffStatus): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (status) {
    case 'active':
      return 'success';
    case 'inactive':
      return 'neutral';
    case 'terminated':
      return 'danger';
    case 'on_leave':
      return 'warning';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for staff status
 */
export function getStatusLabel(status: StaffStatus): string {
  switch (status) {
    case 'active':
      return 'Active';
    case 'inactive':
      return 'Inactive';
    case 'terminated':
      return 'Terminated';
    case 'on_leave':
      return 'On Leave';
    default:
      return status;
  }
}

/**
 * Get display label for staff type
 */
export function getStaffTypeLabel(type: StaffType): string {
  switch (type) {
    case 'teaching':
      return 'Teaching';
    case 'non_teaching':
      return 'Non-Teaching';
    default:
      return type;
  }
}

/**
 * Get display label for gender
 */
export function getGenderLabel(gender: Gender): string {
  switch (gender) {
    case 'male':
      return 'Male';
    case 'female':
      return 'Female';
    case 'other':
      return 'Other';
    default:
      return gender;
  }
}

/**
 * Calculate age from date of birth
 */
export function calculateAge(dateOfBirth: string): number {
  const today = new Date();
  const birthDate = new Date(dateOfBirth);
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--;
  }
  return age;
}

/**
 * Format address to single line display
 */
export function formatAddress(address: StaffAddress | undefined): string {
  if (!address) return '';
  const parts = [
    address.addressLine1,
    address.addressLine2,
    address.city,
    address.state,
    address.pincode,
    address.country,
  ].filter(Boolean);
  return parts.join(', ');
}

/**
 * Calculate tenure in years and months
 */
export function calculateTenure(joinDate: string): string {
  const today = new Date();
  const join = new Date(joinDate);
  let years = today.getFullYear() - join.getFullYear();
  let months = today.getMonth() - join.getMonth();

  if (months < 0) {
    years--;
    months += 12;
  }

  if (today.getDate() < join.getDate()) {
    months--;
    if (months < 0) {
      years--;
      months += 12;
    }
  }

  if (years === 0) {
    return `${months} month${months !== 1 ? 's' : ''}`;
  }
  if (months === 0) {
    return `${years} year${years !== 1 ? 's' : ''}`;
  }
  return `${years} year${years !== 1 ? 's' : ''}, ${months} month${months !== 1 ? 's' : ''}`;
}
