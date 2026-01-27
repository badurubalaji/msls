import { Component, EventEmitter, Input, OnInit, Output, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  SalaryComponent,
  CreateComponentRequest,
  ComponentDropdownItem,
  ComponentType,
  CalculationType,
} from './salary.model';
import { SalaryService } from './salary.service';

@Component({
  selector: 'msls-salary-component-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="form">
      <div class="form-row">
        <div class="form-group">
          <label for="name" class="label required">Component Name</label>
          <input
            id="name"
            type="text"
            formControlName="name"
            class="input"
            placeholder="e.g., Basic Salary"
          />
          @if (form.get('name')?.invalid && form.get('name')?.touched) {
            <span class="error">Name is required</span>
          }
        </div>

        <div class="form-group">
          <label for="code" class="label required">Code</label>
          <input
            id="code"
            type="text"
            formControlName="code"
            class="input"
            placeholder="e.g., BASIC"
            style="text-transform: uppercase;"
          />
          @if (form.get('code')?.invalid && form.get('code')?.touched) {
            <span class="error">Code is required (max 20 characters)</span>
          }
        </div>
      </div>

      <div class="form-group">
        <label for="description" class="label">Description</label>
        <textarea
          id="description"
          formControlName="description"
          class="textarea"
          rows="2"
          placeholder="Optional description for this component..."
        ></textarea>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="componentType" class="label required">Component Type</label>
          <select id="componentType" formControlName="componentType" class="select">
            <option value="earning">Earning</option>
            <option value="deduction">Deduction</option>
          </select>
        </div>

        <div class="form-group">
          <label for="calculationType" class="label required">Calculation Type</label>
          <select
            id="calculationType"
            formControlName="calculationType"
            class="select"
            (change)="onCalculationTypeChange()"
          >
            <option value="fixed">Fixed Amount</option>
            <option value="percentage">Percentage</option>
          </select>
        </div>
      </div>

      @if (form.get('calculationType')?.value === 'percentage') {
        <div class="form-group">
          <label for="percentageOfId" class="label required">Percentage Of</label>
          <select id="percentageOfId" formControlName="percentageOfId" class="select">
            <option value="">Select component...</option>
            @for (comp of earningComponents(); track comp.id) {
              <option [value]="comp.id">{{ comp.name }} ({{ comp.code }})</option>
            }
          </select>
          @if (form.get('percentageOfId')?.invalid && form.get('percentageOfId')?.touched) {
            <span class="error">Please select the base component for percentage calculation</span>
          }
        </div>
      }

      <div class="form-row">
        <div class="form-group">
          <label for="displayOrder" class="label">Display Order</label>
          <input
            id="displayOrder"
            type="number"
            formControlName="displayOrder"
            class="input"
            min="0"
          />
        </div>

        <div class="form-group checkbox-group">
          <label class="checkbox-label">
            <input type="checkbox" formControlName="isTaxable" />
            <span class="checkmark"></span>
            Is Taxable
          </label>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" (click)="cancel.emit()">
          Cancel
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          [disabled]="form.invalid || loading"
        >
          @if (loading) {
            <div class="btn-spinner"></div>
            Saving...
          } @else {
            <i class="fa-solid fa-check"></i>
            {{ component ? 'Update' : 'Create' }}
          }
        </button>
      </div>
    </form>
  `,
  styles: [`
    .form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .label.required::after {
      content: ' *';
      color: #dc2626;
    }

    .input,
    .select,
    .textarea {
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: border-color 0.2s, box-shadow 0.2s;
    }

    .input:focus,
    .select:focus,
    .textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .textarea {
      resize: vertical;
      min-height: 60px;
    }

    .error {
      font-size: 0.75rem;
      color: #dc2626;
    }

    .checkbox-group {
      display: flex;
      align-items: center;
      padding-top: 1.5rem;
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      cursor: pointer;
      font-size: 0.875rem;
      color: #374151;
    }

    .checkbox-label input[type="checkbox"] {
      width: 1rem;
      height: 1rem;
      accent-color: #4f46e5;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      margin-top: 0.5rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to {
        transform: rotate(360deg);
      }
    }

    @media (max-width: 640px) {
      .form-row {
        grid-template-columns: 1fr;
      }

      .checkbox-group {
        padding-top: 0;
      }
    }
  `],
})
export class SalaryComponentFormComponent implements OnInit {
  @Input() component: SalaryComponent | null = null;
  @Input() loading = false;
  @Output() save = new EventEmitter<CreateComponentRequest>();
  @Output() cancel = new EventEmitter<void>();

  private fb = inject(FormBuilder);
  private salaryService = inject(SalaryService);

  earningComponents = signal<ComponentDropdownItem[]>([]);

  form: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.maxLength(100)]],
    code: ['', [Validators.required, Validators.maxLength(20)]],
    description: [''],
    componentType: ['earning' as ComponentType, Validators.required],
    calculationType: ['fixed' as CalculationType, Validators.required],
    percentageOfId: [''],
    isTaxable: [true],
    displayOrder: [0, [Validators.min(0)]],
  });

  ngOnInit(): void {
    this.loadEarningComponents();

    if (this.component) {
      this.form.patchValue({
        name: this.component.name,
        code: this.component.code,
        description: this.component.description || '',
        componentType: this.component.componentType,
        calculationType: this.component.calculationType,
        percentageOfId: this.component.percentageOfId || '',
        isTaxable: this.component.isTaxable,
        displayOrder: this.component.displayOrder,
      });
    }

    this.updatePercentageValidation();
  }

  loadEarningComponents(): void {
    this.salaryService.getComponentsDropdown('earning').subscribe({
      next: components => {
        // Filter out the current component if editing
        const filtered = this.component
          ? components.filter(c => c.id !== this.component!.id)
          : components;
        this.earningComponents.set(filtered);
      },
      error: () => {
        this.earningComponents.set([]);
      },
    });
  }

  onCalculationTypeChange(): void {
    this.updatePercentageValidation();
  }

  private updatePercentageValidation(): void {
    const calcType = this.form.get('calculationType')?.value;
    const percentageOfId = this.form.get('percentageOfId');

    if (calcType === 'percentage') {
      percentageOfId?.setValidators(Validators.required);
    } else {
      percentageOfId?.clearValidators();
      percentageOfId?.setValue('');
    }
    percentageOfId?.updateValueAndValidity();
  }

  onSubmit(): void {
    if (this.form.invalid) return;

    const formValue = this.form.value;
    const request: CreateComponentRequest = {
      name: formValue.name.trim(),
      code: formValue.code.trim().toUpperCase(),
      description: formValue.description?.trim() || undefined,
      componentType: formValue.componentType,
      calculationType: formValue.calculationType,
      percentageOfId: formValue.percentageOfId || undefined,
      isTaxable: formValue.isTaxable,
      displayOrder: formValue.displayOrder,
    };

    this.save.emit(request);
  }
}
