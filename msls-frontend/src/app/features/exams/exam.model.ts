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
