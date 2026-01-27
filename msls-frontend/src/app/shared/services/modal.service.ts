import {
  Injectable,
  signal,
  Type,
  ApplicationRef,
  createComponent,
  EnvironmentInjector,
  inject,
  ComponentRef
} from '@angular/core';
import { Subject, Observable } from 'rxjs';

/** Modal size options */
export type ModalSize = 'sm' | 'md' | 'lg' | 'xl' | 'full';

/** Modal configuration options */
export interface ModalConfig {
  /** Modal size: sm, md, lg, xl, full */
  size?: ModalSize;
  /** Close when clicking the backdrop */
  closeOnBackdrop?: boolean;
  /** Close when pressing Escape key */
  closeOnEscape?: boolean;
  /** Data to pass to the modal component */
  data?: Record<string, unknown>;
  /** Custom CSS class for the modal */
  panelClass?: string;
  /** Whether the modal has a close button */
  showCloseButton?: boolean;
}

/** Reference to an open modal */
export interface ModalRef<T = unknown> {
  /** The component instance */
  instance: T;
  /** Close the modal with an optional result */
  close: (result?: unknown) => void;
  /** Observable that emits when the modal is closed */
  afterClosed$: Observable<unknown>;
}

/** Internal modal entry structure */
interface ModalEntry {
  id: string;
  componentRef: ComponentRef<unknown>;
  config: ModalConfig;
  closeSubject: Subject<unknown>;
}

/** Default modal configuration */
const DEFAULT_CONFIG: ModalConfig = {
  size: 'md',
  closeOnBackdrop: true,
  closeOnEscape: true,
  showCloseButton: true,
};

/**
 * ModalService - Service for programmatically opening modal dialogs.
 *
 * Usage:
 * constructor(private modalService: ModalService) {}
 *
 * openConfirmDialog() {
 *   const modalRef = this.modalService.open(ConfirmDialogComponent, {
 *     size: 'sm',
 *     data: { title: 'Confirm', message: 'Are you sure?' }
 *   });
 *
 *   modalRef.afterClosed$.subscribe(result => {
 *     if (result) { // User confirmed }
 *   });
 * }
 */
@Injectable({ providedIn: 'root' })
export class ModalService {
  private appRef = inject(ApplicationRef);
  private injector = inject(EnvironmentInjector);

  /** Internal signal holding all active modals */
  private modals = signal<ModalEntry[]>([]);

  /** Public readonly accessor for checking if any modals are open */
  readonly hasOpenModals = () => this.modals().length > 0;

  /** Get all open modals (readonly) */
  readonly openModals$ = this.modals.asReadonly();

  /**
   * Open a modal with a component
   * @param component - The component to render in the modal
   * @param config - Modal configuration
   * @returns A ModalRef to control the modal
   */
  open<T>(component: Type<T>, config?: ModalConfig): ModalRef<T> {
    const mergedConfig: ModalConfig = { ...DEFAULT_CONFIG, ...config };
    const id = this.generateId();
    const closeSubject = new Subject<unknown>();

    // Create the component dynamically
    const componentRef = createComponent(component, {
      environmentInjector: this.injector,
    });

    // Pass data to the component if provided
    if (mergedConfig.data && componentRef.instance) {
      const instance = componentRef.instance as Record<string, unknown>;
      Object.entries(mergedConfig.data).forEach(([key, value]) => {
        if (Object.prototype.hasOwnProperty.call(instance, key) || key in instance) {
          instance[key] = value;
        }
      });
    }

    // Attach the component to the application
    this.appRef.attachView(componentRef.hostView);

    // Create the modal entry
    const modalEntry: ModalEntry = {
      id,
      componentRef,
      config: mergedConfig,
      closeSubject,
    };

    // Add to the modals signal
    this.modals.update(modals => [...modals, modalEntry]);

    // Set up keyboard listener for Escape key
    if (mergedConfig.closeOnEscape) {
      this.setupEscapeKeyListener(id);
    }

    // Create the modal reference
    const modalRef: ModalRef<T> = {
      instance: componentRef.instance as T,
      close: (result?: unknown) => this.close(id, result),
      afterClosed$: closeSubject.asObservable(),
    };

    return modalRef;
  }

  /**
   * Close a specific modal by ID
   * @param id - The modal ID
   * @param result - Optional result to pass to subscribers
   */
  close(id: string, result?: unknown): void {
    const modal = this.modals().find(m => m.id === id);
    if (modal) {
      modal.closeSubject.next(result);
      modal.closeSubject.complete();

      // Destroy the component
      this.appRef.detachView(modal.componentRef.hostView);
      modal.componentRef.destroy();

      // Remove from the modals signal
      this.modals.update(modals => modals.filter(m => m.id !== id));
    }
  }

  /**
   * Close all open modals
   */
  closeAll(): void {
    const currentModals = this.modals();
    currentModals.forEach(modal => {
      modal.closeSubject.next(undefined);
      modal.closeSubject.complete();
      this.appRef.detachView(modal.componentRef.hostView);
      modal.componentRef.destroy();
    });
    this.modals.set([]);
  }

  /**
   * Get the topmost modal
   */
  getTopModal(): ModalEntry | undefined {
    const modals = this.modals();
    return modals.length > 0 ? modals[modals.length - 1] : undefined;
  }

  /**
   * Handle backdrop click for a modal
   */
  onBackdropClick(id: string): void {
    const modal = this.modals().find(m => m.id === id);
    if (modal?.config.closeOnBackdrop) {
      this.close(id);
    }
  }

  /**
   * Get modal config by ID
   */
  getModalConfig(id: string): ModalConfig | undefined {
    return this.modals().find(m => m.id === id)?.config;
  }

  /**
   * Get modal component ref by ID
   */
  getModalComponentRef(id: string): ComponentRef<unknown> | undefined {
    return this.modals().find(m => m.id === id)?.componentRef;
  }

  /**
   * Set up keyboard listener for Escape key
   */
  private setupEscapeKeyListener(id: string): void {
    const handler = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        const topModal = this.getTopModal();
        if (topModal?.id === id && topModal.config.closeOnEscape) {
          this.close(id);
          document.removeEventListener('keydown', handler);
        }
      }
    };
    document.addEventListener('keydown', handler);
  }

  /**
   * Generate a unique ID for a modal
   */
  private generateId(): string {
    return `modal-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
  }
}
