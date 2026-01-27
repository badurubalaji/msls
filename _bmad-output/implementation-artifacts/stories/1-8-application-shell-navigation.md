# Story 1.8: Application Shell & Navigation Layout

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As an **admin user**,
I want **a responsive application shell with sidebar navigation**,
So that **I can navigate between different modules easily**.

## Acceptance Criteria

### Desktop View (>1024px)
**Given** an admin user is logged in
**When** they view the application on desktop
**Then** they see a collapsible sidebar with navigation menu
**And** sidebar shows module icons and labels
**And** current route is highlighted in sidebar
**And** sidebar collapse state persists across sessions

### Mobile View (<768px)
**Given** an admin user is on mobile
**When** they view the application
**Then** sidebar is hidden by default
**And** hamburger menu icon in header toggles sidebar
**And** sidebar overlays content when open

### Application Shell
**Given** the application shell
**When** rendered
**Then** header includes school logo, tenant name, user menu
**And** user menu includes profile, settings, logout options
**And** main content area has proper spacing and scroll behavior

## Technical Requirements

### Component Location
`msls-frontend/src/app/layouts/`

### File Structure
```
layouts/
├── main-layout/
│   ├── main-layout.component.ts
│   ├── main-layout.component.html
│   └── main-layout.component.scss
├── sidebar/
│   ├── sidebar.component.ts
│   ├── sidebar.component.html
│   └── sidebar.component.scss
├── header/
│   ├── header.component.ts
│   ├── header.component.html
│   └── header.component.scss
├── user-menu/
│   ├── user-menu.component.ts
│   └── user-menu.component.html
└── nav-item/
    ├── nav-item.component.ts
    └── nav-item.component.html
```

### Navigation Configuration
```typescript
interface NavItem {
  label: string;
  icon: string;
  route?: string;
  children?: NavItem[];
  permissions?: string[];
}

const NAV_ITEMS: NavItem[] = [
  { label: 'Dashboard', icon: 'home', route: '/dashboard' },
  { label: 'Students', icon: 'users', route: '/students' },
  { label: 'Staff', icon: 'briefcase', route: '/staff' },
  { label: 'Academics', icon: 'book-open', route: '/academics', children: [...] },
  { label: 'Finance', icon: 'dollar-sign', route: '/finance', children: [...] },
  { label: 'Settings', icon: 'cog', route: '/settings' },
];
```

### Layout Service
```typescript
@Injectable({ providedIn: 'root' })
export class LayoutService {
  private sidebarCollapsed = signal(false);
  private sidebarOpen = signal(false); // For mobile overlay

  readonly isCollapsed = this.sidebarCollapsed.asReadonly();
  readonly isOpen = this.sidebarOpen.asReadonly();

  toggleCollapse(): void;
  toggleOpen(): void;
  closeOnMobile(): void;
}
```

### Responsive Breakpoints
- Desktop: >= 1024px (sidebar always visible, can collapse)
- Tablet: 768px - 1023px (sidebar collapsed by default)
- Mobile: < 768px (sidebar hidden, opens as overlay)

### Styling Guidelines
- Sidebar width: 256px expanded, 64px collapsed
- Header height: 64px
- Use smooth transitions for collapse/expand
- Use z-index properly for overlay on mobile

## Tasks

1. [ ] Generate main-layout component with ng generate
2. [ ] Generate sidebar component with ng generate
3. [ ] Generate header component with ng generate
4. [ ] Generate user-menu component with ng generate
5. [ ] Generate nav-item component with ng generate
6. [ ] Create LayoutService for sidebar state
7. [ ] Implement sidebar with navigation items
8. [ ] Add collapse/expand functionality
9. [ ] Implement responsive behavior
10. [ ] Add mobile overlay mode
11. [ ] Implement header with logo and user menu
12. [ ] Create user menu dropdown
13. [ ] Add active route highlighting
14. [ ] Persist sidebar state in localStorage
15. [ ] Add keyboard navigation support
16. [ ] Update app.routes.ts with layout wrapper
17. [ ] Write unit tests for components

## Definition of Done

- [ ] All layout components created using ng generate
- [ ] Sidebar collapses on desktop
- [ ] Sidebar works as overlay on mobile
- [ ] Active route is highlighted
- [ ] Sidebar state persists across sessions
- [ ] User menu works correctly
- [ ] Responsive at all breakpoints
- [ ] Unit tests pass
