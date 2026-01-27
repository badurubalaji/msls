/**
 * MSLS Academic Year Form Component
 *
 * Form component for creating and editing academic years.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import {
  MslsInputComponent,
  MslsFormFieldComponent,
  MslsCheckboxComponent,
} from '../../../shared/components';
import { AcademicYear } from './academic-year.model';

@Component({
  selector: 'msls-academic-year-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsInputComponent,
    MslsFormFieldComponent,
    MslsCheckboxComponent,
  ],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="academic-year-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="academicYear">
          <i class="fa-solid" [class.fa-calendar-plus]="!academicYear" [class.fa-pen]="academicYear"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ academicYear ? 'Edit Academic Year' : 'Create Academic Year' }}</h3>
          <p>{{ academicYear ? 'Update academic year details' : 'Set up a new academic year for your institution' }}</p>
        </div>
      </div>

      <msls-form-field
        label="Academic Year Name"
        [required]="true"
        [error]="getFieldError('name') || ''"
        hint="e.g., 2025-26"
      >
        <msls-input
          type="text"
          formControlName="name"
          placeholder="Enter academic year name"
        />
      </msls-form-field>

      <div class="form-row">
        <msls-form-field
          label="Start Date"
          [required]="true"
          [error]="getFieldError('startDate') || ''"
        >
          <input
            type="date"
            class="form-input"
            formControlName="startDate"
          />
        </msls-form-field>

        <msls-form-field
          label="End Date"
          [required]="true"
          [error]="getFieldError('endDate') || ''"
        >
          <input
            type="date"
            class="form-input"
            formControlName="endDate"
          />
        </msls-form-field>
      </div>

      @if (form.errors?.['dateRange']) {
        <div class="form-error">
          <i class="fa-solid fa-circle-exclamation"></i>
          End date must be after start date
        </div>
      }

      <div class="form-checkbox-field">
        <msls-checkbox
          formControlName="isCurrent"
          label="Set as current academic year"
        />
        <p class="field-hint">
          The current academic year will be used as default across all modules
        </p>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button type="submit" class="btn btn--primary" [disabled]="!form.valid || loading">
          <i *ngIf="loading" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!loading && !academicYear" class="fa-solid fa-plus"></i>
          <i *ngIf="!loading && academicYear" class="fa-solid fa-save"></i>
          {{ academicYear ? 'Update Academic Year' : 'Create Academic Year' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .academic-year-form {
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
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
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

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .form-input {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: #ffffff;
      color: #0f172a;
      transition: border-color 0.15s, box-shadow 0.15s;
    }

    .form-input:focus {
      outline: none;
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }

    .form-input:disabled {
      background: #f1f5f9;
      cursor: not-allowed;
    }

    .form-error {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 0.875rem;
      background: #fef2f2;
      border: 1px solid #fecaca;
      border-radius: 0.5rem;
      color: #dc2626;
      font-size: 0.8125rem;
    }

    .form-checkbox-field {
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
    }

    .field-hint {
      margin: 0.5rem 0 0 1.75rem;
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
      border-color: #cbd5e1;
    }

    @media (max-width: 480px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class AcademicYearFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);

  @Input() academicYear: AcademicYear | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<{
    name: string;
    startDate: string;
    endDate: string;
    isCurrent?: boolean;
  }>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  ngOnInit(): void {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Reinitialize form when academicYear input changes
    // This ensures form resets when opening for a new academic year after adding one
    if (changes['academicYear']) {
      this.initForm();
    }
  }

  private initForm(): void {
    this.form = this.fb.group({
      name: [
        this.academicYear?.name || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(50),
        ],
      ],
      startDate: [
        this.academicYear?.startDate ? this.formatDateForInput(this.academicYear.startDate) : '',
        [Validators.required],
      ],
      endDate: [
        this.academicYear?.endDate ? this.formatDateForInput(this.academicYear.endDate) : '',
        [Validators.required],
      ],
      isCurrent: [this.academicYear?.isCurrent || false],
    }, {
      validators: this.dateRangeValidator,
    });
  }

  private formatDateForInput(dateStr: string): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toISOString().split('T')[0];
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
        startDate: value.startDate,
        endDate: value.endDate,
        isCurrent: value.isCurrent || undefined,
      });
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
