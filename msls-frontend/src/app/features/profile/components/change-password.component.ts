/**
 * ChangePasswordComponent - Form for changing user password.
 */
import {
  Component,
  output,
  signal,
  ChangeDetectionStrategy,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
  AbstractControl,
  ValidationErrors,
} from '@angular/forms';

import { ProfileService } from '../../../core/services/profile.service';
import { MslsCardComponent } from '../../../shared/components/card/card.component';
import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsInputComponent } from '../../../shared/components/input/input.component';

/**
 * Custom validator to check if passwords match
 */
function passwordsMatchValidator(control: AbstractControl): ValidationErrors | null {
  const newPassword = control.get('newPassword');
  const confirmPassword = control.get('confirmPassword');

  if (newPassword && confirmPassword && newPassword.value !== confirmPassword.value) {
    confirmPassword.setErrors({ passwordMismatch: true });
    return { passwordMismatch: true };
  }

  return null;
}

/**
 * Password strength validator
 */
function passwordStrengthValidator(control: AbstractControl): ValidationErrors | null {
  const value = control.value;

  if (!value) {
    return null;
  }

  const hasUppercase = /[A-Z]/.test(value);
  const hasLowercase = /[a-z]/.test(value);
  const hasDigit = /\d/.test(value);
  const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(value);
  const isLongEnough = value.length >= 8;

  const errors: ValidationErrors = {};

  if (!isLongEnough) {
    errors['minLength'] = true;
  }
  if (!hasUppercase) {
    errors['noUppercase'] = true;
  }
  if (!hasLowercase) {
    errors['noLowercase'] = true;
  }
  if (!hasDigit) {
    errors['noDigit'] = true;
  }
  if (!hasSpecial) {
    errors['noSpecial'] = true;
  }

  return Object.keys(errors).length > 0 ? errors : null;
}

@Component({
  selector: 'msls-change-password',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsCardComponent,
    MslsButtonComponent,
    MslsInputComponent,
  ],
  template: `
    <msls-card variant="outlined" padding="lg">
      <h3 class="text-lg font-medium text-secondary-900 mb-4">Change Password</h3>

      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <div class="space-y-4">
          <!-- Current Password -->
          <msls-input
            formControlName="currentPassword"
            type="password"
            label="Current Password"
            placeholder="Enter your current password"
            autocomplete="current-password"
            [error]="getError('currentPassword')"
          />

          <!-- New Password -->
          <div>
            <msls-input
              formControlName="newPassword"
              type="password"
              label="New Password"
              placeholder="Enter your new password"
              autocomplete="new-password"
              [error]="getError('newPassword')"
            />

            <!-- Password Requirements -->
            @if (form.get('newPassword')?.touched) {
              <div class="mt-2 text-xs space-y-1">
                <p [class]="getRequirementClass('minLength')">
                  At least 8 characters
                </p>
                <p [class]="getRequirementClass('noUppercase')">
                  At least one uppercase letter
                </p>
                <p [class]="getRequirementClass('noLowercase')">
                  At least one lowercase letter
                </p>
                <p [class]="getRequirementClass('noDigit')">
                  At least one number
                </p>
                <p [class]="getRequirementClass('noSpecial')">
                  At least one special character (!&#64;#$%^&*...)
                </p>
              </div>
            }
          </div>

          <!-- Confirm Password -->
          <msls-input
            formControlName="confirmPassword"
            type="password"
            label="Confirm New Password"
            placeholder="Confirm your new password"
            autocomplete="new-password"
            [error]="getError('confirmPassword')"
          />
        </div>

        <!-- Error Message -->
        @if (errorMessage()) {
          <div class="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p class="text-sm text-red-600">{{ errorMessage() }}</p>
          </div>
        }

        <!-- Submit Button -->
        <div class="mt-6 flex justify-end">
          <msls-button
            type="submit"
            variant="primary"
            [loading]="saving()"
            [disabled]="!form.valid"
          >
            Change Password
          </msls-button>
        </div>
      </form>
    </msls-card>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ChangePasswordComponent {
  /** Emitted when password is successfully changed */
  readonly passwordChanged = output<void>();

  private readonly fb = inject(FormBuilder);
  private readonly profileService = inject(ProfileService);

  /** Form group for password change */
  form: FormGroup;

  /** Saving state */
  readonly saving = signal(false);

  /** Error message from API */
  readonly errorMessage = signal<string | null>(null);

  constructor() {
    this.form = this.fb.group(
      {
        currentPassword: ['', [Validators.required]],
        newPassword: ['', [Validators.required, passwordStrengthValidator]],
        confirmPassword: ['', [Validators.required]],
      },
      { validators: passwordsMatchValidator }
    );
  }

  /**
   * Get error message for a form field
   */
  getError(fieldName: string): string {
    const control = this.form.get(fieldName);
    if (!control || !control.touched || !control.errors) {
      return '';
    }

    if (control.errors['required']) {
      return 'This field is required';
    }
    if (control.errors['passwordMismatch']) {
      return 'Passwords do not match';
    }

    // For newPassword, we show requirements separately
    if (fieldName === 'newPassword' && Object.keys(control.errors).length > 0) {
      return 'Password does not meet requirements';
    }

    return '';
  }

  /**
   * Get CSS class for password requirement indicator
   */
  getRequirementClass(requirement: string): string {
    const control = this.form.get('newPassword');
    if (!control?.value) {
      return 'text-secondary-400';
    }

    // Check if this specific requirement is met
    const errors = control.errors || {};
    const isMet = !errors[requirement];

    return isMet ? 'text-green-600' : 'text-red-600';
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    if (!this.form.valid) {
      return;
    }

    this.saving.set(true);
    this.errorMessage.set(null);

    const formValue = this.form.value;

    this.profileService
      .changePassword({
        currentPassword: formValue.currentPassword,
        newPassword: formValue.newPassword,
        confirmPassword: formValue.confirmPassword,
      })
      .subscribe({
        next: () => {
          this.saving.set(false);
          this.form.reset();
          this.passwordChanged.emit();
        },
        error: (err) => {
          this.saving.set(false);
          this.errorMessage.set(err.message || 'Failed to change password');
        },
      });
  }
}
