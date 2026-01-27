/**
 * MeritService Unit Tests
 */
import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { MeritService } from './merit.service';
import { ApiService } from '../../../core/services/api.service';
import { MeritList, DecisionType } from './merit.model';

describe('MeritService', () => {
  let service: MeritService;
  let apiServiceMock: {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    patch: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
  };

  const mockMeritListEntry = {
    rank: 1,
    applicationId: 'app-uuid-1',
    studentName: 'Test Student',
    score: 85,
    testScore: 85,
    status: 'test_completed',
    parentPhone: '9876543210',
    parentEmail: 'parent@test.com',
  };

  const mockMeritListApiResponse = {
    id: 'merit-uuid-1',
    sessionId: 'session-uuid-1',
    className: 'Class 1',
    testId: 'test-uuid-1',
    generatedAt: '2026-01-15T10:00:00Z',
    cutoffScore: 50,
    entries: [mockMeritListEntry],
    isFinal: false,
    totalCount: 1,
    aboveCutoff: 1,
    createdAt: '2026-01-15T10:00:00Z',
  };

  const mockDecisionApiResponse = {
    id: 'decision-uuid-1',
    applicationId: 'app-uuid-1',
    decision: 'approved',
    decisionDate: '2026-01-20',
    sectionAssigned: 'Section A',
    offerValidUntil: '2026-02-20',
    createdAt: '2026-01-20T10:00:00Z',
    updatedAt: '2026-01-20T10:00:00Z',
  };

  beforeEach(() => {
    apiServiceMock = {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [
        MeritService,
        { provide: ApiService, useValue: apiServiceMock },
      ],
    });
    service = TestBed.inject(MeritService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getMeritList', () => {
    it('should fetch merit list for session and class', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockMeritListApiResponse));

      let result: MeritList | null | undefined;
      service.getMeritList('session-uuid-1', 'Class 1').subscribe(list => {
        result = list;
      });
      tick();

      expect(result).toBeTruthy();
      expect(result?.className).toBe('Class 1');
      expect(result?.entries.length).toBe(1);
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/merit-list', {
        params: { className: 'Class 1' },
      });
    }));

    it('should map API response to frontend model', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockMeritListApiResponse));

      let result: MeritList | null | undefined;
      service.getMeritList('session-uuid-1', 'Class 1').subscribe(list => {
        result = list;
      });
      tick();

      expect(result?.id).toBe('merit-uuid-1');
      expect(result?.sessionId).toBe('session-uuid-1');
      expect(result?.cutoffScore).toBe(50);
      expect(result?.isFinal).toBe(false);
    }));
  });

  describe('listMeritLists', () => {
    it('should fetch all merit lists for a session', fakeAsync(() => {
      const mockResponse = {
        meritLists: [mockMeritListApiResponse],
        total: 1,
      };
      apiServiceMock.get.mockReturnValue(of(mockResponse));

      let result: MeritList[] | undefined;
      service.listMeritLists('session-uuid-1').subscribe(lists => {
        result = lists;
      });
      tick();

      expect(result?.length).toBe(1);
      expect(result?.[0].className).toBe('Class 1');
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/merit-lists');
    }));
  });

  describe('generateMeritList', () => {
    it('should generate merit list with request parameters', fakeAsync(() => {
      const request = {
        className: 'Class 1',
        testId: 'test-uuid-1',
        cutoffScore: 50,
      };
      apiServiceMock.post.mockReturnValue(of(mockMeritListApiResponse));

      let result: MeritList | undefined;
      service.generateMeritList('session-uuid-1', request).subscribe(list => {
        result = list;
      });
      tick();

      expect(result?.className).toBe('Class 1');
      expect(apiServiceMock.post).toHaveBeenCalledWith(
        '/v1/admission-sessions/session-uuid-1/merit-list',
        request
      );
    }));
  });

  describe('finalizeMeritList', () => {
    it('should finalize a merit list', fakeAsync(() => {
      const finalizedResponse = { ...mockMeritListApiResponse, isFinal: true };
      apiServiceMock.post.mockReturnValue(of(finalizedResponse));

      let result: MeritList | undefined;
      service.finalizeMeritList('merit-uuid-1').subscribe(list => {
        result = list;
      });
      tick();

      expect(result?.isFinal).toBe(true);
      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/merit-lists/merit-uuid-1/finalize', {});
    }));
  });

  describe('updateCutoff', () => {
    it('should update cutoff score', fakeAsync(() => {
      const updatedResponse = { ...mockMeritListApiResponse, cutoffScore: 60 };
      apiServiceMock.patch.mockReturnValue(of(updatedResponse));

      let result: MeritList | undefined;
      service.updateCutoff('merit-uuid-1', 60).subscribe(list => {
        result = list;
      });
      tick();

      expect(result?.cutoffScore).toBe(60);
      expect(apiServiceMock.patch).toHaveBeenCalledWith('/v1/merit-lists/merit-uuid-1/cutoff', {
        cutoffScore: 60,
      });
    }));

    it('should allow null cutoff score', fakeAsync(() => {
      const updatedResponse = { ...mockMeritListApiResponse, cutoffScore: null };
      apiServiceMock.patch.mockReturnValue(of(updatedResponse));

      let result: MeritList | undefined;
      service.updateCutoff('merit-uuid-1', null).subscribe(list => {
        result = list;
      });
      tick();

      expect(result?.cutoffScore).toBeNull();
      expect(apiServiceMock.patch).toHaveBeenCalledWith('/v1/merit-lists/merit-uuid-1/cutoff', {
        cutoffScore: null,
      });
    }));
  });

  describe('makeDecision', () => {
    it('should make admission decision for selected', fakeAsync(() => {
      const request = {
        decision: 'selected' as DecisionType,
        sectionAssigned: 'Section A',
        offerValidUntil: '2026-02-20',
      };
      apiServiceMock.post.mockReturnValue(of(mockDecisionApiResponse));

      service.makeDecision('app-uuid-1', request).subscribe(decision => {
        expect(decision.applicationId).toBe('app-uuid-1');
        expect(decision.decision).toBe('selected');
      });
      tick();

      // Backend expects 'approved' not 'selected'
      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/decision', {
        decision: 'approved',
        sectionAssigned: 'Section A',
        offerValidUntil: '2026-02-20',
      });
    }));

    it('should make admission decision for waitlisted', fakeAsync(() => {
      const request = {
        decision: 'waitlisted' as DecisionType,
        waitlistPosition: 5,
      };
      apiServiceMock.post.mockReturnValue(of({ ...mockDecisionApiResponse, decision: 'waitlisted' }));

      service.makeDecision('app-uuid-1', request).subscribe();
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/decision', {
        decision: 'waitlisted',
        waitlistPosition: 5,
      });
    }));

    it('should make admission decision for rejected', fakeAsync(() => {
      const request = {
        decision: 'rejected' as DecisionType,
        rejectionReason: 'Score below cutoff',
      };
      apiServiceMock.post.mockReturnValue(of({ ...mockDecisionApiResponse, decision: 'rejected' }));

      service.makeDecision('app-uuid-1', request).subscribe();
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/decision', {
        decision: 'rejected',
        rejectionReason: 'Score below cutoff',
      });
    }));
  });

  describe('makeBulkDecision', () => {
    it('should make bulk admission decisions', fakeAsync(() => {
      const request = {
        applicationIds: ['app-1', 'app-2', 'app-3'],
        decision: 'selected' as DecisionType,
        sectionAssigned: 'Section A',
      };
      const mockResponse = { successful: 3, failed: 0, decisions: [], errors: [] };
      apiServiceMock.post.mockReturnValue(of(mockResponse));

      service.makeBulkDecision(request).subscribe(response => {
        expect(response.successful).toBeGreaterThan(0);
      });
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/bulk-decision', {
        applicationIds: ['app-1', 'app-2', 'app-3'],
        decision: 'approved',
        sectionAssigned: 'Section A',
      });
    }));
  });

  describe('generateOfferLetter', () => {
    it('should generate offer letter', fakeAsync(() => {
      const mockResponse = {
        url: '/api/v1/applications/app-uuid-1/offer-letter.pdf',
        validUntil: '2026-02-20',
        generatedAt: '2026-01-20T10:00:00Z',
        applicationId: 'app-uuid-1',
      };
      apiServiceMock.post.mockReturnValue(of(mockResponse));

      service.generateOfferLetter('app-uuid-1').subscribe(result => {
        expect(result.url).toContain('offer-letter.pdf');
        expect(result.validUntil).toBe('2026-02-20');
      });
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/offer-letter', {});
    }));
  });

  describe('acceptOffer', () => {
    it('should accept an offer', fakeAsync(() => {
      const request = {
        paymentReference: 'PAY-123',
        remarks: 'Payment completed',
      };
      apiServiceMock.post.mockReturnValue(of({ ...mockDecisionApiResponse, offerAccepted: true }));

      service.acceptOffer('app-uuid-1', request).subscribe(decision => {
        expect(decision.offerAccepted).toBe(true);
      });
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/accept-offer', request);
    }));
  });

  describe('enroll', () => {
    it('should complete enrollment', fakeAsync(() => {
      const request = {
        sectionId: 'section-1',
        rollNumber: '101',
        admissionDate: '2026-04-01',
      };

      const mockResponse = {
        applicationId: 'app-uuid-1',
        studentId: 'student-uuid-1',
        enrollmentNumber: 'ENR-2026-0001',
        admissionDate: '2026-04-01',
        className: 'Class 1',
        sectionName: 'Section A',
        rollNumber: '101',
        status: 'enrolled',
      };
      apiServiceMock.post.mockReturnValue(of(mockResponse));

      service.enroll('app-uuid-1', request).subscribe(result => {
        expect(result.studentId).toBe('student-uuid-1');
        expect(result.enrollmentNumber).toBe('ENR-2026-0001');
      });
      tick();

      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/applications/app-uuid-1/enroll', request);
    }));
  });

  describe('getAvailableClasses', () => {
    it('should return available classes', fakeAsync(() => {
      // The service uses mock data with delay
      let result: string[] | undefined;
      service.getAvailableClasses('session-uuid-1').subscribe(classes => {
        result = classes;
      });

      tick(300); // Delay in service

      expect(result?.length).toBeGreaterThan(0);
      expect(result).toContain('LKG');
    }));
  });

  describe('getSections', () => {
    it('should return available sections for a class', fakeAsync(() => {
      // The service uses mock data with delay
      let result: Array<{ id: string; name: string }> | undefined;
      service.getSections('Class 1').subscribe(sections => {
        result = sections;
      });

      tick(300); // Delay in service

      expect(result?.length).toBeGreaterThan(0);
      expect(result?.[0].name).toBe('Section A');
    }));
  });
});
