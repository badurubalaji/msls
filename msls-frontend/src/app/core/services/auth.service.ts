/**
 * MSLS Authentication Service
 *
 * Manages authentication state, token handling, and user sessions.
 * Uses Angular Signals for reactive state management.
 */

import { Injectable, computed, inject, signal } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, BehaviorSubject, throwError, of } from 'rxjs';
import { catchError, tap, switchMap, finalize, filter, take } from 'rxjs/operators';

import { User } from '../models/user.model';
import {
  AuthResponse,
  AuthTokens,
  LoginRequest,
  RefreshTokenRequest,
} from '../models/auth.model';
import { ApiResponse } from '../models/api-response.model';
import { StorageService, STORAGE_KEYS } from './storage.service';
import { TenantService } from './tenant.service';
import { environment } from '../../../environments/environment';

/**
 * AuthService - Manages authentication state with Angular Signals.
 *
 * Features:
 * - Signal-based reactive state
 * - Automatic token refresh
 * - Request queuing during token refresh
 * - Persistent session storage
 *
 * Usage:
 * constructor(private authService: AuthService) {}
 *
 * login() {
 *   this.authService.login({ email, password, tenantId }).subscribe({
 *     next: (response) => console.log('Logged in:', response.user),
 *     error: (error) => console.error('Login failed:', error)
 *   });
 * }
 *
 * // In templates:
 * @if (authService.isAuthenticated()) {
 *   <p>Welcome, {{ authService.currentUser()?.firstName }}</p>
 * }
 */
