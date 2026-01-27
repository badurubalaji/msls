/**
 * Address Form Component
 *
 * Reusable address form for student current/permanent addresses.
 */

import {
  Component,
  ChangeDetectionStrategy,
  input,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';

import {
  MslsInputComponent,
  MslsSelectComponent,
  MslsFormFieldComponent,
  SelectOption,
} from '../../../../shared/components';

/** Indian states list */
const INDIAN_STATES: SelectOption[] = [
  { value: '', label: 'Select State' },
  { value: 'Andhra Pradesh', label: 'Andhra Pradesh' },
  { value: 'Arunachal Pradesh', label: 'Arunachal Pradesh' },
  { value: 'Assam', label: 'Assam' },
  { value: 'Bihar', label: 'Bihar' },
  { value: 'Chhattisgarh', label: 'Chhattisgarh' },
  { value: 'Goa', label: 'Goa' },
  { value: 'Gujarat', label: 'Gujarat' },
  { value: 'Haryana', label: 'Haryana' },
  { value: 'Himachal Pradesh', label: 'Himachal Pradesh' },
  { value: 'Jharkhand', label: 'Jharkhand' },
  { value: 'Karnataka', label: 'Karnataka' },
  { value: 'Kerala', label: 'Kerala' },
  { value: 'Madhya Pradesh', label: 'Madhya Pradesh' },
  { value: 'Maharashtra', label: 'Maharashtra' },
  { value: 'Manipur', label: 'Manipur' },
  { value: 'Meghalaya', label: 'Meghalaya' },
  { value: 'Mizoram', label: 'Mizoram' },
  { value: 'Nagaland', label: 'Nagaland' },
  { value: 'Odisha', label: 'Odisha' },
  { value: 'Punjab', label: 'Punjab' },
  { value: 'Rajasthan', label: 'Rajasthan' },
  { value: 'Sikkim', label: 'Sikkim' },
  { value: 'Tamil Nadu', label: 'Tamil Nadu' },
  { value: 'Telangana', label: 'Telangana' },
  { value: 'Tripura', label: 'Tripura' },
  { value: 'Uttar Pradesh', label: 'Uttar Pradesh' },
  { value: 'Uttarakhand', label: 'Uttarakhand' },
  { value: 'West Bengal', label: 'West Bengal' },
  // Union Territories
  { value: 'Andaman and Nicobar Islands', label: 'Andaman and Nicobar Islands' },
  { value: 'Chandigarh', label: 'Chandigarh' },
  { value: 'Dadra and Nagar Haveli and Daman and Diu', label: 'Dadra and Nagar Haveli and Daman and Diu' },
  { value: 'Delhi', label: 'Delhi' },
  { value: 'Jammu and Kashmir', label: 'Jammu and Kashmir' },
  { value: 'Ladakh', label: 'Ladakh' },
  { value: 'Lakshadweep', label: 'Lakshadweep' },
  { value: 'Puducherry', label: 'Puducherry' },
];

@Component({
  selector: 'msls-address-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsInputComponent,
    MslsSelectComponent,
    MslsFormFieldComponent,
  ],
  template: `
    <div class="address-form" [formGroup]="addressFormGroup">
      <!-- Address Line 1 -->
      <div class="mb-4">
        <msls-form-field
          label="Address Line 1"
          [required]="required()"
          [error]="getError('addressLine1')"
        >
          <msls-input
            formControlName="addressLine1"
            placeholder="House/Flat No., Building, Street"
          />
        </msls-form-field>
      </div>

      <!-- Address Line 2 -->
      <div class="mb-4">
        <msls-form-field label="Address Line 2" [error]="getError('addressLine2')">
          <msls-input
            formControlName="addressLine2"
            placeholder="Area, Landmark (Optional)"
          />
        </msls-form-field>
      </div>

      <!-- City, State, Postal Code -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <msls-form-field
          label="City"
          [required]="required()"
          [error]="getError('city')"
        >
          <msls-input
            formControlName="city"
            placeholder="Enter city"
          />
        </msls-form-field>

        <msls-form-field
          label="State"
          [required]="required()"
          [error]="getError('state')"
        >
          <msls-select
            formControlName="state"
            [options]="stateOptions"
            placeholder="Select state"
          />
        </msls-form-field>

        <msls-form-field
          label="Postal Code"
          [required]="required()"
          [error]="getError('postalCode')"
          hint="6-digit PIN code"
        >
          <msls-input
            formControlName="postalCode"
            placeholder="Enter PIN code"
            [maxlength]="6"
          />
        </msls-form-field>
      </div>

      <!-- Country -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <msls-form-field
          label="Country"
          [required]="required()"
          [error]="getError('country')"
        >
          <msls-input
            formControlName="country"
            placeholder="Enter country"
            [readonly]="true"
          />
        </msls-form-field>
      </div>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AddressFormComponent {
  // =========================================================================
  // Inputs
  // =========================================================================

  /** Form group name within parent form */
  readonly formGroupName = input.required<string>();

  /** Parent form group */
  readonly parentForm = input.required<FormGroup>();

  /** Whether all fields are required */
  readonly required = input<boolean>(true);

  // =========================================================================
  // Options
  // =========================================================================

  readonly stateOptions = INDIAN_STATES;

  // =========================================================================
  // Computed
  // =========================================================================

  get addressFormGroup(): FormGroup {
    return this.parentForm().get(this.formGroupName()) as FormGroup;
  }

  // =========================================================================
  // Helpers
  // =========================================================================

  getError(controlName: string): string {
    const control = this.addressFormGroup?.get(controlName);
    if (!control || !control.touched || !control.errors) return '';

    if (control.errors['required']) return 'This field is required';
    if (control.errors['maxlength']) {
      return `Maximum ${control.errors['maxlength'].requiredLength} characters allowed`;
    }
    if (control.errors['pattern']) {
      if (controlName === 'postalCode') return 'Must be 6 digits';
      return 'Invalid format';
    }

    return 'Invalid value';
  }
}
