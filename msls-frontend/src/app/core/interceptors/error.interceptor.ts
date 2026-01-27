/**
 * MSLS Error Interceptor
 *
 * Transforms HTTP errors into RFC 7807 ApiError format for consistent
 * error handling throughout the application.
 */

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpErrorResponse, HttpEvent } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { ApiError } from '../models/api-response.model';

/**
 * HTTP status code to error title mapping
 */
const STATUS_TITLES: Record<number, string> = {
  400: 'Bad Request',
  401: 'Unauthorized',
  403: 'Forbidden',
  404: 'Not Found',
  405: 'Method Not Allowed',
  408: 'Request Timeout',
  409: 'Conflict',
  422: 'Unprocessable Entity',
  429: 'Too Many Requests',
  500: 'Internal Server Error',
  502: 'Bad Gateway',
  503: 'Service Unavailable',
  504: 'Gateway Timeout',
};

/**
 * Transform HTTP error into RFC 7807 ApiError format
 */
function transformToApiError(error: HttpErrorResponse): ApiError {
  // If the error already contains an RFC 7807 structure, use it
  if (error.error && typeof error.error === 'object' && 'title' in error.error) {
    return {
      type: error.error.type || `https://httpstatuses.com/${error.status}`,
      title: error.error.title,
      status: error.status,
      detail: error.error.detail || error.error.message,
      instance: error.error.instance || error.url || undefined,
      errors: error.error.errors,
      code: error.error.code,
      traceId: error.error.traceId,
    };
  }

  // Handle network errors (status 0)
  if (error.status === 0) {
    return {
      type: 'https://httpstatuses.com/0',
      title: 'Network Error',
      status: 0,
      detail: 'Unable to connect to the server. Please check your internet connection.',
      instance: error.url || undefined,
    };
  }

  // Extract error message from various response formats
  let detail = 'An unexpected error occurred';
  if (error.error) {
    if (typeof error.error === 'string') {
      detail = error.error;
    } else if (error.error.message) {
      detail = error.error.message;
    } else if (error.error.error) {
      detail = typeof error.error.error === 'string' ? error.error.error : error.error.error.message;
    }
  } else if (error.message) {
    detail = error.message;
  }

  // Build RFC 7807 compliant error
  return {
    type: `https://httpstatuses.com/${error.status}`,
    title: STATUS_TITLES[error.status] || 'Error',
    status: error.status,
    detail,
    instance: error.url || undefined,
  };
}

/**
 * Error Interceptor Function
 *
 * Catches HTTP errors and transforms them into a consistent RFC 7807
 * ApiError format for easier error handling.
 *
 * Features:
 * - Converts all HTTP errors to ApiError format
 * - Preserves existing RFC 7807 error structures
 * - Handles network errors gracefully
 * - Provides meaningful error messages
 */
export const errorInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      // Skip 401 errors - let auth interceptor handle token refresh
      // The auth interceptor will either refresh the token and retry,
      // or logout and redirect to login page
      if (error.status === 401) {
        return throwError(() => error);
      }

      const apiError = transformToApiError(error);

      // Log error for debugging (can be replaced with error logging service)
      console.error('[API Error]', {
        url: request.url,
        method: request.method,
        status: apiError.status,
        title: apiError.title,
        detail: apiError.detail,
      });

      // Create a new HttpErrorResponse with the transformed error
      const transformedError = new HttpErrorResponse({
        error: { error: apiError },
        headers: error.headers,
        status: error.status,
        statusText: error.statusText,
        url: error.url || undefined,
      });

      return throwError(() => transformedError);
    })
  );
};

/**
 * Helper function to extract ApiError from HttpErrorResponse
 * Use this in services/components to get typed error
 */
export function extractApiError(error: unknown): ApiError {
  if (error instanceof HttpErrorResponse && error.error?.error) {
    return error.error.error as ApiError;
  }

  return {
    type: 'https://httpstatuses.com/500',
    title: 'Unknown Error',
    status: 500,
    detail: error instanceof Error ? error.message : 'An unexpected error occurred',
  };
}
