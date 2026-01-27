/**
 * MSLS Admission Application Models
 *
 * TypeScript interfaces for online admission application management.
 */

/**
 * Application stage enum - tracks the lifecycle of an admission application
 */
export type ApplicationStage =
  | 'draft'
  | 'submitted'
  | 'under_review'
  | 'documents_pending'
  | 'documents_verified'
  | 'interview_scheduled'
  | 'interview_completed'
  | 'approved'
  | 'rejected'
  | 'waitlisted'
  | 'admitted';

/**
 * Gender enum
 */
export type Gender = 'male' | 'female' | 'other';

/**
 * Blood group enum
 */
export type BloodGroup = 'A+' | 'A-' | 'B+' | 'B-' | 'O+' | 'O-' | 'AB+' | 'AB-';

/**
 * Parent relation type
 */
export type ParentRelation = 'father' | 'mother' | 'guardian';

/**
 * Document type for application
 */
export type DocumentType =
  | 'birth_certificate'
  | 'transfer_certificate'
  | 'report_card'
  | 'address_proof'
  | 'photo'
  | 'aadhaar_student'
  | 'aadhaar_parent'
  | 'caste_certificate'
  | 'income_certificate'
  | 'medical_certificate'
  | 'other';

/**
 * Stage history entry - tracks stage transitions
 */
export interface StageHistoryEntry {
  stage: ApplicationStage;
  timestamp: string;
  remarks?: string;
  changedBy?: string;
  changedByName?: string;
}

/**
 * Application parent/guardian details
 */
export interface ApplicationParent {
  id: string;
  applicationId: string;
  relation: ParentRelation;
  name: string;
  phone?: string;
  email?: string;
  occupation?: string;
  education?: string;
  annualIncome?: string;
  createdAt: string;
}

/**
 * Application document
 */
export interface ApplicationDocument {
  id: string;
  applicationId: string;
  documentType: DocumentType;
  fileUrl: string;
  fileName: string;
  isVerified: boolean;
  verifiedBy?: string;
  verifiedAt?: string;
  rejectionReason?: string;
  createdAt: string;
}

/**
 * Main admission application entity
 */
export interface AdmissionApplication {
  /** Unique identifier (UUID v7) */
  id: string;

  /** Tenant identifier */
  tenantId: string;

  /** Branch identifier */
  branchId?: string;

  /** Admission session identifier */
  sessionId: string;

  /** Session name for display */
  sessionName?: string;

  /** Enquiry ID if converted from enquiry */
  enquiryId?: string;

  /** Auto-generated application number (e.g., APP-2026-001) */
  applicationNumber: string;

  /** Current stage of the application */
  currentStage: ApplicationStage;

  /** History of stage transitions */
  stageHistory: StageHistoryEntry[];

  /** Class applying for */
  classApplying: string;

  // Student Details
  firstName: string;
  middleName?: string;
  lastName: string;
  dateOfBirth: string;
  gender: Gender;
  bloodGroup?: BloodGroup;
  nationality: string;
  religion?: string;
  category?: string;
  aadhaarNumber?: string;
  photoUrl?: string;

  // Previous School
  previousSchool?: string;
  previousClass?: string;
  previousPercentage?: number;
  transferCertificateUrl?: string;

  // Contact/Address
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;

  // Parents (embedded for list view)
  parents?: ApplicationParent[];

  // Documents (embedded for list view)
  documents?: ApplicationDocument[];

  // Timestamps
  submittedAt?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Request payload for creating an application
 */
export interface CreateApplicationRequest {
  sessionId: string;
  branchId?: string;
  enquiryId?: string;
  classApplying: string;

  // Student Details
  firstName: string;
  middleName?: string;
  lastName: string;
  dateOfBirth: string;
  gender: Gender;
  bloodGroup?: BloodGroup;
  nationality?: string;
  religion?: string;
  category?: string;
  aadhaarNumber?: string;

  // Previous School
  previousSchool?: string;
  previousClass?: string;
  previousPercentage?: number;

  // Contact/Address
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;

