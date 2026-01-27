/**
 * Base environment configuration for MSLS (Multi-School Learning System)
 *
 * This file contains the base configuration that is shared across all environments.
 * Environment-specific files (environment.development.ts, environment.production.ts)
 * will extend or override these values as needed.
 */

export interface Environment {
  /** Whether this is a production build */
  production: boolean;
  /** Application name displayed in the UI */
  appName: string;
  /** Application version from package.json */
  appVersion: string;
  /** Base URL for API calls */
  apiUrl: string;
  /** API version prefix */
  apiVersion: string;
  /** Whether to enable debug logging */
  enableDebugLogging: boolean;
  /** Session timeout in milliseconds (30 minutes default) */
  sessionTimeout: number;
  /** Whether to enable analytics tracking */
  enableAnalytics: boolean;
  /** Maximum file upload size in bytes (10MB default) */
  maxFileUploadSize: number;
  /** Supported file types for uploads */
  supportedFileTypes: string[];
}

/**
 * Base environment configuration
 * This serves as the default and is used during local development
 */
export const environment: Environment = {
  production: false,
  appName: 'MSLS - Multi-School Learning System',
  appVersion: '0.0.1',
  apiUrl: '/api',
  apiVersion: 'v1',
  enableDebugLogging: true,
  sessionTimeout: 30 * 60 * 1000, // 30 minutes
  enableAnalytics: false,
  maxFileUploadSize: 10 * 1024 * 1024, // 10MB
  supportedFileTypes: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/webp',
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/plain',
    'text/csv',
  ],
};
