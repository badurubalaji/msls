/**
 * MSLS Core Interceptors - Barrel Export
 *
 * This file exports all HTTP interceptors from the core/interceptors directory.
 * Interceptors are functional (HttpInterceptorFn) following Angular 17+ patterns.
 */

// Auth Interceptor - Adds Authorization header and handles 401 with refresh
export { authInterceptor } from './auth.interceptor';

// Tenant Interceptor - Adds X-Tenant-ID header
export { tenantInterceptor } from './tenant.interceptor';

// Request ID Interceptor - Adds X-Request-ID header for tracing
export { requestIdInterceptor } from './request-id.interceptor';

// Error Interceptor - Transforms errors to RFC 7807 format
export { errorInterceptor, extractApiError } from './error.interceptor';

// Loading Interceptor - Manages global loading state
export { loadingInterceptor } from './loading.interceptor';
