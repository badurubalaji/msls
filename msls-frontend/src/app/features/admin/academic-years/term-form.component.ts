/**
 * MSLS Term Form Component
 *
 * Form component for creating and editing academic terms/semesters.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import { AcademicYear, AcademicTerm } from './academic-year.model';

@Component({
  selector: 'msls-term-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="term-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="term">
          <i class="fa-solid" [class.fa-calendar-days]="!term" [class.fa-pen]="term"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ term ? 'Edit Term' : 'Add New Term' }}</h3>
          <p *ngIf="academicYear">{{ academicYear.name }}</p>
        </div>
      </div>

      <!-- Term Name -->
      <div class="form-field">
        <label class="field-label">
          Term Name <span class="required">*</span>
        </label>
        <div class="input-wrapper">
          <i class="fa-solid fa-tag input-icon"></i>
          <input
            type="text"
            formControlName="name"
            placeholder="e.g., First Semester, Term 1, Quarter 1"
            class="form-input"
            [class.has-error]="hasError('name')"
          />
        </div>
        <span *ngIf="hasError('name')" class="field-error">{{ getFieldError('name') }}</span>
      </div>

      <!-- Sequence -->
      <div class="form-field">
        <label class="field-label">
          Sequence <span class="required">*</span>
        </label>
        <div class="input-wrapper">
          <i class="fa-solid fa-sort-numeric-up input-icon"></i>
          <input
            type="number"
            formControlName="sequence"
            min="1"
            placeholder="Enter sequence number (1, 2, 3...)"
            class="form-input"
            [class.has-error]="hasError('sequence')"
          />
        </div>
        <span class="field-hint">Order in which the term appears in the academic year</span>
        <span *ngIf="hasError('sequence')" class="field-error">{{ getFieldError('sequence') }}</span>
      </div>

      <!-- Date Row -->
      <div class="form-row">
        <div class="form-field">
          <label class="field-label">
            Start Date <span class="required">*</span>
          </label>
          <div class="input-wrapper">
            <i class="fa-solid fa-calendar input-icon"></i>
            <input
              type="date"
              formControlName="startDate"
              class="form-input"
              [class.has-error]="hasError('startDate')"
              [min]="getMinDate()"
              [max]="getMaxDate()"
            />
          </div>
          <span *ngIf="hasError('startDate')" class="field-error">{{ getFieldError('startDate') }}</span>
        </div>

        <div class="form-field">
          <label class="field-label">
            End Date <span class="required">*</span>
          </label>
          <div class="input-wrapper">
            <i class="fa-solid fa-calendar-check input-icon"></i>
            <input
              type="date"
              formControlName="endDate"
              class="form-input"
              [class.has-error]="hasError('endDate')"
              [min]="getMinDate()"
              [max]="getMaxDate()"
            />
          </div>
          <span *ngIf="hasError('endDate')" class="field-error">{{ getFieldError('endDate') }}</span>
        </div>
      </div>

      <!-- Date Range Errors -->
      <div *ngIf="form.errors?.['dateRange']" class="form-error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        End date must be after start date
      </div>

      <div *ngIf="form.errors?.['outsideRange']" class="form-error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        Term dates must be within the academic year period
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button type="submit" class="btn btn--primary" [disabled]="!form.valid || loading">
          <i *ngIf="loading" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!loading && !term" class="fa-solid fa-plus"></i>
          <i *ngIf="!loading && term" class="fa-solid fa-save"></i>
          {{ term ? 'Update Term' : 'Add Term' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .term-form {
      display: flex;
      flex-direction: column;
      gap: 1.25rem;
    }

    .form-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .form-header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 1.25rem;
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    }

    .form-header-icon.edit-mode {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
    }

    .form-header-text h3 {
      margin: 0;
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
    }

    .form-header-text p {
      margin: 0.25rem 0 0 0;
      font-size: 0.875rem;
      color: #64748b;
    }

    .form-field {
      display: flex;
      flex-direction: column;
    }

    .field-label {
      font-size: 0.8125rem;
      font-weight: 500;
      color: #374151;
      margin-bottom: 0.5rem;
    }

    .required {
      color: #ef4444;
    }

    .input-wrapper {
      position: relative;
    }

    .input-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      font-size: 0.875rem;
      color: #94a3b8;
    }

    .form-input {
      width: 100%;
      height: 2.5rem;
      padding: 0 0.875rem 0 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      color: #0f172a;
      background: white;
      transition: all 0.15s;
    }

    .form-input:focus {
      outline: none;
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }

    .form-input.has-error {
      border-color: #ef4444;
    }

    .field-error {
      margin-top: 0.375rem;
      font-size: 0.75rem;
      color: #ef4444;
    }

    .field-hint {
      margin-top: 0.375rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .form-error-box {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1rem;
      background: #fef2f2;
      border: 1px solid #fecaca;
      border-radius: 0.5rem;
      color: #dc2626;
      font-size: 0.8125rem;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
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
      transition: all 0.15s;
    }

    .btn--primary {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
      color: white;
      border: none;
    }

    .btn--primary:hover:not(:disabled) {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
    }

    .btn--primary:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .btn--secondary {
      background: white;
      color: #64748b;
      border: 1px solid #e2e8f0;
    }

    .btn--secondary:hover {
      background: #f8fafc;
    }

    @media (max-width: 480px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class TermFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);

  @Input() term: AcademicTerm | null = null;
  @Input() academicYear: AcademicYear | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<{
    name: string;
    startDate: string;
    endDate: string;
    sequence?: number;
  }>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  ngOnInit(): void {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Reinitialize form when term or academicYear input changes
    // This ensures form resets when opening for a new term after adding one
    if (changes['term'] || changes['academicYear']) {
      this.initForm();
    }
  }

  hasError(fieldName: string): boolean {
    const field = this.form.get(fieldName);
    return !!field && field.invalid && (field.dirty || field.touched);
  }

  private initForm(): void {
    // Calculate default sequence (next number after existing terms)
    const existingTermsCount = this.academicYear?.terms?.length || 0;
    const defaultSequence = this.term?.sequence || existingTermsCount + 1;

    this.form = this.fb.group({
      name: [
        this.term?.name || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(100),
        ],
      ],
      sequence: [
        defaultSequence,
        [
          Validators.required,
          Validators.min(1),
        ],
      ],
      startDate: [
        this.term?.startDate ? this.formatDateForInput(this.term.startDate) : '',
        [Validators.required],
      ],
      endDate: [
        this.term?.endDate ? this.formatDateForInput(this.term.endDate) : '',
        [Validators.required],
      ],
    }, {
      validators: [this.dateRangeValidator.bind(this), this.academicYearRangeValidator.bind(this)],
    });
  }

  formatDateForInput(dateStr: string): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toISOString().split('T')[0];
  }

  getMinDate(): string {
    return this.academicYear?.startDate
      ? this.formatDateForInput(this.academicYear.startDate)
      : '';
  }

  getMaxDate(): string {
    return this.academicYear?.endDate
      ? this.formatDateForInput(this.academicYear.endDate)
      : '';
  }

  private dateRangeValidator(group: FormGroup): { [key: string]: boolean } | null {
    const startDate = group.get('startDate')?.value;
    const endDate = group.get('endDate')?.value;

    if (startDate && endDate) {
      const start = new Date(startDate);
      const end = new Date(endDate);
      if (end <= start) {
        return { dateRange: true };
      }
    }
    return null;
  }

  private academicYearRangeValidator(group: FormGroup): { [key: string]: boolean } | null {
    if (!this.academicYear) return null;

    const startDate = group.get('startDate')?.value;
    const endDate = group.get('endDate')?.value;

    if (startDate && endDate) {
      const termStart = new Date(startDate);
      const termEnd = new Date(endDate);
      const yearStart = new Date(this.academicYear.startDate);
      const yearEnd = new Date(this.academicYear.endDate);

      if (termStart < yearStart || termEnd > yearEnd) {
        return { outsideRange: true };
      }
    }
    return null;
  }

  getFieldError(field: string): string | null {
    const control = this.form.get(field);
    if (!control || !control.touched || !control.errors) return null;

    if (control.errors['required']) return 'This field is required';
    if (control.errors['minlength']) {
      const min = control.errors['minlength'].requiredLength;
      return `Must be at least ${min} characters`;
    }
    if (control.errors['maxlength']) {
      const max = control.errors['maxlength'].requiredLength;
      return `Must be at most ${max} characters`;
    }
    if (control.errors['min']) {
      return 'Must be at least 1';
    }

    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      const value = this.form.value;
      this.save.emit({
        name: value.name,
        startDate: value.startDate,
        endDate: value.endDate,
        sequence: value.sequence,
      });
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
