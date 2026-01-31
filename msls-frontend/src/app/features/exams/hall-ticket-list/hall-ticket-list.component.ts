import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { ExamService } from '../exam.service';
import {
  HallTicket,
  HallTicketStatus,
  HALL_TICKET_STATUSES,
  GenerateHallTicketsRequest,
  ClassSummary,
} from '../exam.model';
import { ToastService } from '../../../shared/services/toast.service';
import { ApiService } from '../../../core/services/api.service';

interface ClassOption {
  id: string;
  name: string;
}

interface SectionOption {
  id: string;
  name: string;
  classId: string;
}

@Component({
  selector: 'app-hall-ticket-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <button class="back-btn" [routerLink]="['/exams']">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-icon">
            <i class="fa-solid fa-ticket"></i>
          </div>
          <div class="header-text">
            <h1>Hall Tickets</h1>
            <p>{{ examinationName() || 'Loading...' }}</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" (click)="downloadBatchPdf()">
            <i class="fa-solid fa-download"></i>
            <span class="btn-text">Download All</span>
          </button>
          <button class="btn btn-primary" (click)="openGenerateModal()">
            <i class="fa-solid fa-plus"></i>
            <span class="btn-text">Generate</span>
          </button>
        </div>
      </div>

      <!-- Stats Summary -->
      <div class="stats-bar">
        <div class="stat-item">
          <span class="stat-value">{{ totalCount() }}</span>
          <span class="stat-label">Total</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ generatedCount() }}</span>
          <span class="stat-label">Generated</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ downloadedCount() }}</span>
          <span class="stat-label">Downloaded</span>
        </div>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search by name, roll number..."
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
            [ngModel]="classFilter()"
            (ngModelChange)="classFilter.set($event)"
          >
            <option value="">All Classes</option>
            @for (cls of classOptions(); track cls.id) {
              <option [value]="cls.id">{{ cls.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event)"
          >
            <option value="">All Status</option>
            @for (status of hallTicketStatuses; track status.value) {
              <option [value]="status.value">{{ status.label }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading hall tickets...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadHallTickets()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Student</th>
                  <th>Roll Number</th>
                  <th class="hide-mobile">Class</th>
                  <th class="hide-tablet">Admission No</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 120px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (ticket of filteredTickets(); track ticket.id) {
                  <tr>
                    <td class="name-cell">
                      <div class="name-wrapper">
                        <div class="ticket-icon">
                          <i class="fa-solid fa-user-graduate"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ ticket.studentName }}</span>
                          <span class="class-mobile">{{ ticket.className }}</span>
                        </div>
                      </div>
                    </td>
                    <td>
                      <span class="roll-number">{{ ticket.rollNumber }}</span>
                    </td>
                    <td class="hide-mobile">
                      <span class="class-badge">{{ ticket.className }}</span>
                      @if (ticket.sectionName) {
                        <span class="section-badge">{{ ticket.sectionName }}</span>
                      }
                    </td>
                    <td class="hide-tablet">
                      {{ ticket.admissionNumber }}
                    </td>
                    <td>
                      <span class="status-badge" [class]="'status-' + ticket.status">
                        {{ getStatusLabel(ticket.status) }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        title="Download PDF"
                        (click)="downloadPdf(ticket)"
                      >
                        <i class="fa-solid fa-download"></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(ticket)"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="6" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No hall tickets found</p>
                        @if (searchTerm() || classFilter() || statusFilter()) {
                          <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                            Clear Filters
                          </button>
                        } @else {
                          <button class="btn btn-primary btn-sm" (click)="openGenerateModal()">
                            <i class="fa-solid fa-plus"></i>
                            Generate Hall Tickets
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

      <!-- Generate Modal -->
      @if (showGenerateModal()) {
        <div class="modal-overlay" (click)="closeGenerateModal()">
          <div class="modal modal--md" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-ticket"></i>
                Generate Hall Tickets
              </h3>
              <button type="button" class="modal__close" (click)="closeGenerateModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form">
                <div class="form-group">
                  <label for="genClassId">Class (Optional)</label>
                  <select
                    id="genClassId"
                    [(ngModel)]="generateData.classId"
                    name="genClassId"
                    class="form-select"
                  >
                    <option value="">All Classes</option>
                    @for (cls of classOptions(); track cls.id) {
                      <option [value]="cls.id">{{ cls.name }}</option>
                    }
                  </select>
                  <p class="form-hint">Leave empty to generate for all classes in this examination</p>
                </div>
                <div class="form-group">
                  <label for="rollNumberPrefix">Roll Number Prefix (Optional)</label>
                  <input
                    type="text"
                    id="rollNumberPrefix"
                    [(ngModel)]="generateData.rollNumberPrefix"
                    name="rollNumberPrefix"
                    class="form-input"
                    placeholder="e.g., 2026 or EX01"
                  />
                  <p class="form-hint">Default format: YEAR-CLASSCODE-SEQ (e.g., 2026-10A-001)</p>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeGenerateModal()">
                Cancel
              </button>
              <button
                type="button"
                class="btn btn--primary"
                [disabled]="generating()"
                (click)="generateHallTickets()"
              >
                @if (generating()) {
                  <div class="btn-spinner"></div>
                  Generating...
                } @else {
                  <i class="fa-solid fa-ticket"></i>
                  Generate
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
                Delete Hall Ticket
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
                  Are you sure you want to delete the hall ticket for
                  <strong>"{{ ticketToDelete()?.studentName }}"</strong>?
                </p>
                <p class="delete-warning">
                  This action cannot be undone.
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
                (click)="deleteHallTicket()"
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

      <!-- Generation Result Modal -->
      @if (showResultModal()) {
        <div class="modal-overlay" (click)="closeResultModal()">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-check-circle"></i>
                Generation Complete
              </h3>
              <button type="button" class="modal__close" (click)="closeResultModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <div class="result-summary">
                <div class="result-stat">
                  <span class="result-label">Total Students</span>
                  <span class="result-value">{{ generationResult()?.totalStudents }}</span>
                </div>
                <div class="result-stat success">
                  <span class="result-label">Generated</span>
                  <span class="result-value">{{ generationResult()?.generated }}</span>
                </div>
                <div class="result-stat warning">
                  <span class="result-label">Skipped</span>
                  <span class="result-value">{{ generationResult()?.skipped }}</span>
                </div>
                @if (generationResult()?.failed) {
                  <div class="result-stat error">
                    <span class="result-label">Failed</span>
                    <span class="result-value">{{ generationResult()?.failed }}</span>
                  </div>
                }
              </div>
              @if (generationResult()?.errors && generationResult()!.errors!.length > 0) {
                <div class="result-errors">
                  <p class="error-title">Errors:</p>
                  @for (err of generationResult()!.errors; track err) {
                    <p class="error-item">{{ err }}</p>
                  }
                </div>
              }
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--primary" (click)="closeResultModal()">
                Close
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

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s;
    }

    .back-btn:hover {
      background: #e2e8f0;
      color: #1e293b;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: linear-gradient(135deg, #8b5cf6, #7c3aed);
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
      border-color: #8b5cf6;
      box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
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
      border-color: #8b5cf6;
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
      border-top-color: #8b5cf6;
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

    .ticket-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #f3e8ff;
      color: #7c3aed;
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

    .class-mobile {
      display: none;
      font-size: 0.75rem;
      color: #64748b;
    }

    .roll-number {
      font-family: monospace;
      font-weight: 600;
      color: #7c3aed;
    }

    .class-badge, .section-badge {
      display: inline-flex;
      padding: 0.125rem 0.375rem;
      background: #e0e7ff;
      color: #4338ca;
      border-radius: 0.25rem;
      font-size: 0.6875rem;
      font-weight: 500;
      margin-right: 0.25rem;
    }

    .section-badge {
      background: #fef3c7;
      color: #d97706;
    }

    /* Status Badge */
    .status-badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .status-badge.status-generated { background: #dbeafe; color: #1d4ed8; }
    .status-badge.status-printed { background: #fef3c7; color: #d97706; }
    .status-badge.status-downloaded { background: #dcfce7; color: #16a34a; }

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
      color: #8b5cf6;
    }

    .action-btn:disabled {
      opacity: 0.3;
      cursor: not-allowed;
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
      background: linear-gradient(135deg, #8b5cf6, #7c3aed);
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: linear-gradient(135deg, #7c3aed, #6d28d9);
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
      color: #8b5cf6;
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
      background: linear-gradient(135deg, #8b5cf6, #7c3aed);
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

    .btn--primary:hover:not(:disabled) { background: linear-gradient(135deg, #7c3aed, #6d28d9); }
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
      gap: 1.25rem;
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

    .form-input,
    .form-select {
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
    .form-select:focus {
      outline: none;
      border-color: #8b5cf6;
      background: white;
      box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
    }

    .form-hint {
      margin: 0.25rem 0 0;
      font-size: 0.75rem;
      color: #64748b;
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

    /* Result Summary */
    .result-summary {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .result-stat {
      display: flex;
      flex-direction: column;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      text-align: center;
    }

    .result-stat.success { background: #dcfce7; }
    .result-stat.warning { background: #fef3c7; }
    .result-stat.error { background: #fee2e2; }

    .result-label {
      font-size: 0.75rem;
      color: #64748b;
      margin-bottom: 0.25rem;
    }

    .result-value {
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .result-stat.success .result-value { color: #16a34a; }
    .result-stat.warning .result-value { color: #d97706; }
    .result-stat.error .result-value { color: #dc2626; }

    .result-errors {
      padding: 1rem;
      background: #fee2e2;
      border-radius: 0.5rem;
    }

    .error-title {
      margin: 0 0 0.5rem;
      font-weight: 600;
      color: #dc2626;
    }

    .error-item {
      margin: 0.25rem 0;
      font-size: 0.875rem;
      color: #b91c1c;
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

      .class-mobile {
        display: block;
      }
    }
  `],
})
export class HallTicketListComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private examService = inject(ExamService);
  private apiService = inject(ApiService);
  private toastService = inject(ToastService);

  // Constants
  hallTicketStatuses = HALL_TICKET_STATUSES;

  // Route params
  examinationId = '';

  // Data signals
  hallTickets = signal<HallTicket[]>([]);
  examinationName = signal<string>('');
  classOptions = signal<ClassOption[]>([]);

  // State signals
  loading = signal(true);
  generating = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  classFilter = signal<string>('');
  statusFilter = signal<HallTicketStatus | ''>('');

  // Modal state
  showGenerateModal = signal(false);
  showDeleteModal = signal(false);
  showResultModal = signal(false);
  ticketToDelete = signal<HallTicket | null>(null);
  generationResult = signal<{ totalStudents: number; generated: number; skipped: number; failed: number; errors?: string[] } | null>(null);

  // Form data
  generateData = {
    classId: '',
    rollNumberPrefix: '',
  };

  // Computed values
  filteredTickets = computed(() => {
    let result = this.hallTickets();
    const term = this.searchTerm().toLowerCase();
    const classId = this.classFilter();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(
        ticket =>
          ticket.studentName?.toLowerCase().includes(term) ||
          ticket.rollNumber.toLowerCase().includes(term) ||
          ticket.admissionNumber?.toLowerCase().includes(term)
      );
    }

    if (classId) {
      result = result.filter(ticket => ticket.className === this.getClassName(classId));
    }

    if (status) {
      result = result.filter(ticket => ticket.status === status);
    }

    return result;
  });

  totalCount = computed(() => this.hallTickets().length);
  generatedCount = computed(() => this.hallTickets().filter(t => t.status === 'generated').length);
  downloadedCount = computed(() => this.hallTickets().filter(t => t.status === 'downloaded').length);

  ngOnInit(): void {
    this.examinationId = this.route.snapshot.paramMap.get('id') || '';
    if (this.examinationId) {
      this.loadHallTickets();
      this.loadExaminationDetails();
      this.loadClasses();
    }
  }

  loadHallTickets(): void {
    this.loading.set(true);
    this.error.set(null);

    this.examService.getHallTickets(this.examinationId, { limit: 500 }).subscribe({
      next: response => {
        this.hallTickets.set(response.data || []);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load hall tickets. Please try again.');
        this.loading.set(false);
      },
    });
  }

  loadExaminationDetails(): void {
    this.examService.getExamination(this.examinationId).subscribe({
      next: exam => this.examinationName.set(exam.name),
      error: () => console.error('Failed to load examination details'),
    });
  }

  loadClasses(): void {
    this.apiService.get<ClassOption[]>('/classes').subscribe({
      next: classes => this.classOptions.set(classes),
      error: () => console.error('Failed to load classes'),
    });
  }

  openGenerateModal(): void {
    this.generateData = { classId: '', rollNumberPrefix: '' };
    this.showGenerateModal.set(true);
  }

  closeGenerateModal(): void {
    this.showGenerateModal.set(false);
  }

  generateHallTickets(): void {
    this.generating.set(true);

    const request: GenerateHallTicketsRequest = {};
    if (this.generateData.classId) request.classId = this.generateData.classId;
    if (this.generateData.rollNumberPrefix) request.rollNumberPrefix = this.generateData.rollNumberPrefix;

    this.examService.generateHallTickets(this.examinationId, request).subscribe({
      next: result => {
        this.generationResult.set(result);
        this.closeGenerateModal();
        this.showResultModal.set(true);
        this.loadHallTickets();
        this.generating.set(false);
      },
      error: (err) => {
        const message = err?.error?.error || 'Failed to generate hall tickets';
        this.toastService.error(message);
        this.generating.set(false);
      },
    });
  }

  closeResultModal(): void {
    this.showResultModal.set(false);
    this.generationResult.set(null);
  }

  downloadPdf(ticket: HallTicket): void {
    this.examService.downloadHallTicketPdf(this.examinationId, ticket.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `hall_ticket_${ticket.rollNumber}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
        this.loadHallTickets(); // Refresh to update status
      },
      error: () => {
        this.toastService.error('Failed to download hall ticket');
      },
    });
  }

  downloadBatchPdf(): void {
    const classId = this.classFilter() || undefined;
    this.examService.downloadBatchHallTicketsPdf(this.examinationId, classId).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `hall_tickets_${this.examinationName() || 'exam'}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
      },
      error: () => {
        this.toastService.error('Failed to download hall tickets');
      },
    });
  }

  confirmDelete(ticket: HallTicket): void {
    this.ticketToDelete.set(ticket);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.ticketToDelete.set(null);
  }

  deleteHallTicket(): void {
    const ticket = this.ticketToDelete();
    if (!ticket) return;

    this.deleting.set(true);

    this.examService.deleteHallTicket(this.examinationId, ticket.id).subscribe({
      next: () => {
        this.toastService.success('Hall ticket deleted successfully');
        this.closeDeleteModal();
        this.loadHallTickets();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.error || 'Failed to delete hall ticket';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.classFilter.set('');
    this.statusFilter.set('');
  }

  getStatusLabel(status: HallTicketStatus): string {
    return HALL_TICKET_STATUSES.find(s => s.value === status)?.label || status;
  }

  getClassName(classId: string): string {
    return this.classOptions().find(c => c.id === classId)?.name || '';
  }
}
