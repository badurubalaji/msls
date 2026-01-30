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

// ============================================================================
// Period-wise Attendance Models (Story 7.2)
// ============================================================================

export interface PeriodInfo {
  periodSlotId: string;
  timetableEntryId: string;
  periodName: string;
  periodNumber?: number;
  startTime: string;
  endTime: string;
  subjectId?: string;
  subjectName?: string;
  subjectCode?: string;
  staffId?: string;
  staffName?: string;
  isMarked: boolean;
  markedCount?: number;
  totalStudents?: number;
}

export interface TeacherPeriodsResponse {
  sectionId: string;
  sectionName: string;
  className: string;
  date: string;
  dayOfWeek: number;
  dayName: string;
  periods: PeriodInfo[];
}

export interface StudentForPeriodAttendance {
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
}

export interface PeriodAttendance {
  sectionId: string;
  sectionName: string;
  className: string;
  date: string;
  periodSlotId: string;
  periodName: string;
  periodNumber?: number;
  startTime: string;
  endTime: string;
  subjectId?: string;
  subjectName?: string;
  subjectCode?: string;
  timetableEntryId?: string;
  students: StudentForPeriodAttendance[];
  isMarked: boolean;
  canEdit: boolean;
  markedAt?: string;
  markedByName?: string;
  summary: AttendanceSummary;
}

export interface MarkPeriodAttendanceResult {
  sectionId: string;
  periodId: string;
  date: string;
  summary: AttendanceSummary;
  markedAt: string;
  message: string;
}

export interface DailySummaryStudent {
  studentId: string;
  admissionNumber: string;
  rollNumber?: number;
  fullName: string;
  photoUrl?: string;
  periodStatuses: Record<string, StudentAttendanceStatus>;
  totalPeriods: number;
  periodsPresent: number;
  periodsAbsent: number;
  periodsLate: number;
  attendancePercentage: number;
  overallStatus: StudentAttendanceStatus;
}

export interface DailySummarySummary {
  totalStudents: number;
  totalPeriods: number;
  averageAttendance: number;
  fullPresentCount: number;
  absentCount: number;
}

export interface DailySummary {
  sectionId: string;
  sectionName: string;
  className: string;
  date: string;
  dayName: string;
  periods: PeriodInfo[];
  students: DailySummaryStudent[];
  summary: DailySummarySummary;
}

export interface SubjectAttendanceStats {
  studentId: string;
  studentName?: string;
  subjectId: string;
  subjectName: string;
  subjectCode: string;
  totalPeriods: number;
  periodsPresent: number;
  periodsAbsent: number;
  periodsLate: number;
  attendancePercentage: number;
  minimumRequired: number;
  isEligible: boolean;
}

// Request DTOs for period attendance
export interface PeriodAttendanceRecordRequest {
  studentId: string;
  status: StudentAttendanceStatus;
  lateArrivalTime?: string;
  remarks?: string;
}

export interface MarkPeriodAttendanceRequest {
  sectionId: string;
  date: string;
  records: PeriodAttendanceRecordRequest[];
}

// ============================================================================
// Attendance Edit & Audit Models (Story 7.3)
// ============================================================================

export type AttendanceChangeType = 'create' | 'edit';

export interface EditAttendanceRequest {
  status?: StudentAttendanceStatus;
  remarks?: string;
  lateArrivalTime?: string;
  reason: string;
}

export interface EditAttendanceResult {
  attendanceId: string;
  studentId: string;
  date: string;
  status: StudentAttendanceStatus;
  editedAt: string;
  editedBy: string;
  message: string;
}

export interface AttendanceAuditEntry {
  id: string;
  changeType: AttendanceChangeType;
  previousStatus?: StudentAttendanceStatus;
  newStatus: StudentAttendanceStatus;
  previousRemarks?: string;
  newRemarks?: string;
  changeReason: string;
  changedById: string;
  changedByName: string;
  changedAt: string;
}

export interface AttendanceAuditTrail {
  attendanceId: string;
  studentId: string;
  studentName: string;
  date: string;
  auditEntries: AttendanceAuditEntry[];
  totalChanges: number;
}

