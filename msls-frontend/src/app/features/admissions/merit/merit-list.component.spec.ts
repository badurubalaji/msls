/**
 * MeritListComponent Unit Tests
 */
import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { MeritListComponent } from './merit-list.component';
import { MeritService } from './merit.service';
import { AdmissionSessionService } from '../sessions/admission-session.service';
import { ToastService } from '../../../shared/services/toast.service';
import { MeritList, MeritListEntry, ApplicationStatus } from './merit.model';

describe('MeritListComponent', () => {
  let component: MeritListComponent;
  let fixture: ComponentFixture<MeritListComponent>;
  let meritServiceMock: {
    getMeritList: ReturnType<typeof vi.fn>;
    generateMeritList: ReturnType<typeof vi.fn>;
    makeDecision: ReturnType<typeof vi.fn>;
    makeBulkDecision: ReturnType<typeof vi.fn>;
    generateOfferLetter: ReturnType<typeof vi.fn>;
    getAvailableClasses: ReturnType<typeof vi.fn>;
    getSections: ReturnType<typeof vi.fn>;
  };
  let sessionServiceMock: {
    getSessions: ReturnType<typeof vi.fn>;
  };
  let toastServiceMock: {
    success: ReturnType<typeof vi.fn>;
    error: ReturnType<typeof vi.fn>;
    info: ReturnType<typeof vi.fn>;
  };
  let routerMock: {
    navigate: ReturnType<typeof vi.fn>;
  };

  const mockSessions = [
    { id: 'session-1', name: 'Admission 2026', status: 'open' as const },
    { id: 'session-2', name: 'Admission 2027', status: 'upcoming' as const },
  ];

  const mockEntries: MeritListEntry[] = [
    {
      id: 'entry-1',
      applicationId: 'app-1',
      studentName: 'Student One',
      parentName: 'Parent One',
      parentPhone: '9876543210',
      applicationNumber: 'APP-2026-0001',
      classApplying: 'Class 1',
      totalScore: 90,
      maxScore: 100,
      percentage: 90,
      rank: 1,
      status: 'test_completed' as ApplicationStatus,
      subjectScores: [
        { subjectName: 'English', score: 45, maxScore: 50 },
        { subjectName: 'Math', score: 45, maxScore: 50 },
      ],
    },
    {
      id: 'entry-2',
      applicationId: 'app-2',
      studentName: 'Student Two',
      parentName: 'Parent Two',
      parentPhone: '9876543211',
      applicationNumber: 'APP-2026-0002',
      classApplying: 'Class 1',
      totalScore: 85,
      maxScore: 100,
      percentage: 85,
      rank: 2,
      status: 'selected' as ApplicationStatus,
      sectionAssigned: 'Section A',
    },
    {
      id: 'entry-3',
      applicationId: 'app-3',
      studentName: 'Student Three',
      parentName: 'Parent Three',
      parentPhone: '9876543212',
      applicationNumber: 'APP-2026-0003',
      classApplying: 'Class 1',
      totalScore: 75,
      maxScore: 100,
      percentage: 75,
      rank: 3,
      status: 'waitlisted' as ApplicationStatus,
      waitlistPosition: 1,
    },
  ];

  const mockMeritList: MeritList = {
    id: 'merit-1',
    tenantId: 'tenant-1',
    sessionId: 'session-1',
    className: 'Class 1',
    testId: 'test-1',
    testName: 'Entrance Test 2026',
    generatedAt: '2026-01-15T10:00:00Z',
    cutoffScore: 50,
    entries: mockEntries,
    isFinal: false,
    createdAt: '2026-01-15T10:00:00Z',
  };

  beforeEach(async () => {
    meritServiceMock = {
      getMeritList: vi.fn(),
      generateMeritList: vi.fn(),
      makeDecision: vi.fn(),
      makeBulkDecision: vi.fn(),
      generateOfferLetter: vi.fn(),
      getAvailableClasses: vi.fn(),
      getSections: vi.fn(),
    };
    sessionServiceMock = {
      getSessions: vi.fn(),
    };
    toastServiceMock = {
      success: vi.fn(),
      error: vi.fn(),
      info: vi.fn(),
    };
    routerMock = {
      navigate: vi.fn(),
    };

    sessionServiceMock.getSessions.mockReturnValue(of(mockSessions.map(s => ({
      id: s.id,
      name: s.name,
      academicYearId: 'ay-1',
      startDate: '2026-01-01',
      endDate: '2026-03-31',
      status: s.status,
      applicationFee: 500,
      requiredDocuments: [],
      settings: {},
      totalApplications: 100,
      totalSeats: 50,
      filledSeats: 20,
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    }))));

    meritServiceMock.getAvailableClasses.mockReturnValue(of(['LKG', 'UKG', 'Class 1', 'Class 2']));
    meritServiceMock.getMeritList.mockReturnValue(of(mockMeritList));
    meritServiceMock.generateMeritList.mockReturnValue(of(mockMeritList));
    meritServiceMock.makeDecision.mockReturnValue(of({
      id: 'decision-1',
      tenantId: 'tenant-1',
      applicationId: 'app-1',
      decision: 'selected',
      decisionDate: '2026-01-20',
      createdAt: '2026-01-20T00:00:00Z',
      updatedAt: '2026-01-20T00:00:00Z',
    }));
    meritServiceMock.makeBulkDecision.mockReturnValue(of({
      successful: 2,
      failed: 0,
      decisions: [],
      errors: [],
    }));
    meritServiceMock.generateOfferLetter.mockReturnValue(of({
      url: '/offer.pdf',
      validUntil: '2026-02-20',
    }));
    meritServiceMock.getSections.mockReturnValue(of([
      { id: 'section-1', name: 'Section A' },
      { id: 'section-2', name: 'Section B' },
    ]));

    await TestBed.configureTestingModule({
      imports: [MeritListComponent],
      providers: [
        { provide: MeritService, useValue: meritServiceMock },
        { provide: AdmissionSessionService, useValue: sessionServiceMock },
        { provide: ToastService, useValue: toastServiceMock },
        { provide: Router, useValue: routerMock },
        { provide: ActivatedRoute, useValue: {} },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(MeritListComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load sessions on init', fakeAsync(() => {
    fixture.detectChanges();
    tick();

    expect(sessionServiceMock.getSessions).toHaveBeenCalled();
    expect(component.sessionOptions().length).toBe(2);
  }));

  describe('session and class selection', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should load classes when session is selected', fakeAsync(() => {
      component.onSessionChange('session-1');
      tick();

      expect(meritServiceMock.getAvailableClasses).toHaveBeenCalledWith('session-1');
      expect(component.classOptions().length).toBe(4);
    }));

    it('should clear class selection when session changes', fakeAsync(() => {
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();

      expect(component.selectedClass()).toBe('Class 1');

      component.onSessionChange('session-2');
      tick();

      expect(component.selectedClass()).toBeNull();
      expect(component.meritList()).toBeNull();
    }));

    it('should load merit list when class is selected', fakeAsync(() => {
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();

      expect(meritServiceMock.getMeritList).toHaveBeenCalledWith('session-1', 'Class 1');
      expect(component.meritList()).toEqual(mockMeritList);
    }));
  });

  describe('generate merit list', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should generate merit list', fakeAsync(() => {
      component.generateMeritList();
      tick();

      expect(meritServiceMock.generateMeritList).toHaveBeenCalled();
      expect(toastServiceMock.success).toHaveBeenCalled();
    }));

    it('should include cutoff score when set', fakeAsync(() => {
      component.cutoffScore.set(60);
      component.generateMeritList();
      tick();

      expect(meritServiceMock.generateMeritList).toHaveBeenCalledWith('session-1', {
        className: 'Class 1',
        cutoffScore: 60,
      });
    }));

    it('should show error toast on failure', fakeAsync(() => {
      meritServiceMock.generateMeritList.mockReturnValue(throwError(() => new Error('Error')));

      component.generateMeritList();
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to generate merit list');
    }));

    it('should not generate if session or class not selected', () => {
      component.selectedSession.set(null);

      expect(component.canGenerate()).toBe(false);
    });
  });

  describe('computed statistics', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should calculate selected count', () => {
      expect(component.selectedCount()).toBe(1); // One entry with 'selected' status
    });

    it('should calculate waitlisted count', () => {
      expect(component.waitlistedCount()).toBe(1);
    });

    it('should calculate pending count', () => {
      expect(component.pendingCount()).toBe(1); // 'test_completed' status
    });
  });

  describe('entry selection', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should select an entry', () => {
      component.toggleEntrySelection('app-1', true);

      expect(component.selectedEntries()).toContain('app-1');
    });

    it('should deselect an entry', () => {
      component.toggleEntrySelection('app-1', true);
      component.toggleEntrySelection('app-1', false);

      expect(component.selectedEntries()).not.toContain('app-1');
    });

    it('should select all entries', () => {
      component.toggleSelectAll(true);

      expect(component.selectedEntries().length).toBe(3);
      expect(component.allSelected()).toBe(true);
    });

    it('should deselect all entries', () => {
      component.toggleSelectAll(true);
      component.toggleSelectAll(false);

      expect(component.selectedEntries().length).toBe(0);
    });

    it('should clear selection', () => {
      component.toggleSelectAll(true);
      component.clearSelection();

      expect(component.selectedEntries().length).toBe(0);
    });

    it('should detect some selected state', () => {
      component.toggleEntrySelection('app-1', true);

      expect(component.someSelected()).toBe(true);
      expect(component.allSelected()).toBe(false);
    });
  });

  describe('decision modal', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should open single decision modal', () => {
      const entry = mockEntries[0];
      component.openDecisionModal(entry, 'selected');

      expect(component.showDecisionModal()).toBe(true);
      expect(component.selectedEntry()).toEqual(entry);
      expect(component.currentDecision()).toBe('selected');
      expect(component.isBulkDecision()).toBe(false);
    });

    it('should open bulk decision modal', () => {
      component.toggleEntrySelection('app-1', true);
      component.toggleEntrySelection('app-2', true);
      component.openBulkDecisionModal('selected');

      expect(component.showDecisionModal()).toBe(true);
      expect(component.bulkEntries().length).toBe(2);
      expect(component.isBulkDecision()).toBe(true);
    });

    it('should close decision modal and clear state', () => {
      component.openDecisionModal(mockEntries[0], 'selected');
      component.closeDecisionModal();

      expect(component.showDecisionModal()).toBe(false);
      expect(component.selectedEntry()).toBeNull();
    });

    it('should generate correct modal title for single decision', () => {
      component.openDecisionModal(mockEntries[0], 'selected');

      expect(component.decisionModalTitle()).toContain('Approve');
      expect(component.decisionModalTitle()).toContain('Student One');
    });

    it('should generate correct modal title for bulk decision', () => {
      component.toggleSelectAll(true);
      component.openBulkDecisionModal('waitlisted');

      expect(component.decisionModalTitle()).toBe('Waitlist 3 Candidates');
    });
  });

  describe('helper functions', () => {
    it('should get correct rank class', () => {
      expect(component.getRankClass(1)).toBe('rank-1');
      expect(component.getRankClass(2)).toBe('rank-2');
      expect(component.getRankClass(3)).toBe('rank-3');
      expect(component.getRankClass(4)).toBe('rank-other');
      expect(component.getRankClass(100)).toBe('rank-other');
    });

    it('should determine if decision can be made', () => {
      const testCompleted = { ...mockEntries[0], status: 'test_completed' as ApplicationStatus };
      const selected = { ...mockEntries[0], status: 'selected' as ApplicationStatus };

      expect(component.canMakeDecision(testCompleted)).toBe(true);
      expect(component.canMakeDecision(selected)).toBe(false);
    });

    it('should get status config', () => {
      const config = component.getStatusConfig('selected');

      expect(config.label).toBe('Selected');
      expect(config.variant).toBe('success');
    });
  });

  describe('navigation', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should navigate to enrollment page', () => {
      component.goToEnrollment(mockEntries[0]);

      expect(routerMock.navigate).toHaveBeenCalledWith(['/admissions/enrollment', 'app-1']);
    });
  });

  describe('view offer letter', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
      component.onSessionChange('session-1');
      tick();
      component.onClassChange('Class 1');
      tick();
    }));

    it('should generate offer letter', fakeAsync(() => {
      component.viewOfferLetter(mockEntries[1]);
      tick();

      expect(meritServiceMock.generateOfferLetter).toHaveBeenCalledWith('app-2');
      expect(toastServiceMock.info).toHaveBeenCalled();
    }));

    it('should show error on failure', fakeAsync(() => {
      meritServiceMock.generateOfferLetter.mockReturnValue(throwError(() => new Error('Error')));

      component.viewOfferLetter(mockEntries[1]);
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to generate offer letter');
    }));
  });

  describe('cutoff score handling', () => {
    it('should update cutoff score within bounds', () => {
      const mockEvent = { target: { value: '75' } } as unknown as Event;
      component.onCutoffChange(mockEvent);

      expect(component.cutoffScore()).toBe(75);
    });

    it('should clamp cutoff score to 0-100', () => {
      const mockEventNegative = { target: { value: '-10' } } as unknown as Event;
      component.onCutoffChange(mockEventNegative);
      expect(component.cutoffScore()).toBe(0);

      const mockEventOver = { target: { value: '150' } } as unknown as Event;
      component.onCutoffChange(mockEventOver);
      expect(component.cutoffScore()).toBe(100);
    });
  });
});
