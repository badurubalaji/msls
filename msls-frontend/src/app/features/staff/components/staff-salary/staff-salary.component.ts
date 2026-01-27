import { Component, Input, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MslsModalComponent } from '../../../../shared/components/modal/modal.component';
import { SalaryService } from '../../../admin/salary/salary.service';
import {
  StaffSalary,
  ComponentDropdownItem,
  StructureDropdownItem,
  AssignSalaryRequest,
  StaffSalaryComponentInput,
} from '../../../admin/salary/salary.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-staff-salary',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MslsModalComponent],
  template: `
    <div class="salary-container">
      @if (loading()) {
        <div class="loading-state">
          <div class="spinner"></div>
          <span>Loading salary information...</span>
        </div>
      } @else if (error()) {
        <div class="error-state">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button class="btn btn-secondary btn-sm" (click)="loadSalary()">
            <i class="fa-solid fa-refresh"></i>
            Retry
          </button>
        </div>
      } @else if (!currentSalary()) {
        <!-- No Salary Assigned -->
        <div class="empty-state">
          <div class="empty-icon">
            <i class="fa-solid fa-indian-rupee-sign"></i>
          </div>
          <h3>No Salary Assigned</h3>
          <p>This staff member does not have a salary structure assigned yet.</p>
          <button class="btn btn-primary" (click)="openAssignModal()">
            <i class="fa-solid fa-plus"></i>
            Assign Salary
          </button>
        </div>
      } @else {
        <!-- Current Salary Display -->
        <div class="salary-header">
          <div class="salary-summary">
            <div class="salary-amount gross">
              <span class="label">Gross Salary</span>
              <span class="value">{{ formatCurrency(currentSalary()?.grossSalary || '0') }}</span>
            </div>
            <div class="salary-amount net">
              <span class="label">Net Salary</span>
              <span class="value">{{ formatCurrency(currentSalary()?.netSalary || '0') }}</span>
            </div>
            @if (currentSalary()?.ctc) {
              <div class="salary-amount ctc">
                <span class="label">CTC</span>
                <span class="value">{{ formatCurrency(currentSalary()?.ctc || '0') }}</span>
              </div>
            }
          </div>
          <div class="salary-actions">
            <button class="btn btn-secondary btn-sm" (click)="toggleHistory()">
              <i class="fa-solid fa-clock-rotate-left"></i>
              {{ showHistory() ? 'Hide' : 'Show' }} History
            </button>
            <button class="btn btn-primary btn-sm" (click)="openReviseModal()">
              <i class="fa-solid fa-pen"></i>
              Revise Salary
            </button>
          </div>
        </div>

        <div class="salary-meta">
          <span>
            <i class="fa-solid fa-calendar"></i>
            Effective from {{ formatDate(currentSalary()?.effectiveFrom || '') }}
          </span>
          @if (currentSalary()?.structureName) {
            <span>
              <i class="fa-solid fa-sitemap"></i>
              Structure: {{ currentSalary()?.structureName }}
            </span>
          }
        </div>

        <!-- Component Breakdown -->
        <div class="components-section">
          <h4>Salary Breakdown</h4>
          <div class="components-grid">
            <!-- Earnings -->
            <div class="component-card earnings">
              <h5>
                <i class="fa-solid fa-arrow-up"></i>
                Earnings
              </h5>
              <div class="component-list">
                @for (comp of earningComponents(); track comp.id) {
                  <div class="component-item">
                    <span class="comp-name">
                      {{ comp.componentName }}
                      @if (comp.isOverridden) {
                        <i class="fa-solid fa-asterisk overridden" title="Overridden from structure"></i>
                      }
                    </span>
                    <span class="comp-amount">{{ formatCurrency(comp.amount) }}</span>
                  </div>
                }
              </div>
              <div class="component-total">
                <span>Total Earnings</span>
                <span>{{ formatCurrency(totalEarnings()) }}</span>
              </div>
            </div>

            <!-- Deductions -->
            <div class="component-card deductions">
              <h5>
                <i class="fa-solid fa-arrow-down"></i>
                Deductions
              </h5>
              <div class="component-list">
                @for (comp of deductionComponents(); track comp.id) {
                  <div class="component-item">
                    <span class="comp-name">
                      {{ comp.componentName }}
                      @if (comp.isOverridden) {
                        <i class="fa-solid fa-asterisk overridden" title="Overridden from structure"></i>
                      }
                    </span>
                    <span class="comp-amount">{{ formatCurrency(comp.amount) }}</span>
                  </div>
                } @empty {
                  <div class="no-components">No deductions</div>
                }
              </div>
              <div class="component-total">
                <span>Total Deductions</span>
                <span>{{ formatCurrency(totalDeductions()) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Salary History -->
        @if (showHistory()) {
          <div class="history-section">
            <h4>Salary History</h4>
            @if (salaryHistory().length <= 1) {
              <div class="no-history">
                <i class="fa-regular fa-clock"></i>
                <span>No previous salary records</span>
              </div>
            } @else {
              <table class="history-table">
                <thead>
                  <tr>
                    <th>Effective Period</th>
                    <th>Structure</th>
                    <th style="text-align: right;">Gross</th>
                    <th style="text-align: right;">Net</th>
                    <th>Reason</th>
                  </tr>
                </thead>
                <tbody>
                  @for (salary of salaryHistory(); track salary.id) {
                    <tr [class.current]="salary.isCurrent">
                      <td>
                        {{ formatDate(salary.effectiveFrom) }}
                        @if (salary.effectiveTo) {
                          <span> - {{ formatDate(salary.effectiveTo) }}</span>
                        } @else {
                          <span class="current-badge">Current</span>
                        }
                      </td>
                      <td>{{ salary.structureName || '-' }}</td>
                      <td style="text-align: right; font-family: monospace;">
                        {{ formatCurrency(salary.grossSalary) }}
                      </td>
                      <td style="text-align: right; font-family: monospace;">
                        {{ formatCurrency(salary.netSalary) }}
                      </td>
                      <td>{{ salary.revisionReason || '-' }}</td>
                    </tr>
                  }
                </tbody>
              </table>
            }
          </div>
        }
      }

      <!-- Assign/Revise Salary Modal -->
      <msls-modal
        [isOpen]="showModal()"
        [title]="currentSalary() ? 'Revise Salary' : 'Assign Salary'"
        size="lg"
        (closed)="closeModal()"
      >
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="salary-form">
          <div class="form-row">
            <div class="form-group">
              <label class="label">Salary Structure (Optional)</label>
              <select
                formControlName="structureId"
                class="select"
                (change)="onStructureChange()"
              >
                <option value="">Custom (No Structure)</option>
                @for (structure of structures(); track structure.id) {
                  <option [value]="structure.id">{{ structure.name }} ({{ structure.code }})</option>
                }
              </select>
              <span class="hint">Select a predefined structure or define custom components</span>
            </div>

            <div class="form-group">
              <label class="label required">Effective From</label>
              <input
                type="date"
                formControlName="effectiveFrom"
                class="input"
              />
              @if (form.get('effectiveFrom')?.invalid && form.get('effectiveFrom')?.touched) {
                <span class="error">Effective date is required</span>
              }
            </div>
          </div>

          <div class="form-group">
            <label class="label">Revision Reason</label>
            <textarea
              formControlName="revisionReason"
              class="textarea"
              rows="2"
              placeholder="e.g., Annual increment, Promotion, etc."
            ></textarea>
          </div>

          <!-- Components -->
          <div class="components-form-section">
            <div class="section-header">
              <h4>Salary Components</h4>
              <button type="button" class="btn btn-sm btn-secondary" (click)="addComponent()">
                <i class="fa-solid fa-plus"></i>
                Add Component
              </button>
            </div>

            @if (componentsArray.length === 0) {
              <div class="no-components-form">
                <i class="fa-regular fa-money-bill-1"></i>
                <p>No components added</p>
                <span>Add components to define the salary</span>
              </div>
            } @else {
              <div class="components-form-list" formArrayName="components">
                @for (ctrl of componentsArray.controls; track $index; let i = $index) {
                  <div class="component-form-row" [formGroupName]="i">
                    <div class="comp-select">
                      <select formControlName="componentId" class="select">
                        <option value="">Select component...</option>
                        @for (comp of availableComponents(); track comp.id) {
                          <option [value]="comp.id">
                            {{ comp.name }} ({{ comp.componentType === 'earning' ? 'E' : 'D' }})
                          </option>
                        }
                      </select>
                    </div>
                    <div class="comp-amount">
                      <input
                        type="number"
                        formControlName="amount"
                        class="input"
                        placeholder="Amount"
                        min="0"
                        step="0.01"
                      />
                    </div>
                    <div class="comp-override">
                      <label class="checkbox-label">
                        <input type="checkbox" formControlName="isOverridden" />
                        Overridden
                      </label>
                    </div>
                    <button
                      type="button"
                      class="remove-btn"
                      (click)="removeComponent(i)"
                    >
                      <i class="fa-solid fa-xmark"></i>
                    </button>
                  </div>
                }
              </div>
            }
          </div>

          <!-- Summary -->
          <div class="form-summary">
            <div class="summary-item">
              <span>Gross Salary:</span>
              <span class="summary-value">{{ formatCurrency(calculatedGross()) }}</span>
            </div>
            <div class="summary-item">
              <span>Total Deductions:</span>
              <span class="summary-value">{{ formatCurrency(calculatedDeductions()) }}</span>
            </div>
            <div class="summary-item net">
              <span>Net Salary:</span>
              <span class="summary-value">{{ formatCurrency(calculatedNet()) }}</span>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeModal()">
              Cancel
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              [disabled]="form.invalid || saving() || componentsArray.length === 0"
            >
              @if (saving()) {
                <div class="btn-spinner"></div>
                Saving...
              } @else {
                <i class="fa-solid fa-check"></i>
                {{ currentSalary() ? 'Update Salary' : 'Assign Salary' }}
              }
            </button>
          </div>
        </form>
      </msls-modal>
    </div>
  `,
  styles: [`
    .salary-container {
      padding: 0.5rem;
    }

    .loading-state,
    .error-state {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-state {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .empty-state {
      text-align: center;
      padding: 3rem 2rem;
    }

    .empty-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      background: #f1f5f9;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
      color: #64748b;
    }

    .empty-state h3 {
      margin: 0 0 0.5rem;
      color: #1e293b;
    }

    .empty-state p {
      margin: 0 0 1.5rem;
      color: #64748b;
    }

    .salary-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1rem;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .salary-summary {
      display: flex;
      gap: 2rem;
    }

    .salary-amount {
      display: flex;
      flex-direction: column;
    }

    .salary-amount .label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .salary-amount .value {
      font-size: 1.5rem;
      font-weight: 600;
      font-family: monospace;
    }

    .salary-amount.gross .value { color: #16a34a; }
    .salary-amount.net .value { color: #1e293b; }
    .salary-amount.ctc .value { color: #4f46e5; }

    .salary-actions {
      display: flex;
      gap: 0.5rem;
    }

    .salary-meta {
      display: flex;
      gap: 1.5rem;
      margin-bottom: 1.5rem;
      font-size: 0.875rem;
      color: #64748b;
    }

    .salary-meta i {
      margin-right: 0.25rem;
    }

    .components-section h4,
    .history-section h4 {
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .components-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .component-card {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      padding: 1rem;
    }

    .component-card h5 {
      margin: 0 0 0.75rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .component-card.earnings h5 { color: #16a34a; }
    .component-card.deductions h5 { color: #dc2626; }

    .component-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .component-item {
      display: flex;
      justify-content: space-between;
      font-size: 0.875rem;
    }

    .comp-name {
      color: #374151;
    }

    .comp-amount {
      font-family: monospace;
      color: #1e293b;
    }

    .overridden {
      color: #f59e0b;
      font-size: 0.625rem;
      margin-left: 0.25rem;
    }

    .component-total {
      display: flex;
      justify-content: space-between;
      margin-top: 0.75rem;
      padding-top: 0.75rem;
      border-top: 1px solid #e2e8f0;
      font-weight: 600;
      font-size: 0.875rem;
    }

    .no-components {
      color: #64748b;
      font-size: 0.875rem;
      font-style: italic;
    }

    .history-section {
      margin-top: 1.5rem;
      padding-top: 1.5rem;
      border-top: 1px solid #e2e8f0;
    }

    .no-history {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      color: #64748b;
      font-size: 0.875rem;
    }

    .history-table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.875rem;
    }

    .history-table th {
      text-align: left;
      padding: 0.625rem 0.75rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .history-table td {
      padding: 0.75rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .history-table tr.current {
      background: #f0fdf4;
    }

    .current-badge {
      display: inline-flex;
      padding: 0.125rem 0.5rem;
      background: #dcfce7;
      color: #166534;
      border-radius: 9999px;
      font-size: 0.625rem;
      font-weight: 500;
      margin-left: 0.5rem;
    }

    /* Form Styles */
    .salary-form {
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
    }

    .input:focus,
    .select:focus,
    .textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .error {
      font-size: 0.75rem;
      color: #dc2626;
    }

    .components-form-section {
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

    .section-header h4 {
      margin: 0;
    }

    .no-components-form {
      text-align: center;
      padding: 2rem;
      color: #64748b;
    }

    .no-components-form i {
      font-size: 2rem;
      color: #cbd5e1;
      margin-bottom: 0.5rem;
    }

    .no-components-form p {
      margin: 0;
      font-weight: 500;
    }

    .no-components-form span {
      font-size: 0.75rem;
    }

    .components-form-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .component-form-row {
      display: grid;
      grid-template-columns: 1fr 120px auto 40px;
      gap: 0.5rem;
      align-items: center;
      padding: 0.5rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
    }

    .comp-select .select,
    .comp-amount .input {
      width: 100%;
      padding: 0.5rem;
    }

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      font-size: 0.75rem;
      color: #64748b;
      white-space: nowrap;
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
    }

    .remove-btn:hover {
      background: #fef2f2;
    }

    .form-summary {
      display: flex;
      justify-content: flex-end;
      gap: 2rem;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
    }

    .summary-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      color: #64748b;
    }

    .summary-value {
      font-family: monospace;
      font-weight: 600;
      color: #1e293b;
    }

    .summary-item.net .summary-value {
      font-size: 1rem;
      color: #16a34a;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
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

    @media (max-width: 768px) {
      .salary-summary {
        flex-direction: column;
        gap: 1rem;
      }

      .components-grid {
        grid-template-columns: 1fr;
      }

      .form-row {
        grid-template-columns: 1fr;
      }

      .component-form-row {
        grid-template-columns: 1fr;
        gap: 0.25rem;
      }

      .form-summary {
        flex-direction: column;
        gap: 0.5rem;
      }
    }
  `],
})
export class StaffSalaryComponent implements OnInit {
  @Input() staffId!: string;

