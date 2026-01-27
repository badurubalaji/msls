/**
 * MSLS Decision Form Component
 *
 * Form for making admission decisions (approve/waitlist/reject)
 * with section assignment and notes.
 */

import { Component, Input, Output, EventEmitter, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import {
  MslsButtonComponent,
  MslsInputComponent,
  MslsFormFieldComponent,
  MslsSelectComponent,
  MslsBadgeComponent,
  SelectOption,
} from '../../../shared/components';
import { ToastService } from '../../../shared/services/toast.service';
import {
  MeritListEntry,
  DecisionType,
  MakeDecisionRequest,
  BulkDecisionRequest,
  getDecisionConfig,
  DECISION_TYPE_CONFIG,
} from './merit.model';
import { MeritService } from './merit.service';

@Component({
  selector: 'msls-decision-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsButtonComponent,
    MslsInputComponent,
    MslsFormFieldComponent,
    MslsSelectComponent,
  ],
  template: `
    <div class="decision-form">
      <!-- Decision Type Display -->
      <div class="decision-header">
        <div class="decision-type" [class]="'decision-type--' + decisionType">
          <i [class]="getDecisionConfig(decisionType).icon"></i>
          <span>{{ getDecisionConfig(decisionType).label }}</span>
        </div>

        @if (isBulk) {
          <p class="decision-summary">
            You are about to {{ getDecisionActionText() }} <strong>{{ entries.length }}</strong> candidates.
          </p>
        } @else if (entry) {
          <p class="decision-summary">
            {{ getDecisionActionText() | titlecase }} admission for <strong>{{ entry.studentName }}</strong>
            (Rank #{{ entry.rank }}, Score: {{ entry.percentage | number:'1.1-1' }}%)
          </p>
        }
      </div>

      <form [formGroup]="form" (ngSubmit)="onSubmit()">
        <!-- Section Assignment (for selected/approved) -->
        @if (decisionType === 'selected') {
          <msls-form-field
            label="Assign Section"
            [required]="true"
            [error]="getFieldError('sectionAssigned') || ''"
          >
            <msls-select
              formControlName="sectionAssigned"
              [options]="sectionOptions()"
              placeholder="Select section"
            />
          </msls-form-field>
        }

        <!-- Waitlist Position (for waitlisted) -->
        @if (decisionType === 'waitlisted' && !isBulk) {
          <msls-form-field
            label="Waitlist Position"
            hint="Optional - leave blank for auto-assignment"
          >
            <msls-input
              type="number"
              formControlName="waitlistPosition"
              placeholder="e.g., 1"
            />
          </msls-form-field>
        }

        <!-- Rejection Reason (for rejected) -->
        @if (decisionType === 'rejected') {
          <msls-form-field
            label="Rejection Reason"
            [error]="getFieldError('rejectionReason') || ''"
          >
            <msls-select
              formControlName="rejectionReason"
              [options]="rejectionReasonOptions"
              placeholder="Select reason"
            />
          </msls-form-field>
        }

        <!-- Offer Valid Until (for selected) -->
        @if (decisionType === 'selected') {
          <msls-form-field
            label="Offer Valid Until"
            hint="Date by which the offer must be accepted"
            [error]="getFieldError('offerValidUntil') || ''"
          >
            <msls-input
              type="date"
              formControlName="offerValidUntil"
            />
          </msls-form-field>
        }

        <!-- Remarks/Notes -->
        <msls-form-field
          label="Notes/Remarks"
          hint="Optional additional notes"
        >
          <textarea
            formControlName="remarks"
            class="remarks-textarea"
            rows="3"
            placeholder="Add any additional notes..."
          ></textarea>
        </msls-form-field>

        <!-- Bulk Selection Preview -->
        @if (isBulk && entries.length > 0) {
          <div class="bulk-preview">
            <label class="preview-label">Selected Candidates ({{ entries.length }})</label>
            <div class="preview-list">
              @for (e of entries.slice(0, 5); track e.id) {
                <div class="preview-item">
                  <span class="preview-rank">#{{ e.rank }}</span>
                  <span class="preview-name">{{ e.studentName }}</span>
                  <span class="preview-score">{{ e.percentage | number:'1.1-1' }}%</span>
                </div>
              }
              @if (entries.length > 5) {
                <div class="preview-more">
                  +{{ entries.length - 5 }} more candidates
                </div>
              }
            </div>
          </div>
        }

        <!-- Confirmation Checkbox -->
        <label class="confirmation-checkbox">
          <input
            type="checkbox"
            formControlName="confirmed"
          />
          <span class="checkbox-custom"></span>
          <span class="checkbox-text">
            I confirm this decision. {{ getConfirmationText() }}
          </span>
        </label>

        <!-- Form Actions -->
        <div class="form-actions">
          <msls-button
            type="button"
            variant="secondary"
            (click)="onCancel()"
          >
            Cancel
          </msls-button>
          <msls-button
            type="submit"
            [variant]="getButtonVariant()"
            [loading]="loading()"
            [disabled]="!form.valid || !form.get('confirmed')?.value"
          >
            <i [class]="getDecisionConfig(decisionType).icon"></i>
            {{ getSubmitButtonText() }}
          </msls-button>
        </div>
      </form>
    </div>
  `,
  styles: [`
    .decision-form {
      display: flex;
      flex-direction: column;
      gap: 1.25rem;
    }

    .decision-header {
      text-align: center;
      padding-bottom: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .decision-type {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-weight: 600;
      font-size: 1rem;
      margin-bottom: 0.75rem;
    }

    .decision-type--selected {
      background: #dcfce7;
      color: #16a34a;
    }

    .decision-type--waitlisted {
      background: #fef3c7;
      color: #d97706;
    }

    .decision-type--rejected {
      background: #fee2e2;
      color: #dc2626;
    }

    .decision-summary {
      margin: 0;
      font-size: 0.875rem;
      color: #64748b;
    }

    .decision-summary strong {
      color: #0f172a;
    }

    .remarks-textarea {
      width: 100%;
      padding: 0.625rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-family: inherit;
      resize: vertical;
      min-height: 80px;
      transition: border-color 0.15s, box-shadow 0.15s;
    }

    .remarks-textarea:focus {
      outline: none;
      border-color: #3b82f6;
      box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    .bulk-preview {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      padding: 0.75rem 1rem;
    }

    .preview-label {
      display: block;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      margin-bottom: 0.5rem;
    }

    .preview-list {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .preview-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 0.8125rem;
    }

    .preview-rank {
      color: #64748b;
      font-weight: 500;
      min-width: 32px;
    }

    .preview-name {
      flex: 1;
      color: #0f172a;
    }

    .preview-score {
      color: #64748b;
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 0.75rem;
    }

    .preview-more {
      font-size: 0.75rem;
      color: #64748b;
      font-style: italic;
      padding-top: 0.25rem;
    }

    .confirmation-checkbox {
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      padding: 1rem;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: background 0.15s;
    }

    .confirmation-checkbox:hover {
      background: #f1f5f9;
    }

    .confirmation-checkbox input[type="checkbox"] {
      display: none;
    }

    .checkbox-custom {
      width: 1.25rem;
      height: 1.25rem;
      border: 2px solid #cbd5e1;
      border-radius: 0.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      transition: all 0.15s;
      margin-top: 0.125rem;
    }

    .confirmation-checkbox input:checked + .checkbox-custom {
      background: #3b82f6;
      border-color: #3b82f6;
    }

    .confirmation-checkbox input:checked + .checkbox-custom::after {
      content: '';
      width: 0.375rem;
      height: 0.625rem;
      border: solid white;
      border-width: 0 2px 2px 0;
      transform: rotate(45deg);
      margin-bottom: 2px;
    }

    .checkbox-text {
      font-size: 0.875rem;
      color: #334155;
      line-height: 1.5;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }
  `],
})
export class DecisionFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private meritService = inject(MeritService);
  private toastService = inject(ToastService);

  @Input() entry: MeritListEntry | null = null;
  @Input() entries: MeritListEntry[] = [];
  @Input() decisionType: DecisionType = 'selected';
  @Input() isBulk = false;

  @Output() save = new EventEmitter<{ success: boolean }>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;
  loading = signal(false);
  sectionOptions = signal<SelectOption[]>([]);

  rejectionReasonOptions: SelectOption[] = [
    { value: 'low_score', label: 'Score below required threshold' },
    { value: 'no_seats', label: 'No seats available' },
    { value: 'documents_incomplete', label: 'Incomplete documentation' },
    { value: 'age_criteria', label: 'Does not meet age criteria' },
    { value: 'other', label: 'Other reason' },
  ];

  ngOnInit(): void {
    this.initForm();
    this.loadSections();
  }

  private initForm(): void {
    this.form = this.fb.group({
      sectionAssigned: [
        '',
        this.decisionType === 'selected' ? [Validators.required] : [],
      ],
      waitlistPosition: [null],
      rejectionReason: [''],
      offerValidUntil: [this.getDefaultOfferValidDate()],
      remarks: [''],
      confirmed: [false, [Validators.requiredTrue]],
    });
  }

  private getDefaultOfferValidDate(): string {
    const date = new Date();
    date.setDate(date.getDate() + 14); // Default 2 weeks validity
    return date.toISOString().split('T')[0];
  }

  private loadSections(): void {
    const className = this.entry?.classApplying || this.entries[0]?.classApplying || 'Class 1';
    this.meritService.getSections(className).subscribe({
      next: (sections) => {
        this.sectionOptions.set(
          sections.map(s => ({
            value: s.name,
            label: s.name,
          }))
        );
      },
      error: (err) => {
        console.error('Failed to load sections:', err);
      },
    });
  }

  getDecisionConfig(decision: DecisionType) {
    return DECISION_TYPE_CONFIG[decision] || { label: decision, icon: 'fa-solid fa-circle', variant: 'neutral' };
  }

  getDecisionActionText(): string {
    switch (this.decisionType) {
      case 'selected': return 'approve';
      case 'waitlisted': return 'waitlist';
      case 'rejected': return 'reject';
      default: return this.decisionType;
    }
  }

  getConfirmationText(): string {
    switch (this.decisionType) {
      case 'selected':
        return 'An offer letter will be sent to the applicant.';
      case 'waitlisted':
        return 'The applicant will be notified of their waitlist status.';
      case 'rejected':
        return 'The applicant will be notified of the rejection.';
      default:
        return '';
    }
  }

  getButtonVariant(): 'primary' | 'secondary' | 'danger' {
    switch (this.decisionType) {
      case 'selected': return 'primary';
      case 'rejected': return 'danger';
      default: return 'secondary';
    }
  }

  getSubmitButtonText(): string {
    const count = this.isBulk ? ` (${this.entries.length})` : '';
    switch (this.decisionType) {
      case 'selected': return `Approve${count}`;
      case 'waitlisted': return `Waitlist${count}`;
      case 'rejected': return `Reject${count}`;
      default: return `Confirm${count}`;
    }
  }

  getFieldError(field: string): string | null {
    const control = this.form.get(field);
    if (!control || !control.touched || !control.errors) return null;
    if (control.errors['required']) return 'This field is required';
    return null;
  }

  onSubmit(): void {
    if (!this.form.valid) return;

    this.loading.set(true);
    const value = this.form.value;

    if (this.isBulk) {
      const request: BulkDecisionRequest = {
        applicationIds: this.entries.map(e => e.applicationId),
        decision: this.decisionType,
        sectionAssigned: value.sectionAssigned || undefined,
        remarks: value.remarks || undefined,
        offerValidUntil: value.offerValidUntil || undefined,
      };

      this.meritService.makeBulkDecision(request).subscribe({
        next: () => {
          this.loading.set(false);
          this.toastService.success(`Successfully ${this.getDecisionActionText()}ed ${this.entries.length} candidates`);
          this.save.emit({ success: true });
        },
        error: (err) => {
          this.loading.set(false);
          console.error('Failed to make bulk decision:', err);
          this.toastService.error('Failed to process decisions');
          this.save.emit({ success: false });
        },
      });
    } else if (this.entry) {
      const request: MakeDecisionRequest = {
        decision: this.decisionType,
        sectionAssigned: value.sectionAssigned || undefined,
        waitlistPosition: value.waitlistPosition || undefined,
        rejectionReason: value.rejectionReason || undefined,
        remarks: value.remarks || undefined,
        offerValidUntil: value.offerValidUntil || undefined,
      };

      this.meritService.makeDecision(this.entry.applicationId, request).subscribe({
        next: () => {
          this.loading.set(false);
          this.toastService.success(`Successfully ${this.getDecisionActionText()}ed ${this.entry?.studentName}`);
          this.save.emit({ success: true });
        },
        error: (err) => {
          this.loading.set(false);
          console.error('Failed to make decision:', err);
          this.toastService.error('Failed to process decision');
          this.save.emit({ success: false });
        },
      });
    }
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
