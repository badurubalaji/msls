import { Component, input, computed } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Card variant types */
export type CardVariant = 'default' | 'elevated' | 'outlined';

/** Card padding sizes */
export type CardPadding = 'none' | 'sm' | 'md' | 'lg';

/**
 * MslsCardComponent - A versatile card container with header, body, and footer sections.
 *
 * Usage:
 * <msls-card variant="elevated" padding="lg">
 *   <ng-container card-header>Card Title</ng-container>
 *   <ng-container card-body>Card content goes here</ng-container>
 *   <ng-container card-footer>
 *     <button class="btn btn-primary">Action</button>
 *   </ng-container>
 * </msls-card>
 */
@Component({
  selector: 'msls-card',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div [class]="cardClasses()">
      <!-- Card Header (optional) -->
      <div [class]="headerClasses()">
        <ng-content select="[card-header]"></ng-content>
      </div>

      <!-- Card Body -->
      <div [class]="bodyClasses()">
        <ng-content select="[card-body]"></ng-content>
        <ng-content></ng-content>
      </div>

      <!-- Card Footer (optional) -->
      <div [class]="footerClasses()">
        <ng-content select="[card-footer]"></ng-content>
      </div>
    </div>
  `,
})
export class MslsCardComponent {
  /** Card style variant: default, elevated, or outlined */
  variant = input<CardVariant>('default');

  /** Padding size for the card content */
  padding = input<CardPadding>('md');

  /** Whether the card should have hover effect */
  hoverable = input<boolean>(false);

  /** Whether the card is clickable (adds cursor pointer) */
  clickable = input<boolean>(false);

  /** Base card classes */
  private readonly baseClasses = 'bg-white rounded-xl overflow-hidden relative';

  /** Variant classes */
  private readonly variantClasses: Record<CardVariant, string> = {
    default: 'border border-secondary-200 shadow-md shadow-primary-100/30',
    elevated: 'shadow-xl shadow-primary-200/40 border border-primary-100/50 ring-1 ring-primary-100/20',
    outlined: 'border-2 border-primary-200 shadow-sm',
  };

  /** Padding classes for sections */
  private readonly paddingClasses: Record<CardPadding, string> = {
    none: '',
    sm: 'px-4 py-3',
    md: 'px-6 py-4',
    lg: 'px-8 py-6',
  };

  /** Computed CSS classes for the card */
  cardClasses = computed(() => {
    const classes: string[] = [this.baseClasses];

    classes.push(this.variantClasses[this.variant()]);

    if (this.hoverable()) {
      classes.push('transition-all duration-200 hover:shadow-xl hover:-translate-y-0.5');
    }

    if (this.clickable()) {
      classes.push('cursor-pointer select-none');
    }

    return classes.join(' ');
  });

  /** Header classes */
  headerClasses = computed(() => {
    const padding = this.paddingClasses[this.padding()];
    return `border-b border-primary-100 bg-gradient-to-r from-primary-50/80 to-secondary-50 font-semibold text-secondary-900 empty:hidden empty:border-b-0 ${padding}`;
  });

  /** Body classes */
  bodyClasses = computed(() => {
    const padding = this.paddingClasses[this.padding()];
    return `text-secondary-700 ${padding}`;
  });

  /** Footer classes */
  footerClasses = computed(() => {
    const padding = this.paddingClasses[this.padding()];
    return `border-t border-primary-100 bg-gradient-to-r from-secondary-50 to-primary-50/50 empty:hidden empty:border-t-0 ${padding}`;
  });
}
