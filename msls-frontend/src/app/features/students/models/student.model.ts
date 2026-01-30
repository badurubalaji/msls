/**
 * MSLS Student Model
 *
 * Defines student entities and related types for student management.
 */

/** Student status options */
export type StudentStatus = 'active' | 'inactive' | 'transferred' | 'graduated';

/** Gender options */
export type Gender = 'male' | 'female' | 'other';

/** Address type options */
export type AddressType = 'current' | 'permanent';

/**
 * Student address interface
 */
export interface StudentAddress {
  /** Unique address identifier */
  id: string;

  /** Address type: current or permanent */
  addressType: AddressType;

  /** Primary address line */
  addressLine1: string;

  /** Secondary address line (optional) */
  addressLine2?: string;

  /** City */
  city: string;

  /** State/Province */
  state: string;

  /** Postal/ZIP code */
  postalCode: string;

  /** Country */
  country: string;
}

/**
 * Student interface representing a student profile
 */
export interface Student {
  /** Unique student identifier */
  id: string;

  /** Tenant ID for multi-tenancy */
  tenantId: string;

  /** Branch ID */
  branchId: string;

  /** Auto-generated admission number (e.g., MUM-2026-00001) */
  admissionNumber: string;

  /** Student's first name */
  firstName: string;

  /** Student's middle name (optional) */
  middleName?: string;

  /** Student's last name */
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

  /** Aadhaar number (optional, 12 digits) */
  aadhaarNumber?: string;

  /** Photo URL (optional) */
  photoUrl?: string;

  /** Birth certificate URL (optional) */
  birthCertificateUrl?: string;

  /** Student status */
  status: StudentStatus;

  /** Admission date (ISO date string) */
  admissionDate: string;

  /** Current address (optional) */
  currentAddress?: StudentAddress;

  /** Permanent address (optional) */
  permanentAddress?: StudentAddress;

  /** Branch name (for display, optional) */
  branchName?: string;

  /** Class name (for display, optional) */
  className?: string;

  /** Section name (for display, optional) */
  sectionName?: string;

  /** Record creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;

  /** Optimistic lock version */
  version: number;
}

/**
 * DTO for creating a new student
 */
export interface CreateStudentRequest {
  /** Branch ID where the student is being admitted */
  branchId: string;

  /** Student's first name */
  firstName: string;

  /** Student's middle name (optional) */
  middleName?: string;

  /** Student's last name */
  lastName: string;

  /** Date of birth (ISO date string) */
  dateOfBirth: string;

  /** Gender */
  gender: Gender;

  /** Blood group (optional) */
  bloodGroup?: string;

  /** Aadhaar number (optional, 12 digits) */
  aadhaarNumber?: string;

  /** Admission date (optional, defaults to current date) */
  admissionDate?: string;

  /** Current address (optional) */
  currentAddress?: AddressRequest;

  /** Permanent address (optional) */
  permanentAddress?: AddressRequest;

  /** Whether permanent address is same as current */
  sameAsCurrentAddress?: boolean;
}

/**
 * DTO for updating a student
 */
export interface UpdateStudentRequest {
  /** Student's first name (optional) */
  firstName?: string;

  /** Student's middle name (optional) */
  middleName?: string;

  /** Student's last name (optional) */
  lastName?: string;

  /** Date of birth (optional) */
  dateOfBirth?: string;

  /** Gender (optional) */
  gender?: Gender;

  /** Blood group (optional) */
  bloodGroup?: string;

  /** Aadhaar number (optional) */
  aadhaarNumber?: string;

  /** Student status (optional) */
  status?: StudentStatus;

  /** Photo URL (optional) */
  photoUrl?: string;

  /** Birth certificate URL (optional) */
  birthCertificateUrl?: string;

  /** Current address (optional) */
  currentAddress?: AddressRequest;

  /** Permanent address (optional) */
  permanentAddress?: AddressRequest;

  /** Whether permanent address is same as current */
  sameAsCurrentAddress?: boolean;

  /** Version for optimistic locking */
  version: number;
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
  postalCode: string;

  /** Country (defaults to 'India') */
  country?: string;
}

/**
 * Student list filter options
 */
export interface StudentListFilter {
  /** Filter by branch ID */
  branchId?: string;

  /** Filter by class ID (via active enrollment) */
  classId?: string;

  /** Filter by section ID (via active enrollment) */
  sectionId?: string;

  /** Filter by status */
  status?: StudentStatus;

  /** Filter by gender */
  gender?: Gender;

  /** Filter by admission date from (YYYY-MM-DD) */
  admissionFrom?: string;

  /** Filter by admission date to (YYYY-MM-DD) */
  admissionTo?: string;

  /** Search by name or admission number */
  search?: string;

  /** Cursor for pagination */
  cursor?: string;

  /** Number of items per page */
  limit?: number;

  /** Sort by field */
  sortBy?: 'name' | 'admission_number' | 'created_at';

  /** Sort order */
  sortOrder?: 'asc' | 'desc';
}

/**
 * Student list response with pagination
 */
export interface StudentListResponse {
  /** List of students */
  students: Student[];

  /** Cursor for the next page (UUID of the last student) */
  nextCursor?: string;

