import {
  Component,
  inject,
  computed,
  HostListener,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { NavItemComponent } from '../nav-item/nav-item.component';
import { NAV_ITEMS, NavItem } from '../nav-config';
import { LayoutService } from '../../core/services/layout.service';
import { RbacService } from '../../core/services/rbac.service';

/**
 * SidebarComponent - Navigation sidebar for the application.
 *
 * Features:
 * - Collapsible on desktop (256px <-> 64px)
 * - Overlay mode on mobile
 * - Navigation items with nested children support
 * - Active route highlighting
 * - Smooth transition animations
 * - Keyboard navigation support
 *
 * @example
 * ```html
 * <msls-sidebar />
 * ```
 */
@Component({
  selector: 'msls-sidebar',
  standalone: true,
  imports: [CommonModule, NavItemComponent],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  readonly layoutService = inject(LayoutService);
  private readonly rbacService = inject(RbacService);

  /** Navigation items filtered by permissions */
  readonly navItems = computed(() => this.filterNavItems(NAV_ITEMS));

  /**
   * Filter nav items based on user permissions
   */
  private filterNavItems(items: NavItem[]): NavItem[] {
    return items
      .filter(item => this.hasPermissionForItem(item))
      .map(item => {
        if (item.children) {
          return {
            ...item,
            children: this.filterNavItems(item.children),
          };
        }
        return item;
      })
      .filter(item => !item.children || item.children.length > 0);
  }

  /**
   * Check if user has permission for a nav item
   */
  private hasPermissionForItem(item: NavItem): boolean {
    if (!item.permissions || item.permissions.length === 0) {
      return true;
    }
    // User needs any of the specified permissions
    return this.rbacService.hasAnyPermission(item.permissions);
  }

  /** Whether sidebar is collapsed (desktop) */
  readonly isCollapsed = this.layoutService.isCollapsed;

  /** Whether sidebar is open (mobile overlay) */
  readonly isOpen = this.layoutService.isOpen;

  /** Computed CSS classes for the sidebar */
  readonly sidebarClasses = computed(() => {
    const classes: string[] = ['sidebar'];

    if (this.isCollapsed()) {
      classes.push('sidebar--collapsed');
    }

    if (this.isOpen()) {
      classes.push('sidebar--open');
    }

    return classes.join(' ');
  });

  /** Toggle sidebar collapsed state */
  toggleCollapse(): void {
    this.layoutService.toggleCollapse();
  }

  /** Close sidebar on mobile */
  closeMobile(): void {
    this.layoutService.close();
  }

  /** Close mobile sidebar on Escape key */
  @HostListener('document:keydown.escape')
  onEscapeKey(): void {
    if (this.isOpen()) {
      this.closeMobile();
    }
  }
}
