/**
 * Behavioral Incident Models
 *
 * Type definitions for student behavioral tracking.
 * Note: Types are prefixed with "Behavioral" to avoid conflicts with health incident types.
 */

// =============================================================================
// Enums / Types
// =============================================================================

export type BehavioralIncidentType = 'positive_recognition' | 'minor_infraction' | 'major_violation';
export type BehavioralSeverity = 'low' | 'medium' | 'high' | 'critical';
export type FollowUpStatus = 'pending' | 'completed' | 'cancelled';

// =============================================================================
// Interfaces
// =============================================================================

export interface BehavioralIncident {
  id: string;
  studentId: string;
  incidentType: BehavioralIncidentType;
  incidentTypeLabel: string;
  severity: BehavioralSeverity;
  severityLabel: string;
  incidentDate: string;
  incidentTime: string;
  location?: string;
  description: string;
  witnesses?: string[];
  studentResponse?: string;
  actionTaken: string;
  parentMeetingRequired: boolean;
  parentNotified: boolean;
  parentNotifiedAt?: string;
  reportedBy: string;
  reporterName?: string;
  followUps?: FollowUp[];
  createdAt: string;
  updatedAt: string;
}

export interface FollowUp {
  id: string;
  incidentId: string;
  scheduledDate: string;
  scheduledTime?: string;
  participants?: Participant[];
  expectedOutcomes?: string;
  meetingNotes?: string;
  actualOutcomes?: string;
  status: FollowUpStatus;
  statusLabel: string;
  completedAt?: string;
  createdAt: string;
}

export interface Participant {
  name: string;
  role: string;
}

export interface BehaviorSummary {
  totalIncidents: number;
  positiveCount: number;
  minorInfractionCount: number;
  majorViolationCount: number;
  thisMonthCount: number;
  lastMonthCount: number;
  trend: 'improving' | 'declining' | 'stable';
  pendingFollowUps: number;
}

export interface BehavioralIncidentListResponse {
  incidents: BehavioralIncident[];
  total: number;
}

export interface PendingFollowUpItem extends FollowUp {
  studentId: string;
  studentName: string;
  incidentType: string;
  incidentDate: string;
}

export interface PendingFollowUpsResponse {
  followUps: PendingFollowUpItem[];
  total: number;
}

// =============================================================================
// Request DTOs
// =============================================================================

export interface CreateBehavioralIncidentRequest {
  incidentType: BehavioralIncidentType;
  severity?: BehavioralSeverity;
  incidentDate: string;
  incidentTime: string;
  location?: string;
  description: string;
  witnesses?: string[];
  studentResponse?: string;
  actionTaken: string;
  parentMeetingRequired?: boolean;
}

export interface UpdateBehavioralIncidentRequest {
  incidentType?: BehavioralIncidentType;
  severity?: BehavioralSeverity;
  incidentDate?: string;
  incidentTime?: string;
  location?: string;
  description?: string;
  witnesses?: string[];
  studentResponse?: string;
  actionTaken?: string;
  parentMeetingRequired?: boolean;
  parentNotified?: boolean;
}

export interface CreateFollowUpRequest {
  scheduledDate: string;
  scheduledTime?: string;
  participants?: Participant[];
  expectedOutcomes?: string;
}

export interface UpdateFollowUpRequest {
  scheduledDate?: string;
  scheduledTime?: string;
  participants?: Participant[];
  expectedOutcomes?: string;
  meetingNotes?: string;
  actualOutcomes?: string;
  status?: FollowUpStatus;
}

// =============================================================================
// Filter
// =============================================================================

export interface BehavioralIncidentFilter {
  type?: BehavioralIncidentType;
  severity?: BehavioralSeverity;
  dateFrom?: string;
  dateTo?: string;
  limit?: number;
  offset?: number;
}

// =============================================================================
// Options
// =============================================================================

export const BEHAVIORAL_INCIDENT_TYPE_OPTIONS: { value: BehavioralIncidentType; label: string }[] = [
  { value: 'positive_recognition', label: 'Positive Recognition' },
  { value: 'minor_infraction', label: 'Minor Infraction' },
  { value: 'major_violation', label: 'Major Violation' },
];

export const BEHAVIORAL_SEVERITY_OPTIONS: { value: BehavioralSeverity; label: string }[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'critical', label: 'Critical' },
];

export const FOLLOW_UP_STATUS_OPTIONS: { value: FollowUpStatus; label: string }[] = [
  { value: 'pending', label: 'Pending' },
  { value: 'completed', label: 'Completed' },
  { value: 'cancelled', label: 'Cancelled' },
];

// =============================================================================
// Helper Functions
// =============================================================================

export function getBehavioralIncidentTypeLabel(type: BehavioralIncidentType): string {
  const option = BEHAVIORAL_INCIDENT_TYPE_OPTIONS.find(o => o.value === type);
  return option?.label ?? type;
}

export function getBehavioralSeverityLabel(severity: BehavioralSeverity): string {
  const option = BEHAVIORAL_SEVERITY_OPTIONS.find(o => o.value === severity);
  return option?.label ?? severity;
}

export function getFollowUpStatusLabel(status: FollowUpStatus): string {
  const option = FOLLOW_UP_STATUS_OPTIONS.find(o => o.value === status);
  return option?.label ?? status;
}

export function getBehavioralIncidentTypeVariant(type: BehavioralIncidentType): 'success' | 'warning' | 'danger' {
  switch (type) {
    case 'positive_recognition':
      return 'success';
    case 'minor_infraction':
      return 'warning';
    case 'major_violation':
      return 'danger';
    default:
      return 'warning';
  }
}

export function getBehavioralSeverityVariant(severity: BehavioralSeverity): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (severity) {
    case 'low':
      return 'success';
    case 'medium':
      return 'warning';
    case 'high':
    case 'critical':
      return 'danger';
    default:
      return 'neutral';
  }
}

export function getTrendIcon(trend: string): string {
  switch (trend) {
    case 'improving':
      return 'fa-arrow-trend-up';
    case 'declining':
      return 'fa-arrow-trend-down';
    default:
      return 'fa-minus';
  }
}

export function getTrendVariant(trend: string): 'success' | 'warning' | 'danger' | 'neutral' {
  switch (trend) {
    case 'improving':
      return 'success';
    case 'declining':
      return 'danger';
    default:
      return 'neutral';
  }
}
