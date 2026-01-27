import {
  Component,
  inject,
  computed,
  HostListener,
  ChangeDetectionStrategy,
  OnInit,
  OnDestroy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { SidebarComponent } from '../sidebar/sidebar.component';
import { HeaderComponent } from '../header/header.component';
import { LayoutService } from '../../core/services/layout.service';

/** Breakpoint for tablet view (below this is mobile) */
const TABLET_BREAKPOINT = 768;
/** Breakpoint for desktop view (above this sidebar is always visible) */
const DESKTOP_BREAKPOINT = 1024;

/**
 * MainLayoutComponent - Main application layout wrapper.
 *
 * Features:
 * - Responsive layout with sidebar, header, and content area
 * - CSS Grid/Flexbox layout
 * - Automatic sidebar collapse on tablet
 * - Router outlet for content
 * - Smooth transition animations
 *
 * Responsive Behavior:
 * - Desktop (>= 1024px): Sidebar always visible, can collapse
 * - Tablet (768px - 1023px): Sidebar collapsed by default
 * - Mobile (< 768px): Sidebar hidden, opens as overlay
 *
 * @example
 * ```html
 * <msls-main-layout>
 *   <!-- Child routes render here via router-outlet -->
 * </msls-main-layout>
 * ```
 */
@Component({
  selector: 'msls-main-layout',
  standalone: true,
  imports: [CommonModule, RouterOutlet, SidebarComponent, HeaderComponent],
  templateUrl: './main-layout.component.html',
  styleUrl: './main-layout.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MainLayoutComponent implements OnInit, OnDestroy {
  private readonly layoutService = inject(LayoutService);

  /** Whether sidebar is collapsed */
  readonly isCollapsed = this.layoutService.isCollapsed;

  /** Computed CSS classes for the layout */
  readonly layoutClasses = computed(() => {
    const classes: string[] = ['main-layout'];

    if (this.isCollapsed()) {
      classes.push('main-layout--collapsed');
    }

    return classes.join(' ');
  });

  ngOnInit(): void {
    // Set initial collapse state based on viewport
    this.handleResponsiveLayout();
  }

  ngOnDestroy(): void {
    // Clean up if needed
  }

  /** Handle responsive layout changes */
  @HostListener('window:resize')
  onWindowResize(): void {
    this.handleResponsiveLayout();
  }

  /**
   * Handle responsive layout based on viewport width.
   * - Desktop: Restore user's saved preference
   * - Tablet: Auto-collapse sidebar
   * - Mobile: Close overlay if open
   */
  private handleResponsiveLayout(): void {
    const width = window.innerWidth;

    if (width < DESKTOP_BREAKPOINT && width >= TABLET_BREAKPOINT) {
      // Tablet: Auto-collapse if not already collapsed
      if (!this.isCollapsed()) {
        this.layoutService.setCollapsed(true);
      }
    }

    // Mobile: Close sidebar overlay when resizing to larger viewport
    if (width >= DESKTOP_BREAKPOINT) {
      this.layoutService.close();
    }
  }
}
