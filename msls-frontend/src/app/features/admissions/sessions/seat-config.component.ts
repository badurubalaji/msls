/**
 * MSLS Seat Configuration Component
 *
 * Component for managing seat configurations per class for an admission session.
 */

import { Component, Input, Output, EventEmitter, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import {
  MslsButtonComponent,
  MslsInputComponent,
  MslsFormFieldComponent,
  MslsSelectComponent,
  MslsModalComponent,
  SelectOption,
} from '../../../shared/components';
import { ToastService } from '../../../shared/services';
import {
  AdmissionSession,
  AdmissionSeat,
  SeatConfigRequest,
  CLASS_NAMES,
} from './admission-session.model';
import { AdmissionSessionService } from './admission-session.service';

@Component({
  selector: 'msls-seat-config',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsButtonComponent,
    MslsInputComponent,
    MslsFormFieldComponent,
    MslsSelectComponent,
    MslsModalComponent,
  ],
  template: `
    <div class="seat-config">
      <!-- Header -->
      <div class="seat-config__header">
        <div class="seat-config__info">
          <h3 class="seat-config__title">Seat Configuration</h3>
          <p class="seat-config__subtitle">Configure available seats per class for this admission session</p>
        </div>
        <button class="btn btn-primary btn-sm" (click)="openAddModal()">
          <i class="fa-solid fa-plus"></i>
          Add Class
        </button>
      </div>

      <!-- Summary -->
      <div class="seat-summary">
        <div class="summary-item">
          <span class="summary-label">Total Seats</span>
          <span class="summary-value">{{ totalSeats() }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">Filled</span>
          <span class="summary-value summary-value--success">{{ filledSeats() }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">Available</span>
          <span class="summary-value summary-value--primary">{{ availableSeats() }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">Classes</span>
          <span class="summary-value">{{ seats().length }}</span>
        </div>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <p>Loading seat configurations...</p>
        </div>
      } @else if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <p>{{ error() }}</p>
          <button class="btn btn-secondary btn-sm" (click)="loadSeats()">
            Retry
          </button>
        </div>
      } @else {
        <!-- Seat Table -->
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Class</th>
                <th style="width: 100px; text-align: center;">Total Seats</th>
                <th style="width: 100px; text-align: center;">Filled</th>
                <th style="width: 100px; text-align: center;">Available</th>
                <th style="width: 100px; text-align: center;">Waitlist</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (seat of seats(); track seat.id) {
                <tr>
                  <td class="class-cell">{{ seat.className }}</td>
                  <td class="center-cell">{{ seat.totalSeats }}</td>
                  <td class="center-cell">
                    <span class="filled-count">{{ seat.filledSeats }}</span>
                  </td>
                  <td class="center-cell">
                    <span
                      class="available-count"
                      [class.available-count--low]="(seat.totalSeats - seat.filledSeats) < 5"
                    >
                      {{ seat.totalSeats - seat.filledSeats }}
                    </span>
                  </td>
                  <td class="center-cell">{{ seat.waitlistLimit }}</td>
                  <td class="actions-cell">
                    <button
                      class="action-btn"
                      (click)="editSeat(seat)"
                      title="Edit seat configuration"
                    >
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button
                      class="action-btn action-btn--danger"
                      (click)="confirmDeleteSeat(seat)"
                      title="Remove class"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="6" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-rectangle-list"></i>
                      <p>No seat configurations added yet</p>
                      <button class="btn btn-primary btn-sm" (click)="openAddModal()">
                        <i class="fa-solid fa-plus"></i>
                        Add your first class
                      </button>
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      }

      <!-- Footer Actions -->
      <div class="seat-config__footer">
        <msls-button variant="secondary" (clicked)="onClose()">
          Close
        </msls-button>
      </div>

      <!-- Add/Edit Seat Modal -->
      <msls-modal
        [isOpen]="showSeatModal()"
        [title]="editingSeat() ? 'Edit Seat Configuration' : 'Add Class'"
        size="md"
        (closed)="closeSeatModal()"
      >
        <form [formGroup]="seatForm" (ngSubmit)="saveSeat()" class="seat-form">
          <!-- Class Name -->
          <msls-form-field
            label="Class"
            [required]="true"
            [error]="getSeatFieldError('className') || ''"
          >
            @if (editingSeat()) {
              <msls-input
                type="text"
                formControlName="className"
                [disabled]="true"
              />
            } @else {
              <msls-select
                formControlName="className"
                [options]="availableClassOptions()"
                placeholder="Select class"
              />
            }
          </msls-form-field>

          <div class="form-grid">
            <!-- Total Seats -->
            <msls-form-field
              label="Total Seats"
              [required]="true"
              [error]="getSeatFieldError('totalSeats') || ''"
            >
              <msls-input
                type="number"
                formControlName="totalSeats"
                placeholder="e.g., 40"
              />
            </msls-form-field>

            <!-- Waitlist Limit -->
            <msls-form-field
              label="Waitlist Limit"
              [error]="getSeatFieldError('waitlistLimit') || ''"
            >
              <msls-input
                type="number"
                formControlName="waitlistLimit"
                placeholder="e.g., 10"
              />
            </msls-form-field>
          </div>

          <div class="form-actions">
            <msls-button
              type="button"
              variant="secondary"
              (clicked)="closeSeatModal()"
            >
              Cancel
            </msls-button>
            <msls-button
              type="submit"
              variant="primary"
              [loading]="saving()"
              [disabled]="!seatForm.valid"
            >
              {{ editingSeat() ? 'Update' : 'Add Class' }}
            </msls-button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Remove Class"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to remove
            <strong>"{{ seatToDelete()?.className }}"</strong>
            from this admission session?
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteSeat()"
            >
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Removing...
              } @else {
                Remove
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .seat-config {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .seat-config__header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding-bottom: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .seat-config__title {
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.25rem 0;
    }

    .seat-config__subtitle {
      font-size: 0.8125rem;
      color: #64748b;
      margin: 0;
    }

    /* Summary */
    .seat-summary {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
    }

    .summary-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
    }

    .summary-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .summary-value {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin-top: 0.25rem;
    }

    .summary-value--primary {
      color: #4f46e5;
    }

    .summary-value--success {
      color: #16a34a;
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

    .btn-danger {
      background: #dc2626;
      color: white;
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
      padding: 3rem;
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
      padding: 3rem;
      gap: 1rem;
    }

    .error-container i {
      font-size: 2rem;
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
      padding: 0.875rem 1rem;
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

    .class-cell {
      font-weight: 500;
      color: #0f172a;
    }

    .center-cell {
      text-align: center;
    }

    .filled-count {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      padding: 0.25rem 0.5rem;
      background: #dcfce7;
      color: #166534;
      font-size: 0.8125rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .available-count {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      padding: 0.25rem 0.5rem;
      background: #dbeafe;
      color: #1e40af;
      font-size: 0.8125rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .available-count--low {
      background: #fef3c7;
      color: #92400e;
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

    .action-btn--danger:hover {
      background: #fef2f2;
      border-color: #fecaca;
      color: #dc2626;
    }

    /* Empty State */
    .empty-cell {
      padding: 2.5rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #94a3b8;
    }

    .empty-state i {
      font-size: 2rem;
    }

    .empty-state p {
      margin: 0;
      font-size: 0.875rem;
    }

    /* Footer */
    .seat-config__footer {
      display: flex;
      justify-content: flex-end;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    /* Form */
    .seat-form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .form-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      margin-top: 0.5rem;
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
      margin: 0 0 1.5rem 0;
    }

    .delete-confirmation strong {
      color: #0f172a;
    }

    .delete-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
    }

    @media (max-width: 640px) {
      .seat-summary {
        grid-template-columns: repeat(2, 1fr);
      }

      .form-grid {
        grid-template-columns: 1fr;
      }

      .table-container {
        overflow-x: auto;
      }

      .data-table {
        min-width: 500px;
      }
    }
  `],
})
export class SeatConfigComponent implements OnInit {
  private fb = inject(FormBuilder);
  private sessionService = inject(AdmissionSessionService);
  private toastService = inject(ToastService);

