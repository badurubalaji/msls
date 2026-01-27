/**
 * MSLS Branch Form Component
 *
 * Form component for creating and editing branches with organized sections.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import { SelectOption } from '../../../shared/components';
import { Branch, CreateBranchRequest, TIMEZONES } from './branch.model';

@Component({
  selector: 'msls-branch-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
  ],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="branch-form">
      <!-- Basic Information Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--primary">
            <i class="fa-solid fa-building"></i>
          </div>
          <div class="section-title">
            <h3>Basic Information</h3>
            <p>Enter the branch identity details</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">
                Branch Code <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-hashtag input-icon"></i>
                <input
                  type="text"
                  formControlName="code"
                  placeholder="e.g., MAIN, BR01"
                  class="form-input"
                  [class.has-error]="hasError('code')"
                />
              </div>
              @if (hasError('code')) {
                <span class="field-error">{{ getFieldError('code') }}</span>
              }
              <span class="field-hint">Unique identifier for the branch</span>
            </div>
            <div class="form-field">
              <label class="field-label">
                Branch Name <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-font input-icon"></i>
                <input
                  type="text"
                  formControlName="name"
                  placeholder="e.g., Main Campus"
                  class="form-input"
                  [class.has-error]="hasError('name')"
                />
              </div>
              @if (hasError('name')) {
                <span class="field-error">{{ getFieldError('name') }}</span>
              }
            </div>
          </div>
        </div>
      </div>

      <!-- Address Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--info">
            <i class="fa-solid fa-location-dot"></i>
          </div>
          <div class="section-title">
            <h3>Address</h3>
            <p>Physical location of the branch</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field form-field--full">
              <label class="field-label">Street Address</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-road input-icon"></i>
                <input
                  type="text"
                  formControlName="addressLine1"
                  placeholder="123 Main Street"
                  class="form-input"
                />
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="form-field form-field--full">
              <label class="field-label">Address Line 2</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-building input-icon"></i>
                <input
                  type="text"
                  formControlName="addressLine2"
                  placeholder="Apartment, suite, floor (optional)"
                  class="form-input"
                />
              </div>
            </div>
          </div>
          <div class="form-row form-row--3">
            <div class="form-field">
              <label class="field-label">City</label>
              <input
                type="text"
                formControlName="city"
                placeholder="City"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label class="field-label">State / Province</label>
              <input
                type="text"
                formControlName="state"
                placeholder="State"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label class="field-label">Postal Code</label>
              <input
                type="text"
                formControlName="postalCode"
                placeholder="123456"
                class="form-input"
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">Country</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-globe input-icon"></i>
                <input
                  type="text"
                  formControlName="country"
                  placeholder="Country"
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-field">
              <label class="field-label">Timezone</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-clock input-icon"></i>
                <select formControlName="timezone" class="form-input form-select">
                  @for (tz of timezoneOptions; track tz.value) {
                    <option [value]="tz.value">{{ tz.label }}</option>
                  }
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact Information Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--success">
            <i class="fa-solid fa-address-book"></i>
          </div>
          <div class="section-title">
            <h3>Contact Information</h3>
            <p>How to reach this branch</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <div class="form-field">
              <label class="field-label">Phone Number</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-phone input-icon"></i>
                <input
                  type="tel"
                  formControlName="phone"
                  placeholder="+91 98765 43210"
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-field">
              <label class="field-label">Email Address</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-envelope input-icon"></i>
                <input
                  type="email"
                  formControlName="email"
                  placeholder="branch&#64;school.edu"
                  class="form-input"
                  [class.has-error]="hasError('email')"
                />
              </div>
              @if (hasError('email')) {
                <span class="field-error">{{ getFieldError('email') }}</span>
              }
            </div>
          </div>
        </div>
      </div>

      <!-- Primary Branch Toggle -->
      <div class="primary-toggle-card">
        <div class="toggle-content">
          <div class="toggle-icon">
            <i class="fa-solid fa-star"></i>
          </div>
          <div class="toggle-text">
            <span class="toggle-label">Set as Primary Branch</span>
            <span class="toggle-hint">Primary branch is the default for new users and settings</span>
          </div>
        </div>
        <label class="toggle-switch">
          <input type="checkbox" formControlName="isPrimary" />
          <span class="toggle-slider"></span>
        </label>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" class="btn btn-secondary" (click)="onCancel()">
          <i class="fa-solid fa-xmark"></i>
          Cancel
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          [disabled]="!form.valid || loading"
        >
          @if (loading) {
            <span class="btn-spinner"></span>
            Saving...
          } @else {
            <i class="fa-solid fa-check"></i>
            {{ branch ? 'Update Branch' : 'Create Branch' }}
          }
        </button>
      </div>
    </form>
  `,
  styles: [`
    .branch-form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    /* Section Styles */
    .form-section {
      background: #ffffff;
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
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1rem;
      flex-shrink: 0;
    }

    .section-icon--primary {
      background: linear-gradient(135deg, #eef2ff 0%, #e0e7ff 100%);
      color: #4f46e5;
    }

    .section-icon--info {
      background: linear-gradient(135deg, #ecfeff 0%, #cffafe 100%);
      color: #0891b2;
    }

    .section-icon--success {
      background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
      color: #16a34a;
    }

    .section-title h3 {
      margin: 0;
      font-size: 0.9375rem;
      font-weight: 600;
      color: #1e293b;
    }

    .section-title p {
      margin: 0.125rem 0 0;
      font-size: 0.8125rem;
      color: #64748b;
    }

    .section-content {
      padding: 1.25rem;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    /* Form Row */
    .form-row {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .form-row--3 {
      grid-template-columns: repeat(3, 1fr);
    }

    /* Form Field */
    .form-field {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .form-field--full {
      grid-column: 1 / -1;
    }

    .field-label {
      font-size: 0.8125rem;
      font-weight: 500;
      color: #475569;
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
      pointer-events: none;
    }

    .form-input {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: #ffffff;
      color: #1e293b;
      transition: all 0.2s ease;
    }

    .input-wrapper .form-input {
      padding-left: 2.5rem;
    }

    .form-input::placeholder {
      color: #94a3b8;
    }

    .form-input:hover {
      border-color: #cbd5e1;
    }

    .form-input:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .form-input.has-error {
      border-color: #ef4444;
    }

    .form-input.has-error:focus {
      box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
    }

    .form-select {
      cursor: pointer;
      appearance: none;
      background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%2394a3b8' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
      background-position: right 0.75rem center;
      background-repeat: no-repeat;
      background-size: 1.25rem;
      padding-right: 2.5rem;
    }

    .field-error {
      font-size: 0.75rem;
      color: #ef4444;
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .field-hint {
      font-size: 0.75rem;
      color: #94a3b8;
    }

    /* Primary Toggle Card */
    .primary-toggle-card {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1rem 1.25rem;
      background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
      border: 1px solid #fcd34d;
      border-radius: 0.75rem;
    }

    .toggle-content {
      display: flex;
      align-items: center;
      gap: 0.875rem;
    }

    .toggle-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
      color: #d97706;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1rem;
    }

    .toggle-text {
      display: flex;
      flex-direction: column;
    }

    .toggle-label {
      font-size: 0.875rem;
      font-weight: 600;
      color: #92400e;
    }

    .toggle-hint {
      font-size: 0.75rem;
      color: #b45309;
    }

    /* Toggle Switch */
    .toggle-switch {
      position: relative;
      display: inline-block;
      width: 3rem;
      height: 1.625rem;
      cursor: pointer;
    }

    .toggle-switch input {
      opacity: 0;
      width: 0;
      height: 0;
    }

    .toggle-slider {
      position: absolute;
      inset: 0;
      background: #e2e8f0;
      border-radius: 9999px;
      transition: all 0.3s ease;
    }

    .toggle-slider::before {
      content: '';
      position: absolute;
      width: 1.25rem;
      height: 1.25rem;
      left: 0.1875rem;
      bottom: 0.1875rem;
      background: white;
      border-radius: 50%;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
      transition: all 0.3s ease;
    }

    .toggle-switch input:checked + .toggle-slider {
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    }

    .toggle-switch input:checked + .toggle-slider::before {
      transform: translateX(1.375rem);
    }

    /* Form Actions */
    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 0.5rem;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 0.5rem;
      border: none;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
      border-color: #cbd5e1;
    }

    .btn-primary {
      background: linear-gradient(135deg, #4f46e5 0%, #4338ca 100%);
      color: white;
      box-shadow: 0 2px 4px rgba(79, 70, 229, 0.3);
    }

    .btn-primary:hover:not(:disabled) {
      background: linear-gradient(135deg, #4338ca 0%, #3730a3 100%);
      box-shadow: 0 4px 8px rgba(79, 70, 229, 0.4);
      transform: translateY(-1px);
    }

    .btn-primary:disabled {
      opacity: 0.6;
      cursor: not-allowed;
      transform: none;
    }

    .btn-spinner {
      width: 1rem;
      height: 1rem;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-top-color: white;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Responsive */
    @media (max-width: 640px) {
      .form-row,
      .form-row--3 {
        grid-template-columns: 1fr;
      }

      .section-header {
        padding: 0.875rem 1rem;
      }

      .section-content {
        padding: 1rem;
      }

      .primary-toggle-card {
        flex-direction: column;
        gap: 1rem;
        align-items: flex-start;
      }

      .form-actions {
        flex-direction: column-reverse;
      }

      .btn {
        width: 100%;
      }
    }
  `],
})
export class BranchFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);

  @Input() branch: Branch | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<CreateBranchRequest>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  /** Timezone options for select */
  timezoneOptions: SelectOption[] = TIMEZONES.map(tz => ({
    value: tz.value,
    label: tz.label,
  }));

  ngOnInit(): void {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['branch'] && !changes['branch'].firstChange) {
      this.initForm();
    }
  }

  private initForm(): void {
    this.form = this.fb.group({
      code: [
        this.branch?.code || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(20),
        ],
      ],
      name: [
        this.branch?.name || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(200),
        ],
      ],
      addressLine1: [
        this.branch?.addressLine1 || '',
        [Validators.maxLength(255)],
      ],
      addressLine2: [
        this.branch?.addressLine2 || '',
        [Validators.maxLength(255)],
      ],
      city: [
        this.branch?.city || '',
        [Validators.maxLength(100)],
      ],
      state: [
        this.branch?.state || '',
        [Validators.maxLength(100)],
      ],
      postalCode: [
        this.branch?.postalCode || '',
        [Validators.maxLength(20)],
      ],
      country: [
        this.branch?.country || 'India',
        [Validators.maxLength(100)],
      ],
      phone: [
        this.branch?.phone || '',
        [Validators.maxLength(20)],
      ],
      email: [
        this.branch?.email || '',
        [Validators.email, Validators.maxLength(255)],
      ],
      timezone: [
        this.branch?.timezone || 'Asia/Kolkata',
        [Validators.maxLength(50)],
      ],
      isPrimary: [this.branch?.isPrimary || false],
    });
  }

  hasError(field: string): boolean {
    const control = this.form.get(field);
    return !!(control && control.touched && control.errors);
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
    if (control.errors['email']) return 'Please enter a valid email address';

    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      const value = this.form.value;
      const request: CreateBranchRequest = {
        code: value.code,
        name: value.name,
        addressLine1: value.addressLine1 || undefined,
        addressLine2: value.addressLine2 || undefined,
        city: value.city || undefined,
        state: value.state || undefined,
        postalCode: value.postalCode || undefined,
        country: value.country || undefined,
        phone: value.phone || undefined,
        email: value.email || undefined,
        timezone: value.timezone || undefined,
        isPrimary: value.isPrimary,
      };
      this.save.emit(request);
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
