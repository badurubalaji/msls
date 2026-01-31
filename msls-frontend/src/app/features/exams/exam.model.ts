/**
 * Exam module models for Exam Types and related entities.
 */

// ========================================
// Evaluation Type
// ========================================

export type EvaluationType = 'marks' | 'grade';

export const EVALUATION_TYPES: { value: EvaluationType; label: string; icon: string; color: string }[] = [
  { value: 'marks', label: 'Marks Based', icon: 'fa-solid fa-calculator', color: 'bg-blue-100 text-blue-700' },
  { value: 'grade', label: 'Grade Based', icon: 'fa-solid fa-star', color: 'bg-amber-100 text-amber-700' },
];

// ========================================
// Exam Type Models
// ========================================

export interface ExamType {
  id: string;
  name: string;
  code: string;
  description?: string;
  weightage: number;
  evaluationType: EvaluationType;
  defaultMaxMarks: number;
  defaultPassingMarks?: number;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateExamTypeRequest {
  name: string;
  code: string;
  description?: string;
  weightage: number;
  evaluationType: EvaluationType;
  defaultMaxMarks: number;
  defaultPassingMarks?: number;
}

export interface UpdateExamTypeRequest {
  name?: string;
  code?: string;
  description?: string;
  weightage?: number;
  evaluationType?: EvaluationType;
  defaultMaxMarks?: number;
  defaultPassingMarks?: number;
}

export interface ExamTypeListResponse {
  items: ExamType[];
  total: number;
}

export interface ExamTypeFilter {
  isActive?: boolean;
  search?: string;
}

export interface DisplayOrderItem {
  id: string;
  displayOrder: number;
}

export interface UpdateDisplayOrderRequest {
  items: DisplayOrderItem[];
}

export interface ToggleActiveRequest {
  isActive: boolean;
}

// ========================================
// Examination Models
// ========================================

export type ExamStatus = 'draft' | 'scheduled' | 'ongoing' | 'completed' | 'cancelled';

export const EXAM_STATUSES: { value: ExamStatus; label: string; color: string; icon: string }[] = [
  { value: 'draft', label: 'Draft', color: 'bg-gray-100 text-gray-700', icon: 'fa-solid fa-file-pen' },
  { value: 'scheduled', label: 'Scheduled', color: 'bg-blue-100 text-blue-700', icon: 'fa-solid fa-calendar-check' },
  { value: 'ongoing', label: 'Ongoing', color: 'bg-amber-100 text-amber-700', icon: 'fa-solid fa-spinner' },
  { value: 'completed', label: 'Completed', color: 'bg-green-100 text-green-700', icon: 'fa-solid fa-check-circle' },
  { value: 'cancelled', label: 'Cancelled', color: 'bg-red-100 text-red-700', icon: 'fa-solid fa-times-circle' },
];

export interface ClassSummary {
  id: string;
  name: string;
}

export interface AcademicYearSummary {
  id: string;
  name: string;
  isCurrent: boolean;
}

export interface SubjectSummary {
  id: string;
  name: string;
  code: string;
}

export interface ExamSchedule {
  id: string;
  subjectId: string;
  subjectName?: string;
  subjectCode?: string;
  examDate: string;
  startTime: string;
  endTime: string;
  maxMarks: number;
  passingMarks?: number;
  venue?: string;
  notes?: string;
}

export interface Examination {
  id: string;
  name: string;
  examTypeId: string;
  examTypeName?: string;
  academicYearId: string;
  academicYear?: string;
  startDate: string;
  endDate: string;
  status: ExamStatus;
  description?: string;
  classes: ClassSummary[];
  schedules?: ExamSchedule[];
  scheduleCount: number;
  createdAt: string;
  updatedAt: string;
  // For convenience in UI
  examType?: { id: string; name: string };
}

export interface CreateExaminationRequest {
  name: string;
  examTypeId: string;
  academicYearId: string;
  startDate: string;
  endDate: string;
  description?: string;
  classIds: string[];
}

export interface UpdateExaminationRequest {
  name?: string;
  examTypeId?: string;
  academicYearId?: string;
  startDate?: string;
  endDate?: string;
  description?: string;
  classIds?: string[];
}

export interface ExaminationFilter {
  academicYearId?: string;
  examTypeId?: string;
  classId?: string;
  status?: ExamStatus;
  search?: string;
}

export interface CreateScheduleRequest {
  subjectId: string;
  examDate: string;
  startTime: string;
  endTime: string;
  maxMarks: number;
  passingMarks?: number;
  venue?: string;
  notes?: string;
}

export interface UpdateScheduleRequest {
  subjectId?: string;
  examDate?: string;
  startTime?: string;
  endTime?: string;
  maxMarks?: number;
  passingMarks?: number;
  venue?: string;
  notes?: string;
}

// ========================================
// Common Exam Constants
// ========================================

export const DEFAULT_MAX_MARKS = 100;
export const DEFAULT_PASSING_PERCENTAGE = 35;

// Common exam type presets for quick setup
export const EXAM_TYPE_PRESETS: Partial<CreateExamTypeRequest>[] = [
  { name: 'Unit Test', code: 'UT', weightage: 10, evaluationType: 'marks', defaultMaxMarks: 25 },
  { name: 'Mid Term', code: 'MT', weightage: 20, evaluationType: 'marks', defaultMaxMarks: 50 },
  { name: 'Final Term', code: 'FT', weightage: 30, evaluationType: 'marks', defaultMaxMarks: 100 },
  { name: 'Half Yearly', code: 'HY', weightage: 40, evaluationType: 'marks', defaultMaxMarks: 100 },
  { name: 'Annual', code: 'AN', weightage: 60, evaluationType: 'marks', defaultMaxMarks: 100 },
  { name: 'Practical', code: 'PR', weightage: 20, evaluationType: 'marks', defaultMaxMarks: 30 },
  { name: 'Internal Assessment', code: 'IA', weightage: 10, evaluationType: 'grade', defaultMaxMarks: 10 },
];

// ========================================
// Hall Ticket Models
// ========================================

export type HallTicketStatus = 'generated' | 'printed' | 'downloaded';

export const HALL_TICKET_STATUSES: { value: HallTicketStatus; label: string; color: string; icon: string }[] = [
  { value: 'generated', label: 'Generated', color: 'bg-blue-100 text-blue-700', icon: 'fa-solid fa-file-circle-check' },
  { value: 'printed', label: 'Printed', color: 'bg-amber-100 text-amber-700', icon: 'fa-solid fa-print' },
  { value: 'downloaded', label: 'Downloaded', color: 'bg-green-100 text-green-700', icon: 'fa-solid fa-download' },
];

export interface HallTicket {
  id: string;
  examinationId: string;
  studentId: string;
  rollNumber: string;
  qrCodeData: string;
  status: HallTicketStatus;
  generatedAt: string;
  printedAt?: string;
  downloadedAt?: string;
  // Joined fields
  studentName?: string;
  admissionNumber?: string;
  className?: string;
  sectionName?: string;
  examinationName?: string;
}

export interface HallTicketTemplate {
  id: string;
  name: string;
  headerLogoUrl?: string;
  schoolName?: string;
  schoolAddress?: string;
  instructions?: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface HallTicketListResponse {
  data: HallTicket[];
  meta: {
    total: number;
    limit: number;
    offset: number;
  };
}

export interface HallTicketFilter {
  classId?: string;
  sectionId?: string;
  status?: HallTicketStatus;
  search?: string;
  limit?: number;
  offset?: number;
}

export interface GenerateHallTicketsRequest {
  classId?: string;
  sectionId?: string;
  rollNumberPrefix?: string;
}

export interface GenerateHallTicketsResponse {
  totalStudents: number;
  generated: number;
  skipped: number;
  failed: number;
  errors?: string[];
}

export interface CreateHallTicketTemplateRequest {
  name: string;
  headerLogoUrl?: string;
  schoolName?: string;
  schoolAddress?: string;
  instructions?: string;
  isDefault?: boolean;
}

export interface UpdateHallTicketTemplateRequest {
  name?: string;
  headerLogoUrl?: string;
  schoolName?: string;
  schoolAddress?: string;
  instructions?: string;
  isDefault?: boolean;
}

export interface VerifyHallTicketResponse {
  valid: boolean;
  hallTicketId?: string;
  studentName?: string;
  admissionNumber?: string;
  rollNumber?: string;
  examinationName?: string;
  className?: string;
  message: string;
}
