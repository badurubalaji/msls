/**
 * ProfileEditComponent - Form for editing user profile information.
 */
import {
  Component,
  input,
  output,
  signal,
  ChangeDetectionStrategy,
  inject,
  OnInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';

import { ProfileService } from '../../../core/services/profile.service';
import { UserProfile, AVAILABLE_TIMEZONES, AVAILABLE_LOCALES } from '../../../core/models/profile.model';
import { MslsCardComponent } from '../../../shared/components/card/card.component';
import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsInputComponent } from '../../../shared/components/input/input.component';
import { MslsSelectComponent, SelectOption } from '../../../shared/components/select/select.component';

@Component({
  selector: 'msls-profile-edit',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsCardComponent,
    MslsButtonComponent,
    MslsInputComponent,
    MslsSelectComponent,
  ],
  template: `
    <msls-card variant="outlined" padding="lg">
      <h3 class="text-lg font-medium text-secondary-900 mb-4">Personal Information</h3>

      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- First Name -->
          <msls-input
            formControlName="firstName"
            label="First Name"
            placeholder="Enter your first name"
            [error]="getError('firstName')"
          />

          <!-- Last Name -->
          <msls-input
            formControlName="lastName"
            label="Last Name"
            placeholder="Enter your last name"
            [error]="getError('lastName')"
          />

          <!-- Email (readonly) -->
          <div>
            <label class="block text-sm font-medium text-secondary-700 mb-1">Email</label>
            <input
              type="email"
              [value]="profile().email || ''"
              readonly
              class="w-full px-3 py-2 border border-secondary-300 rounded-md shadow-sm
                     bg-secondary-50 text-secondary-600 cursor-not-allowed"
            />
            <p class="mt-1 text-xs text-secondary-500">Email cannot be changed</p>
          </div>

          <!-- Phone -->
          <msls-input
            formControlName="phone"
            label="Phone Number"
            type="tel"
            placeholder="Enter your phone number"
            [error]="getError('phone')"
          />
        </div>

        <!-- Bio -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-secondary-700 mb-1">Bio</label>
          <textarea
            formControlName="bio"
            class="w-full px-3 py-2 border border-secondary-300 rounded-md shadow-sm
                   focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500
                   resize-none"
            rows="3"
            placeholder="Tell us about yourself..."
          ></textarea>
        </div>

        <!-- Timezone and Locale -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
          <div>
            <label class="block text-sm font-medium text-secondary-700 mb-1">Timezone</label>
            <msls-select
              formControlName="timezone"
              [options]="timezoneOptions"
              placeholder="Select timezone"
              [searchable]="true"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-secondary-700 mb-1">Language</label>
            <msls-select
              formControlName="locale"
              [options]="localeOptions"
              placeholder="Select language"
              [searchable]="true"
            />
          </div>
        </div>

        <!-- Submit Button -->
        <div class="mt-6 flex justify-end">
          <msls-button
            type="submit"
            variant="primary"
            [loading]="saving()"
            [disabled]="!form.valid || !form.dirty"
          >
            Save Changes
          </msls-button>
        </div>
      </form>
    </msls-card>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProfileEditComponent implements OnInit {
  /** Current profile data */
  readonly profile = input.required<UserProfile>();

  /** Emitted when profile is successfully updated */
  readonly profileUpdated = output<void>();

  private readonly fb = inject(FormBuilder);
  private readonly profileService = inject(ProfileService);

  /** Form group for profile editing */
  form!: FormGroup;

  /** Saving state */
  readonly saving = signal(false);

  /** Timezone options for select */
  readonly timezoneOptions: SelectOption[] = AVAILABLE_TIMEZONES.map((tz) => ({
    value: tz.value,
    label: tz.label,
  }));

  /** Locale options for select */
  readonly localeOptions: SelectOption[] = AVAILABLE_LOCALES.map((locale) => ({
    value: locale.value,
    label: locale.label,
  }));

  ngOnInit(): void {
    this.initForm();
  }

  /**
   * Initialize the form with profile data
   */
  private initForm(): void {
    const profile = this.profile();

    this.form = this.fb.group({
      firstName: [profile.firstName, [Validators.required, Validators.maxLength(100)]],
      lastName: [profile.lastName, [Validators.required, Validators.maxLength(100)]],
      phone: [profile.phone || '', [Validators.maxLength(20)]],
      bio: [profile.bio || '', [Validators.maxLength(500)]],
      timezone: [profile.timezone || 'UTC'],
      locale: [profile.locale || 'en'],
    });
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
      return `${this.getFieldLabel(fieldName)} is required`;
    }
    if (control.errors['maxlength']) {
      return `${this.getFieldLabel(fieldName)} is too long`;
    }

    return '';
  }

  /**
   * Get human-readable field label
   */
  private getFieldLabel(fieldName: string): string {
    const labels: Record<string, string> = {
      firstName: 'First name',
      lastName: 'Last name',
      phone: 'Phone number',
      bio: 'Bio',
    };
    return labels[fieldName] || fieldName;
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    if (!this.form.valid || !this.form.dirty) {
      return;
    }

    this.saving.set(true);

    const formValue = this.form.value;
    const updateData = {
      firstName: formValue.firstName,
      lastName: formValue.lastName,
      phone: formValue.phone || undefined,
      bio: formValue.bio || undefined,
      timezone: formValue.timezone,
      locale: formValue.locale,
    };

    this.profileService.updateProfile(updateData).subscribe({
      next: () => {
        this.saving.set(false);
        this.form.markAsPristine();
        this.profileUpdated.emit();
      },
      error: () => {
        this.saving.set(false);
      },
    });
  }
}
