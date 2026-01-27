/**
 * MSLS Enrollment Form Component
 *
 * Form for creating or editing student enrollments with polished UI.
 */

import { Component, input, output, inject, OnInit, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

import { EnrollmentService } from '../../services/enrollment.service';
import { Enrollment, CreateEnrollmentRequest, UpdateEnrollmentRequest } from '../../models/enrollment.model';

@Component({
  selector: 'app-enrollment-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="enrollment-form">
      <!-- Form Header -->
      <div class="form-header">
        <div class="form-header-icon" [class.edit-mode]="enrollment()">
          <i class="fa-solid" [class.fa-graduation-cap]="!enrollment()" [class.fa-pen]="enrollment()"></i>
        </div>
        <div class="form-header-text">
          <h3>{{ enrollment() ? 'Edit Enrollment' : 'New Enrollment' }}</h3>
          <p>{{ enrollment() ? 'Update enrollment details' : 'Enroll student in an academic year' }}</p>
        </div>
      </div>

      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <!-- Academic Year -->
        <div class="form-group">
          <label for="academicYearId" class="form-label">
            <i class="fa-solid fa-calendar-alt"></i>
            Academic Year
            <span class="required">*</span>
          </label>
          <select
            id="academicYearId"
            formControlName="academicYearId"
            class="form-select"
            [class.is-invalid]="isInvalid('academicYearId')"
            [disabled]="!!enrollment()"
          >
            <option value="">Select academic year</option>
            @for (year of academicYears(); track year.id) {
              <option [value]="year.id">{{ year.name }}</option>
            }
          </select>
          @if (isInvalid('academicYearId')) {
            <span class="error-text">
              <i class="fa-solid fa-exclamation-circle"></i>
              Academic year is required
            </span>
          }
        </div>

        <!-- Two Column Layout -->
        <div class="form-row">
          <!-- Roll Number -->
          <div class="form-group">
            <label for="rollNumber" class="form-label">
              <i class="fa-solid fa-hashtag"></i>
              Roll Number
            </label>
            <input
              type="text"
              id="rollNumber"
              formControlName="rollNumber"
              placeholder="e.g., A001"
              maxlength="20"
              class="form-input"
            />
            <span class="hint-text">Unique identifier within the class</span>
          </div>

          <!-- Enrollment Date -->
          @if (!enrollment()) {
            <div class="form-group">
              <label for="enrollmentDate" class="form-label">
                <i class="fa-solid fa-calendar-plus"></i>
                Enrollment Date
              </label>
              <input
                type="date"
                id="enrollmentDate"
                formControlName="enrollmentDate"
                class="form-input"
              />
            </div>
          }
        </div>

        <!-- Class & Section Info Box -->
        <div class="info-box">
          <div class="info-box-icon">
            <i class="fa-solid fa-info-circle"></i>
          </div>
          <div class="info-box-content">
            <strong>Class & Section Assignment</strong>
            <p>Class and section assignments will be available once the academic structure is set up. You can update this enrollment later.</p>
          </div>
        </div>

        <!-- Notes -->
        <div class="form-group">
          <label for="notes" class="form-label">
            <i class="fa-solid fa-sticky-note"></i>
            Notes
            <span class="optional">(optional)</span>
          </label>
          <textarea
            id="notes"
            formControlName="notes"
            rows="3"
            placeholder="Any additional notes about this enrollment..."
            class="form-textarea"
          ></textarea>
        </div>

        <!-- Error Message -->
        @if (service.error()) {
          <div class="error-alert">
            <i class="fa-solid fa-exclamation-triangle"></i>
            <span>{{ service.error() }}</span>
          </div>
        }

        <!-- Form Actions -->
        <div class="form-actions">
          <button type="button" class="btn btn--secondary" (click)="onCancel()">
            <i class="fa-solid fa-times"></i>
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn--primary"
            [disabled]="form.invalid || service.loading()"
          >
            @if (service.loading()) {
              <i class="fa-solid fa-spinner fa-spin"></i>
              {{ enrollment() ? 'Updating...' : 'Creating...' }}
            } @else {
              <i class="fa-solid" [class.fa-plus]="!enrollment()" [class.fa-save]="enrollment()"></i>
              {{ enrollment() ? 'Update Enrollment' : 'Create Enrollment' }}
            }
          </button>
        </div>
      </form>
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .enrollment-form {
      padding: 0.5rem;
    }

    /* Form Header */
    .form-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding-bottom: 1.25rem;
      margin-bottom: 1.25rem;
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
      background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
      flex-shrink: 0;
    }

    .form-header-icon.edit-mode {
      background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
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

    /* Form Groups */
    .form-group {
      margin-bottom: 1.25rem;
    }

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    @media (max-width: 480px) {
      .form-row {
        grid-template-columns: 1fr;
      }
    }

    .form-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
      margin-bottom: 0.5rem;
    }

    .form-label i {
      color: #6366f1;
      font-size: 0.75rem;
    }

    .required {
      color: #ef4444;
    }

    .optional {
      font-weight: 400;
      color: #9ca3af;
      font-size: 0.75rem;
    }

    .form-input,
    .form-select,
    .form-textarea {
      width: 100%;
      padding: 0.625rem 0.875rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #0f172a;
      transition: all 0.15s;
    }

    .form-input:focus,
    .form-select:focus,
    .form-textarea:focus {
      outline: none;
      border-color: #8b5cf6;
      box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
    }

    .form-input:disabled,
    .form-select:disabled {
      background: #f8fafc;
      color: #94a3b8;
      cursor: not-allowed;
    }

    .form-input.is-invalid,
    .form-select.is-invalid {
      border-color: #ef4444;
    }

    .form-input.is-invalid:focus,
    .form-select.is-invalid:focus {
      box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
    }

    .form-textarea {
      resize: vertical;
      min-height: 80px;
    }

    .hint-text {
      display: block;
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.375rem;
    }

    .error-text {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
      color: #ef4444;
      margin-top: 0.375rem;
    }

    .error-text i {
      font-size: 0.625rem;
    }

    /* Info Box */
    .info-box {
      display: flex;
      gap: 0.75rem;
      padding: 1rem;
      background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
      border: 1px solid #bae6fd;
      border-radius: 0.75rem;
      margin-bottom: 1.25rem;
    }

    .info-box-icon {
      flex-shrink: 0;
      width: 2rem;
      height: 2rem;
      background: white;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .info-box-icon i {
      color: #0284c7;
      font-size: 0.875rem;
    }

    .info-box-content {
      flex: 1;
    }

    .info-box-content strong {
      display: block;
      font-size: 0.8125rem;
      color: #0c4a6e;
      margin-bottom: 0.25rem;
    }

    .info-box-content p {
      margin: 0;
      font-size: 0.8125rem;
      color: #0369a1;
      line-height: 1.4;
    }

    /* Error Alert */
    .error-alert {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.875rem 1rem;
      background: #fef2f2;
      border: 1px solid #fecaca;
      border-radius: 0.5rem;
      color: #dc2626;
      font-size: 0.875rem;
      margin-bottom: 1.25rem;
    }

    .error-alert i {
      flex-shrink: 0;
    }

    /* Form Actions */
    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1.25rem;
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
      background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
      color: white;
      border: none;
    }

    .btn--primary:hover:not(:disabled) {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
    }

    .btn--primary:disabled {
      opacity: 0.6;
      cursor: not-allowed;
      transform: none;
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

    /* Responsive */
    @media (max-width: 480px) {
      .form-header {
        flex-direction: column;
        text-align: center;
      }

      .form-actions {
        flex-direction: column;
      }

      .form-actions .btn {
        width: 100%;
        justify-content: center;
      }
    }
  `]
})
export class EnrollmentFormComponent implements OnInit {
  protected service = inject(EnrollmentService);
  private fb = inject(FormBuilder);

  /** Student ID */
  studentId = input.required<string>();

  /** Existing enrollment for editing (null for create) */
  enrollment = input<Enrollment | null>(null);

  /** Available academic years */
  academicYears = input<Array<{ id: string; name: string }>>([]);

  /** Event emitted when form is saved successfully */
  saved = output<Enrollment>();

  /** Event emitted when form is cancelled */
  cancelled = output<void>();

  form!: FormGroup;

  constructor() {
    // Initialize form with default values
    this.initForm();

    // Update form when enrollment changes (for edit mode)
    effect(() => {
      const enrollment = this.enrollment();
      if (enrollment) {
        this.populateForm(enrollment);
      } else {
        this.form.reset();
        // Set default enrollment date to today for new enrollments
        this.form.patchValue({
          enrollmentDate: new Date().toISOString().split('T')[0]
        });
      }
    });
  }

  ngOnInit(): void {
    // Form is initialized in constructor
  }

  private initForm(): void {
    this.form = this.fb.group({
      academicYearId: ['', Validators.required],
      classId: [''],
      sectionId: [''],
      rollNumber: ['', Validators.maxLength(20)],
      classTeacherId: [''],
      enrollmentDate: [new Date().toISOString().split('T')[0]],
      notes: [''],
    });
  }

  private populateForm(enrollment: Enrollment): void {
    this.form.patchValue({
      academicYearId: enrollment.academicYear?.id || '',
      classId: enrollment.classId || '',
      sectionId: enrollment.sectionId || '',
      rollNumber: enrollment.rollNumber || '',
      classTeacherId: enrollment.classTeacherId || '',
      notes: enrollment.notes || '',
    });
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

    if (enrollment) {
      // Update existing enrollment
      const request: UpdateEnrollmentRequest = {
        classId: formValue.classId || undefined,
        sectionId: formValue.sectionId || undefined,
        rollNumber: formValue.rollNumber || undefined,
        classTeacherId: formValue.classTeacherId || undefined,
        notes: formValue.notes || undefined,
      };

      this.service.updateEnrollment(this.studentId(), enrollment.id, request).subscribe({
        next: (updated) => this.saved.emit(updated),
        error: () => {}, // Error handled in service
      });
    } else {
      // Create new enrollment
      const request: CreateEnrollmentRequest = {
        academicYearId: formValue.academicYearId,
        classId: formValue.classId || undefined,
        sectionId: formValue.sectionId || undefined,
        rollNumber: formValue.rollNumber || undefined,
        classTeacherId: formValue.classTeacherId || undefined,
        enrollmentDate: formValue.enrollmentDate || undefined,
        notes: formValue.notes || undefined,
      };

      this.service.createEnrollment(this.studentId(), request).subscribe({
        next: (created) => this.saved.emit(created),
        error: () => {}, // Error handled in service
      });
    }
  }

  onCancel(): void {
    this.cancelled.emit();
  }
}
