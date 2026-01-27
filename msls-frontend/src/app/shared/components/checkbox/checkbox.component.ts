import {
  Component,
  input,
  output,
  signal,
  forwardRef,
  computed
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

/**
 * MslsCheckboxComponent - A styled checkbox with label support.
 * Implements ControlValueAccessor for Angular forms integration.
 *
 * Usage:
 * <msls-checkbox
 *   label="Accept terms and conditions"
 *   formControlName="acceptTerms"
 * ></msls-checkbox>
 *
 * Or standalone:
 * <msls-checkbox
 *   label="Remember me"
 *   [(checked)]="rememberMe"
 * ></msls-checkbox>
 */
@Component({
  selector: 'msls-checkbox',
  standalone: true,
  imports: [CommonModule],
  template: `
    <label class="checkbox" [class.checkbox--disabled]="isDisabled()">
      <!-- Hidden native checkbox for accessibility -->
      <input
        type="checkbox"
        class="checkbox__input"
        [checked]="checked()"
        [disabled]="isDisabled()"
        [attr.aria-checked]="checked()"
        (change)="onInputChange($event)"
        (blur)="onBlur()"
      />

      <!-- Custom checkbox box -->
      <span
        class="checkbox__box"
        [class.checkbox__box--checked]="checked()"
        [class.checkbox__box--indeterminate]="indeterminate()"
      >
        @if (checked() && !indeterminate()) {
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        }
        @if (indeterminate()) {
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
          </svg>
        }
      </span>

      <!-- Label text -->
      @if (label()) {
        <span class="checkbox__label">{{ label() }}</span>
      }
      <ng-content></ng-content>
    </label>
  `,
  styleUrl: './checkbox.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => MslsCheckboxComponent),
      multi: true,
    },
  ],
})
export class MslsCheckboxComponent implements ControlValueAccessor {
  /** Label text for the checkbox */
  label = input<string>('');

  /** Disabled state from input */
  disabled = input<boolean>(false);

  /** Indeterminate state (for parent checkboxes) */
  indeterminate = input<boolean>(false);

  /** Internal checked state */
  private _checked = signal(false);

  /** Internal disabled state (from form control) */
  private _disabled = signal(false);

  /** Public checked signal */
  checked = this._checked.asReadonly();

  /** Emitted when checked state changes */
  checkedChange = output<boolean>();

  // ControlValueAccessor callbacks
  private onChange: (value: boolean) => void = () => {};
  private onTouched: () => void = () => {};

  /** Computed disabled state (input OR form control) */
  isDisabled = computed(() => this.disabled() || this._disabled());

  /** Handle input change */
  onInputChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.setChecked(target.checked);
  }

  /** Handle blur */
  onBlur(): void {
    this.onTouched();
  }

  /** Set checked state and notify */
  private setChecked(value: boolean): void {
    this._checked.set(value);
    this.onChange(value);
    this.checkedChange.emit(value);
  }

  /** Toggle checked state */
  toggle(): void {
    if (!this.isDisabled()) {
      this.setChecked(!this._checked());
    }
  }

  // ControlValueAccessor implementation
  writeValue(value: boolean): void {
    this._checked.set(!!value);
  }

  registerOnChange(fn: (value: boolean) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this._disabled.set(isDisabled);
  }
}
