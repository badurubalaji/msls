/**
 * ProfileComponent - Main profile page with tabs for different sections.
 *
 * This component provides a comprehensive profile management interface with:
 * - Profile overview and editing
 * - Avatar upload
 * - Password change
 * - Notification preferences
 */
import { Component, OnInit, ChangeDetectionStrategy, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { ProfileService } from '../../core/services/profile.service';
import { ToastService } from '../../shared/services/toast.service';
import { MslsCardComponent } from '../../shared/components/card/card.component';
import { MslsButtonComponent } from '../../shared/components/button/button.component';
import { MslsSpinnerComponent } from '../../shared/components/spinner/spinner.component';

import { ProfileEditComponent } from './components/profile-edit.component';
import { ChangePasswordComponent } from './components/change-password.component';
import { AvatarUploadComponent } from './components/avatar-upload.component';
import { NotificationPreferencesComponent } from './components/notification-preferences.component';

/** Tab definition */
type ProfileTab = 'overview' | 'security' | 'preferences';

@Component({
  selector: 'msls-profile',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MslsCardComponent,
    MslsButtonComponent,
    MslsSpinnerComponent,
    ProfileEditComponent,
    ChangePasswordComponent,
    AvatarUploadComponent,
    NotificationPreferencesComponent,
  ],
  template: `
    <div class="profile-page">
      <!-- Page Header -->
      <div class="profile-header mb-8 pb-6 border-b-2 border-primary-100">
        <h1 class="text-2xl font-bold bg-gradient-to-r from-primary-600 to-primary-800 bg-clip-text text-transparent">Profile Settings</h1>
        <p class="text-secondary-500 mt-2">Manage your account settings and preferences</p>
      </div>

      @if (profileService.loading() && !profileService.profile()) {
        <div class="flex justify-center items-center py-12">
          <msls-spinner size="lg" />
        </div>
      } @else if (profileService.profile(); as profile) {
        <div class="profile-content grid grid-cols-1 lg:grid-cols-4 gap-6">
          <!-- Left Sidebar - Profile Summary -->
          <div class="lg:col-span-1">
            <msls-card variant="elevated" padding="lg">
              <div class="flex flex-col items-center text-center">
                <!-- Avatar -->
                <msls-avatar-upload
                  [currentAvatarUrl]="profile.avatarUrl || ''"
                  [userName]="profileService.fullName() || ''"
                  (avatarUploaded)="onAvatarUploaded($event)"
                />

                <!-- Name and Email -->
                <h2 class="mt-4 text-lg font-semibold text-secondary-900">
                  {{ profileService.fullName() }}
                </h2>
                <p class="text-secondary-600 text-sm">{{ profile.email }}</p>

                <!-- Status Badge -->
                <span
                  class="mt-2 px-2.5 py-0.5 rounded-full text-xs font-medium"
                  [class]="getStatusClass(profile.status)"
                >
                  {{ profile.status | titlecase }}
                </span>

                <!-- Quick Stats -->
                <div class="mt-6 w-full border-t border-secondary-200 pt-4">
                  <div class="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <p class="text-secondary-500">Timezone</p>
                      <p class="font-medium text-secondary-900">{{ profile.timezone }}</p>
                    </div>
                    <div>
                      <p class="text-secondary-500">Language</p>
                      <p class="font-medium text-secondary-900">{{ profile.locale }}</p>
                    </div>
                  </div>
                </div>

                @if (profile.lastLoginAt) {
                  <div class="mt-4 w-full border-t border-secondary-200 pt-4">
                    <p class="text-xs text-secondary-500">
                      Last login: {{ profile.lastLoginAt | date : 'medium' }}
                    </p>
                  </div>
                }
              </div>
            </msls-card>
          </div>

          <!-- Right Content - Tabs -->
          <div class="lg:col-span-3">
            <!-- Tab Navigation -->
            <div class="mb-6 bg-secondary-50 rounded-xl p-1 border border-secondary-200">
              <nav class="flex space-x-1" aria-label="Profile tabs">
                <button
                  type="button"
                  (click)="activeTab.set('overview')"
                  class="flex-1 py-3 px-4 text-sm font-medium rounded-lg transition-all duration-200"
                  [class]="
                    activeTab() === 'overview'
                      ? 'bg-white text-primary-700 shadow-md shadow-primary-100/50'
                      : 'text-secondary-600 hover:text-secondary-900 hover:bg-white/50'
                  "
                >
                  Profile
                </button>
                <button
                  type="button"
                  (click)="activeTab.set('security')"
                  class="flex-1 py-3 px-4 text-sm font-medium rounded-lg transition-all duration-200"
                  [class]="
                    activeTab() === 'security'
                      ? 'bg-white text-primary-700 shadow-md shadow-primary-100/50'
                      : 'text-secondary-600 hover:text-secondary-900 hover:bg-white/50'
                  "
                >
                  Security
                </button>
                <button
                  type="button"
                  (click)="activeTab.set('preferences')"
                  class="flex-1 py-3 px-4 text-sm font-medium rounded-lg transition-all duration-200"
                  [class]="
                    activeTab() === 'preferences'
                      ? 'bg-white text-primary-700 shadow-md shadow-primary-100/50'
                      : 'text-secondary-600 hover:text-secondary-900 hover:bg-white/50'
                  "
                >
                  Preferences
                </button>
              </nav>
            </div>

            <!-- Tab Content -->
            @switch (activeTab()) {
              @case ('overview') {
                <msls-profile-edit [profile]="profile" (profileUpdated)="onProfileUpdated()" />
              }
              @case ('security') {
                <div class="space-y-6">
                  <msls-change-password (passwordChanged)="onPasswordChanged()" />

                  <!-- 2FA Section -->
                  <msls-card variant="outlined" padding="lg">
                    <div class="flex items-start justify-between">
                      <div>
                        <h3 class="text-lg font-medium text-secondary-900">
                          Two-Factor Authentication
                        </h3>
                        <p class="mt-1 text-sm text-secondary-600">
                          Add an extra layer of security to your account
                        </p>
                      </div>
                      <span
                        class="px-2.5 py-0.5 rounded-full text-xs font-medium"
                        [class]="
                          profile.twoFactorEnabled
                            ? 'bg-green-100 text-green-800'
                            : 'bg-secondary-100 text-secondary-800'
                        "
                      >
                        {{ profile.twoFactorEnabled ? 'Enabled' : 'Disabled' }}
                      </span>
                    </div>
                  </msls-card>

                  <!-- Account Deletion -->
                  <msls-card variant="outlined" padding="lg">
                    <div class="flex items-start justify-between">
                      <div>
                        <h3 class="text-lg font-medium text-red-600">Delete Account</h3>
                        <p class="mt-1 text-sm text-secondary-600">
                          Permanently delete your account and all associated data
                        </p>
                      </div>
                      <msls-button variant="danger" size="sm" (click)="confirmAccountDeletion()">
                        Delete Account
                      </msls-button>
                    </div>
                  </msls-card>
                </div>
              }
              @case ('preferences') {
                <msls-notification-preferences
                  [preferences]="profile.notificationPreferences"
                  (preferencesUpdated)="onPreferencesUpdated()"
                />
              }
            }
          </div>
        </div>
      } @else {
        <msls-card variant="outlined" padding="lg">
          <div class="text-center py-8">
            <p class="text-secondary-600">Failed to load profile. Please try again.</p>
            <msls-button variant="primary" class="mt-4" (click)="loadProfile()">
              Retry
            </msls-button>
          </div>
        </msls-card>
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteConfirmation()) {
        <div
          class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
          (click)="showDeleteConfirmation.set(false)"
        >
          <div
            class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6"
            (click)="$event.stopPropagation()"
          >
            <h3 class="text-lg font-semibold text-secondary-900">Delete Account</h3>
            <p class="mt-2 text-sm text-secondary-600">
              Are you sure you want to delete your account? This action cannot be undone. All your
              data will be permanently removed.
            </p>
            <div class="mt-6 flex justify-end space-x-3">
              <msls-button variant="secondary" (click)="showDeleteConfirmation.set(false)">
                Cancel
              </msls-button>
              <msls-button
                variant="danger"
                [loading]="deletingAccount()"
                (click)="deleteAccount()"
              >
                Delete Account
              </msls-button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: [
    `
      .profile-page {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1.5rem;
      }
    `,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProfileComponent implements OnInit {
  readonly profileService = inject(ProfileService);
  private readonly toastService = inject(ToastService);

  /** Current active tab */
  readonly activeTab = signal<ProfileTab>('overview');

  /** Delete confirmation dialog state */
  readonly showDeleteConfirmation = signal(false);

  /** Account deletion in progress */
  readonly deletingAccount = signal(false);

  ngOnInit(): void {
    this.loadProfile();
  }

  /**
   * Load the user profile
   */
  loadProfile(): void {
    this.profileService.loadProfile().subscribe({
      error: (err) => {
        this.toastService.error('Failed to load profile: ' + err.message);
      },
    });
  }

  /**
   * Handle avatar upload completion
   */
  onAvatarUploaded(avatarUrl: string): void {
    this.toastService.success('Avatar updated successfully');
  }

  /**
   * Handle profile update completion
   */
  onProfileUpdated(): void {
    this.toastService.success('Profile updated successfully');
  }

  /**
   * Handle password change completion
   */
  onPasswordChanged(): void {
    this.toastService.success('Password changed successfully');
  }

  /**
   * Handle preferences update completion
   */
  onPreferencesUpdated(): void {
    this.toastService.success('Preferences updated successfully');
  }

  /**
   * Show account deletion confirmation
   */
  confirmAccountDeletion(): void {
    this.showDeleteConfirmation.set(true);
  }

  /**
   * Delete the account
   */
  deleteAccount(): void {
    this.deletingAccount.set(true);

    this.profileService.requestAccountDeletion().subscribe({
      next: () => {
        this.toastService.success('Account deletion requested. You will be logged out shortly.');
        this.showDeleteConfirmation.set(false);
        // In a real app, you would log out the user here
      },
      error: (err) => {
        this.toastService.error('Failed to delete account: ' + err.message);
        this.deletingAccount.set(false);
      },
    });
  }

  /**
   * Get status badge CSS classes
   */
  getStatusClass(status: string): string {
    switch (status) {
      case 'active':
        return 'bg-gradient-to-r from-success-100 to-success-200 text-success-800 ring-1 ring-success-500/20';
      case 'inactive':
        return 'bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-700 ring-1 ring-secondary-500/20';
      case 'pending':
        return 'bg-gradient-to-r from-warning-100 to-warning-200 text-warning-800 ring-1 ring-warning-500/20';
      case 'suspended':
        return 'bg-gradient-to-r from-danger-100 to-danger-200 text-danger-800 ring-1 ring-danger-500/20';
      default:
        return 'bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-700 ring-1 ring-secondary-500/20';
    }
  }
}
