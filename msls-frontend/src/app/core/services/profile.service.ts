/**
 * MSLS Profile Service
 *
 * Manages user profile operations including profile updates, password changes,
 * avatar uploads, and preference management.
 */

import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient, HttpEventType, HttpEvent, HttpResponse } from '@angular/common/http';
import { Observable, throwError, of } from 'rxjs';
import { catchError, tap, switchMap, filter, map } from 'rxjs/operators';

import {
  UserProfile,
  UpdateProfileRequest,
  ChangePasswordRequest,
  NotificationPreferences,
  UpdatePreferencesRequest,
  AvatarUploadResponse,
  UserPreference,
  UserPreferencesResponse,
  SetPreferenceRequest,
} from '../models/profile.model';
import { ApiResponse } from '../models/api-response.model';
import { environment } from '../../../environments/environment';

/** Upload progress state */
export interface UploadProgress {
  /** Whether upload is in progress */
  uploading: boolean;

  /** Upload progress percentage (0-100) */
  progress: number;

  /** Error message if upload failed */
  error?: string;
}

/**
 * ProfileService - Manages user profile operations with Angular Signals.
 *
 * Features:
 * - Signal-based reactive state for profile data
 * - Avatar upload with progress tracking
 * - Notification preference management
 * - Extended preference storage
 *
 * Usage:
 * constructor(private profileService: ProfileService) {}
 *
 * ngOnInit() {
 *   this.profileService.loadProfile();
 * }
 *
 * // In templates:
 * @if (profileService.profile(); as profile) {
 *   <p>Welcome, {{ profile.fullName }}</p>
 * }
 */
@Injectable({ providedIn: 'root' })
export class ProfileService {
  private http = inject(HttpClient);

  /** API base URL from environment configuration */
  private readonly apiUrl = `${environment.apiUrl}/${environment.apiVersion}`;

