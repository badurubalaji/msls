/**
 * MSLS Applications Management Component
 *
 * Main component for managing admission applications - displays a list with filtering,
 * stage updates, and navigation to detailed views.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';

import { MslsModalComponent, SelectOption } from '../../../shared/components';
import { ToastService } from '../../../shared/services';
import {
  AdmissionApplication,
  ApplicationStage,
  APPLICATION_STAGE_CONFIG,
  ApplicationFilterParams,
} from './application.model';
import { ApplicationService } from './application.service';
import { CLASS_NAMES } from '../sessions/admission-session.model';

@Component({
  selector: 'msls-applications',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule, MslsModalComponent],
  template: `
    <div class="applications-page">
      <div class="applications-card">
        <!-- Header -->
        <div class="applications-header">
          <div class="applications-header__left">
            <h1 class="applications-header__title">Admission Applications</h1>
            <p class="applications-header__subtitle">Manage and review student admission applications</p>
          </div>
          <div class="applications-header__right">
            <button class="btn btn-primary" (click)="createApplication()">
              <i class="fa-solid fa-plus"></i>
              New Application
            </button>
          </div>
        </div>

        <!-- Filters -->
        <div class="filters-section">
          <div class="filters-row">
            <!-- Search -->
            <div class="search-input">
              <i class="fa-solid fa-magnifying-glass search-icon"></i>
              <input
                type="text"
                placeholder="Search by App#, Name..."
                [ngModel]="searchTerm()"
                (ngModelChange)="onSearchChange($event)"
                class="search-field"
              />
            </div>

            <!-- Stage Filter -->
            <div class="filter-group">
              <label class="filter-label">Stage</label>
              <select class="filter-select" [ngModel]="selectedStage()" (ngModelChange)="onStageFilterChange($event)">
                <option value="">All Stages</option>
                @for (stage of stageOptions; track stage.value) {
                  <option [value]="stage.value">{{ stage.label }}</option>
                }
              </select>
            </div>

            <!-- Class Filter -->
            <div class="filter-group">
              <label class="filter-label">Class</label>
              <select class="filter-select" [ngModel]="selectedClass()" (ngModelChange)="onClassFilterChange($event)">
                <option value="">All Classes</option>
                @for (className of classOptions; track className) {
                  <option [value]="className">{{ className }}</option>
                }
              </select>
            </div>

            <!-- Clear Filters -->
            @if (hasActiveFilters()) {
              <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                <i class="fa-solid fa-times"></i>
                Clear Filters
              </button>
            }
          </div>

          <!-- Stage Quick Filters -->
          <div class="stage-tabs">
            <button
              class="stage-tab"
              [class.stage-tab--active]="selectedStage() === ''"
              (click)="onStageFilterChange('')"
            >
              All
              <span class="stage-tab__count">{{ applications().length }}</span>
            </button>
            @for (stage of quickStageFilters; track stage.value) {
              <button
                class="stage-tab"
                [class.stage-tab--active]="selectedStage() === stage.value"
                (click)="onStageFilterChange(stage.value)"
              >
                {{ stage.label }}
                <span class="stage-tab__count">{{ getStageCount(stage.value) }}</span>
              </button>
            }
          </div>
        </div>

        <!-- Loading State -->
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading applications...</p>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <p>{{ error() }}</p>
            <button class="btn btn-secondary" (click)="loadApplications()">Retry</button>
          </div>
        } @else {
          <!-- Table -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 130px">App #</th>
                  <th>Student Name</th>
                  <th>Class</th>
                  <th style="width: 160px">Stage</th>
                  <th style="width: 120px">Date</th>
                  <th style="width: 180px; text-align: right">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (app of filteredApplications(); track app.id) {
                  <tr (click)="viewApplication(app)" class="clickable-row">
                    <td class="app-number-cell">{{ app.applicationNumber }}</td>
                    <td class="name-cell">
                      <div class="student-info">
                        <span class="student-name">{{ app.firstName }} {{ app.lastName }}</span>
                        @if (app.parents && app.parents.length > 0) {
                          <span class="parent-name">{{ getParentLabel(app) }}</span>
                        }
                      </div>
                    </td>
                    <td class="class-cell">{{ app.classApplying }}</td>
                    <td>
                      <span
                        class="stage-badge"
                        [class]="'stage-badge--' + getStageConfig(app.currentStage).variant"
                      >
                        <i [class]="getStageConfig(app.currentStage).icon"></i>
                        {{ getStageConfig(app.currentStage).label }}
                      </span>
                    </td>
                    <td class="date-cell">{{ formatDate(app.submittedAt || app.createdAt) }}</td>
                    <td class="actions-cell" (click)="$event.stopPropagation()">
                      <button class="action-btn" (click)="viewApplication(app)" title="View details">
                        <i class="fa-regular fa-eye"></i>
                      </button>
                      @if (canEditApplication(app)) {
                        <button class="action-btn" (click)="editApplication(app)" title="Edit application">
                          <i class="fa-regular fa-pen-to-square"></i>
                        </button>
                      }
                      <button
                        class="action-btn action-btn--primary"
                        (click)="openStageModal(app)"
                        title="Update stage"
                      >
                        <i class="fa-solid fa-arrow-right-arrow-left"></i>
                      </button>
                      @if (canDeleteApplication(app)) {
                        <button
                          class="action-btn action-btn--danger"
                          (click)="confirmDelete(app)"
                          title="Delete application"
                        >
                          <i class="fa-regular fa-trash-can"></i>
                        </button>
                      }
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="6" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No applications found</p>
                        @if (!hasActiveFilters()) {
                          <button class="btn btn-primary btn-sm" (click)="createApplication()">
                            <i class="fa-solid fa-plus"></i>
                            Create First Application
                          </button>
                        } @else {
                          <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                            <i class="fa-solid fa-times"></i>
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

      <!-- Update Stage Modal -->
      <msls-modal
        [isOpen]="showStageModal()"
        title="Update Application Stage"
        size="md"
        (closed)="closeStageModal()"
      >
        <div class="stage-modal-content">
          @if (selectedApplication()) {
            <div class="current-stage-info">
              <p class="modal-label">Current Stage</p>
              <span
                class="stage-badge stage-badge--lg"
                [class]="'stage-badge--' + getStageConfig(selectedApplication()!.currentStage).variant"
              >
                <i [class]="getStageConfig(selectedApplication()!.currentStage).icon"></i>
                {{ getStageConfig(selectedApplication()!.currentStage).label }}
              </span>
            </div>

            <div class="new-stage-section">
              <label class="modal-label">New Stage</label>
              <select class="modal-select" [ngModel]="newStage()" (ngModelChange)="newStage.set($event)">
                @for (stage of stageOptions; track stage.value) {
                  <option [value]="stage.value" [disabled]="stage.value === selectedApplication()!.currentStage">
                    {{ stage.label }}
                  </option>
                }
              </select>
            </div>

            <div class="remarks-section">
              <label class="modal-label">Remarks (Optional)</label>
              <textarea
                class="modal-textarea"
                rows="3"
                placeholder="Add remarks for this stage change..."
                [ngModel]="stageRemarks()"
                (ngModelChange)="stageRemarks.set($event)"
              ></textarea>
            </div>

            <div class="modal-actions">
              <button class="btn btn-secondary" (click)="closeStageModal()">Cancel</button>
              <button
                class="btn btn-primary"
                [disabled]="!newStage() || newStage() === selectedApplication()!.currentStage || updatingStage()"
                (click)="updateStage()"
              >
                @if (updatingStage()) {
                  <div class="btn-spinner"></div>
                  Updating...
                } @else {
                  Update Stage
                }
              </button>
            </div>
          }
        </div>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Application"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete application
            <strong>{{ applicationToDelete()?.applicationNumber }}</strong
            >?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteApplication()">
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [
    `
      .applications-page {
        padding: 1.5rem;
        max-width: 1400px;
        margin: 0 auto;
      }

      .applications-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 1rem;
        padding: 1.5rem;
      }

      /* Header */
      .applications-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.5rem;
        padding-bottom: 1.5rem;
        border-bottom: 1px solid #e2e8f0;
        flex-wrap: wrap;
        gap: 1rem;
      }

      .applications-header__title {
        font-size: 1.5rem;
        font-weight: 700;
        color: #0f172a;
        margin: 0 0 0.375rem 0;
      }

      .applications-header__subtitle {
        font-size: 0.875rem;
        color: #64748b;
        margin: 0;
      }

      /* Filters */
      .filters-section {
        margin-bottom: 1.5rem;
      }

      .filters-row {
        display: flex;
        gap: 1rem;
        align-items: flex-end;
        flex-wrap: wrap;
        margin-bottom: 1rem;
      }

      .search-input {
        position: relative;
        flex: 1;
        min-width: 200px;
        max-width: 300px;
      }

      .search-icon {
        position: absolute;
        left: 0.875rem;
        top: 50%;
        transform: translateY(-50%);
        color: #94a3b8;
        font-size: 0.875rem;
      }

      .search-field {
        width: 100%;
        padding: 0.625rem 0.875rem 0.625rem 2.5rem;
        font-size: 0.875rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #ffffff;
        color: #0f172a;
        transition: all 0.15s;
      }

      .search-field::placeholder {
        color: #94a3b8;
      }

      .search-field:focus {
        outline: none;
        border-color: #4f46e5;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
      }

      .filter-group {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
      }

      .filter-label {
        font-size: 0.75rem;
        font-weight: 500;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      .filter-select {
        padding: 0.625rem 2rem 0.625rem 0.875rem;
        font-size: 0.875rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #ffffff url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3E%3Cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3E%3C/svg%3E")
          right 0.5rem center / 1.5rem no-repeat;
        color: #0f172a;
        min-width: 150px;
        cursor: pointer;
        appearance: none;
      }

      .filter-select:focus {
        outline: none;
        border-color: #4f46e5;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
      }

      /* Stage Tabs */
      .stage-tabs {
        display: flex;
        gap: 0.5rem;
        flex-wrap: wrap;
        padding-top: 0.5rem;
        border-top: 1px solid #f1f5f9;
      }

      .stage-tab {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        font-size: 0.8125rem;
        font-weight: 500;
        color: #64748b;
        background: transparent;
        border: 1px solid transparent;
        border-radius: 9999px;
        cursor: pointer;
        transition: all 0.15s;
      }

      .stage-tab:hover {
        background: #f8fafc;
        color: #334155;
      }

      .stage-tab--active {
        background: #4f46e5;
        color: #ffffff;
        border-color: #4f46e5;
      }

      .stage-tab--active:hover {
        background: #4338ca;
        color: #ffffff;
      }

      .stage-tab__count {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 1.25rem;
        height: 1.25rem;
        padding: 0 0.375rem;
        font-size: 0.6875rem;
        font-weight: 600;
        border-radius: 9999px;
        background: rgba(0, 0, 0, 0.1);
      }

      .stage-tab--active .stage-tab__count {
        background: rgba(255, 255, 255, 0.2);
      }

      /* Buttons */
      .btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.625rem 1rem;
        font-size: 0.875rem;
        font-weight: 500;
        border-radius: 0.5rem;
        border: none;
        cursor: pointer;
        transition: all 0.15s;
      }

      .btn-sm {
        padding: 0.5rem 0.875rem;
        font-size: 0.8125rem;
      }

      .btn-primary {
        background: #4f46e5;
        color: #ffffff;
      }

      .btn-primary:hover:not(:disabled) {
        background: #4338ca;
      }

      .btn-primary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      .btn-secondary {
        background: #ffffff;
        color: #334155;
        border: 1px solid #e2e8f0;
      }

      .btn-secondary:hover {
        background: #f8fafc;
        border-color: #cbd5e1;
      }

      .btn-danger {
        background: #dc2626;
        color: #ffffff;
      }

      .btn-danger:hover:not(:disabled) {
        background: #b91c1c;
      }

      .btn-danger:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      .btn-spinner {
        width: 1rem;
        height: 1rem;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: #ffffff;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
      }

      /* Loading */
      .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
        gap: 1rem;
      }

      .spinner {
        width: 2rem;
        height: 2rem;
        border: 3px solid #e2e8f0;
        border-top-color: #4f46e5;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
      }

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }

      .loading-container p {
        color: #64748b;
        font-size: 0.875rem;
        margin: 0;
      }

      /* Error */
      .error-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
        gap: 1rem;
      }

      .error-container i {
        font-size: 2.5rem;
        color: #dc2626;
      }

      .error-container p {
        color: #64748b;
        font-size: 0.875rem;
        margin: 0;
      }

      /* Table */
      .table-container {
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        overflow: hidden;
      }

      .data-table {
        width: 100%;
        border-collapse: collapse;
      }

      .data-table thead {
        background: #f8fafc;
        border-bottom: 1px solid #e2e8f0;
      }

      .data-table th {
        padding: 0.75rem 1rem;
        text-align: left;
        font-size: 0.75rem;
        font-weight: 600;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      .data-table td {
        padding: 1rem;
        font-size: 0.875rem;
        border-bottom: 1px solid #f1f5f9;
        color: #334155;
      }

      .data-table tbody tr:last-child td {
        border-bottom: none;
      }

      .clickable-row {
        cursor: pointer;
        transition: background 0.15s;
      }

      .clickable-row:hover {
        background: #f8fafc;
      }

      .app-number-cell {
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 0.8125rem;
        color: #4f46e5;
        font-weight: 500;
      }

      .name-cell {
        font-weight: 500;
        color: #0f172a;
      }

      .student-info {
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
      }

      .student-name {
        font-weight: 500;
        color: #0f172a;
      }

      .parent-name {
        font-size: 0.75rem;
        color: #64748b;
        font-weight: 400;
      }

      .class-cell {
        color: #475569;
      }

      .date-cell {
        color: #64748b;
        font-size: 0.8125rem;
      }

      /* Stage Badge */
      .stage-badge {
        display: inline-flex;
        align-items: center;
        gap: 0.375rem;
        padding: 0.25rem 0.625rem;
        font-size: 0.75rem;
        font-weight: 500;
        border-radius: 9999px;
      }

      .stage-badge--lg {
        padding: 0.375rem 0.875rem;
        font-size: 0.875rem;
      }

      .stage-badge i {
        font-size: 0.625rem;
      }

      .stage-badge--lg i {
        font-size: 0.75rem;
      }

      .stage-badge--primary {
        background: #dbeafe;
        color: #1e40af;
      }

      .stage-badge--success {
        background: #dcfce7;
        color: #166534;
      }

      .stage-badge--warning {
        background: #fef3c7;
        color: #92400e;
      }

      .stage-badge--danger {
        background: #fee2e2;
        color: #991b1b;
      }

      .stage-badge--info {
        background: #e0e7ff;
        color: #3730a3;
      }

      .stage-badge--neutral {
        background: #f1f5f9;
        color: #475569;
      }

      /* Actions */
      .actions-cell {
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
      }

      .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2rem;
        height: 2rem;
        background: transparent;
        border: 1px solid #e2e8f0;
        border-radius: 0.375rem;
        color: #64748b;
        cursor: pointer;
        transition: all 0.15s;
      }

      .action-btn:hover {
        background: #f8fafc;
        border-color: #cbd5e1;
        color: #0f172a;
      }

      .action-btn--primary {
        border-color: #c7d2fe;
        color: #4f46e5;
      }

      .action-btn--primary:hover {
        background: #eef2ff;
        border-color: #a5b4fc;
        color: #4338ca;
      }

      .action-btn--danger:hover {
        background: #fef2f2;
        border-color: #fecaca;
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
        color: #94a3b8;
      }

      .empty-state i {
        font-size: 2.5rem;
      }

      .empty-state p {
        margin: 0;
        font-size: 0.875rem;
      }

      /* Modal Content */
      .stage-modal-content {
        padding: 0.5rem;
      }

      .current-stage-info {
        margin-bottom: 1.5rem;
        padding: 1rem;
        background: #f8fafc;
        border-radius: 0.5rem;
      }

      .modal-label {
        display: block;
        font-size: 0.75rem;
        font-weight: 500;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin-bottom: 0.5rem;
      }

      .new-stage-section {
        margin-bottom: 1.25rem;
      }

      .modal-select {
        width: 100%;
        padding: 0.75rem 2.5rem 0.75rem 0.875rem;
        font-size: 0.875rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #ffffff url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3E%3Cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3E%3C/svg%3E")
          right 0.75rem center / 1.5rem no-repeat;
        color: #0f172a;
        cursor: pointer;
        appearance: none;
      }

      .modal-select:focus {
        outline: none;
        border-color: #4f46e5;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
      }

      .remarks-section {
        margin-bottom: 1.5rem;
      }

      .modal-textarea {
        width: 100%;
        padding: 0.75rem;
        font-size: 0.875rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #ffffff;
        color: #0f172a;
        resize: vertical;
        font-family: inherit;
      }

      .modal-textarea:focus {
        outline: none;
        border-color: #4f46e5;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
      }

      .modal-textarea::placeholder {
        color: #94a3b8;
      }

      .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        padding-top: 1rem;
        border-top: 1px solid #e2e8f0;
      }

      /* Delete Confirmation */
      .delete-confirmation {
        text-align: center;
        padding: 1rem;
      }

      .delete-icon {
        width: 3.5rem;
        height: 3.5rem;
        margin: 0 auto 1rem;
        background: #fef2f2;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .delete-icon i {
        font-size: 1.5rem;
        color: #dc2626;
      }

      .delete-confirmation p {
        color: #475569;
        margin: 0 0 0.5rem 0;
      }

      .delete-confirmation strong {
        color: #dc2626;
      }

      .delete-warning {
        color: #dc2626 !important;
        font-size: 0.8125rem;
        font-weight: 500;
        padding: 0.75rem;
        background: #fef2f2;
        border-radius: 0.5rem;
        margin-top: 1rem !important;
      }

      .delete-actions {
        display: flex;
        justify-content: center;
        gap: 0.75rem;
        margin-top: 1.5rem;
      }

      /* Responsive */
      @media (max-width: 768px) {
        .applications-header {
          flex-direction: column;
          align-items: stretch;
        }

        .applications-header__right {
          width: 100%;
        }

        .applications-header__right .btn {
          width: 100%;
          justify-content: center;
        }

        .filters-row {
          flex-direction: column;
        }

        .search-input {
          max-width: none;
        }

        .filter-group {
          width: 100%;
        }

        .filter-select {
          width: 100%;
        }

        .table-container {
          overflow-x: auto;
        }

        .data-table {
          min-width: 800px;
        }
      }
    `,
  ],
})
export class ApplicationsComponent implements OnInit {
  private router = inject(Router);
  private applicationService = inject(ApplicationService);
  private toastService = inject(ToastService);

  // State signals
  applications = signal<AdmissionApplication[]>([]);
  loading = signal(true);
  error = signal<string | null>(null);

  // Filter signals
  searchTerm = signal('');
  selectedStage = signal<ApplicationStage | ''>('');
  selectedClass = signal('');

  // Modal state
  showStageModal = signal(false);
  showDeleteModal = signal(false);
  selectedApplication = signal<AdmissionApplication | null>(null);
  applicationToDelete = signal<AdmissionApplication | null>(null);

  // Stage update state
  newStage = signal<ApplicationStage | ''>('');
  stageRemarks = signal('');
  updatingStage = signal(false);
  deleting = signal(false);

  // Options
  classOptions = CLASS_NAMES;

  stageOptions: { value: ApplicationStage; label: string }[] = Object.entries(APPLICATION_STAGE_CONFIG).map(
    ([value, config]) => ({
      value: value as ApplicationStage,
      label: config.label,
    })
  );

  quickStageFilters = [
    { value: 'submitted' as ApplicationStage, label: 'Submitted' },
    { value: 'under_review' as ApplicationStage, label: 'Under Review' },
    { value: 'documents_pending' as ApplicationStage, label: 'Documents Pending' },
    { value: 'approved' as ApplicationStage, label: 'Approved' },
    { value: 'rejected' as ApplicationStage, label: 'Rejected' },
  ];

  // Computed filtered applications
  filteredApplications = computed(() => {
    let result = this.applications();
    const search = this.searchTerm().toLowerCase();
    const stage = this.selectedStage();
    const classFilter = this.selectedClass();

    if (search) {
      result = result.filter(
        (app) =>
          app.applicationNumber.toLowerCase().includes(search) ||
          app.firstName.toLowerCase().includes(search) ||
          app.lastName.toLowerCase().includes(search)
      );
    }

    if (stage) {
      result = result.filter((app) => app.currentStage === stage);
    }

    if (classFilter) {
      result = result.filter((app) => app.classApplying === classFilter);
    }

    return result;
  });

  hasActiveFilters = computed(() => {
    return !!this.searchTerm() || !!this.selectedStage() || !!this.selectedClass();
  });

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications(): void {
    this.loading.set(true);
    this.error.set(null);

    this.applicationService.getApplications().subscribe({
      next: (applications) => {
        this.applications.set(applications);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set('Failed to load applications. Please try again.');
        this.loading.set(false);
        console.error('Failed to load applications:', err);
      },
    });
  }

  // Filter handlers
  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  onStageFilterChange(stage: ApplicationStage | ''): void {
    this.selectedStage.set(stage);
  }

  onClassFilterChange(className: string): void {
    this.selectedClass.set(className);
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.selectedStage.set('');
    this.selectedClass.set('');
  }

  getStageCount(stage: ApplicationStage): number {
    return this.applications().filter((a) => a.currentStage === stage).length;
  }

  // Navigation
  createApplication(): void {
    this.router.navigate(['/admissions/applications/new']);
  }

  viewApplication(app: AdmissionApplication): void {
    this.router.navigate(['/admissions/applications', app.id]);
  }

  editApplication(app: AdmissionApplication): void {
    this.router.navigate(['/admissions/applications', app.id, 'edit']);
  }

  // Stage modal
  openStageModal(app: AdmissionApplication): void {
    this.selectedApplication.set(app);
    this.newStage.set('');
    this.stageRemarks.set('');
    this.showStageModal.set(true);
  }

  closeStageModal(): void {
    this.showStageModal.set(false);
    this.selectedApplication.set(null);
  }

  updateStage(): void {
    const app = this.selectedApplication();
    const stage = this.newStage();
    if (!app || !stage) return;

    this.updatingStage.set(true);

    this.applicationService
      .updateStage(app.id, {
        stage: stage as ApplicationStage,
        remarks: this.stageRemarks() || undefined,
      })
      .subscribe({
        next: () => {
          this.toastService.success('Application stage updated successfully');
          this.closeStageModal();
          this.loadApplications();
          this.updatingStage.set(false);
        },
        error: (err) => {
          this.toastService.error('Failed to update application stage');
          this.updatingStage.set(false);
          console.error('Failed to update stage:', err);
        },
      });
  }

  // Delete modal
  confirmDelete(app: AdmissionApplication): void {
    this.applicationToDelete.set(app);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.applicationToDelete.set(null);
  }

  deleteApplication(): void {
    const app = this.applicationToDelete();
    if (!app) return;

    this.deleting.set(true);

    this.applicationService.deleteApplication(app.id).subscribe({
      next: () => {
        this.toastService.success('Application deleted successfully');
        this.closeDeleteModal();
        this.loadApplications();
        this.deleting.set(false);
      },
      error: (err) => {
        this.toastService.error('Failed to delete application');
        this.deleting.set(false);
        console.error('Failed to delete application:', err);
      },
    });
  }

  // Helpers
  getStageConfig(stage: ApplicationStage): { label: string; variant: string; icon: string } {
    return APPLICATION_STAGE_CONFIG[stage] ?? APPLICATION_STAGE_CONFIG.draft;
  }

  getParentLabel(app: AdmissionApplication): string {
    const father = app.parents?.find((p) => p.relation === 'father');
    const mother = app.parents?.find((p) => p.relation === 'mother');
    if (father && mother) {
      return `${father.name} / ${mother.name}`;
    }
    return father?.name || mother?.name || app.parents?.[0]?.name || '';
  }

  formatDate(dateStr: string | undefined): string {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  canEditApplication(app: AdmissionApplication): boolean {
    return ['draft', 'submitted', 'documents_pending'].includes(app.currentStage);
  }

  canDeleteApplication(app: AdmissionApplication): boolean {
    return ['draft', 'rejected'].includes(app.currentStage);
  }
}
