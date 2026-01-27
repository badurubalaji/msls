/**
 * MSLS Document Model
 *
 * Defines document entities and related types for student document management.
 */

/** Document verification status */
export type DocumentStatus = 'pending_verification' | 'verified' | 'rejected';

/**
 * Document type interface representing configurable document types
 */
export interface DocumentType {
  /** Unique document type identifier */
  id: string;

  /** Document type code (e.g., 'birth_certificate', 'aadhaar') */
  code: string;

  /** Display name */
  name: string;

  /** Description (optional) */
  description?: string;

  /** Whether this document is mandatory */
  isMandatory: boolean;

  /** Whether this document type has an expiry date */
  hasExpiry: boolean;

  /** Allowed file extensions (comma-separated) */
  allowedExtensions: string;

  /** Maximum file size in MB */
  maxSizeMb: number;

  /** Sort order for display */
  sortOrder: number;

  /** Whether this type is active */
  isActive: boolean;

  /** Record creation timestamp */
  createdAt: string;

  /** Last update timestamp */
  updatedAt: string;
}

/**
 * Student document interface
 */
export interface StudentDocument {
  /** Unique document identifier */
  id: string;

  /** Student ID */
  studentId: string;

  /** Document type ID */
  documentTypeId: string;

  /** Document type details */
  documentType?: DocumentType;

  /** File URL */
  fileUrl: string;

  /** Original file name */
  fileName: string;

  /** File size in bytes */
  fileSizeBytes: number;

  /** MIME type */
  mimeType: string;

  /** Document number (e.g., Aadhaar number) */
  documentNumber?: string;

  /** Issue date (ISO date string) */
  issueDate?: string;

  /** Expiry date (ISO date string) */
  expiryDate?: string;

  /** Verification status */
  status: DocumentStatus;

  /** Rejection reason (if rejected) */
  rejectionReason?: string;

  /** Verification timestamp */
  verifiedAt?: string;

  /** Verifier user ID */
  verifiedBy?: string;

  /** Verifier name */
  verifierName?: string;

  /** Upload timestamp */
  uploadedAt: string;

  /** Uploader user ID */
  uploadedBy: string;

  /** Uploader name */
  uploaderName?: string;

  /** Last update timestamp */
  updatedAt: string;

  /** Optimistic lock version */
  version: number;
}

/**
 * Document checklist item
 */
export interface DocumentChecklistItem {
  /** Document type */
  documentType: DocumentType;

  /** Uploaded document (if any) */
  document?: StudentDocument;

  /** Whether this document is required */
  isRequired: boolean;

  /** Whether a document has been uploaded */
  isUploaded: boolean;

  /** Whether the document is verified */
  isVerified: boolean;
}

/**
 * Document checklist response
 */
export interface DocumentChecklistResponse {
  /** Checklist items */
  items: DocumentChecklistItem[];

  /** Total required documents */
  totalRequired: number;

  /** Total uploaded documents */
  totalUploaded: number;

  /** Total verified documents */
  totalVerified: number;

  /** Total pending documents */
  totalPending: number;

  /** Total rejected documents */
  totalRejected: number;

  /** Completion percentage */
  completionPercent: number;
}

/**
 * Document type list response
 */
export interface DocumentTypeListResponse {
  /** Document types */
  documentTypes: DocumentType[];

  /** Total count */
  total: number;
}

/**
 * Document list response
 */
export interface DocumentListResponse {
  /** Documents */
  documents: StudentDocument[];

  /** Total count */
  total: number;
}

/**
 * Upload document request
 */
export interface UploadDocumentRequest {
  /** Document type ID */
  documentTypeId: string;

  /** Document number */
  documentNumber?: string;

  /** Issue date (YYYY-MM-DD) */
  issueDate?: string;

  /** Expiry date (YYYY-MM-DD) */
  expiryDate?: string;
}

/**
 * Update document request
 */
export interface UpdateDocumentRequest {
  /** Document number */
  documentNumber?: string;

  /** Issue date (YYYY-MM-DD) */
  issueDate?: string;

  /** Expiry date (YYYY-MM-DD) */
  expiryDate?: string;

  /** Version for optimistic locking */
  version: number;
}

/**
 * Verify document request
 */
export interface VerifyDocumentRequest {
  /** Version for optimistic locking */
  version: number;
}

/**
 * Reject document request
 */
export interface RejectDocumentRequest {
  /** Rejection reason */
  reason: string;

  /** Version for optimistic locking */
  version: number;
}

/**
 * Get status badge variant based on document status
 */
export function getDocumentStatusBadgeVariant(status: DocumentStatus): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (status) {
    case 'verified':
      return 'success';
    case 'pending_verification':
      return 'warning';
    case 'rejected':
      return 'danger';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for document status
 */
export function getDocumentStatusLabel(status: DocumentStatus): string {
  switch (status) {
    case 'verified':
      return 'Verified';
    case 'pending_verification':
      return 'Pending';
    case 'rejected':
      return 'Rejected';
    default:
      return status;
  }
}

/**
 * Format file size for display
 */
export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

/**
 * Check if a file extension is allowed
 */
export function isExtensionAllowed(fileName: string, allowedExtensions: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase() || '';
  const allowed = allowedExtensions.toLowerCase().split(',').map(e => e.trim());
  return allowed.includes(ext);
}

/**
 * Get file icon based on MIME type
 */
export function getFileIcon(mimeType: string): string {
  if (mimeType.startsWith('image/')) return 'fa-file-image';
  if (mimeType === 'application/pdf') return 'fa-file-pdf';
  return 'fa-file';
}
