/**
 * MSLS Enquiry Form Component
 *
 * Form component for creating and editing admission enquiries with organized sections.
 */

import {
  Component,
  input,
  output,
  inject,
  signal,
  effect,
  OnInit,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';

import { EnquiryService } from './enquiry.service';
import {
  Enquiry,
  CreateEnquiryDto,
  UpdateEnquiryDto,
  ENQUIRY_SOURCE_LABELS,
  ENQUIRY_STATUS_CONFIG,
} from './enquiry.model';

interface SelectOption {
  value: string;
  label: string;
}

/**
 * EnquiryFormComponent - Form for creating/editing enquiries with organized sections.
 */
@Component({
  selector: 'msls-enquiry-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="enquiry-form">
      <!-- Student Information Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--primary">
            <i class="fa-solid fa-user-graduate"></i>
          </div>
          <div class="section-title">
            <h3>Student Information</h3>
            <p>Basic details of the prospective student</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">
                Student Name <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-user input-icon"></i>
                <input
                  type="text"
                  formControlName="studentName"
                  placeholder="Enter student's full name"
                  class="form-input"
                  [class.has-error]="hasError('studentName')"
                />
              </div>
              <span *ngIf="hasError('studentName')" class="field-error">
                Student name is required
              </span>
            </div>
            <div class="form-field">
              <label class="field-label">
                Class Applying For <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-graduation-cap input-icon"></i>
                <select
                  formControlName="classApplying"
                  class="form-input form-select"
                  [class.has-error]="hasError('classApplying')"
                >
                  <option value="">Select class</option>
                  <option *ngFor="let opt of classOptions" [value]="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>
              <span *ngIf="hasError('classApplying')" class="field-error">
                Please select a class
              </span>
            </div>
          </div>
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">Date of Birth</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-calendar input-icon"></i>
                <input
                  type="date"
                  formControlName="dateOfBirth"
                  class="form-input"
                />
              </div>
              <span class="field-hint">Used for age verification</span>
            </div>
            <div class="form-field">
              <label class="field-label">Gender</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-venus-mars input-icon"></i>
                <select formControlName="gender" class="form-input form-select">
                  <option value="">Select gender</option>
                  <option *ngFor="let opt of genderOptions" [value]="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Parent/Guardian Information Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--info">
            <i class="fa-solid fa-users"></i>
          </div>
          <div class="section-title">
            <h3>Parent/Guardian Information</h3>
            <p>Contact details for the parent or guardian</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">
                Parent/Guardian Name <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-user-tie input-icon"></i>
                <input
                  type="text"
                  formControlName="parentName"
                  placeholder="Enter parent's full name"
                  class="form-input"
                  [class.has-error]="hasError('parentName')"
                />
              </div>
              <span *ngIf="hasError('parentName')" class="field-error">
                Parent name is required
              </span>
            </div>
            <div class="form-field">
              <label class="field-label">
                Phone Number <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-phone input-icon"></i>
                <input
                  type="tel"
                  formControlName="parentPhone"
                  placeholder="10-digit mobile number"
                  class="form-input"
                  [class.has-error]="hasError('parentPhone')"
                />
              </div>
              <span *ngIf="hasError('parentPhone')" class="field-error">
                <ng-container *ngIf="form.get('parentPhone')?.errors?.['required']">
                  Phone number is required
                </ng-container>
                <ng-container *ngIf="form.get('parentPhone')?.errors?.['pattern']">
                  Enter a valid 10-digit phone number
                </ng-container>
              </span>
              <span *ngIf="!hasError('parentPhone')" class="field-hint">
                For follow-up communication
              </span>
            </div>
          </div>
          <div class="form-row">
            <div class="form-field form-field--full">
              <label class="field-label">Email Address</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-envelope input-icon"></i>
                <input
                  type="email"
                  formControlName="parentEmail"
                  placeholder="parent@example.com"
                  class="form-input"
                  [class.has-error]="hasError('parentEmail')"
                />
              </div>
              <span *ngIf="hasError('parentEmail')" class="field-error">
                Please enter a valid email address
              </span>
              <span *ngIf="!hasError('parentEmail')" class="field-hint">
                Optional - for sending updates and notifications
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Enquiry Details Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--success">
            <i class="fa-solid fa-clipboard-question"></i>
          </div>
          <div class="section-title">
            <h3>Enquiry Details</h3>
            <p>Additional information about the enquiry</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">Source</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-bullhorn input-icon"></i>
                <select formControlName="source" class="form-input form-select">
                  <option *ngFor="let opt of sourceOptions" [value]="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>
              <span class="field-hint">How did they hear about us?</span>
            </div>
            <div class="form-field">
              <label class="field-label">Next Follow-up Date</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-bell input-icon"></i>
                <input
                  type="date"
                  formControlName="followUpDate"
                  class="form-input"
                />
              </div>
              <span class="field-hint">Set a reminder for follow-up</span>
            </div>
          </div>

          <!-- Status (only for edit mode) -->
          <div *ngIf="enquiry()" class="form-row">
            <div class="form-field">
              <label class="field-label">Status</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-flag input-icon"></i>
                <select formControlName="status" class="form-input form-select">
                  <option *ngFor="let opt of statusOptions" [value]="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-field"></div>
          </div>

          <!-- Referral Details (shown when source is referral) -->
          <div *ngIf="form.get('source')?.value === 'referral'" class="form-row">
            <div class="form-field form-field--full">
              <label class="field-label">Referral Details</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-share-nodes input-icon"></i>
                <input
                  type="text"
                  formControlName="referralDetails"
                  placeholder="Name of the person who referred"
                  class="form-input"
                />
              </div>
              <span class="field-hint">Track referral sources for analytics</span>
            </div>
          </div>

          <div class="form-row">
            <div class="form-field form-field--full">
              <label class="field-label">Remarks / Notes</label>
              <div class="input-wrapper input-wrapper--textarea">
                <i class="fa-solid fa-comment input-icon"></i>
                <textarea
                  formControlName="remarks"
                  rows="3"
                  placeholder="Add any notes from the conversation..."
                  class="form-input form-textarea"
                ></textarea>
              </div>
              <span class="field-hint">
                Record important details from initial discussion
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button
          type="submit"
          class="btn btn--primary"
          [disabled]="form.invalid || submitting()"
        >
          <i *ngIf="submitting()" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!submitting() && !enquiry()" class="fa-solid fa-plus"></i>
          <i *ngIf="!submitting() && enquiry()" class="fa-solid fa-save"></i>
          {{ enquiry() ? 'Update Enquiry' : 'Create Enquiry' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .enquiry-form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    /* Form Section */
    .form-section {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
    }

    .section-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.25rem;
      background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
      border-bottom: 1px solid #e2e8f0;
    }

    .section-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.625rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 1rem;
      flex-shrink: 0;
    }

    .section-icon--primary {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
    }

    .section-icon--info {
      background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
    }

    .section-icon--success {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    }

    .section-title h3 {
      margin: 0;
      font-size: 0.9375rem;
      font-weight: 600;
      color: #0f172a;
    }

    .section-title p {
      margin: 0.125rem 0 0 0;
      font-size: 0.8125rem;
      color: #64748b;
    }

    .section-content {
      padding: 1.25rem;
    }

    /* Form Layout */
    .form-row {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .form-row:last-child {
      margin-bottom: 0;
    }

    @media (max-width: 640px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }

    .form-field {
      display: flex;
      flex-direction: column;
    }

    .form-field--full {
      grid-column: 1 / -1;
    }

    /* Field Label */
    .field-label {
      font-size: 0.8125rem;
      font-weight: 500;
      color: #374151;
      margin-bottom: 0.5rem;
    }

    .required {
      color: #ef4444;
    }

    /* Input Wrapper */
    .input-wrapper {
      position: relative;
    }

    .input-wrapper--textarea {
      display: flex;
      align-items: flex-start;
    }

    .input-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      font-size: 0.875rem;
      color: #94a3b8;
      pointer-events: none;
    }

    .input-wrapper--textarea .input-icon {
      top: 0.875rem;
      transform: none;
    }

    /* Form Input */
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

    .form-input.has-error:focus {
      box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
    }

    .form-input::placeholder {
      color: #94a3b8;
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

    .form-textarea {
      height: auto;
      min-height: 5rem;
      padding-top: 0.625rem;
      padding-bottom: 0.625rem;
      resize: vertical;
    }

    /* Field Hint & Error */
    .field-hint {
      display: block;
      margin-top: 0.375rem;
      font-size: 0.75rem;
      color: #94a3b8;
    }

    .field-error {
      display: block;
      margin-top: 0.375rem;
      font-size: 0.75rem;
      color: #ef4444;
    }

    /* Form Actions */
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
      justify-content: center;
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
      color: #374151;
    }

    /* Date input styling */
    input[type="date"]::-webkit-calendar-picker-indicator {
      cursor: pointer;
      opacity: 0.6;
    }

    input[type="date"]::-webkit-calendar-picker-indicator:hover {
      opacity: 1;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EnquiryFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly enquiryService = inject(EnquiryService);

  /** Enquiry to edit (null for new enquiry) */
  enquiry = input<Enquiry | null>(null);

  /** Emitted when form is submitted successfully */
  submitted = output<Enquiry>();

  /** Emitted when form is cancelled */
  cancelled = output<void>();

  /** Form submitting state */
  readonly submitting = signal(false);

  /** Reactive form */
  form!: FormGroup;

  constructor() {
    // Effect to reinitialize form when enquiry input changes
    effect(() => {
      const enquiry = this.enquiry();
      // Only reinitialize if form is already created (not first run)
      if (this.form) {
        this.initForm();
      }
    });
  }

  /** Class options */
  readonly classOptions: SelectOption[] = [
    { value: 'Nursery', label: 'Nursery' },
    { value: 'LKG', label: 'LKG' },
    { value: 'UKG', label: 'UKG' },
    { value: 'Class 1', label: 'Class 1' },
    { value: 'Class 2', label: 'Class 2' },
    { value: 'Class 3', label: 'Class 3' },
    { value: 'Class 4', label: 'Class 4' },
    { value: 'Class 5', label: 'Class 5' },
    { value: 'Class 6', label: 'Class 6' },
    { value: 'Class 7', label: 'Class 7' },
    { value: 'Class 8', label: 'Class 8' },
    { value: 'Class 9', label: 'Class 9' },
    { value: 'Class 10', label: 'Class 10' },
    { value: 'Class 11', label: 'Class 11' },
    { value: 'Class 12', label: 'Class 12' },
  ];

  /** Gender options */
  readonly genderOptions: SelectOption[] = [
    { value: 'male', label: 'Male' },
    { value: 'female', label: 'Female' },
    { value: 'other', label: 'Other' },
  ];

  /** Source options */
  readonly sourceOptions: SelectOption[] = Object.entries(ENQUIRY_SOURCE_LABELS).map(
    ([value, label]) => ({ value, label })
  );

  /** Status options */
  readonly statusOptions: SelectOption[] = Object.entries(ENQUIRY_STATUS_CONFIG).map(
    ([value, config]) => ({ value, label: config.label })
  );

  ngOnInit(): void {
    this.initForm();
  }

  /**
   * Initialize form with default or existing values
   */
  private initForm(): void {
    const enquiry = this.enquiry();

    this.form = this.fb.group({
      studentName: [enquiry?.studentName || '', [Validators.required, Validators.maxLength(200)]],
      classApplying: [enquiry?.classApplying || '', [Validators.required]],
      dateOfBirth: [this.formatDateForInput(enquiry?.dateOfBirth)],
      gender: [enquiry?.gender || ''],
      parentName: [enquiry?.parentName || '', [Validators.required, Validators.maxLength(200)]],
      parentPhone: [
        enquiry?.parentPhone || '',
        [Validators.required, Validators.pattern(/^[0-9]{10}$/)],
      ],
      parentEmail: [enquiry?.parentEmail || '', [Validators.email]],
      source: [enquiry?.source || 'walk_in'],
      referralDetails: [enquiry?.referralDetails || ''],
      followUpDate: [this.formatDateForInput(enquiry?.followUpDate)],
      remarks: [enquiry?.remarks || ''],
      status: [enquiry?.status || 'new'],
    });
  }

  /**
   * Format date string for HTML date input (YYYY-MM-DD)
   */
  private formatDateForInput(dateStr: string | undefined): string {
    if (!dateStr) return '';
    // Handle ISO date format (e.g., "2025-01-24T00:00:00Z")
    if (dateStr.includes('T')) {
      return dateStr.split('T')[0];
    }
    return dateStr;
  }

  /**
   * Check if field has validation error
   */
  hasError(fieldName: string): boolean {
    const field = this.form.get(fieldName);
    return !!field && field.invalid && (field.dirty || field.touched);
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    const formValue = this.form.value;
    const enquiry = this.enquiry();

    if (enquiry) {
      // Update existing enquiry
      const dto: UpdateEnquiryDto = {
        studentName: formValue.studentName,
        classApplying: formValue.classApplying,
        dateOfBirth: formValue.dateOfBirth || undefined,
        gender: formValue.gender || undefined,
        parentName: formValue.parentName,
        parentPhone: formValue.parentPhone,
        parentEmail: formValue.parentEmail || undefined,
        source: formValue.source,
        referralDetails: formValue.referralDetails || undefined,
        followUpDate: formValue.followUpDate || undefined,
        remarks: formValue.remarks || undefined,
        status: formValue.status,
      };

      this.enquiryService.updateEnquiry(enquiry.id, dto).subscribe({
        next: (updated) => {
          this.submitting.set(false);
          this.submitted.emit(updated);
        },
        error: () => {
          this.submitting.set(false);
        },
      });
    } else {
      // Create new enquiry
      const dto: CreateEnquiryDto = {
        studentName: formValue.studentName,
        classApplying: formValue.classApplying,
        dateOfBirth: formValue.dateOfBirth || undefined,
        gender: formValue.gender || undefined,
        parentName: formValue.parentName,
        parentPhone: formValue.parentPhone,
        parentEmail: formValue.parentEmail || undefined,
        source: formValue.source,
        referralDetails: formValue.referralDetails || undefined,
        followUpDate: formValue.followUpDate || undefined,
        remarks: formValue.remarks || undefined,
      };

      this.enquiryService.createEnquiry(dto).subscribe({
        next: (created) => {
          this.submitting.set(false);
          this.submitted.emit(created);
        },
        error: () => {
          this.submitting.set(false);
        },
      });
    }
  }

  /**
   * Handle cancel action
   */
  onCancel(): void {
    this.cancelled.emit();
  }
}
