import { Component, inject, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MslsIconComponent } from '../../shared/components/icon/icon.component';
import { UserMenuComponent } from '../user-menu/user-menu.component';
import { LayoutService } from '../../core/services/layout.service';
import { TenantService } from '../../core/services/tenant.service';

/**
 * HeaderComponent - Application header with navigation controls and user menu.
 *
 * Features:
 * - Logo/brand area
 * - Hamburger menu for mobile (toggles sidebar)
 * - School/tenant name display
 * - User menu component integration
 * - Responsive layout
 *
 * @example
 * ```html
 * <msls-header />
 * ```
 */
@Component({
  selector: 'msls-header',
  standalone: true,
  imports: [CommonModule, MslsIconComponent, UserMenuComponent],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HeaderComponent {
  private readonly layoutService = inject(LayoutService);
  private readonly tenantService = inject(TenantService);

  /** Current tenant information */
  readonly currentTenant = this.tenantService.currentTenant;

  /** Tenant name to display */
  readonly tenantName = computed(() => {
    const tenant = this.currentTenant();
    return tenant?.name ?? 'School Management System';
  });

  /** Toggle mobile sidebar */
  toggleMobileSidebar(): void {
    this.layoutService.toggleOpen();
  }
}
