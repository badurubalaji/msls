/**
 * MSLS Session Form Component
 *
 * Form component for creating and editing admission sessions.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import { SelectOption } from '../../../shared/components';
import {
  AdmissionSession,
  CreateSessionRequest,
  COMMON_DOCUMENTS,
} from './admission-session.model';
import { AdmissionSessionService } from './admission-session.service';

@Component({
  selector: 'msls-session-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="session-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="session">
          <i class="fa-solid" [class.fa-calendar-plus]="!session" [class.fa-pen]="session"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ session ? 'Edit Session' : 'Create New Session' }}</h3>
          <p>Configure admission cycle settings</p>
        </div>
      </div>

      <!-- Session Name -->
      <div class="form-field">
        <label class="field-label">
          Session Name <span class="required">*</span>
        </label>
        <div class="input-wrapper">
          <i class="fa-solid fa-tag input-icon"></i>
          <input
            type="text"
            formControlName="name"
            placeholder="e.g., Admission 2026-27 - Regular"
            class="form-input"
            [class.has-error]="hasError('name')"
          />
        </div>
        <span *ngIf="hasError('name')" class="field-error">{{ getFieldError('name') }}</span>
      </div>

      <div class="form-row">
        <!-- Academic Year -->
        <div class="form-field">
          <label class="field-label">
            Academic Year <span class="required">*</span>
          </label>
          <div class="input-wrapper">
            <i class="fa-solid fa-graduation-cap input-icon"></i>
            <select
              formControlName="academicYearId"
              class="form-input form-select"
              [class.has-error]="hasError('academicYearId')"
            >
              <option value="">Select academic year</option>
              <option *ngFor="let option of academicYearOptions()" [value]="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
          <span *ngIf="hasError('academicYearId')" class="field-error">{{ getFieldError('academicYearId') }}</span>
        </div>

        <!-- Application Fee -->
        <div class="form-field">
          <label class="field-label">Application Fee</label>
          <div class="input-wrapper">
            <i class="fa-solid fa-indian-rupee-sign input-icon"></i>
            <input
              type="number"
              formControlName="applicationFee"
              placeholder="0"
              min="0"
              class="form-input"
              [class.has-error]="hasError('applicationFee')"
            />
          </div>
          <span class="field-hint">Fee charged per application (optional)</span>
          <span *ngIf="hasError('applicationFee')" class="field-error">{{ getFieldError('applicationFee') }}</span>
        </div>
      </div>

      <div class="form-row">
        <!-- Start Date -->
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
            />
          </div>
          <span *ngIf="hasError('startDate')" class="field-error">{{ getFieldError('startDate') }}</span>
        </div>

        <!-- End Date -->
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
            />
          </div>
          <span *ngIf="hasError('endDate')" class="field-error">{{ getFieldError('endDate') }}</span>
        </div>
      </div>

      <!-- Date Range Error -->
      <div *ngIf="form.errors?.['dateRange']" class="form-error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        End date must be after start date
      </div>

      <!-- Required Documents -->
      <div class="documents-section">
        <label class="section-label">
          <i class="fa-solid fa-file-lines"></i>
          Required Documents
        </label>
        <p class="section-hint">Select the documents required for admission applications</p>
        <div class="documents-grid">
          <label *ngFor="let doc of commonDocuments" class="document-checkbox">
            <input
              type="checkbox"
              [checked]="isDocumentSelected(doc)"
              (change)="toggleDocument(doc)"
            />
            <span class="checkbox-custom"></span>
            <span class="checkbox-label">{{ doc }}</span>
          </label>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button type="submit" class="btn btn--primary" [disabled]="!form.valid || loading">
          <i *ngIf="loading" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!loading && !session" class="fa-solid fa-plus"></i>
          <i *ngIf="!loading && session" class="fa-solid fa-save"></i>
          {{ session ? 'Update Session' : 'Create Session' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .session-form {
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

    .documents-section {
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
    }

    .section-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #0f172a;
      margin-bottom: 0.25rem;
    }

    .section-label i {
      color: #64748b;
    }

    .section-hint {
      font-size: 0.75rem;
      color: #64748b;
      margin: 0 0 0.75rem 0;
    }

    .documents-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 0.5rem;
    }

    .document-checkbox {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: background 0.15s;
    }

    .document-checkbox:hover {
      background: #e2e8f0;
    }

    .document-checkbox input[type="checkbox"] {
      display: none;
    }

    .checkbox-custom {
      width: 1.125rem;
      height: 1.125rem;
      border: 2px solid #cbd5e1;
      border-radius: 0.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.15s;
      flex-shrink: 0;
    }

    .document-checkbox input:checked + .checkbox-custom {
      background: #6366f1;
      border-color: #6366f1;
    }

    .document-checkbox input:checked + .checkbox-custom::after {
      content: '';
      width: 0.375rem;
      height: 0.625rem;
      border: solid white;
      border-width: 0 2px 2px 0;
      transform: rotate(45deg);
      margin-bottom: 2px;
    }

    .checkbox-label {
      font-size: 0.8125rem;
      color: #334155;
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

    @media (max-width: 640px) {
      .form-row,
      .documents-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class SessionFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);
  private sessionService = inject(AdmissionSessionService);

  @Input() session: AdmissionSession | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<CreateSessionRequest>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;
  commonDocuments = COMMON_DOCUMENTS;
  selectedDocuments = signal<string[]>([]);

  /** Academic year options for select */
  academicYearOptions = signal<SelectOption[]>([]);

  ngOnInit(): void {
    this.initForm();
    this.loadAcademicYears();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Reinitialize form when session input changes
    if (changes['session'] && !changes['session'].firstChange) {
      this.initForm();
    }
  }

  hasError(fieldName: string): boolean {
    const field = this.form.get(fieldName);
    return !!field && field.invalid && (field.dirty || field.touched);
  }

  private formatDateForInput(dateStr: string | undefined): string {
    if (!dateStr) return '';
    // Handle ISO date strings
    if (dateStr.includes('T')) {
      return dateStr.split('T')[0];
    }
    return dateStr;
  }

  private initForm(): void {
    this.selectedDocuments.set(this.session?.requiredDocuments || []);

    this.form = this.fb.group({
      name: [
        this.session?.name || '',
        [
          Validators.required,
          Validators.minLength(3),
          Validators.maxLength(200),
        ],
      ],
      academicYearId: [
        this.session?.academicYearId || '',
        [Validators.required],
      ],
      startDate: [
        this.formatDateForInput(this.session?.startDate),
        [Validators.required],
      ],
      endDate: [
        this.formatDateForInput(this.session?.endDate),
        [Validators.required],
      ],
      applicationFee: [
        this.session?.applicationFee || 0,
        [Validators.min(0), Validators.max(100000)],
      ],
    }, {
      validators: this.dateRangeValidator.bind(this),
    });
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

  private loadAcademicYears(): void {
    this.sessionService.getAcademicYears().subscribe({
      next: (years) => {
        this.academicYearOptions.set(
          years.map(ay => ({
            value: ay.id,
            label: ay.name + (ay.isCurrent ? ' (Current)' : ''),
          }))
        );
      },
      error: (err) => {
        console.error('Failed to load academic years:', err);
      },
    });
  }

  isDocumentSelected(doc: string): boolean {
    return this.selectedDocuments().includes(doc);
  }

  toggleDocument(doc: string): void {
    const current = this.selectedDocuments();
    if (current.includes(doc)) {
      this.selectedDocuments.set(current.filter(d => d !== doc));
    } else {
      this.selectedDocuments.set([...current, doc]);
    }
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
      return 'Value must be 0 or greater';
    }
    if (control.errors['max']) {
      return 'Value is too large';
    }

    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      const value = this.form.value;
      const request: CreateSessionRequest = {
        name: value.name,
        academicYearId: value.academicYearId,
        startDate: value.startDate,
        endDate: value.endDate,
        applicationFee: Number(value.applicationFee) || 0,
        requiredDocuments: this.selectedDocuments(),
      };
      this.save.emit(request);
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
