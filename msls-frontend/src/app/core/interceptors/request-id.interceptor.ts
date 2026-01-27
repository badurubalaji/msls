/**
 * MSLS Request ID Interceptor
 *
 * Adds a unique X-Request-ID header to each request for request tracing
 * and debugging purposes.
 */

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';

/** Header name for request identification */
const REQUEST_ID_HEADER = 'X-Request-ID';

/**
 * Generate a UUID v4
 * Uses crypto.randomUUID if available, falls back to manual generation
 */
function generateUuid(): string {
  // Use native crypto.randomUUID if available
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID();
  }

  // Fallback UUID generation
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/**
 * Request ID Interceptor Function
 *
 * Automatically adds a unique X-Request-ID header to all outgoing requests.
 * This enables:
 * - Request tracing across services
 * - Correlation of logs
 * - Debugging and monitoring
 *
 * Usage:
 * Requests will automatically include:
 * X-Request-ID: {uuid}
 */
export const requestIdInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
  // Generate unique request ID
  const requestId = generateUuid();

  // Clone request with request ID header
  const modifiedRequest = request.clone({
    setHeaders: {
      [REQUEST_ID_HEADER]: requestId,
    },
  });

  return next(modifiedRequest);
};
