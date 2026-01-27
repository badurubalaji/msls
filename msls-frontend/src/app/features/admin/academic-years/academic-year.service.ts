/**
 * MSLS Academic Year Service
 *
 * HTTP service for academic year, term, and holiday management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  AcademicYear,
  AcademicTerm,
  Holiday,
  CreateAcademicYearRequest,
  UpdateAcademicYearRequest,
  CreateTermRequest,
  UpdateTermRequest,
  CreateHolidayRequest,
  UpdateHolidayRequest,
} from './academic-year.model';

/** Response format for academic year list from backend */
interface AcademicYearListResponse {
  academicYears: AcademicYear[];
  total: number;
}

/** Response format for term list from backend */
interface TermListResponse {
  terms: AcademicTerm[];
  total: number;
}

/** Response format for holiday list from backend */
interface HolidayListResponse {
  holidays: Holiday[];
  total: number;
}

/**
 * AcademicYearService - Handles all academic year-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class AcademicYearService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/academic-years';

  // ============================================================================
  // Academic Year Operations
  // ============================================================================

  /**
   * Get all academic years for the current tenant
   */
  getAcademicYears(): Observable<AcademicYear[]> {
    return this.apiService.get<AcademicYearListResponse>(this.basePath).pipe(
      map(response => response.academicYears || [])
    );
  }

  /**
   * Get a single academic year by ID (includes terms and holidays)
   */
  getAcademicYear(id: string): Observable<AcademicYear> {
    return this.apiService.get<AcademicYear>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new academic year
   */
  createAcademicYear(data: CreateAcademicYearRequest): Observable<AcademicYear> {
    return this.apiService.post<AcademicYear>(this.basePath, data);
  }

  /**
   * Update an existing academic year
   */
  updateAcademicYear(id: string, data: UpdateAcademicYearRequest): Observable<AcademicYear> {
    return this.apiService.put<AcademicYear>(`${this.basePath}/${id}`, data);
  }

  /**
   * Set an academic year as current
   */
  setAsCurrent(id: string): Observable<AcademicYear> {
    return this.apiService.patch<AcademicYear>(`${this.basePath}/${id}/current`, {});
  }

  /**
   * Delete an academic year
   */
  deleteAcademicYear(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  // ============================================================================
  // Term Operations
  // ============================================================================

  /**
   * Get all terms for an academic year
   */
  getTerms(academicYearId: string): Observable<AcademicTerm[]> {
    return this.apiService.get<TermListResponse>(`${this.basePath}/${academicYearId}/terms`).pipe(
      map(response => response.terms || [])
    );
  }

  /**
   * Add a term to an academic year
   */
  createTerm(academicYearId: string, data: CreateTermRequest): Observable<AcademicTerm> {
    return this.apiService.post<AcademicTerm>(`${this.basePath}/${academicYearId}/terms`, data);
  }

  /**
   * Update a term
   */
  updateTerm(academicYearId: string, termId: string, data: UpdateTermRequest): Observable<AcademicTerm> {
    return this.apiService.put<AcademicTerm>(`${this.basePath}/${academicYearId}/terms/${termId}`, data);
  }

  /**
   * Delete a term
   */
  deleteTerm(academicYearId: string, termId: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${academicYearId}/terms/${termId}`);
  }

  // ============================================================================
  // Holiday Operations
  // ============================================================================

  /**
   * Get all holidays for an academic year
   */
  getHolidays(academicYearId: string): Observable<Holiday[]> {
    return this.apiService.get<HolidayListResponse>(`${this.basePath}/${academicYearId}/holidays`).pipe(
      map(response => response.holidays || [])
    );
  }

  /**
   * Add a holiday to an academic year
   */
  createHoliday(academicYearId: string, data: CreateHolidayRequest): Observable<Holiday> {
    return this.apiService.post<Holiday>(`${this.basePath}/${academicYearId}/holidays`, data);
  }

  /**
   * Update a holiday
   */
  updateHoliday(academicYearId: string, holidayId: string, data: UpdateHolidayRequest): Observable<Holiday> {
    return this.apiService.put<Holiday>(`${this.basePath}/${academicYearId}/holidays/${holidayId}`, data);
  }

  /**
   * Delete a holiday
   */
  deleteHoliday(academicYearId: string, holidayId: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${academicYearId}/holidays/${holidayId}`);
  }
}
