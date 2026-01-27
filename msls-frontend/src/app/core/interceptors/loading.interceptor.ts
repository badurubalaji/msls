/**
 * MSLS Loading Interceptor
 *
 * Manages global loading state by tracking HTTP request lifecycle.
 * Calls LoadingService.show() on request start and hide() on completion.
 */

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpEvent } from '@angular/common/http';
import { inject } from '@angular/core';
import { Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';

import { LoadingService } from '../services/loading.service';

/** URLs that should not trigger loading indicator */
const SILENT_URLS: string[] = [
  // Add URLs that should not show loading indicator
  // e.g., '/api/notifications/poll', '/api/heartbeat'
];

/**
 * Check if URL should skip loading indicator
 */
function isSilentUrl(url: string): boolean {
  return SILENT_URLS.some(silentUrl => url.includes(silentUrl));
}

/**
 * Loading Interceptor Function
 *
 * Automatically manages the global loading state for HTTP requests.
 * Increments the loading counter on request start and decrements
 * on completion (success or error).
 *
 * Features:
 * - Tracks multiple concurrent requests
 * - Supports silent URLs that don't trigger loading
 * - Properly handles errors
 */
export const loadingInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
  const loadingService = inject(LoadingService);

  // Skip loading indicator for silent URLs
  if (isSilentUrl(request.url)) {
    return next(request);
  }

  // Show loading indicator
  loadingService.show();

  return next(request).pipe(
    finalize(() => {
      // Hide loading indicator when request completes
      loadingService.hide();
    })
  );
};
