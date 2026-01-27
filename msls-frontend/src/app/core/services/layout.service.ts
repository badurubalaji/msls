import { Injectable, signal, effect } from '@angular/core';

const SIDEBAR_COLLAPSED_KEY = 'sidebar_collapsed';

/**
 * LayoutService - Manages sidebar state and responsive layout behavior.
 *
 * Features:
 * - Sidebar collapse/expand state (persisted to localStorage)
 * - Mobile sidebar overlay open/close state
 * - Automatic persistence of collapse state
 *
 * @example
 * ```typescript
 * export class MyComponent {
 *   private layoutService = inject(LayoutService);
 *
 *   isCollapsed = this.layoutService.isCollapsed;
 *   isOpen = this.layoutService.isOpen;
 *
 *   toggleSidebar() {
 *     this.layoutService.toggleCollapse();
 *   }
 * }
 * ```
 */
@Injectable({ providedIn: 'root' })
export class LayoutService {
  /** Internal signal for sidebar collapsed state (desktop) */
  private readonly _sidebarCollapsed = signal<boolean>(false);

  /** Internal signal for sidebar open state (mobile overlay) */
  private readonly _sidebarOpen = signal<boolean>(false);

  /** Readonly signal - whether sidebar is collapsed (desktop view) */
  readonly isCollapsed = this._sidebarCollapsed.asReadonly();

  /** Readonly signal - whether sidebar is open (mobile overlay view) */
  readonly isOpen = this._sidebarOpen.asReadonly();

  constructor() {
    // Load saved collapse state from localStorage
    this.loadSavedState();

    // Persist collapse state changes to localStorage
    effect(() => {
      const collapsed = this._sidebarCollapsed();
      localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(collapsed));
    });
  }

  /**
   * Toggle sidebar collapsed state (for desktop view).
   * State is automatically persisted to localStorage.
   */
  toggleCollapse(): void {
    this._sidebarCollapsed.update((collapsed) => !collapsed);
  }

  /**
   * Set sidebar collapsed state directly.
   * @param collapsed - Whether sidebar should be collapsed
   */
  setCollapsed(collapsed: boolean): void {
    this._sidebarCollapsed.set(collapsed);
  }

  /**
   * Toggle sidebar open state (for mobile overlay view).
   */
  toggleOpen(): void {
    this._sidebarOpen.update((open) => !open);
  }

  /**
   * Open sidebar (mobile overlay).
   */
  open(): void {
    this._sidebarOpen.set(true);
  }

  /**
   * Close sidebar (mobile overlay).
   * Call this when navigating or clicking outside on mobile.
   */
  close(): void {
    this._sidebarOpen.set(false);
  }

  /**
   * Close sidebar on mobile devices.
   * Useful for closing the sidebar when a navigation item is clicked.
   */
  closeOnMobile(): void {
    if (this._sidebarOpen()) {
      this._sidebarOpen.set(false);
    }
  }

  /**
   * Load saved sidebar state from localStorage.
   */
  private loadSavedState(): void {
    try {
      const saved = localStorage.getItem(SIDEBAR_COLLAPSED_KEY);
      if (saved !== null) {
        this._sidebarCollapsed.set(saved === 'true');
      }
    } catch {
      // localStorage not available (e.g., SSR or privacy mode)
      console.warn('LayoutService: Could not read from localStorage');
    }
  }
}
