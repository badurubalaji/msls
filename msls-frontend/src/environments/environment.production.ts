/**
 * Production environment configuration for MSLS (Multi-School Learning System)
 *
 * This configuration is used for production builds with `ng build --configuration production`.
 * It disables debug logging, enables analytics, and uses production API endpoints.
 */

import { Environment } from './environment';

export const environment: Environment = {
  production: true,
  appName: 'MSLS - Multi-School Learning System',
  appVersion: '0.0.1',
  apiUrl: '/api',
  apiVersion: 'v1',
  enableDebugLogging: false,
  sessionTimeout: 30 * 60 * 1000, // 30 minutes
  enableAnalytics: true,
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