  private salaryService = inject(SalaryService);
  private toastService = inject(ToastService);
  private fb = inject(FormBuilder);

  // State
  loading = signal(true);
  saving = signal(false);
  error = signal<string | null>(null);
  currentSalary = signal<StaffSalary | null>(null);
  salaryHistory = signal<StaffSalary[]>([]);
  showHistory = signal(false);
  showModal = signal(false);

  // Dropdowns
  structures = signal<StructureDropdownItem[]>([]);
  availableComponents = signal<ComponentDropdownItem[]>([]);

  // Form
  form: FormGroup = this.fb.group({
    structureId: [''],
    effectiveFrom: ['', Validators.required],
    revisionReason: [''],
    components: this.fb.array([]),
  });

  get componentsArray(): FormArray {
    return this.form.get('components') as FormArray;
  }

  // Computed values
  earningComponents = computed(() => {
    const salary = this.currentSalary();
    if (!salary?.components) return [];
    return salary.components.filter(c => c.componentType === 'earning');
  });

  deductionComponents = computed(() => {
    const salary = this.currentSalary();
    if (!salary?.components) return [];
    return salary.components.filter(c => c.componentType === 'deduction');
  });

  totalEarnings = computed(() => {
    return this.earningComponents().reduce(
      (sum, c) => sum + parseFloat(c.amount || '0'),
      0
    ).toString();
  });

