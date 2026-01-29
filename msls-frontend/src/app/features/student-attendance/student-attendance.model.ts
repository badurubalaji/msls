/**
 * Student Attendance Models
 */

export type StudentAttendanceStatus = 'present' | 'absent' | 'late' | 'half_day';

export interface TeacherClass {
  sectionId: string;
  sectionName: string;
  sectionCode: string;
  className: string;
  classCode: string;
  studentCount: number;
  isMarkedToday: boolean;
  markedCount?: number;
}

export interface StudentForAttendance {
  studentId: string;
  admissionNumber: string;
  rollNumber?: number;
  firstName: string;
  lastName: string;
  fullName: string;
  photoUrl?: string;
  status?: StudentAttendanceStatus;
  lateArrivalTime?: string;
  remarks?: string;
  last5Days: string[];
}

export interface AttendanceSummary {
  total: number;
  present: number;
  absent: number;
  late: number;
  halfDay: number;
}

export interface ClassAttendance {
  sectionId: string;
  sectionName: string;
  className: string;
  date: string;
  students: StudentForAttendance[];
  isMarked: boolean;
  canEdit: boolean;
  markedAt?: string;
  markedBy?: string;
  markedByName?: string;
  summary: AttendanceSummary;
}

export interface StudentAttendanceRecord {
  id: string;
  studentId: string;
  studentName?: string;
  admissionNumber?: string;
  sectionId: string;
  sectionName?: string;
  attendanceDate: string;
  status: StudentAttendanceStatus;
  statusLabel: string;
  lateArrivalTime?: string;
  remarks?: string;
  markedBy: string;
  markedByName?: string;
  markedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface StudentAttendanceListResponse {
  attendance: StudentAttendanceRecord[];
  nextCursor?: string;
  hasMore: boolean;
  total?: number;
}

export interface MarkAttendanceResult {
  sectionId: string;
  date: string;
  summary: AttendanceSummary;
  markedAt: string;
  message: string;
}

export interface StudentAttendanceSettings {
  id?: string;
  branchId: string;
  branchName?: string;
  editWindowMinutes: number;
  lateThresholdMinutes: number;
  smsOnAbsent: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// Request DTOs
export interface AttendanceRecordRequest {
  studentId: string;
  status: StudentAttendanceStatus;
  lateArrivalTime?: string;
  remarks?: string;
}

export interface MarkClassAttendanceRequest {
  date: string;
  records: AttendanceRecordRequest[];
}

export interface UpdateSettingsRequest {
  branchId: string;
  editWindowMinutes: number;
  lateThresholdMinutes: number;
  smsOnAbsent: boolean;
}

// Filter interfaces
export interface StudentAttendanceFilter {
  sectionId?: string;
  studentId?: string;
  status?: StudentAttendanceStatus;
  dateFrom?: string;
  dateTo?: string;
  cursor?: string;
  limit?: number;
}

// Status label mappings
export const STUDENT_ATTENDANCE_STATUS_LABELS: Record<StudentAttendanceStatus, string> = {
  present: 'Present',
  absent: 'Absent',
  late: 'Late',
  half_day: 'Half Day',
};

export const STUDENT_ATTENDANCE_STATUS_SHORT: Record<StudentAttendanceStatus, string> = {
  present: 'P',
  absent: 'A',
  late: 'L',
  half_day: 'H',
};

export const STUDENT_ATTENDANCE_STATUS_COLORS: Record<StudentAttendanceStatus, string> = {
  present: 'text-green-600 bg-green-100',
  absent: 'text-red-600 bg-red-100',
  late: 'text-yellow-600 bg-yellow-100',
  half_day: 'text-blue-600 bg-blue-100',
};

export const STUDENT_ATTENDANCE_DOT_COLORS: Record<string, string> = {
  P: 'bg-green-500',
  A: 'bg-red-500',
  L: 'bg-yellow-500',
  H: 'bg-blue-500',
};
