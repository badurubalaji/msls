/**
 * Staff Document Management Models
 * Story 5.8: Staff Document Management
 */

// ============================================================================
// Document Type Models
// ============================================================================

export interface DocumentType {
  id: string;
  name: string;
  code: string;
  category: DocumentCategory;
  description?: string;
  is_mandatory: boolean;
  has_expiry: boolean;
  default_validity_months?: number;
  applicable_to: StaffType[];
  is_active: boolean;
  display_order: number;
  created_at: string;
  updated_at: string;
}

export type DocumentCategory = 'identity' | 'education' | 'employment' | 'compliance' | 'other';
export type StaffType = 'teaching' | 'non_teaching';
export type VerificationStatus = 'pending' | 'verified' | 'rejected';

export interface CreateDocumentTypeRequest {
  name: string;
  code: string;
  category: DocumentCategory;
  description?: string;
  is_mandatory: boolean;
  has_expiry: boolean;
  default_validity_months?: number;
  applicable_to?: StaffType[];
  display_order?: number;
}

export interface UpdateDocumentTypeRequest {
  name?: string;
  description?: string;
  is_mandatory?: boolean;
  has_expiry?: boolean;
  default_validity_months?: number;
  applicable_to?: StaffType[];
  is_active?: boolean;
  display_order?: number;
}

export interface DocumentTypeListResponse {
  document_types: DocumentType[];
  total: number;
}

// ============================================================================
// Staff Document Models
// ============================================================================

export interface StaffDocument {
  id: string;
  staff_id: string;
  document_type_id: string;
  document_type?: DocumentType;
  document_number?: string;
  issue_date?: string;
  expiry_date?: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  verification_status: VerificationStatus;
  verified_by?: string;
  verified_at?: string;
  verification_notes?: string;
  rejection_reason?: string;
  remarks?: string;
  is_current: boolean;
  is_expired: boolean;
  is_expiring_soon: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateDocumentRequest {
  document_type_id: string;
  document_number?: string;
  issue_date?: string;
  expiry_date?: string;
  remarks?: string;
}

export interface UpdateDocumentRequest {
  document_number?: string;
  issue_date?: string;
  expiry_date?: string;
  remarks?: string;
}

export interface VerifyDocumentRequest {
  notes?: string;
}

export interface RejectDocumentRequest {
  reason: string;
  notes?: string;
}

export interface DocumentListResponse {
  documents: StaffDocument[];
  next_cursor?: string;
  has_more: boolean;
  total: number;
}

// ============================================================================
// Compliance & Report Models
// ============================================================================

export interface ExpiringDocument {
  document: StaffDocument;
  staff_name: string;
  employee_id: string;
  days_to_expiry: number;
}

export interface ComplianceStats {
  total_staff: number;
  documents_submitted: number;
  pending_verification: number;
  verified: number;
  rejected: number;
  expired: number;
  expiring_in_30_days: number;
  expiring_in_60_days: number;
  expiring_in_90_days: number;
  compliance_percentage: number;
}

export interface DocumentTypeCompliance {
  document_type: DocumentType;
  required: number;
  submitted: number;
  verified: number;
  pending: number;
  rejected: number;
  expired: number;
  compliance_percent: number;
}

export interface StaffComplianceDetail {
  staff_id: string;
  staff_name: string;
  employee_id: string;
  total_required: number;
  submitted: number;
  verified: number;
  pending: number;
  rejected: number;
  expired: number;
  missing_documents: string[];
  compliance_percent: number;
}

export interface ComplianceReportResponse {
  stats: ComplianceStats;
  by_document_type: DocumentTypeCompliance[];
  staff_details?: StaffComplianceDetail[];
}

// ============================================================================
// Helper Functions
// ============================================================================

export const DOCUMENT_CATEGORIES: { value: DocumentCategory; label: string }[] = [
  { value: 'identity', label: 'Identity' },
  { value: 'education', label: 'Education' },
  { value: 'employment', label: 'Employment' },
  { value: 'compliance', label: 'Compliance' },
  { value: 'other', label: 'Other' },
];

export const STAFF_TYPES: { value: StaffType; label: string }[] = [
  { value: 'teaching', label: 'Teaching' },
  { value: 'non_teaching', label: 'Non-Teaching' },
];

export const VERIFICATION_STATUSES: { value: VerificationStatus; label: string; color: string }[] = [
  { value: 'pending', label: 'Pending', color: 'warning' },
  { value: 'verified', label: 'Verified', color: 'success' },
  { value: 'rejected', label: 'Rejected', color: 'danger' },
];

export function getVerificationStatusLabel(status: VerificationStatus): string {
  return VERIFICATION_STATUSES.find(s => s.value === status)?.label ?? status;
}

export function getVerificationStatusColor(status: VerificationStatus): string {
  return VERIFICATION_STATUSES.find(s => s.value === status)?.color ?? 'secondary';
}

export function getCategoryLabel(category: DocumentCategory): string {
  return DOCUMENT_CATEGORIES.find(c => c.value === category)?.label ?? category;
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

export function isDocumentExpiringSoon(expiryDate: string | undefined, days: number = 30): boolean {
  if (!expiryDate) return false;
  const expiry = new Date(expiryDate);
  const threshold = new Date();
  threshold.setDate(threshold.getDate() + days);
  return expiry <= threshold && expiry > new Date();
}

export function getDaysToExpiry(expiryDate: string | undefined): number {
  if (!expiryDate) return -1;
  const expiry = new Date(expiryDate);
  const today = new Date();
  const diffTime = expiry.getTime() - today.getTime();
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}
