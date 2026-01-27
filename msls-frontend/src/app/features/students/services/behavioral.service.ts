/**
 * Behavioral Service
 *
 * Handles student behavioral incident API operations.
 */

import { Injectable, inject, signal } from '@angular/core';
import { Observable, tap, catchError, throwError } from 'rxjs';

import { ApiService } from '../../../core/services';
import {
  BehavioralIncident,
  BehaviorSummary,
  FollowUp,
  BehavioralIncidentListResponse,
  CreateBehavioralIncidentRequest,
  UpdateBehavioralIncidentRequest,
  CreateFollowUpRequest,
  UpdateFollowUpRequest,
  BehavioralIncidentFilter,
  PendingFollowUpsResponse,
} from '../models/behavioral.model';

@Injectable({
  providedIn: 'root',
})
export class BehavioralService {
  private api = inject(ApiService);

  // State
  private _incidents = signal<BehavioralIncident[]>([]);
  private _summary = signal<BehaviorSummary | null>(null);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);

  readonly incidents = this._incidents.asReadonly();
  readonly summary = this._summary.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();

  // ==========================================================================
  // Incident Operations
  // ==========================================================================

  loadIncidents(studentId: string, filter?: BehavioralIncidentFilter): Observable<BehavioralIncidentListResponse> {
    this._loading.set(true);
    this._error.set(null);

    const params: Record<string, string | number> = {};
    if (filter?.type) params['type'] = filter.type;
    if (filter?.severity) params['severity'] = filter.severity;
    if (filter?.dateFrom) params['dateFrom'] = filter.dateFrom;
    if (filter?.dateTo) params['dateTo'] = filter.dateTo;
    if (filter?.limit) params['limit'] = filter.limit;
    if (filter?.offset) params['offset'] = filter.offset;

    return this.api.get<BehavioralIncidentListResponse>(`/students/${studentId}/behavioral-incidents`, { params }).pipe(
      tap((response) => {
        this._incidents.set(response.incidents);
        this._loading.set(false);
      }),
      catchError((error) => {
        this._error.set(error.message);
        this._loading.set(false);
        return throwError(() => error);
      })
    );
  }

  getIncident(studentId: string, incidentId: string): Observable<BehavioralIncident> {
    return this.api.get<BehavioralIncident>(`/students/${studentId}/behavioral-incidents/${incidentId}`);
  }

  createIncident(studentId: string, request: CreateBehavioralIncidentRequest): Observable<BehavioralIncident> {
    return this.api.post<BehavioralIncident>(`/students/${studentId}/behavioral-incidents`, request).pipe(
      tap((incident) => {
        this._incidents.update(list => [incident, ...list]);
        // Update summary counts
        this._summary.update(s => s ? {
          ...s,
          totalIncidents: s.totalIncidents + 1,
          thisMonthCount: s.thisMonthCount + 1,
          positiveCount: incident.incidentType === 'positive_recognition' ? s.positiveCount + 1 : s.positiveCount,
          minorInfractionCount: incident.incidentType === 'minor_infraction' ? s.minorInfractionCount + 1 : s.minorInfractionCount,
          majorViolationCount: incident.incidentType === 'major_violation' ? s.majorViolationCount + 1 : s.majorViolationCount,
        } : null);
      })
    );
  }

  updateIncident(studentId: string, incidentId: string, request: UpdateBehavioralIncidentRequest): Observable<BehavioralIncident> {
    return this.api.put<BehavioralIncident>(`/students/${studentId}/behavioral-incidents/${incidentId}`, request).pipe(
      tap((incident) => {
        this._incidents.update(list => list.map(i => i.id === incidentId ? incident : i));
      })
    );
  }

  deleteIncident(studentId: string, incidentId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/behavioral-incidents/${incidentId}`).pipe(
      tap(() => {
        const deleted = this._incidents().find(i => i.id === incidentId);
        this._incidents.update(list => list.filter(i => i.id !== incidentId));
        // Update summary counts
        if (deleted) {
          this._summary.update(s => s ? {
            ...s,
            totalIncidents: Math.max(0, s.totalIncidents - 1),
            positiveCount: deleted.incidentType === 'positive_recognition' ? Math.max(0, s.positiveCount - 1) : s.positiveCount,
            minorInfractionCount: deleted.incidentType === 'minor_infraction' ? Math.max(0, s.minorInfractionCount - 1) : s.minorInfractionCount,
            majorViolationCount: deleted.incidentType === 'major_violation' ? Math.max(0, s.majorViolationCount - 1) : s.majorViolationCount,
          } : null);
        }
      })
    );
  }

  // ==========================================================================
  // Summary Operations
  // ==========================================================================

  loadSummary(studentId: string): Observable<BehaviorSummary> {
    return this.api.get<BehaviorSummary>(`/students/${studentId}/behavioral-summary`).pipe(
      tap((summary) => {
        this._summary.set(summary);
      }),
      catchError((error) => {
        this._error.set(error.message);
        return throwError(() => error);
      })
    );
  }

  // ==========================================================================
  // Follow-Up Operations
  // ==========================================================================

  createFollowUp(incidentId: string, request: CreateFollowUpRequest): Observable<FollowUp> {
    return this.api.post<FollowUp>(`/behavioral-incidents/${incidentId}/follow-ups`, request).pipe(
      tap((followUp) => {
        // Update the incident's follow-ups list
        this._incidents.update(list => list.map(i => {
          if (i.id === incidentId) {
            return { ...i, followUps: [...(i.followUps || []), followUp] };
          }
          return i;
        }));
        // Update pending follow-ups count
        this._summary.update(s => s ? { ...s, pendingFollowUps: s.pendingFollowUps + 1 } : null);
      })
    );
  }

  updateFollowUp(incidentId: string, followUpId: string, request: UpdateFollowUpRequest): Observable<FollowUp> {
    return this.api.put<FollowUp>(`/behavioral-incidents/${incidentId}/follow-ups/${followUpId}`, request).pipe(
      tap((followUp) => {
        // Update the incident's follow-ups list
        this._incidents.update(list => list.map(i => {
          if (i.id === incidentId && i.followUps) {
            return {
              ...i,
              followUps: i.followUps.map(f => f.id === followUpId ? followUp : f)
            };
          }
          return i;
        }));
        // Update pending count if status changed
        if (request.status && request.status !== 'pending') {
          this._summary.update(s => s ? { ...s, pendingFollowUps: Math.max(0, s.pendingFollowUps - 1) } : null);
        }
      })
    );
  }

  deleteFollowUp(incidentId: string, followUpId: string): Observable<void> {
    return this.api.delete<void>(`/behavioral-incidents/${incidentId}/follow-ups/${followUpId}`).pipe(
      tap(() => {
        // Remove from incident's follow-ups list
        this._incidents.update(list => list.map(i => {
          if (i.id === incidentId && i.followUps) {
            const deleted = i.followUps.find(f => f.id === followUpId);
            if (deleted?.status === 'pending') {
              this._summary.update(s => s ? { ...s, pendingFollowUps: Math.max(0, s.pendingFollowUps - 1) } : null);
            }
            return { ...i, followUps: i.followUps.filter(f => f.id !== followUpId) };
          }
          return i;
        }));
      })
    );
  }

  loadPendingFollowUps(limit = 20, offset = 0): Observable<PendingFollowUpsResponse> {
    return this.api.get<PendingFollowUpsResponse>('/follow-ups/pending', {
      params: { limit, offset }
    });
  }

  // ==========================================================================
  // Reset
  // ==========================================================================

  reset(): void {
    this._incidents.set([]);
    this._summary.set(null);
    this._loading.set(false);
    this._error.set(null);
  }
}
