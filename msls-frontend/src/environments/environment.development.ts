/**
 * Development environment configuration for MSLS (Multi-School Learning System)
 *
 * This configuration is used during local development with `ng serve`.
 * It enables verbose logging, debug tools, and points to the local backend server.
 */

import { Environment } from './environment';

export const environment: Environment = {
  production: false,
  appName: 'MSLS - Multi-School Learning System (Development)',
  appVersion: '0.0.1-dev',
  apiUrl: 'http://localhost:8080/api',
  apiVersion: 'v1',
  enableDebugLogging: true,
  sessionTimeout: 60 * 60 * 1000, // 1 hour for development convenience
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
