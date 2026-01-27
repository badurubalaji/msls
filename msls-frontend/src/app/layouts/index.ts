/**
 * MSLS Layouts Module - Barrel Export
 *
 * This file exports all layout components for the application shell.
 * Import from '@layouts' in your modules.
 */

// =============================================================================
// MAIN LAYOUT COMPONENTS
// =============================================================================

// Main Layout - Application shell wrapper
export { MainLayoutComponent } from './main-layout/main-layout.component';

// Sidebar - Navigation sidebar
export { SidebarComponent } from './sidebar/sidebar.component';

// Header - Application header with user menu
export { HeaderComponent } from './header/header.component';

// User Menu - User profile dropdown
export { UserMenuComponent } from './user-menu/user-menu.component';

// Nav Item - Single navigation item
export { NavItemComponent } from './nav-item/nav-item.component';

// =============================================================================
// NAVIGATION CONFIGURATION
// =============================================================================

export { NAV_ITEMS } from './nav-config';
export type { NavItem } from './nav-config';
