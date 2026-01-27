/**
 * MSLS Holiday Form Component
 *
 * Form component for creating and editing holidays.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import { AcademicYear, Holiday, HolidayType, HOLIDAY_TYPES } from './academic-year.model';

@Component({
  selector: 'msls-holiday-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="holiday-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="holiday">
          <i class="fa-solid" [class.fa-umbrella-beach]="!holiday" [class.fa-pen]="holiday"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ holiday ? 'Edit Holiday' : 'Add New Holiday' }}</h3>
          <p *ngIf="academicYear">{{ academicYear.name }}</p>
        </div>
      </div>

      <!-- Holiday Name -->
      <div class="form-field">
        <label class="field-label">
          Holiday Name <span class="required">*</span>
        </label>
        <div class="input-wrapper">
          <i class="fa-solid fa-tag input-icon"></i>
          <input
            type="text"
            formControlName="name"
            placeholder="e.g., Diwali, Christmas, Republic Day"
            class="form-input"
            [class.has-error]="hasError('name')"
          />
        </div>
        <span *ngIf="hasError('name')" class="field-error">{{ getFieldError('name') }}</span>
      </div>

      <!-- Date and Type Row -->
      <div class="form-row">
        <div class="form-field">
          <label class="field-label">
            Date <span class="required">*</span>
          </label>
          <div class="input-wrapper">
            <i class="fa-solid fa-calendar input-icon"></i>
            <input
              type="date"
              formControlName="date"
              class="form-input"
              [class.has-error]="hasError('date')"
              [min]="getMinDate()"
              [max]="getMaxDate()"
            />
          </div>
          <span *ngIf="hasError('date')" class="field-error">{{ getFieldError('date') }}</span>
        </div>

        <div class="form-field">
          <label class="field-label">Type</label>
          <div class="input-wrapper">
            <i class="fa-solid fa-shapes input-icon"></i>
            <select formControlName="type" class="form-input form-select">
              <option *ngFor="let opt of holidayTypeOptions" [value]="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Date Range Error -->
      <div *ngIf="form.errors?.['outsideRange']" class="form-error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        Holiday date must be within the academic year period
      </div>

      <!-- Optional Holiday Toggle -->
      <div class="form-field-toggle">
        <label class="toggle-wrapper">
          <input type="checkbox" formControlName="isOptional" class="toggle-input" />
          <span class="toggle-slider"></span>
          <span class="toggle-label">Optional Holiday</span>
        </label>
        <p class="field-hint">
          Optional holidays allow staff to choose whether to take the day off
        </p>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button type="submit" class="btn btn--primary" [disabled]="!form.valid || loading">
          <i *ngIf="loading" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!loading && !holiday" class="fa-solid fa-plus"></i>
          <i *ngIf="!loading && holiday" class="fa-solid fa-save"></i>
          {{ holiday ? 'Update Holiday' : 'Add Holiday' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .holiday-form {
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
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
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

    .form-select {
      cursor: pointer;
      appearance: none;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='%2394a3b8'%3E%3Cpath fill-rule='evenodd' d='M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z' clip-rule='evenodd'/%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 0.75rem center;
      background-size: 1rem;
      padding-right: 2.5rem;
    }

    .field-error {
      margin-top: 0.375rem;
      font-size: 0.75rem;
      color: #ef4444;
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

    .form-field-toggle {
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
    }

    .toggle-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      cursor: pointer;
    }

    .toggle-input {
      position: absolute;
      opacity: 0;
      width: 0;
      height: 0;
    }

    .toggle-slider {
      position: relative;
      width: 2.75rem;
      height: 1.5rem;
      background: #cbd5e1;
      border-radius: 1rem;
      transition: all 0.2s;
    }

    .toggle-slider::before {
      content: '';
      position: absolute;
      top: 0.125rem;
      left: 0.125rem;
      width: 1.25rem;
      height: 1.25rem;
      background: white;
      border-radius: 50%;
      transition: all 0.2s;
      box-shadow: 0 1px 3px rgba(0,0,0,0.2);
    }

    .toggle-input:checked + .toggle-slider {
      background: #6366f1;
    }

    .toggle-input:checked + .toggle-slider::before {
      transform: translateX(1.25rem);
    }

    .toggle-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .field-hint {
      margin: 0.5rem 0 0 3.5rem;
      font-size: 0.75rem;
      color: #64748b;
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
export class HolidayFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);

  @Input() holiday: Holiday | null = null;
  @Input() academicYear: AcademicYear | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<{
    name: string;
    date: string;
    type?: HolidayType;
    isOptional?: boolean;
  }>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  // Convert HOLIDAY_TYPES to select options format
  holidayTypeOptions = HOLIDAY_TYPES.map(t => ({
    value: t.value,
    label: t.label,
  }));

  ngOnInit(): void {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Reinitialize form when holiday or academicYear input changes
    // This ensures form resets when opening for a new holiday after adding one
    if (changes['holiday'] || changes['academicYear']) {
      this.initForm();
    }
  }

  hasError(fieldName: string): boolean {
    const field = this.form.get(fieldName);
    return !!field && field.invalid && (field.dirty || field.touched);
  }

  private initForm(): void {
    this.form = this.fb.group({
      name: [
        this.holiday?.name || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(200),
        ],
      ],
      date: [
        this.holiday?.date ? this.formatDateForInput(this.holiday.date) : '',
        [Validators.required],
      ],
      type: [
        this.holiday?.type || 'public',
      ],
      isOptional: [
        this.holiday?.isOptional || false,
      ],
    }, {
      validators: this.academicYearRangeValidator.bind(this),
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

  private academicYearRangeValidator(group: FormGroup): { [key: string]: boolean } | null {
    if (!this.academicYear) return null;

    const date = group.get('date')?.value;

    if (date) {
      const holidayDate = new Date(date);
      const yearStart = new Date(this.academicYear.startDate);
      const yearEnd = new Date(this.academicYear.endDate);

      if (holidayDate < yearStart || holidayDate > yearEnd) {
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

    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      const value = this.form.value;
      this.save.emit({
        name: value.name,
        date: value.date,
        type: value.type || undefined,
        isOptional: !!value.isOptional,
      });
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
