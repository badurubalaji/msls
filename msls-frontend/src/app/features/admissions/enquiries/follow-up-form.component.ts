/**
 * MSLS Follow-up Form Component
 *
 * Form component for adding follow-ups to admission enquiries.
 */

import {
  Component,
  input,
  output,
  inject,
  signal,
  OnInit,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';

import { MslsButtonComponent } from '../../../shared/components/button/button.component';
import { MslsSelectComponent, SelectOption } from '../../../shared/components/select/select.component';

import { EnquiryService } from './enquiry.service';
import {
  EnquiryFollowUp,
  CreateFollowUpDto,
  ContactMode,
  FollowUpOutcome,
  CONTACT_MODE_LABELS,
  FOLLOW_UP_OUTCOME_LABELS,
} from './enquiry.model';

/**
 * FollowUpFormComponent - Form for adding follow-ups.
 */
@Component({
  selector: 'msls-follow-up-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsButtonComponent,
    MslsSelectComponent,
  ],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-5">
      <!-- Follow-up Date -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: #374151;">
          Follow-up Date <span style="color: #ef4444;">*</span>
        </label>
        <input
          type="date"
          formControlName="followUpDate"
          class="w-full h-10 px-3 rounded-lg border text-sm focus:outline-none focus:ring-2 transition-colors"
          [style.border-color]="hasError('followUpDate') ? '#ef4444' : '#e2e8f0'"
          style="background-color: #ffffff; color: #1e293b;"
        />
        @if (hasError('followUpDate')) {
          <p class="mt-1 text-xs" style="color: #ef4444;">
            Follow-up date is required
          </p>
        }
      </div>

      <!-- Contact Mode -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: #374151;">
          Contact Mode <span style="color: #ef4444;">*</span>
        </label>
        <msls-select
          formControlName="contactMode"
          [options]="contactModeOptions"
          placeholder="Select contact mode"
        />
        @if (hasError('contactMode')) {
          <p class="mt-1 text-xs" style="color: #ef4444;">
            Please select a contact mode
          </p>
        }
      </div>

      <!-- Outcome -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: #374151;">
          Outcome
        </label>
        <msls-select
          formControlName="outcome"
          [options]="outcomeOptions"
          placeholder="Select outcome"
        />
      </div>

      <!-- Notes -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: #374151;">
          Notes
        </label>
        <textarea
          formControlName="notes"
          rows="3"
          placeholder="Add notes about the follow-up..."
          class="w-full px-3 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2 transition-colors resize-none"
          style="border-color: #e2e8f0; background-color: #ffffff; color: #1e293b;"
        ></textarea>
      </div>

      <!-- Next Follow-up Date -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: #374151;">
          Next Follow-up Date
        </label>
        <input
          type="date"
          formControlName="nextFollowUp"
          class="w-full h-10 px-3 rounded-lg border text-sm focus:outline-none focus:ring-2 transition-colors"
          style="border-color: #e2e8f0; background-color: #ffffff; color: #1e293b;"
        />
        <p class="mt-1 text-xs" style="color: #64748b;">
          Schedule the next follow-up date if needed
        </p>
      </div>

      <!-- Form Actions -->
      <div class="flex justify-end gap-3 pt-4" style="border-top: 1px solid #e2e8f0;">
        <msls-button variant="ghost" type="button" (click)="onCancel()">
          Cancel
        </msls-button>
        <msls-button
          variant="primary"
          type="submit"
          [loading]="submitting()"
          [disabled]="form.invalid || submitting()"
        >
          <i class="fa-solid fa-plus mr-2"></i>
          Add Follow-up
        </msls-button>
      </div>
    </form>
  `,
  styles: [`
    :host {
      display: block;
    }

    input:focus,
    textarea:focus {
      border-color: #3b82f6;
      box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    input[type="date"]::-webkit-calendar-picker-indicator {
      cursor: pointer;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FollowUpFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly enquiryService = inject(EnquiryService);

  /** Enquiry ID to add follow-up to */
  enquiryId = input.required<string>();

  /** Emitted when form is submitted successfully */
  submitted = output<EnquiryFollowUp>();

  /** Emitted when form is cancelled */
  cancelled = output<void>();

  /** Form submitting state */
  readonly submitting = signal(false);

  /** Reactive form */
  form!: FormGroup;

  /** Contact mode options */
  readonly contactModeOptions: SelectOption[] = Object.entries(CONTACT_MODE_LABELS).map(
    ([value, label]) => ({ value, label })
  );

  /** Outcome options */
  readonly outcomeOptions: SelectOption[] = [
    { value: '', label: 'Select outcome' },
    ...Object.entries(FOLLOW_UP_OUTCOME_LABELS).map(([value, label]) => ({
      value,
      label,
    })),
  ];

  ngOnInit(): void {
    this.initForm();
  }

  /**
   * Initialize form with defaults
   */
  private initForm(): void {
    // Default to today's date
    const today = new Date().toISOString().split('T')[0];

    this.form = this.fb.group({
      followUpDate: [today, [Validators.required]],
      contactMode: ['phone', [Validators.required]],
      outcome: [''],
      notes: [''],
      nextFollowUp: [''],
    });
  }

  /**
   * Check if field has validation error
   */
  hasError(fieldName: string): boolean {
    const field = this.form.get(fieldName);
    return !!field && field.invalid && (field.dirty || field.touched);
  }

  /**
   * Handle form submission
   */
  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    const formValue = this.form.value;

    const dto: CreateFollowUpDto = {
      followUpDate: formValue.followUpDate,
      contactMode: formValue.contactMode as ContactMode,
      outcome: formValue.outcome || undefined,
      notes: formValue.notes || undefined,
      nextFollowUp: formValue.nextFollowUp || undefined,
    };

    this.enquiryService.addFollowUp(this.enquiryId(), dto).subscribe({
      next: (followUp) => {
        this.submitting.set(false);
        this.submitted.emit(followUp);
      },
      error: () => {
        this.submitting.set(false);
      },
    });
  }

  /**
   * Handle cancel action
   */
  onCancel(): void {
    this.cancelled.emit();
  }
}
