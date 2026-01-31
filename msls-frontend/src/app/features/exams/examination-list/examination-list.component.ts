import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ExamService } from '../exam.service';
import {
  Examination,
  ExamType,
  ExamStatus,
  EXAM_STATUSES,
  CreateExaminationRequest,
  UpdateExaminationRequest,
  ClassSummary,
  AcademicYearSummary,
} from '../exam.model';
import { ToastService } from '../../../shared/services/toast.service';
import { ApiService } from '../../../core/services/api.service';

interface ClassOption {
  id: string;
  name: string;
}

@Component({
  selector: 'app-examination-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-calendar-check"></i>
          </div>
          <div class="header-text">
            <h1>Examinations</h1>
            <p>Create and manage examination schedules</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" (click)="showHelpPanel.set(!showHelpPanel())">
            <i class="fa-solid fa-circle-question"></i>
            <span class="btn-text">Help</span>
          </button>
          <button class="btn btn-primary" (click)="openCreateModal()">
            <i class="fa-solid fa-plus"></i>
            <span class="btn-text">New Examination</span>
          </button>
        </div>
      </div>

      <!-- Help Panel -->
      @if (showHelpPanel()) {
        <div class="help-panel">
          <div class="help-header">
            <h3><i class="fa-solid fa-graduation-cap"></i> Understanding Examinations</h3>
            <button class="help-close" (click)="showHelpPanel.set(false)">
              <i class="fa-solid fa-xmark"></i>
            </button>
          </div>
          <div class="help-grid">
            <div class="help-card">
              <div class="help-card-icon blue">
                <i class="fa-solid fa-clipboard-list"></i>
              </div>
              <h4>What are Examinations?</h4>
              <p>Examinations are scheduled assessment events. Create an examination, add subject schedules, then publish to notify students and teachers.</p>
            </div>
            <div class="help-card">
              <div class="help-card-icon purple">
                <i class="fa-solid fa-route"></i>
              </div>
              <h4>Workflow</h4>
              <ul>
                <li><strong>Draft:</strong> Create and edit freely</li>
                <li><strong>Scheduled:</strong> Published, visible to all</li>
                <li><strong>Ongoing:</strong> Exams in progress</li>
                <li><strong>Completed:</strong> Finished exams</li>
              </ul>
            </div>
            <div class="help-card">
              <div class="help-card-icon green">
                <i class="fa-solid fa-calendar-days"></i>
              </div>
              <h4>Schedules</h4>
              <p>Add subject-wise schedules with date, time, max marks, and venue. Each subject can only be scheduled once per examination.</p>
            </div>
            <div class="help-card">
              <div class="help-card-icon orange">
                <i class="fa-solid fa-lightbulb"></i>
              </div>
              <h4>Tips</h4>
              <ul>
                <li>Add at least one schedule before publishing</li>
                <li>Unpublish to make edits if needed</li>
                <li>Classes determine which students see the exam</li>
              </ul>
            </div>
          </div>
        </div>
      }

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search examinations..."
            [ngModel]="searchTerm()"
            (ngModelChange)="searchTerm.set($event)"
            class="search-input"
          />
          @if (searchTerm()) {
            <button class="clear-search" (click)="searchTerm.set('')">
              <i class="fa-solid fa-xmark"></i>
            </button>
          }
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event)"
          >
            <option value="">All Status</option>
            @for (status of examStatuses; track status.value) {
              <option [value]="status.value">{{ status.label }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="examTypeFilter()"
            (ngModelChange)="examTypeFilter.set($event)"
          >
            <option value="">All Exam Types</option>
            @for (type of examTypes(); track type.id) {
              <option [value]="type.id">{{ type.name }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Stats Summary -->
      <div class="stats-bar">
        <div class="stat-item">
          <span class="stat-value">{{ examinations().length }}</span>
          <span class="stat-label">Total</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ draftCount() }}</span>
          <span class="stat-label">Draft</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ scheduledCount() }}</span>
          <span class="stat-label">Scheduled</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ ongoingCount() }}</span>
          <span class="stat-label">Ongoing</span>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading examinations...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadExaminations()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Examination</th>
                  <th class="hide-tablet">Type</th>
                  <th>Dates</th>
                  <th class="hide-mobile">Classes</th>
                  <th class="hide-tablet">Schedules</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 160px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (exam of filteredExaminations(); track exam.id) {
                  <tr>
                    <td class="name-cell">
                      <div class="name-wrapper">
                        <div class="exam-icon" [class]="'status-' + exam.status">
                          <i [class]="getStatusIcon(exam.status)"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ exam.name }}</span>
                          @if (exam.examTypeName) {
                            <span class="type-badge-mobile">{{ exam.examTypeName }}</span>
                          }
                          @if (exam.description) {
                            <span class="description">{{ exam.description }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td class="hide-tablet">
                      @if (exam.examTypeName) {
                        <span class="type-badge">{{ exam.examTypeName }}</span>
                      }
                    </td>
                    <td class="date-cell">
                      <div class="date-range">
                        <span class="date-start">{{ formatDate(exam.startDate) }}</span>
                        <i class="fa-solid fa-arrow-right date-arrow"></i>
                        <span class="date-end">{{ formatDate(exam.endDate) }}</span>
                      </div>
                    </td>
                    <td class="hide-mobile">
                      @if (exam.classes && exam.classes.length > 0) {
                        <div class="class-badges">
                          @for (cls of exam.classes.slice(0, 3); track cls.id) {
                            <span class="class-badge">{{ cls.name }}</span>
                          }
                          @if (exam.classes.length > 3) {
                            <span class="class-badge more">+{{ exam.classes.length - 3 }}</span>
                          }
                        </div>
                      } @else {
                        <span class="no-classes">No classes</span>
                      }
                    </td>
                    <td class="hide-tablet">
                      <span class="schedule-count">
                        {{ exam.schedules?.length || 0 }} subjects
                      </span>
                    </td>
                    <td>
                      <span class="status-badge" [class]="'status-' + exam.status">
                        {{ getStatusLabel(exam.status) }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      @if (exam.status === 'draft') {
                        <button
                          class="action-btn action-btn--success"
                          title="Publish"
                          (click)="publishExam(exam)"
                          [disabled]="!exam.schedules?.length"
                        >
                          <i class="fa-solid fa-paper-plane"></i>
                        </button>
                      } @else if (exam.status === 'scheduled') {
                        <button
                          class="action-btn"
                          title="Unpublish"
                          (click)="unpublishExam(exam)"
                        >
                          <i class="fa-solid fa-rotate-left"></i>
                        </button>
                      }
                      <button
                        class="action-btn"
                        title="View/Edit Schedules"
                        [routerLink]="['/exams', exam.id, 'schedules']"
                      >
                        <i class="fa-solid fa-calendar-days"></i>
                      </button>
                      @if (exam.status === 'scheduled' || exam.status === 'ongoing') {
                        <button
                          class="action-btn"
                          title="Hall Tickets"
                          [routerLink]="['/exams', exam.id, 'hall-tickets']"
                        >
                          <i class="fa-solid fa-ticket"></i>
                        </button>
                      }
                      <button
                        class="action-btn"
                        title="Edit"
                        (click)="editExamination(exam)"
                        [disabled]="exam.status !== 'draft'"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(exam)"
                        [disabled]="exam.status !== 'draft'"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="7" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No examinations found</p>
                        @if (searchTerm() || statusFilter() || examTypeFilter()) {
                          <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                            Clear Filters
                          </button>
                        } @else {
                          <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                            <i class="fa-solid fa-plus"></i>
                            Create First Examination
                          </button>
                        }
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Create/Edit Modal -->
      @if (showFormModal()) {
        <div class="modal-overlay" (click)="closeFormModal()">
          <div class="modal modal--lg" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-calendar-check"></i>
                {{ modalTitle() }}
              </h3>
              <button type="button" class="modal__close" (click)="closeFormModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form">
                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-info-circle"></i>
                    Basic Information
                  </h4>
                  <div class="form-group">
                    <label for="name">Examination Name <span class="required">*</span></label>
                    <input
                      type="text"
                      id="name"
                      [(ngModel)]="formData.name"
                      name="name"
                      class="form-input"
                      placeholder="e.g., Mid-Term Examination 2025"
                      required
                    />
                  </div>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="examTypeId">Exam Type <span class="required">*</span></label>
                      <select
                        id="examTypeId"
                        [(ngModel)]="formData.examTypeId"
                        name="examTypeId"
                        class="form-select"
                        required
                      >
                        <option value="">Select Exam Type</option>
                        @for (type of examTypes(); track type.id) {
                          <option [value]="type.id">{{ type.name }}</option>
                        }
                      </select>
                    </div>
                    <div class="form-group">
                      <label for="academicYearId">Academic Year <span class="required">*</span></label>
                      <select
                        id="academicYearId"
                        [(ngModel)]="formData.academicYearId"
                        name="academicYearId"
                        class="form-select"
                        required
                      >
                        <option value="">Select Academic Year</option>
                        @for (year of academicYears(); track year.id) {
                          <option [value]="year.id">{{ year.name }} {{ year.isCurrent ? '(Current)' : '' }}</option>
                        }
                      </select>
                    </div>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-calendar"></i>
                    Examination Period
                  </h4>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="startDate">Start Date <span class="required">*</span></label>
                      <input
                        type="date"
                        id="startDate"
                        [(ngModel)]="formData.startDate"
                        name="startDate"
                        class="form-input"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label for="endDate">End Date <span class="required">*</span></label>
                      <input
                        type="date"
                        id="endDate"
                        [(ngModel)]="formData.endDate"
                        name="endDate"
                        class="form-input"
                        required
                      />
                    </div>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-users"></i>
                    Applicable Classes
                  </h4>
                  <div class="class-selector">
                    @for (cls of classOptions(); track cls.id) {
                      <label class="class-checkbox">
                        <input
                          type="checkbox"
                          [checked]="formData.classIds.includes(cls.id)"
                          (change)="toggleClass(cls.id)"
                        />
                        <span class="checkbox-label">{{ cls.name }}</span>
                      </label>
                    }
                  </div>
                  @if (formData.classIds.length === 0) {
                    <p class="form-hint warning">Please select at least one class</p>
                  }
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-align-left"></i>
                    Additional Details
                  </h4>
                  <div class="form-group">
                    <label for="description">Description</label>
                    <textarea
                      id="description"
                      [(ngModel)]="formData.description"
                      name="description"
                      class="form-textarea"
                      rows="3"
                      placeholder="Optional description..."
                    ></textarea>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeFormModal()">
                Cancel
              </button>
              <button
                type="button"
                class="btn btn--primary"
                [disabled]="saving() || !isFormValid()"
                (click)="saveExamination()"
              >
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ editingExamination() ? 'Saving...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid" [class.fa-check]="editingExamination()" [class.fa-plus]="!editingExamination()"></i>
                  {{ editingExamination() ? 'Save Changes' : 'Create Examination' }}
                }
              </button>
            </div>
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteModal()) {
        <div class="modal-overlay" (click)="closeDeleteModal()">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-trash"></i>
                Delete Examination
              </h3>
              <button type="button" class="modal__close" (click)="closeDeleteModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <div class="delete-confirmation">
                <div class="delete-icon">
                  <i class="fa-solid fa-triangle-exclamation"></i>
                </div>
                <p>
                  Are you sure you want to delete
                  <strong>"{{ examToDelete()?.name }}"</strong>?
                </p>
                <p class="delete-warning">
                  This will also delete all associated schedules. This action cannot be undone.
                </p>
              </div>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeDeleteModal()">
                Cancel
              </button>
              <button
                type="button"
                class="btn btn--danger"
                [disabled]="deleting()"
                (click)="deleteExamination()"
              >
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
        </div>
      }
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
      gap: 1rem;
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
      background: linear-gradient(135deg, #10b981, #059669);
      color: white;
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
      gap: 0.5rem;
    }

    /* Help Panel */
    .help-panel {
      background: linear-gradient(135deg, #ecfdf5, #d1fae5);
      border: 1px solid #a7f3d0;
      border-radius: 1rem;
      padding: 1.5rem;
      margin-bottom: 1.5rem;
    }

    .help-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }

    .help-header h3 {
      margin: 0;
      font-size: 1rem;
      font-weight: 600;
      color: #047857;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .help-close {
      background: white;
      border: none;
      width: 2rem;
      height: 2rem;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s;
    }

    .help-close:hover {
      background: #e2e8f0;
      color: #1e293b;
    }

    .help-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 1rem;
    }

    .help-card {
      background: white;
      border-radius: 0.75rem;
      padding: 1rem;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    }

    .help-card-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 0.75rem;
    }

    .help-card-icon.blue { background: #dbeafe; color: #1d4ed8; }
    .help-card-icon.purple { background: #f3e8ff; color: #7c3aed; }
    .help-card-icon.green { background: #dcfce7; color: #16a34a; }
    .help-card-icon.orange { background: #ffedd5; color: #ea580c; }

    .help-card h4 {
      margin: 0 0 0.5rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #1e293b;
    }

    .help-card p, .help-card ul {
      margin: 0;
      font-size: 0.8125rem;
      color: #64748b;
      line-height: 1.5;
    }

    .help-card ul {
      padding-left: 1rem;
    }

    .help-card li {
      margin-bottom: 0.25rem;
    }

    /* Filters */
    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
    }

    .search-box {
      flex: 1;
      min-width: 200px;
      max-width: 400px;
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      color: #9ca3af;
    }

    .search-input {
      width: 100%;
      padding: 0.625rem 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: border-color 0.2s, box-shadow 0.2s;
    }

    .search-input:focus {
      outline: none;
      border-color: #10b981;
      box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
    }

    .clear-search {
      position: absolute;
      right: 0.5rem;
      top: 50%;
      transform: translateY(-50%);
      background: none;
      border: none;
      color: #9ca3af;
      cursor: pointer;
      padding: 0.25rem;
      transition: color 0.2s;
    }

    .clear-search:hover {
      color: #6b7280;
    }

    .filter-select {
      padding: 0.625rem 2rem 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
      transition: border-color 0.2s;
    }

    .filter-select:focus {
      outline: none;
      border-color: #10b981;
    }

    /* Stats Bar */
    .stats-bar {
      display: flex;
      gap: 2rem;
      margin-bottom: 1rem;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.75rem;
    }

    .stat-item {
      display: flex;
      flex-direction: column;
    }

    .stat-value {
      font-size: 1.25rem;
      font-weight: 600;
      color: #1e293b;
    }

    .stat-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    /* Content Card */
    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .table-container {
      overflow-x: auto;
    }

    .loading-container, .error-container {
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
      border-top-color: #10b981;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Data Table */
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
      white-space: nowrap;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #374151;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    /* Name Cell */
    .name-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .exam-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .exam-icon.status-draft { background: #f1f5f9; color: #64748b; }
    .exam-icon.status-scheduled { background: #dbeafe; color: #1d4ed8; }
    .exam-icon.status-ongoing { background: #fef3c7; color: #d97706; }
    .exam-icon.status-completed { background: #dcfce7; color: #16a34a; }
    .exam-icon.status-cancelled { background: #fee2e2; color: #dc2626; }

    .name-content {
      display: flex;
      flex-direction: column;
      min-width: 0;
    }

    .name {
      font-weight: 500;
      color: #1e293b;
    }

    .type-badge-mobile {
      display: none;
      font-size: 0.75rem;
      color: #64748b;
    }

    .description {
      font-size: 0.75rem;
      color: #64748b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 200px;
    }

    .type-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      background: #f1f5f9;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      color: #475569;
    }

    /* Date Cell */
    .date-range {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
    }

    .date-arrow {
      color: #9ca3af;
      font-size: 0.625rem;
    }

    .date-start, .date-end {
      white-space: nowrap;
    }

    /* Class Badges */
    .class-badges {
      display: flex;
      flex-wrap: wrap;
      gap: 0.25rem;
    }

    .class-badge {
      display: inline-flex;
      padding: 0.125rem 0.375rem;
      background: #e0e7ff;
      color: #4338ca;
      border-radius: 0.25rem;
      font-size: 0.6875rem;
      font-weight: 500;
    }

    .class-badge.more {
      background: #f1f5f9;
      color: #64748b;
    }

    .no-classes {
      color: #9ca3af;
      font-size: 0.75rem;
    }

    .schedule-count {
      font-size: 0.875rem;
      color: #64748b;
    }

    /* Status Badge */
    .status-badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .status-badge.status-draft { background: #f1f5f9; color: #64748b; }
    .status-badge.status-scheduled { background: #dbeafe; color: #1d4ed8; }
    .status-badge.status-ongoing { background: #fef3c7; color: #d97706; }
    .status-badge.status-completed { background: #dcfce7; color: #16a34a; }
    .status-badge.status-cancelled { background: #fee2e2; color: #dc2626; }

    /* Actions */
    .actions-cell {
      text-align: right;
      white-space: nowrap;
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
      color: #10b981;
    }

    .action-btn:disabled {
      opacity: 0.3;
      cursor: not-allowed;
    }

    .action-btn--success:hover:not(:disabled) {
      background: #dcfce7;
      color: #16a34a;
    }

    .action-btn--danger:hover:not(:disabled) {
      background: #fef2f2;
      color: #dc2626;
    }

    /* Empty State */
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
      background: linear-gradient(135deg, #10b981, #059669);
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: linear-gradient(135deg, #059669, #047857);
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

    /* Modal Styles */
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      padding: 1rem;
    }

    .modal {
      background: white;
      border-radius: 1.25rem;
      width: 100%;
      max-height: 90vh;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

    .modal--sm { max-width: 28rem; }
    .modal--md { max-width: 36rem; }
    .modal--lg { max-width: 48rem; }

    .modal__header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .modal__header h3 {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .modal__header h3 i {
      color: #10b981;
    }

    .modal__close {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s ease;
    }

    .modal__close:hover {
      background: #e2e8f0;
      color: #334155;
    }

    .modal__body {
      padding: 1.5rem;
      overflow-y: auto;
      flex: 1;
    }

    .modal__footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1.25rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
      border-radius: 0 0 1.25rem 1.25rem;
    }

    .btn--primary {
      background: linear-gradient(135deg, #10b981, #059669);
      color: white;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--primary:hover:not(:disabled) { background: linear-gradient(135deg, #059669, #047857); }
    .btn--primary:disabled { opacity: 0.5; cursor: not-allowed; }

    .btn--secondary {
      background: #f1f5f9;
      color: #475569;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--secondary:hover { background: #e2e8f0; }

    .btn--danger {
      background: #dc2626;
      color: white;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--danger:hover:not(:disabled) { background: #b91c1c; }
    .btn--danger:disabled { opacity: 0.5; cursor: not-allowed; }

    /* Form Styles */
    .form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    .form-section {
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .form-section:last-of-type {
      border-bottom: none;
      padding-bottom: 0;
    }

    .form-section-title {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .form-section-title i {
      color: #10b981;
    }

    .form-row {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .form-group label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .required {
      color: #dc2626;
    }

    .form-input,
    .form-select,
    .form-textarea {
      width: 100%;
      padding: 0.75rem 1rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      background: #f8fafc;
      color: #0f172a;
      transition: all 0.2s ease;
    }

    .form-input:focus,
    .form-select:focus,
    .form-textarea:focus {
      outline: none;
      border-color: #10b981;
      background: white;
      box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
    }

    .form-textarea {
      resize: vertical;
      min-height: 80px;
    }

    .form-hint {
      margin: 0.5rem 0 0;
      font-size: 0.75rem;
      color: #64748b;
    }

    .form-hint.warning {
      color: #dc2626;
    }

    /* Class Selector */
    .class-selector {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
    }

    .class-checkbox {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 0.75rem;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .class-checkbox:has(input:checked) {
      background: #ecfdf5;
      border-color: #10b981;
    }

    .class-checkbox input {
      accent-color: #10b981;
    }

    .checkbox-label {
      font-size: 0.875rem;
      color: #374151;
    }

    /* Delete Confirmation */
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

    /* Responsive Styles */
    @media (max-width: 1024px) {
      .hide-tablet {
        display: none;
      }
    }

    @media (max-width: 768px) {
      .page {
        padding: 1rem;
      }

      .page-header {
        flex-direction: column;
        align-items: stretch;
      }

      .header-actions {
        justify-content: flex-end;
      }

      .btn-text {
        display: none;
      }

      .filters-bar {
        flex-direction: column;
      }

      .search-box {
        max-width: 100%;
      }

      .filter-group {
        width: 100%;
      }

      .filter-select {
        width: 100%;
      }

      .stats-bar {
        gap: 1rem;
        justify-content: space-between;
      }

      .hide-mobile {
        display: none;
      }

      .type-badge-mobile {
        display: block;
      }

      .form-row {
        grid-template-columns: 1fr;
      }

      .help-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class ExaminationListComponent implements OnInit {
  private examService = inject(ExamService);
  private apiService = inject(ApiService);
  private toastService = inject(ToastService);

  // Constants
  examStatuses = EXAM_STATUSES;

  // Data signals
  examinations = signal<Examination[]>([]);
  examTypes = signal<ExamType[]>([]);
  academicYears = signal<AcademicYearSummary[]>([]);
  classOptions = signal<ClassOption[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<ExamStatus | ''>('');
  examTypeFilter = signal<string>('');
  showHelpPanel = signal(false);

  // Modal state
  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingExamination = signal<Examination | null>(null);
  examToDelete = signal<Examination | null>(null);

  // Form data
  formData = this.getEmptyFormData();

  // Computed values
  modalTitle = computed(() => {
    const editing = this.editingExamination();
    return editing ? `Edit: ${editing.name}` : 'Create New Examination';
  });

  filteredExaminations = computed(() => {
    let result = this.examinations();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();
    const typeId = this.examTypeFilter();

    if (term) {
      result = result.filter(
        exam => exam.name.toLowerCase().includes(term)
      );
    }

    if (status) {
      result = result.filter(exam => exam.status === status);
    }

    if (typeId) {
      result = result.filter(exam => exam.examTypeId === typeId);
    }

    return result;
  });

  draftCount = computed(() => this.examinations().filter(e => e.status === 'draft').length);
  scheduledCount = computed(() => this.examinations().filter(e => e.status === 'scheduled').length);
  ongoingCount = computed(() => this.examinations().filter(e => e.status === 'ongoing').length);

  ngOnInit(): void {
    this.loadExaminations();
    this.loadExamTypes();
    this.loadAcademicYears();
    this.loadClasses();
  }

  loadExaminations(): void {
    this.loading.set(true);
    this.error.set(null);

    this.examService.getExaminations({}).subscribe({
      next: exams => {
        this.examinations.set(exams);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load examinations. Please try again.');
        this.loading.set(false);
      },
    });
  }

  loadExamTypes(): void {
    this.examService.getExamTypes({}).subscribe({
      next: types => this.examTypes.set(types),
      error: () => console.error('Failed to load exam types'),
    });
  }

  loadAcademicYears(): void {
    this.apiService.get<{ id: string; name: string; isCurrent: boolean }[]>('/academic-years').subscribe({
      next: years => this.academicYears.set(years),
      error: () => console.error('Failed to load academic years'),
    });
  }

  loadClasses(): void {
    this.apiService.get<ClassOption[]>('/classes').subscribe({
      next: classes => this.classOptions.set(classes),
      error: () => console.error('Failed to load classes'),
    });
  }

  openCreateModal(): void {
    this.editingExamination.set(null);
    this.formData = this.getEmptyFormData();

    // Pre-select current academic year
    const currentYear = this.academicYears().find(y => y.isCurrent);
    if (currentYear) {
      this.formData.academicYearId = currentYear.id;
    }

    this.showFormModal.set(true);
  }

  editExamination(exam: Examination): void {
    this.editingExamination.set(exam);
    this.formData = {
      name: exam.name,
      examTypeId: exam.examTypeId,
      academicYearId: exam.academicYearId,
      startDate: exam.startDate,
      endDate: exam.endDate,
      description: exam.description || '',
      classIds: exam.classes?.map(c => c.id) || [],
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingExamination.set(null);
  }

  toggleClass(classId: string): void {
    const index = this.formData.classIds.indexOf(classId);
    if (index === -1) {
      this.formData.classIds.push(classId);
    } else {
      this.formData.classIds.splice(index, 1);
    }
  }

  isFormValid(): boolean {
    return !!(
      this.formData.name &&
      this.formData.examTypeId &&
      this.formData.academicYearId &&
      this.formData.startDate &&
      this.formData.endDate &&
      this.formData.classIds.length > 0
    );
  }

  saveExamination(): void {
    if (!this.isFormValid()) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    const editing = this.editingExamination();

    const operation = editing
      ? this.examService.updateExamination(editing.id, {
          name: this.formData.name,
          examTypeId: this.formData.examTypeId,
          academicYearId: this.formData.academicYearId,
          startDate: this.formData.startDate,
          endDate: this.formData.endDate,
          description: this.formData.description || undefined,
          classIds: this.formData.classIds,
        })
      : this.examService.createExamination({
          name: this.formData.name,
          examTypeId: this.formData.examTypeId,
          academicYearId: this.formData.academicYearId,
          startDate: this.formData.startDate,
          endDate: this.formData.endDate,
          description: this.formData.description || undefined,
          classIds: this.formData.classIds,
        });

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Examination updated successfully' : 'Examination created successfully'
        );
        this.closeFormModal();
        this.loadExaminations();
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || (editing ? 'Failed to update examination' : 'Failed to create examination');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  publishExam(exam: Examination): void {
    this.examService.publishExamination(exam.id).subscribe({
      next: () => {
        this.toastService.success('Examination published successfully');
        this.loadExaminations();
      },
      error: (err) => {
        const message = err?.error?.message || 'Failed to publish examination';
        this.toastService.error(message);
      },
    });
  }

  unpublishExam(exam: Examination): void {
    this.examService.unpublishExamination(exam.id).subscribe({
      next: () => {
        this.toastService.success('Examination reverted to draft');
        this.loadExaminations();
      },
      error: (err) => {
        const message = err?.error?.message || 'Failed to unpublish examination';
        this.toastService.error(message);
      },
    });
  }

  confirmDelete(exam: Examination): void {
    this.examToDelete.set(exam);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.examToDelete.set(null);
  }

  deleteExamination(): void {
    const exam = this.examToDelete();
    if (!exam) return;

    this.deleting.set(true);

    this.examService.deleteExamination(exam.id).subscribe({
      next: () => {
        this.toastService.success('Examination deleted successfully');
        this.closeDeleteModal();
        this.loadExaminations();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || 'Failed to delete examination';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('');
    this.examTypeFilter.set('');
  }

  getStatusLabel(status: ExamStatus): string {
    return EXAM_STATUSES.find(s => s.value === status)?.label || status;
  }

  getStatusIcon(status: ExamStatus): string {
    return EXAM_STATUSES.find(s => s.value === status)?.icon || 'fa-solid fa-circle';
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', { day: '2-digit', month: 'short' });
  }

  private getEmptyFormData() {
    return {
      name: '',
      examTypeId: '',
      academicYearId: '',
      startDate: '',
      endDate: '',
      description: '',
      classIds: [] as string[],
    };
  }
}
