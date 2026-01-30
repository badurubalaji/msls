import { Component, input, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Icon size types */
export type IconSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';

/** Available icon names */
export type IconName =
  | 'academic-cap'
  | 'adjustments-horizontal'
  | 'adjustments-vertical'
  | 'archive-box'
  | 'arrow-down'
  | 'arrow-left'
  | 'arrow-path'
  | 'arrow-right'
  | 'arrow-right-on-rectangle'
  | 'arrow-up'
  | 'bars-3'
  | 'bell'
  | 'bolt'
  | 'bell-alert'
  | 'book-open'
  | 'bookmark'
  | 'briefcase'
  | 'building-office'
  | 'calendar'
  | 'calendar-days'
  | 'banknotes'
  | 'chart-bar'
  | 'chart-bar-square'
  | 'check'
  | 'check-circle'
  | 'chevron-down'
  | 'chevron-left'
  | 'chevron-right'
  | 'chevron-up'
  | 'clipboard'
  | 'clipboard-document-list'
  | 'clipboard-document-check'
  | 'clock'
  | 'cog'
  | 'currency-dollar'
  | 'document'
  | 'document-text'
  | 'ellipsis-horizontal'
  | 'ellipsis-vertical'
  | 'envelope'
  | 'exclamation-circle'
  | 'exclamation-triangle'
  | 'eye'
  | 'eye-slash'
  | 'flag'
  | 'funnel'
  | 'globe-alt'
  | 'heart'
  | 'home'
  | 'identification'
  | 'information-circle'
  | 'key'
  | 'lock-closed'
  | 'lock-open'
  | 'login'
  | 'logout'
  | 'magnifying-glass'
  | 'map-pin'
  | 'minus'
  | 'pencil'
  | 'phone'
  | 'photo'
  | 'plus'
  | 'question-mark-circle'
  | 'shield-check'
  | 'star'
  | 'trophy'
  | 'award'
  | 'list-check'
  | 'file-pen'
  | 'user-plus'
  | 'comment-dots'
  | 'pen-to-square'
  | 'trash'
  | 'user'
  | 'user-circle'
  | 'user-group'
  | 'users'
  | 'x-circle'
  | 'x-mark'
  | string; // Allow custom icon names

/**
 * MslsIconComponent - Reusable icon component using Font Awesome.
 *
 * @example
 * ```html
 * <msls-icon name="user" size="md" />
 * <msls-icon name="check-circle" size="lg" class="text-success-500" />
 * ```
 */
@Component({
  selector: 'msls-icon',
  standalone: true,
  imports: [CommonModule],
  template: `
    <i
      [class]="iconClasses()"
      [attr.aria-label]="label() || null"
      [attr.aria-hidden]="isDecorative()"
      role="img"
    ></i>
  `,
  styles: [`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }

    i {
      line-height: 1;
    }

    .msls-icon--xs { font-size: 0.75rem; }
    .msls-icon--sm { font-size: 1rem; }
    .msls-icon--md { font-size: 1.25rem; }
    .msls-icon--lg { font-size: 1.5rem; }
    .msls-icon--xl { font-size: 2rem; }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MslsIconComponent {
  /** Icon name */
  readonly name = input.required<IconName>();

  /** Icon size: xs, sm, md, lg, xl */
  readonly size = input<IconSize>('md');

  /** Accessible label for screen readers */
  readonly label = input<string>('');

  /** Computed CSS classes based on size and icon name */
  readonly iconClasses = computed(() => {
    const faClass = this.getFontAwesomeClass(this.name());
    const sizeClass = `msls-icon--${this.size()}`;
    return `${faClass} ${sizeClass}`;
  });

  /** Whether icon is decorative (no label) */
  readonly isDecorative = computed(() => !this.label());

  /** Map icon names to Font Awesome classes */
  private readonly iconMap: Record<string, string> = {
    'academic-cap': 'fa-solid fa-graduation-cap',
    'adjustments-horizontal': 'fa-solid fa-sliders',
    'adjustments-vertical': 'fa-solid fa-sliders',
    'archive-box': 'fa-solid fa-box-archive',
    'arrow-down': 'fa-solid fa-arrow-down',
    'arrow-left': 'fa-solid fa-arrow-left',
    'arrow-path': 'fa-solid fa-rotate',
    'arrow-right': 'fa-solid fa-arrow-right',
    'arrow-right-on-rectangle': 'fa-solid fa-right-from-bracket',
    'arrow-up': 'fa-solid fa-arrow-up',
    'bars-3': 'fa-solid fa-bars',
    'bell': 'fa-regular fa-bell',
    'bolt': 'fa-solid fa-bolt',
    'bell-alert': 'fa-solid fa-bell',
    'book-open': 'fa-solid fa-book-open',
    'bookmark': 'fa-regular fa-bookmark',
    'briefcase': 'fa-solid fa-briefcase',
    'building-office': 'fa-solid fa-building',
    'calendar': 'fa-regular fa-calendar',
    'calendar-days': 'fa-solid fa-calendar-days',
    'banknotes': 'fa-solid fa-money-bills',
    'chart-bar': 'fa-solid fa-chart-bar',
    'chart-bar-square': 'fa-solid fa-chart-column',
    'check': 'fa-solid fa-check',
    'check-circle': 'fa-solid fa-circle-check',
    'chevron-down': 'fa-solid fa-chevron-down',
    'chevron-left': 'fa-solid fa-chevron-left',
    'chevron-right': 'fa-solid fa-chevron-right',
    'chevron-up': 'fa-solid fa-chevron-up',
    'clipboard': 'fa-regular fa-clipboard',
    'clipboard-document-list': 'fa-solid fa-clipboard-list',
    'clipboard-document-check': 'fa-solid fa-clipboard-check',
    'document-check': 'fa-solid fa-file-circle-check',
    'file-signature': 'fa-solid fa-file-signature',
    'rectangle-stack': 'fa-solid fa-layer-group',
    'arrows-right-left': 'fa-solid fa-arrows-left-right',
    'table-cells': 'fa-solid fa-table-cells',
    'clock': 'fa-regular fa-clock',
    'cog': 'fa-solid fa-gear',
    'currency-dollar': 'fa-solid fa-dollar-sign',
    'document': 'fa-regular fa-file',
    'document-text': 'fa-regular fa-file-lines',
    'ellipsis-horizontal': 'fa-solid fa-ellipsis',
    'ellipsis-vertical': 'fa-solid fa-ellipsis-vertical',
    'envelope': 'fa-regular fa-envelope',
    'exclamation-circle': 'fa-solid fa-circle-exclamation',
    'exclamation-triangle': 'fa-solid fa-triangle-exclamation',
    'eye': 'fa-regular fa-eye',
    'eye-slash': 'fa-regular fa-eye-slash',
    'flag': 'fa-regular fa-flag',
    'funnel': 'fa-solid fa-filter',
    'globe-alt': 'fa-solid fa-globe',
    'heart': 'fa-regular fa-heart',
    'home': 'fa-solid fa-house',
    'identification': 'fa-solid fa-id-badge',
    'information-circle': 'fa-solid fa-circle-info',
    'key': 'fa-solid fa-key',
    'lock-closed': 'fa-solid fa-lock',
    'lock-open': 'fa-solid fa-lock-open',
    'login': 'fa-solid fa-right-to-bracket',
    'logout': 'fa-solid fa-right-from-bracket',
    'magnifying-glass': 'fa-solid fa-magnifying-glass',
    'map-pin': 'fa-solid fa-location-dot',
    'minus': 'fa-solid fa-minus',
    'pencil': 'fa-solid fa-pen',
    'phone': 'fa-solid fa-phone',
    'photo': 'fa-regular fa-image',
    'plus': 'fa-solid fa-plus',
    'question-mark-circle': 'fa-solid fa-circle-question',
    'shield-check': 'fa-solid fa-shield-halved',
    'star': 'fa-regular fa-star',
    'trophy': 'fa-solid fa-trophy',
    'award': 'fa-solid fa-award',
    'list-check': 'fa-solid fa-list-check',
    'file-pen': 'fa-solid fa-file-pen',
    'user-plus': 'fa-solid fa-user-plus',
    'comment-dots': 'fa-solid fa-comment-dots',
    'pen-to-square': 'fa-solid fa-pen-to-square',
    'trash': 'fa-regular fa-trash-can',
    'user': 'fa-regular fa-user',
    'user-circle': 'fa-regular fa-circle-user',
    'user-group': 'fa-solid fa-user-group',
    'users': 'fa-solid fa-users',
    'x-circle': 'fa-solid fa-circle-xmark',
    'x-mark': 'fa-solid fa-xmark',
  };

  /** Get Font Awesome class for the given icon name */
  private getFontAwesomeClass(name: string): string {
    // Check if it's a direct Font Awesome class (starts with fa-)
    if (name.startsWith('fa-')) {
      return name;
    }
    // Otherwise, look up in the map
    return this.iconMap[name] || 'fa-solid fa-question';
  }
}
