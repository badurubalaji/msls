import {
  Component,
  input,
  output,
  signal,
  computed,
  ElementRef,
  AfterViewInit,
  OnDestroy,
  ViewChild,
  inject
} from '@angular/core';
import { CommonModule } from '@angular/common';

/** Modal size options */
export type ModalSize = 'sm' | 'md' | 'lg' | 'xl' | 'full';

/**
 * MslsModalComponent - A modal dialog with focus trapping and keyboard support.
 *
 * Usage (declarative):
 * <msls-modal [isOpen]="showModal" (closed)="showModal = false" size="md">
 *   <ng-container modal-header>Modal Title</ng-container>
 *   <ng-container modal-body>Modal content here</ng-container>
 *   <ng-container modal-footer>
 *     <button class="btn btn-primary" (click)="showModal = false">Close</button>
 *   </ng-container>
 * </msls-modal>
 *
 * For programmatic modals, use ModalService instead.
 */
@Component({
  selector: 'msls-modal',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './modal.component.html',
  styleUrl: './modal.component.scss'
})
export class MslsModalComponent implements AfterViewInit, OnDestroy {
  private elementRef = inject(ElementRef);

  /** Whether the modal is open */
  isOpen = input<boolean>(false);

  /** Modal size */
  size = input<ModalSize>('md');

  /** Close when clicking the backdrop */
  closeOnBackdrop = input<boolean>(true);

  /** Close when pressing Escape key */
  closeOnEscape = input<boolean>(true);

  /** Whether to show the close button */
  showCloseButton = input<boolean>(true);

  /** Modal title (alternative to slot) */
  title = input<string>('');

  /** Emitted when the modal should close */
  closed = output<void>();

  /** Reference to the modal content for focus trapping */
  @ViewChild('modalContent') modalContent!: ElementRef<HTMLElement>;

  /** Store the element that had focus before the modal opened */
  private previouslyFocusedElement: HTMLElement | null = null;

  /** All focusable elements in the modal */
  private focusableElements: HTMLElement[] = [];

  /** Computed size class */
  sizeClass = computed(() => `modal__panel--${this.size()}`);

  ngAfterViewInit(): void {
    if (this.isOpen()) {
      this.setupFocusTrap();
    }
  }

  ngOnDestroy(): void {
    this.restoreFocus();
  }

  /** Handle backdrop click */
  onBackdropClick(event: MouseEvent): void {
    // Only close if clicking directly on the backdrop, not the modal content
    if (event.target === event.currentTarget && this.closeOnBackdrop()) {
      this.close();
    }
  }

  /** Handle keyboard events */
  onKeydown(event: KeyboardEvent): void {
    if (event.key === 'Escape' && this.closeOnEscape()) {
      this.close();
      return;
    }

    // Focus trap
    if (event.key === 'Tab') {
      this.handleTabKey(event);
    }
  }

  /** Close the modal */
  close(): void {
    this.closed.emit();
    this.restoreFocus();
  }

  /** Set up focus trap when modal opens */
  setupFocusTrap(): void {
    // Store currently focused element
    this.previouslyFocusedElement = document.activeElement as HTMLElement;

    // Get all focusable elements
    setTimeout(() => {
      this.updateFocusableElements();
      // Focus the first focusable element or the modal itself
      if (this.focusableElements.length > 0) {
        this.focusableElements[0].focus();
      } else if (this.modalContent) {
        this.modalContent.nativeElement.focus();
      }
    });
  }

  /** Update the list of focusable elements */
  private updateFocusableElements(): void {
    if (!this.modalContent) return;

    const focusableSelectors = [
      'button:not([disabled])',
      'a[href]',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])',
    ].join(', ');

    this.focusableElements = Array.from(
      this.modalContent.nativeElement.querySelectorAll(focusableSelectors)
    ) as HTMLElement[];
  }

  /** Handle Tab key for focus trapping */
  private handleTabKey(event: KeyboardEvent): void {
    this.updateFocusableElements();

    if (this.focusableElements.length === 0) {
      event.preventDefault();
      return;
    }

    const firstElement = this.focusableElements[0];
    const lastElement = this.focusableElements[this.focusableElements.length - 1];

    if (event.shiftKey) {
      // Shift + Tab: go to last element if on first
      if (document.activeElement === firstElement) {
        event.preventDefault();
        lastElement.focus();
      }
    } else {
      // Tab: go to first element if on last
      if (document.activeElement === lastElement) {
        event.preventDefault();
        firstElement.focus();
      }
    }
  }

  /** Restore focus to the previously focused element */
  private restoreFocus(): void {
    if (this.previouslyFocusedElement) {
      this.previouslyFocusedElement.focus();
      this.previouslyFocusedElement = null;
    }
  }
}
