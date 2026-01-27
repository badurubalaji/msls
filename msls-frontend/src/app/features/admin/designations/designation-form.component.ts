import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Designation, CreateDesignationRequest, DESIGNATION_LEVELS } from './designation.model';
import { DepartmentService } from '../departments/department.service';
import { DepartmentDropdownItem } from '../departments/department.model';

interface SelectOption {
  value: string | number;
  label: string;
}

@Component({
  selector: 'msls-designation-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="designation-form">
      <!-- Basic Information Section -->
      <div class="form-section">
        <div class="section-header">
          <div class="section-icon section-icon--primary">
            <i class="fa-solid fa-user-tie"></i>
          </div>
          <div class="section-title">
            <h3>Designation Details</h3>
            <p>Define the position title and hierarchy level</p>
          </div>
        </div>
        <div class="section-content">
          <div class="form-row">
            <!-- Name Field -->
            <div class="form-field form-field--full">
              <label class="field-label">
                Designation Name <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-briefcase input-icon"></i>
                <input
                  type="text"
                  formControlName="name"
                  placeholder="e.g., Principal, Teacher, Accountant"
                  class="form-input"
                  [class.has-error]="hasError('name')"
                />
              </div>
              @if (hasError('name')) {
                <span class="field-error">{{ getFieldError('name') }}</span>
              }
            </div>
          </div>

          <div class="form-row">
            <!-- Level Field -->
            <div class="form-field">
              <label class="field-label">
                Hierarchy Level <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <i class="fa-solid fa-layer-group input-icon"></i>
                <select
                  formControlName="level"
                  class="form-input form-select"
                  [class.has-error]="hasError('level')"
                >
                  @for (level of levelOptions; track level.value) {
                    <option [value]="level.value">{{ level.label }}</option>
                  }
                </select>
              </div>
              @if (hasError('level')) {
                <span class="field-error">{{ getFieldError('level') }}</span>
              }
              <span class="field-hint">Level 1 is highest (Executive), Level 10 is lowest (Entry)</span>
            </div>

            <!-- Department Field -->
            <div class="form-field">
              <label class="field-label">Department</label>
              <div class="input-wrapper">
                <i class="fa-solid fa-building-user input-icon"></i>
                <select
                  formControlName="departmentId"
                  class="form-input form-select"
                >
                  <option value="">No specific department</option>
                  @for (dept of departmentOptions; track dept.value) {
                    <option [value]="dept.value">{{ dept.label }}</option>
                  }
                </select>
              </div>
              <span class="field-hint">Leave empty if designation applies to all departments</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Status Toggle -->
      <div class="status-toggle-card">
        <div class="toggle-content">
          <div class="toggle-icon">
            <i class="fa-solid fa-toggle-on"></i>
          </div>
          <div class="toggle-text">
            <span class="toggle-label">Active Status</span>
            <span class="toggle-hint">Inactive designations won't appear in staff forms</span>
          </div>
        </div>
        <label class="toggle-switch">
          <input type="checkbox" formControlName="isActive" />
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
            {{ designation ? 'Update' : 'Create' }} Designation
          }
        </button>
      </div>
    </form>
  `,
  styles: [`
    .designation-form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    .form-section {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .section-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.5rem;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .section-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1rem;
    }

    .section-icon--primary {
      background: #eef2ff;
      color: #4f46e5;
    }

    .section-title h3 {
      margin: 0;
      font-size: 1rem;
      font-weight: 600;
      color: #1e293b;
    }

    .section-title p {
      margin: 0.25rem 0 0;
      font-size: 0.875rem;
      color: #64748b;
    }

    .section-content {
      padding: 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .form-row {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .form-field {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .form-field--full {
      grid-column: 1 / -1;
    }

    .field-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
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
      color: #9ca3af;
      font-size: 0.875rem;
      pointer-events: none;
    }

    .form-input {
      width: 100%;
      padding: 0.625rem 0.875rem 0.625rem 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: all 0.2s;
      background: white;
    }

    .form-input:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .form-input.has-error {
      border-color: #ef4444;
    }

    .form-select {
      appearance: none;
      cursor: pointer;
      padding-right: 2.5rem;
      background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
      background-position: right 0.5rem center;
      background-repeat: no-repeat;
      background-size: 1.5em 1.5em;
    }

    .field-error {
      font-size: 0.75rem;
      color: #ef4444;
    }

    .field-hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .status-toggle-card {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1rem 1.5rem;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
    }

    .toggle-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .toggle-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #dcfce7;
      color: #16a34a;
      font-size: 1rem;
    }

    .toggle-text {
      display: flex;
      flex-direction: column;
    }

    .toggle-label {
      font-weight: 500;
      color: #1e293b;
    }

    .toggle-hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .toggle-switch {
      position: relative;
      display: inline-block;
      width: 48px;
      height: 26px;
    }

    .toggle-switch input {
      opacity: 0;
      width: 0;
      height: 0;
    }

    .toggle-slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: #cbd5e1;
      transition: 0.3s;
      border-radius: 26px;
    }

    .toggle-slider:before {
      position: absolute;
      content: "";
      height: 20px;
      width: 20px;
      left: 3px;
      bottom: 3px;
      background-color: white;
      transition: 0.3s;
      border-radius: 50%;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    input:checked + .toggle-slider {
      background-color: #4f46e5;
    }

    input:checked + .toggle-slider:before {
      transform: translateX(22px);
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
      transition: all 0.2s;
      border: none;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: #4338ca;
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

    @keyframes spin {
      to {
        transform: rotate(360deg);
      }
    }

    @media (max-width: 640px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class DesignationFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);
  private departmentService = inject(DepartmentService);

  @Input() designation: Designation | null = null;
  @Input() loading = false;
  @Output() save = new EventEmitter<CreateDesignationRequest>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;
  levelOptions: SelectOption[] = DESIGNATION_LEVELS;
  departmentOptions: SelectOption[] = [];

  ngOnInit(): void {
    this.initForm();
    this.loadDepartments();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['designation'] && !changes['designation'].firstChange) {
      this.initForm();
    }
  }

  private initForm(): void {
    this.form = this.fb.group({
      name: [
        this.designation?.name || '',
        [Validators.required, Validators.maxLength(100)],
      ],
      level: [
        this.designation?.level || 5,
        [Validators.required, Validators.min(1), Validators.max(10)],
      ],
      departmentId: [this.designation?.departmentId || ''],
      isActive: [this.designation?.isActive ?? true],
    });
  }

  private loadDepartments(): void {
    this.departmentService.getDepartmentsDropdown().subscribe({
      next: departments => {
        this.departmentOptions = departments.map((d: DepartmentDropdownItem) => ({
          value: d.id,
          label: `${d.name} (${d.code})`,
        }));
      },
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
    if (control.errors['min']) return 'Level must be at least 1';
    if (control.errors['max']) return 'Level must be at most 10';
    if (control.errors['maxlength']) {
      const max = control.errors['maxlength'].requiredLength;
      return `Must be at most ${max} characters`;
    }

    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      const value = this.form.value;
      const request: CreateDesignationRequest = {
        name: value.name,
        level: Number(value.level),
        departmentId: value.departmentId || undefined,
        isActive: value.isActive,
      };
      this.save.emit(request);
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
