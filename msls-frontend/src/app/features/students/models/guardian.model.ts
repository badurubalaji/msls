/**
 * Guardian and Emergency Contact Models
 */

/**
 * Guardian relation types
 */
export type GuardianRelation =
  | 'father'
  | 'mother'
  | 'grandfather'
  | 'grandmother'
  | 'uncle'
  | 'aunt'
  | 'sibling'
  | 'guardian'
  | 'other';

/**
 * Guardian interface
 */
export interface Guardian {
  id: string;
  studentId: string;
  relation: GuardianRelation;
  firstName: string;
  lastName: string;
  fullName: string;
  phone: string;
  email?: string;
  occupation?: string;
  annualIncome?: string;
  education?: string;
  isPrimary: boolean;
  hasPortalAccess: boolean;
  userId?: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Guardian list response
 */
export interface GuardianListResponse {
  guardians: Guardian[];
  total: number;
}

/**
 * Create guardian request
 */
export interface CreateGuardianRequest {
  relation: GuardianRelation;
  firstName: string;
  lastName: string;
  phone: string;
  email?: string;
  occupation?: string;
  annualIncome?: string;
  education?: string;
  isPrimary?: boolean;
  hasPortalAccess?: boolean;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
}

/**
 * Update guardian request
 */
export interface UpdateGuardianRequest {
  relation?: GuardianRelation;
  firstName?: string;
  lastName?: string;
  phone?: string;
  email?: string;
  occupation?: string;
  annualIncome?: string;
  education?: string;
  isPrimary?: boolean;
  hasPortalAccess?: boolean;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
}

/**
 * Emergency contact interface
 */
export interface EmergencyContact {
  id: string;
  studentId: string;
  name: string;
  relation: string;
  phone: string;
  alternatePhone?: string;
  priority: number;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Emergency contact list response
 */
export interface EmergencyContactListResponse {
  contacts: EmergencyContact[];
  total: number;
}

/**
 * Create emergency contact request
 */
export interface CreateEmergencyContactRequest {
  name: string;
  relation: string;
  phone: string;
  alternatePhone?: string;
  priority?: number;
  notes?: string;
}

/**
 * Update emergency contact request
 */
export interface UpdateEmergencyContactRequest {
  name?: string;
  relation?: string;
  phone?: string;
  alternatePhone?: string;
  priority?: number;
  notes?: string;
}

/**
 * Guardian relation labels for display
 */
export const GUARDIAN_RELATION_LABELS: Record<GuardianRelation, string> = {
  father: 'Father',
  mother: 'Mother',
  grandfather: 'Grandfather',
  grandmother: 'Grandmother',
  uncle: 'Uncle',
  aunt: 'Aunt',
  sibling: 'Sibling',
  guardian: 'Legal Guardian',
  other: 'Other',
};

/**
 * Get guardian relation label
 */
export function getGuardianRelationLabel(relation: GuardianRelation): string {
  return GUARDIAN_RELATION_LABELS[relation] || relation;
}

/**
 * Guardian relation options for dropdowns
 */
export const GUARDIAN_RELATION_OPTIONS: { value: GuardianRelation; label: string }[] = [
  { value: 'father', label: 'Father' },
  { value: 'mother', label: 'Mother' },
  { value: 'grandfather', label: 'Grandfather' },
  { value: 'grandmother', label: 'Grandmother' },
  { value: 'uncle', label: 'Uncle' },
  { value: 'aunt', label: 'Aunt' },
  { value: 'sibling', label: 'Sibling' },
  { value: 'guardian', label: 'Legal Guardian' },
  { value: 'other', label: 'Other' },
];
