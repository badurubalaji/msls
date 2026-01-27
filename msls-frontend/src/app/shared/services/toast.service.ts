import { Injectable, signal } from '@angular/core';

/** Toast variant types */
export type ToastVariant = 'success' | 'error' | 'warning' | 'info';

/** Toast message interface */
export interface Toast {
  id: string;
  message: string;
  variant: ToastVariant;
  duration: number;
  dismissible: boolean;
}

/** Configuration for creating a toast */
export interface ToastConfig {
  message: string;
  variant?: ToastVariant;
  duration?: number;
  dismissible?: boolean;
}

/** Default toast duration in milliseconds */
const DEFAULT_DURATION = 5000;

/**
 * ToastService - Manages toast notifications with signal-based state.
 *
 * Usage:
 * constructor(private toastService: ToastService) {}
 *
 * showSuccess() {
 *   this.toastService.success('Operation completed successfully!');
 * }
 *
 * showError() {
 *   this.toastService.error('Something went wrong!', 10000);
 * }
 */
@Injectable({ providedIn: 'root' })
export class ToastService {
  /** Internal signal holding all active toasts */
  private toasts = signal<Toast[]>([]);

  /** Public readonly accessor for toasts */
  readonly toasts$ = this.toasts.asReadonly();

  /** Map of toast IDs to their timeout handles */
  private timeouts = new Map<string, ReturnType<typeof setTimeout>>();

  /**
   * Show a success toast
   * @param message - The message to display
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  success(message: string, duration: number = DEFAULT_DURATION): string {
    return this.show({ message, variant: 'success', duration });
  }

  /**
   * Show an error toast
   * @param message - The message to display
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  error(message: string, duration: number = DEFAULT_DURATION): string {
    return this.show({ message, variant: 'error', duration });
  }

  /**
   * Show a warning toast
   * @param message - The message to display
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  warning(message: string, duration: number = DEFAULT_DURATION): string {
    return this.show({ message, variant: 'warning', duration });
  }

  /**
   * Show an info toast
   * @param message - The message to display
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  info(message: string, duration: number = DEFAULT_DURATION): string {
    return this.show({ message, variant: 'info', duration });
  }

  /**
   * Show a toast with custom configuration
   * @param config - Toast configuration
   * @returns The ID of the created toast
   */
  show(config: ToastConfig): string {
    const toast: Toast = {
      id: this.generateId(),
      message: config.message,
      variant: config.variant ?? 'info',
      duration: config.duration ?? DEFAULT_DURATION,
      dismissible: config.dismissible ?? true,
    };

    // Add toast to the list
    this.toasts.update(toasts => [...toasts, toast]);

    // Set up auto-dismiss if duration > 0
    if (toast.duration > 0) {
      const timeout = setTimeout(() => {
        this.dismiss(toast.id);
      }, toast.duration);
      this.timeouts.set(toast.id, timeout);
    }

    return toast.id;
  }

  /**
   * Dismiss a specific toast by ID
   * @param id - The ID of the toast to dismiss
   */
  dismiss(id: string): void {
    // Clear the timeout if it exists
    const timeout = this.timeouts.get(id);
    if (timeout) {
      clearTimeout(timeout);
      this.timeouts.delete(id);
    }

    // Remove the toast from the list
    this.toasts.update(toasts => toasts.filter(t => t.id !== id));
  }

  /**
   * Dismiss all toasts
   */
  dismissAll(): void {
    // Clear all timeouts
    this.timeouts.forEach(timeout => clearTimeout(timeout));
    this.timeouts.clear();

    // Clear all toasts
    this.toasts.set([]);
  }

  /**
   * Generate a unique ID for a toast
   */
  private generateId(): string {
    return `toast-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
  }
}
