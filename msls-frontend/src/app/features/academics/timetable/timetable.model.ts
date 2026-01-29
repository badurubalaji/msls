/**
 * Timetable module models for Shifts, Day Patterns, and Period Slots.
 */

// ========================================
// Period Slot Type
// ========================================

export type PeriodSlotType = 'regular' | 'short' | 'assembly' | 'break' | 'lunch' | 'activity' | 'zero_period';

export const PERIOD_SLOT_TYPES: { value: PeriodSlotType; label: string; icon: string; color: string }[] = [
  { value: 'regular', label: 'Regular', icon: 'fa-book', color: 'bg-blue-100 text-blue-700' },
  { value: 'short', label: 'Short Period', icon: 'fa-clock', color: 'bg-purple-100 text-purple-700' },
  { value: 'assembly', label: 'Assembly', icon: 'fa-users', color: 'bg-amber-100 text-amber-700' },
  { value: 'break', label: 'Break', icon: 'fa-mug-hot', color: 'bg-green-100 text-green-700' },
  { value: 'lunch', label: 'Lunch', icon: 'fa-utensils', color: 'bg-orange-100 text-orange-700' },
  { value: 'activity', label: 'Activity', icon: 'fa-futbol', color: 'bg-teal-100 text-teal-700' },
  { value: 'zero_period', label: 'Zero Period', icon: 'fa-sun', color: 'bg-gray-100 text-gray-700' },
];

// ========================================
// Shift Models
// ========================================

export interface Shift {
  id: string;
  branchId: string;
  branchName?: string;
  name: string;
  code: string;
  startTime: string;
  endTime: string;
  description?: string;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateShiftRequest {
  branchId: string;
  name: string;
  code: string;
  startTime: string;
  endTime: string;
  description?: string;
  displayOrder?: number;
}

export interface UpdateShiftRequest {
  name?: string;
  code?: string;
  startTime?: string;
  endTime?: string;
  description?: string;
  displayOrder?: number;
  isActive?: boolean;
}

export interface ShiftListResponse {
  shifts: Shift[];
  total: number;
}

export interface ShiftFilter {
  branchId?: string;
  isActive?: boolean;
}

// ========================================
// Day Pattern Models
// ========================================

export interface DayPattern {
  id: string;
  name: string;
  code: string;
  description?: string;
  totalPeriods: number;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDayPatternRequest {
  name: string;
  code: string;
  description?: string;
  totalPeriods?: number;
  displayOrder?: number;
}

export interface UpdateDayPatternRequest {
  name?: string;
  code?: string;
  description?: string;
  totalPeriods?: number;
  displayOrder?: number;
  isActive?: boolean;
}

export interface DayPatternListResponse {
  dayPatterns: DayPattern[];
  total: number;
}

export interface DayPatternFilter {
  isActive?: boolean;
}

// ========================================
// Day Pattern Assignment Models
// ========================================

export const DAYS_OF_WEEK = [
  { value: 0, label: 'Sunday', short: 'Sun' },
  { value: 1, label: 'Monday', short: 'Mon' },
  { value: 2, label: 'Tuesday', short: 'Tue' },
  { value: 3, label: 'Wednesday', short: 'Wed' },
  { value: 4, label: 'Thursday', short: 'Thu' },
  { value: 5, label: 'Friday', short: 'Fri' },
  { value: 6, label: 'Saturday', short: 'Sat' },
];

export interface DayPatternAssignment {
  id: string;
  branchId: string;
  dayOfWeek: number;
  dayPatternId?: string;
  dayPatternName?: string;
  isWorkingDay: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateDayPatternAssignmentRequest {
  dayPatternId?: string | null;
  isWorkingDay?: boolean;
}

export interface DayPatternAssignmentListResponse {
  assignments: DayPatternAssignment[];
}

// ========================================
// Period Slot Models
// ========================================

export interface PeriodSlot {
  id: string;
  branchId: string;
  branchName?: string;
  name: string;
  periodNumber?: number;
  slotType: PeriodSlotType;
  startTime: string;
  endTime: string;
  durationMinutes: number;
  dayPatternId?: string;
  dayPatternName?: string;
  shiftId?: string;
  shiftName?: string;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreatePeriodSlotRequest {
  branchId: string;
  name: string;
  periodNumber?: number;
  slotType: PeriodSlotType;
  startTime: string;
  endTime: string;
  durationMinutes: number;
  dayPatternId?: string;
  shiftId?: string;
  displayOrder?: number;
}

export interface UpdatePeriodSlotRequest {
  name?: string;
  periodNumber?: number;
  slotType?: PeriodSlotType;
  startTime?: string;
  endTime?: string;
  durationMinutes?: number;
  dayPatternId?: string;
  shiftId?: string;
  displayOrder?: number;
  isActive?: boolean;
}

export interface PeriodSlotListResponse {
  periodSlots: PeriodSlot[];
  total: number;
}

export interface PeriodSlotFilter {
  branchId?: string;
  dayPatternId?: string;
  shiftId?: string;
  slotType?: PeriodSlotType;
  isActive?: boolean;
}
