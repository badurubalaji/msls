import {
  Component,
  input,
  output,
  signal,
  computed,
  forwardRef,
  ElementRef,
  HostListener,
  inject
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ControlValueAccessor,
  NG_VALUE_ACCESSOR,
  FormsModule
} from '@angular/forms';

/** Option interface for select options */
export interface SelectOption {
  value: string | number;
  label: string;
  disabled?: boolean;
  group?: string;
}

/**
 * MslsSelectComponent - A custom select dropdown with search functionality.
 * Implements ControlValueAccessor for Angular forms integration.
 *
 * Usage:
 * <msls-select
 *   [options]="options"
 *   placeholder="Select an option"
 *   [searchable]="true"
 *   formControlName="mySelect"
 * ></msls-select>
 */
@Component({
  selector: 'msls-select',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './select.component.html',
  styleUrl: './select.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => MslsSelectComponent),
      multi: true,
    },
  ],
})
export class MslsSelectComponent implements ControlValueAccessor {
  private elementRef = inject(ElementRef);

  /** Available options */
  options = input<SelectOption[]>([]);

  /** Placeholder text */
  placeholder = input<string>('Select...');

  /** Allow multiple selection */
  multiple = input<boolean>(false);

  /** Enable search functionality */
  searchable = input<boolean>(false);

  /** Disabled state from input */
  disabled = input<boolean>(false);

  /** Whether the dropdown is currently open */
  isOpen = signal(false);

  /** Internal disabled state (from form control) */
  private _disabled = signal(false);

  /** Search query for filtering options */
  searchQuery = signal('');

  /** Internal value */
  private _value = signal<string | number | (string | number)[] | null>(null);

  /** Emitted when value changes */
  valueChange = output<string | number | (string | number)[] | null>();

  // ControlValueAccessor callbacks
  private onChange: (value: unknown) => void = () => {};
  private onTouched: () => void = () => {};

  /** Computed disabled state (input OR form control) */
  isDisabled = computed(() => this.disabled() || this._disabled());

  /** Computed filtered options based on search query */
  filteredOptions = computed(() => {
    const query = this.searchQuery().toLowerCase().trim();
    const opts = this.options();

    if (!query) {
      return opts;
    }

    return opts.filter(opt =>
      opt.label.toLowerCase().includes(query)
    );
  });

  /** Get display text for the selected value(s) */
  displayValue = computed(() => {
    const value = this._value();
    const opts = this.options();

    if (value === null || value === undefined) {
      return '';
    }

    if (this.multiple() && Array.isArray(value)) {
      const labels = value
        .map(v => opts.find(o => o.value === v)?.label)
        .filter(Boolean);
      return labels.join(', ');
    }

    return opts.find(o => o.value === value)?.label ?? '';
  });

  /** Check if an option is selected */
  isSelected(option: SelectOption): boolean {
    const value = this._value();

    if (value === null || value === undefined) {
      return false;
    }

    if (this.multiple() && Array.isArray(value)) {
      return value.includes(option.value);
    }

    return value === option.value;
  }

  /** Toggle dropdown open state */
  toggle(): void {
    if (this.isDisabled()) return;

    if (this.isOpen()) {
      this.close();
    } else {
      this.open();
    }
  }

  /** Open the dropdown */
  open(): void {
    if (this.isDisabled()) return;
    this.isOpen.set(true);
    this.searchQuery.set('');
  }

  /** Close the dropdown */
  close(): void {
    this.isOpen.set(false);
    this.searchQuery.set('');
    this.onTouched();
  }

  /** Select an option */
  selectOption(option: SelectOption): void {
    if (option.disabled) return;

    if (this.multiple()) {
      const currentValue = this._value();
      let newValue: (string | number)[];

      if (Array.isArray(currentValue)) {
        if (currentValue.includes(option.value)) {
          // Remove if already selected
          newValue = currentValue.filter(v => v !== option.value);
        } else {
          // Add if not selected
          newValue = [...currentValue, option.value];
        }
      } else {
        newValue = [option.value];
      }

      this.setValue(newValue);
    } else {
      this.setValue(option.value);
      this.close();
    }
  }

  /** Clear the selection */
  clear(event: Event): void {
    event.stopPropagation();
    this.setValue(this.multiple() ? [] : null);
  }

  /** Handle search input */
  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.searchQuery.set(target.value);
  }

  /** Set value and notify form control */
  private setValue(value: string | number | (string | number)[] | null): void {
    this._value.set(value);
    this.onChange(value);
    this.valueChange.emit(value);
  }

  /** Close dropdown when clicking outside */
  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.close();
    }
  }

  /** Handle keyboard navigation */
  @HostListener('keydown', ['$event'])
  onKeydown(event: KeyboardEvent): void {
    if (event.key === 'Escape') {
      this.close();
    } else if (event.key === 'Enter' || event.key === ' ') {
      if (!this.isOpen()) {
        event.preventDefault();
        this.open();
      }
    }
  }

  // ControlValueAccessor implementation
  writeValue(value: string | number | (string | number)[] | null): void {
    this._value.set(value);
  }

  registerOnChange(fn: (value: unknown) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this._disabled.set(isDisabled);
  }

  /** Track function for ngFor */
  trackByValue(index: number, option: SelectOption): string | number {
    return option.value;
  }
}
