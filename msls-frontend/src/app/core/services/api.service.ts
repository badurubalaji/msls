/**
 * MSLS API Service
 *
 * Base HTTP client wrapper providing typed API methods with standardized
 * request/response handling.
 */

import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams, HttpContext } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiResponse, PaginatedResponse, PaginationParams } from '../models/api-response.model';

/** Environment configuration (to be replaced with actual environment) */
const API_BASE_URL = '/api/v1';

/** Request options interface */
export interface RequestOptions {
  /** HTTP headers */
  headers?: HttpHeaders | Record<string, string | string[]>;

  /** URL query parameters */
  params?: HttpParams | Record<string, string | number | boolean | string[]>;

  /** HTTP context */
  context?: HttpContext;

  /** Whether to include credentials */
  withCredentials?: boolean;

  /** Response type */
  responseType?: 'json';

  /** Report progress */
  reportProgress?: boolean;
}

/**
 * ApiService - Base HTTP client wrapper with typed methods.
 *
 * Provides standardized CRUD operations with automatic response unwrapping.
 *
 * Usage:
 * constructor(private apiService: ApiService) {}
 *
 * getUsers() {
 *   return this.apiService.get<User[]>('/users');
 * }
 *
 * createUser(user: CreateUserDto) {
 *   return this.apiService.post<User>('/users', user);
 * }
 */
@Injectable({ providedIn: 'root' })
export class ApiService {
  private http = inject(HttpClient);

  /** Base URL for API requests */
  private baseUrl = API_BASE_URL;

  /**
   * Perform a GET request
   * @param endpoint - API endpoint (relative to base URL)
   * @param options - Request options
   */
  get<T>(endpoint: string, options?: RequestOptions): Observable<T> {
    return this.http
      .get<ApiResponse<T>>(this.buildUrl(endpoint), this.buildOptions(options))
      .pipe(map(response => this.extractData(response)));
  }

  /**
   * Perform a GET request with pagination
   * @param endpoint - API endpoint (relative to base URL)
   * @param pagination - Pagination parameters
   * @param options - Additional request options
   */
  getList<T>(
    endpoint: string,
    pagination?: PaginationParams,
    options?: RequestOptions
  ): Observable<PaginatedResponse<T>> {
    const params = this.buildPaginationParams(pagination, options?.params);
    return this.http.get<PaginatedResponse<T>>(this.buildUrl(endpoint), {
      ...this.buildOptions(options),
      params,
    });
  }

  /**
   * Perform a POST request
   * @param endpoint - API endpoint (relative to base URL)
   * @param body - Request body
   * @param options - Request options
   */
  post<T>(endpoint: string, body: unknown, options?: RequestOptions): Observable<T> {
    return this.http
      .post<ApiResponse<T>>(this.buildUrl(endpoint), body, this.buildOptions(options))
      .pipe(map(response => this.extractData(response)));
  }

  /**
   * Perform a PUT request
   * @param endpoint - API endpoint (relative to base URL)
   * @param body - Request body
   * @param options - Request options
   */
  put<T>(endpoint: string, body: unknown, options?: RequestOptions): Observable<T> {
    return this.http
      .put<ApiResponse<T>>(this.buildUrl(endpoint), body, this.buildOptions(options))
      .pipe(map(response => this.extractData(response)));
  }

  /**
   * Perform a PATCH request
   * @param endpoint - API endpoint (relative to base URL)
   * @param body - Request body
   * @param options - Request options
   */
  patch<T>(endpoint: string, body: unknown, options?: RequestOptions): Observable<T> {
    return this.http
      .patch<ApiResponse<T>>(this.buildUrl(endpoint), body, this.buildOptions(options))
      .pipe(map(response => this.extractData(response)));
  }

  /**
   * Perform a DELETE request
   * @param endpoint - API endpoint (relative to base URL)
   * @param options - Request options
   */
  delete<T>(endpoint: string, options?: RequestOptions): Observable<T> {
    return this.http
      .delete<ApiResponse<T>>(this.buildUrl(endpoint), this.buildOptions(options))
      .pipe(map(response => {
        // Handle 204 No Content responses (empty body)
        if (response === null || response === undefined) {
          return undefined as T;
        }
        return this.extractData(response);
      }));
  }

