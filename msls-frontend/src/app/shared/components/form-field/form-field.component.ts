import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * MslsFormFieldComponent - Wraps form inputs with label, hint text, and error messages.
 *
 * Usage:
 * <msls-form-field label="Email" hint="Enter your email" [required]="true" [error]="emailError()">
 *   <input class="input" type="email" />
 * </msls-form-field>
 */
@Component({
  selector: 'msls-form-field',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './form-field.component.html',
  styleUrl: './form-field.component.scss'
})
export class MslsFormFieldComponent {
  /** The label text displayed above the input */
  label = input<string>('');

  /** Hint text displayed below the input (when no error) */
  hint = input<string>('');

  /** Error message to display (only shown when provided) */
  error = input<string>('');

  /** Whether the field is required (shows asterisk indicator) */
  required = input<boolean>(false);

  /** Custom ID for the form field (for accessibility) */
  fieldId = input<string>('');
}
