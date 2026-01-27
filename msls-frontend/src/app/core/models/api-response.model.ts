/**
 * MSLS API Response Models
 *
 * Defines standard API response wrappers and error handling types
 * following RFC 7807 Problem Details for HTTP APIs.
 */

/**
 * Pagination metadata for list responses
 */
export interface PaginationMeta {
  /** Current page number (1-indexed) */
  page: number;

  /** Number of items per page */
  pageSize: number;

  /** Total number of items across all pages */
  totalItems: number;

  /** Total number of pages */
  totalPages: number;

  /** Whether there is a next page */
  hasNextPage: boolean;

  /** Whether there is a previous page */
  hasPreviousPage: boolean;
}

/**
 * Pagination request parameters
 */
export interface PaginationParams {
  /** Page number to retrieve (1-indexed) */
  page?: number;

  /** Number of items per page */
  pageSize?: number;

  /** Field to sort by */
  sortBy?: string;

  /** Sort direction */
  sortOrder?: 'asc' | 'desc';
}

/**
 * RFC 7807 Problem Details error format
 */
export interface ApiError {
  /** URI reference identifying the problem type */
  type?: string;

  /** Short human-readable summary */
  title: string;

  /** HTTP status code */
  status: number;

  /** Human-readable explanation specific to this occurrence */
  detail?: string;

  /** URI reference identifying the specific occurrence */
  instance?: string;

  /** Validation errors keyed by field name */
  errors?: Record<string, string[]>;

  /** Error code for programmatic handling */
  code?: string;

  /** Trace ID for debugging */
  traceId?: string;
}

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T> {
  /** Whether the request was successful */
  success: boolean;

  /** Response payload when successful */
  data?: T;

  /** Error details when unsuccessful */
  error?: ApiError;

  /** Response metadata */
  meta?: {
    /** Server timestamp */
    timestamp: string;

    /** Request duration in milliseconds */
    duration?: number;

    /** API version */
    version?: string;
  };
}

/**
 * Paginated list response
 */
export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  /** Pagination metadata */
  pagination: PaginationMeta;
}

/**
 * Type guard to check if response is successful
 */
export function isSuccessResponse<T>(response: ApiResponse<T>): response is ApiResponse<T> & { data: T } {
  return response.success && response.data !== undefined;
}

/**
 * Type guard to check if response is an error
 */
export function isErrorResponse<T>(response: ApiResponse<T>): response is ApiResponse<T> & { error: ApiError } {
  return !response.success && response.error !== undefined;
}

/**
 * Create a success response
 */
export function createSuccessResponse<T>(data: T): ApiResponse<T> {
  return {
    success: true,
    data,
    meta: {
      timestamp: new Date().toISOString(),
    },
  };
}

/**
 * Create an error response
 */
export function createErrorResponse<T>(error: ApiError): ApiResponse<T> {
  return {
    success: false,
    error,
    meta: {
      timestamp: new Date().toISOString(),
    },
  };
}
