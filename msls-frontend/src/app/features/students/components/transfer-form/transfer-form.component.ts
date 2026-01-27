/**
 * MSLS Transfer Form Component
 *
 * Modal form for processing student transfers and dropouts.
 */

import { Component, input, output, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

import { EnrollmentService } from '../../services/enrollment.service';
import { Enrollment, TransferRequest, DropoutRequest } from '../../models/enrollment.model';

export type TransferFormMode = 'transfer' | 'dropout';

@Component({
  selector: 'app-transfer-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="transfer-form">
      <!-- Modal Backdrop -->
      <div
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity z-40"
        (click)="onCancel()"
      ></div>

      <!-- Modal Panel -->
      <div class="fixed inset-0 z-50 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
            <!-- Header -->
            <div class="mb-4">
              <h3 class="text-lg font-semibold text-gray-900">
                {{ isTransfer() ? 'Process Transfer' : 'Process Dropout' }}
              </h3>
              <p class="mt-1 text-sm text-gray-500">
                {{ isTransfer()
                  ? 'This will mark the student as transferred and set their status to inactive.'
                  : 'This will mark the student as dropped out and set their status to inactive.'
                }}
              </p>
            </div>

            <!-- Warning -->
            <div class="mb-4 rounded-md bg-yellow-50 p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-yellow-800">Attention</h3>
                  <p class="mt-1 text-sm text-yellow-700">
                    This action cannot be undone. The student's enrollment will be permanently marked as
                    {{ isTransfer() ? 'transferred' : 'dropped out' }}.
                  </p>
                </div>
              </div>
            </div>

            <!-- Form -->
            <form [formGroup]="form" (ngSubmit)="onSubmit()">
              <!-- Date Field -->
              <div class="mb-4">
                <label for="date" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ isTransfer() ? 'Transfer Date' : 'Dropout Date' }} <span class="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  id="date"
                  formControlName="date"
                  [max]="today"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  [class.border-red-300]="isInvalid('date')"
                />
                @if (isInvalid('date')) {
                  <p class="mt-1 text-sm text-red-600">
                    {{ isTransfer() ? 'Transfer date is required' : 'Dropout date is required' }}
                  </p>
                }
              </div>

              <!-- Reason Field -->
              <div class="mb-4">
                <label for="reason" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ isTransfer() ? 'Transfer Reason' : 'Dropout Reason' }} <span class="text-red-500">*</span>
                </label>
                <textarea
                  id="reason"
                  formControlName="reason"
                  rows="3"
                  [placeholder]="isTransfer() ? 'e.g., Family relocation, Admission to another school...' : 'e.g., Financial difficulties, Personal reasons...'"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  [class.border-red-300]="isInvalid('reason')"
                ></textarea>
                @if (isInvalid('reason')) {
                  <p class="mt-1 text-sm text-red-600">
                    {{ isTransfer() ? 'Transfer reason is required' : 'Dropout reason is required' }}
                  </p>
                }
              </div>

              <!-- Error Message -->
              @if (service.error()) {
                <div class="mb-4 rounded-md bg-red-50 p-4">
                  <p class="text-sm text-red-800">{{ service.error() }}</p>
                </div>
              }

              <!-- Form Actions -->
              <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                <button
                  type="submit"
                  [disabled]="form.invalid || service.loading()"
                  class="inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-white shadow-sm sm:col-start-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  [ngClass]="isTransfer() ? 'bg-orange-600 hover:bg-orange-500' : 'bg-red-600 hover:bg-red-500'"
                >
                  @if (service.loading()) {
                    <span class="flex items-center">
                      <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                      </svg>
                      Processing...
                    </span>
                  } @else {
                    {{ isTransfer() ? 'Confirm Transfer' : 'Confirm Dropout' }}
                  }
                </button>
                <button
                  type="button"
                  class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                  (click)="onCancel()"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }
  `]
})
export class TransferFormComponent implements OnInit {
  protected service = inject(EnrollmentService);
  private fb = inject(FormBuilder);

  /** Student ID */
  studentId = input.required<string>();

  /** The enrollment being processed */
  enrollment = input.required<Enrollment>();

  /** Form mode: 'transfer' or 'dropout' */
  mode = input<TransferFormMode>('transfer');

  /** Event emitted when form is saved successfully */
  saved = output<Enrollment>();

  /** Event emitted when form is cancelled */
  cancelled = output<void>();

  form!: FormGroup;
  today = new Date().toISOString().split('T')[0];

  ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.fb.group({
      date: [this.today, Validators.required],
      reason: ['', Validators.required],
    });
  }

  isTransfer(): boolean {
    return this.mode() === 'transfer';
  }

  isInvalid(field: string): boolean {
    const control = this.form.get(field);
    return !!(control && control.invalid && (control.dirty || control.touched));
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const formValue = this.form.value;
    const enrollment = this.enrollment();

    if (this.isTransfer()) {
      const request: TransferRequest = {
        transferDate: formValue.date,
        transferReason: formValue.reason,
      };

      this.service.processTransfer(this.studentId(), enrollment.id, request).subscribe({
        next: (updated) => this.saved.emit(updated),
        error: () => {}, // Error handled in service
      });
    } else {
      const request: DropoutRequest = {
        dropoutDate: formValue.date,
        dropoutReason: formValue.reason,
      };

      this.service.processDropout(this.studentId(), enrollment.id, request).subscribe({
        next: (updated) => this.saved.emit(updated),
        error: () => {}, // Error handled in service
      });
    }
  }

  onCancel(): void {
    this.cancelled.emit();
  }
}
