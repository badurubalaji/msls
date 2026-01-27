/**
 * Promotion Models
 *
 * TypeScript interfaces for student promotion/retention processing.
 */

import type { AcademicYearRef } from './enrollment.model';

// Re-export for convenience
export type { AcademicYearRef };

// ============================================================================
// Enums
// ============================================================================

export type PromotionDecision = 'pending' | 'promote' | 'retain' | 'transfer';
export type BatchStatus = 'draft' | 'processing' | 'completed' | 'cancelled';

// ============================================================================
// Promotion Rules
// ============================================================================

export interface PromotionRule {
  id: string;
  classId: string;
  minAttendancePct?: number;
  minOverallMarksPct?: number;
  minSubjectsPassed: number;
  autoPromoteOnCriteria: boolean;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRuleRequest {
  classId: string;
  minAttendancePct?: number;
  minOverallMarksPct?: number;
  minSubjectsPassed?: number;
  autoPromoteOnCriteria?: boolean;
}

export interface RuleListResponse {
  rules: PromotionRule[];
  total: number;
}

export interface PromotionBatch {
  id: string;
  fromAcademicYear?: AcademicYearRef;
  toAcademicYear?: AcademicYearRef;
  fromClassId: string;
  fromSectionId?: string;
  toClassId?: string;
  status: BatchStatus;
  totalStudents: number;
  promotedCount: number;
  retainedCount: number;
  transferredCount: number;
  pendingCount: number;
  processedAt?: string;
  processedBy?: string;
  cancelledAt?: string;
  cancelledBy?: string;
  cancellationReason?: string;
  notes?: string;
  createdAt: string;
  createdBy?: string;
}

export interface CreateBatchRequest {
  fromAcademicYearId: string;
  toAcademicYearId: string;
  fromClassId: string;
  fromSectionId?: string;
  toClassId?: string;
  notes?: string;
}

export interface BatchListResponse {
  batches: PromotionBatch[];
  total: number;
}

// ============================================================================
// Promotion Records
// ============================================================================

export interface StudentRef {
  id: string;
  admissionNumber: string;
  fullName: string;
  photoUrl?: string;
}

export interface PromotionRecord {
  id: string;
  batchId: string;
  student?: StudentRef;
  fromEnrollmentId: string;
  toEnrollmentId?: string;
  decision: PromotionDecision;
  toClassId?: string;
  toSectionId?: string;
  rollNumber?: string;
  autoDecided: boolean;
  decisionReason?: string;
  attendancePct?: number;
  overallMarksPct?: number;
  subjectsPassed?: number;
  overrideBy?: string;
  overrideAt?: string;
  overrideReason?: string;
  retentionReason?: string;
  transferDestination?: string;
  createdAt: string;
  updatedAt: string;
}

export interface RecordsSummary {
  totalStudents: number;
  pendingCount: number;
  promoteCount: number;
  retainCount: number;
  transferCount: number;
  autoDecided: number;
  manualDecided: number;
}

export interface RecordListResponse {
  records: PromotionRecord[];
  total: number;
  summary?: RecordsSummary;
}

export interface UpdateRecordRequest {
  decision?: PromotionDecision;
  toClassId?: string;
  toSectionId?: string;
  rollNumber?: string;
  overrideReason?: string;
  retentionReason?: string;
  transferDestination?: string;
}

export interface BulkUpdateRecordsRequest {
  recordIds: string[];
  decision: PromotionDecision;
  toClassId?: string;
  toSectionId?: string;
  reason?: string;
}

// ============================================================================
// Processing
// ============================================================================

export interface ProcessBatchRequest {
  generateRollNumbers: boolean;
}

export interface CancelBatchRequest {
  reason: string;
}

// ============================================================================
// Report
// ============================================================================

export interface PromotionReportRow {
  studentAdmissionNo: string;
  studentName: string;
  fromClass: string;
  fromSection: string;
  decision: string;
  toClass: string;
  toSection: string;
  rollNumber: string;
  attendancePct: string;
  marksPct: string;
  reason: string;
}

export interface PromotionReportResponse {
  rows: PromotionReportRow[];
  total: number;
}

// ============================================================================
// Helper Functions
// ============================================================================

export function getDecisionBadgeVariant(decision: PromotionDecision): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (decision) {
    case 'promote':
      return 'success';
    case 'retain':
      return 'warning';
    case 'transfer':
      return 'danger';
    case 'pending':
    default:
      return 'neutral';
  }
}

export function getDecisionLabel(decision: PromotionDecision): string {
  switch (decision) {
    case 'promote':
      return 'Promote';
    case 'retain':
      return 'Retain';
    case 'transfer':
      return 'Transfer';
    case 'pending':
    default:
      return 'Pending';
  }
}

export function getBatchStatusBadgeVariant(status: BatchStatus): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (status) {
    case 'completed':
      return 'success';
    case 'processing':
      return 'warning';
    case 'cancelled':
      return 'danger';
    case 'draft':
    default:
      return 'neutral';
  }
}

export function getBatchStatusLabel(status: BatchStatus): string {
  switch (status) {
    case 'completed':
      return 'Completed';
    case 'processing':
      return 'Processing';
    case 'cancelled':
      return 'Cancelled';
    case 'draft':
    default:
      return 'Draft';
  }
}
