import { Component, EventEmitter, Input, OnInit, Output, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  SalaryStructure,
  CreateStructureRequest,
  ComponentDropdownItem,
  StructureComponentInput,
} from './salary.model';
import { SalaryService } from './salary.service';
import { DesignationService } from '../designations/designation.service';
import { DesignationDropdownItem } from '../designations/designation.model';

@Component({
  selector: 'msls-salary-structure-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="form">
      <div class="form-row">
        <div class="form-group">
          <label for="name" class="label required">Structure Name</label>
          <input
            id="name"
            type="text"
            formControlName="name"
            class="input"
            placeholder="e.g., Teacher Grade A"
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
            placeholder="e.g., TGA"
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
          placeholder="Optional description for this structure..."
        ></textarea>
      </div>

      <div class="form-group">
        <label for="designationId" class="label">Default Designation</label>
        <select id="designationId" formControlName="designationId" class="select">
          <option value="">None (Generic Structure)</option>
          @for (designation of designations(); track designation.id) {
            <option [value]="designation.id">{{ designation.name }} (Level {{ designation.level }})</option>
          }
        </select>
        <span class="hint">Optionally link this structure to a specific designation</span>
      </div>

      <!-- Components Section -->
      <div class="components-section">
        <div class="section-header">
          <h3>Salary Components</h3>
          <button type="button" class="btn btn-sm btn-secondary" (click)="addComponent()">
            <i class="fa-solid fa-plus"></i>
            Add Component
          </button>
        </div>

        @if (componentsArray.length === 0) {
          <div class="empty-components">
            <i class="fa-regular fa-money-bill-1"></i>
            <p>No components added yet</p>
            <span>Add at least one component to define this salary structure</span>
          </div>
        } @else {
          <div class="components-list" formArrayName="components">
            @for (compCtrl of componentsArray.controls; track $index; let i = $index) {
              <div class="component-row" [formGroupName]="i">
                <div class="component-select">
                  <select formControlName="componentId" class="select">
                    <option value="">Select component...</option>
                    @for (comp of availableComponents(); track comp.id) {
                      <option [value]="comp.id">
                        {{ comp.name }} ({{ comp.code }}) - {{ comp.componentType === 'earning' ? 'Earning' : 'Deduction' }}
                      </option>
                    }
                  </select>
                </div>
                <div class="component-amount">
                  <input
                    type="number"
                    formControlName="amount"
                    class="input"
                    placeholder="Amount"
                    min="0"
                    step="0.01"
                  />
                </div>
                <div class="component-percentage">
                  <input
                    type="number"
                    formControlName="percentage"
                    class="input"
                    placeholder="%"
                    min="0"
                    max="100"
                    step="0.01"
                  />
                </div>
                <button
                  type="button"
                  class="remove-btn"
                  title="Remove"
                  (click)="removeComponent(i)"
                >
                  <i class="fa-solid fa-xmark"></i>
                </button>
              </div>
            }
          </div>
        }

        @if (form.get('components')?.invalid && form.get('components')?.touched) {
          <span class="error">At least one component is required</span>
        }
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" (click)="cancel.emit()">
          Cancel
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          [disabled]="form.invalid || loading || componentsArray.length === 0"
        >
          @if (loading) {
            <div class="btn-spinner"></div>
            Saving...
          } @else {
            <i class="fa-solid fa-check"></i>
            {{ structure ? 'Update' : 'Create' }}
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

    .hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .error {
      font-size: 0.75rem;
      color: #dc2626;
    }

    .components-section {
      margin-top: 0.5rem;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
    }

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }

    .section-header h3 {
      margin: 0;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .empty-components {
      text-align: center;
      padding: 2rem;
      color: #64748b;
    }

    .empty-components i {
      font-size: 2rem;
      color: #cbd5e1;
      margin-bottom: 0.5rem;
    }

    .empty-components p {
      margin: 0;
      font-weight: 500;
    }

    .empty-components span {
      font-size: 0.75rem;
    }

    .components-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .component-row {
      display: grid;
      grid-template-columns: 1fr 120px 80px 40px;
      gap: 0.5rem;
      align-items: center;
      padding: 0.5rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
    }

    .component-select .select,
    .component-amount .input,
    .component-percentage .input {
      width: 100%;
      padding: 0.5rem;
      font-size: 0.875rem;
    }

    .remove-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border: none;
      background: transparent;
      color: #dc2626;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: background 0.2s;
    }

    .remove-btn:hover {
      background: #fef2f2;
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

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
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
      to { transform: rotate(360deg); }
    }

    @media (max-width: 640px) {
      .form-row {
        grid-template-columns: 1fr;
      }

      .component-row {
        grid-template-columns: 1fr 1fr 40px;
        grid-template-rows: auto auto;
      }

      .component-select {
        grid-column: 1 / -2;
      }
    }
  `],
})
export class SalaryStructureFormComponent implements OnInit {
  @Input() structure: SalaryStructure | null = null;
  @Input() loading = false;
  @Output() save = new EventEmitter<CreateStructureRequest>();
  @Output() cancel = new EventEmitter<void>();

  private fb = inject(FormBuilder);
  private salaryService = inject(SalaryService);
  private designationService = inject(DesignationService);

  availableComponents = signal<ComponentDropdownItem[]>([]);
  designations = signal<DesignationDropdownItem[]>([]);

  form: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.maxLength(100)]],
    code: ['', [Validators.required, Validators.maxLength(20)]],
    description: [''],
    designationId: [''],
    components: this.fb.array([], Validators.minLength(1)),
  });

  get componentsArray(): FormArray {
    return this.form.get('components') as FormArray;
  }

  ngOnInit(): void {
    this.loadComponents();
    this.loadDesignations();

    if (this.structure) {
      this.form.patchValue({
        name: this.structure.name,
        code: this.structure.code,
        description: this.structure.description || '',
        designationId: this.structure.designationId || '',
      });

      // Populate existing components
      if (this.structure.components) {
        this.structure.components.forEach(comp => {
          this.componentsArray.push(
            this.fb.group({
              componentId: [comp.componentId, Validators.required],
              amount: [comp.amount || ''],
              percentage: [comp.percentage || ''],
            })
          );
        });
      }
    }
  }

  loadComponents(): void {
    this.salaryService.getComponentsDropdown().subscribe({
      next: components => this.availableComponents.set(components),
      error: () => this.availableComponents.set([]),
    });
  }

  loadDesignations(): void {
    this.designationService.getDesignationsDropdown().subscribe({
      next: designations => this.designations.set(designations),
      error: () => this.designations.set([]),
    });
  }

  addComponent(): void {
    this.componentsArray.push(
      this.fb.group({
        componentId: ['', Validators.required],
        amount: [''],
        percentage: [''],
      })
    );
  }

  removeComponent(index: number): void {
    this.componentsArray.removeAt(index);
  }

  onSubmit(): void {
    if (this.form.invalid || this.componentsArray.length === 0) return;

    const formValue = this.form.value;
    const components: StructureComponentInput[] = formValue.components.map(
      (comp: { componentId: string; amount: string; percentage: string }) => ({
        componentId: comp.componentId,
        amount: comp.amount || undefined,
        percentage: comp.percentage || undefined,
      })
    );

    const request: CreateStructureRequest = {
      name: formValue.name.trim(),
      code: formValue.code.trim().toUpperCase(),
      description: formValue.description?.trim() || undefined,
      designationId: formValue.designationId || undefined,
      components,
    };

    this.save.emit(request);
  }
}
