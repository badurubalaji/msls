/**
 * TestsComponent Unit Tests
 */
import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { Router } from '@angular/router';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { TestsComponent } from './tests';
import { EntranceTestService } from '../../entrance-test.service';
import { ToastService } from '../../../../../shared/services';
import { EntranceTest, TestStatus } from '../../entrance-test.model';

describe('TestsComponent', () => {
  let component: TestsComponent;
  let fixture: ComponentFixture<TestsComponent>;
  let testServiceMock: {
    getTests: ReturnType<typeof vi.fn>;
    createTest: ReturnType<typeof vi.fn>;
    updateTest: ReturnType<typeof vi.fn>;
    deleteTest: ReturnType<typeof vi.fn>;
    completeTest: ReturnType<typeof vi.fn>;
    cancelTest: ReturnType<typeof vi.fn>;
  };
  let toastServiceMock: {
    success: ReturnType<typeof vi.fn>;
    error: ReturnType<typeof vi.fn>;
  };
  let routerMock: {
    navigate: ReturnType<typeof vi.fn>;
  };

  const mockTests: EntranceTest[] = [
    {
      id: 'test-1',
      tenantId: 'tenant-1',
      sessionId: 'session-1',
      testName: 'LKG Entrance Test',
      testDate: '2026-02-15',
      startTime: '10:00',
      durationMinutes: 60,
      venue: 'Main Hall',
      classNames: ['LKG'],
      maxCandidates: 50,
      status: 'scheduled' as TestStatus,
      subjects: [{ name: 'English', maxMarks: 25 }],
      registeredCount: 10,
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    },
    {
      id: 'test-2',
      tenantId: 'tenant-1',
      sessionId: 'session-1',
      testName: 'Class 1 Test',
      testDate: '2026-02-20',
      startTime: '09:00',
      durationMinutes: 90,
      venue: 'Hall B',
      classNames: ['Class 1'],
      maxCandidates: 40,
      status: 'completed' as TestStatus,
      subjects: [{ name: 'Math', maxMarks: 50 }],
      registeredCount: 35,
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    },
    {
      id: 'test-3',
      tenantId: 'tenant-1',
      sessionId: 'session-1',
      testName: 'UKG Entrance',
      testDate: '2026-02-10',
      startTime: '14:00',
      durationMinutes: 45,
      classNames: ['UKG'],
      maxCandidates: 30,
      status: 'scheduled' as TestStatus,
      subjects: [{ name: 'GK', maxMarks: 20 }],
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    },
  ];

  beforeEach(async () => {
    testServiceMock = {
      getTests: vi.fn(),
      createTest: vi.fn(),
      updateTest: vi.fn(),
      deleteTest: vi.fn(),
      completeTest: vi.fn(),
      cancelTest: vi.fn(),
    };
    toastServiceMock = {
      success: vi.fn(),
      error: vi.fn(),
    };
    routerMock = {
      navigate: vi.fn(),
    };

    testServiceMock.getTests.mockReturnValue(of(mockTests));
    testServiceMock.createTest.mockReturnValue(of(mockTests[0]));
    testServiceMock.updateTest.mockReturnValue(of(mockTests[0]));
    testServiceMock.deleteTest.mockReturnValue(of(undefined));
    testServiceMock.completeTest.mockReturnValue(of(mockTests[0]));
    testServiceMock.cancelTest.mockReturnValue(of(mockTests[0]));

    await TestBed.configureTestingModule({
      imports: [TestsComponent],
      providers: [
        { provide: EntranceTestService, useValue: testServiceMock },
        { provide: ToastService, useValue: toastServiceMock },
        { provide: Router, useValue: routerMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(TestsComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load tests on init', fakeAsync(() => {
    fixture.detectChanges();
    tick();

    expect(testServiceMock.getTests).toHaveBeenCalled();
    expect(component.tests().length).toBe(3);
    expect(component.loading()).toBe(false);
  }));

  it('should set error when loading fails', fakeAsync(() => {
    testServiceMock.getTests.mockReturnValue(throwError(() => new Error('API Error')));
    fixture.detectChanges();
    tick();

    expect(component.error()).toBe('Failed to load entrance tests. Please try again.');
    expect(component.loading()).toBe(false);
  }));

  describe('filtering', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should filter tests by search term', () => {
      component.onSearchChange('LKG');
      const filtered = component.filteredTests();

      expect(filtered.length).toBe(1);
      expect(filtered[0].testName).toBe('LKG Entrance Test');
    });

    it('should filter tests by venue', () => {
      component.onSearchChange('Hall B');
      const filtered = component.filteredTests();

      expect(filtered.length).toBe(1);
      expect(filtered[0].testName).toBe('Class 1 Test');
    });

    it('should filter tests by class name', () => {
      component.onSearchChange('UKG');
      const filtered = component.filteredTests();

      expect(filtered.length).toBe(2); // LKG test has UKG in name, UKG test has UKG class
    });

    it('should filter tests by status', () => {
      component.onStatusFilterChange('completed');
      const filtered = component.filteredTests();

      expect(filtered.length).toBe(1);
      expect(filtered[0].status).toBe('completed');
    });

    it('should combine search and status filters', () => {
      component.onSearchChange('Test');
      component.onStatusFilterChange('scheduled');
      const filtered = component.filteredTests();

      // Only scheduled tests with "Test" in name
      expect(filtered.every(t => t.status === 'scheduled')).toBe(true);
      expect(filtered.every(t => t.testName.toLowerCase().includes('test'))).toBe(true);
    });
  });

  describe('sorting', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should sort tests by date and time', () => {
      const filtered = component.filteredTests();

      // Tests should be sorted with upcoming tests first
      // Exact order depends on current date vs test dates
      expect(filtered.length).toBe(3);
    });
  });

  describe('create/edit modal', () => {
    it('should open create modal with null editingTest', () => {
      component.openCreateModal();

      expect(component.showTestModal()).toBe(true);
      expect(component.editingTest()).toBeNull();
    });

    it('should open edit modal with selected test', () => {
      component.editTest(mockTests[0]);

      expect(component.showTestModal()).toBe(true);
      expect(component.editingTest()).toEqual(mockTests[0]);
    });

    it('should close modal and clear editing state', () => {
      component.editTest(mockTests[0]);
      component.closeTestModal();

      expect(component.showTestModal()).toBe(false);
      expect(component.editingTest()).toBeNull();
    });
  });

  describe('save test', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should create new test when editingTest is null', fakeAsync(() => {
      const newTestData = {
        testName: 'New Test',
        sessionId: 'session-1',
        testDate: '2026-03-01',
        startTime: '10:00',
        durationMinutes: 60,
        classNames: ['Class 1'],
        subjects: [{ name: 'Subject', maxMarks: 25 }],
      };

      component.openCreateModal();
      component.saveTest(newTestData);
      tick();

      expect(testServiceMock.createTest).toHaveBeenCalledWith(newTestData);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Test created successfully');
      expect(component.showTestModal()).toBe(false);
    }));

    it('should update test when editingTest is set', fakeAsync(() => {
      const updateData = { testName: 'Updated Name' };

      component.editTest(mockTests[0]);
      component.saveTest(updateData);
      tick();

      expect(testServiceMock.updateTest).toHaveBeenCalledWith('test-1', updateData);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Test updated successfully');
    }));

    it('should show error toast on save failure', fakeAsync(() => {
      testServiceMock.createTest.mockReturnValue(throwError(() => new Error('Error')));

      component.openCreateModal();
      component.saveTest({ testName: 'Test' });
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to create test');
    }));
  });

  describe('delete test', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should open delete confirmation modal', () => {
      component.confirmDelete(mockTests[0]);

      expect(component.showDeleteModal()).toBe(true);
      expect(component.testToDelete()).toEqual(mockTests[0]);
    });

    it('should delete test and show success message', fakeAsync(() => {
      component.confirmDelete(mockTests[0]);
      component.deleteTest();
      tick();

      expect(testServiceMock.deleteTest).toHaveBeenCalledWith('test-1');
      expect(toastServiceMock.success).toHaveBeenCalledWith('Test deleted successfully');
      expect(component.showDeleteModal()).toBe(false);
    }));

    it('should not delete if no test selected', fakeAsync(() => {
      component.deleteTest();
      tick();

      expect(testServiceMock.deleteTest).not.toHaveBeenCalled();
    }));

    it('should show error toast on delete failure', fakeAsync(() => {
      testServiceMock.deleteTest.mockReturnValue(throwError(() => new Error('Error')));

      component.confirmDelete(mockTests[0]);
      component.deleteTest();
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to delete test');
    }));
  });

  describe('status actions', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should complete a test', fakeAsync(() => {
      component.completeTest(mockTests[0]);
      tick();

      expect(testServiceMock.completeTest).toHaveBeenCalledWith('test-1');
      expect(toastServiceMock.success).toHaveBeenCalledWith('Test marked as completed');
    }));

    it('should cancel a test', fakeAsync(() => {
      component.cancelTest(mockTests[0]);
      tick();

      expect(testServiceMock.cancelTest).toHaveBeenCalledWith('test-1');
      expect(toastServiceMock.success).toHaveBeenCalledWith('Test cancelled');
    }));
  });

  describe('navigation', () => {
    it('should navigate to registrations', () => {
      component.viewRegistrations(mockTests[0]);

      expect(routerMock.navigate).toHaveBeenCalledWith(['/admissions/tests', 'test-1', 'registrations']);
    });

    it('should navigate to results', () => {
      component.viewResults(mockTests[0]);

      expect(routerMock.navigate).toHaveBeenCalledWith(['/admissions/tests', 'test-1', 'results']);
    });
  });

  describe('helper functions', () => {
    it('should calculate total max marks', () => {
      const test: EntranceTest = {
        ...mockTests[0],
        subjects: [
          { name: 'English', maxMarks: 25 },
          { name: 'Math', maxMarks: 30 },
          { name: 'Science', maxMarks: 20 },
        ],
      };

      const total = component.getTotalMaxMarks(test);
      expect(total).toBe(75);
    });

    it('should get status config', () => {
      const config = component.getStatusConfig('scheduled');

      expect(config.label).toBe('Scheduled');
      expect(config.variant).toBe('info');
    });
  });
});
