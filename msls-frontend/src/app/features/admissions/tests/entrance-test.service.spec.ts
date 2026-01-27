/**
 * EntranceTestService Unit Tests
 */
import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { EntranceTestService } from './entrance-test.service';
import { ApiService } from '../../../core/services/api.service';
import { EntranceTest, CreateTestDto, TestStatus } from './entrance-test.model';

describe('EntranceTestService', () => {
  let service: EntranceTestService;
  let apiServiceMock: {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
  };

  const mockTest: EntranceTest = {
    id: 'test-uuid-1',
    tenantId: 'tenant-1',
    sessionId: 'session-uuid-1',
    testName: 'LKG Entrance Test',
    testDate: '2026-02-15',
    startTime: '10:00',
    durationMinutes: 60,
    venue: 'Main Hall',
    classNames: ['LKG', 'UKG'],
    maxCandidates: 50,
    status: 'scheduled' as TestStatus,
    subjects: [
      { name: 'English', maxMarks: 25 },
      { name: 'Math', maxMarks: 25 },
    ],
    registeredCount: 10,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  };

  const mockApiResponse = {
    id: 'test-uuid-1',
    sessionId: 'session-uuid-1',
    testName: 'LKG Entrance Test',
    testDate: '2026-02-15',
    startTime: '10:00',
    durationMinutes: 60,
    venue: 'Main Hall',
    classNames: ['LKG', 'UKG'],
    maxCandidates: 50,
    status: 'scheduled',
    subjects: [
      { subject: 'English', maxMarks: 25 },
      { subject: 'Math', maxMarks: 25 },
    ],
    registeredCount: 10,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  };

  beforeEach(() => {
    apiServiceMock = {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [
        EntranceTestService,
        { provide: ApiService, useValue: apiServiceMock },
      ],
    });
    service = TestBed.inject(EntranceTestService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getTests', () => {
    it('should fetch all tests', fakeAsync(() => {
      const mockResponse = { tests: [mockApiResponse], total: 1 };
      apiServiceMock.get.mockReturnValue(of(mockResponse));

      let result: EntranceTest[] | undefined;
      service.getTests().subscribe(tests => {
        result = tests;
      });
      tick();

      expect(result?.length).toBe(1);
      expect(result?.[0].testName).toBe('LKG Entrance Test');
      expect(result?.[0].subjects[0].name).toBe('English');
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/entrance-tests', expect.any(Object));
    }));

    it('should handle filters', fakeAsync(() => {
      const mockResponse = { tests: [mockApiResponse], total: 1 };
      apiServiceMock.get.mockReturnValue(of(mockResponse));

      service.getTests({ status: 'scheduled', sessionId: 'session-1' }).subscribe();
      tick();

      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/entrance-tests', {
        params: { status: 'scheduled', sessionId: 'session-1' },
      });
    }));

    it('should update loading state', fakeAsync(() => {
      const mockResponse = { tests: [], total: 0 };
      apiServiceMock.get.mockReturnValue(of(mockResponse));

      expect(service.loading()).toBe(false);

      service.getTests().subscribe();
      tick();

      expect(service.loading()).toBe(false);
    }));
  });

  describe('getTest', () => {
    it('should fetch a single test by ID', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockApiResponse));

      let result: EntranceTest | undefined;
      service.getTest('test-uuid-1').subscribe(test => {
        result = test;
      });
      tick();

      expect(result?.id).toBe('test-uuid-1');
      expect(result?.testName).toBe('LKG Entrance Test');
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/entrance-tests/test-uuid-1');
    }));

    it('should update selectedTest signal', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockApiResponse));

      service.getTest('test-uuid-1').subscribe();
      tick();

      expect(service.selectedTest()?.id).toBe('test-uuid-1');
    }));
  });

  describe('createTest', () => {
    it('should create a new test with correct payload transformation', fakeAsync(() => {
      const createDto: CreateTestDto = {
        testName: 'New Test',
        sessionId: 'session-uuid-1',
        testDate: '2026-03-01',
        startTime: '09:00',
        durationMinutes: 90,
        classNames: ['Class 1'],
        subjects: [{ name: 'Subject A', maxMarks: 50 }],
      };

      apiServiceMock.post.mockReturnValue(of(mockApiResponse));

      let result: EntranceTest | undefined;
      service.createTest(createDto).subscribe(test => {
        result = test;
      });
      tick();

      expect(result?.testName).toBe('LKG Entrance Test');
      expect(apiServiceMock.post).toHaveBeenCalledWith('/v1/entrance-tests', expect.objectContaining({
        testName: 'New Test',
        subjects: [{ subject: 'Subject A', maxMarks: 50 }],
      }));
    }));
  });

  describe('updateTest', () => {
    it('should update a test with correct payload transformation', fakeAsync(() => {
      const updateData = {
        testName: 'Updated Test',
        subjects: [{ name: 'Updated Subject', maxMarks: 30 }],
      };

      apiServiceMock.put.mockReturnValue(of(mockApiResponse));

      service.updateTest('test-uuid-1', updateData).subscribe();
      tick();

      expect(apiServiceMock.put).toHaveBeenCalledWith('/v1/entrance-tests/test-uuid-1', expect.objectContaining({
        testName: 'Updated Test',
        subjects: [{ subject: 'Updated Subject', maxMarks: 30 }],
      }));
    }));
  });

  describe('deleteTest', () => {
    it('should delete a test', fakeAsync(() => {
      apiServiceMock.delete.mockReturnValue(of(null));

      let completed = false;
      service.deleteTest('test-uuid-1').subscribe(() => {
        completed = true;
      });
      tick();

      expect(completed).toBe(true);
      expect(apiServiceMock.delete).toHaveBeenCalledWith('/v1/entrance-tests/test-uuid-1');
    }));

    it('should remove deleted test from tests signal', fakeAsync(() => {
      // First, populate tests
      const mockResponse = { tests: [mockApiResponse], total: 1 };
      apiServiceMock.get.mockReturnValue(of(mockResponse));
      service.getTests().subscribe();
      tick();

      expect(service.tests().length).toBe(1);

      // Then delete
      apiServiceMock.delete.mockReturnValue(of(null));
      service.deleteTest('test-uuid-1').subscribe();
      tick();

      expect(service.tests().length).toBe(0);
    }));
  });

  describe('completeTest', () => {
    it('should mark test as completed', fakeAsync(() => {
      apiServiceMock.put.mockReturnValue(of({ ...mockApiResponse, status: 'completed' }));

      service.completeTest('test-uuid-1').subscribe();
      tick();

      expect(apiServiceMock.put).toHaveBeenCalledWith('/v1/entrance-tests/test-uuid-1', { status: 'completed' });
    }));
  });

  describe('cancelTest', () => {
    it('should mark test as cancelled', fakeAsync(() => {
      apiServiceMock.put.mockReturnValue(of({ ...mockApiResponse, status: 'cancelled' }));

      service.cancelTest('test-uuid-1').subscribe();
      tick();

      expect(apiServiceMock.put).toHaveBeenCalledWith('/v1/entrance-tests/test-uuid-1', { status: 'cancelled' });
    }));
  });

  describe('subject field transformation', () => {
    it('should transform backend "subject" to frontend "name" on response', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockApiResponse));

      let result: EntranceTest | undefined;
      service.getTest('test-uuid-1').subscribe(test => {
        result = test;
      });
      tick();

      // Backend returns 'subject', frontend should have 'name'
      expect(result?.subjects[0].name).toBe('English');
      expect((result?.subjects[0] as any).subject).toBeUndefined();
    }));
  });
});
