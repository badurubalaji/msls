import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Spinner size types */
export type SpinnerSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';

/** Spinner variant types */
export type SpinnerVariant = 'primary' | 'secondary' | 'white';

/**
 * MslsSpinnerComponent - Animated loading indicator.
 *
 * @example
 * ```html
 * <msls-spinner size="md" />
 * <msls-spinner size="lg" variant="white" label="Loading content..." />
 * ```
 */
@Component({
  selector: 'msls-spinner',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './spinner.component.html',
  styleUrl: './spinner.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MslsSpinnerComponent {
  /** Spinner size: xs, sm, md, lg, xl */
  readonly size = input<SpinnerSize>('md');

  /** Spinner variant: primary, secondary, white */
  readonly variant = input<SpinnerVariant>('primary');

  /** Accessible label for screen readers */
  readonly label = input<string>('Loading');

  /** Computed CSS classes based on size and variant */
  readonly spinnerClasses = computed(() => {
    const classes: string[] = ['msls-spinner'];

    // Size classes
    classes.push(`msls-spinner--${this.size()}`);

    // Variant classes
    classes.push(`msls-spinner--${this.variant()}`);

    return classes.join(' ');
  });
}
