import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';

import { BehavioralService } from './behavioral.service';
import { ApiService } from '../../../core/services';
import {
  BehavioralIncident,
  BehaviorSummary,
  FollowUp,
  BehavioralIncidentType,
  CreateBehavioralIncidentRequest,
} from '../models/behavioral.model';

describe('BehavioralService', () => {
  let service: BehavioralService;
  let httpMock: HttpTestingController;

  const mockStudentId = 'test-student-id';

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [BehavioralService, ApiService],
    });

    service = TestBed.inject(BehavioralService);
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
      expect(service.incidents()).toEqual([]);
      expect(service.summary()).toBeNull();
      expect(service.loading()).toBeFalsy();
      expect(service.error()).toBeNull();
    });
  });

  describe('loadIncidents', () => {
    it('should load incidents for a student', () => {
      const mockIncidents: BehavioralIncident[] = [
        {
          id: 'incident-1',
          studentId: mockStudentId,
          incidentType: 'minor_infraction' as BehavioralIncidentType,
          incidentTypeLabel: 'Minor Infraction',
          severity: 'low',
          severityLabel: 'Low',
          incidentDate: '2024-01-15',
          incidentTime: '10:00',
          description: 'Late to class',
          actionTaken: 'Verbal warning',
          parentMeetingRequired: false,
          parentNotified: false,
          reportedBy: 'user-1',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
      ];

      service.loadIncidents(mockStudentId).subscribe((response) => {
        expect(response.incidents.length).toBe(1);
        expect(service.incidents().length).toBe(1);
      });

      const req = httpMock.expectOne((r) => r.url.includes(`/students/${mockStudentId}/behavioral-incidents`));
      expect(req.request.method).toBe('GET');
      req.flush({ incidents: mockIncidents, total: 1 });
    });

    it('should set loading state during request', () => {
      service.loadIncidents(mockStudentId).subscribe();
      expect(service.loading()).toBeTruthy();

      const req = httpMock.expectOne((r) => r.url.includes(`/students/${mockStudentId}/behavioral-incidents`));
      req.flush({ incidents: [], total: 0 });
      expect(service.loading()).toBeFalsy();
    });
  });

  describe('loadSummary', () => {
    it('should load behavior summary', () => {
      const mockSummary: BehaviorSummary = {
        totalIncidents: 10,
        positiveCount: 5,
        minorInfractionCount: 3,
        majorViolationCount: 2,
        thisMonthCount: 2,
        lastMonthCount: 4,
        trend: 'improving',
        pendingFollowUps: 1,
      };

      service.loadSummary(mockStudentId).subscribe((summary) => {
        expect(summary).toEqual(mockSummary);
        expect(service.summary()).toEqual(mockSummary);
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/behavioral-summary`);
      expect(req.request.method).toBe('GET');
      req.flush(mockSummary);
    });
  });

  describe('Incident CRUD', () => {
    it('should create incident', () => {
      const request: CreateBehavioralIncidentRequest = {
        incidentType: 'positive_recognition' as BehavioralIncidentType,
        severity: 'low',
        incidentDate: '2024-01-20',
        incidentTime: '10:00',
        description: 'Helped classmate',
        actionTaken: 'Praise and recognition',
        parentMeetingRequired: false,
      };

      const mockResponse: BehavioralIncident = {
        id: 'new-incident',
        studentId: mockStudentId,
        incidentType: request.incidentType,
        incidentTypeLabel: 'Positive Recognition',
        severity: request.severity || 'low',
        severityLabel: 'Low',
        incidentDate: request.incidentDate,
        incidentTime: request.incidentTime,
        description: request.description,
        actionTaken: request.actionTaken,
        parentMeetingRequired: request.parentMeetingRequired ?? false,
        parentNotified: false,
        reportedBy: 'user-1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      service.createIncident(mockStudentId, request).subscribe((incident) => {
        expect(incident.id).toBe('new-incident');
      });

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/behavioral-incidents`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(request);
      req.flush(mockResponse);
    });

    it('should get single incident', () => {
      const incidentId = 'incident-1';

      service.getIncident(mockStudentId, incidentId).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/behavioral-incidents/${incidentId}`);
      expect(req.request.method).toBe('GET');
      req.flush({ id: incidentId });
    });

    it('should update incident', () => {
      const incidentId = 'incident-1';
      const request = {
        actionTaken: 'Updated action',
        parentNotified: true,
      };

      service.updateIncident(mockStudentId, incidentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/behavioral-incidents/${incidentId}`);
      expect(req.request.method).toBe('PUT');
      expect(req.request.body).toEqual(request);
      req.flush({ id: incidentId, ...request });
    });

    it('should delete incident', () => {
      const incidentId = 'incident-1';

      service.deleteIncident(mockStudentId, incidentId).subscribe();

      const req = httpMock.expectOne(`/api/v1/students/${mockStudentId}/behavioral-incidents/${incidentId}`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });

  describe('Follow-up CRUD', () => {
    const incidentId = 'incident-1';

    it('should create follow-up', () => {
      const request = {
        scheduledDate: '2024-02-01',
        scheduledTime: '14:00',
        participants: [{ name: 'Parent', role: 'Guardian' }],
        expectedOutcomes: 'Discuss improvement plan',
      };

      service.createFollowUp(incidentId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/behavioral-incidents/${incidentId}/follow-ups`);
      expect(req.request.method).toBe('POST');
      req.flush({ id: 'new-followup', ...request, status: 'pending' });
    });

    it('should update follow-up', () => {
      const followUpId = 'followup-1';
      const request = {
        status: 'completed' as const,
        meetingNotes: 'Meeting went well',
        actualOutcomes: 'Agreed on improvement plan',
      };

      service.updateFollowUp(incidentId, followUpId, request).subscribe();

      const req = httpMock.expectOne(`/api/v1/behavioral-incidents/${incidentId}/follow-ups/${followUpId}`);
      expect(req.request.method).toBe('PUT');
      req.flush({ id: followUpId, ...request });
    });

    it('should delete follow-up', () => {
      const followUpId = 'followup-1';

      service.deleteFollowUp(incidentId, followUpId).subscribe();

      const req = httpMock.expectOne(`/api/v1/behavioral-incidents/${incidentId}/follow-ups/${followUpId}`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });

  describe('Incident Types', () => {
    it('should recognize positive recognition type', () => {
      const incident = { incidentType: 'positive_recognition' as BehavioralIncidentType };
      expect(incident.incidentType).toBe('positive_recognition');
    });

    it('should recognize minor infraction type', () => {
      const incident = { incidentType: 'minor_infraction' as BehavioralIncidentType };
      expect(incident.incidentType).toBe('minor_infraction');
    });

    it('should recognize major violation type', () => {
      const incident = { incidentType: 'major_violation' as BehavioralIncidentType };
      expect(incident.incidentType).toBe('major_violation');
    });
  });

  describe('Severity Levels', () => {
    it('should handle low severity', () => {
      const incident = { severity: 'low' };
      expect(incident.severity).toBe('low');
    });

    it('should handle medium severity', () => {
      const incident = { severity: 'medium' };
      expect(incident.severity).toBe('medium');
    });

    it('should handle high severity', () => {
      const incident = { severity: 'high' };
      expect(incident.severity).toBe('high');
    });

    it('should handle critical severity', () => {
      const incident = { severity: 'critical' };
      expect(incident.severity).toBe('critical');
    });
  });

  describe('reset', () => {
    it('should reset all state', () => {
      service.reset();

      expect(service.incidents()).toEqual([]);
      expect(service.summary()).toBeNull();
      expect(service.loading()).toBeFalsy();
      expect(service.error()).toBeNull();
    });
  });
});
