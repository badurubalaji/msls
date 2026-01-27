/**
 * MSLS Entrance Test Models
 *
 * TypeScript interfaces for entrance test management including
 * tests, registrations, results, and related DTOs.
 */

/**
 * Test status enum
 */
export type TestStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

/**
 * Registration status enum
 */
export type RegistrationStatus = 'registered' | 'appeared' | 'absent' | 'result_pending' | 'passed' | 'failed';

/**
 * Result verdict enum
 */
export type ResultVerdict = 'pass' | 'fail' | 'pending';

/**
 * Subject marks interface for test results
 */
export interface SubjectMarks {
  subjectName: string;
  maxMarks: number;
  obtainedMarks: number;
}

/**
 * Entrance test entity
 */
export interface EntranceTest {
  /** Unique identifier (UUID v7) */
  id: string;

  /** Tenant identifier */
  tenantId: string;

  /** Admission session identifier */
  sessionId: string;

  /** Test name */
  testName: string;

  /** Date of the test */
  testDate: string;

  /** Start time (HH:mm format) */
  startTime: string;

  /** Duration in minutes */
  durationMinutes: number;

  /** Venue/location of the test */
  venue?: string;

  /** Classes this test is for */
  classNames: string[];

  /** Maximum candidates allowed */
  maxCandidates: number;

  /** Current status */
  status: TestStatus;

  /** Subjects included in the test with max marks */
  subjects: TestSubject[];

  /** Count of registered candidates */
  registeredCount?: number;

  /** Record timestamps */
  createdAt: string;
  updatedAt: string;
}

/**
 * Test subject configuration
 */
export interface TestSubject {
  name: string;
  maxMarks: number;
  passingMarks?: number;
}

/**
 * Test registration entity
 */
export interface TestRegistration {
  /** Unique identifier */
  id: string;

  /** Tenant identifier */
  tenantId: string;

  /** Test identifier */
  testId: string;

  /** Application identifier */
  applicationId: string;

  /** Generated roll number */
  rollNumber?: string;

  /** Registration status */
  status: RegistrationStatus;

  /** Marks obtained per subject */
  marks: SubjectMarks[];

  /** Total marks obtained */
  totalMarks?: number;

  /** Maximum possible marks */
  maxMarks?: number;

  /** Percentage obtained */
  percentage?: number;

  /** Final result */
  result?: ResultVerdict;

  /** Remarks/comments */
  remarks?: string;

  /** Record timestamps */
  createdAt: string;
  updatedAt: string;

  /** Nested application details */
  application?: RegistrationApplication;
}

/**
 * Application details for registration list
 */
export interface RegistrationApplication {
  id: string;
  applicationNumber: string;
  studentName: string;
  parentName: string;
  parentPhone: string;
  classApplying: string;
  status: string;
}

/**
 * DTO for creating an entrance test
 */
export interface CreateTestDto {
  testName: string;
  sessionId: string;
  testDate: string;
  startTime: string;
  durationMinutes: number;
  venue?: string;
  classNames: string[];
  maxCandidates?: number;
  subjects: TestSubject[];
}

/**
 * DTO for updating an entrance test
 */
export type UpdateTestDto = Partial<CreateTestDto>;

/**
 * DTO for registering a candidate
 */
export interface RegisterCandidateDto {
  applicationId: string;
}

/**
 * DTO for submitting test results
 */
export interface SubmitResultsDto {
  registrationId: string;
  marks: SubjectMarks[];
  remarks?: string;
}

/**
 * Bulk result submission DTO
 */
export interface BulkResultsDto {
  results: SubmitResultsDto[];
}

/**
 * Filter parameters for listing tests
 */
export interface TestFilterParams {
  sessionId?: string;
  status?: TestStatus;
  className?: string;
  fromDate?: string;
  toDate?: string;
}

/**
 * Status configuration for display
 */
export interface TestStatusConfig {
  label: string;
  variant: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  icon: string;
}

/**
 * Test status display configuration
 */
export const TEST_STATUS_CONFIG: Record<TestStatus, TestStatusConfig> = {
  scheduled: {
    label: 'Scheduled',
    variant: 'info',
    icon: 'fa-solid fa-calendar-check',
  },
  in_progress: {
    label: 'In Progress',
    variant: 'warning',
    icon: 'fa-solid fa-hourglass-half',
  },
  completed: {
    label: 'Completed',
    variant: 'success',
    icon: 'fa-solid fa-check-circle',
  },
  cancelled: {
    label: 'Cancelled',
    variant: 'danger',
    icon: 'fa-solid fa-times-circle',
  },
};

/**
 * Registration status display configuration
 */
export const REGISTRATION_STATUS_CONFIG: Record<RegistrationStatus, TestStatusConfig> = {
  registered: {
    label: 'Registered',
    variant: 'info',
    icon: 'fa-solid fa-user-check',
  },
  appeared: {
    label: 'Appeared',
    variant: 'primary',
    icon: 'fa-solid fa-clipboard-check',
  },
  absent: {
    label: 'Absent',
    variant: 'danger',
    icon: 'fa-solid fa-user-times',
  },
  result_pending: {
    label: 'Result Pending',
    variant: 'warning',
    icon: 'fa-solid fa-clock',
  },
  passed: {
    label: 'Passed',
    variant: 'success',
    icon: 'fa-solid fa-trophy',
  },
  failed: {
    label: 'Failed',
    variant: 'danger',
    icon: 'fa-solid fa-times',
  },
};

/**
 * Calculate total marks from subject marks array
 */
export function calculateTotalMarks(marks: SubjectMarks[]): number {
  return marks.reduce((sum, m) => sum + (m.obtainedMarks || 0), 0);
}

/**
 * Calculate max marks from subject marks array
 */
export function calculateMaxMarks(marks: SubjectMarks[]): number {
  return marks.reduce((sum, m) => sum + (m.maxMarks || 0), 0);
}

/**
 * Calculate percentage from marks
 */
export function calculatePercentage(obtained: number, max: number): number {
  if (max === 0) return 0;
  return Math.round((obtained / max) * 10000) / 100;
}

/**
 * Format time string for display
 */
export function formatTime(time: string): string {
  const [hours, minutes] = time.split(':');
  const h = parseInt(hours, 10);
  const ampm = h >= 12 ? 'PM' : 'AM';
  const h12 = h % 12 || 12;
  return `${h12}:${minutes} ${ampm}`;
}

/**
 * Format date string for display
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
}
