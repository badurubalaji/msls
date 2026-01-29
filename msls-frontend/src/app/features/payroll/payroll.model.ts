/**
 * Payroll model interfaces for the frontend
 * Story 5.6: Payroll Processing
 */

// Pay Run Status
export type PayRunStatus = 'draft' | 'processing' | 'calculated' | 'approved' | 'finalized' | 'reversed';

// Payslip Status
export type PayslipStatus = 'calculated' | 'adjusted' | 'approved' | 'paid';

/**
 * Pay Run interface representing a monthly payroll batch
 */
export interface PayRun {
  id: string;
  payPeriodMonth: number;
  payPeriodYear: number;
  branchId?: string;
  branchName?: string;
  status: PayRunStatus;
  totalStaff: number;
  totalGross: string;
  totalDeductions: string;
  totalNet: string;
  calculatedAt?: string;
  approvedAt?: string;
  approvedBy?: string;
  approvedByName?: string;
  finalizedAt?: string;
  finalizedBy?: string;
  finalizedByName?: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
  createdBy?: string;
}

/**
 * Payslip interface representing individual staff pay record
 */
export interface Payslip {
  id: string;
  payRunId: string;
  staffId: string;
  staffName?: string;
  staffEmployeeId?: string;
  staffSalaryId?: string;
  workingDays: number;
  presentDays: number;
  leaveDays: number;
  absentDays: number;
  lopDays: number;
  grossSalary: string;
  totalEarnings: string;
  totalDeductions: string;
  netSalary: string;
  lopDeduction: string;
  status: PayslipStatus;
  paymentDate?: string;
  paymentReference?: string;
  components?: PayslipComponent[];
  createdAt: string;
  updatedAt: string;
}

/**
 * Payslip component breakdown
 */
export interface PayslipComponent {
  id: string;
  componentId: string;
  componentName: string;
  componentCode: string;
  componentType: 'earning' | 'deduction';
  amount: string;
  isProrated: boolean;
}

/**
 * Pay Run list response
 */
export interface PayRunListResponse {
  payRuns: PayRun[];
  total: number;
}

/**
 * Payslip list response
 */
export interface PayslipListResponse {
  payslips: Payslip[];
  total: number;
}

/**
 * Create pay run request DTO
 */
export interface CreatePayRunRequest {
  payPeriodMonth: number;
  payPeriodYear: number;
  branchId?: string;
  notes?: string;
}

/**
 * Adjust payslip request DTO
 */
export interface AdjustPayslipRequest {
  components: PayslipComponentInput[];
  notes?: string;
}

/**
 * Payslip component input for adjustment
 */
export interface PayslipComponentInput {
  componentId: string;
  amount: string;
}

/**
 * Pay run summary by department
 */
export interface PayRunSummary {
  totalStaff: number;
  totalGross: string;
  totalDeductions: string;
  totalNet: string;
  departmentSummaries: DepartmentSummary[];
}

/**
 * Department-wise summary
 */
export interface DepartmentSummary {
  departmentId: string;
  departmentName: string;
  staffCount: number;
  totalGross: string;
  totalDeductions: string;
  totalNet: string;
}

/**
 * Get status badge variant based on pay run status
 */
export function getPayRunStatusVariant(status: PayRunStatus): 'success' | 'warning' | 'danger' | 'neutral' | 'info' {
  switch (status) {
    case 'draft':
      return 'neutral';
    case 'processing':
      return 'info';
    case 'calculated':
      return 'warning';
    case 'approved':
      return 'info';
    case 'finalized':
      return 'success';
    case 'reversed':
      return 'danger';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for pay run status
 */
export function getPayRunStatusLabel(status: PayRunStatus): string {
  switch (status) {
    case 'draft':
      return 'Draft';
    case 'processing':
      return 'Processing';
    case 'calculated':
      return 'Calculated';
    case 'approved':
      return 'Approved';
    case 'finalized':
      return 'Finalized';
    case 'reversed':
      return 'Reversed';
    default:
      return status;
  }
}

/**
 * Get status badge variant based on payslip status
 */
export function getPayslipStatusVariant(status: PayslipStatus): 'success' | 'warning' | 'danger' | 'neutral' | 'info' {
  switch (status) {
    case 'calculated':
      return 'neutral';
    case 'adjusted':
      return 'warning';
    case 'approved':
      return 'info';
    case 'paid':
      return 'success';
    default:
      return 'neutral';
  }
}

/**
 * Get display label for payslip status
 */
export function getPayslipStatusLabel(status: PayslipStatus): string {
  switch (status) {
    case 'calculated':
      return 'Calculated';
    case 'adjusted':
      return 'Adjusted';
    case 'approved':
      return 'Approved';
    case 'paid':
      return 'Paid';
    default:
      return status;
  }
}

/**
 * Get month name from month number
 */
export function getMonthName(month: number): string {
  const months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];
  return months[month - 1] || '';
}

/**
 * Format pay period display
 */
export function formatPayPeriod(month: number, year: number): string {
  return `${getMonthName(month)} ${year}`;
}

/**
 * Format currency for display
 */
export function formatCurrency(amount: string | number): string {
  const num = typeof amount === 'string' ? parseFloat(amount || '0') : amount;
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(num);
}
