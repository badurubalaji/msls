import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Badge variant types */
export type BadgeVariant = 'primary' | 'success' | 'warning' | 'danger' | 'error' | 'info' | 'neutral';

/** Badge size types */
export type BadgeSize = 'sm' | 'md';

/**
 * MslsBadgeComponent - Reusable badge component for status indicators and labels.
 *
 * @example
 * ```html
 * <msls-badge variant="success">Active</msls-badge>
 * <msls-badge variant="warning" size="sm">Pending</msls-badge>
 * ```
 */
@Component({
  selector: 'msls-badge',
  standalone: true,
  imports: [CommonModule],
  template: `
    <span [class]="badgeClasses()">
      @if (dot()) {
        <span [class]="dotClasses()"></span>
      }
      <span><ng-content /></span>
    </span>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MslsBadgeComponent {
  /** Badge variant: success, warning, error, info, neutral */
  readonly variant = input<BadgeVariant>('neutral');

  /** Badge size: sm, md */
  readonly size = input<BadgeSize>('md');

  /** Show dot indicator before text */
  readonly dot = input<boolean>(false);

  /** Base classes */
  private readonly baseClasses = 'inline-flex items-center gap-1.5 font-medium rounded-full whitespace-nowrap';

  /** Size classes */
  private readonly sizeClasses: Record<BadgeSize, string> = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-2.5 py-1 text-xs',
  };

  /** Variant classes */
  private readonly variantClasses: Record<BadgeVariant, string> = {
    primary: 'bg-gradient-to-r from-primary-100 to-primary-200 text-primary-800 ring-1 ring-inset ring-primary-500/30 shadow-sm shadow-primary-500/10',
    success: 'bg-gradient-to-r from-success-100 to-success-200 text-success-800 ring-1 ring-inset ring-success-500/30 shadow-sm shadow-success-500/10',
    warning: 'bg-gradient-to-r from-warning-100 to-warning-200 text-warning-800 ring-1 ring-inset ring-warning-500/30 shadow-sm shadow-warning-500/10',
    danger: 'bg-gradient-to-r from-danger-100 to-danger-200 text-danger-800 ring-1 ring-inset ring-danger-500/30 shadow-sm shadow-danger-500/10',
    error: 'bg-gradient-to-r from-danger-100 to-danger-200 text-danger-800 ring-1 ring-inset ring-danger-500/30 shadow-sm shadow-danger-500/10',
    info: 'bg-gradient-to-r from-accent-100 to-accent-200 text-accent-800 ring-1 ring-inset ring-accent-500/30 shadow-sm shadow-accent-500/10',
    neutral: 'bg-gradient-to-r from-secondary-100 to-secondary-200 text-secondary-800 ring-1 ring-inset ring-secondary-500/30 shadow-sm',
  };

  /** Dot color classes */
  private readonly dotColorClasses: Record<BadgeVariant, string> = {
    primary: 'bg-primary-500',
    success: 'bg-success-500',
    warning: 'bg-warning-500',
    danger: 'bg-danger-500',
    error: 'bg-danger-500',
    info: 'bg-accent-500',
    neutral: 'bg-secondary-500',
  };

  /** Computed CSS classes based on variant and size */
  readonly badgeClasses = computed(() => {
    return `${this.baseClasses} ${this.sizeClasses[this.size()]} ${this.variantClasses[this.variant()]}`;
  });

  /** Computed dot classes */
  readonly dotClasses = computed(() => {
    return `h-1.5 w-1.5 rounded-full ${this.dotColorClasses[this.variant()]}`;
  });
}
