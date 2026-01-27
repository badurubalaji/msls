/**
 * Student Service Tests
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';

import { StudentService } from './student.service';
import { Student, CreateStudentRequest, StudentListResponse } from '../models/student.model';

describe('StudentService', () => {
  let service: StudentService;
  let httpMock: HttpTestingController;

  const mockStudent: Student = {
    id: 'student-1',
    tenantId: 'tenant-1',
    branchId: 'branch-1',
    admissionNumber: 'MUM-2026-00001',
    firstName: 'John',
    middleName: 'William',
    lastName: 'Doe',
    fullName: 'John William Doe',
    initials: 'JD',
    dateOfBirth: '2015-05-15',
    gender: 'male',
    bloodGroup: 'O+',
    status: 'active',
    admissionDate: '2026-01-15',
    createdAt: '2026-01-15T10:00:00Z',
    updatedAt: '2026-01-15T10:00:00Z',
    version: 1,
  };

  const mockListResponse: StudentListResponse = {
    students: [mockStudent],
    nextCursor: 'cursor-2',
    hasMore: true,
    total: 50,
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [StudentService],
    });

    service = TestBed.inject(StudentService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('loadStudents', () => {
    it('should load students and update state', () => {
      service.loadStudents().subscribe((response) => {
        expect(response.students.length).toBe(1);
        expect(response.total).toBe(50);
      });

      const req = httpMock.expectOne('/api/v1/students');
      expect(req.request.method).toBe('GET');
      req.flush({ success: true, data: mockListResponse });

      expect(service.students().length).toBe(1);
      expect(service.totalCount()).toBe(50);
      expect(service.hasMore()).toBe(true);
      expect(service.loading()).toBe(false);
    });

    it('should apply filters to the request', () => {
      service.loadStudents({
        branchId: 'branch-1',
        status: 'active',
        search: 'John',
      }).subscribe();

      const req = httpMock.expectOne((request) =>
        request.url === '/api/v1/students' &&
        request.params.get('branchId') === 'branch-1' &&
        request.params.get('status') === 'active' &&
        request.params.get('search') === 'John'
      );
      expect(req.request.method).toBe('GET');
      req.flush({ success: true, data: mockListResponse });
    });

    it('should append students when append is true', () => {
      // First load
      service.loadStudents().subscribe();
      const req1 = httpMock.expectOne('/api/v1/students');
      req1.flush({ success: true, data: mockListResponse });

      // Second load with append
      service.loadStudents({ cursor: 'cursor-2' }, true).subscribe();
      const req2 = httpMock.expectOne((request) =>
        request.url === '/api/v1/students' &&
        request.params.get('cursor') === 'cursor-2'
      );

      const newStudent = { ...mockStudent, id: 'student-2', admissionNumber: 'MUM-2026-00002' };
      req2.flush({
        success: true,
        data: { students: [newStudent], nextCursor: null, totalCount: 50 },
      });

      expect(service.students().length).toBe(2);
    });

    it('should handle errors', () => {
      service.loadStudents().subscribe({
        error: (err) => {
          expect(err).toBeTruthy();
        },
      });

      const req = httpMock.expectOne('/api/v1/students');
      req.flush(
        { success: false, error: { title: 'Error', status: 500 } },
        { status: 500, statusText: 'Server Error' }
      );

      expect(service.error()).toBeTruthy();
      expect(service.loading()).toBe(false);
    });
  });

  describe('getStudent', () => {
    it('should load a single student', () => {
      service.getStudent('student-1').subscribe((student) => {
        expect(student.id).toBe('student-1');
        expect(student.fullName).toBe('John William Doe');
      });

      const req = httpMock.expectOne('/api/v1/students/student-1');
      expect(req.request.method).toBe('GET');
      req.flush({ success: true, data: mockStudent });

      expect(service.selectedStudent()?.id).toBe('student-1');
    });
  });

  describe('createStudent', () => {
    it('should create a student and add to list', () => {
      const createRequest: CreateStudentRequest = {
        branchId: 'branch-1',
        firstName: 'Jane',
        lastName: 'Doe',
        dateOfBirth: '2016-03-20',
        gender: 'female',
      };

      const newStudent = {
        ...mockStudent,
        id: 'student-2',
        firstName: 'Jane',
        fullName: 'Jane Doe',
      };

      service.createStudent(createRequest).subscribe((student) => {
        expect(student.firstName).toBe('Jane');
      });

      const req = httpMock.expectOne('/api/v1/students');
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(createRequest);
      req.flush({ success: true, data: newStudent });

      expect(service.students()[0].firstName).toBe('Jane');
      expect(service.selectedStudent()?.firstName).toBe('Jane');
    });
  });

  describe('updateStudent', () => {
    it('should update a student', () => {
      // First, populate the list
      service.loadStudents().subscribe();
      const loadReq = httpMock.expectOne('/api/v1/students');
      loadReq.flush({ success: true, data: mockListResponse });

      // Then update
      const updatedStudent = { ...mockStudent, firstName: 'Johnny' };
      service.updateStudent('student-1', { firstName: 'Johnny', version: 1 }).subscribe((student) => {
        expect(student.firstName).toBe('Johnny');
      });

      const req = httpMock.expectOne('/api/v1/students/student-1');
      expect(req.request.method).toBe('PUT');
      req.flush({ success: true, data: updatedStudent });

      expect(service.students()[0].firstName).toBe('Johnny');
    });
  });

  describe('deleteStudent', () => {
    it('should delete a student and remove from list', () => {
      // First, populate the list
      service.loadStudents().subscribe();
      const loadReq = httpMock.expectOne('/api/v1/students');
      loadReq.flush({ success: true, data: mockListResponse });

      expect(service.students().length).toBe(1);

      // Then delete
      service.deleteStudent('student-1').subscribe();

      const req = httpMock.expectOne('/api/v1/students/student-1');
      expect(req.request.method).toBe('DELETE');
      req.flush(null);

      expect(service.students().length).toBe(0);
      expect(service.totalCount()).toBe(49);
    });
  });

  describe('uploadPhoto', () => {
    it('should upload a photo and update student', () => {
      const file = new File(['photo'], 'photo.jpg', { type: 'image/jpeg' });
      const updatedStudent = { ...mockStudent, photoUrl: 'http://example.com/photo.jpg' };

      service.uploadPhoto('student-1', file).subscribe((student) => {
        expect(student.photoUrl).toBe('http://example.com/photo.jpg');
      });

      const req = httpMock.expectOne('/api/v1/students/student-1/photo');
      expect(req.request.method).toBe('POST');
      req.flush({ success: true, data: updatedStudent });
    });
  });

  describe('getNextAdmissionNumber', () => {
    it('should get next admission number preview', () => {
      service.getNextAdmissionNumber('branch-1').subscribe((result) => {
        expect(result.admissionNumber).toBe('MUM-2026-00002');
      });

      const req = httpMock.expectOne((request) =>
        request.url === '/api/v1/students/next-admission-number' &&
        request.params.get('branchId') === 'branch-1'
      );
      expect(req.request.method).toBe('GET');
      req.flush({ success: true, data: { admissionNumber: 'MUM-2026-00002' } });
    });
  });

  describe('state management', () => {
    it('should reset all state', () => {
      // Load some data first
      service.loadStudents().subscribe();
      const req = httpMock.expectOne('/api/v1/students');
      req.flush({ success: true, data: mockListResponse });

      expect(service.students().length).toBe(1);

      // Reset
      service.reset();

      expect(service.students().length).toBe(0);
      expect(service.selectedStudent()).toBeNull();
      expect(service.totalCount()).toBe(0);
      expect(service.hasMore()).toBe(false);
    });

    it('should clear selection', () => {
      service.selectStudent(mockStudent);
      expect(service.selectedStudent()).toBeTruthy();

      service.clearSelection();
      expect(service.selectedStudent()).toBeNull();
    });

    it('should clear error', () => {
      service.loadStudents().subscribe({ error: () => {} });
      const req = httpMock.expectOne('/api/v1/students');
      req.flush(
        { success: false, error: { title: 'Error', status: 500 } },
        { status: 500, statusText: 'Server Error' }
      );

      expect(service.error()).toBeTruthy();

      service.clearError();
      expect(service.error()).toBeNull();
    });
  });
});
