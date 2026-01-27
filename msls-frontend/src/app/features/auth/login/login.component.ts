/**
 * LoginComponent - User authentication login page for MSLS.
 *
 * Features:
 * - Email/password authentication
 * - Remember me functionality
 * - Password visibility toggle
 * - Form validation with error messages
 * - Loading state during submission
 * - Error handling for various auth failures
 * - Redirect to originally requested URL after login
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  signal,
  computed,
  OnInit,
  OnDestroy,
  effect,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router, RouterLink, ActivatedRoute } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';

import { AuthService } from '../../../core/services/auth.service';
import { TenantService } from '../../../core/services/tenant.service';
import {
  MslsInputComponent,
  MslsCheckboxComponent,
} from '../../../shared/components';

/**
 * LoginComponent - Authentication login page.
 *
 * Provides a professional login form with:
 * - Email and password validation
 * - Remember me option
 * - Password visibility toggle
 * - Error message display
 * - Loading state management
 * - Redirect handling after successful login
 */
@Component({
  selector: 'msls-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    MslsInputComponent,
    MslsCheckboxComponent,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent implements OnInit, OnDestroy {
  private readonly authService = inject(AuthService);
  private readonly tenantService = inject(TenantService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly fb = inject(FormBuilder);

  /** Destroy subject for subscription cleanup */
  private readonly destroy$ = new Subject<void>();

  /** Login form with email, password, and rememberMe fields */
  loginForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(8)]],
    rememberMe: [false],
  });

  /** Loading state during form submission */
  isLoading = signal(false);

  /** Effect to disable/enable form based on loading state */
  private readonly loadingEffect = effect(() => {
    if (this.isLoading()) {
      this.loginForm.disable();
    } else {
      this.loginForm.enable();
    }
  });

  /** Error message to display */
  error = signal<string | null>(null);

  /** Password visibility toggle state */
  showPassword = signal(false);

  /** Return URL for post-login redirect */
  private returnUrl = '/dashboard';

  /** Current year for copyright display */
  readonly currentYear = new Date().getFullYear();

  /** Computed error messages for form fields */
  emailError = computed(() => {
    const control = this.loginForm.get('email');
    if (control?.touched && control?.errors) {
      if (control.errors['required']) {
        return 'Email is required';
      }
      if (control.errors['email']) {
        return 'Please enter a valid email address';
      }
    }
    return '';
  });

  passwordError = computed(() => {
    const control = this.loginForm.get('password');
    if (control?.touched && control?.errors) {
      if (control.errors['required']) {
        return 'Password is required';
      }
      if (control.errors['minlength']) {
        return 'Password must be at least 8 characters';
      }
    }
    return '';
  });

  /** Computed state for form validity */
  isFormValid = computed(() => {
    return this.loginForm.valid && !this.isLoading();
  });

  ngOnInit(): void {
    // Get return URL from route params
    this.route.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params) => {
      this.returnUrl = params['returnUrl'] || '/dashboard';
    });

    // Clear error when form changes
    this.loginForm.valueChanges.pipe(takeUntil(this.destroy$)).subscribe(() => {
      if (this.error()) {
        this.error.set(null);
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /**
   * Toggle password visibility
   */
  togglePasswordVisibility(): void {
    this.showPassword.update((show) => !show);
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    // Mark all fields as touched to show validation errors
    this.loginForm.markAllAsTouched();

    if (this.loginForm.invalid) {
      return;
    }

    this.isLoading.set(true);
    this.error.set(null);

    const { email, password, rememberMe } = this.loginForm.value;

    // Use demo tenant for development
    // TODO: Implement proper tenant resolution (e.g., from subdomain or login form)
    // Note: We use a hardcoded tenant ID here because stored tenant might be stale after DB reset
    const tenantId = '61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4';

    this.authService
      .login({
        email: email!,
        password: password!,
        tenantId,
        rememberMe: rememberMe ?? false,
      })
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          // Navigate to the return URL after successful login
          this.router.navigateByUrl(this.returnUrl);
        },
        error: (error: HttpErrorResponse | Error) => {
          this.isLoading.set(false);
          this.handleLoginError(error);
        },
      });
  }

  /**
   * Handle login errors and display appropriate messages
   */
  private handleLoginError(error: HttpErrorResponse | Error): void {
    if (error instanceof HttpErrorResponse) {
      switch (error.status) {
        case 401:
          this.error.set('Invalid email or password. Please try again.');
          break;
        case 423:
          this.error.set(
            'Your account has been locked. Please contact support.'
          );
          break;
        case 429:
          this.error.set('Too many login attempts. Please try again later.');
          break;
        case 0:
          this.error.set(
            'Unable to connect to the server. Please check your network connection.'
          );
          break;
        default:
          this.error.set(
            error.error?.error?.detail ||
              error.error?.message ||
              'An error occurred during login. Please try again.'
          );
      }
    } else {
      this.error.set(error.message || 'An unexpected error occurred.');
    }
  }
}
