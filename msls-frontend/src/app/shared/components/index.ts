/**
 * MSLS Shared Components - Barrel Export
 *
 * This file exports all reusable components from the shared/components directory.
 * Import from '@shared/components' in your feature modules.
 */

// =============================================================================
// ATOMIC COMPONENTS (Story 1.5)
// =============================================================================

// Button Component
export { MslsButtonComponent } from './button/button.component';
export type { ButtonVariant, ButtonSize, ButtonType } from './button/button.component';

// Input Component
export { MslsInputComponent } from './input/input.component';
export type { InputType } from './input/input.component';

// Badge Component
export { MslsBadgeComponent } from './badge/badge.component';
export type { BadgeVariant, BadgeSize } from './badge/badge.component';

// Avatar Component
export { MslsAvatarComponent } from './avatar/avatar.component';
export type { AvatarSize } from './avatar/avatar.component';

// Icon Component
export { MslsIconComponent } from './icon/icon.component';
export type { IconSize, IconName } from './icon/icon.component';

// Spinner Component
export { MslsSpinnerComponent } from './spinner/spinner.component';
export type { SpinnerSize, SpinnerVariant } from './spinner/spinner.component';

// =============================================================================
// MOLECULE COMPONENTS (Story 1.6)
// =============================================================================

// Form Field Component
export { MslsFormFieldComponent } from './form-field/form-field.component';

// Card Component
export { MslsCardComponent } from './card/card.component';
export type { CardVariant, CardPadding } from './card/card.component';

// Dropdown Component
export { MslsDropdownComponent } from './dropdown/dropdown.component';
export type { DropdownTrigger, DropdownPosition } from './dropdown/dropdown.component';

// Toast Component
export { MslsToastComponent } from './toast/toast.component';

// Modal Component
export { MslsModalComponent } from './modal/modal.component';
export type { ModalSize } from './modal/modal.component';

// Select Component
export { MslsSelectComponent } from './select/select.component';
export type { SelectOption } from './select/select.component';

// Checkbox Component
export { MslsCheckboxComponent } from './checkbox/checkbox.component';

// Table Component
export { MslsTableComponent, MslsTableCellDirective } from './table/table.component';
export type { TableColumn, SortDirection, SortEvent } from './table/table.component';
