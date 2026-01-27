# Story 1.6: Design System - Molecule Components

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As a **developer**,
I want **reusable molecule components (FormField, Card, Dropdown, Toast, Modal)**,
So that **complex UI patterns are consistent across the application**.

## Acceptance Criteria

### FormField Component
**Given** a developer needs a form field with label and validation
**When** they use the MslsFormField component
**Then** it wraps input with label, hint text, and error message
**And** error message appears only when field is touched and invalid
**And** required indicator (*) appears for required fields

### Card Component
**Given** a developer needs a card container
**When** they use the MslsCard component
**Then** it supports header, body, footer sections
**And** it supports variants: default, elevated, outlined

### Toast Notifications
**Given** a developer needs notifications
**When** they use the ToastService
**Then** toasts appear in top-right corner
**And** variants exist: success (green), error (red), warning (amber), info (blue)
**And** toasts auto-dismiss after configurable duration
**And** toasts can be manually dismissed

### Modal Component
**Given** a developer needs a modal dialog
**When** they use the MslsModal component
**Then** it supports sizes: sm, md, lg, xl
**And** it traps focus within modal
**And** it closes on Escape key or backdrop click (configurable)

## Technical Requirements

### Component Location
`msls-frontend/src/app/shared/components/`

### Component Structure
```
shared/components/
├── form-field/
│   ├── form-field.component.ts
│   ├── form-field.component.html
│   └── form-field.component.spec.ts
├── card/
│   ├── card.component.ts
│   ├── card.component.html
│   └── card.component.spec.ts
├── dropdown/
│   ├── dropdown.component.ts
│   ├── dropdown.component.html
│   └── dropdown.component.spec.ts
├── toast/
│   ├── toast.component.ts
│   ├── toast.component.html
│   ├── toast.service.ts
│   └── toast.component.spec.ts
├── modal/
│   ├── modal.component.ts
│   ├── modal.component.html
│   ├── modal.service.ts
│   └── modal.component.spec.ts
├── select/
│   ├── select.component.ts
│   ├── select.component.html
│   └── select.component.spec.ts
├── checkbox/
│   ├── checkbox.component.ts
│   └── checkbox.component.spec.ts
└── table/
    ├── table.component.ts
    ├── table.component.html
    └── table.component.spec.ts
```

### FormField Component API
```typescript
@Component({
  selector: 'msls-form-field',
  standalone: true,
})
export class MslsFormFieldComponent {
  label = input<string>('');
  hint = input<string>('');
  error = input<string>('');
  required = input<boolean>(false);
}
```

### Card Component API
```typescript
@Component({
  selector: 'msls-card',
  standalone: true,
})
export class MslsCardComponent {
  variant = input<'default' | 'elevated' | 'outlined'>('default');
  padding = input<'none' | 'sm' | 'md' | 'lg'>('md');
}

// With content projection:
// <msls-card>
//   <ng-container card-header>Title</ng-container>
//   <ng-container card-body>Content</ng-container>
//   <ng-container card-footer>Actions</ng-container>
// </msls-card>
```

### Toast Service API
```typescript
@Injectable({ providedIn: 'root' })
export class ToastService {
  private toasts = signal<Toast[]>([]);
  readonly toasts$ = this.toasts.asReadonly();

  success(message: string, duration?: number): void;
  error(message: string, duration?: number): void;
  warning(message: string, duration?: number): void;
  info(message: string, duration?: number): void;
  dismiss(id: string): void;
  dismissAll(): void;
}
```

### Modal Service API
```typescript
@Injectable({ providedIn: 'root' })
export class ModalService {
  open<T>(component: Type<T>, config?: ModalConfig): ModalRef<T>;
  closeAll(): void;
}

interface ModalConfig {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  closeOnBackdrop?: boolean;
  closeOnEscape?: boolean;
  data?: any;
}

interface ModalRef<T> {
  instance: T;
  close(result?: any): void;
  afterClosed$: Observable<any>;
}
```

### Dropdown Component API
```typescript
@Component({
  selector: 'msls-dropdown',
  standalone: true,
})
export class MslsDropdownComponent {
  trigger = input.required<'click' | 'hover'>('click');
  position = input<'bottom-start' | 'bottom-end' | 'top-start' | 'top-end'>('bottom-start');

  isOpen = signal(false);
}
```

### Select Component API
```typescript
@Component({
  selector: 'msls-select',
  standalone: true,
})
export class MslsSelectComponent implements ControlValueAccessor {
  options = input<SelectOption[]>([]);
  placeholder = input<string>('Select...');
  multiple = input<boolean>(false);
  searchable = input<boolean>(false);
  disabled = input<boolean>(false);
}
```

## Tasks

1. [ ] Generate form-field component with ng generate
2. [ ] Implement label, hint, error display
3. [ ] Add required indicator styling
4. [ ] Generate card component with ng generate
5. [ ] Implement card with header/body/footer slots
6. [ ] Add card variants (default, elevated, outlined)
7. [ ] Generate toast component with ng generate
8. [ ] Create ToastService with signal-based state
9. [ ] Implement toast container positioning
10. [ ] Add toast auto-dismiss logic
11. [ ] Generate modal component with ng generate
12. [ ] Create ModalService for programmatic modals
13. [ ] Implement focus trap in modal
14. [ ] Add keyboard navigation (Escape to close)
15. [ ] Generate dropdown component with ng generate
16. [ ] Implement dropdown positioning
17. [ ] Generate select component with ng generate
18. [ ] Implement ControlValueAccessor for select
19. [ ] Add search functionality to select
20. [ ] Generate checkbox component with ng generate
21. [ ] Generate table component with ng generate
22. [ ] Update barrel export (index.ts)
23. [ ] Write unit tests for all components
24. [ ] Write tests for services

## Definition of Done

- [ ] All molecule components created using ng generate
- [ ] Components follow Angular 21 standalone patterns
- [ ] Components use Signals for state
- [ ] Toast and Modal services work correctly
- [ ] Form components implement ControlValueAccessor
- [ ] Focus trap works in modal
- [ ] Unit tests pass with >80% coverage
- [ ] Components exported from shared module