@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private router = inject(Router);
  private storageService = inject(StorageService);
  private tenantService = inject(TenantService);

  /** API base URL from environment configuration */
  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}`;

  /** Internal signals for auth state */
  private _currentUser = signal<User | null>(null);
  private _accessToken = signal<string | null>(null);
  private _refreshToken = signal<string | null>(null);

  /** Token refresh state management */
  private isRefreshing = false;
  private refreshTokenSubject = new BehaviorSubject<string | null>(null);

  /** Public readonly signals */
  readonly currentUser = this._currentUser.asReadonly();
  readonly accessToken = this._accessToken.asReadonly();
  readonly refreshToken = this._refreshToken.asReadonly();

  /** Computed authentication state */
  readonly isAuthenticated = computed(() => !!this._accessToken());

  /** Computed user display name */
  readonly displayName = computed(() => {
    const user = this._currentUser();
    if (!user) return null;
    return `${user.firstName} ${user.lastName}`.trim();
  });

  constructor() {
    this.initializeFromStorage();
  }

  /**
   * Initialize auth state from storage on service creation
   */
  initializeFromStorage(): void {
    const accessToken = this.storageService.getItem<string>(STORAGE_KEYS.ACCESS_TOKEN);
    const refreshToken = this.storageService.getItem<string>(STORAGE_KEYS.REFRESH_TOKEN);
    const user = this.storageService.getItem<User>(STORAGE_KEYS.CURRENT_USER);

    if (accessToken) {
      this._accessToken.set(accessToken);
    }

    if (refreshToken) {
      this._refreshToken.set(refreshToken);
    }

    if (user) {
      this._currentUser.set(user);
    }
  }

  /**
   * Authenticate user with credentials
   * @param credentials - Login request with email, password, and tenantId
   */
  login(credentials: LoginRequest): Observable<AuthResponse> {
    // Transform to snake_case for backend API
    const payload = {
      email: credentials.email,
      password: credentials.password,
      tenant_id: credentials.tenantId,
    };

    return this.http
      .post<ApiResponse<unknown>>(`${this.apiUrl}/auth/login`, payload)
      .pipe(
        tap(response => {
          if (response.success && response.data) {
            const authResponse = this.mapLoginResponse(response.data);
            this.handleAuthSuccess(authResponse);
          }
        }),
        switchMap(response => {
          if (response.success && response.data) {
            return of(this.mapLoginResponse(response.data));
          }
          return throwError(() => response.error);
        }),
        catchError(this.handleError)
      );
  }

  /**
   * Map backend login response (snake_case) to frontend model (camelCase)
   */
  private mapLoginResponse(data: unknown): AuthResponse {
    const raw = data as Record<string, unknown>;
    const userRaw = raw['user'] as Record<string, unknown>;

    // Map roles from backend format (array of objects with 'name') to frontend format
    const rolesArray = userRaw['roles'] as Array<Record<string, unknown>> | undefined;
    const roles = rolesArray?.map(r => r['name'] as string) || [];

    // Permissions come as string array from backend
    const permissions = (userRaw['permissions'] as string[]) || [];

    return {
      user: {
        id: userRaw['id'] as string,
        email: userRaw['email'] as string,
        firstName: userRaw['first_name'] as string,
        lastName: userRaw['last_name'] as string,
        tenantId: userRaw['tenant_id'] as string,
        status: (userRaw['status'] as string) as User['status'],
        createdAt: userRaw['created_at'] as string,
        updatedAt: userRaw['created_at'] as string, // Use created_at as fallback
        roles: roles as User['roles'],
        permissions: permissions as User['permissions'],
      },
      tokens: {
        accessToken: raw['access_token'] as string,
        refreshToken: raw['refresh_token'] as string,
        expiresIn: raw['expires_in'] as number,
      },
    };
  }

  /**
   * Log out the current user
   * @param redirectToLogin - Whether to redirect to login page
   */
  logout(redirectToLogin: boolean = true): Observable<void> {
    const refreshToken = this._refreshToken();
    const accessToken = this._accessToken();

    // Notify server first (fire and forget), then clear local state
    // Need to include Authorization header manually since we're about to clear tokens
    if (refreshToken && accessToken) {
      this.http
        .post(
          `${this.apiUrl}/auth/logout`,
          { refresh_token: refreshToken },
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          }
        )
        .pipe(catchError(() => of(null)))
        .subscribe();
    }

    // Clear local state after making the API call
    this.clearTokens();

    if (redirectToLogin) {
      this.router.navigate(['/login']);
    }

    return of(undefined);
  }

  /**
   * Refresh the access token using the refresh token
   */
  refreshAccessToken(): Observable<AuthTokens> {
    const currentRefreshToken = this._refreshToken();

    if (!currentRefreshToken) {
      this.logout();
      return throwError(() => new Error('No refresh token available'));
    }

    // If already refreshing, wait for the result
    if (this.isRefreshing) {
      return this.refreshTokenSubject.pipe(
        filter(token => token !== null),
        take(1),
        switchMap(token => {
          if (token) {
            return of({
              accessToken: token,
              refreshToken: this._refreshToken() || '',
              expiresIn: 3600,
            } as AuthTokens);
          }
          return throwError(() => new Error('Token refresh failed'));
        })
      );
    }

    this.isRefreshing = true;
    this.refreshTokenSubject.next(null);

    // Backend expects snake_case: refresh_token
    const request = { refresh_token: currentRefreshToken };

    return this.http
      .post<ApiResponse<unknown>>(`${this.apiUrl}/auth/refresh`, request)
      .pipe(
        tap(response => {
          if (response.success && response.data) {
            const tokens = this.mapRefreshResponse(response.data);
            this.setTokens(tokens);
            this.refreshTokenSubject.next(tokens.accessToken);
          }
        }),
        switchMap(response => {
          if (response.success && response.data) {
            return of(this.mapRefreshResponse(response.data));
          }
          return throwError(() => response.error);
        }),
        catchError(error => {
          this.logout();
          return throwError(() => error);
        }),
        finalize(() => {
          this.isRefreshing = false;
        })
      );
  }

  /**
   * Map backend refresh response (snake_case) to frontend model (camelCase)
   */
  private mapRefreshResponse(data: unknown): AuthTokens {
    const raw = data as Record<string, unknown>;
    return {
      accessToken: raw['access_token'] as string,
      refreshToken: raw['refresh_token'] as string,
      expiresIn: raw['expires_in'] as number,
    };
  }

  /**
   * Set authentication tokens
   * @param tokens - Auth tokens to store
   */
  setTokens(tokens: AuthTokens): void {
    this._accessToken.set(tokens.accessToken);
    this._refreshToken.set(tokens.refreshToken);

    this.storageService.setItem(STORAGE_KEYS.ACCESS_TOKEN, tokens.accessToken);
    this.storageService.setItem(STORAGE_KEYS.REFRESH_TOKEN, tokens.refreshToken);

    // Store token expiry
    const expiryTime = Date.now() + tokens.expiresIn * 1000;
    this.storageService.setItem(STORAGE_KEYS.TOKEN_EXPIRY, expiryTime.toString());
  }

  /**
   * Clear all authentication tokens and user data
   */
  clearTokens(): void {
    this._accessToken.set(null);
    this._refreshToken.set(null);
    this._currentUser.set(null);

    this.storageService.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    this.storageService.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    this.storageService.removeItem(STORAGE_KEYS.TOKEN_EXPIRY);
    this.storageService.removeItem(STORAGE_KEYS.CURRENT_USER);
  }

  /**
   * Set the current user
   * @param user - User to set
   */
  setCurrentUser(user: User): void {
    this._currentUser.set(user);
    this.storageService.setItem(STORAGE_KEYS.CURRENT_USER, user);
  }

  /**
   * Check if token is expired or will expire soon
   * @param bufferSeconds - Buffer time in seconds (default: 60)
   */
  isTokenExpired(bufferSeconds: number = 60): boolean {
    const expiryStr = this.storageService.getItem<string>(STORAGE_KEYS.TOKEN_EXPIRY);
    if (!expiryStr) return true;

    const expiry = parseInt(expiryStr, 10);
    const now = Date.now();
    const buffer = bufferSeconds * 1000;

    return now >= expiry - buffer;
  }

  /**
   * Check if user has a specific permission
   * @param permission - Permission to check
   */
  hasPermission(permission: string): boolean {
    const user = this._currentUser();
    return user?.permissions?.includes(permission as never) ?? false;
  }

  /**
   * Check if user has a specific role
   * @param role - Role to check
   */
  hasRole(role: string): boolean {
    const user = this._currentUser();
    return user?.roles?.includes(role as never) ?? false;
  }

  /**
   * Handle successful authentication
   */
  private handleAuthSuccess(response: AuthResponse): void {
    this.setTokens(response.tokens);
    this.setCurrentUser(response.user);

    // Set tenant if available from user
    if (response.user.tenantId) {
      this.tenantService.setTenantId(response.user.tenantId);
    }
  }

  /**
   * Handle HTTP errors
   */
  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'An error occurred';

    if (error.error?.error) {
      errorMessage = error.error.error.detail || error.error.error.title;
    } else if (error.message) {
      errorMessage = error.message;
    }

    return throwError(() => new Error(errorMessage));
  }

  /**
   * Observable to wait for token refresh completion
   * Used by interceptors to queue requests during refresh
   */
  get refreshTokenSubject$(): Observable<string | null> {
    return this.refreshTokenSubject.asObservable();
  }

  /**
   * Check if token refresh is in progress
   */
  get isRefreshingToken(): boolean {
    return this.isRefreshing;
  }
}
