/**
 * MSLS Public Application Status Check Component
 *
 * Allows parents/guardians to check their admission application status
 * using their application number and registered phone number.
 * This is a public page - no authentication required.
 */

import { Component, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, NgForm } from '@angular/forms';

import { ApplicationService } from '../../application.service';
import {
  StatusCheckRequest,
  StatusCheckResponse,
  APPLICATION_STAGE_CONFIG,
  ApplicationStage,
} from '../../application.model';

@Component({
  selector: 'app-public-status-check',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './public-status-check.html',
  styleUrl: './public-status-check.scss',
})
export class PublicStatusCheck {
  private readonly applicationService = inject(ApplicationService);

  // Form state
  applicationNumber = signal('');
  phone = signal('');

  // UI state
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  result = signal<StatusCheckResponse | null>(null);

  // Computed status config
  statusConfig = computed(() => {
    const status = this.result()?.status;
    if (!status) return null;
    return APPLICATION_STAGE_CONFIG[status as ApplicationStage];
  });

  /**
   * Handle form submission
   */
  onSubmit(form: NgForm): void {
    if (form.invalid) {
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set(null);
    this.result.set(null);

    const request: StatusCheckRequest = {
      applicationNumber: this.applicationNumber().trim(),
      phone: this.phone().trim(),
    };

    this.applicationService.checkStatus(request).subscribe({
      next: (response) => {
        this.result.set(response);
        this.isLoading.set(false);
      },
      error: (error) => {
        this.errorMessage.set(
          error.message || 'Unable to find application. Please verify your application number and phone number.'
        );
        this.isLoading.set(false);
      },
    });
  }

  /**
   * Reset the form to check another application
   */
  reset(): void {
    this.applicationNumber.set('');
    this.phone.set('');
    this.result.set(null);
    this.errorMessage.set(null);
  }

  /**
   * Format date for display
   */
  formatDate(dateString: string | undefined): string {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
}
