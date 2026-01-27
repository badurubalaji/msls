import {
  Component,
  input,
  output,
  signal,
  computed,
  forwardRef,
  ChangeDetectionStrategy,
  ElementRef,
  viewChild,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ControlValueAccessor,
  NG_VALUE_ACCESSOR,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { MslsIconComponent } from '../icon/icon.component';

/** Input type attribute */
export type InputType = 'text' | 'email' | 'password' | 'number' | 'tel' | 'search' | 'date';

/**
 * MslsInputComponent - Reusable input component with ControlValueAccessor support.
 *
 * @example
 * ```html
 * <msls-input
 *   label="Email"
 *   type="email"
 *   placeholder="Enter your email"
 *   [error]="emailError()"
 *   prefixIcon="envelope"
 *   [(ngModel)]="email"
 * />
 * ```
 */
@Component({
  selector: 'msls-input',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, MslsIconComponent],
  templateUrl: './input.component.html',
  styleUrl: './input.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => MslsInputComponent),
      multi: true,
    },
  ],
})
export class MslsInputComponent implements ControlValueAccessor {
  /** Input type: text, email, password, number, tel, search */
  readonly type = input<InputType>('text');

  /** Placeholder text */
  readonly placeholder = input<string>('');

  /** Label text displayed above the input */
  readonly label = input<string>('');

  /** Hint text displayed below the input */
  readonly hint = input<string>('');

  /** Error message displayed below the input */
  readonly error = input<string>('');

  /** Prefix icon name (Heroicons) */
  readonly prefixIcon = input<string>('');

  /** Suffix icon name (Heroicons) */
  readonly suffixIcon = input<string>('');

  /** Disabled state */
  readonly disabled = input<boolean>(false);

  /** Required field indicator */
  readonly required = input<boolean>(false);

  /** Readonly state */
  readonly readonly = input<boolean>(false);

  /** Autocomplete attribute */
  readonly autocomplete = input<string>('off');

  /** Input ID (auto-generated if not provided) */
  readonly inputId = input<string>(`msls-input-${Math.random().toString(36).substring(2, 9)}`);

  /** Maximum length */
  readonly maxlength = input<number | null>(null);

  /** Minimum length */
  readonly minlength = input<number | null>(null);

  /** Pattern for validation */
  readonly pattern = input<string>('');

  /** Reference to the native input element */
  readonly inputRef = viewChild<ElementRef<HTMLInputElement>>('inputElement');

  /** Internal value signal */
  readonly value = signal<string>('');

  /** Internal disabled state from ControlValueAccessor */
  readonly isDisabledByForm = signal<boolean>(false);

  /** Computed disabled state */
  readonly isDisabled = computed(() => this.disabled() || this.isDisabledByForm());

  /** Computed error state */
  readonly hasError = computed(() => !!this.error());

  /** Computed container classes */
  readonly containerClasses = computed(() => {
    const classes: string[] = ['msls-input'];

    if (this.hasError()) {
      classes.push('msls-input--error');
    }

    if (this.isDisabled()) {
      classes.push('msls-input--disabled');
    }

    if (this.prefixIcon()) {
      classes.push('msls-input--has-prefix');
    }

    if (this.suffixIcon()) {
      classes.push('msls-input--has-suffix');
    }

    return classes.join(' ');
  });

  // ControlValueAccessor callbacks
  private onChange: (value: string) => void = () => {};
  private onTouched: () => void = () => {};

  /** Write value from form control */
  writeValue(value: string): void {
    this.value.set(value ?? '');
  }

  /** Register onChange callback */
  registerOnChange(fn: (value: string) => void): void {
    this.onChange = fn;
  }

  /** Register onTouched callback */
  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  /** Set disabled state from form control */
  setDisabledState(isDisabled: boolean): void {
    this.isDisabledByForm.set(isDisabled);
  }

  /** Handle input event */
  onInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.value.set(target.value);
    this.onChange(target.value);
  }

  /** Handle blur event */
  onBlur(): void {
    this.onTouched();
  }

  /** Focus the input programmatically */
  focus(): void {
    this.inputRef()?.nativeElement.focus();
  }

  /** Blur the input programmatically */
  blur(): void {
    this.inputRef()?.nativeElement.blur();
  }
}
