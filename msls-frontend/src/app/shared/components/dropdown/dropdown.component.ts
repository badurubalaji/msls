import {
  Component,
  input,
  signal,
  output,
  ElementRef,
  inject,
  OnDestroy,
  HostListener,
  computed
} from '@angular/core';
import { CommonModule } from '@angular/common';

/** Dropdown trigger types */
export type DropdownTrigger = 'click' | 'hover';

/** Dropdown position types */
export type DropdownPosition = 'bottom-start' | 'bottom-end' | 'top-start' | 'top-end';

/**
 * MslsDropdownComponent - A dropdown menu with configurable trigger and positioning.
 *
 * Usage:
 * <msls-dropdown trigger="click" position="bottom-start">
 *   <button dropdown-trigger class="btn btn-secondary">Open Menu</button>
 *   <div dropdown-content>
 *     <a href="#">Menu Item 1</a>
 *     <a href="#">Menu Item 2</a>
 *   </div>
 * </msls-dropdown>
 */
@Component({
  selector: 'msls-dropdown',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dropdown.component.html',
  styleUrl: './dropdown.component.scss'
})
export class MslsDropdownComponent implements OnDestroy {
  private elementRef = inject(ElementRef);

  /** How the dropdown is triggered: click or hover */
  trigger = input<DropdownTrigger>('click');

  /** Position of the dropdown content relative to trigger */
  position = input<DropdownPosition>('bottom-start');

  /** Whether the dropdown is currently open */
  isOpen = signal(false);

  /** Emitted when dropdown opens */
  opened = output<void>();

  /** Emitted when dropdown closes */
  closed = output<void>();

  /** Computed position classes */
  positionClasses = computed(() => {
    const pos = this.position();
    return {
      'dropdown__content--bottom-start': pos === 'bottom-start',
      'dropdown__content--bottom-end': pos === 'bottom-end',
      'dropdown__content--top-start': pos === 'top-start',
      'dropdown__content--top-end': pos === 'top-end',
    };
  });

  private hoverTimeout: ReturnType<typeof setTimeout> | null = null;

  /** Toggle dropdown open state */
  toggle(): void {
    if (this.isOpen()) {
      this.close();
    } else {
      this.open();
    }
  }

  /** Open the dropdown */
  open(): void {
    this.isOpen.set(true);
    this.opened.emit();
  }

  /** Close the dropdown */
  close(): void {
    this.isOpen.set(false);
    this.closed.emit();
  }

  /** Handle trigger click */
  onTriggerClick(): void {
    if (this.trigger() === 'click') {
      this.toggle();
    }
  }

  /** Handle mouse enter on dropdown */
  onMouseEnter(): void {
    if (this.trigger() === 'hover') {
      if (this.hoverTimeout) {
        clearTimeout(this.hoverTimeout);
        this.hoverTimeout = null;
      }
      this.open();
    }
  }

  /** Handle mouse leave on dropdown */
  onMouseLeave(): void {
    if (this.trigger() === 'hover') {
      this.hoverTimeout = setTimeout(() => {
        this.close();
      }, 150);
    }
  }

  /** Close dropdown when clicking outside */
  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.close();
    }
  }

  /** Close dropdown on Escape key */
  @HostListener('document:keydown.escape')
  onEscapeKey(): void {
    if (this.isOpen()) {
      this.close();
    }
  }

  ngOnDestroy(): void {
    if (this.hoverTimeout) {
      clearTimeout(this.hoverTimeout);
    }
  }
}
