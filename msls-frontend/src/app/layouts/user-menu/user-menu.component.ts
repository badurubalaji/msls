import { Component, inject, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { MslsIconComponent } from '../../shared/components/icon/icon.component';
import { MslsAvatarComponent } from '../../shared/components/avatar/avatar.component';
import { MslsDropdownComponent } from '../../shared/components/dropdown/dropdown.component';
import { AuthService } from '../../core/services/auth.service';

/**
 * UserMenuComponent - User profile dropdown menu.
 *
 * Features:
 * - User avatar display (image or initials)
 * - User name and email display
 * - Dropdown menu with: Profile, Settings, Logout
 * - Logout functionality with redirect
 * - Responsive design
 *
 * @example
 * ```html
 * <msls-user-menu />
 * ```
 */
@Component({
  selector: 'msls-user-menu',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MslsIconComponent,
    MslsAvatarComponent,
    MslsDropdownComponent,
  ],
  templateUrl: './user-menu.component.html',
  styleUrl: './user-menu.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserMenuComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  /** Current authenticated user */
  readonly currentUser = this.authService.currentUser;

  /** User display name */
  readonly displayName = this.authService.displayName;

  /** User full name for avatar */
  readonly userFullName = computed(() => {
    const user = this.currentUser();
    if (!user) return '';
    return `${user.firstName} ${user.lastName}`.trim();
  });

  /** User email */
  readonly userEmail = computed(() => {
    return this.currentUser()?.email ?? '';
  });

  /** User initials for avatar fallback */
  readonly userInitials = computed(() => {
    const name = this.userFullName();
    if (!name) return '';

    const parts = name.split(/\s+/);
    if (parts.length === 1) {
      return parts[0].charAt(0).toUpperCase();
    }
    return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
  });

  /** Navigate to profile page */
  goToProfile(): void {
    this.router.navigate(['/profile']);
  }

  /** Navigate to settings page */
  goToSettings(): void {
    this.router.navigate(['/settings']);
  }

  /** Log out the current user */
  logout(): void {
    this.authService.logout(true).subscribe({
      complete: () => {
        // Logout observable completes after clearing tokens
        // Navigation to login is handled by the service
      },
    });
  }
}