  totalDeductions = computed(() => {
    return this.deductionComponents().reduce(
      (sum, c) => sum + parseFloat(c.amount || '0'),
      0
    ).toString();
  });

  // Calculated values for form
  calculatedGross = computed(() => {
    const components = this.form.get('components')?.value || [];
    const allComps = this.availableComponents();
    return components.reduce((sum: number, c: { componentId: string; amount: string }) => {
      const comp = allComps.find(ac => ac.id === c.componentId);
      if (comp?.componentType === 'earning') {
        return sum + parseFloat(c.amount || '0');
      }
      return sum;
    }, 0).toString();
  });

  calculatedDeductions = computed(() => {
    const components = this.form.get('components')?.value || [];
    const allComps = this.availableComponents();
    return components.reduce((sum: number, c: { componentId: string; amount: string }) => {
      const comp = allComps.find(ac => ac.id === c.componentId);
      if (comp?.componentType === 'deduction') {
        return sum + parseFloat(c.amount || '0');
      }
      return sum;
    }, 0).toString();
  });

  calculatedNet = computed(() => {
    const gross = parseFloat(this.calculatedGross());
    const deductions = parseFloat(this.calculatedDeductions());
    return (gross - deductions).toString();
  });

  ngOnInit(): void {
    this.loadSalary();
    this.loadDropdowns();
  }

