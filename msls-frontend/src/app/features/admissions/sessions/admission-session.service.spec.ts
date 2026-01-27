/**
 * AdmissionSessionService Unit Tests
 */
import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { AdmissionSessionService } from './admission-session.service';
import { ApiService } from '../../../core/services/api.service';
import { AdmissionSession, CreateSessionRequest, SeatConfigRequest } from './admission-session.model';

describe('AdmissionSessionService', () => {
  let service: AdmissionSessionService;
  let apiServiceMock: {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    patch: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
  };

  const mockApiSession = {
    id: 'session-uuid-1',
    branchId: 'branch-1',
    academicYearId: 'ay-uuid-1',
    name: 'Admission 2026',
    description: 'Admissions for academic year 2026',
    startDate: '2026-01-01',
    endDate: '2026-03-31',
    status: 'open',
    applicationFee: 500,
    requiredDocuments: ['birth_certificate', 'photo'],
    settings: {
      allowOnlineApplication: true,
      notifyOnApplication: true,
    },
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    stats: {
      totalApplications: 150,
      approvedCount: 50,
      pendingCount: 80,
      rejectedCount: 20,
      totalSeats: 200,
      filledSeats: 50,
      availableSeats: 150,
    },
  };

  const mockApiSeat = {
    id: 'seat-uuid-1',
    sessionId: 'session-uuid-1',
    className: 'Class 1',
    totalSeats: 40,
    filledSeats: 15,
    availableSeats: 25,
    waitlistLimit: 10,
    reservedSeats: { SC: 10, ST: 5 },
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  };

  const mockAcademicYear = {
    id: 'ay-uuid-1',
    name: '2025-2026',
    startDate: '2025-04-01',
    endDate: '2026-03-31',
    isCurrent: true,
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
        AdmissionSessionService,
        { provide: ApiService, useValue: apiServiceMock },
      ],
    });
    service = TestBed.inject(AdmissionSessionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getSessions', () => {
    it('should fetch all sessions with academic year names', fakeAsync(() => {
      const sessionsResponse = { sessions: [mockApiSession], total: 1 };
      const academicYearsResponse = { academicYears: [mockAcademicYear], total: 1 };

      apiServiceMock.get.mockImplementation((url: string) => {
        if (url.includes('academic-years')) {
          return of(academicYearsResponse);
        }
        return of(sessionsResponse);
      });

      let result: AdmissionSession[] | undefined;
      service.getSessions().subscribe(sessions => {
        result = sessions;
      });
      tick();

      expect(result?.length).toBe(1);
      expect(result?.[0].name).toBe('Admission 2026');
      expect(result?.[0].academicYearName).toBe('2025-2026');
    }));

    it('should map API response to frontend model', fakeAsync(() => {
      const sessionsResponse = { sessions: [mockApiSession], total: 1 };
      const academicYearsResponse = { academicYears: [mockAcademicYear], total: 1 };

      apiServiceMock.get.mockImplementation((url: string) => {
        if (url.includes('academic-years')) {
          return of(academicYearsResponse);
        }
        return of(sessionsResponse);
      });

      let result: AdmissionSession[] | undefined;
      service.getSessions().subscribe(sessions => {
        result = sessions;
      });
      tick();

      const session = result?.[0];
      expect(session?.id).toBe('session-uuid-1');
      expect(session?.status).toBe('open');
      expect(session?.applicationFee).toBe(500);
      expect(session?.requiredDocuments).toContain('birth_certificate');
      expect(session?.totalApplications).toBe(150);
    }));
  });

  describe('getSession', () => {
    it('should fetch a single session by ID', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockApiSession));

      let result: AdmissionSession | undefined;
      service.getSession('session-uuid-1').subscribe(session => {
        result = session;
      });
      tick();

      expect(result?.id).toBe('session-uuid-1');
      expect(result?.name).toBe('Admission 2026');
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1');
    }));
  });

  describe('createSession', () => {
    it('should create a new session', fakeAsync(() => {
      const createRequest: CreateSessionRequest = {
        name: 'New Session',
        academicYearId: 'ay-uuid-1',
        startDate: '2026-04-01',
        endDate: '2026-06-30',
        applicationFee: 600,
        requiredDocuments: ['photo'],
        settings: { allowOnlineApplication: true },
      };

      apiServiceMock.post.mockReturnValue(of(mockApiSession));

      let result: AdmissionSession | undefined;
      service.createSession(createRequest).subscribe(session => {
        result = session;
      });
      tick();

      expect(result?.name).toBe('Admission 2026');
      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/admission-sessions', createRequest);
    }));
  });

  describe('updateSession', () => {
    it('should update an existing session', fakeAsync(() => {
      const updateRequest = {
        name: 'Updated Session',
        applicationFee: 700,
      };

      apiServiceMock.put.mockReturnValue(of({ ...mockApiSession, name: 'Updated Session' }));

      let result: AdmissionSession | undefined;
      service.updateSession('session-uuid-1', updateRequest).subscribe(session => {
        result = session;
      });
      tick();

      expect(result?.id).toBe('session-uuid-1');
      expect(apiServiceMock.put).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1', updateRequest);
    }));
  });

  describe('changeStatus', () => {
    it('should change session status to open', fakeAsync(() => {
      apiServiceMock.patch.mockReturnValue(of(mockApiSession));

      let result: AdmissionSession | undefined;
      service.changeStatus('session-uuid-1', { status: 'open' }).subscribe(session => {
        result = session;
      });
      tick();

      expect(result?.status).toBe('open');
      expect(apiServiceMock.patch).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/status', {
        status: 'open',
      });
    }));

    it('should change session status to closed', fakeAsync(() => {
      apiServiceMock.patch.mockReturnValue(of({ ...mockApiSession, status: 'closed' }));

      service.changeStatus('session-uuid-1', { status: 'closed' }).subscribe();
      tick();

      expect(apiServiceMock.patch).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/status', {
        status: 'closed',
      });
    }));
  });

  describe('deleteSession', () => {
    it('should delete a session', fakeAsync(() => {
      apiServiceMock.delete.mockReturnValue(of(null));

      let completed = false;
      service.deleteSession('session-uuid-1').subscribe(() => {
        completed = true;
      });
      tick();

      expect(completed).toBe(true);
      expect(apiServiceMock.delete).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1');
    }));
  });

  describe('seat operations', () => {
    describe('getSeats', () => {
      it('should fetch seats for a session', fakeAsync(() => {
        const response = { seats: [mockApiSeat], total: 1 };
        apiServiceMock.get.mockReturnValue(of(response));

        let result: any;
        service.getSeats('session-uuid-1').subscribe(seats => {
          result = seats;
        });
        tick();

        expect(result.length).toBe(1);
        expect(result[0].className).toBe('Class 1');
        expect(result[0].totalSeats).toBe(40);
        expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/seats');
      }));
    });

    describe('addSeat', () => {
      it('should add seat configuration', fakeAsync(() => {
        const seatRequest: SeatConfigRequest = {
          className: 'Class 2',
          totalSeats: 35,
          waitlistLimit: 5,
          reservedSeats: { SC: 8 },
        };

        apiServiceMock.post.mockReturnValue(of(mockApiSeat));

        let result: any;
        service.addSeat('session-uuid-1', seatRequest).subscribe(seat => {
          result = seat;
        });
        tick();

        expect(result.className).toBe('Class 1');
        expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/admission-sessions/session-uuid-1/seats', seatRequest);
      }));
    });

    describe('updateSeat', () => {
      it('should update seat configuration', fakeAsync(() => {
        const seatRequest: SeatConfigRequest = {
          className: 'Class 1',
          totalSeats: 45,
          waitlistLimit: 12,
        };

        apiServiceMock.put.mockReturnValue(of(mockApiSeat));

        let result: any;
        service.updateSeat('session-uuid-1', 'seat-uuid-1', seatRequest).subscribe(seat => {
          result = seat;
        });
        tick();

        expect(result.totalSeats).toBe(40);
        expect(apiServiceMock.put).toHaveBeenCalledWith(
          '/v1/admission-sessions/session-uuid-1/seats/seat-uuid-1',
          seatRequest
        );
      }));
    });

    describe('deleteSeat', () => {
      it('should delete seat configuration', fakeAsync(() => {
        apiServiceMock.delete.mockReturnValue(of(null));

        let completed = false;
        service.deleteSeat('session-uuid-1', 'seat-uuid-1').subscribe(() => {
          completed = true;
        });
        tick();

        expect(completed).toBe(true);
        expect(apiServiceMock.delete).toHaveBeenCalledWith(
          '/v1/admission-sessions/session-uuid-1/seats/seat-uuid-1'
        );
      }));
    });
  });

  describe('getAcademicYears', () => {
    it('should fetch academic years', fakeAsync(() => {
      const response = { academicYears: [mockAcademicYear], total: 1 };
      apiServiceMock.get.mockReturnValue(of(response));

      let result: any;
      service.getAcademicYears().subscribe(years => {
        result = years;
      });
      tick();

      expect(result.length).toBe(1);
      expect(result[0].name).toBe('2025-2026');
      expect(result[0].isCurrent).toBe(true);
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/academic-years');
    }));
  });
});
