/**
 * Student List Component Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { signal } from '@angular/core';

import { StudentListComponent } from './student-list.component';
import { StudentService } from '../../services/student.service';
import { Student, StudentListResponse } from '../../models/student.model';
import { of } from 'rxjs';

describe('StudentListComponent', () => {
  let component: StudentListComponent;
  let fixture: ComponentFixture<StudentListComponent>;
  let studentServiceMock: Partial<StudentService>;
  let router: Router;

  const mockStudent: Student = {
    id: 'student-1',
    tenantId: 'tenant-1',
    branchId: 'branch-1',
    admissionNumber: 'MUM-2026-00001',
    firstName: 'John',
    lastName: 'Doe',
    fullName: 'John Doe',
    initials: 'JD',
    dateOfBirth: '2015-05-15',
    gender: 'male',
    status: 'active',
    admissionDate: '2026-01-15',
    createdAt: '2026-01-15T10:00:00Z',
    updatedAt: '2026-01-15T10:00:00Z',
    version: 1,
  };

  const mockListResponse: StudentListResponse = {
    students: [mockStudent],
    nextCursor: undefined,
    hasMore: false,
    total: 1,
  };

  beforeEach(async () => {
    // Create mock signals
    const studentsSignal = signal<Student[]>([]);
    const loadingSignal = signal<boolean>(false);
    const errorSignal = signal<string | null>(null);
    const totalCountSignal = signal<number>(0);

    studentServiceMock = {
      students: studentsSignal.asReadonly(),
      loading: loadingSignal.asReadonly(),
      error: errorSignal.asReadonly(),
      totalCount: totalCountSignal.asReadonly(),
      hasMore: signal(false).asReadonly(),
      isEmpty: signal(true).asReadonly(),
      loadStudents: vi.fn().mockImplementation(() => {
        studentsSignal.set(mockListResponse.students);
        totalCountSignal.set(mockListResponse.students.length);
        return of(mockListResponse);
      }),
      refresh: vi.fn().mockReturnValue(of(mockListResponse)),
      loadMore: vi.fn().mockReturnValue(null),
    };

    await TestBed.configureTestingModule({
      imports: [StudentListComponent],
      providers: [
        { provide: StudentService, useValue: studentServiceMock },
        provideRouter([]),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(StudentListComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load students on init', () => {
    fixture.detectChanges();
    expect(studentServiceMock.loadStudents).toHaveBeenCalled();
  });

  it('should apply search filter', () => {
    fixture.detectChanges();

    component.onSearchChange('John');

    expect(studentServiceMock.loadStudents).toHaveBeenCalled();
    expect(component.searchTerm()).toBe('John');
  });

  it('should apply status filter', () => {
    fixture.detectChanges();

    component.onStatusChange('active');

    expect(studentServiceMock.loadStudents).toHaveBeenCalled();
    expect(component.selectedStatus()).toBe('active');
  });

  it('should clear filters', () => {
    fixture.detectChanges();

    component.searchTerm.set('test');
    component.selectedStatus.set('active');

    component.clearFilters();

    expect(component.searchTerm()).toBe('');
    expect(component.selectedStatus()).toBe('');
    expect(studentServiceMock.loadStudents).toHaveBeenCalled();
  });

  it('should navigate to student detail on row click', () => {
    const navigateSpy = vi.spyOn(router, 'navigate');
    fixture.detectChanges();

    component.onRowClick({ id: 'student-1' });

    expect(navigateSpy).toHaveBeenCalledWith(['/students', 'student-1']);
  });

  it('should navigate to add student form', () => {
    const navigateSpy = vi.spyOn(router, 'navigate');
    fixture.detectChanges();

    component.onAddStudent();

    expect(navigateSpy).toHaveBeenCalledWith(['/students', 'new']);
  });

  it('should navigate to edit student', () => {
    const navigateSpy = vi.spyOn(router, 'navigate');
    const event = new MouseEvent('click');
    vi.spyOn(event, 'stopPropagation');

    fixture.detectChanges();

    component.onEditStudent(event, mockStudent);

    expect(event.stopPropagation).toHaveBeenCalled();
    expect(navigateSpy).toHaveBeenCalledWith(['/students', 'student-1', 'edit']);
  });

  it('should have students loaded from service', () => {
    fixture.detectChanges();

    const students = component.students();

    expect(students.length).toBe(1);
    expect(students[0].id).toBe('student-1');
    expect(students[0].fullName).toBe('John Doe');
    expect(students[0].admissionNumber).toBe('MUM-2026-00001');
    expect(students[0].status).toBe('active');
  });

  it('should get correct status variant', () => {
    expect(component.getStatusVariant('active')).toBe('success');
    expect(component.getStatusVariant('inactive')).toBe('neutral');
    expect(component.getStatusVariant('transferred')).toBe('warning');
    expect(component.getStatusVariant('graduated')).toBe('info');
  });

  it('should get correct status label', () => {
    expect(component.getStatusLabel('active')).toBe('Active');
    expect(component.getStatusLabel('inactive')).toBe('Inactive');
    expect(component.getStatusLabel('transferred')).toBe('Transferred');
    expect(component.getStatusLabel('graduated')).toBe('Graduated');
  });

  it('should format class section correctly', () => {
    const studentWithClass = {
      ...mockStudent,
      className: 'Class 5',
      sectionName: 'A',
    };

    expect(component.formatClassSection(studentWithClass)).toBe('Class 5 - A');
    expect(component.formatClassSection(mockStudent)).toBe('—');
  });

  describe('Stats Methods', () => {
    it('should count active students', () => {
      fixture.detectChanges();
      expect(component.getActiveCount()).toBe(1); // mockStudent is active
    });

    it('should count graduated students', () => {
      fixture.detectChanges();
      expect(component.getGraduatedCount()).toBe(0);
    });

    it('should count new students this month', () => {
      fixture.detectChanges();
      // mockStudent has admissionDate '2026-01-15' which is this month (January 2026)
      expect(component.getNewThisMonth()).toBe(1);
    });
  });
});
