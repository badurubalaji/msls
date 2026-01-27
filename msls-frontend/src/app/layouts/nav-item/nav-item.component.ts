import {
  Component,
  input,
  signal,
  computed,
  inject,
  ChangeDetectionStrategy,
  HostListener,
  ElementRef,
  ViewChild,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { NavItem } from '../nav-config';
import { LayoutService } from '../../core/services/layout.service';

/**
 * NavItemComponent - Single navigation item for the sidebar.
 *
 * Features:
 * - Icon and label display
 * - RouterLinkActive for highlighting active routes
 * - Expandable children support
 * - Collapse-aware display (icon-only when sidebar collapsed)
 * - Flyout menu for collapsed parent items
 * - Keyboard navigation support
 */
@Component({
  selector: 'msls-nav-item',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './nav-item.component.html',
  styleUrl: './nav-item.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavItemComponent {
  private readonly router = inject(Router);
  private readonly layoutService = inject(LayoutService);

  @ViewChild('parentButton') parentButton!: ElementRef<HTMLButtonElement>;

  /** Navigation item configuration */
  readonly item = input.required<NavItem>();

  /** Whether the sidebar is collapsed (icon-only mode) */
  readonly collapsed = input<boolean>(false);

  /** Nesting level for indentation */
  readonly level = input<number>(0);

  /** Internal state for expanded children */
  readonly isExpanded = signal<boolean>(false);

  /** Whether to show the flyout menu (collapsed state) */
  readonly showFlyout = signal<boolean>(false);

  /** Flyout position for fixed positioning */
  readonly flyoutPosition = signal<{ top: number; left: number }>({ top: 0, left: 0 });

  /** Hover state for parent button */
  private parentHovered = signal<boolean>(false);

  /** Hover state for flyout menu */
  private flyoutHovered = signal<boolean>(false);

  /** Icon mapping from custom names to Font Awesome classes */
  private readonly iconMap: Record<string, string> = {
    // Dashboard & Home
    'home': 'fa-solid fa-house',
    'dashboard': 'fa-solid fa-gauge-high',
    'chart-bar': 'fa-solid fa-chart-bar',

    // Users & People
    'users': 'fa-solid fa-users',
    'user': 'fa-solid fa-user',
    'user-group': 'fa-solid fa-user-group',
    'academic-cap': 'fa-solid fa-graduation-cap',
    'user-circle': 'fa-solid fa-user-circle',

    // Admin & Settings
    'cog': 'fa-solid fa-cog',
    'cog-6-tooth': 'fa-solid fa-gear',
    'adjustments-horizontal': 'fa-solid fa-sliders',
    'shield-check': 'fa-solid fa-shield-halved',
    'key': 'fa-solid fa-key',
    'lock-closed': 'fa-solid fa-lock',
    'flag': 'fa-solid fa-flag',

    // Documents & Content
    'document-text': 'fa-solid fa-file-lines',
    'document': 'fa-solid fa-file',
    'clipboard-document-list': 'fa-solid fa-clipboard-list',
    'folder': 'fa-solid fa-folder',
    'folder-open': 'fa-solid fa-folder-open',

    // Calendar & Time
    'calendar': 'fa-solid fa-calendar',
    'calendar-days': 'fa-solid fa-calendar-days',
    'clock': 'fa-solid fa-clock',

    // Money & Finance
    'currency-dollar': 'fa-solid fa-dollar-sign',
    'banknotes': 'fa-solid fa-money-bill',
    'credit-card': 'fa-solid fa-credit-card',

    // Communication
    'chat-bubble-left': 'fa-solid fa-comment',
    'envelope': 'fa-solid fa-envelope',
    'bell': 'fa-solid fa-bell',

    // Other
    'book-open': 'fa-solid fa-book-open',
    'building-office': 'fa-solid fa-building',
    'truck': 'fa-solid fa-truck',
    'beaker': 'fa-solid fa-flask',
    'presentation-chart-line': 'fa-solid fa-chart-line',
    'question-mark-circle': 'fa-solid fa-circle-question',
    'information-circle': 'fa-solid fa-circle-info',

    // Arrows & Navigation
    'chevron-down': 'fa-solid fa-chevron-down',
    'chevron-right': 'fa-solid fa-chevron-right',
    'chevron-left': 'fa-solid fa-chevron-left',
    'arrow-right': 'fa-solid fa-arrow-right',
  };

  /** Whether this item has children */
  readonly hasChildren = computed(() => {
    const children = this.item().children;
    return children && children.length > 0;
  });

  /** Whether this item is a direct link (has route, no children) */
  readonly isLink = computed(() => {
    return !!this.item().route && !this.hasChildren();
  });

  /** Computed CSS classes for the nav item */
  readonly itemClasses = computed(() => {
    const classes: string[] = ['nav-item'];

    if (this.collapsed()) {
      classes.push('nav-item--collapsed');
    }

    if (this.hasChildren()) {
      classes.push('nav-item--parent');
    }

    if (this.isExpanded()) {
      classes.push('nav-item--expanded');
    }

    if (this.level() > 0) {
      classes.push(`nav-item--level-${this.level()}`);
    }

    return classes.join(' ');
  });

  /** Get Font Awesome icon class from icon name */
  getIconClass(iconName: string | undefined): string {
    if (!iconName) return 'fa-solid fa-circle';
    return this.iconMap[iconName] || `fa-solid fa-${iconName}`;
  }

  /** Toggle expanded state for parent items */
  toggleExpanded(): void {
    if (this.hasChildren()) {
      this.isExpanded.update((expanded) => !expanded);
    }
  }

  /** Handle click on nav item */
  onItemClick(event: Event): void {
    const item = this.item();

    if (this.hasChildren()) {
      // Parent item - toggle expansion (when not collapsed)
      event.preventDefault();
      if (!this.collapsed()) {
        this.toggleExpanded();
      }
      // When collapsed, clicking does nothing - flyout handles navigation
    } else if (item.route) {
      // Link item - navigate and close mobile sidebar
      this.layoutService.closeOnMobile();
    }
  }

  /** Handle parent button hover (for collapsed flyout) */
  onParentHover(isHovering: boolean, event?: MouseEvent): void {
    if (this.collapsed()) {
      this.parentHovered.set(isHovering);
      if (isHovering && event) {
        this.updateFlyoutPosition(event.currentTarget as HTMLElement);
      }
      this.updateFlyoutVisibility();
    }
  }

  /** Update flyout position based on button element */
  private updateFlyoutPosition(buttonElement: HTMLElement): void {
    const rect = buttonElement.getBoundingClientRect();
    this.flyoutPosition.set({
      top: rect.top,
      left: rect.right + 8, // 8px gap
    });
  }

  /** Handle flyout hover */
  onFlyoutHover(isHovering: boolean): void {
    this.flyoutHovered.set(isHovering);
    this.updateFlyoutVisibility();
  }

  /** Update flyout visibility based on hover states */
  private updateFlyoutVisibility(): void {
    // Small delay to allow mouse movement between parent and flyout
    setTimeout(() => {
      this.showFlyout.set(this.parentHovered() || this.flyoutHovered());
    }, 50);
  }

  /** Handle child item click in flyout */
  onChildClick(): void {
    this.showFlyout.set(false);
    this.layoutService.closeOnMobile();
  }

  /** Handle keyboard navigation */
  @HostListener('keydown', ['$event'])
  onKeyDown(event: KeyboardEvent): void {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      this.onItemClick(event);
    }

    if (this.hasChildren() && !this.collapsed()) {
      if (event.key === 'ArrowRight' && !this.isExpanded()) {
        event.preventDefault();
        this.isExpanded.set(true);
      } else if (event.key === 'ArrowLeft' && this.isExpanded()) {
        event.preventDefault();
        this.isExpanded.set(false);
      }
    }
  }

  /** Check if any child route is active */
  isChildActive(): boolean {
    const children = this.item().children;
    if (!children) return false;

    return children.some((child) => {
      if (child.route) {
        return this.router.isActive(child.route, {
          paths: 'subset',
          queryParams: 'ignored',
          fragment: 'ignored',
          matrixParams: 'ignored',
        });
      }
      return false;
    });
  }
}