  /**
   * Perform a raw GET request (without response unwrapping)
   * @param endpoint - API endpoint
   * @param options - Request options
   */
  getRaw<T>(endpoint: string, options?: RequestOptions): Observable<T> {
    return this.http.get<T>(this.buildUrl(endpoint), this.buildOptions(options));
  }

  /**
   * Perform a raw POST request (without response unwrapping)
   * @param endpoint - API endpoint
   * @param body - Request body
   * @param options - Request options
   */
  postRaw<T>(endpoint: string, body: unknown, options?: RequestOptions): Observable<T> {
    return this.http.post<T>(this.buildUrl(endpoint), body, this.buildOptions(options));
  }

  /**
   * Download a file as Blob
   * @param endpoint - API endpoint
   * @param options - Request options
   */
  getBlob(endpoint: string, options?: RequestOptions): Observable<Blob> {
    return this.http.get(this.buildUrl(endpoint), {
      ...this.buildOptions(options),
      responseType: 'blob',
    });
  }

  /**
   * Upload a file
   * @param endpoint - API endpoint
   * @param file - File to upload
   * @param fieldName - Form field name for the file
   * @param additionalData - Additional form data
   */
  uploadFile<T>(
    endpoint: string,
    file: File,
    fieldName: string = 'file',
    additionalData?: Record<string, string>
  ): Observable<T> {
    const formData = new FormData();
    formData.append(fieldName, file, file.name);

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value);
      });
    }

    return this.http
      .post<ApiResponse<T>>(this.buildUrl(endpoint), formData, {
        reportProgress: true,
      })
      .pipe(map(response => this.extractData(response)));
  }

  /**
   * Build full URL from endpoint
   */
  private buildUrl(endpoint: string): string {
    // If endpoint is already a full URL, return as-is
    if (endpoint.startsWith('http://') || endpoint.startsWith('https://')) {
      return endpoint;
    }

    // Ensure endpoint starts with /
    const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    return `${this.baseUrl}${normalizedEndpoint}`;
  }

  /**
   * Build request options
   */
  private buildOptions(options?: RequestOptions): RequestOptions {
    return {
      headers: options?.headers,
      params: this.buildHttpParams(options?.params),
      context: options?.context,
      withCredentials: options?.withCredentials ?? false,
      reportProgress: options?.reportProgress,
    };
  }

  /**
   * Build HttpParams from various input types
   */
  private buildHttpParams(
    params?: HttpParams | Record<string, string | number | boolean | string[]>
  ): HttpParams | undefined {
    if (!params) return undefined;

    if (params instanceof HttpParams) {
      return params;
    }

    let httpParams = new HttpParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          value.forEach(v => {
            httpParams = httpParams.append(key, v);
          });
        } else {
          httpParams = httpParams.set(key, String(value));
        }
      }
    });

    return httpParams;
  }

  /**
   * Build pagination params
   */
  private buildPaginationParams(
    pagination?: PaginationParams,
    existingParams?: HttpParams | Record<string, string | number | boolean | string[]>
  ): HttpParams {
    let params = this.buildHttpParams(existingParams) || new HttpParams();

    if (pagination) {
      if (pagination.page !== undefined) {
        params = params.set('page', String(pagination.page));
      }
      if (pagination.pageSize !== undefined) {
        params = params.set('pageSize', String(pagination.pageSize));
      }
      if (pagination.sortBy) {
        params = params.set('sortBy', pagination.sortBy);
      }
      if (pagination.sortOrder) {
        params = params.set('sortOrder', pagination.sortOrder);
      }
    }

    return params;
  }

  /**
   * Extract data from API response
   */
  private extractData<T>(response: ApiResponse<T>): T {
    if (response.success && response.data !== undefined) {
      return response.data;
    }
    throw response.error || new Error('Unknown error occurred');
  }
}
