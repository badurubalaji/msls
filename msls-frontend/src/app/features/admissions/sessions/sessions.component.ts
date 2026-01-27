/**
 * MSLS Admission Sessions Component
 *
 * Main component for managing admission sessions - displays a list of sessions with CRUD operations.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsModalComponent } from '../../../shared/components';
import { ToastService } from '../../../shared/services';
import {
  AdmissionSession,
  CreateSessionRequest,
  SessionStatus,
  getStatusConfig,
} from './admission-session.model';
import { AdmissionSessionService } from './admission-session.service';
import { SessionFormComponent } from './session-form.component';
import { SeatConfigComponent } from './seat-config.component';

@Component({
  selector: 'msls-sessions',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsModalComponent,
    SessionFormComponent,
    SeatConfigComponent,
    DatePipe,
  ],
  template: `
    <div class="sessions-page">
      <div class="sessions-card">
        <!-- Header -->
        <div class="sessions-header">
          <div class="sessions-header__left">
            <h1 class="sessions-header__title">Admission Sessions</h1>
            <p class="sessions-header__subtitle">
              Configure and manage admission cycles for your school
            </p>
          </div>
          <div class="sessions-header__right">
            <div class="search-input">
              <i class="fa-solid fa-magnifying-glass search-icon"></i>
              <input
                type="text"
                placeholder="Search sessions..."
                [ngModel]="searchTerm()"
                (ngModelChange)="onSearchChange($event)"
                class="search-field"
              />
            </div>
            <button class="btn btn-primary" (click)="openCreateModal()">
              <i class="fa-solid fa-plus"></i>
              New Session
            </button>
          </div>
        </div>

        <!-- Loading State -->
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading admission sessions...</p>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <p>{{ error() }}</p>
            <button class="btn btn-secondary" (click)="loadSessions()">
              Retry
            </button>
          </div>
        } @else {
          <!-- Table -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Session Name</th>
                  <th style="width: 100px;">Academic Year</th>
                  <th style="width: 120px;">Start Date</th>
                  <th style="width: 120px;">End Date</th>
                  <th style="width: 100px; text-align: center;">Status</th>
                  <th style="width: 140px; text-align: center;">Seats</th>
                  <th style="width: 200px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (session of filteredSessions(); track session.id) {
                  <tr>
                    <td class="name-cell">
                      <div class="session-name">
                        {{ session.name }}
                        @if (session.applicationFee > 0) {
                          <span class="fee-badge">
                            <i class="fa-solid fa-indian-rupee-sign"></i>
                            {{ session.applicationFee }}
                          </span>
                        }
                      </div>
                      @if (session.totalApplications > 0) {
                        <span class="applications-count">
                          {{ session.totalApplications }} applications
                        </span>
                      }
                    </td>
                    <td class="year-cell">{{ session.academicYearName }}</td>
                    <td class="date-cell">{{ session.startDate | date:'dd MMM yyyy' }}</td>
                    <td class="date-cell">{{ session.endDate | date:'dd MMM yyyy' }}</td>
                    <td class="status-cell">
                      <span
                        class="badge"
                        [class]="getStatusConfig(session.status).class"
                      >
                        <i class="fa-solid" [class]="getStatusConfig(session.status).icon"></i>
                        {{ getStatusConfig(session.status).label }}
                      </span>
                    </td>
                    <td class="seats-cell">
                      <div class="seats-info">
                        <span class="seats-filled">{{ session.filledSeats }}</span>
                        <span class="seats-separator">/</span>
                        <span class="seats-total">{{ session.totalSeats }}</span>
                      </div>
                      @if (session.totalSeats > 0) {
                        <div class="seats-bar">
                          <div
                            class="seats-bar__fill"
                            [style.width.%]="(session.filledSeats / session.totalSeats) * 100"
                          ></div>
                        </div>
                      }
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        (click)="openSeatConfig(session)"
                        title="Configure seats"
                      >
                        <i class="fa-solid fa-chair"></i>
                      </button>
                      @if (session.status === 'upcoming') {
                        <button
                          class="action-btn action-btn--success"
                          (click)="openSession(session)"
                          title="Open admissions"
                        >
                          <i class="fa-solid fa-door-open"></i>
                        </button>
                      }
                      @if (session.status === 'open') {
                        <button
                          class="action-btn action-btn--warning"
                          (click)="closeSession(session)"
                          title="Close admissions"
                        >
                          <i class="fa-solid fa-door-closed"></i>
                        </button>
                      }
                      <button
                        class="action-btn"
                        (click)="editSession(session)"
                        title="Edit session"
                        [disabled]="session.status === 'closed'"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        (click)="confirmDelete(session)"
                        title="Delete session"
                        [disabled]="session.status !== 'upcoming' || session.totalApplications > 0"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="7" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-calendar-plus"></i>
                        <p>No admission sessions found</p>
                        @if (!searchTerm()) {
                          <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                            <i class="fa-solid fa-plus"></i>
                            Create your first session
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

      <!-- Create/Edit Session Modal -->
      <msls-modal
        [isOpen]="showSessionModal()"
        [title]="editingSession() ? 'Edit Session' : 'Create Admission Session'"
        size="lg"
        (closed)="closeSessionModal()"
      >
        <msls-session-form
          [session]="editingSession()"
          [loading]="saving()"
          (save)="saveSession($event)"
          (cancel)="closeSessionModal()"
        />
      </msls-modal>

      <!-- Seat Configuration Modal -->
      <msls-modal
        [isOpen]="showSeatModal()"
        [title]="selectedSession()?.name + ' - Seat Configuration'"
        size="xl"
        (closed)="closeSeatModal()"
      >
        @if (selectedSession()) {
          <msls-seat-config
            [session]="selectedSession()!"
            (close)="closeSeatModal()"
            (seatsChanged)="onSeatsChanged()"
          />
        }
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Session"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete
            <strong>"{{ sessionToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteSession()"
            >
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

      <!-- Status Change Confirmation Modal -->
      <msls-modal
        [isOpen]="showStatusModal()"
        [title]="statusAction() === 'open' ? 'Open Admissions' : 'Close Admissions'"
        size="sm"
        (closed)="closeStatusModal()"
      >
        <div class="status-confirmation">
          <div
            class="status-icon"
            [class.status-icon--success]="statusAction() === 'open'"
            [class.status-icon--warning]="statusAction() === 'close'"
          >
            <i
              class="fa-solid"
              [class.fa-door-open]="statusAction() === 'open'"
              [class.fa-door-closed]="statusAction() === 'close'"
            ></i>
          </div>
          @if (statusAction() === 'open') {
            <p>
              Are you sure you want to <strong>open</strong> admissions for
              <strong>"{{ sessionToChangeStatus()?.name }}"</strong>?
            </p>
            <p class="status-note">
              This will allow students to submit applications for this session.
            </p>
          } @else {
            <p>
              Are you sure you want to <strong>close</strong> admissions for
              <strong>"{{ sessionToChangeStatus()?.name }}"</strong>?
            </p>
            <p class="status-note">
              No new applications will be accepted after closing.
            </p>
          }
          <div class="status-actions">
            <button class="btn btn-secondary" (click)="closeStatusModal()">
              Cancel
            </button>
            <button
              class="btn"
              [class.btn-success]="statusAction() === 'open'"
              [class.btn-warning]="statusAction() === 'close'"
              [disabled]="changingStatus()"
              (click)="confirmStatusChange()"
            >
              @if (changingStatus()) {
                <div class="btn-spinner"></div>
                Processing...
              } @else {
                {{ statusAction() === 'open' ? 'Open Admissions' : 'Close Admissions' }}
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .sessions-page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .sessions-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    /* Header */
    .sessions-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .sessions-header__title {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0 0 0.375rem 0;
    }

    .sessions-header__subtitle {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
    }

    .sessions-header__right {
      display: flex;
      gap: 0.75rem;
      align-items: center;
    }

    /* Search Input */
    .search-input {
      position: relative;
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
      padding: 0.625rem 0.875rem 0.625rem 2.5rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #0f172a;
      width: 220px;
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
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: white;
      color: #334155;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    .btn-success {
      background: #16a34a;
      color: white;
    }

    .btn-success:hover:not(:disabled) {
      background: #15803d;
    }

    .btn-warning {
      background: #ea580c;
      color: white;
    }

    .btn-warning:hover:not(:disabled) {
      background: #c2410c;
    }

    .btn-danger {
      background: #dc2626;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
    }

    .btn:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 1rem;
      height: 1rem;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-top-color: white;
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
      to { transform: rotate(360deg); }
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

    .data-table tbody tr {
      transition: background 0.15s;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .name-cell {
      font-weight: 500;
      color: #0f172a;
    }

    .session-name {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .fee-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.125rem;
      padding: 0.125rem 0.5rem;
      background: #f0fdf4;
      color: #166534;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .fee-badge i {
      font-size: 0.625rem;
    }

    .applications-count {
      display: block;
      font-size: 0.75rem;
      color: #64748b;
      font-weight: 400;
      margin-top: 0.25rem;
    }

    .year-cell {
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      font-size: 0.8125rem;
      color: #64748b;
      background: #f1f5f9;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      display: inline-block;
    }

    .date-cell {
      color: #64748b;
      white-space: nowrap;
    }

    .status-cell {
      text-align: center;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.25rem 0.625rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .badge i {
      font-size: 0.625rem;
    }

    .badge-blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .badge-green {
      background: #dcfce7;
      color: #166534;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #475569;
    }

    /* Seats */
    .seats-cell {
      text-align: center;
    }

    .seats-info {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.25rem;
      font-size: 0.875rem;
    }

    .seats-filled {
      font-weight: 600;
      color: #16a34a;
    }

    .seats-separator {
      color: #94a3b8;
    }

    .seats-total {
      color: #64748b;
    }

    .seats-bar {
      width: 80px;
      height: 4px;
      background: #e2e8f0;
      border-radius: 2px;
      margin: 0.375rem auto 0;
      overflow: hidden;
    }

    .seats-bar__fill {
      height: 100%;
      background: #4f46e5;
      border-radius: 2px;
      transition: width 0.3s;
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

    .action-btn:hover:not(:disabled) {
      background: #f8fafc;
      border-color: #cbd5e1;
      color: #0f172a;
    }

    .action-btn:disabled {
      opacity: 0.4;
      cursor: not-allowed;
    }

    .action-btn--success {
      border-color: #bbf7d0;
      color: #16a34a;
    }

    .action-btn--success:hover:not(:disabled) {
      background: #f0fdf4;
      border-color: #86efac;
      color: #15803d;
    }

    .action-btn--warning {
      border-color: #fed7aa;
      color: #ea580c;
    }

    .action-btn--warning:hover:not(:disabled) {
      background: #fff7ed;
      border-color: #fdba74;
      color: #c2410c;
    }

    .action-btn--danger:hover:not(:disabled) {
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

    /* Status Confirmation */
    .status-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .status-icon {
      width: 3.5rem;
      height: 3.5rem;
      margin: 0 auto 1rem;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .status-icon--success {
      background: #dcfce7;
    }

    .status-icon--success i {
      font-size: 1.5rem;
      color: #16a34a;
    }

    .status-icon--warning {
      background: #fff7ed;
    }

    .status-icon--warning i {
      font-size: 1.5rem;
      color: #ea580c;
    }

    .status-confirmation p {
      color: #475569;
      margin: 0 0 0.5rem 0;
    }

    .status-confirmation strong {
      color: #0f172a;
    }

    .status-note {
      font-size: 0.8125rem;
      color: #64748b !important;
      padding: 0.75rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      margin-top: 1rem !important;
    }

    .status-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }

    /* Responsive */
    @media (max-width: 768px) {
      .sessions-header {
        flex-direction: column;
        align-items: stretch;
      }

      .sessions-header__right {
        flex-direction: column;
      }

      .search-field {
        width: 100%;
      }

      .btn-primary {
        width: 100%;
        justify-content: center;
      }

      .table-container {
        overflow-x: auto;
      }

      .data-table {
        min-width: 900px;
      }
    }
  `],
})
export class SessionsComponent implements OnInit {
  private sessionService = inject(AdmissionSessionService);
  private toastService = inject(ToastService);

  // State signals
  sessions = signal<AdmissionSession[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  changingStatus = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');

  // Modal state
  showSessionModal = signal(false);
  showSeatModal = signal(false);
  showDeleteModal = signal(false);
  showStatusModal = signal(false);
  editingSession = signal<AdmissionSession | null>(null);
  selectedSession = signal<AdmissionSession | null>(null);
  sessionToDelete = signal<AdmissionSession | null>(null);
  sessionToChangeStatus = signal<AdmissionSession | null>(null);
  statusAction = signal<'open' | 'close'>('open');

  // Computed filtered sessions
  filteredSessions = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.sessions();

    return this.sessions().filter(
      session =>
        session.name.toLowerCase().includes(term) ||
        session.academicYearName?.toLowerCase().includes(term)
    );
  });

  // Helper function for template
  getStatusConfig = getStatusConfig;

  ngOnInit(): void {
    this.loadSessions();
  }

  loadSessions(): void {
    this.loading.set(true);
    this.error.set(null);

    this.sessionService.getSessions().subscribe({
      next: sessions => {
        this.sessions.set(sessions);
        this.loading.set(false);
      },
      error: err => {
        this.error.set('Failed to load admission sessions. Please try again.');
        this.loading.set(false);
        console.error('Failed to load sessions:', err);
      },
    });
  }

  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  // Create/Edit Modal
  openCreateModal(): void {
    this.editingSession.set(null);
    this.showSessionModal.set(true);
  }

  editSession(session: AdmissionSession): void {
    this.editingSession.set(session);
    this.showSessionModal.set(true);
  }

  closeSessionModal(): void {
    this.showSessionModal.set(false);
    this.editingSession.set(null);
  }

  saveSession(data: CreateSessionRequest): void {
    this.saving.set(true);

    const editing = this.editingSession();
    const operation = editing
      ? this.sessionService.updateSession(editing.id, data)
      : this.sessionService.createSession(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Session updated successfully' : 'Session created successfully'
        );
        this.closeSessionModal();
        this.loadSessions();
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update session' : 'Failed to create session'
        );
        this.saving.set(false);
        console.error('Failed to save session:', err);
      },
    });
  }

  // Seat Configuration Modal
  openSeatConfig(session: AdmissionSession): void {
    this.selectedSession.set(session);
    this.showSeatModal.set(true);
  }

  closeSeatModal(): void {
    this.showSeatModal.set(false);
    this.selectedSession.set(null);
  }

  onSeatsChanged(): void {
    this.loadSessions();
  }

  // Status Change
  openSession(session: AdmissionSession): void {
    this.sessionToChangeStatus.set(session);
    this.statusAction.set('open');
    this.showStatusModal.set(true);
  }

  closeSession(session: AdmissionSession): void {
    this.sessionToChangeStatus.set(session);
    this.statusAction.set('close');
    this.showStatusModal.set(true);
  }

  closeStatusModal(): void {
    this.showStatusModal.set(false);
    this.sessionToChangeStatus.set(null);
  }

  confirmStatusChange(): void {
    const session = this.sessionToChangeStatus();
    if (!session) return;

    this.changingStatus.set(true);
    const newStatus: SessionStatus = this.statusAction() === 'open' ? 'open' : 'closed';

    this.sessionService.changeStatus(session.id, { status: newStatus }).subscribe({
      next: () => {
        this.toastService.success(
          this.statusAction() === 'open'
            ? 'Admissions opened successfully'
            : 'Admissions closed successfully'
        );
        this.closeStatusModal();
        this.loadSessions();
        this.changingStatus.set(false);
      },
      error: err => {
        this.toastService.error('Failed to update session status');
        this.changingStatus.set(false);
        console.error('Failed to change status:', err);
      },
    });
  }

  // Delete Modal
  confirmDelete(session: AdmissionSession): void {
    this.sessionToDelete.set(session);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.sessionToDelete.set(null);
  }

  deleteSession(): void {
    const session = this.sessionToDelete();
    if (!session) return;

    this.deleting.set(true);

    this.sessionService.deleteSession(session.id).subscribe({
      next: () => {
        this.toastService.success('Session deleted successfully');
        this.closeDeleteModal();
        this.loadSessions();
        this.deleting.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete session');
        this.deleting.set(false);
        console.error('Failed to delete session:', err);
      },
    });
  }
}
