/**
 * MSLS Loading Service
 *
 * Manages global loading state with a request counter pattern.
 * Uses Angular Signals for reactive state management.
 */

import { Injectable, computed, signal } from '@angular/core';

/**
 * LoadingService - Manages global loading state with request counting.
 *
 * The service uses a counter to track multiple concurrent requests,
 * ensuring the loading indicator remains visible until all requests complete.
 *
 * Usage:
 * constructor(private loadingService: LoadingService) {}
 *
 * // In templates:
 * @if (loadingService.isLoading()) {
 *   <msls-spinner />
 * }
 *
 * // Manual control:
 * this.loadingService.show();
 * await someAsyncOperation();
 * this.loadingService.hide();
 */
@Injectable({ providedIn: 'root' })
export class LoadingService {
  /** Internal signal tracking the number of active requests */
  private _activeRequests = signal<number>(0);

  /** Public readonly signal for the number of active requests */
  readonly activeRequests = this._activeRequests.asReadonly();

  /** Computed signal that is true when any requests are active */
  readonly isLoading = computed(() => this._activeRequests() > 0);

  /**
   * Increment the active request counter
   * Call this when starting an HTTP request or async operation
   */
  show(): void {
    this._activeRequests.update(count => count + 1);
  }

  /**
   * Decrement the active request counter
   * Call this when an HTTP request or async operation completes
   * The counter will never go below 0
   */
  hide(): void {
    this._activeRequests.update(count => Math.max(0, count - 1));
  }

  /**
   * Reset the counter to 0
   * Use this to force-clear the loading state
   */
  reset(): void {
    this._activeRequests.set(0);
  }

  /**
   * Execute an async operation with automatic loading state management
   * @param operation - The async operation to execute
   * @returns The result of the operation
   */
  async withLoading<T>(operation: () => Promise<T>): Promise<T> {
    this.show();
    try {
      return await operation();
    } finally {
      this.hide();
    }
  }
}
