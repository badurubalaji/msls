/**
 * MSLS Admission Reports & Analytics Models
 *
 * Defines interfaces for admission reports including dashboard stats,
 * conversion funnel, class-wise reports, and filter parameters.
 */

/**
 * Dashboard statistics summary
 */
export interface DashboardStats {
  /** Total enquiries received */
  totalEnquiries: number;

  /** Total applications submitted */
  totalApplications: number;

  /** Number of approved applications */
  approved: number;

  /** Number of enrolled students */
  enrolled: number;

  /** Number of pending applications */
  pending: number;

  /** Number of rejected applications */
  rejected: number;

  /** Conversion rates between stages */
  conversionRates: ConversionRates;
}

/**
 * Conversion rates between different admission stages
 */
export interface ConversionRates {
  /** Percentage of enquiries that became applications */
  enquiryToApplication: number;

  /** Percentage of applications that were approved */
  applicationToApproved: number;

  /** Percentage of approved applications that enrolled */
  approvedToEnrolled: number;
}

/**
 * Funnel stage for visualization
 */
export interface FunnelStage {
  /** Stage name */
  name: string;

  /** Count at this stage */
  count: number;

  /** Percentage of total (for progress bar width) */
  percentage: number;

  /** Color variant for the stage */
  variant: 'info' | 'primary' | 'warning' | 'success';

  /** Icon class for the stage */
  icon: string;
}

/**
 * Class-wise admission report data
 */
export interface ClassWiseReport {
  /** Class name (e.g., LKG, UKG, Class 1) */
  className: string;

  /** Total seats available for the class */
  totalSeats: number;

  /** Number of applications received */
  applications: number;

  /** Number of approved applications */
  approved: number;

  /** Number of students enrolled */
  enrolled: number;

  /** Number of students on waitlist */
  waitlisted: number;

  /** Number of vacant seats */
  vacant: number;

  /** Fill percentage */
  fillPercentage: number;
}

/**
 * Class-wise report response from API
 */
export interface ClassWiseReportResponse {
  /** Array of class-wise data */
  classes: ClassWiseReport[];

  /** Session information */
  sessionName?: string;

  /** Report generation timestamp */
  generatedAt?: string;
}

/**
 * Source analysis data
 */
export interface SourceAnalysis {
  /** Source name */
  source: string;

  /** Source display label */
  label: string;

  /** Count from this source */
  count: number;

  /** Percentage of total */
  percentage: number;
}

/**
 * Daily trend data point
 */
export interface DailyTrendPoint {
  /** Date string (YYYY-MM-DD) */
  date: string;

  /** Number of enquiries on this date */
  enquiries: number;

  /** Number of applications on this date */
  applications: number;

  /** Number of enrollments on this date */
  enrollments: number;
}

/**
 * Filter parameters for admission reports
 */
export interface ReportFilterParams {
  /** Session ID to filter by */
  sessionId?: string;

  /** Branch ID to filter by */
  branchId?: string;

  /** Start date for date range filter */
  fromDate?: string;

  /** End date for date range filter */
  toDate?: string;

  /** Class to filter by */
  className?: string;
}

/**
 * Summary card configuration
 */
export interface SummaryCard {
  /** Card title */
  title: string;

  /** Main value to display */
  value: number;

  /** Icon class (Font Awesome) */
  icon: string;

  /** Background color class */
  bgColor: string;

  /** Icon color class */
  iconColor: string;

  /** Text color class */
  textColor: string;

  /** Optional subtitle (e.g., percentage change) */
  subtitle?: string;

  /** Trend indicator */
  trend?: 'up' | 'down' | 'neutral';

  /** Trend value (e.g., "+12%") */
  trendValue?: string;
}

/**
 * Export format options
 */
export type ExportFormat = 'excel' | 'pdf';

/**
 * Export request parameters
 */
export interface ExportRequest {
  /** Report type to export */
  reportType: 'dashboard' | 'class-wise' | 'funnel' | 'source-analysis';

  /** Export format */
  format: ExportFormat;

  /** Filter parameters */
  filters?: ReportFilterParams;
}

/**
 * Admission session for filter dropdown
 */
export interface AdmissionSessionOption {
  id: string;
  name: string;
  status: 'upcoming' | 'open' | 'closed';
}

/**
 * Get status badge configuration for fill percentage
 */
export function getFillStatusConfig(fillPercentage: number): {
  label: string;
  variant: 'success' | 'warning' | 'danger' | 'info';
} {
  if (fillPercentage >= 90) {
    return { label: 'Almost Full', variant: 'danger' };
  } else if (fillPercentage >= 70) {
    return { label: 'Filling Fast', variant: 'warning' };
  } else if (fillPercentage >= 40) {
    return { label: 'Available', variant: 'info' };
  } else {
    return { label: 'Open', variant: 'success' };
  }
}

/**
 * Default funnel stages configuration
 */
export const DEFAULT_FUNNEL_STAGES: Omit<FunnelStage, 'count' | 'percentage'>[] = [
  { name: 'Enquiries', variant: 'info', icon: 'fa-solid fa-question-circle' },
  { name: 'Applications', variant: 'primary', icon: 'fa-solid fa-file-alt' },
  { name: 'Approved', variant: 'warning', icon: 'fa-solid fa-check-circle' },
  { name: 'Enrolled', variant: 'success', icon: 'fa-solid fa-user-graduate' },
];

/**
 * Summary card configurations
 */
export const SUMMARY_CARD_CONFIGS: Omit<SummaryCard, 'value'>[] = [
  {
    title: 'Total Enquiries',
    icon: 'fa-solid fa-question-circle',
    bgColor: 'bg-blue-50',
    iconColor: 'text-blue-500',
    textColor: 'text-blue-700',
  },
  {
    title: 'Applications',
    icon: 'fa-solid fa-file-alt',
    bgColor: 'bg-purple-50',
    iconColor: 'text-purple-500',
    textColor: 'text-purple-700',
  },
  {
    title: 'Approved',
    icon: 'fa-solid fa-check-circle',
    bgColor: 'bg-amber-50',
    iconColor: 'text-amber-500',
    textColor: 'text-amber-700',
  },
  {
    title: 'Enrolled',
    icon: 'fa-solid fa-user-graduate',
    bgColor: 'bg-emerald-50',
    iconColor: 'text-emerald-500',
    textColor: 'text-emerald-700',
  },
];
