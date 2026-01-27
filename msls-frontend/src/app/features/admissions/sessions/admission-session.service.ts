/**
 * MSLS Admission Session Service
 *
 * HTTP service for admission session management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable, forkJoin, of } from 'rxjs';
import { map, switchMap, catchError } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  AdmissionSession,
  AdmissionSeat,
  CreateSessionRequest,
  UpdateSessionRequest,
  ChangeStatusRequest,
  SeatConfigRequest,
  AcademicYear,
} from './admission-session.model';

/** API response for session list */
interface SessionListResponse {
  sessions: ApiSession[];
  total: number;
}

/** API response for seat list */
interface SeatListResponse {
  seats: ApiSeat[];
  total: number;
}

/** API response for academic year list */
interface AcademicYearListResponse {
  academicYears: ApiAcademicYear[];
  total: number;
}

/** Session as returned from API */
interface ApiSession {
  id: string;
  branchId?: string;
  academicYearId?: string;
  name: string;
  description?: string;
  startDate: string;
  endDate: string;
  status: string;
  applicationFee: number;
  requiredDocuments: string[];
  settings: {
    allowOnlineApplication?: boolean;
    notifyOnApplication?: boolean;
    autoConfirmPayment?: boolean;
    maxApplicationsPerDay?: number;
    instructions?: string;
  };
  createdAt: string;
  updatedAt: string;
  seats?: ApiSeat[];
  stats?: {
    totalApplications: number;
    approvedCount: number;
    pendingCount: number;
    rejectedCount: number;
    totalSeats: number;
    filledSeats: number;
    availableSeats: number;
  };
}

/** Seat as returned from API */
interface ApiSeat {
  id: string;
  sessionId: string;
  className: string;
  totalSeats: number;
  filledSeats: number;
  availableSeats: number;
  waitlistLimit: number;
  reservedSeats: Record<string, number>;
  createdAt: string;
  updatedAt: string;
}

/** Academic year as returned from API */
interface ApiAcademicYear {
  id: string;
  name: string;
  startDate: string;
  endDate: string;
  isCurrent: boolean;
}

/**
 * AdmissionSessionService - Handles all admission session-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class AdmissionSessionService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/admission-sessions';
  private readonly academicYearsPath = '/academic-years';

  /**
   * Get all admission sessions with resolved academic year names
   */
  getSessions(): Observable<AdmissionSession[]> {
    return forkJoin({
      sessions: this.apiService.get<SessionListResponse>(`${this.basePath}?includeSeats=true`),
      academicYears: this.getAcademicYears().pipe(catchError(() => of([]))),
    }).pipe(
      map(({ sessions, academicYears }) => {
        const yearMap = new Map(academicYears.map(ay => [ay.id, ay.name]));
        return (sessions.sessions || []).map(s => ({
          ...this.mapApiSessionToModel(s),
          academicYearName: s.academicYearId ? yearMap.get(s.academicYearId) || 'Unknown' : 'Not Set',
        }));
      })
    );
  }

  /**
   * Get a single session by ID
   */
  getSession(id: string): Observable<AdmissionSession> {
    return this.apiService.get<ApiSession>(`${this.basePath}/${id}`).pipe(
      map(session => this.mapApiSessionToModel(session))
    );
  }

  /**
   * Create a new admission session
   */
  createSession(data: CreateSessionRequest): Observable<AdmissionSession> {
    return this.apiService.post<ApiSession>(this.basePath, data).pipe(
      map(session => this.mapApiSessionToModel(session))
    );
  }

  /**
   * Update an existing session
   */
  updateSession(id: string, data: UpdateSessionRequest): Observable<AdmissionSession> {
    return this.apiService.put<ApiSession>(`${this.basePath}/${id}`, data).pipe(
      map(session => this.mapApiSessionToModel(session))
    );
  }

  /**
   * Change session status (open/close)
   */
  changeStatus(id: string, data: ChangeStatusRequest): Observable<AdmissionSession> {
    return this.apiService.patch<ApiSession>(`${this.basePath}/${id}/status`, data).pipe(
      map(session => this.mapApiSessionToModel(session))
    );
  }

  /**
   * Delete an admission session
   */
  deleteSession(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }

  /**
   * Get seat configurations for a session
   */
  getSeats(sessionId: string): Observable<AdmissionSeat[]> {
    return this.apiService.get<SeatListResponse>(`${this.basePath}/${sessionId}/seats`).pipe(
      map(response => response.seats.map(s => this.mapApiSeatToModel(s)))
    );
  }

  /**
   * Add seat configuration
   */
  addSeat(sessionId: string, data: SeatConfigRequest): Observable<AdmissionSeat> {
    return this.apiService.post<ApiSeat>(`${this.basePath}/${sessionId}/seats`, data).pipe(
      map(seat => this.mapApiSeatToModel(seat))
    );
  }

  /**
   * Update seat configuration
   */
  updateSeat(sessionId: string, seatId: string, data: SeatConfigRequest): Observable<AdmissionSeat> {
    return this.apiService.put<ApiSeat>(`${this.basePath}/${sessionId}/seats/${seatId}`, data).pipe(
      map(seat => this.mapApiSeatToModel(seat))
    );
  }

  /**
   * Delete seat configuration
   */
  deleteSeat(sessionId: string, seatId: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${sessionId}/seats/${seatId}`);
  }

  /**
   * Get academic years for dropdown
   */
  getAcademicYears(): Observable<AcademicYear[]> {
    return this.apiService.get<AcademicYearListResponse>(this.academicYearsPath).pipe(
      map(response => (response.academicYears || []).map(ay => ({
        id: ay.id,
        name: ay.name,
        startDate: ay.startDate,
        endDate: ay.endDate,
        isCurrent: ay.isCurrent,
      })))
    );
  }

  /**
   * Map API session response to frontend model
   */
  private mapApiSessionToModel(session: ApiSession): AdmissionSession {
    return {
      id: session.id,
      name: session.name,
      academicYearId: session.academicYearId || '',
      branchId: session.branchId,
      startDate: session.startDate,
      endDate: session.endDate,
      status: session.status as 'upcoming' | 'open' | 'closed',
      applicationFee: session.applicationFee,
      requiredDocuments: session.requiredDocuments || [],
      settings: session.settings || {},
      totalApplications: session.stats?.totalApplications || 0,
      totalSeats: session.stats?.totalSeats || 0,
      filledSeats: session.stats?.filledSeats || 0,
      createdAt: session.createdAt,
      updatedAt: session.updatedAt,
    };
  }

  /**
   * Map API seat response to frontend model
   */
  private mapApiSeatToModel(seat: ApiSeat): AdmissionSeat {
    return {
      id: seat.id,
      sessionId: seat.sessionId,
      className: seat.className,
      totalSeats: seat.totalSeats,
      filledSeats: seat.filledSeats,
      waitlistLimit: seat.waitlistLimit,
      reservedSeats: seat.reservedSeats || {},
      createdAt: seat.createdAt,
      updatedAt: seat.updatedAt,
    };
  }
}
