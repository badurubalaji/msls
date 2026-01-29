/**
 * Assignment Form Component
 * Story 5.7: Teacher Subject Assignment
 *
 * Form for creating and editing teacher assignments.
 */

import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { AssignmentService } from '../../assignment.service';
import {
  Assignment,
  CreateAssignmentRequest,
  UpdateAssignmentRequest,
} from '../../assignment.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-assignment-form',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <button class="back-btn" routerLink="/assignments">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-icon">
            <i class="fa-solid fa-chalkboard-user"></i>
          </div>
          <div class="header-text">
            <h1>{{ isEditMode() ? 'Edit' : 'New' }} Assignment</h1>
            <p>{{ isEditMode() ? 'Update teacher assignment details' : 'Assign a subject to a teacher' }}</p>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading...</span>
          </div>
        } @else {
          <form (ngSubmit)="onSubmit()" class="form">
            <div class="form-section">
              <h3 class="section-title">Assignment Details</h3>

              <div class="form-row">
                <div class="form-group">
                  <label class="label required">Teacher</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.staffId"
                    name="staffId"
                    required
                    [disabled]="isEditMode()"
                  >
                    <option value="">Select Teacher</option>
                    @for (teacher of teachers(); track teacher.id) {
                      <option [value]="teacher.id">{{ teacher.name }} ({{ teacher.employeeId }})</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label class="label required">Subject</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.subjectId"
                    name="subjectId"
                    required
                    [disabled]="isEditMode()"
                  >
                    <option value="">Select Subject</option>
                    @for (subject of subjects(); track subject.id) {
                      <option [value]="subject.id">{{ subject.name }}</option>
                    }
                  </select>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="label required">Class</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.classId"
                    name="classId"
                    required
                    [disabled]="isEditMode()"
                    (ngModelChange)="onClassChange($event)"
                  >
                    <option value="">Select Class</option>
                    @for (cls of classes(); track cls.id) {
                      <option [value]="cls.id">{{ cls.name }}</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label class="label">Section</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.sectionId"
                    name="sectionId"
                    [disabled]="isEditMode()"
                  >
                    <option value="">All Sections</option>
                    @for (section of filteredSections(); track section.id) {
                      <option [value]="section.id">{{ section.name }}</option>
                    }
                  </select>
                  <span class="hint">Leave empty to assign for all sections</span>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="label required">Academic Year</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.academicYearId"
                    name="academicYearId"
                    required
                    [disabled]="isEditMode()"
                  >
                    <option value="">Select Academic Year</option>
                    @for (year of academicYears(); track year.id) {
                      <option [value]="year.id">{{ year.name }}</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label class="label required">Periods per Week</label>
                  <input
                    type="number"
                    class="input"
                    [(ngModel)]="formData.periodsPerWeek"
                    name="periodsPerWeek"
                    min="0"
                    max="50"
                    required
                  />
                  <span class="hint">Number of teaching periods assigned per week</span>
                </div>
              </div>
            </div>

            <div class="form-section">
              <h3 class="section-title">Effective Period</h3>

              <div class="form-row">
                <div class="form-group">
                  <label class="label required">Effective From</label>
                  <input
                    type="date"
                    class="input"
                    [(ngModel)]="formData.effectiveFrom"
                    name="effectiveFrom"
                    required
                  />
                </div>
                <div class="form-group">
                  <label class="label">Effective To</label>
                  <input
                    type="date"
                    class="input"
                    [(ngModel)]="formData.effectiveTo"
                    name="effectiveTo"
                  />
                  <span class="hint">Leave empty for ongoing assignment</span>
                </div>
              </div>
            </div>

            <div class="form-section">
              <h3 class="section-title">Additional Options</h3>

              <div class="form-group">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    [(ngModel)]="formData.isClassTeacher"
                    name="isClassTeacher"
                  />
                  Set as Class Teacher
                </label>
                <span class="hint">Only one teacher can be class teacher per class-section. This will replace existing class teacher if any.</span>
              </div>

              <div class="form-group">
                <label class="label">Remarks</label>
                <textarea
                  class="textarea"
                  [(ngModel)]="formData.remarks"
                  name="remarks"
                  rows="3"
                  placeholder="Optional notes about this assignment..."
                ></textarea>
              </div>

              @if (isEditMode()) {
                <div class="form-group">
                  <label class="label">Status</label>
                  <select
                    class="select"
                    [(ngModel)]="formData.status"
                    name="status"
                  >
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                  </select>
                </div>
              }
            </div>

            <div class="form-actions">
              <button type="button" class="btn btn-secondary" routerLink="/assignments">
                Cancel
              </button>
              <button type="submit" class="btn btn-primary" [disabled]="saving()">
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ isEditMode() ? 'Updating...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid fa-check"></i>
                  {{ isEditMode() ? 'Update Assignment' : 'Create Assignment' }}
                }
              </button>
            </div>
          </form>
        }
      </div>
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 900px;
      margin: 0 auto;
    }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
      background: white;
      color: #64748b;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.2s;
    }

    .back-btn:hover {
      background: #f1f5f9;
      color: #1e293b;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: #dbeafe;
      color: #2563eb;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    .loading-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .form {
      display: flex;
      flex-direction: column;
      gap: 2rem;
    }

    .form-section {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .section-title {
      margin: 0;
      font-size: 1rem;
      font-weight: 600;
      color: #1e293b;
      padding-bottom: 0.75rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .label.required::after {
      content: ' *';
      color: #dc2626;
    }

    .select,
    .input,
    .textarea {
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: all 0.2s;
    }

    .select:focus,
    .input:focus,
    .textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .select:disabled,
    .input:disabled {
      background: #f1f5f9;
      cursor: not-allowed;
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
      cursor: pointer;
    }

    .checkbox-label input {
      width: 1rem;
      height: 1rem;
      cursor: pointer;
    }

    .hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @media (max-width: 640px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class AssignmentFormComponent implements OnInit {
  private assignmentService = inject(AssignmentService);
  private toastService = inject(ToastService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  // State
  loading = signal(false);
  saving = signal(false);
  isEditMode = signal(false);
  assignmentId = signal<string | null>(null);

  // Reference data
  academicYears = signal<{ id: string; name: string }[]>([]);
  classes = signal<{ id: string; name: string }[]>([]);
  sections = signal<{ id: string; name: string; classId: string }[]>([]);
  subjects = signal<{ id: string; name: string }[]>([]);
  teachers = signal<{ id: string; name: string; employeeId: string }[]>([]);
  filteredSections = signal<{ id: string; name: string }[]>([]);

  // Form data
  formData: Partial<CreateAssignmentRequest & UpdateAssignmentRequest & { status?: string }> = {
    staffId: '',
    subjectId: '',
    classId: '',
    sectionId: undefined,
    academicYearId: '',
    periodsPerWeek: 5,
    isClassTeacher: false,
    effectiveFrom: new Date().toISOString().split('T')[0],
    effectiveTo: undefined,
    remarks: '',
    status: 'active',
  };

  ngOnInit(): void {
    this.loadReferenceData();

    // Check for edit mode
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEditMode.set(true);
      this.assignmentId.set(id);
      this.loadAssignment(id);
    }

    // Check for query params (from unassigned subjects)
    const queryParams = this.route.snapshot.queryParams;
    if (queryParams['subjectId']) {
      this.formData.subjectId = queryParams['subjectId'];
    }
    if (queryParams['classId']) {
      this.formData.classId = queryParams['classId'];
      this.onClassChange(queryParams['classId']);
    }
    if (queryParams['sectionId']) {
      this.formData.sectionId = queryParams['sectionId'];
    }
  }

  loadReferenceData(): void {
    // TODO: Load from API
    this.academicYears.set([
      { id: '1', name: '2024-25' },
      { id: '2', name: '2025-26' },
    ]);
    this.classes.set([
      { id: '1', name: 'Class 1' },
      { id: '2', name: 'Class 2' },
      { id: '3', name: 'Class 3' },
    ]);
    this.sections.set([
      { id: '1', name: 'Section A', classId: '1' },
      { id: '2', name: 'Section B', classId: '1' },
      { id: '3', name: 'Section A', classId: '2' },
    ]);
    this.subjects.set([
      { id: '1', name: 'Mathematics' },
      { id: '2', name: 'English' },
      { id: '3', name: 'Science' },
    ]);
    this.teachers.set([
      { id: '1', name: 'John Smith', employeeId: 'EMP001' },
      { id: '2', name: 'Jane Doe', employeeId: 'EMP002' },
    ]);

    if (this.academicYears().length > 0 && !this.formData.academicYearId) {
      this.formData.academicYearId = this.academicYears()[0].id;
    }
  }

  loadAssignment(id: string): void {
    this.loading.set(true);

    this.assignmentService.getAssignment(id).subscribe({
      next: assignment => {
        this.formData = {
          staffId: assignment.staffId,
          subjectId: assignment.subjectId,
          classId: assignment.classId,
          sectionId: assignment.sectionId || undefined,
          academicYearId: assignment.academicYearId,
          periodsPerWeek: assignment.periodsPerWeek,
          isClassTeacher: assignment.isClassTeacher,
          effectiveFrom: assignment.effectiveFrom,
          effectiveTo: assignment.effectiveTo || undefined,
          remarks: assignment.remarks || '',
          status: assignment.status,
        };

        if (assignment.classId) {
          this.onClassChange(assignment.classId);
        }

        this.loading.set(false);
      },
      error: () => {
        this.toastService.error('Failed to load assignment');
        this.router.navigate(['/assignments']);
      },
    });
  }

  onClassChange(classId: string): void {
    this.filteredSections.set(
      this.sections().filter(s => s.classId === classId)
    );
    if (!this.isEditMode()) {
      this.formData.sectionId = undefined;
    }
  }

  onSubmit(): void {
    if (!this.formData.staffId || !this.formData.subjectId ||
        !this.formData.classId || !this.formData.academicYearId) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);

    if (this.isEditMode() && this.assignmentId()) {
      const updateData: UpdateAssignmentRequest = {
        periodsPerWeek: this.formData.periodsPerWeek,
        isClassTeacher: this.formData.isClassTeacher,
        effectiveFrom: this.formData.effectiveFrom,
        effectiveTo: this.formData.effectiveTo || undefined,
        status: this.formData.status as 'active' | 'inactive',
        remarks: this.formData.remarks,
      };

      this.assignmentService.updateAssignment(this.assignmentId()!, updateData).subscribe({
        next: () => {
          this.toastService.success('Assignment updated successfully');
          this.router.navigate(['/assignments']);
        },
        error: err => {
          const message = err?.error?.message || 'Failed to update assignment';
          this.toastService.error(message);
          this.saving.set(false);
        },
      });
    } else {
      const createData: CreateAssignmentRequest = {
        staffId: this.formData.staffId!,
        subjectId: this.formData.subjectId!,
        classId: this.formData.classId!,
        sectionId: this.formData.sectionId,
        academicYearId: this.formData.academicYearId!,
        periodsPerWeek: this.formData.periodsPerWeek!,
        isClassTeacher: this.formData.isClassTeacher || false,
        effectiveFrom: this.formData.effectiveFrom!,
        effectiveTo: this.formData.effectiveTo,
        remarks: this.formData.remarks,
      };

      this.assignmentService.createAssignment(createData).subscribe({
        next: () => {
          this.toastService.success('Assignment created successfully');
          this.router.navigate(['/assignments']);
        },
        error: err => {
          const message = err?.error?.message || 'Failed to create assignment';
          this.toastService.error(message);
          this.saving.set(false);
        },
      });
    }
  }
}
