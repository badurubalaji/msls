/**
 * Health Service
 *
 * Handles student health records API operations.
 */

import { Injectable, inject, signal } from '@angular/core';
import { Observable, tap, catchError, throwError } from 'rxjs';

import { ApiService } from '../../../core/services';
import {
  HealthProfile,
  HealthSummary,
  Allergy,
  AllergyListResponse,
  ChronicCondition,
  ConditionListResponse,
  Medication,
  MedicationListResponse,
  Vaccination,
  VaccinationListResponse,
  MedicalIncident,
  IncidentListResponse,
  CreateHealthProfileRequest,
  CreateAllergyRequest,
  CreateConditionRequest,
  CreateMedicationRequest,
  CreateVaccinationRequest,
  CreateIncidentRequest,
} from '../models/health.model';

@Injectable({
  providedIn: 'root',
})
export class HealthService {
  private api = inject(ApiService);

  // State
  private _healthSummary = signal<HealthSummary | null>(null);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);

  readonly healthSummary = this._healthSummary.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();

  // ==========================================================================
  // Health Summary
  // ==========================================================================

  loadHealthSummary(studentId: string): Observable<HealthSummary> {
    this._loading.set(true);
    this._error.set(null);

    return this.api.get<HealthSummary>(`/students/${studentId}/health`).pipe(
      tap((summary) => {
        this._healthSummary.set(summary);
        this._loading.set(false);
      }),
      catchError((error) => {
        this._error.set(error.message);
        this._loading.set(false);
        return throwError(() => error);
      })
    );
  }

  // ==========================================================================
  // Health Profile
  // ==========================================================================

  getHealthProfile(studentId: string): Observable<HealthProfile> {
    return this.api.get<HealthProfile>(`/students/${studentId}/health/profile`);
  }

  saveHealthProfile(studentId: string, request: CreateHealthProfileRequest): Observable<HealthProfile> {
    return this.api.put<HealthProfile>(`/students/${studentId}/health/profile`, request).pipe(
      tap((profile) => {
        this._healthSummary.update(s => s ? { ...s, profile } : null);
      })
    );
  }

  // ==========================================================================
  // Allergies
  // ==========================================================================

  listAllergies(studentId: string, activeOnly = false): Observable<AllergyListResponse> {
    const params = activeOnly ? '?active=true' : '';
    return this.api.get<AllergyListResponse>(`/students/${studentId}/health/allergies${params}`);
  }

  getAllergy(studentId: string, allergyId: string): Observable<Allergy> {
    return this.api.get<Allergy>(`/students/${studentId}/health/allergies/${allergyId}`);
  }

  createAllergy(studentId: string, request: CreateAllergyRequest): Observable<Allergy> {
    return this.api.post<Allergy>(`/students/${studentId}/health/allergies`, request).pipe(
      tap((allergy) => {
        this._healthSummary.update(s => s ? {
          ...s,
          allergies: [...s.allergies, allergy]
        } : null);
      })
    );
  }

  updateAllergy(studentId: string, allergyId: string, request: Partial<CreateAllergyRequest>): Observable<Allergy> {
    return this.api.put<Allergy>(`/students/${studentId}/health/allergies/${allergyId}`, request).pipe(
      tap((allergy) => {
        this._healthSummary.update(s => s ? {
          ...s,
          allergies: s.allergies.map(a => a.id === allergyId ? allergy : a)
        } : null);
      })
    );
  }

  deleteAllergy(studentId: string, allergyId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/health/allergies/${allergyId}`).pipe(
      tap(() => {
        this._healthSummary.update(s => s ? {
          ...s,
          allergies: s.allergies.filter(a => a.id !== allergyId)
        } : null);
      })
    );
  }

  // ==========================================================================
  // Chronic Conditions
  // ==========================================================================

  listConditions(studentId: string, activeOnly = false): Observable<ConditionListResponse> {
    const params = activeOnly ? '?active=true' : '';
    return this.api.get<ConditionListResponse>(`/students/${studentId}/health/conditions${params}`);
  }

  getCondition(studentId: string, conditionId: string): Observable<ChronicCondition> {
    return this.api.get<ChronicCondition>(`/students/${studentId}/health/conditions/${conditionId}`);
  }

  createCondition(studentId: string, request: CreateConditionRequest): Observable<ChronicCondition> {
    return this.api.post<ChronicCondition>(`/students/${studentId}/health/conditions`, request).pipe(
      tap((condition) => {
        this._healthSummary.update(s => s ? {
          ...s,
          conditions: [...s.conditions, condition]
        } : null);
      })
    );
  }

  updateCondition(studentId: string, conditionId: string, request: Partial<CreateConditionRequest>): Observable<ChronicCondition> {
    return this.api.put<ChronicCondition>(`/students/${studentId}/health/conditions/${conditionId}`, request).pipe(
      tap((condition) => {
        this._healthSummary.update(s => s ? {
          ...s,
          conditions: s.conditions.map(c => c.id === conditionId ? condition : c)
        } : null);
      })
    );
  }

  deleteCondition(studentId: string, conditionId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/health/conditions/${conditionId}`).pipe(
      tap(() => {
        this._healthSummary.update(s => s ? {
          ...s,
          conditions: s.conditions.filter(c => c.id !== conditionId)
        } : null);
      })
    );
  }

  // ==========================================================================
  // Medications
  // ==========================================================================

  listMedications(studentId: string, activeOnly = false): Observable<MedicationListResponse> {
    const params = activeOnly ? '?active=true' : '';
    return this.api.get<MedicationListResponse>(`/students/${studentId}/health/medications${params}`);
  }

  getMedication(studentId: string, medicationId: string): Observable<Medication> {
    return this.api.get<Medication>(`/students/${studentId}/health/medications/${medicationId}`);
  }

  createMedication(studentId: string, request: CreateMedicationRequest): Observable<Medication> {
    return this.api.post<Medication>(`/students/${studentId}/health/medications`, request).pipe(
      tap((medication) => {
        this._healthSummary.update(s => s ? {
          ...s,
          medications: [...s.medications, medication]
        } : null);
      })
    );
  }

  updateMedication(studentId: string, medicationId: string, request: Partial<CreateMedicationRequest>): Observable<Medication> {
    return this.api.put<Medication>(`/students/${studentId}/health/medications/${medicationId}`, request).pipe(
      tap((medication) => {
        this._healthSummary.update(s => s ? {
          ...s,
          medications: s.medications.map(m => m.id === medicationId ? medication : m)
        } : null);
      })
    );
  }

  deleteMedication(studentId: string, medicationId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/health/medications/${medicationId}`).pipe(
      tap(() => {
        this._healthSummary.update(s => s ? {
          ...s,
          medications: s.medications.filter(m => m.id !== medicationId)
        } : null);
      })
    );
  }

  // ==========================================================================
  // Vaccinations
  // ==========================================================================

  listVaccinations(studentId: string): Observable<VaccinationListResponse> {
    return this.api.get<VaccinationListResponse>(`/students/${studentId}/health/vaccinations`);
  }

  getVaccination(studentId: string, vaccinationId: string): Observable<Vaccination> {
    return this.api.get<Vaccination>(`/students/${studentId}/health/vaccinations/${vaccinationId}`);
  }

  createVaccination(studentId: string, request: CreateVaccinationRequest): Observable<Vaccination> {
    return this.api.post<Vaccination>(`/students/${studentId}/health/vaccinations`, request).pipe(
      tap((vaccination) => {
        this._healthSummary.update(s => s ? {
          ...s,
          vaccinations: [...s.vaccinations, vaccination]
        } : null);
      })
    );
  }

  updateVaccination(studentId: string, vaccinationId: string, request: Partial<CreateVaccinationRequest>): Observable<Vaccination> {
    return this.api.put<Vaccination>(`/students/${studentId}/health/vaccinations/${vaccinationId}`, request).pipe(
      tap((vaccination) => {
        this._healthSummary.update(s => s ? {
          ...s,
          vaccinations: s.vaccinations.map(v => v.id === vaccinationId ? vaccination : v)
        } : null);
      })
    );
  }

  deleteVaccination(studentId: string, vaccinationId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/health/vaccinations/${vaccinationId}`).pipe(
      tap(() => {
        this._healthSummary.update(s => s ? {
          ...s,
          vaccinations: s.vaccinations.filter(v => v.id !== vaccinationId)
        } : null);
      })
    );
  }

  // ==========================================================================
  // Medical Incidents
  // ==========================================================================

  listIncidents(studentId: string, limit?: number): Observable<IncidentListResponse> {
    const params = limit ? `?limit=${limit}` : '';
    return this.api.get<IncidentListResponse>(`/students/${studentId}/health/incidents${params}`);
  }

  getIncident(studentId: string, incidentId: string): Observable<MedicalIncident> {
    return this.api.get<MedicalIncident>(`/students/${studentId}/health/incidents/${incidentId}`);
  }

  createIncident(studentId: string, request: CreateIncidentRequest): Observable<MedicalIncident> {
    return this.api.post<MedicalIncident>(`/students/${studentId}/health/incidents`, request).pipe(
      tap((incident) => {
        this._healthSummary.update(s => s ? {
          ...s,
          recentIncidents: [incident, ...s.recentIncidents].slice(0, 5)
        } : null);
      })
    );
  }

  updateIncident(studentId: string, incidentId: string, request: Partial<CreateIncidentRequest>): Observable<MedicalIncident> {
    return this.api.put<MedicalIncident>(`/students/${studentId}/health/incidents/${incidentId}`, request).pipe(
      tap((incident) => {
        this._healthSummary.update(s => s ? {
          ...s,
          recentIncidents: s.recentIncidents.map(i => i.id === incidentId ? incident : i)
        } : null);
      })
    );
  }

  deleteIncident(studentId: string, incidentId: string): Observable<void> {
    return this.api.delete<void>(`/students/${studentId}/health/incidents/${incidentId}`).pipe(
      tap(() => {
        this._healthSummary.update(s => s ? {
          ...s,
          recentIncidents: s.recentIncidents.filter(i => i.id !== incidentId)
        } : null);
      })
    );
  }

  // ==========================================================================
  // Reset
  // ==========================================================================

  reset(): void {
    this._healthSummary.set(null);
    this._loading.set(false);
    this._error.set(null);
  }
}
