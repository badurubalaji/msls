/**
 * MSLS Role Form Component
 *
 * Form component for creating and editing roles.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import {
  MslsInputComponent,
  MslsFormFieldComponent,
} from '../../../shared/components';
import { Role } from '../../../core/models';

@Component({
  selector: 'msls-role-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsInputComponent,
    MslsFormFieldComponent,
  ],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="role-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="role">
          <i class="fa-solid" [class.fa-user-shield]="!role" [class.fa-pen]="role"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ role ? 'Edit Role' : 'Create New Role' }}</h3>
          <p>{{ role ? 'Update role details' : 'Define a new role for your organization' }}</p>
        </div>
      </div>

      <msls-form-field
        label="Role Name"
        [required]="true"
        [error]="getFieldError('name') || ''"
      >
        <msls-input
          type="text"
          formControlName="name"
          placeholder="Enter role name"
        />
      </msls-form-field>

      <msls-form-field
        label="Description"
        [error]="getFieldError('description') || ''"
      >
        <textarea
          class="form-textarea"
          formControlName="description"
          placeholder="Enter role description"
          rows="3"
        ></textarea>
      </msls-form-field>

      <div class="form-actions">
        <button type="button" class="btn btn--secondary" (click)="onCancel()">
          <i class="fa-solid fa-times"></i>
          Cancel
        </button>
        <button type="submit" class="btn btn--primary" [disabled]="!form.valid || loading">
          <i *ngIf="loading" class="fa-solid fa-spinner fa-spin"></i>
          <i *ngIf="!loading && !role" class="fa-solid fa-plus"></i>
          <i *ngIf="!loading && role" class="fa-solid fa-save"></i>
          {{ role ? 'Update Role' : 'Create Role' }}
        </button>
      </div>
    </form>
  `,
  styles: [`
    .role-form {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-lg);
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

    .form-textarea {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      resize: vertical;
      min-height: 80px;
      font-family: inherit;
      color: #0f172a;
      background: white;
      transition: all 0.15s;
    }

    .form-textarea:focus {
      outline: none;
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
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
  `],
})
export class RoleFormComponent implements OnInit, OnChanges {
  private fb = inject(FormBuilder);

  @Input() role: Role | null = null;
  @Input() loading = false;

  @Output() save = new EventEmitter<{
    name: string;
    description?: string;
  }>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  ngOnInit(): void {
    this.initForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Reinitialize form when role input changes
    // This ensures form resets when opening for a new role after editing one
    if (changes['role']) {
      this.initForm();
    }
  }

  private initForm(): void {
    this.form = this.fb.group({
      name: [
        this.role?.name || '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(100),
        ],
      ],
      description: [
        this.role?.description || '',
        [Validators.maxLength(500)],
      ],
    });
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
        description: value.description || undefined,
      });
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