  loadSalary(): void {
    this.loading.set(true);
    this.error.set(null);

    this.salaryService.getStaffSalary(this.staffId).subscribe({
      next: salary => {
        this.currentSalary.set(salary);
        this.loading.set(false);
        this.loadHistory();
      },
      error: err => {
        // No salary assigned is not an error
        if (err.status === 404) {
          this.currentSalary.set(null);
          this.loading.set(false);
        } else {
          this.error.set('Failed to load salary information');
          this.loading.set(false);
        }
      },
    });
  }

  loadHistory(): void {
    this.salaryService.getStaffSalaryHistory(this.staffId).subscribe({
      next: history => this.salaryHistory.set(history),
      error: () => this.salaryHistory.set([]),
    });
  }

  loadDropdowns(): void {
    this.salaryService.getStructuresDropdown().subscribe({
      next: structures => this.structures.set(structures),
      error: () => this.structures.set([]),
    });

    this.salaryService.getComponentsDropdown().subscribe({
      next: components => this.availableComponents.set(components),
      error: () => this.availableComponents.set([]),
    });
  }

  toggleHistory(): void {
    this.showHistory.set(!this.showHistory());
  }

  openAssignModal(): void {
    this.resetForm();
    this.showModal.set(true);
  }

  openReviseModal(): void {
    this.resetForm();
    const current = this.currentSalary();
    if (current) {
      this.form.patchValue({
        structureId: current.structureId || '',
        effectiveFrom: new Date().toISOString().split('T')[0],
      });

      // Populate components
      current.components?.forEach(comp => {
        this.componentsArray.push(
          this.fb.group({
            componentId: [comp.componentId, Validators.required],
            amount: [comp.amount, [Validators.required, Validators.min(0)]],
            isOverridden: [comp.isOverridden],
          })
        );
      });
    }
    this.showModal.set(true);
  }

