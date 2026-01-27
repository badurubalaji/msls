/**
 * MSLS Branch Models
 *
 * TypeScript interfaces for branch management.
 */

/**
 * Branch entity returned from API
 */
export interface Branch {
  id: string;
  code: string;
  name: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country: string;
  phone?: string;
  email?: string;
  logoUrl?: string;
  timezone: string;
  isPrimary: boolean;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * Request payload for creating a branch
 */
export interface CreateBranchRequest {
  code: string;
  name: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
  phone?: string;
  email?: string;
  timezone?: string;
  isPrimary?: boolean;
}

/**
 * Request payload for updating a branch
 */
export type UpdateBranchRequest = Partial<CreateBranchRequest>;

/**
 * Common timezones for dropdown
 */
export const TIMEZONES: { value: string; label: string }[] = [
  { value: 'Asia/Kolkata', label: 'India Standard Time (IST)' },
  { value: 'Asia/Dubai', label: 'Gulf Standard Time (GST)' },
  { value: 'Asia/Singapore', label: 'Singapore Time (SGT)' },
  { value: 'Asia/Hong_Kong', label: 'Hong Kong Time (HKT)' },
  { value: 'Asia/Tokyo', label: 'Japan Standard Time (JST)' },
  { value: 'Europe/London', label: 'Greenwich Mean Time (GMT)' },
  { value: 'Europe/Paris', label: 'Central European Time (CET)' },
  { value: 'America/New_York', label: 'Eastern Time (ET)' },
  { value: 'America/Chicago', label: 'Central Time (CT)' },
  { value: 'America/Los_Angeles', label: 'Pacific Time (PT)' },
  { value: 'Australia/Sydney', label: 'Australian Eastern Time (AET)' },
];
