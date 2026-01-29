/**
 * Teacher Assignment List Component
 * Story 5.7: Teacher Subject Assignment
 *
 * Displays list of teacher assignments with filtering and actions.
 */

import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { MslsModalComponent } from '../../../../shared/components/modal/modal.component';
import { AssignmentService } from '../../assignment.service';
import {
  Assignment,
  CreateAssignmentRequest,
  getAssignmentStatusLabel,
  formatClassSection,
} from '../../assignment.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-assignment-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule, MslsModalComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-chalkboard-user"></i>
          </div>
          <div class="header-text">
            <h1>Teacher Assignments</h1>
            <p>Manage subject-teacher assignments and class responsibilities</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" routerLink="workload">
            <i class="fa-solid fa-chart-bar"></i>
            Workload Report
          </button>
          <button class="btn btn-primary" (click)="openCreateModal()">
            <i class="fa-solid fa-plus"></i>
            New Assignment
          </button>
        </div>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <label class="filter-label">Academic Year</label>
          <select
            class="filter-select"
            [ngModel]="academicYearFilter()"
            (ngModelChange)="academicYearFilter.set($event); loadAssignments()"
          >
            <option value="">Select Year</option>
            @for (year of academicYears(); track year.id) {
              <option [value]="year.id">{{ year.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Class</label>
          <select
            class="filter-select"
            [ngModel]="classFilter()"
            (ngModelChange)="classFilter.set($event); loadAssignments()"
          >
            <option value="">All Classes</option>
            @for (cls of classes(); track cls.id) {
              <option [value]="cls.id">{{ cls.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Subject</label>
          <select
            class="filter-select"
            [ngModel]="subjectFilter()"
            (ngModelChange)="subjectFilter.set($event); loadAssignments()"
          >
            <option value="">All Subjects</option>
            @for (subject of subjects(); track subject.id) {
              <option [value]="subject.id">{{ subject.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Status</label>
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event); loadAssignments()"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
        <div class="filter-group checkbox-filter">
          <label class="checkbox-label">
            <input
              type="checkbox"
              [ngModel]="classTeacherFilter()"
              (ngModelChange)="classTeacherFilter.set($event); loadAssignments()"
            />
            Class Teachers Only
          </label>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading assignments...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadAssignments()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Teacher</th>
                <th>Subject</th>
                <th>Class/Section</th>
                <th style="text-align: center;">Periods/Week</th>
                <th style="width: 100px;">Class Teacher</th>
                <th style="width: 100px;">Status</th>
                <th>Effective From</th>
                <th style="width: 120px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (assignment of assignments(); track assignment.id) {
                <tr>
                  <td class="teacher-cell">
                    <div class="teacher-wrapper">
                      <div class="teacher-avatar">
                        {{ getInitials(assignment.staffName || '') }}
                      </div>
                      <div class="teacher-info">
                        <span class="teacher-name">{{ assignment.staffName }}</span>
                        <span class="teacher-id">{{ assignment.staffEmployeeId }}</span>
                      </div>
                    </div>
                  </td>
                  <td>
                    <div class="subject-cell">
                      <span class="subject-name">{{ assignment.subjectName }}</span>
                      <span class="subject-code">{{ assignment.subjectCode }}</span>
                    </div>
                  </td>
                  <td>{{ formatClassSection(assignment.className || '', assignment.sectionName) }}</td>
                  <td class="periods-cell">
                    <span class="periods-badge">{{ assignment.periodsPerWeek }}</span>
                  </td>
                  <td style="text-align: center;">
                    @if (assignment.isClassTeacher) {
                      <span class="badge badge-blue">
                        <i class="fa-solid fa-star"></i>
                        Yes
                      </span>
                    } @else {
                      <span class="text-muted">-</span>
                    }
                  </td>
                  <td>
                    <span class="badge" [class]="'badge-' + (assignment.status === 'active' ? 'green' : 'gray')">
                      {{ getAssignmentStatusLabel(assignment.status) }}
                    </span>
                  </td>
                  <td>{{ formatDate(assignment.effectiveFrom) }}</td>
                  <td class="actions-cell">
                    <button class="action-btn" title="Edit" (click)="editAssignment(assignment)">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    @if (assignment.status === 'active') {
                      <button
                        class="action-btn action-btn--danger"
                        title="Deactivate"
                        (click)="confirmDeactivate(assignment)"
                      >
                        <i class="fa-regular fa-circle-xmark"></i>
                      </button>
                    }
                    <button
                      class="action-btn action-btn--danger"
                      title="Delete"
                      (click)="confirmDelete(assignment)"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="8" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-rectangle-list"></i>
                      <p>No assignments found</p>
                      <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                        <i class="fa-solid fa-plus"></i>
                        Create Assignment
                      </button>
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>

          <!-- Pagination -->
          @if (hasMore()) {
            <div class="pagination">
              <button class="btn btn-secondary" (click)="loadMore()" [disabled]="loadingMore()">
                @if (loadingMore()) {
                  <div class="btn-spinner"></div>
                  Loading...
                } @else {
                  Load More
                }
              </button>
            </div>
          }
        }
      </div>

      <!-- Create Assignment Modal -->
      <msls-modal
        [isOpen]="showCreateModal()"
        title="Create Assignment"
        size="lg"
        (closed)="closeCreateModal()"
      >
        <form (ngSubmit)="createAssignment()" class="create-form">
          <div class="form-row">
            <div class="form-group">
              <label class="label required">Teacher</label>
              <select
                class="select"
                [(ngModel)]="newAssignment.staffId"
                name="staffId"
                required
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
                [(ngModel)]="newAssignment.subjectId"
                name="subjectId"
                required
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
                [(ngModel)]="newAssignment.classId"
                name="classId"
                required
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
                [(ngModel)]="newAssignment.sectionId"
                name="sectionId"
              >
                <option value="">All Sections</option>
                @for (section of filteredSections(); track section.id) {
                  <option [value]="section.id">{{ section.name }}</option>
                }
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="label required">Academic Year</label>
              <select
                class="select"
                [(ngModel)]="newAssignment.academicYearId"
                name="academicYearId"
                required
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
                [(ngModel)]="newAssignment.periodsPerWeek"
                name="periodsPerWeek"
                min="0"
                max="50"
                required
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="label required">Effective From</label>
              <input
                type="date"
                class="input"
                [(ngModel)]="newAssignment.effectiveFrom"
                name="effectiveFrom"
                required
              />
            </div>
            <div class="form-group">
              <label class="label">Effective To</label>
              <input
                type="date"
                class="input"
                [(ngModel)]="newAssignment.effectiveTo"
                name="effectiveTo"
              />
              <span class="hint">Leave empty for ongoing assignment</span>
            </div>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                [(ngModel)]="newAssignment.isClassTeacher"
                name="isClassTeacher"
              />
              Set as Class Teacher
            </label>
            <span class="hint">Only one teacher can be class teacher per class-section</span>
          </div>

          <div class="form-group">
            <label class="label">Remarks</label>
            <textarea
              class="textarea"
              [(ngModel)]="newAssignment.remarks"
              name="remarks"
              rows="2"
              placeholder="Optional notes..."
            ></textarea>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeCreateModal()">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary" [disabled]="creating()">
              @if (creating()) {
                <div class="btn-spinner"></div>
                Creating...
              } @else {
                <i class="fa-solid fa-plus"></i>
                Create Assignment
              }
            </button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Assignment"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the assignment for
            <strong>{{ assignmentToDelete()?.staffName }}</strong>
            teaching <strong>{{ assignmentToDelete()?.subjectName }}</strong>?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteAssignment()">
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                <i class="fa-solid fa-trash"></i>
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
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

    .header-actions {
      display: flex;
      gap: 0.75rem;
    }

    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
      align-items: flex-end;
    }

    .filter-group {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .checkbox-filter {
      justify-content: flex-end;
    }

    .filter-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
    }

    .filter-select {
      padding: 0.5rem 2rem 0.5rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
      min-width: 160px;
    }

    .filter-select:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      color: #374151;
      cursor: pointer;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .loading-container,
    .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
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

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table th {
      text-align: left;
      padding: 0.875rem 1rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #374151;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .teacher-cell {
      min-width: 200px;
    }

    .teacher-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .teacher-avatar {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 50%;
      background: linear-gradient(135deg, #4f46e5, #7c3aed);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 600;
      flex-shrink: 0;
    }

    .teacher-info {
      display: flex;
      flex-direction: column;
    }

    .teacher-name {
      font-weight: 500;
      color: #1e293b;
    }

    .teacher-id {
      font-size: 0.75rem;
      color: #64748b;
    }

    .subject-cell {
      display: flex;
      flex-direction: column;
    }

    .subject-name {
      font-weight: 500;
    }

    .subject-code {
      font-size: 0.75rem;
      color: #64748b;
    }

    .periods-cell {
      text-align: center;
    }

    .periods-badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      height: 1.75rem;
      padding: 0 0.5rem;
      background: #f1f5f9;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #475569;
    }

    .text-muted {
      color: #94a3b8;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge i {
      font-size: 0.625rem;
    }

    .badge-gray { background: #f1f5f9; color: #64748b; }
    .badge-blue { background: #dbeafe; color: #1e40af; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-yellow { background: #fef3c7; color: #92400e; }
    .badge-red { background: #fef2f2; color: #991b1b; }

    .actions-cell {
      text-align: right;
    }

    .action-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border: none;
      background: transparent;
      color: #64748b;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .action-btn:hover:not(:disabled) {
      background: #f1f5f9;
      color: #4f46e5;
    }

    .action-btn--danger:hover:not(:disabled) {
      background: #fef2f2;
      color: #dc2626;
    }

    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    .pagination {
      display: flex;
      justify-content: center;
      padding: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    /* Form Styles */
    .create-form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
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
    }

    .select:focus,
    .input:focus,
    .textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
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

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
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

    .btn-danger {
      background: #dc2626;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
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

    /* Delete Modal */
    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      border-radius: 50%;
      background: #fef2f2;
      color: #dc2626;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }

    .delete-confirmation p {
      margin: 0 0 0.5rem;
      color: #374151;
    }

    .delete-warning {
      font-size: 0.875rem;
      color: #64748b;
    }

    .delete-actions {
      display: flex;
      gap: 0.75rem;
      justify-content: center;
      margin-top: 1.5rem;
    }

    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        gap: 1rem;
      }

      .filters-bar {
        flex-direction: column;
      }

      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class AssignmentListComponent implements OnInit {
  private assignmentService = inject(AssignmentService);
  private toastService = inject(ToastService);
  private router = inject(Router);

  // State signals
  assignments = signal<Assignment[]>([]);
  loading = signal(true);
  loadingMore = signal(false);
  creating = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  hasMore = signal(false);
  cursor = signal<string | null>(null);

  // Filter signals
  academicYearFilter = signal<string>('');
  classFilter = signal<string>('');
  subjectFilter = signal<string>('');
  statusFilter = signal<string>('active');
  classTeacherFilter = signal(false);

  // Modal signals
  showCreateModal = signal(false);
  showDeleteModal = signal(false);
  assignmentToDelete = signal<Assignment | null>(null);

  // Reference data (would be loaded from API)
  academicYears = signal<{ id: string; name: string }[]>([]);
  classes = signal<{ id: string; name: string }[]>([]);
  sections = signal<{ id: string; name: string; classId: string }[]>([]);
  subjects = signal<{ id: string; name: string }[]>([]);
  teachers = signal<{ id: string; name: string; employeeId: string }[]>([]);
  filteredSections = signal<{ id: string; name: string }[]>([]);

  // New assignment form
  newAssignment: Partial<CreateAssignmentRequest> = {
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
  };

  ngOnInit(): void {
    this.loadReferenceData();
    this.loadAssignments();
  }

  loadReferenceData(): void {
    // TODO: Load from API - using mock data for now
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

    // Set default academic year
    if (this.academicYears().length > 0) {
      this.academicYearFilter.set(this.academicYears()[0].id);
    }
  }

  loadAssignments(): void {
    this.loading.set(true);
    this.error.set(null);
    this.cursor.set(null);

    this.assignmentService
      .getAssignments({
        academicYearId: this.academicYearFilter() || undefined,
        classId: this.classFilter() || undefined,
        subjectId: this.subjectFilter() || undefined,
        status: this.statusFilter() || undefined,
        isClassTeacher: this.classTeacherFilter() || undefined,
        limit: 20,
      })
      .subscribe({
        next: response => {
          this.assignments.set(response.assignments);
          this.hasMore.set(response.hasMore);
          this.cursor.set(response.nextCursor || null);
          this.loading.set(false);
        },
        error: () => {
          this.error.set('Failed to load assignments. Please try again.');
          this.loading.set(false);
        },
      });
  }

  loadMore(): void {
    if (!this.cursor() || this.loadingMore()) return;

    this.loadingMore.set(true);

    this.assignmentService
      .getAssignments({
        academicYearId: this.academicYearFilter() || undefined,
        classId: this.classFilter() || undefined,
        subjectId: this.subjectFilter() || undefined,
        status: this.statusFilter() || undefined,
        isClassTeacher: this.classTeacherFilter() || undefined,
        cursor: this.cursor() || undefined,
        limit: 20,
      })
      .subscribe({
        next: response => {
          this.assignments.update(current => [...current, ...response.assignments]);
          this.hasMore.set(response.hasMore);
          this.cursor.set(response.nextCursor || null);
          this.loadingMore.set(false);
        },
        error: () => {
          this.toastService.error('Failed to load more assignments');
          this.loadingMore.set(false);
        },
      });
  }

  onClassChange(classId: string): void {
    this.filteredSections.set(
      this.sections().filter(s => s.classId === classId)
    );
    this.newAssignment.sectionId = undefined;
  }

  openCreateModal(): void {
    this.newAssignment = {
      staffId: '',
      subjectId: '',
      classId: '',
      sectionId: undefined,
      academicYearId: this.academicYearFilter() || (this.academicYears().length > 0 ? this.academicYears()[0].id : ''),
      periodsPerWeek: 5,
      isClassTeacher: false,
      effectiveFrom: new Date().toISOString().split('T')[0],
      effectiveTo: undefined,
      remarks: '',
    };
    this.filteredSections.set([]);
    this.showCreateModal.set(true);
  }

  closeCreateModal(): void {
    this.showCreateModal.set(false);
  }

  createAssignment(): void {
    if (!this.newAssignment.staffId || !this.newAssignment.subjectId ||
        !this.newAssignment.classId || !this.newAssignment.academicYearId) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.creating.set(true);

    this.assignmentService.createAssignment(this.newAssignment as CreateAssignmentRequest).subscribe({
      next: () => {
        this.toastService.success('Assignment created successfully');
        this.closeCreateModal();
        this.loadAssignments();
        this.creating.set(false);
      },
      error: err => {
        const message = err?.error?.message || 'Failed to create assignment';
        this.toastService.error(message);
        this.creating.set(false);
      },
    });
  }

  editAssignment(assignment: Assignment): void {
    this.router.navigate(['/assignments', assignment.id, 'edit']);
  }

  confirmDeactivate(assignment: Assignment): void {
    // Update status to inactive
    this.assignmentService.updateAssignment(assignment.id, { status: 'inactive' }).subscribe({
      next: () => {
        this.toastService.success('Assignment deactivated');
        this.loadAssignments();
      },
      error: () => {
        this.toastService.error('Failed to deactivate assignment');
      },
    });
  }

  confirmDelete(assignment: Assignment): void {
    this.assignmentToDelete.set(assignment);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.assignmentToDelete.set(null);
  }

  deleteAssignment(): void {
    const assignment = this.assignmentToDelete();
    if (!assignment) return;

    this.deleting.set(true);

    this.assignmentService.deleteAssignment(assignment.id).subscribe({
      next: () => {
        this.toastService.success('Assignment deleted successfully');
        this.closeDeleteModal();
        this.loadAssignments();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete assignment');
        this.deleting.set(false);
      },
    });
  }

  getInitials(name: string): string {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .substring(0, 2)
      .toUpperCase();
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }

  formatClassSection = formatClassSection;
  getAssignmentStatusLabel = getAssignmentStatusLabel;
}
