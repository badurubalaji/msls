/**
 * MSLS Tenant Interceptor
 *
 * Adds X-Tenant-ID header to all API requests for multi-tenant support.
 */

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpEvent } from '@angular/common/http';
import { inject } from '@angular/core';
import { Observable } from 'rxjs';

import { TenantService } from '../services/tenant.service';

/** Header name for tenant identification */
const TENANT_HEADER = 'X-Tenant-ID';

/**
 * Tenant Interceptor Function
 *
 * Automatically adds the X-Tenant-ID header to all outgoing requests
 * if a tenant is set in the TenantService.
 *
 * Usage:
 * Requests will automatically include:
 * X-Tenant-ID: {tenantId}
 */
export const tenantInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
  const tenantService = inject(TenantService);

  // Get current tenant ID
  const tenantId = tenantService.currentTenantId();

  // If no tenant ID, proceed without modification
  if (!tenantId) {
    return next(request);
  }

  // Clone request with tenant header
  const modifiedRequest = request.clone({
    setHeaders: {
      [TENANT_HEADER]: tenantId,
    },
  });

  return next(modifiedRequest);
};
