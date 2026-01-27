/**
 * NotificationPreferencesComponent - Component for managing notification preferences.
 */
import {
  Component,
  input,
  output,
  signal,
  ChangeDetectionStrategy,
  inject,
  OnInit,
  effect,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { ProfileService } from '../../../core/services/profile.service';
import { NotificationPreferences } from '../../../core/models/profile.model';
import { MslsCardComponent } from '../../../shared/components/card/card.component';
import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsCheckboxComponent } from '../../../shared/components/checkbox/checkbox.component';

@Component({
  selector: 'msls-notification-preferences',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsCardComponent,
    MslsButtonComponent,
    MslsCheckboxComponent,
  ],
  template: `
    <msls-card variant="outlined" padding="lg">
      <h3 class="text-lg font-medium text-secondary-900 mb-4">Notification Preferences</h3>
      <p class="text-sm text-secondary-600 mb-6">
        Choose how you want to receive notifications about your account activity.
      </p>

      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <div class="space-y-4">
          <!-- Email Notifications -->
          <div
            class="flex items-start justify-between p-4 rounded-lg border border-secondary-200 hover:border-secondary-300 transition-colors"
          >
            <div>
              <h4 class="text-sm font-medium text-secondary-900">Email Notifications</h4>
              <p class="text-sm text-secondary-600 mt-1">
                Receive notifications via email about important updates, reminders, and announcements.
              </p>
            </div>
            <msls-checkbox formControlName="email" />
          </div>

          <!-- Push Notifications -->
          <div
            class="flex items-start justify-between p-4 rounded-lg border border-secondary-200 hover:border-secondary-300 transition-colors"
          >
            <div>
              <h4 class="text-sm font-medium text-secondary-900">Push Notifications</h4>
              <p class="text-sm text-secondary-600 mt-1">
                Receive push notifications in your browser for real-time alerts and updates.
              </p>
            </div>
            <msls-checkbox formControlName="push" />
          </div>

          <!-- SMS Notifications -->
          <div
            class="flex items-start justify-between p-4 rounded-lg border border-secondary-200 hover:border-secondary-300 transition-colors"
          >
            <div>
              <h4 class="text-sm font-medium text-secondary-900">SMS Notifications</h4>
              <p class="text-sm text-secondary-600 mt-1">
                Receive text messages for urgent notifications and security alerts.
              </p>
            </div>
            <msls-checkbox formControlName="sms" />
          </div>
        </div>

        <!-- Submit Button -->
        <div class="mt-6 flex justify-end">
          <msls-button
            type="submit"
            variant="primary"
            [loading]="saving()"
            [disabled]="!form.dirty"
          >
            Save Preferences
          </msls-button>
        </div>
      </form>
    </msls-card>

    <!-- Additional Preference Categories -->
    <msls-card variant="outlined" padding="lg" class="mt-6">
      <h3 class="text-lg font-medium text-secondary-900 mb-4">Notification Categories</h3>
      <p class="text-sm text-secondary-600 mb-6">
        Fine-tune which types of notifications you want to receive.
      </p>

      <div class="space-y-3">
        @for (category of notificationCategories; track category.key) {
          <div
            class="flex items-center justify-between py-2 border-b border-secondary-100 last:border-0"
          >
            <div>
              <h4 class="text-sm font-medium text-secondary-900">{{ category.label }}</h4>
              <p class="text-xs text-secondary-500">{{ category.description }}</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                [checked]="category.enabled"
                (change)="toggleCategory(category.key)"
                class="sr-only peer"
              />
              <div
                class="w-11 h-6 bg-secondary-200 peer-focus:outline-none peer-focus:ring-4
                       peer-focus:ring-primary-300 rounded-full peer
                       peer-checked:after:translate-x-full peer-checked:after:border-white
                       after:content-[''] after:absolute after:top-[2px] after:left-[2px]
                       after:bg-white after:border-secondary-300 after:border after:rounded-full
                       after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"
              ></div>
            </label>
          </div>
        }
      </div>
    </msls-card>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NotificationPreferencesComponent implements OnInit {
  /** Current notification preferences */
  readonly preferences = input<NotificationPreferences | undefined>();

  /** Emitted when preferences are successfully updated */
  readonly preferencesUpdated = output<void>();

  private readonly fb = inject(FormBuilder);
  private readonly profileService = inject(ProfileService);

  /** Form group for preferences */
  form!: FormGroup;

  /** Saving state */
  readonly saving = signal(false);

  /** Notification categories for fine-grained control */
  readonly notificationCategories = [
    {
      key: 'account_activity',
      label: 'Account Activity',
      description: 'Login alerts, password changes, and security updates',
      enabled: true,
    },
    {
      key: 'course_updates',
      label: 'Course Updates',
      description: 'New assignments, grades, and course announcements',
      enabled: true,
    },
    {
      key: 'schedule_reminders',
      label: 'Schedule Reminders',
      description: 'Upcoming classes, events, and deadlines',
      enabled: true,
    },
    {
      key: 'marketing',
      label: 'Marketing & Promotions',
      description: 'News, tips, and promotional offers',
      enabled: false,
    },
  ];

  constructor() {
    // Initialize form with effect to update when preferences change
    effect(() => {
      const prefs = this.preferences();
      if (prefs && this.form) {
        this.form.patchValue({
          email: prefs.email,
          push: prefs.push,
          sms: prefs.sms,
        }, { emitEvent: false });
        this.form.markAsPristine();
      }
    });
  }

  ngOnInit(): void {
    const prefs = this.preferences() || { email: true, push: true, sms: false };

    this.form = this.fb.group({
      email: [prefs.email],
      push: [prefs.push],
      sms: [prefs.sms],
    });
  }

  /**
   * Toggle a notification category
   */
  toggleCategory(key: string): void {
    const category = this.notificationCategories.find((c) => c.key === key);
    if (category) {
      category.enabled = !category.enabled;
      // In a real app, you would save this to extended preferences
      this.profileService
        .setExtendedPreference({
          category: 'notifications',
          key,
          value: category.enabled,
        })
        .subscribe();
    }
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    if (!this.form.dirty) {
      return;
    }

    this.saving.set(true);

    const formValue = this.form.value;

    this.profileService
      .updateNotificationPreferences({
        email: formValue.email,
        push: formValue.push,
        sms: formValue.sms,
      })
      .subscribe({
        next: () => {
          this.saving.set(false);
          this.form.markAsPristine();
          this.preferencesUpdated.emit();
        },
        error: () => {
          this.saving.set(false);
        },
      });
  }
}
