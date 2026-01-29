/**
 * MSLS Timetable Service
 *
 * HTTP service for timetable management API calls including shifts, day patterns, and period slots.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  Shift,
  ShiftListResponse,
  ShiftFilter,
  CreateShiftRequest,
  UpdateShiftRequest,
  DayPattern,
  DayPatternListResponse,
  DayPatternFilter,
  CreateDayPatternRequest,
  UpdateDayPatternRequest,
  DayPatternAssignment,
  DayPatternAssignmentListResponse,
  UpdateDayPatternAssignmentRequest,
  PeriodSlot,
  PeriodSlotListResponse,
  PeriodSlotFilter,
  CreatePeriodSlotRequest,
  UpdatePeriodSlotRequest,
} from './timetable.model';

/**
 * TimetableService - Handles all timetable-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class TimetableService {
  private readonly apiService = inject(ApiService);

  // ========================================
  // Shift Methods
  // ========================================

  getShifts(filter?: ShiftFilter): Observable<Shift[]> {
    const params = this.buildShiftFilterParams(filter);
    return this.apiService.get<ShiftListResponse>('/shifts', { params }).pipe(
      map(response => response.shifts || [])
    );
  }

  getShiftsWithTotal(filter?: ShiftFilter): Observable<ShiftListResponse> {
    const params = this.buildShiftFilterParams(filter);
    return this.apiService.get<ShiftListResponse>('/shifts', { params });
  }

  getShift(id: string): Observable<Shift> {
    return this.apiService.get<Shift>(`/shifts/${id}`);
  }

  createShift(data: CreateShiftRequest): Observable<Shift> {
    return this.apiService.post<Shift>('/shifts', data);
  }

  updateShift(id: string, data: UpdateShiftRequest): Observable<Shift> {
    return this.apiService.put<Shift>(`/shifts/${id}`, data);
  }

  deleteShift(id: string): Observable<void> {
    return this.apiService.delete<void>(`/shifts/${id}`);
  }

  private buildShiftFilterParams(filter?: ShiftFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.branchId) params['branch_id'] = filter.branchId;
    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);

    return params;
  }

  // ========================================
  // Day Pattern Methods
  // ========================================

  getDayPatterns(filter?: DayPatternFilter): Observable<DayPattern[]> {
    const params = this.buildDayPatternFilterParams(filter);
    return this.apiService.get<DayPatternListResponse>('/day-patterns', { params }).pipe(
      map(response => response.dayPatterns || [])
    );
  }

  getDayPatternsWithTotal(filter?: DayPatternFilter): Observable<DayPatternListResponse> {
    const params = this.buildDayPatternFilterParams(filter);
    return this.apiService.get<DayPatternListResponse>('/day-patterns', { params });
  }

  getDayPattern(id: string): Observable<DayPattern> {
    return this.apiService.get<DayPattern>(`/day-patterns/${id}`);
  }

  createDayPattern(data: CreateDayPatternRequest): Observable<DayPattern> {
    return this.apiService.post<DayPattern>('/day-patterns', data);
  }

  updateDayPattern(id: string, data: UpdateDayPatternRequest): Observable<DayPattern> {
    return this.apiService.put<DayPattern>(`/day-patterns/${id}`, data);
  }

  deleteDayPattern(id: string): Observable<void> {
    return this.apiService.delete<void>(`/day-patterns/${id}`);
  }

  private buildDayPatternFilterParams(filter?: DayPatternFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);

    return params;
  }

  // ========================================
  // Day Pattern Assignment Methods
  // ========================================

  getDayPatternAssignments(branchId: string): Observable<DayPatternAssignment[]> {
    return this.apiService.get<DayPatternAssignmentListResponse>(
      '/day-pattern-assignments',
      { params: { branch_id: branchId } }
    ).pipe(
      map(response => response.assignments || [])
    );
  }

  updateDayPatternAssignment(
    branchId: string,
    dayOfWeek: number,
    data: UpdateDayPatternAssignmentRequest
  ): Observable<DayPatternAssignment> {
    return this.apiService.put<DayPatternAssignment>(
      `/day-pattern-assignments/${dayOfWeek}`,
      data,
      { params: { branch_id: branchId } }
    );
  }

  // ========================================
  // Period Slot Methods
  // ========================================

  getPeriodSlots(filter?: PeriodSlotFilter): Observable<PeriodSlot[]> {
    const params = this.buildPeriodSlotFilterParams(filter);
    return this.apiService.get<PeriodSlotListResponse>('/period-slots', { params }).pipe(
      map(response => response.periodSlots || [])
    );
  }

  getPeriodSlotsWithTotal(filter?: PeriodSlotFilter): Observable<PeriodSlotListResponse> {
    const params = this.buildPeriodSlotFilterParams(filter);
    return this.apiService.get<PeriodSlotListResponse>('/period-slots', { params });
  }

  getPeriodSlot(id: string): Observable<PeriodSlot> {
    return this.apiService.get<PeriodSlot>(`/period-slots/${id}`);
  }

  createPeriodSlot(data: CreatePeriodSlotRequest): Observable<PeriodSlot> {
    return this.apiService.post<PeriodSlot>('/period-slots', data);
  }

  updatePeriodSlot(id: string, data: UpdatePeriodSlotRequest): Observable<PeriodSlot> {
    return this.apiService.put<PeriodSlot>(`/period-slots/${id}`, data);
  }

  deletePeriodSlot(id: string): Observable<void> {
    return this.apiService.delete<void>(`/period-slots/${id}`);
  }

  private buildPeriodSlotFilterParams(filter?: PeriodSlotFilter): Record<string, string> {
    const params: Record<string, string> = {};
    if (!filter) return params;

    if (filter.branchId) params['branch_id'] = filter.branchId;
    if (filter.dayPatternId) params['day_pattern_id'] = filter.dayPatternId;
    if (filter.shiftId) params['shift_id'] = filter.shiftId;
    if (filter.slotType) params['slot_type'] = filter.slotType;
    if (filter.isActive !== undefined) params['is_active'] = String(filter.isActive);

    return params;
  }
}
