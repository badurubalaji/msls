/**
 * MSLS Admission Enquiry Models
 *
 * Defines interfaces for admission enquiry management including
 * enquiries, follow-ups, and related DTOs.
 */

/**
 * Enquiry status enum
 */
export type EnquiryStatus = 'new' | 'contacted' | 'interested' | 'converted' | 'closed';

/**
 * Enquiry source enum
 */
export type EnquirySource = 'walk_in' | 'phone' | 'website' | 'referral' | 'advertisement' | 'social_media' | 'other';

/**
 * Contact mode for follow-ups
 */
export type ContactMode = 'phone' | 'email' | 'whatsapp' | 'in_person' | 'other';

/**
 * Follow-up outcome
 */
export type FollowUpOutcome = 'interested' | 'not_interested' | 'follow_up_required' | 'converted' | 'no_response';

/**
 * Admission enquiry interface
 */
export interface Enquiry {
  /** Unique identifier (UUID v7) */
  id: string;

  /** Tenant identifier */
  tenantId: string;

  /** Branch identifier */
  branchId?: string;

  /** Admission session identifier */
  sessionId?: string;

  /** Auto-generated enquiry number (e.g., ENQ-2026-001) */
  enquiryNumber: string;

  /** Student's name */
  studentName: string;

  /** Student's date of birth */
  dateOfBirth?: string;

  /** Student's gender */
  gender?: string;

  /** Class the student is applying for */
  classApplying: string;

  /** Parent/guardian name */
  parentName: string;

  /** Parent/guardian phone number */
  parentPhone: string;

  /** Parent/guardian email */
  parentEmail?: string;

  /** Source of enquiry */
  source: EnquirySource;

  /** Referral details if source is referral */
  referralDetails?: string;

  /** Remarks/notes from conversation */
  remarks?: string;

  /** Current status of enquiry */
  status: EnquiryStatus;

  /** Next follow-up date */
  followUpDate?: string;

  /** User assigned to handle this enquiry */
  assignedTo?: string;

  /** ID of application if converted */
  convertedApplicationId?: string;

  /** Record timestamps */
  createdAt: string;
  updatedAt: string;
  createdBy?: string;
}

/**
 * Enquiry follow-up record
 */
export interface EnquiryFollowUp {
  /** Unique identifier */
  id: string;

  /** Tenant identifier */
  tenantId: string;

  /** Parent enquiry identifier */
  enquiryId: string;

  /** Date of follow-up */
  followUpDate: string;

  /** Mode of contact */
  contactMode: ContactMode;

  /** Notes from the follow-up */
  notes?: string;

  /** Outcome of the follow-up */
  outcome?: FollowUpOutcome;

  /** Next scheduled follow-up date */
  nextFollowUp?: string;

  /** Record timestamps */
  createdAt: string;
  createdBy?: string;
}

/**
 * DTO for creating a new enquiry
 */
export interface CreateEnquiryDto {
  studentName: string;
  dateOfBirth?: string;
  gender?: string;
  classApplying: string;
  parentName: string;
  parentPhone: string;
  parentEmail?: string;
  source?: EnquirySource;
  referralDetails?: string;
  remarks?: string;
  followUpDate?: string;
  sessionId?: string;
  branchId?: string;
  assignedTo?: string;
}

/**
 * DTO for updating an enquiry
 */
export interface UpdateEnquiryDto extends Partial<CreateEnquiryDto> {
  status?: EnquiryStatus;
}

/**
 * DTO for creating a follow-up
 */
export interface CreateFollowUpDto {
  followUpDate: string;
  contactMode: ContactMode;
  notes?: string;
  outcome?: FollowUpOutcome;
  nextFollowUp?: string;
}

/**
 * DTO for converting enquiry to application
 */
export interface ConvertEnquiryDto {
  sessionId?: string;
  branchId?: string;
  remarks?: string;
}

/**
 * Filter parameters for listing enquiries
 */
export interface EnquiryFilterParams {
  status?: EnquiryStatus;
  classApplying?: string;
  source?: EnquirySource;
  fromDate?: string;
  toDate?: string;
  search?: string;
  assignedTo?: string;
  branchId?: string;
  sessionId?: string;
}

/**
 * Status configuration for display
 */
export interface StatusConfig {
  label: string;
  variant: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  icon: string;
}

/**
 * Status display configuration map
 */
export const ENQUIRY_STATUS_CONFIG: Record<EnquiryStatus, StatusConfig> = {
  new: {
    label: 'New',
    variant: 'info',
    icon: 'fa-solid fa-sparkles',
  },
  contacted: {
    label: 'Contacted',
    variant: 'primary',
    icon: 'fa-solid fa-phone',
  },
  interested: {
    label: 'Interested',
    variant: 'success',
    icon: 'fa-solid fa-thumbs-up',
  },
  converted: {
    label: 'Converted',
    variant: 'success',
    icon: 'fa-solid fa-check-circle',
  },
  closed: {
    label: 'Closed',
    variant: 'neutral',
    icon: 'fa-solid fa-times-circle',
  },
};

/**
 * Source display labels
 */
export const ENQUIRY_SOURCE_LABELS: Record<EnquirySource, string> = {
  walk_in: 'Walk-in',
  phone: 'Phone Call',
  website: 'Website',
  referral: 'Referral',
  advertisement: 'Advertisement',
  social_media: 'Social Media',
  other: 'Other',
};

/**
 * Contact mode labels
 */
export const CONTACT_MODE_LABELS: Record<ContactMode, string> = {
  phone: 'Phone',
  email: 'Email',
  whatsapp: 'WhatsApp',
  in_person: 'In Person',
  other: 'Other',
};

/**
 * Follow-up outcome labels
 */
export const FOLLOW_UP_OUTCOME_LABELS: Record<FollowUpOutcome, string> = {
  interested: 'Interested',
  not_interested: 'Not Interested',
  follow_up_required: 'Follow-up Required',
  converted: 'Converted',
  no_response: 'No Response',
};
