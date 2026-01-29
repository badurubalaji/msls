import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClassService } from '../services/class.service';
import { StreamService } from '../services/stream.service';
import { Class, CreateClassRequest, UpdateClassRequest, Stream, CLASS_LEVELS, ClassLevel } from '../academic.model';
import { ToastService } from '../../../shared/services/toast.service';
import { BranchService } from '../../admin/branches/branch.service';
import { Branch } from '../../admin/branches/branch.model';

@Component({
  selector: 'msls-classes',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-graduation-cap"></i>
          </div>
          <div class="header-text">
            <h1>Classes</h1>
            <p>Manage academic classes and their configurations</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          <span class="btn-text">Add Class</span>
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search classes..."
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
            [ngModel]="levelFilter()"
            (ngModelChange)="levelFilter.set($event)"
          >
            <option value="">All Levels</option>
            @for (level of classLevels; track level.value) {
              <option [value]="level.value">{{ level.label }}</option>
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
            <span>Loading classes...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadClasses()">
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
                  <th>Class</th>
                  <th class="hide-mobile">Code</th>
                  <th class="hide-tablet">Level</th>
                  <th class="hide-tablet">Streams</th>
                  <th class="hide-mobile" style="width: 100px; text-align: center;">Sections</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 140px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (cls of filteredClasses(); track cls.id) {
                  <tr>
                    <td class="name-cell">
                      <div class="name-wrapper">
                        <div class="class-icon">
                          <i class="fa-solid fa-graduation-cap"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ cls.name }}</span>
                          <span class="code-mobile">{{ cls.code }}</span>
                          @if (cls.description) {
                            <span class="description">{{ cls.description }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td class="code-cell hide-mobile">
                      <span class="code-badge">{{ cls.code }}</span>
                    </td>
                    <td class="level-cell hide-tablet">
                      @if (cls.level) {
                        <span class="level-badge" [class]="'level-' + cls.level">
                          {{ getLevelLabel(cls.level) }}
                        </span>
                      } @else {
                        <span class="text-muted">-</span>
                      }
                    </td>
                    <td class="streams-cell hide-tablet">
                      @if (cls.hasStreams && cls.streams?.length) {
                        <div class="stream-tags">
                          @for (stream of cls.streams; track stream.id) {
                            <span class="stream-tag">{{ stream.code }}</span>
                          }
                        </div>
                      } @else if (cls.hasStreams) {
                        <span class="text-muted">Stream-based</span>
                      } @else {
                        <span class="text-muted">-</span>
                      }
                    </td>
                    <td class="sections-cell hide-mobile">
                      <span class="section-count">{{ cls.sections?.length || 0 }}</span>
                    </td>
                    <td>
                      <span
                        class="badge"
                        [class.badge-green]="cls.isActive"
                        [class.badge-gray]="!cls.isActive"
                      >
                        {{ cls.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        title="Edit"
                        (click)="editClass(cls)"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn"
                        title="Toggle Status"
                        (click)="toggleStatus(cls)"
                      >
                        <i
                          class="fa-solid"
                          [class.fa-toggle-on]="cls.isActive"
                          [class.fa-toggle-off]="!cls.isActive"
                        ></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(cls)"
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
                        <p>No classes found</p>
                        @if (searchTerm() || statusFilter() !== 'all' || levelFilter()) {
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

      <!-- Class Form Modal -->
      @if (showClassModal()) {
        <div class="modal-overlay" (click)="closeClassModal()">
          <div class="modal modal--lg" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-graduation-cap"></i>
                {{ modalTitle() }}
              </h3>
              <button type="button" class="modal__close" (click)="closeClassModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form" (ngSubmit)="saveClass()">
                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-info-circle"></i>
                    Basic Information
                  </h4>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="className">Class Name <span class="required">*</span></label>
                      <input
                        type="text"
                        id="className"
                        [(ngModel)]="formData.name"
                        name="name"
                        class="form-input"
                        placeholder="e.g., Class 10"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label for="classCode">Code <span class="required">*</span></label>
                      <input
                        type="text"
                        id="classCode"
                        [(ngModel)]="formData.code"
                        name="code"
                        class="form-input"
                        placeholder="e.g., X"
                        required
                      />
                    </div>
                  </div>

                  <div class="form-row">
                    <div class="form-group">
                      <label for="branchId">Branch <span class="required">*</span></label>
                      <select
                        id="branchId"
                        [(ngModel)]="formData.branchId"
                        name="branchId"
                        class="form-input"
                        [disabled]="!!editingClass()"
                        required
                      >
                        <option value="">Select Branch</option>
                        @for (branch of branches(); track branch.id) {
                          <option [value]="branch.id">{{ branch.name }}</option>
                        }
                      </select>
                    </div>
                    <div class="form-group">
                      <label for="level">Level</label>
                      <select id="level" [(ngModel)]="formData.level" name="level" class="form-input">
                        <option value="">Select Level</option>
                        @for (level of classLevels; track level.value) {
                          <option [value]="level.value">{{ level.label }}</option>
                        }
                      </select>
                    </div>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-sliders"></i>
                    Configuration
                  </h4>
                  <div class="form-row">
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
                    <div class="form-group checkbox-group">
                      <label class="checkbox-label">
                        <input
                          type="checkbox"
                          [(ngModel)]="formData.hasStreams"
                          name="hasStreams"
                        />
                        <span class="checkbox-text">Has Streams</span>
                      </label>
                    </div>
                  </div>

                  @if (formData.hasStreams) {
                    <div class="form-group">
                      <label>Assign Streams</label>
                      <div class="stream-checkboxes">
                        @for (stream of streams(); track stream.id) {
                          <label class="checkbox-label stream-checkbox">
                            <input
                              type="checkbox"
                              [checked]="isStreamSelected(stream.id)"
                              (change)="toggleStream(stream.id)"
                            />
                            <span class="checkbox-text">{{ stream.name }} ({{ stream.code }})</span>
                          </label>
                        }
                      </div>
                    </div>
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
                      placeholder="Optional description about this class..."
                    ></textarea>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeClassModal()">
                Cancel
              </button>
              <button type="button" class="btn btn--primary" [disabled]="saving()" (click)="saveClass()">
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ editingClass() ? 'Saving...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid" [class.fa-check]="editingClass()" [class.fa-plus]="!editingClass()"></i>
                  {{ editingClass() ? 'Save Changes' : 'Create Class' }}
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
                Delete Class
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
                  <strong>"{{ classToDelete()?.name }}"</strong>?
                </p>
                <p class="delete-warning">
                  This action cannot be undone. Classes with sections cannot be deleted.
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
                (click)="deleteClass()"
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
      background: #eef2ff;
      color: #4f46e5;
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
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
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
      border-color: #4f46e5;
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

    .class-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #eef2ff;
      color: #4f46e5;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

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

    .description {
      font-size: 0.75rem;
      color: #64748b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 200px;
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

    .level-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .level-nursery { background: #fef3c7; color: #92400e; }
    .level-primary { background: #dbeafe; color: #1e40af; }
    .level-middle { background: #dcfce7; color: #166534; }
    .level-secondary { background: #e0e7ff; color: #3730a3; }
    .level-senior_secondary { background: #fae8ff; color: #86198f; }

    .stream-tags {
      display: flex;
      gap: 0.25rem;
      flex-wrap: wrap;
    }

    .stream-tag {
      display: inline-flex;
      padding: 0.125rem 0.375rem;
      background: #f1f5f9;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      font-weight: 500;
      color: #475569;
    }

    .text-muted {
      color: #9ca3af;
      font-size: 0.875rem;
    }

    .sections-cell {
      text-align: center;
    }

    .section-count {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      height: 1.5rem;
      padding: 0 0.5rem;
      background: #f1f5f9;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
      color: #475569;
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
      color: #4f46e5;
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
      color: #6366f1;
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

    .form-hint {
      font-size: 0.75rem;
      color: #9ca3af;
    }

    .form-group input,
    .form-group select,
    .form-group textarea {
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: border-color 0.2s, box-shadow 0.2s;
    }

    .form-group input:focus,
    .form-group select:focus,
    .form-group textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .form-group input:disabled,
    .form-group select:disabled {
      background: #f8fafc;
      color: #64748b;
      cursor: not-allowed;
    }

    .checkbox-group {
      justify-content: flex-end;
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      cursor: pointer;
    }

    .checkbox-label input[type="checkbox"] {
      width: 1.125rem;
      height: 1.125rem;
      accent-color: #4f46e5;
    }

    .checkbox-text {
      font-size: 0.875rem;
      color: #374151;
    }

    .stream-checkboxes {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
      gap: 0.75rem;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
    }

    .stream-checkbox {
      padding: 0.5rem 0.75rem;
      background: white;
      border-radius: 0.375rem;
      border: 1px solid #e2e8f0;
      transition: border-color 0.2s, background 0.2s;
    }

    .stream-checkbox:hover {
      border-color: #c7d2fe;
      background: #eef2ff;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
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
          color: #4f46e5;
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
      background: #4f46e5;
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

      &:hover:not(:disabled) { background: #4338ca; }
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
        border-color: #4f46e5;
        background: white;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
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

      .checkbox-group {
        justify-content: flex-start;
      }

      .stream-checkboxes {
        grid-template-columns: 1fr;
      }

      .form-actions {
        flex-direction: column-reverse;
      }

      .form-actions .btn {
        width: 100%;
        justify-content: center;
      }

      .delete-actions {
        flex-direction: column-reverse;
      }

      .delete-actions .btn {
        width: 100%;
        justify-content: center;
      }
    }
  `],
})
export class ClassesComponent implements OnInit {
  private classService = inject(ClassService);
  private streamService = inject(StreamService);
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);

  // Constants
  classLevels = CLASS_LEVELS;

  // Data signals
  classes = signal<Class[]>([]);
  streams = signal<Stream[]>([]);
  branches = signal<Branch[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');
  levelFilter = signal<string>('');

  // Modal state
  showClassModal = signal(false);
  showDeleteModal = signal(false);
  editingClass = signal<Class | null>(null);
  classToDelete = signal<Class | null>(null);

  // Form data
  formData: {
    name: string;
    code: string;
    branchId: string;
    level: ClassLevel | '';
    displayOrder: number;
    description: string;
    hasStreams: boolean;
    streamIds: string[];
  } = this.getEmptyFormData();

  // Computed modal title
  modalTitle = computed(() => {
    const editing = this.editingClass();
    return editing ? `Edit Class: ${editing.name}` : 'Create New Class';
  });

  // Computed filtered list
  filteredClasses = computed(() => {
    let result = this.classes();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();
    const level = this.levelFilter();

    if (term) {
      result = result.filter(
        cls =>
          cls.name.toLowerCase().includes(term) ||
          cls.code.toLowerCase().includes(term)
      );
    }

    if (status === 'active') {
      result = result.filter(cls => cls.isActive);
    } else if (status === 'inactive') {
      result = result.filter(cls => !cls.isActive);
    }

    if (level) {
      result = result.filter(cls => cls.level === level);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadClasses();
    this.loadStreams();
    this.loadBranches();
  }

  loadClasses(): void {
    this.loading.set(true);
    this.error.set(null);

    this.classService.getClasses().subscribe({
      next: classes => {
        this.classes.set(classes);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load classes. Please try again.');
        this.loading.set(false);
      },
    });
  }

  loadStreams(): void {
    this.streamService.getStreams(true).subscribe({
      next: streams => this.streams.set(streams),
      error: () => console.error('Failed to load streams'),
    });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: branches => this.branches.set(branches),
      error: () => console.error('Failed to load branches'),
    });
  }

  openCreateModal(): void {
    this.editingClass.set(null);
    this.formData = this.getEmptyFormData();
    // Set default branch if only one exists
    const branchList = this.branches();
    if (branchList.length === 1) {
      this.formData.branchId = branchList[0].id;
    }
    this.showClassModal.set(true);
  }

  editClass(cls: Class): void {
    this.editingClass.set(cls);
    this.formData = {
      name: cls.name,
      code: cls.code,
      branchId: cls.branchId,
      level: cls.level || '',
      displayOrder: cls.displayOrder,
      description: cls.description || '',
      hasStreams: cls.hasStreams,
      streamIds: cls.streams?.map(s => s.id) || [],
    };
    this.showClassModal.set(true);
  }

  closeClassModal(): void {
    this.showClassModal.set(false);
    this.editingClass.set(null);
  }

  saveClass(): void {
    if (!this.formData.name || !this.formData.code || !this.formData.branchId) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    const editing = this.editingClass();

    const data: CreateClassRequest | UpdateClassRequest = {
      name: this.formData.name,
      code: this.formData.code,
      displayOrder: this.formData.displayOrder,
      description: this.formData.description || undefined,
      hasStreams: this.formData.hasStreams,
      streamIds: this.formData.hasStreams ? this.formData.streamIds : undefined,
    };

    if (this.formData.level) {
      data.level = this.formData.level as ClassLevel;
    }

    if (!editing) {
      (data as CreateClassRequest).branchId = this.formData.branchId;
    }

    const operation = editing
      ? this.classService.updateClass(editing.id, data as UpdateClassRequest)
      : this.classService.createClass(data as CreateClassRequest);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Class updated successfully' : 'Class created successfully'
        );
        this.closeClassModal();
        this.loadClasses();
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || (editing ? 'Failed to update class' : 'Failed to create class');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  toggleStatus(cls: Class): void {
    this.classService.updateClass(cls.id, { isActive: !cls.isActive }).subscribe({
      next: () => {
        this.toastService.success(
          `Class ${cls.isActive ? 'deactivated' : 'activated'} successfully`
        );
        this.loadClasses();
      },
      error: () => {
        this.toastService.error('Failed to update class status');
      },
    });
  }

  confirmDelete(cls: Class): void {
    this.classToDelete.set(cls);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.classToDelete.set(null);
  }

  deleteClass(): void {
    const cls = this.classToDelete();
    if (!cls) return;

    this.deleting.set(true);

    this.classService.deleteClass(cls.id).subscribe({
      next: () => {
        this.toastService.success('Class deleted successfully');
        this.closeDeleteModal();
        this.loadClasses();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.message || 'Failed to delete class. It may have sections.';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
    this.levelFilter.set('');
  }

  getLevelLabel(level: string): string {
    const found = CLASS_LEVELS.find(l => l.value === level);
    return found ? found.label.split(' ')[0] : level;
  }

  isStreamSelected(streamId: string): boolean {
    return this.formData.streamIds.includes(streamId);
  }

  toggleStream(streamId: string): void {
    const index = this.formData.streamIds.indexOf(streamId);
    if (index === -1) {
      this.formData.streamIds.push(streamId);
    } else {
      this.formData.streamIds.splice(index, 1);
    }
  }

  private getEmptyFormData() {
    return {
      name: '',
      code: '',
      branchId: '',
      level: '' as ClassLevel | '',
      displayOrder: 0,
      description: '',
      hasStreams: false,
      streamIds: [] as string[],
    };
  }
}
