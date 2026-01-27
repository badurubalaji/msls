import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';

import { HealthService } from './health.service';
import { ApiService } from '../../../core/services';
import {
  HealthSummary,
  HealthProfile,
  Allergy,
  ChronicCondition,
  Medication,
  Vaccination,
  MedicalIncident,
  CreateHealthProfileRequest,
  CreateAllergyRequest,
} from '../models/health.model';

describe('HealthService', () => {
  let service: HealthService;
  let httpMock: HttpTestingController;

  const mockStudentId = 'test-student-id';

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [HealthService, ApiService],
    });

    service = TestBed.inject(HealthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('Service Creation', () => {
    it('should be created', () => {
      expect(service).toBeTruthy();
    });

    it('should have initial state', () => {
      expect(service.healthSummary()).toBeNull();
      expect(service.loading()).toBeFalsy();
      expect(service.error()).toBeNull();
    });
  });

  describe('loadHealthSummary', () => {
    it('should load health summary for a student', () => {
      const mockSummary: HealthSummary = {
        profile: undefined,
        allergies: [],
        conditions: [],
        medications: [],
        vaccinations: [],
        recentIncidents: [],
      };

      service.loadHealthSummary(mockStudentId).subscribe((summary) => {
        expect(summary).toEqual(mockSummary);
        expect(service.healthSummary()).toEqual(mockSummary);
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health`);
      expect(req.request.method).toBe('GET');
      req.flush(mockSummary);
    });

    it('should set loading state during request', () => {
      service.loadHealthSummary(mockStudentId).subscribe();
      expect(service.loading()).toBeTruthy();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health`);
      req.flush({});
      expect(service.loading()).toBeFalsy();
    });
  });

  describe('Health Profile', () => {
    it('should get health profile', () => {
      const mockProfile: HealthProfile = {
        id: 'profile-1',
        studentId: mockStudentId,
        bloodGroup: 'A+',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      service.getHealthProfile(mockStudentId).subscribe((profile) => {
        expect(profile).toEqual(mockProfile);
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/profile`);
      expect(req.request.method).toBe('GET');
      req.flush(mockProfile);
    });

    it('should save health profile', () => {
      const request: CreateHealthProfileRequest = { bloodGroup: 'O+', heightCm: 150, weightKg: 45 };
      const mockProfile: HealthProfile = {
        id: 'profile-1',
        studentId: mockStudentId,
        bloodGroup: 'O+',
        heightCm: 150,
        weightKg: 45,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      service.saveHealthProfile(mockStudentId, request).subscribe((profile) => {
        expect(profile).toEqual(mockProfile);
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/profile`);
      expect(req.request.method).toBe('PUT');
      expect(req.request.body).toEqual(request);
      req.flush(mockProfile);
    });
  });

  describe('Allergies', () => {
    it('should list allergies', () => {
      const mockResponse = {
        allergies: [
          { id: 'allergy-1', allergen: 'Peanuts', severity: 'severe' },
        ],
        total: 1,
      };

      service.listAllergies(mockStudentId).subscribe((response) => {
        expect(response.allergies.length).toBe(1);
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/allergies`);
      expect(req.request.method).toBe('GET');
      req.flush(mockResponse);
    });

    it('should list active allergies only', () => {
      service.listAllergies(mockStudentId, true).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/allergies?active=true`);
      expect(req.request.method).toBe('GET');
      req.flush({ allergies: [], total: 0 });
    });

    it('should create allergy', () => {
      const request: CreateAllergyRequest = {
        allergen: 'Peanuts',
        allergyType: 'food',
        severity: 'severe',
      };

      service.createAllergy(mockStudentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/allergies`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(request);
      req.flush({ id: 'new-allergy', ...request });
    });

    it('should delete allergy', () => {
      const allergyId = 'allergy-1';

      service.deleteAllergy(mockStudentId, allergyId).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/allergies/${allergyId}`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });

  describe('Conditions', () => {
    it('should create condition', () => {
      const request = {
        conditionName: 'Asthma',
        conditionType: 'respiratory' as const,
        severity: 'moderate' as const,
      };

      service.createCondition(mockStudentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/conditions`);
      expect(req.request.method).toBe('POST');
      req.flush({ id: 'new-condition', ...request });
    });
  });

  describe('Medications', () => {
    it('should create medication', () => {
      const request = {
        medicationName: 'Ventolin',
        dosage: '2 puffs',
        frequency: 'as_needed' as const,
        route: 'inhaler' as const,
        startDate: '2024-01-01',
        administeredAtSchool: false,
      };

      service.createMedication(mockStudentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/medications`);
      expect(req.request.method).toBe('POST');
      req.flush({ id: 'new-medication', ...request });
    });
  });

  describe('Vaccinations', () => {
    it('should list vaccinations', () => {
      service.listVaccinations(mockStudentId).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/vaccinations`);
      expect(req.request.method).toBe('GET');
      req.flush({ vaccinations: [], total: 0 });
    });

    it('should create vaccination', () => {
      const request = {
        vaccineName: 'MMR',
        doseNumber: 1,
        administeredDate: '2024-01-15',
        hadReaction: false,
      };

      service.createVaccination(mockStudentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/vaccinations`);
      expect(req.request.method).toBe('POST');
      req.flush({ id: 'new-vaccination', ...request });
    });
  });

  describe('Medical Incidents', () => {
    it('should list incidents', () => {
      service.listIncidents(mockStudentId).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/incidents`);
      expect(req.request.method).toBe('GET');
      req.flush({ incidents: [], total: 0 });
    });

    it('should list incidents with limit', () => {
      service.listIncidents(mockStudentId, 5).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/incidents?limit=5`);
      expect(req.request.method).toBe('GET');
      req.flush({ incidents: [], total: 0 });
    });

    it('should create incident', () => {
      const request = {
        incidentDate: '2024-01-20',
        incidentTime: '10:30',
        incidentType: 'injury' as const,
        description: 'Fall on playground',
        actionTaken: 'Applied ice pack',
        firstAidGiven: true,
        parentNotified: false,
        hospitalVisitRequired: false,
        studentSentHome: false,
        followUpRequired: false,
      };

      service.createIncident(mockStudentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/health/incidents`);
      expect(req.request.method).toBe('POST');
      req.flush({ id: 'new-incident', ...request });
    });
  });

  describe('reset', () => {
    it('should reset all state', () => {
      service.reset();

      expect(service.healthSummary()).toBeNull();
      expect(service.loading()).toBeFalsy();
      expect(service.error()).toBeNull();
    });
  });
});