  closeModal(): void {
    this.showModal.set(false);
    this.resetForm();
  }

  resetForm(): void {
    this.form.reset({
      structureId: '',
      effectiveFrom: new Date().toISOString().split('T')[0],
      revisionReason: '',
    });
    this.componentsArray.clear();
  }

  onStructureChange(): void {
    const structureId = this.form.get('structureId')?.value;
    if (!structureId) return;

    // Load structure and populate components
    this.salaryService.getStructure(structureId).subscribe({
      next: structure => {
        this.componentsArray.clear();
        structure.components?.forEach(comp => {
          this.componentsArray.push(
            this.fb.group({
              componentId: [comp.componentId, Validators.required],
              amount: [comp.amount || '', [Validators.required, Validators.min(0)]],
              isOverridden: [false],
            })
          );
        });
      },
    });
  }

  addComponent(): void {
    this.componentsArray.push(
      this.fb.group({
        componentId: ['', Validators.required],
        amount: ['', [Validators.required, Validators.min(0)]],
        isOverridden: [false],
      })
    );
  }

  removeComponent(index: number): void {
    this.componentsArray.removeAt(index);
  }

  onSubmit(): void {
    if (this.form.invalid || this.componentsArray.length === 0) return;

    this.saving.set(true);
    const formValue = this.form.value;

    const components: StaffSalaryComponentInput[] = formValue.components.map(
      (c: { componentId: string; amount: string; isOverridden: boolean }) => ({
        componentId: c.componentId,
        amount: c.amount,
        isOverridden: c.isOverridden || false,
      })
    );

    const request: AssignSalaryRequest = {
      staffId: this.staffId,
      structureId: formValue.structureId || undefined,
      effectiveFrom: formValue.effectiveFrom,
      components,
      revisionReason: formValue.revisionReason?.trim() || undefined,
    };

    this.salaryService.assignSalary(request).subscribe({
      next: () => {
        this.toastService.success(
          this.currentSalary() ? 'Salary updated successfully' : 'Salary assigned successfully'
        );
        this.closeModal();
        this.loadSalary();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error('Failed to save salary');
        this.saving.set(false);
      },
    });
  }

  formatCurrency(amount: string): string {
    const num = parseFloat(amount || '0');
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(num);
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }
}
