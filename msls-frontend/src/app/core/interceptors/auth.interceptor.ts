/**
 * MSLS Auth Interceptor
 *
 * Adds Authorization header to requests and handles 401 responses
 * with automatic token refresh and request retry.
 */

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpErrorResponse, HttpEvent } from '@angular/common/http';
import { inject } from '@angular/core';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, filter, switchMap, take, finalize } from 'rxjs/operators';

import { AuthService } from '../services/auth.service';

/** URLs that should skip authentication */
const PUBLIC_URLS = [
  '/auth/login',
  '/auth/register',
  '/auth/refresh',
  '/auth/forgot-password',
  '/public/',
];

/** Subject for queuing requests during token refresh */
let isRefreshing = false;
let refreshTokenSubject: BehaviorSubject<string | null> = new BehaviorSubject<string | null>(null);

/**
 * Check if URL should skip authentication
 */
function isPublicUrl(url: string): boolean {
  return PUBLIC_URLS.some(publicUrl => url.includes(publicUrl));
}

/**
 * Add authorization header to request
 */
function addAuthHeader(request: HttpRequest<unknown>, token: string): HttpRequest<unknown> {
  return request.clone({
    setHeaders: {
      Authorization: `Bearer ${token}`,
    },
  });
}

/**
 * Auth Interceptor Function
 *
 * Features:
 * - Adds Bearer token to Authorization header
 * - Handles 401 responses with token refresh
 * - Queues requests during token refresh
 * - Retries failed requests after successful refresh
 */
export const authInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
  const authService = inject(AuthService);

  // Skip authentication for public URLs
  if (isPublicUrl(request.url)) {
    return next(request);
  }

  // Get current access token
  const accessToken = authService.accessToken();

  // Add auth header if token exists
  if (accessToken) {
    request = addAuthHeader(request, accessToken);
  }

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      // Handle 401 Unauthorized - attempt token refresh
      if (error.status === 401 && !isPublicUrl(request.url)) {
        console.log('[Auth] Received 401 for:', request.url);
        return handle401Error(request, next, authService);
      }

      return throwError(() => error);
    })
  );
};

/**
 * Handle 401 error with token refresh
 */
function handle401Error(
  request: HttpRequest<unknown>,
  next: HttpHandlerFn,
  authService: AuthService
): Observable<HttpEvent<unknown>> {
  // Check if we have a refresh token
  const refreshToken = authService.refreshToken();
  if (!refreshToken) {
    console.log('[Auth] No refresh token available, logging out');
    authService.logout(true);
    return throwError(() => new Error('No refresh token available'));
  }

  // If already refreshing, queue the request
  if (isRefreshing) {
    console.log('[Auth] Token refresh in progress, queuing request:', request.url);
    return refreshTokenSubject.pipe(
      filter(token => token !== null),
      take(1),
      switchMap(token => {
        console.log('[Auth] Retrying queued request with new token:', request.url);
        return next(addAuthHeader(request, token!));
      })
    );
  }

  // Start token refresh
  console.log('[Auth] Starting token refresh for:', request.url);
  isRefreshing = true;
  refreshTokenSubject.next(null);

  return authService.refreshAccessToken().pipe(
    switchMap(tokens => {
      console.log('[Auth] Token refresh successful, retrying request:', request.url);
      refreshTokenSubject.next(tokens.accessToken);
      return next(addAuthHeader(request, tokens.accessToken));
    }),
    catchError(refreshError => {
      // Token refresh failed - logout and redirect
      console.error('[Auth] Token refresh failed, logging out:', refreshError);
      authService.logout(true);
      return throwError(() => refreshError);
    }),
    finalize(() => {
      isRefreshing = false;
    })
  );
}
