/**
 * MSLS Authentication Models
 *
 * Defines authentication-related interfaces for login, tokens, and responses.
 */

import { User } from './user.model';

/**
 * Authentication tokens returned after successful login
 */
export interface AuthTokens {
  /** JWT access token for API requests */
  accessToken: string;

  /** Refresh token for obtaining new access tokens */
  refreshToken: string;

  /** Access token expiration time in seconds */
  expiresIn: number;

  /** Token type (typically 'Bearer') */
  tokenType?: string;
}

/**
 * Login request payload
 */
export interface LoginRequest {
  /** User's email address */
  email: string;

  /** User's password */
  password: string;

  /** Tenant identifier for multi-tenant authentication */
  tenantId: string;

  /** Remember me option for extended session */
  rememberMe?: boolean;
}

/**
 * Registration request payload
 */
export interface RegisterRequest {
  /** User's email address */
  email: string;

  /** User's password */
  password: string;

  /** Tenant identifier */
  tenantId: string;

  /** User's first name */
  firstName: string;

  /** User's last name */
  lastName: string;

  /** User's phone number (optional) */
  phone?: string;
}

/**
 * Password reset request payload
 */
export interface PasswordResetRequest {
  /** User's email address */
  email: string;

  /** Tenant identifier */
  tenantId: string;
}

/**
 * Password change request payload
 */
export interface PasswordChangeRequest {
  /** Current password */
  currentPassword: string;

  /** New password */
  newPassword: string;

  /** Confirm new password */
  confirmPassword: string;
}

/**
 * Refresh token request payload
 */
export interface RefreshTokenRequest {
  /** Current refresh token */
  refreshToken: string;
}

/**
 * Authentication response after successful login
 */
export interface AuthResponse {
  /** Authenticated user details */
  user: User;

  /** Authentication tokens */
  tokens: AuthTokens;
}

/**
 * Token verification response
 */
export interface TokenVerifyResponse {
  /** Whether the token is valid */
  valid: boolean;

  /** User if token is valid */
  user?: User;

  /** Token expiration timestamp */
  expiresAt?: string;
}
