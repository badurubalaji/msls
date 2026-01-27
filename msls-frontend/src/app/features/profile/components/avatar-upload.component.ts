/**
 * AvatarUploadComponent - Component for uploading and displaying user avatar.
 */
import {
  Component,
  input,
  output,
  signal,
  computed,
  ChangeDetectionStrategy,
  inject,
  ElementRef,
  viewChild,
} from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileService } from '../../../core/services/profile.service';
import { MslsAvatarComponent } from '../../../shared/components/avatar/avatar.component';
import { MslsButtonComponent } from '../../../shared/components/button/button.component';

/** Maximum file size in bytes (2MB) */
const MAX_FILE_SIZE = 2 * 1024 * 1024;

/** Allowed MIME types */
const ALLOWED_TYPES = ['image/jpeg', 'image/png'];

@Component({
  selector: 'msls-avatar-upload',
  standalone: true,
  imports: [CommonModule, MslsAvatarComponent, MslsButtonComponent],
  template: `
    <div class="avatar-upload relative">
      <!-- Avatar Display -->
      <div class="relative inline-block">
        <msls-avatar
          [src]="currentAvatarUrl()"
          [name]="userName()"
          size="xl"
        />

        <!-- Upload Overlay -->
        <button
          type="button"
          class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-40
                 rounded-full opacity-0 hover:opacity-100 transition-opacity cursor-pointer"
          (click)="triggerFileInput()"
          [disabled]="uploading()"
        >
          @if (uploading()) {
            <!-- Progress Indicator -->
            <div class="text-white text-center">
              <svg
                class="w-6 h-6 animate-spin mx-auto"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              <span class="text-xs mt-1 block">{{ uploadProgress() }}%</span>
            </div>
          } @else {
            <!-- Camera Icon -->
            <svg
              class="w-6 h-6 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          }
        </button>
      </div>

      <!-- Hidden File Input -->
      <input
        #fileInput
        type="file"
        accept="image/jpeg,image/png"
        class="hidden"
        (change)="onFileSelected($event)"
      />

      <!-- Upload Button (alternative) -->
      <div class="mt-3 text-center">
        <msls-button
          variant="ghost"
          size="sm"
          [loading]="uploading()"
          (click)="triggerFileInput()"
        >
          {{ uploading() ? 'Uploading...' : 'Change Photo' }}
        </msls-button>
      </div>

      <!-- Error Message -->
      @if (errorMessage()) {
        <p class="mt-2 text-xs text-red-600 text-center">{{ errorMessage() }}</p>
      }

      <!-- Help Text -->
      <p class="mt-1 text-xs text-secondary-500 text-center">
        JPG or PNG. Max 2MB.
      </p>
    </div>
  `,
  styles: [
    `
      .avatar-upload {
        display: flex;
        flex-direction: column;
        align-items: center;
      }
    `,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AvatarUploadComponent {
  /** Current avatar URL */
  readonly currentAvatarUrl = input<string>('');

  /** User name for fallback initials */
  readonly userName = input<string>('');

  /** Emitted when avatar is successfully uploaded */
  readonly avatarUploaded = output<string>();

  private readonly profileService = inject(ProfileService);

  /** Reference to file input element */
  private readonly fileInputRef = viewChild<ElementRef<HTMLInputElement>>('fileInput');

  /** Error message */
  readonly errorMessage = signal<string | null>(null);

  /** Upload state from service */
  readonly uploading = computed(() => this.profileService.uploadProgress().uploading);
  readonly uploadProgress = computed(() => this.profileService.uploadProgress().progress);

  /**
   * Trigger the hidden file input
   */
  triggerFileInput(): void {
    this.fileInputRef()?.nativeElement.click();
  }

  /**
   * Handle file selection
   */
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (!file) {
      return;
    }

    // Reset error
    this.errorMessage.set(null);

    // Validate file type
    if (!ALLOWED_TYPES.includes(file.type)) {
      this.errorMessage.set('Please select a JPG or PNG image.');
      input.value = '';
      return;
    }

    // Validate file size
    if (file.size > MAX_FILE_SIZE) {
      this.errorMessage.set('Image must be less than 2MB.');
      input.value = '';
      return;
    }

    // Upload the file
    this.profileService.uploadAvatar(file).subscribe({
      next: (response) => {
        this.avatarUploaded.emit(response.avatarUrl);
        input.value = '';
      },
      error: (err) => {
        this.errorMessage.set(err.message || 'Failed to upload avatar.');
        input.value = '';
      },
    });
  }
}
