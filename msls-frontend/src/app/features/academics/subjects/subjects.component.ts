import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SubjectService } from '../services/subject.service';
import { Subject, CreateSubjectRequest, UpdateSubjectRequest, SUBJECT_TYPES, SubjectType } from '../academic.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-subjects',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-book"></i>
          </div>
          <div class="header-text">
            <h1>Subjects</h1>
            <p>Manage academic subjects and their configurations</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          <span class="btn-text">Add Subject</span>
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search subjects..."
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
            [ngModel]="typeFilter()"
            (ngModelChange)="typeFilter.set($event)"
          >
            <option value="">All Types</option>
            @for (type of subjectTypes; track type.value) {
              <option [value]="type.value">{{ type.label }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event)"
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading subjects...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadSubjects()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <!-- Desktop Table View -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Subject</th>
                  <th class="hide-mobile">Code</th>
                  <th class="hide-tablet">Type</th>
                  <th class="hide-tablet">Max Marks</th>
                  <th class="hide-mobile">Passing</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 140px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (subject of filteredSubjects(); track subject.id) {
                  <tr>
                    <td class="name-cell">
                      <div class="name-wrapper">
                        <div class="subject-icon" [class]="'type-' + subject.subjectType">
                          <i class="fa-solid fa-book"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ subject.name }}</span>
                          <span class="code-mobile">{{ subject.code }}</span>
                          @if (subject.shortName) {
                            <span class="short-name">{{ subject.shortName }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td class="code-cell hide-mobile">
                      <span class="code-badge">{{ subject.code }}</span>
                    </td>
                    <td class="type-cell hide-tablet">
                      <span class="type-badge" [class]="'type-' + subject.subjectType">
                        {{ getTypeLabel(subject.subjectType) }}
                      </span>
                    </td>
                    <td class="marks-cell hide-tablet">
                      <span class="marks-value">{{ subject.maxMarks }}</span>
                    </td>
                    <td class="marks-cell hide-mobile">
                      <span class="marks-value">{{ subject.passingMarks }}</span>
                    </td>
                    <td>
                      <span
                        class="badge"
                        [class.badge-green]="subject.isActive"
                        [class.badge-gray]="!subject.isActive"
                      >
                        {{ subject.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        title="Edit"
                        (click)="editSubject(subject)"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn"
                        title="Toggle Status"
                        (click)="toggleStatus(subject)"
                      >
                        <i
                          class="fa-solid"
                          [class.fa-toggle-on]="subject.isActive"
                          [class.fa-toggle-off]="!subject.isActive"
                        ></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(subject)"
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
                        <p>No subjects found</p>
                        @if (searchTerm() || statusFilter() !== 'all' || typeFilter()) {
                          <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                            Clear Filters
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

      <!-- Subject Form Modal -->
      @if (showSubjectModal()) {
        <div class="modal-overlay" (click)="closeSubjectModal()">
          <div class="modal modal--lg" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-book"></i>
                {{ modalTitle() }}
              </h3>
              <button type="button" class="modal__close" (click)="closeSubjectModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form" (ngSubmit)="saveSubject()">
                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-info-circle"></i>
                    Basic Information
                  </h4>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="subjectName">Subject Name <span class="required">*</span></label>
                      <input
                        type="text"
                        id="subjectName"
                        [(ngModel)]="formData.name"
                        name="name"
                        class="form-input"
                        placeholder="e.g., Mathematics"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label for="subjectCode">Code <span class="required">*</span></label>
                      <input
                        type="text"
                        id="subjectCode"
                        [(ngModel)]="formData.code"
                        name="code"
                        class="form-input"
                        placeholder="e.g., MATH"
                        required
                      />
                    </div>
                  </div>

                  <div class="form-row">
                    <div class="form-group">
                      <label for="shortName">Short Name</label>
                      <input
                        type="text"
                        id="shortName"
                        [(ngModel)]="formData.shortName"
                        name="shortName"
                        class="form-input"
                        placeholder="e.g., Math"
                      />
                    </div>
                    <div class="form-group">
                      <label for="subjectType">Type <span class="required">*</span></label>
                      <select
                        id="subjectType"
                        [(ngModel)]="formData.subjectType"
                        name="subjectType"
                        class="form-input"
                        required
                      >
                        <option value="">Select Type</option>
                        @for (type of subjectTypes; track type.value) {
                          <option [value]="type.value">{{ type.label }}</option>
                        }
                      </select>
                    </div>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-chart-line"></i>
                    Marks Configuration
                  </h4>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="maxMarks">Maximum Marks</label>
                      <input
                        type="number"
                        id="maxMarks"
                        [(ngModel)]="formData.maxMarks"
                        name="maxMarks"
                        class="form-input"
                        min="0"
                        placeholder="100"
                      />
                    </div>
                    <div class="form-group">
                      <label for="passingMarks">Passing Marks</label>
                      <input
                        type="number"
                        id="passingMarks"
                        [(ngModel)]="formData.passingMarks"
                        name="passingMarks"
                        class="form-input"
                        min="0"
                        placeholder="35"
                      />
                    </div>
                  </div>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="creditHours">Credit Hours</label>
                      <input
                        type="number"
                        id="creditHours"
                        [(ngModel)]="formData.creditHours"
                        name="creditHours"
                        class="form-input"
                        min="0"
                        step="0.5"
                        placeholder="0"
                      />
                    </div>
                    <div class="form-group">
                      <label for="displayOrder">Display Order</label>
                      <input
                        type="number"
                        id="displayOrder"
                        [(ngModel)]="formData.displayOrder"
                        name="displayOrder"
                        class="form-input"
                        min="0"
                        placeholder="0"
                      />
                    </div>
                  </div>
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
                      placeholder="Optional description about this subject..."
                    ></textarea>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeSubjectModal()">
                Cancel
              </button>
              <button type="button" class="btn btn--primary" [disabled]="saving()" (click)="saveSubject()">
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ editingSubject() ? 'Saving...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid" [class.fa-check]="editingSubject()" [class.fa-plus]="!editingSubject()"></i>
                  {{ editingSubject() ? 'Save Changes' : 'Create Subject' }}
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
                Delete Subject
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
                  <strong>"{{ subjectToDelete()?.name }}"</strong>?
                </p>
                <p class="delete-warning">
                  This action cannot be undone. Subjects assigned to classes cannot be deleted.
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
                (click)="deleteSubject()"
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
      background: #dcfce7;
      color: #16a34a;
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
      border-color: #16a34a;
      box-shadow: 0 0 0 3px rgba(22, 163, 74, 0.1);
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
      border-color: #16a34a;
    }

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
      border-top-color: #16a34a;
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

    .name-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .subject-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .subject-icon.type-core { background: #dbeafe; color: #1e40af; }
    .subject-icon.type-elective { background: #fae8ff; color: #86198f; }
    .subject-icon.type-language { background: #dcfce7; color: #166534; }
    .subject-icon.type-co_curricular { background: #fef3c7; color: #92400e; }
    .subject-icon.type-vocational { background: #e0e7ff; color: #3730a3; }

    .name-content {
      display: flex;
      flex-direction: column;
      min-width: 0;
    }

    .name {
      font-weight: 500;
      color: #1e293b;
    }

    .code-mobile {
      display: none;
      font-size: 0.75rem;
      color: #64748b;
      font-family: monospace;
    }

    .short-name {
      font-size: 0.75rem;
      color: #64748b;
    }

    .code-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      background: #f1f5f9;
      border-radius: 0.25rem;
      font-family: monospace;
      font-size: 0.75rem;
      color: #475569;
    }

    .type-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .type-badge.type-core { background: #dbeafe; color: #1e40af; }
    .type-badge.type-elective { background: #fae8ff; color: #86198f; }
    .type-badge.type-language { background: #dcfce7; color: #166534; }
    .type-badge.type-co_curricular { background: #fef3c7; color: #92400e; }
    .type-badge.type-vocational { background: #e0e7ff; color: #3730a3; }

    .marks-value {
      font-family: monospace;
      font-weight: 500;
    }

    .text-muted {
      color: #9ca3af;
      font-size: 0.875rem;
    }

    .badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }

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

    .action-btn:hover {
      background: #f1f5f9;
      color: #16a34a;
    }

    .action-btn--danger:hover {
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
      background: #16a34a;
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: #15803d;
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
      color: #16a34a;
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

      h3 {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;

        i {
          color: #16a34a;
        }
      }
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

      &:hover {
        background: #e2e8f0;
        color: #334155;
      }
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
      background: #16a34a;
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

      &:hover:not(:disabled) { background: #15803d; }
      &:disabled { opacity: 0.5; cursor: not-allowed; }
    }

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

      &:hover { background: #e2e8f0; }
    }

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

      &:hover:not(:disabled) { background: #b91c1c; }
      &:disabled { opacity: 0.5; cursor: not-allowed; }
    }

    .form-input,
    .form-textarea {
      width: 100%;
      padding: 0.75rem 1rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      background: #f8fafc;
      color: #0f172a;
      transition: all 0.2s ease;

      &:focus {
        outline: none;
        border-color: #16a34a;
        background: white;
        box-shadow: 0 0 0 3px rgba(22, 163, 74, 0.1);
      }

      &:disabled {
        background: #f1f5f9;
        color: #64748b;
        cursor: not-allowed;
      }

      &::placeholder {
        color: #94a3b8;
      }
    }

    .form-textarea {
      resize: vertical;
      min-height: 80px;
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

      .btn-text {
        display: none;
      }

      .btn-primary {
        justify-content: center;
        padding: 0.75rem 1rem;
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

      .hide-mobile {
        display: none;
      }

      .code-mobile {
        display: block;
      }

      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class SubjectsComponent implements OnInit {
  private subjectService = inject(SubjectService);
  private toastService = inject(ToastService);

  // Constants
  subjectTypes = SUBJECT_TYPES;

  // Data signals
  subjects = signal<Subject[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');
  typeFilter = signal<string>('');

  // Modal state
  showSubjectModal = signal(false);
  showDeleteModal = signal(false);
  editingSubject = signal<Subject | null>(null);
  subjectToDelete = signal<Subject | null>(null);

  // Form data
  formData: {
    name: string;
    code: string;
    shortName: string;
    subjectType: SubjectType | '';
    maxMarks: number;
    passingMarks: number;
    creditHours: number;
    displayOrder: number;
    description: string;
  } = this.getEmptyFormData();

  // Computed modal title
  modalTitle = computed(() => {
    const editing = this.editingSubject();
    return editing ? `Edit Subject: ${editing.name}` : 'Create New Subject';
  });

  // Computed filtered list
  filteredSubjects = computed(() => {
    let result = this.subjects();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();
    const type = this.typeFilter();

    if (term) {
      result = result.filter(
        subject =>
          subject.name.toLowerCase().includes(term) ||
          subject.code.toLowerCase().includes(term) ||
          (subject.shortName?.toLowerCase().includes(term) ?? false)
      );
    }

    if (status === 'active') {
      result = result.filter(subject => subject.isActive);
    } else if (status === 'inactive') {
      result = result.filter(subject => !subject.isActive);
    }

    if (type) {
      result = result.filter(subject => subject.subjectType === type);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadSubjects();
  }

  loadSubjects(): void {
    this.loading.set(true);
    this.error.set(null);

    this.subjectService.getSubjects().subscribe({
      next: subjects => {
        this.subjects.set(subjects);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load subjects. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingSubject.set(null);
    this.formData = this.getEmptyFormData();
    this.showSubjectModal.set(true);
  }

  editSubject(subject: Subject): void {
    this.editingSubject.set(subject);
    this.formData = {
      name: subject.name,
      code: subject.code,
      shortName: subject.shortName || '',
      subjectType: subject.subjectType,
      maxMarks: subject.maxMarks,
      passingMarks: subject.passingMarks,
      creditHours: subject.creditHours,
      displayOrder: subject.displayOrder,
      description: subject.description || '',
    };
    this.showSubjectModal.set(true);
  }

  closeSubjectModal(): void {
    this.showSubjectModal.set(false);
    this.editingSubject.set(null);
  }

  saveSubject(): void {
    if (!this.formData.name || !this.formData.code || !this.formData.subjectType) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    const editing = this.editingSubject();

    const data: CreateSubjectRequest | UpdateSubjectRequest = {
      name: this.formData.name,
      code: this.formData.code,
      shortName: this.formData.shortName || undefined,
      subjectType: this.formData.subjectType as SubjectType,
      maxMarks: this.formData.maxMarks,
      passingMarks: this.formData.passingMarks,
      creditHours: this.formData.creditHours,
      displayOrder: this.formData.displayOrder,
      description: this.formData.description || undefined,
    };

    const operation = editing
      ? this.subjectService.updateSubject(editing.id, data as UpdateSubjectRequest)
      : this.subjectService.createSubject(data as CreateSubjectRequest);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Subject updated successfully' : 'Subject created successfully'
        );
        this.closeSubjectModal();
        this.loadSubjects();
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || (editing ? 'Failed to update subject' : 'Failed to create subject');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  toggleStatus(subject: Subject): void {
    this.subjectService.updateSubject(subject.id, { isActive: !subject.isActive }).subscribe({
      next: () => {
        this.toastService.success(
          `Subject ${subject.isActive ? 'deactivated' : 'activated'} successfully`
        );
        this.loadSubjects();
      },
      error: () => {
        this.toastService.error('Failed to update subject status');
      },
    });
  }

  confirmDelete(subject: Subject): void {
    this.subjectToDelete.set(subject);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.subjectToDelete.set(null);
  }

  deleteSubject(): void {
    const subject = this.subjectToDelete();
    if (!subject) return;

    this.deleting.set(true);

    this.subjectService.deleteSubject(subject.id).subscribe({
      next: () => {
        this.toastService.success('Subject deleted successfully');
        this.closeDeleteModal();
        this.loadSubjects();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || 'Failed to delete subject. It may be assigned to classes.';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
    this.typeFilter.set('');
  }

  getTypeLabel(type: string): string {
    const found = SUBJECT_TYPES.find(t => t.value === type);
    return found ? found.label : type;
  }

  private getEmptyFormData() {
    return {
      name: '',
      code: '',
      shortName: '',
      subjectType: '' as SubjectType | '',
      maxMarks: 100,
      passingMarks: 35,
      creditHours: 0,
      displayOrder: 0,
      description: '',
    };
  }
}