  /** Whether there are more results */
  hasMore: boolean;

  /** Total count of students matching the filter */
  total: number;
}

/**
 * Get status badge variant based on student status
 */
export function getStatusBadgeVariant(status: StudentStatus): 'success' | 'warning' | 'info' | 'neutral' {
  switch (status) {
    case 'active':
      return 'success';
    case 'inactive':
      return 'neutral';
    case 'transferred':
      return 'warning';
    case 'graduated':
      return 'info';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for student status
 */
export function getStatusLabel(status: StudentStatus): string {
  switch (status) {
    case 'active':
      return 'Active';
    case 'inactive':
      return 'Inactive';
    case 'transferred':
      return 'Transferred';
    case 'graduated':
      return 'Graduated';
    default:
      return status;
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
 * Format date of birth to display age
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
export function formatAddress(address: StudentAddress | undefined): string {
  if (!address) return '';
  const parts = [
    address.addressLine1,
    address.addressLine2,
    address.city,
    address.state,
    address.postalCode,
    address.country,
  ].filter(Boolean);
  return parts.join(', ');
}

// =========================================================================
// Bulk Operations Types
// =========================================================================

/** Bulk operation type */
export type BulkOperationType = 'send_sms' | 'send_email' | 'update_status' | 'export';

/** Bulk operation status */
export type BulkOperationStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';

/** Bulk operation item status */
export type BulkItemStatus = 'pending' | 'success' | 'failed' | 'skipped';

/**
 * Bulk operation response
 */
export interface BulkOperation {
  id: string;
  operationType: BulkOperationType;
  status: BulkOperationStatus;
  totalCount: number;
  processedCount: number;
  successCount: number;
  failureCount: number;
  resultUrl?: string;
  errorMessage?: string;
  startedAt?: string;
  completedAt?: string;
  createdAt: string;
  items?: BulkOperationItem[];
}

/**
 * Bulk operation item
 */
export interface BulkOperationItem {
  id: string;
  studentId: string;
  studentName?: string;
  status: BulkItemStatus;
  errorMessage?: string;
  processedAt?: string;
}

/**
 * Bulk status update request
 */
export interface BulkStatusUpdateRequest {
  studentIds: string[];
  newStatus: StudentStatus;
}

/**
 * Export request
 */
export interface ExportRequest {
  studentIds: string[];
  format: 'xlsx' | 'csv';
  columns?: string[];
}

/**
 * Export column definition
 */
export interface ExportColumn {
  key: string;
  label: string;
  selected: boolean;
}

/**
 * Default export columns
 */
export const DEFAULT_EXPORT_COLUMNS: ExportColumn[] = [
  { key: 'admission_number', label: 'Admission Number', selected: true },
  { key: 'first_name', label: 'First Name', selected: true },
  { key: 'last_name', label: 'Last Name', selected: true },
  { key: 'full_name', label: 'Full Name', selected: false },
  { key: 'gender', label: 'Gender', selected: false },
  { key: 'date_of_birth', label: 'Date of Birth', selected: false },
  { key: 'class', label: 'Class', selected: true },
  { key: 'section', label: 'Section', selected: true },
  { key: 'roll_number', label: 'Roll Number', selected: false },
  { key: 'guardian_name', label: 'Guardian Name', selected: true },
  { key: 'guardian_phone', label: 'Guardian Phone', selected: true },
  { key: 'status', label: 'Status', selected: true },
  { key: 'branch', label: 'Branch', selected: false },
  { key: 'blood_group', label: 'Blood Group', selected: false },
  { key: 'admission_date', label: 'Admission Date', selected: false },
  { key: 'address', label: 'Address', selected: false },
  { key: 'city', label: 'City', selected: false },
  { key: 'state', label: 'State', selected: false },
];

/**
 * Get bulk operation status badge variant
 */
export function getBulkStatusVariant(
  status: BulkOperationStatus
): 'success' | 'warning' | 'danger' | 'info' | 'neutral' {
  switch (status) {
    case 'completed':
      return 'success';
    case 'processing':
      return 'info';
    case 'pending':
      return 'neutral';
    case 'failed':
      return 'danger';
    case 'cancelled':
      return 'warning';
    default:
      return 'neutral';
  }
}

/**
 * Get bulk operation status label
 */
export function getBulkStatusLabel(status: BulkOperationStatus): string {
  switch (status) {
    case 'completed':
      return 'Completed';
    case 'processing':
      return 'Processing';
    case 'pending':
      return 'Pending';
    case 'failed':
      return 'Failed';
    case 'cancelled':
      return 'Cancelled';
    default:
      return status;
  }
}

// =========================================================================
// Import Types
// =========================================================================

/**
 * Import error details
 */
export interface ImportError {
  /** Row number in the file */
  row: number;
  /** Column name where error occurred */
  column?: string;
  /** Error message */
  message: string;
}

/**
 * Import result
 */
export interface ImportResult {
  /** Total rows processed */
  totalRows: number;
  /** Number of successful imports */
  successCount: number;
  /** Number of failed imports */
  failedCount: number;
  /** List of errors encountered */
  errors?: ImportError[];
  /** IDs of created students */
  createdIds?: string[];
}
