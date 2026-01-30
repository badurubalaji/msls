import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ExamService } from '../exam.service';
import { ExamType, CreateExamTypeRequest, EVALUATION_TYPES, EXAM_TYPE_PRESETS } from '../exam.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'app-exam-type-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-file-signature"></i>
          </div>
          <div class="header-text">
            <h1>Exam Types</h1>
            <p>Configure different types of examinations and their weightage</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" (click)="showHelpPanel.set(!showHelpPanel())">
            <i class="fa-solid fa-circle-question"></i>
            <span class="btn-text">Help</span>
          </button>
          <button class="btn btn-primary" (click)="openCreateModal()">
            <i class="fa-solid fa-plus"></i>
            <span class="btn-text">Add Exam Type</span>
          </button>
        </div>
      </div>

      <!-- Help Panel -->
      @if (showHelpPanel()) {
        <div class="help-panel">
          <div class="help-header">
            <h3><i class="fa-solid fa-graduation-cap"></i> Understanding Exam Types</h3>
            <button class="help-close" (click)="showHelpPanel.set(false)">
              <i class="fa-solid fa-xmark"></i>
            </button>
          </div>
          <div class="help-grid">
            <div class="help-card">
              <div class="help-card-icon blue">
                <i class="fa-solid fa-clipboard-list"></i>
              </div>
              <h4>What are Exam Types?</h4>
              <p>Exam types define assessment categories like Unit Tests, Mid-Terms, Finals, Practicals, and Projects. Each type can have its own evaluation method and contribution to final grades.</p>
            </div>
            <div class="help-card">
              <div class="help-card-icon purple">
                <i class="fa-solid fa-star"></i>
              </div>
              <h4>Evaluation Types</h4>
              <ul>
                <li><strong>Marks-based:</strong> Numeric scores (85/100) for traditional exams</li>
                <li><strong>Grade-based:</strong> Letter grades (A+, A, B) for subjective assessments</li>
              </ul>
            </div>
            <div class="help-card">
              <div class="help-card-icon green">
                <i class="fa-solid fa-scale-balanced"></i>
              </div>
              <h4>Weightage</h4>
              <p>Determines contribution to final result. Example: UT(20%) + Mid(30%) + Final(50%) = Weighted average for final marks.</p>
            </div>
            <div class="help-card">
              <div class="help-card-icon orange">
                <i class="fa-solid fa-lightbulb"></i>
              </div>
              <h4>Tips</h4>
              <ul>
                <li>Ensure total weightage equals 100%</li>
                <li>Use unique codes for easy identification</li>
                <li>Deactivate instead of delete to preserve data</li>
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
            placeholder="Search exam types..."
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
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="typeFilter()"
            (ngModelChange)="typeFilter.set($event)"
          >
            <option value="">All Evaluation Types</option>
            @for (type of evaluationTypes; track type.value) {
              <option [value]="type.value">{{ type.label }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Stats Summary -->
      <div class="stats-bar">
        <div class="stat-item">
          <span class="stat-value">{{ examTypes().length }}</span>
          <span class="stat-label">Total Types</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ activeCount() }}</span>
          <span class="stat-label">Active</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ totalWeightage() }}%</span>
          <span class="stat-label" [class.warning]="totalWeightage() !== 100">Total Weightage</span>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading exam types...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadExamTypes()">
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
                  <th style="width: 50px;">Order</th>
                  <th>Exam Type</th>
                  <th class="hide-mobile">Code</th>
                  <th class="hide-tablet">Evaluation</th>
                  <th class="hide-tablet">Max Marks</th>
                  <th>Weightage</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 140px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (examType of filteredExamTypes(); track examType.id; let i = $index; let first = $first; let last = $last) {
                  <tr>
                    <td class="order-cell">
                      <div class="order-buttons">
                        <button
                          class="order-btn"
                          [disabled]="first"
                          (click)="moveUp(i)"
                          title="Move Up"
                        >
                          <i class="fa-solid fa-chevron-up"></i>
                        </button>
                        <span class="order-number">{{ examType.displayOrder }}</span>
                        <button
                          class="order-btn"
                          [disabled]="last"
                          (click)="moveDown(i)"
                          title="Move Down"
                        >
                          <i class="fa-solid fa-chevron-down"></i>
                        </button>
                      </div>
                    </td>
                    <td class="name-cell">
                      <div class="name-wrapper">
                        <div class="type-icon" [class]="'eval-' + examType.evaluationType">
                          <i [class]="examType.evaluationType === 'marks' ? 'fa-solid fa-calculator' : 'fa-solid fa-star'"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ examType.name }}</span>
                          <span class="code-mobile">{{ examType.code }}</span>
                          @if (examType.description) {
                            <span class="description">{{ examType.description }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td class="code-cell hide-mobile">
                      <span class="code-badge">{{ examType.code }}</span>
                    </td>
                    <td class="hide-tablet">
                      <span class="eval-badge" [class]="'eval-' + examType.evaluationType">
                        <i [class]="examType.evaluationType === 'marks' ? 'fa-solid fa-calculator' : 'fa-solid fa-star'"></i>
                        {{ getEvaluationLabel(examType.evaluationType) }}
                      </span>
                    </td>
                    <td class="marks-cell hide-tablet">
                      <span class="marks-value">{{ examType.defaultMaxMarks }}</span>
                    </td>
                    <td class="weightage-cell">
                      <div class="weightage-bar">
                        <div class="weightage-fill" [style.width.%]="examType.weightage"></div>
                      </div>
                      <span class="weightage-value">{{ examType.weightage }}%</span>
                    </td>
                    <td>
                      <span
                        class="badge"
                        [class.badge-green]="examType.isActive"
                        [class.badge-gray]="!examType.isActive"
                      >
                        {{ examType.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        title="Edit"
                        (click)="editExamType(examType)"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn"
                        title="Toggle Status"
                        (click)="toggleStatus(examType)"
                      >
                        <i
                          class="fa-solid"
                          [class.fa-toggle-on]="examType.isActive"
                          [class.fa-toggle-off]="!examType.isActive"
                          [class.text-green]="examType.isActive"
                        ></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(examType)"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="8" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No exam types found</p>
                        @if (searchTerm() || statusFilter() !== 'all' || typeFilter()) {
                          <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                            Clear Filters
                          </button>
                        } @else {
                          <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                            <i class="fa-solid fa-plus"></i>
                            Create First Exam Type
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
                <i class="fa-solid fa-file-signature"></i>
                {{ modalTitle() }}
              </h3>
              <button type="button" class="modal__close" (click)="closeFormModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <!-- Quick Presets (Create mode only) -->
              @if (!editingExamType()) {
                <div class="presets-section">
                  <h4 class="presets-title">
                    <i class="fa-solid fa-magic-wand-sparkles"></i>
                    Quick Presets
                  </h4>
                  <div class="presets-list">
                    @for (preset of presets; track preset.code) {
                      <button
                        type="button"
                        class="preset-btn"
                        (click)="applyPreset(preset)"
                      >
                        {{ preset.name }}
                      </button>
                    }
                  </div>
                </div>
              }

              <form class="form">
                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-info-circle"></i>
                    Basic Information
                  </h4>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="name">Name <span class="required">*</span></label>
                      <input
                        type="text"
                        id="name"
                        [(ngModel)]="formData.name"
                        name="name"
                        class="form-input"
                        placeholder="e.g., Unit Test"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label for="code">Code <span class="required">*</span></label>
                      <input
                        type="text"
                        id="code"
                        [(ngModel)]="formData.code"
                        name="code"
                        class="form-input uppercase"
                        placeholder="e.g., UT"
                        required
                      />
                    </div>
                  </div>
                  <div class="form-group">
                    <label for="description">Description</label>
                    <textarea
                      id="description"
                      [(ngModel)]="formData.description"
                      name="description"
                      class="form-textarea"
                      rows="2"
                      placeholder="Optional description..."
                    ></textarea>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-sliders"></i>
                    Evaluation Settings
                  </h4>
                  <div class="form-group">
                    <label>Evaluation Type <span class="required">*</span></label>
                    <div class="radio-group">
                      @for (evalType of evaluationTypes; track evalType.value) {
                        <label class="radio-card" [class.selected]="formData.evaluationType === evalType.value">
                          <input
                            type="radio"
                            [(ngModel)]="formData.evaluationType"
                            name="evaluationType"
                            [value]="evalType.value"
                          />
                          <div class="radio-card-content">
                            <i [class]="evalType.icon"></i>
                            <span class="radio-label">{{ evalType.label }}</span>
                            <span class="radio-desc">{{ evalType.value === 'marks' ? 'Numeric scores' : 'Letter grades' }}</span>
                          </div>
                        </label>
                      }
                    </div>
                  </div>
                  <div class="form-row">
                    <div class="form-group">
                      <label for="defaultMaxMarks">Default Max Marks <span class="required">*</span></label>
                      <input
                        type="number"
                        id="defaultMaxMarks"
                        [(ngModel)]="formData.defaultMaxMarks"
                        name="defaultMaxMarks"
                        class="form-input"
                        min="1"
                        placeholder="100"
                        required
                      />
                    </div>
                    <div class="form-group">
                      <label for="defaultPassingMarks">Default Passing Marks</label>
                      <input
                        type="number"
                        id="defaultPassingMarks"
                        [(ngModel)]="formData.defaultPassingMarks"
                        name="defaultPassingMarks"
                        class="form-input"
                        min="0"
                        placeholder="35"
                      />
                    </div>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-scale-balanced"></i>
                    Weightage
                  </h4>
                  <div class="form-group">
                    <label for="weightage">Weightage (%) <span class="required">*</span></label>
                    <div class="weightage-input-group">
                      <input
                        type="range"
                        id="weightage"
                        [(ngModel)]="formData.weightage"
                        name="weightage"
                        class="weightage-slider"
                        min="0"
                        max="100"
                        step="5"
                      />
                      <span class="weightage-display">{{ formData.weightage }}%</span>
                    </div>
                    <p class="form-hint">
                      Contribution to final result calculation. Current total: {{ totalWeightage() }}%
                    </p>
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
                (click)="saveExamType()"
              >
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ editingExamType() ? 'Saving...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid" [class.fa-check]="editingExamType()" [class.fa-plus]="!editingExamType()"></i>
                  {{ editingExamType() ? 'Save Changes' : 'Create Exam Type' }}
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
                Delete Exam Type
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
                  <strong>"{{ examTypeToDelete()?.name }}"</strong>?
                </p>
                <p class="delete-warning">
                  This action cannot be undone. Consider deactivating instead to preserve historical data.
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
                (click)="deleteExamType()"
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
      background: linear-gradient(135deg, #818cf8, #6366f1);
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
      background: linear-gradient(135deg, #eef2ff, #e0e7ff);
      border: 1px solid #c7d2fe;
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
      color: #3730a3;
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

    .help-card p {
      margin: 0;
      font-size: 0.8125rem;
      color: #64748b;
      line-height: 1.5;
    }

    .help-card ul {
      margin: 0;
      padding-left: 1rem;
      font-size: 0.8125rem;
      color: #64748b;
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
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
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
      border-color: #6366f1;
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

    .stat-label.warning {
      color: #dc2626;
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
      border-top-color: #6366f1;
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

    /* Order Cell */
    .order-cell {
      text-align: center;
    }

    .order-buttons {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.125rem;
    }

    .order-btn {
      width: 1.5rem;
      height: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border: none;
      background: transparent;
      color: #9ca3af;
      cursor: pointer;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      transition: all 0.2s;
    }

    .order-btn:hover:not(:disabled) {
      background: #e2e8f0;
      color: #6366f1;
    }

    .order-btn:disabled {
      opacity: 0.3;
      cursor: not-allowed;
    }

    .order-number {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
    }

    /* Name Cell */
    .name-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .type-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .type-icon.eval-marks { background: #dbeafe; color: #1d4ed8; }
    .type-icon.eval-grade { background: #fef3c7; color: #d97706; }

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

    .eval-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .eval-badge.eval-marks { background: #dbeafe; color: #1d4ed8; }
    .eval-badge.eval-grade { background: #fef3c7; color: #d97706; }

    .marks-value {
      font-family: monospace;
      font-weight: 500;
    }

    /* Weightage Cell */
    .weightage-cell {
      min-width: 120px;
    }

    .weightage-bar {
      width: 100%;
      height: 6px;
      background: #e2e8f0;
      border-radius: 3px;
      overflow: hidden;
      margin-bottom: 0.25rem;
    }

    .weightage-fill {
      height: 100%;
      background: linear-gradient(90deg, #6366f1, #818cf8);
      border-radius: 3px;
      transition: width 0.3s ease;
    }

    .weightage-value {
      font-size: 0.75rem;
      font-weight: 500;
      color: #6366f1;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }

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

    .action-btn:hover {
      background: #f1f5f9;
      color: #6366f1;
    }

    .action-btn--danger:hover {
      background: #fef2f2;
      color: #dc2626;
    }

    .text-green {
      color: #16a34a !important;
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
      background: linear-gradient(135deg, #6366f1, #4f46e5);
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: linear-gradient(135deg, #4f46e5, #4338ca);
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

      h3 {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;

        i {
          color: #6366f1;
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
      background: linear-gradient(135deg, #6366f1, #4f46e5);
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

      &:hover:not(:disabled) { background: linear-gradient(135deg, #4f46e5, #4338ca); }
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

    /* Presets */
    .presets-section {
      background: #f0f9ff;
      border: 1px solid #bae6fd;
      border-radius: 0.75rem;
      padding: 1rem;
      margin-bottom: 1.5rem;
    }

    .presets-title {
      margin: 0 0 0.75rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #0369a1;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .presets-list {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
    }

    .preset-btn {
      padding: 0.375rem 0.75rem;
      font-size: 0.8125rem;
      background: white;
      border: 1px solid #bae6fd;
      border-radius: 9999px;
      color: #0369a1;
      cursor: pointer;
      transition: all 0.2s;
    }

    .preset-btn:hover {
      background: #0ea5e9;
      border-color: #0ea5e9;
      color: white;
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

      i {
        color: #6366f1;
      }
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
        border-color: #6366f1;
        background: white;
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
      }

      &::placeholder {
        color: #94a3b8;
      }
    }

    .form-input.uppercase {
      text-transform: uppercase;
    }

    .form-textarea {
      resize: vertical;
      min-height: 60px;
    }

    .form-hint {
      margin: 0.5rem 0 0;
      font-size: 0.75rem;
      color: #64748b;
    }

    /* Radio Cards */
    .radio-group {
      display: flex;
      gap: 1rem;
    }

    .radio-card {
      flex: 1;
      cursor: pointer;
    }

    .radio-card input {
      display: none;
    }

    .radio-card-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.5rem;
      padding: 1rem;
      background: #f8fafc;
      border: 2px solid #e2e8f0;
      border-radius: 0.75rem;
      transition: all 0.2s;

      i {
        font-size: 1.5rem;
        color: #64748b;
      }
    }

    .radio-card.selected .radio-card-content {
      background: #eef2ff;
      border-color: #6366f1;

      i {
        color: #6366f1;
      }
    }

    .radio-card:hover .radio-card-content {
      border-color: #c7d2fe;
    }

    .radio-label {
      font-weight: 500;
      color: #1e293b;
    }

    .radio-desc {
      font-size: 0.75rem;
      color: #64748b;
    }

    /* Weightage Slider */
    .weightage-input-group {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .weightage-slider {
      flex: 1;
      height: 6px;
      -webkit-appearance: none;
      background: #e2e8f0;
      border-radius: 3px;
      outline: none;

      &::-webkit-slider-thumb {
        -webkit-appearance: none;
        width: 18px;
        height: 18px;
        background: linear-gradient(135deg, #6366f1, #4f46e5);
        border-radius: 50%;
        cursor: pointer;
        box-shadow: 0 2px 6px rgba(99, 102, 241, 0.4);
      }
    }

    .weightage-display {
      min-width: 50px;
      text-align: center;
      font-weight: 600;
      color: #6366f1;
      font-size: 1.125rem;
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

      .code-mobile {
        display: block;
      }

      .form-row {
        grid-template-columns: 1fr;
      }

      .radio-group {
        flex-direction: column;
      }

      .help-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class ExamTypeListComponent implements OnInit {
  private examService = inject(ExamService);
  private toastService = inject(ToastService);

  // Constants
  evaluationTypes = EVALUATION_TYPES;
  presets = EXAM_TYPE_PRESETS;

  // Data signals
  examTypes = signal<ExamType[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');
  typeFilter = signal<string>('');
  showHelpPanel = signal(false);

  // Modal state
  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingExamType = signal<ExamType | null>(null);
  examTypeToDelete = signal<ExamType | null>(null);

  // Form data
  formData = this.getEmptyFormData();

  // Computed values
  modalTitle = computed(() => {
    const editing = this.editingExamType();
    return editing ? `Edit: ${editing.name}` : 'Create New Exam Type';
  });

  filteredExamTypes = computed(() => {
    let result = this.examTypes();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();
    const type = this.typeFilter();

    if (term) {
      result = result.filter(
        examType =>
          examType.name.toLowerCase().includes(term) ||
          examType.code.toLowerCase().includes(term)
      );
    }

    if (status === 'active') {
      result = result.filter(examType => examType.isActive);
    } else if (status === 'inactive') {
      result = result.filter(examType => !examType.isActive);
    }

    if (type) {
      result = result.filter(examType => examType.evaluationType === type);
    }

    return result;
  });

  activeCount = computed(() => this.examTypes().filter(t => t.isActive).length);

  totalWeightage = computed(() =>
    this.examTypes()
      .filter(t => t.isActive)
      .reduce((sum, t) => sum + t.weightage, 0)
  );

  ngOnInit(): void {
    this.loadExamTypes();
  }

  loadExamTypes(): void {
    this.loading.set(true);
    this.error.set(null);

    this.examService.getExamTypes({}).subscribe({
      next: types => {
        this.examTypes.set(types);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load exam types. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingExamType.set(null);
    this.formData = this.getEmptyFormData();
    this.showFormModal.set(true);
  }

  editExamType(examType: ExamType): void {
    this.editingExamType.set(examType);
    this.formData = {
      name: examType.name,
      code: examType.code,
      description: examType.description || '',
      evaluationType: examType.evaluationType,
      defaultMaxMarks: examType.defaultMaxMarks,
      defaultPassingMarks: examType.defaultPassingMarks || 0,
      weightage: examType.weightage,
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingExamType.set(null);
  }

  applyPreset(preset: Partial<CreateExamTypeRequest>): void {
    this.formData = {
      ...this.formData,
      name: preset.name || '',
      code: preset.code || '',
      evaluationType: preset.evaluationType || 'marks',
      defaultMaxMarks: preset.defaultMaxMarks || 100,
      weightage: preset.weightage || 0,
    };
  }

  isFormValid(): boolean {
    return !!(
      this.formData.name &&
      this.formData.code &&
      this.formData.evaluationType &&
      this.formData.defaultMaxMarks > 0
    );
  }

  saveExamType(): void {
    if (!this.isFormValid()) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    const editing = this.editingExamType();

    const data: CreateExamTypeRequest = {
      name: this.formData.name,
      code: this.formData.code.toUpperCase(),
      description: this.formData.description || undefined,
      evaluationType: this.formData.evaluationType,
      defaultMaxMarks: this.formData.defaultMaxMarks,
      defaultPassingMarks: this.formData.defaultPassingMarks || undefined,
      weightage: this.formData.weightage,
    };

    const operation = editing
      ? this.examService.updateExamType(editing.id, data)
      : this.examService.createExamType(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Exam type updated successfully' : 'Exam type created successfully'
        );
        this.closeFormModal();
        this.loadExamTypes();
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.detail || (editing ? 'Failed to update exam type' : 'Failed to create exam type');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  toggleStatus(examType: ExamType): void {
    this.examService.toggleExamTypeActive(examType.id, !examType.isActive).subscribe({
      next: () => {
        this.toastService.success(
          `Exam type ${examType.isActive ? 'deactivated' : 'activated'} successfully`
        );
        this.loadExamTypes();
      },
      error: () => {
        this.toastService.error('Failed to update exam type status');
      },
    });
  }

  confirmDelete(examType: ExamType): void {
    this.examTypeToDelete.set(examType);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.examTypeToDelete.set(null);
  }

  deleteExamType(): void {
    const examType = this.examTypeToDelete();
    if (!examType) return;

    this.deleting.set(true);

    this.examService.deleteExamType(examType.id).subscribe({
      next: () => {
        this.toastService.success('Exam type deleted successfully');
        this.closeDeleteModal();
        this.loadExamTypes();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.detail || 'Failed to delete exam type';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  moveUp(index: number): void {
    if (index <= 0) return;
    this.swapItems(index, index - 1);
  }

  moveDown(index: number): void {
    const items = this.filteredExamTypes();
    if (index >= items.length - 1) return;
    this.swapItems(index, index + 1);
  }

  private swapItems(fromIndex: number, toIndex: number): void {
    const items = [...this.filteredExamTypes()];
    [items[fromIndex], items[toIndex]] = [items[toIndex], items[fromIndex]];

    const orderUpdates = items.map((item, index) => ({
      id: item.id,
      displayOrder: index + 1,
    }));

    // Optimistic update
    const allItems = [...this.examTypes()];
    const fromItem = items[toIndex];
    const toItem = items[fromIndex];
    const fromIdx = allItems.findIndex(i => i.id === fromItem.id);
    const toIdx = allItems.findIndex(i => i.id === toItem.id);
    if (fromIdx !== -1 && toIdx !== -1) {
      [allItems[fromIdx], allItems[toIdx]] = [allItems[toIdx], allItems[fromIdx]];
      this.examTypes.set(allItems);
    }

    this.examService.updateDisplayOrder({ items: orderUpdates }).subscribe({
      error: () => {
        this.toastService.error('Failed to update display order');
        this.loadExamTypes();
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
    this.typeFilter.set('');
  }

  getEvaluationLabel(type: string): string {
    const found = EVALUATION_TYPES.find(t => t.value === type);
    return found?.label || type;
  }

  private getEmptyFormData() {
    return {
      name: '',
      code: '',
      description: '',
      evaluationType: 'marks' as 'marks' | 'grade',
      defaultMaxMarks: 100,
      defaultPassingMarks: 0,
      weightage: 0,
    };
  }
}
