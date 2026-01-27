/**
 * MSLS Application Review Models
 *
 * TypeScript interfaces for application review and document verification.
 */

/**
 * Application status enum
 */
export type ApplicationStatus =
  | 'submitted'
  | 'under_review'
  | 'documents_pending'
  | 'documents_verified'
  | 'test_scheduled'
  | 'test_completed'
  | 'approved'
  | 'rejected'
  | 'waitlisted'
  | 'enrolled';

/**
 * Document verification status
 */
export type DocumentStatus = 'pending' | 'verified' | 'rejected';

/**
 * Review type enum
 */
export type ReviewType = 'document_verification' | 'academic_review' | 'interview' | 'final_decision';

/**
 * Application document interface
 */
export interface ApplicationDocument {
  id: string;
  name: string;
  type: string;
  fileName: string;
  fileUrl: string;
  fileSize: number;
  uploadedAt: string;
  status: DocumentStatus;
  verifiedBy?: string;
  verifiedAt?: string;
  rejectionReason?: string;
}

/**
 * Application review record
 */
export interface ApplicationReview {
  id: string;
  tenantId: string;
  applicationId: string;
  reviewerId: string;
  reviewerName?: string;
  reviewType: ReviewType;
  status: 'approved' | 'rejected' | 'pending';
  comments?: string;
  createdAt: string;
}

/**
 * Admission application interface
 */
export interface AdmissionApplication {
  id: string;
  tenantId: string;
  sessionId: string;
  branchId?: string;
  applicationNumber: string;

  /** Student Details */
  studentName: string;
  dateOfBirth: string;
  gender: string;
  bloodGroup?: string;
  nationality?: string;
  religion?: string;
  category?: string;
  aadharNumber?: string;
  previousSchool?: string;
  previousClass?: string;

  /** Admission Details */
  classApplying: string;
  academicYear: string;

  /** Parent/Guardian Details */
  parentName: string;
  parentPhone: string;
  parentEmail?: string;
  parentOccupation?: string;
  parentAddress?: string;
  fatherName?: string;
  motherName?: string;

  /** Address */
  address: string;
  city?: string;
  state?: string;
  pincode?: string;

  /** Documents */
  documents: ApplicationDocument[];

  /** Status & Reviews */
  status: ApplicationStatus;
  reviews: ApplicationReview[];

  /** Test Information */
  testId?: string;
  testResult?: 'pass' | 'fail' | 'pending';
  testScore?: number;
  testPercentage?: number;

  /** Timestamps */
  submittedAt: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * DTO for adding a review
 */
export interface AddReviewDto {
  reviewType: ReviewType;
  status: 'approved' | 'rejected' | 'pending';
  comments?: string;
}

/**
 * DTO for verifying a document
 */
export interface VerifyDocumentDto {
  status: DocumentStatus;
  rejectionReason?: string;
}

/**
 * DTO for updating application status
 */
export interface UpdateStatusDto {
  status: ApplicationStatus;
  comments?: string;
}

/**
 * Application filter parameters
 */
export interface ApplicationFilterParams {
  sessionId?: string;
  status?: ApplicationStatus;
  classApplying?: string;
  search?: string;
  fromDate?: string;
  toDate?: string;
}

/**
 * Status configuration for display
 */
export interface ApplicationStatusConfig {
  label: string;
  variant: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  icon: string;
  description: string;
}

/**
 * Application status display configuration
 */
export const APPLICATION_STATUS_CONFIG: Record<ApplicationStatus, ApplicationStatusConfig> = {
  submitted: {
    label: 'Submitted',
    variant: 'info',
    icon: 'fa-solid fa-paper-plane',
    description: 'Application has been submitted',
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
    description: 'Some documents need attention',
  },
  documents_verified: {
    label: 'Documents Verified',
    variant: 'success',
    icon: 'fa-solid fa-file-circle-check',
    description: 'All documents verified',
  },
  test_scheduled: {
    label: 'Test Scheduled',
    variant: 'info',
    icon: 'fa-solid fa-calendar-check',
    description: 'Entrance test has been scheduled',
  },
  test_completed: {
    label: 'Test Completed',
    variant: 'primary',
    icon: 'fa-solid fa-clipboard-check',
    description: 'Entrance test completed',
  },
  approved: {
    label: 'Approved',
    variant: 'success',
    icon: 'fa-solid fa-check-circle',
    description: 'Application approved for admission',
  },
  rejected: {
    label: 'Rejected',
    variant: 'danger',
    icon: 'fa-solid fa-times-circle',
    description: 'Application has been rejected',
  },
  waitlisted: {
    label: 'Waitlisted',
    variant: 'warning',
    icon: 'fa-solid fa-hourglass-half',
    description: 'Application is on waitlist',
  },
  enrolled: {
    label: 'Enrolled',
    variant: 'success',
    icon: 'fa-solid fa-graduation-cap',
    description: 'Student has been enrolled',
  },
};

/**
 * Document status display configuration
 */
export const DOCUMENT_STATUS_CONFIG: Record<DocumentStatus, { label: string; variant: string; icon: string }> = {
  pending: {
    label: 'Pending',
    variant: 'warning',
    icon: 'fa-solid fa-clock',
  },
  verified: {
    label: 'Verified',
    variant: 'success',
    icon: 'fa-solid fa-check-circle',
  },
  rejected: {
    label: 'Rejected',
    variant: 'danger',
    icon: 'fa-solid fa-times-circle',
  },
};

/**
 * Review type labels
 */
export const REVIEW_TYPE_LABELS: Record<ReviewType, string> = {
  document_verification: 'Document Verification',
  academic_review: 'Academic Review',
  interview: 'Interview',
  final_decision: 'Final Decision',
};

/**
 * Required documents for admission (configurable per session)
 */
export const DEFAULT_REQUIRED_DOCUMENTS: string[] = [
  'Birth Certificate',
  'Transfer Certificate',
  'Previous School Report Card',
  'Address Proof',
  'Passport Size Photos',
  'Aadhar Card (Student)',
  'Aadhar Card (Parent/Guardian)',
];

/**
 * Format file size for display
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Check if all documents are verified
 */
export function allDocumentsVerified(documents: ApplicationDocument[]): boolean {
  return documents.length > 0 && documents.every(doc => doc.status === 'verified');
}

/**
 * Check if any document is rejected
 */
export function hasRejectedDocuments(documents: ApplicationDocument[]): boolean {
  return documents.some(doc => doc.status === 'rejected');
}

/**
 * Get document verification summary
 */
export function getDocumentSummary(documents: ApplicationDocument[]): {
  total: number;
  verified: number;
  pending: number;
  rejected: number;
} {
  return {
    total: documents.length,
    verified: documents.filter(d => d.status === 'verified').length,
    pending: documents.filter(d => d.status === 'pending').length,
    rejected: documents.filter(d => d.status === 'rejected').length,
  };
}
