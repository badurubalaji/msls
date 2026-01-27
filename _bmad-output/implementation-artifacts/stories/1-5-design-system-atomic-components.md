# Story 1.5: Design System - Atomic Components

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As a **developer**,
I want **reusable atomic UI components (Button, Input, Badge, Avatar, Icon)**,
So that **all features have consistent visual styling and behavior**.

## Acceptance Criteria

### Button Component
**Given** a developer needs a button in their feature
**When** they use the MslsButton component
**Then** it supports variants: primary, secondary, danger, ghost, outline
**And** it supports sizes: sm, md, lg
**And** it supports states: default, hover, focus, disabled, loading
**And** it meets WCAG 2.1 AA contrast requirements
**And** focus ring is visible for keyboard navigation

### Input Component
**Given** a developer needs an input field
**When** they use the MslsInput component
**Then** it supports types: text, email, password, number, tel, search
**And** it displays error state with message below
**And** it supports prefix/suffix icons
**And** it meets minimum 44x44px touch target on mobile

### Badge & Avatar Components
**Given** a developer needs a badge or avatar
**When** they use MslsBadge or MslsAvatar components
**Then** Badge supports variants: success, warning, error, info, neutral
**And** Avatar supports sizes and fallback initials

## Technical Requirements

### Component Location
`msls-frontend/src/app/shared/components/`

### Component Structure
```
shared/components/
├── button/
│   ├── button.component.ts
│   ├── button.component.html
│   ├── button.component.scss
│   └── button.component.spec.ts
├── input/
│   ├── input.component.ts
│   ├── input.component.html
│   ├── input.component.scss
│   └── input.component.spec.ts
├── badge/
│   ├── badge.component.ts
│   └── badge.component.spec.ts
├── avatar/
│   ├── avatar.component.ts
│   └── avatar.component.spec.ts
├── icon/
│   ├── icon.component.ts
│   └── icon.component.spec.ts
├── spinner/
│   ├── spinner.component.ts
│   └── spinner.component.spec.ts
└── index.ts
```

### Button Component API
```typescript
@Component({
  selector: 'msls-button',
  standalone: true,
})
export class MslsButtonComponent {
  variant = input<'primary' | 'secondary' | 'danger' | 'ghost' | 'outline'>('primary');
  size = input<'sm' | 'md' | 'lg'>('md');
  loading = input<boolean>(false);
  disabled = input<boolean>(false);
  type = input<'button' | 'submit' | 'reset'>('button');
}
```

### Input Component API
```typescript
@Component({
  selector: 'msls-input',
  standalone: true,
})
export class MslsInputComponent implements ControlValueAccessor {
  type = input<'text' | 'email' | 'password' | 'number' | 'tel' | 'search'>('text');
  placeholder = input<string>('');
  label = input<string>('');
  hint = input<string>('');
  error = input<string>('');
  prefixIcon = input<string>('');
  suffixIcon = input<string>('');
  disabled = input<boolean>(false);
  required = input<boolean>(false);
}
```

### Badge Component API
```typescript
@Component({
  selector: 'msls-badge',
  standalone: true,
})
export class MslsBadgeComponent {
  variant = input<'success' | 'warning' | 'error' | 'info' | 'neutral'>('neutral');
  size = input<'sm' | 'md'>('md');
}
```

### Avatar Component API
```typescript
@Component({
  selector: 'msls-avatar',
  standalone: true,
})
export class MslsAvatarComponent {
  src = input<string>('');
  name = input<string>('');
  size = input<'sm' | 'md' | 'lg' | 'xl'>('md');
}
```

### Icon Component
Use Heroicons (outline) via inline SVG or icon library.

### Styling Guidelines
- Use Tailwind CSS custom theme colors
- Use CSS custom properties for dynamic theming
- Ensure keyboard focus is visible (ring-2)
- Support dark mode (future)

## Tasks

1. [ ] Generate button component with ng generate
2. [ ] Implement button variants and sizes
3. [ ] Add loading spinner to button
4. [ ] Generate input component with ng generate
5. [ ] Implement ControlValueAccessor for input
6. [ ] Add prefix/suffix icon support to input
7. [ ] Add error state styling to input
8. [ ] Generate badge component with ng generate
9. [ ] Implement badge variants
10. [ ] Generate avatar component with ng generate
11. [ ] Implement avatar with image and fallback initials
12. [ ] Generate icon component with ng generate
13. [ ] Set up Heroicons integration
14. [ ] Generate spinner component
15. [ ] Create barrel export (index.ts)
16. [ ] Write unit tests for all components
17. [ ] Verify WCAG AA compliance

## Definition of Done

- [ ] All atomic components created using ng generate
- [ ] Components follow Angular 21 standalone patterns
- [ ] Components use Signals for inputs
- [ ] Input component implements ControlValueAccessor
- [ ] All variants and sizes work correctly
- [ ] Focus states visible for keyboard navigation
- [ ] Unit tests pass with >80% coverage
- [ ] Components exported from shared module
