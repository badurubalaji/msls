import { Component, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastService, Toast, ToastVariant } from '../../services/toast.service';

/**
 * MslsToastComponent - Displays toast notifications in the top-right corner.
 *
 * Usage:
 * Add this component once in your app root:
 * <msls-toast></msls-toast>
 *
 * Then use ToastService to show toasts:
 * this.toastService.success('Done!');
 */
@Component({
  selector: 'msls-toast',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './toast.component.html',
  styleUrl: './toast.component.scss'
})
export class MslsToastComponent {
  private toastService = inject(ToastService);

  /** All active toasts */
  toasts = this.toastService.toasts$;

  /** Dismiss a toast */
  dismiss(id: string): void {
    this.toastService.dismiss(id);
  }

  /** Get icon for toast variant */
  getIcon(variant: ToastVariant): string {
    switch (variant) {
      case 'success':
        return 'M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z';
      case 'error':
        return 'M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z';
      case 'warning':
        return 'M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z';
      case 'info':
      default:
        return 'M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z';
    }
  }

  /** Get CSS class for toast variant */
  getVariantClass(variant: ToastVariant): string {
    return `toast--${variant}`;
  }

  /** Track function for ngFor */
  trackByToastId(index: number, toast: Toast): string {
    return toast.id;
  }
}
