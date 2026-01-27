/**
 * Staff Attendance Models
 */

export type AttendanceStatus = 'present' | 'absent' | 'half_day' | 'on_leave' | 'holiday';
export type HalfDayType = 'first_half' | 'second_half';
export type RegularizationStatus = 'pending' | 'approved' | 'rejected';

export interface Attendance {
  id: string;
  staffId: string;
  staffName?: string;
  employeeId?: string;
  attendanceDate: string;
  status: AttendanceStatus;
  checkInTime?: string;
  checkOutTime?: string;
  isLate: boolean;
  lateMinutes: number;
  halfDayType?: HalfDayType;
  remarks?: string;
  markedBy?: string;
  markedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface TodayAttendance {
  status: 'not_marked' | 'checked_in' | 'checked_out';
  attendance?: Attendance;
  canCheckIn: boolean;
  canCheckOut: boolean;
}

export interface AttendanceSummary {
  month: string;
  year: number;
  totalDays: number;
  presentDays: number;
  absentDays: number;
  halfDays: number;
  leaveDays: number;
  holidayDays: number;
  lateDays: number;
  totalLateMinutes: number;
}

export interface AttendanceListResponse {
  attendance: Attendance[];
  nextCursor?: string;
  hasMore: boolean;
  total?: number;
}

export interface Regularization {
  id: string;
  staffId: string;
  staffName?: string;
  employeeId?: string;
  attendanceId?: string;
  requestDate: string;
  requestedStatus: AttendanceStatus;
  reason: string;
  supportingDocumentUrl?: string;
  status: RegularizationStatus;
  reviewedBy?: string;
  reviewedAt?: string;
  rejectionReason?: string;
  createdAt: string;
  updatedAt: string;
}

export interface RegularizationListResponse {
  regularizations: Regularization[];
  nextCursor?: string;
  hasMore: boolean;
  total?: number;
}

export interface AttendanceSettings {
  id?: string;
  branchId: string;
  branchName?: string;
  workStartTime: string;
  workEndTime: string;
  lateThresholdMinutes: number;
  halfDayThresholdHours: number;
  allowSelfCheckout: boolean;
  requireRegularizationApproval: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// Request DTOs
export interface CheckInRequest {
  staffId: string;
  halfDayType?: HalfDayType;
  remarks?: string;
}

export interface CheckOutRequest {
  staffId: string;
  remarks?: string;
}

export interface MarkAttendanceRequest {
  staffId: string;
  attendanceDate: string;
  status: AttendanceStatus;
  checkInTime?: string;
  checkOutTime?: string;
  halfDayType?: HalfDayType;
  remarks?: string;
}

export interface RegularizationRequest {
  staffId: string;
  requestDate: string;
  requestedStatus: AttendanceStatus;
  reason: string;
  supportingDocumentUrl?: string;
}

export interface RegularizationReviewRequest {
  rejectionReason?: string;
}

export interface UpdateSettingsRequest {
  branchId: string;
  workStartTime: string;
  workEndTime: string;
  lateThresholdMinutes: number;
  halfDayThresholdHours: number;
  allowSelfCheckout: boolean;
  requireRegularizationApproval: boolean;
}

// Filter interfaces
export interface AttendanceFilter {
  staffId?: string;
  branchId?: string;
  departmentId?: string;
  status?: AttendanceStatus;
  dateFrom?: string;
  dateTo?: string;
  cursor?: string;
  limit?: number;
}

export interface RegularizationFilter {
  staffId?: string;
  status?: RegularizationStatus;
  dateFrom?: string;
  dateTo?: string;
  cursor?: string;
  limit?: number;
}

// Status label mappings
export const ATTENDANCE_STATUS_LABELS: Record<AttendanceStatus, string> = {
  present: 'Present',
  absent: 'Absent',
  half_day: 'Half Day',
  on_leave: 'On Leave',
  holiday: 'Holiday',
};

export const REGULARIZATION_STATUS_LABELS: Record<RegularizationStatus, string> = {
  pending: 'Pending',
  approved: 'Approved',
  rejected: 'Rejected',
};

export const HALF_DAY_TYPE_LABELS: Record<HalfDayType, string> = {
  first_half: 'First Half',
  second_half: 'Second Half',
};
