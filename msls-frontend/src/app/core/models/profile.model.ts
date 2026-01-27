/**
 * MSLS Profile Models
 *
 * Defines profile-related interfaces for user profile management.
 */

/**
 * Notification preferences for the user
 */
export interface NotificationPreferences {
  /** Email notifications enabled */
  email: boolean;

  /** Push notifications enabled */
  push: boolean;

  /** SMS notifications enabled */
  sms: boolean;
}

/**
 * User profile with extended information
 */
export interface UserProfile {
  /** Unique user identifier */
  id: string;

  /** Tenant ID for multi-tenancy support */
  tenantId: string;

  /** User's email address */
  email?: string;

  /** User's phone number */
  phone?: string;

  /** User's first name */
  firstName: string;

  /** User's last name */
  lastName: string;

  /** User's full name (computed) */
  fullName?: string;

  /** URL to user's avatar image */
  avatarUrl?: string;

  /** User's biography or description */
  bio?: string;

  /** User's preferred timezone (IANA format) */
  timezone: string;

  /** User's preferred locale/language code */
  locale: string;

  /** Notification preferences */
  notificationPreferences?: NotificationPreferences;

  /** Account status */
  status: 'active' | 'inactive' | 'pending' | 'suspended';

  /** Whether email is verified */
  emailVerifiedAt?: string;

  /** Whether phone is verified */
  phoneVerifiedAt?: string;

  /** Last login timestamp */
  lastLoginAt?: string;

  /** Whether 2FA is enabled */
  twoFactorEnabled: boolean;

  /** Account creation timestamp */
  createdAt: string;

  /** User's assigned roles */
  roles?: ProfileRole[];

  /** User's permissions */
  permissions?: string[];
}

/**
 * Role information in profile response
 */
export interface ProfileRole {
  /** Role ID */
  id: string;

  /** Role name */
  name: string;

  /** Role description */
  description?: string;

  /** Whether this is a system role */
  isSystem: boolean;
}

/**
 * Profile update request
 */
export interface UpdateProfileRequest {
  /** First name */
  firstName?: string;

  /** Last name */
  lastName?: string;

  /** Phone number */
  phone?: string;

  /** Bio/description */
  bio?: string;

  /** Timezone (IANA format) */
  timezone?: string;

  /** Locale code */
  locale?: string;
}

/**
 * Password change request
 */
export interface ChangePasswordRequest {
  /** Current password */
  currentPassword: string;

  /** New password */
  newPassword: string;

  /** Confirm new password */
  confirmPassword: string;
}

/**
 * Notification preferences update request
 */
export interface UpdatePreferencesRequest {
  /** Email notifications enabled */
  email?: boolean;

  /** Push notifications enabled */
  push?: boolean;

  /** SMS notifications enabled */
  sms?: boolean;
}

/**
 * Avatar upload response
 */
export interface AvatarUploadResponse {
  /** URL to the uploaded avatar */
  avatarUrl: string;
}

/**
 * Extended user preference
 */
export interface UserPreference {
  /** Preference category */
  category: string;

  /** Preference key */
  key: string;

  /** Preference value */
  value: unknown;
}

/**
 * Extended preferences response
 */
export interface UserPreferencesResponse {
  /** List of preferences */
  preferences: UserPreference[];
}

/**
 * Set preference request
 */
export interface SetPreferenceRequest {
  /** Preference category */
  category: string;

  /** Preference key */
  key: string;

  /** Preference value */
  value: unknown;
}

/**
 * Available timezones for selection
 */
export const AVAILABLE_TIMEZONES = [
  { value: 'UTC', label: 'UTC' },
  { value: 'America/New_York', label: 'Eastern Time (US & Canada)' },
  { value: 'America/Chicago', label: 'Central Time (US & Canada)' },
  { value: 'America/Denver', label: 'Mountain Time (US & Canada)' },
  { value: 'America/Los_Angeles', label: 'Pacific Time (US & Canada)' },
  { value: 'Europe/London', label: 'London' },
  { value: 'Europe/Paris', label: 'Paris' },
  { value: 'Europe/Berlin', label: 'Berlin' },
  { value: 'Asia/Tokyo', label: 'Tokyo' },
  { value: 'Asia/Shanghai', label: 'Shanghai' },
  { value: 'Asia/Kolkata', label: 'Mumbai, Kolkata' },
  { value: 'Asia/Dubai', label: 'Dubai' },
  { value: 'Australia/Sydney', label: 'Sydney' },
  { value: 'Pacific/Auckland', label: 'Auckland' },
] as const;

/**
 * Available locales for selection
 */
export const AVAILABLE_LOCALES = [
  { value: 'en', label: 'English' },
  { value: 'en-US', label: 'English (US)' },
  { value: 'en-GB', label: 'English (UK)' },
  { value: 'es', label: 'Espanol' },
  { value: 'fr', label: 'Francais' },
  { value: 'de', label: 'Deutsch' },
  { value: 'ja', label: 'Japanese' },
  { value: 'zh', label: 'Chinese' },
  { value: 'hi', label: 'Hindi' },
  { value: 'ar', label: 'Arabic' },
] as const;