  /** Internal signals for profile state */
  private _profile = signal<UserProfile | null>(null);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);
  private _uploadProgress = signal<UploadProgress>({
    uploading: false,
    progress: 0,
  });

  /** Public readonly signals */
  readonly profile = this._profile.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();
  readonly uploadProgress = this._uploadProgress.asReadonly();

  /** Computed profile helpers */
  readonly fullName = computed(() => {
    const profile = this._profile();
    if (!profile) return null;
    return profile.fullName || `${profile.firstName} ${profile.lastName}`.trim();
  });

  readonly avatarUrl = computed(() => {
    const profile = this._profile();
    return profile?.avatarUrl || null;
  });

  readonly notificationPreferences = computed(() => {
    const profile = this._profile();
    return (
      profile?.notificationPreferences || {
        email: true,
        push: true,
        sms: false,
      }
    );
  });

  /**
   * Load the current user's profile
   */
  loadProfile(): Observable<UserProfile> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.get<ApiResponse<unknown>>(`${this.apiUrl}/profile`).pipe(
      tap((response) => {
        if (response.success && response.data) {
          const profile = this.mapProfileResponse(response.data);
          this._profile.set(profile);
        }
        this._loading.set(false);
      }),
      switchMap((response) => {
        if (response.success && response.data) {
          return of(this.mapProfileResponse(response.data));
        }
        return throwError(() => response.error);
      }),
      catchError((err) => {
        this._loading.set(false);
        this._error.set(err.message || 'Failed to load profile');
        return throwError(() => err);
      })
    );
  }

  /**
   * Map backend profile response (snake_case) to frontend model (camelCase)
   */
  private mapProfileResponse(data: unknown): UserProfile {
    const raw = data as Record<string, unknown>;

    // Map roles from backend format
    const rolesRaw = raw['roles'] as Array<Record<string, unknown>> | undefined;
    const roles = rolesRaw?.map(r => ({
      id: r['id'] as string,
      name: r['name'] as string,
      description: r['description'] as string | undefined,
      isSystem: r['is_system'] as boolean,
    })) || [];

    // Map notification preferences
    const notifPrefs = raw['notification_preferences'] as Record<string, boolean> | undefined;

    return {
      id: raw['id'] as string,
      tenantId: raw['tenant_id'] as string,
      email: raw['email'] as string | undefined,
      phone: raw['phone'] as string | undefined,
      firstName: raw['first_name'] as string,
      lastName: raw['last_name'] as string,
      fullName: raw['full_name'] as string | undefined,
      avatarUrl: raw['avatar_url'] as string | undefined,
      bio: raw['bio'] as string | undefined,
      timezone: raw['timezone'] as string,
      locale: raw['locale'] as string,
      notificationPreferences: notifPrefs ? {
        email: notifPrefs['email'] ?? true,
        push: notifPrefs['push'] ?? true,
        sms: notifPrefs['sms'] ?? false,
      } : undefined,
      status: raw['status'] as UserProfile['status'],
      emailVerifiedAt: raw['email_verified_at'] as string | undefined,
      phoneVerifiedAt: raw['phone_verified_at'] as string | undefined,
      lastLoginAt: raw['last_login_at'] as string | undefined,
      twoFactorEnabled: raw['two_factor_enabled'] as boolean,
      createdAt: raw['created_at'] as string,
      roles,
      permissions: raw['permissions'] as string[] | undefined,
    };
  }

  /**
   * Update the current user's profile
   * @param data - Profile update data
   */
  updateProfile(data: UpdateProfileRequest): Observable<UserProfile> {
    this._loading.set(true);
    this._error.set(null);

    // Transform to snake_case for backend
    const payload: Record<string, unknown> = {};
    if (data.firstName !== undefined) payload['first_name'] = data.firstName;
    if (data.lastName !== undefined) payload['last_name'] = data.lastName;
    if (data.phone !== undefined) payload['phone'] = data.phone;
    if (data.bio !== undefined) payload['bio'] = data.bio;
    if (data.timezone !== undefined) payload['timezone'] = data.timezone;
    if (data.locale !== undefined) payload['locale'] = data.locale;

    return this.http.put<ApiResponse<unknown>>(`${this.apiUrl}/profile`, payload).pipe(
      tap((response) => {
        if (response.success && response.data) {
          const profile = this.mapProfileResponse(response.data);
          this._profile.set(profile);
        }
        this._loading.set(false);
      }),
      switchMap((response) => {
        if (response.success && response.data) {
          return of(this.mapProfileResponse(response.data));
        }
        return throwError(() => response.error);
      }),
      catchError((err) => {
        this._loading.set(false);
        this._error.set(err.message || 'Failed to update profile');
        return throwError(() => err);
      })
    );
  }

  /**
   * Change the current user's password
   * @param data - Password change data
   */
  changePassword(data: ChangePasswordRequest): Observable<{ message: string }> {
    this._loading.set(true);
    this._error.set(null);

    // Transform to snake_case for backend
    const payload = {
      current_password: data.currentPassword,
      new_password: data.newPassword,
      confirm_password: data.confirmPassword,
    };

    return this.http
      .put<ApiResponse<{ message: string }>>(`${this.apiUrl}/profile/password`, payload)
      .pipe(
        tap(() => {
          this._loading.set(false);
        }),
        switchMap((response) => {
          if (response.success && response.data) {
            return of(response.data);
          }
          return throwError(() => response.error);
        }),
        catchError((err) => {
          this._loading.set(false);
          this._error.set(err.message || 'Failed to change password');
          return throwError(() => err);
        })
      );
  }

  /**
   * Upload a new avatar image
   * @param file - The image file to upload
   */
  uploadAvatar(file: File): Observable<AvatarUploadResponse> {
    // Reset upload progress
    this._uploadProgress.set({ uploading: true, progress: 0 });
    this._error.set(null);

    const formData = new FormData();
    formData.append('avatar', file, file.name);

    return this.http
      .post<ApiResponse<AvatarUploadResponse>>(`${this.apiUrl}/profile/avatar`, formData, {
        reportProgress: true,
        observe: 'events',
      })
      .pipe(
        tap((event: HttpEvent<ApiResponse<AvatarUploadResponse>>) => {
          if (event.type === HttpEventType.UploadProgress && event.total) {
            const progress = Math.round((100 * event.loaded) / event.total);
            this._uploadProgress.update((state) => ({ ...state, progress }));
          }
        }),
        filter((event): event is HttpResponse<ApiResponse<AvatarUploadResponse>> =>
          event.type === HttpEventType.Response
        ),
        map((event) => event.body!),
        tap((response) => {
          if (response?.success && response.data) {
            // Update profile with new avatar URL
            this._profile.update((profile) =>
              profile ? { ...profile, avatarUrl: response.data?.avatarUrl } : null
            );
          }
          this._uploadProgress.set({ uploading: false, progress: 100 });
        }),
        switchMap((response) => {
          if (response?.success && response.data) {
            return of(response.data);
          }
          return throwError(() => response?.error || new Error('Upload failed'));
        }),
        catchError((err) => {
          this._uploadProgress.set({ uploading: false, progress: 0, error: err.message });
          this._error.set(err.message || 'Failed to upload avatar');
          return throwError(() => err);
        })
      );
  }

  /**
   * Get notification preferences
   */
  getNotificationPreferences(): Observable<NotificationPreferences> {
    return this.http
      .get<ApiResponse<NotificationPreferences>>(`${this.apiUrl}/profile/preferences`)
      .pipe(
        switchMap((response) => {
          if (response.success && response.data) {
            return of(response.data);
          }
          return throwError(() => response.error);
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to get preferences');
          return throwError(() => err);
        })
      );
  }

  /**
   * Update notification preferences
   * @param data - Preferences update data
   */
  updateNotificationPreferences(data: UpdatePreferencesRequest): Observable<NotificationPreferences> {
    this._loading.set(true);
    this._error.set(null);

    return this.http
      .put<ApiResponse<NotificationPreferences>>(`${this.apiUrl}/profile/preferences`, data)
      .pipe(
        tap((response) => {
          if (response.success && response.data) {
            // Update profile with new preferences
            this._profile.update((profile) =>
              profile ? { ...profile, notificationPreferences: response.data } : null
            );
          }
          this._loading.set(false);
        }),
        switchMap((response) => {
          if (response.success && response.data) {
            return of(response.data);
          }
          return throwError(() => response.error);
        }),
        catchError((err) => {
          this._loading.set(false);
          this._error.set(err.message || 'Failed to update preferences');
          return throwError(() => err);
        })
      );
  }

  /**
   * Request account deletion
   */
  requestAccountDeletion(): Observable<{ message: string }> {
    this._loading.set(true);
    this._error.set(null);

    return this.http.delete<ApiResponse<{ message: string }>>(`${this.apiUrl}/profile`).pipe(
      tap(() => {
        this._loading.set(false);
      }),
      switchMap((response) => {
        if (response.success && response.data) {
          return of(response.data);
        }
        return throwError(() => response.error);
      }),
      catchError((err) => {
        this._loading.set(false);
        this._error.set(err.message || 'Failed to request account deletion');
        return throwError(() => err);
      })
    );
  }

  /**
   * Get extended user preferences
   * @param category - Optional category filter
   */
  getExtendedPreferences(category?: string): Observable<UserPreference[]> {
    const options = category ? { params: { category } } : {};

    return this.http
      .get<ApiResponse<UserPreferencesResponse>>(`${this.apiUrl}/profile/preferences/extended`, options)
      .pipe(
        switchMap((response) => {
          if (response.success && response.data) {
            return of(response.data.preferences);
          }
          return throwError(() => response.error);
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to get extended preferences');
          return throwError(() => err);
        })
      );
  }

  /**
   * Set an extended user preference
   * @param data - Preference data
   */
  setExtendedPreference(data: SetPreferenceRequest): Observable<{ message: string }> {
    return this.http
      .post<ApiResponse<{ message: string }>>(`${this.apiUrl}/profile/preferences/extended`, data)
      .pipe(
        switchMap((response) => {
          if (response.success && response.data) {
            return of(response.data);
          }
          return throwError(() => response.error);
        }),
        catchError((err) => {
          this._error.set(err.message || 'Failed to set preference');
          return throwError(() => err);
        })
      );
  }

  /**
   * Delete an extended user preference
   * @param category - Preference category
   * @param key - Preference key
   */
  deleteExtendedPreference(category: string, key: string): Observable<void> {
    return this.http
      .delete<void>(`${this.apiUrl}/profile/preferences/extended`, {
        params: { category, key },
      })
      .pipe(
        catchError((err) => {
          this._error.set(err.message || 'Failed to delete preference');
          return throwError(() => err);
        })
      );
  }

  /**
   * Clear the current profile state
   */
  clearProfile(): void {
    this._profile.set(null);
    this._error.set(null);
  }

  /**
   * Clear any error state
   */
  clearError(): void {
    this._error.set(null);
  }
}
