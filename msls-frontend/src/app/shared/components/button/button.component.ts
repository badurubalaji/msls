import { Component, input, output, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Button variant types */
export type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'outline';

/** Button size types */
export type ButtonSize = 'sm' | 'md' | 'lg';

/** Button type attribute */
export type ButtonType = 'button' | 'submit' | 'reset';

/**
 * MslsButtonComponent - Reusable button component with multiple variants and sizes.
 *
 * @example
 * ```html
 * <msls-button variant="primary" size="md">Click me</msls-button>
 * <msls-button variant="danger" [loading]="true">Saving...</msls-button>
 * ```
 */
@Component({
  selector: 'msls-button',
  standalone: true,
  imports: [CommonModule],
  template: `
    <button
      [type]="type()"
      [disabled]="isDisabled()"
      [class]="buttonClasses()"
      [attr.aria-busy]="loading() ? 'true' : null"
      [attr.aria-disabled]="isDisabled() ? 'true' : null"
      (click)="onClick($event)"
    >
      @if (loading()) {
        <span class="absolute inset-0 flex items-center justify-center">
          <svg
            class="h-5 w-5 animate-spin"
            viewBox="0 0 24 24"
            fill="none"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="3"
            />
            <path
              class="opacity-75"
              d="M12 2C6.47715 2 2 6.47715 2 12"
              stroke="currentColor"
              stroke-width="3"
              stroke-linecap="round"
            />
          </svg>
        </span>
      }
      <span [class]="loading() ? 'invisible' : 'inline-flex items-center gap-2'">
        <ng-content />
      </span>
    </button>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MslsButtonComponent {
  /** Button variant: primary, secondary, danger, ghost, outline */
  readonly variant = input<ButtonVariant>('primary');

  /** Button size: sm, md, lg */
  readonly size = input<ButtonSize>('md');

  /** Shows loading spinner and disables button */
  readonly loading = input<boolean>(false);

  /** Disables the button */
  readonly disabled = input<boolean>(false);

  /** HTML button type attribute */
  readonly type = input<ButtonType>('button');

  /** Full width button */
  readonly fullWidth = input<boolean>(false);

  /** Icon only button (square aspect ratio) */
  readonly iconOnly = input<boolean>(false);

  /** Emitted when the button is clicked */
  readonly clicked = output<MouseEvent>();

  /** Computed disabled state (disabled or loading) */
  readonly isDisabled = computed(() => this.disabled() || this.loading());

  /** Handle button click */
  onClick(event: MouseEvent): void {
    if (!this.isDisabled()) {
      this.clicked.emit(event);
    }
  }

  /** Base classes for all buttons */
  private readonly baseClasses = 'relative inline-flex items-center justify-center font-medium rounded-lg border-2 border-transparent transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';

  /** Size variant classes */
  private readonly sizeClasses: Record<ButtonSize, string> = {
    sm: 'h-8 px-3 text-xs gap-1.5',
    md: 'h-10 px-4 text-sm gap-2',
    lg: 'h-12 px-6 text-base gap-2.5',
  };

  /** Icon-only size classes */
  private readonly iconOnlySizeClasses: Record<ButtonSize, string> = {
    sm: 'h-8 w-8 p-0',
    md: 'h-10 w-10 p-0',
    lg: 'h-12 w-12 p-0',
  };

  /** Variant classes */
  private readonly variantClasses: Record<ButtonVariant, string> = {
    primary: 'bg-gradient-to-r from-primary-600 to-primary-700 text-white hover:from-primary-700 hover:to-primary-800 active:from-primary-800 active:to-primary-900 focus:ring-primary-500 shadow-md shadow-primary-500/30 hover:shadow-lg hover:shadow-primary-500/40',
    secondary: 'bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-800 hover:from-secondary-200 hover:to-secondary-300 active:from-secondary-300 active:to-secondary-400 focus:ring-secondary-400 shadow-sm',
    danger: 'bg-gradient-to-r from-danger-500 to-danger-600 text-white hover:from-danger-600 hover:to-danger-700 active:from-danger-700 active:to-danger-800 focus:ring-danger-500 shadow-md shadow-danger-500/30',
    ghost: 'bg-transparent text-primary-600 hover:bg-primary-50 hover:text-primary-700 active:bg-primary-100 focus:ring-primary-400',
    outline: 'bg-transparent border-primary-500 text-primary-600 hover:bg-primary-50 hover:border-primary-600 active:bg-primary-100 focus:ring-primary-500',
  };

  /** Computed CSS classes based on variant, size, and state */
  readonly buttonClasses = computed(() => {
    const classes: string[] = [this.baseClasses];

    // Size classes
    if (this.iconOnly()) {
      classes.push(this.iconOnlySizeClasses[this.size()]);
    } else {
      classes.push(this.sizeClasses[this.size()]);
    }

    // Variant classes
    classes.push(this.variantClasses[this.variant()]);

    // Full width
    if (this.fullWidth()) {
      classes.push('w-full');
    }

    return classes.join(' ');
  });
}
