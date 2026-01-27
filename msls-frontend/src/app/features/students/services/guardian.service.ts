/**
 * Guardian Service
 *
 * Handles guardian and emergency contact API operations.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { Observable, tap, catchError, throwError } from 'rxjs';

import { ApiService } from '../../../core/services';
import {
  Guardian,
  GuardianListResponse,
  CreateGuardianRequest,
  UpdateGuardianRequest,
  EmergencyContact,
  EmergencyContactListResponse,
  CreateEmergencyContactRequest,
  UpdateEmergencyContactRequest,
} from '../models/guardian.model';

@Injectable({
  providedIn: 'root',
})
export class GuardianService {
  private api = inject(ApiService);

  // =========================================================================
  // Guardian State
  // =========================================================================

  private _guardians = signal<Guardian[]>([]);
  private _selectedGuardian = signal<Guardian | null>(null);
  private _guardiansLoading = signal<boolean>(false);
  private _guardiansError = signal<string | null>(null);

  readonly guardians = this._guardians.asReadonly();
  readonly selectedGuardian = this._selectedGuardian.asReadonly();
  readonly guardiansLoading = this._guardiansLoading.asReadonly();
  readonly guardiansError = this._guardiansError.asReadonly();

  readonly primaryGuardian = computed(() =>
    this._guardians().find(g => g.isPrimary) || null
  );

  // =========================================================================
  // Emergency Contact State
  // =========================================================================

  private _emergencyContacts = signal<EmergencyContact[]>([]);
  private _contactsLoading = signal<boolean>(false);
  private _contactsError = signal<string | null>(null);

  readonly emergencyContacts = this._emergencyContacts.asReadonly();
  readonly contactsLoading = this._contactsLoading.asReadonly();
  readonly contactsError = this._contactsError.asReadonly();

  // =========================================================================
  // Guardian API Methods
  // =========================================================================

  /**
   * Load all guardians for a student.
   */
  loadGuardians(studentId: string): Observable<GuardianListResponse> {
    this._guardiansLoading.set(true);
    this._guardiansError.set(null);

    return this.api.get<GuardianListResponse>(`/students/${studentId}/guardians`).pipe(
      tap((response) => {
        this._guardians.set(response.guardians);
        this._guardiansLoading.set(false);
      }),
      catchError((error) => {
        this._guardiansError.set(error.message);
        this._guardiansLoading.set(false);
        return throwError(() => error);
      })
    );
  }

  /**
   * Get a guardian by ID.
   */
  getGuardian(studentId: string, guardianId: string): Observable<Guardian> {
    return this.api.get<Guardian>(`/students/${studentId}/guardians/${guardianId}`).pipe(
      tap((guardian) => {
        this._selectedGuardian.set(guardian);
      })
    );
  }

  /**
   * Create a new guardian.
   */
  createGuardian(studentId: string, request: CreateGuardianRequest): Observable<Guardian> {
    return this.api.post<Guardian>(`/students/${studentId}/guardians`, request).pipe(
      tap((guardian) => {
        this._guardians.update((current) => [...current, guardian]);
        // If new guardian is primary, update other guardians
        if (guardian.isPrimary) {
          this._guardians.update((current) =>
            current.map(g => g.id === guardian.id ? g : { ...g, isPrimary: false })
          );
        }
      })
    );
  }

  /**
   * Update a guardian.
   */
  updateGuardian(
    studentId: string,
    guardianId: string,
    request: UpdateGuardianRequest
  ): Observable<Guardian> {
    return this.api.put<Guardian>(`/students/${studentId}/guardians/${guardianId}`, request).pipe(
      tap((guardian) => {
        this._guardians.update((current) =>
          current.map((g) => (g.id === guardianId ? guardian : g))
        );
        // If updated guardian is now primary, update others
        if (guardian.isPrimary) {
          this._guardians.update((current) =>
            current.map(g => g.id === guardian.id ? g : { ...g, isPrimary: false })
          );
        }
        if (this._selectedGuardian()?.id === guardianId) {
          this._selectedGuardian.set(guardian);
        }
      })
    );
  }

  /**
   * Delete a guardian.
   */
  deleteGuardian(studentId: string, guardianId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/guardians/${guardianId}`).pipe(
      tap(() => {
        this._guardians.update((current) => current.filter((g) => g.id !== guardianId));
        if (this._selectedGuardian()?.id === guardianId) {
          this._selectedGuardian.set(null);
        }
      })
    );
  }

  /**
   * Set a guardian as primary.
   */
  setPrimaryGuardian(studentId: string, guardianId: string): Observable<Guardian> {
    return this.api.post<Guardian>(
      `/students/${studentId}/guardians/${guardianId}/set-primary`,
      {}
    ).pipe(
      tap((guardian) => {
        // Update all guardians - set the selected one as primary, others as not primary
        this._guardians.update((current) =>
          current.map((g) => ({
            ...g,
            isPrimary: g.id === guardianId,
          }))
        );
      })
    );
  }

  // =========================================================================
  // Emergency Contact API Methods
  // =========================================================================

  /**
   * Load all emergency contacts for a student.
   */
  loadEmergencyContacts(studentId: string): Observable<EmergencyContactListResponse> {
    this._contactsLoading.set(true);
    this._contactsError.set(null);

    return this.api.get<EmergencyContactListResponse>(
      `/students/${studentId}/emergency-contacts`
    ).pipe(
      tap((response) => {
        this._emergencyContacts.set(response.contacts);
        this._contactsLoading.set(false);
      }),
      catchError((error) => {
        this._contactsError.set(error.message);
        this._contactsLoading.set(false);
        return throwError(() => error);
      })
    );
  }

  /**
   * Get an emergency contact by ID.
   */
  getEmergencyContact(studentId: string, contactId: string): Observable<EmergencyContact> {
    return this.api.get<EmergencyContact>(
      `/students/${studentId}/emergency-contacts/${contactId}`
    );
  }

  /**
   * Create a new emergency contact.
   */
  createEmergencyContact(
    studentId: string,
    request: CreateEmergencyContactRequest
  ): Observable<EmergencyContact> {
    return this.api.post<EmergencyContact>(
      `/students/${studentId}/emergency-contacts`,
      request
    ).pipe(
      tap((contact) => {
        this._emergencyContacts.update((current) =>
          [...current, contact].sort((a, b) => a.priority - b.priority)
        );
      })
    );
  }

  /**
   * Update an emergency contact.
   */
  updateEmergencyContact(
    studentId: string,
    contactId: string,
    request: UpdateEmergencyContactRequest
  ): Observable<EmergencyContact> {
    return this.api.put<EmergencyContact>(
      `/students/${studentId}/emergency-contacts/${contactId}`,
      request
    ).pipe(
      tap((contact) => {
        this._emergencyContacts.update((current) =>
          current.map((c) => (c.id === contactId ? contact : c))
            .sort((a, b) => a.priority - b.priority)
        );
      })
    );
  }

  /**
   * Delete an emergency contact.
   */
  deleteEmergencyContact(studentId: string, contactId: string): Observable<void> {
    return this.api.delete<void>(
      `/students/${studentId}/emergency-contacts/${contactId}`
    ).pipe(
      tap(() => {
        this._emergencyContacts.update((current) =>
          current.filter((c) => c.id !== contactId)
        );
      })
    );
  }

  // =========================================================================
  // State Reset
  // =========================================================================

  /**
   * Reset all guardian state.
   */
  resetGuardians(): void {
    this._guardians.set([]);
    this._selectedGuardian.set(null);
    this._guardiansLoading.set(false);
    this._guardiansError.set(null);
  }

  /**
   * Reset all emergency contact state.
   */
  resetEmergencyContacts(): void {
    this._emergencyContacts.set([]);
    this._contactsLoading.set(false);
    this._contactsError.set(null);
  }

  /**
   * Reset all state.
   */
  reset(): void {
    this.resetGuardians();
    this.resetEmergencyContacts();
  }
}