export interface EditWindowStatus {
  attendanceId: string;
  markedAt: string;
  windowEndAt: string;
  windowMinutes: number;
  remainingMinutes: number;
  isWithinWindow: boolean;
  canEdit: boolean;
  editDeniedReason?: string;
  isOriginalMarker: boolean;
  requiresAdminEdit: boolean;
}

// ============================================================================
// Calendar & Reports Models (Stories 7.4-7.8)
// ============================================================================

export interface CalendarDay {
  date: string;
  dayOfWeek: number;
  status?: StudentAttendanceStatus;
  isHoliday: boolean;
  isWeekend: boolean;
  holidayName?: string;
  remarks?: string;
}

export interface MonthlySummary {
  workingDays: number;
  present: number;
  absent: number;
  late: number;
  halfDay: number;
  holidays: number;
  percentage: number;
}

export interface MonthlyCalendar {
  studentId: string;
  studentName: string;
  year: number;
  month: number;
  monthName: string;
  days: CalendarDay[];
  summary: MonthlySummary;
  classAverage: number;
  trend: 'improving' | 'declining' | 'stable';
}

export interface StudentSummary {
  studentId: string;
  studentName: string;
  sectionId: string;
  sectionName: string;
  dateFrom: string;
  dateTo: string;
  summary: MonthlySummary;
  classAverage: number;
  trend: 'improving' | 'declining' | 'stable';
  trendPercentage?: number;
}

export interface StudentReportEntry {
  studentId: string;
  admissionNumber: string;
  fullName: string;
  rollNumber?: number;
  status: StudentAttendanceStatus;
  statusLabel: string;
  remarks?: string;
}

export interface ClassReport {
  sectionId: string;
  sectionName: string;
  className: string;
  date: string;
  students: StudentReportEntry[];
  summary: AttendanceSummary;
  attendanceRate: number;
}

export interface MonthlyStudentReport {
  studentId: string;
  admissionNumber: string;
  fullName: string;
  rollNumber?: number;
  dailyStatus: Record<string, StudentAttendanceStatus>;
  present: number;
  absent: number;
  late: number;
  percentage: number;
}

export interface ClassMonthlySummary {
  totalStudents: number;
  averageAttendance: number;
  studentsAbove90: number;
  studentsBelow75: number;
  studentsBelow60: number;
}

export interface MonthlyClassReport {
  sectionId: string;
  sectionName: string;
  className: string;
  year: number;
  month: number;
  monthName: string;
  workingDays: number;
  dates: string[];
  students: MonthlyStudentReport[];
  summary: ClassMonthlySummary;
}

export interface LowAttendanceStudent {
  studentId: string;
  admissionNumber: string;
  fullName: string;
  className: string;
  sectionName: string;
  attendanceRate: number;
  daysAbsent: number;
  lastPresent?: string;
  consecutiveAbsent?: number;
}

export interface ClassAttendanceBreakdown {
  className: string;
  sectionName: string;
  sectionId: string;
  totalStudents: number;
  attendanceRate: number;
  belowThreshold?: number;
}

export interface LowAttendanceDashboard {
  dateFrom: string;
  dateTo: string;
  threshold: number;
  criticalThreshold: number;
  totalStudents: number;
  belowThreshold: number;
  chronicAbsentees: number;
  overallAttendanceRate: number;
  students: LowAttendanceStudent[];
  classBreakdown?: ClassAttendanceBreakdown[];
}

export interface UnmarkedClassInfo {
  sectionId: string;
  sectionName: string;
  className: string;
  teacherId?: string;
  teacherName?: string;
  studentCount: number;
  isEscalated: boolean;
}

export interface UnmarkedAttendance {
  date: string;
  deadline: string;
  isPostDeadline: boolean;
  unmarkedClasses: UnmarkedClassInfo[];
  totalClasses: number;
  markedClasses: number;
}

export interface DailyReportSummary {
  date: string;
  overallAttendanceRate: number;
  totalStudents: number;
  totalPresent: number;
  totalAbsent: number;
  classesMarked: number;
  classesTotal: number;
  lowAttendanceClasses: ClassAttendanceBreakdown[];
  generatedAt: string;
}