  @Input() session!: AdmissionSession;
  @Output() close = new EventEmitter<void>();
  @Output() seatsChanged = new EventEmitter<void>();

  // State signals
  seats = signal<AdmissionSeat[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);

  // Modal state
  showSeatModal = signal(false);
  showDeleteModal = signal(false);
  editingSeat = signal<AdmissionSeat | null>(null);
  seatToDelete = signal<AdmissionSeat | null>(null);

  seatForm!: FormGroup;

  // Computed values
  totalSeats = computed(() =>
    this.seats().reduce((sum, s) => sum + s.totalSeats, 0)
  );

  filledSeats = computed(() =>
    this.seats().reduce((sum, s) => sum + s.filledSeats, 0)
  );

  availableSeats = computed(() =>
    this.totalSeats() - this.filledSeats()
  );

  /** Available class options (excluding already added) */
  availableClassOptions = computed<SelectOption[]>(() => {
    const usedClasses = new Set(this.seats().map(s => s.className));
    return CLASS_NAMES
      .filter(c => !usedClasses.has(c))
      .map(c => ({ value: c, label: c }));
  });

  ngOnInit(): void {
    this.initSeatForm();
    this.loadSeats();
  }

  private initSeatForm(): void {
    this.seatForm = this.fb.group({
      className: ['', [Validators.required]],
      totalSeats: [40, [Validators.required, Validators.min(1), Validators.max(1000)]],
      waitlistLimit: [10, [Validators.min(0), Validators.max(100)]],
    });
  }