  // Parent/Guardian Info
  fatherName?: string;
  fatherPhone?: string;
  fatherEmail?: string;
  fatherOccupation?: string;
  motherName?: string;
  motherPhone?: string;
  motherEmail?: string;
  motherOccupation?: string;
  guardianName?: string;
  guardianPhone?: string;
  guardianEmail?: string;
  guardianRelation?: string;
}

/**
 * Request payload for updating an application
 */
export type UpdateApplicationRequest = Partial<CreateApplicationRequest>;

/**
 * Request payload for adding/updating parent
 */
export interface ParentRequest {
  relation: ParentRelation;
  name: string;
  phone?: string;
  email?: string;
  occupation?: string;
  education?: string;
  annualIncome?: string;
}

/**
 * Request payload for updating application stage
 */
export interface UpdateStageRequest {
  stage: ApplicationStage;
  remarks?: string;
}

/**
 * Request for document upload
 */
export interface UploadDocumentRequest {
  documentType: DocumentType;
  file: File;
}

/**
 * Filter parameters for listing applications
 */
export interface ApplicationFilterParams {
  stage?: ApplicationStage;
  classApplying?: string;
  sessionId?: string;
  branchId?: string;
  search?: string;
  fromDate?: string;
  toDate?: string;
}

/**
 * Stage configuration for display
 */
export interface StageConfig {
  label: string;
  variant: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  icon: string;
  description: string;
}

/**
 * Stage display configuration map
 */
export const APPLICATION_STAGE_CONFIG: Record<ApplicationStage, StageConfig> = {
  draft: {
    label: 'Draft',
    variant: 'neutral',
    icon: 'fa-solid fa-file-pen',
    description: 'Application saved as draft',
  },
  submitted: {
    label: 'Submitted',
    variant: 'info',
    icon: 'fa-solid fa-paper-plane',
    description: 'Application submitted for review',
  },
  under_review: {
    label: 'Under Review',
    variant: 'primary',
    icon: 'fa-solid fa-magnifying-glass',
    description: 'Application is being reviewed',
  },
  documents_pending: {
    label: 'Documents Pending',
    variant: 'warning',
    icon: 'fa-solid fa-file-circle-exclamation',
    description: 'Additional documents required',
  },
  documents_verified: {
    label: 'Documents Verified',
    variant: 'success',
    icon: 'fa-solid fa-file-circle-check',
    description: 'All documents verified',
  },
  interview_scheduled: {
    label: 'Interview Scheduled',
    variant: 'primary',
    icon: 'fa-solid fa-calendar-check',
    description: 'Interview has been scheduled',
  },
  interview_completed: {
    label: 'Interview Completed',
    variant: 'info',
    icon: 'fa-solid fa-comments',
    description: 'Interview has been conducted',
  },
  approved: {
    label: 'Approved',
    variant: 'success',
    icon: 'fa-solid fa-circle-check',
    description: 'Application approved for admission',
  },
  rejected: {
    label: 'Rejected',
    variant: 'danger',
    icon: 'fa-solid fa-circle-xmark',
    description: 'Application rejected',
  },
  waitlisted: {
    label: 'Waitlisted',
    variant: 'warning',
    icon: 'fa-solid fa-clock',
    description: 'Application on waitlist',
  },
  admitted: {
    label: 'Admitted',
    variant: 'success',
    icon: 'fa-solid fa-graduation-cap',
    description: 'Student admitted',
  },
};

/**
 * Document type labels
 */
export const DOCUMENT_TYPE_LABELS: Record<DocumentType, string> = {
  birth_certificate: 'Birth Certificate',
  transfer_certificate: 'Transfer Certificate',
  report_card: 'Previous School Report Card',
  address_proof: 'Address Proof',
  photo: 'Passport Size Photo',
  aadhaar_student: 'Aadhaar Card (Student)',
  aadhaar_parent: 'Aadhaar Card (Parent/Guardian)',
  caste_certificate: 'Caste Certificate',
  income_certificate: 'Income Certificate',
  medical_certificate: 'Medical Certificate',
  other: 'Other Document',
};

/**
 * Parent relation labels
 */
export const PARENT_RELATION_LABELS: Record<ParentRelation, string> = {
  father: 'Father',
  mother: 'Mother',
  guardian: 'Guardian',
};

/**
 * Gender labels
 */
export const GENDER_LABELS: Record<Gender, string> = {
  male: 'Male',
  female: 'Female',
  other: 'Other',
};

/**
 * Blood group labels
 */
export const BLOOD_GROUP_OPTIONS: BloodGroup[] = ['A+', 'A-', 'B+', 'B-', 'O+', 'O-', 'AB+', 'AB-'];

/**
 * Common categories for reservation
 */
export const CATEGORY_OPTIONS = [
  { value: 'general', label: 'General' },
  { value: 'sc', label: 'SC (Scheduled Caste)' },
  { value: 'st', label: 'ST (Scheduled Tribe)' },
  { value: 'obc', label: 'OBC (Other Backward Class)' },
  { value: 'ews', label: 'EWS (Economically Weaker Section)' },
];

/**
 * Get stage configuration
 */
export function getStageConfig(stage: ApplicationStage): StageConfig {
  return APPLICATION_STAGE_CONFIG[stage] ?? APPLICATION_STAGE_CONFIG.draft;
}

/**
 * Get badge variant from stage
 */
export function getStageBadgeVariant(stage: ApplicationStage): string {
  return APPLICATION_STAGE_CONFIG[stage]?.variant ?? 'neutral';
}

/**
 * Request payload for public status check
 */
export interface StatusCheckRequest {
  applicationNumber: string;
  phone: string;
}

/**
 * Response for public status check
 */
export interface StatusCheckResponse {
  applicationNumber: string;
  studentName: string;
  className: string;
  status: ApplicationStage;
  submittedAt?: string;
  reviewNotes?: string;
}