  loadSeats(): void {
    this.loading.set(true);
    this.error.set(null);

    this.sessionService.getSeats(this.session.id).subscribe({
      next: (seats) => {
        this.seats.set(seats);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set('Failed to load seat configurations');
        this.loading.set(false);
        console.error('Failed to load seats:', err);
      },
    });
  }

  openAddModal(): void {
    this.editingSeat.set(null);
    this.seatForm.reset({ totalSeats: 40, waitlistLimit: 10 });
    this.seatForm.get('className')?.enable();
    this.showSeatModal.set(true);
  }

  editSeat(seat: AdmissionSeat): void {
    this.editingSeat.set(seat);
    this.seatForm.patchValue({
      className: seat.className,
      totalSeats: seat.totalSeats,
      waitlistLimit: seat.waitlistLimit,
    });
    this.seatForm.get('className')?.disable();
    this.showSeatModal.set(true);
  }

  closeSeatModal(): void {
    this.showSeatModal.set(false);
    this.editingSeat.set(null);
    this.seatForm.get('className')?.enable();
  }

  saveSeat(): void {
    if (!this.seatForm.valid) return;

    this.saving.set(true);
    const value = this.seatForm.getRawValue();
    const request: SeatConfigRequest = {
      className: value.className,
      totalSeats: Number(value.totalSeats),
      waitlistLimit: Number(value.waitlistLimit) || 10,
    };

    const editing = this.editingSeat();
    const operation = editing
      ? this.sessionService.updateSeat(this.session.id, editing.id, request)
      : this.sessionService.addSeat(this.session.id, request);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Seat configuration updated' : 'Class added successfully'
        );
        this.closeSeatModal();
        this.loadSeats();
        this.seatsChanged.emit();
        this.saving.set(false);
      },
      error: (err) => {
        this.toastService.error(
          editing ? 'Failed to update seat configuration' : 'Failed to add class'
        );
        this.saving.set(false);
        console.error('Failed to save seat:', err);
      },
    });
  }

  confirmDeleteSeat(seat: AdmissionSeat): void {
    this.seatToDelete.set(seat);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.seatToDelete.set(null);
  }

  deleteSeat(): void {
    const seat = this.seatToDelete();
    if (!seat) return;

    this.deleting.set(true);

    this.sessionService.deleteSeat(this.session.id, seat.id).subscribe({
      next: () => {
        this.toastService.success('Class removed successfully');
        this.closeDeleteModal();
        this.loadSeats();
        this.seatsChanged.emit();
        this.deleting.set(false);
      },
      error: (err) => {
        this.toastService.error('Failed to remove class');
        this.deleting.set(false);
        console.error('Failed to delete seat:', err);
      },
    });
  }

  getSeatFieldError(field: string): string | null {
    const control = this.seatForm.get(field);
    if (!control || !control.touched || !control.errors) return null;

    if (control.errors['required']) return 'This field is required';
    if (control.errors['min']) return 'Value must be at least 1';
    if (control.errors['max']) return 'Value is too large';

    return null;
  }

  onClose(): void {
    this.close.emit();
  }
}
